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
	"strings"

	"github.com/urfave/cli"

	catalogModels "github.com/trustedanalytics-ng/tap-catalog/models"
	"github.com/trustedanalytics-ng/tap-cli/cli/actions"
	"github.com/trustedanalytics-ng/tap-cli/cli/converter"
)

type bindingFunction func(srcInstance, dstInstance actions.BindableInstance) error

func bindingCommands(instanceType catalogModels.InstanceType) TapCommand {
	var instanceTypeString = strings.ToLower(string(instanceType))

	var name string
	var nameFlag = cli.StringFlag{
		Name:        "name",
		Usage:       "`name` of " + instanceTypeString,
		Destination: &name,
	}
	var isDst bool
	var isDstFlag = cli.BoolFlag{
		Name:        "is-dst",
		Usage:       "Show bindings for which, current " + instanceTypeString + " is a destination",
		Destination: &isDst,
	}
	var dstName string
	var dstFlag = cli.StringFlag{
		Name:        "dst-name",
		Usage:       "name of `destination instance`",
		Destination: &dstName,
	}
	var srcName string
	var srcFlag = cli.StringFlag{
		Name:        "src-name",
		Usage:       "name of `source instance`",
		Destination: &srcName,
	}

	var listBindingCommand = TapCommand{
		Name:          "list",
		Usage:         "list bindings",
		RequiredFlags: []cli.Flag{nameFlag, isDstFlag},
		MainAction: func(c *cli.Context) error {
			return newOAuth2Service().GetInstanceBindings(
				actions.BindableInstance{Name: name, Type: instanceType})
		},
	}

	var createBindingCommand = TapCommand{
		Name:             "create",
		Usage:            "bind " + instanceTypeString + " to any instance",
		RequiredFlags:    []cli.Flag{nameFlag},
		AlternativeFlags: []cli.Flag{dstFlag, srcFlag},
		MainAction: func(c *cli.Context) error {
			return binding(name, dstName, srcName, instanceType, newOAuth2Service().BindInstance)
		},
	}

	var deleteBindingCommand = TapCommand{
		Name:             "delete",
		Usage:            "unbind " + instanceTypeString + " from any instance",
		RequiredFlags:    []cli.Flag{nameFlag},
		AlternativeFlags: []cli.Flag{dstFlag, srcFlag},
		MainAction: func(c *cli.Context) error {
			return binding(name, dstName, srcName, instanceType, newOAuth2Service().UnbindInstance)
		},
	}

	return TapCommand{
		Name:  "binding",
		Usage: "binding context commands",
		MainAction: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
		Subcommands: []TapCommand{
			listBindingCommand,
			createBindingCommand,
			deleteBindingCommand,
		},
	}
}

func binding(name, dstName, srcName string, instanceType catalogModels.InstanceType, bindingFunc bindingFunction) error {
	if dstName != "" {
		return bindingFunc(
			actions.BindableInstance{Name: name, Type: instanceType},
			actions.BindableInstance{Name: dstName, Type: converter.InstanceTypeBoth})
	} else if srcName != "" {
		return bindingFunc(
			actions.BindableInstance{Name: srcName, Type: converter.InstanceTypeBoth},
			actions.BindableInstance{Name: name, Type: instanceType})
	}
	return errors.New("dstName and srcName cannot be empty at the same time. Verification for alternative flags probably failed.")
}
