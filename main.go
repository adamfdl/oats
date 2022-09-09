package main

import (
	"flag"
	"log"
	"os"
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

	report, err := exec(fileBytes)
	if err != nil {
		log.Fatal(err)
	}

	areAllSuccess := report.generate()
	if !areAllSuccess {
		log.Fatalf("There are failed test suites")
	}
}
