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

package actions

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/assertions"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics-ng/tap-api-service/models"
	"github.com/trustedanalytics-ng/tap-cli/api"
)

const sampleOfferingName = "sampleOffering"
const sampleManifestFilename = "testPurposeManifest.json"

var fakeErr error = errors.New("some fake error")

func TestCreateOffering(t *testing.T) {
	Convey("Test offering creation", t, func() {
		actionsConfig, mockCtrl := setupActionsTest(t)

		Convey("shall return error when manifest file does not exist", func() {
			err := actionsConfig.CreateOffering("someUnexistingFile.json")

			So(err.Error(), assertions.ShouldContainSubstring, "no such file or directory")
		})

		Convey("shall return error when mistaken manifest provided", func() {
			manifest := `{
			  "broker_name" : "test-broker",
			  "services" : [],
			  "template" : "UNEXPECTED STRING"
			}`
			ioutil.WriteFile(sampleManifestFilename, []byte(manifest), 0644)

			err := actionsConfig.CreateOffering(sampleManifestFilename)

			So(err.Error(), assertions.ShouldContainSubstring, "cannot unmarshal string")
		})

		Convey("shall proceed correctly with correct manifest", func() {
			manifest := `{
			  "broker_name" : "test-broker",
			  "services" : [{"Plans":[]}],
			  "template" : {}
			}`
			ioutil.WriteFile(sampleManifestFilename, []byte(manifest), 0644)
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().CreateOffer(gomock.Any()).Return(nil, nil)

			err := actionsConfig.CreateOffering(sampleManifestFilename)

			So(err, assertions.ShouldBeNil)
		})

		Reset(func() {
			os.Remove(sampleManifestFilename)
			mockCtrl.Finish()
		})
	})
}

func TestGetOffering(t *testing.T) {
	Convey("Get offering info", t, func() {
		actionsConfig, mockCtrl := setupActionsTest(t)

		Convey("shall return error when offering does not exist", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().GetOfferings().Return([]models.Offering{}, nil)

			err := actionsConfig.GetOffering("offeringThatDoesNotExist")

			So(err.Error(), assertions.ShouldContainSubstring, "Could not find offering with such name")
		})

		Convey("shall succeed when offering exist", func() {
			fakeOfferings := []models.Offering{}
			fakeOfferings = append(fakeOfferings, models.Offering{Name: sampleOfferingName})
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().GetOfferings().Return(fakeOfferings, nil)

			err := actionsConfig.GetOffering(sampleOfferingName)

			So(err, assertions.ShouldBeNil)
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestListOfferings(t *testing.T) {
	Convey("Get offerings list", t, func() {
		actionsConfig, mockCtrl := setupActionsTest(t)

		Convey("shall return error when API call failed", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().GetOfferings().Return(nil, fakeErr)

			err := actionsConfig.ListOfferings()

			So(err, assertions.ShouldEqual, fakeErr)
		})

		Convey("shall not report errror when succeeded", func() {
			sampleOfferings := []models.Offering{}
			sampleOfferings = append(sampleOfferings, models.Offering{Name: sampleOfferingName})
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().GetOfferings().Return(sampleOfferings, nil)

			err := actionsConfig.ListOfferings()

			So(err, assertions.ShouldBeNil)
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestDeleteOffering(t *testing.T) {
	Convey("Delete offering", t, func() {
		actionsConfig, mockCtrl := setupActionsTest(t)

		Convey("shall return error when fetching service ID impossible", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().GetOfferings().Return([]models.Offering{}, nil)

			err := actionsConfig.DeleteOffering("notExistingOffering")

			So(err.Error(), assertions.ShouldContainSubstring, "Cannot fetch service id")
		})

		Convey("shall return error when offering removal failed", func() {
			sampleOfferings := []models.Offering{}
			sampleOfferings = append(sampleOfferings, models.Offering{Name: sampleOfferingName})
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().GetOfferings().Return(sampleOfferings, nil)
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().DeleteOffering(gomock.Any()).Return(fakeErr)

			err := actionsConfig.DeleteOffering(sampleOfferingName)

			So(err.Error(), assertions.ShouldContainSubstring, "Cannot delete offering")
		})

		Convey("shall return nothing when offering removed correctly", func() {
			sampleOfferings := []models.Offering{}
			sampleOfferings = append(sampleOfferings, models.Offering{Name: sampleOfferingName})
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().GetOfferings().Return(sampleOfferings, nil)
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).EXPECT().DeleteOffering(gomock.Any()).Return(nil)

			err := actionsConfig.DeleteOffering(sampleOfferingName)

			So(err, assertions.ShouldBeNil)
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}
