package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type testResult int

const (
	testResultPassed testResult = iota
	testResultFailed
)

var testResultToString = map[testResult]string{
	testResultPassed: "PASSED",
	testResultFailed: "FAILED",
}

func (t testResult) String() string {
	return testResultToString[t]
}

type testSuiteReport struct {
	PathName      string
	Description   string
	Operation     string
	Status        testResult
	FailureReason string
}

func (t *testSuiteReport) Fail(reason string) {
	t.Status = testResultFailed
	t.FailureReason = reason
}

func (t *testSuiteReport) Pass() {
	t.Status = testResultPassed
	t.FailureReason = "-"
}

type Report struct {
	TestSuites []testSuiteReport
}

func (r Report) generateFailingTestDescriptions() {

}

func (r Report) generate() {

        data := [][]string{}

        for _, testSuiteReport := range r.TestSuites {
            data = append(data, []string{
                testSuiteReport.PathName,
                testSuiteReport.Operation,
                testSuiteReport.Description,
                testSuiteReport.Status.String(),
            })
        }

        table := tablewriter.NewWriter(os.Stdout)
        for i,  row := range data {

            // Row 3 is STATUS -- Refer to SetHeader function
            if row[3] == "FAILED" {

                    // The function `Rich` also appends data to tablewriter
                    // so we do not need to manually append again
                    table.Rich(data[i], []tablewriter.Colors{
                        {},
                        {},
                        {},
                        {
                            // If test is failing, we mark the cell as RED color
                            tablewriter.Bold, tablewriter.BgRedColor,
                        },
                    })

            } else {
                    table.Rich(data[i], []tablewriter.Colors{
                        {},
                        {},
                        {},
                        {
                            // If test is failing, we mark the cell as GREEN color
                            tablewriter.Bold, tablewriter.FgGreenColor,
                        },
                    })

            }

        }

        table.SetHeader([]string{"Path", "Operation", "Desc", "Status"})
        table.SetAutoMergeCellsByColumnIndex([]int{0})
        table.SetRowLine(true)
        table.Render()

}
