package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDefaultDisallowVolatile judges the columns top-level ALTER ... ADD
// COLUMN statements add, whatever relation they name, for a DEFAULT holding a
// function call anywhere in its expression. Any call counts, volatile or
// not: now() and lower('a') both do, CURRENT_TIMESTAMP does not. The DEFAULT
// is read from the column's raw default and from each DEFAULT constraint it
// carries. CREATE TABLE is not judged, nor is an ALTER whose relation has no
// name, nor a statement in a routine body.
//
// Each finding reports ColumnDefaultVolatile with Name the column and Table
// the relation's name without its schema, in the order the columns are
// added. It anchors the whole statement and sets Search to the column name:
// Bytebase's V1 advisor reported the line of the name's first occurrence in
// the statement text. Ported from Bytebase's V1 rule
// COLUMN_DEFAULT_DISALLOW_VOLATILE.
func ColumnDefaultDisallowVolatile(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := statementNode(stmt).(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		table := relname(n.Relation)
		if table == "" {
			continue
		}
		for _, cmd := range alterCommands(n) {
			if ast.AlterTableType(cmd.Subtype) != ast.AT_AddColumn {
				continue
			}
			col, ok := cmd.Def.(*ast.ColumnDef)
			if !ok || col == nil || !hasFunctionCallDefault(col) {
				continue
			}
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{
				Code:   v1review.ColumnDefaultVolatile,
				Name:   col.Colname,
				Table:  table,
				Search: col.Colname,
				Start:  start,
				End:    end,
			})
		}
	}
	return findings
}

func hasFunctionCallDefault(col *ast.ColumnDef) bool {
	if col.RawDefault != nil && containsFunctionCall(col.RawDefault) {
		return true
	}
	for _, c := range columnConstraints(col) {
		if c.Contype == ast.CONSTR_DEFAULT && c.RawExpr != nil && containsFunctionCall(c.RawExpr) {
			return true
		}
	}
	return false
}

func containsFunctionCall(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if _, ok := n.(*ast.FuncCall); ok {
			found = true
			return false
		}
		return true
	})
	return found
}
