package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// parser/prepare.go: prepared statements, EXECUTE AS and TRANSLATE, all
// engine-verified against the pinned StarRocks 3.4.10.

func TestPrepare(t *testing.T) {
	stmt := parseOne(t, `PREPARE p FROM SELECT a FROM t`).(*ast.PrepareStmt)
	if stmt.Name.String() != "p" {
		t.Errorf("Name = %q, want p", stmt.Name.String())
	}
	if _, ok := stmt.Query.(*ast.SelectStmt); !ok {
		t.Errorf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}

	// prepareSql's string is SQL as well, so its text is parsed too.
	stmt = parseOne(t, `PREPARE p FROM 'SELECT 1'`).(*ast.PrepareStmt)
	if stmt.Query != nil || stmt.SQL != "SELECT 1" {
		t.Errorf("stmt = %+v, want the string SELECT 1", stmt)
	}
	for _, sql := range []string{
		`PREPARE p FROM 'a'`,
		`PREPARE p FROM`,
		`PREPARE FROM SELECT 1`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestExecuteAndDeallocate(t *testing.T) {
	stmt := parseOne(t, `EXECUTE p USING @a, @b`).(*ast.ExecuteStmt)
	if stmt.Name.String() != "p" || len(stmt.Using) != 2 || stmt.Using[1] != "b" {
		t.Errorf("stmt = %+v, want p USING @a, @b", stmt)
	}
	for _, sql := range []string{`DEALLOCATE PREPARE p`, `DROP PREPARE p`} {
		stmt := parseOne(t, sql).(*ast.DeallocateStmt)
		if stmt.Name.String() != "p" {
			t.Errorf("%q: Name = %q, want p", sql, stmt.Name.String())
		}
	}
	for _, sql := range []string{`EXECUTE`, `EXECUTE p USING a`, `DEALLOCATE p`} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestExecuteAs(t *testing.T) {
	stmt := parseOne(t, `EXECUTE AS jack@'%'`).(*ast.ExecuteAsStmt)
	if stmt.User.Username != "jack" || stmt.User.Host != "%" || stmt.WithNoRevert {
		t.Errorf("stmt = %+v, want jack@%%", stmt)
	}
	stmt = parseOne(t, `EXECUTE AS jack WITH NO REVERT`).(*ast.ExecuteAsStmt)
	if !stmt.WithNoRevert {
		t.Error("WithNoRevert = false, want true")
	}
	// EXECUTE AS is the statement EXECUTE opens when AS follows; without it
	// the word names a prepared statement.
	if _, errs := parseForTest(`EXECUTE AS`); len(errs) == 0 {
		t.Error("EXECUTE AS parsed, want a syntax error")
	}
}

func TestTranslate(t *testing.T) {
	stmt := parseOne(t, `TRANSLATE presto SELECT a FROM t`).(*ast.TranslateStmt)
	if stmt.Dialect != "presto" || stmt.SQL != "SELECT a FROM t" {
		t.Errorf("stmt = %+v, want presto with the SQL text", stmt)
	}

	// translateSQL is `.+` over raw tokens, so anything after the dialect is
	// the SQL text, whether or not it is SQL.
	stmt = parseOne(t, `TRANSLATE presto SELECT a FROM`).(*ast.TranslateStmt)
	if stmt.SQL != "SELECT a FROM" {
		t.Errorf("SQL = %q, want the text as written", stmt.SQL)
	}

	// At the end of the input there is no token for `.+`, and the pinned
	// server refuses a bare `TRANSLATE a`; its terminator is one, so
	// `TRANSLATE a;` is accepted.
	if _, errs := parseForTest(`TRANSLATE presto`); len(errs) == 0 {
		t.Error("TRANSLATE presto parsed, want a syntax error")
	}
	if _, err := Parse(`TRANSLATE presto;`); err != nil {
		t.Errorf("TRANSLATE presto; : %v", err)
	}
}
