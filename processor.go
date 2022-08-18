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
			PathName:    pathName,
			Description: test.Description,
			Operation:   httpMethod,
		}

		statusCode, body, err := testProcessor.httpExecuter.execute(httpMethod, pathName, test.Request)
		if err != nil {
			reason := fmt.Sprintf("Failed HTTP request. Error is: %v", err)
			testSuite.Fail(reason)
			testSuites = append(testSuites, testSuite)
			continue
		}

		if test.Response.HTTPStatus != statusCode {
			reason := fmt.Sprintf("Status code mismatch\n\nExpected %d, got: %d", test.Response.HTTPStatus, statusCode)
			testSuite.Fail(reason)
			testSuites = append(testSuites, testSuite)
			continue
		}

		// TODO@adam: Response object has not been validated
                bodyIsEqual, err := compareResponse([]byte(test.Response.Body), body)
                if err != nil {
			reason := fmt.Sprintf("Possible bug in oats project. Error is: %v", err)
			testSuite.Fail(reason)
			testSuites = append(testSuites, testSuite)
			continue
                }

                if !bodyIsEqual {
			reason := fmt.Sprintf("Body mismatch. Expected %s got: %s", test.Response.Body, string(body))
			testSuite.Fail(reason)
			testSuites = append(testSuites, testSuite)
			continue
                }

		testSuite.Pass()
		testSuites = append(testSuites, testSuite)
	}

	return testSuites
}

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
