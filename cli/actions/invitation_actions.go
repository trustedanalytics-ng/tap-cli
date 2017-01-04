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

	"github.com/trustedanalytics/tap-cli/cli/printer"
)

func (a *ActionsConfig) SendInvitation(email string) error {
	if _, err := a.ApiService.SendInvitation(email); err != nil {
		fmt.Printf("Sending invitation to email %s failed\n", email)
		return err
	}
	fmt.Printf("User %q successfully invited\n", email)
	return nil
}

func (a *ActionsConfig) ResendInvitation(email string) error {
	if err := a.ApiService.ResendInvitation(email); err != nil {
		fmt.Printf("Resending invitation to email %s failed\n", email)
		return err
	}
	fmt.Printf("User %q successfully invited\n", email)
	return nil
}

func (a *ActionsConfig) ListInvitations() error {
	invitations, err := a.ApiService.GetInvitations()
	if err != nil {
		fmt.Println("Listing invitations failed")
		return err
	}
	printInvitations(invitations)
	return nil
}

func printInvitations(invitations []string) {
	printableInvitations := []printer.Printable{}
	for _, inv := range invitations {
		printableInvitations = append(printableInvitations, printer.PrintableInvitation{Email: inv})
	}
	printer.PrintTable(printableInvitations)
}

func (a *ActionsConfig) DeleteInvitation(email string) error {
	if err := a.ApiService.DeleteInvitation(email); err != nil {
		fmt.Printf("Deleting invitation of user %s failed\n", email)
		return err
	}
	fmt.Printf("Invitation for user %q successfully removed\n", email)
	return nil
}
