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
