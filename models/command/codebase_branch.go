package command

type CreateCodebaseBranch struct {
	Name          string  `json:"name" valid:"Required;Match(/^[a-z0-9][a-z0-9-.]*[a-z0-9]$/)"`
	Commit        string  `json:"commit"`
	Username      string  `json:"username"`
	Version       *string `json:"startVersioningFrom,omitempty"`
	Build         *string `json:"build,omitempty"`
	Release       bool    `json:"release"`
	MasterVersion *string `json:"masterVersion,omitempty"`
}

type BranchCriteria struct {
	Status *string
}
