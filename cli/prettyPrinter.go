package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/trustedanalytics/tapng-cli/api"
	consoleServiceModels "github.com/trustedanalytics/tapng-console-service/models"
)

func createAndRenderTable(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, row := range rows {
		table.Append(row)
	}

	table.Render()
}

func printCatalog(catalog []consoleServiceModels.Service) {

	header := []string{"NAME", "PLAN", "DESCRIPTION"}
	rows := [][]string{}

	for _, service := range catalog {
		line := []string{service.Entity.Label, service.Entity.ServicePlans[0].Entity.Name, service.Entity.Description}
		rows = append(rows, line)
	}

	createAndRenderTable(header, rows)

}

func printCredentials(creds api.Credentials) {

	header := []string{"API", "USERNAME"}
	rows := [][]string{}

	rows = append(rows, []string{creds.Address, creds.Username})

	createAndRenderTable(header, rows)
}

func printApplicationInstances(applications []consoleServiceModels.ApplicationInstance) {

	header := []string{"NAME", "IMAGE STATE", "STATE", "REPLICATION"}
	rows := [][]string{}

	for _, app := range applications {
		rows = append(rows, []string{app.Name, fmt.Sprintf("%s", app.ImageState), fmt.Sprintf("%s", app.State), strconv.Itoa(app.Replication)})
	}

	createAndRenderTable(header, rows)
}

func printServices(services []consoleServiceModels.ServiceInstance) {

	header := []string{"NAME", "SERVICE", "PLAN", "STATE"}
	rows := [][]string{}

	for _, service := range services {
		rows = append(rows, []string{service.Name, service.ServiceName, service.ServicePlanName, fmt.Sprintf("%s", service.State)})
	}

	createAndRenderTable(header, rows)
}
