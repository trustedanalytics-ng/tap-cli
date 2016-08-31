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
	"encoding/json"
	"io/ioutil"
	"os"
)

//TODO: make credsPath read optionally from environment or options
var cliConfigDir string = os.Getenv("HOME") + "/.tap-cli"
var credsPath string = cliConfigDir + "/credentials.json"

const PERMISSIONS os.FileMode = 0744

type Credentials struct {
	Address   string `json:"address"`
	Username  string `json:"username"`
	Token     string `json:"token"`
	TokenType string `json:"type"`
	ExpiresIn int    `json:"expires"`
}

func GetCredentials() (Credentials, error) {

	creds := Credentials{}

	b, err := ioutil.ReadFile(credsPath)
	if err != nil {
		return creds, err
	}

	err = json.Unmarshal(b, &creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

func SetCredentials(creds Credentials) error {

	jsonBytes, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	os.MkdirAll(cliConfigDir, PERMISSIONS)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(credsPath, jsonBytes, PERMISSIONS)
	return err
}
