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
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

type TapCommand struct {
	Name              string
	Usage             string
	Aliases           []string
	Subcommands       []TapCommand
	DefaultSubcommand *TapCommand
	MainAction        func(c *cli.Context) error
	AlternativeFlags  []cli.Flag
	OptionalFlags     []cli.Flag
	RequiredFlags     []cli.Flag
}

func (tc TapCommand) ToCliCommand() cli.Command {
	requiredFlags := tc.RequiredFlags
	optionalFlags := tc.OptionalFlags
	alternativeFlags := tc.AlternativeFlags
	if tc.DefaultSubcommand != nil {
		requiredFlags = tc.DefaultSubcommand.RequiredFlags
		optionalFlags = tc.DefaultSubcommand.OptionalFlags
		alternativeFlags = tc.DefaultSubcommand.AlternativeFlags
	}

	return cli.Command{
		Name:        tc.Name,
		HelpName:    tc.Name,
		Usage:       tc.Usage,
		Aliases:     tc.Aliases,
		Subcommands: toCommands(tc.Subcommands, tc.DefaultSubcommand),
		ArgsUsage:   getArgsUsage(requiredFlags, alternativeFlags, optionalFlags),
		Flags:       sumFlags(requiredFlags, optionalFlags, alternativeFlags, GetCommonFlags()),
		Action: func(c *cli.Context) error {
			if err := handleCommonFlags(c); err != nil {
				return err
			}
			for _, rf := range requiredFlags {
				if name, exists := checkIfRequiredFlagExists(c, rf); !exists {
					fmt.Println("\nMISSING PARAMETER: '--" + name + "'\n\nCommand usage:")
					cli.ShowCommandHelp(c, c.Command.Name)
					cli.OsExiter(requiredFlagMissingExitCode)
				}
			}

			if len(tc.AlternativeFlags) == 1 {
				printApplicationBugInfo("Only one alternative flag specified.")
				cli.OsExiter(onlyOneFlagInAlternative)
			} else if len(tc.AlternativeFlags) > 1 {
				handleAlternativeFlags(c, tc.AlternativeFlags)
			}

			if tc.DefaultSubcommand != nil {
				// A this moment there should be only command parameters. If user tried to enter
				// a command, it should be already parsed as one of the subcommands.
				if c.NArg() > 0 && !strings.HasPrefix(c.Args()[0], "--") {
					cli.ShowCommandHelp(c, tc.DefaultSubcommand.Name)
					return nil
				}
				return tc.DefaultSubcommand.MainAction(c)
			} else if tc.MainAction == nil {
				return nil
			}
			return tc.MainAction(c)
		},
	}
}

func handleAlternativeFlags(c *cli.Context, alternativeFlags []cli.Flag) {
	amount := 0
	allFlags := []string{}
	specifiedFlags := []string{}
	for _, af := range alternativeFlags {
		name, exists := checkIfRequiredFlagExists(c, af)
		if exists {
			amount++
			specifiedFlags = append(specifiedFlags, name)
		}
		allFlags = append(allFlags, name)
	}

	if amount == 0 {
		fmt.Println("\nMISSING PARAMETER. You need to specify one of alternative flags (" + strings.Join(allFlags, " OR ") + ") \n\nCommand usage:") //tc.AlternativeFlags
		cli.ShowCommandHelp(c, c.Command.Name)
		cli.OsExiter(alternativeFlagMissingExitCode)
	}
	if amount != 1 {
		fmt.Println("\nWRONG PARAMETER. Cannot use more then one alternative flags (" + strings.Join(specifiedFlags, " AND ") + ") in the same time\n\nCommand usage:") //tc.AlternativeFlags
		cli.ShowCommandHelp(c, c.Command.Name)
		cli.OsExiter(alternativeFlagTooManyExitCode)
	}
}

func getArgsUsage(required, alternative, optional []cli.Flag) string {
	argUsage := []string{}
	requiredArgUsage := argsUsageFromFlags(required, " ")
	if len(required) > 0 {
		argUsage = append(argUsage, requiredArgUsage)
	}

	alternativeArgUsage := argsUsageFromFlags(alternative, "|")
	if len(alternative) > 0 {
		argUsage = append(argUsage, alternativeArgUsage)
	}

	optionalArgUsage := "[" + argsUsageFromFlags(optional, " ") + "]"
	if len(optional) > 0 {
		argUsage = append(argUsage, optionalArgUsage)
	}

	return strings.Join(argUsage, " ")
}

func argsUsageFromFlags(flags []cli.Flag, separator string) (argsUsage string) {
	for _, f := range flags {
		argsUsage += extractFlagUse(f) + separator
	}
	return strings.TrimSuffix(argsUsage, separator)
}

func extractFlagUse(flag cli.Flag) string {
	stringifiedFlag := flag.String()
	//this returns string in a form: "--api API     TAP API you would like to use" (used in OPTIONS section)
	//I parse this string to comply to placeholders names used in this section
	splitted := strings.Split(stringifiedFlag, "\t")
	if _, ok := flag.(cli.BoolFlag); ok {
		return splitted[0]
	}
	splitted = strings.SplitN(splitted[0], " ", 2)
	return splitted[0] + "=<" + splitted[1] + ">"
}

func toCommands(tapCommands []TapCommand, defaultCommand *TapCommand) (commands []cli.Command) {
	commands = make([]cli.Command, len(tapCommands))
	for i, tc := range tapCommands {
		commands[i] = tc.ToCliCommand()
		if defaultCommand != nil && tc.Name == defaultCommand.Name {
			commands[i].HelpName = "[" + tc.Name + "]"
		}
	}
	return
}
