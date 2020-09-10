package context

import (
	"github.com/astaxie/beego"
	"strconv"
)

var (
	Namespace          = beego.AppConfig.String("cicdNamespace")
	EDPVersion         = beego.AppConfig.String("edpVersion")
	Tenant             = beego.AppConfig.String("edpName")
	BasePath           = beego.AppConfig.String("basePath")
	DiagramPageEnabled = convertToBool()
)

func convertToBool() bool {
	s := beego.AppConfig.String("diagramPageEnabled")
	if s == "" {
		return false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return b
}
