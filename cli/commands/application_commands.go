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
	"strconv"

	"github.com/urfave/cli"
)

func listApplicationsCommand() cli.Command {
	return cli.Command{
		Name:      "applications",
		ArgsUsage: "",
		Aliases:   []string{"apps"},
		Usage:     "list applications",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListApplications()
		},
	}
}

func getApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "application",
		ArgsUsage: "<applicationName>",
		Aliases:   []string{"a"},
		Usage:     "application instance details",
		Flags:     GetCommonFlags(),
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

func pushApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "push",
		ArgsUsage: "(archive_path)",
		Usage: "create application from archive provided or from compressed current directory by default,\n" +
			"\tmanifest should be in current working directory",
		Flags: GetCommonFlags(),
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

func deleteApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "delete",
		ArgsUsage: "<applicationName>",
		Aliases:   []string{"d"},
		Usage:     "delete application",
		Flags:     GetCommonFlags(),
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

func startApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "start",
		ArgsUsage: "<applicationName>",
		Usage:     "start application with single instance",
		Flags:     GetCommonFlags(),
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
		Flags:     GetCommonFlags(),
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

func restartApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "restart",
		ArgsUsage: "<applicationName>",
		Usage:     "restart application",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().RestartApplication(c.Args().First())
		},
	}
}

func scaleApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "scale",
		ArgsUsage: "<applicationName> <instances>",
		Aliases:   []string{"sc"},
		Usage:     "scale application",
		Flags:     GetCommonFlags(),
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
