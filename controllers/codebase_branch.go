package controllers

import (
	"edp-admin-console/context"
	validation2 "edp-admin-console/controllers/validation"
	"edp-admin-console/models/command"
	"edp-admin-console/service"
	cbs "edp-admin-console/service/codebasebranch"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"go.uber.org/zap"
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
		log.Error("Failed to validate request data", zap.String("err", errMsg.Message))
		c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%s/overview", context.BasePath, appName), 302)
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

	log.Debug("Request data to create CR for codebase branch is valid",
		zap.String("branch", branchInfo.Name),
		zap.String("commit hash", branchInfo.Commit))

	exist := c.CodebaseService.ExistCodebaseAndBranch(appName, branchInfo.Name)

	if exist {
		c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%s/overview?errorExistingBranch=%s#branchExistsModal", context.BasePath,
			appName, url.PathEscape(branchInfo.Name)), 302)
		return
	}

	cb, err := c.BranchService.CreateCodebaseBranch(branchInfo, appName)
	if err != nil {
		c.Abort("500")
		return
	}

	log.Info("BranchRelease resource is saved into cluster", zap.String("name", cb.Name))
	c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%s/overview?%s=%s#branchSuccessModal", context.BasePath, appName,
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
	log.Debug("delete codebase branch method is invoked",
		zap.String("codebase name", cn),
		zap.String("branch name", bn))
	if err := c.BranchService.Delete(cn, bn); err != nil {
		if dberror.CodebaseBranchErrorOccurred(err) {
			cberr := err.(dberror.RemoveCodebaseBranchRestriction)
			f := beego.NewFlash()
			f.Error(cberr.Message)
			f.Store(&c.Controller)
			log.Error(cberr.Message, zap.Error(err))
			c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%v/overview?name=%v#branchIsUsedSuccessModal", context.BasePath, cn, bn), 302)
			return
		}
		log.Error("delete process is failed", zap.Error(err))
		c.Abort("500")
		return
	}
	log.Info("delete codebase branch method is finished",
		zap.String("codebase name", cn),
		zap.String("branch name", bn))
	c.Redirect(fmt.Sprintf("%s/admin/edp/codebase/%v/overview?name=%v#branchDeletedSuccessModal", context.BasePath, cn, bn), 302)
}
