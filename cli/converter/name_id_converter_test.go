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

package converter

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-api-service/models"
	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-cli/cli/test"
)

func getFakeServices() []models.Offering {
	result := []models.Offering{}
	result = append(result, test.NewFakeOffering(map[string]string{"name": "name_1", "offering_id": "offering_id_1", "plan_name": "plan_1", "plan_id": "plan_id_1"}))
	result = append(result, test.NewFakeOffering(map[string]string{"name": "name_2", "offering_id": "offering_id_2", "plan_name": "plan_2", "plan_id": "plan_id_2"}))
	result = append(result, test.NewFakeOffering(map[string]string{"name": "name_3", "offering_id": "offering_id_3", "plan_name": "plan_3", "plan_id": "plan_id_3"}))

	return result
}

func GetFakeServiceInstances() []models.ServiceInstance {
	result := []models.ServiceInstance{}
	result = append(result, models.ServiceInstance{Id: "1", Name: "instance1", Type: "SERVICE"})
	result = append(result, models.ServiceInstance{Id: "2", Name: "instance2", Type: "SERVICE"})
	result = append(result, models.ServiceInstance{Id: "3", Name: "instance3", Type: "SERVICE"})

	return result
}

func TestConvertFunction(t *testing.T) {
	Convey("Test convert method", t, func() {
		apiConfig, mockCtrl := test.SetApiAndLoginServiceMocks(t)
		fakeServices := getFakeServices()

		Convey("Should fail when GetOfferings return err", func() {
			fakeErr := errors.New("Error_msg")
			apiConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return([]models.Offering{}, fakeErr)

			_, _, err := FetchServiceAndPlanID(apiConfig, "service_name", "service_plan")

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, fakeErr)
		})
		Convey("Should fail when given plan doesn't exit", func() {
			apiConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return(fakeServices, nil)

			_, _, err := FetchServiceAndPlanID(apiConfig, "name_1", "wrong_plan_name")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find plan: 'wrong_plan_name' for service: 'name_1'")
		})
		Convey("Should fail when given service doesn't exist", func() {
			apiConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return(fakeServices, nil)

			_, _, err := FetchServiceAndPlanID(apiConfig, "wrong_label_name", "plan_1")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find service: 'wrong_label_name'")
		})
		Convey("Should pass when service guid and plan guid returned succesfully", func() {
			apiConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return(fakeServices, nil)

			serviceID, planID, err := FetchServiceAndPlanID(apiConfig, "name_3", "plan_3")

			So(err, ShouldBeNil)
			So(serviceID, ShouldEqual, "offering_id_3")
			So(planID, ShouldEqual, "plan_id_3")
		})
		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestGetServiceID(t *testing.T) {
	Convey("Test getServiceID", t, func() {
		apiConfig, mockCtrl := test.SetApiAndLoginServiceMocks(t)
		fakeServices := getFakeServices()

		Convey("Should fail when GetOfferings returns error", func() {
			fakeErr := errors.New("Error_msg")
			apiConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return([]models.Offering{}, fakeErr)

			_, err := GetOfferingID(apiConfig, "service_name")

			So(err, ShouldNotBeNil)
		})
		Convey("Should fail when given service doesn't exist", func() {
			apiConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return(fakeServices, nil)

			_, err := GetOfferingID(apiConfig, "wrong_label_name")

			So(err, ShouldNotBeNil)
		})
		Convey("Should pass when service guid returned succesfully", func() {
			apiConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return(fakeServices, nil)

			serviceID, err := GetOfferingID(apiConfig, "name_3")

			So(err, ShouldBeNil)
			So(serviceID, ShouldEqual, "offering_id_3")
		})
		Reset(func() {
			mockCtrl.Finish()
		})
	})
}
