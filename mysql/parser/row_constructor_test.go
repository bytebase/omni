package parser

import (
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

func TestParseRowConstructor_Bare(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"(1, 2)", "{ROW :loc 0 :items ({INT_LIT :val 1 :loc 1} {INT_LIT :val 2 :loc 4})}"},
		{"(1, 2, 3)", "{ROW :loc 0 :items ({INT_LIT :val 1 :loc 1} {INT_LIT :val 2 :loc 4} {INT_LIT :val 3 :loc 7})}"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ast.NodeToString(parseExpr(t, tt.input))
			if got != tt.want {
				t.Errorf("parseExpr(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseRowConstructor_Nested(t *testing.T) {
	expr := parseExpr(t, "((1, 2), (3, 4))")
	row, ok := expr.(*ast.RowExpr)
	if !ok {
		t.Fatalf("expected *ast.RowExpr, got %T", expr)
	}
	if len(row.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(row.Items))
	}
	for i, item := range row.Items {
		if _, ok := item.(*ast.RowExpr); !ok {
			t.Errorf("item %d: expected *ast.RowExpr, got %T", i, item)
		}
	}
}

func TestParseRowConstructor_SingleStaysParen(t *testing.T) {
	got := ast.NodeToString(parseExpr(t, "(1)"))
	want := "{PAREN :loc 0 :expr {INT_LIT :val 1 :loc 1}}"
	if got != want {
		t.Errorf("parseExpr(\"(1)\"):\n  got:  %s\n  want: %s", got, want)
	}
}

func TestParseRowConstructor_Reject(t *testing.T) {
	for _, sql := range []string{
		"SELECT ()",
		"SELECT (1,)",
		"SELECT (1, 2,)",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q): expected parse error, got nil", sql)
			}
		})
	}
}

func TestParseRowConstructor_RowKeyword(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t WHERE ROW(a, b) = ROW(1, 2)",
		"SELECT * FROM t WHERE ROW(a, b) IN ((1, 2), (3, 4))",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): unexpected parse error: %v", sql, err)
			}
		})
	}
}

func TestParseRowConstructor_InStatements(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t WHERE (a, b) = (1, 2)",
		"SELECT * FROM t WHERE (a, b) IN ((1, 2), (3, 4))",
		"SELECT * FROM t WHERE (a, b) > (1, 2)",
		"SELECT * FROM t WHERE (a, b) NOT IN ((1, 2))",
		"SELECT * FROM t WHERE (a, b) <=> (1, 2)",
		"SELECT * FROM t WHERE (a, b) IN (SELECT a, b FROM u)",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): unexpected parse error: %v", sql, err)
			}
		})
	}
}
