package v1review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// columnAnchor is the node a cross-statement finding takes its line from,
// and the index of the statement holding it.
type columnAnchor struct {
	statement int
	loc       ast.Loc
}

// objectName is a table's key as the V1 column rules built it: schema.name,
// the schema defaulting to currentDatabase, and the bare name when both are
// empty. A nil name is the empty string.
func objectName(name *ast.ObjectName, currentDatabase string) string {
	if name == nil {
		return ""
	}
	schema := name.Schema
	if schema == "" {
		schema = currentDatabase
	}
	if schema == "" {
		return name.Name
	}
	return schema + "." + name.Name
}

func columnDefs(list *ast.List) []*ast.ColumnDef {
	if list == nil {
		return nil
	}
	result := make([]*ast.ColumnDef, 0, len(list.Items))
	for _, item := range list.Items {
		if col, ok := item.(*ast.ColumnDef); ok {
			result = append(result, col)
		}
	}
	return result
}

// actionColumns is the columns an ALTER TABLE action defines: its column
// list, then its single column, which may be nil.
func actionColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	return append(columnDefs(cmd.ColumnDefs), cmd.ColumnDef)
}

func alterTableCmds(stmt *ast.AlterTableStmt) []*ast.AlterTableCmd {
	if stmt.Actions == nil {
		return nil
	}
	result := make([]*ast.AlterTableCmd, 0, len(stmt.Actions.Items))
	for _, item := range stmt.Actions.Items {
		if cmd, ok := item.(*ast.AlterTableCmd); ok {
			result = append(result, cmd)
		}
	}
	return result
}

func columnHasConstraint(col *ast.ColumnDef, typ ast.ConstraintType) bool {
	if col.Constraints == nil {
		return false
	}
	for _, item := range col.Constraints.Items {
		if c, ok := item.(*ast.ColumnConstraint); ok && c.Type == typ {
			return true
		}
	}
	return false
}

// typeName is a type's name parts joined by dots and upper-cased.
func typeName(tn *ast.TypeName) string {
	return strings.ToUpper(strings.Join(listStrings(tn.Names), "."))
}

// firstTypeModInt is a type's first modifier as an integer: an integer
// literal, or a numeric literal whose whole text is an integer.
func firstTypeModInt(tn *ast.TypeName) (int, bool) {
	if tn.TypeMods == nil || len(tn.TypeMods.Items) == 0 {
		return 0, false
	}
	switch n := tn.TypeMods.Items[0].(type) {
	case *ast.Integer:
		return int(n.Ival), true
	case *ast.Float:
		v, err := strconv.Atoi(n.Fval)
		return v, err == nil
	}
	return 0, false
}

// lastIdentifier is the text after a key's last dot, which is what V1
// printed as a column's name, so a column name holding a dot prints only
// its tail.
func lastIdentifier(key string) string {
	return key[strings.LastIndex(key, ".")+1:]
}
