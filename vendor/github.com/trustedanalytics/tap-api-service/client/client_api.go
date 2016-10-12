package client

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/trustedanalytics/tap-api-service/models"
	userManagement "github.com/trustedanalytics/tap-api-service/user-management-connector"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	containerBrokerModels "github.com/trustedanalytics/tap-container-broker/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
	"github.com/trustedanalytics/tap-go-common/logger"
)

var logger = logger_wrapper.InitLogger("client")

type TapApiServiceApi interface {
	GetPlatformInfo() (models.PlatformInfo, error)

	GetInstanceBindings(instanceId string) (models.InstanceBindings, error)
	BindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error)
	UnbindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error)

	CreateApplicationInstance(blob multipart.File, manifest models.Manifest) (catalogModels.Application, error)
	CreateOffer(serviceWithTemplate models.ServiceDeploy) ([]catalogModels.Service, error)
	CreateServiceInstance(serviceId string, instance models.Instance) (containerBrokerModels.MessageResponse, error)

	DeleteOffering(serviceId string) error
	DeleteServiceInstance(instanceId string) error
	DeleteApplicationInstance(instanceId string) error

	GetCatalog() ([]models.Service, error)
	GetApplicationInstance(instanceId string) (models.ApplicationInstance, error)
	GetServiceInstance(instanceId string) (models.ServiceInstance, error)
	GetInstanceLogs(instanceId string) (map[string]string, error)
	GetInstanceCredentials(instanceId string) ([]containerBrokerModels.DeploymentEnvs, error)

	ListApplicationInstances() ([]models.ApplicationInstance, error)
	ListServiceInstances() ([]models.ServiceInstance, error)

	ScaleApplicationInstance(instanceId string, replication int) (containerBrokerModels.MessageResponse, error)
	ScaleServiceInstance(instanceId string, replication int) (containerBrokerModels.MessageResponse, error)

	InviteUser(email string) (userManagement.InvitationResponse, error)
	DeleteUser(email string) error
}

func NewTapApiServiceApiWithOAuth2(address, tokenType, token string) (TapApiServiceApi, error) {
	client, _, err := brokerHttp.GetHttpClient()
	if err != nil {
		return nil, err
	}
	return &TapApiServiceApiOAuth2Connector{address, tokenType, token, client}, nil
}

type TapApiServiceApiOAuth2Connector struct {
	Address   string
	TokenType string
	Token     string
	Client    *http.Client
}

func (c *TapApiServiceApiOAuth2Connector) getApiOAuth2Connector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		OAuth2: &brokerHttp.OAuth2{c.TokenType, c.Token},
		Client: c.Client,
		Url:    url,
	}
}

func (c *TapApiServiceApiOAuth2Connector) CreateServiceInstance(serviceId string, instance models.Instance) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s", c.Address, serviceId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PostModel(connector, instance, http.StatusAccepted, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) CreateOffer(serviceWithTemplate models.ServiceDeploy) ([]catalogModels.Service, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/offering", c.Address))
	result := &[]catalogModels.Service{}
	_, err := brokerHttp.PostModel(connector, serviceWithTemplate, http.StatusAccepted, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) DeleteOffering(serviceId string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/offering/%s", c.Address, serviceId))
	_, err := brokerHttp.DeleteModel(connector, http.StatusAccepted)
	return err
}

func (c *TapApiServiceApiOAuth2Connector) DeleteServiceInstance(instanceId string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s", c.Address, instanceId))
	_, err := brokerHttp.DeleteModel(connector, http.StatusAccepted)
	return err
}

func (c *TapApiServiceApiOAuth2Connector) GetPlatformInfo() (models.PlatformInfo, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/platform_info", c.Address))
	result := &models.PlatformInfo{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetCatalog() ([]models.Service, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/catalog", c.Address))
	result := &[]models.Service{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetInstanceLogs(instanceId string) (map[string]string, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/logs/%s", c.Address, instanceId))
	result := make(map[string]string)
	_, err := brokerHttp.GetModel(connector, http.StatusOK, &result)
	return result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetInstanceCredentials(instanceId string) ([]containerBrokerModels.DeploymentEnvs, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s/credentials", c.Address, instanceId))
	result := []containerBrokerModels.DeploymentEnvs{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, &result)
	return result, err
}

func (c *TapApiServiceApiOAuth2Connector) ListServiceInstances() ([]models.ServiceInstance, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services", c.Address))
	result := &[]models.ServiceInstance{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetServiceInstance(instanceId string) (models.ServiceInstance, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s", c.Address, instanceId))
	result := &models.ServiceInstance{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) ScaleServiceInstance(instanceId string, replication int) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s/scale", c.Address, instanceId))
	body := containerBrokerModels.ScaleInstanceRequest{
		Replicas: replication,
	}
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PutModel(connector, body, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetInstanceBindings(instanceId string) (models.InstanceBindings, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/bindings/%s", c.Address, instanceId))
	result := &models.InstanceBindings{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) BindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/bind/%s/%s", c.Address, srcInstanceId, dstInstanceId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PostModel(connector, "", http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) UnbindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/unbind/%s/%s", c.Address, srcInstanceId, dstInstanceId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PostModel(connector, "", http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) InviteUser(email string) (userManagement.InvitationResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users", c.Address))
	body := userManagement.InvitationRequest{
		Email: email,
	}
	result := &userManagement.InvitationResponse{}
	_, err := brokerHttp.PostModel(connector, body, http.StatusCreated, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) DeleteUser(email string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users", c.Address))
	body := userManagement.InvitationRequest{
		Email: email,
	}
	_, err := brokerHttp.DeleteModelWithBody(connector, body, http.StatusNoContent)
	return err
}
