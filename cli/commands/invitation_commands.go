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

func invitationsCommand() TapCommand {
	var email string
	var emailFlag = cli.StringFlag{
		Name:        "email",
		Usage:       "`user email`",
		Destination: &email,
	}

	confirmed := false
	var confirmationFlag = cli.BoolFlag{
		Name:        "yes",
		Usage:       "use with caution when want to suppress removal confirmation",
		Destination: &confirmed,
	}

	var listInvitationsCommand = TapCommand{
		Name:  "list",
		Usage: "list pending invitations",
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ListInvitations()
		},
	}

	var sendInvitationCommand = TapCommand{
		Name:          "send",
		Usage:         "invite new user to TAP",
		RequiredFlags: []cli.Flag{emailFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().SendInvitation(email)
		},
	}

	var deleteInvitationCommand = TapCommand{
		Name:          "delete",
		Usage:         "delete invitation for given `email`",
		RequiredFlags: []cli.Flag{emailFlag},
		OptionalFlags: []cli.Flag{confirmationFlag},
		MainAction: func(c *cli.Context) error {
			if !confirmed {
				err := removalConfirmationPrompt("invitation for " + email)
				cli.HandleExitCoder(err)
			}
			return newOAuth2Service().DeleteInvitation(email)
		},
	}

	var resendInvitationCommand = TapCommand{
		Name:          "resend",
		Usage:         "resend invitation for user",
		RequiredFlags: []cli.Flag{emailFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().ResendInvitation(email)
		},
	}

	return TapCommand{
		Name:  "invitation",
		Usage: "user invitation context commands",
		Subcommands: []TapCommand{
			listInvitationsCommand,
			sendInvitationCommand,
			resendInvitationCommand,
			deleteInvitationCommand,
		},
		DefaultSubcommand: &listInvitationsCommand,
	}
}
