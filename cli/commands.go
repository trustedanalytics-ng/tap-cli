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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"

	"github.com/trustedanalytics/tap-api-service/client"
	"github.com/trustedanalytics/tap-cli/api"
	commonHttp "github.com/trustedanalytics/tap-go-common/http"
	"github.com/trustedanalytics/tap-go-common/logger"
)

var loggerVerbosity string

func getCommonFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "verbosity,v",
			Usage: fmt.Sprintf("logger verbosity [%s,%s,%s,%s,%s,%s] (default: %q)", logger.LevelCritical, logger.LevelError,
				logger.LevelWarning, logger.LevelNotice, logger.LevelInfo, logger.LevelDebug, defaultLogLevel),
			Destination: &loggerVerbosity,
		},
	}
}

func handleCommonFlags(c *cli.Context) error {
	if loggerVerbosity == "" {
		loggerVerbosity = c.GlobalString("verbosity")
	}
	if err := client.SetLoggerLevel(loggerVerbosity); err != nil {
		return err
	}
	if err := commonHttp.SetLoggerLevel(loggerVerbosity); err != nil {
		return err
	}
	return nil
}

func sumFlags(a []cli.Flag, b []cli.Flag) []cli.Flag {
	res := []cli.Flag{}
	for _, flag := range a {
		res = append(res, flag)
	}
	for _, flag := range b {
		res = append(res, flag)
	}
	return res
}

func validateArgs(c *cli.Context, mustCount int) *cli.ExitError {
	if c.NArg() != mustCount {
		return cli.NewExitError("not enough args: \n"+c.Command.Name+" "+c.Command.ArgsUsage, 1)
	}
	return nil
}

func validateAndSplitEnvFlags(envs cli.StringSlice) (map[string]string, *cli.ExitError) {
	result := make(map[string]string)
	for _, env := range envs {
		splittedEnv := strings.Split(env, "=")
		if len(splittedEnv) < 2 || splittedEnv[0] == "" {
			return result, cli.NewExitError("use NAME=VALUE format for env: \n"+env, 1)
		}
		key := splittedEnv[0]
		value := strings.TrimPrefix(env, key+"=")
		result[key] = value
	}
	return result, nil
}

func sendInvitationCommand() cli.Command {
	return cli.Command{
		Name:      "invite",
		Usage:     "invite new user to TAP",
		ArgsUsage: "<email>",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().SendInvitation(c.Args().First())
		},
	}
}

func resendInvitationCommand() cli.Command {
	return cli.Command{
		Name:      "reinvite",
		Usage:     "resend invitation for user",
		ArgsUsage: "<email>",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().ResendInvitation(c.Args().First())
		},
	}
}

func listUsersCommand() cli.Command {
	return cli.Command{
		Name:  "users",
		Usage: "list platform users",
		Flags: getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListUsers()
		},
	}
}

func listInvitationsCommand() cli.Command {
	return cli.Command{
		Name:    "invitations",
		Usage:   "list pending invitations",
		Aliases: []string{"invs"},
		Flags:   getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListInvitations()
		},
	}
}

func deleteInvitationCommand() cli.Command {
	return cli.Command{
		Name:      "delete-invitation",
		Usage:     "delete invitation",
		Aliases:   []string{"di"},
		ArgsUsage: "<email>",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().DeleteInvitation(c.Args().First())
		},
	}
}

func deleteUserCommand() cli.Command {
	return cli.Command{
		Name:      "delete-user",
		Usage:     "delete user from TAP",
		Aliases:   []string{"du"},
		ArgsUsage: "<email>",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().DeleteUser(c.Args().First())
		},
	}
}

func changeCurrentUserPasswordCommand() cli.Command {
	return cli.Command{
		Name:      "chpasswd",
		Usage:     "change password of currently logged user",
		ArgsUsage: "<currentPassword> <newPassword>",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return newOAuth2Service().ChangeCurrentUserPassword(c.Args().First(), c.Args().Get(1))
		},
	}
}

func loginCommand() cli.Command {
	return cli.Command{
		Name:      "login",
		Usage:     "login to TAP",
		ArgsUsage: "<address> <username> <password>",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 3)
			if err != nil {
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
		Flags:   getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().Target()
		},
	}
}

func catalogCommand() cli.Command {
	return cli.Command{
		Name:    "catalog",
		Aliases: []string{"o"},
		Usage:   "list available offerings",
		Flags:   getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().Catalog()
		},
	}
}

func createOfferingCommand() cli.Command {
	return cli.Command{
		Name:      "create-offering",
		Aliases:   []string{"co"},
		ArgsUsage: "<path to json with service definition>",
		Usage:     "create new offering",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().CreateOffer(c.Args().First())
		},
	}
}

func deleteOfferingCommand() cli.Command {
	return cli.Command{
		Name:      "delete-offering",
		ArgsUsage: "<offering_custom_name>",
		Aliases:   []string{"do"},
		Usage:     "delete offering",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().DeleteOffering(c.Args().First())
		},
	}
}

func createServiceCommand() cli.Command {
	envsFlag := cli.StringSlice{}
	return cli.Command{
		Name:      "create-service",
		ArgsUsage: "<service_name> <plan_name> <custom_name>",
		Aliases:   []string{"cs"},
		Usage:     "create instance of service",
		Flags: sumFlags(getCommonFlags(),
			[]cli.Flag{
				cli.StringSliceFlag{
					Name:  "env, e",
					Usage: "pass envs in format: `NAME=VALUE` this flag can be used multiple times",
					Value: &envsFlag,
				},
			},
		),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 3)
			if err != nil {
				return err
			}

			envs, err := validateAndSplitEnvFlags(envsFlag)
			if err != nil {
				return err
			}

			return newOAuth2Service().CreateServiceInstance(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2), envs)
		},
	}
}

func deleteServiceCommand() cli.Command {
	return cli.Command{
		Name:      "delete-service",
		ArgsUsage: "<service_custom_name>",
		Aliases:   []string{"ds"},
		Usage:     "delete instance of service",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}
			return newOAuth2Service().DeleteInstance(c.Args().Get(0))
		},
	}
}

func exposeServiceCommand() cli.Command {
	return cli.Command{
		Name:      "expose-service",
		ArgsUsage: "<service_custom_name>, <should_expose>",
		Aliases:   []string{"expose"},
		Usage:     "expose service ports",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			exposed, parseErr := strconv.ParseBool(c.Args().Get(1))
			if parseErr != nil {
				return errors.New("exposed argument has to be a boolean value: true/false")
			}

			return newOAuth2Service().ExposeService(c.Args().First(), exposed)
		},
	}
}

func bindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "bind-instance",
		ArgsUsage: "<src_instance_name>, <dst_instance_name>",
		Aliases:   []string{"bind"},
		Usage:     "bind instance to another",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return newOAuth2Service().BindInstance(c.Args().First(), c.Args().Get(1))
		},
	}
}

func unbindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "unbind-instance",
		ArgsUsage: "<src_instance_name>, <dst_instance_name>",
		Aliases:   []string{"unbind"},
		Usage:     "unbind instance from another",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return newOAuth2Service().UnbindInstance(c.Args().First(), c.Args().Get(1))
		},
	}
}

func listInstanceBindingsCommand() cli.Command {
	return cli.Command{
		Name:      "bindings",
		ArgsUsage: "<instanceName>",
		Usage:     "list bindings",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().GetApplicationBindings(c.Args().First())
		},
	}
}

func listApplicationsCommand() cli.Command {
	return cli.Command{
		Name:      "applications",
		ArgsUsage: "",
		Aliases:   []string{"apps"},
		Usage:     "list applications",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListApplications()
		},
	}
}

func pushApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "push",
		ArgsUsage: "(archive_path)",
		Usage: "create application from archive provided or from compressed current directory by default,\n" +
			"\tmanifest should be in current working directory",
		Flags: getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			if _, err := os.Stat("manifest.json"); os.IsNotExist(err) {
				return fmt.Errorf("manifest.json does not exist: create one with metadata about your application")
			}

			err := validateArgs(c, 1)
			if err != nil {
				return newOAuth2Service().CompressCwdAndPushAsApplication()
			}

			return newOAuth2Service().PushApplication(c.Args().First())
		},
	}
}

func listServicesCommand() cli.Command {
	return cli.Command{
		Name:      "services",
		ArgsUsage: "",
		Aliases:   []string{"svcs"},
		Usage:     "list all service instances",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListServices()
		},
	}
}

func scaleApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "scale",
		ArgsUsage: "<applicationName> <instances>",
		Aliases:   []string{"sc"},
		Usage:     "scale application",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			i, errr := strconv.Atoi(c.Args().Get(1))
			if errr != nil {
				return cli.NewExitError(errr.Error(), -1)
			}

			return newOAuth2Service().ScaleApplication(c.Args().First(), i)
		},
	}
}

func startApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "start",
		ArgsUsage: "<applicationName>",
		Usage:     "start application with single instance",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().StartApplication(c.Args().First())
		},
	}
}

func stopApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "stop",
		ArgsUsage: "<applicationName>",
		Usage:     "stop all application instances",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().StopApplication(c.Args().First())
		},
	}
}

func deleteApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "delete",
		ArgsUsage: "<applicationName>",
		Aliases:   []string{"d"},
		Usage:     "delete application",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().DeleteApplication(c.Args().First())
		},
	}
}

func getInstanceLogsCommand() cli.Command {
	return cli.Command{
		Name:      "logs",
		ArgsUsage: "<instanceName>",
		Aliases:   []string{"log"},
		Usage:     "get logs for all containers in instance",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().GetInstanceLogs(c.Args().First())
		},
	}
}

func getInstanceCredentialsCommand() cli.Command {
	return cli.Command{
		Name:      "credentials",
		ArgsUsage: "<instanceName>",
		Aliases:   []string{"creds"},
		Usage:     "get credentials for all containers in service instance",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().GetInstanceCredentials(c.Args().First())
		},
	}
}

func getApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "application",
		ArgsUsage: "<applicationName>",
		Aliases:   []string{"a"},
		Usage:     "application instance details",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().GetApplication(c.Args().First())
		},
	}
}

func getServiceCommand() cli.Command {
	return cli.Command{
		Name:      "service",
		ArgsUsage: "<serviceName>",
		Aliases:   []string{"s"},
		Usage:     "service instance details",
		Flags:     getCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().GetService(c.Args().First())
		},
	}
}

func newBasicAuthService(address string, username string, password string) *ActionsConfig {
	if !isProcotolSet(address) {
		address = "https://" + address
	}
	apiConnector, err := client.NewTapApiServiceLoginApiWithBasicAuth(address, username, password)
	if err != nil {
		panic(err)
	}
	return &ActionsConfig{api.Config{nil, apiConnector}}
}

func isProcotolSet(address string) bool {
	index := strings.Index(address[0:], "://")
	return index != -1
}

func newOAuth2Service() *ActionsConfig {
	a := &ActionsConfig{api.Config{}}

	creds, err := a.GetCredentials()
	if err != nil {
		if os.IsNotExist(err) {
			panic("Please login first!")
		}
		panic(err.Error())
	}

	apiConnector, err := client.NewTapApiServiceApiWithOAuth2(creds.Address, creds.TokenType, creds.Token)
	if err != nil {
		panic(err.Error())
	}

	a.ApiService = apiConnector
	return a
}
