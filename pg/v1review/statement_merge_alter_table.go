package v1review

import (
	"cmp"
	"slices"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementMergeAlterTable judges the statements together and reports each
// table that a CREATE TABLE creates and later statements alter more than
// once in all: one RedundantAlterTable finding per table, with Count the
// CREATE and its ALTERs, Table the name the CREATE gave, and Statement the
// last ALTER. An ALTER counts toward the most recent CREATE of the table,
// a second CREATE starts the count over, and an ALTER of a table no
// earlier statement created is not counted. A table is keyed by its
// schema, public when unqualified, and its name joined with ".", so names
// holding a "." can collide, as they did in Bytebase's V1 advisor; and an
// ALTER INDEX, ALTER SEQUENCE, ALTER VIEW or other ALTER that shares the
// ALTER TABLE node counts too, while a RENAME does not. Findings come in the order of their
// Statement and anchor nothing: the V1 advisor placed each on its last
// statement's last line, relative to that statement, and sorted them by
// that line. Ported from Bytebase's V1 rule STATEMENT_MERGE_ALTER_TABLE.
func StatementMergeAlterTable(stmts []ast.Node) []v1review.Finding {
	type table struct {
		name      string
		count     int
		statement int
	}
	tables := make(map[string]table)
	for i, stmt := range stmts {
		switch n := unwrapRaw(stmt).(type) {
		case *ast.CreateStmt:
			if n.Relation != nil && n.Relation.Relname != "" {
				tables[tableKey(n.Relation)] = table{name: n.Relation.Relname, count: 1, statement: i}
			}
		case *ast.AlterTableStmt:
			if n.Relation == nil || n.Relation.Relname == "" {
				continue
			}
			key := tableKey(n.Relation)
			if t, ok := tables[key]; ok {
				t.count++
				t.statement = i
				tables[key] = t
			}
		}
	}
	var findings []v1review.Finding
	for _, t := range tables {
		if t.count > 1 {
			findings = append(findings, v1review.Finding{Code: v1review.RedundantAlterTable, Table: t.name, Count: t.count, Statement: t.statement})
		}
	}
	// Each ALTER names one table, so no two findings share a Statement.
	slices.SortFunc(findings, func(a, b v1review.Finding) int { return cmp.Compare(a.Statement, b.Statement) })
	return findings
}

// tableKey is the key Bytebase's V1 advisor kept a table under: its schema,
// public when unqualified, and its name joined with ".".
func tableKey(rv *ast.RangeVar) string {
	schema := rv.Schemaname
	if schema == "" {
		schema = "public"
	}
	return schema + "." + rv.Relname
}
