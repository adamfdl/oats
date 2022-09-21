package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

type XTestSuiteRequest struct {
	PathParam  map[string]string      `json:"pathParam,omitempty"`
	QueryParam map[string]string      `json:"queryParam,omitempty"`
	Header     map[string]string      `json:"header,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"` // TODO@adam: Possibly refactor this to be map[string]interface{}
}

func (xtsr XTestSuiteRequest) Validate() error {
	return nil
}

type XTestSuiteResponse struct {
	//TODO: Probably should refactor these fields as pointers
	Body       string `json:"body,omitempty"`
	HTTPStatus int    `json:"http-status,omitempty"`
}

func (xtsr XTestSuiteResponse) ShouldSkipBodyValidation() bool {
	return xtsr.Body == ""
}

type XTestSuite struct {
	Description string             `json:"description"`
	Request     XTestSuiteRequest  `json:"request"`
	Response    XTestSuiteResponse `json:"response"`
}

type Operation struct {
	XTestSuites []XTestSuite `json:"x-test-suite,omitempty"`
}

func (o *Operation) HasTestSuites() bool {
	return len(o.XTestSuites) != 0
}

type PathItem struct {
	Post *Operation `json:"post"`
	Get  *Operation `json:"get"`
}

type Paths map[string]*PathItem

func (p Paths) Validate() error {
	if len(p) == 0 {
		return errors.New("There is no path to process")
	}
	return nil
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type Servers []*Server

func (s *Servers) Validate() error {
	if s == nil {
		return errors.New("Should specify server")
	}
	if len(*s) != 1 {
		//TODO@adam: For MVP, this can only parse 1 server. If it has more than one servers, it will fail.
		return errors.New("Server should only be 1, no more no less")
	}
	return nil
}

type Spec struct {
	Servers Servers `json:"servers"`
	Paths   Paths   `json:"paths"`
}

func (s Spec) Validate() error {
	if err := s.Servers.Validate(); err != nil {
		return err
	}
	if err := s.Paths.Validate(); err != nil {
		return err
	}
	return nil
}

type testProcessor interface {
	processTestSuites(httpMethod, pathName string, xTestSuites []XTestSuite) []testSuiteReport
}

func parseAndValidateSpec(fileBytes []byte) (s Spec, err error) {
	openApi, err := openapi3.NewLoader().LoadFromData(fileBytes)
	if err != nil {
		return
	}

	ctx := context.Background()
	if err = openApi.Validate(ctx); err != nil {
		return
	}

	// Validate OpenApi spec format
	openApiJSON, err := openApi.MarshalJSON()
	if err != nil {
		return
	}

	if err = json.Unmarshal(openApiJSON, &s); err != nil {
		return
	}

	// Validate Oat's spec (X-Test-Suites)
	if err = s.Validate(); err != nil {
		return
	}

	return s, err
}

func exec(fileBytes []byte) (Report, error) {

	s, err := parseAndValidateSpec(fileBytes)
	if err != nil {
		return Report{}, err
	}

	// Init http client once, so we can test this layer
	httpClient := newHTTPClient(s.Servers[0].URL)

	// Init test processor with httpClient as the dependency
	testProcessor := newTestProcessor(httpClient)

	// Inject dependencies and execute test
	return execWithReporter(s, testProcessor)
}

func execWithReporter(s Spec, testProcessor testProcessor) (Report, error) {
	report := Report{
		TestSuites: []testSuiteReport{},
	}

	// In OpenApi
	for pathName, path := range s.Paths {
		if path.Get != nil {
			if path.Get.HasTestSuites() {
				testReports := testProcessor.processTestSuites(http.MethodGet, pathName, path.Get.XTestSuites)
				report.TestSuites = append(report.TestSuites, testReports...)
			}
		}

		if path.Post != nil && path.Post.HasTestSuites() {
			if path.Post.HasTestSuites() {
				testReports := testProcessor.processTestSuites(http.MethodPost, pathName, path.Post.XTestSuites)
				report.TestSuites = append(report.TestSuites, testReports...)
			}
		}

	}

	return report, nil
}
