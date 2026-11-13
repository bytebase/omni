package v1review

import (
	"slices"
	"strconv"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexTotalNumberLimit reports the tables a script may have given an index
// that end up with more than limit indexes. final is the schema after the
// script, as Bytebase simulated it; the count is the number of indexes it
// holds for the table, the primary key among them. A limit of 0 reports
// every such table with an index, a negative one every such table.
//
// A top-level statement marks its table when it is a CREATE TABLE, with or
// without keys; a CREATE INDEX of any kind; or an ALTER TABLE that adds,
// modifies or changes a column to one with the PRIMARY KEY or UNIQUE
// attribute, or adds a PRIMARY KEY, UNIQUE, INDEX or FULLTEXT definition. An
// added SPATIAL or VECTOR index or FOREIGN KEY, a key inside the
// parenthesized ADD form, and a statement that only drops or renames do not
// mark it.
//
// Tables are marked by name as written, without the database qualifier, so
// two spellings of one table are two findings; each is looked up in final,
// in the schema of the empty name, as Schema.Table does, and skipped when
// final does not hold it. When the script marks any table and final has no
// snapshot the function returns v1review.ErrSchemaMissing.
//
// Table is the name as written and Value the number of the table's indexes. A finding
// lies in the last statement that marked its table and has no range: V1
// anchored it at the first character of that statement's text that is not
// white space. Findings are in the order of those statements; the V1
// advisor ordered them by that line, leaving the order of two on one line
// to chance. Ported from Bytebase's V1 rule INDEX_TOTAL_NUMBER_LIMIT.
func IndexTotalNumberLimit(stmts []ast.Node, limit int, final v1review.Schema) ([]v1review.Finding, error) {
	type mark struct {
		table     string
		statement int
	}
	marked := make(map[string]mark)
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table != nil {
				marked[n.Table.Name] = mark{n.Table.Name, i}
			}
		case *ast.CreateIndexStmt:
			if n.Table != nil {
				marked[n.Table.Name] = mark{n.Table.Name, i}
			}
		case *ast.AlterTableStmt:
			if n.Table != nil && slices.ContainsFunc(n.Commands, addsIndex) {
				marked[n.Table.Name] = mark{n.Table.Name, i}
			}
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
		indexes, found := final.TableIndexes("", m.table)
		if !found || len(indexes) <= limit {
			continue
		}
		findings = append(findings, v1review.Finding{Code: v1review.IndexCountExceeded, Table: m.table, Value: strconv.Itoa(len(indexes)), Statement: m.statement})
	}
	return findings, nil
}

// addsIndex reports whether an ALTER TABLE action is one the V1 total index
// rule took to add an index.
func addsIndex(cmd *ast.AlterTableCmd) bool {
	if cmd == nil {
		return false
	}
	keyed := func(column *ast.ColumnDef) bool {
		return column != nil && slices.ContainsFunc(column.Constraints, func(c *ast.ColumnConstraint) bool {
			return c != nil && (c.Type == ast.ColConstrPrimaryKey || c.Type == ast.ColConstrUnique)
		})
	}
	switch cmd.Type {
	case ast.ATAddColumn:
		return slices.ContainsFunc(addedColumns(cmd), keyed)
	case ast.ATModifyColumn, ast.ATChangeColumn:
		return keyed(cmd.Column)
	case ast.ATAddConstraint, ast.ATAddIndex:
		if cmd.Constraint == nil {
			return false
		}
		switch cmd.Constraint.Type {
		case ast.ConstrPrimaryKey, ast.ConstrUnique, ast.ConstrIndex, ast.ConstrFulltextIndex:
			return true
		default:
		}
	default:
	}
	return false
}
