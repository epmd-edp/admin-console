package controllers

import (
	validation2 "edp-admin-console/controllers/validation"
	"edp-admin-console/service"
	"edp-admin-console/util"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"net/http"
)

type RepositoryRestController struct {
	beego.Controller
	AppService service.CodebaseService
}

type RepoData struct {
	Url      string `json:"url,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (this *RepositoryRestController) IsGitRepoAvailable() {
	var repo RepoData
	err := json.NewDecoder(this.Ctx.Request.Body).Decode(&repo)
	if err != nil {
		http.Error(this.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	errMsg := validRepoRequestData(repo)
	if errMsg != nil {
		http.Error(this.Ctx.ResponseWriter, errMsg.Message, errMsg.StatusCode)
		return
	}

	this.Data["json"] = util.IsGitRepoAvailable(repo.Url, repo.Login, repo.Password)
	this.ServeJSON()
}

func validRepoRequestData(repo RepoData) *validation2.ErrMsg {
	valid := validation.Validation{}

	_, err := valid.Valid(repo)
	if err != nil {
		return &validation2.ErrMsg{"An error has occurred while validating application's form fields.", http.StatusInternalServerError}
	}

	if valid.Errors == nil {
		return nil
	}

	return &validation2.ErrMsg{string(validation2.CreateErrorResponseBody(valid)), http.StatusBadRequest}
}
