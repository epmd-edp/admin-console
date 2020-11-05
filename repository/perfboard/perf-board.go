package perfboard

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IPerfServer interface {
	GetPerfServers() ([]*query.PerfServer, error)
	GetPerfServerName(id int) (*query.PerfServer, error)
	GetCodebaseDataSources(codebaseId int) ([]string, error)
}

type PerfServer struct {
}

const selectCodebaseDataSources = `select type
									from perf_data_sources pds
											 left join codebase_perf_data_sources cpds on pds.id = cpds.data_source_id
									where cpds.codebase_id = ?;`

func (PerfServer) GetPerfServers() ([]*query.PerfServer, error) {
	o := orm.NewOrm()
	var servers []*query.PerfServer
	_, err := o.QueryTable(new(query.PerfServer)).
		OrderBy("name").
		All(&servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (PerfServer) GetPerfServerName(id int) (*query.PerfServer, error) {
	o := orm.NewOrm()
	ps := &query.PerfServer{}
	err := o.QueryTable(new(query.PerfServer)).
		Filter("id", id).
		One(ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (PerfServer) GetCodebaseDataSources(codebaseId int) ([]string, error) {
	o := orm.NewOrm()
	var ds []string
	_, err := o.Raw(selectCodebaseDataSources, codebaseId).QueryRows(&ds)
	if err != nil {
		return nil, err
	}
	return ds, nil
}
