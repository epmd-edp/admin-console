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

package repository

import (
	"edp-admin-console/models"
	"edp-admin-console/models/dto"
	"edp-admin-console/models/query"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type ICDPipelineRepository interface {
	GetCDPipelineByName(pipelineName string) (*query.CDPipeline, error)
	GetCDPipelines(criteria query.CDPipelineCriteria) ([]*query.CDPipeline, error)
	GetStage(cdPipelineName, stageName string) (*models.StageView, error)
	GetCodebaseAndBranchName(codebaseId, branchId int) (*dto.CodebaseBranchDTO, error)
	GetQualityGates(stageId int64) ([]query.QualityGate, error)
	GetCDPipelinesUsingApplication(codebaseName string) ([]string, error)
	GetCDPipelinesUsingAutotest(codebaseName string) ([]string, error)
	GetCDPipelinesUsingLibrary(codebaseName string) ([]string, error)
	SelectMaxOrderBetweenStages(pipeName string) (*int, error)
	SelectStageOrder(pipeName, stageName string) (*int, error)
	SelectCDPipelinesUsingInputStageAsSource(pipeName, stageName string) ([]string, error)
	GetCDPipelinesUsingApplicationAndBranch(codebase, branch string) ([]string, error)
	GetCDPipelinesUsingAutotestAndBranch(codebase, branch string) ([]string, error)
	GetCDPipelinesUsingLibraryAndBranch(codebase, branch string) ([]string, error)
	GetAllCodebaseDockerStreams() ([]string, error)
}

const (
	SelectStageByCDPipelineAndStageNames = "select cs.name stage_name, " +
		"	cp.name pipeline_name, " +
		"	cs.description, " +
		"	cs.trigger_type, " +
		"	c.name app_name, " +
		"	cb.name branch_name, " +
		"	in_cds.oc_image_stream_name  input_image_stream, " +
		"	out_cds.oc_image_stream_name output_image_stream, " +
		"	cs.\"order\", " +
		"	cs.id " +
		"from cd_stage cs " +
		"	left join stage_codebase_docker_stream scds on cs.id = scds.cd_stage_id " +
		"	left join codebase_docker_stream in_cds on scds.input_codebase_docker_stream_id = in_cds.id " +
		"	left join codebase_docker_stream out_cds on scds.output_codebase_docker_stream_id = out_cds.id " +
		"	left join cd_pipeline_docker_stream cpds on cs.cd_pipeline_id = cpds.cd_pipeline_id " +
		"	left join codebase_docker_stream cds on cpds.codebase_docker_stream_id = cds.id " +
		"	left join codebase_branch cb on cb.id = cds.codebase_branch_id " +
		"   left join codebase_branch in_cbb on in_cbb.id = in_cds.codebase_branch_id " +
		"	left join codebase c on cb.codebase_id = c.id " +
		"   left join codebase in_c on in_cbb.codebase_id = in_c.id " +
		"	left join cd_pipeline cp on cs.cd_pipeline_id = cp.id " +
		"where cp.name = ? " +
		"  and cs.name = ? " +
		"  and c.id = in_c.id;"
	SelectCodebaseAndBranchName = "select c.name codebase_name, cb.name codebase_branch_name " +
		"	from codebase c " +
		"left join codebase_branch cb on c.id = cb.codebase_id " +
		"where c.type = 'autotests' " +
		"  and c.id = ? " +
		"and cb.id = ? ;"
	SelectDockerStreamName = "select cds.id, cds.oc_image_stream_name " +
		"	from cd_pipeline cp " +
		"left join cd_pipeline_docker_stream cpds on cp.id = cpds.cd_pipeline_id " +
		"left join codebase_docker_stream cds on cpds.codebase_docker_stream_id = cds.id " +
		"where cp.name = ? and cds.codebase_branch_id = ? ;"
	SelectCDPipelineByCodebaseName = "select cp.name " +
		"	from cd_pipeline cp " +
		"left join cd_pipeline_docker_stream cpds on cp.id = cpds.cd_pipeline_id " +
		"left join codebase_docker_stream cds on cpds.codebase_docker_stream_id = cds.id " +
		"left join codebase_branch cb on cds.codebase_branch_id = cb.id " +
		"left join codebase c on cb.codebase_id = c.id " +
		"	where c.name = ? ;"
	SelectCDPipelineByAutotestName = "select cp.name " +
		"	from cd_pipeline cp " +
		"left join cd_stage cs on cp.id = cs.cd_pipeline_id " +
		"left join quality_gate_stage qgs on cs.id = qgs.cd_stage_id " +
		"left join codebase c on qgs.codebase_id = c.id " +
		"where qgs.quality_gate = 'autotests' " +
		"  and c.name = ? " +
		"group by cp.name;"
	SelectCDPipelineByLibraryName = "select cp.name " +
		"	from cd_pipeline cp " +
		"left join cd_stage cs on cp.id = cs.cd_pipeline_id " +
		"left join codebase_branch cb on cs.codebase_branch_id = cb.id " +
		"left join codebase c on cb.codebase_id = c.id " +
		"where c.type = 'library' " +
		"  and c.name = ? " +
		"group by cp.name ;"
	selectMaxOrderBetweenStages = "select max(cs.\"order\") " +
		"from cd_pipeline cp " +
		"left join cd_stage cs on cp.id = cs.cd_pipeline_id " +
		"where cp.name = ? ;"
	selectStageOrder = "select cs.\"order\" " +
		"from cd_pipeline cp " +
		"left join cd_stage cs on cp.id = cs.cd_pipeline_id " +
		"where cp.name = ? and cs.name = ? ;"
	selectSourceStage = "select cp_out.name " +
		"	from cd_stage cs_in " +
		"left join cd_pipeline cp_in on cs_in.cd_pipeline_id = cp_in.id " +
		"left join stage_codebase_docker_stream scds_in on cs_in.id = scds_in.cd_stage_id " +
		"left join stage_codebase_docker_stream scds_out " +
		"on scds_in.output_codebase_docker_stream_id = scds_out.input_codebase_docker_stream_id " +
		"left join cd_stage cs_out on cs_out.id = scds_out.cd_stage_id " +
		"left join cd_pipeline cp_out on cp_out.id = cs_out.cd_pipeline_id " +
		"where cp_in.name = ? " +
		"  and cs_in.name = ? " +
		"  and not cp_out.name = ? ;"
	selectCDPipelinesUsingCodebaseAndBranch = "select cp.name " +
		"	from cd_pipeline cp " +
		"left join cd_pipeline_docker_stream cpds on cp.id = cpds.cd_pipeline_id " +
		"left join codebase_docker_stream cds on cpds.codebase_docker_stream_id = cds.id " +
		"left join codebase_branch cb on cds.codebase_branch_id = cb.id " +
		"left join codebase c on cb.codebase_id = c.id " +
		"where c.name = ? " +
		"  and cb.name = ? ;"
	selectCDPipelineUsingAutotestAndBranch = "select cp.name " +
		"	from cd_pipeline cp " +
		"left join cd_stage cs on cp.id = cs.cd_pipeline_id " +
		"left join quality_gate_stage qgs on cs.id = qgs.cd_stage_id " +
		"left join codebase_branch cb on qgs.codebase_branch_id = cb.id " +
		"left join codebase c on cb.codebase_id = c.id " +
		"where qgs.quality_gate = 'autotests' " +
		"  and c.name = ? " +
		"  and cb.name = ? " +
		"group by cp.name;"
	selectCDPipelineUsingLibraryAndBranch = "select cp.name " +
		"	from cd_pipeline cp " +
		"left join cd_stage cs on cp.id = cs.cd_pipeline_id " +
		"left join codebase_branch cb on cs.codebase_branch_id = cb.id " +
		"left join codebase c on cb.codebase_id = c.id " +
		"where c.type = 'library' " +
		"  and c.name = ? " +
		"  and cb.name = ? " +
		"group by cp.name;"
	selectCodebaseDockerStream      = "select oc_image_stream_name from codebase_docker_stream;"
	selectStageCodebaseDockerStream = "select " +
		" cd_stage_id, cds.oc_image_stream_name as input, cds1.oc_image_stream_name as output" +
		" from stage_codebase_docker_stream scds " +
		" left join codebase_docker_stream cds on scds.input_codebase_docker_stream_id = cds.id" +
		" left join codebase_docker_stream cds1 on scds.output_codebase_docker_stream_id = cds1.id " +
		" where cd_stage_id = ?;"
)

type CDPipelineRepository struct {
	ICDPipelineRepository
}

func (r CDPipelineRepository) GetCDPipelineByName(pipelineName string) (*query.CDPipeline, error) {
	o := orm.NewOrm()
	cdPipeline := query.CDPipeline{Name: pipelineName}

	err := o.Read(&cdPipeline, "Name")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&cdPipeline, "CodebaseDockerStream", false, 100, 0, "Id")
	if err != nil {
		return nil, err
	}

	branches, err := loadRelatedCodebaseDockerStreams(cdPipeline.CodebaseDockerStream)
	if err != nil {
		return nil, err
	}
	cdPipeline.CodebaseBranch = branches

	if err = loadRelatedCodebases(cdPipeline.CodebaseBranch); err != nil {
		return nil, err
	}

	if err = loadRelatedDockerStreams(pipelineName, cdPipeline.CodebaseBranch); err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&cdPipeline, "Stage", false, 100, 0, "Name")
	if err != nil {
		return nil, err
	}

	if err = loadRelatedQualityGates(cdPipeline.Stage); err != nil {
		return nil, err
	}

	for _, s := range cdPipeline.Stage {
		if err := loadRelatedAutotest(s.QualityGates); err != nil {
			return nil, err
		}

		if err = loadRelatedBranch(s.QualityGates); err != nil {
			return nil, err
		}

		if err := loadRelatedSource(s); err != nil {
			return nil, err
		}
		_, err = o.LoadRelated(s, "JobProvisioning", false, 100, 0, "Name")
		if err != nil {
			return nil, err
		}
	}

	_, err = o.LoadRelated(&cdPipeline, "ThirdPartyService", false, 100, 0, "Name")
	if err != nil {
		return nil, err
	}

	err = loadRelatedActionLogForCDPipeline(&cdPipeline)
	if err != nil {
		return nil, err
	}

	return &cdPipeline, nil
}

func loadRelatedSource(s *query.Stage) error {
	if s.SourceCodebaseBranchId == nil {
		s.Source = query.Source{
			Type:    "default",
			Library: nil,
		}
		return nil
	}

	if err := loadRelatedSourceLibrary(s); err != nil {
		return err
	}

	return nil
}

func loadRelatedSourceLibrary(s *query.Stage) error {
	o := orm.NewOrm()
	b := query.CodebaseBranch{}
	err := o.QueryTable(new(query.CodebaseBranch)).
		RelatedSel().
		Filter("id", *s.SourceCodebaseBranchId).
		One(&b)
	if err != nil {
		return err
	}

	s.Source = query.Source{
		Type: "library",
		Library: &query.SourceLibrary{
			Name:   b.Codebase.Name,
			Branch: b.Name,
		},
	}

	return nil
}

func loadRelatedCodebaseDockerStreams(dockerStreams []*query.CodebaseDockerStream) ([]*query.CodebaseBranch, error) {
	var branches []*query.CodebaseBranch
	o := orm.NewOrm()

	for _, dockerStream := range dockerStreams {
		_, err := o.LoadRelated(dockerStream, "CodebaseBranch", false, 100, 0, "Id")
		if err != nil {
			return nil, err
		}
		branches = append(branches, dockerStream.CodebaseBranch)
	}

	return branches, nil
}

func loadRelatedCodebases(branches []*query.CodebaseBranch) error {
	o := orm.NewOrm()

	for _, branch := range branches {
		_, err := o.LoadRelated(branch, "Codebase", false, 100, 0, "Name")
		if err != nil {
			return err
		}
	}

	return nil
}

func loadRelatedDockerStreams(pipelineName string, branches []*query.CodebaseBranch) error {
	o := orm.NewOrm()
	for _, branch := range branches {
		var dockerStream query.CodebaseDockerStream
		err := o.Raw(SelectDockerStreamName, pipelineName, branch.Id).QueryRow(&dockerStream)
		if err != nil {
			return err
		}

		branch.CodebaseDockerStream = []*query.CodebaseDockerStream{&dockerStream}
	}

	return nil
}

func loadRelatedQualityGates(stages []*query.Stage) error {
	for i, stage := range stages {
		o := orm.NewOrm()

		_, err := o.QueryTable(new(query.QualityGate)).
			Filter("cd_stage_id", stage.Id).
			All(&stage.QualityGates)
		if err != nil {
			return err
		}

		stages[i] = stage
	}

	return nil
}

func loadRelatedAutotest(gates []query.QualityGate) error {
	for i, gate := range gates {
		if gate.QualityGateType == "autotests" {
			o := orm.NewOrm()

			codebase := query.Codebase{}
			err := o.QueryTable(new(query.Codebase)).
				Filter("id", gate.CodebaseId).
				One(&codebase)
			if err != nil {
				return err
			}

			gates[i].Autotest = &codebase
		}
	}

	return nil
}

func loadRelatedBranch(gates []query.QualityGate) error {
	for i, gate := range gates {
		if gate.QualityGateType == "autotests" {
			o := orm.NewOrm()

			branch := query.CodebaseBranch{}
			err := o.QueryTable(new(query.CodebaseBranch)).
				Filter("id", gate.CodebaseBranchId).
				One(&branch)
			if err != nil {
				return err
			}

			gates[i].Branch = &branch
		}
	}

	return nil
}

func (r CDPipelineRepository) GetCDPipelines(criteria query.CDPipelineCriteria) ([]*query.CDPipeline, error) {
	o := orm.NewOrm()
	var pipelines []*query.CDPipeline

	qs := o.QueryTable(new(query.CDPipeline))

	if criteria.Status != "" {
		qs = qs.Filter("status", criteria.Status)
	}

	_, err := qs.OrderBy("name").All(&pipelines)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	for _, p := range pipelines {
		if err = loadRelatedStage(p); err != nil {
			return nil, err
		}

		if _, err = o.LoadRelated(p, "CodebaseDockerStream"); err != nil {
			return nil, err
		}

		for _, s := range p.Stage {
			ds, err := r.getStageCodebaseDockerStream(s.Id)
			if err != nil {
				return nil, err
			}
			s.StageCodebaseDockerStream = ds
		}

		if err = loadRelatedQualityGates(p.Stage); err != nil {
			return nil, err
		}
	}

	return pipelines, nil
}

func loadRelatedStage(pipeline *query.CDPipeline) error {
	o := orm.NewOrm()
	qs := o.QueryTable(new(query.Stage))
	_, err := qs.Filter("cd_pipeline_id", pipeline.Id).
		OrderBy("Name").
		All(&pipeline.Stage)
	return err
}

func loadRelatedActionLogForCDPipeline(cdPipeline *query.CDPipeline) error {
	o := orm.NewOrm()

	_, err := o.QueryTable(new(query.ActionLog)).
		Filter("cdPipeline__cd_pipeline_id", cdPipeline.Id).
		OrderBy("LastTimeUpdate").
		Distinct().
		All(&cdPipeline.ActionLog, "LastTimeUpdate", "UserName",
			"Message", "Action", "Result")

	return err
}

func (r CDPipelineRepository) GetStage(cdPipelineName, stageName string) (*models.StageView, error) {
	o := orm.NewOrm()
	var stage models.StageView
	var maps []orm.Params

	_, err := o.Raw(SelectStageByCDPipelineAndStageNames, cdPipelineName, stageName).Values(&maps)

	if err != nil {
		return nil, err
	}

	if maps == nil {
		return nil, nil
	}

	for index, row := range maps {
		if index == 0 {
			stage.Name = row["stage_name"].(string)
			stage.CDPipeline = row["pipeline_name"].(string)
			stage.Description = row["description"].(string)
			stage.TriggerType = row["trigger_type"].(string)
			stage.Order = row["order"].(string)

			id, err := strconv.ParseInt(row["id"].(string), 10, 64)
			if err != nil {
				return nil, err
			}

			stage.Id = id
		}
		stage.Applications = append(stage.Applications, models.ApplicationStage{
			Name:       row["app_name"].(string),
			BranchName: row["branch_name"].(string),
			InputIs:    row["input_image_stream"].(string),
			OutputIs:   row["output_image_stream"].(string),
		})
	}

	return &stage, nil
}

func (CDPipelineRepository) GetCodebaseAndBranchName(codebaseId, branchId int) (*dto.CodebaseBranchDTO, error) {
	o := orm.NewOrm()

	result := dto.CodebaseBranchDTO{}
	var maps []orm.Params

	_, err := o.Raw(SelectCodebaseAndBranchName, codebaseId, branchId).Values(&maps)
	if err != nil {
		return nil, err
	}

	for _, row := range maps {
		result.AppName = row["codebase_name"].(string)
		result.BranchName = row["codebase_branch_name"].(string)
	}

	return &result, nil
}

func (CDPipelineRepository) GetQualityGates(stageId int64) ([]query.QualityGate, error) {
	o := orm.NewOrm()

	var gates []query.QualityGate
	_, err := o.QueryTable(new(query.QualityGate)).
		Filter("cd_stage_id", stageId).
		All(&gates)
	if err != nil {
		return nil, err
	}

	err = loadRelatedAutotest(gates)
	if err != nil {
		return nil, err
	}

	err = loadRelatedBranch(gates)
	if err != nil {
		return nil, err
	}

	for _, gate := range gates {
		if gate.QualityGateType == "autotests" && gate.Autotest.GitServerId != nil {
			err := loadRelatedGitServerName(gate.Autotest)
			if err != nil {
				return nil, err
			}
		}
	}

	return gates, nil
}

func (CDPipelineRepository) GetCDPipelinesUsingApplication(codebaseName string) ([]string, error) {
	o := orm.NewOrm()
	var name []string
	_, err := o.Raw(SelectCDPipelineByCodebaseName, codebaseName).QueryRows(&name)
	if err != nil {
		return nil, err
	}
	return name, nil
}

func (CDPipelineRepository) GetCDPipelinesUsingAutotest(codebaseName string) ([]string, error) {
	o := orm.NewOrm()
	var name []string
	_, err := o.Raw(SelectCDPipelineByAutotestName, codebaseName).QueryRows(&name)
	if err != nil {
		return nil, err
	}
	return name, nil
}

func (CDPipelineRepository) GetCDPipelinesUsingLibrary(codebaseName string) ([]string, error) {
	o := orm.NewOrm()
	var name []string
	_, err := o.Raw(SelectCDPipelineByLibraryName, codebaseName).QueryRows(&name)
	if err != nil {
		return nil, err
	}
	return name, nil
}

func (CDPipelineRepository) SelectMaxOrderBetweenStages(pipeName string) (*int, error) {
	o := orm.NewOrm()
	var c int
	if err := o.Raw(selectMaxOrderBetweenStages, pipeName).QueryRow(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (CDPipelineRepository) SelectStageOrder(pipeName, stageName string) (*int, error) {
	o := orm.NewOrm()
	var c int
	if err := o.Raw(selectStageOrder, pipeName, stageName).QueryRow(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (CDPipelineRepository) SelectCDPipelinesUsingInputStageAsSource(pipeName, stageName string) ([]string, error) {
	o := orm.NewOrm()
	var p []string
	if _, err := o.Raw(selectSourceStage, pipeName, stageName, pipeName).QueryRows(&p); err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (CDPipelineRepository) GetCDPipelinesUsingApplicationAndBranch(codebase, branch string) ([]string, error) {
	o := orm.NewOrm()
	var p []string
	if _, err := o.Raw(selectCDPipelinesUsingCodebaseAndBranch, codebase, branch).QueryRows(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func (CDPipelineRepository) GetCDPipelinesUsingAutotestAndBranch(codebase, branch string) ([]string, error) {
	o := orm.NewOrm()
	var p []string
	if _, err := o.Raw(selectCDPipelineUsingAutotestAndBranch, codebase, branch).QueryRows(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func (CDPipelineRepository) GetCDPipelinesUsingLibraryAndBranch(codebase, branch string) ([]string, error) {
	o := orm.NewOrm()
	var p []string
	if _, err := o.Raw(selectCDPipelineUsingLibraryAndBranch, codebase, branch).QueryRows(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func (CDPipelineRepository) GetAllCodebaseDockerStreams() ([]string, error) {
	o := orm.NewOrm()
	var cds []string
	if _, err := o.Raw(selectCodebaseDockerStream).QueryRows(&cds); err != nil {
		return nil, err
	}
	return cds, nil
}

func (CDPipelineRepository) getStageCodebaseDockerStream(stageId int) ([]query.StageCodebaseDockerStream, error) {
	o := orm.NewOrm()
	var ds []query.StageCodebaseDockerStream
	if _, err := o.Raw(selectStageCodebaseDockerStream, stageId).QueryRows(&ds); err != nil {
		return nil, err
	}
	return ds, nil
}
