package cli

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	catalogModels "github.com/trustedanalytics/tapng-catalog/models"
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

func printApplications(applications []catalogModels.Application) {

	header := []string{"APP_ID", "IMAGE_ID", "TEMPLATE_ID", "REPLICATION"}
	rows := [][]string{}

	for _, app := range applications {
		rows = append(rows, []string{app.Id, app.ImageId, app.TemplateId, fmt.Sprintf("%d", app.Replication)})
	}

	createAndRenderTable(header, rows)
}

func printAppInstance(instance catalogModels.Instance, replication int) {

	header := []string{"INSTANCE_NAME", "INSTANCE_ID", "REPLICATION"}
	rows := [][]string{
		[]string{instance.Name, instance.Id, fmt.Sprintf("%d", replication)},
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
