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
	"github.com/urfave/cli"
)

func serviceCommand() TapCommand {

	var serviceName string
	var serviceNameFlag = cli.StringFlag{
		Name:        "name",
		Usage:       "`serviceName`",
		Destination: &serviceName,
	}

	var offeringName string
	var offeringNameFlag = cli.StringFlag{
		Name:        "offering",
		Usage:       "`offeringName`",
		Destination: &offeringName,
	}

	var planName string
	var planNameFlag = cli.StringFlag{
		Name:        "plan",
		Usage:       "`planName`",
		Destination: &planName,
	}

	var envs cli.StringSlice
	var envsFlag = cli.StringSliceFlag{
		Name:  "envs",
		Usage: "pass envs in format: `NAME=VALUE` this flag can be used multiple times",
		Value: &envs,
	}

	confirmed := false
	var confirmationFlag = cli.BoolFlag{
		Name:        "yes",
		Usage:       "use with caution when want to suppress removal confirmation",
		Destination: &confirmed,
	}

	var listServiceCommand = TapCommand{
		Name:  "list",
		Usage: "list services",
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ListServices()
		},
	}

	var serviceInfoCommand = TapCommand{
		Name:          "info",
		Usage:         "service instance details",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetService(serviceName)
		},
	}

	var createServiceCommand = TapCommand{
		Name:          "create",
		Usage:         "create new service instance",
		RequiredFlags: []cli.Flag{serviceNameFlag, offeringNameFlag, planNameFlag},
		OptionalFlags: []cli.Flag{envsFlag},
		MainAction: func(c *cli.Context) error {
			splitEnvs, err := validateAndSplitEnvFlags(envs)
			if err != nil {
				return err
			}
			return newOAuth2Service().CreateServiceInstance(offeringName, planName, serviceName, splitEnvs)
		},
	}

	var deleteServiceCommand = TapCommand{
		Name:          "delete",
		Usage:         "delete service instance",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		OptionalFlags: []cli.Flag{confirmationFlag},
		MainAction: func(c *cli.Context) error {
			if !confirmed {
				err := removalConfirmationPrompt("service instance " + serviceName)
				cli.HandleExitCoder(err)
			}
			return newOAuth2Service().DeleteService(serviceName)
		},
	}

	var startServiceCommand = TapCommand{
		Name:          "start",
		Usage:         "start service instance",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().StartService(serviceName)
		},
	}

	var stopServiceCommand = TapCommand{
		Name:          "stop",
		Usage:         "stop service instance",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().StopService(serviceName)
		},
	}

	var restartServiceCommand = TapCommand{
		Name:          "restart",
		Usage:         "restart service instance",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().RestartService(serviceName)
		},
	}

	var serviceLogsShowCommand = TapCommand{
		Name:          "show",
		Usage:         "show service instances's logs",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetInstanceLogs(serviceName)
		},
	}

	var serviceLogsCommand = TapCommand{
		Name:        "logs",
		Usage:       "service instances's logs",
		Subcommands: []TapCommand{serviceLogsShowCommand},
		MainAction: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
	}

	var serviceCredentialsShowCommand = TapCommand{
		Name:          "show",
		Usage:         "show service instances's credentials",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetServiceCredentials(serviceName)
		},
	}

	var serviceCredentialsCommand = TapCommand{
		Name:        "credentials",
		Usage:       "service instances's credentials",
		Subcommands: []TapCommand{serviceCredentialsShowCommand},
		MainAction: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
	}

	var exposeServiceCommand = TapCommand{
		Name:          "expose",
		Usage:         "expose service instance under externally available URL",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ExposeService(serviceName, true)
		},
	}

	var unexposeServiceCommand = TapCommand{
		Name:          "unexpose",
		Usage:         "unexpose service instance and remove externally available URL",
		RequiredFlags: []cli.Flag{serviceNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ExposeService(serviceName, false)
		},
	}

	return TapCommand{
		Name:  "service",
		Usage: "service context commands",
		Subcommands: []TapCommand{
			listServiceCommand,
			serviceInfoCommand,
			createServiceCommand,
			deleteServiceCommand,
			startServiceCommand,
			stopServiceCommand,
			restartServiceCommand,
			serviceLogsCommand,
			serviceCredentialsCommand,
			exposeServiceCommand,
			unexposeServiceCommand,
		},
		MainAction: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
	}
}
