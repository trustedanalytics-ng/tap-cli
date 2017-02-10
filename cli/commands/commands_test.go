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

package commands

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli"

	"github.com/trustedanalytics-ng/tap-api-service/client"
	"github.com/trustedanalytics-ng/tap-cli/cli/test"
)

func init() {
	test.SwitchToTestCredentialsFile()
}

func TestSumFlags(t *testing.T) {
	testCases := []struct {
		a   []cli.Flag
		b   []cli.Flag
		res []cli.Flag
	}{
		{[]cli.Flag{}, []cli.Flag{}, []cli.Flag{}},
		{[]cli.Flag{cli.StringSliceFlag{Name: "abc"}, cli.StringFlag{Name: "cba"}}, []cli.Flag{cli.BoolFlag{Name: "qwe"}}, []cli.Flag{cli.StringSliceFlag{Name: "abc"}, cli.StringFlag{Name: "cba"}, cli.BoolFlag{Name: "qwe"}}},
		{[]cli.Flag{}, []cli.Flag{cli.BoolFlag{Name: "qwe"}}, []cli.Flag{cli.BoolFlag{Name: "qwe"}}},
		{[]cli.Flag{cli.StringSliceFlag{Name: "abc"}, cli.StringFlag{Name: "cba"}}, []cli.Flag{}, []cli.Flag{cli.StringSliceFlag{Name: "abc"}, cli.StringFlag{Name: "cba"}}},
	}

	Convey("For set of test cases sumFlags should return proper responses", t, func() {
		for _, tc := range testCases {
			Convey(fmt.Sprintf("For set of flags %v and %v sumFlags should result in %v", tc.a, tc.b, tc.res), func() {
				response := sumFlags(tc.a, tc.b)
				So(response, ShouldResemble, tc.res)
			})
		}
	})
}

func TestApiAndLoginServiceSetters(t *testing.T) {
	Convey("Test OAuth2 login2", t, func() {
		Convey("Should fail when no credentials.json file", func() {
			test.DeleteTestCredentialsFile()

			So(func() {
				newOAuth2Service()
			}, ShouldPanicWith, "Please login first!")
		})
		Convey("Should fail when wrong format in credentials.json file", func() {
			wrongContent := "@"
			test.FillCredentialsTestFile(wrongContent)

			So(func() {
				newOAuth2Service()
			}, ShouldPanicWith, "invalid character '"+wrongContent+"' looking for beginning of value")
		})
	})
}

func TestNewBasicAuthService(t *testing.T) {
	Convey("Should trim ending slash if provided", t, func() {
		basicAuth := newBasicAuthService("myaddress.com/", "user", "password", false)
		basicCreds := basicAuth.ApiServiceLogin.(*client.TapApiServiceApiBasicAuthConnector)

		So(basicCreds.Address, ShouldEqual, "https://myaddress.com")
	})
	Convey("Should add https address if address not provided", t, func() {
		basicAuth := newBasicAuthService("myaddress.com", "user", "password", false)
		basicCreds := basicAuth.ApiServiceLogin.(*client.TapApiServiceApiBasicAuthConnector)

		So(basicCreds.Address, ShouldEqual, "https://myaddress.com")
	})
	Convey("Should not add https", t, func() {
		Convey("when there is http:// ", func() {
			basicAuth := newBasicAuthService("http://myaddress.com", "user", "password", false)
			basicCreds := basicAuth.ApiServiceLogin.(*client.TapApiServiceApiBasicAuthConnector)

			So(basicCreds.Address, ShouldEqual, "http://myaddress.com")
		})
		Convey("when there is ftp:// ", func() {
			basicAuth := newBasicAuthService("ftp://myaddress.com", "user", "password", false)
			basicCreds := basicAuth.ApiServiceLogin.(*client.TapApiServiceApiBasicAuthConnector)

			So(basicCreds.Address, ShouldEqual, "ftp://myaddress.com")
		})
	})
}

func TestValidateAndSplitEnvFlags(t *testing.T) {
	Convey("Should successfully call validateAndSplitEnvFlags", t, func() {
		envsFlag := cli.StringSlice{
			"myenv=value", "my_env=", "special=Sign=123",
		}

		result, err := validateAndSplitEnvFlags(envsFlag)
		So(err, ShouldBeNil)
		So(len(result), ShouldEqual, 3)
		So(result, ShouldContainKey, "myenv")
		So(result, ShouldContainKey, "my_env")
		So(result, ShouldContainKey, "special")
		So(result["myenv"], ShouldEqual, "value")
		So(result["my_env"], ShouldEqual, "")
		So(result["special"], ShouldEqual, "Sign=123")
	})

	Convey("Should return error when key is empty", t, func() {
		envsFlag := cli.StringSlice{
			"=value",
		}

		_, err := validateAndSplitEnvFlags(envsFlag)
		So(err, ShouldNotBeNil)
	})

	Convey("Should return error when there is only value", t, func() {
		envsFlag := cli.StringSlice{
			"value",
		}

		_, err := validateAndSplitEnvFlags(envsFlag)
		So(err, ShouldNotBeNil)
	})
}
