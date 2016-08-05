package cli

import (
	"github.com/urfave/cli"
	"strconv"
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

func CreateOfferingCommand() cli.Command {
	return cli.Command{
		Name:      "create-offering",
		Aliases:   []string{"co"},
		ArgsUsage: "<path to json with service definition>",
		Usage:     "create new offering",
		Action: func(c *cli.Context) error {
			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return CreateOffer(c.Args().First())
		},
	}
}

func CreateServiceCommand() cli.Command {
	return cli.Command{
		Name:      "create-service",
		ArgsUsage: "<service_id> <plan_id> <custom_name>",
		Aliases:   []string{"cs"},
		Usage:     "create instance of service",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 3)
			if err != nil {
				return err
			}

			return CreateInstance(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
		},
	}
}

func BindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "bind-instance",
		ArgsUsage: "<src_instance_id>, <dst_instance_id>",
		Aliases:   []string{"bind"},
		Usage:     "bind instance to another",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return BindInstance(c.Args().First(), c.Args().Get(1))
		},
	}
}

func UnbindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "unbind-instance",
		ArgsUsage: "<src_instance_id>, <dst_instance_id>",
		Aliases:   []string{"unbind"},
		Usage:     "unbind instance from another",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return UnbindInstance(c.Args().First(), c.Args().Get(1))
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
		ArgsUsage: "<archive_path>",
		Usage:     "create application from archive, manifest should be in current working directory",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return PushApplication(c.Args().First())
		},
	}
}

func ListServicesCommand() cli.Command {
	return cli.Command{
		Name:      "services",
		ArgsUsage: "",
		Aliases:   []string{"s"},
		Usage:     "list services",
		Action: func(c *cli.Context) error {
			return ListServices()
		},
	}
}

func ScaleApplicationCommand() cli.Command {
	return cli.Command{
		Name:      "scale",
		ArgsUsage: "<instanceId> <instances>",
		Aliases:   []string{"sc"},
		Usage:     "scale application",
		Action: func(c *cli.Context) error {

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			i, errr := strconv.Atoi(c.Args().Get(1))
			if errr != nil {
				return cli.NewExitError(errr.Error(), -1)
			}

			return ScaleApplication(c.Args().First(), i)
		},
	}
}
