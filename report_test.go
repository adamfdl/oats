package main

import (
        "errors"
        "testing"
)

func TestReportGenerateTable(t *testing.T) {
        testSuites := []testSuiteReport{
	        {
			PathName:    "/users/{id}",
			Operation:   "GET",
			Description: "[POS] Should return success",
			Status:      testResultPassed,
		},
	        {
			PathName:    "/users/{id}",
			Operation:   "PATCH",
			Description: "[POS] Should also be successful",
			Status:      testResultPassed,
		},
	        {
			PathName:    "/users",
			Operation:   "PUT",
			Description: "[NEG] Should return error because of BANK_LINKING_ERROR",
			Status: testResultFailed,
                        ResultDetails: resultDetails{
                                HTTPStatusCodes: actualExpectStatusCodes{
                                        Expected: 200,
                                        Actual: 400,
                                },
                                Body: actualExpectBody{
                                    Expected: `{"status": "ok"}`,
                                    Actual: `{"status": "failed"}`,
                                },
                        },
		},
	        {
			PathName:    "/users",
			Operation:   "POST",
			Description: "[NEG] Should return error because of BANK_LINKING_ERROR",
			Status: testResultFailed,
                        Err: errors.New("ECONNREFUSED: dial tcp 127.0.0.1:3002: connect: connection refused"),
		},
	}

        report := Report{
            TestSuites: testSuites,
        }

        report.generate()
}
