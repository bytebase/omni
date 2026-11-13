package v1review

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// identifier is a name a statement defines and the Loc the V1 advisor took
// its line from.
type identifier struct {
	name string
	loc  ast.Loc
}

// identifiers returns the names a statement defines, in the order a
// depth-first walk meets them: created objects and their columns,
// constraints, and partitions; aliases and CTEs; rename targets; and PL/SQL
// declarations, nested statements included. References are not collected,
// nor are a column's named NOT NULL constraint, a subpartition, or a
// cluster's columns.
// A name is as the parser stores it — an unquoted name upper-cased, a
// quoted one as written — and without its schema qualifier.
//
// The Loc is the one the V1 advisor took its line from. It is the name's
// own for a created object held as an *ast.ObjectName, schema qualifier
// included, and for a table alias; otherwise it is the node that declares
// the name, such as the column definition, the constraint, the RENAME
// clause, or the PL/SQL declaration. A view's column aliases and its WITH
// CHECK OPTION constraint carry the CREATE VIEW statement's Loc, as does a
// created object whose name the AST holds as a string.
//
// A name is collected once per statement for each Loc start it is found
// at, so a CTE named like one of its columns yields one identifier.
func identifiers(node ast.Node) []identifier {
	var result []identifier
	seen := make(map[string]bool)
	add := func(name string, loc ast.Loc) {
		if name == "" {
			return
		}
		key := fmt.Sprintf("%s:%d", name, loc.Start)
		if seen[key] {
			return
		}
		seen[key] = true
		result = append(result, identifier{name: name, loc: loc})
	}
	addObject := func(name *ast.ObjectName) {
		if name != nil {
			add(name.Name, name.Loc)
		}
	}
	addStrings := func(list *ast.List, loc ast.Loc) {
		for _, name := range listStrings(list) {
			add(name, loc)
		}
	}
	// A TableConstraint node also names the target of DROP CONSTRAINT, so
	// constraint names are taken from the statements that define them.
	addTableConstraints := func(list *ast.List) {
		for _, c := range tableConstraints(list) {
			add(c.Name, c.Loc)
		}
	}
	ast.Inspect(node, func(n ast.Node) bool {
		add(createdName(n))
		switch x := n.(type) {
		case *ast.CreateTableStmt:
			addTableConstraints(x.Constraints)
		case *ast.CreateViewStmt:
			addStrings(x.Columns, x.Loc)
			addTableConstraints(x.Constraints)
			add(x.ConstraintName, x.Loc)
		case *ast.AdminDDLStmt:
			if x.Action == "CREATE" {
				addObject(x.Name)
			}
		case *ast.RenameStmt:
			addObject(x.NewName)
		case *ast.AlterTableCmd:
			switch x.Action {
			case ast.AT_ADD_CONSTRAINT:
				if x.Constraint != nil {
					add(x.Constraint.Name, x.Constraint.Loc)
				}
			case ast.AT_RENAME, ast.AT_RENAME_COLUMN, ast.AT_RENAME_CONSTRAINT, ast.AT_RENAME_PARTITION:
				add(x.NewName, x.Loc)
			}
		case *ast.AlterIndexStmt:
			add(x.NewName, x.Loc)
		case *ast.AlterTriggerStmt:
			add(x.NewName, x.Loc)
		case *ast.ColumnDef:
			add(x.Name, x.Loc)
		case *ast.ColumnConstraint:
			add(x.Name, x.Loc)
		case *ast.PartitionDef:
			add(x.Name, x.Loc)
		case *ast.ResTarget:
			add(x.Name, x.Loc)
		case *ast.Alias:
			add(x.Name, x.Loc)
			addStrings(x.Cols, x.Loc)
		case *ast.CTE:
			add(x.Name, x.Loc)
			addStrings(x.Columns, x.Loc)
		case *ast.Parameter:
			add(x.Name, x.Loc)
		case *ast.PLSQLVarDecl:
			add(x.Name, x.Loc)
		case *ast.PLSQLCursorDecl:
			add(x.Name, x.Loc)
		case *ast.PLSQLTypeDecl:
			add(x.Name, x.Loc)
		case *ast.PLSQLSubtypeDecl:
			add(x.Name, x.Loc)
		case *ast.PLSQLLoop:
			// Loops with an Iteration declare their variables as PLSQLIterand nodes.
			if x.Iteration == nil {
				add(x.Iterator, x.Loc)
			}
		case *ast.PLSQLIterand:
			add(x.Name, x.Loc)
		case *ast.PLSQLForall:
			add(x.Index, x.Loc)
		}
		return true
	})
	return result
}

// createdName returns the Name field of any Create*Stmt node, found by the
// node's type name as the V1 advisor found it, so every CREATE statement
// the AST has is judged without being listed here. A Create*Stmt with no
// Name field yields nothing.
func createdName(node ast.Node) (string, ast.Loc) {
	value := reflect.ValueOf(node)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return "", ast.Loc{}
	}
	elem := value.Elem()
	if elem.Kind() != reflect.Struct {
		return "", ast.Loc{}
	}
	typeName := elem.Type().Name()
	if !strings.HasPrefix(typeName, "Create") || !strings.HasSuffix(typeName, "Stmt") {
		return "", ast.Loc{}
	}
	nameField := elem.FieldByName("Name")
	if !nameField.IsValid() {
		return "", ast.Loc{}
	}
	switch name := nameField.Interface().(type) {
	case *ast.ObjectName:
		if name != nil {
			return name.Name, name.Loc
		}
	case string:
		if locField := elem.FieldByName("Loc"); locField.IsValid() {
			if loc, ok := locField.Interface().(ast.Loc); ok {
				return name, loc
			}
		}
		return name, ast.Loc{}
	}
	return "", ast.Loc{}
}

// listStrings returns the text of a list's *ast.String and *ast.ColumnRef
// items, a column reference by its column alone.
func listStrings(list *ast.List) []string {
	if list == nil {
		return nil
	}
	result := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		switch n := item.(type) {
		case *ast.String:
			result = append(result, n.Str)
		case *ast.ColumnRef:
			result = append(result, n.Column)
		}
	}
	return result
}

func tableConstraints(list *ast.List) []*ast.TableConstraint {
	if list == nil {
		return nil
	}
	result := make([]*ast.TableConstraint, 0, len(list.Items))
	for _, item := range list.Items {
		if c, ok := item.(*ast.TableConstraint); ok {
			result = append(result, c)
		}
	}
	return result
}
