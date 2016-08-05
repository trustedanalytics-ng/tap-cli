package models

type Manifest struct {
	Name      string `json:"name"`
	ImageType string `json:"type"`
	Instances int    `json:"instances"`
}
