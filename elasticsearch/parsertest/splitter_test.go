package parsertest

import (
	"testing"

	"github.com/stretchr/testify/require"

	es "github.com/bytebase/omni/elasticsearch"
)

func TestSplitMultiSQL(t *testing.T) {
	type statementCase struct {
		Text  string      `yaml:"text,omitempty"`
		Start es.Position `yaml:"start,omitempty"`
		End   es.Position `yaml:"end,omitempty"`
		Range es.Range    `yaml:"range,omitempty"`
	}
	type testCase struct {
		Description   string          `yaml:"description,omitempty"`
		Statement     string          `yaml:"statement,omitempty"`
		ExpectedCount int             `yaml:"expectedCount"`
		Statements    []statementCase `yaml:"statements,omitempty"`
	}

	cases := loadYAML[testCase](t, "splitter.yaml")

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			got, err := es.SplitMultiSQL(tc.Statement)
			require.NoErrorf(t, err, "description: %s", tc.Description)
			require.Lenf(t, got, tc.ExpectedCount, "description: %s", tc.Description)

			for i, want := range tc.Statements {
				if i >= len(got) {
					break
				}
				stmt := got[i]
				if want.Text != "" {
					require.Equalf(t, want.Text, stmt.Text, "description: %s, statement[%d].text", tc.Description, i)
				}
				if want.Start.Line != 0 || want.Start.Column != 0 {
					require.Equalf(t, want.Start.Line, stmt.Start.Line, "description: %s, statement[%d].start.line", tc.Description, i)
					require.Equalf(t, want.Start.Column, stmt.Start.Column, "description: %s, statement[%d].start.column", tc.Description, i)
				}
				if want.End.Line != 0 || want.End.Column != 0 {
					require.Equalf(t, want.End.Line, stmt.End.Line, "description: %s, statement[%d].end.line", tc.Description, i)
					require.Equalf(t, want.End.Column, stmt.End.Column, "description: %s, statement[%d].end.column", tc.Description, i)
				}
				if want.Range.Start != 0 || want.Range.End != 0 {
					require.Equalf(t, want.Range.Start, stmt.Range.Start, "description: %s, statement[%d].range.start", tc.Description, i)
					require.Equalf(t, want.Range.End, stmt.Range.End, "description: %s, statement[%d].range.end", tc.Description, i)
				}
			}
		})
	}
}
