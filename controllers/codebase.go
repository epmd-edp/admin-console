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
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	cbs "edp-admin-console/service/codebasebranch"
	ec "edp-admin-console/service/edp-component"
	"edp-admin-console/util"
	"edp-admin-console/util/auth"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"errors"
	"fmt"
	"html/template"

	"github.com/astaxie/beego"
	"go.uber.org/zap"
)

type CodebaseController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	EDPTenantService service.EDPTenantService
	BranchService    cbs.CodebaseBranchService
	GitServerService service.GitServerService
	EDPComponent     ec.EDPComponentService
}

const (
	paramWaitingForBranch = "waitingforbranch"
)

func (c *CodebaseController) GetCodebaseOverviewPage() {
	codebaseName := c.GetString(":codebaseName")
	codebase, err := c.CodebaseService.GetCodebaseByName(codebaseName)
	if err != nil {
		c.Abort("500")
		return
	}

	if codebase == nil {
		c.Abort("404")
		return
	}

	err = c.createBranchLinks(*codebase, context.Tenant)
	if err != nil {
		log.Error("an error has occurred while creating link to Git Server", zap.Error(err))
		c.Abort("500")
		return
	}

	codebase.CodebaseBranch = addCodebaseBranchInProgressIfAny(codebase.CodebaseBranch, c.GetString(paramWaitingForBranch))

	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["error"] != "" {
		c.Data["ErrorBranch"] = flash.Data["error"]
	}
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Codebase"] = codebase
	c.Data["Type"] = codebase.Type
	contextRoles := c.GetSession("realm_roles").([]string)
	c.Data["HasRights"] = auth.IsAdmin(contextRoles)
	switch codebase.Type {
	case "application":
		{
			c.Data["TypeCaption"] = "Application"
			c.Data["TypeSingular"] = "application"
		}
	case "autotests":
		{
			c.Data["TypeCaption"] = "Autotests codebase"
			c.Data["TypeSingular"] = "autotests codebase"
		}
	case "library":
		{
			c.Data["TypeCaption"] = "Library"
			c.Data["TypeSingular"] = "library"
		}
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["BasePath"] = context.BasePath
	c.Data["DiagramPageEnabled"] = context.DiagramPageEnabled
	c.TplName = "codebase_overview.html"
}

func addCodebaseBranchInProgressIfAny(branches []*query.CodebaseBranch, branchInProgress string) []*query.CodebaseBranch {
	if branchInProgress != "" {
		for _, branch := range branches {
			if branch.Name == branchInProgress {
				return branches
			}
		}

		log.Debug("Adding branch which is going to be created to the list.",
			zap.String("name", branchInProgress))
		branch := query.CodebaseBranch{
			Name:   branchInProgress,
			Status: "inactive",
		}
		branches = append(branches, &branch)
	}
	return branches
}

func (c CodebaseController) createBranchLinks(codebase query.Codebase, tenant string) error {
	if codebase.Strategy == consts.ImportStrategy {
		return c.createLinksForGitProvider(codebase, tenant)
	}
	return c.createLinksForGerritProvider(codebase, tenant)
}

func (c CodebaseController) createLinksForGitProvider(codebase query.Codebase, tenant string) error {
	g, err := c.GitServerService.GetGitServer(*codebase.GitServer)
	if err != nil {
		return err
	}
	if g == nil {
		return errors.New(fmt.Sprintf("unexpected behaviour. couldn't find %v GitServer in DB", *codebase.GitServer))
	}

	jc, err := c.EDPComponent.GetEDPComponent(consts.Jenkins)
	if err != nil {
		return err
	}
	if jc == nil {
		return fmt.Errorf("jenkin link can't be created for %v codebase because of edp-component %v is absent in DB",
			codebase.Name, consts.Jenkins)
	}

	for i, b := range codebase.CodebaseBranch {
		codebase.CodebaseBranch[i].VCSLink = util.CreateGitLink(g.Hostname, *codebase.GitProjectPath, b.Name)
		codebase.CodebaseBranch[i].CICDLink = util.CreateCICDApplicationLink(jc.Url, codebase.Name,
			util.ProcessNameToKubernetesConvention(b.Name))
	}

	return nil
}

func (c CodebaseController) createLinksForGerritProvider(codebase query.Codebase, tenant string) error {
	cj, err := c.EDPComponent.GetEDPComponent(consts.Jenkins)
	if err != nil {
		return err
	}

	if cj == nil {
		return fmt.Errorf("jenkin link can't be created for %v codebase because of edp-component %v is absent in DB",
			codebase.Name, consts.Jenkins)
	}

	cg, err := c.EDPComponent.GetEDPComponent(consts.Gerrit)
	if err != nil {
		return err
	}

	if cg == nil {
		return fmt.Errorf("gerrit link can't be created for %v codebase because of edp-component %v is absent in DB",
			codebase.Name, consts.Gerrit)
	}

	for i, b := range codebase.CodebaseBranch {
		codebase.CodebaseBranch[i].VCSLink = util.CreateGerritLink(cg.Url, codebase.Name, b.Name)
		codebase.CodebaseBranch[i].CICDLink = util.CreateCICDApplicationLink(cj.Url, codebase.Name,
			util.ProcessNameToKubernetesConvention(b.Name))
	}

	return nil
}

func (c *CodebaseController) Delete() {
	flash := beego.NewFlash()
	cn := c.GetString("name")
	log.Debug("delete codebase method is invoked", zap.String("name", cn))
	ct := c.GetString("codebase-type")
	if err := c.CodebaseService.Delete(cn, ct); err != nil {
		if dberror.CodebaseIsUsed(err) {
			cerr := err.(dberror.CodebaseIsUsedByCDPipeline)
			flash.Error(cerr.Message)
			flash.Store(&c.Controller)
			log.Error(cerr.Message, zap.Error(err))
			c.Redirect(createCodebaseIsUsedURL(cerr.Codebase, ct), 302)
			return
		}
		log.Error("delete process is failed", zap.Error(err))
		c.Abort("500")
		return
	}
	log.Info("delete codebase method is finished", zap.String("name", cn))
	c.Redirect(createCodebaseIsDeletedURL(cn, ct), 302)
}

func createCodebaseIsUsedURL(codebaseName, codebaseType string) string {
	if codebaseType == consts.Autotest {
		codebaseType = "autotest"
	}
	return fmt.Sprintf("%s/admin/edp/%v/overview?codebase=%v#codebaseIsUsed", context.BasePath,
		codebaseType, codebaseName)
}

func createCodebaseIsDeletedURL(codebaseName, codebaseType string) string {
	if codebaseType == consts.Autotest {
		codebaseType = "autotest"
	}
	return fmt.Sprintf("%s/admin/edp/%v/overview?codebase=%v#codebaseIsDeleted", context.BasePath, codebaseType, codebaseName)
}

func (c *CodebaseController) GetEditCodebasePage() {
	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["error"] != "" {
		c.Data["CodebaseUpdateError"] = flash.Data["error"]
	}

	n := c.GetString(":name")
	log.Debug("start executing GetEditCodebasePage method", zap.String("name", n))

	codebase, err := c.CodebaseService.GetCodebaseByName(n)
	if err != nil {
		log.Error("couldn't get codebase from db", zap.Error(err))
		c.Abort("500")
		return
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["BasePath"] = context.BasePath
	c.Data["Codebase"] = codebase
	c.Data["Type"] = codebase.Type
	c.Data["DiagramPageEnabled"] = context.DiagramPageEnabled
	c.TplName = "edit_codebase.html"
}

func (c *CodebaseController) Update() {
	flash := beego.NewFlash()
	name := c.GetString("name")
	log.Debug("start executing Update method", zap.String("name", name))

	cc := command.UpdateCodebaseCommand{
		Name:               name,
		CommitMessageRegex: c.GetString("commitMessagePattern"),
		TicketNameRegex:    c.GetString("ticketNamePattern"),
	}

	errMsg := validation.ValidateCodebaseUpdateRequestData(cc)
	if errMsg != nil {
		log.Error("Codebase update request data is invalid", zap.String("err", errMsg.Message))
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("%v/admin/edp/codebase/%v/update", context.BasePath, cc.Name), 302)
		return
	}

	codebase, err := c.CodebaseService.Update(cc)
	if err != nil {
		log.Error("couldn't update codebase", zap.Error(err))
		c.Abort("500")
		return
	}

	c.Redirect(fmt.Sprintf("%v/admin/edp/%v/overview#codebaseUpdateSuccessModal",
		context.BasePath, getType(codebase.Spec.Type)), 302)
}

func getType(codebaseType string) string {
	if codebaseType == "autotests" {
		return "autotest"
	}
	return codebaseType
}
