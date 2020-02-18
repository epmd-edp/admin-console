package command

type CreateCodebaseBranch struct {
	Name     string  `json:"name" valid:"Required;Match(/^[a-z0-9][a-z0-9-.]*[a-z0-9]$/)"`
	Commit   string  `json:"commit"`
	Username string  `json:"username"`
	Version  *string `json:"startVersioningFrom"`
	Build    *string `json:"build"`
}

type BranchCriteria struct {
	Status *string
}
