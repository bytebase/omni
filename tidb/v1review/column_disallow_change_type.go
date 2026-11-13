package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowChangeType reports each top-level ALTER TABLE statement
// naming a table that changes a column's type, once however many it
// changes, anchored at the whole statement.
//
// The actions judged are CHANGE COLUMN, whose old name is looked up, and
// MODIFY COLUMN, whose column's name is; an action whose definition has no
// type is not judged. The column is looked up in schema, in the schema named
// "", by the table's name without a qualifier, column names compared without
// regard to case; a column schema lacks is not judged. Its type changes when
// its stored type and the new type differ once each is lower-cased and a
// bare integer type is given MySQL's display width, BOOL, BOOLEAN and a bare
// DECIMAL, NUMERIC or FIXED canonicalized, as is a bare BIT, BINARY or YEAR.
// The new type is rendered as ColumnTypeDisallowList renders it, followed by
// " unsigned" when UNSIGNED or ZEROFILL and " zerofill" when ZEROFILL. Ported
// from Bytebase's V1 rule COLUMN_DISALLOW_CHANGE_TYPE.
func ColumnDisallowChangeType(stmts []ast.Node, schema v1review.Schema) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := stmt.(*ast.AlterTableStmt)
		if !ok || alter.Table == nil {
			continue
		}
		for _, cmd := range alter.Commands {
			if typeChanged(schema, alter.Table.Name, cmd) {
				start, end := span(alter.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.ColumnTypeChanged, Table: alter.Table.Name, Start: start, End: end})
				break
			}
		}
	}
	return findings
}

// typeChanged reports whether a CHANGE or MODIFY COLUMN action gives a column
// schema holds another type.
func typeChanged(schema v1review.Schema, table string, cmd *ast.AlterTableCmd) bool {
	if cmd == nil || cmd.Column == nil || cmd.Column.TypeName == nil {
		return false
	}
	var name string
	switch cmd.Type {
	case ast.ATChangeColumn:
		name = cmd.Name
	case ast.ATModifyColumn:
		name = cmd.Column.Name
	default:
		return false
	}
	if name == "" {
		return false
	}
	column := schema.TableColumn("", table, name, false)
	if column == nil {
		return false
	}
	return normalizeColumnType(column.GetType()) != normalizeColumnType(columnType(cmd.Column.TypeName))
}

// columnType renders a type in lower case with its UNSIGNED and ZEROFILL, ""
// for no type.
func columnType(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	rendered := compactType(dt)
	if dt.Unsigned || dt.Zerofill {
		rendered += " unsigned"
	}
	if dt.Zerofill {
		rendered += " zerofill"
	}
	return rendered
}

// normalizeColumnType lower-cases a type and canonicalizes the bare forms
// MySQL stores with a width, precision or alias resolved.
func normalizeColumnType(tp string) string {
	switch lower := strings.ToLower(tp); lower {
	case "boolean", "bool":
		return "tinyint(1)"
	case "decimal", "numeric", "fixed":
		return "decimal(10,0)"
	case "bit":
		return "bit(1)"
	case "binary":
		return "binary(1)"
	case "year":
		return "year(4)"
	case "tinyint":
		return "tinyint(4)"
	case "tinyint unsigned":
		return "tinyint(4) unsigned"
	case "smallint":
		return "smallint(6)"
	case "smallint unsigned":
		return "smallint(6) unsigned"
	case "mediumint":
		return "mediumint(9)"
	case "mediumint unsigned":
		return "mediumint(9) unsigned"
	case "int":
		return "int(11)"
	case "int unsigned":
		return "int(11) unsigned"
	case "bigint":
		return "bigint(20)"
	case "bigint unsigned":
		return "bigint(20) unsigned"
	default:
		return lower
	}
}
