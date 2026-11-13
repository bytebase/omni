package v1review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// SchemaBackwardCompatibility judges which top-level statements may break
// the existing data or code. It takes every statement of the change in one
// call, since an object the change itself created earlier may be dropped or
// altered freely, and sets Statement on each finding. A statement nested in a
// routine body or an IF or BEGIN block is not judged. currentDatabase is
// the database the change runs against, used as written.
//
// The change's tables are remembered by database.schema.table as written,
// case kept, the database defaulting to currentDatabase and the schema to
// dbo, and an empty part (no currentDatabase) left out; its schemas by
// currentDatabase.schema with the schema lower-cased (the AUTHORIZATION
// owner when unnamed); its databases by name lower-cased. A drop forgets
// the object, so dropping it again is reported.
//
// Reported, in statement and then list or action order:
//   - IncompatibleDropTable, IncompatibleDropSchema,
//     IncompatibleDropDatabase: each object a DROP TABLE, DROP SCHEMA or
//     DROP DATABASE lists that the change did not create, Name holding its
//     remembered key; anchored at the statement.
//   - For an ALTER TABLE of a table the change did not create, by action:
//     IncompatibleDropColumn for a DROP COLUMN, Name holding its columns
//     lower-cased and joined by ", "; IncompatibleAlterColumn for every
//     ALTER COLUMN, Name holding the column lower-cased, possibly empty;
//     IncompatibleAddForeignKeyNoCheck or IncompatibleAddCheckNoCheck for
//     a foreign key or check added WITH NOCHECK or a NOCHECK CONSTRAINT
//     action carrying one; these anchor the action. A constraint added
//     with neither WITH CHECK nor WITH NOCHECK is IncompatibleAddPrimaryKey,
//     IncompatibleAddUniqueKey, IncompatibleAddCheck or
//     IncompatibleAddForeignKey, anchored at the statement; one added WITH
//     CHECK is not reported.
//   - IncompatibleRenameTable: an EXEC of sp_rename, any letter case, with
//     no schema qualifier, whose first argument is a string literal;
//     anchored at the statement.
//
// Ported from Bytebase's V1 rule SCHEMA_BACKWARD_COMPATIBILITY.
func SchemaBackwardCompatibility(stmts []ast.Node, currentDatabase string) []v1review.Finding {
	c := &compatibility{
		currentDatabase: currentDatabase,
		newTables:       map[string]bool{},
		newSchemas:      map[string]bool{},
		newDatabases:    map[string]bool{},
	}
	for i, stmt := range stmts {
		c.stmt = i
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Name != nil {
				c.newTables[normalizeTableRef(n.Name, currentDatabase, "dbo")] = true
			}
		case *ast.CreateSchemaStmt:
			schema := strings.ToLower(n.Name)
			if schema == "" {
				schema = strings.ToLower(n.Authorization)
			}
			c.newSchemas[currentDatabase+"."+schema] = true
		case *ast.CreateDatabaseStmt:
			c.newDatabases[strings.ToLower(n.Name)] = true
		case *ast.DropStmt:
			c.drop(n)
		case *ast.AlterTableStmt:
			c.alterTable(n)
		case *ast.ExecStmt:
			if isSpRename(n) {
				c.add(v1review.IncompatibleRenameTable, "", n.Loc)
			}
		}
	}
	return c.findings
}

type compatibility struct {
	currentDatabase string
	newTables       map[string]bool
	newSchemas      map[string]bool
	newDatabases    map[string]bool
	stmt            int
	findings        []v1review.Finding
}

func (c *compatibility) add(code v1review.Code, name string, loc ast.Loc) {
	start, end := span(loc)
	c.findings = append(c.findings, v1review.Finding{Code: code, Name: name, Start: start, End: end, Statement: c.stmt})
}

func (c *compatibility) drop(n *ast.DropStmt) {
	var code v1review.Code
	var created map[string]bool
	var key func(*ast.TableRef) string
	switch n.ObjectType {
	case ast.DropTable:
		code, created = v1review.IncompatibleDropTable, c.newTables
		key = func(ref *ast.TableRef) string { return normalizeTableRef(ref, c.currentDatabase, "dbo") }
	case ast.DropSchema:
		code, created = v1review.IncompatibleDropSchema, c.newSchemas
		key = func(ref *ast.TableRef) string { return c.currentDatabase + "." + strings.ToLower(ref.Object) }
	case ast.DropDatabase:
		code, created = v1review.IncompatibleDropDatabase, c.newDatabases
		key = func(ref *ast.TableRef) string { return strings.ToLower(ref.Object) }
	default:
		return
	}
	if n.Names == nil {
		return
	}
	for _, item := range n.Names.Items {
		ref, ok := item.(*ast.TableRef)
		if !ok {
			continue
		}
		name := key(ref)
		if !created[name] {
			c.add(code, name, n.Loc)
		}
		delete(created, name)
	}
}

func (c *compatibility) alterTable(n *ast.AlterTableStmt) {
	if n.Name == nil || c.newTables[normalizeTableRef(n.Name, c.currentDatabase, "dbo")] || n.Actions == nil {
		return
	}
	for _, item := range n.Actions.Items {
		action, ok := item.(*ast.AlterTableAction)
		if !ok {
			continue
		}
		switch action.Type {
		case ast.ATDropColumn:
			var columns []string
			if action.ColName != "" {
				columns = append(columns, strings.ToLower(action.ColName))
			}
			if action.Names != nil {
				for _, nameItem := range action.Names.Items {
					if ref, ok := nameItem.(*ast.TableRef); ok {
						columns = append(columns, strings.ToLower(ref.Object))
					}
				}
			}
			if len(columns) > 0 {
				c.add(v1review.IncompatibleDropColumn, strings.Join(columns, ", "), action.Loc)
			}
		case ast.ATAlterColumn:
			column := strings.ToLower(action.ColName)
			if action.Column != nil {
				column = strings.ToLower(action.Column.Name)
			}
			c.add(v1review.IncompatibleAlterColumn, column, action.Loc)
		case ast.ATAddConstraint:
			if action.Constraint == nil {
				continue
			}
			if strings.EqualFold(action.WithCheck, "NOCHECK") {
				c.addNoCheck(action)
				continue
			}
			if strings.EqualFold(action.WithCheck, "CHECK") {
				continue
			}
			var code v1review.Code
			switch action.Constraint.Type {
			case ast.ConstraintPrimaryKey:
				code = v1review.IncompatibleAddPrimaryKey
			case ast.ConstraintUnique:
				code = v1review.IncompatibleAddUniqueKey
			case ast.ConstraintCheck:
				code = v1review.IncompatibleAddCheck
			case ast.ConstraintForeignKey:
				code = v1review.IncompatibleAddForeignKey
			default:
				continue
			}
			c.add(code, "", n.Loc)
		case ast.ATNocheckConstraint:
			if action.Constraint != nil {
				c.addNoCheck(action)
			}
		}
	}
}

// addNoCheck reports a foreign key or check constraint added without
// checking the existing rows.
func (c *compatibility) addNoCheck(action *ast.AlterTableAction) {
	switch action.Constraint.Type {
	case ast.ConstraintForeignKey:
		c.add(v1review.IncompatibleAddForeignKeyNoCheck, "", action.Loc)
	case ast.ConstraintCheck:
		c.add(v1review.IncompatibleAddCheckNoCheck, "", action.Loc)
	}
}

// isSpRename reports whether n executes sp_rename, unqualified by a schema,
// with a string literal first argument.
func isSpRename(n *ast.ExecStmt) bool {
	if n.Name == nil || n.Name.Schema != "" || strings.ToLower(n.Name.Object) != "sp_rename" || n.Args.Len() == 0 {
		return false
	}
	arg, ok := n.Args.Items[0].(*ast.ExecArg)
	if !ok {
		return false
	}
	lit, ok := arg.Value.(*ast.Literal)
	return ok && lit.Type == ast.LitString
}

// normalizeTableRef joins the database, schema and object of ref with dots,
// case kept, the database defaulting to fallbackDB and the schema to
// fallbackSchema, and an empty part left out.
func normalizeTableRef(ref *ast.TableRef, fallbackDB, fallbackSchema string) string {
	if ref == nil {
		return ""
	}
	db, schema := ref.Database, ref.Schema
	if db == "" {
		db = fallbackDB
	}
	if schema == "" {
		schema = fallbackSchema
	}
	var parts []string
	for _, part := range []string{db, schema, ref.Object} {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}
