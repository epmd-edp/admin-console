package jira_server

import (
	"edp-admin-console/models/query"
	jiraserver "edp-admin-console/repository/jira-server"
	"edp-admin-console/service/logger"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

type JiraServer struct {
	IJiraServer jiraserver.IJiraServer
}

//GetJiraServers gets all Jira Servers from DB
func (s JiraServer) GetJiraServers() ([]*query.JiraServer, error) {
	log.Debug("start fetching Jira servers from DB")
	servers, err := s.IJiraServer.GetJiraServers()
	if err != nil {
		return nil, err
	}
	log.Info("Jira servers have been retrieved", zap.Any("servers", servers))
	return servers, nil
}
