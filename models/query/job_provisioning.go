package query

type JobProvisioning struct {
	Id    int    `json:"id" orm:"column(id)"`
	Name  string `json:"name" orm:"column(name)"`
	Scope string `json:"scope" orm:"column(scope)"`
	Stage *Stage `orm:"reverse(one)"`
}

func (c *JobProvisioning) TableName() string {
	return "job_provisioning"
}

type JobProvisioningCriteria struct {
	Scope *string
}
