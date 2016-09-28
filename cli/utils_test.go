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

package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"

	"github.com/golang/mock/gomock"
	"github.com/trustedanalytics/tap-cli/api"
	"io"
	"io/ioutil"
	"strconv"
)

func setApiAndLoginServiceMocks(t *testing.T) *ActionsConfig {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	a := api.NewMockTapApiServiceApi(mockCtrl)
	b := api.NewMockTapApiServiceLoginApi(mockCtrl)
	return &ActionsConfig{api.Config{a, b}}
}

func fillCredentialsFile(content string) {
	ioutil.WriteFile(api.CredsPath, []byte(content), api.PERMISSIONS)
}

func newFakeOffering(m map[string]string) models.Service {
	return models.Service{
		models.ServiceEntity{
			Label:        m["label"],
			ServicePlans: []models.ServicePlan{{models.ServicePlanEntity{Name: m["name"]}, models.Metadata{}}},
			Description:  m["desc"],
			State:        m["state"],
		},
		models.Metadata{Guid: "RANDOM_GUID"}}
}

func newFakeAppInstance(m map[string]string) models.ApplicationInstance {
	createdOn, _ := strconv.Atoi(m["ob"])
	updatedOn, _ := strconv.Atoi(m["ub"])
	rep, _ := strconv.Atoi(m["replication"])
	return models.ApplicationInstance{
		catalogModels.Instance{
			Name:       m["name"],
			AuditTrail: catalogModels.AuditTrail{int64(createdOn), m["cb"], int64(updatedOn), m["ub"]},
			State:      catalogModels.InstanceState(m["instance_state"]),
		},
		rep,
		catalogModels.ImageState(m["image_state"]),
		[]string{m["urls"]},
		catalogModels.ImageType("fakeType"),
		m["memory"],
		m["quota"],
		0}
}

//return service.Entity.UniqueId, plan.Entity.UniqueId, nil
func newFakeService(m map[string]string) models.Service {
	return models.Service{
		models.ServiceEntity{
			Label:    m["label"],
			UniqueId: m["service_id"],
			ServicePlans: []models.ServicePlan{
				{
					models.ServicePlanEntity{
						Name:     m["plan_name"],
						UniqueId: m["plan_id"],
					},
					models.Metadata{},
				},
			},
		},
		models.Metadata{m["service_id"]},
	}
}

func captureStdout(f func()) string {
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
