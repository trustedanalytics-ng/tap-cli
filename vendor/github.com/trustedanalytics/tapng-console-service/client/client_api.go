package client

import (
	"fmt"
	"net/http"

	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-console-service/models"
	containerBrokerModels "github.com/trustedanalytics/tapng-container-broker/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
	"github.com/trustedanalytics/tapng-go-common/logger"
	"mime/multipart"
)

var logger = logger_wrapper.InitLogger("client")

type TapConsoleServiceApi interface {
	GetCatalog() ([]models.Service, error)
	CreateOffer(serviceWithTemplate models.ServiceDeploy) (catalogModels.Service, error)
	CreateInstance(serviceId string, instance models.Instance) (containerBrokerModels.MessageResponse, error)
	ListApplications() ([]catalogModels.Application, error)
	CreateApplication(blob multipart.File, manifest models.Manifest) (catalogModels.Application, error)
	ListServicesInstances() ([]models.ServiceInstance, error)
	BindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error)
	UnbindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error)
}

func NewTapConsoleServiceApiWithBasicAuth(address, username, password string) (*TapConsoleServiceApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithBasicAuth()
	if err != nil {
		return nil, err
	}
	return &TapConsoleServiceApiConnector{address, username, password, client}, nil
}

func NewTapConsoleServiceApiWithSSLAndBasicAuth(address, username, password, certPemFile, keyPemFile, caPemFile string) (*TapConsoleServiceApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &TapConsoleServiceApiConnector{address, username, password, client}, nil
}

type TapConsoleServiceApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

func (c *TapConsoleServiceApiConnector) getApiConnector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		BasicAuth: &brokerHttp.BasicAuth{c.Username, c.Password},
		Client:    c.Client,
		Url:       url,
	}
}

func (c *TapConsoleServiceApiConnector) GetCatalog() ([]models.Service, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/catalog", c.Address))
	result := &[]models.Service{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapConsoleServiceApiConnector) CreateOffer(serviceWithTemplate models.ServiceDeploy) (catalogModels.Service, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/offering", c.Address))
	result := &catalogModels.Service{}
	err := brokerHttp.AddModel(connector, serviceWithTemplate, http.StatusCreated, result)
	return *result, err
}

func (c *TapConsoleServiceApiConnector) CreateInstance(serviceId string, instance models.Instance) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/instances/%s", c.Address, serviceId))
	result := &containerBrokerModels.MessageResponse{}
	err := brokerHttp.AddModel(connector, instance, http.StatusCreated, result)
	return *result, err
}

func (c *TapConsoleServiceApiConnector) ListServicesInstances() ([]models.ServiceInstance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/instances", c.Address))
	result := &[]models.ServiceInstance{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapConsoleServiceApiConnector) BindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/bind/%s/%s", c.Address, srcInstanceId, dstInstanceId))
	result := &containerBrokerModels.MessageResponse{}
	err := brokerHttp.AddModel(connector, "", http.StatusOK, result)
	return *result, err
}

func (c *TapConsoleServiceApiConnector) UnbindInstance(srcInstanceId, dstInstanceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/unbind/%s/%s", c.Address, srcInstanceId, dstInstanceId))
	result := &containerBrokerModels.MessageResponse{}
	err := brokerHttp.AddModel(connector, "", http.StatusOK, result)
	return *result, err
}
