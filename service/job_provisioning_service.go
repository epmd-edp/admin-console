package service

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type JobProvisioning struct {
	IJobProvisioningRepository repository.IJobProvisioningRepository
}

//GetAllJobsProvisioning gets all job provisioning entries from DB
func (s JobProvisioning) GetAllJobProvisioners(criteria query.JobProvisioningCriteria) ([]*query.JobProvisioning, error) {
	log.Debug("Start fetching all available job provisioning entries...")
	p, err := s.IJobProvisioningRepository.GetAllJobProvisioners(criteria)
	if err != nil {
		return nil, errors.Wrap(err, "an error has occurred while fetching job provisioning entities from DB")
	}
	log.Info("Fetched Job Provisioning", zap.Any("job provision", p))
	return p, nil
}
