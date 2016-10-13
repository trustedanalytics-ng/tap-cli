package models

import (
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
)

type Instance struct {
	Name     string                           `json:"name"`
	Type     catalogModels.InstanceType       `json:"type"`
	ClassId  string                           `json:"classId"`
	Bindings []catalogModels.InstanceBindings `json:"bindings"`
	Metadata []catalogModels.Metadata         `json:"metadata"`
}

type InstanceBindings struct {
	Resources []InstanceBindingsResource `json:"resources"`
}

type InstanceBindingsResource struct {
	InstanceBindingsEntity `json:"entity"`
}

type InstanceBindingRequest struct {
	ApplicationId string `json:"application_id"`
	ServiceId     string `json:"service_id"`
}

type InstanceBindingsEntity struct {
	AppGuid             string `json:"app_guid"`
	ServiceInstanceGuid string `json:"service_instance_guid"`
	ServiceInstanceName string `json:"service_instance_name"`
}
