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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
	containerBrokerModels "github.com/trustedanalytics/tap-container-broker/models"
)

func announceSuccessfulOperation() {
	fmt.Println("OK")
}

type ActionsConfig struct {
	api.Config
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
		return err
	}

	fmt.Println("User password successfully changed")

	return nil
}

func (a *ActionsConfig) SendInvitation(email string) error {
	if _, err := a.ApiService.SendInvitation(email); err != nil {
		return err
	}

	fmt.Printf("User %q successfully invited\n", email)

	return nil
}

func (a *ActionsConfig) ResendInvitation(email string) error {
	if err := a.ApiService.ResendInvitation(email); err != nil {
		return err
	}

	fmt.Printf("User %q successfully invited\n", email)

	return nil
}

func (a *ActionsConfig) ListUsers() error {
	users, err := a.ApiService.GetUsers()
	if err != nil {
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
		return err
	}

	for _, email := range invitations {
		fmt.Println(email)
	}

	return nil
}

func (a *ActionsConfig) DeleteInvitation(email string) error {
	if err := a.ApiService.DeleteInvitation(email); err != nil {
		return err
	}

	fmt.Printf("Invitation for user %q successfully removed\n", email)

	return nil
}

func (a *ActionsConfig) DeleteUser(email string) error {
	if err := a.ApiService.DeleteUser(email); err != nil {
		return err
	}

	fmt.Printf("User %q successfully removed\n", email)

	return nil
}

func (a *ActionsConfig) Catalog() error {
	servicesList, err := a.ApiService.GetOfferings()
	if err != nil {
		return err
	}

	printCatalog(servicesList)

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

	printCredentials(creds)

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
				serviceId, planId, err := convertServiceAndPlanNameToId(a, dependency.ServiceName, dependency.PlanName)
				if err != nil {
					return err
				}
				plan.Dependencies[k].ServiceId = serviceId
				plan.Dependencies[k].PlanId = planId
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
	serviceID, err := getOfferingID(a, serviceName)
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
	serviceId, planId, err := convertServiceAndPlanNameToId(a, serviceName, planName)
	if err != nil {
		return err
	}

	//TODO DPNG-11398: this should be move to api-service
	instanceBody := apiServiceModels.Instance{}
	instanceBody.Type = catalogModels.InstanceTypeService
	instanceBody.ClassId = serviceId
	planMeta := catalogModels.Metadata{Id: catalogModels.OFFERING_PLAN_ID, Value: planId}
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

func (a *ActionsConfig) DeleteInstance(serviceName string) error {
	instanceId, _, err := convertInstance(a, catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		return err
	}

	if err = a.ApiService.DeleteServiceInstance(instanceId); err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) GetApplicationBindings(instanceName string) error {
	instanceId, instanceType, err := convertInstance(a, InstanceTypeBoth, instanceName)
	if err != nil {
		return err
	}

	var bindings apiServiceModels.InstanceBindings
	if instanceType == catalogModels.InstanceTypeApplication {
		bindings, err = a.ApiService.GetApplicationBindings(instanceId)
	} else if instanceType == catalogModels.InstanceTypeService {
		bindings, err = a.ApiService.GetServiceBindings(instanceId)
	}
	if err != nil {
		return err
	}

	printInstancesBindings(bindings)

	return nil
}

func (a *ActionsConfig) BindInstance(srcInstanceName, dstInstanceName string) error {
	srcInstanceId, srcInstanceType, err := convertInstance(a, InstanceTypeBoth, srcInstanceName)
	if err != nil {
		return err
	}

	dstInstanceId, dstInstanceType, err := convertInstance(a, InstanceTypeBoth, dstInstanceName)
	if err != nil {
		return err
	}

	instanceBinding := apiServiceModels.InstanceBindingRequest{}
	if srcInstanceType == catalogModels.InstanceTypeApplication {
		instanceBinding.ApplicationId = srcInstanceId
	} else if srcInstanceType == catalogModels.InstanceTypeService {
		instanceBinding.ServiceId = srcInstanceId
	}

	if dstInstanceType == catalogModels.InstanceTypeApplication {
		_, err = a.ApiService.BindToApplicationInstance(instanceBinding, dstInstanceId)
	} else if dstInstanceType == catalogModels.InstanceTypeService {
		_, err = a.ApiService.BindToServiceInstance(instanceBinding, dstInstanceId)
	}
	if err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) UnbindInstance(srcInstanceName, dstInstanceName string) error {
	srcInstanceId, srcInstanceType, err := convertInstance(a, InstanceTypeBoth, srcInstanceName)
	if err != nil {
		return err
	}

	dstInstanceId, dstInstanceType, err := convertInstance(a, InstanceTypeBoth, dstInstanceName)
	if err != nil {
		return err
	}

	instanceBinding := apiServiceModels.InstanceBindingRequest{}
	if srcInstanceType == catalogModels.InstanceTypeApplication {
		instanceBinding.ApplicationId = srcInstanceId
	} else if srcInstanceType == catalogModels.InstanceTypeService {
		instanceBinding.ServiceId = srcInstanceId
	}

	if dstInstanceType == catalogModels.InstanceTypeApplication {
		_, err = a.ApiService.UnbindFromApplicationInstance(instanceBinding, dstInstanceId)
	} else if dstInstanceType == catalogModels.InstanceTypeService {
		_, err = a.ApiService.UnbindFromServiceInstance(instanceBinding, dstInstanceId)
	}
	if err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) ListApplications() error {
	applicationInstances, err := a.ApiService.ListApplicationInstances()
	if err != nil {
		return err
	}

	printApplicationInstances(applicationInstances)

	return nil
}

func (a *ActionsConfig) GetApplication(applicationName string) error {
	applicationId, err := getApplicationID(a, applicationName)
	if err != nil {
		return err
	}

	applicationInstance, err := a.ApiService.GetApplicationInstance(applicationId)
	if err != nil {
		return err
	}

	printFormattedDetails(applicationInstance)

	return nil
}

func (a *ActionsConfig) GetService(serviceName string) error {
	instanceId, _, err := convertInstance(a, catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		return err
	}
	serviceInstance, err := a.ApiService.GetServiceInstance(instanceId)
	if err != nil {
		return err
	}

	printFormattedDetails(serviceInstance)

	return nil
}

func (a *ActionsConfig) ListServices() error {
	services, err := a.ApiService.ListServiceInstances()
	if err != nil {
		return err
	}

	printServices(services)

	return nil
}

func (a *ActionsConfig) ScaleApplication(applicationName string, replication int) error {
	instanceId, _, err := convertInstance(a, catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		return err
	}

	message, err := a.ApiService.ScaleApplicationInstance(instanceId, replication)
	if err != nil {
		return err
	}
	fmt.Println(message.Message)

	return nil
}

func (a *ActionsConfig) StartApplication(instanceId string) error {
	return a.ScaleApplication(instanceId, 1)
}

func (a *ActionsConfig) StopApplication(instanceId string) error {
	return a.ScaleApplication(instanceId, 0)
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

	if err := convertBindingsList(a, manifest.Bindings); err != nil {
		return err
	}

	app, err := a.ApiService.CreateApplicationInstance(blob, manifest)
	if err != nil {
		return err
	}

	printApplication([]catalogModels.Application{app})

	if manifest.Instances != 1 {
		a.ScaleApplication(app.Id, manifest.Instances)
		printApplication([]catalogModels.Application{app})
	}

	return nil
}

func (a *ActionsConfig) CompressCwdAndPushAsApplication() error {
	folder, err := os.Getwd()
	if err != nil {
		return err
	}
	archivePath, err := createApplicationArchive(folder)
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
	instanceId, instanceType, err := convertInstance(a, InstanceTypeBoth, instanceName)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	logs := make(map[string]string)
	if instanceType == catalogModels.InstanceTypeApplication {
		logs, err = a.ApiService.GetApplicationLogs(instanceId)
		if err != nil {
			return err
		}
	}
	if instanceType == catalogModels.InstanceTypeService {
		logs, err = a.ApiService.GetServiceLogs(instanceId)
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
	instanceId, instanceType, err := convertInstance(a, InstanceTypeBoth, instanceName)
	if err != nil {
		return err
	}

	creds := []containerBrokerModels.ContainerCredenials{}
	if instanceType == catalogModels.InstanceTypeService {
		creds, err = a.ApiService.GetInstanceCredentials(instanceId)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%q is not a service\n", instanceName)
	}

	for _, cred := range creds {
		printFormattedDetails(cred)
		fmt.Println()
	}

	return nil
}

func (a *ActionsConfig) DeleteApplication(applicationName string) error {
	instanceId, _, err := convertInstance(a, catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		return err
	}

	if err = a.ApiService.DeleteApplicationInstance(instanceId); err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}
