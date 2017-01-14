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
	"errors"
	"flag"
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli"
)

func TestExtractFlagUse(t *testing.T) {
	Convey("Test ExtractFlagUse", t, func() {
		Convey("Should create ArgUse with proper placeholder", func() {
			testFlag := cli.StringFlag{
				Name:  "test-flag-name",
				Usage: "test flag usage with `PLACEHOLDER with Spaces`",
			}
			flagUse := extractFlagUse(testFlag)

			So(flagUse, ShouldEqual, "--test-flag-name=<PLACEHOLDER with Spaces>")
		})

		Convey("Should create ArgUse with default placeholder", func() {
			testFlag := cli.StringFlag{
				Name:  "test-flag-name",
				Usage: "test flag usage without PLACEHOLDER",
			}
			flagUse := extractFlagUse(testFlag)

			So(flagUse, ShouldEqual, "--test-flag-name=<value>")
		})
	})
}

func TestGetArgsUsage(t *testing.T) {
	Convey("Test GetArgsUsage", t, func() {
		requiredFlag1 := cli.StringFlag{
			Name:  "required-flag-name-1",
			Usage: "required `flag` usage ",
		}
		requiredFlag1ArgUse := "--required-flag-name-1=<flag>"
		requiredFlag2 := cli.StringFlag{
			Name:  "required-flag-name-2",
			Usage: "required `flag` usage",
		}
		requiredFlag2ArgUse := "--required-flag-name-2=<flag>"
		optionalFlag := cli.StringFlag{
			Name:  "optional-flag-name",
			Usage: "optional `flag` usage",
		}
		optionalFlagArgUse := "--optional-flag-name=<flag>"

		Convey("Should create proper ArgUse", func() {
			testCommand := TapCommand{
				OptionalFlags: []cli.Flag{optionalFlag},
				RequiredFlags: []cli.Flag{requiredFlag1, requiredFlag2},
			}.ToCliCommand()

			So(testCommand.ArgsUsage, ShouldEqual, requiredFlag1ArgUse+" "+requiredFlag2ArgUse+" ["+optionalFlagArgUse+"]")
		})

		Convey("Should create ArgUse without optional part, when no optional flags specified", func() {
			testCommand := TapCommand{
				RequiredFlags: []cli.Flag{requiredFlag1, requiredFlag2},
			}.ToCliCommand()

			So(testCommand.ArgsUsage, ShouldEqual, requiredFlag1ArgUse+" "+requiredFlag2ArgUse)
		})

		Convey("Should create ArgUse without required part, when no required flags specified", func() {
			testCommand := TapCommand{
				OptionalFlags: []cli.Flag{optionalFlag},
			}.ToCliCommand()

			So(testCommand.ArgsUsage, ShouldEqual, "["+optionalFlagArgUse+"]")
		})
	})
}

func TestToCliCommand(t *testing.T) {
	Convey("Test ToCliCommand", t, func() {
		Convey("Should create command with proper name", func() {
			test_command_name := "test command name"
			testCommand := TapCommand{
				Name: test_command_name,
			}.ToCliCommand()

			So(testCommand.Name, ShouldEqual, test_command_name)
		})

		Convey("Should create command with proper usage", func() {
			test_command_usage := "test command usage"
			testCommand := TapCommand{
				Usage: test_command_usage,
			}.ToCliCommand()

			So(testCommand.Usage, ShouldEqual, test_command_usage)
		})

		Convey("Should create command with proper aliases", func() {
			alias1 := "alias1"
			alias2 := "alias2"
			testCommand := TapCommand{
				Aliases: []string{alias1, alias2},
			}.ToCliCommand()

			So(testCommand.Aliases, ShouldContain, alias1)
			So(testCommand.Aliases, ShouldContain, alias2)
		})

		Convey("Should create command with proper subcomands", func() {
			testSubCommand1 := TapCommand{
				Name: "subcomand1",
			}
			testSubCommand2 := TapCommand{
				Name: "subcomand2",
			}
			testCommand := TapCommand{
				Subcommands: []TapCommand{testSubCommand1, testSubCommand2},
			}.ToCliCommand()

			So(testCommand.Subcommands[0], ShouldCliCommandResemble, testSubCommand1.ToCliCommand())
			So(testCommand.Subcommands[1], ShouldCliCommandResemble, testSubCommand2.ToCliCommand())
		})

		Convey("Should create command with proper flags and argUse (required and optional)", func() {
			requiredFlag1 := cli.StringFlag{
				Name:  "required-flag-name-1",
				Usage: "required `flag` usage ",
			}
			requiredFlag1ArgUse := "--required-flag-name-1=<flag>"
			requiredFlag2 := cli.StringFlag{
				Name:  "required-flag-name-2",
				Usage: "required `flag` usage",
			}
			requiredFlag2ArgUse := "--required-flag-name-2=<flag>"
			optionalFlag := cli.StringFlag{
				Name:  "optional-flag-name",
				Usage: "optional `flag` usage",
			}
			optionalFlagArgUse := "--optional-flag-name=<flag>"

			testCommand := TapCommand{
				OptionalFlags: []cli.Flag{optionalFlag},
				RequiredFlags: []cli.Flag{requiredFlag1, requiredFlag2},
			}.ToCliCommand()

			So(testCommand.ArgsUsage, ShouldEqual, requiredFlag1ArgUse+" "+requiredFlag2ArgUse+" ["+optionalFlagArgUse+"]")
			So(testCommand.Flags, ShouldContain, requiredFlag1)
			So(testCommand.Flags, ShouldContain, requiredFlag2)
			So(testCommand.Flags, ShouldContain, optionalFlag)
		})

		Convey("Should run mainAction", func() {
			errorMessage := "Fake Error for Test purpose"
			testCommand := TapCommand{
				MainAction: func(c *cli.Context) error {
					return errors.New(errorMessage)
				},
			}.ToCliCommand()

			err := executeCommandAction(testCommand)

			So(err.Error(), ShouldEqual, errorMessage)
		})

		Convey("Should check String flag value and", func() {
			cli.OsExiter = mockExiter
			var flagDestination string
			testFlag := cli.StringFlag{
				Name:        "test-string-flag-name",
				Destination: &flagDestination,
			}
			Convey("throw error if value is unspecified", func() {
				testCommand := TapCommand{
					Name:          "test-command-name",
					RequiredFlags: []cli.Flag{testFlag},
				}.ToCliCommand()

				So(func() { executeCommandActionWithFlag(testCommand, testFlag) }, ShouldPanicWith, requiredFlagMissingExitCode)
			})

			Convey("work OK if default Value is specified", func() {
				testFlag.Value = "default specified"

				testCommand := TapCommand{
					RequiredFlags: []cli.Flag{testFlag},
				}.ToCliCommand()

				So(func() { executeCommandActionWithFlag(testCommand, testFlag) }, ShouldNotPanic)
			})

			Convey("work OK for optional flags (no checking)", func() {
				testCommand := TapCommand{
					OptionalFlags: []cli.Flag{testFlag},
				}.ToCliCommand()

				So(func() { executeCommandActionWithFlag(testCommand, testFlag) }, ShouldNotPanic)
			})

			Reset(func() {
				cli.OsExiter = os.Exit
			})
		})

		Convey("Should check int flag value and", func() {
			cli.OsExiter = mockExiter
			var flagDestination int
			testFlag := cli.IntFlag{
				Name:        "test-int-flag-name",
				Destination: &flagDestination,
			}

			Convey("work OK if default Value is specified", func() {
				testFlag.Value = 5

				testCommand := TapCommand{
					RequiredFlags: []cli.Flag{testFlag},
				}.ToCliCommand()

				So(func() { executeCommandActionWithFlag(testCommand, testFlag) }, ShouldNotPanic)
			})

			Convey("throw error if default Value is set to 0 (unsupported case)", func() {
				testFlag.Value = 0

				testCommand := TapCommand{
					RequiredFlags: []cli.Flag{testFlag},
				}.ToCliCommand()

				So(func() { executeCommandActionWithFlag(testCommand, testFlag) }, ShouldPanicWith, requiredFlagMissingExitCode)
			})

			Reset(func() {
				cli.OsExiter = os.Exit
			})
		})
	})
}

func ShouldCliCommandResemble(actual interface{}, expected ...interface{}) string {
	act, ok := actual.(cli.Command)
	if ok == false {
		return "Not a CliCommand"
	}

	exp, ok := expected[0].(cli.Command)
	if ok == false {
		return "Not a CliCommand"
	}

	//ShouldResemble cannot compare pointer to methods if they ar not nil
	// -> as it is in code of deepValueEqual in reflect package: "Can't do better than this:"
	act_ptr := reflect.ValueOf(act.Action).Pointer()
	exp_ptr := reflect.ValueOf(exp.Action).Pointer()
	if act_ptr != exp_ptr {
		return "Different Actions"
	}
	act.Action = nil
	exp.Action = nil

	return ShouldResemble(act, exp)
}

func mockExiter(code int) {
	panic(code)
}

func executeCommandAction(command cli.Command) error {
	return executeCommandActionWithFlagSet(command, nil)
}

func executeCommandActionWithFlag(command cli.Command, testFlag cli.Flag) error {
	set := flag.FlagSet{}
	testFlag.Apply(&set)

	return executeCommandActionWithFlagSet(command, &set)
}

func executeCommandActionWithFlagSet(command cli.Command, set *flag.FlagSet) error {
	fun, _ := command.Action.(func(c *cli.Context) error)

	app := cli.NewApp()
	ctx := cli.NewContext(app, set, nil)
	ctx.Command = command
	loggerVerbosity = DefaultLogLevel
	return fun(ctx)
}
