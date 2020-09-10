package util

import (
	"fmt"
	"strings"
)

func CreateNativeProjectLink(domain, project string) string {
	return fmt.Sprintf("%v/console/project/%v/overview", domain, project)
}

func CreateNonNativeProjectLink(domain, namespace string) string {
	return fmt.Sprintf("%v/#/overview?namespace=%v", domain, namespace)
}

func CreateNativeDockerStreamLink(domain, namespace, stream string) string {
	return fmt.Sprintf("%v/console/project/%v/browse/images/%v", domain, namespace, stream)
}

func CreateNonNativeDockerStreamLink(domain, stream string) string {
	return fmt.Sprintf("%v/%v/", domain, stream)
}

func CreateCICDApplicationLink(domain, codebase, branch string) string {
	return fmt.Sprintf("%v/job/%s/view/%s", domain, codebase, strings.ToUpper(branch))
}

func CreateCICDPipelineLink(domain, pipelineName string) string {
	return fmt.Sprintf("%v/job/%v-%v", domain, pipelineName, "cd-pipeline")
}

func CreateGerritLink(domain, codebaseName, branchName string) string {
	return fmt.Sprintf("%v/gitweb?p=%s.git;a=shortlog;h=refs/heads/%s", domain, codebaseName, branchName)
}

func CreateGitLink(hostname, path, branch string) string {
	return fmt.Sprintf("https://%s%s/commits/%s", hostname, path, branch)
}

func CreateGitlabCILink(domain, relativePath string) string {
	return fmt.Sprintf("https://%v%v/pipelines?scope=branches&page=1", domain, relativePath)
}
