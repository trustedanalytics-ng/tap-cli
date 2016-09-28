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

	"github.com/trustedanalytics/tap-api-service/models"
	consoleServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
)

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
		fmt.Println("Authentication failed")
		return err
	} else if err != nil {
		fmt.Println("Error connecting:", err)
		return err
	}

	creds.Token = loginResp.AccessToken
	creds.TokenType = loginResp.TokenType
	creds.ExpiresIn = loginResp.ExpiresIn

	err = a.SetCredentials(creds)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Authentication succeeded")

	return nil
}

func (a *ActionsConfig) InviteUser(email string) error {

	_, err := a.ApiService.InviteUser(email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("User %s successfully invited\n", email)

	return nil
}

func (a *ActionsConfig) DeleteUser(email string) error {

	err := a.ApiService.DeleteUser(email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("User %s successfully removed\n", email)

	return nil
}

func (a *ActionsConfig) Catalog() error {

	servicesList, err := a.ApiService.GetCatalog()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printCatalog(servicesList)

	return nil

}

func (a *ActionsConfig) Target() error {
	creds, err := a.GetCredentials()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Please login first!")
		} else {
			fmt.Println(err)
		}
		return err
	}

	printCredentials(creds)

	return nil
}

func (a *ActionsConfig) CreateOffer(jsonFilename string) error {

	b, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		fmt.Println(err)
		return err
	}

	serviceWithTemplate := models.ServiceDeploy{}

	err = json.Unmarshal(b, &serviceWithTemplate)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i, service := range serviceWithTemplate.Services {
		for j, plan := range service.Plans {
			for k, dependency := range plan.Dependencies {
				serviceId, planId, err := convertServiceAndPlanNameToId(a, dependency.ServiceName, dependency.PlanName)
				if err != nil {
					fmt.Println(err)
					return err
				}
				plan.Dependencies[k].ServiceId = serviceId
				plan.Dependencies[k].PlanId = planId
			}
			service.Plans[j] = plan
		}
		serviceWithTemplate.Services[i] = service
	}

	_, err = a.ApiService.CreateOffer(serviceWithTemplate)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

func (a *ActionsConfig) DeleteOffering(serviceName string) error {

	serviceID, err := getServiceID(a, serviceName)
	if err != nil {
		err = errors.New("Cannot fetch service id: " + err.Error())
		fmt.Println(err)
		return err
	}

	if err = a.ApiService.DeleteOffering(serviceID); err != nil {
		err = errors.New("Cannot delete offering: " + err.Error())
		fmt.Println(err)
		return err
	}

	fmt.Printf("OK")
	return nil
}

func (a *ActionsConfig) CreateServiceInstance(serviceName, planName, customName string) error {

	serviceId, planId, err := convertServiceAndPlanNameToId(a, serviceName, planName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//TODO DPNG-11398: this should be move to api-service
	instanceBody := models.Instance{}
	instanceBody.Type = catalogModels.InstanceTypeService
	planMeta := catalogModels.Metadata{Id: catalogModels.OFFERING_PLAN_ID, Value: planId}
	instanceBody.Metadata = append(instanceBody.Metadata, planMeta)
	instanceBody.Name = customName

	_, err = a.ApiService.CreateServiceInstance(serviceId, instanceBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

func (a *ActionsConfig) DeleteInstance(serviceName string) error {

	instanceId, err := convertInstance(a, catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = a.ApiService.DeleteServiceInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

func (a *ActionsConfig) GetInstanceBindings(instanceName string) error {

	instanceId, err := convertInstance(a, InstanceTypeBoth, instanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	bindings, err := a.ApiService.GetInstanceBindings(instanceId)
	if err != nil {
		fmt.Printf("ERROR: %v", err.Error())
		return err
	}

	printInstancesBindings(bindings)
	return nil
}

func (a *ActionsConfig) BindInstance(srcInstanceName, dstInstanceName string) error {

	srcInstanceId, err := convertInstance(a, InstanceTypeBoth, srcInstanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dstInstanceId, err := convertInstance(a, InstanceTypeBoth, dstInstanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = a.ApiService.BindInstance(srcInstanceId, dstInstanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil

}

func (a *ActionsConfig) UnbindInstance(srcInstanceName, dstInstanceName string) error {

	srcInstanceId, err := convertInstance(a, InstanceTypeBoth, srcInstanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dstInstanceId, err := convertInstance(a, InstanceTypeBoth, dstInstanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = a.ApiService.UnbindInstance(srcInstanceId, dstInstanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil

}

func (a *ActionsConfig) ListApplications() error {

	applicationInstances, err := a.ApiService.ListApplicationInstances()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printApplicationInstances(applicationInstances)
	return nil
}

func (a *ActionsConfig) GetApplication(applicationName string) error {

	instanceId, err := convertInstance(a, catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	applicationInstance, err := a.ApiService.GetApplicationInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printInstanceDetails(applicationInstance)
	fmt.Println("OK")
	return nil
}

func (a *ActionsConfig) GetService(serviceName string) error {

	instanceId, err := convertInstance(a, catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	serviceInstance, err := a.ApiService.GetServiceInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printInstanceDetails(serviceInstance)
	fmt.Println("OK")
	return nil
}

func (a *ActionsConfig) ListServices() error {

	services, err := a.ApiService.ListServiceInstances()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printServices(services)
	return nil
}

func (a *ActionsConfig) ScaleApplication(applicationName string, replication int) error {

	instanceId, err := convertInstance(a, catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	message, err := a.ApiService.ScaleApplicationInstance(instanceId, replication)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return err
	}
	defer blob.Close()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return err
	}

	manifestBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/manifest.json", pwd))
	if err != nil {
		fmt.Println(err)
		return err
	}

	manifest := consoleServiceModels.Manifest{}
	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		fmt.Println(err)
		return err
	}

	app, err := a.ApiService.CreateApplicationInstance(blob, manifest)
	if err != nil {
		fmt.Printf("ERROR: %v", err.Error())
		return err
	}

	printApplication([]catalogModels.Application{app})

	if manifest.Instances != 1 {
		a.ScaleApplication(app.Id, manifest.Instances)
		printApplication([]catalogModels.Application{app})
	}

	fmt.Println("OK")
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

	instanceId, err := convertInstance(a, InstanceTypeBoth, instanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	logs, err := a.ApiService.GetInstanceLogs(instanceId)
	if err != nil {
		fmt.Printf("ERROR: %v", err.Error())
		return err
	}

	for container, log := range logs {
		fmt.Printf("%s:\n\n%s\n", container, log)
	}

	fmt.Println("OK")
	return nil
}

func (a *ActionsConfig) DeleteApplication(applicationName string) error {

	instanceId, err := convertInstance(a, catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = a.ApiService.DeleteApplicationInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}
