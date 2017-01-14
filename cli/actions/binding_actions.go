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
	"errors"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	"github.com/trustedanalytics/tap-cli/cli/printer"
)

type bindingOperationType string

const (
	bind   bindingOperationType = "bind"
	unbind bindingOperationType = "unbind"
)

type BindableInstance struct {
	Name string
	Type catalogModels.InstanceType
}

func (a *ActionsConfig) BindInstance(srcInstance, dstInstance BindableInstance) error {
	return a.changeInstanceBinding(bind, srcInstance, dstInstance)
}

func (a *ActionsConfig) UnbindInstance(srcInstance, dstInstance BindableInstance) error {
	return a.changeInstanceBinding(unbind, srcInstance, dstInstance)
}

func (a *ActionsConfig) changeInstanceBinding(operationType bindingOperationType, srcInstance, dstInstance BindableInstance) error {
	srcInstanceID, srcInstanceType, err := converter.FetchInstanceIDandType(a.Config, srcInstance.Type, srcInstance.Name)
	if err != nil {
		return err
	}

	dstInstanceID, dstInstanceType, err := converter.FetchInstanceIDandType(a.Config, dstInstance.Type, dstInstance.Name)
	if err != nil {
		return err
	}

	instanceBinding := apiServiceModels.InstanceBindingRequest{}
	if srcInstanceType == catalogModels.InstanceTypeApplication {
		instanceBinding.ApplicationId = srcInstanceID
	} else if srcInstanceType == catalogModels.InstanceTypeService {
		instanceBinding.ServiceId = srcInstanceID
	}

	if operationType == bind && dstInstanceType == catalogModels.InstanceTypeApplication {
		_, err = a.ApiService.BindToApplicationInstance(instanceBinding, dstInstanceID)
	} else if operationType == bind && dstInstanceType == catalogModels.InstanceTypeService {
		_, err = a.ApiService.BindToServiceInstance(instanceBinding, dstInstanceID)
	} else if operationType == unbind {
		err = a.handleUnbindOperation(srcInstanceID, srcInstanceType, dstInstanceType, dstInstanceID)
	} else {
		err = errors.New("Cannot " + string(operationType) + " instance of type: " + string(dstInstanceType))
	}

	if err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) GetInstanceBindings(instance BindableInstance) error {
	instanceID, instanceType, err := converter.FetchInstanceIDandType(a.Config, instance.Type, instance.Name)
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

	printBindings(bindings)
	return nil
}

func printBindings(bindings apiServiceModels.InstanceBindings) {
	printableBindings := []printer.Printable{}
	for _, resource := range bindings.Resources {
		printableBindings = append(printableBindings, printer.PrintableResource{InstanceBindingsResource: resource})
	}
	printer.PrintTable(printableBindings)
}

func (a *ActionsConfig) handleUnbindOperation(srcID string, srcType catalogModels.InstanceType, dstType catalogModels.InstanceType, dstID string) error {
	var err error
	if srcType == catalogModels.InstanceTypeApplication && dstType == catalogModels.InstanceTypeApplication {
		_, err = a.ApiService.UnbindApplicationFromApplicationInstance(srcID, dstID)
	}
	if srcType == catalogModels.InstanceTypeService && dstType == catalogModels.InstanceTypeApplication {
		_, err = a.ApiService.UnbindServiceFromApplicationInstance(srcID, dstID)
	}
	if srcType == catalogModels.InstanceTypeApplication && dstType == catalogModels.InstanceTypeService {
		_, err = a.ApiService.UnbindApplicationFromServiceInstance(srcID, dstID)
	}
	if srcType == catalogModels.InstanceTypeService && dstType == catalogModels.InstanceTypeService {
		_, err = a.ApiService.UnbindServiceFromServiceInstance(srcID, dstID)
	}
	if err != nil {
		return err
	}
	return nil
}
