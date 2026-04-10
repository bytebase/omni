package parsertest

import (
	"testing"

	"github.com/bytebase/omni/elasticsearch/analysis"
	"github.com/stretchr/testify/require"
)

func TestExtractPredicateFields(t *testing.T) {
	type testCase struct {
		Description string   `yaml:"description"`
		Body        string   `yaml:"body"`
		WantFields  []string `yaml:"wantFields"`
	}

	cases := loadYAML[testCase](t, "masking_predicate_fields.yaml")

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			// Drive through AnalyzeRequest so the predicate walker is exercised
			// end-to-end via the body analysis path.
			result := analysis.AnalyzeRequest("POST", "/my-index/_search", tc.Body)
			if len(tc.WantFields) > 0 {
				require.ElementsMatch(t, tc.WantFields, result.PredicateFields)
			} else {
				require.Empty(t, result.PredicateFields)
			}
		})
	}
}

func TestAnalyzeRequest(t *testing.T) {
	type testCase struct {
		Description         string                    `yaml:"description"`
		Method              string                    `yaml:"method"`
		URL                 string                    `yaml:"url"`
		Body                string                    `yaml:"body"`
		WantAPI             analysis.MaskableAPI      `yaml:"wantAPI"`
		WantIndex           string                    `yaml:"wantIndex"`
		WantBlocked         []analysis.BlockedFeature `yaml:"wantBlocked"`
		WantPredicateFields []string                  `yaml:"wantPredicateFields"`
	}

	cases := loadYAML[testCase](t, "masking_analyze_request.yaml")

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			result := analysis.AnalyzeRequest(tc.Method, tc.URL, tc.Body)
			require.Equal(t, tc.WantAPI, result.API)
			if tc.WantIndex != "" {
				require.Equal(t, tc.WantIndex, result.Index)
			}
			if tc.WantBlocked != nil {
				require.Equal(t, tc.WantBlocked, result.BlockedFeatures)
			} else {
				require.Empty(t, result.BlockedFeatures)
			}
			if tc.WantPredicateFields != nil {
				require.ElementsMatch(t, tc.WantPredicateFields, result.PredicateFields)
			}
		})
	}
}
