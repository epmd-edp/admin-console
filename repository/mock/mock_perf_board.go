package mock

import (
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockPerfBoard struct {
	mock.Mock
}

func (m MockPerfBoard) GetPerfServers() ([]*query.PerfServer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*query.PerfServer), args.Error(1)
}
