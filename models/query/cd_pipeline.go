package query

type CDPipeline struct {
	Id                    int                                                     `json:"id" orm:"column(id)"`
	Name                  string                                                  `json:"name" orm:"column(name)"`
	Status                string                                                  `json:"status" orm:"column(status)"`
	JenkinsLink           string                                                  `json:"jenkinsLink" orm:"-"`
	CodebaseBranch        []*CodebaseBranch                                       `json:"codebaseBranches" orm:"rel(m2m);rel_table(cd_pipeline_codebase_branch)"`
	Stage                 []*Stage                                                `json:"cd_stage" orm:"reverse(many)"`
	ThirdPartyService     []*ThirdPartyService                                    `json:"services" orm:"rel(m2m);rel_table(cd_pipeline_third_party_service)"`
	CodebaseStageMatrix   map[CDCodebaseStageMatrixKey]CDCodebaseStageMatrixValue `json:"-" orm:"-"`
	ApplicationsToPromote []string                                                `json:"applicationsToPromote" orm:"-"`
	CodebaseDockerStream  []*CodebaseDockerStream                                 `json:"image_stream" orm:"rel(m2m);rel_table(cd_pipeline_docker_stream)"`
	ActionLog             []*ActionLog                                            `json:"-" orm:"rel(m2m);rel_table(cd_pipeline_action_log)"`
}

type CDCodebaseStageMatrixKey struct {
	CodebaseBranch *CodebaseBranch `json:"codebaseBranch"`
	Stage          *Stage          `json:"stage"`
}

type CDCodebaseStageMatrixValue struct {
	DockerVersion string `json:"dockerVersion"`
}

func (c *CDPipeline) GetCDCodebaseStageMatrixValue(codebase *CodebaseBranch, stage *Stage) CDCodebaseStageMatrixValue {
	return c.CodebaseStageMatrix[CDCodebaseStageMatrixKey{
		CodebaseBranch: codebase,
		Stage:          stage,
	}]
}

type CDPipelineCriteria struct {
	Status Status
}

func (cb *CDPipeline) TableName() string {
	return "cd_pipeline"
}
