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
	b, err := strconv.ParseBool(beego.AppConfig.String("diagramPageEnabled"))
	if err != nil {
		panic(err)
	}
	return b
}
