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
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

func validateArgs(c *cli.Context, mustCount int) *cli.ExitError {
	if c.NArg() != mustCount {
		return cli.NewExitError("not enough args: \n"+c.Command.Name+" "+c.Command.ArgsUsage, 1)
	}
	return nil
}

func LoginCommand() cli.Command {
	return cli.Command{
		Name:      "login",
		Usage:     "login to TAP",
		ArgsUsage: "<address> <username> <password>",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 3)
			if err != nil {
				return err
			}

			return Login(c.Args().First(), c.Args().Get(1), c.Args().Get(2))
		},
	}
}

func TargetCommand() cli.Command {
	return cli.Command{
		Name:  "target",
		Usage: "print actual credentials",
		Action: func(c *cli.Context) error {
			return Target()
		},
	}
}

func CatalogCommand() cli.Command {
	return cli.Command{
		Name:  "catalog",
		Usage: "list available offerings",
		Action: func(c *cli.Context) error {
			return Catalog()
		},
	}
}

func CreateOfferingCommand() cli.Command {
	return cli.Command{
		Name:      "create-offering",
		Aliases:   []string{"co"},
		ArgsUsage: "<path to json with service definition>",
		Usage:     "create new offering",
		Action: func(c *cli.Context) error {
			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return CreateOffer(c.Args().First())
		},
	}
}

func CreateServiceCommand() cli.Command {
	return cli.Command{
		Name:      "create-service",
		ArgsUsage: "<service_name> <plan_name> <custom_name>",
		Aliases:   []string{"cs"},
		Usage:     "create instance of service",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 3)
			if err != nil {
				return err
			}

			return CreateServiceInstance(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
		},
	}
}

func DeleteServiceCommand() cli.Command {
	return cli.Command{
		Name:      "delete-service",
		ArgsUsage: "<service_custom_name>",
		Aliases:   []string{"ds"},
		Usage:     "delete instance of service",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}
			return DeleteInstance(c.Args().Get(0))
		},
	}
}

func BindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "bind-instance",
		ArgsUsage: "<src_instance_name>, <dst_instance_name>",
		Aliases:   []string{"bind"},
		Usage:     "bind instance to another",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return BindInstance(c.Args().First(), c.Args().Get(1))
		},
	}
}

func UnbindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "unbind-instance",
		ArgsUsage: "<src_instance_name>, <dst_instance_name>",
		Aliases:   []string{"unbind"},
		Usage:     "unbind instance from another",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return UnbindInstance(c.Args().First(), c.Args().Get(1))
		},
	}
}

func ListInstanceBindingsCommand() cli.Command {
	return cli.Command{
		Name:      "bindings",
		ArgsUsage: "<instanceName>",
		Usage:     "list bindings",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return GetInstanceBindings(c.Args().First())
		},
	}
}

func ListApplicationsCommand() cli.Command {
	return cli.Command{
		Name:      "applications",
		ArgsUsage: "",
		Aliases:   []string{"apps"},
		Usage:     "list applications",
		Action: func(c *cli.Context) error {
			return ListApplications()
		},
	}
}

func PushApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "push",
		ArgsUsage: "(archive_path)",
		Usage: "create application from archive provided or from compressed current directory by default,\n" +
			"\tmanifest should be in current working directory",
		Action: func(c *cli.Context) error {

			if _, err := os.Stat("manifest.json"); os.IsNotExist(err) {
				fmt.Println("manifest.json does dot exist")
				fmt.Println("Create one with metadata about your application.")
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return CompressCwdAndPushAsApplication()
			}

			return PushApplication(c.Args().First())
		},
	}
}

func ListServicesCommand() cli.Command {
	return cli.Command{
		Name:      "services",
		ArgsUsage: "",
		Aliases:   []string{"svcs"},
		Usage:     "list all service instances",
		Action: func(c *cli.Context) error {
			return ListServices()
		},
	}
}

func ScaleApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "scale",
		ArgsUsage: "<applicationName> <instances>",
		Aliases:   []string{"sc"},
		Usage:     "scale application",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			i, errr := strconv.Atoi(c.Args().Get(1))
			if errr != nil {
				return cli.NewExitError(errr.Error(), -1)
			}

			return ScaleApplication(c.Args().First(), i)
		},
	}
}

func StartApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "start",
		ArgsUsage: "<applicationName>",
		Usage:     "start application with single instance",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return StartApplication(c.Args().First())
		},
	}
}

func StopApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "stop",
		ArgsUsage: "<applicationName>",
		Usage:     "stop all application instances",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return StopApplication(c.Args().First())
		},
	}
}

func DeleteApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "delete",
		ArgsUsage: "<applicationName>",
		Aliases:   []string{"d"},
		Usage:     "delete application",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return DeleteApplication(c.Args().First())
		},
	}
}

func GetInstanceLogsCommand() cli.Command {
	return cli.Command{
		Name:      "logs",
		ArgsUsage: "<instanceName>",
		Aliases:   []string{"log"},
		Usage:     "get logs for all containers in instance",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return GetInstanceLogs(c.Args().First())
		},
	}
}

func GetApplicationsCommand() cli.Command {
	return cli.Command{
		Name:      "application",
		ArgsUsage: "<applicationName>",
		Aliases:   []string{"a"},
		Usage:     "application instance details",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return GetApplication(c.Args().First())
		},
	}
}

func GetServiceCommand() cli.Command {
	return cli.Command{
		Name:      "service",
		ArgsUsage: "<serviceName>",
		Aliases:   []string{"s"},
		Usage:     "service instance details",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return GetService(c.Args().First())
		},
	}
}
