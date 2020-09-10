package mock

import (
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockCodebase struct {
	mock.Mock
}

func (m MockCodebase) GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error) {
	panic("implement me")
}

func (m MockCodebase) FindCodebaseByName(name string) bool {
	panic("implement me")
}

func (m MockCodebase) FindCodebaseByProjectPath(gitProjectPath *string) bool {
	panic("implement me")
}

func (m MockCodebase) GetCodebaseByName(name string) (*query.Codebase, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	c := args.Get(0).(query.Codebase)
	return &c, args.Error(1)
}

func (m MockCodebase) GetCodebaseById(id int) (*query.Codebase, error) {
	panic("implement me")
}

func (m MockCodebase) ExistActiveBranch(dockerStreamName string) (bool, error) {
	panic("implement me")
}

func (m MockCodebase) ExistCodebaseAndBranch(cbName, brName string) bool {
	panic("implement me")
}

func (m MockCodebase) SelectApplicationToPromote(cdPipelineId int) ([]*query.ApplicationsToPromote, error) {
	panic("implement me")
}
