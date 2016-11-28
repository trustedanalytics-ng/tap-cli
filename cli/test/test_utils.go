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

package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
)

func SetApiAndLoginServiceMocks(t *testing.T) api.Config {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	apiServiceMock := api.NewMockTapApiServiceApi(mockCtrl)
	apiServiceLoginMock := api.NewMockTapApiServiceLoginApi(mockCtrl)
	return api.Config{ApiService: apiServiceMock, ApiServiceLogin: apiServiceLoginMock}
}

func FillCredentialsFile(content string) {
	ioutil.WriteFile(api.CredsPath, []byte(content), api.PERMISSIONS)
}

func NewFakeOffering(m map[string]string) models.Service {
	return models.Service{
		Entity: models.ServiceEntity{
			Label: m["label"],
			ServicePlans: []models.ServicePlan{
				{Entity: models.ServicePlanEntity{Name: m["name"]}, Metadata: models.Metadata{}}},
			Description: m["desc"],
			State:       m["state"],
		},
		Metadata: models.Metadata{Guid: "RANDOM_GUID"}}
}

func NewFakeAppInstance(m map[string]string) models.ApplicationInstance {
	createdOn, _ := strconv.Atoi(m["ob"])
	updatedOn, _ := strconv.Atoi(m["ub"])
	rep, _ := strconv.Atoi(m["replication"])
	appInstance := models.ApplicationInstance{
		Instance: catalogModels.Instance{
			Name: m["name"],
			AuditTrail: catalogModels.AuditTrail{
				CreatedOn:     int64(createdOn),
				CreatedBy:     m["cb"],
				LastUpdatedOn: int64(updatedOn),
				LastUpdateBy:  m["ub"]},
			State: catalogModels.InstanceState(m["instance_state"]),
		},
		Replication:      rep,
		ImageState:       catalogModels.ImageState(m["image_state"]),
		Urls:             []string{m["urls"]},
		ImageType:        catalogModels.ImageType("fakeType"),
		Memory:           m["memory"],
		DiskQuota:        m["quota"],
		RunningInstances: 0}

	if value, exist := m[catalogModels.LAST_STATE_CHANGE_REASON]; exist {
		appInstance.Metadata = []catalogModels.Metadata{
			{Id: catalogModels.LAST_STATE_CHANGE_REASON, Value: value},
		}
	}
	return appInstance
}

func NewFakeService(m map[string]string) models.Service {
	return models.Service{
		Entity: models.ServiceEntity{
			Label:    m["label"],
			UniqueId: m["service_id"],
			ServicePlans: []models.ServicePlan{
				{
					Entity: models.ServicePlanEntity{
						Name:     m["plan_name"],
						UniqueId: m["plan_id"],
					},
					Metadata: models.Metadata{},
				},
			},
		},
		Metadata: models.Metadata{Guid: m["service_id"]},
	}
}

func CaptureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
