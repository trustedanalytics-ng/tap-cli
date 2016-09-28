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

package cli

import (
	"errors"
	"fmt"

	catalogModels "github.com/trustedanalytics/tap-catalog/models"
)

func convertServiceAndPlanNameToId(a *ActionsConfig, serviceName, planName string) (string, string, error) {

	catalog, err := a.ApiService.GetCatalog()
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

func convertInstance(a *ActionsConfig, instanceType catalogModels.InstanceType, instanceName string) (string, error) {

	if instanceType == InstanceTypeBoth || instanceType == catalogModels.InstanceTypeService {
		serviceInstances, err := a.ApiService.ListServiceInstances()
		if err == nil {
			for _, instance := range serviceInstances {
				if instance.Name == instanceName {
					return instance.Id, nil
				}
			}
		}
	}
	if instanceType == InstanceTypeBoth || instanceType == catalogModels.InstanceTypeApplication {
		applicationInstances, err := a.ApiService.ListApplicationInstances()
		if err == nil {
			for _, instance := range applicationInstances {
				if instance.Name == instanceName {
					return instance.Id, nil
				}
			}
		}
	}

	return "", errors.New("cannot find instance with name: " + instanceName)
}

func getServiceID(a *ActionsConfig, serviceName string) (string, error) {

	services, err := a.ApiService.GetCatalog()
	if err != nil {
		return "", errors.New("Cannot fetch offering list: " + err.Error())
	}

	for _, service := range services {
		if service.Entity.Label == serviceName {
			return service.Entity.UniqueId, nil
		}
	}

	return "", fmt.Errorf("Service %s not found", serviceName)
}
