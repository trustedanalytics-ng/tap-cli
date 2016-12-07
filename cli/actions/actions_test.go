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
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-api-service/models"
	"github.com/trustedanalytics/tap-api-service/uaa-connector"
	"github.com/trustedanalytics/tap-api-service/user-management-connector"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-cli/cli/printer"
	"github.com/trustedanalytics/tap-cli/cli/test"
)

var url string = "fake_url"
var login string = "fake_admin"
var pass string = "fake_password"

var expectedUaaRes = uaa_connector.LoginResponse{
	AccessToken:  "fake_access_token",
	RefreshToken: "fake_refresh_token",
	TokenType:    "fake_token_type",
	ExpiresIn:    10,
	Scope:        "fake_scope",
	Jti:          "fake_jti",
}

var expectedCredsFileContent = "{" +
	"\"address\":\"" + url + "\"," +
	"\"username\":\"" + login + "\"," +
	"\"token\":\"" + expectedUaaRes.AccessToken + "\"," +
	"\"type\":\"" + expectedUaaRes.TokenType + "\"," +
	"\"expires\":" + strconv.Itoa(expectedUaaRes.ExpiresIn) +
	"}"

func prepareLoginMock(c ActionsConfig, res uaa_connector.LoginResponse, status int, err error) {
	c.ApiServiceLogin.(*api.MockTapApiServiceLoginApi).
		EXPECT().
		Login().
		Return(res, status, err)
}

func prepareIntroduceMock(c ActionsConfig, err error) {
	c.ApiServiceLogin.(*api.MockTapApiServiceLoginApi).
		EXPECT().
		Introduce().
		Return(err)
}

func TestLoginActions(t *testing.T) {
	actionsConfig := ActionsConfig{test.SetApiAndLoginServiceMocks(t)}
	Convey("Test Login command", t, func() {

		actionsConfig.ApiServiceLogin.(*api.MockTapApiServiceLoginApi).
			EXPECT().
			GetLoginCredentials().
			Return(url, login, pass)

		Convey("Should fail when user unauthorized", func() {
			someErr := errors.New("Authentication failed")
			prepareIntroduceMock(actionsConfig, nil)
			prepareLoginMock(actionsConfig, expectedUaaRes, http.StatusUnauthorized, someErr)

			err := actionsConfig.Login()

			So(err.Error(), ShouldContainSubstring, someErr.Error())
		})
		Convey("Should fail when connecting error occurs", func() {
			someErr := errors.New("server error")
			prepareIntroduceMock(actionsConfig, nil)
			prepareLoginMock(actionsConfig, expectedUaaRes, http.StatusInternalServerError, someErr)

			err := actionsConfig.Login()

			So(err.Error(), ShouldContainSubstring, someErr.Error())
		})
		Convey("Should fail when incompatibility detected", func() {
			prepareIntroduceMock(actionsConfig, nil)
			prepareLoginMock(actionsConfig, expectedUaaRes, http.StatusNotFound, nil)

			err := actionsConfig.Login()

			So(err.Error(), ShouldContainSubstring, "incompatibility detected")
		})
		Convey("Should fail when we do not talk with TAP", func() {
			someErr := errors.New("anything")
			prepareIntroduceMock(actionsConfig, someErr)

			err := actionsConfig.Login()

			So(err, ShouldEqual, someErr)
		})
		Convey("Should pass when credentials succesfully saved", func() {
			prepareIntroduceMock(actionsConfig, nil)
			prepareLoginMock(actionsConfig, expectedUaaRes, http.StatusOK, nil)

			stdout := test.CaptureStdout(func() {
				actionsConfig.Login()
			})

			b, err := ioutil.ReadFile(api.CredsPath)
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, string(expectedCredsFileContent))
			So(stdout, ShouldContainSubstring, "Authentication succeeded")
		})
	})
}

func TestSendInvitationCommand(t *testing.T) {
	actionsConfig := ActionsConfig{test.SetApiAndLoginServiceMocks(t)}
	Convey("Test Login command", t, func() {
		Convey("Should fail when inviting error occurs", func() {
			errorMsg := "cannot invite"
			test.FillCredentialsFile(expectedCredsFileContent)

			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				SendInvitation(login).
				Return(user_management_connector.InvitationResponse{}, errors.New(errorMsg))

			err := actionsConfig.SendInvitation(login)

			So(err, ShouldNotBeNil)
		})

		Convey("Should pass when user invited", func() {
			test.FillCredentialsFile(expectedCredsFileContent)

			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				SendInvitation(login).
				Return(user_management_connector.InvitationResponse{}, nil)

			stdout := test.CaptureStdout(func() {
				actionsConfig.SendInvitation(login)
			})

			So(stdout, ShouldContainSubstring, fmt.Sprintf("User %q successfully invited", login))
		})
	})
}

func TestDeleteUserCommand(t *testing.T) {
	actionsConfig := ActionsConfig{test.SetApiAndLoginServiceMocks(t)}
	test.FillCredentialsFile(expectedCredsFileContent)

	Convey("Test delete-user command", t, func() {
		errorMsg := "error message"

		Convey("Should fail when no user exists", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				DeleteUser(login).
				Return(errors.New(errorMsg))

			err := actionsConfig.DeleteUser(login)

			So(err, ShouldNotBeNil)
		})

		Convey("Should pass when user exists", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				DeleteUser(login).
				Return(nil)

			stdout := test.CaptureStdout(func() {
				actionsConfig.DeleteUser(login)
			})

			So(stdout, ShouldContainSubstring, fmt.Sprintf("User %q successfully removed\n", login))
		})
	})
}

func TestCatalogCommand(t *testing.T) {
	actionsConfig := ActionsConfig{test.SetApiAndLoginServiceMocks(t)}
	test.FillCredentialsFile(expectedCredsFileContent)

	fakeOff1 := test.NewFakeOffering(map[string]string{"name": "OFFERING_1", "offering_id": "offering_id_1", "plan_name": "PLAN_1", "plan_id": "plan_id_1", "desc": "DESC_1", "state": "READY"})
	fakeOff2 := test.NewFakeOffering(map[string]string{"name": "OFFERING_2", "offering_id": "offering_id_2", "plan_name": "PLAN_2", "plan_id": "plan_id_2", "desc": "DESC_2", "state": "READY"})
	fakeOff3 := test.NewFakeOffering(map[string]string{"name": "OFFERING_3", "offering_id": "offering_id_3", "plan_name": "PLAN_3", "plan_id": "plan_id_3", "desc": "DESC_3", "state": "READY"})

	Convey("Test catalog command", t, func() {
		Convey("Should pretty print offerings list", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				GetOfferings().
				Return([]models.Offering{fakeOff1, fakeOff2, fakeOff3}, nil)

			stdout := test.CaptureStdout(func() {
				actionsConfig.Catalog()
			})

			lines := strings.Split(stdout, "\n")
			for _, val := range []string{"NAME", "PLAN", "DESCRIPTION", "STATE"} {
				So(lines[1], ShouldContainSubstring, val)
			}
			for _, val := range []string{"OFFERING_1", "PLAN_1", "DESC_1", "READY"} {
				So(lines[3], ShouldContainSubstring, val)
			}
			for _, val := range []string{"OFFERING_2", "PLAN_2", "DESC_2", "READY"} {
				So(lines[4], ShouldContainSubstring, val)
			}
			for _, val := range []string{"OFFERING_3", "PLAN_3", "DESC_3", "READY"} {
				So(lines[5], ShouldContainSubstring, val)
			}
		})
	})
}

func TestListApplicationsCommand(t *testing.T) {
	actionsConfig := ActionsConfig{test.SetApiAndLoginServiceMocks(t)}
	test.FillCredentialsFile(expectedCredsFileContent)

	header := []string{"NAME", "IMAGE STATE", "STATE", "REPLICATION", "MEMORY", "DISK", "URLS", "CREATED BY", "CREATE", "UPDATED BY", "UPDATE", "MESSAGE"}
	fakeApp1Params := map[string]string{
		"name":           "App_1",
		"image_state":    "fake_State_1",
		"instance_state": "fake_state_1",
		"replication":    "1",
		"memory":         "128m",
		"quota":          "1G",
		"urls":           "fake_url_1",
		"cb":             "user_1",
		"co":             "1",
		"ub":             "user_2",
		catalogModels.LAST_STATE_CHANGE_REASON: "message",
	}

	fakeApp2Params := map[string]string{
		"name":           "App_2",
		"image_state":    "fake_State_2",
		"instance_state": "fake_state_2",
		"replication":    "2",
		"memory":         "128m",
		"quota":          "2G",
		"urls":           "fake_url_2",
		"cb":             "user_3",
		"co":             "3",
		"ub":             "user_4",
		"uo":             "4",
	}

	fakeApp3Params := map[string]string{
		"name":           "App_3",
		"image_state":    "fake_State_3",
		"instance_state": "fake_state_3",
		"replication":    "3",
		"memory":         "128m",
		"quota":          "3G",
		"urls":           "fake_url_3",
		"cb":             "user_5",
		"co":             "5",
		"ub":             "user_6",
		"uo":             "6",
	}

	fakeApp1 := test.NewFakeAppInstance(fakeApp1Params)
	fakeApp2 := test.NewFakeAppInstance(fakeApp2Params)
	fakeApp3 := test.NewFakeAppInstance(fakeApp3Params)
	Convey("Test list applications command", t, func() {
		Convey("Should pretty print applications list", func() {
			actionsConfig.ApiService.(*api.MockTapApiServiceApi).
				EXPECT().
				ListApplicationInstances().
				Return([]models.ApplicationInstance{fakeApp1, fakeApp2, fakeApp3}, nil)

			stdout := test.CaptureStdout(func() {
				actionsConfig.ListApplications()
			})

			lines := strings.Split(stdout, "\n")
			for _, val := range header {
				So(lines[1], ShouldContainSubstring, val)
			}
			for key, val := range fakeApp1Params {
				if key == catalogModels.LAST_STATE_CHANGE_REASON {
					So(lines[3], ShouldContainSubstring, printer.LastMessageMark)
				} else {
					So(lines[3], ShouldContainSubstring, val)
				}
			}
			for _, val := range fakeApp2Params {
				So(lines[4], ShouldContainSubstring, val)
			}
			for _, val := range fakeApp3Params {
				So(lines[5], ShouldContainSubstring, val)
			}
		})
	})
}

func TestUnbindInstance(t *testing.T) {
	actionsConfig := ActionsConfig{test.SetApiAndLoginServiceMocks(t)}

	fakeApplications := GetFakeApplicationInstances()
	fakeServices := GetFakeServiceInstances()

	Convey("Testing UnbindInstance", t, func() {
		actionsConfig.ApiService.(*api.MockTapApiServiceApi).
			EXPECT().
			ListApplicationInstances().
			Return(fakeApplications, nil).AnyTimes()
		actionsConfig.ApiService.(*api.MockTapApiServiceApi).
			EXPECT().
			ListServiceInstances().
			Return(fakeServices, nil).AnyTimes()

		Convey("When unbind destination is application", func() {
			Convey("When unbind source is service", func() {
				actionsConfig.ApiService.(*api.MockTapApiServiceApi).
					EXPECT().
					UnbindServiceFromApplicationInstance(fakeServices[0].Id, fakeApplications[0].Id).
					Return(http.StatusAccepted, nil)

				var err error
				stdout := test.CaptureStdout(func() {
					err = actionsConfig.UnbindInstance(fakeServices[0].Name, fakeApplications[0].Name)
				})

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("success message should be printed", func() {
					assertSuccessMessage(stdout)
				})
			})

			Convey("When unbind source is application", func() {
				actionsConfig.ApiService.(*api.MockTapApiServiceApi).
					EXPECT().
					UnbindApplicationFromApplicationInstance(fakeApplications[0].Id, fakeApplications[1].Id).
					Return(http.StatusAccepted, nil)

				var err error
				stdout := test.CaptureStdout(func() {
					err = actionsConfig.UnbindInstance(fakeApplications[0].Name, fakeApplications[1].Name)
				})

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("success message should be printed", func() {
					assertSuccessMessage(stdout)
				})
			})
		})

		Convey("When unbind destination is service", func() {
			Convey("When unbind source is service", func() {
				actionsConfig.ApiService.(*api.MockTapApiServiceApi).
					EXPECT().
					UnbindServiceFromServiceInstance(fakeServices[0].Id, fakeServices[1].Id).
					Return(http.StatusAccepted, nil)

				var err error
				stdout := test.CaptureStdout(func() {
					err = actionsConfig.UnbindInstance(fakeServices[0].Name, fakeServices[1].Name)
				})

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("success message should be printed", func() {
					assertSuccessMessage(stdout)
				})
			})

			Convey("When unbind source is application", func() {
				actionsConfig.ApiService.(*api.MockTapApiServiceApi).
					EXPECT().
					UnbindApplicationFromServiceInstance(fakeApplications[0].Id, fakeServices[0].Id).
					Return(http.StatusAccepted, nil)

				var err error
				stdout := test.CaptureStdout(func() {
					err = actionsConfig.UnbindInstance(fakeApplications[0].Name, fakeServices[0].Name)
				})

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("success message should be printed", func() {
					assertSuccessMessage(stdout)
				})
			})

			Convey("When unbind operation is not successful", func() {
				actionsConfig.ApiService.(*api.MockTapApiServiceApi).
					EXPECT().
					UnbindApplicationFromServiceInstance(fakeApplications[0].Id, fakeServices[0].Id).
					Return(http.StatusInternalServerError, errors.New("some error"))

				err := actionsConfig.UnbindInstance(fakeApplications[0].Name, fakeServices[0].Name)

				Convey("err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
			Convey("When unbind source does not exist", func() {
				fakeName := "fake_id"
				actionsConfig.ApiService.(*api.MockTapApiServiceApi).
					EXPECT().
					ListApplicationInstances().
					Return(fakeApplications, nil).AnyTimes()

				err := actionsConfig.UnbindInstance(fakeName, fakeServices[0].Name)

				Convey("err should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

}

func assertSuccessMessage(stdout string) {
	lines := strings.Split(stdout, "\n")
	So(len(lines), ShouldBeGreaterThan, 0)
	if len(lines) <= 0 {
		return
	}

	So(lines[0], ShouldEqual, SuccessMessage)
}

func GetFakeApplicationInstances() []models.ApplicationInstance {
	fakeApp1Params := map[string]string{
		"name": "App_1",
		"id":   "ID132482375",
	}

	fakeApp2Params := map[string]string{
		"name": "App_2",
		"id":   "ID623463247",
	}

	fakeApp1 := test.NewFakeAppInstance(fakeApp1Params)
	fakeApp2 := test.NewFakeAppInstance(fakeApp2Params)

	return []models.ApplicationInstance{fakeApp1, fakeApp2}
}

func GetFakeServiceInstances() []models.ServiceInstance {
	result := []models.ServiceInstance{}
	result = append(result, models.ServiceInstance{Instance: catalogModels.Instance{Id: "1", Name: "instance1", Type: "SERVICE"}})
	result = append(result, models.ServiceInstance{Instance: catalogModels.Instance{Id: "2", Name: "instance2", Type: "SERVICE"}})
	result = append(result, models.ServiceInstance{Instance: catalogModels.Instance{Id: "3", Name: "instance3", Type: "SERVICE"}})

	return result
}
