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
	"bytes"
	"fmt"
	"mime/multipart"
)

import (
	"encoding/json"
	"io"
	"net/http"

	"errors"
	"io/ioutil"

	"github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	containerBrokerModels "github.com/trustedanalytics/tap-container-broker/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

func (c *TapApiServiceApiOAuth2Connector) CreateApplicationInstance(blob multipart.File, manifest models.Manifest) (catalogModels.Application, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications", c.Address))
	result := catalogModels.Application{}

	contentType, bodyBuf, err := c.prepareApplicationCreationForm(blob, manifest)
	if err != nil {
		logger.Error("ERROR: Preparing application creation form failed", err)
		return result, err
	}

	req, _ := http.NewRequest("POST", connector.Url, bodyBuf)
	req.Header.Add("Authorization", brokerHttp.GetOAuth2Header(connector.OAuth2))
	brokerHttp.SetContentType(req, contentType)

	logger.Infof("Doing: POST %v", connector.Url)
	resp, err := connector.Client.Do(req)
	if err != nil {
		logger.Error("ERROR: Make http request POST", err)
		return result, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("ERROR: Make http request POST", err)
		return result, err
	}

	if resp.StatusCode != http.StatusAccepted {
		return result, errors.New("Wrong response code! - data:" + string(data))
	}

	json.Unmarshal(data, &result)
	return result, nil
}

func (c *TapApiServiceApiOAuth2Connector) DeleteApplicationInstance(instanceId string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications/%s", c.Address, instanceId))
	_, err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return err
}

func (c *TapApiServiceApiOAuth2Connector) ListApplicationInstances() ([]models.ApplicationInstance, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications", c.Address))
	result := &[]models.ApplicationInstance{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) ScaleApplicationInstance(instanceId string, replication int) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications/%s/scale", c.Address, instanceId))
	body := models.ScaleApplicationRequest{
		Replicas: replication,
	}
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PutModel(connector, body, http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) StartApplicationInstance(instanceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications/%s/start", c.Address, instanceId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PutModel(connector, "", http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) StopApplicationInstance(instanceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications/%s/stop", c.Address, instanceId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PutModel(connector, "", http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) RestartApplicationInstance(instanceId string) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications/%s/restart", c.Address, instanceId))
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PutModel(connector, "", http.StatusOK, result)
	return *result, err
}

func (c *TapApiServiceApiOAuth2Connector) prepareApplicationCreationForm(blob multipart.File, manifest models.Manifest) (string, *bytes.Buffer, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := c.createBlobFormFile(blob, bodyWriter)
	if err != nil {
		return "", bodyBuf, err
	}
	c.createManifestFormFile(manifest, bodyWriter)
	if err != nil {
		return "", bodyBuf, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return contentType, bodyBuf, nil
}

func (c *TapApiServiceApiOAuth2Connector) createBlobFormFile(blob multipart.File, bodyWriter *multipart.Writer) error {
	blobWriter, err := bodyWriter.CreateFormFile("blob", "blob.tar.gz")
	if err != nil {
		logger.Error("Error creating blob file field")
		return err
	}
	size, err := io.Copy(blobWriter, blob)
	if err != nil {
		logger.Error("Error copying blob to buffer")
		return err
	}
	logger.Infof("Written %v bytes of blob to buffer", size)
	return nil
}

func (c *TapApiServiceApiOAuth2Connector) createManifestFormFile(manifest models.Manifest,
	bodyWriter *multipart.Writer) error {

	manifestWriter, err := bodyWriter.CreateFormFile("manifest", "manifest.json")
	if err != nil {
		logger.Error("Error creating manifest file field")
		return err
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		logger.Error("Error marshalling manifest.json")
		return err
	}
	size, err := manifestWriter.Write(manifestBytes)
	if err != nil {
		logger.Error("Error writing manifest to buffer")
		return err
	}
	logger.Infof("Written %v bytes of manifest to buffer", size)
	return nil
}

func (c *TapApiServiceApiOAuth2Connector) DeleteApplication(instanceId string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications/%s", c.Address, instanceId))
	_, err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return err
}

func (c *TapApiServiceApiOAuth2Connector) GetApplicationInstance(applicationId string) (models.ApplicationInstance, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v3/applications/%s", c.Address, applicationId))
	result := &models.ApplicationInstance{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}
