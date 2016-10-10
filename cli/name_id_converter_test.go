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
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
)

func getFakeServices() []models.Service {
	result := []models.Service{}
	result = append(result, newFakeService(map[string]string{"label": "label_1", "service_id": "service_guid_1", "plan_name": "plan_1", "plan_id": "plan_guid_1"}))
	result = append(result, newFakeService(map[string]string{"label": "label_2", "service_id": "service_guid_2", "plan_name": "plan_2", "plan_id": "plan_guid_2"}))
	result = append(result, newFakeService(map[string]string{"label": "label_3", "service_id": "service_guid_3", "plan_name": "plan_3", "plan_id": "plan_guid_3"}))

	return result
}

func getFakeInstances() []models.ServiceInstance {
	result := []models.ServiceInstance{}
	result = append(result, models.ServiceInstance{Instance: catalogModels.Instance{Id: "1", Name: "instance1", Type: "SERVICE"}})
	result = append(result, models.ServiceInstance{Instance: catalogModels.Instance{Id: "2", Name: "instance2", Type: "SERVICE"}})
	result = append(result, models.ServiceInstance{Instance: catalogModels.Instance{Id: "3", Name: "instance3", Type: "SERVICE"}})

	return result
}

func TestConvertFunction(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)

	fakeServices := getFakeServices()

	Convey("Test convert method", t, func() {
		Convey("Should fail when GetCatalog return err", func() {
			fakeErr := errors.New("Error_msg")
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetCatalog().
				Return([]models.Service{}, fakeErr)

			_, _, err := convertServiceAndPlanNameToId(actionsConfig, "service_name", "service_plan")

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, fakeErr)
		})
		Convey("Should fail when given plan doesn't exit", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetCatalog().
				Return(fakeServices, nil)

			_, _, err := convertServiceAndPlanNameToId(actionsConfig, "label_1", "wrong_plan_name")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find plan: 'wrong_plan_name' for service: 'label_1'")
		})
		Convey("Should fail when given service doesn't exist", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetCatalog().
				Return(fakeServices, nil)

			_, _, err := convertServiceAndPlanNameToId(actionsConfig, "wrong_label_name", "plan_1")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find service: 'wrong_label_name'")
		})
		Convey("Should pass when service guid and plan guid returned succesfully", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetCatalog().
				Return(fakeServices, nil)

			serviceID, planID, err := convertServiceAndPlanNameToId(actionsConfig, "label_3", "plan_3")

			So(err, ShouldBeNil)
			So(serviceID, ShouldEqual, "service_guid_3")
			So(planID, ShouldEqual, "plan_guid_3")
		})
	})
}

func TestGetServiceID(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)

	fakeServices := getFakeServices()

	Convey("Test getServiceID", t, func() {
		Convey("Should fail when GetCatalog returns error", func() {
			fakeErr := errors.New("Error_msg")
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetCatalog().
				Return([]models.Service{}, fakeErr)

			_, err := getServiceID(actionsConfig, "service_name")

			So(err, ShouldNotBeNil)
		})
		Convey("Should fail when given service doesn't exist", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetCatalog().
				Return(fakeServices, nil)

			_, err := getServiceID(actionsConfig, "wrong_label_name")

			So(err, ShouldNotBeNil)
		})
		Convey("Should pass when service guid returned succesfully", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetCatalog().
				Return(fakeServices, nil)

			serviceID, err := getServiceID(actionsConfig, "label_3")

			So(err, ShouldBeNil)
			So(serviceID, ShouldEqual, "service_guid_3")
		})
	})
}

func TestConvertBindingsList(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)

	fakeInstances := getFakeInstances()

	Convey("When ListServiceInstances returns error", t, func() {
		fakeErr := errors.New("Error_msg")
		actionsConfig.ApiService.(*api.MockTapApiServiceApi).
			EXPECT().
			ListServiceInstances().
			Return(nil, fakeErr)
		sampleBindingList := []string{"instance1"}

		Convey("convertBindingList should return error", func() {
			err := convertBindingsList(actionsConfig, sampleBindingList)
			So(err, ShouldNotBeNil)
		})
	})

	Convey(fmt.Sprintf("When ListServiceInstance returns %v", fakeInstances), t, func() {
		actionsConfig.ApiService.(*api.MockTapApiServiceApi).
			EXPECT().
			ListServiceInstances().
			Return(fakeInstances, nil)

		testCases := []struct {
			bindingList   []string
			isError       bool
			convertedList []string
		}{
			{[]string{"instance2", "XXXX"}, true, []string{}},
			{[]string{"instance2", "instance4"}, true, []string{}},
			{[]string{"", "instance4"}, true, []string{}},
			{[]string{"instance1", "instance3"}, false, []string{"1", "3"}},
			{[]string{"instance1", "instance3", "instance2"}, false, []string{"1", "3", "2"}},
			{[]string{"instance2"}, false, []string{"2"}},
			{[]string{}, false, []string{}},
		}

		for _, tc := range testCases {
			Convey(fmt.Sprintf("convertBindingList should return proper response for %v", tc.bindingList), func() {
				err := convertBindingsList(actionsConfig, tc.bindingList)

				if tc.isError {
					Convey("error should not be nil", func() {
						So(err, ShouldNotBeNil)
					})
				} else {
					Convey("error should be nil", func() {
						So(err, ShouldBeNil)
					})
					Convey("bindingList should be properly converted", func() {
						So(tc.bindingList, ShouldResemble, tc.convertedList)
					})
				}
			})
		}
	})
}
