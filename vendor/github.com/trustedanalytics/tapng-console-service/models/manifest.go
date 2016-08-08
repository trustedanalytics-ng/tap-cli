package models

import (
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
)

type Manifest struct {
	Name      string                  `json:"name"`
	ImageType catalogModels.ImageType `json:"type"`
	Instances int                     `json:"instances"`
}
