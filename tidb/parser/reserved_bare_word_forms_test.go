package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TIDB_CURRENT_TSO wherever a name is read, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-12.
//
// The word is reserved (the ReservedKeyword %token block, parser.y:277) and reaches the
// grammar only as a function name (FunctionNameOptionalBraces, parser.y:8148), so no
// Identifier is spelled with it. omni's keyword table has no token for it, so the lexer
// scans a bare identifier, and parseIdent is what refuses it (tidbReservedBareWords).
// Before, only USE did, so `EXECUTE TIDB_CURRENT_TSO` parsed. use_dbname_forms_test.go
// holds USE.
//
// Not in this table, because omni still disagrees, by reading the word in an expression
// as any other name: `SELECT tidb_current_tso(1)` and `SELECT tidb_current_tso .a FROM t`
// are 1064s on v8.5.5. So is `CREATE TABLE t (a INT) ENGINE = tidb_current_tso`, whose
// name parseIdent does not read.

var reservedBareWordCases = []tableSampleCase{
	// A statement name is an Identifier (PreparedStmt, ExecuteStmt and DeallocateStmt,
	// parser.y:9141-9203).
	{"EXECUTE TIDB_CURRENT_TSO", false},
	{"EXECUTE tidb_current_tso", false},
	{"EXECUTE TIDB_CURRENT_TSO USING @a", false},
	{"PREPARE TIDB_CURRENT_TSO FROM 'select 1'", false},
	{"DEALLOCATE PREPARE TIDB_CURRENT_TSO", false},
	{"DROP PREPARE TIDB_CURRENT_TSO", false},
	{"EXECUTE `tidb_current_tso`", true},
	{"EXECUTE `tidb_current_tso` USING @a", true},
	{"PREPARE `tidb_current_tso` FROM 'select 1'", true},
	{"EXECUTE s1", true},
	{"EXECUTE s1 USING @tidb_current_tso", true},
	// So is every other name.
	{"CREATE TABLE tidb_current_tso (a int)", false},
	{"CREATE TABLE t9 (tidb_current_tso int)", false},
	{"CREATE DATABASE tidb_current_tso", false},
	{"CREATE VIEW tidb_current_tso AS SELECT 1", false},
	{"CREATE SEQUENCE tidb_current_tso", false},
	{"CREATE INDEX tidb_current_tso ON t (a)", false},
	{"CREATE INDEX i9 ON t (tidb_current_tso)", false},
	{"CREATE USER tidb_current_tso", false},
	{"DROP TABLE tidb_current_tso", false},
	{"TRUNCATE TABLE tidb_current_tso", false},
	{"ANALYZE TABLE tidb_current_tso", false},
	{"SHOW CREATE TABLE tidb_current_tso", false},
	{"DESC tidb_current_tso", false},
	{"LOCK TABLES tidb_current_tso READ", false},
	{"ALTER TABLE t RENAME TO tidb_current_tso", false},
	{"ALTER TABLE t ADD COLUMN tidb_current_tso INT", false},
	{"ALTER TABLE t DROP COLUMN tidb_current_tso", false},
	{"ALTER TABLE t ADD INDEX tidb_current_tso (a)", false},
	{"SELECT * FROM tidb_current_tso", false},
	{"SELECT 1 FROM tidb_current_tso t", false},
	{"SELECT a FROM t tidb_current_tso", false},
	{"SELECT a FROM t AS tidb_current_tso", false},
	{"SELECT a FROM t JOIN t tidb_current_tso", false},
	{"SELECT * FROM t PARTITION (tidb_current_tso)", false},
	{"SELECT * FROM t USE INDEX (tidb_current_tso)", false},
	{"SELECT 1 AS tidb_current_tso", false},
	{"SELECT 1 tidb_current_tso", false},
	{"SELECT a tidb_current_tso, b FROM t", false},
	{"SELECT TIDB_CURRENT_TSO AS TIDB_CURRENT_TSO", false},
	{"SELECT 1 FROM t WINDOW tidb_current_tso AS ()", false},
	{"WITH tidb_current_tso AS (SELECT 1) SELECT 1", false},
	{"INSERT INTO t (tidb_current_tso) VALUES (1)", false},
	{"UPDATE t SET tidb_current_tso = 1", false},
	{"SET tidb_current_tso = 1", false},
	{"SET NAMES tidb_current_tso", false},
	{"SAVEPOINT tidb_current_tso", false},
	{"CALL tidb_current_tso()", false},
	// Backquoted it is an ordinary name.
	{"CREATE TABLE `tidb_current_tso` (a int)", true},
	{"CREATE TABLE t9 (`tidb_current_tso` int)", true},
	{"SELECT a AS `tidb_current_tso` FROM t", true},
	// So it is beside a dot, where upstream's lexer looks no keyword up
	// (isTokenIdentifier, misc.go:1107): the dot touches the word before it, and any
	// spaces may stand between a dot and the word after it.
	{"SELECT tidb_current_tso.a FROM t", true},
	{"SELECT tidb_current_tso. a FROM t", true},
	{"SELECT t.tidb_current_tso FROM t", true},
	{"SELECT t. tidb_current_tso FROM t", true},
	{"SELECT tidb_current_tso.* FROM t", true},
	{"SELECT * FROM tidb_current_tso.t", true},
	{"SELECT * FROM t.tidb_current_tso", true},
	{"SELECT a.tidb_current_tso() FROM t", true},
	{"SELECT tidb_current_tso.f() FROM t", true},
	// A user variable and a system variable hold the word inside their own token.
	{"SET @tidb_current_tso = 1", true},
	{"SELECT @tidb_current_tso", true},
	{"SELECT @@tidb_current_tso", true},
	// And it stays the function, with its empty parentheses or without them.
	{"SELECT TIDB_CURRENT_TSO", true},
	{"SELECT TIDB_CURRENT_TSO()", true},
	{"SELECT tidb_current_tso ()", true},
	{"SELECT tidb_current_tso(), 1", true},
	{"SELECT tidb_current_tso FROM t", true},
	{"SELECT tidb_current_tso + 1", true},
	{"SELECT tidb_current_tso AS x", true},
	{"SELECT CAST(tidb_current_tso AS CHAR)", true},
	{"SELECT COUNT(tidb_current_tso) FROM t", true},
	{"SELECT * FROM t WHERE a = tidb_current_tso", true},
	{"SELECT * FROM t ORDER BY tidb_current_tso", true},
	{"INSERT INTO t (a) VALUES (tidb_current_tso)", true},
	{"UPDATE t SET a = tidb_current_tso", true},
	{"SET @a = tidb_current_tso()", true},
	{"DO TIDB_CURRENT_TSO", true},
}

// TestReservedBareWordForms holds omni to the recorded verdicts without a server.
func TestReservedBareWordForms(t *testing.T) {
	for _, c := range reservedBareWordCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestReservedBareWordsAreNoKeywords holds what tidbReservedBareWords is for: a word in
// it has no keyword token. One that gains a token is refused by the token's category,
// and belongs in the keyword sweep instead.
func TestReservedBareWordsAreNoKeywords(t *testing.T) {
	for word := range tidbReservedBareWords {
		if _, ok := keywords[word]; ok {
			t.Errorf("%q has a keyword token now: drop it from tidbReservedBareWords", word)
		}
	}
}

// TestExecuteNameTree holds the name EXECUTE reads when it is written backquoted or is
// an ordinary word, since refusing the reserved one must not change either.
func TestExecuteNameTree(t *testing.T) {
	for sql, want := range map[string]string{
		"EXECUTE `tidb_current_tso`": "tidb_current_tso",
		"EXECUTE s1 USING @a":        "s1",
	} {
		list, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
			continue
		}
		if got := list.Items[0].(*nodes.ExecuteStmt).Name; got != want {
			t.Errorf("Parse(%q) names %q, want %q", sql, got, want)
		}
	}
}
