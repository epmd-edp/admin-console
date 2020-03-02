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

package filters

import (
	ctx "context"
	appCtx "edp-admin-console/context"
	"encoding/json"
	bgCtx "github.com/astaxie/beego/context"
	"github.com/coreos/go-oidc"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

func AuthFilter(context *bgCtx.Context) {
	log.Println("Start auth filter..")
	tsRaw := context.Input.Session("token_source")
	if tsRaw == nil {
		log.Println("There are no token source in the session")
		startAuth(context)
		return
	}
	ts := tsRaw.(oauth2.TokenSource)
	token, err := ts.Token()
	if err != nil {
		log.Printf("Token source presented in the session is not valid")
		startAuth(context)
		return
	}
	idToken, err := appCtx.GetAuthConfig().Verifier.Verify(ctx.Background(), token.AccessToken)
	if err != nil {
		log.Printf("Token presented in the session is not valid")
		startAuth(context)
		return
	}
	realmRoles := getRealmRoles(context, idToken)
	log.Printf("Roles %s has been retrieved from the token", realmRoles)
	resourceRoles := getResourceAccessValues(context, idToken)
	log.Printf("ResourceAccess %s has been retrieved from the token", resourceRoles)
	context.Output.Session("resource_access", resourceRoles)
	context.Output.Session("realm_roles", realmRoles)
	username := getUserInfoFromToken(context, idToken, "name")
	log.Printf("Username {%s} has been fetched from token", username)
	context.Output.Session("username", username)
}

func getRealmRoles(context *bgCtx.Context, token *oidc.IDToken) []string {
	log.Printf("Start to check roles ...")
	var claim map[string]*json.RawMessage
	err := token.Claims(&claim)
	if err != nil {
		log.Printf("Error has been occurred during the parsing token %+v", token)
		context.Abort(200, "500")
	}
	var realmAccess map[string]*[]string
	err = json.Unmarshal(*claim["realm_access"], &realmAccess)
	if err != nil {
		log.Printf("Error has been occurred during the parsing token %+v", token)
		context.Abort(200, "500")
	}

	return *realmAccess["roles"]
}

func getResourceAccessValues(context *bgCtx.Context, token *oidc.IDToken) map[string][]string {
	log.Printf("Start to check roles ...")
	var claim map[string]*json.RawMessage
	err := token.Claims(&claim)
	if err != nil {
		log.Printf("Error has been occurred during the parsing token %+v", token)
		context.Abort(200, "500")
	}
	var resourceAccess map[string]*map[string][]string
	err = json.Unmarshal(*claim["resource_access"], &resourceAccess)
	if err != nil {
		log.Printf("Error has been occurred during the parsing token %+v", token)
		context.Abort(200, "500")
	}

	instances := make(map[string][]string, len(resourceAccess))
	for key, value := range resourceAccess {
		var val = *value
		instances[key] = val["roles"]
	}
	return instances
}

func startAuth(context *bgCtx.Context) {
	authConfig := appCtx.GetAuthConfig()
	state := uuid.NewV4().String()
	log.Printf("State %s has been generated, saved in the session and added in the auth request", state)
	context.Output.Session(authConfig.StateAuthKey, state)
	if context.Request.Method == "GET" {
		context.Output.Session("request_path", context.Request.URL.Path)
	}
	context.Redirect(http.StatusFound, authConfig.Oauth2Config.AuthCodeURL(state))
}

func getUserInfoFromToken(context *bgCtx.Context, token *oidc.IDToken, userKey string) string {
	var claim map[string]*json.RawMessage
	err := token.Claims(&claim)
	if err != nil {
		log.Printf("Error has been occurred during the parsing token %+v", token)
		context.Abort(200, "500")
	}
	return strings.Replace(string(*claim[userKey]), "\"", "", -1)
}
