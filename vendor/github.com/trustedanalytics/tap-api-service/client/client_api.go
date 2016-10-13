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

	GetApplicationBindings(applicationId string) (models.InstanceBindings, error)
	GetServiceBindings(serviceId string) (models.InstanceBindings, error)
	BindToApplicationInstance(bindingRequest models.InstanceBindingRequest, applicationId string) (containerBrokerModels.MessageResponse, error)
	BindToServiceInstance(bindingRequest models.InstanceBindingRequest, serviceId string) (containerBrokerModels.MessageResponse, error)
	UnbindFromApplicationInstance(bindingRequest models.InstanceBindingRequest, applicationId string) (int, error)
	UnbindFromServiceInstance(bindingRequest models.InstanceBindingRequest, serviceId string) (int, error)

	CreateApplicationInstance(blob multipart.File, manifest models.Manifest) (catalogModels.Application, error)
	CreateOffer(serviceWithTemplate models.ServiceDeploy) ([]catalogModels.Service, error)
	CreateServiceInstance(instance models.Instance) (containerBrokerModels.MessageResponse, error)

	DeleteOffering(serviceId string) error
	DeleteServiceInstance(instanceId string) error
	DeleteApplicationInstance(instanceId string) error

	GetOfferings() ([]models.Service, error)
	GetOffering(offeringId string) (models.Service, error)
	GetApplicationInstance(applicationId string) (models.ApplicationInstance, error)
	GetServiceInstance(serviceId string) (models.ServiceInstance, error)
	GetApplicationLogs(applicationId string) (map[string]string, error)
	GetServiceLogs(serviceId string) (map[string]string, error)
	GetInstanceCredentials(instanceId string) ([]containerBrokerModels.DeploymentEnvs, error)

	ListApplicationInstances() ([]models.ApplicationInstance, error)
	ListServiceInstances() ([]models.ServiceInstance, error)

	ScaleApplicationInstance(applicationId string, replication int) (containerBrokerModels.MessageResponse, error)
	ScaleServiceInstance(serviceId string, replication int) (containerBrokerModels.MessageResponse, error)

	GetInvitations() ([]string, error)
	SendInvitation(email string) (userManagement.InvitationResponse, error)
	ResendInvitation(email string) error
	DeleteInvitation(email string) error

	GetUsers() ([]userManagement.UaaUser, error)
	ChangeCurrentUserPassword(password, newPassword string) error
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

func (c *TapApiServiceApiOAuth2Connector) CreateServiceInstance(instance models.Instance) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services", c.Address))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PostModel(connector, instance, http.StatusAccepted, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) CreateOffer(serviceWithTemplate models.ServiceDeploy) ([]catalogModels.Service, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/offerings", c.Address))
	result := &[]catalogModels.Service{}
	_, err := brokerHttp.PostModel(connector, serviceWithTemplate, http.StatusAccepted, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) DeleteOffering(serviceId string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/offerings/%s", c.Address, serviceId))
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

func (c *TapApiServiceApiOAuth2Connector) GetOfferings() ([]models.Service, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/offerings", c.Address))
	result := &[]models.Service{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetOffering(offeringId string) (models.Service, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/offerings/%s", c.Address, offeringId))
	result := &models.Service{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetApplicationLogs(applicationId string) (map[string]string, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s/logs", c.Address, applicationId))
	result := make(map[string]string)
	_, err := brokerHttp.GetModel(connector, http.StatusOK, &result)
	return result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetServiceLogs(serviceId string) (map[string]string, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s/logs", c.Address, serviceId))
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

func (c *TapApiServiceApiOAuth2Connector) GetServiceInstance(serviceId string) (models.ServiceInstance, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s", c.Address, serviceId))
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

func (c *TapApiServiceApiOAuth2Connector) GetApplicationBindings(applicationId string) (models.InstanceBindings, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s/bindings", c.Address, applicationId))
	result := &models.InstanceBindings{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetServiceBindings(serviceId string) (models.InstanceBindings, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s/bindings", c.Address, serviceId))
	result := &models.InstanceBindings{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) BindToApplicationInstance(bindingRequest models.InstanceBindingRequest, applicationId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s/bindings", c.Address, applicationId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PostModel(connector, bindingRequest, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) BindToServiceInstance(bindingRequest models.InstanceBindingRequest, serviceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s/bindings", c.Address, serviceId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PostModel(connector, bindingRequest, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) UnbindFromApplicationInstance(bindingRequest models.InstanceBindingRequest, applicationId string) (int, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s/bindings", c.Address, applicationId))
	return brokerHttp.DeleteModelWithBody(connector, bindingRequest, http.StatusNoContent)
}

func (c *TapApiServiceApiOAuth2Connector) UnbindFromServiceInstance(bindingRequest models.InstanceBindingRequest, serviceId string) (int, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/services/%s/bindings", c.Address, serviceId))
	return brokerHttp.DeleteModelWithBody(connector, bindingRequest, http.StatusOK)
}

func (c *TapApiServiceApiOAuth2Connector) SendInvitation(email string) (userManagement.InvitationResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users/invitations", c.Address))
	body := userManagement.InvitationRequest{
		Email: email,
	}
	result := &userManagement.InvitationResponse{}
	_, err := brokerHttp.PostModel(connector, body, http.StatusCreated, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) ResendInvitation(email string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users/invitations/resend", c.Address))
	body := userManagement.InvitationRequest{
		Email: email,
	}
	_, err := brokerHttp.PostModel(connector, body, http.StatusCreated, "")
	return err
}

func (c *TapApiServiceApiOAuth2Connector) GetInvitations() ([]string, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users/invitations", c.Address))
	result := []string{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, &result)
	return result, err
}

func (c *TapApiServiceApiOAuth2Connector) GetUsers() ([]userManagement.UaaUser, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users", c.Address))
	result := []userManagement.UaaUser{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, &result)
	return result, err
}

func (c *TapApiServiceApiOAuth2Connector) DeleteInvitation(email string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users/invitations", c.Address))
	body := userManagement.InvitationRequest{
		Email: email,
	}
	_, err := brokerHttp.DeleteModelWithBody(connector, body, http.StatusNoContent)
	return err
}

func (c *TapApiServiceApiOAuth2Connector) DeleteUser(email string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users", c.Address))
	body := userManagement.InvitationRequest{
		Email: email,
	}
	_, err := brokerHttp.DeleteModelWithBody(connector, body, http.StatusNoContent)
	return err
}

func (c *TapApiServiceApiOAuth2Connector) ChangeCurrentUserPassword(password, newPassword string) error {
	body := userManagement.ChangePasswordRequest{
		CurrentPasswd: password,
		NewPasswd:     newPassword,
	}
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/users/current/password", c.Address))
	_, err := brokerHttp.PutModel(connector, body, http.StatusOK, "")
	return err
}
