package repository

import (
	"edp-admin-console/models/query"

	"github.com/astaxie/beego/orm"
)

type IJobProvisioningRepository interface {
	GetAllJobProvisioners(criteria query.JobProvisioningCriteria) ([]*query.JobProvisioning, error)
}

type JobProvisioning struct {
}

func (JobProvisioning) GetAllJobProvisioners(criteria query.JobProvisioningCriteria) ([]*query.JobProvisioning, error) {
	o := orm.NewOrm()
	var jobsProvisioning []*query.JobProvisioning

	qs := o.QueryTable(new(query.JobProvisioning))

	if criteria.Scope != nil && *criteria.Scope != "" {
		qs = qs.Filter("scope", criteria.Scope)
	}

	_, err := qs.OrderBy("name").All(&jobsProvisioning)
	if err != nil {
		return nil, err
	}

	return jobsProvisioning, nil
}
