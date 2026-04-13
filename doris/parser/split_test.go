package parser

import (
	"reflect"
	"testing"
)

// runSplitCases is a table-driven helper that asserts Split(input) returns
// exactly the expected Segment list.
func runSplitCases(t *testing.T, cases []struct {
	name  string
	input string
	want  []Segment
}) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Split(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Split(%q) mismatch\n got: %+v\nwant: %+v", c.input, got, c.want)
			}
		})
	}
}

func TestSplit_Empty(t *testing.T) {
	if got := Split(""); got != nil {
		t.Errorf("Split(\"\") = %+v, want nil", got)
	}
}

func TestSplit_Simple(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"one with semi", "SELECT 1;", []Segment{{"SELECT 1", 0, 8}}},
		{"one no semi", "SELECT 1", []Segment{{"SELECT 1", 0, 8}}},
		{"two stmts", "SELECT 1;SELECT 2;", []Segment{
			{"SELECT 1", 0, 8},
			{"SELECT 2", 9, 17},
		}},
		{"three stmts", "A;B;C;", []Segment{
			{"A", 0, 1},
			{"B", 2, 3},
			{"C", 4, 5},
		}},
		{"trailing whitespace", "SELECT 1;  ", []Segment{{"SELECT 1", 0, 8}}},
		{"leading whitespace", "  SELECT 1;", []Segment{{"  SELECT 1", 0, 10}}},
	}
	runSplitCases(t, cases)
}

func TestSplit_EmptyStatements(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"lone semi", ";", nil},
		{"triple semi", ";;;", nil},
		{"leading semi", ";SELECT 1;", []Segment{{"SELECT 1", 1, 9}}},
		{"trailing semi", "SELECT 1;;", []Segment{{"SELECT 1", 0, 8}}},
		{"embedded semi", "SELECT 1;;SELECT 2;", []Segment{
			{"SELECT 1", 0, 8},
			{"SELECT 2", 10, 18},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_Positions(t *testing.T) {
	// Verify byte offsets are accurate for a representative multi-statement input.
	input := "SELECT 1; SELECT 2; SELECT 3"
	got := Split(input)
	want := []Segment{
		{"SELECT 1", 0, 8},
		{" SELECT 2", 9, 18},
		{" SELECT 3", 19, 28},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split mismatch\n got: %+v\nwant: %+v", got, want)
	}
}

func TestSplit_SemisInStrings(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"single-quoted", "SELECT 'a;b';", []Segment{{"SELECT 'a;b'", 0, 12}}},
		{"double-quoted", `SELECT "a;b";`, []Segment{{`SELECT "a;b"`, 0, 12}}},
		{"doubled-quote in string", "SELECT 'a''b;c';", []Segment{{"SELECT 'a''b;c'", 0, 15}}},
		{"backtick ident with semi", "SELECT `a;b`;", []Segment{{"SELECT `a;b`", 0, 12}}},
		{"hex literal", "SELECT X'3B';", []Segment{{"SELECT X'3B'", 0, 12}}},
		{"bit literal", "SELECT B'1';", []Segment{{"SELECT B'1'", 0, 11}}},
	}
	runSplitCases(t, cases)
}

func TestSplit_Comments(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"dash line before", "-- comment\nSELECT 1;", []Segment{
			{"-- comment\nSELECT 1", 0, 19},
		}},
		{"slash-slash line before", "// comment\nSELECT 1;", []Segment{
			{"// comment\nSELECT 1", 0, 19},
		}},
		{"hash comment", "# comment\nSELECT 1;", []Segment{
			{"# comment\nSELECT 1", 0, 18},
		}},
		{"block before", "/* comment */ SELECT 1;", []Segment{
			{"/* comment */ SELECT 1", 0, 22},
		}},
		{"inline block with semi", "SELECT /* ; */ 1;", []Segment{
			{"SELECT /* ; */ 1", 0, 16},
		}},
		{"comment-only filtered", "-- just a comment\n", nil},
		{"block comment only filtered", "/* only */", nil},
	}
	runSplitCases(t, cases)
}

func TestSplit_TCLBegin(t *testing.T) {
	// BEGIN followed by ;, TRANSACTION, WORK, or EOF is a TCL statement.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"bare begin semi", "BEGIN; SELECT 1; COMMIT;", []Segment{
			{"BEGIN", 0, 5},
			{" SELECT 1", 6, 15},
			{" COMMIT", 16, 23},
		}},
		{"begin transaction", "BEGIN TRANSACTION; SELECT 1; COMMIT;", []Segment{
			{"BEGIN TRANSACTION", 0, 17},
			{" SELECT 1", 18, 27},
			{" COMMIT", 28, 35},
		}},
		{"begin work", "BEGIN WORK; COMMIT;", []Segment{
			{"BEGIN WORK", 0, 10},
			{" COMMIT", 11, 18},
		}},
		{"begin at eof", "BEGIN", []Segment{
			{"BEGIN", 0, 5},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_CompoundBlock(t *testing.T) {
	// BEGIN followed by non-TCL tokens is a compound block — the whole
	// BEGIN..END is one statement.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"simple block", "BEGIN SELECT 1; END;", []Segment{
			{"BEGIN SELECT 1; END", 0, 19},
		}},
		{"two stmts in block", "BEGIN INSERT INTO t VALUES (1); INSERT INTO t VALUES (2); END;", []Segment{
			{"BEGIN INSERT INTO t VALUES (1); INSERT INTO t VALUES (2); END", 0, 61},
		}},
		{"nested block", "BEGIN BEGIN SELECT 1; END; END;", []Segment{
			{"BEGIN BEGIN SELECT 1; END; END", 0, 30},
		}},
		{"block then top-level", "BEGIN SELECT 1; END; SELECT 2;", []Segment{
			{"BEGIN SELECT 1; END", 0, 19},
			{" SELECT 2", 20, 29},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_UnclosedBlock(t *testing.T) {
	// BEGIN without a matching END should produce a best-effort single
	// segment covering the whole input (including internal semicolons).
	got := Split("BEGIN SELECT 1;")
	want := []Segment{{"BEGIN SELECT 1;", 0, 15}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split(unclosed BEGIN) = %+v, want %+v", got, want)
	}
}

func TestSplit_TrailingWhitespaceAfterSemi(t *testing.T) {
	// Whitespace-only content after the last `;` should NOT produce a
	// trailing segment (filtered by Empty()).
	got := Split("SELECT 1;   \n\t")
	want := []Segment{{"SELECT 1", 0, 8}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split mismatch\n got: %+v\nwant: %+v", got, want)
	}
}

func TestSplit_Mixed(t *testing.T) {
	// Mixed DML statements with whitespace.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{
			"create insert select",
			"CREATE TABLE t (id INT); INSERT INTO t VALUES (1); SELECT * FROM t",
			[]Segment{
				{"CREATE TABLE t (id INT)", 0, 23},
				{" INSERT INTO t VALUES (1)", 24, 49},
				{" SELECT * FROM t", 50, 66},
			},
		},
		{
			"whitespace between",
			"  SELECT 1 ;  SELECT 2  ",
			[]Segment{
				{"  SELECT 1 ", 0, 11},
				{"  SELECT 2  ", 12, 24},
			},
		},
	}
	runSplitCases(t, cases)
}
