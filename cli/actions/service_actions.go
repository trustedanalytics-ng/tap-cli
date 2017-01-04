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

package actions

import (
	"fmt"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	"github.com/trustedanalytics/tap-cli/cli/printer"
	containerBrokerModels "github.com/trustedanalytics/tap-container-broker/models"
)

func (a *ActionsConfig) CreateServiceInstance(serviceName, planName, customName string, envs map[string]string) error {
	serviceID, planID, err := converter.FetchServiceAndPlanID(a.Config, serviceName, planName)
	if err != nil {
		return err
	}

	//TODO DPNG-11398: this should be move to api-service
	instanceBody := apiServiceModels.ServiceInstanceRequest{}
	instanceBody.Type = catalogModels.InstanceTypeService
	instanceBody.OfferingId = serviceID
	planMeta := catalogModels.Metadata{Id: catalogModels.OFFERING_PLAN_ID, Value: planID}
	instanceBody.Metadata = append(instanceBody.Metadata, planMeta)
	instanceBody.Name = customName
	for key, value := range envs {
		instanceBody.Metadata = append(instanceBody.Metadata, catalogModels.Metadata{
			Id:    key,
			Value: value,
		})
	}

	if _, err = a.ApiService.CreateServiceInstance(instanceBody); err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) GetService(serviceName string) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		return err
	}
	serviceInstance, err := a.ApiService.GetServiceInstance(instanceID)
	if err != nil {
		return err
	}

	printer.PrintFormattedJSON(serviceInstance)

	return nil
}

func (a *ActionsConfig) ListServices() error {
	services, err := a.ApiService.ListServiceInstances()
	if err != nil {
		fmt.Println("Retrieving services list failed")
		return err
	}
	printServices(services)
	return nil
}

func printServices(services []apiServiceModels.ServiceInstance) {
	printableServices := []printer.Printable{}
	for _, s := range services {
		printableServices = append(printableServices, printer.PrintableService{ServiceInstance: s})
	}
	printer.PrintTable(printableServices)
}

func (a *ActionsConfig) DeleteService(serviceName string) error {
	return a.deleteInstance(a.ApiService.DeleteServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) StartService(serviceName string) error {
	return a.changeState(a.ApiService.StartServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) RestartService(serviceName string) error {
	return a.changeState(a.ApiService.RestartServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) StopService(serviceName string) error {
	return a.changeState(a.ApiService.StopServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) GetServiceCredentials(instanceName string) error {
	instanceID, instanceType, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, instanceName)
	if err != nil {
		return err
	}

	creds := []containerBrokerModels.ContainerCredenials{}
	if instanceType == catalogModels.InstanceTypeService {
		creds, err = a.ApiService.GetInstanceCredentials(instanceID)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%q is not a service\n", instanceName)
	}

	for _, cred := range creds {
		printer.PrintFormattedJSON(cred)
		fmt.Println()
	}

	return nil
}

func (a *ActionsConfig) ExposeService(serviceID string, shouldExpose bool) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, serviceID)
	if err != nil {
		return err
	}

	hosts, _, err := a.ApiService.ExposeService(instanceID, shouldExpose)
	if err != nil {
		return err
	}

	printer.PrintFormattedJSON(hosts)
	return nil
}
