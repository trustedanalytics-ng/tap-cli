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

import (
	"strings"

	"github.com/urfave/cli"
)

type TapCommand struct {
	Name          string
	Usage         string
	Aliases       []string
	Subcommands   []TapCommand
	MainAction    func(c *cli.Context) error
	OptionalFlags []cli.Flag
	RequiredFlags []cli.Flag
}

func (tc TapCommand) ToCliCommand() cli.Command {
	return cli.Command{
		Name:        tc.Name,
		Usage:       tc.Usage,
		Aliases:     tc.Aliases,
		Subcommands: toCommands(tc.Subcommands),
		ArgsUsage:   getArgsUsage(tc.RequiredFlags, tc.OptionalFlags),
		Flags:       sumFlags(sumFlags(tc.RequiredFlags, tc.OptionalFlags), GetCommonFlags()),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}
			for _, rf := range tc.RequiredFlags {
				flag, ok := rf.(cli.StringFlag)
				if ok {
					checkRequiredStringFlag(flag, c)
				}
				//TODO: Add support for other types of flags
			}
			if tc.MainAction == nil {
				return nil
			}
			return tc.MainAction(c)
		},
	}
}

func getArgsUsage(required, optional []cli.Flag) (argUsageString string) {
	argUsageString = argsUsageFromFlags(required)
	if len(optional) > 0 {
		if len(required) > 0 {
			argUsageString += " "
		}
		argUsageString += "[" + argsUsageFromFlags(optional) + "]"
	}
	return
}

func argsUsageFromFlags(flags []cli.Flag) (argsUsage string) {
	for _, f := range flags {
		argsUsage += extractFlagUse(f) + " "
	}
	return strings.TrimSuffix(argsUsage, " ")
}

func extractFlagUse(flag cli.Flag) string {
	stringifiedFlag := flag.String()
	//this returns string in a form: "--api API     TAP API you would like to use" (used in OPTIONS section)
	//I parse this string to comply to placeholders names used in this section
	splitted := strings.Split(stringifiedFlag, "\t")
	splitted = strings.SplitN(splitted[0], " ", 2)
	return splitted[0] + "=<" + splitted[1] + ">"
}

func toCommands(tapCommands []TapCommand) (commands []cli.Command) {
	commands = make([]cli.Command, len(tapCommands))
	for i, tc := range tapCommands {
		commands[i] = tc.ToCliCommand()
	}
	return
}
