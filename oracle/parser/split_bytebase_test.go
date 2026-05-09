package parser

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type bytebaseSplitCase struct {
	Description string                    `yaml:"description"`
	Input       string                    `yaml:"input"`
	Error       string                    `yaml:"error,omitempty"`
	Result      []bytebaseStatementResult `yaml:"result,omitempty"`
}

type bytebaseStatementResult struct {
	Text  string              `yaml:"text"`
	Range bytebaseRangeResult `yaml:"range"`
	Empty bool                `yaml:"empty"`
}

type bytebaseRangeResult struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

func TestSplitBytebasePLSQLBoundaryCompatibility(t *testing.T) {
	data, err := os.ReadFile("testdata/bytebase_plsql_split.yaml")
	if err != nil {
		t.Fatalf("read Bytebase PLSQL split fixture: %v", err)
	}

	var cases []bytebaseSplitCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal Bytebase PLSQL split fixture: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("Bytebase PLSQL split fixture is empty")
	}

	checked := 0
	skippedErrors := 0
	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			if tc.Error != "" {
				skippedErrors++
				t.Skip("Bytebase parser-level syntax error expectation is outside lexical splitter scope")
			}
			checked++

			got := Split(tc.Input)
			if len(got) != len(tc.Result) {
				t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(tc.Result))
			}
			for i, seg := range got {
				want := tc.Result[i]
				gotComparable := trimBytebaseTrailingHidden(seg.Text)
				wantComparable := trimBytebaseTrailingHidden(want.Text)
				if gotComparable != wantComparable {
					t.Fatalf("segment[%d] comparable text = %q, want %q; raw got %q, raw want %q", i, gotComparable, wantComparable, seg.Text, want.Text)
				}
				if seg.ByteStart != want.Range.Start {
					t.Fatalf("segment[%d] start = %d, want %d", i, seg.ByteStart, want.Range.Start)
				}
				if gotComparableEnd := seg.ByteStart + len(gotComparable); gotComparableEnd != want.Range.End {
					t.Fatalf("segment[%d] comparable end = %d, want %d; raw range [%d,%d]", i, gotComparableEnd, want.Range.End, seg.ByteStart, seg.ByteEnd)
				}
				if seg.Empty() != want.Empty {
					t.Fatalf("segment[%d] empty = %v, want %v", i, seg.Empty(), want.Empty)
				}
				if startsWithStatementDelimiter(seg.Text) {
					t.Fatalf("segment[%d] starts with a statement delimiter: %q", i, seg.Text)
				}
			}
		})
	}

	if checked == 0 {
		t.Fatal("no successful Bytebase split cases checked")
	}
	if skippedErrors == 0 {
		t.Fatal("fixture shape changed: expected at least one parser-level error case")
	}
}

func trimBytebaseTrailingHidden(text string) string {
	end := len(text)
	for {
		trimmed := trimRightSpace(text, end)
		if trimmed != end {
			end = trimmed
			continue
		}
		if blockStart := bytebaseTrailingBlockCommentStart(text, end); blockStart >= 0 {
			end = blockStart
			continue
		}
		if lineStart := bytebaseTrailingLineCommentStart(text, end); lineStart >= 0 {
			end = lineStart
			continue
		}
		return text[:end]
	}
}

func bytebaseTrailingBlockCommentStart(text string, end int) int {
	if end < 4 || text[end-2] != '*' || text[end-1] != '/' {
		return -1
	}
	for i := end - 4; i >= 0; i-- {
		if text[i] != '/' || text[i+1] != '*' {
			continue
		}
		next, ok := splitSkipBlockComment(text, i)
		if ok && next == end {
			return i
		}
	}
	return -1
}

func bytebaseTrailingLineCommentStart(text string, end int) int {
	lineStart := lineStartOffset(text, end)
	for i := end - 2; i >= lineStart; i-- {
		if text[i] == '-' && text[i+1] == '-' {
			return i
		}
	}
	return -1
}

func startsWithStatementDelimiter(text string) bool {
	i := 0
	for i < len(text) {
		switch text[i] {
		case ' ', '\t', '\n', '\r', '\f':
			i++
			continue
		case '-':
			if i+1 < len(text) && text[i+1] == '-' {
				i = splitSkipLineComment(text, i)
				continue
			}
		case '/':
			if i+1 < len(text) && text[i+1] == '*' {
				next, ok := splitSkipBlockComment(text, i)
				if ok {
					i = next
					continue
				}
			}
			return isSlashDelimiterLine(text, i, i+1)
		case ';':
			return true
		}
		return false
	}
	return false
}
