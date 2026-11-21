package analysis_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestRewriteResultLimitCorpus(t *testing.T) {
	type limitCase struct {
		Stmt  string `yaml:"stmt"`
		Limit int    `yaml:"limit"`
		Want  string `yaml:"want"`
	}

	raw, err := os.ReadFile("testdata/limit.yaml")
	require.NoError(t, err)
	var cases []limitCase
	require.NoError(t, yaml.Unmarshal(raw, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		got, err := analysis.RewriteResultLimit(tc.Stmt, tc.Limit)
		require.NoErrorf(t, err, "stmt: %s", tc.Stmt)
		require.Equalf(t, tc.Want, got, "stmt: %s", tc.Stmt)
	}
}
