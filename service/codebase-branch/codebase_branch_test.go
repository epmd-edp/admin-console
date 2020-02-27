package codebase_branch

import (
	"edp-admin-console/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type ReleaseBranchRepositoryMock struct {
	mock.Mock
}

func (mock *ReleaseBranchRepositoryMock) GetAllReleaseBranchesByAppName(appName string) ([]models.ReleaseBranchView, error) {
	args := mock.Called(appName)
	return args.Get(0).([]models.ReleaseBranchView), args.Error(1)
}

func (mock *ReleaseBranchRepositoryMock) GetAllReleaseBranches(branchFilterCriteria models.BranchCriteria) ([]models.ReleaseBranchView, error) {
	args := mock.Called(branchFilterCriteria)
	return args.Get(0).([]models.ReleaseBranchView), args.Error(1)
}

func (mock *ReleaseBranchRepositoryMock) GetReleaseBranch(appName, branchName string) (*models.ReleaseBranchView, error) {
	args := mock.Called(appName, branchName)
	return args.Get(0).(*models.ReleaseBranchView), args.Error(1)
}

func TestShouldReturnBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetBranchByCodebaseAndName", mock.Anything, mock.Anything, mock.Anything).Return(&models.ReleaseBranchView{}, nil)
	branchService := CodebaseBranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetReleaseBranch(mock.Anything, mock.Anything)
	assert.Nil(t, err)
	assert.NotNil(t, branchEntity)
}

func TestShouldReturnErrorBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetBranchByCodebaseAndName", mock.Anything, mock.Anything, mock.Anything).Return(&models.ReleaseBranchView{}, errors.New("internal error"))
	branchService := CodebaseBranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetReleaseBranch(mock.Anything, mock.Anything)
	assert.NotNil(t, err)
	assert.Nil(t, branchEntity)
}

func TestShouldReturnAllBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetCodebaseBranchesByCriteria", mock.Anything, mock.Anything).Return(createBranchEntities(), nil)
	branchService := CodebaseBranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetAllReleaseBranches(models.BranchCriteria{})
	assert.Nil(t, err)
	assert.NotNil(t, branchEntity)
}

func TestShouldReturnErrorAllBranch(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetCodebaseBranchesByCriteria", mock.Anything, mock.Anything).Return(createBranchEntities(), errors.New("internal error"))
	branchService := CodebaseBranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetAllReleaseBranches(models.BranchCriteria{})
	assert.NotNil(t, err)
	assert.Nil(t, branchEntity)
}

func TestShouldReturnAllReleaseBranchesByAppName(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetAllReleaseBranchesByAppName", mock.Anything).Return([]models.ReleaseBranchView{}, nil)
	branchService := CodebaseBranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetAllReleaseBranchesByAppName(mock.Anything)
	assert.Nil(t, err)
	assert.NotNil(t, branchEntity)
}

func TestShouldReturnErrorAllReleaseBranchesByAppName(t *testing.T) {
	repositoryMock := new(ReleaseBranchRepositoryMock)
	repositoryMock.On("GetAllReleaseBranchesByAppName", mock.Anything).Return([]models.ReleaseBranchView{}, errors.New("internal error"))
	branchService := CodebaseBranchService{IReleaseBranchRepository: repositoryMock}

	branchEntity, err := branchService.GetAllReleaseBranchesByAppName(mock.Anything)
	assert.NotNil(t, err)
	assert.Nil(t, branchEntity)
}

func createBranchEntities() []models.ReleaseBranchView {
	return []models.ReleaseBranchView{
		{
			Name:            "fake-name",
			Event:           "fake-event",
			Username:        "fake-username",
			DetailedMessage: "fake-message",
			UpdatedAt:       time.Now(),
		},
	}
}
