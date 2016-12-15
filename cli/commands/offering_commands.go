package commands

import "github.com/urfave/cli"

func listOfferings() cli.Command {
	return cli.Command{
		Name:    "catalog",
		Aliases: []string{"o"},
		Usage:   "list available offerings",
		Flags:   GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			return newOAuth2Service().ListOfferings()
		},
	}
}

func createOfferingCommand() cli.Command {
	return cli.Command{
		Name:      "create-offering",
		Aliases:   []string{"co"},
		ArgsUsage: "<path to json with service definition>",
		Usage:     "create new offering",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().CreateOffering(c.Args().First())
		},
	}
}

func deleteOfferingCommand() cli.Command {
	return cli.Command{
		Name:      "delete-offering",
		ArgsUsage: "<offering_custom_name>",
		Aliases:   []string{"do"},
		Usage:     "delete offering",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().DeleteOffering(c.Args().First())
		},
	}
}
