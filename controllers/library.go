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

type LibraryController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
	CodebaseService  service.CodebaseService
	BranchService    service.CodebaseBranchService
	GitServerService service.GitServerService
	SlaveService     service.SlaveService
	JobProvisioning  service.JobProvisioning

	IntegrationStrategies []string
	BuildTools            []string
	VersioningTypes       []string
	DeploymentScript      []string
}

func (c *LibraryController) GetLibraryListPage() {
	flash := beego.ReadFromRequest(&c.Controller)
	codebases, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Type: query.Library,
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
	c.Data["Type"] = query.Library
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "codebase.html"
}

func (c *LibraryController) GetCreatePage() {
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

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = auth.IsAdmin(c.GetSession("realm_roles").([]string))
	c.Data["IsVcsEnabled"] = isVcsEnabled
	c.Data["Type"] = query.Library
	c.Data["CodeBaseIntegrationStrategy"] = true
	c.Data["IntegrationStrategies"] = c.IntegrationStrategies
	c.Data["JenkinsSlaves"] = s
	c.Data["BuildTools"] = c.BuildTools
	c.Data["JobProvisioners"] = p
	c.Data["VersioningTypes"] = c.VersioningTypes
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "create_library.html"
}

func (c *LibraryController) Create() {
	flash := beego.NewFlash()
	codebase := c.extractLibraryRequestData()
	errMsg := validateLibraryRequestData(codebase)
	if errMsg != nil {
		log.Info("Failed to validate library request data", "err", errMsg.Message)
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect("/admin/edp/library/create", 302)
		return
	}
	logLibraryRequestData(codebase)

	createdObject, err := c.CodebaseService.CreateCodebase(codebase)

	if err != nil {
		if err.Error() == "CODEBASE_ALREADY_EXISTS" {
			flash.Error("Library %s is already exists.", codebase.Name)
			flash.Store(&c.Controller)
			c.Redirect("/admin/edp/library/create", 302)
			return
		}
		c.Abort("500")
		return
	}

	log.Info("Library object is saved into cluster", "library", createdObject.Name)
	flash.Success("Library object is created.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/admin/edp/library/overview?%s=%s#codebaseSuccessModal", paramWaitingForCodebase, codebase.Name), 302)
}

func (c *LibraryController) extractLibraryRequestData() command.CreateCodebase {
	library := command.CreateCodebase{
		Lang:             c.GetString("appLang"),
		BuildTool:        c.GetString("buildTool"),
		Strategy:         strings.ToLower(c.GetString("strategy")),
		Type:             "library",
		JenkinsSlave:     c.GetString("jenkinsSlave"),
		JobProvisioning:  c.GetString("jobProvisioning"),
		DeploymentScript: c.GetString("deploymentScript"),
	}

	library.Versioning.Type = c.GetString("versioningType")
	startVersioningFrom := c.GetString("startVersioningFrom")
	library.Versioning.StartFrom = util.GetStringOrNil(startVersioningFrom)

	if o := OtherLanguage; library.Lang == OtherLanguage {
		library.Framework = &o
	}

	if library.Strategy == strings.ToLower(ImportStrategy) {
		library.GitServer = c.GetString("gitServer")
		gitRepoPath := c.GetString("gitRelativePath")
		library.GitUrlPath = &gitRepoPath
		library.Name = path.Base(*library.GitUrlPath)
	} else {
		library.Name = c.GetString("nameOfApp")
		library.GitServer = "gerrit"
	}

	repoUrl := c.GetString("gitRepoUrl")
	if repoUrl != "" {
		library.Repository = &command.Repository{
			Url: repoUrl,
		}

		isRepoPrivate, _ := c.GetBool("isRepoPrivate", false)
		if isRepoPrivate {
			library.Repository.Login = c.GetString("repoLogin")
			library.Repository.Password = c.GetString("repoPassword")
		}
	}

	vcsLogin := c.GetString("vcsLogin")
	vcsPassword := c.GetString("vcsPassword")
	if vcsLogin != "" && vcsPassword != "" {
		library.Vcs = &command.Vcs{
			Login:    vcsLogin,
			Password: vcsPassword,
		}
	}
	library.Username = c.Ctx.Input.Session("username").(string)
	return library
}

func validateLibraryRequestData(library command.CreateCodebase) *validation2.ErrMsg {
	valid := validation.Validation{}

	_, err := valid.Valid(library)

	if library.Strategy == strings.ToLower(ImportStrategy) {
		valid.Match(library.GitUrlPath, regexp.MustCompile("^\\/.*$"), "Spec.GitUrlPath")
	}

	if library.Strategy == "clone" && library.Repository != nil {
		_, err = valid.Valid(library.Repository)

		isAvailable := util.IsGitRepoAvailable(library.Repository.Url, library.Repository.Login, library.Repository.Password)

		if !isAvailable {
			err := &validation.Error{Key: "repository", Message: "Repository doesn't exist or invalid login and password."}
			valid.Errors = append(valid.Errors, err)
		}
	}

	if library.Vcs != nil {
		_, err = valid.Valid(library.Vcs)
	}

	if err != nil {
		return &validation2.ErrMsg{"An internal error has occurred on server while validating autotest's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &validation2.ErrMsg{string(validation2.CreateErrorResponseBody(valid)), http.StatusBadRequest}
}

func logLibraryRequestData(library command.CreateCodebase) {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Request data to create codebase CR is valid. name=%s, strategy=%s, lang=%s, buildTool=%s",
		library.Name, library.Strategy, library.Lang, library.BuildTool))

	if library.Repository != nil {
		result.WriteString(fmt.Sprintf(", repositoryUrl=%s, repositoryLogin=%s", library.Repository.Url, library.Repository.Login))
	}

	if library.Vcs != nil {
		result.WriteString(fmt.Sprintf(", vcsLogin=%s", library.Vcs.Login))
	}

	log.Info(result.String())
}
