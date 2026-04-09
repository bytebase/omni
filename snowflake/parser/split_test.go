package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	//        0123456789...
	//        s1=[0,8] ; s2=[9,18] ; s3=[19,28]
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
		{"doubled-quote in string", "SELECT 'a''b;c';", []Segment{{"SELECT 'a''b;c'", 0, 15}}},
		{"dollar-quoted", "SELECT $$a;b$$;", []Segment{{"SELECT $$a;b$$", 0, 14}}},
		{"hex literal", "SELECT X'3B';", []Segment{{"SELECT X'3B'", 0, 12}}},
		{"quoted ident with semi", `SELECT "a;b";`, []Segment{{`SELECT "a;b"`, 0, 12}}},
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
		{"slash line before", "// comment\nSELECT 1;", []Segment{
			{"// comment\nSELECT 1", 0, 19},
		}},
		{"block before", "/* comment */ SELECT 1;", []Segment{
			{"/* comment */ SELECT 1", 0, 22},
		}},
		{"inline block", "SELECT /* a; b */ 1;", []Segment{
			{"SELECT /* a; b */ 1", 0, 19},
		}},
		{"nested block", "/* outer /* inner */ still */ SELECT 1;", []Segment{
			{"/* outer /* inner */ still */ SELECT 1", 0, 38},
		}},
		{"comment-only filtered", "-- just a comment\n", nil},
		{"block comment only filtered", "/* only */", nil},
	}
	runSplitCases(t, cases)
}

func TestSplit_TCLBegin(t *testing.T) {
	// BEGIN followed by ;, TRANSACTION, WORK, NAME, or EOF is a TCL
	// statement, not a Snowflake Scripting block.
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
		{"begin name", "BEGIN NAME tx1; COMMIT;", []Segment{
			{"BEGIN NAME tx1", 0, 14},
			{" COMMIT", 15, 22},
		}},
		{"begin at eof", "BEGIN", []Segment{
			{"BEGIN", 0, 5},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_ScriptingBegin(t *testing.T) {
	// BEGIN followed by non-TCL tokens (INSERT, SELECT, DECLARE, etc.) is
	// a Snowflake Scripting block — the whole BEGIN..END is one statement.
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

func TestSplit_DeclareBegin(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"standalone declare begin", "DECLARE x INTEGER; BEGIN x := 1; END;", []Segment{
			{"DECLARE x INTEGER; BEGIN x := 1; END", 0, 36},
		}},
		{"multiple decls", "DECLARE x INT; y FLOAT; BEGIN x := 1; END;", []Segment{
			{"DECLARE x INT; y FLOAT; BEGIN x := 1; END", 0, 41},
		}},
		{"declare then top-level", "DECLARE x INT; BEGIN SELECT 1; END; SELECT 2;", []Segment{
			{"DECLARE x INT; BEGIN SELECT 1; END", 0, 34},
			{" SELECT 2", 35, 44},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_InlineProcedureBody(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"create task begin end", "CREATE TASK t AS BEGIN SELECT 1; END;", []Segment{
			{"CREATE TASK t AS BEGIN SELECT 1; END", 0, 36},
		}},
		{"create task declare begin end", "CREATE TASK t AS DECLARE x FLOAT; BEGIN x := 1; END;", []Segment{
			{"CREATE TASK t AS DECLARE x FLOAT; BEGIN x := 1; END", 0, 51},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_DollarQuotedBody(t *testing.T) {
	// Dollar-quoted strings are opaque to the splitter — the entire $$..$$
	// block lexes as one tokString, so internal ; BEGIN END are invisible.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"proc with dollar body", "CREATE PROCEDURE p() AS $$ BEGIN SELECT 1; END; $$;", []Segment{
			{"CREATE PROCEDURE p() AS $$ BEGIN SELECT 1; END; $$", 0, 50},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_UnclosedBlock(t *testing.T) {
	// BEGIN without a matching END should produce a best-effort single
	// segment covering the whole input.
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

func TestSplit_LegacyCorpus(t *testing.T) {
	corpusDir := "testdata/legacy"
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("failed to read corpus dir %s: %v", corpusDir, err)
	}

	sqlCount := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		sqlCount++

		t.Run(entry.Name(), func(t *testing.T) {
			path := filepath.Join(corpusDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			segments := Split(string(data))

			// Every corpus file must produce at least one segment.
			if len(segments) == 0 {
				t.Errorf("%s: Split returned 0 segments", entry.Name())
				return
			}

			// Every segment must have valid byte range and non-empty text.
			for i, seg := range segments {
				if seg.ByteStart >= seg.ByteEnd {
					t.Errorf("%s: segment[%d] ByteStart=%d >= ByteEnd=%d",
						entry.Name(), i, seg.ByteStart, seg.ByteEnd)
				}
				if seg.Empty() {
					t.Errorf("%s: segment[%d] Text is empty or whitespace-only: %q",
						entry.Name(), i, seg.Text)
				}
				// Text must equal the raw slice.
				if seg.Text != string(data)[seg.ByteStart:seg.ByteEnd] {
					t.Errorf("%s: segment[%d] Text != input[%d:%d]",
						entry.Name(), i, seg.ByteStart, seg.ByteEnd)
				}
			}
		})
	}

	if sqlCount == 0 {
		t.Errorf("found 0 .sql files in %s — corpus missing?", corpusDir)
	}
	if sqlCount != 27 {
		t.Logf("note: found %d .sql files (expected 27 from the legacy lift)", sqlCount)
	}
}

func TestSplit_TaskScriptingHasTwoSegments(t *testing.T) {
	// task_scripting.sql contains two CREATE TASK statements, each with
	// an inline BEGIN..END body. The legacy (dumb) splitter would produce
	// 5+ fragments. F3's state machine must produce exactly 2 segments.
	data, err := os.ReadFile("testdata/legacy/task_scripting.sql")
	if err != nil {
		t.Fatalf("read task_scripting.sql: %v", err)
	}
	segments := Split(string(data))
	if len(segments) != 2 {
		t.Errorf("Split(task_scripting.sql) returned %d segments, want 2:", len(segments))
		for i, seg := range segments {
			t.Errorf("  [%d] [%d,%d] %q", i, seg.ByteStart, seg.ByteEnd, seg.Text)
		}
	}
}
