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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	"github.com/trustedanalytics/tap-cli/cli/printer"
)

func (a *ActionsConfig) CreateOffering(jsonFilename string) error {
	b, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		return err
	}

	serviceWithTemplate := apiServiceModels.ServiceDeploy{}

	err = json.Unmarshal(b, &serviceWithTemplate)
	if err != nil {
		return err
	}

	for i, service := range serviceWithTemplate.Services {
		for j, plan := range service.Plans {
			for k, dependency := range plan.Dependencies {
				serviceID, planID, err := converter.FetchServiceAndPlanID(
					a.Config, dependency.ServiceName, dependency.PlanName)
				if err != nil {
					return err
				}
				plan.Dependencies[k].ServiceId = serviceID
				plan.Dependencies[k].PlanId = planID
			}
			service.Plans[j] = plan
		}
		serviceWithTemplate.Services[i] = service
	}

	if _, err = a.ApiService.CreateOffer(serviceWithTemplate); err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) GetOffering(name string) error {
	offeringsList, err := a.ApiService.GetOfferings()
	if err != nil {
		fmt.Println("Retrieving catalog failed")
		return err
	}

	for _, of := range offeringsList {
		if of.Name == name {
			return printSingleOffering(of)
		}
	}

	return errors.New("Could not find offering with such name")
}

func printSingleOffering(of apiServiceModels.Offering) error {
	marshalled, err := json.MarshalIndent(of, "", "  ")
	if err != nil {
		fmt.Println("Could not marshal fetched data")
		return err
	}
	fmt.Println(string(marshalled))
	return nil
}

func (a *ActionsConfig) ListOfferings() error {
	offeringsList, err := a.ApiService.GetOfferings()
	if err != nil {
		fmt.Println("Retrieving catalog failed")
		return err
	}
	printOfferings(offeringsList)
	return nil
}

func printOfferings(offerings []apiServiceModels.Offering) {
	printableOfferings := []printer.Printable{}
	for _, of := range offerings {
		printableOfferings = append(printableOfferings, printer.PrintableOffering{Offering: of})
	}
	printer.PrintTable(printableOfferings)
}

func (a *ActionsConfig) DeleteOffering(serviceName string) error {
	serviceID, err := converter.GetOfferingID(a.Config, serviceName)
	if err != nil {
		return fmt.Errorf("Cannot fetch service id: %v", err.Error())
	}

	if err = a.ApiService.DeleteOffering(serviceID); err != nil {
		return fmt.Errorf("Cannot delete offering: %v", err.Error())
	}

	announceSuccessfulOperation()

	return nil
}
