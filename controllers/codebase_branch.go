package controllers

import (
	validation2 "edp-admin-console/controllers/validation"
	"edp-admin-console/models/command"
	"edp-admin-console/service"
	cbs "edp-admin-console/service/codebasebranch"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"net/http"
	"net/url"
	"regexp"
)

type BranchController struct {
	beego.Controller
	CodebaseService service.CodebaseService
	BranchService   cbs.CodebaseBranchService
}

func (c *BranchController) CreateCodebaseBranch() {
	branchInfo := c.extractCodebaseBranchRequestData()
	appName := c.GetString(":codebaseName")
	errMsg := validCodebaseBranchRequestData(branchInfo)
	if errMsg != nil {
		log.Info("Failed to validate request data", "err", errMsg.Message)
		c.Redirect(fmt.Sprintf("/admin/edp/codebase/%s/overview", appName), 302)
		return
	}

	if branchInfo.Release {
		mv := c.GetString("masterVersion")
		mp := c.GetString("snapshotStaticField")
		masterVersion := util.GetVersionOrNil(mv, mp)
		err := c.BranchService.UpdateCodebaseBranch(appName, "master", masterVersion)
		if err != nil {
			c.Abort("500")
			return
		}
	}

	log.Info("Request data to create CR for codebase branch is valid",
		"branch", branchInfo.Name, "commit hash", branchInfo.Commit)

	exist := c.CodebaseService.ExistCodebaseAndBranch(appName, branchInfo.Name)

	if exist {
		c.Redirect(fmt.Sprintf("/admin/edp/codebase/%s/overview?errorExistingBranch=%s#branchExistsModal",
			appName, url.PathEscape(branchInfo.Name)), 302)
		return
	}

	cb, err := c.BranchService.CreateCodebaseBranch(branchInfo, appName)
	if err != nil {
		c.Abort("500")
		return
	}

	log.Info("BranchRelease resource is saved into cluster", "name", cb.Name)
	c.Redirect(fmt.Sprintf("/admin/edp/codebase/%s/overview?%s=%s#branchSuccessModal", appName,
		paramWaitingForBranch, url.PathEscape(branchInfo.Name)), 302)
}

func (c *BranchController) extractCodebaseBranchRequestData() command.CreateCodebaseBranch {
	cb := command.CreateCodebaseBranch{
		Name:     c.GetString("name"),
		Commit:   c.GetString("commit"),
		Username: c.Ctx.Input.Session("username").(string),
	}

	vf := c.GetString("version")
	px := c.GetString("versioningPostfix")
	cb.Version = util.GetVersionOrNil(vf, px)

	cb.Build = &consts.DefaultBuildNumber

	r, _ := c.GetBool("releaseBranch", false)
	cb.Release = r

	return cb
}

func validCodebaseBranchRequestData(requestData command.CreateCodebaseBranch) *validation2.ErrMsg {
	valid := validation.Validation{}
	_, err := valid.Valid(requestData)

	if len(requestData.Commit) != 0 {
		valid.Match(requestData.Commit, regexp.MustCompile("\\b([a-f0-9]{40})\\b"), "Commit.Match")
	}

	if err != nil {
		return &validation2.ErrMsg{"An internal error has occurred on server while validating branch's form fields.",
			http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &validation2.ErrMsg{string(validation2.CreateErrorResponseBody(valid)), http.StatusBadRequest}
}

func (c *BranchController) Delete() {
	cn := c.GetString("codebase-name")
	bn := c.GetString("name")
	rl := log.WithValues("codebase name", cn, "branch name", bn)
	rl.V(2).Info("delete codebase branch method is invoked")
	if err := c.BranchService.Delete(cn, bn); err != nil {
		if dberror.CodebaseBranchErrorOccurred(err) {
			cberr := err.(dberror.RemoveCodebaseBranchRestriction)
			f := beego.NewFlash()
			f.Error(cberr.Message)
			f.Store(&c.Controller)
			log.Error(err, cberr.Message)
			c.Redirect(fmt.Sprintf("/admin/edp/codebase/%v/overview?name=%v#branchIsUsedSuccessModal", cn, bn), 302)
			return
		}
		log.Error(err, "delete process is failed")
		c.Abort("500")
		return
	}
	rl.V(2).Info("delete codebase branch method is finished")
	c.Redirect(fmt.Sprintf("/admin/edp/codebase/%v/overview?name=%v#branchDeletedSuccessModal", cn, bn), 302)
}
