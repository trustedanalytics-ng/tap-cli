package actions

import (
	"errors"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	"github.com/trustedanalytics/tap-cli/cli/printer"
)

type bindingOperationType string

const (
	BIND   bindingOperationType = "bind"
	UNBIND bindingOperationType = "unbind"
)

func (a *ActionsConfig) BindInstance(srcInstanceName, dstInstanceName string) error {
	return a.changeInstanceBinding(BIND, srcInstanceName, dstInstanceName)
}

func (a *ActionsConfig) UnbindInstance(srcInstanceName, dstInstanceName string) error {
	return a.changeInstanceBinding(UNBIND, srcInstanceName, dstInstanceName)
}

func (a *ActionsConfig) changeInstanceBinding(operationType bindingOperationType, srcInstanceName string, dstInstanceName string) error {
	srcInstanceID, srcInstanceType, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, srcInstanceName)
	if err != nil {
		return err
	}

	dstInstanceID, dstInstanceType, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, dstInstanceName)
	if err != nil {
		return err
	}

	instanceBinding := apiServiceModels.InstanceBindingRequest{}
	if srcInstanceType == catalogModels.InstanceTypeApplication {
		instanceBinding.ApplicationId = srcInstanceID
	} else if srcInstanceType == catalogModels.InstanceTypeService {
		instanceBinding.ServiceId = srcInstanceID
	}

	if operationType == BIND && dstInstanceType == catalogModels.InstanceTypeApplication {
		_, err = a.ApiService.BindToApplicationInstance(instanceBinding, dstInstanceID)
	} else if operationType == BIND && dstInstanceType == catalogModels.InstanceTypeService {
		_, err = a.ApiService.BindToServiceInstance(instanceBinding, dstInstanceID)
	} else if operationType == UNBIND && dstInstanceType == catalogModels.InstanceTypeApplication {
		_, err = a.ApiService.UnbindFromApplicationInstance(instanceBinding, dstInstanceID)
	} else if operationType == UNBIND && dstInstanceType == catalogModels.InstanceTypeService {
		_, err = a.ApiService.UnbindFromServiceInstance(instanceBinding, dstInstanceID)
	} else {
		err = errors.New("Cannot " + string(operationType) + " instance of type: " + string(dstInstanceType))
	}

	if err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) GetInstanceBindings(instanceName string) error {
	instanceID, instanceType, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, instanceName)
	if err != nil {
		return err
	}

	var bindings apiServiceModels.InstanceBindings
	if instanceType == catalogModels.InstanceTypeApplication {
		bindings, err = a.ApiService.GetApplicationBindings(instanceID)
	} else if instanceType == catalogModels.InstanceTypeService {
		bindings, err = a.ApiService.GetServiceBindings(instanceID)
	}
	if err != nil {
		return err
	}

	printer.PrintInstancesBindings(bindings)

	return nil
}
