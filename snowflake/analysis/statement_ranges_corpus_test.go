package analysis_test

import (
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestStatementRangesCorpus(t *testing.T) {
	type lineChar struct {
		Line      int `yaml:"line"`
		Character int `yaml:"character"`
	}
	type rangeCase struct {
		Statement string `yaml:"statement"`
		Ranges    []struct {
			Start lineChar `yaml:"start"`
			End   lineChar `yaml:"end"`
		} `yaml:"ranges"`
	}

	paths, err := filepath.Glob(filepath.Join("testdata", "statement-ranges", "*.yaml"))
	require.NoError(t, err)
	require.NotEmpty(t, paths)

	for _, path := range paths {
		raw, err := os.ReadFile(path)
		require.NoError(t, err)
		var cases []rangeCase
		require.NoError(t, yaml.Unmarshal(raw, &cases))
		require.NotEmpty(t, cases)

		for i, tc := range cases {
			if tc.Statement == "" {
				continue
			}
			got := analysis.StatementRanges(tc.Statement)
			require.Lenf(t, got, len(tc.Ranges), "%s case %d", path, i+1)
			for j, want := range tc.Ranges {
				line, char := utf16LineChar(tc.Statement, got[j].Start)
				require.Equalf(t, want.Start, lineChar{Line: line, Character: char}, "%s case %d range %d start", path, i+1, j)
				line, char = utf16LineChar(tc.Statement, got[j].End)
				require.Equalf(t, want.End, lineChar{Line: line, Character: char}, "%s case %d range %d end", path, i+1, j)
			}
		}
	}
}

func utf16LineChar(script string, offset int) (int, int) {
	line, char := 0, 0
	for i := 0; i < offset && i < len(script); {
		r, size := utf8.DecodeRuneInString(script[i:])
		if r == '\n' {
			line++
			char = 0
		} else if r > 0xFFFF {
			char += 2
		} else {
			char++
		}
		i += size
	}
	return line, char
}
