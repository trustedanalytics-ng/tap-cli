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

	"github.com/olekukonko/tablewriter"
)

const emptyListMsg = "(empty list)"

func PrintFormattedJSON(instance interface{}) {
	prettyJSON, err := json.MarshalIndent(instance, "", "    ")
	if err == nil {
		fmt.Print(string(prettyJSON))
	} else {
		fmt.Print("Error marshaling JSON")
	}
}

func PrintTable(items []Printable) {
	rows := [][]string{}
	if len(items) < 1 {
		createAndRenderTable(nil, append(rows, []string{emptyListMsg}))
		return
	}
	header := items[0].Headers()
	for _, i := range items {
		rows = append(rows, i.StandarizedData())
	}
	createAndRenderTable(header, rows)
}

func createAndRenderTable(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(rows)
	table.Render()
}
