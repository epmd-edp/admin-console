package mock

import (
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockJiraServer struct {
	mock.Mock
}

func (m MockJiraServer) GetJiraServers() ([]*query.JiraServer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*query.JiraServer), args.Error(1)
}
