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
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/trustedanalytics/tap-api-service/client"
	"github.com/trustedanalytics/tap-cli/api"
	"github.com/trustedanalytics/tap-cli/cli/actions"
	commonHttp "github.com/trustedanalytics/tap-go-common/http"
	"github.com/trustedanalytics/tap-go-common/logger"
)

const DefaultLogLevel = logger.LevelCritical
const requiredFlagMissingExitCode = 3
const errorReadingPassword = 4

var loggerVerbosity string

func GetCommands() []cli.Command {
	return []cli.Command{
		loginCommand(),
		TapInfoCommand,
		listOfferings(),
		createOfferingCommand(),
		deleteOfferingCommand(),
		createServiceCommand(),
		deleteServiceCommand(),
		startServiceCommand(),
		stopServiceCommand(),
		restartServiceCommand(),
		exposeServiceCommand(),
		listInstanceBindingsCommand(),
		bindInstanceCommand(),
		unbindInstanceCommand(),
		pushApplicationCommand(),
		listApplicationsCommand(),
		getApplicationCommand(),
		listServicesCommand(),
		getServiceCommand(),
		scaleApplicationCommand(),
		startApplicationCommand(),
		stopApplicationCommand(),
		restartApplicationCommand(),
		getInstanceLogsCommand(),
		getServiceCredentialsCommand(),
		deleteApplicationCommand(),
		sendInvitationCommand(),
		resendInvitationCommand(),
		listUsersCommand(),
		listInvitationsCommand(),
		deleteInvitationCommand(),
		deleteUserCommand(),
		changeCurrentUserPasswordCommand(),
	}
}

func GetCommonFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "verbosity,v",
			Usage: fmt.Sprintf("logger verbosity [%s,%s,%s,%s,%s,%s]", logger.LevelCritical, logger.LevelError,
				logger.LevelWarning, logger.LevelNotice, logger.LevelInfo, logger.LevelDebug),
			Destination: &loggerVerbosity,
			Value:       DefaultLogLevel,
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
	res = append(res, a...)
	res = append(res, b...)
	return res
}

func checkRequiredStringFlag(flag cli.StringFlag, ctx *cli.Context) {
	value := *flag.Destination
	if value == "" {
		fmt.Println("\nMISSING PARAMETER: '", flag.Name, "'\n\nCommand usage:")
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		cli.OsExiter(requiredFlagMissingExitCode)
	}
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

func newOAuth2Service() *actions.ActionsConfig {
	a := &actions.ActionsConfig{Config: api.Config{}}

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
