package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockPath = map[string]*PathItem{
	"/users/{id}": &PathItem{
		Get: &Operation{
			XTestSuites: []XTestSuite{
				{
					Description: "[POS] Should return success",
					Request:     XTestSuiteRequest{},
					Response:    XTestSuiteResponse{},
				},
			},
		},
	},
}

var mockServers = []*Server{
	{
		URL:         "http://mock.com",
		Description: "Local dev",
	},
}

var mockSpec = Spec{
	Paths:   mockPath,
	Servers: mockServers,
}

type mockTestProcessor struct{}

func (mockTestProcessor) processTestSuites(httpMethod ,pathName string, xTestSuites []XTestSuite) []testSuiteReport {
	testSuiteReports := []testSuiteReport{
		{
			PathName:    "/users/{id}",
			Operation:   "GET",
			Description: "[POS] Should return success",
			Status:      testResultPassed,
		},
	}
	return testSuiteReports
}

func TestSpec(t *testing.T) {

	report, err := execWithReporter(mockSpec, mockTestProcessor{})
	assert.NoError(t, err)

	assert.Equal(t, "/users/{id}", report.TestSuites[0].PathName)
	assert.Equal(t, "GET", report.TestSuites[0].Operation)
	assert.Equal(t, "[POS] Should return success", report.TestSuites[0].Description)
	assert.Equal(t, testResultPassed, report.TestSuites[0].Status)
}
