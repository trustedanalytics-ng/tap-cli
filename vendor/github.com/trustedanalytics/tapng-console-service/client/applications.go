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

	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
	templateRepositoryModels "github.com/trustedanalytics/tapng-template-repository/model"
	"io/ioutil"
)

func (c *TapConsoleServiceApiConnector) ListApplications() ([]catalogModels.Application, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/applications", c.Address))
	result := &[]catalogModels.Application{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapConsoleServiceApiConnector) CreateApplication(blob multipart.File, image catalogModels.Image,
	template templateRepositoryModels.Template) (catalogModels.Application, error) {

	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/applications", c.Address))
	result := catalogModels.Application{}

	contentType, bodyBuf, err := c.prepareApplicationCreationForm(blob, image, template)
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
	json.Unmarshal(data, &result)

	return result, nil
}

func (c *TapConsoleServiceApiConnector) prepareApplicationCreationForm(blob multipart.File, image catalogModels.Image,
	template templateRepositoryModels.Template) (string, *bytes.Buffer, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := c.createBlobFormFile(blob, bodyWriter)
	if err != nil {
		return "", bodyBuf, err
	}
	c.createImageMetaFormFile(image, bodyWriter)
	if err != nil {
		return "", bodyBuf, err
	}
	c.createTemplateMetaFormFile(template, bodyWriter)
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

func (c *TapConsoleServiceApiConnector) createImageMetaFormFile(image catalogModels.Image,
	bodyWriter *multipart.Writer) error {

	imageMetaWriter, err := bodyWriter.CreateFormFile("image", "image.json")
	if err != nil {
		logger.Error("Error creating image file field")
		return err
	}
	imageBytes, err := json.Marshal(image)
	if err != nil {
		logger.Error("Error marshalling image.json")
		return err
	}
	size, err := imageMetaWriter.Write(imageBytes)
	if err != nil {
		logger.Error("Error writing image metadata to buffer")
		return err
	}
	logger.Infof("Written %v bytes of image to buffer", size)
	return nil
}

func (c *TapConsoleServiceApiConnector) createTemplateMetaFormFile(template templateRepositoryModels.Template,
	bodyWriter *multipart.Writer) error {

	templateWriter, err := bodyWriter.CreateFormFile("template", "template.json")
	if err != nil {
		logger.Error("Error creating image file field")
		return err
	}
	templateBytes, err := json.Marshal(template)
	if err != nil {
		logger.Error("Error marshalling template.json")
		return err
	}
	size, err := templateWriter.Write(templateBytes)
	if err != nil {
		logger.Error("Error writing template metadata to buffer")
		return err
	}
	logger.Infof("Written %v bytes of template to buffer", size)
	return nil
}
