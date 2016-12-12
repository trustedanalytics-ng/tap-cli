package actions

import (
	"fmt"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/cli/converter"
	"github.com/trustedanalytics/tap-cli/cli/printer"
)

func (a *ActionsConfig) CreateServiceInstance(serviceName, planName, customName string, envs map[string]string) error {
	serviceID, planID, err := converter.FetchServiceAndPlanID(a.Config, serviceName, planName)
	if err != nil {
		return err
	}

	//TODO DPNG-11398: this should be move to api-service
	instanceBody := apiServiceModels.Instance{}
	instanceBody.Type = catalogModels.InstanceTypeService
	instanceBody.ClassId = serviceID
	planMeta := catalogModels.Metadata{Id: catalogModels.OFFERING_PLAN_ID, Value: planID}
	instanceBody.Metadata = append(instanceBody.Metadata, planMeta)
	instanceBody.Name = customName
	for key, value := range envs {
		instanceBody.Metadata = append(instanceBody.Metadata, catalogModels.Metadata{
			Id:    key,
			Value: value,
		})
	}

	if _, err = a.ApiService.CreateServiceInstance(instanceBody); err != nil {
		return err
	}

	announceSuccessfulOperation()

	return nil
}

func (a *ActionsConfig) GetService(serviceName string) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, catalogModels.InstanceTypeService, serviceName)
	if err != nil {
		return err
	}
	serviceInstance, err := a.ApiService.GetServiceInstance(instanceID)
	if err != nil {
		return err
	}

	printer.PrintFormattedJSON(serviceInstance)

	return nil
}

func (a *ActionsConfig) ListServices() error {
	services, err := a.ApiService.ListServiceInstances()
	if err != nil {
		fmt.Println("Retrieving services list failed")
		return err
	}
	printServices(services)
	return nil
}

func printServices(services []apiServiceModels.ServiceInstance) {
	printableServices := []printer.Printable{}
	for _, s := range services {
		printableServices = append(printableServices, printer.PrintableService{ServiceInstance: s})
	}
	printer.PrintTable(printableServices)
}

func (a *ActionsConfig) DeleteService(serviceName string) error {
	return a.deleteInstance(a.ApiService.DeleteServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) StartService(serviceName string) error {
	return a.changeState(a.ApiService.StartServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) RestartService(serviceName string) error {
	return a.changeState(a.ApiService.RestartServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) StopService(serviceName string) error {
	return a.changeState(a.ApiService.StopServiceInstance, catalogModels.InstanceTypeService, serviceName)
}

func (a *ActionsConfig) ExposeService(serviceID string, shouldExpose bool) error {
	instanceID, _, err := converter.FetchInstanceIDandType(a.Config, converter.InstanceTypeBoth, serviceID)
	if err != nil {
		return err
	}

	hosts, _, err := a.ApiService.ExposeService(instanceID, shouldExpose)
	if err != nil {
		return err
	}

	printer.PrintFormattedJSON(hosts)
	return nil
}
