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

package controllers

import (
	"edp-admin-console/controllers/validation"
	"edp-admin-console/models/command"
	edperror "edp-admin-console/models/error"
	"edp-admin-console/service/cd_pipeline"
	dberror "edp-admin-console/util/error/db-errors"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

type CDPipelineRestController struct {
	beego.Controller
	CDPipelineService cd_pipeline.CDPipelineService
}

func (c *CDPipelineRestController) Prepare() {
	c.EnableXSRF = false
}

func (c *CDPipelineRestController) GetCDPipelineByName() {
	pipelineName := c.GetString(":name")
	cdPipeline, err := c.CDPipelineService.GetCDPipelineByName(pipelineName)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if cdPipeline == nil {
		nonAppMsg := fmt.Sprintf("Please check CD Pipeline name. It seems there's not %s pipeline.", pipelineName)
		http.Error(c.Ctx.ResponseWriter, nonAppMsg, http.StatusNotFound)
		return
	}

	c.Data["json"] = cdPipeline
	c.ServeJSON()
}

func (c *CDPipelineRestController) GetStage() {
	pipelineName := c.GetString(":pipelineName")
	stageName := c.GetString(":stageName")

	stage, err := c.CDPipelineService.GetStage(pipelineName, stageName)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if stage == nil {
		http.Error(c.Ctx.ResponseWriter, "Please check request data.", http.StatusNotFound)
		return
	}

	c.Data["json"] = stage
	c.ServeJSON()
}

func (c *CDPipelineRestController) CreateCDPipeline() {
	var cdPipelineCreateCommand command.CDPipelineCommand
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&cdPipelineCreateCommand)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	errMsg := validation.ValidateCDPipelineRequest(cdPipelineCreateCommand)
	if errMsg != nil {
		log.Error("Failed to validate request data", zap.String("err", errMsg.Message))
		http.Error(c.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}
	log.Info("Request data is receieved to create CD pipeline",
		zap.String("pipeline", cdPipelineCreateCommand.Name),
		zap.Any("applications", cdPipelineCreateCommand.Applications),
		zap.Any("stages", cdPipelineCreateCommand.Stages),
		zap.Any("services", cdPipelineCreateCommand.ThirdPartyServices))

	_, pipelineErr := c.CDPipelineService.CreatePipeline(cdPipelineCreateCommand)
	if pipelineErr != nil {

		switch pipelineErr.(type) {
		case *edperror.CDPipelineExistsError:
			http.Error(c.Ctx.ResponseWriter, fmt.Sprintf("cd pipeline %v is already exists", cdPipelineCreateCommand.Name), http.StatusFound)
			return
		case *edperror.NonValidRelatedBranchError:
			http.Error(c.Ctx.ResponseWriter, fmt.Sprintf("one or more applications have non valid branches: %v", cdPipelineCreateCommand.Applications), http.StatusBadRequest)
			return
		default:
			http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
}

func (c *CDPipelineRestController) UpdateCDPipeline() {
	var pipelineUpdateCommand command.CDPipelineCommand
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&pipelineUpdateCommand)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	pipelineUpdateCommand.Name = c.GetString(":name")

	errMsg := validation.ValidateCDPipelineUpdateRequestData(pipelineUpdateCommand)
	if errMsg != nil {
		log.Error("Request data is not valid", zap.String("err", errMsg.Message))
		http.Error(c.Ctx.ResponseWriter, errMsg.Message, http.StatusBadRequest)
		return
	}
	log.Info("Request data is received to update CD pipeline",
		zap.String("pipeline", pipelineUpdateCommand.Name),
		zap.Any("applications", pipelineUpdateCommand.Applications))

	err = c.CDPipelineService.UpdatePipeline(pipelineUpdateCommand)
	if err != nil {

		switch err.(type) {
		case *edperror.CDPipelineDoesNotExistError:
			http.Error(c.Ctx.ResponseWriter, fmt.Sprintf("cd pipeline %v doesn't exist", pipelineUpdateCommand.Name), http.StatusNotFound)
			return
		case *edperror.NonValidRelatedBranchError:
			http.Error(c.Ctx.ResponseWriter, fmt.Sprintf("one or more applications have non valid branches: %v", pipelineUpdateCommand.Name), http.StatusNotFound)
			return
		default:
			http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	c.Ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
}

func (c *CDPipelineRestController) DeleteCDStage() {
	var sc command.DeleteStageCommand
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&sc); err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("request to delete cd stage has been retrieved",
		zap.String("pipeline", sc.CDPipelineName),
		zap.String("stage", sc.Name))
	if err := c.CDPipelineService.DeleteCDStage(sc.CDPipelineName, sc.Name); err != nil {
		if dberror.StageErrorOccurred(err) {
			serr := err.(dberror.RemoveStageRestriction)
			log.Error(serr.Message, zap.Error(err))
			http.Error(c.Ctx.ResponseWriter, serr.Message, http.StatusConflict)
			return
		}
		log.Error("delete process is failed", zap.Error(err))
		http.Error(c.Ctx.ResponseWriter, "delete process is failed", http.StatusInternalServerError)
		return
	}
	log.Debug("delete cd stage method is finished",
		zap.String("pipeline", sc.CDPipelineName),
		zap.String("stage", sc.Name))
	location := fmt.Sprintf("%s/%s", c.Ctx.Input.URL(), uuid.NewV4().String())
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Ctx.Output.Header("Location", location)
}
