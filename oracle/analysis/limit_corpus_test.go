package analysis

import (
	"os"
	"testing"

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
			got := rewriteFor19c(t, tc.Stmt, tc.Limit)
			require.Equal(t, tc.Want, got)
		})
	}
}

func TestRewriteResultLimitSkipsOnlySimpleDualSelect(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      string
	}{
		{
			name:      "simple dual select",
			statement: "SELECT 1 FROM DUAL",
			want:      "SELECT 1 FROM DUAL",
		},
		{
			name:      "dual select with scalar subquery",
			statement: "SELECT (SELECT COUNT(*) FROM employees) FROM DUAL",
			want:      "SELECT (SELECT COUNT(*) FROM employees) FROM DUAL FETCH NEXT 10 ROWS ONLY",
		},
		{
			name:      "dual select star keeps historical no-skip behavior",
			statement: "SELECT * FROM DUAL",
			want:      "SELECT * FROM DUAL FETCH NEXT 10 ROWS ONLY",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := rewriteFor19c(t, tc.statement, 10)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRewriteResultLimitKeepsNonSelectStatementsUnchanged(t *testing.T) {
	tests := []struct {
		name          string
		statement     string
		engineVersion string
	}{
		{name: "12c update", statement: "UPDATE employees SET salary = salary + 1", engineVersion: "19.0.0.0.0"},
		{name: "11g delete", statement: "DELETE FROM employees WHERE id = 1", engineVersion: "11.2.0.4.0"},
		{name: "leading whitespace insert", statement: "\n\tINSERT INTO employees(id) VALUES (1)", engineVersion: "19.0.0.0.0"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RewriteResultLimit(tc.statement, 10, tc.engineVersion)
			require.NoError(t, err)
			require.Equal(t, tc.statement, got)
		})
	}
}
