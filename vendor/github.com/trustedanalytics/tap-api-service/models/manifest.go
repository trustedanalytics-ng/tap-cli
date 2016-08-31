package models

import (
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
)

type Manifest struct {
	Name      string                  `json:"name"`
	ImageType catalogModels.ImageType `json:"type"`
	Instances int                     `json:"instances"`
}
