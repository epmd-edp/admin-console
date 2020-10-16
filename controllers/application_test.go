package controllers

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/mock"
	"edp-admin-console/service"
	jiraservice "edp-admin-console/service/jira-server"
	"edp-admin-console/service/perfboard"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
	"github.com/epmd-edp/codebase-operator/v2/pkg/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCreateApplicationPageMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	sm := new(mock.MockSlave)
	mjp := new(mock.MockJobProvision)
	mjs := new(mock.MockJiraServer)
	mpb := new(mock.MockPerfBoard)

	apc := ApplicationController{
		EDPTenantService: service.EDPTenantService{},
		SlaveService:     service.SlaveService{ISlaveRepository: sm},
		JobProvisioning:  service.JobProvisioning{IJobProvisioningRepository: mjp},
		JiraServer:       jiraservice.JiraServer{IJiraServer: mjs},
		PerfService:      perfboard.PerfBoard{PerfRepo: mpb},
	}

	beego.AppConfig.Set("vcsIntegrationEnabled", "false")
	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)

	apc.Ctx = input.Context
	apc.Ctx.Input.Context.Input = input
	apc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	apc.Ctx.Input.CruSession.Flush()
	apc.Data = map[interface{}]interface{}{}

	sm.On("GetAllSlaves").Return(
		[]*query.JenkinsSlave{
			{
				Id:   1,
				Name: "fake-slave",
			},
		}, nil)

	mjp.On("GetAllJobProvisioners", query.JobProvisioningCriteria{Scope: util.GetStringP("ci")}).Return(
		[]*query.JobProvisioning{
			{
				Id:   1,
				Name: "fake-job-provison",
			},
		}, nil)

	mjs.On("GetJiraServers").Return(
		[]*query.JiraServer{
			{
				Id:   1,
				Name: "fake-jira-server",
			},
		}, nil)

	mpb.On("GetPerfServers").Return(
		[]*query.PerfServer{
			{
				Id:   1,
				Name: "fake-perf-server",
			},
		}, nil)

	apc.GetCreateApplicationPage()
}

func TestGetCreateApplicationPageMethod_GetPerfServersShouldReturnError(t *testing.T) {
	sm := new(mock.MockSlave)
	mjp := new(mock.MockJobProvision)
	mjs := new(mock.MockJiraServer)
	mpb := new(mock.MockPerfBoard)

	apc := ApplicationController{
		EDPTenantService: service.EDPTenantService{},
		SlaveService:     service.SlaveService{ISlaveRepository: sm},
		JobProvisioning:  service.JobProvisioning{IJobProvisioningRepository: mjp},
		JiraServer:       jiraservice.JiraServer{IJiraServer: mjs},
		PerfService:      perfboard.PerfBoard{PerfRepo: mpb},
	}

	beego.AppConfig.Set("vcsIntegrationEnabled", "false")
	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)

	apc.Ctx = input.Context
	apc.Ctx.Input.Context.Input = input
	apc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	apc.Ctx.Input.CruSession.Flush()
	apc.Data = map[interface{}]interface{}{}

	sm.On("GetAllSlaves").Return(
		[]*query.JenkinsSlave{
			{
				Id:   1,
				Name: "fake-slave",
			},
		}, nil)

	mjp.On("GetAllJobProvisioners", query.JobProvisioningCriteria{Scope: util.GetStringP("ci")}).Return(
		[]*query.JobProvisioning{
			{
				Id:   1,
				Name: "fake-job-provison",
			},
		}, nil)

	mjs.On("GetJiraServers").Return(
		[]*query.JiraServer{
			{
				Id:   1,
				Name: "fake-jira-server",
			},
		}, nil)

	mpb.On("GetPerfServers").Return(nil, errors.New("failed"))

	assert.Panics(t, apc.GetCreateApplicationPage)
}

func TestExtractApplicationRequestDataMethod_ReturnsDtoWithPerfField(t *testing.T) {
	apc := ApplicationController{}
	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)
	input.SetParam("perfServer", "stub-name")
	apc.Ctx = input.Context
	apc.Ctx.Input.Context.Input = input
	apc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	apc.Ctx.Input.CruSession.Flush()
	apc.Ctx.Input.CruSession.Set("username", "stub-value")

	c := apc.extractApplicationRequestData()
	assert.Equal(t, "stub-name", c.Perf.Name)
}

func TestExtractApplicationRequestDataMethod_ReturnsDtoWithoutPerfField(t *testing.T) {
	apc := ApplicationController{}
	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)
	apc.Ctx = input.Context
	apc.Ctx.Input.Context.Input = input
	apc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	apc.Ctx.Input.CruSession.Flush()
	apc.Ctx.Input.CruSession.Set("username", "stub-value")

	c := apc.extractApplicationRequestData()
	assert.Nil(t, c.Perf)
}
