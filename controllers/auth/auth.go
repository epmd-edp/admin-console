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

package auth

import (
	"context"
	ctx "edp-admin-console/context"
	"fmt"
	"edp-admin-console/service/logger"
	"github.com/astaxie/beego"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

type AuthController struct {
	beego.Controller
}

func (this *AuthController) Callback() {
	authConfig := ctx.GetAuthConfig()
	log.Info("Start callback flow...")
	queryState := this.Ctx.Input.Query("state")
	log.Info("State has been retrieved from query param", zap.String("queryState", queryState))
	sessionState := this.Ctx.Input.Session(authConfig.StateAuthKey)
	log.Info("State has been retrieved from the session", zap.Any("sessionState", sessionState))
	if queryState != sessionState {
		log.Info("State does not match")
		this.Abort("400")
		return
	}

	authCode := this.Ctx.Input.Query("code")
	log.Info("Authorization code has been retrieved from query param")
	token, err := authConfig.Oauth2Config.Exchange(context.Background(), authCode)

	if err != nil {
		log.Info("Failed to exchange token with code", zap.String("code", authCode))
		this.Abort("500")
		return
	}
	log.Info("Authorization code has been successfully exchanged with token")

	ts := authConfig.Oauth2Config.TokenSource(context.Background(), token)

	this.Ctx.Output.Session("token_source", ts)
	log.Info("Token source has been saved to the session")
	path := this.getRedirectPath()
	this.Redirect(path, 302)
}

func (this *AuthController) getRedirectPath() string {
	requestPath := this.Ctx.Input.Session("request_path")
	if requestPath == nil {
		return fmt.Sprintf("%s/admin/edp/overview", ctx.BasePath)
	}
	this.Ctx.Output.Session("request_path", nil)
	return requestPath.(string)
}
