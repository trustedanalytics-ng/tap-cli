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
	"os"

	"github.com/urfave/cli"

	tapngCli "github.com/trustedanalytics/tapng-cli/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "TAP CLI"
	app.Usage = "client for managing TAP"
	app.Version = "0.8" //todo: should probably change before release

	app.Commands = []cli.Command{
		tapngCli.LoginCommand(),
		tapngCli.TargetCommand(),
		tapngCli.CatalogCommand(),
		tapngCli.CreateOfferingCommand(),
		tapngCli.CreateServiceCommand(),
		tapngCli.DeleteServiceCommand(),
		tapngCli.ListInstanceBindingsCommand(),
		tapngCli.BindInstanceCommand(),
		tapngCli.UnbindInstanceCommand(),
		tapngCli.PushApplicationCommand(),
		tapngCli.ListApplicationsCommand(),
		tapngCli.GetApplicationsCommand(),
		tapngCli.ListServicesCommand(),
		tapngCli.ScaleApplicationCommand(),
		tapngCli.StartApplicationCommand(),
		tapngCli.StopApplicationCommand(),
		tapngCli.GetInstanceLogsCommand(),
		tapngCli.DeleteApplicationCommand(),
	}

	app.Run(os.Args)
}
