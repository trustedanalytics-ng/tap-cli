package commands

import (
	"strings"

	"fmt"
	"github.com/howeyc/gopass"
	"github.com/trustedanalytics/tap-api-service/client"
	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-cli/cli/actions"
	"github.com/urfave/cli"
)

var TapInfoCommand = cli.Command{
	Name:  "[info]",
	Usage: "prints info about current api and user",
	Flags: GetCommonFlags(),
	Action: func(c *cli.Context) error {
		if err := handleCommonFlags(c); err != nil {
			return err
		}
		return newOAuth2Service().Target()
	},
}

func loginCommand() cli.Command {

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

	return cli.Command{
		Name:      "login",
		Usage:     "login to TAP. If you don't provide password you'll be promped for it.",
		ArgsUsage: "--" + apiFlag.Name + "=<api address> --" + usernameFlag.Name + "=<username> [--" + passwordFlag.Name + "=<password>]",
		Flags:     sumFlags([]cli.Flag{apiFlag, usernameFlag, passwordFlag}, GetCommonFlags()),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}
			checkRequiredStringFlag(apiFlag, c)
			checkRequiredStringFlag(usernameFlag, c)
			if password == "" {
				password = promptForPassword()
			}
			return newBasicAuthService(apiUrl, username, password).Login()
		},
	}
}

func promptForPassword() string {
	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println("Error reading password: ", err)
		cli.OsExiter(errorReadingPassword)
	}
	password := string(pass)
	if password == "" {
		fmt.Println("Password cannot be empty")
		cli.OsExiter(errorReadingPassword)
	}
	return password
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
