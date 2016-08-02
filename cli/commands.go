package cli

import (
	"github.com/urfave/cli"
)

func validateArgs(c *cli.Context, mustCount int) *cli.ExitError {
	if c.NArg() != mustCount {
		return cli.NewExitError("not enough args: \n"+c.Command.Name+" "+c.Command.ArgsUsage, 1)
	}
	return nil
}

func LoginCommand() cli.Command {
	return cli.Command{
		Name:      "login",
		Usage:     "login to TAPNG",
		ArgsUsage: "<address> <username> <password>",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 3)
			if err != nil {
				return err
			}

			return Login(c.Args().First(), c.Args().Get(1), c.Args().Get(2))
		},
	}
}

func TargetCommand() cli.Command {
	return cli.Command{
		Name:  "target",
		Usage: "print actual credentials",
		Action: func(c *cli.Context) error {
			return Target()
		},
	}
}

func CatalogCommand() cli.Command {
	return cli.Command{
		Name:  "catalog",
		Usage: "list available services",
		Action: func(c *cli.Context) error {
			return Catalog()
		},
	}
}

func DeployCommand() cli.Command {
	return cli.Command{
		Name:      "deploy",
		ArgsUsage: "<path to json with service definition>",
		Usage:     "create new service",
		Action: func(c *cli.Context) error {
			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return Deploy(c.Args().First())
		},
	}
}

func CreateServiceCommand() cli.Command {
	return cli.Command{
		Name:      "create-service",
		ArgsUsage: "<service_id>",
		Aliases:   []string{"cs"},
		Usage:     "create instance of service",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return CreateInstance(c.Args().First())
		},
	}
}

func ListApplicationsCommand() cli.Command {
	return cli.Command{
		Name:      "applications",
		ArgsUsage: "",
		Aliases:   []string{"apps"},
		Usage:     "list applications",
		Action: func(c *cli.Context) error {
			return ListApplications()
		},
	}
}

func PushApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "push",
		ArgsUsage: "<archive_path> <image_json_path> <template_json_path>",
		Usage:     "create application from archive",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 3)
			if err != nil {
				return err
			}

			return PushApplication(c.Args().First(), c.Args().Get(1), c.Args().Get(2))
		},
	}
}
