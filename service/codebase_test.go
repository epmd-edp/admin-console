package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/mock"
	"edp-admin-console/service/perfboard"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCodebaseByNameMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	mCodebase := new(mock.MockCodebase)
	cs := CodebaseService{
		ICodebaseRepository: mCodebase,
	}

	mCodebase.On("GetCodebaseByName", "stub-name").Return(
		query.Codebase{}, nil)

	c, err := cs.GetCodebaseByName("stub-name")
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestGetCodebaseByNameMethod_ShouldBeExecutedWithError(t *testing.T) {
	mCodebase := new(mock.MockCodebase)
	cs := CodebaseService{
		ICodebaseRepository: mCodebase,
	}

	mCodebase.On("GetCodebaseByName", "stub-name").Return(
		nil, errors.New("stub-msg"))

	c, err := cs.GetCodebaseByName("stub-name")
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestGetCodebaseByNameMethod_WithPerfServerShouldBeExecutedSuccessfully(t *testing.T) {
	mCodebase := new(mock.MockCodebase)
	pbrm := new(mock.MockPerfBoard)
	cs := CodebaseService{
		ICodebaseRepository: mCodebase,
		PerfService: perfboard.PerfBoard{
			PerfRepo: pbrm,
		},
	}

	mCodebase.On("GetCodebaseByName", "stub-name").Return(
		query.Codebase{
			Id:           1,
			PerfServerId: getIntP(1),
			Perf:         nil,
		}, nil)

	pbrm.On("GetPerfServerName", 1).Return(
		&query.PerfServer{
			Id:   1,
			Name: "fake-server",
		}, nil)

	pbrm.On("GetCodebaseDataSources", 1).Return(
		[]string{"ds1"}, nil)

	c, err := cs.GetCodebaseByName("stub-name")
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestGetCodebaseByNameMethod_WithPerfServerPerfServerShouldNotBeFound(t *testing.T) {
	mCodebase := new(mock.MockCodebase)
	pbrm := new(mock.MockPerfBoard)
	cs := CodebaseService{
		ICodebaseRepository: mCodebase,
		PerfService: perfboard.PerfBoard{
			PerfRepo: pbrm,
		},
	}

	mCodebase.On("GetCodebaseByName", "stub-name").Return(
		query.Codebase{
			Id:           1,
			PerfServerId: getIntP(1),
			Perf:         nil,
		}, nil)

	pbrm.On("GetPerfServerName", 1).Return(nil, errors.New("failed"))

	c, err := cs.GetCodebaseByName("stub-name")
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestGetCodebaseByNameMethod_WithPerfServerDataSourceShouldNotBeFound(t *testing.T) {
	mCodebase := new(mock.MockCodebase)
	pbrm := new(mock.MockPerfBoard)
	cs := CodebaseService{
		ICodebaseRepository: mCodebase,
		PerfService: perfboard.PerfBoard{
			PerfRepo: pbrm,
		},
	}

	mCodebase.On("GetCodebaseByName", "stub-name").Return(
		query.Codebase{
			Id:           1,
			PerfServerId: getIntP(1),
			Perf:         nil,
		}, nil)

	pbrm.On("GetPerfServerName", 1).Return(
		&query.PerfServer{
			Id:   1,
			Name: "fake-server",
		}, nil)

	pbrm.On("GetCodebaseDataSources", 1).Return(nil, errors.New("failed"))

	c, err := cs.GetCodebaseByName("stub-name")
	assert.Error(t, err)
	assert.Nil(t, c)
}

func getIntP(val int) *int {
	return &val
}
