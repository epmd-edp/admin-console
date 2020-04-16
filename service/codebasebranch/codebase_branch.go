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

package codebasebranch

import (
	"edp-admin-console/context"
	"edp-admin-console/k8s"
	"edp-admin-console/models/command"
	"edp-admin-console/models/query"
	"edp-admin-console/repository"
	"edp-admin-console/service/logger"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"fmt"
	edpv1alpha1 "github.com/epmd-edp/codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)

var log = logger.GetLogger()

type CodebaseBranchService struct {
	Clients                  k8s.ClientSet
	IReleaseBranchRepository repository.ICodebaseBranchRepository
	ICDPipelineRepository    repository.ICDPipelineRepository
	ICodebaseRepository      repository.ICodebaseRepository
	CodebaseBranchValidation map[string]func(string, string) ([]string, error)
}

func (s *CodebaseBranchService) CreateCodebaseBranch(branchInfo command.CreateCodebaseBranch, appName string) (*edpv1alpha1.CodebaseBranch, error) {
	log.Debug("start creating CodebaseBranch CR",
		zap.String("codebase", appName), zap.String("branch", branchInfo.Name))
	edpRestClient := s.Clients.EDPRestClient

	cb := util.ProcessNameToKubernetesConvention(branchInfo.Name)

	releaseBranchCR, err := getReleaseBranchCR(edpRestClient, cb, appName, context.Namespace)
	if err != nil {
		return nil, errors.Wrapf(err, "an error has occurred while getting %v CodebaseBranch CR from cluster", cb)
	}

	if releaseBranchCR != nil {
		return nil, fmt.Errorf("CodebaseBranch %v already exists", cb)
	}

	c, err := util.GetCodebaseCR(s.Clients.EDPRestClient, appName)
	if err != nil {
		return nil, err
	}

	branch := &edpv1alpha1.CodebaseBranch{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v2.edp.epam.com/v1alpha1",
			Kind:       "CodebaseBranch",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", appName, cb),
			Namespace: context.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "v2.edp.epam.com/v1alpha1",
					Kind:               consts.CodebaseKind,
					Name:               c.Name,
					UID:                c.UID,
					BlockOwnerDeletion: newTrue(),
				},
			},
		},
		Spec: edpv1alpha1.CodebaseBranchSpec{
			BranchName:   branchInfo.Name,
			FromCommit:   branchInfo.Commit,
			Version:      branchInfo.Version,
			Release:      branchInfo.Release,
			CodebaseName: appName,
		},
		Status: edpv1alpha1.CodebaseBranchStatus{
			Build:               branchInfo.Build,
			Status:              "initialized",
			LastTimeUpdated:     time.Now(),
			LastSuccessfulBuild: nil,
			Username:            branchInfo.Username,
			Action:              "codebase_branch_registration",
			Result:              "success",
			Value:               "inactive",
		},
	}

	result := &edpv1alpha1.CodebaseBranch{}
	err = edpRestClient.Post().Namespace(context.Namespace).Resource("codebasebranches").Body(branch).Do().Into(result)
	if err != nil {
		return &edpv1alpha1.CodebaseBranch{}, errors.Wrap(err, "an error has occurred while creating CodebaseBranch CR in cluster")
	}
	return result, nil
}

func newTrue() *bool {
	b := true
	return &b
}

func (s *CodebaseBranchService) UpdateCodebaseBranch(appName, branchName string, version *string) error {
	log.Debug("start updating CodebaseBranch CR",
		zap.String("version", *version),
		zap.String("branch", branchName))
	edpRestClient := s.Clients.EDPRestClient
	br, err := getReleaseBranchCR(edpRestClient, branchName, appName, context.Namespace)
	if err != nil {
		return err
	}

	br.Spec.Version = version
	bytes, err := util.EncodeStructToBytes(br)
	if err != nil {
		return err
	}

	err = edpRestClient.Patch(types.MergePatchType).
		Namespace(context.Namespace).
		Resource(consts.CodebaseBranchPlural).
		Name(fmt.Sprintf("%v-%v", appName, branchName)).
		Body(bytes).
		Do().Error()
	if err != nil {
		return errors.Wrapf(err, "couldn't update codebase branch %v from cluster", branchName)
	}
	log.Info("codebase branch has been updated",
		zap.String("name", branchName),
		zap.String("version", *version),
		zap.String("appName", appName))
	return nil
}

func (s *CodebaseBranchService) GetCodebaseBranchesByCriteria(criteria query.CodebaseBranchCriteria) ([]query.CodebaseBranch, error) {
	codebaseBranches, err := s.IReleaseBranchRepository.GetCodebaseBranchesByCriteria(criteria)
	if err != nil {
		return nil, errors.Wrap(err, "an error has occurred while getting branch entities")
	}
	return codebaseBranches, nil
}

func getReleaseBranchCR(edpRestClient *rest.RESTClient, branchName string, appName string, namespace string) (*edpv1alpha1.CodebaseBranch, error) {
	result := &edpv1alpha1.CodebaseBranch{}
	err := edpRestClient.Get().Namespace(namespace).Resource("codebasebranches").Name(fmt.Sprintf("%s-%s", appName, branchName)).Do().Into(result)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Debug("CodebaseBranch doesn't exist in cluster", zap.String("branch", branchName))
			return nil, nil
		}

		return nil, errors.Wrapf(err, "an error has occurred while getting CodebaseBranch CR from cluster")
	}

	return result, nil
}

func (s *CodebaseBranchService) Delete(codebase, branch string) error {
	log.Debug("start executing service codebase branch delete method",
		zap.String("name", codebase),
		zap.String("branch", branch))
	if err := s.canCodebaseBranchBeDeleted(codebase, branch); err != nil {
		return err
	}

	crbn := fmt.Sprintf("%v-%v", codebase, branch)
	if err := s.deleteCodebaseBranch(crbn); err != nil {
		return err
	}
	log.Info("codebase branch has been marked for deletion",
		zap.String("name", codebase),
		zap.String("branch", branch))
	return nil
}

func (s *CodebaseBranchService) canCodebaseBranchBeDeleted(codebase, branch string) error {
	c, err := s.ICodebaseRepository.GetCodebaseByName(codebase)
	if err != nil {
		return err
	}
	p, err := s.CodebaseBranchValidation[string(c.Type)](codebase, branch)
	if err != nil {
		return err
	}
	if p != nil {
		return dberror.RemoveCodebaseBranchRestriction{
			Status:  dberror.StatusReasonCodebaseBranchIsUsedByCDPipeline,
			Message: fmt.Sprintf("%v CodebaseBranch is used in %v CD Pipeline(s)", branch, strings.Join(p, ",")),
		}
	}
	return nil
}

func (s *CodebaseBranchService) deleteCodebaseBranch(name string) error {
	cb := &edpv1alpha1.CodebaseBranch{}
	err := s.Clients.EDPRestClient.Delete().
		Namespace(context.Namespace).
		Resource(consts.CodebaseBranchPlural).
		Name(name).
		Do().Into(cb)
	if err != nil {
		return errors.Wrapf(err, "couldn't delete codebase branch %v from cluster", name)
	}
	return nil
}
