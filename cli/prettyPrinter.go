package cli

import (
	consoleServiceModels "github.com/trustedanalytics/tapng-console-service/models"
	"github.com/olekukonko/tablewriter"
	"os"
	"github.com/trustedanalytics/tapng-cli/api"
)

func createAndRenderTable(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _,row := range rows {
		table.Append(row)
	}

	table.Render()
}

func printCatalog(catalog []consoleServiceModels.Service) {

	header := []string{"NAME","PLAN","DESCRIPTION"}
	rows := [][]string{}

	for _, service := range catalog {
		line := []string{service.Entity.Label, service.Entity.ServicePlans[0].Entity.Name, service.Entity.Description}
		rows = append(rows, line)
	}

	createAndRenderTable(header, rows)

}

func printCredentials(creds api.Credentials) {

	header := []string{"API","USERNAME"}
	rows := [][]string{}

	rows = append(rows, []string{creds.Address, creds.Username})

	createAndRenderTable(header, rows)
}