/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"edp-admin-console/context"
	"edp-admin-console/k8s"
	"edp-admin-console/models"
	"edp-admin-console/models/command"
	edperror "edp-admin-console/models/error"
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	cbs "edp-admin-console/service/codebasebranch"
	"edp-admin-console/service/logger"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"fmt"
	edpv1alpha1 "github.com/epmd-edp/codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	"strings"
	"time"
)

var clog = logger.GetLogger()

type CodebaseService struct {
	Clients               k8s.ClientSet
	ICodebaseRepository   repository.ICodebaseRepository
	ICDPipelineRepository repository.ICDPipelineRepository
	BranchService         cbs.CodebaseBranchService
}

func (s CodebaseService) CreateCodebase(codebase command.CreateCodebase) (*edpv1alpha1.Codebase, error) {
	clog.Info("start creating Codebase resource", zap.String("name", codebase.Name))

	codebaseCr, err := util.GetCodebaseCR(s.Clients.EDPRestClient, codebase.Name)
	if err != nil {
		clog.Info("an error has occurred while fetching Codebase CR from cluster",
			zap.String("name", codebase.Name))
		return nil, err
	}
	if codebaseCr != nil {
		clog.Info("codebase already exists in cluster", zap.String("name", codebaseCr.Name))
		return nil, edperror.NewCodebaseAlreadyExistsError()
	}

	if s.findCodebaseByName(codebase.Name) {
		clog.Info("Codebase already exists in DB", zap.String("name", codebase.Name))
		return nil, edperror.NewCodebaseAlreadyExistsError()
	}

	if s.findCodebaseByProjectPath(codebase.GitUrlPath) {
		clog.Info("Codebase with the same gitUrlPath already exists in DB",
			zap.String("gitUrlPath", *codebase.GitUrlPath))
		return nil, edperror.NewCodebaseWithGitUrlPathAlreadyExistsError()
	}

	edpClient := s.Clients.EDPRestClient
	coreClient := s.Clients.CoreClient

	c := &edpv1alpha1.Codebase{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v2.edp.epam.com/v1alpha1",
			Kind:       consts.CodebaseKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       codebase.Name,
			Namespace:  context.Namespace,
			Finalizers: []string{"foregroundDeletion"},
		},
		Spec: convertData(codebase),
		Status: edpv1alpha1.CodebaseStatus{
			Available:       false,
			LastTimeUpdated: time.Now(),
			Status:          "initialized",
			Username:        codebase.Username,
			Action:          "codebase_registration",
			Result:          "success",
			Value:           "inactive",
		},
	}
	clog.Debug("CR was generated. Waiting to save ...", zap.String("name", c.Name))

	if err := createTempSecrets(context.Namespace, codebase, coreClient); err != nil {
		return nil, err
	}

	result := &edpv1alpha1.Codebase{}
	err = edpClient.Post().Namespace(context.Namespace).Resource(consts.CodebasePlural).Body(c).Do().Into(result)
	if err != nil {
		clog.Error("an error has occurred while creating codebase resource in cluster", zap.Error(err))
		return &edpv1alpha1.Codebase{}, err
	}

	p := setCodebaseBranchCr(codebase.Versioning.Type, codebase.Username, codebase.Versioning.StartFrom, codebase.DefaultBranch)

	if _, err = s.BranchService.CreateCodebaseBranch(p, codebase.Name); err != nil {
		clog.Error("an error has been occurred during the master branch creation", zap.Error(err))
		return &edpv1alpha1.Codebase{}, err
	}
	return result, nil
}

func (s *CodebaseService) GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error) {
	codebases, err := s.ICodebaseRepository.GetCodebasesByCriteria(criteria)
	if err != nil {
		clog.Error("an error has occurred while getting codebase objects", zap.Error(err))
		return nil, err
	}
	clog.Debug("fetched codebases", zap.Int("count", len(codebases)))

	return codebases, nil
}

func (s CodebaseService) GetCodebaseByName(name string) (*query.Codebase, error) {
	c, err := s.ICodebaseRepository.GetCodebaseByName(name)
	if err != nil {
		return nil, errors.Wrapf(err, "an error has occurred while getting %v codebase from db", name)
	}
	clog.Info("codebase has been fetched from db", zap.String("name", c.Name))
	return c, nil
}

func (s *CodebaseService) findCodebaseByName(name string) bool {
	exist := s.ICodebaseRepository.FindCodebaseByName(name)
	clog.Debug("codebase exists", zap.Bool("exists", exist), zap.String("name", name))
	return exist
}

func (s *CodebaseService) findCodebaseByProjectPath(gitProjectPath *string) bool {
	if gitProjectPath == nil {
		return false
	}
	exist := s.ICodebaseRepository.FindCodebaseByProjectPath(gitProjectPath)
	clog.Debug("codebase exists", zap.Bool("exists", exist), zap.String("url", *gitProjectPath))
	return exist
}

func (s CodebaseService) ExistCodebaseAndBranch(cbName, brName string) bool {
	return s.ICodebaseRepository.ExistCodebaseAndBranch(cbName, brName)
}

func createSecret(namespace string, secret *v1.Secret, coreClient *coreV1Client.CoreV1Client) (*v1.Secret, error) {
	createdSecret, err := coreClient.Secrets(namespace).Create(secret)
	if err != nil {
		clog.Error("an error has occurred while saving secret", zap.Error(err))
		return &v1.Secret{}, err
	}
	return createdSecret, nil
}

func createTempSecrets(namespace string, codebase command.CreateCodebase, coreClient *coreV1Client.CoreV1Client) error {
	if codebase.Repository != nil && (codebase.Repository.Login != "" && codebase.Repository.Password != "") {
		repoSecretName := fmt.Sprintf("repository-codebase-%s-temp", codebase.Name)
		tempRepoSecret := getSecret(repoSecretName, codebase.Repository.Login, codebase.Repository.Password)

		if _, err := createSecret(namespace, tempRepoSecret, coreClient); err != nil {
			clog.Error("an error has occurred while creating repository secret", zap.Error(err))
			return err
		}
		clog.Info("repository secret for codebase was created", zap.String("codebase", codebase.Name))
	}

	if codebase.Vcs != nil {
		vcsSecretName := fmt.Sprintf("vcs-autouser-codebase-%s-temp", codebase.Name)
		tempVcsSecret := getSecret(vcsSecretName, codebase.Vcs.Login, codebase.Vcs.Password)

		if _, err := createSecret(namespace, tempVcsSecret, coreClient); err != nil {
			clog.Error("an error has occurred while creating vcs secret", zap.Error(err))
			return err
		}
		clog.Info("VCS secret for codebase was created", zap.String("codebase", codebase.Name))
	}

	return nil
}

func getSecret(name string, username string, password string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: map[string]string{
			"username": username,
			"password": password,
		},
	}
}

func convertData(codebase command.CreateCodebase) edpv1alpha1.CodebaseSpec {
	cs := edpv1alpha1.CodebaseSpec{
		Lang:                 codebase.Lang,
		Framework:            codebase.Framework,
		BuildTool:            codebase.BuildTool,
		Strategy:             edpv1alpha1.Strategy(codebase.Strategy),
		Type:                 codebase.Type,
		GitServer:            codebase.GitServer,
		JenkinsSlave:         codebase.JenkinsSlave,
		JobProvisioning:      codebase.JobProvisioning,
		DeploymentScript:     codebase.DeploymentScript,
		JiraServer:           codebase.JiraServer,
		CommitMessagePattern: codebase.CommitMessageRegex,
		TicketNamePattern:    codebase.TicketNameRegex,
		CiTool:               codebase.CiTool,
	}
	if cs.Strategy == "import" {
		cs.GitUrlPath = codebase.GitUrlPath
	}
	if codebase.Framework != nil {
		cs.Framework = codebase.Framework
	}
	if codebase.Repository != nil {
		cs.Repository = &edpv1alpha1.Repository{
			Url: codebase.Repository.Url,
		}
	}
	if codebase.Route != nil {
		cs.Route = &edpv1alpha1.Route{
			Site: codebase.Route.Site,
		}
		if len(codebase.Route.Path) > 0 {
			cs.Route.Path = codebase.Route.Path
		}
	}
	if codebase.Database != nil {
		cs.Database = &edpv1alpha1.Database{
			Kind:     codebase.Database.Kind,
			Version:  codebase.Database.Version,
			Capacity: codebase.Database.Capacity,
			Storage:  codebase.Database.Storage,
		}
	}
	if codebase.TestReportFramework != nil {
		cs.TestReportFramework = codebase.TestReportFramework
	}
	if codebase.Description != nil {
		cs.Description = codebase.Description
	}
	cs.Versioning.Type = edpv1alpha1.VersioningType(codebase.Versioning.Type)
	cs.Versioning.StartFrom = codebase.Versioning.StartFrom
	return cs
}

func (s CodebaseService) CheckBranch(apps []models.CDPipelineApplicationCommand) (bool, error) {
	for _, app := range apps {
		exist, err := s.ICodebaseRepository.ExistActiveBranch(app.InputDockerStream)
		if err != nil {
			clog.Error("an error has occurred while checking status of branch", zap.Error(err))
			return false, err
		}

		if !exist {
			return false, nil
		}
	}
	return true, nil
}

func (s CodebaseService) GetApplicationsToPromote(cdPipelineId int) ([]string, error) {
	appsToPromote, err := s.ICodebaseRepository.SelectApplicationToPromote(cdPipelineId)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred while fetching Ids of applications which shoould be promoted: %v", err)
	}
	return s.selectApplicationNames(appsToPromote)
}

func (s CodebaseService) selectApplicationNames(applicationsToPromote []*query.ApplicationsToPromote) ([]string, error) {
	var result []string
	for _, app := range applicationsToPromote {
		codebase, err := s.ICodebaseRepository.GetCodebaseById(app.CodebaseId)
		if err != nil {
			return nil, fmt.Errorf("an error has occurred while fetching Codebase by Id %v: %v", app.CodebaseId, err)
		}
		result = append(result, codebase.Name)
	}
	clog.Debug("Applications to promote have been fetched", zap.Any("applications", result))
	return result, nil
}

func (s CodebaseService) Delete(name, codebaseType string) error {
	clog.Debug("start executing service delete method", zap.String("codebase", name))
	cdp, err := s.getCdPipelinesUsingCodebase(name, codebaseType)
	if err != nil {
		return err
	}
	if cdp != nil {
		p := strings.Join(cdp[:], ",")
		return dberror.CodebaseIsUsedByCDPipeline{
			Status:   dberror.StatusReasonCodebaseIsUsedByCDPipeline,
			Message:  fmt.Sprintf("%v %v is used by %v CD Pipeline(s). couldn't delete.", codebaseType, name, p),
			Codebase: name,
			Pipeline: p,
		}
	}

	if err := s.deleteCodebase(name); err != nil {
		return err
	}
	clog.Info("end executing service codebase delete method", zap.String("codebase", name))
	return nil
}

func (s CodebaseService) getCdPipelinesUsingCodebase(name, codebaseType string) ([]string, error) {
	if consts.Application == codebaseType {
		cdp, err := s.ICDPipelineRepository.GetCDPipelinesUsingApplication(name)
		if err != nil {
			return nil, err
		}
		return cdp, nil
	} else if consts.Autotest == codebaseType {
		cdp, err := s.ICDPipelineRepository.GetCDPipelinesUsingAutotest(name)
		if err != nil {
			return nil, err
		}
		return cdp, nil
	} else {
		cdp, err := s.ICDPipelineRepository.GetCDPipelinesUsingLibrary(name)
		if err != nil {
			return nil, err
		}
		return cdp, nil
	}
}

func (s CodebaseService) deleteCodebase(name string) error {
	clog.Debug("start executing codebase delete request", zap.String("codebase", name))
	r := &edpv1alpha1.Codebase{}
	err := s.Clients.EDPRestClient.Delete().
		Namespace(context.Namespace).
		Resource(consts.CodebasePlural).
		Name(name).
		Do().Into(r)
	if err != nil {
		return errors.Wrapf(err, "couldn't delete codebase %v from cluster", name)
	}
	clog.Debug("end executing codebase delete request", zap.String("codebase", name))
	return nil
}

func setCodebaseBranchCr(vt string, username string, version *string, defaultBranch string) command.CreateCodebaseBranch {
	if vt == consts.DefaultVersioningType {
		return command.CreateCodebaseBranch{
			Name:     defaultBranch,
			Username: username,
			Build:    &consts.DefaultBuildNumber,
		}
	}

	return command.CreateCodebaseBranch{
		Name:     defaultBranch,
		Username: username,
		Version:  version,
		Build:    &consts.DefaultBuildNumber,
	}
}

func (s *CodebaseService) Update(command command.UpdateCodebaseCommand) (*edpv1alpha1.Codebase, error) {
	log.Debug("start executing Update method fort codebase", zap.String("name", command.Name))
	c, err := util.GetCodebaseCR(s.Clients.EDPRestClient, command.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get codebase from cluster %v", command.Name)
	}

	c.Spec.CommitMessagePattern = &command.CommitMessageRegex
	c.Spec.TicketNamePattern = &command.TicketNameRegex
	log.Debug("new values",
		zap.String("commitMessagePattern", *c.Spec.CommitMessagePattern),
		zap.String("ticketNamePattern", *c.Spec.TicketNamePattern))

	if err := s.executeUpdateRequest(c); err != nil {
		return nil, err
	}
	log.Info("codebase has been updated", zap.String("name", c.Name))
	return c, nil
}

func (s *CodebaseService) executeUpdateRequest(c *edpv1alpha1.Codebase) error {
	err := s.Clients.EDPRestClient.Put().
		Namespace(context.Namespace).
		Resource("codebases").
		Name(c.Name).
		Body(c).
		Do().
		Into(c)
	if err != nil {
		return errors.Wrap(err, "an error has occurred while updating Codebase CR in cluster")
	}
	return nil
}
