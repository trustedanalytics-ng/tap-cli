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
	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-console-service/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
	"io/ioutil"
)

func (c *TapConsoleServiceApiConnector) ListApplications() ([]catalogModels.Application, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/applications", c.Address))
	result := &[]catalogModels.Application{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapConsoleServiceApiConnector) CreateApplication(blob multipart.File, manifest models.Manifest) (catalogModels.Instance, error) {

	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/applications", c.Address))
	result := catalogModels.Instance{}

	contentType, bodyBuf, err := c.prepareApplicationCreationForm(blob, manifest)
	if err != nil {
		logger.Error("ERROR: Preparing application creation form failed", err)
		return result, err
	}

	req, _ := http.NewRequest("POST", connector.Url, bodyBuf)
	brokerHttp.AddBasicAuth(req, connector.BasicAuth)
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

	if resp.StatusCode != http.StatusCreated {
		return result, errors.New(string(data))
	}

	json.Unmarshal(data, &result)
	return result, nil
}

func (c *TapConsoleServiceApiConnector) prepareApplicationCreationForm(blob multipart.File, manifest models.Manifest) (string, *bytes.Buffer, error) {
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

func (c *TapConsoleServiceApiConnector) createBlobFormFile(blob multipart.File, bodyWriter *multipart.Writer) error {
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

func (c *TapConsoleServiceApiConnector) createManifestFormFile(manifest models.Manifest,
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

func (c *TapConsoleServiceApiConnector) DeleteApplication(instanceId string) error {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/applications/%s", c.Address, instanceId))
	err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return err
}
