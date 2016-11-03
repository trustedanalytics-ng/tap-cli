/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package client

import (
	"fmt"
	"net/http"

	"errors"

	uaa "github.com/trustedanalytics/tap-api-service/uaa-connector"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

type TapApiServiceLoginApi interface {
	Login() (uaa.LoginResponse, int, error)
	GetApiServiceHealth() error
	GetLoginCredentials() (Address, Username, Password string)
}

type TapApiServiceApiBasicAuthConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

func NewTapApiServiceLoginApiWithBasicAuth(address, username, password string) (TapApiServiceLoginApi, error) {
	client, _, err := brokerHttp.GetHttpClient()
	if err != nil {
		return nil, err
	}
	return &TapApiServiceApiBasicAuthConnector{address, username, password, client}, nil
}

func NewTapApiServiceLoginApiWithSSLAndBasicAuth(address, username, password, certPemFile, keyPemFile, caPemFile string) (TapApiServiceLoginApi, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &TapApiServiceApiBasicAuthConnector{address, username, password, client}, nil
}

func (c *TapApiServiceApiBasicAuthConnector) getApiBasicAuthConnector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		BasicAuth: &brokerHttp.BasicAuth{c.Username, c.Password},
		Client:    c.Client,
		Url:       url,
	}
}

func (c *TapApiServiceApiBasicAuthConnector) GetLoginCredentials() (Address, Username, Password string) {
	return c.Address, c.Username, c.Password
}

func (c *TapApiServiceApiBasicAuthConnector) Login() (uaa.LoginResponse, int, error) {
	connector := c.getApiBasicAuthConnector(fmt.Sprintf("%s/api/v2/login", c.Address))
	result := &uaa.LoginResponse{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapApiServiceApiBasicAuthConnector) GetApiServiceHealth() error {
	connector := c.getApiBasicAuthConnector(fmt.Sprintf("%s/api/v2/healthz", c.Address))
	status, _, err := brokerHttp.RestGET(connector.Url, brokerHttp.GetBasicAuthHeader(connector.BasicAuth), connector.Client)
	if status != http.StatusOK {
		err = errors.New("Invalid health status: " + string(status))
	}
	return err
}
