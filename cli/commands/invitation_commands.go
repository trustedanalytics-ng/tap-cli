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

func listInvitationsCommand() cli.Command {
	return cli.Command{
		Name:    "invitations",
		Usage:   "list pending invitations",
		Aliases: []string{"invs"},
		Flags:   GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListInvitations()
		},
	}
}

func sendInvitationCommand() cli.Command {
	return cli.Command{
		Name:      "invite",
		Usage:     "invite new user to TAP",
		ArgsUsage: "<email>",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().SendInvitation(c.Args().First())
		},
	}
}

func deleteInvitationCommand() cli.Command {
	return cli.Command{
		Name:      "delete-invitation",
		Usage:     "delete invitation",
		Aliases:   []string{"di"},
		ArgsUsage: "<email>",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().DeleteInvitation(c.Args().First())
		},
	}
}

func resendInvitationCommand() cli.Command {
	return cli.Command{
		Name:      "reinvite",
		Usage:     "resend invitation for user",
		ArgsUsage: "<email>",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().ResendInvitation(c.Args().First())
		},
	}
}
