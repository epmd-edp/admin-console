/*
 * Copyright 2019 EPAM Systems.
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
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	ec "edp-admin-console/service/edp-component"
	"edp-admin-console/util"
	"edp-admin-console/util/consts"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"log"
	"strings"
)

type CodebaseController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	EDPTenantService service.EDPTenantService
	BranchService    service.CodebaseBranchService
	GitServerService service.GitServerService
	EDPComponent     ec.EDPComponentService
}

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
		log.Printf("An error has occurred while creating link to Git Server: %v", err)
		c.Abort("500")
		return
	}

	codebase.CodebaseBranch = addCodebaseBranchInProgressIfAny(codebase.CodebaseBranch, c.GetString(paramWaitingForBranch))

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Codebase"] = codebase
	c.Data["Type"] = codebase.Type
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
	c.TplName = "codebase_overview.html"
}

func addCodebaseBranchInProgressIfAny(branches []*query.CodebaseBranch, branchInProgress string) []*query.CodebaseBranch {
	if branchInProgress != "" {
		for _, branch := range branches {
			if branch.Name == branchInProgress {
				return branches
			}
		}

		log.Println("Adding branch " + branchInProgress + " which is going to be created to the list.")
		branch := query.CodebaseBranch{
			Name:   branchInProgress,
			Status: "inactive",
		}
		branches = append(branches, &branch)
	}
	return branches
}

func (c CodebaseController) createBranchLinks(codebase query.Codebase, tenant string) error {
	if codebase.Strategy == strings.ToLower(ImportStrategy) {
		return c.createLinksForGitProvider(codebase, tenant)
	}
	return c.createLinksForGerritProvider(codebase, tenant)
}

func (c CodebaseController) createLinksForGitProvider(codebase query.Codebase, tenant string) error {
	w := beego.AppConfig.String("dnsWildcard")
	g, err := c.GitServerService.GetGitServer(*codebase.GitServer)
	if err != nil {
		return err
	}

	if g == nil {
		return errors.New(fmt.Sprintf("unexpected behaviour. couldn't find %v GitServer in DB", *codebase.GitServer))
	}

	for i, b := range codebase.CodebaseBranch {
		codebase.CodebaseBranch[i].VCSLink = util.CreateGitLink(g.Hostname, *codebase.GitProjectPath, b.Name)
		j := fmt.Sprintf("https://%s-%s-edp-cicd.%s", consts.Jenkins, tenant, w)
		codebase.CodebaseBranch[i].CICDLink = util.CreateCICDApplicationLink(j, codebase.Name, b.Name)
	}

	return nil
}

func (c CodebaseController) createLinksForGerritProvider(codebase query.Codebase, tenant string) error {
	cj, err := c.EDPComponent.GetEDPComponent(consts.Jenkins)
	if err != nil {
		return err
	}

	cg, err := c.EDPComponent.GetEDPComponent(consts.Gerrit)
	if err != nil {
		return err
	}

	for i, b := range codebase.CodebaseBranch {
		codebase.CodebaseBranch[i].VCSLink = util.CreateGerritLink(cg.Url, codebase.Name, b.Name)
		codebase.CodebaseBranch[i].CICDLink = util.CreateCICDApplicationLink(cj.Url, codebase.Name, b.Name)
	}

	return nil
}
