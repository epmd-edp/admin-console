package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGitlabCILinkMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	link := CreateGitlabCILink("stub-domain", "stub-relative-path")
	assert.Equal(t, "https://stub-domainstub-relative-path/pipelines?scope=branches&page=1", link)
}
