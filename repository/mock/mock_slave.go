package mock

import (
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockSlave struct {
	mock.Mock
}

func (m MockSlave) GetAllSlaves() ([]*query.JenkinsSlave, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*query.JenkinsSlave), args.Error(1)
}
