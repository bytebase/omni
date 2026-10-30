package parser

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// TestExecCommentBoundary_Oracle proves against a live MySQL 8.0 container how
// the syntax of an executable comment bounds tokens. The lexer reads the body
// in place, so "/*!" and "*/" end the token before them the way whitespace
// would, and a token that swallows the comment's closing "*/" leaves the
// comment open, which MySQL rejects.
func TestExecCommentBoundary_Oracle(t *testing.T) {
	o := startParserOracle(t)

	for _, sql := range []string{
		// "*/" ends the number and the identifier before it: SELECT 1 2 and
		// SELECT 'v' x y, both rejected.
		"/*!SELECT 1*/2",
		"/*!SELECT 'v' x*/y",
		// A string or a line comment holding the closing "*/".
		"SELECT 'x' /*! , 'a*/'",
		"SELECT 1 /*!50000 + 1 -- c */\n",
		// Adjacent strings on either side of "*/" concatenate.
		"/*!SELECT 'a'*/'b'",
		"SELECT /*!50000 1, */ 2",
	} {
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
