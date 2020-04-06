package main

import (
	_ "edp-admin-console/routers"
	_ "edp-admin-console/template_function"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
