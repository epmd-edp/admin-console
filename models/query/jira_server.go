package query

type JiraServer struct {
	Id        int    `json:"id" orm:"column(id)"`
	Name      string `json:"name" orm:"column(name)"`
	Available bool   `json:"available" orm:"column(available)"`
}

func (c *JiraServer) TableName() string {
	return "jira_server"
}
