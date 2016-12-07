package actions

import (
	"fmt"

	userManagement "github.com/trustedanalytics/tap-api-service/user-management-connector"
	"github.com/trustedanalytics/tap-cli/cli/printer"
)

func (a *ActionsConfig) ChangeCurrentUserPassword(currentPassword, newPassword string) error {
	if err := a.ApiService.ChangeCurrentUserPassword(currentPassword, newPassword); err != nil {
		fmt.Println("Changing user password failed")
		return err
	}
	fmt.Println("User password successfully changed")
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
