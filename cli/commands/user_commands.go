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

func listUsersCommand() cli.Command {
	return cli.Command{
		Name:  "users",
		Usage: "list platform users",
		Flags: GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListUsers()
		},
	}
}

func deleteUserCommand() cli.Command {
	return cli.Command{
		Name:      "delete-user",
		Usage:     "delete user from TAP",
		Aliases:   []string{"du"},
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

			return newOAuth2Service().DeleteUser(c.Args().First())
		},
	}
}

func changeCurrentUserPasswordCommand() cli.Command {
	return cli.Command{
		Name:      "chpasswd",
		Usage:     "change password of currently logged user",
		ArgsUsage: "<currentPassword> <newPassword>",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return newOAuth2Service().ChangeCurrentUserPassword(c.Args().First(), c.Args().Get(1))
		},
	}
}
