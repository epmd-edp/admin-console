package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/models/command"
	"edp-admin-console/service"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type BranchController struct {
	beego.Controller
	CodebaseService service.CodebaseService
	BranchService   service.CodebaseBranchService
}

func (c *BranchController) CreateCodebaseBranch() {
	branchInfo := c.extractCodebaseBranchRequestData()
	errMsg := validCodebaseBranchRequestData(branchInfo)
	appName := c.GetString(":codebaseName")
	if errMsg != nil {
		log.Printf("Failed to validate request data: %s", errMsg.Message)
		c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%s/overview", context.BasePath, appName), 302)
		return
	}
	log.Printf("Request data to create CR for codebase branch is valid: {Branch Name: %s, Commit Hash: %s}", branchInfo.Name, branchInfo.Commit)

	exist := c.CodebaseService.ExistCodebaseAndBranch(appName, branchInfo.Name)

	if exist {
		c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%s/overview?errorExistingBranch=%s#branchExistsModal", context.BasePath, appName, branchInfo.Name), 302)
		return
	}

	cb, err := c.BranchService.CreateCodebaseBranch(branchInfo, appName)
	if err != nil {
		c.Abort("500")
		return
	}

	log.Printf("BranchRelease resource is saved into k8s: %s", cb)
	c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%s/overview?%s=%s#branchSuccessModal", context.BasePath, appName, paramWaitingForBranch, branchInfo.Name), 302)
}

func (c *BranchController) extractCodebaseBranchRequestData() command.CreateCodebaseBranch {
	return command.CreateCodebaseBranch{
		Name:     c.GetString("name"),
		Commit:   c.GetString("commit"),
		Username: c.Ctx.Input.Session("username").(string),
	}
}

func validCodebaseBranchRequestData(requestData command.CreateCodebaseBranch) *ErrMsg {
	valid := validation.Validation{}
	_, err := valid.Valid(requestData)

	if len(requestData.Commit) != 0 {
		valid.Match(requestData.Commit, regexp.MustCompile("\\b([a-f0-9]{40})\\b"), "Commit.Match")
	}

	if err != nil {
		return &ErrMsg{"An internal error has occurred on server while validating branch's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &ErrMsg{string(createErrorResponseBody(valid)), http.StatusBadRequest}
}
