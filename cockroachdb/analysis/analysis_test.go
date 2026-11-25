package analysis

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/enginetest/compat"
)

func TestSplit(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []string
	}{
		{statement: "SELECT 1;", want: []string{"SELECT 1;"}},
		{statement: "SELECT\n\n\n\n1;", want: []string{"SELECT\n\n\n\n1;"}},
		{
			statement: "SELECT\n\n\n\n1;\n\n\nSELECT\n\n\n\n2;",
			want:      []string{"SELECT\n\n\n\n1;", "SELECT\n\n\n\n2;"},
		},
		{
			statement: "SELECT\n\n\n\n1;\n\n;\n\n\nSELECT\n\n\n\n2;",
			want:      []string{"SELECT\n\n\n\n1;", "", "SELECT\n\n\n\n2;"},
		},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			if diff := cmp.Diff(tc.want, Split(tc.statement)); diff != "" {
				t.Errorf("Split() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSplitKeepsEmptySegments(t *testing.T) {
	got := Split(`SELECT 1; ; SELECT 2;`)
	want := []string{`SELECT 1;`, ``, `SELECT 2;`}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Split() mismatch (-want +got):\n%s", diff)
	}
}

func TestSplitDoesNotParse(t *testing.T) {
	got := Split(`SELECT 1; SELECT FROM; SELECT 2;`)
	want := []string{`SELECT 1;`, `SELECT FROM;`, `SELECT 2;`}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Split() mismatch (-want +got):\n%s", diff)
	}
}

func TestSplitSegments(t *testing.T) {
	type segment struct {
		Text       string
		Start, End int
		Empty      bool
	}
	for _, tc := range []struct {
		script string
		want   []segment
	}{
		{script: "SELECT '你好';\nSELECT 2;\n  ", want: []segment{
			{Text: "SELECT '你好';", Start: 0, End: 16},
			{Text: "\nSELECT 2;", Start: 16, End: 26},
			{Text: "\n  ", Start: 26, End: 29, Empty: true},
		}},
		{script: ";; -- c", want: []segment{
			{Text: ";", Start: 0, End: 1, Empty: true},
			{Text: ";", Start: 1, End: 2, Empty: true},
			{Text: " -- c", Start: 2, End: 7, Empty: true},
		}},
		// A semicolon inside a BEGIN ATOMIC body does not end the statement.
		{script: "CREATE FUNCTION f() RETURNS INT LANGUAGE SQL BEGIN ATOMIC SELECT 1; SELECT 2; END; SELECT 3", want: []segment{
			{Text: "CREATE FUNCTION f() RETURNS INT LANGUAGE SQL BEGIN ATOMIC SELECT 1; SELECT 2; END;", Start: 0, End: 82},
			{Text: " SELECT 3", Start: 82, End: 91},
		}},
		{script: "BEGIN; SELECT 1; COMMIT;", want: []segment{
			{Text: "BEGIN;", Start: 0, End: 6},
			{Text: " SELECT 1;", Start: 6, End: 16},
			{Text: " COMMIT;", Start: 16, End: 24},
		}},
		// After a lexical error the rest of the script is one segment.
		{script: "SELECT 'unterminated; SELECT 2;", want: []segment{
			{Text: "SELECT 'unterminated; SELECT 2;", Start: 0, End: 31},
		}},
	} {
		t.Run(tc.script, func(t *testing.T) {
			var got []segment
			for _, s := range SplitSegments(tc.script) {
				got = append(got, segment{Text: s.Text, Start: s.ByteStart, End: s.ByteEnd, Empty: s.Empty()})
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("SplitSegments() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSplitKeepsFunctionBodies(t *testing.T) {
	got := Split("CREATE FUNCTION f() RETURNS INT LANGUAGE SQL BEGIN ATOMIC SELECT 1; END; SELECT 2;")
	want := []string{"CREATE FUNCTION f() RETURNS INT LANGUAGE SQL BEGIN ATOMIC SELECT 1; END;", "SELECT 2;"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Split() mismatch (-want +got):\n%s", diff)
	}
}

func TestParseDoBlock(t *testing.T) {
	parsed, err := Parse(`DO $$ BEGIN DELETE FROM t; END $$`)
	if err != nil {
		t.Fatalf("Parse(DO block): %v", err)
	}
	if len(parsed) != 1 {
		t.Fatalf("Parse(DO block) returned %d statements, want 1", len(parsed))
	}
}

func TestParseScript(t *testing.T) {
	const script = "SELECT 1;\n\n;\n\nCREATE TABLE t (id INT)"
	parsed, err := ParseScript(script)
	if err != nil {
		t.Fatalf("ParseScript: %v", err)
	}
	if len(parsed) != 2 {
		t.Fatalf("ParseScript returned %d statements, want 2", len(parsed))
	}
	if parsed[0].Text() != "SELECT 1" {
		t.Errorf("parsed[0].Text() = %q, want %q", parsed[0].Text(), "SELECT 1")
	}
	wantSplit := []string{`SELECT 1;`, ``, `SELECT FROM;`, `SELECT 2;`}
	if diff := cmp.Diff(wantSplit, SplitScript(`SELECT 1; ; SELECT FROM; SELECT 2;`)); diff != "" {
		t.Errorf("SplitScript() mismatch (-want +got):\n%s", diff)
	}
}

// corpus returns the statements of the parser's compat.tsv, each with whether
// the pinned CockroachDB accepts it.
func corpus(t *testing.T) []compat.Row {
	t.Helper()
	f, err := compat.ReadFile(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	return f.Rows
}

func TestSplitSegmentsAgreesWithParser(t *testing.T) {
	var accepted []string
	for _, row := range corpus(t) {
		if row.Expected == compat.Accept {
			accepted = append(accepted, strings.TrimRight(row.Statement, "; \n"))
		}
	}
	for i := 0; i+1 < len(accepted); i++ {
		script := accepted[i] + ";\n" + accepted[i+1] + ";"
		parsed, err := ParseScript(script)
		if err != nil {
			continue
		}
		statements := 0
		for _, segment := range SplitSegments(script) {
			if !segment.Empty() {
				statements++
			}
		}
		if statements != len(parsed) {
			t.Errorf("SplitSegments(%q) holds %d statements, the parser reads %d", script, statements, len(parsed))
		}
	}
}
