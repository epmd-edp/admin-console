package jira_server

import (
	"edp-admin-console/models/query"
	"github.com/astaxie/beego/orm"
)

type IJiraServer interface {
	GetJiraServers() ([]*query.JiraServer, error)
}

type JiraServer struct {
	IJiraServer
}

func (JiraServer) GetJiraServers() ([]*query.JiraServer, error) {
	o := orm.NewOrm()
	var servers []*query.JiraServer
	_, err := o.QueryTable(new(query.JiraServer)).
		OrderBy("name").
		All(&servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}
