package util

import (
	"github.com/astaxie/beego"
	"log"
	"strings"
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

func GetStringOrNil(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}