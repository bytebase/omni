package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
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
		{"around dollar body", ";CREATE PROCEDURE p() AS $$ a; b $$;;", []Segment{
			{"", 0, 0},
			{"CREATE PROCEDURE p() AS $$ a; b $$", 1, 35},
			{"", 36, 36},
		}},
		{"around scripting block", "; BEGIN SELECT 1; END; ;", []Segment{
			{"", 0, 0},
			{" BEGIN SELECT 1; END", 1, 21},
			{" ", 22, 23},
		}},
		{"around declare block", "DECLARE x INT; BEGIN x := 1; END;;-- c", []Segment{
			{"DECLARE x INT; BEGIN x := 1; END", 0, 32},
			{"", 33, 33},
			{"-- c", 34, 38},
		}},
		{"after task body", "CREATE TASK t AS BEGIN SELECT 1; END;\n;", []Segment{
			{"CREATE TASK t AS BEGIN SELECT 1; END", 0, 36},
			{"\n", 37, 38},
		}},
	}
	runSplitCases(t, cases)
}

// TestSplit_EmptyMarked pins which segments Empty() reports: a chunk of
// whitespace, a comment, or nothing is empty; a statement is not. Each
// segment's Text is its byte range of the input.
func TestSplit_EmptyMarked(t *testing.T) {
	cases := []struct {
		input string
		empty []bool
	}{
		{";", []bool{true}},
		{";;", []bool{true, true}},
		{"SELECT 1;;SELECT 2", []bool{false, true, false}},
		{"-- c\n", []bool{true}},
		{"// c\n", []bool{true}},
		{"/* c */;SELECT 1", []bool{true, false}},
		{"/* a /* nested */ b */;", []bool{true}},
		{"SELECT 1; -- done", []bool{false, true}},
		{"SELECT 1;  ", []bool{false, true}},
		{"SELECT 1;", []bool{false}},
		{"SELECT $$;$$; ;", []bool{false, true}},
		{"BEGIN; ; COMMIT;", []bool{false, true, false}},
		{"EXECUTE IMMEDIATE $$ BEGIN RETURN 1; END $$; /* c */", []bool{false, true}},
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

// TestParseSkipsEmptySegments pins that the parse entry points skip the
// segments Split marks Empty: a lone ';' or a comment adds no statement, and
// the statements keep their absolute positions.
func TestParseSkipsEmptySegments(t *testing.T) {
	input := "; SELECT 1;; /* c */; SELECT 2; -- done"
	for name, result := range map[string]*ParseResult{
		"strict":      ParseStrict(input),
		"best effort": ParseBestEffort(input),
	} {
		if len(result.Errors) != 0 {
			t.Fatalf("%s: errors = %+v, want none", name, result.Errors)
		}
		if len(result.File.Stmts) != 2 {
			t.Fatalf("%s: %d statements, want 2", name, len(result.File.Stmts))
		}
		for i, want := range []string{"SELECT 1", "SELECT 2"} {
			loc := ast.NodeLoc(result.File.Stmts[i])
			if got := input[loc.Start:loc.End]; got != want {
				t.Errorf("%s: statement %d = %q, want %q", name, i, got, want)
			}
		}
	}
	if result := ParseStrict(";; -- only\n/* c */;"); len(result.Errors) != 0 || len(result.File.Stmts) != 0 {
		t.Errorf("ParseStrict(empty segments) = %d statements, errors %+v; want none", len(result.File.Stmts), result.Errors)
	}
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
		{"comment-only", "-- just a comment\n", []Segment{{"-- just a comment\n", 0, 18}}},
		{"slash comment only", "// note\n", []Segment{{"// note\n", 0, 8}}},
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

func TestSplit_MultiLineSingleQuotedBody(t *testing.T) {
	// A single-quoted routine body that SPANS NEWLINES and embeds ';', '{', '}'
	// (a Java / JavaScript / SQL handler) must stay in ONE segment. Split runs its
	// lexer in multilineString mode so the opening ' only ends at the matching
	// unescaped ' — the embedded ';' on a continuation line is NOT a terminator.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"java body", "CREATE FUNCTION f() RETURNS VARCHAR LANGUAGE JAVA AS\n'class T {\n  String e(String x){ return x; }\n}';", []Segment{
			{"CREATE FUNCTION f() RETURNS VARCHAR LANGUAGE JAVA AS\n'class T {\n  String e(String x){ return x; }\n}'", 0, 100},
		}},
		{"js body with semicolons and braces", "CREATE FUNCTION f(d double) RETURNS double LANGUAGE JAVASCRIPT AS '\n  if (D <= 0) { return 1; }\n  return D;\n';", []Segment{
			{"CREATE FUNCTION f(d double) RETURNS double LANGUAGE JAVASCRIPT AS '\n  if (D <= 0) { return 1; }\n  return D;\n'", 0, 109},
		}},
		{"doubled-quote escape spanning newline", "CREATE FUNCTION f() RETURNS VARCHAR AS 'a;\nit''s b;\nc';", []Segment{
			{"CREATE FUNCTION f() RETURNS VARCHAR AS 'a;\nit''s b;\nc'", 0, 54},
		}},
		{"multi-line body then a following statement", "CREATE FUNCTION f() RETURNS INT AS 'x;\ny';\nSELECT 2;", []Segment{
			{"CREATE FUNCTION f() RETURNS INT AS 'x;\ny'", 0, 41},
			{"\nSELECT 2", 42, 51},
		}},
	}
	runSplitCases(t, cases)
}

func TestSplit_StringStateRegression(t *testing.T) {
	// The multilineString change must NOT alter segmentation of ordinary inputs:
	// a single-line string with an embedded ';' is still one segment, a $$ body
	// with ';' is unchanged, and a ';'-separated pair of plain statements still
	// splits into two.
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{"single-line string", "SELECT 'a;b';", []Segment{{"SELECT 'a;b'", 0, 12}}},
		{"single-line string two stmts", "SELECT 'a;b'; SELECT 2;", []Segment{
			{"SELECT 'a;b'", 0, 12},
			{" SELECT 2", 13, 22},
		}},
		{"dollar body with semicolons", "CREATE PROCEDURE p() AS $$ body; with; semicolons $$;", []Segment{
			{"CREATE PROCEDURE p() AS $$ body; with; semicolons $$", 0, 52},
		}},
		{"plain pair", "SELECT 1; SELECT 2;", []Segment{
			{"SELECT 1", 0, 8},
			{" SELECT 2", 9, 18},
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

func TestSplit_UnterminatedMultiLineString(t *testing.T) {
	// LOOP-GUARD: an unterminated single-quoted body (no closing ') in
	// multilineString mode must not loop — Split treats EOF as terminating and
	// returns one best-effort segment covering the whole input.
	got := Split("CREATE FUNCTION f() RETURNS INT AS 'x;\ny")
	want := []Segment{{"CREATE FUNCTION f() RETURNS INT AS 'x;\ny", 0, 40}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Split(unterminated multi-line string) = %+v, want %+v", got, want)
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

			// Every segment must have a valid byte range; one holding a
			// statement has a non-empty one. Segments holding only
			// whitespace or comments are kept and marked Empty().
			for i, seg := range segments {
				if seg.ByteStart > seg.ByteEnd || (!seg.Empty() && seg.ByteStart == seg.ByteEnd) {
					t.Errorf("%s: segment[%d] ByteStart=%d, ByteEnd=%d",
						entry.Name(), i, seg.ByteStart, seg.ByteEnd)
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
	var segments []Segment
	for _, seg := range Split(string(data)) {
		if !seg.Empty() {
			segments = append(segments, seg)
		}
	}
	if len(segments) != 2 {
		t.Errorf("Split(task_scripting.sql) returned %d segments, want 2:", len(segments))
		for i, seg := range segments {
			t.Errorf("  [%d] [%d,%d] %q", i, seg.ByteStart, seg.ByteEnd, seg.Text)
		}
	}
}
