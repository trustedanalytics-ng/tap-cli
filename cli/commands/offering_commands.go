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

func offeringCommand() TapCommand {
	var name string
	var nameFlag = cli.StringFlag{
		Name:        "name",
		Usage:       "`name of offering`",
		Destination: &name,
	}

	var manifestPath string
	var manifestFlag = cli.StringFlag{
		Name:        "manifest",
		Usage:       "`path to json` with service definition",
		Value:       "manifest.json",
		Destination: &manifestPath,
	}

	confirmed := false
	var confirmationFlag = cli.BoolFlag{
		Name:        "yes",
		Usage:       "use with caution when want to suppress removal confirmation",
		Destination: &confirmed,
	}

	var infoOfferingCommand = TapCommand{
		Name:          "info",
		Usage:         "show information about specific offering",
		RequiredFlags: []cli.Flag{nameFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetOffering(name)
		},
	}

	var listOfferingsCommand = TapCommand{
		Name:  "list",
		Usage: "list available offerings",
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ListOfferings()
		},
	}

	var createOfferingCommand = TapCommand{
		Name:          "create",
		Usage:         "create new offering",
		RequiredFlags: []cli.Flag{manifestFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().CreateOffering(manifestPath)
		},
	}

	var deleteOfferingCommand = TapCommand{
		Name:          "delete",
		Usage:         "delete offering",
		RequiredFlags: []cli.Flag{nameFlag},
		OptionalFlags: []cli.Flag{confirmationFlag},
		MainAction: func(c *cli.Context) error {
			if !confirmed {
				err := removalConfirmationPrompt("offering " + name)
				cli.HandleExitCoder(err)
			}
			return newOAuth2Service().DeleteOffering(name)
		},
	}

	return TapCommand{
		Name:  "offering",
		Usage: "offering context commands",
		Subcommands: []TapCommand{
			infoOfferingCommand,
			listOfferingsCommand,
			createOfferingCommand,
			deleteOfferingCommand,
		},
		DefaultSubcommand: &infoOfferingCommand,
	}
}
