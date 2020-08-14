package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
)

type ThirdPartyServiceController struct {
	beego.Controller
	ThirdPartyService service.ThirdPartyService
}

func (s *ThirdPartyServiceController) GetServicePage() {
	services, err := s.ThirdPartyService.GetAllServices()
	if err != nil {
		s.Abort("500")
		return
	}

	s.Data["EDPVersion"] = context.EDPVersion
	s.Data["Username"] = s.Ctx.Input.Session("username")
	s.Data["Services"] = services
	s.Data["Type"] = "services"
	s.Data["BasePath"] = context.BasePath
	s.Data["DiagramPageEnabled"] = context.DiagramPageEnabled
	s.TplName = "service.html"
}
