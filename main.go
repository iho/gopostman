package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const usage = `gopostman - Generate testify suite integration tests from a Postman collection.

Usage:
  gopostman <collection.json> <output.go> [suite-name] [package-name]

Arguments:
  collection.json   Path to the Postman collection v2.1 JSON file.
  output.go         Destination Go file to write (created or overwritten).
  suite-name        Name of the testify suite struct (default: Suite).
  package-name      Go package name for the output file (default: derived from output dir).

Example:
  gopostman collection.json ./tests/integration_gen_test.go Suite integration_test
`

func main() {
	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	collectionPath := os.Args[1]
	outputPath := os.Args[2]

	suiteName := "Suite"
	if len(os.Args) >= 4 {
		suiteName = os.Args[3]
	}

	pkgName := ""
	if len(os.Args) >= 5 {
		pkgName = os.Args[4]
	}
	if pkgName == "" {
		pkgName = derivePkg(outputPath)
	}

	data, err := os.ReadFile(collectionPath)
	if err != nil {
		fatalf("reading collection: %v", err)
	}

	var collection Collection
	if err := json.Unmarshal(data, &collection); err != nil {
		fatalf("parsing collection JSON: %v", err)
	}

	if dir := filepath.Dir(outputPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fatalf("creating output directory: %v", err)
		}
	}

	if err := generate(collection, outputPath, pkgName, suiteName); err != nil {
		fatalf("generating code: %v", err)
	}

	fmt.Printf("generated %s (package %s, suite %s)\n", outputPath, pkgName, suiteName)
}

// derivePkg returns a package name based on the output file's directory name.
func derivePkg(outputPath string) string {
	dir := filepath.Base(filepath.Dir(outputPath))
	if dir == "." || dir == "" {
		return "integration_test"
	}
	return dir
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
