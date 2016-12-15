package commands

import "github.com/urfave/cli"

func listInstanceBindingsCommand() cli.Command {
	return cli.Command{
		Name:      "bindings",
		ArgsUsage: "<instanceName>",
		Usage:     "list bindings",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 1)
			if err != nil {
				return err
			}

			return newOAuth2Service().GetInstanceBindings(c.Args().First())
		},
	}
}

func bindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "bind-instance",
		ArgsUsage: "<src_instance_name>, <dst_instance_name>",
		Aliases:   []string{"bind"},
		Usage:     "bind instance to another",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return newOAuth2Service().BindInstance(c.Args().First(), c.Args().Get(1))
		},
	}
}

func unbindInstanceCommand() cli.Command {
	return cli.Command{
		Name:      "unbind-instance",
		ArgsUsage: "<src_instance_name>, <dst_instance_name>",
		Aliases:   []string{"unbind"},
		Usage:     "unbind instance from another",
		Flags:     GetCommonFlags(),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}

			err := validateArgs(c, 2)
			if err != nil {
				return err
			}

			return newOAuth2Service().UnbindInstance(c.Args().First(), c.Args().Get(1))
		},
	}
}
