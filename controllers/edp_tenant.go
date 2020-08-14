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
	"edp-admin-console/service"
	ec "edp-admin-console/service/edp-component"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
)

type EDPTenantController struct {
	beego.Controller
	EDPTenantService service.EDPTenantService
	EDPComponent     ec.EDPComponentService
}

func (c *EDPTenantController) GetEDPComponents() {
	comp, err := c.EDPComponent.GetEDPComponents()
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["InputURL"] = strings.TrimSuffix(c.Ctx.Input.URL(), "/"+context.Tenant)
	c.Data["EDPTenantName"] = context.Tenant
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["EDPComponents"] = comp
	c.Data["Type"] = "overview"
	c.Data["BasePath"] = context.BasePath
	c.Data["DiagramPageEnabled"] = context.DiagramPageEnabled
	c.TplName = "edp_components.html"
}

func (this *EDPTenantController) GetVcsIntegrationValue() {
	isVcsEnabled, err := this.EDPTenantService.GetVcsIntegrationValue()

	if err != nil {
		if err.Error() == "NOT_FOUND" {
			http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Data["json"] = isVcsEnabled
	this.ServeJSON()
}
