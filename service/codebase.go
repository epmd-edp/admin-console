/*
 * Copyright 2019 EPAM Systems.
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
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"fmt"
	edpv1alpha1 "github.com/epmd-edp/codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strings"
	"time"
)

var clog = logf.Log.WithName("codebase-service")

type CodebaseService struct {
	Clients               k8s.ClientSet
	ICodebaseRepository   repository.ICodebaseRepository
	ICDPipelineRepository repository.ICDPipelineRepository
	BranchService         CodebaseBranchService
}

func (s CodebaseService) CreateCodebase(codebase command.CreateCodebase) (*edpv1alpha1.Codebase, error) {
	clog.Info("start creating Codebase resource", "name", codebase.Name)

	codebaseCr, err := util.GetCodebaseCR(s.Clients.EDPRestClient, codebase.Name)
	if err != nil {
		clog.Info("an error has occurred while fetching Codebase CR from cluster", "name", codebase.Name)
		return nil, err
	}

	if codebaseCr != nil {
		clog.Info("codebase already exists in cluster", "name", codebaseCr.Name)
		return nil, errors.New("CODEBASE_ALREADY_EXISTS")
	}

	codebaseDb, err := s.GetCodebaseByName(codebase.Name)
	if err != nil {
		clog.Info("an error has occurred while fetching Codebase entity from DB: %s", codebase.Name)
		return nil, err
	}

	if codebaseDb != nil {
		clog.Info("Codebase is already exists in DB", "name", codebaseDb.Name)
		return nil, errors.New("CODEBASE_ALREADY_EXISTS")
	}

	edpClient := s.Clients.EDPRestClient
	coreClient := s.Clients.CoreClient

	crd := &edpv1alpha1.Codebase{
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
	clog.Info("CR was generated. Waiting to save ...", "cr", crd)

	err = createTempSecrets(context.Namespace, codebase, coreClient)

	if err != nil {
		return nil, err
	}

	result := &edpv1alpha1.Codebase{}
	err = edpClient.Post().Namespace(context.Namespace).Resource(consts.CodebasePlural).Body(crd).Do().Into(result)

	if err != nil {
		clog.Error(err, "an error has occurred while creating codebase resource in cluster")
		return &edpv1alpha1.Codebase{}, err
	}

	_, err = s.BranchService.CreateCodebaseBranch(command.CreateCodebaseBranch{
		Name:     "master",
		Username: codebase.Username,
	}, codebase.Name)
	if err != nil {
		clog.Error(err, "an error has been occurred during the master branch creation")
		return &edpv1alpha1.Codebase{}, err
	}
	return result, nil
}

func (s *CodebaseService) GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error) {
	codebases, err := s.ICodebaseRepository.GetCodebasesByCriteria(criteria)
	if err != nil {
		clog.Error(err, "an error has occurred while getting codebase objects")
		return nil, err
	}
	clog.Info("fetched codebases", "count", len(codebases), "codebase", codebases)

	return codebases, nil
}

func (s CodebaseService) GetCodebaseByName(name string) (*query.Codebase, error) {
	codebase, err := s.ICodebaseRepository.GetCodebaseByName(name)
	if err != nil {
		clog.Error(err, "an error has occurred while getting codebase object", "name", name)
		return nil, err
	}
	clog.Info("fetched codebase info", "codebase", codebase)

	return codebase, nil
}

func (s CodebaseService) ExistCodebaseAndBranch(cbName, brName string) bool {
	return s.ICodebaseRepository.ExistCodebaseAndBranch(cbName, brName)
}

func createSecret(namespace string, secret *v1.Secret, coreClient *coreV1Client.CoreV1Client) (*v1.Secret, error) {
	createdSecret, err := coreClient.Secrets(namespace).Create(secret)
	if err != nil {
		clog.Error(err, "an error has occurred while saving secret")
		return &v1.Secret{}, err
	}
	return createdSecret, nil
}

func createTempSecrets(namespace string, codebase command.CreateCodebase, coreClient *coreV1Client.CoreV1Client) error {
	if codebase.Repository != nil && (codebase.Repository.Login != "" && codebase.Repository.Password != "") {
		repoSecretName := fmt.Sprintf("repository-codebase-%s-temp", codebase.Name)
		tempRepoSecret := getSecret(repoSecretName, codebase.Repository.Login, codebase.Repository.Password)

		if _, err := createSecret(namespace, tempRepoSecret, coreClient); err != nil {
			clog.Error(err, "an error has occurred while creating repository secret")
			return err
		}
		clog.Info("repository secret for codebase was created", "codebase", codebase.Name)
	}

	if codebase.Vcs != nil {
		vcsSecretName := fmt.Sprintf("vcs-autouser-codebase-%s-temp", codebase.Name)
		tempVcsSecret := getSecret(vcsSecretName, codebase.Vcs.Login, codebase.Vcs.Password)

		if _, err := createSecret(namespace, tempVcsSecret, coreClient); err != nil {
			clog.Error(err, "an error has occurred while creating vcs secret")
			return err
		}
		clog.Info("VCS secret for codebase was created", "codebase", codebase.Name)
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
	s := edpv1alpha1.CodebaseSpec{
		Lang:             codebase.Lang,
		Framework:        codebase.Framework,
		BuildTool:        codebase.BuildTool,
		Strategy:         edpv1alpha1.Strategy(codebase.Strategy),
		Type:             codebase.Type,
		GitServer:        codebase.GitServer,
		JenkinsSlave:     codebase.JenkinsSlave,
		JobProvisioning:  codebase.JobProvisioning,
		DeploymentScript: codebase.DeploymentScript,
	}

	if s.Strategy == "import" {
		s.GitUrlPath = codebase.GitUrlPath
	}

	if codebase.Framework != nil {
		s.Framework = codebase.Framework
	}

	if codebase.Repository != nil {
		s.Repository = &edpv1alpha1.Repository{
			Url: codebase.Repository.Url,
		}
	}

	if codebase.Route != nil {
		s.Route = &edpv1alpha1.Route{
			Site: codebase.Route.Site,
		}
		if len(codebase.Route.Path) > 0 {
			s.Route.Path = codebase.Route.Path
		}
	}

	if codebase.Database != nil {
		s.Database = &edpv1alpha1.Database{
			Kind:     codebase.Database.Kind,
			Version:  codebase.Database.Version,
			Capacity: codebase.Database.Capacity,
			Storage:  codebase.Database.Storage,
		}
	}

	if codebase.TestReportFramework != nil {
		s.TestReportFramework = codebase.TestReportFramework
	}

	if codebase.Description != nil {
		s.Description = codebase.Description
	}

	s.Versioning.Type = edpv1alpha1.VersioningType(codebase.Versioning.Type)
	s.Versioning.StartFrom = codebase.Versioning.StartFrom

	return s
}

func (s CodebaseService) CheckBranch(apps []models.CDPipelineApplicationCommand) (bool, error) {
	for _, app := range apps {
		exist, err := s.ICodebaseRepository.ExistActiveBranch(app.InputDockerStream)
		if err != nil {
			clog.Error(err, "an error has occurred while checking status of branch")
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

	clog.Info("fetched Application to promote", "applications", result)

	return result, nil
}

func (s CodebaseService) Delete(name, codebaseType string) error {
	clog.Info("start executing service delete method", "codebase", name)
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
	clog.Info("end executing service codebase delete method", "codebase", name)
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
	clog.Info("start executing codebase delete request", "codebase", name)
	r := &edpv1alpha1.Codebase{}
	err := s.Clients.EDPRestClient.Delete().
		Namespace(context.Namespace).
		Resource(consts.CodebasePlural).
		Name(name).
		Do().Into(r)
	if err != nil {
		return errors.Wrapf(err, "couldn't delete codebase %v from cluster", name)
	}
	clog.Info("end executing codebase delete request", "codebase", name)
	return nil
}
