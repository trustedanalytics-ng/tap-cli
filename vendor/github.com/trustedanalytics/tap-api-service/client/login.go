package client

import (
	"fmt"
	"net/http"

	"errors"

	uaa "github.com/trustedanalytics/tap-api-service/uaa-connector"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

type TapConsoleServiceLoginApi interface {
	Login() (uaa.LoginResponse, int, error)
	GetConsoleServiceHealth() error
}

type TapConsoleServiceApiBasicAuthConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

func NewTapConsoleServiceLoginApiWithBasicAuth(address, username, password string) (TapConsoleServiceLoginApi, error) {
	client, _, err := brokerHttp.GetHttpClient()
	if err != nil {
		return nil, err
	}
	return &TapConsoleServiceApiBasicAuthConnector{address, username, password, client}, nil
}

func NewTapConsoleServiceLoginApiWithSSLAndBasicAuth(address, username, password, certPemFile, keyPemFile, caPemFile string) (TapConsoleServiceLoginApi, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &TapConsoleServiceApiBasicAuthConnector{address, username, password, client}, nil
}

func (c *TapConsoleServiceApiBasicAuthConnector) getApiBasicAuthConnector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		BasicAuth: &brokerHttp.BasicAuth{c.Username, c.Password},
		Client:    c.Client,
		Url:       url,
	}
}

func (c *TapConsoleServiceApiBasicAuthConnector) Login() (uaa.LoginResponse, int, error) {
	connector := c.getApiBasicAuthConnector(fmt.Sprintf("%s/api/v1/login", c.Address))
	result := &uaa.LoginResponse{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapConsoleServiceApiBasicAuthConnector) GetConsoleServiceHealth() error {
	connector := c.getApiBasicAuthConnector(fmt.Sprintf("%s/api/v1/healthz", c.Address))
	status, _, err := brokerHttp.RestGET(connector.Url, brokerHttp.GetBasicAuthHeader(connector.BasicAuth), connector.Client)
	if status != http.StatusOK {
		err = errors.New("Invalid health status: " + string(status))
	}
	return err
}
