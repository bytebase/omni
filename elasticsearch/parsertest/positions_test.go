package parsertest

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"

	es "github.com/bytebase/omni/elasticsearch"
	"github.com/bytebase/omni/elasticsearch/parser"
	"github.com/bytebase/omni/review"
)

// TestPositions pins positions on text that holds multi-byte characters,
// U+FFFD, or bytes that are not UTF-8. Each failed before:
//   - a statement after an invalid byte started too early, because the
//     splitter counted each invalid byte as the three bytes of U+FFFD;
//   - error columns counted bytes, so a multi-byte character before an
//     error moved it right, and an error on the line of the error before it
//     counted that error's byte a second time;
//   - GetStatementRanges dropped every range after U+FFFD or an invalid byte.
func TestPositions(t *testing.T) {
	t.Run("statement after invalid bytes", func(t *testing.T) {
		stmts, err := es.SplitMultiSQL("GET /a\n{\"k\": \"\xff\xff\"}\nGET /b")
		require.NoError(t, err)
		require.Len(t, stmts, 2)
		require.Equal(t, es.Position{Line: 1, Column: 1}, stmts[0].Start)
		require.Equal(t, es.Position{Line: 2, Column: 12}, stmts[0].End)
		require.Equal(t, es.Position{Line: 3, Column: 1}, stmts[1].Start)
		require.Equal(t, es.Position{Line: 3, Column: 7}, stmts[1].End)
	})
	t.Run("error after a multi-byte character, twice on one line", func(t *testing.T) {
		res, err := es.ParseElasticsearchREST("GET /é\n{\"é\": x}")
		require.NoError(t, err)
		require.Len(t, res.Errors, 2)
		require.Equal(t, es.Position{Line: 1, Column: 8}, res.Errors[0].Position)
		require.Equal(t, es.Position{Line: 1, Column: 8}, res.Errors[1].Position)
		diags, err := es.Diagnose("GET /é\n{\"é\": x}")
		require.NoError(t, err)
		require.Len(t, diags, 2)
		require.Equal(t, es.DiagnosticRange{Start: es.Position{Line: 1, Column: 8}, End: es.Position{Line: 1, Column: 9}}, diags[0].Range)
	})
	for _, in := range []string{"GET /�\nGET /b", "GET /\xff\nGET /b"} {
		t.Run(fmt.Sprintf("range after %q", in[5:len(in)-7]), func(t *testing.T) {
			ranges, err := es.GetStatementRanges(in)
			require.NoError(t, err)
			require.Equal(t, []es.LSPRange{{
				Start: es.LSPPosition{Line: 1, Character: 0},
				End:   es.LSPPosition{Line: 1, Character: 6},
			}}, ranges)
		})
	}
	t.Run("ranges count UTF-16 code units", func(t *testing.T) {
		ranges, err := es.GetStatementRanges("GET /😀\nGET /b😀")
		require.NoError(t, err)
		require.Equal(t, []es.LSPRange{
			{Start: es.LSPPosition{Line: 0, Character: 0}, End: es.LSPPosition{Line: 0, Character: 7}},
			{Start: es.LSPPosition{Line: 1, Character: 0}, End: es.LSPPosition{Line: 1, Character: 8}},
		}, ranges)
	})
}

// TestPositionsAgreeWithReviewIndex checks every position each entry point
// reports against review.Index, the conversion AGENTS.md names for positions
// that reach users, over every prefix and one-character substitution of
// scripts that hold multi-byte characters and bytes that are not UTF-8.
// SplitMultiSQL reports review.Index's positions; ParseElasticsearchREST,
// Diagnose and GetQuerySpan its line made 0-based; GetStatementRanges its
// line and column made 0-based, with each character outside the BMP counted
// twice, as UTF-16 does.
func TestPositionsAgreeWithReviewIndex(t *testing.T) {
	scripts := append([]string{
		"GET /é\n{\"é\": x, \"中\": [1 2]}\nPOST /😀\n{\"a\": \"\xff\"}\nPUT /b\n{\"c\": \"�\"} x\n",
		"GET /a\xff\xfe\nGET /b\n{\"x\": \"\xe4\xb8\"}\n\n# 中文\nDELETE /😀/é\n{} {",
	}, noPanicScripts...)
	substitutes := []string{"é", "中", "😀", "�", "\xff", "\x80", "\xe4\xb8", "\x00", "\n", "x", "{"}
	seen := map[string]bool{}
	checked := 0
	for _, script := range scripts {
		for i := 0; i <= len(script); i++ {
			for _, in := range []string{script[:i], script[:i] + substitutes[i%len(substitutes)] + script[min(i+1, len(script)):]} {
				if seen[in] {
					continue
				}
				seen[in] = true
				if err := checkPositions(in); err != nil {
					t.Fatalf("%q: %v", in, err)
				}
				checked++
			}
		}
	}
	if checked < 1000 {
		t.Fatalf("checked %d inputs", checked)
	}
}

func checkPositions(in string) error {
	index := review.Index(in)

	stmts, err := es.SplitMultiSQL(in)
	if err == nil {
		for i, s := range stmts {
			if s.Range == (es.Range{}) && s.Text == in {
				// No request parsed: the whole text is one statement, with no position.
				continue
			}
			if l, c := index.Position(s.Range.Start); s.Start != (es.Position{Line: l, Column: c}) {
				return fmt.Errorf("statement %d starts at %+v, review.Index gives %d:%d", i, s.Start, l, c)
			}
			if l, c := index.Position(s.Range.End); s.End != (es.Position{Line: l, Column: c}) {
				return fmt.Errorf("statement %d ends at %+v, review.Index gives %d:%d", i, s.End, l, c)
			}
		}
	}

	raw, err := parser.Parse(in)
	if err != nil {
		return nil
	}
	res, err := es.ParseElasticsearchREST(in)
	if err != nil {
		return fmt.Errorf("ParseElasticsearchREST: %v", err)
	}
	if len(res.Errors) != len(raw.Errors) {
		return fmt.Errorf("%d errors, the parser recorded %d", len(res.Errors), len(raw.Errors))
	}
	diags, err := es.Diagnose(in)
	if err != nil {
		return fmt.Errorf("Diagnose: %v", err)
	}
	for i, e := range res.Errors {
		offset := raw.Errors[i].ByteOffset
		if offset >= len(in) {
			// An error at the end of the text stands for its last character.
			_, size := utf8.DecodeLastRuneInString(in)
			offset = len(in) - size
		}
		l, c := index.Position(offset)
		want := es.Position{Line: l - 1, Column: c}
		if e.Position != want {
			return fmt.Errorf("error %d (%q at %d) is at %+v, review.Index gives %+v", i, e.Message, raw.Errors[i].ByteOffset, e.Position, want)
		}
		if diags[i].Range.Start != want || diags[i].Range.End != (es.Position{Line: want.Line, Column: want.Column + 1}) {
			return fmt.Errorf("diagnostic %d is at %+v, review.Index gives %+v", i, diags[i].Range, want)
		}
	}
	if len(res.Errors) > 0 {
		_, err := es.GetQuerySpan(in)
		want := fmt.Sprintf("syntax error at line %d, column %d: %s", res.Errors[0].Position.Line, res.Errors[0].Position.Column, res.Errors[0].Message)
		if err == nil || err.Error() != want {
			return fmt.Errorf("GetQuerySpan: %v, want %s", err, want)
		}
	}

	ranges, err := es.GetStatementRanges(in)
	if err != nil {
		return fmt.Errorf("GetStatementRanges: %v", err)
	}
	var offsets [][2]int
	for _, r := range res.Requests {
		if r.EndOffset > r.StartOffset {
			offsets = append(offsets, [2]int{r.StartOffset, r.EndOffset})
		}
	}
	if len(ranges) != len(offsets) {
		return fmt.Errorf("%d ranges for %d requests", len(ranges), len(offsets))
	}
	lspAt := func(offset int) es.LSPPosition {
		l, c := index.Position(offset)
		units := c - 1
		lineStart := offset
		for lineStart > 0 && in[lineStart-1] != '\n' {
			lineStart--
		}
		// The characters that start before offset, as review.Index counts
		// them, each outside the BMP once more.
		for pos := lineStart; pos < offset && pos < len(in); {
			r, size := utf8.DecodeRuneInString(in[pos:])
			if r > 0xFFFF {
				units++
			}
			pos += size
		}
		return es.LSPPosition{Line: uint32(l - 1), Character: uint32(units)}
	}
	for i, r := range ranges {
		want := es.LSPRange{Start: lspAt(offsets[i][0]), End: lspAt(offsets[i][1])}
		if r != want {
			return fmt.Errorf("range %d is %+v, review.Index gives %+v", i, r, want)
		}
	}
	return nil
}
