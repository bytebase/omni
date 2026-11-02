package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	gomysql "github.com/go-sql-driver/mysql"
)

// execNestingCases are TiDB v8.5.5's verdicts on executable comments inside
// executable comments and on the feature list of a /*T![...] gate. TiDB does
// not nest them: an opener inside one is dropped, and the first "*/" outside a
// token ends both, so "SELECT 1 /*! + 1 /*! + 2 /* x */ + 3 */" returns 7 and
// a second "*/" is a 1064. A feature list is names of letters, digits, '_' and
// '$' separated by single commas; anything else leaves "/*T!" an executable
// comment whose body starts at the "[", which TiDB rejects.
var execNestingCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT 1 /*! + 1 /*! + 2 /* x */ + 3 */", true},
	{"SELECT 1 /*! + 1 /*! + 2 /* x */ */", true},
	{"SELECT 1 /*! + 1 /*! + 2 */", true},
	{"SELECT 1 /*T! + 1 /*T! + 2 /* x */ + 3 */", true},
	{"SELECT 1 /*! + 1 /*T![clustered_index] + 2 */", true},
	{"SELECT 1 /*! + 1 /*T![no_such_feature] + 2 */ + 3 */", true},
	{"SELECT 1 /*T![clustered_index] + 1 /*T![clustered_index] + 2 */", true},
	{"SELECT 1 /*! + 1 /*!50000 + 2 */", true},
	{"SELECT 1 /*! + 1 /*T! + 2 */", true},
	{"SELECT 1 /*! + 1 /*! + 2 '*/' */", true},
	{"SELECT 1 /*! + 1 /*! /*! + 2 */", true},
	{"SELECT 1 /*! + 1 */ /*! + 2 */", true},
	{"SELECT 1 /*! + 1 /*! + 2 */ */", false},
	{"SELECT 1 /*! + 1 /*! + 2 */ + 3 */", false},
	{"SELECT 1 /*!50000 + 1 /*!50000 + 2 */ */", false},
	{"SELECT 1 /*T! + 1 /*! + 2 */ */", false},
	{"SELECT 1 /*T![clustered_index] + 1 /*T![clustered_index] + 2 */ */", false},
	{"SELECT 1 /*! + 1 /*T![clustered_index + 2 */ + 3 */", false},

	{"SELECT 1 /*T![clustered_index,auto_rand] + 2 */", true},
	{"SELECT 1 /*T![clustered_index]+ 2 */", true},
	{"SELECT 1 /*T![no_such_feature] + 2 */ + 3", true},
	{"SELECT 1 /*T![clustered_index,no_such_feature] + 2 */ + 3", true},
	{"SELECT 1 /*T![Clustered_Index] + 2 */", true},
	{"SELECT 1 /*T![clustered_index$] + 2 */", true},
	{"SELECT 1 /*T![9] + 2 */", true},
	{"SELECT 1 /*T![clustered_indexé] + 2 */", true},
	{"SELECT 1 /*T![é] + 2 */", true},
	{"SELECT 1 /*T![clustered_index] , 2", true},
	{"/*T![clustered_index] WITH c AS (SELECT 1 AS a) SELECT a FROM c", true},
	{"SELECT 1 /*T![clustered_index + 2 */", false},
	{"SELECT 1 /*T![clustered_index + 2 */ + 3", false},
	{"SELECT 1 /*T![ */ + 3", false},
	{"SELECT 1 /*T![no_such_feature + 2 */ + 3", false},
	{"SELECT 1 /*T![clustered_index ] + 2 */", false},
	{"SELECT 1 /*T![ clustered_index] + 2 */", false},
	{"SELECT 1 /*T![] + 2 */", false},
	{"SELECT 1 /*T![clustered_index,] + 2 */", false},
	{"SELECT 1 /*T![clustered_index, auto_rand] + 2 */", false},
	{"SELECT 1 /*T![clustered-index] + 2 */", false},
}

// TestExecCommentsDoNotNest pins omni to the verdicts in execNestingCases,
// which TestExecCommentNesting_Oracle takes from TiDB.
func TestExecCommentsDoNotNest(t *testing.T) {
	for _, tc := range execNestingCases {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("%q: Parse error = %v, want accept %v", tc.sql, err, tc.accept)
		}
	}
}

// TestMalformedFeatureListErrorsAtBracket: TiDB reads "/*T!" before a feature
// list it cannot read as an executable comment whose body starts at the "[",
// and reports the 1064 there (column 14).
func TestMalformedFeatureListErrorsAtBracket(t *testing.T) {
	sql := "SELECT 1 /*T![clustered_index ] + 2 */"
	_, err := Parse(sql)
	var pe *ParseError
	if !errors.As(err, &pe) || pe.Position != strings.Index(sql, "[") {
		t.Fatalf("Parse(%q) error = %v, want a syntax error at the '[' (%d)", sql, err, strings.Index(sql, "["))
	}
}

// TestExecCommentNesting_Oracle proves execNestingCases against a live TiDB
// v8.5.5. A statement the server neither runs nor rejects with a syntax error
// fails the test rather than counting as either verdict.
func TestExecCommentNesting_Oracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range execNestingCases {
		err := tc.canExecute(c.sql)
		var me *gomysql.MySQLError
		if err != nil && !errors.As(err, &me) {
			t.Fatalf("%q: TiDB gave no verdict: %v", c.sql, err)
		}
		tidbAccepts := err == nil || me.Number != 1064
		if tidbAccepts != c.accept {
			t.Errorf("%q: TiDB accepts=%v (err=%v), the case says %v", c.sql, tidbAccepts, err, c.accept)
		}
		if _, perr := Parse(c.sql); (perr == nil) != tidbAccepts {
			t.Errorf("%q: omni accepts=%v (err=%v), TiDB accepts=%v", c.sql, perr == nil, perr, tidbAccepts)
		}
	}
}

// TestExecCommentNestingSplit_Oracle proves against a live TiDB that the first
// "*/" ends executable comments nested in one another, so a ';' after it ends
// the statement: the server inserts three rows and then drops t2, and omni's
// Split and Parse describe the same two statements.
func TestExecCommentNestingSplit_Oracle(t *testing.T) {
	tc := startTiDB(t)
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer conn.Close()
	for _, q := range []string{"CREATE DATABASE IF NOT EXISTS exec_nest_1", "USE exec_nest_1", "DROP TABLE IF EXISTS t1, t2", "CREATE TABLE t1 (id INT)", "CREATE TABLE t2 (id INT)"} {
		if _, err := conn.ExecContext(tc.ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS exec_nest_1")

	statement := "INSERT INTO t1 VALUES (1) /*! , (2) /*! , (3) */ ; DROP TABLE t2; -- */"
	if _, err := conn.ExecContext(tc.ctx, statement); err != nil {
		t.Fatalf("server rejected %q: %v", statement, err)
	}
	var rows, t2 int
	if err := conn.QueryRowContext(tc.ctx, "SELECT COUNT(*) FROM t1").Scan(&rows); err != nil {
		t.Fatalf("count t1: %v", err)
	}
	if err := conn.QueryRowContext(tc.ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'exec_nest_1' AND table_name = 't2'").Scan(&t2); err != nil {
		t.Fatalf("check t2: %v", err)
	}
	if rows != 3 || t2 != 0 {
		t.Fatalf("server inserted %d rows and kept t2=%d; the oracle expects 3 rows and t2 dropped", rows, t2)
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
	ins, ok := list.Items[0].(*ast.InsertStmt)
	if !ok || len(ins.Values) != 3 {
		t.Errorf("statement 1: got %T with %d rows, want an INSERT of 3 rows", list.Items[0], valuesLen(ins))
	}
	if _, ok := list.Items[1].(*ast.DropTableStmt); !ok {
		t.Errorf("statement 2: got %T, want *ast.DropTableStmt", list.Items[1])
	}
}

// TestExecCommentNestingSplit: Split ends nested executable comments where the
// lexer does, at the first "*/", so a ';' after it is a statement boundary
// even when a later "*/" would have closed the outer comment had they nested
// (TestExecCommentNestingSplit_Oracle proves the boundary against TiDB).
func TestExecCommentNestingSplit(t *testing.T) {
	sql := "SELECT 1 /*! + 1 /*! + 2 */ ; SELECT 2 */"
	segs := Split(sql)
	if len(segs) != 2 || segs[0].Text != "SELECT 1 /*! + 1 /*! + 2 */ " {
		t.Fatalf("Split(%q) = %+v, want the statement before the ';' and the one after it", sql, segs)
	}
}

func valuesLen(ins *ast.InsertStmt) int {
	if ins == nil {
		return 0
	}
	return len(ins.Values)
}

// TestFeatureGate pins where a /*T![feature] gate ends: at the first "*/" the
// lexer reads outside a token, past a nested executable comment's opener and a
// nested ordinary comment, or nowhere, when TiDB drops the opener.
func TestFeatureGate(t *testing.T) {
	cases := []struct {
		name  string
		sql   string
		gate  bool
		close int // offset of the closing "*/", or -1 for none
	}{
		{"plain gate", "/*T![clustered_index] SELECT 1 */", true, 31},
		{"quoted close does not end the gate", "/*T![clustered_index] SELECT '*/' */", true, 34},
		{"quoted opener adds no nesting", "/*T![clustered_index] SELECT '/*' */", true, 34},
		{"nested comment is consumed whole", "/*T![clustered_index] /* c */ SELECT 1 */", true, 39},
		{"nested executable opener is dropped", "/*T![clustered_index] SELECT 1 /*! + 2 */ */", true, 39},
		{"line comment takes the close with it", "/*T![clustered_index] SELECT 1 -- */", true, -1},
		{"no close at all", "/*T![clustered_index] SELECT 1", true, -1},
		{"unsupported feature is an ordinary comment", "/*T![no_such_feature] SELECT 1 */", false, 0},
		{"malformed feature list is no gate", "/*T![clustered_index ] SELECT 1 */", false, 0},
		{"no feature list", "/*T! SELECT 1 */", false, 0},
		{"MySQL executable comment", "/*! SELECT 1 */", false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, closeStart, ok := FeatureGate(tc.sql, 0)
			if ok != tc.gate {
				t.Fatalf("FeatureGate(%q) ok = %v, want %v", tc.sql, ok, tc.gate)
			}
			if !ok {
				return
			}
			if want := strings.Index(tc.sql, "]") + 1; body != want {
				t.Errorf("FeatureGate(%q) body = %d, want %d", tc.sql, body, want)
			}
			if closeStart != tc.close {
				t.Errorf("FeatureGate(%q) close = %d, want %d", tc.sql, closeStart, tc.close)
			}
		})
	}
}

// TestParseExprReportsLexingErrors: ParseExpr is as strict as Parse about a
// token that runs off the end of the input.
func TestParseExprReportsLexingErrors(t *testing.T) {
	for _, sql := range []string{"'abc", "a + 'abc", "`a", "a /* c", "1 /*! + 'x"} {
		if _, err := ParseExpr(sql); err == nil {
			t.Errorf("ParseExpr(%q) accepted an unterminated token", sql)
		}
	}
	if _, err := ParseExpr("a + INTERVAL 1 DAY"); err != nil {
		t.Errorf("ParseExpr: %v", err)
	}
}
