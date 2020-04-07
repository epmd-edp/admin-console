package controllers

import (
	"edp-admin-console/context"
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (this *ErrorController) Error500() {
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["BasePath"] = context.BasePath
	this.TplName = "error/error_500.html"
}

func (this *ErrorController) Error403() {
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["BasePath"] = context.BasePath
	this.TplName = "error/error_403.html"
}

func (this *ErrorController) Error404() {
	this.Data["EDPVersion"] = context.EDPVersion
	this.Data["Username"] = this.Ctx.Input.Session("username")
	this.Data["BasePath"] = context.BasePath
	this.TplName = "error/error_404.html"
}
