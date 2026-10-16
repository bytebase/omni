package parser

import (
	"errors"
	"testing"
)

// unclosedParenCases end inside an open parenthesis. SQL Server rejects each
// one at the end of input (Msg 102), and accepts it once close is appended.
// Each case reaches a production that used to drop the error from its
// closing expect(')') or from the nested parse inside the parentheses, so
// Parse returned the statement as if the parenthesis had been closed.
var unclosedParenCases = []struct {
	name  string
	sql   string
	close string
}{
	// Parenthesized expressions and subqueries.
	{"paren/nested", "SELECT ((((1", "))))"},
	{"paren/select-list", "SELECT (1", ")"},
	{"paren/binary-operand", "SELECT 1 + (2", ")"},
	{"paren/where", "SELECT * FROM t WHERE (a = 1", ")"},
	{"paren/update-set", "UPDATE t SET a = (1", ")"},
	{"subquery/scalar", "SELECT (SELECT 1", ")"},
	{"subquery/in", "SELECT * FROM t WHERE a IN (SELECT a FROM t", ")"},
	{"subquery/any", "SELECT * FROM t WHERE a > ANY (SELECT a FROM t", ")"},
	{"subquery/exists", "SELECT * FROM t WHERE EXISTS (SELECT 1", ")"},

	// Built-in functions with their own productions.
	{"func/cast", "SELECT CAST(1 AS int", ")"},
	{"func/try-cast", "SELECT TRY_CAST(1 AS int", ")"},
	{"func/convert", "SELECT CONVERT(int, 1", ")"},
	{"func/convert-style", "SELECT CONVERT(varchar(10), 1, 1", ")"},
	{"func/try-convert", "SELECT TRY_CONVERT(int, 1", ")"},
	{"func/coalesce", "SELECT COALESCE(1, 2", ")"},
	{"func/nullif", "SELECT NULLIF(1, 2", ")"},
	{"func/iif", "SELECT IIF(1 = 1, 1, 0", ")"},
	{"func/count-star", "SELECT COUNT(*", ")"},

	// Window and ordered-set clauses.
	{"over/partition-by", "SELECT COUNT(*) OVER (PARTITION BY a", ")"},
	{"over/frame", "SELECT SUM(a) OVER (ORDER BY a ROWS UNBOUNDED PRECEDING", ")"},
	{"within-group", "SELECT STRING_AGG(a, ',') WITHIN GROUP (ORDER BY a", ")"},

	// Query clauses.
	{"openjson/with", "SELECT * FROM OPENJSON(N'[]') WITH (a int", ")"},
	{"openjson/with-type-length", "SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10", "))"},
	{"for-xml/root", "SELECT a FROM t FOR XML RAW, ROOT('r'", ")"},
	{"for-xml/xmlschema", "SELECT a FROM t FOR XML AUTO, XMLSCHEMA('urn:x'", ")"},
	{"for-json/root", "SELECT a FROM t FOR JSON PATH, ROOT('r'", ")"},

	// Statements with a parenthesized operand.
	{"declare/default", "DECLARE @x int = (1", ")"},
	{"declare/table-var", "DECLARE @t TABLE (a int", ")"},
	{"declare/cursor-query", "DECLARE c CURSOR FOR SELECT (1", ")"},
	{"return", "CREATE PROCEDURE p AS RETURN (1", ")"},
	{"waitfor/receive", "WAITFOR (RECEIVE * FROM q", ")"},
	{"receive/top", "RECEIVE TOP (1", ") * FROM q"},
	{"exec/string", "EXEC ('SELECT 1'", ")"},
	{"insert-exec/string", "INSERT INTO t EXEC ('SELECT 1'", ")"},
	{"raiserror/args", "RAISERROR('x', 16, 1", ")"},
	{"commit/delayed-durability", "COMMIT TRAN WITH (DELAYED_DURABILITY = ON", ")"},
}

// unbalancedParenCases lose a ')' in the middle of the statement, so a later
// ',' or keyword arrives where the inner parenthesis has to close. SQL Server
// rejects bad and accepts good.
var unbalancedParenCases = []struct {
	name string
	bad  string
	good string
}{
	{"values-row", "SELECT * FROM (VALUES (1, 2, (3, 4)) AS t(a, b)", "SELECT * FROM (VALUES (1, 2), (3, 4)) AS t(a, b)"},
	{"grouping-sets", "SELECT a FROM t GROUP BY GROUPING SETS ((a, (a, b), ())", "SELECT a FROM t GROUP BY GROUPING SETS ((a), (a, b), ())"},
	{"count-star", "SELECT COUNT(* FROM t", "SELECT COUNT(*) FROM t"},
	{"type-precision", "SELECT CAST(1 AS decimal((10, 2))", "SELECT CAST(1 AS decimal(10, 2))"},
	{"static-method-arg", "SELECT geography::Point((1, 2, 4326)", "SELECT geography::Point(1, 2, 4326)"},
	{"openjson-type-length", "SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10 '$.a')", "SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10) '$.a')"},
	{"receive-top", "WAITFOR (RECEIVE TOP (1 * FROM q)", "WAITFOR (RECEIVE TOP (1) * FROM q)"},
	{"create-table/identity", "CREATE TABLE t (a int IDENTITY(1, 1, b int)", "CREATE TABLE t (a int IDENTITY(1, 1), b int)"},
	{"create-table/column-check", "CREATE TABLE t (a int CHECK (a > 0, b int)", "CREATE TABLE t (a int CHECK (a > 0), b int)"},
	{"create-table/table-check", "CREATE TABLE t (a int, CHECK (a > 0, b int)", "CREATE TABLE t (a int, CHECK (a > 0), b int)"},
}

// unclosedParenOperandCases open a parenthesis where SQL Server takes only a
// constant or a variable, so it rejects the '(' itself, closed or not. omni
// still parses a full expression in these positions; until it stops, these
// pin that the unclosed form at least is an error rather than a statement.
var unclosedParenOperandCases = []struct {
	name string
	sql  string
}{
	{"set/option-value", "SET TEXTSIZE (1"},
	{"fetch/absolute", "FETCH ABSOLUTE (1"},
	{"waitfor/delay", "WAITFOR DELAY ('00:00:01'"},
	{"waitfor/timeout", "WAITFOR (RECEIVE * FROM q), TIMEOUT (1"},
	{"exec/arg", "EXEC p @a = (1"},
	{"insert-exec/arg", "INSERT INTO t EXEC p @a = (1"},
	{"raiserror/arg", "RAISERROR('x', 16, (1"},
	{"throw/state", "THROW 50000, 'x', (1"},
	{"over/schema-qualified-func", "SELECT dbo.f(1) OVER (ORDER BY a"},
}

// TestUnclosedParenRejected pins Parse's strictness: an unclosed parenthesis
// is a ParseError at the end of input, never a statement.
func TestUnclosedParenRejected(t *testing.T) {
	for _, tc := range unclosedParenCases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				n := 0
				if list != nil {
					n = len(list.Items)
				}
				t.Fatalf("Parse(%q) = %d statements, err %v; want a *ParseError", tc.sql, n, err)
			}
			if pe.Position != len(tc.sql) {
				t.Errorf("Parse(%q) error at %d, want end of input %d: %v", tc.sql, pe.Position, len(tc.sql), err)
			}
			if _, err := Parse(tc.sql + tc.close); err != nil {
				t.Errorf("Parse(%q) closed: %v", tc.sql+tc.close, err)
			}
		})
	}
	for _, tc := range unbalancedParenCases {
		t.Run("unbalanced/"+tc.name, func(t *testing.T) {
			var pe *ParseError
			if _, err := Parse(tc.bad); !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) err %v; want a *ParseError", tc.bad, err)
			}
			if _, err := Parse(tc.good); err != nil {
				t.Errorf("Parse(%q): %v", tc.good, err)
			}
		})
	}
	for _, tc := range unclosedParenOperandCases {
		t.Run("operand/"+tc.name, func(t *testing.T) {
			var pe *ParseError
			if _, err := Parse(tc.sql); !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) err %v; want a *ParseError", tc.sql, err)
			}
		})
	}
}

// TestUnclosedParenOracle checks the expectations above against SQL Server
// (SET PARSEONLY ON) and SqlScriptDOM: both reject every unclosed or
// unbalanced form and accept its closed counterpart.
func TestUnclosedParenOracle(t *testing.T) {
	oracle := startParserOracle(t)
	h := getHarness(t)

	check := func(t *testing.T, sql string, wantAccept bool) {
		t.Helper()
		parseErr, err := oracle.parseError(sql)
		if err != nil {
			t.Fatalf("oracle: %v", err)
		}
		if got := parseErr == nil; got != wantAccept {
			t.Errorf("SQL Server accept=%v (err %v), want %v: %s", got, parseErr, wantAccept, sql)
		}
		errs, _ := h.Shape(t, sql)["errors"].([]any)
		if got := len(errs) == 0; got != wantAccept {
			t.Errorf("SqlScriptDOM accept=%v (errors %v), want %v: %s", got, errs, wantAccept, sql)
		}
	}
	for _, tc := range unclosedParenCases {
		t.Run(tc.name, func(t *testing.T) {
			check(t, tc.sql, false)
			check(t, tc.sql+tc.close, true)
		})
	}
	for _, tc := range unbalancedParenCases {
		t.Run("unbalanced/"+tc.name, func(t *testing.T) {
			check(t, tc.bad, false)
			check(t, tc.good, true)
		})
	}
	for _, tc := range unclosedParenOperandCases {
		t.Run("operand/"+tc.name, func(t *testing.T) {
			check(t, tc.sql, false)
		})
	}
}
