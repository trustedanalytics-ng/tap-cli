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

func userCommand() TapCommand {
	var username string
	var usernameFlag = cli.StringFlag{
		Name:        "name",
		Usage:       "`username`",
		Destination: &username,
	}

	var currentPass string
	var currentPassFlag = cli.StringFlag{
		Name:        "current",
		Usage:       "`current password`",
		Destination: &currentPass,
	}

	var newPass string
	var newPassFlag = cli.StringFlag{
		Name:        "new",
		Usage:       "`new password`",
		Destination: &newPass,
	}

	confirmed := false
	var confirmationFlag = cli.BoolFlag{
		Name:        "yes",
		Usage:       "use with caution when want to suppress removal confirmation",
		Destination: &confirmed,
	}

	var listUsersCommand = TapCommand{
		Name:  "list",
		Usage: "list platform users",
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ListUsers()
		},
	}

	var deleteUserCommand = TapCommand{
		Name:          "delete",
		Usage:         "delete user from TAP",
		OptionalFlags: []cli.Flag{confirmationFlag},
		RequiredFlags: []cli.Flag{usernameFlag},
		MainAction: func(c *cli.Context) error {
			if !confirmed {
				err := removalConfirmationPrompt("user " + username)
				cli.HandleExitCoder(err)
			}
			return newOAuth2Service().DeleteUser(username)
		},
	}

	var changeCurrentUserPasswordCommand = TapCommand{
		Name:          "passwd",
		Usage:         "change password of currently logged user",
		OptionalFlags: []cli.Flag{currentPassFlag, newPassFlag},
		MainAction: func(c *cli.Context) error {
			if currentPass == "" {
				currentPass = promptForSensitive("Current Password")
			}
			if newPass == "" {
				newPass = promptForSensitive("New Password")
			}
			return newOAuth2Service().ChangeCurrentUserPassword(currentPass, newPass)
		},
	}

	return TapCommand{
		Name:  "user",
		Usage: "user context commands",
		MainAction: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
		Subcommands: []TapCommand{
			listUsersCommand,
			deleteUserCommand,
			changeCurrentUserPasswordCommand,
			invitationsCommand(),
		},
	}
}
