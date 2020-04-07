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

package pipeline

import (
	"edp-admin-console/context"
	"edp-admin-console/controllers/validation"
	"edp-admin-console/models"
	"edp-admin-console/models/command"
	edperror "edp-admin-console/models/error"
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	"edp-admin-console/service/cd_pipeline"
	cbs "edp-admin-console/service/codebasebranch"
	ec "edp-admin-console/service/edp-component"
	"edp-admin-console/service/logger"
	"edp-admin-console/service/platform"
	"edp-admin-console/util"
	"edp-admin-console/util/auth"
	"edp-admin-console/util/consts"
	dberror "edp-admin-console/util/error/db-errors"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	edppipelinesv1alpha1 "github.com/epmd-edp/cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"regexp"
	"sort"
)

var log = logger.GetLogger()

type CDPipelineController struct {
	beego.Controller
	CodebaseService   service.CodebaseService
	PipelineService   cd_pipeline.CDPipelineService
	EDPTenantService  service.EDPTenantService
	BranchService     cbs.CodebaseBranchService
	ThirdPartyService service.ThirdPartyService
	EDPComponent      ec.EDPComponentService
}

const (
	paramWaitingForCdPipeline = "waitingforcdpipeline"
)

func (c *CDPipelineController) GetContinuousDeliveryPage() {
	applications, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		Status: query.Active,
		Type:   query.App,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	branches, err := c.BranchService.GetCodebaseBranchesByCriteria(query.CodebaseBranchCriteria{
		Status: "active",
	})
	if err != nil {
		c.Abort("500")
		return
	}

	cdPipelines, err := c.PipelineService.GetAllPipelines(query.CDPipelineCriteria{})
	if err != nil {
		c.Abort("500")
		return
	}

	err = c.createJenkinsLinks(cdPipelines)
	if err != nil {
		c.Abort("500")
		return
	}

	cdPipelines = addCdPipelineInProgressIfAny(cdPipelines, c.GetString(paramWaitingForCdPipeline))

	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}
	contextRoles := c.GetSession("realm_roles").([]string)
	c.Data["ActiveApplicationsAndBranches"] = len(applications) > 0 && len(branches) > 0
	c.Data["CDPipelines"] = cdPipelines
	c.Data["Applications"] = applications
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["HasRights"] = auth.IsAdmin(contextRoles)
	c.Data["Type"] = "delivery"
	c.Data["BasePath"] = context.BasePath
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "continuous_delivery.html"
}

func (c *CDPipelineController) GetCreateCDPipelinePage() {
	flash := beego.ReadFromRequest(&c.Controller)
	apps, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.App,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	groovyLibs, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.Library,
		Language:     "groovy-pipeline",
	})
	if err != nil {
		c.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}

	services, err := c.ThirdPartyService.GetAllServices()
	if err != nil {
		c.Abort("500")
		return
	}

	autotests, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.Autotests,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	autotests = filterAutotestsWithActiveBranches(autotests)

	c.Data["Services"] = services
	c.Data["Apps"] = apps
	c.Data["GroovyLibs"] = groovyLibs
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Type"] = "delivery"
	c.Data["Autotests"] = autotests
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["BasePath"] = context.BasePath
	c.TplName = "create_cd_pipeline.html"
}

func filterAutotestsWithActiveBranches(autotests []*query.Codebase) []*query.Codebase {
	if autotests == nil {
		return autotests
	}

	for i, autotest := range autotests {
		if len(autotest.CodebaseBranch) == 0 {
			autotests = append(autotests[:i], autotests[i+1:]...)
		}
	}

	return autotests
}

func (c *CDPipelineController) GetEditCDPipelinePage() {
	flash := beego.ReadFromRequest(&c.Controller)
	pipelineName := c.GetString(":name")

	cdPipeline, err := c.PipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		c.Abort("500")
		return
	}

	applications, err := c.CodebaseService.GetCodebasesByCriteria(query.CodebaseCriteria{
		BranchStatus: query.Active,
		Status:       query.Active,
		Type:         query.App,
	})
	if err != nil {
		c.Abort("500")
		return
	}

	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}

	c.Data["CDPipeline"] = cdPipeline
	c.Data["Apps"] = applications
	c.Data["Type"] = "delivery"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["BasePath"] = context.BasePath
	c.TplName = "edit_cd_pipeline.html"
}

func (c *CDPipelineController) UpdateCDPipeline() {
	flash := beego.NewFlash()
	appNameCheckboxes := c.GetStrings("app")
	pipelineName := c.GetString(":name")

	pipelineUpdateCommand := command.CDPipelineCommand{
		Name:                 pipelineName,
		Applications:         c.convertApplicationWithBranchesData(appNameCheckboxes),
		ApplicationToApprove: c.getApplicationsToPromoteFromRequest(appNameCheckboxes),
	}

	errMsg := validation.ValidateCDPipelineUpdateRequestData(pipelineUpdateCommand)
	if errMsg != nil {
		log.Info("Request data is not valid", zap.String("err", errMsg.Message))
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/%s/update", context.BasePath, pipelineName), 302)
		return
	}
	log.Debug("Request data is received to update CD pipeline",
		zap.String("pipeline", pipelineName),
		zap.Any("applications", pipelineUpdateCommand.Applications),
		zap.Any("stages", pipelineUpdateCommand.Stages),
		zap.Any("services", pipelineUpdateCommand.ThirdPartyServices))

	err := c.PipelineService.UpdatePipeline(pipelineUpdateCommand)
	if err != nil {

		switch err.(type) {
		case *edperror.CDPipelineDoesNotExistError:
			flash.Error(fmt.Sprintf("cd pipeline %v doesn't exist", pipelineName))
			flash.Store(&c.Controller)
			c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/%s/update", context.BasePath, pipelineName), http.StatusFound)
			return
		case *edperror.NonValidRelatedBranchError:
			flash.Error(fmt.Sprintf("one or more applications have non valid branches: %v", pipelineUpdateCommand.Applications))
			flash.Store(&c.Controller)
			c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/%s/update", context.BasePath, pipelineName), http.StatusBadRequest)
			return
		default:
			c.Abort("500")
			return
		}
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/overview#cdPipelineEditSuccessModal", context.BasePath), 302)
}

func (c *CDPipelineController) CreateCDPipeline() {
	flash := beego.NewFlash()
	appNameCheckboxes := c.GetStrings("app")
	pipelineName := c.GetString("pipelineName")
	serviceCheckboxes := c.GetStrings("service")
	stages := retrieveStagesFromRequest(c)

	cdPipelineCreateCommand := command.CDPipelineCommand{
		Name:                 pipelineName,
		Applications:         c.convertApplicationWithBranchesData(appNameCheckboxes),
		ThirdPartyServices:   serviceCheckboxes,
		Stages:               stages,
		ApplicationToApprove: c.getApplicationsToPromoteFromRequest(appNameCheckboxes),
		Username:             c.Ctx.Input.Session("username").(string),
	}

	errMsg := validation.ValidateCDPipelineRequest(cdPipelineCreateCommand)
	if errMsg != nil {
		log.Error("Request data is not valid", zap.String("err", errMsg.Message))
		flash.Error(errMsg.Message)
		flash.Store(&c.Controller)
		c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/create", context.BasePath), 302)
		return
	}
	log.Debug("Request data is received to create CD pipeline",
		zap.String("pipeline", pipelineName),
		zap.Any("applications", cdPipelineCreateCommand.Applications),
		zap.Any("stages", cdPipelineCreateCommand.Stages),
		zap.Any("services", cdPipelineCreateCommand.ThirdPartyServices))

	_, pipelineErr := c.PipelineService.CreatePipeline(cdPipelineCreateCommand)
	if pipelineErr != nil {

		switch pipelineErr.(type) {
		case *edperror.CDPipelineExistsError:
			flash.Error(fmt.Sprintf("cd pipeline %v is already exists", cdPipelineCreateCommand.Name))
			flash.Store(&c.Controller)
			c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/create", context.BasePath), http.StatusFound)
			return
		case *edperror.NonValidRelatedBranchError:
			flash.Error(fmt.Sprintf("one or more applications have non valid branches: %v", cdPipelineCreateCommand.Applications))
			flash.Store(&c.Controller)
			c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/create", context.BasePath), http.StatusBadRequest)
			return
		default:
			c.Abort("500")
			return
		}

	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Redirect(fmt.Sprintf("%s/admin/edp/cd-pipeline/overview?%s=%s#cdPipelineSuccessModal", context.BasePath, paramWaitingForCdPipeline, pipelineName), 302)
}

func (c *CDPipelineController) GetCDPipelineOverviewPage() {
	pipelineName := c.GetString(":pipelineName")

	cdPipeline, err := c.PipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		c.Abort("500")
		return
	}

	if cdPipeline == nil {
		c.Abort("404")
		return
	}

	if err := c.createOneJenkinsLink(cdPipeline); err != nil {
		c.Abort("500")
		return
	}

	if err := c.createDockerImageLinks(cdPipeline.CodebaseDockerStream); err != nil {
		c.Abort("500")
		return
	}

	if err := c.createPlatformLinks(cdPipeline.Stage, cdPipeline.Name); err != nil {
		log.Error("an error has occurred while creating platform links", zap.Error(err))
		c.Abort("500")
		return
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if flash.Data["success"] != "" {
		c.Data["Success"] = true
	}
	if flash.Data["error"] != "" {
		c.Data["Error"] = flash.Data["error"]
	}
	contextRoles := c.GetSession("realm_roles").([]string)
	c.Data["CDPipeline"] = cdPipeline
	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Type"] = "delivery"
	c.Data["IsOpenshift"] = platform.IsOpenshift()
	c.Data["BasePath"] = context.BasePath
	c.Data["HasRights"] = auth.IsAdmin(contextRoles)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "cd_pipeline_overview.html"
}

func retrieveStagesFromRequest(this *CDPipelineController) []command.CDStageCommand {
	var stages []command.CDStageCommand

	for index, stageName := range this.GetStrings("stageName") {
		stgSrc := edppipelinesv1alpha1.Source{}
		name := this.GetString(stageName + "-pipelineLibraryName")
		if name == "default" {
			stgSrc.Type = "default"
		} else {
			stgSrc.Type = "library"
			stgSrc.Library = edppipelinesv1alpha1.Library{
				Name:   name,
				Branch: this.GetString(stageName + "-pipelineLibraryBranch"),
			}
		}
		stageRequest := command.CDStageCommand{
			Name:        stageName,
			Description: this.GetString(stageName + "-stageDesc"),
			TriggerType: this.GetString(stageName + "-triggerType"),
			Source:      stgSrc,
			Order:       index,
		}

		for _, stepName := range this.GetStrings(stageName + "-stageStepName") {
			qualityGateRequest := edppipelinesv1alpha1.QualityGate{
				QualityGateType: this.GetString(stageName + "-" + stepName + "-stageQualityGateType"),
				StepName:        stepName,
			}

			if qualityGateRequest.QualityGateType == "autotests" {
				autotestName := this.GetString(stageName + "-" + stepName + "-stageAutotests")
				qualityGateRequest.AutotestName = &autotestName
				stageName := this.GetString(stageName + "-" + stepName + "-stageBranch")
				qualityGateRequest.BranchName = &stageName
			}

			stageRequest.QualityGates = append(stageRequest.QualityGates, qualityGateRequest)
		}

		stages = append(stages, stageRequest)
	}

	sort.SliceStable(stages, func(i, j int) bool {
		return stages[i].Order < stages[j].Order
	})

	log.Info("Stages are fetched from request", zap.Any("stages", stages))
	return stages
}

func (c *CDPipelineController) convertApplicationWithBranchesData(appNameCheckboxes []string) []models.CDPipelineApplicationCommand {
	var applicationWithBranches []models.CDPipelineApplicationCommand
	for _, appName := range appNameCheckboxes {
		applicationWithBranches = append(applicationWithBranches, models.CDPipelineApplicationCommand{
			ApplicationName:   appName,
			InputDockerStream: c.GetString(appName),
		})
	}
	return applicationWithBranches
}

func addCdPipelineInProgressIfAny(cdPipelines []*query.CDPipeline, pipelineInProgress string) []*query.CDPipeline {
	if pipelineInProgress != "" {
		for _, pipeline := range cdPipelines {
			if pipeline.Name == pipelineInProgress {
				return cdPipelines
			}
		}

		log.Debug("Adding CD Pipeline which is going to be created to the list",
			zap.String("name", pipelineInProgress))
		pipeline := query.CDPipeline{
			Name:   pipelineInProgress,
			Status: "inactive",
		}
		cdPipelines = append(cdPipelines, &pipeline)
	}
	return cdPipelines
}

func (c *CDPipelineController) getApplicationsToPromoteFromRequest(appNameCheckboxes []string) []string {
	var applicationsToPromote []string
	for _, appName := range appNameCheckboxes {
		promote, _ := c.GetBool(appName+"-promote", false)
		if promote {
			applicationsToPromote = append(applicationsToPromote, appName)
		}
	}
	return applicationsToPromote
}

func (c *CDPipelineController) createOneJenkinsLink(cdPipeline *query.CDPipeline) error {
	edc, err := c.EDPComponent.GetEDPComponent(consts.Jenkins)
	if err != nil {
		return err
	}

	if edc == nil {
		return fmt.Errorf("jenkins link can't be created for %v cd pipeline because of edp-component %v is absent in DB",
			cdPipeline.Name, consts.Jenkins)
	}

	cdPipeline.JenkinsLink = util.CreateCICDPipelineLink(edc.Url, cdPipeline.Name)
	log.Info("Created CD Pipeline Jenkins link", zap.String("jenkins link", cdPipeline.JenkinsLink))
	return nil
}

func (c *CDPipelineController) createDockerImageLinks(stream []*query.CodebaseDockerStream) error {
	if platform.IsOpenshift() {
		return c.createNativeDockerImageLinks(stream)
	}
	return c.createNonNativeDockerImageLinks(stream)
}

func (c *CDPipelineController) createNativeDockerImageLinks(s []*query.CodebaseDockerStream) error {
	co, err := c.EDPComponent.GetEDPComponent(consts.Openshift)
	if err != nil {
		return err
	}

	if co == nil {
		return fmt.Errorf("openshift link can't be created because of edp-component %v is absent in DB", consts.Openshift)
	}

	cj, err := c.EDPComponent.GetEDPComponent(consts.Jenkins)
	if err != nil {
		return err
	}

	if cj == nil {
		return fmt.Errorf("jenkins link can't be created because of edp-component %v is absent in DB", consts.Jenkins)
	}

	for i, v := range s {
		s[i].CICDLink = util.CreateCICDApplicationLink(cj.Url, v.CodebaseBranch.Codebase.Name,
			util.ProcessNameToKubernetesConvention(v.CodebaseBranch.Name))
		s[i].ImageLink = util.CreateNativeDockerStreamLink(co.Url, context.Namespace,
			getImageStreamName(v.CodebaseBranch.Release, v.OcImageStreamName))
	}

	return nil
}

func getImageStreamName(release bool, imageStream string) string {
	if release {
		return regexp.MustCompile(`\/([0-9]\d*)\.([0-9]\d*)`).ReplaceAllString(imageStream, "-$1-$2")
	}
	return imageStream
}

func (c *CDPipelineController) createNonNativeDockerImageLinks(s []*query.CodebaseDockerStream) error {
	cd, err := c.EDPComponent.GetEDPComponent(consts.DockerRegistry)
	if err != nil {
		return err
	}

	if cd == nil {
		return fmt.Errorf("docker registry link can't be created because of edp-component %v is absent in DB",
			consts.DockerRegistry)
	}

	cj, err := c.EDPComponent.GetEDPComponent(consts.Jenkins)
	if err != nil {
		return err
	}

	if cj == nil {
		return fmt.Errorf("jenkins link can't be created because of edp-component %v is absent in DB", consts.Jenkins)
	}

	for i, v := range s {
		s[i].ImageLink = util.CreateNonNativeDockerStreamLink(cd.Url,
			getImageStreamName(v.CodebaseBranch.Release, v.OcImageStreamName))
		s[i].CICDLink = util.CreateCICDApplicationLink(cj.Url, v.CodebaseBranch.Codebase.Name,
			util.ProcessNameToKubernetesConvention(v.CodebaseBranch.Name))
	}

	return nil
}

func (c *CDPipelineController) createPlatformLinks(stages []*query.Stage, cdPipelineName string) error {
	log.Debug("Start creating Platform links forCD Pipeline", zap.String("name", cdPipelineName))

	if len(stages) == 0 {
		return errors.New("stages can't be an empty or nil")
	}

	if platform.IsOpenshift() {
		return c.createNativePlatformLinks(stages, cdPipelineName)
	}
	return c.createNonNativePlatformLinks(stages, cdPipelineName)
}

func (c *CDPipelineController) createNativePlatformLinks(stages []*query.Stage, cdPipelineName string) error {
	log.Debug("Start creating Openshift Platform links forCD Pipeline", zap.String("name", cdPipelineName))

	edc, err := c.EDPComponent.GetEDPComponent(consts.Openshift)
	if err != nil {
		return err
	}

	if edc == nil {
		return fmt.Errorf("openshift link can't be created because of edp-component %v is absent in DB", consts.Openshift)
	}

	for i, v := range stages {
		stages[i].PlatformProjectLink = util.CreateNativeProjectLink(edc.Url, v.PlatformProjectName)
	}

	return nil
}

func (c *CDPipelineController) createNonNativePlatformLinks(stages []*query.Stage, cdPipelineName string) error {
	log.Debug("Start creating Kubernetes Platform links for CD Pipeline", zap.String("name", cdPipelineName))

	edc, err := c.EDPComponent.GetEDPComponent(consts.Kubernetes)
	if err != nil {
		return err
	}

	if edc == nil {
		log.Debug("Creating Kubernetes Platform links has been skipped forCD Pipeline",
			zap.String("name", cdPipelineName))
		return nil
	}

	for i, v := range stages {
		stages[i].PlatformProjectLink = util.CreateNonNativeProjectLink(edc.Url, v.PlatformProjectName)
	}

	return nil
}

func (c *CDPipelineController) createJenkinsLinks(cdPipelines []*query.CDPipeline) error {
	if len(cdPipelines) == 0 {
		log.Info("There're no CD Pipelines. Generating Jenkins links are skipped.")
		return nil
	}

	edc, err := c.EDPComponent.GetEDPComponent(consts.Jenkins)
	if err != nil {
		return err
	}

	if edc == nil {
		return fmt.Errorf("jenkins links can't be created because of edp-component %v is absent in DB", consts.Jenkins)
	}

	for index, pipeline := range cdPipelines {
		cdPipelines[index].JenkinsLink = util.CreateCICDPipelineLink(edc.Url, pipeline.Name)
		log.Debug("Created Jenkins link", zap.String("link", pipeline.JenkinsLink))
	}
	return nil
}

func (c CDPipelineController) DeleteCDStage() {
	flash := beego.NewFlash()
	pn := c.GetString("pipeline")
	sn := c.GetString("name")
	o, err := c.GetInt("order")
	if err != nil {
		c.Abort("500")
		return
	}
	log.Debug("request to delete cd stage has been retrieved",
		zap.String("pipeline", pn),
		zap.String("stage", sn),
		zap.Int("order", o))

	if o == 0 {
		if err := c.PipelineService.DeleteCDPipeline(pn); err != nil {
			if dberror.CDPipelineErrorOccurred(err) {
				perr := err.(dberror.RemoveCDPipelineRestriction)
				flash.Error(perr.Message)
				flash.Store(&c.Controller)
				log.Error(perr.Message, zap.Error(err))
				c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/overview?name=%v#cdPipelineIsUsedAsSource", pn), 302)
				return
			}
			log.Error("cd pipeline delete process is failed", zap.Error(err))
			c.Abort("500")
			return
		}
		log.Debug("delete cd stage method is finished",
			zap.String("pipeline", pn),
			zap.String("stage", sn))
		c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/overview?name=%v#cdPipelineDeletedSuccessModal", pn), 302)
	}

	if err := c.PipelineService.DeleteCDStage(pn, sn); err != nil {
		if dberror.StageErrorOccurred(err) {
			serr := err.(dberror.RemoveStageRestriction)
			flash.Error(serr.Message)
			flash.Store(&c.Controller)
			log.Error(serr.Message, zap.Error(err))
			c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/%v/overview?stage=%v#stageIsUsedAsSource", pn, sn), 302)
			return
		}
		log.Error("cd stage delete process is failed", zap.Error(err))
		c.Abort("500")
		return
	}
	log.Debug("delete cd stage method is finished",
		zap.String("pipeline", pn),
		zap.String("stage", sn))
	c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/%v/overview?stage=%v#stageSuccessModal", pn, sn), 302)
}

func (c CDPipelineController) DeleteCDPipeline() {
	n := c.GetString("name")
	log.Debug("request to delete cd pipeline has been received", zap.String("name", n))
	if err := c.PipelineService.DeleteCDPipeline(n); err != nil {
		flash := beego.NewFlash()
		if dberror.CDPipelineErrorOccurred(err) {
			perr := err.(dberror.RemoveCDPipelineRestriction)
			flash.Error(perr.Message)
			flash.Store(&c.Controller)
			log.Error(perr.Message, zap.Error(err))
			c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/overview?name=%v#cdPipelineIsUsedAsSource", n), 302)
			return
		}
		log.Error("cd pipeline delete process is failed", zap.Error(err))
		c.Abort("500")
		return
	}
	log.Debug("delete cd pipeline method is finished", zap.String("name", n))
	c.Redirect(fmt.Sprintf("/admin/edp/cd-pipeline/overview?name=%v#cdPipelineDeletedSuccessModal", n), 302)
}
