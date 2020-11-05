package perfboard

import (
	"edp-admin-console/models/query"
	"edp-admin-console/repository/perfboard"
	"edp-admin-console/service/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

func (s PerfBoard) GetPerfServerName(id int) (*query.PerfServer, error) {
	log.Info("start retrieving PERF server from DB", zap.Int("id", id))
	ps, err := s.PerfRepo.GetPerfServerName(id)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get PERF server from DB")
	}
	return ps, nil
}

func (s PerfBoard) GetCodebaseDataSources(codebaseId int) ([]string, error) {
	log.Info("start retrieving PERF data sources", zap.Int("codebase id", codebaseId))
	ds, err := s.PerfRepo.GetCodebaseDataSources(codebaseId)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get PERF data sources from DB")
	}
	return ds, nil
}
