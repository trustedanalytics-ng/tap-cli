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
