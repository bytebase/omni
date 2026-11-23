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
		{"trailing whitespace", "SELECT 1;  ", []Segment{{"SELECT 1", 0, 8}, {"  ", 9, 11}}},
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
		{"lone semi", ";", []Segment{{"", 0, 0}}},
		{"double semi", ";;", []Segment{{"", 0, 0}, {"", 1, 1}}},
		{"triple semi", ";;;", []Segment{{"", 0, 0}, {"", 1, 1}, {"", 2, 2}}},
		{"spaced semis", " ; ;", []Segment{{" ", 0, 1}, {" ", 2, 3}}},
		{"leading semi", ";SELECT 1;", []Segment{{"", 0, 0}, {"SELECT 1", 1, 9}}},
		{"trailing semi", "SELECT 1;;", []Segment{{"SELECT 1", 0, 8}, {"", 9, 9}}},
		{"embedded semi", "SELECT 1;;SELECT 2;", []Segment{
			{"SELECT 1", 0, 8},
			{"", 9, 9},
			{"SELECT 2", 10, 18},
		}},
	}
	runSplitCases(t, cases)
}

// TestSplit_EmptyMarked pins which segments Empty() reports: a chunk of
// whitespace, a comment, or nothing is empty; a statement is not.
func TestSplit_EmptyMarked(t *testing.T) {
	cases := []struct {
		input string
		empty []bool
	}{
		{";", []bool{true}},
		{";;", []bool{true, true}},
		{"SELECT 1;;SELECT 2", []bool{false, true, false}},
		{"-- c\n", []bool{true}},
		{"/* c */;SELECT 1", []bool{true, false}},
		{"SELECT 1; -- done", []bool{false, true}},
		{"SELECT 1;  ", []bool{false, true}},
		{"SELECT 1;", []bool{false}},
	}
	for _, c := range cases {
		segments := Split(c.input)
		got := make([]bool, len(segments))
		for i, segment := range segments {
			got[i] = segment.Empty()
			if segment.Text != c.input[segment.ByteStart:segment.ByteEnd] {
				t.Errorf("Split(%q)[%d].Text = %q, want input[%d:%d]", c.input, i, segment.Text, segment.ByteStart, segment.ByteEnd)
			}
		}
		if !reflect.DeepEqual(got, c.empty) {
			t.Errorf("Split(%q) Empty() = %v, want %v", c.input, got, c.empty)
		}
	}
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
		// Neither '//' nor '#' opens a comment (settled against the pinned
		// server, see the lexer's skipWhitespaceAndComments), so nothing is
		// hidden on their line and the segment boundaries do not move.
		{"slash-slash is not a comment", "// comment\nSELECT 1;", []Segment{
			{"// comment\nSELECT 1", 0, 19},
		}},
		{"hash is not a comment", "# comment\nSELECT 1;", []Segment{
			{"# comment\nSELECT 1", 0, 18},
		}},
		{"block before", "/* comment */ SELECT 1;", []Segment{
			{"/* comment */ SELECT 1", 0, 22},
		}},
		{"inline block with semi", "SELECT /* ; */ 1;", []Segment{
			{"SELECT /* ; */ 1", 0, 16},
		}},
		{"comment-only", "-- just a comment\n", []Segment{{"-- just a comment\n", 0, 18}}},
		{"block comment only", "/* only */", []Segment{{"/* only */", 0, 10}}},
		{"comment-only chunk between", "SELECT 1; /* c */; SELECT 2", []Segment{
			{"SELECT 1", 0, 8},
			{" /* c */", 9, 17},
			{" SELECT 2", 18, 27},
		}},
		{"comment-only tail without semi", "SELECT 1; -- done", []Segment{
			{"SELECT 1", 0, 8},
			{" -- done", 9, 17},
		}},
		// A buffer of '//' or '#' holds tokens, so it is a statement chunk,
		// not an empty one; only the engines' own comment forms are empty.
		{"slash-slash only kept", "// note\n", []Segment{{"// note\n", 0, 8}}},
		{"hash only kept", "# note\n", []Segment{{"# note\n", 0, 7}}},
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
	// Whitespace-only content after the last `;` is a trailing segment that
	// Empty() reports.
	got := Split("SELECT 1;   \n\t")
	want := []Segment{{"SELECT 1", 0, 8}, {"   \n\t", 9, 14}}
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

func TestSplit_BeginWithLabelIsTCL(t *testing.T) {
	segs := Split("BEGIN WITH LABEL load_1; COMMIT")
	if len(segs) != 2 {
		t.Fatalf("Split returned %d segments, want 2: %+v", len(segs), segs)
	}

	// A compound block whose first statement is a CTE must stay one segment.
	segs = Split("BEGIN WITH c AS (SELECT 1) SELECT * FROM c; END")
	if len(segs) != 1 {
		t.Fatalf("Split returned %d segments, want 1 (compound block): %+v", len(segs), segs)
	}
}
