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

package actions

import (
	"fmt"

	userManagement "github.com/trustedanalytics-ng/tap-api-service/user-management-connector"
	"github.com/trustedanalytics-ng/tap-cli/cli/printer"
)

func (a *ActionsConfig) ChangeCurrentUserPassword(currentPassword, newPassword string) error {
	if err := a.ApiService.ChangeCurrentUserPassword(currentPassword, newPassword); err != nil {
		fmt.Println("Changing user password failed")
		return err
	}
	fmt.Println("User password successfully changed.\nPlease remember to login again now.")
	return nil
}

func (a *ActionsConfig) ListUsers() error {
	users, err := a.ApiService.GetUsers()
	if err != nil {
		fmt.Println("Listing users failed")
		return err
	}
	printUsers(users)
	return nil
}

func printUsers(users []userManagement.UaaUser) {
	printableUsers := []printer.Printable{}
	for _, user := range users {
		printableUsers = append(printableUsers, printer.PrintableUser{UaaUser: user})
	}
	printer.PrintTable(printableUsers)
}

func (a *ActionsConfig) DeleteUser(email string) error {
	if err := a.ApiService.DeleteUser(email); err != nil {
		fmt.Printf("Deleting user %s failed\n", email)
		return err
	}
	fmt.Printf("User %q successfully removed\n", email)
	return nil
}
