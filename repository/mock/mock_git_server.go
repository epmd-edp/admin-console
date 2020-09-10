package mock

import (
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockGitServer struct {
	mock.Mock
}

func (m MockGitServer) GetGitServersByCriteria(criteria query.GitServerCriteria) ([]*query.GitServer, error) {
	panic("implement me")
}

func (m MockGitServer) GetGitServerByName(name string) (*query.GitServer, error) {
	args := m.Called(name)
	gs := args.Get(0).(query.GitServer)
	return &gs, args.Error(1)
}
