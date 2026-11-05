package parser

import (
	"context"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// execCommentBoundaryCases show how the syntax of an executable comment
// bounds tokens. The lexer reads the body in place, so "/*!" and "*/" end the
// token before them the way whitespace would, and a token that swallows the
// comment's closing "*/" leaves the comment open, which the server rejects.
//
// MySQL 8.0 and OceanBase 4.3.5 agree on these and on the quoted-delimiter
// lists below: each statement either parses with the bodies in place, or
// fails the same way with them read as plain comments (exec_fallback.go), or
// parses with a comment left open, which OceanBase refuses without reading
// it again. So the MySQL 8.0 oracle guards them as well; the
// _OceanBase tests hold them to the engine itself. Executable comments whose
// bodies do not parse are OceanBase's own and answer to it alone, in
// execFallbackCases.
var execCommentBoundaryCases = []string{
	// "*/" ends the number and the identifier before it: SELECT 1 2 and
	// SELECT 'v' x y, both rejected, and the plain reading leaves 2 and y.
	"/*!SELECT 1*/2",
	"/*!SELECT 'v' x*/y",
	// A string or a line comment holding the closing "*/".
	"SELECT 'x' /*! , 'a*/'",
	"SELECT 1 /*!50000 + 1 -- c */\n",
	// Adjacent strings on either side of "*/" concatenate.
	"/*!SELECT 'a'*/'b'",
	"SELECT /*!50000 1, */ 2",
}

// TestExecCommentBoundary_Oracle proves execCommentBoundaryCases against a
// live MySQL 8.0 container.
func TestExecCommentBoundary_Oracle(t *testing.T) {
	o := startParserOracle(t)

	for _, sql := range execCommentBoundaryCases {
		assertOracleMatch(t, o, sql)
	}

	sql := "/*!SELECT 'a'*/'b'"
	var server string
	if err := o.db.QueryRowContext(o.ctx, sql).Scan(&server); err != nil {
		t.Fatalf("server: %q: %v", sql, err)
	}
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	lit, ok := list.Items[0].(*ast.SelectStmt).TargetList[0].(*ast.StringLit)
	if !ok || lit.Value != server {
		t.Errorf("Parse(%q) target = %#v, want the string %q the server returns", sql, list.Items[0].(*ast.SelectStmt).TargetList[0], server)
	}
}

// TestExecCommentQuotedDelimiters_Oracle proves against the live container that
// omni ends an executable comment where the server does when a "/*" or "*/"
// sits inside a string, a quoted identifier or a line comment.
func TestExecCommentQuotedDelimiters_Oracle(t *testing.T) {
	o := startParserOracle(t)

	for _, list := range [][]string{execCommentQuotedAccepted, execCommentQuotedHeldClose} {
		for _, sql := range list {
			assertOracleMatch(t, o, sql)
		}
	}
}

// TestExecCommentBoundary_OceanBase proves execCommentBoundaryCases and the
// quoted-delimiter lists against a live OceanBase 4.3.5 (DSN-gated, see
// oceanBaseOracle), through PREPARE.
func TestExecCommentBoundary_OceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, list := range [][]string{execCommentBoundaryCases, execCommentQuotedAccepted, execCommentQuotedHeldClose} {
		for _, sql := range list {
			engine := prepareAcceptsOceanBase(t, ctx, db, sql)
			if _, err := Parse(sql); (err == nil) != engine {
				t.Errorf("DIVERGENCE %q: omni accepts=%v (%v), OceanBase accepts=%v", sql, err == nil, err, engine)
			}
		}
	}
	for _, sql := range execCommentQuotedHeldClose {
		if prepareAcceptsOceanBase(t, ctx, db, sql) {
			t.Errorf("%q: OceanBase accepts a comment whose close is held inside a token", sql)
		}
	}
}
