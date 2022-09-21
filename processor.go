package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type httpExecuter interface {
	execute(httpMethod, pathName string, xTestSuiteRequest XTestSuiteRequest) (statusCode int, body []byte, err error)
}

type processTestSuites struct {
	httpExecuter httpExecuter
}

func newTestProcessor(httpExecuter httpExecuter) processTestSuites {
	return processTestSuites{
		httpExecuter: httpExecuter,
	}
}

func (testProcessor processTestSuites) processTestSuites(httpMethod, pathName string, xTestSuites []XTestSuite) []testSuiteReport {

	testSuites := []testSuiteReport{}

	for _, test := range xTestSuites {

		testSuite := testSuiteReport{
			PathName:      pathName,
			Description:   test.Description,
			Operation:     httpMethod,
			ResultDetails: resultDetails{},
		}

		statusCode, body, err := testProcessor.httpExecuter.execute(httpMethod, pathName, test.Request)
		if err != nil {
			testSuite.FailWithError(err)
			testSuites = append(testSuites, testSuite)
			continue
		}

		// Set the expected and actual results for the report
		testSuite.ResultDetails.SetActualExpectHTTPStatus(test.Response.HTTPStatus, statusCode)
		testSuite.ResultDetails.SetActualExpectBody(test.Response.Body, string(body))

		// `test.Response.ShouldSkipBodyValidation` should be run first before anything
		// so that the value is set before we do assertions.
		//
		// If we do assertions first there might be a chance that this
		// `test.Response.ShouldSkipBodyValidation` would be skipped
		if test.Response.ShouldSkipBodyValidation() {
			// Modify the test suite report so it does not show
			// body assertions
			testSuite.ShouldSkipBodyValidation = true
			if test.Description == "[NEG] Invalid otp" {
				fmt.Println(testSuite.ShouldSkipBodyValidation)
			}

		} else {
			bodyIsEqual, err := compareResponse([]byte(test.Response.Body), body)
			if err != nil {
				testSuite.FailWithError(err)
				testSuites = append(testSuites, testSuite)
				continue
			}

			if !bodyIsEqual {
				testSuite.Fail()
				testSuites = append(testSuites, testSuite)
				continue
			}
		}

		if test.Response.HTTPStatus != statusCode {
			testSuite.Fail()
			testSuites = append(testSuites, testSuite)
			continue
		}

		testSuite.Pass()
		testSuites = append(testSuites, testSuite)
	}

	return testSuites
}

//TODO: Refactor this code placement
func compareResponse(a, b []byte) (bool, error) {

	var json1, json2 interface{}
	if err := json.Unmarshal(a, &json1); err != nil {
		return false, err
	}

	if err := json.Unmarshal(b, &json2); err != nil {
		return false, err
	}

	return reflect.DeepEqual(json1, json2), nil
}
