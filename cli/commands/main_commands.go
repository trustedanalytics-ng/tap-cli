package commands

import (
	"strings"

	"github.com/trustedanalytics/tap-api-service/client"
	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-cli/cli/actions"
	"github.com/urfave/cli"
)

func loginCommand() cli.Command {
	return cli.Command{
		Name:      "login",
		Usage:     "login to TAP. You can omitt address if it was set as target previously",
		ArgsUsage: "[<address>] <username> <password>",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 3)
			//if there are less than 3 args...
			if err != nil {
				a := &actions.ActionsConfig{Config: api.Config{}}
				creds, errcreds := a.GetCredentials()
				//...and we have credentials..
				if errcreds == nil && creds.Address != "" {
					err := validateArgs(c, 2)
					if err != nil {
						return err
					}

					return newBasicAuthService(creds.Address, c.Args().First(), c.Args().Get(1)).Login()
				}
				return err
			}

			return newBasicAuthService(c.Args().First(), c.Args().Get(1), c.Args().Get(2)).Login()
		},
	}
}

func targetCommand() cli.Command {
	return cli.Command{
		Name:    "target",
		Aliases: []string{"t"},
		Usage:   "print actual credentials",
		Flags:   GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().Target()
		},
	}
}

func newBasicAuthService(address string, username string, password string) *actions.ActionsConfig {
	address = trimEndingSlash(address)
	if !isProcotolSet(address) {
		address = "https://" + address
	}
	apiConnector, err := client.NewTapApiServiceLoginApiWithBasicAuth(address, username, password)
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
