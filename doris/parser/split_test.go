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
		{"comment-only filtered", "-- just a comment\n", nil},
		{"block comment only filtered", "/* only */", nil},
		// A buffer of '//' or '#' holds tokens, so it is a statement chunk and
		// Split keeps it; only the engines' own comment forms are filtered.
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

// A procedure's free text is opaque (DorisParser.g4:43 leaves it to `.*?`),
// so Parse reads it up to the boundary Split draws: an unclosed BEGIN in it
// takes the statements after it, a block it closes leaves them their own, and a
// ';' inside its parentheses separates nothing (SEMICOLON is otherwise only a
// statement delimiter, DorisParser.g4:29-33).
func TestProcedureTextFollowsSplit(t *testing.T) {
	// A BEGIN in the free text opens a block only where a statement can begin,
	// since it is also a column name: ALTER PROCEDURE a() zzz BEGIN x; is the
	// text zzz and a column named BEGIN.
	for _, sql := range []string{
		"ALTER PROCEDURE a() BEGIN SELECT 1; END; SELECT 1;",
		"CREATE PROCEDURE p(IN x INT, OUT y VARCHAR(50)) BEGIN SELECT 1; END; SELECT 1;",
		"CREATE PROCEDURE p() COMMENT 'c' BEGIN SELECT 1; END; SELECT 1;",
		"ALTER PROCEDURE a() COMMENT 'c' BEGIN SELECT 1; END; SELECT 2;",
		"ALTER PROCEDURE a(x VARCHAR(50)) BEGIN SELECT 1; END; SELECT 2;",
		"ALTER PROCEDURE a() BEGIN TRANSACTION; SELECT 1;",
		"ALTER PROCEDURE a(); SELECT 1;",
		"CREATE PROCEDURE a(x;y); SELECT 1;",
		"ALTER PROCEDURE a(x;y) zzz; SELECT 1;",
		"REPLACE PROC a(x;y); SELECT 1;",
		"ALTER PROCEDURE a(BEGIN x); SELECT 1;",
		// The `.*?` closes at the first ')', so a '(' inside it nests nothing.
		"ALTER PROCEDURE a((x); SELECT 1;",
		// A hint span is read past here as the parser reads past it.
		"/*+ hint */ CREATE PROCEDURE a(x;y); SELECT 1;",
		"CREATE PROCEDURE /*+ hint */ a(x;y); SELECT 1;",
	} {
		segments := Split(sql)
		if len(segments) == 0 {
			t.Errorf("Split(%q) found no segment", sql)
		}
		for _, segment := range segments {
			if _, err := Parse(segment.Text); err != nil {
				t.Errorf("Parse(Split(%q) segment %q) = %v, want it accepted", sql, segment.Text, err)
			}
		}
	}

	// ALTER and REPLACE reach a procedure only through PROCEDURE or PROC, and
	// the name is multipartIdentifier (identifier ('.' identifier)*) whose first
	// part a reserved keyword cannot be, so other statements keep their
	// separators even when their text leaves a parenthesis open. So does a
	// procedure's own parenthesis when it never closes.
	for _, sql := range []string{
		"ALTER TABLE t ADD COLUMN c DECIMAL(10; SELECT 1;",
		"ALTER PROCEDURE a zzz (; SELECT 1;",
		"ALTER PROCEDURE a. (; SELECT 1;",
		"ALTER PROCEDURE SELECT (; SELECT 1;",
		"ALTER PROCEDURE a(x; SELECT 1;",
		// END CASE with no CASE open closes what it can, here the block, so
		// the statement after it keeps its own segment.
		"BEGIN SELECT 1; END CASE; SELECT 2;",
	} {
		if got := Split(sql); len(got) != 2 {
			t.Errorf("Split(%q) = %+v, want 2 segments", sql, got)
		}
	}
}

// END IF, END LOOP and END CASE close their own statement, not the block, so a
// procedure with control flow in its body stays one statement, as it does for
// Parse.
func TestSplitControlFlowEndInBody(t *testing.T) {
	for _, sql := range []string{
		"CREATE PROCEDURE p() BEGIN IF a THEN b; END IF; END; SELECT 1;",
		"CREATE PROCEDURE p() BEGIN LOOP SET a = 1; END LOOP; END; SELECT 1;",
		"CREATE PROCEDURE p() BEGIN CASE WHEN a THEN SET b = 1; END CASE; END; SELECT 1;",
		"ALTER PROCEDURE a() BEGIN WHILE b DO SET a = 1; END WHILE; END; SELECT 1;",
		// An expression's CASE ends at a bare END, which closes the CASE and
		// not the block, and a CASE branch holding a BEGIN gets its END back.
		"REPLACE PROCEDURE p() BEGIN SELECT CASE WHEN 1 THEN 2 END; SELECT 3; END; SELECT 4;",
		"CREATE PROCEDURE p() BEGIN CASE WHEN 1 THEN BEGIN SELECT 1; END; END CASE; END; SELECT 2;",
		// CASE, BEGIN and END after a dot are field names, not constructs.
		"CREATE PROCEDURE p() BEGIN SELECT t.case FROM t; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT t.begin FROM t; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT t.end FROM t; END; SELECT 2;",
		// So are BEGIN and END standing inside a statement, since both are
		// non-reserved words Doris reads as identifiers.
		"CREATE PROCEDURE p() BEGIN SELECT begin FROM t; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT end FROM t; END; SELECT 2;",
		// A BEGIN a control-flow statement puts before its body still opens a
		// block, and a named terminator inside a CASE closes its own statement.
		"CREATE PROCEDURE p() BEGIN IF a THEN BEGIN SELECT 1; END; END IF; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN LOOP BEGIN SELECT 1; END; END LOOP; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN lbl: BEGIN SELECT 1; END; END; SELECT 2;",
		// NULL and CURRENT_DATE end an operand, so the CASE's END closes it.
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN 1 THEN NULL END; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN 1 THEN CURRENT_DATE END; END; SELECT 2;",
		// A unary operator and a predicate operator leave an operand due too.
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN ! end THEN CASE WHEN 1 THEN 2 END ELSE 3 END; END; SELECT 9;",
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN 'x' REGEXP end THEN CASE WHEN 1 THEN 2 END ELSE 3 END; END; SELECT 9;",
		// A declaration's own body opens after what it ends with, and a type
		// parenthesis closes the parameter list early.
		"CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR SQLEXCEPTION BEGIN SELECT 1; END; SELECT 2; END; SELECT 3;",
		"CREATE PROCEDURE p(x VARCHAR(50), y INT) BEGIN SELECT 1; END; SELECT 2;",
		// BEGIN after a statement's first word names a column.
		"ALTER PROCEDURE a() SELECT begin FROM t; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN CASE WHEN a THEN IF b THEN SET x = 1; END IF; END CASE; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN CASE WHEN a THEN LOOP SET x = 1; END LOOP; END CASE; END; SELECT 2;",
		// An END naming a column inside a CASE branch, and an expression CASE
		// whose item is aliased LOOP: the END closes the CASE the word after it
		// is not the terminator of.
		"CREATE PROCEDURE p() BEGIN CASE WHEN a THEN SELECT end FROM t; END CASE; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN 1 THEN 2 END loop; END; SELECT 2;",
		// A CASE nested in an expression's branch is an expression too, and
		// begin and end stand in it as column names.
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN 1 THEN CASE WHEN 2 THEN 3 END END; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN 1 THEN begin ELSE 2 END; END; SELECT 2;",
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN end THEN CASE WHEN 1 THEN 2 END ELSE 3 END; END; SELECT 9;",
		"CREATE PROCEDURE p() BEGIN SELECT CASE WHEN 1 THEN 2 END + 1; END; SELECT 2;",
	} {
		segments := Split(sql)
		if len(segments) != 2 {
			t.Errorf("Split(%q) = %d segments, want 2", sql, len(segments))
			continue
		}
		if _, err := Parse(segments[0].Text); err != nil {
			t.Errorf("Parse(%q) = %v, want it accepted", segments[0].Text, err)
		}
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

func TestSplit_BeginWithLabelIsTCL(t *testing.T) {
	// BEGIN WITH LABEL is the transaction form: the ; after it must split.
	segs := Split("BEGIN WITH LABEL load_1; COMMIT")
	if len(segs) != 2 {
		t.Fatalf("Split returned %d segments, want 2: %+v", len(segs), segs)
	}
	if segs[0].Text != "BEGIN WITH LABEL load_1" || segs[1].Text != " COMMIT" {
		t.Errorf("segments = %q, %q", segs[0].Text, segs[1].Text)
	}

	// A compound block whose first statement is a CTE must stay one segment:
	// WITH alone does not make BEGIN a transaction.
	segs = Split("BEGIN WITH c AS (SELECT 1) SELECT * FROM c; END")
	if len(segs) != 1 {
		t.Fatalf("Split returned %d segments, want 1 (compound block): %+v", len(segs), segs)
	}
}
