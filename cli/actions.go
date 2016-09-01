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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-api-service/models"
	consoleServiceModels "github.com/trustedanalytics/tap-api-service/models"
)

func Login(address string, username string, password string) error {
	creds := api.Credentials{}
	creds.Address = address
	creds.Username = username

	fmt.Println("Authenticating...")

	err := api.InitBasicAuthConnection(address, username, password)
	if err != nil {
		fmt.Println("error creating connection:", err)
		return err
	}

	loginResp, status, err := api.ConnectionConfig.ConsoleServiceLoginApi.Login()
	if status == http.StatusUnauthorized {
		fmt.Println("Authentication failed")
		return err
	} else if err != nil {
		fmt.Println("Error connecting: ", err)
		return err
	}

	creds.Token = loginResp.AccessToken
	creds.TokenType = loginResp.TokenType
	creds.ExpiresIn = loginResp.ExpiresIn

	err = api.SetCredentials(creds)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Authentication succeeded")

	return nil
}

func InviteUser(email string) error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = api.ConnectionConfig.ConsoleServiceApi.InviteUser(email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("User %s successfully invited\n", email)

	return nil
}

func DeleteUser(email string) error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = api.ConnectionConfig.ConsoleServiceApi.DeleteUser(email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("User %s successfully removed\n", email)

	return nil
}

func Catalog() error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	servicesList, err := api.ConnectionConfig.ConsoleServiceApi.GetCatalog()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printCatalog(servicesList)

	return nil

}

func Target() error {
	creds, err := api.GetCredentials()
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

func CreateOffer(jsonFilename string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

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

	_, err = api.ConnectionConfig.ConsoleServiceApi.CreateOffer(serviceWithTemplate)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

func CreateServiceInstance(serviceName, planName, customName string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	serviceId, planId, err := convert(serviceName, planName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceBody := models.Instance{}
	instanceBody.Type = catalogModels.InstanceTypeService
	planMeta := catalogModels.Metadata{Id: "plan", Value: planId}
	instanceBody.Metadata = append(instanceBody.Metadata, planMeta)
	instanceBody.Name = customName

	_, err = api.ConnectionConfig.ConsoleServiceApi.CreateServiceInstance(serviceId, instanceBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil

}

func DeleteInstance(serviceName string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = api.ConnectionConfig.ConsoleServiceApi.DeleteServiceInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

func GetInstanceBindings(instanceName string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(InstanceTypeBoth, instanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	bindings, err := api.ConnectionConfig.ConsoleServiceApi.GetInstanceBindings(instanceId)
	if err != nil {
		fmt.Printf("ERROR: %v", err.Error())
		return err
	}

	printInstancesBindings(bindings)
	return nil
}

func BindInstance(srcInstanceName, dstInstanceName string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	srcInstanceId, err := convertInstance(InstanceTypeBoth, srcInstanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dstInstanceId, err := convertInstance(InstanceTypeBoth, dstInstanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = api.ConnectionConfig.ConsoleServiceApi.BindInstance(srcInstanceId, dstInstanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil

}

func UnbindInstance(srcInstanceId, dstInstanceId string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = api.ConnectionConfig.ConsoleServiceApi.UnbindInstance(srcInstanceId, dstInstanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil

}

func ListApplications() error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	applicationInstances, err := api.ConnectionConfig.ConsoleServiceApi.ListApplicationInstances()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printApplicationInstances(applicationInstances)
	return nil
}

func GetApplication(applicationName string) error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	applicationInstance, err := api.ConnectionConfig.ConsoleServiceApi.GetApplicationInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printInstanceDetails(applicationInstance)
	fmt.Println("OK")
	return nil
}

func GetService(serviceName string) error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	serviceInstance, err := api.ConnectionConfig.ConsoleServiceApi.GetServiceInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printInstanceDetails(serviceInstance)
	fmt.Println("OK")
	return nil
}

func ListServices() error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	services, err := api.ConnectionConfig.ConsoleServiceApi.ListServiceInstances()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printServices(services)
	return nil
}

func ScaleApplication(applicationName string, replication int) error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	message, err := api.ConnectionConfig.ConsoleServiceApi.ScaleApplicationInstance(instanceId, replication)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(message.Message)

	return nil
}

func StartApplication(instanceId string) error {
	return ScaleApplication(instanceId, 1)
}

func StopApplication(instanceId string) error {
	return ScaleApplication(instanceId, 0)
}

func PushApplication(blob_path string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

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

	appInstance, err := api.ConnectionConfig.ConsoleServiceApi.CreateApplicationInstance(blob, manifest)
	if err != nil {
		fmt.Printf("ERROR: %v", err.Error())
		return err
	}

	printApplicationInstances([]models.ApplicationInstance{appInstance})

	if manifest.Instances != 1 {
		ScaleApplication(appInstance.Id, manifest.Instances)
		printApplicationInstances([]models.ApplicationInstance{appInstance})
	}

	fmt.Println("OK")
	return nil
}

func CompressCwdAndPushAsApplication() error {
	folder, err := os.Getwd()
	if err != nil {
		return err
	}
	archivePath, err := createApplicationArchive(folder)
	if err != nil {
		return err
	}
	err = PushApplication(archivePath)
	err2 := os.Remove(archivePath)
	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func GetInstanceLogs(instanceName string) error {

	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(InstanceTypeBoth, instanceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	logs, err := api.ConnectionConfig.ConsoleServiceApi.GetInstanceLogs(instanceId)
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

func DeleteApplication(applicationName string) error {
	err := api.InitOAuth2Connection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(catalogModels.InstanceTypeApplication, applicationName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = api.ConnectionConfig.ConsoleServiceApi.DeleteApplicationInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

