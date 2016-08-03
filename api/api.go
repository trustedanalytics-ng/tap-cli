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
