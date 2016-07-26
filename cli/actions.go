package cli

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/trustedanalytics/tapng-console-service/models"
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-cli/api"
	"os"
)


func Login(address string, username string, password string) error{
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