package perfboard

import "edp-admin-console/models/query"

type IPerfServer interface {
	GetPerfServers() ([]*query.PerfServer, error)
}

type PerfServer struct {
}

func (PerfServer) GetPerfServers() ([]*query.PerfServer, error) {
	return []*query.PerfServer{
		{
			Id:   1,
			Name: "epam-perf",
		},
	}, nil
}
