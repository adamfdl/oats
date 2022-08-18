package main

import (
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

func (r Report) generate() {

        data := [][]string{}

        for _, testSuiteReport := range r.TestSuites {
            data = append(data, []string{
                testSuiteReport.PathName,
                testSuiteReport.Operation,
                testSuiteReport.Description,
                testSuiteReport.Status.String(),
                testSuiteReport.FailureReason,
            })
        }

        table := tablewriter.NewWriter(os.Stdout)
        table.SetHeader([]string{"PATH", "OPERATION", "DESC", "STATUS", "FAILURE REASON"})
        table.SetAutoMergeCellsByColumnIndex([]int{0})
        table.SetRowLine(true)
        table.AppendBulk(data)
        table.Render()
}
