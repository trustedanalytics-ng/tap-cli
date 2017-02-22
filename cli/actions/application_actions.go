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
	"fmt"
	"io/ioutil"
	"os"
	"time"

	apiServiceModels "github.com/trustedanalytics-ng/tap-api-service/models"
	catalogModels "github.com/trustedanalytics-ng/tap-catalog/models"
	"github.com/trustedanalytics-ng/tap-cli/cli/archiver"
	"github.com/trustedanalytics-ng/tap-cli/cli/converter"
	"github.com/trustedanalytics-ng/tap-cli/cli/printer"
)

func (a *ActionsConfig) PushApplication(blobPath string, pushTimeout time.Duration) error {
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

	app, err := a.ApiService.CreateApplicationInstance(blob, manifest, pushTimeout)
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

func (a *ActionsConfig) CompressCwdAndPushAsApplication(pushTimeout time.Duration) error {
	folder, err := os.Getwd()
	if err != nil {
		return err
	}
	archivePath, err := archiver.CreateApplicationArchive(folder)
	if err != nil {
		return err
	}
	err = a.PushApplication(archivePath, pushTimeout)
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
	return a.deleteInstance(a.ApiService.DeleteApplicationInstance, catalogModels.InstanceTypeApplication, applicationName)
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
