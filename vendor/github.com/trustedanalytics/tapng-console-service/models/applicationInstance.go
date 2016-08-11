package models

import (
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
)

type ApplicationInstance struct {
	catalogModels.Instance
	Replication int                      `json:"replication"`
	ImageState  catalogModels.ImageState `json:"imageState"`
}
