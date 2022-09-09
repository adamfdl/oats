package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type httpClientMock struct{}

func (mock httpClientMock) execute(_, _ string, _ XTestSuiteRequest) (statusCode int, body []byte, err error) {
	return 200, []byte(`{ "status": "ok" }`), nil
}

func TestProcessGetTestSuite(t *testing.T) {

	xTestSuite := XTestSuite{
		Description: "[POS] Should return success",
		Request: XTestSuiteRequest{
			QueryParam: map[string]string{
				"date": "2022-01-01",
			},
			Header: map[string]string{
				"X-Business-ID": "mock-biz-id",
			},
			PathParam: map[string]string{
				"id": "1",
			},
		},
		Response: XTestSuiteResponse{
			HTTPStatus: http.StatusOK,
			Body:       `{ "status": "ok" }`,
		},
	}

	xTestSuites := []XTestSuite{xTestSuite}
	testProcessor := newTestProcessor(httpClientMock{})

	testSuites := testProcessor.processTestSuites("GET", "/users", xTestSuites)
	assert.Equal(t, "[POS] Should return success", testSuites[0].Description)
	assert.Equal(t, "GET", testSuites[0].Operation)
	assert.Equal(t, testResultPassed, testSuites[0].Status)
	assert.Equal(t, "/users", testSuites[0].PathName)
}
