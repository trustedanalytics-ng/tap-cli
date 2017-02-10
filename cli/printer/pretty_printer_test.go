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
	"strings"
	"testing"

	"strconv"

	"github.com/trustedanalytics-ng/tap-cli/cli/test"
)

func TestThatPrintTable_handlesEmptyList(t *testing.T) {
	stdout := test.CaptureStdout(func() {
		PrintTable([]Printable{})
	})
	assertThatContainsCaseInsensitive(t, stdout, "empty")
}

func TestThatPrintTable_printsPrintableHeaders(t *testing.T) {
	printables := createExamplaryPrintableList()
	stdout := test.CaptureStdout(func() {
		PrintTable(printables)
	})
	assertThatContainsCaseInsensitive(t, stdout, header1)
	assertThatContainsCaseInsensitive(t, stdout, header2)
}

func TestThatPrintTable_printsAllPrintableData(t *testing.T) {
	printables := createExamplaryPrintableList()
	stdout := test.CaptureStdout(func() {
		PrintTable(printables)
	})
	for vs := range values1 {
		assertThatContainsCaseInsensitive(t, stdout, values1[vs])
	}
	for vi := range values2 {
		assertThatContainsCaseInsensitive(t, stdout, strconv.Itoa(values2[vi]))
	}
}

func assertThatContainsCaseInsensitive(t *testing.T, txt string, substring string) {
	if !strings.Contains(strings.ToLower(txt), strings.ToLower(substring)) {
		t.Log(txt + " does not contain " + substring)
		t.Fail()
	}
}

const header1 = "column1 header"
const header2 = "column2 header"

var values1 = [3]string{"Ala", "ma", "kota"}
var values2 = [3]int{1, 2, 3}

type printableTestItem struct {
	Value1 string
	Value2 int
}

func (pti printableTestItem) Headers() []string {
	return []string{header1, header2}
}
func (pti printableTestItem) StandarizedData() []string {
	return []string{pti.Value1, strconv.Itoa(pti.Value2)}
}

func createExamplaryPrintableList() []Printable {
	printables := []Printable{}
	for i := 0; i < 3; i++ {
		printables = append(printables, printableTestItem{Value1: values1[i], Value2: values2[i]})
	}
	return printables
}
