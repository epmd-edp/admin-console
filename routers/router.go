/*
 * Copyright 2019 EPAM Systems.
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

package routers

import (
	"edp-admin-console/context"
	"edp-admin-console/controllers"
	"edp-admin-console/filters"
	"edp-admin-console/k8s"
	"edp-admin-console/repository"
	edpComponentRepo "edp-admin-console/repository/edp-component"
	"edp-admin-console/service"
	edpComponentService "edp-admin-console/service/edp-component"
	"edp-admin-console/util"
	"fmt"
	"log"

	"github.com/astaxie/beego"
)

const (
	integrationStrategies = "integrationStrategies"
	buildTools            = "buildTools"
	testReportTools       = "testReportTools"
	deploymentScript      = "deploymentScript"

	CreateStrategy = "Create"
)

func init() {
	log.Printf("Start application in %s mode with %s EDP version...", beego.AppConfig.String("runmode"), context.EDPVersion)
	authEnabled, err := beego.AppConfig.Bool("keycloakAuthEnabled")
	if err != nil {
		log.Printf("Cannot read property keycloakAuthEnabled: %v. Set default: true", err)
		authEnabled = true
	}

	if authEnabled {
		context.InitAuth()
		beego.Router(fmt.Sprintf("%s/auth/callback", context.BasePath), &controllers.AuthController{}, "get:Callback")
		beego.InsertFilter(fmt.Sprintf("%s/admin/*", context.BasePath), beego.BeforeRouter, filters.AuthFilter)
		beego.InsertFilter(fmt.Sprintf("%s/api/v1/edp/*", context.BasePath), beego.BeforeRouter, filters.AuthRestFilter)
		beego.InsertFilter(fmt.Sprintf("%s/admin/edp/*", context.BasePath), beego.BeforeRouter, filters.RoleAccessControlFilter)
		beego.InsertFilter(fmt.Sprintf("%s/api/v1/edp/*", context.BasePath), beego.BeforeRouter, filters.RoleAccessControlRestFilter)
	} else {
		beego.InsertFilter(fmt.Sprintf("%s/*", context.BasePath), beego.BeforeRouter, filters.StubAuthFilter)
	}

	dbEnable, err := beego.AppConfig.Bool("dbEnabled")
	if err != nil {
		log.Printf("Cannot read property dbEnabled: %v. Set default: true", err)
		dbEnable = true
	}

	if dbEnable {
		context.InitDb()
	}

	clients := k8s.CreateOpenShiftClients()
	codebaseRepository := repository.CodebaseRepository{}
	branchRepository := repository.CodebaseBranchRepository{}
	pipelineRepository := repository.CDPipelineRepository{}
	serviceRepository := repository.ServiceCatalogRepository{}
	gitServerRepository := repository.GitServerRepository{}
	sr := repository.SlaveRepository{}
	pr := repository.JobProvisioning{}
	ecr := edpComponentRepo.EDPComponent{}

	ecs := edpComponentService.EDPComponentService{IEDPComponent: ecr}
	edpService := service.EDPTenantService{Clients: clients}
	clusterService := service.ClusterService{Clients: clients}
	branchService := service.CodebaseBranchService{Clients: clients, IReleaseBranchRepository: branchRepository}
	codebaseService := service.CodebaseService{
		Clients:             clients,
		ICodebaseRepository: codebaseRepository,
		BranchService:       branchService,
	}
	pipelineService := service.CDPipelineService{
		Clients:               clients,
		ICDPipelineRepository: pipelineRepository,
		CodebaseService:       codebaseService,
		BranchService:         branchService,
		EDPComponent:          ecs,
	}
	thirdPartyService := service.ThirdPartyService{IServiceCatalogRepository: serviceRepository}
	gitServerService := service.GitServerService{IGitServerRepository: gitServerRepository}
	ss := service.SlaveService{ISlaveRepository: sr}
	ps := service.JobProvisioning{IJobProvisioningRepository: pr}

	beego.ErrorController(&controllers.ErrorController{})
	beego.Router(fmt.Sprintf("%s/", context.BasePath), &controllers.MainController{EDPTenantService: edpService}, "get:Index")
	beego.SetStaticPath(fmt.Sprintf("%s/static", context.BasePath), "static")

	integrationStrategies := util.GetValuesFromConfig(integrationStrategies)
	if integrationStrategies == nil {
		log.Fatalf("integrationStrategies config variable is empty.")
	}

	buildTools := util.GetValuesFromConfig(buildTools)
	if buildTools == nil {
		log.Fatalf("buildTools config variable is empty.")
	}

	testReportTools := util.GetValuesFromConfig(testReportTools)
	if testReportTools == nil {
		log.Fatalf("testReportTools config variable is empty.")
	}

	ds := util.GetValuesFromConfig(deploymentScript)
	if ds == nil {
		log.Fatalf("deploymentScript config variable is empty.")
	}

	is := make([]string, len(integrationStrategies))
	copy(is, integrationStrategies)

	appc := controllers.ApplicationController{
		CodebaseService:  codebaseService,
		EDPTenantService: edpService,
		BranchService:    branchService,
		GitServerService: gitServerService,
		SlaveService:     ss,
		JobProvisioning:  ps,

		IntegrationStrategies: is,
		BuildTools:            buildTools,
		DeploymentScript:      ds,
	}

	autis := make([]string, len(integrationStrategies))
	copy(autis, integrationStrategies)

	autc := controllers.AutotestsController{
		EDPTenantService: edpService,
		CodebaseService:  codebaseService,
		BranchService:    branchService,
		GitServerService: gitServerService,
		SlaveService:     ss,
		JobProvisioning:  ps,

		IntegrationStrategies: util.RemoveElByValue(autis, CreateStrategy),
		BuildTools:            buildTools,
		TestReportTools:       testReportTools,
		DeploymentScript:      ds,
	}

	lc := controllers.LibraryController{
		EDPTenantService: edpService,
		CodebaseService:  codebaseService,
		GitServerService: gitServerService,
		SlaveService:     ss,
		JobProvisioning:  ps,

		IntegrationStrategies: is,
		BuildTools:            buildTools,
		DeploymentScript:      ds,
	}

	ec := controllers.EDPTenantController{
		EDPTenantService: edpService,
		EDPComponent:     ecs,
	}

	cc := controllers.CodebaseController{
		CodebaseService:  codebaseService,
		EDPTenantService: edpService,
		BranchService:    branchService,
		GitServerService: gitServerService,
		EDPComponent:     ecs,
	}

	cpc := controllers.CDPipelineController{
		CodebaseService:   codebaseService,
		PipelineService:   pipelineService,
		EDPTenantService:  edpService,
		BranchService:     branchService,
		ThirdPartyService: thirdPartyService,
		EDPComponent:      ecs,
		JobProvisioning:   ps,
	}

	adminEdpNamespace := beego.NewNamespace(fmt.Sprintf("%s/admin/edp", context.BasePath),
		beego.NSRouter("/overview", &ec, "get:GetEDPComponents"),
		beego.NSRouter("/application/overview", &appc, "get:GetApplicationsOverviewPage"),
		beego.NSRouter("/application/create", &appc, "get:GetCreateApplicationPage"),
		beego.NSRouter("/application", &appc, "post:CreateApplication"),

		beego.NSRouter("/cd-pipeline/overview", &cpc, "get:GetContinuousDeliveryPage"),
		beego.NSRouter("/cd-pipeline/create", &cpc, "get:GetCreateCDPipelinePage"),
		beego.NSRouter("/cd-pipeline/:name/update", &cpc, "get:GetEditCDPipelinePage"),
		beego.NSRouter("/cd-pipeline", &cpc, "post:CreateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:name/update", &cpc, "post:UpdateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:pipelineName/overview", &cpc, "get:GetCDPipelineOverviewPage"),
		beego.NSRouter("/autotest/overview", &autc, "get:GetAutotestsOverviewPage"),
		beego.NSRouter("/autotest/create", &autc, "get:GetCreateAutotestsPage"),
		beego.NSRouter("/autotest", &autc, "post:CreateAutotests"),

		beego.NSRouter("/codebase/:codebaseName/overview", &cc, "get:GetCodebaseOverviewPage"),
		beego.NSRouter("/codebase/:codebaseName/branch", &controllers.BranchController{BranchService: branchService, CodebaseService: codebaseService}, "post:CreateCodebaseBranch"),

		beego.NSRouter("/library/overview", &lc, "get:GetLibraryListPage"),
		beego.NSRouter("/library/create", &lc, "get:GetCreatePage"),
		beego.NSRouter("/library", &lc, "post:Create"),
		beego.NSRouter("/service/overview", &controllers.ThirdPartyServiceController{ThirdPartyService: thirdPartyService}, "get:GetServicePage"),
	)
	beego.AddNamespace(adminEdpNamespace)

	apiV1EdpNamespace := beego.NewNamespace(fmt.Sprintf("%s/api/v1/edp", context.BasePath),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "post:CreateCodebase"),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebases"),
		beego.NSRouter("/codebase/:codebaseName", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebase"),
		beego.NSRouter("/vcs", &ec, "get:GetVcsIntegrationValue"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetCDPipelineByName"),
		beego.NSRouter("/cd-pipeline/:pipelineName/stage/:stageName", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetStage"),
		beego.NSRouter("/cd-pipeline", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "post:CreateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "put:UpdateCDPipeline"),
	)
	beego.AddNamespace(apiV1EdpNamespace)

	apiV1Namespace := beego.NewNamespace(fmt.Sprintf("%s/api/v1", context.BasePath),
		beego.NSRouter("/storage-class", &controllers.OpenshiftRestController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository/available", &controllers.RepositoryRestController{}, "post:IsGitRepoAvailable"),
	)
	beego.AddNamespace(apiV1Namespace)
}
