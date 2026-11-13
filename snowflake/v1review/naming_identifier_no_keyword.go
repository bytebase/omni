package v1review

import (
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIdentifierNoKeyword judges the column names top-level statements
// introduce against the V1 keyword table (v1Keywords): each column a
// CREATE TABLE declares, CREATE TABLE ... AS SELECT included, and the new
// name an ALTER TABLE ... RENAME COLUMN gives. No other identifier is
// judged: not a table, constraint or index name, not a column that ALTER
// TABLE ... ADD COLUMN adds, and nothing inside a Snowflake Scripting
// block, a task's body, or under EXPLAIN.
//
// A name is matched as Snowflake resolves it and then upper-cased, so a
// quoted name matches whatever its case: "connect" is reported though
// Snowflake keeps it distinct from the keyword. Finding.Name is not the
// matched text but the name in its source form: an unquoted name in the
// case it was written in, a quoted name with its quotes.
//
// A CREATE TABLE finding anchors the first item of the column list, be it
// a column, an out-of-line constraint or an inline index, so every finding
// from one statement shares that anchor. A RENAME COLUMN finding anchors
// the new name. NameIsKeyword is the only code. Ported from Bytebase's V1
// rule NAMING_IDENTIFIER_NO_KEYWORD.
func NamingIdentifierNoKeyword(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	declaredColumns(stmts, func(name ast.Ident, anchor ast.Loc) {
		if !isV1Keyword(name.Normalize()) {
			return
		}
		start, end := span(anchor)
		findings = append(findings, v1review.Finding{Code: v1review.NameIsKeyword, Name: name.String(), Start: start, End: end})
	})
	return findings
}

// declaredColumns calls visit, in statement order, with each column name
// the V1 identifier rules judge and the range they anchor it at; see
// NamingIdentifierNoKeyword for both.
func declaredColumns(stmts []ast.Node, visit func(name ast.Ident, anchor ast.Loc)) {
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			anchor := columnListLoc(n)
			for _, column := range n.Columns {
				visit(column.Name, anchor)
			}
		case *ast.AlterTableStmt:
			for _, a := range n.Actions {
				if a == nil || a.Kind != ast.AlterTableRenameColumn {
					continue
				}
				visit(a.NewColName, a.NewColName.Loc)
			}
		}
	}
}

// columnListLoc returns the range of the item that starts first in a
// CREATE TABLE column list, among its columns, out-of-line constraints and
// inline indexes, and NoLoc when none has a known start.
func columnListLoc(stmt *ast.CreateTableStmt) ast.Loc {
	first := ast.NoLoc()
	merge := func(loc ast.Loc) {
		if loc.Start >= 0 && (first.Start < 0 || loc.Start < first.Start) {
			first = loc
		}
	}
	for _, column := range stmt.Columns {
		merge(column.Loc)
	}
	for _, constraint := range stmt.Constraints {
		merge(constraint.Loc)
	}
	for _, index := range stmt.Indexes {
		merge(index.Loc)
	}
	return first
}
