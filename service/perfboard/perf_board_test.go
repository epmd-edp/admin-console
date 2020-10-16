package perfboard

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPerfServersMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	pbrm := new(mock.MockPerfBoard)
	pb := PerfBoard{
		PerfRepo: pbrm,
	}

	pbrm.On("GetPerfServers").Return(
		[]*query.PerfServer{
			{
				Id:   1,
				Name: "fake-server",
			},
		}, nil)

	ps, err := pb.GetPerfServers()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ps))
}

func TestGetPerfServersMethod_ShouldReturnError(t *testing.T) {
	pbrm := new(mock.MockPerfBoard)
	pb := PerfBoard{
		PerfRepo: pbrm,
	}

	pbrm.On("GetPerfServers").Return(nil, errors.New("failed"))

	_, err := pb.GetPerfServers()
	assert.Error(t, err)
}
