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

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	tapCli "github.com/trustedanalytics/tap-cli/cli"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", r)
			os.Exit(1)
		}
	}()

	app := cli.NewApp()
	app.Name = "TAP CLI"
	app.Usage = "client for managing TAP"
	app.Version = "0.8.0"

	app.Commands = []cli.Command{
		tapCli.LoginCommand(),
		tapCli.TargetCommand(),
		tapCli.CatalogCommand(),
		tapCli.CreateOfferingCommand(),
		tapCli.DeleteOfferingCommand(),
		tapCli.CreateServiceCommand(),
		tapCli.DeleteServiceCommand(),
		tapCli.ListInstanceBindingsCommand(),
		tapCli.BindInstanceCommand(),
		tapCli.UnbindInstanceCommand(),
		tapCli.PushApplicationCommand(),
		tapCli.ListApplicationsCommand(),
		tapCli.GetApplicationCommand(),
		tapCli.ListServicesCommand(),
		tapCli.GetServiceCommand(),
		tapCli.ScaleApplicationCommand(),
		tapCli.StartApplicationCommand(),
		tapCli.StopApplicationCommand(),
		tapCli.GetInstanceLogsCommand(),
		tapCli.GetInstanceCredentialsCommand(),
		tapCli.DeleteApplicationCommand(),
		tapCli.SendInvitationCommand(),
		tapCli.ResendInvitationCommand(),
		tapCli.ListUsersCommand(),
		tapCli.ListInvitationsCommand(),
		tapCli.DeleteInvitationCommand(),
		tapCli.DeleteUserCommand(),
		tapCli.ChangeCurrentUserPasswordCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
