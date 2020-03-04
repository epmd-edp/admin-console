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
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"fmt"
	edpv1alpha1 "github.com/epmd-edp/codebase-operator/v2/pkg/apis/edp/v1alpha1"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strings"
	"time"
)

var log = logf.Log.WithName("codebase-branch-service")

type CodebaseBranchService struct {
	Clients                  k8s.ClientSet
	IReleaseBranchRepository repository.ICodebaseBranchRepository
	ICDPipelineRepository    repository.ICDPipelineRepository
	ICodebaseRepository      repository.ICodebaseRepository
	CodebaseBranchValidation map[string]func(string, string) ([]string, error)
}

func (s *CodebaseBranchService) CreateCodebaseBranch(branchInfo command.CreateCodebaseBranch, appName string) (*edpv1alpha1.CodebaseBranch, error) {
	log.V(2).Info("start creating CodebaseBranch CR", "codebase", appName, "branch", branchInfo.Name)
	edpRestClient := s.Clients.EDPRestClient

	releaseBranchCR, err := getReleaseBranchCR(edpRestClient, branchInfo.Name, appName, context.Namespace)
	if err != nil {
		return nil, errors.Wrapf(err, "an error has occurred while getting %v CodebaseBranch CR from cluster", branchInfo.Name)
	}

	if releaseBranchCR != nil {
		return nil, fmt.Errorf("CodebaseBranch %v already exists", branchInfo.Name)
	}

	c, err := util.GetCodebaseCR(s.Clients.EDPRestClient, appName)
	if err != nil {
		return nil, err
	}

	spec := convertBranchInfoData(branchInfo, appName)
	branch := &edpv1alpha1.CodebaseBranch{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v2.edp.epam.com/v1alpha1",
			Kind:       "CodebaseBranch",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", appName, branchInfo.Name),
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
		Spec: spec,
		Status: edpv1alpha1.CodebaseBranchStatus{
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
	log.V(2).Info("start updating CodebaseBranch CR", "version", version, "branch",
		branchName, "version", version)
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

	log.Info("codebase branch has been updated", "name", branchName, "version", version, "appName", appName)
	return nil
}

func (s *CodebaseBranchService) GetCodebaseBranchesByCriteria(criteria query.CodebaseBranchCriteria) ([]query.CodebaseBranch, error) {
	codebaseBranches, err := s.IReleaseBranchRepository.GetCodebaseBranchesByCriteria(criteria)
	if err != nil {
		return nil, errors.Wrap(err, "an error has occurred while getting branch entities")
	}
	return codebaseBranches, nil
}

func convertBranchInfoData(branchInfo command.CreateCodebaseBranch, appName string) edpv1alpha1.CodebaseBranchSpec {
	return edpv1alpha1.CodebaseBranchSpec{
		BranchName:   branchInfo.Name,
		FromCommit:   branchInfo.Commit,
		Version:      branchInfo.Version,
		Build:        branchInfo.Build,
		Release:      branchInfo.Release,
		CodebaseName: appName,
	}
}

func getReleaseBranchCR(edpRestClient *rest.RESTClient, branchName string, appName string, namespace string) (*edpv1alpha1.CodebaseBranch, error) {
	result := &edpv1alpha1.CodebaseBranch{}
	err := edpRestClient.Get().Namespace(namespace).Resource("codebasebranches").Name(fmt.Sprintf("%s-%s", appName, branchName)).Do().Into(result)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.V(2).Info("CodebaseBranch doesn't exist in cluster", "branch", branchName)
			return nil, nil
		}

		return nil, errors.Wrapf(err, "an error has occurred while getting CodebaseBranch CR from cluster")
	}

	return result, nil
}

func (s *CodebaseBranchService) Delete(codebase, branch string) error {
	log.V(2).Info("start executing service codebase branch delete method",
		"name", codebase, "branch", branch)
	if err := s.canCodebaseBranchBeDeleted(codebase, branch); err != nil {
		return err
	}

	crbn := fmt.Sprintf("%v-%v", codebase, branch)
	if err := s.deleteCodebaseBranch(crbn); err != nil {
		return err
	}
	log.Info("codebase branch has been marked for deletion",
		"name", codebase, "branch", branch)
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
			Status: dberror.StatusReasonCodebaseBranchIsUsedByCDPipeline,
			Message: fmt.Sprintf("%v %v CodebaseBranch is used in %v CD Pipeline(s)",
				c.Type, branch, strings.Join(p, ",")),
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
