package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseRange_PositionsAreAbsolute: a segment parsed in place reports
// every Loc and error position as an offset into the whole script, and
// parses nothing outside its range.
func TestParseRange_PositionsAreAbsolute(t *testing.T) {
	script := "SELECT 1 FROM dual;\nSELECT 2 FROM dual;\nSELECT 3 FROM dual"
	start := strings.Index(script, "SELECT 2")
	end := start + len("SELECT 2 FROM dual")

	list, err := ParseRange(script, start, end)
	if err != nil {
		t.Fatalf("ParseRange: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Items = %d, want 1: the range must not reach the neighbouring statements", len(list.Items))
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		t.Fatalf("item = %T, want *ast.RawStmt", list.Items[0])
	}
	if loc := ast.NodeLoc(raw.Stmt); loc.Start != start || loc.End != end {
		t.Errorf("statement Loc = [%d,%d), want [%d,%d)", loc.Start, loc.End, start, end)
	}

	bad := "SELECT 1 FROM dual;\nSELECT 2 FROM;\nSELECT 3 FROM dual"
	bstart := strings.Index(bad, "SELECT 2")
	_, err = ParseRange(bad, bstart, bstart+len("SELECT 2 FROM"))
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("ParseRange(bad) err = %v, want *ParseError", err)
	}
	if want := bstart + len("SELECT 2 FROM"); pe.Position != want {
		t.Errorf("error Position = %d, want %d (end of the bad segment in script coordinates)", pe.Position, want)
	}
}

// TestParseRange_MatchesPaddedParse: parsing a range gives the same tree and
// positions the old space-padded whole-string parse produced, for a script
// mixing SQL, a PL/SQL block, and a wrapped procedure.
func TestParseRange_MatchesPaddedParse(t *testing.T) {
	script := strings.Join([]string{
		"CREATE TABLE t (id NUMBER)",
		"BEGIN NULL; END",
		"CREATE PROCEDURE p WRAPPED\na000000\nabcd",
		"SELECT id FROM t WHERE id = 1",
	}, ";\n")
	for _, seg := range Split(script) {
		got, err := ParseRange(script, seg.ByteStart, seg.ByteEnd)
		if err != nil {
			t.Fatalf("ParseRange(%q): %v", seg.Text, err)
		}
		want, err := Parse(strings.Repeat(" ", seg.ByteStart) + seg.Text)
		if err != nil {
			t.Fatalf("padded Parse(%q): %v", seg.Text, err)
		}
		if len(got.Items) != 1 || len(want.Items) != 1 {
			t.Fatalf("%q: items ranged=%d padded=%d, want 1 each", seg.Text, len(got.Items), len(want.Items))
		}
		g, w := got.Items[0].(*ast.RawStmt).Stmt, want.Items[0].(*ast.RawStmt).Stmt
		if ast.NodeLoc(g) != ast.NodeLoc(w) {
			t.Errorf("%q: Loc ranged=%v padded=%v", seg.Text, ast.NodeLoc(g), ast.NodeLoc(w))
		}
		if gp, ok := g.(*ast.CreateProcedureStmt); ok {
			wp := w.(*ast.CreateProcedureStmt)
			if gp.WrappedSource != wp.WrappedSource {
				t.Errorf("wrapped source ranged=%q padded=%q", gp.WrappedSource, wp.WrappedSource)
			}
			if strings.Contains(gp.WrappedSource, "SELECT") {
				t.Errorf("wrapped source ran past its segment: %q", gp.WrappedSource)
			}
		}
	}
}

// TestParseRange_ClampsBounds: out-of-range bounds are clamped rather than
// panicking, and an empty range parses to an empty list.
func TestParseRange_ClampsBounds(t *testing.T) {
	list, err := ParseRange("SELECT 1 FROM dual", 5, 100)
	if err == nil && list != nil && len(list.Items) != 0 {
		t.Logf("clamped range parsed %d items", len(list.Items))
	}
	list, err = ParseRange("SELECT 1 FROM dual", 30, 40)
	if err != nil || (list != nil && len(list.Items) != 0) {
		t.Errorf("empty range: list=%+v err=%v, want no items and no error", list, err)
	}
}
