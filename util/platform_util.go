package util

import (
	"strings"

	"github.com/astaxie/beego"
	"log"
)

func GetValuesFromConfig(name string) []string {
	values := beego.AppConfig.String(name)
	if values == "" {
		log.Printf("'%v' env variable is empty.", name)
		return nil
	}

	s := strings.Split(values, ",")
	log.Printf("Fetched data from %v env variable: %v", name, s)
	return s
}

func GetStringP(val string) *string {
	return &val
}
