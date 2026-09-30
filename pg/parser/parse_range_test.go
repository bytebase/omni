package parser

import (
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParseRange parses the segments of one script in place and checks
// that every position is an offset into the script, that text past the
// segment is not read, and that a segment's own error carries its
// absolute offset.
func TestParseRange(t *testing.T) {
	script := "SELECT 1;\nUPDATE t SET a = 'x' WHERE b = 2;\nSELEC 3;\nSELECT 'unterminated"
	at := func(text string) int { return strings.Index(script, text) }
	type segment struct{ start, end int }
	segments := []segment{
		{0, at("\nUPDATE")},
		{at("\nUPDATE"), at("\nSELEC")},
		{at("\nSELEC"), at("\nSELECT '")},
		{at("\nSELECT '"), len(script)},
	}

	list, err := ParseRange(script, segments[1].start, segments[1].end)
	if err != nil {
		t.Fatalf("ParseRange(segment 1): %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("segment 1 parsed %d statements, want 1", len(list.Items))
	}
	raw := list.Items[0].(*nodes.RawStmt)
	if want := (nodes.Loc{Start: at("UPDATE"), End: at("2;") + 1}); raw.Loc != want {
		t.Errorf("segment 1 Loc = %+v, want %+v", raw.Loc, want)
	}
	upd, ok := raw.Stmt.(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("segment 1 statement is %T, want *UpdateStmt", raw.Stmt)
	}
	if upd.Relation.Loc.Start != at("t SET") {
		t.Errorf("relation Loc.Start = %d, want %d", upd.Relation.Loc.Start, at("t SET"))
	}

	// The segment before the failing ones parses although the script as a
	// whole does not: text past end is not read.
	if _, err := ParseRange(script, segments[0].start, segments[0].end); err != nil {
		t.Errorf("ParseRange(segment 0): %v", err)
	}
	if _, err := Parse(script); err == nil {
		t.Error("Parse(script) succeeded, want the syntax error of segment 2")
	}

	// A failing segment reports its error at an absolute offset.
	_, err = ParseRange(script, segments[2].start, segments[2].end)
	var perr *ParseError
	if !errors.As(err, &perr) {
		t.Fatalf("ParseRange(segment 2) = %v, want a *ParseError", err)
	}
	if perr.Position != at("SELEC 3") {
		t.Errorf("segment 2 error Position = %d, want %d", perr.Position, at("SELEC 3"))
	}
	_, err = ParseRange(script, segments[3].start, segments[3].end)
	if !errors.As(err, &perr) || perr.Position < segments[3].start {
		t.Errorf("ParseRange(segment 3) = %v, want an unterminated-string error inside the segment", err)
	}

	// Bounds are clamped rather than trusted: reversed, negative, or past
	// the end, none of them reads outside the text.
	if _, err := ParseRange(script, 5, 2); err != nil {
		t.Errorf("ParseRange(5, 2): %v", err)
	}
	if _, err := ParseRange(script, -1, 0); err != nil {
		t.Errorf("ParseRange(-1, 0): %v", err)
	}
	if _, err := ParseRange(script, -5, -1); err != nil {
		t.Errorf("ParseRange(-5, -1): %v", err)
	}
	if _, err := ParseRange(script, -3, 9); err != nil {
		t.Errorf("ParseRange(-3, 9): %v", err)
	}
	if list, err := ParseRange("", -1, 1); err != nil || (list != nil && len(list.Items) != 0) {
		t.Errorf("ParseRange(\"\", -1, 1) = %v, %v; want nothing", list, err)
	}
	if _, err := ParseRange(script, 0, len(script)+10); err == nil {
		t.Error("ParseRange(0, past end) succeeded, want the script's error")
	}
}

// TestParseRangeReadsNothingPastEnd checks that a construct recognized by
// looking ahead, a psql metacommand, is judged on the range alone: a
// trailing backslash followed by a letter past end is not a metacommand.
func TestParseRangeReadsNothingPastEnd(t *testing.T) {
	script := "SELECT 1\\x"
	if _, err := ParseRange(script, 0, 9); err == nil {
		t.Error("ParseRange over a range ending in a backslash succeeded, want the syntax error Parse gives the same text")
	}
	if _, err := Parse(script[:9]); err == nil {
		t.Error("Parse of the same text succeeded")
	}
	if _, err := ParseRange("SELECT 1;\n\\echo hi\nSELECT 2", 9, 19); err != nil {
		t.Errorf("a metacommand inside the range: %v", err)
	}

	// A range that ends inside the byte order mark holds one stray byte,
	// not a BOM.
	if _, err := ParseRange("\xEF\xBB\xBFSELECT 1", 0, 1); err == nil {
		t.Error("ParseRange over the first BOM byte succeeded, want a syntax error")
	}

	// COPY FROM STDIN takes its inline data only when the rest of the
	// line inside the range is blank; a byte past end does not decide it.
	copyScript := "COPY t FROM STDIN;X"
	list, err := ParseRange(copyScript, 0, 18)
	if err != nil {
		t.Fatalf("ParseRange(COPY): %v", err)
	}
	whole, err := Parse(copyScript[:18])
	if err != nil {
		t.Fatalf("Parse(COPY): %v", err)
	}
	if got, want := list.Items[0].(*nodes.RawStmt).Loc, whole.Items[0].(*nodes.RawStmt).Loc; got != want {
		t.Errorf("COPY Loc in range = %+v, whole text = %+v", got, want)
	}
}

// TestParseRangeSkipsBOMOnlyAtStart checks the byte order mark is trivia
// only at the start of the script.
func TestParseRangeSkipsBOMOnlyAtStart(t *testing.T) {
	script := "\xEF\xBB\xBFSELECT 1;SELECT 2"
	if _, err := ParseRange(script, 0, 12); err != nil {
		t.Errorf("ParseRange with a leading BOM: %v", err)
	}
	if _, err := ParseRange(script, 12, len(script)); err != nil {
		t.Errorf("ParseRange(second segment): %v", err)
	}
}
