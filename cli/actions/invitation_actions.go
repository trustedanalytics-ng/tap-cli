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
