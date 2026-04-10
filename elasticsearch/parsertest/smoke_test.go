package parsertest

import "testing"

// TestGoldenFilesLoadable verifies that all 8 YAML golden files from bytebase
// can be read and unmarshaled without error. This is schema-agnostic — it
// deserializes into map[string]any to catch file corruption or YAML syntax
// issues without depending on any parser types.
func TestGoldenFilesLoadable(t *testing.T) {
	files := []string{
		"parse-elasticsearch-rest.yaml",
		"splitter.yaml",
		"statement-ranges.yaml",
		"masking_classify_api.yaml",
		"masking_analyze_body.yaml",
		"masking_predicate_fields.yaml",
		"query_span.yaml",
		"masking_analyze_request.yaml",
	}
	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			loadYAML[map[string]any](t, f)
		})
	}
}
