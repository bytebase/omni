package v1review

import (
	"slices"
	"strconv"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementMergeAlterTable reports each table that more than one top-level
// statement touches, counting from the table's last CREATE TABLE (which counts
// one and discards earlier counts) or else from its first ALTER TABLE; each
// later ALTER TABLE adds one. Tables are keyed and printed by their bare name,
// case-sensitively, so db1.t and db2.t are one table. Findings are ordered by
// the table's last statement, carry no range, and set Statement to that
// statement's index (the V1 advisor printed its first content line).
// RedundantAlterTable is the only code; Value is the statements counted.
// Ported from Bytebase's V1 rule STATEMENT_MERGE_ALTER_TABLE.
func StatementMergeAlterTable(stmts []ast.Node) []v1review.Finding {
	type tableStatements struct {
		count int
		last  int
	}
	tables := map[string]*tableStatements{}
	var order []string
	for i, stmt := range stmts {
		var table *ast.TableRef
		creates := false
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table, creates = n.Table, true
		case *ast.AlterTableStmt:
			table = n.Table
		}
		if table == nil {
			continue
		}
		t, ok := tables[table.Name]
		if !ok {
			t = &tableStatements{}
			tables[table.Name] = t
			order = append(order, table.Name)
		}
		if creates {
			t.count = 0
		}
		t.count++
		t.last = i
	}
	var findings []v1review.Finding
	for _, name := range order {
		if t := tables[name]; t.count > 1 {
			findings = append(findings, v1review.Finding{Code: v1review.RedundantAlterTable, Table: name, Value: strconv.Itoa(t.count), Statement: t.last})
		}
	}
	slices.SortFunc(findings, func(a, b v1review.Finding) int { return a.Statement - b.Statement })
	return findings
}
