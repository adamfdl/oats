package main

import (
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
			Operation:   "POST",
			Description: "[NEG] Should return error because of BANK_LINKING_ERROR",
			Status: testResultFailed,
                        FailureReason: "Expected this, got that",
		},
	}

        report := Report{
            TestSuites: testSuites,
        }

        report.generate()
}
