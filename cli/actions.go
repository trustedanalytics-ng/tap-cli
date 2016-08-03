package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-cli/api"
	"github.com/trustedanalytics/tapng-console-service/models"
	templateRepositoryModels "github.com/trustedanalytics/tapng-template-repository/model"
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

func Deploy(jsonFilename string) error {

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

	_, err = api.ConnectionConfig.ConsoleServiceApi.Deploy(serviceWithTemplate)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OK")
	return nil
}

func CreateInstance(serviceId string) error {

	err := api.InitConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	instanceBody := models.Instance{}
	instanceBody.Type = catalogModels.InstanceTypeService

	_, err = api.ConnectionConfig.ConsoleServiceApi.CreateInstance(serviceId, instanceBody)
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

func PushApplication(blob_path, image_path, template_path string) error {

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

	imageBytes, err := ioutil.ReadFile(image_path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	image := catalogModels.Image{}
	err = json.Unmarshal(imageBytes, &image)
	if err != nil {
		fmt.Println(err)
		return err
	}

	templateBytes, err := ioutil.ReadFile(template_path)
	template := templateRepositoryModels.Template{}
	err = json.Unmarshal(templateBytes, &template)
	if err != nil {
		fmt.Println(err)
		return err
	}

	application, err := api.ConnectionConfig.ConsoleServiceApi.CreateApplication(blob, image, template)
	if err != nil {
		fmt.Println(err)
		return err
	}

	printApplications([]catalogModels.Application{application})

	fmt.Println("OK")
	return nil

}
