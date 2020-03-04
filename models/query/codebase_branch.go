package query

type CodebaseBranch struct {
	Id                   int                     `json:"id" orm:"column(id)"`
	Name                 string                  `json:"branchName" orm:"column(name)"`
	FromCommit           string                  `json:"from_commit" orm:"column(from_commit)"`
	Status               string                  `json:"status" orm:"column(status)"`
	Version              *string                 `json:"version" orm:"column(version)"`
	Build                *string                 `json:"build_number" orm:"column(build_number)"`
	LastSuccessBuild     *string                 `json:"last_success_build" orm:"column(last_success_build)"`
	VCSLink              string                  `json:"branchLink" orm:"-"`
	CICDLink             string                  `json:"jenkinsLink" orm:"-"`
	AppName              string                  `json:"appName" orm:"-"`
	Release              bool                    `json:"release" orm:"column(release)"`
	Codebase             *Codebase               `json:"-" orm:"rel(fk)"`
	CodebaseDockerStream []*CodebaseDockerStream `json:"codebaseDockerStream" orm:"reverse(many)"`
}

func (cb *CodebaseBranch) TableName() string {
	return "codebase_branch"
}

type CodebaseBranchCriteria struct {
	Status Status
}
