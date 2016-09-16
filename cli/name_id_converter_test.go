package cli

import (
	"testing"
	"errors"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-api-service/models"
)

func TestConvertFunction(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)

	fakeService1 := newFakeService(map[string]string{"label": "label_1","service_id":"service_guid_1","plan_name":"plan_1","plan_id":"plan_guid_1"})
	fakeService2 := newFakeService(map[string]string{"label": "label_2","service_id":"service_guid_2","plan_name":"plan_2","plan_id":"plan_guid_2"})
	fakeService3 := newFakeService(map[string]string{"label": "label_3","service_id":"service_guid_3","plan_name":"plan_3","plan_id":"plan_guid_3"})


	Convey("Test convert method", t, func() {
		Convey("Should fail when GetCatalog return err", func() {
			fakeErr:=errors.New("Error_msg")
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				GetCatalog().
				Return([]models.Service{}, fakeErr)

				_,_,err := convert(actionsConfig, "service_name", "service_plan")

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, fakeErr)
		})
		Convey("Should fail when given plan dosen't exits", func() {
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				GetCatalog().
				Return([]models.Service{fakeService1,fakeService2,fakeService3}, nil)

			_,_,err := convert(actionsConfig, "label_1", "wrong_plan_name")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find plan: 'wrong_plan_name' for service: 'label_1'")
		})
		Convey("Should fail when given service dosen't exits", func() {
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				GetCatalog().
				Return([]models.Service{fakeService1,fakeService2,fakeService3}, nil)

			_,_,err := convert(actionsConfig, "wrong_label_name", "plan_1")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find service: 'wrong_label_name'")
		})
		Convey("Should pass when service guid and plan guid returned succesfully", func() {
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				GetCatalog().
				Return([]models.Service{fakeService1,fakeService2,fakeService3}, nil)

			serviceID,planID,err := convert(actionsConfig, "label_3", "plan_3")

			So(err, ShouldBeNil)
			So(serviceID, ShouldEqual, "service_guid_3")
			So(planID, ShouldEqual, "plan_guid_3")
		})
	})
}

