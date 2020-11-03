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

package pipeline

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/mock"
	"edp-admin-console/service"
	"edp-admin-console/service/cd_pipeline"
	"edp-admin-console/util"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEditCDPipelinePageMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	cdpm := new(mock.MockCdPipeline)
	mCodebase := new(mock.MockCodebase)
	mjp := new(mock.MockJobProvision)

	c := CDPipelineController{
		PipelineService: cd_pipeline.CDPipelineService{
			ICDPipelineRepository: cdpm,
		},
		CodebaseService: service.CodebaseService{
			ICodebaseRepository: mCodebase,
		},
		JobProvisioning: service.JobProvisioning{
			IJobProvisioningRepository: mjp,
		},
	}

	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)
	input.SetParam(":name", "stub-name")
	input.Context.Request.Header = http.Header{
		"Cookie": []string{"BEEGO_FLASH=test"},
	}
	c.Ctx = input.Context
	c.Ctx.Input.Context.Input = input
	c.Ctx.Input.CruSession = &session.CookieSessionStore{}
	c.Ctx.Input.CruSession.Flush()
	c.Ctx.Input.CruSession.Set("realm_roles", []string{"stub-value"})
	c.Data = map[interface{}]interface{}{}

	cdpm.On("GetCDPipelineByName", "stub-name").Return(nil, nil)

	mCodebase.On("GetCodebasesByCriteria", query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.App,
	}).Return([]*query.Codebase{
		{
			Id:   1,
			Name: "fake-name",
		},
	}, nil)

	mCodebase.On("GetCodebasesByCriteria", query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.Library,
		Language:     "groovy-pipeline",
	}).Return([]*query.Codebase{
		{
			Id:   1,
			Name: "fake-name",
		},
	}, nil)

	mCodebase.On("GetCodebasesByCriteria", query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.Autotests,
	}).Return([]*query.Codebase{
		{
			Id:   1,
			Name: "fake-name",
		},
	}, nil)

	mjp.On("GetAllJobProvisioners", query.JobProvisioningCriteria{Scope: util.GetStringP("cd")}).Return(
		[]*query.JobProvisioning{
			{
				Id:   1,
				Name: "fake-job-provison",
			},
		}, nil)

	c.GetEditCDPipelinePage()
}
