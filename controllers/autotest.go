package controllers

import (
	"edp-admin-console/context"
	validation2 "edp-admin-console/controllers/validation"
	"edp-admin-console/models/command"
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"edp-admin-console/util/auth"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"html/template"
	"net/http"
	"path"
	"regexp"
	"strings"
)

type AutotestsController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	EDPTenantService service.EDPTenantService
	BranchService    service.CodebaseBranchService
	GitServerService service.GitServerService
	SlaveService     service.SlaveService
	JobProvisioning  service.JobProvisioning

	IntegrationStrategies []string
	BuildTools            []string
	VersioningTypes       []string
	DeploymentScript      []string
}

const (
	ImportStrategy = "Import"
)

func (c *AutotestsController) CreateAutotests() {
	flash := beego.NewFlash()
	codebase := c.extractAutotestsRequestData()
	errMsg := validateAutotestsRequestData(codebase)
	if errMsg != nil {
		log.Info("Failed to validate autotests request data", "err", errMsg.Message)
		flash := beego.NewFlash()
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect("/admin/edp/autotest/create", 302)
		return
	}
	logAutotestsRequestData(codebase)

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Autotests %s is already exists.", codebase.Name)
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/autotest/create", 302)
			return
		}
		c.Abort("500")
		return
	}

	log.Info("Autotests object is saved into cluster", "name", createdObject.Name)
	flash.Success("Autotests object is created.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/admin/edp/autotest/overview?%s=%s#codebaseSuccessModal", paramWaitingForCodebase, codebase.Name), 302)
}

func logAutotestsRequestData(autotests command.CreateCodebase) {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Request data to create codebase CR is valid. name=%s, strategy=%s, lang=%s, buildTool=%s, testReportFramework=%s",
		autotests.Name, autotests.Strategy, autotests.Lang, autotests.BuildTool, *autotests.TestReportFramework))

	if autotests.Repository != nil {
		result.WriteString(fmt.Sprintf(", repositoryUrl=%s, repositoryLogin=%s", autotests.Repository.Url, autotests.Repository.Login))
	}

	if autotests.Vcs != nil {
		result.WriteString(fmt.Sprintf(", vcsLogin=%s", autotests.Vcs.Login))
	}

	log.Info(result.String())
}

func (c *AutotestsController) extractAutotestsRequestData() command.CreateCodebase {
	codebase := command.CreateCodebase{
		Lang:             c.GetString("appLang"),
		BuildTool:        c.GetString("buildTool"),
		Strategy:         strings.ToLower(c.GetString("strategy")),
		Type:             "autotests",
		JenkinsSlave:     c.GetString("jenkinsSlave"),
		JobProvisioning:  c.GetString("jobProvisioning"),
		DeploymentScript: c.GetString("deploymentScript"),
	}

	codebase.Versioning.Type = c.GetString("versioningType")
	startVersioningFrom := c.GetString("startVersioningFrom")
	codebase.Versioning.StartFrom = util.GetStringOrNil(startVersioningFrom)

	if o := OtherLanguage; codebase.Lang == OtherLanguage {
		codebase.Framework = &o
	}

	if codebase.Strategy == strings.ToLower(ImportStrategy) {
		codebase.GitServer = c.GetString("gitServer")
		gitRepoPath := c.GetString("gitRelativePath")
		codebase.GitUrlPath = &gitRepoPath
		codebase.Name = path.Base(*codebase.GitUrlPath)
	} else {
		codebase.Name = c.GetString("nameOfApp")
		codebase.GitServer = "gerrit"
	}

	testReportFramework := c.GetString("testReportFramework")
	if testReportFramework != "" {
		codebase.TestReportFramework = &testReportFramework
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

	description := c.GetString("description")
	if description != "" {
		codebase.Description = &description
	}

	vcsLogin := c.GetString("vcsLogin")
	vcsPassword := c.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		codebase.Vcs = &command.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}
	codebase.Username = c.Ctx.Input.Session("username").(string)
	return codebase
}

func validateAutotestsRequestData(autotests command.CreateCodebase) *validation2.ErrMsg {
	valid := validation.Validation{}

	_, err := valid.Valid(autotests)

	if autotests.Strategy == strings.ToLower(ImportStrategy) {
		valid.Match(autotests.GitUrlPath, regexp.MustCompile("^\\/.*$"), "Spec.GitUrlPath")
	}

	if autotests.Repository != nil {
		_, err = valid.Valid(autotests.Repository)

		isAvailable := util.IsGitRepoAvailable(autotests.Repository.Url, autotests.Repository.Login, autotests.Repository.Password)

		if !isAvailable {
			err := &validation.Error{Key: "repository", Message: "Repository doesn't exist or invalid login and password."}
			valid.Errors = append(valid.Errors, err)
		}
	}

	if autotests.Description == nil {
		err := &validation.Error{Key: "description", Message: "Description field can't be empty."}
		valid.Errors = append(valid.Errors, err)
	}

	if autotests.Vcs != nil {
		_, err = valid.Valid(autotests.Vcs)
	}

	if err != nil {
		return &validation2.ErrMsg{"An internal error has occurred on server while validating autotests's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &validation2.ErrMsg{string(validation2.CreateErrorResponseBody(valid)), http.StatusBadRequest}
}

func (c *AutotestsController) GetCreateAutotestsPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}

	isVcsEnabled, err := c.EDPTenantService.GetVcsIntegrationValue()
	if err != nil {
		c.Abort("500")
		return
	}

	contains := doesIntegrationStrategiesContainImportStrategy(c.IntegrationStrategies)
	if contains {
		log.Info("Import strategy is used.")

		gitServers, err := c.GitServerService.GetServers(query.GitServerCriteria{Available: true})
		if err != nil {
			c.Abort("500")
			return
		}
		log.Info("Fetched Git Servers", "git servers", gitServers)

		c.Data["GitServers"] = gitServers
	}

	s, err := c.SlaveService.GetAllSlaves()
	if err != nil {
		c.Abort("500")
		return
	}

	p, err := c.JobProvisioning.GetAllJobProvisioners()
	if err != nil {
		c.Abort("500")
		return
	}

	log.Info("Create strategy is removed from list due to Autotest")

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = auth.IsAdmin(c.GetSession("realm_roles").([]string))
	c.Data["IsVcsEnabled"] = isVcsEnabled
	c.Data["Type"] = query.Autotests
	c.Data["IntegrationStrategies"] = c.IntegrationStrategies
	c.Data["CodeBaseIntegrationStrategy"] = true
	c.Data["JenkinsSlaves"] = s
	c.Data["BuildTools"] = c.BuildTools
	c.Data["JobProvisioners"] = p
	c.Data["VersioningTypes"] = c.VersioningTypes
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "create_autotest.html"
}

func (c *AutotestsController) GetAutotestsOverviewPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	codebases, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Type: query.Autotests,
	})
	codebases = addCodebaseInProgressIfAny(codebases, c.GetString(paramWaitingForCodebase))
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
	c.Data["Codebases"] = codebases
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = auth.IsAdmin(c.GetSession("realm_roles").([]string))
	c.Data["Type"] = query.Autotests
	c.Data["VersioningTypes"] = c.VersioningTypes
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "codebase.html"
}
