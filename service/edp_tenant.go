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

package service

import (
	"edp-admin-console/k8s"
	"github.com/astaxie/beego"
	"log"
	"strconv"
)

type EDPTenantService struct {
	Clients k8s.ClientSet
}

func (this EDPTenantService) GetVcsIntegrationValue() (bool, error) {
	res, err := strconv.ParseBool(beego.AppConfig.String("vcsIntegrationEnabled"))
	if err != nil {
		log.Printf("failed to read VCS_INTEGRATION_ENABLED value: %s", err)
		return false, err
	}

	return res, nil
}
