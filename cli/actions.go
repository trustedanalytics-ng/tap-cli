package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-cli/api"
	"github.com/trustedanalytics/tapng-console-service/models"
	consoleServiceModels "github.com/trustedanalytics/tapng-console-service/models"
)

func Login(address string, username string, password string) error {
	creds := api.Credentials{}
	creds.Address = address
	creds.Username = username
	creds.Password = password

	err := api.SetCredentials(creds)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Authenticating...")

	err = api.InitConnection()
	if err != nil {
		fmt.Println("error creating connection:", err)
		return err
	}

	_, err = api.ConnectionConfig.ConsoleServiceApi.GetCatalog()
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return err
	}

	fmt.Println("OK")

	return nil
}

func Catalog() error {
	err := api.InitConnection()
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

	err := api.InitConnection()
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

func CreateInstance(serviceName, planName, customName string) error {

	err := api.InitConnection()
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

	_, err = api.ConnectionConfig.ConsoleServiceApi.CreateInstance(serviceId, instanceBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil

}

func DeleteInstance(serviceName string) error {

	err := api.InitConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceId, err := convertInstance(serviceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = api.ConnectionConfig.ConsoleServiceApi.DeleteInstance(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

func BindInstance(srcInstanceId, dstInstanceId string) error {

	err := api.InitConnection()
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

	err := api.InitConnection()
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
	err := api.InitConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	applications, err := api.ConnectionConfig.ConsoleServiceApi.ListApplications()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printApplications(applications)
	return nil
}

func ListServices() error {
	err := api.InitConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	services, err := api.ConnectionConfig.ConsoleServiceApi.ListServicesInstances()
	if err != nil {
		fmt.Println(err)
		return err
	}

	printServices(services)
	return nil
}

func ScaleApplication(instanceId string, replication int) error {
	err := api.InitConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	message, err := api.ConnectionConfig.ConsoleServiceApi.ScaleInstance(instanceId, replication)
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

	err := api.InitConnection()
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

	appInstance, err := api.ConnectionConfig.ConsoleServiceApi.CreateApplication(blob, manifest)
	if err != nil {
		fmt.Printf("ERROR: %v", err.Error())
		return err
	}

	printAppInstance(appInstance, 1)

	if manifest.Instances != 1 {
		ScaleApplication(appInstance.Id, manifest.Instances)
		printAppInstance(appInstance, manifest.Instances)
	}

	fmt.Println("OK")
	return nil
}

func GetInstanceLogs(instanceId string) error {

	err := api.InitConnection()
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

func DeleteApplication(instanceId string) error {
	err := api.InitConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = api.ConnectionConfig.ConsoleServiceApi.DeleteApplication(instanceId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}
