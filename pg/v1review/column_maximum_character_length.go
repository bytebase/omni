package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnMaximumCharacterLength judges the column types top-level statements
// write for a fixed-length character type longer than maximum: each column a
// CREATE TABLE defines, each column an ALTER ... ADD COLUMN adds, and the new
// type an ALTER COLUMN ... TYPE gives, the ALTER forms whatever relation they
// name. A type is judged when the last part of its name, lower-cased, is
// char, character, bpchar or nchar, so CHAR(n) and CHARACTER(n) are and
// VARCHAR(n) is not. Its length is its first type modifier when the parser
// stores that as an integer, which it does for CHAR(n), CHARACTER(n) and
// NCHAR(n) but not for a type written by a generic name such as bpchar(n),
// whose length is never read. A statement whose relation has no name is
// skipped, and a statement in a routine body is not judged.
//
// Each statement reports at most one ColumnCharTooLong, for its first such
// column. Name is the column, Table the relation's name and Value its schema
// as written, "" when unqualified. The finding anchors the whole statement
// and sets Search to the column name: Bytebase's V1 advisor reported the line
// of the name's first occurrence in the statement text. Ported from
// Bytebase's V1 rule COLUMN_MAXIMUM_CHARACTER_LENGTH.
func ColumnMaximumCharacterLength(stmts []ast.Node, maximum int) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var rel *ast.RangeVar
		var loc ast.Loc
		var columns []typedColumn
		switch n := statementNode(stmt).(type) {
		case *ast.CreateStmt:
			rel, loc = n.Relation, n.Loc
			for _, col := range tableColumns(n) {
				if col.TypeName != nil {
					columns = append(columns, typedColumn{col.Colname, col.TypeName})
				}
			}
		case *ast.AlterTableStmt:
			rel, loc = n.Relation, n.Loc
			columns = alteredColumnTypes(n)
		default:
			continue
		}
		if rel == nil || rel.Relname == "" {
			continue
		}
		for _, col := range columns {
			if fixedCharLength(col.typ) > maximum {
				start, end := span(loc)
				findings = append(findings, v1review.Finding{
					Code:   v1review.ColumnCharTooLong,
					Name:   col.name,
					Table:  rel.Relname,
					Value:  rel.Schemaname,
					Search: col.name,
					Start:  start,
					End:    end,
				})
				break
			}
		}
	}
	return findings
}

// typedColumn is a column a statement gives a type, with that type.
type typedColumn struct {
	name string
	typ  *ast.TypeName
}

// alteredColumnTypes returns, in order, the columns an ALTER adds with ADD
// COLUMN and those it retypes with ALTER COLUMN ... TYPE, each with the type
// it writes; an action without a type is left out.
func alteredColumnTypes(n *ast.AlterTableStmt) []typedColumn {
	var columns []typedColumn
	for _, cmd := range alterCommands(n) {
		col, ok := cmd.Def.(*ast.ColumnDef)
		if !ok || col == nil || col.TypeName == nil {
			continue
		}
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddColumn:
			columns = append(columns, typedColumn{col.Colname, col.TypeName})
		case ast.AT_AlterColumnType:
			columns = append(columns, typedColumn{cmd.Name, col.TypeName})
		default:
		}
	}
	return columns
}

// fixedCharLength returns a fixed-length character type's length, 0 for any
// other type or when the length is not an integer literal.
func fixedCharLength(tn *ast.TypeName) int {
	switch strings.ToLower(lastTypeName(tn)) {
	case "char", "character", "bpchar", "nchar":
	default:
		return 0
	}
	if tn.Typmods == nil || len(tn.Typmods.Items) == 0 {
		return 0
	}
	if n, ok := tn.Typmods.Items[0].(*ast.Integer); ok {
		return int(n.Ival)
	}
	return 0
}
