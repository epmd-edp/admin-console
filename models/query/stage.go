package query

type Stage struct {
	Id                        int                         `json:"id" orm:"column(id)"`
	Name                      string                      `json:"name" orm:"column(name)"`
	Description               string                      `json:"description" orm:"column(description)"`
	TriggerType               string                      `json:"triggerType" orm:"column(trigger_type)"`
	Order                     int                         `json:"order" orm:"column(order)"`
	PlatformProjectLink       string                      `json:"platformProjectLink" orm:"-"`
	PlatformProjectName       string                      `json:"platformProjectName" orm:"-"`
	CDPipeline                *CDPipeline                 `json:"-" orm:"rel(fk);column(cd_pipeline_id)"`
	QualityGates              []QualityGate               `json:"qualityGates" orm:"-"`
	SourceCodebaseBranchId    *int                        `json:"-" orm:"codebase_branch_id;column(codebase_branch_id)"`
	Source                    Source                      `json:"source" orm:"-"`
	JobProvisioning           *JobProvisioning            `orm:"null;rel(one)"`
	StageCodebaseDockerStream []StageCodebaseDockerStream `json:"stageCodebaseDockerStream" orm:"-"`
}

type Source struct {
	Type    string         `json:"type"`
	Library *SourceLibrary `json:"library"`
}

type SourceLibrary struct {
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

type QualityGate struct {
	Id               int             `json:"id" orm:"column(id)"`
	QualityGateType  string          `json:"qualityGateType" orm:"column(quality_gate)"`
	StepName         string          `json:"stepName" orm:"column(step_name)"`
	CdStageId        *int            `json:"cdStageId" orm:"column(cd_stage_id)"`
	CodebaseId       *int            `json:"-" orm:"column(codebase_id)"`
	CodebaseBranchId *int            `json:"branchId" orm:"column(codebase_branch_id)"`
	Autotest         *Codebase       `json:"autotest" orm:"-"`
	Branch           *CodebaseBranch `json:"codebaseBranch" orm:"-"`
}

type StageCodebaseDockerStream struct {
	CdStageId                    int    `json:"cdStageId" orm:"column(cd_stage_id)"`
	InputCodebaseDockerStreamId  string `json:"inputCodebaseDockerStreamId" orm:"column(input)"`
	OutputCodebaseDockerStreamId string `json:"outputCodebaseDockerStreamId" orm:"column(output)"`
}

func (cb *Stage) TableName() string {
	return "cd_stage"
}

func (qg *QualityGate) TableName() string {
	return "quality_gate_stage"
}

func (cb *StageCodebaseDockerStream) TableName() string {
	return "stage_codebase_docker_stream"
}
