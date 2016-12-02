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
	"net/http"
	"os"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-cli/cli/archiver"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	"github.com/trustedanalytics/tap-cli/cli/printer"
	containerBrokerModels "github.com/trustedanalytics/tap-container-broker/models"
)

const SuccessMessage = "OK"

type ActionsConfig struct {
	api.Config
}

func announceSuccessfulOperation() {
	fmt.Println("OK")
}

func (a *ActionsConfig) Login() error {
	address, username, _ := a.ApiServiceLogin.GetLoginCredentials()
	creds := api.Credentials{}
	creds.Address = address
	creds.Username = username

	fmt.Println("Authenticating...")

	loginResp, status, err := a.ApiServiceLogin.Login()
	if status == http.StatusUnauthorized {
		return fmt.Errorf("Authentication failed")
	} else if err != nil {
		return fmt.Errorf("Authentication failed: %v", err)
	}

	creds.Token = loginResp.AccessToken
	creds.TokenType = loginResp.TokenType
	creds.ExpiresIn = loginResp.ExpiresIn

	if err = a.SetCredentials(creds); err != nil {
		return err
	}

	fmt.Println("Authentication succeeded")

	return nil
}

func (a *ActionsConfig) ChangeCurrentUserPassword(currentPassword, newPassword string) error {
	if err := a.ApiService.ChangeCurrentUserPassword(currentPassword, newPassword); err != nil {
		fmt.Println("Changing user password failed")
		return err
	}

	fmt.Println("User password successfully changed")

	return nil
}

func (a *ActionsConfig) SendInvitation(email string) error {
	if _, err := a.ApiService.SendInvitation(email); err != nil {
		fmt.Printf("Sending invitation to email %s failed\n", email)
		return err
	}

	fmt.Printf("User %q successfully invited\n", email)

	return nil
}

func (a *ActionsConfig) ResendInvitation(email string) error {
	if err := a.ApiService.ResendInvitation(email); err != nil {
		fmt.Printf("Resending invitation to email %s failed\n", email)
		return err
	}

	fmt.Printf("User %q successfully invited\n", email)

	return nil
}

func (a *ActionsConfig) ListUsers() error {
	users, err := a.ApiService.GetUsers()
	if err != nil {
		fmt.Println("Listing users failed")
		return err
	}

	for _, user := range users {
		fmt.Println(user.Username)
	}

	return nil
}

func (a *ActionsConfig) ListInvitations() error {
	invitations, err := a.ApiService.GetInvitations()
	if err != nil {
		fmt.Println("Listing invitations failed")
		return err
	}

	for _, email := range invitations {
		fmt.Println(email)
	}

	return nil
}

func (a *ActionsConfig) DeleteInvitation(email string) error {
	if err := a.ApiService.DeleteInvitation(email); err != nil {
		fmt.Printf("Deleting invitation of user %s failed\n", email)
		return err
	}

	fmt.Printf("Invitation for user %q successfully removed\n", email)

	return nil
}

func (a *ActionsConfig) DeleteUser(email string) error {
	if err := a.ApiService.DeleteUser(email); err != nil {
		fmt.Printf("Deleting user %s failed\n", email)
		return err
	}

	fmt.Printf("User %q successfully removed\n", email)

	return nil
}

func (a *ActionsConfig) Catalog() error {
	servicesList, err := a.ApiService.GetOfferings()
	if err != nil {
		fmt.Println("Retrieving catalog failed")
		return err
	}

	printer.PrintCatalog(servicesList)

	return nil
}

func (a *ActionsConfig) Target() error {
	creds, err := a.GetCredentials()
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Please login first!")
		}
		return err
	}

	printer.PrintCredentials(creds)

	return nil
}

func (a *ActionsConfig) CreateOffer(jsonFilename string) error {
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

func (a *ActionsConfig) CreateServiceInstance(serviceName, planName, customName string, envs map[string]string) error {
	serviceID, planID, err := converter.FetchServiceAndPlanID(a.Config, serviceName, planName)
	if err != nil {
		return err
	}

	//TODO DPNG-11398: this should be move to api-service
	instanceBody := apiServiceModels.Instance{}
	instanceBody.Type = catalogModels.InstanceTypeService
	instanceBody.ClassId = serviceID
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

func (a *ActionsConfig) DeleteService(serviceName string) error {
	return a.deleteInstance(catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) DeleteApplication(applicationName string) error {
	return a.deleteInstance(catalogModels.InstanceTypeApplication, applicationName)
}

func (a *ActionsConfig) deleteInstance(instanceType catalogModels.InstanceType, instanceName string) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, instanceType, instanceName)
	if err != nil {
		return err
	}
	if err = a.ApiService.DeleteApplicationInstance(instanceID); err != nil {
		return err
	}
	announceSuccessfulOperation()
	return nil
}

func (a *ActionsConfig) ExposeService(serviceID string, shouldExpose bool) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, serviceID)
	if err != nil {
		return err
	}

	if hosts, _, err := a.ApiService.ExposeService(instanceID, shouldExpose); err != nil {
		return err
	} else {
		printer.PrintFormattedDetails(hosts)
		return nil
	}
}

func (a *ActionsConfig) ListApplications() error {
	applicationInstances, err := a.ApiService.ListApplicationInstances()
	if err != nil {
		fmt.Println("Retrieving applications list failed")
		return err
	}

	printer.PrintApplicationInstances(applicationInstances)

	return nil
}

func (a *ActionsConfig) GetApplication(applicationName string) error {
	applicationID, err := converter.GetApplicationID(a.Config, applicationName)
	if err != nil {
		return err
	}

	applicationInstance, err := a.ApiService.GetApplicationInstance(applicationID)
	if err != nil {
		return err
	}

	printer.PrintFormattedDetails(applicationInstance)

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

	printer.PrintFormattedDetails(serviceInstance)

	return nil
}

func (a *ActionsConfig) ListServices() error {
	services, err := a.ApiService.ListServiceInstances()
	if err != nil {
		fmt.Println("Retrieving services list failed")
		return err
	}

	printer.PrintServices(services)

	return nil
}

func (a *ActionsConfig) ScaleApplication(applicationName string, replication int) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		return err
	}

	message, err := a.ApiService.ScaleApplicationInstance(instanceID, replication)
	if err != nil {
		return err
	}
	fmt.Println(message.Message)
	return nil
}

func (a *ActionsConfig) PushApplication(blob_path string) error {
	blob, err := os.Open(blob_path)
	if err != nil {
		return err
	}
	defer blob.Close()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	manifestBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/manifest.json", pwd))
	if err != nil {
		return err
	}

	manifest := apiServiceModels.Manifest{}
	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		return err
	}

	app, err := a.ApiService.CreateApplicationInstance(blob, manifest)
	if err != nil {
		return err
	}

	printer.PrintApplication([]catalogModels.Application{app})
	return nil
}

func (a *ActionsConfig) CompressCwdAndPushAsApplication() error {
	folder, err := os.Getwd()
	if err != nil {
		return err
	}
	archivePath, err := archiver.CreateApplicationArchive(folder)
	if err != nil {
		return err
	}
	err = a.PushApplication(archivePath)
	err2 := os.Remove(archivePath)
	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func (a *ActionsConfig) GetInstanceLogs(instanceName string) error {
	instanceID, instanceType, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, instanceName)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	logs := make(map[string]string)
	if instanceType == catalogModels.InstanceTypeApplication {
		logs, err = a.ApiService.GetApplicationLogs(instanceID)
		if err != nil {
			return err
		}
	}
	if instanceType == catalogModels.InstanceTypeService {
		logs, err = a.ApiService.GetServiceLogs(instanceID)
		if err != nil {
			return err
		}
	}

	for container, log := range logs {
		fmt.Printf("%s:\n\n%s\n", container, log)
	}

	return nil
}

func (a *ActionsConfig) GetInstanceCredentials(instanceName string) error {
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
		printer.PrintFormattedDetails(cred)
		fmt.Println()
	}

	return nil
}
