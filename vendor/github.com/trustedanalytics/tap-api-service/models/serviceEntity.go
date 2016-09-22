package models

type Service struct {
	Entity   ServiceEntity `json:"entity"`
	Metadata Metadata      `json:"metadata"`
}

type ServiceEntity struct {
	Label             string        `json:"label"`
	Provider          string        `json:"provider"`
	Url               string        `json:"url"`
	Description       string        `json:"description"`
	LongDescription   string        `json:"long_Description"`
	Version           string        `json:"version"`
	InfoUrl           string        `json:"info_url"`
	Active            bool          `json:"active"`
	Bindable          bool          `json:"bindable"`
	UniqueId          string        `json:"unique_id"`
	Extra             string        `json:"extra"`
	Tags              []string      `json:"tags"`
	Requires          []string      `json:"requires"`
	DocumentationUrl  string        `json:"documentation_url"`
	ServiceBrokerGuid string        `json:"service_broker_guid"`
	PlanUpdateable    bool          `json:"plan_updateable"`
	ServicePlansUrl   string        `json:"service_plans_url"`
	State             string        `json:"state"`
	ServicePlans      []ServicePlan `json:"service_plans"`
}

type ServicePlan struct {
	Entity   ServicePlanEntity `json:"entity"`
	Metadata Metadata          `json:"metadata"`
}

type ServicePlanEntity struct {
	Name                string   `json:"name"`
	Free                bool     `json:"free"`
	Description         string   `json:"description"`
	ServiceGuid         string   `json:"service_guid"`
	Extra               string   `json:"extra"`
	UniqueId            string   `json:"unique_id"`
	Public              bool     `json:"public"`
	Active              bool     `json:"active"`
	ServiceUrl          string   `json:"service_url"`
	Service             string   `json:"service"`
	ServiceInstancesUrl string   `json:"service_instances_url"`
	Metadata            Metadata `json:"metadata"`
}

type Metadata struct {
	Guid string `json:"guid"`
}
