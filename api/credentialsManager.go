package api

import (
	"encoding/json"
	"io/ioutil"

	"os"
)

//TODO: make credsPath read optionally from environment or options
var cliConfigDir string = os.Getenv("HOME") + "/.tapng-cli"
var credsPath string = cliConfigDir + "/credentials.json"

const PERMISSIONS os.FileMode = 0644

type Credentials struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
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
