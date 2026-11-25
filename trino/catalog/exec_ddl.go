package catalog

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// The statements Exec runs act as Trino 482's do on its memory connector,
// as testdata/simulate.tsv records them: what a statement refuses is
// Trino's errorName, and what it stores is what SHOW CREATE TABLE, SHOW
// CREATE VIEW and information_schema read back.
//
// Trino keeps every name in lower case, quoted or not, so a name is folded
// whole here (unlike Normalize, which keeps a quoted name's case).
//
// trino: core/trino-main/src/main/java/io/trino/execution (the
// DataDefinitionTask of each statement) and
// plugin/trino-memory/src/main/java/io/trino/plugin/memory/MemoryMetadata.java

// fold is the name Trino stores for an identifier.
func fold(id *ast.Identifier) string {
	if id == nil {
		return ""
	}
	return strings.ToLower(id.Value)
}

// objName is a fully qualified catalog.schema.name.
type objName struct{ catalog, schema, name string }

func (n objName) String() string { return n.catalog + "." + n.schema + "." + n.name }

// schemaString is the catalog.schema a name is in.
func (n objName) schemaString() string { return n.catalog + "." + n.schema }

// qualify completes a table or view name with the session's catalog and
// schema, as MetadataUtil.createQualifiedObjectName does.
func (c *Catalog) qualify(qn *ast.QualifiedName) (objName, error) {
	var parts []string
	if qn != nil {
		for _, p := range qn.Parts {
			parts = append(parts, fold(p))
		}
	}
	switch len(parts) {
	case 1:
		if c.currentCatalog == "" {
			return objName{}, errorf("MISSING_CATALOG_NAME", "Catalog must be specified when session catalog is not set")
		}
		if c.currentSchema == "" {
			return objName{}, errorf("MISSING_SCHEMA_NAME", "Schema must be specified when session schema is not set")
		}
		return objName{c.currentCatalog, c.currentSchema, parts[0]}, nil
	case 2:
		if c.currentCatalog == "" {
			return objName{}, errorf("MISSING_CATALOG_NAME", "Catalog must be specified when session catalog is not set")
		}
		return objName{c.currentCatalog, parts[0], parts[1]}, nil
	case 3:
		return objName{parts[0], parts[1], parts[2]}, nil
	default:
		return objName{}, errorf("SYNTAX_ERROR", "Too many dots in table name: %s", strings.Join(parts, "."))
	}
}

// qualifySchema completes a schema name with the session's catalog.
func (c *Catalog) qualifySchema(qn *ast.QualifiedName) (catalog, schema string, err error) {
	var parts []string
	if qn != nil {
		for _, p := range qn.Parts {
			parts = append(parts, fold(p))
		}
	}
	switch len(parts) {
	case 1:
		if c.currentCatalog == "" {
			return "", "", errorf("MISSING_CATALOG_NAME", "Catalog must be specified when session catalog is not set")
		}
		return c.currentCatalog, parts[0], nil
	case 2:
		return parts[0], parts[1], nil
	default:
		return "", "", errorf("SYNTAX_ERROR", "Too many parts in schema name: %s", strings.Join(parts, "."))
	}
}

// database returns the catalog a name is in.
func (c *Catalog) database(name string) (*Database, error) {
	db := c.databases[name]
	if db == nil {
		return nil, errorf("CATALOG_NOT_FOUND", "Catalog '%s' not found", name)
	}
	return db, nil
}

// schemaOf returns the schema a name is in.
func (c *Catalog) schemaOf(n objName) (*Schema, error) {
	db, err := c.database(n.catalog)
	if err != nil {
		return nil, err
	}
	s := db.schemas[n.schema]
	if s == nil {
		return nil, errorf("SCHEMA_NOT_FOUND", "Schema '%s' does not exist", n.schema)
	}
	return s, nil
}

// lookup returns the table or view a name holds, nil for none: a missing
// catalog or schema holds nothing, as for Trino's getTableHandle.
func (c *Catalog) lookup(n objName) (*Table, *View) {
	db := c.databases[n.catalog]
	if db == nil {
		return nil, nil
	}
	s := db.schemas[n.schema]
	if s == nil {
		return nil, nil
	}
	return s.tables[n.name], s.views[n.name]
}

// exec runs one statement.
func (c *Catalog) exec(stmt ast.Node, sql string) error {
	switch n := stmt.(type) {
	case *parser.CreateSchemaStmt:
		return c.createSchema(n)
	case *parser.DropStmt:
		return c.drop(n)
	case *parser.AlterSchemaStmt:
		return c.alterSchema(n)
	case *parser.CreateTableStmt:
		return c.createTable(n, sql)
	case *parser.CreateTableAsStmt:
		return c.createTableAs(n)
	case *parser.AlterTableStmt:
		return c.alterTable(n, sql)
	case *parser.CreateViewStmt:
		return c.createView(n, sql)
	case *parser.AlterViewStmt:
		return c.alterView(n)
	case *parser.CreateMaterializedViewStmt:
		return c.createMaterializedView(n)
	case *parser.AlterMaterializedViewStmt:
		return c.alterMaterializedView(n)
	case *parser.CommentStmt:
		return c.comment(n)
	case *parser.TruncateStmt:
		return c.truncate(n)
	}
	// Anything else leaves the schema as it is.
	return nil
}

// propertyError refuses the first property of a WITH list: the memory
// connector declares no schema, table, view or column property.
func propertyError(code, what, catalog string, props []*parser.Property) error {
	if len(props) == 0 {
		return nil
	}
	return errorf(code, "Catalog '%s' %s property '%s' does not exist", catalog, what, fold(props[0].Name))
}

func (c *Catalog) createSchema(n *parser.CreateSchemaStmt) error {
	catalog, schema, err := c.qualifySchema(n.Name)
	if err != nil {
		return err
	}
	db, err := c.database(catalog)
	if err != nil {
		return err
	}
	// The properties are checked before the schema is looked up, a role
	// after.
	if err := propertyError("INVALID_SCHEMA_PROPERTY", "schema", catalog, n.Properties); err != nil {
		return err
	}
	if db.schemas[schema] != nil {
		if n.IfNotExists {
			return nil
		}
		return errorf("SCHEMA_ALREADY_EXISTS", "Schema '%s.%s' already exists", catalog, schema)
	}
	if err := principalError(n.Authorization); err != nil {
		return err
	}
	db.EnsureSchema(schema)
	return nil
}

// principalError refuses a ROLE principal of CREATE SCHEMA: the memory
// connector supports no role.
func principalError(p *parser.Principal) error {
	if p != nil && p.Kind == parser.PrincipalRole && p.Name != nil {
		return errorf("NOT_SUPPORTED", "This connector does not support roles")
	}
	return nil
}

func (c *Catalog) alterSchema(n *parser.AlterSchemaStmt) error {
	catalog, schema, err := c.qualifySchema(n.Name)
	if err != nil {
		return err
	}
	// A missing catalog holds no schema.
	db := c.databases[catalog]
	var s *Schema
	if db != nil {
		s = db.schemas[schema]
	}
	if s == nil {
		return errorf("SCHEMA_NOT_FOUND", "Schema '%s.%s' does not exist", catalog, schema)
	}
	if n.Kind == parser.AlterSchemaSetAuthorization {
		return errorf("NOT_SUPPORTED", "Catalog does not support permission management: %s", catalog)
	}
	target := fold(n.NewName)
	if db.schemas[target] != nil {
		return errorf("SCHEMA_ALREADY_EXISTS", "Target schema '%s.%s' already exists", catalog, target)
	}
	delete(db.schemas, schema)
	s.Name = target
	db.schemas[target] = s
	return nil
}

func (c *Catalog) drop(n *parser.DropStmt) error {
	if n.Object == parser.DropSchema {
		catalog, schema, err := c.qualifySchema(n.Name)
		if err != nil {
			return err
		}
		db := c.databases[catalog]
		var s *Schema
		if db != nil {
			s = db.schemas[schema]
		}
		if s == nil {
			if n.IfExists {
				return nil
			}
			return errorf("SCHEMA_NOT_FOUND", "Schema '%s.%s' does not exist", catalog, schema)
		}
		if n.Behavior != parser.DropBehaviorCascade && len(s.tables)+len(s.views) > 0 {
			return errorf("SCHEMA_NOT_EMPTY", "Cannot drop non-empty schema '%s'", schema)
		}
		delete(db.schemas, schema)
		return nil
	}
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	t, v := c.lookup(name)
	// A relation of another kind is refused even with IF EXISTS; a missing
	// view or materialized view is a user error, a missing table not found.
	switch n.Object {
	case parser.DropTable:
		switch {
		case t != nil:
			delete(t.schema.tables, t.Name)
		case v != nil:
			return errorf("GENERIC_USER_ERROR", "Table '%s' does not exist, but a view with that name exists. Did you mean DROP VIEW %s?", name, name)
		case !n.IfExists:
			return errorf("TABLE_NOT_FOUND", "Table '%s' does not exist", name)
		}
	case parser.DropView:
		switch {
		case v != nil:
			delete(v.schema.views, v.Name)
		case t != nil:
			return errorf("GENERIC_USER_ERROR", "View '%s' does not exist, but a table with that name exists. Did you mean DROP TABLE %s?", name, name)
		case !n.IfExists:
			return errorf("GENERIC_USER_ERROR", "View '%s' does not exist", name)
		}
	case parser.DropMaterializedView:
		switch {
		case t != nil || v != nil:
			return errorf("GENERIC_USER_ERROR", "Materialized view '%s' does not exist, but a table or view with that name exists", name)
		case !n.IfExists:
			return errorf("GENERIC_USER_ERROR", "Materialized view '%s' does not exist", name)
		}
	}
	return nil
}

// newRelationTarget checks that a table or view may be created at name:
// its catalog and schema exist and no relation holds the name. It returns
// the schema, or nil when IF NOT EXISTS met an existing table.
func (c *Catalog) newRelationTarget(name objName, ifNotExists bool, what string) (*Schema, error) {
	s, err := c.schemaOf(name)
	if err != nil {
		return nil, err
	}
	if s.tables[name.name] != nil || s.views[name.name] != nil {
		if ifNotExists {
			return nil, nil
		}
		return nil, errorf("TABLE_ALREADY_EXISTS", "%s '%s' already exists", what, name)
	}
	return s, nil
}

func (c *Catalog) createTable(n *parser.CreateTableStmt, sql string) error {
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	if n.OrReplace {
		if _, err := c.schemaOf(name); err != nil {
			return err
		}
		return errorf("NOT_SUPPORTED", "This connector does not support replacing tables")
	}
	s, err := c.newRelationTarget(name, n.IfNotExists, "Table")
	if s == nil || err != nil {
		return err
	}
	var cols []*Column
	seen := map[string]bool{}
	for _, el := range n.Elements {
		if el.Like != nil {
			src, err := c.qualify(el.Like.Source)
			if err != nil {
				return err
			}
			t, _ := c.lookup(src)
			if t == nil {
				return errorf("TABLE_NOT_FOUND", "LIKE table '%s' does not exist", src)
			}
			for _, col := range t.columns {
				if seen[col.Name] {
					return errorf("DUPLICATE_COLUMN_NAME", "Column name '%s' specified more than once", col.Name)
				}
				seen[col.Name] = true
				cp := *col
				cols = append(cols, &cp)
			}
			continue
		}
		col, err := columnOf(el.Column, name.catalog, sql)
		if err != nil {
			return err
		}
		if seen[col.Name] {
			return errorf("DUPLICATE_COLUMN_NAME", "Column name '%s' specified more than once", col.Name)
		}
		seen[col.Name] = true
		cols = append(cols, col)
	}
	if err := propertyError("INVALID_TABLE_PROPERTY", "table", name.catalog, n.Properties); err != nil {
		return err
	}
	t := s.AddTable(name.name, cols...)
	t.Comment = n.Comment
	return nil
}

// columnOf reads a column definition: its type, default and properties,
// checked in that order.
func columnOf(def *parser.ColumnDefinition, catalog, sql string) (*Column, error) {
	name := fold(def.Name)
	typ, err := typeName(def.Type)
	if err != nil {
		return nil, typeError(err, def.Type, def.Name.Value)
	}
	col := &Column{Name: name, Type: typ, Nullable: !def.NotNull, Comment: def.Comment}
	if def.Default != nil {
		if col.Default, err = defaultText(def.Default, typ, sql); err != nil {
			return nil, err
		}
	}
	if err := propertyError("INVALID_COLUMN_PROPERTY", "column", catalog, def.Properties); err != nil {
		return nil, err
	}
	return col, nil
}

// typeError is the refusal of a column type: an unknown name, or Trino's
// own error for a known type's bad parameters.
func typeError(err error, dt *parser.DataType, column string) error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return errorf("TYPE_NOT_FOUND", "Unknown type '%s' for column '%s'", dt.String(), column)
}

// exprText is an expression's source text.
func exprText(e parser.Expr, sql string) string {
	span := e.Span()
	if span.Start < 0 || span.End > len(sql) || span.Start > span.End {
		return ""
	}
	return sql[span.Start:span.End]
}

func (c *Catalog) createTableAs(n *parser.CreateTableAsStmt) error {
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	if n.OrReplace {
		if _, err := c.schemaOf(name); err != nil {
			return err
		}
	}
	s, err := c.newRelationTarget(name, n.IfNotExists || n.OrReplace, "Table")
	if err != nil {
		return err
	}
	if s == nil && !n.OrReplace {
		return nil
	}
	if err := propertyError("INVALID_TABLE_PROPERTY", "table", name.catalog, n.Properties); err != nil {
		return err
	}
	cols, err := c.queryColumns(n.Query)
	if err != nil {
		return err
	}
	if n.ColumnAliases != nil {
		if len(n.ColumnAliases) != len(cols) {
			return errorf("MISMATCHED_COLUMN_ALIASES", "Column alias list has %d entries but relation has %d columns", len(n.ColumnAliases), len(cols))
		}
		for i, a := range n.ColumnAliases {
			cols[i].Name = fold(a)
		}
	}
	if err := namedColumns(cols); err != nil {
		return err
	}
	if n.OrReplace {
		return errorf("NOT_SUPPORTED", "This connector does not support replacing tables")
	}
	t := s.AddTable(name.name, cols...)
	t.Comment = n.Comment
	return nil
}

// namedColumns refuses a query output column with no name, or a name twice.
func namedColumns(cols []*Column) error {
	seen := map[string]bool{}
	for i, col := range cols {
		if col.Name == "" {
			return errorf("MISSING_COLUMN_NAME", "Column name not specified at position %d", i+1)
		}
		if seen[col.Name] {
			return errorf("DUPLICATE_COLUMN_NAME", "Column name '%s' specified more than once", col.Name)
		}
		seen[col.Name] = true
	}
	return nil
}

// createView runs CREATE VIEW: a table of the name, or a view without OR
// REPLACE, refuses it before its query is analyzed; a missing schema, after.
func (c *Catalog) createView(n *parser.CreateViewStmt, sql string) error {
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	t, v := c.lookup(name)
	if t != nil {
		return errorf("TABLE_ALREADY_EXISTS", "Table already exists: '%s'", name)
	}
	if v != nil && !n.OrReplace {
		return errorf("TABLE_ALREADY_EXISTS", "View already exists: '%s'", name)
	}
	cols, err := c.queryColumns(n.Query)
	if err != nil {
		return err
	}
	if err := namedColumns(cols); err != nil {
		return err
	}
	if err := propertyError("INVALID_VIEW_PROPERTY", "view", name.catalog, n.Properties); err != nil {
		return err
	}
	db, err := c.database(name.catalog)
	if err != nil {
		return err
	}
	s := db.schemas[name.schema]
	if s == nil {
		return errorf("NOT_FOUND", "Schema %s not found", name.schema)
	}
	nv := s.AddView(name.name, cols...)
	// A view's definition is read in the view's catalog and schema, so it
	// is kept only when those are the session's it was written in.
	q := n.Query.Span()
	if name.catalog == c.currentCatalog && name.schema == c.currentSchema && q.Start >= 0 && q.Start <= q.End && q.End <= len(sql) {
		nv.Definition = sql[q.Start:q.End]
	}
	nv.Comment = n.Comment
	nv.Invoker = n.Security == parser.ViewSecurityInvoker
	return nil
}

// createMaterializedView runs CREATE MATERIALIZED VIEW, which the memory
// connector refuses once the statement is analyzed: its GRACE PERIOD, its
// query, and its properties.
func (c *Catalog) createMaterializedView(n *parser.CreateMaterializedViewStmt) error {
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	if iv, ok := n.GracePeriod.(*parser.IntervalLiteral); ok {
		if _, err := strconv.ParseFloat(strings.TrimSpace(iv.Value), 64); err != nil {
			return errorf("INVALID_LITERAL", "Invalid INTERVAL %s value: %s", iv.From, iv.Value)
		}
	}
	cols, err := c.queryColumns(n.Query)
	if err != nil {
		return err
	}
	if err := namedColumns(cols); err != nil {
		return err
	}
	if err := propertyError("INVALID_MATERIALIZED_VIEW_PROPERTY", "materialized view", name.catalog, n.Properties); err != nil {
		return err
	}
	if _, err := c.database(name.catalog); err != nil {
		return err
	}
	if n.GracePeriod != nil {
		return errorf("NOT_SUPPORTED", "Catalog '%s' does not support GRACE PERIOD", name.catalog)
	}
	return errorf("NOT_SUPPORTED", "This connector does not support creating materialized views")
}

// alterMaterializedView runs ALTER MATERIALIZED VIEW. The memory
// connector holds none, so each form refuses: SET PROPERTIES on its
// properties, the others for want of the materialized view, even with IF
// EXISTS when a table or view holds the name.
func (c *Catalog) alterMaterializedView(n *parser.AlterMaterializedViewStmt) error {
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	if n.Kind == parser.AlterMVSetProperties {
		if err := propertyError("INVALID_MATERIALIZED_VIEW_PROPERTY", "materialized view", name.catalog, n.Properties); err != nil {
			return err
		}
	}
	t, v := c.lookup(name)
	if n.Kind == parser.AlterMVRename && n.IfExists && t == nil && v == nil {
		return nil
	}
	return errorf("TABLE_NOT_FOUND", "Materialized view '%s' does not exist", name)
}

func (c *Catalog) alterView(n *parser.AlterViewStmt) error {
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	if n.Kind == parser.AlterViewSetAuthorization {
		if _, err := c.database(name.catalog); err != nil {
			return err
		}
	}
	t, v := c.lookup(name)
	if v == nil {
		if t != nil {
			return errorf("TABLE_NOT_FOUND", "View '%s' does not exist, but a table with that name exists. Did you mean ALTER TABLE %s RENAME TO ...?", name, name)
		}
		return errorf("TABLE_NOT_FOUND", "View '%s' does not exist", name)
	}
	switch n.Kind {
	case parser.AlterViewSetAuthorization:
		if err := principalRole(n.Authorization); err != nil {
			return err
		}
		return errorf("NOT_SUPPORTED", "Catalog does not support permission management: %s", name.catalog)
	case parser.AlterViewRefresh:
		return nil
	}
	if len(n.NewName.Parts) > 3 {
		return errorf("SYNTAX_ERROR", "Too many dots in name: %s", n.NewName.Normalize())
	}
	target, err := c.qualify(n.NewName)
	if err != nil {
		return err
	}
	if target.catalog != name.catalog {
		if _, err := c.database(target.catalog); err != nil {
			return err
		}
		return errorf("NOT_SUPPORTED", "View rename across catalogs is not supported")
	}
	ts, err := c.schemaOf(target)
	if err != nil {
		return err
	}
	if ts.tables[target.name] != nil || ts.views[target.name] != nil {
		return errorf("GENERIC_USER_ERROR", "Target view '%s' already exists", target)
	}
	if ts != v.schema {
		// Its definition was written for its old schema, where it would
		// no longer be read.
		v.Definition = ""
	}
	delete(v.schema.views, v.Name)
	v.Name, v.schema = target.name, ts
	ts.views[target.name] = v
	return nil
}

func (c *Catalog) comment(n *parser.CommentStmt) error {
	var text *string
	if !n.IsNull {
		s := n.Comment
		text = &s
	}
	if n.Object == parser.CommentOnColumn {
		if len(n.Name.Parts) < 2 {
			return errorf("MISSING_TABLE", "Table must be specified")
		}
		rel := &ast.QualifiedName{Parts: n.Name.Parts[:len(n.Name.Parts)-1]}
		if len(rel.Parts) > 3 {
			// a.b.c.d names column d of table a.b.c.
			rel.Parts = rel.Parts[len(rel.Parts)-3:]
		}
		name, err := c.qualify(rel)
		if err != nil {
			return err
		}
		t, v := c.lookup(name)
		colName := fold(n.Name.Parts[len(n.Name.Parts)-1])
		var cols []*Column
		switch {
		case t != nil:
			cols = t.columns
		case v != nil:
			cols = v.columns
		default:
			return errorf("TABLE_NOT_FOUND", "Table does not exist: %s", name)
		}
		for _, col := range cols {
			if col.Name == colName {
				col.Comment = text
				return nil
			}
		}
		return errorf("COLUMN_NOT_FOUND", "Column does not exist: %s", colName)
	}
	name, err := c.qualify(n.Name)
	if err != nil {
		return err
	}
	t, v := c.lookup(name)
	if n.Object == parser.CommentOnView {
		if v == nil {
			if t != nil {
				return errorf("TABLE_NOT_FOUND", "View '%s' does not exist, but a table with that name exists. Did you mean COMMENT ON TABLE %s IS ...?", name, name)
			}
			return errorf("TABLE_NOT_FOUND", "View '%s' does not exist", name)
		}
		v.Comment = text
		return nil
	}
	if t == nil {
		if v != nil {
			return errorf("TABLE_NOT_FOUND", "Table '%s' does not exist, but a view with that name exists. Did you mean COMMENT ON VIEW %s IS ...?", name, name)
		}
		return errorf("TABLE_NOT_FOUND", "Table does not exist: %s", name)
	}
	t.Comment = text
	return nil
}

func (c *Catalog) truncate(n *parser.TruncateStmt) error {
	name, err := c.qualify(n.Target)
	if err != nil {
		return err
	}
	t, v := c.lookup(name)
	switch {
	case v != nil:
		return errorf("NOT_SUPPORTED", "Cannot truncate a view")
	case t == nil:
		return errorf("TABLE_NOT_FOUND", "Table '%s' does not exist", name)
	}
	return nil
}
