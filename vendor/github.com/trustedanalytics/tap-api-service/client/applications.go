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

func (c *TapConsoleServiceApiOAuth2Connector) CreateApplicationInstance(blob multipart.File, manifest models.Manifest) (catalogModels.Application, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications", c.Address))
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

func (c *TapConsoleServiceApiOAuth2Connector) DeleteApplicationInstance(instanceId string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s", c.Address, instanceId))
	_, err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return err
}

func (c *TapConsoleServiceApiOAuth2Connector) ListApplicationInstances() ([]models.ApplicationInstance, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications", c.Address))
	result := &[]models.ApplicationInstance{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapConsoleServiceApiOAuth2Connector) ScaleApplicationInstance(instanceId string, replication int) (containerBrokerModels.MessageResponse, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s/scale", c.Address, instanceId))
	body := containerBrokerModels.ScaleInstanceRequest{
		Replicas: replication,
	}
	result := &containerBrokerModels.MessageResponse{}
	_, err := brokerHttp.PutModel(connector, body, http.StatusOK, result)
	return *result, err
}

func (c *TapConsoleServiceApiOAuth2Connector) prepareApplicationCreationForm(blob multipart.File, manifest models.Manifest) (string, *bytes.Buffer, error) {
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

func (c *TapConsoleServiceApiOAuth2Connector) createBlobFormFile(blob multipart.File, bodyWriter *multipart.Writer) error {
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

func (c *TapConsoleServiceApiOAuth2Connector) createManifestFormFile(manifest models.Manifest,
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

func (c *TapConsoleServiceApiOAuth2Connector) DeleteApplication(instanceId string) error {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s", c.Address, instanceId))
	_, err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return err
}

func (c *TapConsoleServiceApiOAuth2Connector) GetApplicationInstance(instanceId string) (models.ApplicationInstance, error) {
	connector := c.getApiOAuth2Connector(fmt.Sprintf("%s/api/v1/applications/%s", c.Address, instanceId))
	result := &models.ApplicationInstance{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}
