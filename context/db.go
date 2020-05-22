/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package context

import (
	"edp-admin-console/models/query"
	"edp-admin-console/service/logger"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

func InitDb() {
	err := orm.RegisterDriver("postgres", orm.DRPostgres)
	checkErr(err)

	pgUser := beego.AppConfig.String("pgUser")
	pgPassword := beego.AppConfig.String("pgPassword")
	pgHost := beego.AppConfig.String("pgHost")
	pgDatabase := beego.AppConfig.String("pgDatabase")
	pgPort := beego.AppConfig.String("pgPort")
	pgSchema := Tenant

	params := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s search_path=%s sslmode=disable",
		pgUser, pgPassword, pgHost, pgPort, pgDatabase, pgSchema)

	err = orm.RegisterDataBase("default", "postgres", params)
	checkErr(err)
	log.Info("Connection to database is established.",
		zap.String("host", pgHost),
		zap.String("port", pgPort))

	db, err := orm.GetDB("default")
	checkErr(err)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		pgDatabase, driver)
	checkErr(err)
	err = m.Up()
	checkErr(err)
	debug, err := beego.AppConfig.Bool("ormDebug")
	if err != nil {
		log.Info("Cannot read orm debug config. Set to false", zap.Error(err))
		debug = false
	}
	orm.Debug = debug
	orm.RegisterModel(new(query.Codebase), new(query.ActionLog), new(query.CodebaseBranch), new(query.ThirdPartyService),
		new(query.CDPipeline), new(query.JobProvisioning), new(query.Stage), new(query.QualityGate), new(query.ApplicationsToPromote),
		new(query.CodebaseDockerStream), new(query.GitServer), new(query.JenkinsSlave),
		new(query.EDPComponent), new(query.JiraServer))
}

func checkErr(err error) {
	if err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	if err.Error() == "no change" {
		log.Info("Warning from db migration", zap.Error(err))
	} else {
		log.Fatal("An error has occurred during migration", zap.Error(err))
	}
}
