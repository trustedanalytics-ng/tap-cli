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

	"github.com/urfave/cli"

	catalogModels "github.com/trustedanalytics-ng/tap-catalog/models"
)

func applicationCommand() TapCommand {

	var applicationName string
	var applicationNameFlag = cli.StringFlag{
		Name:        "name",
		Usage:       "`application name`",
		Destination: &applicationName,
	}

	const manifestFileName = "manifest.json"

	var archivePath string
	var archivePathFlag = cli.StringFlag{
		Name:        "archive-path",
		Usage:       "`path to archive with application`",
		Destination: &archivePath,
	}

	var replicas int
	var replicasFlag = cli.IntFlag{
		Name:        "replicas",
		Usage:       "`number of replicas`",
		Destination: &replicas,
	}

	confirmed := false
	var confirmationFlag = cli.BoolFlag{
		Name:        "yes",
		Usage:       "use with caution when want to suppress removal confirmation",
		Destination: &confirmed,
	}

	var listApplicationsCommand = TapCommand{
		Name:  "list",
		Usage: "list applications",
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ListApplications()
		},
	}

	var getApplicationCommand = TapCommand{
		Name:          "info",
		Usage:         "application instance details",
		RequiredFlags: []cli.Flag{applicationNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetApplication(applicationName)
		},
	}

	var pushApplicationCommand = TapCommand{
		Name: "push",
		Usage: "create application from compressed current directory (by default) or from indicated tar archive,\n" +
			"\tmanifest should be in current working directory",
		OptionalFlags: []cli.Flag{archivePathFlag},
		MainAction: func(c *cli.Context) error {
			if _, err := os.Stat(manifestFileName); os.IsNotExist(err) {
				return fmt.Errorf(manifestFileName + " does not exist: create one with metadata about your application")
			}
			if "" == archivePath {
				return newOAuth2Service().CompressCwdAndPushAsApplication()
			}
			return newOAuth2Service().PushApplication(archivePath)
		},
	}

	var deleteApplicationCommand = TapCommand{
		Name:          "delete",
		Usage:         "delete application",
		RequiredFlags: []cli.Flag{applicationNameFlag},
		OptionalFlags: []cli.Flag{confirmationFlag},
		MainAction: func(c *cli.Context) error {
			if !confirmed {
				err := removalConfirmationPrompt("application " + applicationName)
				cli.HandleExitCoder(err)
			}
			return newOAuth2Service().DeleteApplication(applicationName)
		},
	}

	var startApplicationCommand = TapCommand{
		Name:          "start",
		Usage:         "start application",
		RequiredFlags: []cli.Flag{applicationNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().StartApplication(applicationName)
		},
	}

	var stopApplicationCommand = TapCommand{
		Name:          "stop",
		Usage:         "stop application",
		RequiredFlags: []cli.Flag{applicationNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().StopApplication(applicationName)
		},
	}

	var restartApplicationCommand = TapCommand{
		Name:          "restart",
		Usage:         "restart application",
		RequiredFlags: []cli.Flag{applicationNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().RestartApplication(applicationName)
		},
	}

	var scaleApplicationCommand = TapCommand{
		Name:          "scale",
		Usage:         "scale application",
		RequiredFlags: []cli.Flag{applicationNameFlag, replicasFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ScaleApplication(applicationName, replicas)
		},
	}

	var applicationLogsShowCommand = TapCommand{
		Name:          "show",
		Usage:         "show application logs",
		RequiredFlags: []cli.Flag{applicationNameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetInstanceLogs(applicationName)
		},
	}

	var getInstanceLogsCommand = TapCommand{
		Name:              "logs",
		Usage:             "get logs for all containers in instance",
		Subcommands:       []TapCommand{applicationLogsShowCommand},
		DefaultSubcommand: &applicationLogsShowCommand,
	}

	return TapCommand{
		Name:  "application",
		Usage: "application context commands",
		Subcommands: []TapCommand{
			listApplicationsCommand,
			getApplicationCommand,
			pushApplicationCommand,
			deleteApplicationCommand,
			startApplicationCommand,
			stopApplicationCommand,
			restartApplicationCommand,
			scaleApplicationCommand,
			getInstanceLogsCommand,
			bindingCommands(catalogModels.InstanceTypeApplication),
		},
		DefaultSubcommand: &getApplicationCommand,
	}
}
