package v1review

import (
	"slices"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementMergeAlterTable reports each table that more than one top-level
// statement alters. Tables are keyed by name alone, case-sensitively, with
// any database qualifier dropped. An ALTER TABLE counts one; a CREATE TABLE
// resets the table's count to one, discarding earlier ALTERs. Findings come
// in the order of each table's last counted statement. Table is the name,
// Count the number of statements, and Statement the index of the last of
// them; findings anchor nothing within it (the V1 advisor placed them on
// that statement's first line). RedundantAlterTable is the only code. Ported
// from Bytebase's V1 rule STATEMENT_MERGE_ALTER_TABLE, which sorted by that
// statement's line and so ordered tables whose last statements share a line
// nondeterministically; this keeps statement order.
func StatementMergeAlterTable(stmts []ast.Node) []v1review.Finding {
	type table struct {
		count, last int
	}
	tables := map[string]*table{}
	var order []string
	note := func(name string, i int, reset bool) {
		t, ok := tables[name]
		if !ok {
			t = &table{}
			tables[name] = t
			order = append(order, name)
		}
		if reset {
			t.count = 0
		}
		t.count++
		t.last = i
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table != nil {
				note(n.Table.Name, i, true)
			}
		case *ast.AlterTableStmt:
			if n.Table != nil {
				note(n.Table.Name, i, false)
			}
		}
	}
	var findings []v1review.Finding
	for _, name := range order {
		if t := tables[name]; t.count > 1 {
			findings = append(findings, v1review.Finding{Code: v1review.RedundantAlterTable, Table: name, Count: t.count, Statement: t.last})
		}
	}
	slices.SortFunc(findings, func(a, b v1review.Finding) int { return a.Statement - b.Statement })
	return findings
}
