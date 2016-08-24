package models

import (
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
)

type ApplicationInstance struct {
	catalogModels.Instance
	Replication      int                      `json:"replication"`
	ImageState       catalogModels.ImageState `json:"imageState"`
	Urls             []string                 `json:"urls"`
	Memory           string                   `json:"memory"`
	DiskQuota        string                   `json:"disk_quota"`
	RunningInstances int                      `json:"running_instances"`
}
