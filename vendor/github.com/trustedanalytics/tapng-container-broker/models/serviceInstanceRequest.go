package models

type CreateInstanceRequest struct {
	Image      string `json:"image"`
	TemplateId string `json:"template_id"`
}

type ScaleInstanceRequest struct {
	Replicas int `json:"replicas"`
}
