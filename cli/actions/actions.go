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
	"net/http"
	"os"

	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-cli/cli/printer"
)

const successMessage = "OK"

type ActionsConfig struct {
	api.Config
}

func announceSuccessfulOperation() {
	fmt.Println(successMessage)
}

func (a *ActionsConfig) Login() error {
	address, username, _ := a.ApiServiceLogin.GetLoginCredentials()
	creds := api.Credentials{}
	creds.Address = address
	creds.Username = username

	fmt.Println("Authenticating...")

	err := a.ApiServiceLogin.Introduce()
	if err != nil {
		return err
	}

	loginResp, status, err := a.ApiServiceLogin.Login()
	if status == http.StatusUnauthorized {
		return fmt.Errorf("Authentication failed")
	} else if status == http.StatusNotFound {
		return fmt.Errorf("CLI <-> API service incompatibility detected. Check your CLI version")
	} else if err != nil {
		return fmt.Errorf("Authentication failed: %v", err)
	}

	creds.Token = loginResp.AccessToken
	creds.TokenType = loginResp.TokenType
	creds.ExpiresIn = loginResp.ExpiresIn

	if err = a.SetCredentials(creds); err != nil {
		return err
	}

	fmt.Println("Authentication succeeded")

	return nil
}

func (a *ActionsConfig) Target() error {
	creds, err := a.GetCredentials()
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Please login first!")
		}
		return err
	}
	printCredentials(creds)
	return nil
}

func printCredentials(creds api.Credentials) {
	printableCredentials := []printer.Printable{printer.PrintableCredentials{Credentials: creds}}
	printer.PrintTable(printableCredentials)
}
