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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/trustedanalytics-ng/tap-api-service/client"
	"github.com/trustedanalytics-ng/tap-cli/api"
	"github.com/trustedanalytics-ng/tap-cli/cli/actions"
	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
	"github.com/trustedanalytics-ng/tap-go-common/logger"
)

const DefaultLogLevel = logger.LevelCritical
const (
	requiredFlagMissingExitCode    = 3
	errorReadingPassword           = 4
	flagDestinationNil             = 5
	onlyOneFlagInAlternative       = 6
	alternativeFlagMissingExitCode = 7
	alternativeFlagTooManyExitCode = 8
	flagTypeNotSupported           = 9
)

var loggerVerbosity string

func GetCommands() []cli.Command {
	defaultInfoCommand := TapInfoCommand()
	return toCommands([]TapCommand{
		loginCommand(),
		defaultInfoCommand,
		offeringCommand(),
		serviceCommand(),
		applicationCommand(),
		userCommand(),
	}, &defaultInfoCommand)
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

func sumFlags(flags ...[]cli.Flag) []cli.Flag {
	res := []cli.Flag{}
	for _, flag := range flags {
		res = append(res, flag...)
	}
	return res
}

func checkIfRequiredFlagExists(c *cli.Context, flag cli.Flag) (string, bool) {
	sFlag, ok := flag.(cli.StringFlag)
	if ok {
		if sFlag.Destination == nil {
			printMissingDestinationForFlagError(sFlag.Name)
		}
		if c.IsSet(sFlag.Name) {
			return sFlag.Name, true
		}
		// checking for default
		if sFlag.Value != "" {
			return sFlag.Name, true
		}
		return sFlag.Name, false
	}

	bFlag, ok := flag.(cli.BoolFlag)
	if ok {
		if bFlag.Destination == nil {
			printMissingDestinationForFlagError(bFlag.Name)
		}
		// bool Flag cannot have default values
		return bFlag.Name, c.IsSet(bFlag.Name)
	}

	iFlag, ok := flag.(cli.IntFlag)
	if ok {
		if iFlag.Destination == nil {
			printMissingDestinationForFlagError(iFlag.Name)
		}
		if c.IsSet(iFlag.Name) {
			return iFlag.Name, true
		}
		// checking for default. Int Flag cannot by set to 0 by default.
		if iFlag.Value != 0 {
			return iFlag.Name, true
		}
		return iFlag.Name, false
	}

	printApplicationBugInfo("Flag type not supported.")
	cli.OsExiter(flagTypeNotSupported)
	return "", false
}

func printMissingDestinationForFlagError(flagName string) {
	printApplicationBugInfo(flagName + " Destination not set.")
	cli.OsExiter(flagDestinationNil)
}

func UnrecognizedCommand(command string) {
	fmt.Println("\nUNRECOGNIZED COMMAND: '" + command + "'")
}

func PrintHelpMsg() {
	fmt.Println("\tIf you need help use 'help' command.")
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

func removalConfirmationPrompt(resourceName string) error {
	fmt.Printf("Are you sure you want to delete %s? [y/N]: ", resourceName)
	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))
	if text != "y" && text != "yes" {
		return cli.NewExitError("Canceled", -1)
	}
	return nil
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

	apiConnector, err := client.NewTapApiServiceApiWithOAuth2AndCustomSSLValidation(creds.Address, creds.TokenType, creds.Token, creds.SkipSSLValidation)
	if err != nil {
		panic(err.Error())
	}

	a.ApiService = apiConnector
	return a
}

func printApplicationBugInfo(msg string) {
	fmt.Println(msg + " This is a bug in the application. Please contact your administrator.")
}
