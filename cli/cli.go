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
	"os"

	"github.com/urfave/cli"

	"github.com/trustedanalytics-ng/tap-cli/cli/commands"
)

func Run() error {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	cli.SubcommandHelpTemplate = SubcommandHelpTemplate

	app := cli.NewApp()
	app.Name = "TAP CLI"
	app.Usage = "client for managing TAP"
	app.Version = "0.8.0"
	app.Commands = commands.GetCommands()
	app.Flags = commands.GetCommonFlags()

	app.Action = func(c *cli.Context) error {
		if len(c.Args()) > 0 {
			commands.UnrecognizedCommand(c.Args()[0])
		} else {
			commands.TapInfoCommand().MainAction(c)
		}
		commands.PrintHelpMsg()
		return nil
	}

	//override version flag to change its shortcut from v to V (reserved for verbosity global flag)
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	return app.Run(os.Args)
}
