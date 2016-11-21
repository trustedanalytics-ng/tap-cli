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

package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	consoleServiceModels "github.com/trustedanalytics/tap-api-service/models"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
)

const timeFormatter = "Jan 02 15:04"
const LastMessageMark = "..."

func PrintCatalog(catalog []consoleServiceModels.Offering) {
	header := []string{"NAME", "PLAN", "DESCRIPTION", "STATE"}
	rows := [][]string{}

	for _, offering := range catalog {
		planNames := []string{}
		for _, planName := range offering.OfferingPlans {
			planNames = append(planNames, planName.Name)
		}
		line := []string{offering.Name, strings.Join(planNames, ", "), offering.Description, offering.State}
		rows = append(rows, line)
	}

	createAndRenderTable(header, rows)

}

func PrintCredentials(creds api.Credentials) {

	header := []string{"API", "USERNAME"}
	rows := [][]string{}

	rows = append(rows, []string{creds.Address, creds.Username})

	createAndRenderTable(header, rows)
}

func PrintApplicationInstances(applications []consoleServiceModels.ApplicationInstance) {

	header := []string{"NAME", "IMAGE STATE", "STATE", "REPLICATION", "MEMORY", "DISK", "URLS", "CREATED BY", "CREATE", "UPDATED BY", "UPDATE", "MESSAGE"}
	rows := [][]string{}

	for _, app := range applications {
		rows = append(rows, []string{
			app.Name, fmt.Sprintf("%s", app.ImageState), fmt.Sprintf("%s", app.State),
			strconv.Itoa(app.Replication), app.Memory, app.DiskQuota, strings.Join(app.Urls, ","),
			app.AuditTrail.CreatedBy, time.Unix(app.AuditTrail.CreatedOn, 0).Format(timeFormatter),
			app.AuditTrail.LastUpdateBy, time.Unix(app.AuditTrail.LastUpdatedOn, 0).Format(timeFormatter),
			getLastMessageMark(app.Metadata),
		})
	}

	createAndRenderTable(header, rows)
}

func PrintApplication(applications []catalogModels.Application) {

	header := []string{"NAME", "IMAGE ID", "DESCRIPTION", "REPLICATION", "CREATED BY", "CREATE", "UPDATED BY", "UPDATE"}
	rows := [][]string{}

	for _, app := range applications {
		rows = append(rows, []string{app.Name, fmt.Sprintf("%s", app.ImageId), fmt.Sprintf("%s", app.Description),
			strconv.Itoa(app.Replication), app.AuditTrail.CreatedBy, time.Unix(app.AuditTrail.CreatedOn, 0).Format(timeFormatter),
			app.AuditTrail.LastUpdateBy, time.Unix(app.AuditTrail.LastUpdatedOn, 0).Format(timeFormatter)})
	}

	createAndRenderTable(header, rows)
}

func PrintServices(services []consoleServiceModels.ServiceInstance) {

	header := []string{"NAME", "SERVICE", "PLAN", "STATE", "CREATED BY", "CREATE", "UPDATED BY", "UPDATE", "MESSAGE"}
	rows := [][]string{}

	for _, service := range services {
		rows = append(rows, []string{
			service.Name, service.ServiceName, service.ServicePlanName, fmt.Sprintf("%s", service.State),
			service.AuditTrail.CreatedBy, time.Unix(service.AuditTrail.CreatedOn, 0).Format(timeFormatter),
			service.AuditTrail.LastUpdateBy, time.Unix(service.AuditTrail.LastUpdatedOn, 0).Format(timeFormatter),
			getLastMessageMark(service.Metadata),
		})
	}

	createAndRenderTable(header, rows)
}

func PrintFormattedDetails(instance interface{}) {
	prettyJSON, err := json.MarshalIndent(instance, "", "    ")
	if err == nil {
		fmt.Print(string(prettyJSON))
	} else {
		fmt.Print("Error marshaling JSON")
	}
}

func PrintInstancesBindings(bindings consoleServiceModels.InstanceBindings) {

	header := []string{"BINDING NAME", "BINDING ID"}
	rows := [][]string{}

	for _, resource := range bindings.Resources {
		rows = append(rows, []string{resource.ServiceInstanceName, resource.ServiceInstanceGuid})
	}

	createAndRenderTable(header, rows)
}

func createAndRenderTable(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(rows)
	table.Render()
}

func getLastMessageMark(metadata []catalogModels.Metadata) string {
	if catalogModels.GetValueFromMetadata(metadata, catalogModels.LAST_STATE_CHANGE_REASON) != "" {
		return LastMessageMark
	}
	return ""
}
