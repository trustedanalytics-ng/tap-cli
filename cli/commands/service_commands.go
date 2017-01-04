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
	"errors"
	"strconv"

	"github.com/urfave/cli"
)

func listServicesCommand() cli.Command {
	return cli.Command{
		Name:      "services",
		ArgsUsage: "",
		Aliases:   []string{"svcs"},
		Usage:     "list all service instances",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListServices()
		},
	}
}

func getServiceCommand() cli.Command {
	return cli.Command{
		Name:      "service",
		ArgsUsage: "<serviceName>",
		Aliases:   []string{"s"},
		Usage:     "service instance details",
		Flags:     GetCommonFlags(),
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

func createServiceCommand() cli.Command {
	envsFlag := cli.StringSlice{}
	return cli.Command{
		Name:      "create-service",
		ArgsUsage: "<service_name> <plan_name> <custom_name>",
		Aliases:   []string{"cs"},
		Usage:     "create instance of service",
		Flags: sumFlags(GetCommonFlags(),
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
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}
			return newOAuth2Service().DeleteService(c.Args().Get(0))
		},
	}
}

func startServiceCommand() cli.Command {
	return cli.Command{
		Name:      "service-start",
		ArgsUsage: "<service_custom_name>",
		Usage:     "start service",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().StartService(c.Args().First())
		},
	}
}

func stopServiceCommand() cli.Command {
	return cli.Command{
		Name:      "service-stop",
		ArgsUsage: "<service_custom_name>",
		Usage:     "stop all service instances",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().StopService(c.Args().First())
		},
	}
}

func restartServiceCommand() cli.Command {
	return cli.Command{
		Name:      "service-restart",
		ArgsUsage: "<service_custom_name>",
		Usage:     "restart service",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().RestartService(c.Args().First())
		},
	}
}

func getServiceCredentialsCommand() cli.Command {
	return cli.Command{
		Name:      "credentials",
		ArgsUsage: "<instanceName>",
		Aliases:   []string{"creds"},
		Usage:     "get credentials for all containers in service instance",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().GetServiceCredentials(c.Args().First())
		},
	}
}

func exposeServiceCommand() cli.Command {
	return cli.Command{
		Name:      "expose-service",
		ArgsUsage: "<service_custom_name>, <should_expose>",
		Aliases:   []string{"expose"},
		Usage:     "expose service ports",
		Flags:     GetCommonFlags(),
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
