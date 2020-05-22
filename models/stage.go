package models

import "edp-admin-console/models/query"

type StageView struct {
	Id              int64               `json:"-"`
	Name            string              `json:"name"`
	CDPipeline      string              `json:"cdPipeline"`
	Description     string              `json:"description"`
	TriggerType     string              `json:"triggerType"`
	Order           string              `json:"order"`
	Applications    []ApplicationStage  `json:"applications"`
	QualityGates    []query.QualityGate `json:"qualityGates"`
	JobProvisioning string              `json:"jobProvisioning"`
}

type ApplicationStage struct {
	Name       string `json:"name"`
	BranchName string `json:"branchName"`
	InputIs    string `json:"inputIs"`
	OutputIs   string `json:"outputIs"`
}

type CDPipelineStageView struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	TriggerType          string `json:"triggerType"`
	QualityGate          string `json:"qualityGate"`
	JenkinsStepName      string `json:"jenkinsStepName"`
	Order                int    `json:"order"`
	OpenshiftProjectLink string `json:"openshiftProjectLink"`
}
