package controllers

import (
	"edp-admin-console/models/query"
	"edp-admin-console/util/consts"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
