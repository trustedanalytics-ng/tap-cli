package models

type CreateInstanceRequest struct {
	IsServiceBroker bool `json:"isServiceBroker"`
}

type ScaleInstanceRequest struct {
	Replicas int `json:"replicas"`
}
