package v1review

import (
	"slices"
	"strconv"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexTotalNumberLimit reports the tables a script may have given an index
// that end up with more than limit indexes. final is the schema after the
// script, as Bytebase simulated it; the count is the number of indexes it
// holds for the table, the primary key among them. A limit of 0 reports
// every such table with an index, a negative one every such table.
//
// A top-level statement marks its table when it is a CREATE INDEX of any
// kind; a CREATE TABLE that defines a PRIMARY KEY or UNIQUE constraint, on a
// column or on the table; or an ALTER TABLE with an ADD COLUMN whose column
// has such a constraint or an ADD of such a table constraint. A CREATE TABLE
// without one, an added EXCLUDE or FOREIGN KEY constraint, CREATE TABLE AS,
// and a statement in a CREATE SCHEMA or in a routine body do not mark it.
//
// Tables are marked by schema and name as the parser stores them, the
// schema "public" when the statement names none; each is looked up in final
// as Schema.Table does, and skipped when final does not hold it. When the
// script marks any table and final has no snapshot the function returns
// v1review.ErrSchemaMissing.
//
// Schema and Table are the marked names and Value the number of the table's indexes. A
// finding lies in the last statement that marked its table, held in
// Statement, and anchors that whole statement; Bytebase's V1 advisor
// reported the line of that statement's last character that is not blank or
// a semicolon. Findings are in the order of those statements; the V1
// advisor ordered them by line, leaving the order of two on one line to
// chance. Ported from Bytebase's V1 rule INDEX_TOTAL_NUMBER_LIMIT.
func IndexTotalNumberLimit(stmts []ast.Node, limit int, final v1review.Schema) ([]v1review.Finding, error) {
	type table struct{ schema, name string }
	type mark struct {
		table     table
		statement int
		loc       ast.Loc
	}
	marked := make(map[table]mark)
	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		var rv *ast.RangeVar
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if createsKey(n) {
				rv, loc = n.Relation, n.Loc
			}
		case *ast.AlterTableStmt:
			if addsKey(n) {
				rv, loc = n.Relation, n.Loc
			}
		case *ast.IndexStmt:
			rv, loc = n.Relation, n.Loc
		}
		if name := relname(rv); name != "" {
			t := table{schemaOrDefault(rv), name}
			marked[t] = mark{t, i, loc}
		}
	}
	if len(marked) == 0 {
		return nil, nil
	}
	if final.Metadata == nil {
		return nil, v1review.ErrSchemaMissing
	}
	marks := make([]mark, 0, len(marked))
	for _, m := range marked {
		marks = append(marks, m)
	}
	slices.SortFunc(marks, func(a, b mark) int { return a.statement - b.statement })
	var findings []v1review.Finding
	for _, m := range marks {
		indexes, found := final.TableIndexes(m.table.schema, m.table.name)
		if !found || len(indexes) <= limit {
			continue
		}
		start, end := span(m.loc)
		findings = append(findings, v1review.Finding{Code: v1review.IndexCountExceeded, Schema: m.table.schema, Table: m.table.name, Value: strconv.Itoa(len(indexes)), Statement: m.statement, Start: start, End: end})
	}
	return findings, nil
}

// isKey reports whether a constraint is one the V1 total index rule took to
// add an index.
func isKey(c *ast.Constraint) bool {
	return c.Contype == ast.CONSTR_PRIMARY || c.Contype == ast.CONSTR_UNIQUE
}

// createsKey reports whether a CREATE TABLE defines a PRIMARY KEY or UNIQUE
// constraint, on a column or on the table.
func createsKey(n *ast.CreateStmt) bool {
	if n.TableElts == nil {
		return false
	}
	for _, item := range n.TableElts.Items {
		switch e := item.(type) {
		case *ast.ColumnDef:
			if slices.ContainsFunc(columnConstraints(e), isKey) {
				return true
			}
		case *ast.Constraint:
			if isKey(e) {
				return true
			}
		}
	}
	return false
}

// addsKey reports whether an ALTER TABLE adds a column with a PRIMARY KEY or
// UNIQUE constraint, or adds such a table constraint.
func addsKey(n *ast.AlterTableStmt) bool {
	if n.Cmds == nil {
		return false
	}
	for _, item := range n.Cmds.Items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddColumn:
			if col, ok := cmd.Def.(*ast.ColumnDef); ok && slices.ContainsFunc(columnConstraints(col), isKey) {
				return true
			}
		case ast.AT_AddConstraint:
			if c, ok := cmd.Def.(*ast.Constraint); ok && isKey(c) {
				return true
			}
		}
	}
	return false
}
