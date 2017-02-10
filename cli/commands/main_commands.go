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
	"strings"

	"fmt"
	"github.com/howeyc/gopass"
	"github.com/trustedanalytics-ng/tap-api-service/client"
	"github.com/trustedanalytics-ng/tap-cli/api"
	"github.com/trustedanalytics-ng/tap-cli/cli/actions"
	"github.com/urfave/cli"
)

func TapInfoCommand() TapCommand {
	return TapCommand{
		Name:  "info",
		Usage: "prints info about current api and user",
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().Target()
		},
	}
}

func loginCommand() TapCommand {

	var apiUrl string
	var apiFlag = cli.StringFlag{
		Name:        "api",
		Usage:       "TAP `API` you would like to use",
		Destination: &apiUrl,
	}

	var username string
	var usernameFlag = cli.StringFlag{
		Name:        "username",
		Destination: &username,
	}

	var password string
	var passwordFlag = cli.StringFlag{
		Name:        "password",
		Usage:       "user `PASSWORD`",
		Destination: &password,
	}

	var skipSSLValidation bool
	var skipSSLValidationFlag = cli.BoolFlag{
		Name:        "skip-ssl-validation",
		Usage:       "When specified, the SSL certificate of the API service will not be validated",
		Destination: &skipSSLValidation,
	}

	return TapCommand{
		Name:          "login",
		Usage:         "login to TAP. If you don't provide password you'll be prompted for it.",
		OptionalFlags: []cli.Flag{passwordFlag, skipSSLValidationFlag},
		RequiredFlags: []cli.Flag{apiFlag, usernameFlag},
		MainAction: func(c *cli.Context) error {
			if password == "" {
				password = promptForSensitive("Password")
			}
			return newBasicAuthService(apiUrl, username, password, skipSSLValidation).Login(skipSSLValidation)
		},
	}
}

func promptForSensitive(name string) string {
	fmt.Printf("%s: ", name)
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", name, err)
		cli.OsExiter(errorReadingPassword)
	}
	password := string(pass)
	if password == "" {
		fmt.Printf("%s cannot be empty\n", name)
		cli.OsExiter(errorReadingPassword)
	}
	return password
}

func newBasicAuthService(address string, username string, password string, skipSSLValidation bool) *actions.ActionsConfig {
	address = trimEndingSlash(address)
	if !isProcotolSet(address) {
		address = "https://" + address
	}
	apiConnector, err := client.NewTapApiServiceLoginApiWithSSLValidationAndBasicAuth(address, username, password, skipSSLValidation)
	if err != nil {
		panic(err)
	}
	return &actions.ActionsConfig{Config: api.Config{ApiService: nil, ApiServiceLogin: apiConnector}}
}

func trimEndingSlash(str string) string {
	return strings.TrimSuffix(str, "/")
}

func isProcotolSet(address string) bool {
	index := strings.Index(address[0:], "://")
	return index != -1
}
