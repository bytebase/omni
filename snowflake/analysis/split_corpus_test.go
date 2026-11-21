package analysis_test

import (
	"os"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestSplitStatementsCorpus(t *testing.T) {
	type lineCol struct {
		Line   int `yaml:"line"`
		Column int `yaml:"column"`
	}
	type splitCase struct {
		Description string `yaml:"description"`
		Input       string `yaml:"input"`
		Result      []struct {
			Text  string  `yaml:"text"`
			Start lineCol `yaml:"start"`
			End   lineCol `yaml:"end"`
			Range struct {
				Start int `yaml:"start"`
				End   int `yaml:"end"`
			} `yaml:"range"`
			Empty bool `yaml:"empty"`
		} `yaml:"result"`
	}

	raw, err := os.ReadFile("testdata/split.yaml")
	require.NoError(t, err)
	var cases []splitCase
	require.NoError(t, yaml.Unmarshal(raw, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			got := analysis.SplitStatements(tc.Input)
			require.Len(t, got, len(tc.Result))
			for i, want := range tc.Result {
				require.Equal(t, want.Text, got[i].Text)
				require.Equal(t, want.Empty, got[i].Empty)
				require.Equal(t, want.Range.Start, got[i].Start)
				require.Equal(t, want.Range.End, got[i].End)
				line, col := lineCol1(tc.Input, got[i].Start)
				require.Equal(t, want.Start, lineCol{Line: line, Column: col})
				line, col = lineCol1(tc.Input, got[i].End)
				require.Equal(t, want.End, lineCol{Line: line, Column: col})
			}
		})
	}
}

func TestSplitStatementsEmptyDropped(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"SELECT 1;;SELECT 2;", 2},
		{"SELECT 1;  ;SELECT 2;", 2},
		{";", 0},
		{";;", 0},
		{"SELECT 1;", 1},
	}
	for _, tc := range cases {
		got := analysis.SplitStatements(tc.in)
		require.Len(t, got, tc.want, "input %q", tc.in)
		for i, st := range got {
			require.NotEmpty(t, st.Text, "stmt %d of %q", i, tc.in)
		}
	}
}

func lineCol1(script string, offset int) (int, int) {
	line, col := 1, 1
	for i := 0; i < offset && i < len(script); {
		r, size := utf8.DecodeRuneInString(script[i:])
		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
		i += size
	}
	return line, col
}
