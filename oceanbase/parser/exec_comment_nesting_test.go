package parser

import (
	"context"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// execNestingCases are OceanBase 4.3.5's verdicts on executable comments inside
// executable comments. They do not nest: an opener inside one is dropped, and
// the first "*/" outside a token ends both, so "SELECT 1 /*! + 1 /*! + 2 /* x */
// + 3 */" returns 7 and a second "*/" is a 1064. MySQL 8.0, MariaDB 11.8.8,
// OceanBase 4.3.5 and TiDB v8.5.5 all read them this way.
var execNestingCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT 1 /*! + 1 /*! + 2 /* x */ + 3 */", true},
	{"SELECT 1 /*! + 1 /*! + 2 /* x */ */", true},
	{"SELECT 1 /*! + 1 /*! + 2 */", true},
	{"SELECT 1 /*!50000 + 1 /*!50100 + 2 */", true},
	{"SELECT 1 /*! + 1 /*! + 2 '*/' */", true},
	{"SELECT 1 /*! + 1 /*! /*! + 2 */", true},
	{"SELECT 1 /*! + 1 */ /*! + 2 */", true},
	{"SELECT 1 /*! + 1 /* c */ + 2 */", true},
	{"/*!50000 /*!50100 SELECT 1 */", true},
	{"SELECT 1 /*! + 1 /*! + 2 */ */", false},
	{"SELECT 1 /*! + 1 /*! + 2 */ + 3 */", false},
	{"SELECT 1 /*!50000 + 1 /*!50000 + 2 */ */", false},
	{"/*!50000 /*!50100 SELECT 1 */ */", false},
	{"SELECT 1 /*M! + 1 /*! + 2 */ */", false},
	{"SELECT 1 /*! + 1 /*M! + 2 */", false},
	{"SELECT 1 /*!99999 + 1 /*! + 2 */", true},
	{"SELECT 1 /*!99999 + 1 /*! + 2 */ */", false},
}

// TestExecCommentsDoNotNest pins omni to execNestingCases, which
// TestExecCommentNesting_OceanBase takes from the engine.
func TestExecCommentsDoNotNest(t *testing.T) {
	for _, c := range execNestingCases {
		if _, err := Parse(c.sql); (err == nil) != c.accept {
			t.Errorf("%q: Parse error = %v, want accept %v", c.sql, err, c.accept)
		}
	}
}

// TestExecCommentCloseDropsNestedOpener: the close of an executable comment is
// the first "*/" outside a token, past a nested executable comment's opener,
// not the close of a nested comment that opener would start.
func TestExecCommentCloseDropsNestedOpener(t *testing.T) {
	sql := "/*! SELECT 1 /*! + 2 */ */"
	if got, want := NewLexer(sql).execCommentClose(3), strings.Index(sql, "*/"); got != want {
		t.Errorf("execCommentClose(%q) = %d, want %d", sql, got, want)
	}
}

// TestExecCommentNestingSplit: Split ends nested executable comments where the
// lexer does, so a ';' after the first "*/" is a statement boundary even when a
// later "*/" would have closed the outer comment had they nested.
func TestExecCommentNestingSplit(t *testing.T) {
	sql := "SELECT 1 /*! + 1 /*! + 2 */ ; SELECT 2 */"
	segs := Split(sql)
	if len(segs) != 2 || segs[0].Text != "SELECT 1 /*! + 1 /*! + 2 */ " {
		t.Fatalf("Split(%q) = %+v, want the statement before the ';' and the one after it", sql, segs)
	}
}

// TestExecCommentNesting_OceanBase proves execNestingCases against a live OceanBase 4.3.5.
func TestExecCommentNesting_OceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range execNestingCases {
		ok, err := engineAccepts(ctx, db, c.sql)
		if err != nil {
			t.Fatalf("%q: OceanBase gave no verdict: %v", c.sql, err)
		}
		if ok != c.accept {
			t.Errorf("%q: the engine accepts=%v, the case says %v", c.sql, ok, c.accept)
		}
		if _, err := Parse(c.sql); (err == nil) != ok {
			t.Errorf("%q: omni accepts=%v (err=%v), the engine accepts=%v", c.sql, err == nil, err, ok)
		}
	}
}

// TestExecCommentNestingSplit_OceanBase proves against a live OceanBase 4.3.5 that a ';' after the
// first "*/" of nested executable comments ends the statement: the server
// inserts three rows and then drops the second table, and omni's Split and
// Parse describe the same two statements.
func TestExecCommentNestingSplit_OceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	setup := []string{"DROP DATABASE IF EXISTS exec_nest", "CREATE DATABASE exec_nest", "USE exec_nest", "CREATE TABLE t1 (id INT)", "CREATE TABLE t2 (id INT)"}
	defer conn.ExecContext(ctx, "DROP DATABASE IF EXISTS exec_nest")
	for _, q := range setup {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	statement := "INSERT INTO t1 VALUES (1) /*! , (2) /*! , (3) */ ; DROP TABLE t2; -- */"
	if _, err := conn.ExecContext(ctx, statement); err != nil {
		t.Fatalf("the server rejected %q: %v", statement, err)
	}
	var rows, left int
	if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM t1").Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 't2'").Scan(&left); err != nil {
		t.Fatal(err)
	}
	if rows != 3 || left != 0 {
		t.Fatalf("the server inserted %d rows and kept %d of t2; the oracle expects 3 rows and t2 dropped", rows, left)
	}
	if segs := Split(statement); len(segs) != 2 {
		t.Fatalf("Split(%q) returned %d segments, want 2", statement, len(segs))
	}
	list, err := Parse(statement)
	if err != nil {
		t.Fatalf("Parse(%q): %v", statement, err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("Parse(%q) returned %d statements, want 2", statement, len(list.Items))
	}
	if ins, ok := list.Items[0].(*nodes.InsertStmt); !ok || len(ins.Values) != 3 {
		t.Errorf("statement 1: got %T, want an INSERT of 3 rows", list.Items[0])
	}
	if _, ok := list.Items[1].(*nodes.DropTableStmt); !ok {
		t.Errorf("statement 2: got %T, want a DROP TABLE", list.Items[1])
	}
}
