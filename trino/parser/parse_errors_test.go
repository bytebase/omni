package parser

import (
	"errors"
	"testing"
)

// TestParse_ReturnsParseErrors pins the error contract every omni engine
// shares: Parse returns nil on success, otherwise an error that errors.As
// resolves to the first *ParseError with a byte Position, while AllErrors
// still yields every error a multi-statement script produced.
func TestParse_ReturnsParseErrors(t *testing.T) {
	if _, err := Parse("SELECT 1"); err != nil {
		t.Fatalf("Parse(valid) = %v, want nil", err)
	}

	file, err := Parse("SELECT 1 ))); SELECT 2; SELECT 3 )))")
	if err == nil {
		t.Fatal("Parse(two bad segments) returned nil error")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("errors.As(*ParseError) failed on %T", err)
	}
	if pe.Position != len("SELECT 1 ") {
		t.Errorf("first error Position = %d, want %d (the first ')')", pe.Position, len("SELECT 1 "))
	}
	if pe.Message == "" {
		t.Error("first error has an empty Message")
	}
	if got := AllErrors(err); len(got) < 2 {
		t.Errorf("AllErrors = %d errors, want at least 2 (one per bad segment): %v", len(got), got)
	}
	if file == nil || len(file.Stmts) != 1 {
		t.Errorf("File keeps the statements that parsed: got %v", file)
	}
	if AllErrors(nil) != nil || AllErrors(errors.New("other")) != nil {
		t.Error("AllErrors must be nil for nil and for non-parse errors")
	}
}

// TestParseBestEffort_ObjectKeywordDrainsSegment pins what an invalid CREATE or
// DROP object keyword leaves behind in the best-effort parse: the syntax error
// at that keyword, and the lexical errors of the rest of the segment, which the
// lazy lexer reaches only because the rejection drains it. Strict Parse lexes
// the whole input itself, so it reports the same two.
func TestParseBestEffort_ObjectKeywordDrainsSegment(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"CREATE INDEX 'unterminated", []string{"syntax error at or near INDEX", "unterminated string literal"}},
		{"DROP FOO 'unterminated", []string{"syntax error at or near FOO", "unterminated string literal"}},
	} {
		res := ParseBestEffort(tc.sql)
		if len(res.Errors) != len(tc.want) {
			t.Fatalf("ParseBestEffort(%q) = %v, want %v", tc.sql, res.Errors, tc.want)
		}
		for i, want := range tc.want {
			if res.Errors[i].Message != want {
				t.Errorf("ParseBestEffort(%q) error %d = %q, want %q", tc.sql, i, res.Errors[i].Message, want)
			}
		}
		if _, err := Parse(tc.sql); len(AllErrors(err)) != len(tc.want) {
			t.Errorf("Parse(%q) = %v, want %v", tc.sql, AllErrors(err), tc.want)
		}
	}
}

// TestParseBestEffort_KeepsExplainedNodeWithDeferredCheck pins the node a
// deferred check (see Parser.deferred) leaves behind: the explained statement is
// invalid, and best effort still returns the EXPLAIN node around it, as it does
// for an inner parse error.
func TestParseBestEffort_KeepsExplainedNodeWithDeferredCheck(t *testing.T) {
	for _, sql := range []string{
		"EXPLAIN SELECT TRIM(FROM 1.5)",
		"EXPLAIN GRANT SELECT TO a",
	} {
		res := ParseBestEffort(sql)
		if len(res.Errors) != 1 {
			t.Fatalf("ParseBestEffort(%q) errors = %v, want 1", sql, res.Errors)
		}
		if len(res.File.Stmts) != 1 {
			t.Fatalf("ParseBestEffort(%q) statements = %v, want the EXPLAIN node", sql, res.File.Stmts)
		}
		if _, ok := res.File.Stmts[0].(*ExplainStmt); !ok {
			t.Errorf("ParseBestEffort(%q) statement is %T, want *ExplainStmt", sql, res.File.Stmts[0])
		}
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want the deferred check's", sql)
		}
	}
}

// TestParse_AlignedErrorPositionsAndForms pins the byte positions and newly
// accepted non-reserved forms closed against Trino 482's compat.tsv.
func TestParse_AlignedErrorPositionsAndForms(t *testing.T) {
	// Valid Trino forms enabled by non-reserved disambiguation, UESCAPE, and
	// DOLLAR_STRING function definitions.
	valid := []string{
		"SELECT U&'!0041' UESCAPE '!'",
		"SELECT JSON_OBJECT(RETURNING : 1)",
		"SELECT JSON_OBJECT(RETURNING VALUE 1)",
		"SELECT * FROM TABLE(f(copartition(a, b), c))",
		"SELECT * FROM TABLE(f(TABLE(a) ORDER BY (x + 1) DESC))",
		"SELECT * FROM t MATCH_RECOGNIZE (a, b)",
		"SELECT a OVER(ROWS UNBOUNDED + 1 PRECEDING)",
		"SELECT a OVER(ROWS CURRENT + 1 PRECEDING)",
		"CREATE TABLE t (a varchar DEFAULT 'x' NOT NULL)",
		"CREATE FUNCTION f() RETURNS int LANGUAGE PYTHON AS $$def f(): return 1;$$; SELECT 1",
	}
	for _, sql := range valid {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) = %v, want nil", sql, err)
		}
	}

	// Aligned error positions across every closed family.
	rejectAt := []struct {
		sql string
		pos int
	}{
		// AS $$...$$ function definition
		{"CREATE FUNCTION a() RETURNS TIMESTAMP AS $$a$$", 41},
		{"CREATE FUNCTION a() RETURNS TIMESTAMP AS", 40},
		{"CREATE FUNCTION a() RETURNS TIMESTAMP AS $$a$$ zzz", 47},
		{"CREATE FUNCTION a() RETURNS TIMESTAMP LANGUAGE SQL AS $$a$$", 54},
		// parameterDeclaration with non-reserved RETURNS / further error
		{"CREATE FUNCTION a(RETURNS TIMESTAMP RETURN 1.5", 36},
		{"DROP FUNCTION a(a zzz TIMESTAMP)", 22},
		// WITH FUNCTION missing name & subquery CTE fallback
		{"WITH FUNCTION() RETURNS TIMESTAMP RETURN 1.5 SELECT 1.5", 14},
		{"SELECT * FROM (WITH FUNCTION f(x int) RETURNS int RETURN x SELECT f(1) AS c) t", 29},
		// EXTRACT field & U&'...' UESCAPE validation
		{"SELECT EXTRACT(a FROM 1.5)", 7},
		{"SELECT EXTRACT(a FROM 1.5) zzz", 7},
		{"SELECT U&'a' UESCAPE 'a'", 7},
		{"SELECT U&'a' UESCAPE", 20},
		{"SELECT U&'a' UESCAPE 'a' zzz", 7},
		// LISTAGG non-string separator
		{"SELECT LISTAGG(1.5,) WITHIN GROUP(ORDER BY 1.5)", 19},
		// JSON_OBJECT KEY / RETURNING / ROW field
		{"SELECT JSON_OBJECT(KEY zzz 1.5 VALUE 1.5)", 27},
		{"SELECT JSON_OBJECT(RETURNING(TIMESTAMP))", 39},
		{"SELECT JSON_OBJECT(RETURNING < TIMESTAMP >)", 41},
		{"SELECT JSON_OBJECT(RETURNING < TIMESTAMP, TIMESTAMP >)", 40},
		{"SELECT JSON_OBJECT(RETURNING ROW(a TIMESTAMP zzz))", 45},
		// TABLE(...) COPARTITION & ORDER BY()
		{"SELECT 1.5 FROM TABLE(a(COPARTITION(a, a), a, a)))", 49},
		{"SELECT 1.5 FROM TABLE(a(TABLE(a) ORDER BY()))", 43},
		// MATCH_RECOGNIZE alias fallback, OVER(MEASURES), AFTER MATCH SKIP TO FIRST/LAST
		{"SELECT 1.5 FROM a MATCH_RECOGNIZE(zzz PATTERN(a) DEFINE a AS 1.5)", 38},
		{"SELECT * FROM t MATCH_RECOGNIZE (DEFINE A AS true)", 40},
		{"SELECT a OVER(MEASURES 1.5 a RANGE UNBOUNDED PRECEDING)", 27},
		{"SELECT a OVER(RANGE UNBOUNDED PRECEDING AFTER MATCH SKIP TO FIRST a zzz)", 68},
		{"SELECT a OVER(RANGE UNBOUNDED PRECEDING AFTER MATCH SKIP TO LAST a zzz)", 67},
		// frameBound UNBOUNDED as non-reserved identifier expression
		{"SELECT a OVER(RANGE BETWEEN UNBOUNDED AND UNBOUNDED PRECEDING)", 61},
		{"SELECT a OVER(ROWS BETWEEN UNBOUNDED AND UNBOUNDED PRECEDING)", 60},
		{"SELECT a OVER(GROUPS BETWEEN UNBOUNDED AND UNBOUNDED PRECEDING)", 62},
		// NOT negated predicate
		{"SELECT 1.5 NOT zzz LIKE 1.5", 15},
		{"MERGE INTO a USING a ON 1.5 NOT MATCHED THEN INSERT VALUES(1.5)", 32},
	}
	for _, tc := range rejectAt {
		_, err := Parse(tc.sql)
		if err == nil {
			t.Errorf("Parse(%q) = nil, want reject@%d", tc.sql, tc.pos)
			continue
		}
		errs := AllErrors(err)
		if len(errs) == 0 || errs[0].Position != tc.pos {
			t.Errorf("Parse(%q) first error pos = %v, want %d (errs=%v)", tc.sql, errs, tc.pos, errs)
		}
	}
}
