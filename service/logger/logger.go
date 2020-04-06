package logger

import (
	"github.com/astaxie/beego"
	"go.uber.org/zap"
)

const (
	debugVerbosity = "debugVerbosity"
)

func GetLogger() *zap.Logger {
	if b, _ := beego.AppConfig.Bool(debugVerbosity); b == true {
		log, _ := zap.NewDevelopment()
		return log
	}
	log, _ := zap.NewProduction()
	return log
}
