package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {

	var openApiFilePath string
	flag.StringVar(&openApiFilePath, "f", "", "-f path/to/openapi.yaml")
	flag.Parse()

	if openApiFilePath == "" {
		log.Fatalf("No file is specified. Ex oatsy -f path/to/openapi.yaml")
	}

	fileBytes, err := os.ReadFile(openApiFilePath)
	if err != nil {
		log.Fatalf("Failed to read file, err: %v", err)
	}

	// Load file
	openApiSpec, err := openapi3.NewLoader().LoadFromData(fileBytes)
	if err != nil {
		log.Fatalf("Not a yaml file, err: %v", err)
	}

	// Validate the spec
	ctx := context.Background()
	if err := openApiSpec.Validate(ctx); err != nil {
		log.Fatalf("Bad OpenApi spec, please provide a correct OpenApi format, err: %v", err)
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
