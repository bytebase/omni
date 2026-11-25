package catalog

import (
	"strings"

	"github.com/bytebase/omni/trino/parser"
)

// alterTable runs ALTER TABLE, with the checks Trino makes in the order it
// makes them (recorded in testdata/simulate.tsv).
//
// trino: core/trino-main/src/main/java/io/trino/execution/RenameTableTask,
// AddColumnTask, DropColumnTask, RenameColumnTask, SetColumnTypeTask,
// DropNotNullConstraintTask, SetDefaultValueTask, DropDefaultValueTask,
// SetPropertiesTask, SetTableAuthorizationTask, and TableExecute
func (c *Catalog) alterTable(n *parser.AlterTableStmt, sql string) error {
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	if n.Kind == parser.AlterTableSetProperties {
		// The properties are checked before the table is looked up.
		if err := propertyError("INVALID_TABLE_PROPERTY", "table", name.catalog, n.Properties); err != nil {
			return err
		}
	}
	if n.Kind == parser.AlterTableSetAuthorization {
		if _, err := c.database(name.catalog); err != nil {
			return err
		}
	}
	t, v := c.lookup(name)
	if t == nil && v == nil {
		if n.IfExists {
			return nil
		}
		return errorf("TABLE_NOT_FOUND", "Table '%s' does not exist", name)
	}
	switch n.Kind {
	case parser.AlterTableRename:
		if v != nil {
			return errorf("GENERIC_USER_ERROR", "Table '%s' does not exist, but a view with that name exists. Did you mean ALTER VIEW %s RENAME TO ...?", name, name)
		}
		return c.renameTable(t, name, n)
	case parser.AlterTableSetAuthorization:
		if err := principalRole(n.Authorization); err != nil {
			return err
		}
		return errorf("NOT_SUPPORTED", "Catalog does not support permission management: %s", name.catalog)
	case parser.AlterTableSetProperties:
		return nil
	case parser.AlterTableExecute:
		if v != nil {
			return errorf("NOT_SUPPORTED", "ALTER TABLE EXECUTE is not supported for views")
		}
		return errorf("PROCEDURE_NOT_FOUND", "Table procedure not registered: %s", strings.ToUpper(n.Procedure.Value))
	}
	if t == nil {
		// The memory connector fails on a view here with an internal
		// error (GENERIC_INTERNAL_ERROR, a NullPointerException in
		// MemoryMetadata); omni says what other connectors do.
		return errorf("TABLE_NOT_FOUND", "Table '%s' does not exist", name)
	}
	switch n.Kind {
	case parser.AlterTableAddColumn:
		return addColumn(t, n, name.catalog, sql)
	case parser.AlterTableDropColumn:
		col, err := columnOrNested(t, n.Column)
		if col == nil || err != nil {
			if err == nil && !n.ColumnIfExists {
				err = errorf("COLUMN_NOT_FOUND", "Column '%s' does not exist", n.Column.Normalize())
			}
			return err
		}
		return errorf("NOT_SUPPORTED", "This connector does not support dropping columns")
	case parser.AlterTableRenameColumn:
		col, err := columnOrNested(t, n.Column)
		if col == nil || err != nil {
			if err == nil && !n.ColumnIfExists {
				err = errorf("COLUMN_NOT_FOUND", "Column '%s' does not exist", n.Column.Normalize())
			}
			return err
		}
		target := fold(n.RenamedColumn)
		if t.GetColumn(target) != nil {
			return errorf("COLUMN_ALREADY_EXISTS", "Column '%s' already exists", target)
		}
		delete(t.colByName, col.Name)
		col.Name = target
		for i, cc := range t.columns {
			t.colByName[cc.Name] = i
		}
		return nil
	case parser.AlterTableAlterColumn:
		return alterColumn(t, n, sql)
	}
	return nil
}

// principalRole refuses a ROLE principal of SET AUTHORIZATION, as no role
// exists; a principal spelled ROLE with no name is a user named role.
func principalRole(p *parser.Principal) error {
	if p != nil && p.Kind == parser.PrincipalRole && p.Name != nil {
		return errorf("ROLE_NOT_FOUND", "Role '%s' does not exist", fold(p.Name))
	}
	return nil
}

// columnOrNested finds the column a DROP or RENAME COLUMN names. A
// dotted name reaches into a row field, which no column of the memory
// connector's tables here has.
func columnOrNested(t *Table, qn interface{ NormalizedParts() []string }) (*Column, error) {
	parts := qn.NormalizedParts()
	col := t.GetColumn(strings.ToLower(parts[0]))
	if col == nil || len(parts) == 1 {
		return col, nil
	}
	if !strings.HasPrefix(col.Type, "row(") {
		return nil, errorf("NOT_SUPPORTED", "Cannot resolve field '%s' within %s type", strings.ToLower(parts[1]), col.Type)
	}
	return nil, errorf("NOT_SUPPORTED", "This connector does not support dropping fields")
}

func (c *Catalog) renameTable(t *Table, name objName, n *parser.AlterTableStmt) error {
	if len(n.NewName.Parts) > 3 {
		return errorf("SYNTAX_ERROR", "Too many dots in name: %s", n.NewName.Normalize())
	}
	// An unqualified target is in the source's schema.
	target := objName{name.catalog, name.schema, ""}
	parts := n.NewName.Parts
	switch len(parts) {
	case 1:
		target.name = fold(parts[0])
	case 2:
		target.schema, target.name = fold(parts[0]), fold(parts[1])
	case 3:
		target = objName{fold(parts[0]), fold(parts[1]), fold(parts[2])}
	}
	if target.catalog != name.catalog {
		if _, err := c.database(target.catalog); err != nil {
			return err
		}
		return errorf("NOT_SUPPORTED", "Table rename across catalogs is not supported")
	}
	ts, err := c.schemaOf(target)
	if err != nil {
		return err
	}
	if ts.tables[target.name] != nil || ts.views[target.name] != nil {
		return errorf("TABLE_ALREADY_EXISTS", "Target table '%s' already exists", target)
	}
	delete(t.schema.tables, t.Name)
	t.Name, t.schema = target.name, ts
	ts.tables[target.name] = t
	return nil
}

func addColumn(t *Table, n *parser.AlterTableStmt, catalog, sql string) error {
	col, err := columnOf(n.NewColumn, catalog, sql)
	if err != nil {
		return err
	}
	if t.GetColumn(col.Name) != nil {
		if n.ColumnIfNotExists {
			return nil
		}
		return errorf("COLUMN_ALREADY_EXISTS", "Column '%s' already exists", col.Name)
	}
	switch n.Position {
	case parser.ColumnPositionFirst:
		return errorf("NOT_SUPPORTED", "This connector does not support adding columns with FIRST clause")
	case parser.ColumnPositionAfter:
		return errorf("NOT_SUPPORTED", "This connector does not support adding columns with AFTER clause")
	}
	t.addColumn(col)
	return nil
}

func alterColumn(t *Table, n *parser.AlterTableStmt, sql string) error {
	col, err := columnOrNested(t, n.AlterColumn)
	if err != nil {
		return err
	}
	if col == nil {
		return errorf("COLUMN_NOT_FOUND", "Column '%s' does not exist", n.AlterColumn.Normalize())
	}
	switch n.ColumnAction {
	case parser.AlterColumnSetDataType:
		if _, err := typeName(n.NewColumnType); err != nil {
			return errorf("TYPE_NOT_FOUND", "Unknown type '%s' for column '%s'", n.NewColumnType.String(), col.Name)
		}
		return errorf("NOT_SUPPORTED", "This connector does not support setting column types")
	case parser.AlterColumnDropNotNull:
		if col.Nullable {
			return errorf("NOT_SUPPORTED", "Column is already nullable")
		}
		col.Nullable = true
	case parser.AlterColumnSetDefault:
		text, err := defaultText(n.NewDefault, col.Type, sql)
		if err != nil {
			return err
		}
		col.Default = text
	case parser.AlterColumnDropDefault:
		if col.Default == "" {
			return errorf("GENERIC_USER_ERROR", "Column '%s' does not have a default value", col.Name)
		}
		col.Default = ""
	}
	return nil
}
