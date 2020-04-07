package context

import (
	"github.com/astaxie/beego"
)

var (
	Namespace  = beego.AppConfig.String("cicdNamespace")
	EDPVersion = beego.AppConfig.String("edpVersion")
	Tenant     = beego.AppConfig.String("edpName")
	BasePath   = beego.AppConfig.String("basePath")
)
