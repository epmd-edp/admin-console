package perfboard

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/perfboard"
	"edp-admin-console/service/logger"
	"github.com/pkg/errors"
)

var log = logger.GetLogger()

type PerfBoard struct {
	PerfRepo perfboard.IPerfServer
}

func (s PerfBoard) GetPerfServers() ([]*query.PerfServer, error) {
	log.Info("start retrieving PERF servers from DB")
	ps, err := s.PerfRepo.GetPerfServers()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get PERF server from DB")
	}
	return ps, nil
}
