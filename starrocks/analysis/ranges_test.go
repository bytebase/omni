package analysis

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// rangesCorpus mirrors Bytebase's statement-range fixture, which records
// 0-based line and UTF-16 character positions.
type rangesCorpus struct {
	Statement string `yaml:"statement"`
	Ranges    []struct {
		Start yamlLineCharacter `yaml:"start"`
		End   yamlLineCharacter `yaml:"end"`
	} `yaml:"ranges"`
}

func TestStatementRangesCorpus(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "statement-ranges", "*.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("statement-range corpus is empty")
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var corpus []rangesCorpus
		if err := yaml.Unmarshal(data, &corpus); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		for _, tc := range corpus {
			if tc.Statement == "" {
				continue
			}
			t.Run(path+"/"+firstLine(tc.Statement), func(t *testing.T) {
				ranges := StatementRanges(tc.Statement)
				if len(ranges) != len(tc.Ranges) {
					t.Fatalf("StatementRanges returned %d ranges, want %d (%+v)", len(ranges), len(tc.Ranges), ranges)
				}
				for i, want := range tc.Ranges {
					line, character := utf16LineCharacter(tc.Statement, ranges[i].Start)
					if line != want.Start.Line || character != want.Start.Character {
						t.Errorf("range %d start = %d:%d, want %d:%d", i, line, character, want.Start.Line, want.Start.Character)
					}
					line, character = utf16LineCharacter(tc.Statement, ranges[i].End)
					if line != want.End.Line || character != want.End.Character {
						t.Errorf("range %d end = %d:%d, want %d:%d", i, line, character, want.End.Line, want.End.Character)
					}
				}
			})
		}
	}
}

func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}
	return s
}
