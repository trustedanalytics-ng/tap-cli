/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package converter

import (
	"errors"
	"fmt"
	"strings"

	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
)

func FetchServiceAndPlanID(apiConfig api.Config, serviceName, planName string) (string, string, error) {

	catalog, err := apiConfig.ApiService.GetOfferings()
	if err != nil {
		return "", "", err
	}

	for _, service := range catalog {

		if service.Entity.Label == serviceName {
			for _, plan := range service.Entity.ServicePlans {

				if plan.Entity.Name == planName {
					return service.Entity.UniqueId, plan.Entity.UniqueId, nil
				}
			}
			return "", "", errors.New("cannot find plan: '" + planName + "' for service: '" + serviceName + "'")
		}
	}

	return "", "", errors.New("cannot find service: '" + serviceName + "'")
}

const (
	InstanceTypeBoth catalogModels.InstanceType = "BOTH"
)

func FetchInstanceIDandType(apiConfig api.Config, instanceType catalogModels.InstanceType, instanceName string) (string, catalogModels.InstanceType, error) {

	if instanceType == InstanceTypeBoth || instanceType == catalogModels.InstanceTypeService {
		serviceInstances, err := apiConfig.ApiService.ListServiceInstances()
		if err == nil {
			for _, instance := range serviceInstances {
				if instance.Name == instanceName {
					return instance.Id, catalogModels.InstanceTypeService, nil
				}
			}
		}
	}
	if instanceType == InstanceTypeBoth || instanceType == catalogModels.InstanceTypeApplication {
		applicationInstances, err := apiConfig.ApiService.ListApplicationInstances()
		if err == nil {
			for _, instance := range applicationInstances {
				if instance.Name == instanceName {
					return instance.Id, catalogModels.InstanceTypeApplication, nil
				}
			}
		}
	}

	return "", "", errors.New("cannot find instance with name: " + instanceName)
}

func GetOfferingID(apiConfig api.Config, serviceName string) (string, error) {

	services, err := apiConfig.ApiService.GetOfferings()
	if err != nil {
		return "", errors.New("cannot fetch offering list: " + err.Error())
	}

	for _, service := range services {
		if service.Entity.Label == serviceName {
			return service.Entity.UniqueId, nil
		}
	}

	return "", fmt.Errorf("service %s not found", serviceName)
}

func ConvertBindingsList(apiConfig api.Config, bindings []string) error {
	instances, err := apiConfig.ApiService.ListServiceInstances()
	if err != nil {
		return errors.New("cannot fetch service intances: " + err.Error())
	}

	notFoundIDs := []string{}
	for i := 0; i < len(bindings); i++ {
		found := false
		for _, instance := range instances {
			if instance.Name == bindings[i] {
				bindings[i] = instance.Id
				found = true
				break
			}
		}
		if !found {
			notFoundIDs = append(notFoundIDs, bindings[i])
		}
	}

	if len(notFoundIDs) != 0 {
		return errors.New("following service instances don't exist: " + strings.Join(notFoundIDs, ", "))
	}

	return nil
}

func GetApplicationID(apiConfig api.Config, applicationName string) (string, error) {

	applications, err := apiConfig.ApiService.ListApplicationInstances()
	if err != nil {
		return "", errors.New("Cannot fetch applications list: " + err.Error())
	}

	for _, app := range applications {
		if app.Name == applicationName {
			return app.Id, nil
		}
	}

	return "", fmt.Errorf("Application %s not found", applicationName)
}
