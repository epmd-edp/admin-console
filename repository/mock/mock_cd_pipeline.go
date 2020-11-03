package mock

import (
	"edp-admin-console/models"
	"edp-admin-console/models/dto"
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockCdPipeline struct {
	mock.Mock
}

func (m MockCdPipeline) GetCDPipelineByName(pipelineName string) (*query.CDPipeline, error) {
	args := m.Called(pipelineName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	p := args.Get(0).(query.CDPipeline)
	return &p, args.Error(1)
}
func (m MockCdPipeline) GetCDPipelines(criteria query.CDPipelineCriteria) ([]*query.CDPipeline, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetStage(cdPipelineName, stageName string) (*models.StageView, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetCodebaseAndBranchName(codebaseId, branchId int) (*dto.CodebaseBranchDTO, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetQualityGates(stageId int64) ([]query.QualityGate, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetCDPipelinesUsingApplication(codebaseName string) ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetCDPipelinesUsingAutotest(codebaseName string) ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetCDPipelinesUsingLibrary(codebaseName string) ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) SelectMaxOrderBetweenStages(pipeName string) (*int, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) SelectStageOrder(pipeName, stageName string) (*int, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) SelectCDPipelinesUsingInputStageAsSource(pipeName, stageName string) ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetCDPipelinesUsingApplicationAndBranch(codebase, branch string) ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetCDPipelinesUsingAutotestAndBranch(codebase, branch string) ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetCDPipelinesUsingLibraryAndBranch(codebase, branch string) ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) GetAllCodebaseDockerStreams() ([]string, error) {
	panic("implement me!!!")
}
func (m MockCdPipeline) SelectCountStages(pipeName string) (*int, error) {
	panic("implement me!!!")
}
