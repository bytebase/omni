package parsertest

import (
	"testing"

	"github.com/stretchr/testify/require"

	es "github.com/bytebase/omni/elasticsearch"
)

func TestGetStatementRanges(t *testing.T) {
	type testCase struct {
		Description string      `yaml:"description,omitempty"`
		Statement   string      `yaml:"statement,omitempty"`
		Result      []es.LSPRange `yaml:"result,omitempty"`
	}

	var (
		filename = "statement-ranges.yaml"
		record   = false
	)

	cases := loadYAML[testCase](t, filename)

	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			got, err := es.GetStatementRanges(tc.Statement)
			require.NoErrorf(t, err, "description: %s", tc.Description)
			if record {
				cases[i].Result = got
			} else {
				require.Equalf(t, tc.Result, got, "description: %s", tc.Description)
			}
		})
	}

	if record {
		writeYAML(t, filename, cases)
	}
}
