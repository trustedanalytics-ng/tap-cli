package models

import (
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	templateRepositoryModels "github.com/trustedanalytics/tapng-template-repository/model"
)

type ServiceDeploy struct {
	Template templateRepositoryModels.Template `json:"template"`
	Service  catalogModels.Service             `json:"service"`
}
