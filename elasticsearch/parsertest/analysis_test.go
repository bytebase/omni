package parsertest

import (
	"testing"

	"github.com/bytebase/omni/elasticsearch/analysis"
	"github.com/stretchr/testify/require"
)

func TestClassifyMaskableAPI(t *testing.T) {
	type testCase struct {
		Description string              `yaml:"description"`
		Method      string              `yaml:"method"`
		URL         string              `yaml:"url"`
		WantAPI     analysis.MaskableAPI `yaml:"wantAPI"`
		WantIndex   string              `yaml:"wantIndex"`
	}

	cases := loadYAML[testCase](t, "masking_classify_api.yaml")

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			api, index := analysis.ClassifyMaskableAPI(tc.Method, tc.URL)
			require.Equal(t, tc.WantAPI, api, "API classification")
			require.Equal(t, tc.WantIndex, index, "index extraction")
		})
	}
}

func TestAnalyzeRequestBody(t *testing.T) {
	type testCase struct {
		Description        string                   `yaml:"description"`
		Body               string                   `yaml:"body"`
		WantBlocked        []analysis.BlockedFeature `yaml:"wantBlocked"`
		WantSourceDisabled bool                     `yaml:"wantSourceDisabled"`
		WantSourceFields   []string                 `yaml:"wantSourceFields"`
		WantFields         []string                 `yaml:"wantFields"`
		WantHighlight      []string                 `yaml:"wantHighlight"`
		WantSort           []string                 `yaml:"wantSort"`
		WantInnerHits      bool                     `yaml:"wantInnerHits"`
	}

	cases := loadYAML[testCase](t, "masking_analyze_body.yaml")

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			result := analysis.AnalyzeRequest("POST", "/my-index/_search", tc.Body)

			if tc.WantBlocked != nil {
				require.Equal(t, tc.WantBlocked, result.BlockedFeatures)
			} else {
				require.Empty(t, result.BlockedFeatures)
			}
			require.Equal(t, tc.WantSourceDisabled, result.SourceDisabled)
			if tc.WantSourceFields != nil {
				require.Equal(t, tc.WantSourceFields, result.SourceFields)
			}
			if tc.WantFields != nil {
				require.Equal(t, tc.WantFields, result.RequestedFields)
			}
			if tc.WantHighlight != nil {
				require.ElementsMatch(t, tc.WantHighlight, result.HighlightFields)
			}
			if tc.WantSort != nil {
				require.Equal(t, tc.WantSort, result.SortFields)
			}
			require.Equal(t, tc.WantInnerHits, result.HasInnerHits)
		})
	}
}
