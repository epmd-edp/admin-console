package mock

import (
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockJobProvision struct {
	mock.Mock
}

func (m MockJobProvision) GetAllJobProvisioners(criteria query.JobProvisioningCriteria) ([]*query.JobProvisioning, error) {
	args := m.Called(criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*query.JobProvisioning), args.Error(1)
}
