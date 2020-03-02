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

package routers

import (
	"edp-admin-console/context"
	"edp-admin-console/controllers"
	"edp-admin-console/controllers/auth"
	cdPipeController "edp-admin-console/controllers/cd-pipeline"
	"edp-admin-console/filters"
	"edp-admin-console/k8s"
	"edp-admin-console/repository"
	edpComponentRepo "edp-admin-console/repository/edp-component"
	"edp-admin-console/service"
	"edp-admin-console/service/cd_pipeline"
	cbs "edp-admin-console/service/codebasebranch"
	edpComponentService "edp-admin-console/service/edp-component"
	"edp-admin-console/util"
	"github.com/astaxie/beego"
	"log"
)

const (
	integrationStrategies = "integrationStrategies"
	buildTools            = "buildTools"
	versioningTypes       = "versioningTypes"
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
		beego.Router("/auth/callback", &auth.AuthController{}, "get:Callback")
		beego.InsertFilter("/admin/*", beego.BeforeRouter, filters.AuthFilter)
		beego.InsertFilter("/api/v1/edp/*", beego.BeforeRouter, filters.AuthRestFilter)
		beego.InsertFilter("/admin/edp/*", beego.BeforeRouter, filters.RoleAccessControlFilter)
		beego.InsertFilter("/api/v1/edp/*", beego.BeforeRouter, filters.RoleAccessControlRestFilter)
	} else {
		beego.InsertFilter("/*", beego.BeforeRouter, filters.StubAuthFilter)
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
	branchService := cbs.CodebaseBranchService{
		Clients:                  clients,
		IReleaseBranchRepository: branchRepository,
		ICDPipelineRepository:    pipelineRepository,
		ICodebaseRepository:      codebaseRepository,
		CodebaseBranchValidation: map[string]func(string, string) ([]string, error){
			"application": pipelineRepository.GetCDPipelinesUsingApplicationAndBranch,
			"autotests":   pipelineRepository.GetCDPipelinesUsingAutotestAndBranch,
			"library":     pipelineRepository.GetCDPipelinesUsingLibraryAndBranch,
		},
	}
	codebaseService := service.CodebaseService{
		Clients:               clients,
		ICodebaseRepository:   codebaseRepository,
		ICDPipelineRepository: pipelineRepository,
		BranchService:         branchService,
	}
	pipelineService := cd_pipeline.CDPipelineService{
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
	beego.Router("/", &controllers.MainController{EDPTenantService: edpService}, "get:Index")

	integrationStrategies := util.GetValuesFromConfig(integrationStrategies)
	if integrationStrategies == nil {
		log.Fatalf("integrationStrategies config variable is empty.")
	}

	buildTools := util.GetValuesFromConfig(buildTools)
	if buildTools == nil {
		log.Fatalf("buildTools config variable is empty.")
	}

	vt := util.GetValuesFromConfig(versioningTypes)
	if vt == nil {
		log.Fatalf("versioningTypes config variable is empty.")
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
		VersioningTypes:       vt,
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
		VersioningTypes:       vt,
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
		VersioningTypes:       vt,
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

	cpc := cdPipeController.CDPipelineController{
		CodebaseService:   codebaseService,
		PipelineService:   pipelineService,
		EDPTenantService:  edpService,
		BranchService:     branchService,
		ThirdPartyService: thirdPartyService,
		EDPComponent:      ecs,
	}

	cbc := controllers.BranchController{
		BranchService:   branchService,
		CodebaseService: codebaseService,
	}

	adminEdpNamespace := beego.NewNamespace("/admin/edp",
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
		beego.NSRouter("/codebase", &cc, "post:Delete"),
		beego.NSRouter("/codebase/branch/delete", &cbc, "post:Delete"),
		beego.NSRouter("/stage", &cpc, "post:DeleteCDStage"),
		beego.NSRouter("/cd-pipeline/delete", &cpc, "post:DeleteCDPipeline"),
		beego.NSRouter("/codebase/:codebaseName/branch", &cbc, "post:CreateCodebaseBranch"),

		beego.NSRouter("/library/overview", &lc, "get:GetLibraryListPage"),
		beego.NSRouter("/library/create", &lc, "get:GetCreatePage"),
		beego.NSRouter("/library", &lc, "post:Create"),
		beego.NSRouter("/service/overview", &controllers.ThirdPartyServiceController{ThirdPartyService: thirdPartyService}, "get:GetServicePage"),
	)
	beego.AddNamespace(adminEdpNamespace)

	apiV1EdpNamespace := beego.NewNamespace("/api/v1/edp",
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "post:CreateCodebase"),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebases"),
		beego.NSRouter("/codebase/:codebaseName", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "get:GetCodebase"),
		beego.NSRouter("/vcs", &ec, "get:GetVcsIntegrationValue"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetCDPipelineByName"),
		beego.NSRouter("/cd-pipeline/:pipelineName/stage/:stageName", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "get:GetStage"),
		beego.NSRouter("/cd-pipeline", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "post:CreateCDPipeline"),
		beego.NSRouter("/cd-pipeline/:name", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "put:UpdateCDPipeline"),
		beego.NSRouter("/codebase", &controllers.CodebaseRestController{CodebaseService: codebaseService}, "delete:Delete"),
		beego.NSRouter("/stage", &controllers.CDPipelineRestController{CDPipelineService: pipelineService}, "delete:DeleteCDStage"),
	)
	beego.AddNamespace(apiV1EdpNamespace)

	apiV1Namespace := beego.NewNamespace("/api/v1",
		beego.NSRouter("/storage-class", &controllers.OpenshiftRestController{ClusterService: clusterService}, "get:GetAllStorageClasses"),
		beego.NSRouter("/repository/available", &controllers.RepositoryRestController{}, "post:IsGitRepoAvailable"),
	)
	beego.AddNamespace(apiV1Namespace)
}
