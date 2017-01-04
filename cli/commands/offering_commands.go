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

import "github.com/urfave/cli"

//catalog -> offering
func offeringCommands() TapCommand {
	return TapCommand{
		Name:  "offering",
		Usage: "offering context commands",
		MainAction: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
		Subcommands: []TapCommand{
			infoOfferingCommand(),
			listOfferingsCommand(),
			createOfferingCommand(),
			deleteOfferingCommand(),
		},
	}
}

func infoOfferingCommand() TapCommand {
	var name string
	var nameFlag = cli.StringFlag{
		Name:        "name",
		Usage:       "`name of offering` you would like to display",
		Destination: &name,
	}

	return TapCommand{
		Name:          "info",
		Usage:         "show information about specific offering",
		RequiredFlags: []cli.Flag{nameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetOffering(name)
		},
	}
}

func listOfferingsCommand() TapCommand {
	return TapCommand{
		Name:  "list",
		Usage: "list available offerings",
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ListOfferings()
		},
	}
}

func createOfferingCommand() TapCommand {
	var manifestPath string
	var manifestFlag = cli.StringFlag{
		Name:        "manifest",
		Usage:       "`path to json` with service definition",
		Value:       "manifest.json",
		Destination: &manifestPath,
	}

	return TapCommand{
		Name:          "create",
		Usage:         "create new offering",
		RequiredFlags: []cli.Flag{manifestFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().CreateOffering(manifestPath)
		},
	}
}

func deleteOfferingCommand() TapCommand {
	var name string
	var nameFlag = cli.StringFlag{
		Name:        "name",
		Usage:       "`name of offering` you would like to delete",
		Destination: &name,
	}

	return TapCommand{
		Name:          "delete",
		Usage:         "delete offering",
		RequiredFlags: []cli.Flag{nameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().DeleteOffering(name)
		},
	}
}
