package models

import (
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
)

type ServiceInstance struct {
	catalogModels.Instance
	ServiceName     string `json:"serviceName"`
	ServicePlanName string `json:"planName"`
}