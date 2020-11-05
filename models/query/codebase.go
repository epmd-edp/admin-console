package query

type Codebase struct {
	Id                   int               `json:"id" orm:"column(id)"`
	Name                 string            `json:"name" orm:"column(name)"`
	Language             string            `json:"language" orm:"column(language)"`
	BuildTool            string            `json:"build_tool" orm:"column(build_tool)"`
	Framework            string            `json:"framework" orm:"column(framework)"`
	Strategy             string            `json:"strategy" orm:"column(strategy)"`
	GitUrl               string            `json:"git_url" orm:"column(repository_url)"`
	RouteSite            string            `json:"route_site" orm:"column(route_site)"`
	RoutePath            string            `json:"route_path" orm:"column(route_path)"`
	DbKind               string            `json:"db_kind" orm:"column(database_kind)"`
	DbVersion            string            `json:"db_version" orm:"column(database_version)"`
	DbCapacity           string            `json:"db_capacity" orm:"column(database_capacity)"`
	DbStorage            string            `json:"db_storage" orm:"column(database_storage)"`
	Type                 CodebaseType      `json:"type" orm:"column(type)"`
	Status               Status            `json:"status" orm:"column(status)"`
	TestReportFramework  string            `json:"testReportFramework" orm:"column(test_report_framework)"`
	Description          string            `json:"description" orm:"column(description)"`
	CodebaseBranch       []*CodebaseBranch `json:"codebase_branch" orm:"reverse(many)"`
	ActionLog            []*ActionLog      `json:"-" orm:"rel(m2m);rel_table(codebase_action_log)"`
	GitServerId          *int              `json:"-" orm:"column(git_server_id)"`
	GitServer            *string           `json:"gitServer" orm:"-"`
	GitProjectPath       *string           `json:"gitProjectPath" orm:"column(git_project_path)"`
	JenkinsSlaveId       *int              `json:"-" orm:"column(jenkins_slave_id)"`
	JenkinsSlave         string            `json:"jenkinsSlave" orm:"-"`
	JobProvisioningId    *int              `json:"-" orm:"column(job_provisioning_id)"`
	JobProvisioning      string            `json:"jobProvisioning" orm:"-"`
	DeploymentScript     string            `json:"deploymentScript" orm:"deployment_script"`
	VersioningType       string            `json:"versioningType" orm:"versioning_type"`
	StartVersioningFrom  *string           `json:"startFrom" orm:"start_versioning_from"`
	JiraServerId         *int              `json:"-" orm:"column(jira_server_id)"`
	JiraServer           *string           `json:"jiraServer" orm:"-"`
	CommitMessagePattern string            `json:"commitMessagePattern" orm:"commit_message_pattern"`
	TicketNamePattern    string            `json:"ticketNamePattern" orm:"ticket_name_pattern"`
	CiTool               string            `json:"ciTool" orm:"ci_tool"`
	PerfServerId         *int              `json:"-" orm:"column(perf_server_id)"`
	Perf                 *Perf             `json:"perf" orm:"-"`
}

type Perf struct {
	Name        string   `json:"name"`
	DataSources []string `json:"dataSources"`
}

func (c *Codebase) TableName() string {
	return "codebase"
}

type CodebaseCriteria struct {
	BranchStatus Status
	Status       Status
	Type         CodebaseType
	Language     CodebaseLanguage
}

type CodebaseType string
type CodebaseLanguage string

const (
	App       CodebaseType = "application"
	Autotests CodebaseType = "autotests"
	Library   CodebaseType = "library"
)

var CodebaseTypes = map[string]CodebaseType{
	"application": App,
	"autotests":   Autotests,
	"library":     Library,
}

type ApplicationsToPromote struct {
	Id           int `orm:"column(id)"`
	CdPipelineId int `orm:"column(cd_pipeline_id)"`
	CodebaseId   int `orm:"column(codebase_id)"`
}

func (c *ApplicationsToPromote) TableName() string {
	return "applications_to_promote"
}
