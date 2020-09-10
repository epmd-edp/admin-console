package mock

import (
	"edp-admin-console/models/query"
	"github.com/stretchr/testify/mock"
)

type MockEDPComponent struct {
	mock.Mock
}

func (m MockEDPComponent) GetEDPComponent(componentType string) (*query.EDPComponent, error) {
	args := m.Called(componentType)
	ec := args.Get(0).(query.EDPComponent)
	return &ec, args.Error(1)
}

func (m MockEDPComponent) GetEDPComponents() ([]*query.EDPComponent, error) {
	panic("implement me")
}
