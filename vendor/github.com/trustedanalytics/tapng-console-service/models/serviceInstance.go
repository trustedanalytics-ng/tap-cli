package models

import (
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
)

type ServiceInstance struct {
	catalogModels.Instance
	ServiceName     string `json:"name"`
	ServicePlanName string `json:"planName"`
}
