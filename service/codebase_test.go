package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/mock"
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
