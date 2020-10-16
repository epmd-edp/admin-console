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

func TestGetCreatePageMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	sm := new(mock.MockSlave)
	mjp := new(mock.MockJobProvision)
	mjs := new(mock.MockJiraServer)
	mpb := new(mock.MockPerfBoard)

	lc := LibraryController{
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

	lc.Ctx = input.Context
	lc.Ctx.Input.Context.Input = input
	lc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	lc.Ctx.Input.CruSession.Flush()
	lc.Ctx.Input.CruSession.Set("realm_roles", []string{"stub-value"})
	lc.Data = map[interface{}]interface{}{}

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

	lc.GetCreatePage()
}

func TestGetCreatePageMethod_ShouldReturnError(t *testing.T) {
	sm := new(mock.MockSlave)
	mjp := new(mock.MockJobProvision)
	mjs := new(mock.MockJiraServer)
	mpb := new(mock.MockPerfBoard)

	lc := LibraryController{
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

	lc.Ctx = input.Context
	lc.Ctx.Input.Context.Input = input
	lc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	lc.Ctx.Input.CruSession.Flush()
	lc.Ctx.Input.CruSession.Set("realm_roles", []string{"stub-value"})
	lc.Data = map[interface{}]interface{}{}

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

	assert.Panics(t, lc.GetCreatePage)
}

func TestExtractLibraryRequestData_ReturnsDtoWithPerfField(t *testing.T) {
	lc := LibraryController{}
	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)
	input.SetParam("perfServer", "stub-name")
	lc.Ctx = input.Context
	lc.Ctx.Input.Context.Input = input
	lc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	lc.Ctx.Input.CruSession.Flush()
	lc.Ctx.Input.CruSession.Set("username", "stub-value")

	c := lc.extractLibraryRequestData()
	assert.Equal(t, "stub-name", c.Perf.Name)
}

func TestExtractLibraryRequestData_ReturnsDtoWithoutPerfField(t *testing.T) {
	lc := LibraryController{}
	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)
	lc.Ctx = input.Context
	lc.Ctx.Input.Context.Input = input
	lc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	lc.Ctx.Input.CruSession.Flush()
	lc.Ctx.Input.CruSession.Set("username", "stub-value")

	c := lc.extractLibraryRequestData()
	assert.Nil(t, c.Perf)
}
