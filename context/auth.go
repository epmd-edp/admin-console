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

package context

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
)

type KeycloakParameters struct {
	KeycloakURL      string
	ClientId         string
	ClientSecret     string
	Host             string
	CallBackEndpoint string
	StateAuthKey     string
}

type AuthConfig struct {
	Oauth2Config oauth2.Config
	Verifier     *oidc.IDTokenVerifier
	StateAuthKey string
}

var authConfig AuthConfig

func InitAuth() {
	parameters := getParameters()
	log.Println(fmt.Sprintf("Keycloak has been retrieved: %s. \n"+
		"ClientId has been retrieved: %s. \n"+
		"ClientSecret has been retrieved. \n"+
		"Host has been retrieved: %s. \n"+
		"CallBackEndpoint has been retrieved: %s. \n"+
		"StateAuthKey has been retrieved: %s.",
		parameters.KeycloakURL, parameters.ClientId,
		parameters.Host, parameters.CallBackEndpoint, parameters.StateAuthKey))

	provider, err := oidc.NewProvider(context.Background(), parameters.KeycloakURL)

	if err != nil {
		log.Fatal(err)
		return
	}

	oauth2Config := oauth2.Config{
		ClientID:     parameters.ClientId,
		ClientSecret: parameters.ClientSecret,
		RedirectURL:  parameters.Host + parameters.CallBackEndpoint,
		Endpoint:     provider.Endpoint(),
	}

	oidcConfig := &oidc.Config{
		ClientID: parameters.ClientId,
	}

	verifier := provider.Verifier(oidcConfig)

	authConfig = AuthConfig{oauth2Config, verifier, parameters.StateAuthKey}
}

func getParameters() KeycloakParameters {
	//todo add checking all variables
	return KeycloakParameters{
		beego.AppConfig.String("keycloakURL"),
		beego.AppConfig.String("clientId"),
		beego.AppConfig.String("clientSecret"),
		beego.AppConfig.String("host"),
		beego.AppConfig.String("callBackEndpoint"),
		beego.AppConfig.String("stateAuthKey")}
}

func GetAuthConfig() AuthConfig {
	return authConfig
}
