package models

import (
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
)

type Instance struct {
	Name     string                           `json:"name"`
	Type     catalogModels.InstanceType       `json:"type"`
	Bindings []catalogModels.InstanceBindings `json:"bindings"`
	Metadata []catalogModels.Metadata         `json:"metadata"`
}
