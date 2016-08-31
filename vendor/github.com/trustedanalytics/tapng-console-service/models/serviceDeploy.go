package models

import (
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	templateRepositoryModels "github.com/trustedanalytics/tapng-template-repository/model"
)

type ServiceDeploy struct {
	BrokerName string                            `json:"broker_name"`
	Template   templateRepositoryModels.Template `json:"template"`
	Services   []catalogModels.Service           `json:"services"`
}
