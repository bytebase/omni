package analysis

import (
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

// The corpus in testdata/statement-ranges is Bytebase's TiDB statement-range
// fixture, which the editor fixture was recorded from. It is reproduced here,
// never re-recorded against this implementation.
//
// A range is asserted at the position the editor sees it: 0-based lines and
// UTF-16 code-unit columns, which is what this package's byte ranges become at
// the boundary that reports them.

type rangesCorpus struct {
	Statement string `yaml:"statement"`
	Ranges    []struct {
		Start yamlLineCharacter `yaml:"start"`
		End   yamlLineCharacter `yaml:"end"`
	} `yaml:"ranges"`
}

type yamlLineCharacter struct {
	Line      int `yaml:"line"`
	Character int `yaml:"character"`
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

// utf16LineCharacter converts a byte offset in script to the editor's position:
// a 0-based line and a 0-based UTF-16 code-unit column.
func utf16LineCharacter(script string, offset int) (int, int) {
	line, character := 0, 0
	for i := 0; i < offset && i < len(script); {
		c := script[i]
		if c == '\n' {
			line++
			character = 0
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(script[i:])
		if r <= 0xFFFF {
			character++
		} else {
			character += 2
		}
		i += size
	}
	return line, character
}
