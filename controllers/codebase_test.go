package controllers

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/mock"
	"edp-admin-console/service"
	edpComponent "edp-admin-console/service/edp-component"
	"edp-admin-console/util/consts"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCodebaseOverviewPage_ShouldBeExecutedSuccessfully(t *testing.T) {
	mCodebase := new(mock.MockCodebase)
	mGitServer := new(mock.MockGitServer)
	mEdpComponent := new(mock.MockEDPComponent)
	cc := CodebaseController{
		CodebaseService: service.CodebaseService{
			ICodebaseRepository: mCodebase,
		},
		GitServerService: service.GitServerService{
			IGitServerRepository: mGitServer,
		},
		EDPComponent: edpComponent.EDPComponentService{
			IEDPComponent: mEdpComponent,
		},
	}

	input := context.NewInput()
	input.Context = context.NewContext()
	r, _ := http.NewRequest("GET", "stub-request", nil)
	input.Context.Reset(httptest.NewRecorder(), r)
	input.SetParam(":codebaseName", "stub-name")
	input.SetParam("waitingforbranch", "stub-name")
	input.Context.Request.Header = http.Header{
		"Cookie": []string{"BEEGO_FLASH=test"},
	}
	cc.Ctx = input.Context
	cc.Ctx.Input.Context.Input = input
	cc.Ctx.Input.CruSession = &session.CookieSessionStore{}
	cc.Ctx.Input.CruSession.Flush()
	cc.Ctx.Input.CruSession.Set("username", "stub-value")
	cc.Ctx.Input.CruSession.Set("realm_roles", []string{"stub-value"})
	cc.Data = map[interface{}]interface{}{}

	gs := "stub-name"
	pp := "/stub/project/path"
	mCodebase.On("GetCodebaseByName", "stub-name").Return(
		query.Codebase{
			Name:           "test",
			Type:           "application",
			Strategy:       consts.ImportStrategy,
			GitServer:      &gs,
			GitProjectPath: &pp,
			CodebaseBranch: []*query.CodebaseBranch{
				{
					Name: "stub-name",
				},
			},
		}, nil)
	mGitServer.On("GetGitServerByName", "stub-name").Return(
		query.GitServer{
			Id:        0,
			Name:      "stub-name",
			Hostname:  "stub-host",
			Available: true,
		}, nil)
	mEdpComponent.On("GetEDPComponent", consts.Jenkins).Return(
		query.EDPComponent{
			Id:      0,
			Type:    consts.Jenkins,
			Url:     "stub-url",
			Icon:    "stub-icon",
			Visible: false,
		}, nil)

	cc.GetCodebaseOverviewPage()
}

func TestGetCiLinkMethod_ShouldReturnJenkinsCiLink(t *testing.T) {
	c := query.Codebase{
		Name:   "stub-name",
		CiTool: consts.JenkinsCITool,
	}
	l := getCiLink(c, "jenkins-stub-host", "stub-name", "git-stub-host")
	assert.Equal(t, "jenkins-stub-host/job/stub-name/view/STUB-NAME", l)
}

func TestGetCiLinkMethod_ShouldReturnGitlabCiLink(t *testing.T) {
	url := "/stub"
	c := query.Codebase{
		GitProjectPath: &url,
		CiTool:         "GitlabCI",
	}
	l := getCiLink(c, "jenkins-stub-host", "stub-name", "git-stub-host")
	assert.Equal(t, "https://git-stub-host/stub/pipelines?scope=branches&page=1", l)
}
