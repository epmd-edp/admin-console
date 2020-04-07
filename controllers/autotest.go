package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/models/command"
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
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
		log.Printf("Failed to validate autotests request data: %s", errMsg.Message)
		flash := beego.NewFlash()
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("%s/admin/edp/autotest/create", context.BasePath), 302)
		return
	}
	logAutotestsRequestData(codebase)

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Autotests %s is already exists.", codebase.Name)
			flash.Store(&c.Controller)
			c.Redirect(fmt.Sprintf("%s/admin/edp/autotest/create", context.BasePath), 302)
			return
		}
		c.Abort("500")
		return
	}

	log.Printf("Autotests object is saved into k8s: %s", createdObject)
	flash.Success("Autotests object is created.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("%s/admin/edp/autotest/overview?%s=%s#codebaseSuccessModal", context.BasePath, paramWaitingForCodebase, codebase.Name), 302)
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

	log.Println(result.String())
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

func validateAutotestsRequestData(autotests command.CreateCodebase) *ErrMsg {
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
		return &ErrMsg{"An internal error has occurred on server while validating autotests's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
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
		log.Println("Import strategy is used.")

		gitServers, err := c.GitServerService.GetServers(query.GitServerCriteria{Available: true})
		if err != nil {
			c.Abort("500")
			return
		}
		log.Printf("Fetched Git Servers: %v", gitServers)

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

	log.Println("Create strategy is removed from list due to Autotest")

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = isAdmin(c.GetSession("realm_roles").([]string))
	c.Data["IsVcsEnabled"] = isVcsEnabled
	c.Data["Type"] = query.Autotests
	c.Data["IntegrationStrategies"] = c.IntegrationStrategies
	c.Data["CodeBaseIntegrationStrategy"] = true
	c.Data["JenkinsSlaves"] = s
	c.Data["BuildTools"] = c.BuildTools
	c.Data["JobProvisioners"] = p
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["BasePath"] = context.BasePath
	c.TplName = "create_autotest.html"
}

func (c *AutotestsController) GetAutotestsOverviewPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["success"] != "" {
		c.Data["Success"] = true
	}

	codebases, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Type: query.Autotests,
	})
	codebases = addCodebaseInProgressIfAny(codebases, c.GetString(paramWaitingForCodebase))
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["Codebases"] = codebases
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = isAdmin(c.GetSession("realm_roles").([]string))
	c.Data["Type"] = query.Autotests
	c.Data["BasePath"] = context.BasePath
	c.TplName = "codebase.html"
}
