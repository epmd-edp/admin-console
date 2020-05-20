/*
 * Copyright 2020 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repository

import (
	"edp-admin-console/models/query"

	"github.com/astaxie/beego/orm"
)

type ICodebaseRepository interface {
	GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error)
	GetCodebaseByName(name string) (*query.Codebase, error)
	GetCodebaseById(id int) (*query.Codebase, error)
	ExistActiveBranch(dockerStreamName string) (bool, error)
	ExistCodebaseAndBranch(cbName, brName string) bool
	SelectApplicationToPromote(cdPipelineId int) ([]*query.ApplicationsToPromote, error)
	FindCodebaseByName(name string) bool
	FindCodebaseByProjectPath(gitProjectPath *string) bool
}

type CodebaseRepository struct {
	ICodebaseRepository
}

func (CodebaseRepository) GetCodebasesByCriteria(criteria query.CodebaseCriteria) ([]*query.Codebase, error) {
	o := orm.NewOrm()
	var codebases []*query.Codebase

	qs := o.QueryTable(new(query.Codebase))

	if criteria.Type != "" {
		qs = qs.Filter("type", criteria.Type)
	}
	if criteria.Status != "" {
		qs = qs.Filter("status", criteria.Status)
	}

	if criteria.Language != "" {
		qs = qs.Filter("language", criteria.Language)
	}

	_, err := qs.OrderBy("name").
		All(&codebases)

	for _, c := range codebases {

		err = loadRelatedActionLog(c)
		if err != nil {
			return nil, err
		}

		err = loadRelatedCodebaseBranch(c, criteria.BranchStatus)
		if err != nil {
			return nil, err
		}

		err = loadRelatedCodebaseDockerStream(c.CodebaseBranch)
		if err != nil {
			return nil, err
		}

		for _, branch := range c.CodebaseBranch {
			err := loadRelatedBranches(branch)
			if err != nil {
				return nil, err
			}
		}

		if c.GitServerId != nil {
			err = loadRelatedGitServerName(c)
			if err != nil {
				return nil, err
			}
		}

		if c.JiraServerId != nil {
			if err := loadRelatedJiraServerName(c); err != nil {
				return nil, err
			}
		}

		if c.JenkinsSlaveId != nil {
			err := loadRelatedJenkinsSlaveName(c)
			if err != nil {
				return nil, err
			}
		}

	}
	return codebases, err
}

func (CodebaseRepository) FindCodebaseByName(name string) bool {
	return orm.NewOrm().QueryTable(new(query.Codebase)).Filter("name", name).Exist()
}

func (CodebaseRepository) FindCodebaseByProjectPath(gitProjectPath *string) bool {
	return orm.NewOrm().QueryTable(new(query.Codebase)).Filter("git_project_path", *gitProjectPath).Exist()
}

func (CodebaseRepository) GetCodebaseByName(name string) (*query.Codebase, error) {
	o := orm.NewOrm()
	codebase := query.Codebase{Name: name}

	err := o.Read(&codebase, "Name")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	err = loadRelatedActionLog(&codebase)

	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&codebase, "CodebaseBranch", false, 100, 0, "Name")
	if err != nil {
		return nil, err
	}

	if codebase.GitServerId != nil {
		err = loadRelatedGitServerName(&codebase)
		if err != nil {
			return nil, err
		}
	}

	if codebase.JiraServerId != nil {
		if err := loadRelatedJiraServerName(&codebase); err != nil {
			return nil, err
		}
	}

	if codebase.JenkinsSlaveId != nil {
		err = loadRelatedJenkinsSlaveName(&codebase)
		if err != nil {
			return nil, err
		}
	}

	if codebase.JobProvisioningId != nil {
		err = loadRelatedJobProvisioner(&codebase)
		if err != nil {
			return nil, err
		}
	}

	return &codebase, nil
}

func (CodebaseRepository) GetCodebaseById(id int) (*query.Codebase, error) {
	o := orm.NewOrm()
	codebase := query.Codebase{Id: id}

	err := o.Read(&codebase, "Id")

	if err == orm.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	err = loadRelatedActionLog(&codebase)
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(&codebase, "CodebaseBranch", false, 100, 0, "Name")
	if err != nil {
		return nil, err
	}

	if codebase.GitServerId != nil {
		err = loadRelatedGitServerName(&codebase)
		if err != nil {
			return nil, err
		}
	}

	if codebase.JiraServerId != nil {
		if err := loadRelatedJiraServerName(&codebase); err != nil {
			return nil, err
		}
	}

	if codebase.JenkinsSlaveId != nil {
		err = loadRelatedJenkinsSlaveName(&codebase)
		if err != nil {
			return nil, err
		}
	}

	return &codebase, nil
}

func loadRelatedActionLog(codebase *query.Codebase) error {
	o := orm.NewOrm()

	_, err := o.QueryTable(new(query.ActionLog)).
		Filter("codebase__codebase_id", codebase.Id).
		OrderBy("LastTimeUpdate").
		Distinct().
		All(&codebase.ActionLog, "LastTimeUpdate", "UserName",
			"Message", "Action", "Result")

	return err
}

func loadRelatedCodebaseBranch(codebase *query.Codebase, status query.Status) error {
	o := orm.NewOrm()

	qs := o.QueryTable(new(query.CodebaseBranch))

	if status != "" {
		qs = qs.Filter("status", status)
	}

	_, err := qs.Filter("codebase_id", codebase.Id).
		OrderBy("Name").
		All(&codebase.CodebaseBranch, "Id", "Name", "FromCommit", "Status")

	return err
}

func loadRelatedCodebaseDockerStream(branches []*query.CodebaseBranch) error {
	o := orm.NewOrm()

	for _, branch := range branches {
		qs := o.QueryTable(new(query.CodebaseDockerStream))

		_, err := qs.Filter("codebase_branch_id", branch.Id).
			All(&branch.CodebaseDockerStream, "Id", "OcImageStreamName")
		if err != nil {
			return err
		}
	}
	return nil
}

func loadRelatedBranches(branch *query.CodebaseBranch) error {
	o := orm.NewOrm()

	for _, dockerStream := range branch.CodebaseDockerStream {
		_, err := o.LoadRelated(dockerStream, "CodebaseBranch", false, 100, 0, "Name")
		if err != nil {
			return err
		}
	}
	return nil
}

func loadRelatedGitServerName(codebase *query.Codebase) error {
	o := orm.NewOrm()

	server := query.GitServer{}
	err := o.QueryTable(new(query.GitServer)).
		Filter("id", codebase.GitServerId).
		One(&server)
	if err != nil {
		return err
	}

	codebase.GitServer = &server.Name

	return nil
}

func loadRelatedJiraServerName(codebase *query.Codebase) error {
	o := orm.NewOrm()
	server := query.JiraServer{}
	err := o.QueryTable(new(query.JiraServer)).
		Filter("id", codebase.JiraServerId).
		One(&server)
	if err != nil {
		return err
	}
	codebase.JiraServer = &server.Name
	return nil
}

func loadRelatedJenkinsSlaveName(c *query.Codebase) error {
	o := orm.NewOrm()

	s := query.JenkinsSlave{}
	err := o.QueryTable(new(query.JenkinsSlave)).
		Filter("id", c.JenkinsSlaveId).
		One(&s)
	if err != nil {
		return err
	}

	c.JenkinsSlave = s.Name

	return nil
}

func loadRelatedJobProvisioner(c *query.Codebase) error {
	o := orm.NewOrm()

	s := query.JobProvisioning{}
	err := o.QueryTable(new(query.JobProvisioning)).
		Filter("id", c.JobProvisioningId).
		One(&s)
	if err != nil {
		return err
	}

	c.JobProvisioning = s.Name

	return nil
}

func (CodebaseRepository) ExistActiveBranch(dockerStreamName string) (bool, error) {
	o := orm.NewOrm()

	var dockerStream query.CodebaseDockerStream

	err := o.QueryTable(new(query.CodebaseDockerStream)).Filter("ocImageStreamName", dockerStreamName).One(&dockerStream)
	if err != nil {
		return false, err
	}

	_, err = o.LoadRelated(&dockerStream, "CodebaseBranch", false, 100, 0, "Name")
	if err != nil {
		return false, err
	}

	return dockerStream.CodebaseBranch.Status == "active", nil
}

func (CodebaseRepository) ExistCodebaseAndBranch(cbName, brName string) bool {
	return orm.NewOrm().QueryTable(new(query.Codebase)).
		Filter("name", cbName).
		Filter("CodebaseBranch__name", brName).
		Exist()
}

func (CodebaseRepository) SelectApplicationToPromote(cdPipelineId int) ([]*query.ApplicationsToPromote, error) {
	o := orm.NewOrm()
	var applicationsToPromote []*query.ApplicationsToPromote

	_, err := o.QueryTable(new(query.ApplicationsToPromote)).
		Filter("cd_pipeline_id", cdPipelineId).All(&applicationsToPromote)

	return applicationsToPromote, err
}
