package models

type CreateInstanceRequest struct {
	IsServiceBroker bool   `json:"isServiceBroker"`
	Image           string `json:"image"`
}

type ScaleInstanceRequest struct {
	Replicas int `json:"replicas"`
}
