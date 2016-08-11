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

package api

import (
	"errors"
	"os"

	"github.com/trustedanalytics/tapng-console-service/client"
)

type Config struct {
	ConsoleServiceApi client.TapConsoleServiceApi
}

var ConnectionConfig *Config

func InitConnection() error {

	creds, err := GetCredentials()
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Please login first!")
		}
		return err
	}

	apiConnector, err := client.NewTapConsoleServiceApiWithBasicAuth(
		"http://"+creds.Address,
		creds.Username,
		creds.Password,
	)
	if err != nil {
		return err
	}

	ConnectionConfig = &Config{}
	ConnectionConfig.ConsoleServiceApi = apiConnector

	return nil
}
