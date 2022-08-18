package main

import (
	"errors"
	"net/http"
)

type XTestSuiteRequest struct {
	PathParam []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"pathParam,omitempty"`
	QueryParam map[string]string `json:"queryParam,omitempty"`
	Header     map[string]string `json:"header,omitempty"`
        Body       string            `json:"body,omitempty"` // TODO@adam: Possibly refactor this to be map[string]interface{}
}

func (xtsr XTestSuiteRequest) Validate() error {
	return nil
}

type XTestSuiteResponse struct {
	Body       string `json:"body,omitempty"`
	HTTPStatus int    `json:"http-status,omitempty"`
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
        if err :=  s.Paths.Validate(); err != nil {
                return err
        }
        return nil
}

type testProcessor interface {
	processTestSuites(httpMethod, pathName string, xTestSuites []XTestSuite) []testSuiteReport
}

func exec(s Spec) (Report, error) {
	httpClient := newHTTPClient(s.Servers[0].URL)
	testProcessor := newTestProcessor(httpClient)
	return execWithReporter(s, testProcessor)
}

func execWithReporter(s Spec, testProcessor testProcessor) (Report, error) {
        if err := s.Validate(); err != nil {
                return Report{}, err
        }

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

		// TODO@adam: Process Delete operation

		// TODO@adam: Process Patch operation
	}

	return report, nil
}
