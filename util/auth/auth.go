package auth

import (
	"edp-admin-console/util"
	"github.com/astaxie/beego"
)

func IsAdmin(contextRoles []string) bool {
	if contextRoles == nil {
		return false
	}
	return util.Contains(contextRoles, beego.AppConfig.String("adminRole"))
}
