package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type httpClient struct {
	baseURL string
}

func newHTTPClient(baseURL string) httpClient {
	return httpClient{
		baseURL: baseURL,
	}
}

func (h httpClient) execute(httpMethod, pathName string, xTestSuiteRequest XTestSuiteRequest) (statusCode int, body []byte, err error) {
	if httpMethod == http.MethodGet {
		return h.get(pathName, xTestSuiteRequest)
	} else if httpMethod == http.MethodPost {
		return h.post(pathName, xTestSuiteRequest)
	} else {
		return 0, nil, err
	}
}

// get method will execute HTTP request, if the test suite request has a body, it will ignore
func (h httpClient) get(pathName string, xTestSuiteRequest XTestSuiteRequest) (statusCode int, body []byte, err error) {

	// Construct path params
	// What this piece of code will do is e.g "/users/{id}/order/{order_id}" -> "/users/1/order/2"
	for key, value := range xTestSuiteRequest.PathParam {
		stringToReplace := fmt.Sprintf("{%s}", key)
		if strings.Contains(pathName, stringToReplace) {
			pathName = strings.Replace(pathName, stringToReplace, value, 1)
		}
	}

	client := resty.New().SetBaseURL(h.baseURL)
	// TODO@adam: Query param and headers needs validation
	// - Is is it eempty
	// - Does the key have a valid type
	// Set query param
	response, err := client.R().
		SetQueryParams(xTestSuiteRequest.QueryParam).
		SetHeaders(xTestSuiteRequest.Header).
		Get(pathName)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode(), response.Body(), nil
}

func (h httpClient) post(pathName string, xTestSuiteRequest XTestSuiteRequest) (statusCode int, body []byte, err error) {

	// Construct path params
	// What this piece of code will do is e.g "/users/{id}/order/{order_id}" -> "/users/1/order/2"
	for key, value := range xTestSuiteRequest.PathParam {
		stringToReplace := fmt.Sprintf("{%s}", key)
		if strings.Contains(pathName, stringToReplace) {
			pathName = strings.Replace(pathName, stringToReplace, value, 1)
		}
	}

	client := resty.New().SetBaseURL(h.baseURL)

	response, respErr := client.R().
		SetBody(xTestSuiteRequest.Body).
		SetHeaders(xTestSuiteRequest.Header).
		Post(pathName)
	if respErr != nil {
		return 0, nil, respErr
	}

	return response.StatusCode(), response.Body(), nil
}
