package actions

import (
	"fmt"

	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	containerBrokerModels "github.com/trustedanalytics/tap-container-broker/models"
)

type stateChangingFunction func(string) (containerBrokerModels.MessageResponse, error)
type deletingFunction func(string) error

func (a *ActionsConfig) deleteInstance(df deletingFunction, instanceType catalogModels.InstanceType, instanceName string) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, instanceType, instanceName)
	if err != nil {
		return err
	}
	if err = df(instanceID); err != nil {
		return err
	}
	announceSuccessfulOperation()
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
