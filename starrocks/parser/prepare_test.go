package parser

import (
	"strings"
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

	// prepareSql's string is SQL as well, so its text is parsed too and the
	// tree kept beside it.
	stmt = parseOne(t, `PREPARE p FROM 'SELECT 1'`).(*ast.PrepareStmt)
	if _, ok := stmt.Query.(*ast.SelectStmt); !ok || stmt.SQL != "SELECT 1" {
		t.Errorf("stmt = %+v, want the string SELECT 1 and its tree", stmt)
	}
	for _, sql := range []string{
		`PREPARE p FROM 'a'`,
		`PREPARE p FROM`,
		`PREPARE FROM SELECT 1`,
		// prepareSql admits SINGLE_QUOTED_TEXT, and the string holds one
		// statement; the server refuses each of these.
		`PREPARE p FROM "SELECT 1"`,
		`PREPARE p FROM ''`,
		`PREPARE p FROM '-- c'`,
		`PREPARE p FROM 'SELECT 1; SELECT 2'`,
		// The name is one identifier.
		`PREPARE a.b FROM SELECT 1`,
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
	for _, sql := range []string{
		`EXECUTE`,
		`EXECUTE p USING a`,
		`EXECUTE a.b`,
		`DEALLOCATE p`,
		`DEALLOCATE PREPARE a.b`,
		`DEALLOCATE x y`,
	} {
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

func TestPrepareStringSharesTheNestingDepth(t *testing.T) {
	// The string's statement is parsed at the depth the PREPARE sits at, so a
	// wrapper cannot hide a depth the limit would catch outside one.
	inner := "SELECT " + strings.Repeat("(", 10001) + "1" + strings.Repeat(")", 10001)
	if _, err := Parse("PREPARE p FROM '" + inner + "'"); err == nil {
		t.Error("a PREPARE string deeper than the limit parsed, want an error")
	}
}

func TestPrepareStringKeepsItsTree(t *testing.T) {
	// A string body is SQL too, and its tree is kept, so a query span sees the
	// tables `PREPARE p FROM 'SELECT * FROM secret'` reads. The tree's Locs
	// are offsets into the string's decoded text, not the statement.
	stmt := parseOne(t, `PREPARE p FROM 'SELECT * FROM secret'`).(*ast.PrepareStmt)
	if stmt.SQL != `SELECT * FROM secret` {
		t.Errorf("SQL = %q", stmt.SQL)
	}
	sel, ok := stmt.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}
	if len(sel.From) != 1 {
		t.Errorf("From = %v, want the one table", sel.From)
	}
}

func TestPrepareStringTreeIsAtItsOwnBytes(t *testing.T) {
	// An unescaped literal is parsed where it stands, so a table the prepared
	// SELECT reads sits at its own bytes in the statement.
	const sql = `PREPARE p FROM 'SELECT * FROM secret'`
	stmt := parseOne(t, sql).(*ast.PrepareStmt)
	sel, ok := stmt.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}
	want := strings.Index(sql, "secret")
	if got := ast.NodeLoc(sel.From[0]); got.Start != want || got.End != want+len("secret") {
		t.Errorf("table Loc = %+v, want {%d %d}", got, want, want+len("secret"))
	}
	if got := sql[ast.NodeLoc(sel).Start:ast.NodeLoc(sel).End]; got != "SELECT * FROM secret" {
		t.Errorf("the tree slices %q, want the prepared statement", got)
	}
	// An escaped literal's decoded text is not the statement byte for byte, so
	// its tree is placed through the escape table: the table is still at the
	// bytes that spell it, and a node after an escape keeps its own.
	const escaped = `PREPARE p FROM 'SELECT ''x'' FROM secret'`
	stmt = parseOne(t, escaped).(*ast.PrepareStmt)
	sel, ok = stmt.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Query = %T for an escaped literal, want *ast.SelectStmt", stmt.Query)
	}
	want = strings.Index(escaped, "secret")
	if got := ast.NodeLoc(sel.From[0]); got.Start != want || got.End != want+len("secret") {
		t.Errorf("table Loc = %+v, want {%d %d}", got, want, want+len("secret"))
	}
	if got := escaped[ast.NodeLoc(sel).Start:ast.NodeLoc(sel).End]; got != `SELECT ''x'' FROM secret` {
		t.Errorf("the tree slices %q, want the prepared statement", got)
	}
}
