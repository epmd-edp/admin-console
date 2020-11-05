package util

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
)

func CreateNativeProjectLink(domain, project string) string {
	replacer := strings.NewReplacer("{namespace}", project)
	return fmt.Sprintf("%v%v", domain, replacer.Replace(beego.AppConfig.String("projectMaskUrl")))
}

func CreateNativeDockerStreamLink(domain, namespace, stream string) string {
	replacer := strings.NewReplacer("{namespace}", namespace, "{stream}", stream)
	return fmt.Sprintf("%v%v", domain, replacer.Replace(beego.AppConfig.String("imageStreamMaskUrl")))
}

func CreateNonNativeDockerStreamLink(domain, stream string) string {
	replacer := strings.NewReplacer("{stream}", stream)
	return fmt.Sprintf("%v%v", domain, replacer.Replace(beego.AppConfig.String("imageStreamMaskUrl")))
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
