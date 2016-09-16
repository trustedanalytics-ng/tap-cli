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
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"errors"

	"github.com/trustedanalytics/tap-api-service/models"
	"github.com/trustedanalytics/tap-api-service/uaa-connector"

	"github.com/trustedanalytics/tap-cli/api"
	"io/ioutil"
	"strconv"
	"strings"
	"github.com/trustedanalytics/tap-api-service/user-management-connector"
)

var url string = "fake_url"
var login string = "fake_admin"
var pass string = "fake_password"

var expectedUaaRes = uaa_connector.LoginResponse{
	"fake_access_token",
	"fake_refresh_token",
	"fake_token_type",
	10,
	"fake_scope",
	"fake_jti",
}

var expectedCredsFileContent = "{" +
	"\"address\":\"" + url + "\"," +
	"\"username\":\"" + login + "\"," +
	"\"token\":\"" + expectedUaaRes.AccessToken + "\"," +
	"\"type\":\"" + expectedUaaRes.TokenType + "\"," +
	"\"expires\":" + strconv.Itoa(expectedUaaRes.ExpiresIn) +
	"}"

func TestLoginCommand(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)
	Convey("Test Login command", t, func() {
		Convey("Should fail when user unauthorized", func() {
			actionsConfig.ApiServiceLogin.(*api.MockTapConsoleServiceLoginApi).
				EXPECT().
				Login().
				Return(expectedUaaRes, http.StatusUnauthorized, errors.New(""))

			stdout := captureStdout(func() {
				actionsConfig.Login(url, login, pass)
			})
			So(stdout, ShouldContainSubstring, "Authentication failed")
		})
		Convey("Should fail when connecting error occurs", func() {
			loginErrorMsg := "server error"
			actionsConfig.ApiServiceLogin.(*api.MockTapConsoleServiceLoginApi).
				EXPECT().
				Login().
				Return(expectedUaaRes, http.StatusInternalServerError, errors.New(loginErrorMsg))

			stdout := captureStdout(func() {
				actionsConfig.Login(url, login, pass)
			})

			So(stdout, ShouldContainSubstring, "Error connecting: " + loginErrorMsg)
		})
		Convey("Should pass when credentials succesfully saved", func() {
			actionsConfig.ApiServiceLogin.(*api.MockTapConsoleServiceLoginApi).
				EXPECT().
				Login().
				Return(expectedUaaRes, http.StatusOK, nil)

			stdout := captureStdout(func() {
				actionsConfig.Login(url, login, pass)
			})

			b, err := ioutil.ReadFile(api.CredsPath)
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, string(expectedCredsFileContent))
			So(stdout, ShouldContainSubstring, "Authentication succeeded")
		})
	})
}

func TestInviteUserCommand(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)
	Convey("Test Login command", t, func() {
		Convey("Should fail when inviting error occurs", func() {
			errorMsg := "cannot invite given user"
			fillCredentialsFile(expectedCredsFileContent)

			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				InviteUser(login).
				Return(user_management_connector.InvitationResponse{}, errors.New(errorMsg))

			stdout := captureStdout(func() {
				actionsConfig.InviteUser(login)
			})

			So(stdout, ShouldContainSubstring, errorMsg)
		})

		Convey("Should pass when user invited", func() {
			fillCredentialsFile(expectedCredsFileContent)

			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				InviteUser(login).
				Return(user_management_connector.InvitationResponse{}, nil)

			stdout := captureStdout(func() {
				actionsConfig.InviteUser(login)
			})

			So(stdout, ShouldContainSubstring, "User " + login + " successfully invited")
		})
	})
}

func TestDeleteUserCommand(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)
	fillCredentialsFile(expectedCredsFileContent)

	Convey("Test delete-user command", t, func() {
		errorMsg := "error message"

		Convey("Should fail when no user exists", func() {
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				DeleteUser(login).
				Return(errors.New(errorMsg))

			stdout := captureStdout(func() {
				actionsConfig.DeleteUser(login)
			})

			So(stdout, ShouldContainSubstring, errorMsg)
		})

		Convey("Should pass when user exists", func() {
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				DeleteUser(login).
				Return(nil)

			stdout := captureStdout(func() {
				actionsConfig.DeleteUser(login)
			})

			So(stdout, ShouldContainSubstring, "User " + login + " successfully removed\n")
		})
	})
}

func TestCatalogCommand(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)
	fillCredentialsFile(expectedCredsFileContent)

	fakeOff1 := newFakeOffering(map[string]string{"label":"OFFERING_1", "name":"PLAN_1", "desc":"DESC_1"})

	fakeOff2 := newFakeOffering(map[string]string{"label":"OFFERING_2", "name":"PLAN_2", "desc":"DESC_2"})

	fakeOff3 := newFakeOffering(map[string]string{"label":"OFFERING_3", "name":"PLAN_3", "desc":"DESC_3"})

	Convey("Test catalog command", t, func() {
		Convey("Should pretty print offerings list", func() {
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				GetCatalog().
				Return([]models.Service{fakeOff1, fakeOff2, fakeOff3}, nil)

			stdout := captureStdout(func() {
				actionsConfig.Catalog()
			})

			lines := strings.Split(stdout, "\n")
			for _, val := range []string{"NAME", "PLAN", "DESCRIPTION"} {
				So(lines[1], ShouldContainSubstring, val)
			}
			for _, val := range []string{"OFFERING_1", "PLAN_1", "DESC_1"} {
				So(lines[3], ShouldContainSubstring, val)
			}
			for _, val := range []string{"OFFERING_2", "PLAN_2", "DESC_2"} {
				So(lines[4], ShouldContainSubstring, val)
			}
			for _, val := range []string{"OFFERING_3", "PLAN_3", "DESC_3"} {
				So(lines[5], ShouldContainSubstring, val)
			}
		})
	})
}

func TestListApplicationsCommand(t *testing.T) {
	actionsConfig := setApiAndLoginServiceMocks(t)
	fillCredentialsFile(expectedCredsFileContent)

	header := []string{"NAME", "IMAGE STATE", "STATE", "REPLICATION", "MEMORY", "DISK", "URLS", "CREATED BY", "CREATE", "UPDATED BY", "UPDATE"}

	fakeApp1Params := map[string]string{
		"name":"App_1",
		"image_state":"fake_State_1",
		"instance_state":"fake_state_1",
		"replication":"1",
		"memory":"128m",
		"quota":"1G",
		"urls":"fake_url_1",
		"cb":"user_1",
		"co":"1",
		"ub":"user_2",
		"uo":"2",
	}

	fakeApp2Params := map[string]string{
		"name":"App_2",
		"image_state":"fake_State_2",
		"instance_state":"fake_state_2",
		"replication":"2",
		"memory":"128m",
		"quota":"2G",
		"urls":"fake_url_2",
		"cb":"user_3",
		"co":"3",
		"ub":"user_4",
		"uo":"4",
	}

	fakeApp3Params := map[string]string{
		"name":"App_3",
		"image_state":"fake_State_3",
		"instance_state":"fake_state_3",
		"replication":"3",
		"memory":"128m",
		"quota":"3G",
		"urls":"fake_url_3",
		"cb":"user_5",
		"co":"5",
		"ub":"user_6",
		"uo":"6",
	}

	fakeApp1 := newFakeAppInstance(fakeApp1Params)
	fakeApp2 := newFakeAppInstance(fakeApp2Params)
	fakeApp3 := newFakeAppInstance(fakeApp3Params)
	Convey("Test list applications command", t, func() {
		Convey("Should pretty print applications list", func() {
			actionsConfig.ApiService.(*api.MockTapConsoleServiceApi).
				EXPECT().
				ListApplicationInstances().
				Return([]models.ApplicationInstance{fakeApp1, fakeApp2, fakeApp3}, nil)

			stdout := captureStdout(func() {
				actionsConfig.ListApplications()
			})

			lines := strings.Split(stdout, "\n")
			for _, val := range header {
				So(lines[1], ShouldContainSubstring, val)
			}
			for _, val := range fakeApp1Params {
				So(lines[3], ShouldContainSubstring, val)
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
