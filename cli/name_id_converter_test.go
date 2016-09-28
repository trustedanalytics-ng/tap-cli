package cli

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-api-service/models"
	"github.com/trustedanalytics/tap-cli/api"
)

func getFakeServices() []models.Service {
	result := []models.Service{}
	result = append(result, newFakeService(map[string]string{"label": "label_1", "service_id": "service_guid_1", "plan_name": "plan_1", "plan_id": "plan_guid_1"}))
	result = append(result, newFakeService(map[string]string{"label": "label_2", "service_id": "service_guid_2", "plan_name": "plan_2", "plan_id": "plan_guid_2"}))
	result = append(result, newFakeService(map[string]string{"label": "label_3", "service_id": "service_guid_3", "plan_name": "plan_3", "plan_id": "plan_guid_3"}))

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
