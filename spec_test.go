package main

import (
	"log"
	"os"
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

func (mockTestProcessor) processTestSuites(httpMethod, pathName string, xTestSuites []XTestSuite) []testSuiteReport {
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

/*
   These tests below are going to be used for testing openapi spec parsing
*/

func TestParseOpenApi_AllPositive(t *testing.T) {

	fileBytes, err := os.ReadFile("test-data/openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to read file, err: %v", err)
	}

	_, err = parseAndValidateSpec(fileBytes)
	assert.NoError(t, err)
}

func TestParseOpenApi_MultipleServers(t *testing.T) {

	fileBytes, err := os.ReadFile("test-data/openapi_multiple_servers.yaml")
	if err != nil {
		log.Fatalf("Failed to read file, err: %v", err)
	}

	_, err = parseAndValidateSpec(fileBytes)
	assert.Error(t, err)
}

func TestParseOpenApi(t *testing.T) {

	testSuites := []struct {
		description string
		filePath    string
		shouldFail  bool
	}{
		{
			description: "[POS] Vanilla test case",
			filePath:    "test-data/openapi.yaml",
			shouldFail:  false,
		},
		{
			description: "[NEG] Should fail because of multiple servers",
			filePath:    "test-data/openapi_multiple_servers.yaml",
			shouldFail:  true,
		},
	}

	for _, test := range testSuites {
		t.Run(test.description, func(t *testing.T) {
			fileBytes, _ := os.ReadFile(test.filePath)
			_, err := parseAndValidateSpec(fileBytes)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})

	}
}
