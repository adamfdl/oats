package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteGetRequest(t *testing.T) {

	testSuite := XTestSuiteRequest{
		QueryParam: map[string]string{
			"date": "2022-01-01",
		},
		Header: map[string]string{
			"X-Business-ID": "mock-biz-id",
		},
		PathParam: map[string]string{
			"id": "1",
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/1", r.URL.Path)
		assert.Equal(t, "mock-biz-id", r.Header.Get("X-Business-ID"))
		assert.Equal(t, "2022-01-01", r.URL.Query().Get("date"))

		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	httpClient := newHTTPClient(s.URL)
	statusCode, _, err := httpClient.get("/users/{id}", testSuite)
	assert.Nil(t, err)
	assert.Equal(t, 200, statusCode)
}

func TestExecutePostRequest(t *testing.T) {

	testSuite := XTestSuiteRequest{
		Header: map[string]string{
			"X-Business-ID": "mock-biz-id",
		},
		PathParam: map[string]string{
			"id":       "1",
			"order_id": "2",
		},
		Body: map[string]interface{}{
			"token": "mock_token",
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		j := map[string]string{}
		err := json.Unmarshal(body, &j)
		assert.NoError(t, err)
		assert.Equal(t, "mock_token", j["token"])

		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/users/1/order_id/2", r.URL.Path)
		assert.Equal(t, "mock-biz-id", r.Header.Get("X-Business-ID"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(` { "status": "ok" } `))
	}))
	defer s.Close()

	httpClient := newHTTPClient(s.URL)
	statusCode, body, err := httpClient.post("/users/{id}/order_id/{order_id}", testSuite)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	j := map[string]string{}
	err = json.Unmarshal(body, &j)
	assert.NoError(t, err)

	assert.Equal(t, "ok", j["status"])
}
