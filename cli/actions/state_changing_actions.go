package actions

import (
	"fmt"

	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	containerBrokerModels "github.com/trustedanalytics/tap-container-broker/models"
)

type stateChangingFunction func(string) (containerBrokerModels.MessageResponse, error)

func (a *ActionsConfig) StartService(serviceName string) error {
	return a.changeState(a.ApiService.StartServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) StartApplication(applicationName string) error {
	return a.changeState(a.ApiService.StartApplicationInstance, catalogModels.InstanceTypeApplication, applicationName)
}

func (a *ActionsConfig) RestartService(serviceName string) error {
	return a.changeState(a.ApiService.RestartServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) RestartApplication(applicationName string) error {
	return a.changeState(a.ApiService.RestartApplicationInstance, catalogModels.InstanceTypeApplication, applicationName)
}

func (a *ActionsConfig) StopService(serviceName string) error {
	return a.changeState(a.ApiService.StopServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) StopApplication(applicationName string) error {
	return a.changeState(a.ApiService.StopApplicationInstance, catalogModels.InstanceTypeApplication, applicationName)
}

func (a *ActionsConfig) changeState(scf stateChangingFunction, instanceType catalogModels.InstanceType, instanceName string) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, instanceType, instanceName)
	if err != nil {
		return err
	}

	message, err := scf(instanceID)
	if err != nil {
		return err
	}

	fmt.Println(message.Message)
	return nil
}
