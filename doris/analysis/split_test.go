package analysis

import (
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

type yamlLineCol struct {
	Line   int `yaml:"line"`
	Column int `yaml:"column"`
}

type yamlLineCharacter struct {
	Line      int `yaml:"line"`
	Character int `yaml:"character"`
}

// splitCorpus mirrors Bytebase's split fixture, which records the 1-based
// line/column positions and the byte range its adapter reports.
type splitCorpus struct {
	Description string `yaml:"description"`
	Input       string `yaml:"input"`
	Result      []struct {
		Text     string      `yaml:"text"`
		Baseline int         `yaml:"baseline"`
		Start    yamlLineCol `yaml:"start"`
		End      yamlLineCol `yaml:"end"`
		Range    struct {
			Start int `yaml:"start"`
			End   int `yaml:"end"`
		} `yaml:"range"`
		Empty bool `yaml:"empty"`
	} `yaml:"result"`
}

func TestSplitStatementsCorpus(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "split.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var corpus []splitCorpus
	if err := yaml.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	if len(corpus) == 0 {
		t.Fatal("split corpus is empty")
	}
	for _, tc := range corpus {
		t.Run(tc.Description, func(t *testing.T) {
			statements := SplitStatements(tc.Input)
			if len(statements) != len(tc.Result) {
				t.Fatalf("SplitStatements(%q) returned %d statements, want %d", tc.Input, len(statements), len(tc.Result))
			}
			for i, want := range tc.Result {
				got := statements[i]
				if got.Text != want.Text {
					t.Errorf("statement %d Text = %q, want %q", i, got.Text, want.Text)
				}
				if got.Empty != want.Empty {
					t.Errorf("statement %d Empty = %v, want %v", i, got.Empty, want.Empty)
				}
				if got.Start != want.Range.Start || got.End != want.Range.End {
					t.Errorf("statement %d range = [%d,%d), want [%d,%d)", i, got.Start, got.End, want.Range.Start, want.Range.End)
				}
				if line, column := lineColumn1(tc.Input, got.Start); line != want.Start.Line || column != want.Start.Column {
					t.Errorf("statement %d start = %d:%d, want %d:%d", i, line, column, want.Start.Line, want.Start.Column)
				}
				if line, column := lineColumn1(tc.Input, got.End); line != want.End.Line || column != want.End.Column {
					t.Errorf("statement %d end = %d:%d, want %d:%d", i, line, column, want.End.Line, want.End.Column)
				}
			}
		})
	}
}

func TestSplitStatementsEdges(t *testing.T) {
	for _, tc := range []struct {
		script string
		want   []Statement
	}{
		{"", nil},
		{";", nil},
		{";;;", nil},
		{"-- just a comment", nil},
		{"SELECT 1", []Statement{{Text: "SELECT 1", Start: 0, End: 8}}},
		{"SELECT 1;", []Statement{{Text: "SELECT 1;", Start: 0, End: 9}}},
		{" SELECT 1 ; ", []Statement{{Text: " SELECT 1 ;", Start: 0, End: 11}}},
		{"BEGIN SELECT 1; END; SELECT 2;", []Statement{
			{Text: "BEGIN SELECT 1; END;", Start: 0, End: 20},
			{Text: " SELECT 2;", Start: 20, End: 30},
		}},
	} {
		got := SplitStatements(tc.script)
		if len(got) != len(tc.want) {
			t.Errorf("SplitStatements(%q) = %+v, want %+v", tc.script, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("SplitStatements(%q)[%d] = %+v, want %+v", tc.script, i, got[i], tc.want[i])
			}
		}
	}
}

// lineColumn1 converts a byte offset to a 1-based line and 1-based rune
// column, the position convention Bytebase's adapter reports.
func lineColumn1(script string, offset int) (int, int) {
	line, column := 1, 1
	for i := 0; i < offset && i < len(script); {
		r, size := utf8.DecodeRuneInString(script[i:])
		if r == '\n' {
			line++
			column = 1
		} else {
			column++
		}
		i += size
	}
	return line, column
}

// utf16LineCharacter converts a byte offset to a 0-based line and 0-based
// UTF-16 character, the position convention of the LSP statement ranges.
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
