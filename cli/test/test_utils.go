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

func NewFakeAppInstance(m map[string]string) models.ApplicationInstance {
	createdOn, _ := strconv.Atoi(m["ob"])
	updatedOn, _ := strconv.Atoi(m["ub"])
	rep, _ := strconv.Atoi(m["replication"])
	appInstance := models.ApplicationInstance{
		Id:   m["id"],
		Name: m["name"],
		AuditTrail: catalogModels.AuditTrail{
			CreatedOn:     int64(createdOn),
			CreatedBy:     m["cb"],
			LastUpdatedOn: int64(updatedOn),
			LastUpdateBy:  m["ub"]},
		State:            catalogModels.InstanceState(m["instance_state"]),
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

func NewFakeOffering(m map[string]string) models.Offering {
	return models.Offering{
		Name:        m["name"],
		Id:          m["offering_id"],
		Description: m["desc"],
		State:       m["state"],
		OfferingPlans: []models.OfferingPlan{
			{
				Name: m["plan_name"],
				Id:   m["plan_id"],
			},
		},
		Metadata: []catalogModels.Metadata{},
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
