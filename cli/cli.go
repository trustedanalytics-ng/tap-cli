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

	"github.com/urfave/cli"

	"github.com/trustedanalytics/tap-go-common/logger"
)

const defaultLogLevel = logger.LevelCritical

func getCommands() []cli.Command {
	return []cli.Command{
		loginCommand(),
		targetCommand(),
		catalogCommand(),
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
		getInstanceCredentialsCommand(),
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

func Run() error {
	app := cli.NewApp()
	app.Name = "TAP CLI"
	app.Usage = "client for managing TAP"
	app.Version = "0.8.0"
	app.Commands = getCommands()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "verbosity,v",
			Value: defaultLogLevel,
			Usage: fmt.Sprintf("logger verbosity [%s,%s,%s,%s,%s,%s]", logger.LevelCritical, logger.LevelError,
				logger.LevelWarning, logger.LevelNotice, logger.LevelInfo, logger.LevelDebug),
		},
	}

	//override version flag to change its shortcut from v to V (reserved for verbosity global flag)
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	return app.Run(os.Args)
}
