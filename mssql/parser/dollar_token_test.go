package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestDollarPseudoColumns covers the $-prefixed pseudo-columns SQL Server
// accepts as column references: $action (MERGE OUTPUT, BYT-9813), $IDENTITY,
// $ROWGUID, and the graph-table pseudo-columns, in bare and table-qualified
// forms. Context restrictions ($action outside MERGE OUTPUT) are binding-time
// errors in the engine (Msg 207), not parse errors, so the parser accepts
// them anywhere a column reference is legal — same as SqlScriptDOM.
func TestDollarPseudoColumns(t *testing.T) {
	accept := []string{
		// $action in MERGE OUTPUT — the shape that motivated BYT-9813.
		`MERGE INTO dst AS d
USING src AS s ON d.k = s.k
WHEN MATCHED THEN UPDATE SET d.v = s.v
WHEN NOT MATCHED THEN INSERT (k, v) VALUES (s.k, s.v)
OUTPUT $action, INSERTED.k, INSERTED.v;`,
		"MERGE INTO t AS T USING s AS S ON T.Id = S.Id WHEN MATCHED THEN UPDATE SET T.Name = S.Name OUTPUT $ACTION;",
		"MERGE INTO t AS T USING s AS S ON T.Id = S.Id WHEN MATCHED THEN UPDATE SET T.Name = S.Name OUTPUT $action AS act, INSERTED.Id INTO @changes;",
		"INSERT INTO t (a) OUTPUT $action VALUES (1);",
		// Engine parses these; the "only in MERGE OUTPUT" restriction is Msg 207
		// at binding, not a syntax error.
		"SELECT $action FROM t;",
		// $IDENTITY / $ROWGUID pseudo-columns, bare and qualified.
		"SELECT $IDENTITY FROM t;",
		"SELECT $ROWGUID FROM t;",
		"SELECT t.$IDENTITY FROM t;",
		"SELECT t.$rowguid FROM t;",
		// Graph pseudo-columns.
		"SELECT $node_id FROM Person;",
		"SELECT Person.$node_id FROM Person;",
		"SELECT $from_id, $to_id FROM Likes;",
		"SELECT $edge_id FROM Likes;",
		// $PARTITION system partition function.
		"SELECT $PARTITION.pf1(10);",
		"SELECT $partition.pf1(o.OrderDate) FROM Orders o;",
		// `$` inside identifiers and bracketed names are ordinary columns.
		"SELECT a$b FROM t;",
		"SELECT [$action] FROM t;",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): unexpected error: %v", sql, err)
			}
		})
	}

	// Unknown pseudo-columns are parse errors, matching the engine's Msg 126
	// "Invalid pseudocolumn". A bare `$` is rejected (SqlScriptDOM behavior;
	// the engine lexes a lone `$` as money 0 — deliberate divergence, see
	// lexDollar). $PARTITION requires the .function(args) suffix.
	reject := []string{
		"SELECT $foo FROM t;",
		"SELECT t.$foo FROM t;",
		"SELECT $ FROM t;",
		"SELECT $PARTITION FROM t;",
		"SELECT $PARTITION.pf1 FROM t;",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q): expected parse error, got nil", sql)
			}
		})
	}
}

// TestDollarPseudoColumnAST pins the AST shape: pseudo-columns are ColumnRef
// nodes (reusing the existing node — no new consumer surface).
func TestDollarPseudoColumnAST(t *testing.T) {
	expr := parseExprT(t, "$action")
	ref, ok := expr.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("expected *ast.ColumnRef, got %T", expr)
	}
	if ref.Column != "$action" {
		t.Errorf("Column = %q, want %q", ref.Column, "$action")
	}

	expr = parseExprT(t, "t.$IDENTITY")
	ref, ok = expr.(*ast.ColumnRef)
	if !ok {
		t.Fatalf("expected *ast.ColumnRef, got %T", expr)
	}
	if ref.Table != "t" || ref.Column != "$IDENTITY" {
		t.Errorf("got Table=%q Column=%q, want t / $IDENTITY", ref.Table, ref.Column)
	}
}

// TestMoneyLiterals covers T-SQL money constants: a currency symbol followed
// by digits with an optional decimal point. Non-$ Unicode currency symbols
// (£, €, ¥) are also money constants in the engine — previously they were
// mis-lexed as identifiers.
func TestMoneyLiterals(t *testing.T) {
	accept := []string{
		"SELECT $12;",
		"SELECT $12.5;",
		"SELECT -$4.78;",
		"SELECT $.5;",
		"SELECT £10;",
		"SELECT €7;",
		"SELECT ¥100;",
		"INSERT INTO t (price) VALUES ($19.99);",
		"SELECT * FROM t WHERE price > $100;",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): unexpected error: %v", sql, err)
			}
		})
	}
}

// TestMoneyLiteralAST pins the AST shape: LitMoney with the raw text
// (currency symbol included) in Str.
func TestMoneyLiteralAST(t *testing.T) {
	for input, want := range map[string]string{
		"$12.50": "$12.50",
		"£10":    "£10",
	} {
		expr := parseExprT(t, input)
		lit, ok := expr.(*ast.Literal)
		if !ok {
			t.Fatalf("parseExpr(%q): expected *ast.Literal, got %T", input, expr)
		}
		if lit.Type != ast.LitMoney {
			t.Errorf("parseExpr(%q): Type = %d, want LitMoney", input, lit.Type)
		}
		if lit.Str != want {
			t.Errorf("parseExpr(%q): Str = %q, want %q", input, lit.Str, want)
		}
	}
}

// TestIdentityColRowGuidCol covers the keyword-only column references
// IDENTITYCOL and ROWGUIDCOL, bare and table-qualified.
func TestIdentityColRowGuidCol(t *testing.T) {
	accept := []string{
		"SELECT IDENTITYCOL FROM t;",
		"SELECT ROWGUIDCOL FROM t;",
		"SELECT t.IDENTITYCOL FROM t;",
		"SELECT t.ROWGUIDCOL FROM t;",
		"SELECT * FROM t WHERE IDENTITYCOL = 5;",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): unexpected error: %v", sql, err)
			}
		})
	}
}

// parseExprT parses `SELECT <input>` and returns the first target expression.
func parseExprT(t *testing.T, input string) ast.ExprNode {
	t.Helper()
	list, err := Parse("SELECT " + input)
	if err != nil {
		t.Fatalf("Parse(SELECT %s): %v", input, err)
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt, got %T", list.Items[0])
	}
	target, ok := sel.TargetList.Items[0].(*ast.ResTarget)
	if !ok {
		t.Fatalf("expected *ast.ResTarget, got %T", sel.TargetList.Items[0])
	}
	expr, ok := target.Val.(ast.ExprNode)
	if !ok {
		t.Fatalf("expected ast.ExprNode, got %T", target.Val)
	}
	return expr
}
