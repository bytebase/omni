package parsertest

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

// loadYAML loads a YAML file from testdata/ and unmarshals into a slice of T.
// Each feature node's test file defines its own case struct with YAML tags
// matching that golden file's schema.
func loadYAML[T any](t *testing.T, filename string) []T {
	t.Helper()
	data, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatalf("read %s: %v", filename, err)
	}
	var cases []T
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal %s: %v", filename, err)
	}
	if len(cases) == 0 {
		t.Fatalf("%s: expected at least one test case, got 0", filename)
	}
	return cases
}

// writeYAML writes cases back to a YAML file in testdata/ (record mode).
// Feature test files use this to regenerate goldens during development.
func writeYAML[T any](t *testing.T, filename string, cases []T) {
	t.Helper()
	data, err := yaml.Marshal(cases)
	if err != nil {
		t.Fatalf("marshal %s: %v", filename, err)
	}
	if err := os.WriteFile("testdata/"+filename, data, 0644); err != nil {
		t.Fatalf("write %s: %v", filename, err)
	}
}
