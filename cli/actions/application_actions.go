package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/cli/archiver"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	"github.com/trustedanalytics/tap-cli/cli/printer"
)

func (a *ActionsConfig) PushApplication(blobPath string) error {
	blob, err := os.Open(blobPath)
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

	printApplication(app)
	return nil
}

func printApplication(app catalogModels.Application) {
	printableApplications := []printer.Printable{printer.PrintableRecentlyPushedApplication{Application: app}}
	printer.PrintTable(printableApplications)
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

func (a *ActionsConfig) GetApplication(applicationName string) error {
	applicationID, err := converter.GetApplicationID(a.Config, applicationName)
	if err != nil {
		return err
	}

	applicationInstance, err := a.ApiService.GetApplicationInstance(applicationID)
	if err != nil {
		return err
	}

	printer.PrintFormattedJSON(applicationInstance)

	return nil
}

func (a *ActionsConfig) ListApplications() error {
	applicationInstances, err := a.ApiService.ListApplicationInstances()
	if err != nil {
		fmt.Println("Retrieving applications list failed")
		return err
	}
	printApplications(applicationInstances)
	return nil
}

func printApplications(applications []apiServiceModels.ApplicationInstance) {
	printableApplications := []printer.Printable{}
	for _, app := range applications {
		printableApplications = append(printableApplications, printer.PrintableApplication{ApplicationInstance: app})
	}
	printer.PrintTable(printableApplications)
}

func (a *ActionsConfig) DeleteApplication(applicationName string) error {
	return a.deleteInstance(catalogModels.InstanceTypeApplication, applicationName)
}

func (a *ActionsConfig) StartApplication(applicationName string) error {
	return a.changeState(a.ApiService.StartApplicationInstance, catalogModels.InstanceTypeApplication, applicationName)
}

func (a *ActionsConfig) RestartApplication(applicationName string) error {
	return a.changeState(a.ApiService.RestartApplicationInstance, catalogModels.InstanceTypeApplication, applicationName)
}

func (a *ActionsConfig) StopApplication(applicationName string) error {
	return a.changeState(a.ApiService.StopApplicationInstance, catalogModels.InstanceTypeApplication, applicationName)
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
