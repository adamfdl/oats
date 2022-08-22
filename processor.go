package main

import (
	"encoding/json"
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

		if test.Response.HTTPStatus != statusCode {
			testSuite.Fail()
			testSuites = append(testSuites, testSuite)
			continue
		}

		// TODO@adam: Response object has not been validated
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
