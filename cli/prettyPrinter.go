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

package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
		planNames := []string{}
		for _, planName := range service.Entity.ServicePlans {
			planNames = append(planNames, planName.Entity.Name)
		}
		line := []string{service.Entity.Label, strings.Join(planNames, ", "), service.Entity.Description}
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

	header := []string{"NAME", "IMAGE STATE", "STATE", "REPLICATION", "MEMORY", "DISK", "URLS"}
	rows := [][]string{}

	for _, app := range applications {
		rows = append(rows, []string{app.Name, fmt.Sprintf("%s", app.ImageState), fmt.Sprintf("%s", app.State),
			strconv.Itoa(app.Replication), app.Memory, app.DiskQuota, strings.Join(app.Urls, ",")})
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

func printInstancesBindings(bindings consoleServiceModels.InstanceBindings) {

	header := []string{"BINDING NAME", "BINDING ID"}
	rows := [][]string{}

	for _, resource := range bindings.Resources {
		rows = append(rows, []string{resource.ServiceInstanceName, resource.ServiceInstanceGuid})
	}

	createAndRenderTable(header, rows)
}
