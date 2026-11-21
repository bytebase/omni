package analysis

import (
	"os"
	"testing"

	"github.com/bytebase/omni/mssql"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type limitCorpusCase struct {
	Stmt  string `yaml:"stmt"`
	Limit int    `yaml:"limit"`
	Want  string `yaml:"want"`
}

func TestLimitCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/limit.yaml")
	require.NoError(t, err)
	var cases []limitCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		t.Run(tc.Stmt, func(t *testing.T) {
			got, err := mssql.StatementWithResultLimit(tc.Stmt, tc.Limit)
			require.NoError(t, err)
			require.Equal(t, tc.Want, got)
		})
	}
}
