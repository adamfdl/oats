package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {

	// Load file
	openApiSpec, err := openapi3.NewLoader().LoadFromFile("openapi.yaml")
	if err != nil {
		log.Fatal("Failed to load OpenApi spec file")
	}

	// Validate the spec
	ctx := context.Background()
	if err := openApiSpec.Validate(ctx); err != nil {
		log.Fatal("Bad OpenApi spec")
	}

	openApiJSON, err := openApiSpec.MarshalJSON()
	if err != nil {
		log.Fatal("Failed to marshal")
	}

	var s Spec
	if err := json.Unmarshal(openApiJSON, &s); err != nil {
		log.Fatalf("Failed to marshal OpenApi spec, most likely bad X-Test-Suite format\nErr: %v", err)
	}

	report, err := exec(s)
	if err != nil {
		log.Fatalf("Failed to process test suite. Err: %v", err)
	}

	areAllSuccess := report.generate()
	if !areAllSuccess {
		log.Fatalf("There are failed test suites")
	}
}
