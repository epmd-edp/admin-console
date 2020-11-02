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

package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/controllers/validation"
	"edp-admin-console/models/command"
	edperror "edp-admin-console/models/error"
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	cbs "edp-admin-console/service/codebasebranch"
	jiraservice "edp-admin-console/service/jira-server"
	"edp-admin-console/service/logger"
	"edp-admin-console/service/platform"
	"edp-admin-console/util"
	"edp-admin-console/util/auth"
	"fmt"
	"html/template"
	"strings"

	"github.com/astaxie/beego"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

type ApplicationController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	EDPTenantService service.EDPTenantService
	BranchService    cbs.CodebaseBranchService
	GitServerService service.GitServerService
	SlaveService     service.SlaveService
	JobProvisioning  service.JobProvisioning
	JiraServer       jiraservice.JiraServer

	IntegrationStrategies []string
	BuildTools            []string
	VersioningTypes       []string
	DeploymentScript      []string
	CiTools               []string
}

const (
	paramWaitingForCodebase = "waitingforcodebase"
	scope                   = "ci"
)

func (c *ApplicationController) GetApplicationsOverviewPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	applications, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Type: query.App,
	})
	applications = addCodebaseInProgressIfAny(applications, c.GetString(paramWaitingForCodebase))
	if err != nil {
		c.Abort("500")
		return
	}

	if flash.Data["success"] != "" {
		c.Data["Success"] = true
	}
	if flash.Data["error"] != "" {
		c.Data["DeletionError"] = flash.Data["error"]
	}
	contextRoles := c.GetSession("realm_roles").([]string)
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = auth.IsAdmin(contextRoles)
	c.Data["Codebases"] = applications
	c.Data["Type"] = query.App
	c.Data["VersioningTypes"] = c.VersioningTypes
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["BasePath"] = context.BasePath
	c.Data["DiagramPageEnabled"] = context.DiagramPageEnabled
	c.TplName = "codebase.html"
}

func addCodebaseInProgressIfAny(codebases []*query.Codebase, codebaseInProgress string) []*query.Codebase {
	if codebaseInProgress != "" {
		for _, codebase := range codebases {
			if codebase.Name == codebaseInProgress {
				return codebases
			}
		}

		log.Debug("adding codebase which is going to be created to the list",
			zap.String("name", codebaseInProgress))
		app := query.Codebase{
			Name:   codebaseInProgress,
			Status: query.Inactive,
		}
		codebases = append(codebases, &app)
	}
	return codebases
}

func (c *ApplicationController) GetCreateApplicationPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	isVcsEnabled, err := c.EDPTenantService.GetVcsIntegrationValue()
	if err != nil {
		c.Abort("500")
		return
	}
	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}

	if doesIntegrationStrategiesContainImportStrategy(c.IntegrationStrategies) {
		log.Info("Import strategy is used.")
		gitServers, err := c.GitServerService.GetServers(query.GitServerCriteria{Available: true})
		if err != nil {
			c.Abort("500")
			return
		}
		log.Debug("fetched Git Servers", zap.Any("git servers", gitServers))
		c.Data["GitServers"] = gitServers
	}

	s, err := c.SlaveService.GetAllSlaves()
	if err != nil {
		c.Abort("500")
		return
	}

	p, err := c.JobProvisioning.GetAllJobProvisioners(query.JobProvisioningCriteria{Scope: util.GetStringP(scope)})
	if err != nil {
		c.Abort("500")
		return
	}

	servers, err := c.JiraServer.GetJiraServers()
	if err != nil {
		log.Error(err.Error())
		c.Abort("500")
		return
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["IsVcsEnabled"] = isVcsEnabled
	c.Data["Type"] = query.App
	c.Data["CodeBaseIntegrationStrategy"] = true
	c.Data["IntegrationStrategies"] = c.IntegrationStrategies
	c.Data["JenkinsSlaves"] = s
	c.Data["BuildTools"] = c.BuildTools
	c.Data["JobProvisioners"] = p
	c.Data["VersioningTypes"] = c.VersioningTypes
	c.Data["DeploymentScripts"] = c.DeploymentScript
	c.Data["IsOpenshift"] = platform.IsOpenshift()
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["BasePath"] = context.BasePath
	c.Data["JiraServer"] = servers
	c.Data["DiagramPageEnabled"] = context.DiagramPageEnabled
	c.Data["CiTools"] = c.CiTools
	c.TplName = "create_application.html"
}

func doesIntegrationStrategiesContainImportStrategy(integrationStrategies []string) bool {
	return contains(integrationStrategies, "import")
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == strings.ToLower(n) {
			return true
		}
	}
	return false
}

func (c *ApplicationController) CreateApplication() {
	flash := beego.NewFlash()
	codebase := c.extractApplicationRequestData()
	errMsg := validation.ValidCodebaseRequestData(codebase)
	if errMsg != nil {
		log.Error("failed to validate request data", zap.String("message", errMsg.Message))
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("%s/admin/edp/application/create", context.BasePath), 302)
		return
	}
	ld := validation.CreateCodebaseLogRequestData(codebase)
	log.Info(ld.String())

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)
	if err != nil {
		c.checkError(err, flash, codebase.Name, codebase.GitUrlPath)
		return
	}

	log.Info("application object is saved into cluster", zap.String("name", createdObject.Name))
	flash.Success("Application object is created.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("%s/admin/edp/application/overview?%s=%s#codebaseSuccessModal", context.BasePath, paramWaitingForCodebase, codebase.Name), 302)
}

func (c *ApplicationController) checkError(err error, flash *beego.FlashData, name string, url *string) {
	switch err.(type) {
	case *edperror.CodebaseAlreadyExistsError:
		flash.Error("Application %v already exists.", name)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("%s/admin/edp/application/create", context.BasePath), 302)
	case *edperror.CodebaseWithGitUrlPathAlreadyExistsError:
		flash.Error("Application %v with %v project path already exists.", name, *url)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("%s/admin/edp/application/create", context.BasePath), 302)
	default:
		log.Error("couldn't create codebase", zap.Error(err))
		c.Abort("500")
	}
}

func (c *ApplicationController) extractApplicationRequestData() command.CreateCodebase {
	codebase := command.CreateCodebase{
		Lang:             c.GetString("appLang"),
		BuildTool:        c.GetString("buildTool"),
		Strategy:         strings.ToLower(c.GetString("strategy")),
		Type:             "application",
		DeploymentScript: c.GetString("deploymentScript"),
		CiTool:           c.GetString("ciTool"),
		Name:             c.GetString("appName"),
		DefaultBranch:    c.GetString("defaultBranchName"),
	}

	if js := c.GetString("jenkinsSlave"); len(js) > 0 {
		codebase.JenkinsSlave = &js
	}

	if jp := c.GetString("jobProvisioning"); len(jp) > 0 {
		codebase.JobProvisioning = &jp
	}

	if s := c.GetString("jiraServer"); len(s) > 0 {
		codebase.JiraServer = &s
	}

	if v := c.GetString("commitMessagePattern"); len(v) > 0 {
		codebase.CommitMessageRegex = &v
	}

	if v := c.GetString("ticketNamePattern"); len(v) > 0 {
		codebase.TicketNameRegex = &v
	}

	codebase.Versioning.Type = c.GetString("versioningType")
	startVersioningFrom := c.GetString("startVersioningFrom")
	sp := c.GetString("snapshotStaticField")

	codebase.Versioning.StartFrom = util.GetVersionOrNil(startVersioningFrom, sp)

	if codebase.Strategy == "import" {
		codebase.GitServer = c.GetString("gitServer")
		gitRepoPath := c.GetString("gitRelativePath")
		codebase.GitUrlPath = &gitRepoPath
	} else {
		codebase.GitServer = "gerrit"
	}

	framework := c.GetString("framework")
	codebase.Framework = &framework

	isMultiModule, _ := c.GetBool("isMultiModule", false)
	codebase.MultiModule = isMultiModule

	if isMultiModule {
		multimoduleApp := fmt.Sprintf("%s-multimodule", *codebase.Framework)
		codebase.Framework = &multimoduleApp
	}

	repoUrl := c.GetString("gitRepoUrl")
	if repoUrl != "" {
		codebase.Repository = &command.Repository{
			Url: repoUrl,
		}

		isRepoPrivate, _ := c.GetBool("isRepoPrivate", false)
		if isRepoPrivate {
			codebase.Repository.Login = c.GetString("repoLogin")
			codebase.Repository.Password = c.GetString("repoPassword")
		}
	}

	vcsLogin := c.GetString("vcsLogin")
	vcsPassword := c.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		codebase.Vcs = &command.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}

	needRoute, _ := c.GetBool("needRoute", false)
	if needRoute {
		codebase.Route = &command.Route{
			Site: c.GetString("routeSite"),
		}
		if len(c.GetString("routePath")) > 0 {
			codebase.Route.Path = c.GetString("routePath")
		}
	}

	needDb, _ := c.GetBool("needDb", false)
	if needDb {
		codebase.Database = &command.Database{
			Kind:     c.GetString("database"),
			Version:  c.GetString("dbVersion"),
			Capacity: c.GetString("dbCapacity") + c.GetString("capacityExt"),
			Storage:  c.GetString("dbPersistentStorage"),
		}
	}
	codebase.Username = c.Ctx.Input.Session("username").(string)
	return codebase
}
