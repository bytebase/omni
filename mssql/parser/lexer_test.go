package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestLexerLongIdentifierAcceptedForReviewLayer pins that mssql/parser's lexer
// accepts identifiers longer than SQL Server's 128-character limit rather than
// reporting a syntax error. The length limit is intentionally enforced by the
// review layer (mssql/review WT403 / NameTooLong) so users receive a named
// lint finding instead of a parse failure.
func TestLexerLongIdentifierAcceptedForReviewLayer(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   string
		sql  string
	}{
		{
			name: "bare_128",
			id:   strings.Repeat("a", 128),
			sql:  "SELECT " + strings.Repeat("a", 128),
		},
		{
			name: "bare_129",
			id:   strings.Repeat("a", 129),
			sql:  "SELECT " + strings.Repeat("a", 129),
		},
		{
			name: "bracketed_129",
			id:   strings.Repeat("b", 129),
			sql:  "SELECT [" + strings.Repeat("b", 129) + "]",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q): unexpected error: %v", tc.sql, err)
			}
			if list == nil || len(list.Items) != 1 {
				t.Fatalf("Parse(%q): got %d statements, want 1", tc.sql, list.Len())
			}
			sel, ok := list.Items[0].(*ast.SelectStmt)
			if !ok || sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
				t.Fatalf("Parse(%q): unexpected AST %#v", tc.sql, list.Items[0])
			}
			rt, ok := sel.TargetList.Items[0].(*ast.ResTarget)
			if !ok {
				t.Fatalf("Parse(%q): target = %T, want *ast.ResTarget", tc.sql, sel.TargetList.Items[0])
			}
			col, ok := rt.Val.(*ast.ColumnRef)
			if !ok || col.Column != tc.id {
				t.Fatalf("Parse(%q): column = %#v, want Column=%q", tc.sql, rt.Val, tc.id)
			}
		})
	}
}
