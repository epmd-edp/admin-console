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

func TestGetPerfServerNameMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	pbrm := new(mock.MockPerfBoard)
	pb := PerfBoard{
		PerfRepo: pbrm,
	}

	pbrm.On("GetPerfServerName", 1).Return(
		&query.PerfServer{
			Id:   1,
			Name: "fake-server",
		}, nil)

	ps, err := pb.GetPerfServerName(1)
	assert.NoError(t, err)
	assert.NotNil(t, ps)
}

func TestGetPerfServerNameMethod_ShouldReturnError(t *testing.T) {
	pbrm := new(mock.MockPerfBoard)
	pb := PerfBoard{
		PerfRepo: pbrm,
	}

	pbrm.On("GetPerfServerName", 1).Return(nil, errors.New("failed"))

	ps, err := pb.GetPerfServerName(1)
	assert.Error(t, err)
	assert.Nil(t, ps)
}

func TestGetCodebaseDataSourcesMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	pbrm := new(mock.MockPerfBoard)
	pb := PerfBoard{
		PerfRepo: pbrm,
	}

	pbrm.On("GetCodebaseDataSources", 1).Return(
		[]string{"ds1"}, nil)

	ds, err := pb.GetCodebaseDataSources(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ds))
}

func TestGetCodebaseDataSourcesMethod_ShouldReturnError(t *testing.T) {
	pbrm := new(mock.MockPerfBoard)
	pb := PerfBoard{
		PerfRepo: pbrm,
	}

	pbrm.On("GetCodebaseDataSources", 1).Return(nil, errors.New("failed"))

	ds, err := pb.GetCodebaseDataSources(1)
	assert.Error(t, err)
	assert.Nil(t, ds)
}
