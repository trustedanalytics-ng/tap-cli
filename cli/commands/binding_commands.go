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
