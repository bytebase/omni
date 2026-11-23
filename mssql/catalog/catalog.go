// Package catalog is SQL Server's DDL catalog: it loads a schema from SQL or
// from the metadata the driver syncs, renders the definition back as T-SQL,
// diffs two snapshots, and generates the migration that takes one to the
// other.
//
// The rendering and diff judgment are Bytebase's schema/mssql ported nearly
// verbatim, so the DDL a dump emits and the migrations a release generates stay
// byte-for-byte what they were before the move.
package catalog

import (
	"context"
	"encoding/json"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

const (
	defaultSchema = "dbo"
	noAction      = "NO ACTION"
)

// Catalog is one SQL Server database: its schemas and the objects in them.
// Unlike pg and mysql, SQL Server's objects are already the shape Bytebase
// syncs, so the catalog holds metadata.DatabaseSchemaMetadata and applies DDL to
// it.
type Catalog struct {
	database                          string
	defaultSchema                     string
	metadataDefaultSchema             string
	schemas                           []*metadata.SchemaMetadata
	schemaByName                      map[string]*metadata.SchemaMetadata
	schemasByFold                     map[string][]*metadata.SchemaMetadata
	tableByName                       map[tableKey]*metadata.TableMetadata
	indexCounter                      int
	databaseOrder                     []string
	databaseByName                    map[string]*Catalog
	skipViewIndexes                   bool
	omitUnnamedDefaultConstraintNames bool
	lenientMetadataLoad               bool
	// loadingScript is set while LoadSQL runs a desired-state script, whose
	// objects need not come in the order SQL Server would create them, so a
	// view's query is not bound against what precedes it.
	loadingScript bool
	// viewColumns holds the columns a view was bound to when Exec created or
	// altered it, which SQL Server keeps until the view is altered or
	// refreshed.
	viewColumns map[*metadata.ViewMetadata][]qcol
	// fromMetadata is set once a synced snapshot is loaded, which records
	// neither full-text indexes, HIDDEN columns, nor system versioning, so
	// the checks that read them are inconclusive.
	fromMetadata bool
	// fulltextTables and hiddenColumns record what Exec ran: the tables
	// with a full-text index, and the HIDDEN columns.
	fulltextTables        map[tableKey]bool
	hiddenColumns         map[*metadata.ColumnMetadata]bool
	systemTimeTables      map[tableKey][2]string
	systemVersionedTables map[tableKey]bool
}

type tableKey struct {
	schema string
	table  string
}

// New returns an empty catalog.
func New() *Catalog {
	return &Catalog{
		defaultSchema:         defaultSchema,
		schemaByName:          make(map[string]*metadata.SchemaMetadata),
		schemasByFold:         make(map[string][]*metadata.SchemaMetadata),
		tableByName:           make(map[tableKey]*metadata.TableMetadata),
		databaseByName:        make(map[string]*Catalog),
		systemTimeTables:      make(map[tableKey][2]string),
		systemVersionedTables: make(map[tableKey]bool),
	}
}

// DatabaseName returns the database the catalog describes.
func (c *Catalog) DatabaseName() string { return c.database }

// DefaultSchema returns the schema an unqualified name resolves to.
func (c *Catalog) DefaultSchema() string { return c.defaultSchema }

// SetDefaultSchema sets the schema an unqualified name resolves to.
func (c *Catalog) SetDefaultSchema(schemaName string) {
	if schemaName != "" {
		c.defaultSchema = schemaName
	}
}

// Databases returns the loaded or registered database names in load order.
func (c *Catalog) Databases() []string {
	if len(c.databaseOrder) > 0 {
		out := make([]string, len(c.databaseOrder))
		copy(out, c.databaseOrder)
		return out
	}
	if c.database != "" {
		return []string{c.database}
	}
	return nil
}

// DatabaseCatalog returns the catalog for the named database, or the receiver
// when database is empty or matches the primary database.
func (c *Catalog) DatabaseCatalog(database string) *Catalog {
	if c == nil {
		return nil
	}
	if database == "" {
		return c
	}
	if c.databaseByName != nil {
		if sub, ok := c.databaseByName[normalize(database)]; ok {
			return sub
		}
	}
	if strings.EqualFold(c.database, database) {
		return c
	}
	return nil
}

// Schemas returns the schemas in the named database (or the primary database
// when database is empty).
func (c *Catalog) Schemas(database string) []*metadata.SchemaMetadata {
	target := c.DatabaseCatalog(database)
	if target == nil {
		return nil
	}
	return target.schemas
}

// ListDatabaseSchemaNames returns the schema names in the named database in
// load order.
func (c *Catalog) ListDatabaseSchemaNames(database string) []string {
	target := c.DatabaseCatalog(database)
	if target == nil {
		return nil
	}
	return target.ListSchemaNames()
}

// GetDatabaseSchema returns the named schema in the named database, or nil.
func (c *Catalog) GetDatabaseSchema(database, schemaName string) *metadata.SchemaMetadata {
	target := c.DatabaseCatalog(database)
	if target == nil {
		return nil
	}
	if schemaName == "" {
		schemaName = target.defaultSchema
	}
	return target.GetSchema(schemaName)
}

// GetDatabaseTable returns the named table in the named database and schema,
// or nil.
func (c *Catalog) GetDatabaseTable(database, schemaName, tableName string) *metadata.TableMetadata {
	target := c.DatabaseCatalog(database)
	if target == nil {
		return nil
	}
	if schemaName == "" {
		schemaName = target.defaultSchema
	}
	return target.GetTable(schemaName, tableName)
}

// GetView returns the named view in the named schema of the primary database,
// or nil.
func (c *Catalog) GetView(schemaName, viewName string) *metadata.ViewMetadata {
	return c.GetDatabaseView("", schemaName, viewName)
}

// GetDatabaseView returns the named view in the named database and schema, or
// nil.
func (c *Catalog) GetDatabaseView(database, schemaName, viewName string) *metadata.ViewMetadata {
	schema := c.GetDatabaseSchema(database, schemaName)
	if schema == nil {
		return nil
	}
	for _, view := range schema.Views {
		if strings.EqualFold(view.Name, viewName) {
			return view
		}
	}
	return nil
}

// GetSchema returns the named schema, or nil.
func (c *Catalog) GetSchema(name string) *metadata.SchemaMetadata {
	if name == "" {
		name = c.defaultSchema
	}
	if schema := c.schemaByName[name]; schema != nil {
		return schema
	}
	matches := c.schemasByFold[normalize(name)]
	if len(matches) == 1 {
		return matches[0]
	}
	return nil
}

// ListSchemaNames returns the schema names in load order.
func (c *Catalog) ListSchemaNames() []string {
	names := make([]string, 0, len(c.schemas))
	for _, schema := range c.schemas {
		names = append(names, schema.Name)
	}
	return names
}

// GetTable returns the named table in the named schema, or nil.
func (c *Catalog) GetTable(schemaName, tableName string) *metadata.TableMetadata {
	schema := c.GetSchema(schemaName)
	if schema == nil {
		return nil
	}
	if table := c.tableByName[tableKey{schema: schema.Name, table: tableName}]; table != nil {
		return table
	}
	var match *metadata.TableMetadata
	for _, table := range schema.Tables {
		if strings.EqualFold(table.Name, tableName) {
			if match != nil {
				return nil
			}
			match = table
		}
	}
	return match
}

// Metadata returns a copy of the catalog as the metadata snapshot Bytebase
// syncs and stores.
func (c *Catalog) Metadata() *metadata.DatabaseSchemaMetadata {
	meta := &metadata.DatabaseSchemaMetadata{
		Name:          c.database,
		DefaultSchema: c.metadataDefaultSchema,
		Schemas:       make([]*metadata.SchemaMetadata, 0, len(c.schemas)),
	}
	for _, schema := range c.schemas {
		meta.Schemas = append(meta.Schemas, cloneSchema(schema))
	}
	return meta
}

func (c *Catalog) registerDatabaseName(database string) {
	if database == "" {
		return
	}
	for _, existing := range c.databaseOrder {
		if strings.EqualFold(existing, database) {
			return
		}
	}
	c.databaseOrder = append(c.databaseOrder, database)
}

func (c *Catalog) loadPrimaryMetadata(meta *metadata.DatabaseSchemaMetadata, database string) {
	c.fromMetadata = true
	c.database = database
	c.defaultSchema = defaultSchema
	c.metadataDefaultSchema = meta.GetDefaultSchema()
	if meta.GetDefaultSchema() != "" {
		c.defaultSchema = meta.GetDefaultSchema()
	}
	c.schemas = nil
	c.schemaByName = make(map[string]*metadata.SchemaMetadata)
	c.schemasByFold = make(map[string][]*metadata.SchemaMetadata)
	c.tableByName = make(map[tableKey]*metadata.TableMetadata)
	c.systemTimeTables = make(map[tableKey][2]string)
	c.systemVersionedTables = make(map[tableKey]bool)
	for _, schema := range meta.GetSchemas() {
		c.addSchema(cloneSchema(schema))
	}
}

// LoadMetadata replaces the catalog with the snapshot. Nothing can fail to
// install: the catalog stores the snapshot's objects as they are, and a later
// Exec or Definition reads them through the same helpers.
func (c *Catalog) LoadMetadata(ctx context.Context, meta *metadata.DatabaseSchemaMetadata) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.databaseOrder = nil
	c.databaseByName = make(map[string]*Catalog)
	c.loadPrimaryMetadata(meta, meta.GetName())
	if c.database != "" {
		c.registerDatabaseName(c.database)
		c.databaseByName[normalize(c.database)] = c
	}
	return nil
}

// LoadMetadataInto loads meta under the given database name. When the receiver
// has no primary database yet (or database matches the primary database), it
// populates the receiver as the primary database; otherwise it stores the
// snapshot as an additional database for cross-database resolution while
// keeping the primary database intact. A nil meta registers the database name
// without schemas.
func (c *Catalog) LoadMetadataInto(ctx context.Context, meta *metadata.DatabaseSchemaMetadata, database string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if c.databaseByName == nil {
		c.databaseByName = make(map[string]*Catalog)
	}
	if meta == nil {
		if database != "" {
			if c.database == "" {
				c.database = database
			}
			c.registerDatabaseName(database)
			if _, ok := c.databaseByName[normalize(database)]; !ok {
				sub := New()
				sub.database = database
				c.databaseByName[normalize(database)] = sub
			}
		}
		return nil
	}
	if c.database == "" || strings.EqualFold(c.database, database) {
		c.loadPrimaryMetadata(meta, database)
		if database != "" {
			c.registerDatabaseName(database)
			c.databaseByName[normalize(database)] = c
		}
		return nil
	}
	sub := New()
	sub.loadPrimaryMetadata(meta, database)
	if database != "" {
		c.registerDatabaseName(database)
		c.databaseByName[normalize(database)] = sub
	}
	return nil
}

// ExecOptions selects how Exec treats a failing statement.
type ExecOptions struct {
	// ContinueOnError continues after a statement fails, collecting every
	// error. When false execution stops at the first failure.
	ContinueOnError bool
}

// ExecResult is the outcome of one statement.
type ExecResult struct {
	Index   int
	SQL     string
	Line    int
	Skipped bool
	Error   error
}

// Exec parses sql and applies the CREATE statements Bytebase's metadata
// extraction models to the catalog. Batches separated by GO are handled by the
// parser; GO itself is a no-op. ALTER, DROP and other statements are skipped,
// exactly as Bytebase's extraction skips them, so loading a desired-state
// script yields the same snapshot on both sides.
func (c *Catalog) Exec(sql string, opts *ExecOptions) ([]ExecResult, error) {
	stmts, err := mssql.Parse(sql)
	if err != nil {
		return nil, err
	}
	continueOnError := false
	if opts != nil {
		continueOnError = opts.ContinueOnError
	}
	results := make([]ExecResult, 0, len(stmts))
	for i := range stmts {
		stmt := &stmts[i]
		result := ExecResult{Index: i, SQL: stmt.Text, Line: stmt.Start.Line}
		if stmt.Empty() {
			result.Skipped = true
			results = append(results, result)
			continue
		}
		if _, ok := stmt.AST.(*ast.GoStmt); ok {
			result.Skipped = true
			results = append(results, result)
			continue
		}
		if err := c.execStatement(stmt.AST, sql); err != nil {
			result.Error = err
			results = append(results, result)
			if !continueOnError {
				return results, nil
			}
			continue
		}
		results = append(results, result)
	}
	return results, nil
}

func loadSQLWithOptions(sql string, skipViewIndexes, omitUnnamedDefaultConstraintNames bool) (*Catalog, error) {
	c := New()
	c.skipViewIndexes = skipViewIndexes
	c.omitUnnamedDefaultConstraintNames = omitUnnamedDefaultConstraintNames
	c.lenientMetadataLoad = skipViewIndexes
	c.loadingScript = true
	results, err := c.Exec(sql, nil)
	c.loadingScript = false
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		c.getOrCreateSchema(defaultSchema)
	}
	for _, result := range results {
		if result.Error != nil {
			return c, result.Error
		}
	}
	c.sortSchemas()
	return c, nil
}

// LoadSQL parses sql into a new catalog.
func LoadSQL(sql string) (*Catalog, error) {
	return loadSQLWithOptions(sql, false, false)
}

// GetDatabaseMetadata parses SQL Server DDL and returns the schema it defines.
func GetDatabaseMetadata(schemaText string) (*metadata.DatabaseSchemaMetadata, error) {
	// Skip indexed views and omit synthetic DF_<table>_<col> default names
	// (only needed by forward LoadSQL/migration synthesis) so round-tripping
	// extracted schemas neither fails on view indexes nor collides across
	// underscore-split table/column names.
	c, err := loadSQLWithOptions(schemaText, true, true)
	if err != nil {
		return nil, err
	}
	return c.Metadata(), nil
}

// sortSchemas orders schemas by name, as Bytebase's metadata extraction does.
func (c *Catalog) sortSchemas() {
	slices.SortStableFunc(c.schemas, func(a, b *metadata.SchemaMetadata) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func (c *Catalog) execStatement(node ast.Node, source string) error {
	switch n := node.(type) {
	case *ast.CreateSchemaStmt:
		return c.execCreateSchema(n, source)
	case *ast.AlterSchemaStmt:
		return c.execAlterSchema(n)
	case *ast.CreateTableStmt:
		return c.execCreateTable(n, source)
	case *ast.CreateTableAsSelectStmt:
		return c.execCreateTableAsSelect(n)
	case *ast.CreateIndexStmt:
		return c.execCreateIndex(n, source)
	case *ast.CreateSpatialIndexStmt:
		return c.execCreateSpatialIndex(n)
	case *ast.CreateXmlIndexStmt:
		return c.execCreateXMLIndex(n)
	case *ast.CreateSelectiveXmlIndexStmt:
		return c.execCreateSelectiveXMLIndex(n)
	case *ast.AlterIndexStmt:
		return c.execAlterIndex(n)
	case *ast.CreateViewStmt:
		return c.execCreateOrAlterView(n, source)
	case *ast.CreateFulltextIndexStmt:
		c.execCreateFulltextIndex(n)
	case *ast.CreateProcedureStmt:
		return c.execCreateOrAlterProcedure(n, source)
	case *ast.CreateFunctionStmt:
		return c.execCreateOrAlterFunction(n, source)
	case *ast.CreateSequenceStmt:
		return c.execCreateSequence(n, source)
	case *ast.AlterSequenceStmt:
		return c.execAlterSequence(n, source)
	case *ast.CreateTriggerStmt:
		return c.execCreateOrAlterTrigger(n, source)
	case *ast.CreateSynonymStmt:
		return c.execCreateSynonym(n)
	case *ast.AlterTableStmt:
		return c.execAlterTable(n, source)
	case *ast.DropStmt:
		return c.execDrop(n)
	case *ast.TruncateStmt:
		return c.execTruncateTable(n)
	case *ast.ExecStmt:
		c.execExtendedProperty(n, source)
	}
	return nil
}

// extractColumn extracts a column definition and its inline constraints into
// table.

func (c *Catalog) extractColumn(n *ast.ColumnDef, table *metadata.TableMetadata, schemaName, source string) {
	column := &metadata.ColumnMetadata{
		Name:     n.Name,
		Position: int32(len(table.Columns) + 1),
		Nullable: true,
	}
	// HIDDEN is kept by the column itself, so it goes when the column or
	// its table does, and follows a rename.
	if n.Hidden {
		if c.hiddenColumns == nil {
			c.hiddenColumns = map[*metadata.ColumnMetadata]bool{}
		}
		c.hiddenColumns[column] = true
	}
	if n.Computed != nil {
		expression := strings.TrimSpace(nodeText(n.Computed.Expr, source))
		if !strings.HasPrefix(expression, "(") {
			expression = "(" + expression + ")"
		}
		column.Type = "AS " + expression
		if n.Computed.Persisted {
			column.Type += " PERSISTED"
		}
		if n.Computed.NotNull {
			column.Nullable = false
		} else if n.Computed.Persisted && isNonNullLiteralExpr(n.Computed.Expr) {
			column.Nullable = false
		}
	} else if n.DataType != nil {
		column.Type = nodeText(n.DataType, source)
	}
	if n.Identity != nil {
		column.IsIdentity = true
		column.IdentitySeed = n.Identity.Seed
		column.IdentityIncrement = n.Identity.Increment
		column.Nullable = false
	}
	if n.Nullable != nil {
		column.Nullable = !n.Nullable.NotNull
	}
	if n.DefaultExpr != nil {
		column.Default = nodeText(n.DefaultExpr, source)
	}
	if n.Collation != "" {
		column.Collation = n.Collation
	}

	if n.Constraints != nil {
		for _, item := range n.Constraints.Items {
			constraint, ok := item.(*ast.ConstraintDef)
			if !ok {
				continue
			}
			switch constraint.Type {
			case ast.ConstraintPrimaryKey, ast.ConstraintUnique:
				primary := constraint.Type == ast.ConstraintPrimaryKey
				cols := []string{column.Name}
				if constraint.Columns != nil {
					// A key written beside a column may name its own columns,
					// as `ALTER TABLE a ADD x a UNIQUE(x)` does: the list keys
					// x, not the column it is written beside.
					cols = stringNodeList(constraint.Columns)
					if primary {
						for _, name := range cols {
							if col := columnByName(table, name); col != nil {
								col.Nullable = false
							}
						}
					}
				}
				if primary && (constraint.Columns == nil || slices.Contains(cols, column.Name)) {
					column.Nullable = false
				}
				prefix := "UQ"
				if primary {
					prefix = "PK"
				}
				table.Indexes = append(table.Indexes, c.constraintIndexMetadata(
					constraint,
					constraintPrefix(prefix, table.Name),
					cols,
					primary,
				))
			case ast.ConstraintDefault:
				if constraint.Expr != nil {
					column.Default = nodeText(constraint.Expr, source)
					column.DefaultConstraintName = constraint.Name
				}
			case ast.ConstraintCheck:
				c.appendCheckConstraint(table, constraint, source)
			case ast.ConstraintForeignKey:
				c.appendForeignKey(table, constraint, []string{column.Name}, c.defaultSchema)
			}
		}
	}

	if column.Default != "" && column.DefaultConstraintName == "" && !c.omitUnnamedDefaultConstraintNames {
		column.DefaultConstraintName = "DF_" + table.Name + "_" + column.Name
	}
	table.Columns = append(table.Columns, column)
}

func constraintPrefix(kind, table string) string {
	return kind + "_" + table
}

func constraintName(prefix string, counter int) string {
	return prefix + "_" + strconv.Itoa(counter)
}

func (c *Catalog) extractTableConstraint(n *ast.ConstraintDef, table *metadata.TableMetadata, schemaName, source string) {
	switch n.Type {
	case ast.ConstraintPrimaryKey:
		index := c.constraintIndexMetadata(n, constraintPrefix("PK", table.Name), nil, true)
		for _, name := range index.Expressions {
			if column := columnByName(table, name); column != nil {
				column.Nullable = false
			}
		}
		table.Indexes = append(table.Indexes, index)
	case ast.ConstraintUnique:
		table.Indexes = append(table.Indexes, c.constraintIndexMetadata(n, constraintPrefix("UQ", table.Name), nil, false))
	case ast.ConstraintCheck:
		c.appendCheckConstraint(table, n, source)
	case ast.ConstraintForeignKey:
		c.appendForeignKey(table, n, stringNodeList(n.Columns), c.defaultSchema)
	}
}

func (c *Catalog) execExtendedProperty(n *ast.ExecStmt, source string) {
	if n.Name == nil || !strings.EqualFold(n.Name.Object, "sp_addextendedproperty") {
		return
	}
	if n.Args == nil {
		return
	}
	named := make(map[string]string)
	var positional []string
	for _, item := range n.Args.Items {
		arg, ok := item.(*ast.ExecArg)
		if !ok || arg.Value == nil {
			continue
		}
		var value string
		switch v := arg.Value.(type) {
		case *ast.Literal:
			value = v.Str
		default:
			value = nodeText(v, source)
		}
		if arg.Name != "" {
			named[strings.ToLower(strings.TrimPrefix(arg.Name, "@"))] = value
			continue
		}
		positional = append(positional, value)
	}
	var name, comment, schemaName, tableName, columnName, objectType string
	if len(named) > 0 {
		name = named["name"]
		comment = named["value"]
		schemaName = named["schema"]
		if schemaName == "" {
			schemaName = named["level0name"]
		}
		tableName = named["table"]
		if tableName == "" {
			tableName = named["level1name"]
		}
		objectType = named["table"]
		if objectType == "" {
			objectType = named["level1type"]
		}
		columnName = named["column"]
		if columnName == "" {
			columnName = named["level2name"]
		}
	} else {
		if len(positional) < 6 {
			return
		}
		name = positional[0]
		comment = positional[1]
		for i := 2; i < len(positional)-1; i += 2 {
			switch strings.ToUpper(positional[i]) {
			case "SCHEMA":
				schemaName = positional[i+1]
			case "TABLE", "VIEW", "FUNCTION", "PROCEDURE", "SEQUENCE":
				objectType = strings.ToUpper(positional[i])
				tableName = positional[i+1]
			case "COLUMN":
				columnName = positional[i+1]
			}
		}
	}
	if !strings.EqualFold(name, "MS_Description") || schemaName == "" || tableName == "" {
		return
	}
	c.applyCommentForType(schemaName, objectType, tableName, columnName, comment)
}

func (c *Catalog) applyCommentForType(schemaName, objectType, objectName, columnName, comment string) {
	if objectType == "" {
		objectType = "TABLE"
	}
	switch strings.ToUpper(objectType) {
	case "TABLE":
		c.applyComment(schemaName, objectName, columnName, comment)
	case "VIEW":
		if schema := c.GetSchema(schemaName); schema != nil {
			for _, view := range schema.Views {
				if strings.EqualFold(view.Name, objectName) {
					view.Comment = comment
					return
				}
			}
		}
	case "FUNCTION":
		if schema := c.GetSchema(schemaName); schema != nil {
			for _, function := range schema.Functions {
				if strings.EqualFold(function.Name, objectName) {
					function.Comment = comment
					return
				}
			}
		}
	case "PROCEDURE":
		if schema := c.GetSchema(schemaName); schema != nil {
			for _, procedure := range schema.Procedures {
				if strings.EqualFold(procedure.Name, objectName) {
					procedure.Comment = comment
					return
				}
			}
		}
	case "SEQUENCE":
		if schema := c.GetSchema(schemaName); schema != nil {
			for _, sequence := range schema.Sequences {
				if strings.EqualFold(sequence.Name, objectName) {
					sequence.Comment = comment
					return
				}
			}
		}
	}
}

func (c *Catalog) applyComment(schemaName, tableName, columnName, comment string) {
	table := c.resolveOrCreateTable(schemaName, tableName)
	if columnName != "" {
		for _, column := range table.Columns {
			if strings.EqualFold(column.Name, columnName) {
				column.Comment = comment
				return
			}
		}
		return
	}
	table.Comment = comment
}

func (c *Catalog) constraintIndexMetadata(n *ast.ConstraintDef, prefix string, columns []string, primary bool) *metadata.IndexMetadata {
	name := n.Name
	if name == "" && prefix != "" {
		c.indexCounter++
		name = constraintName(prefix, c.indexCounter)
	}
	index := &metadata.IndexMetadata{
		Name:         name,
		Primary:      primary,
		Unique:       true,
		IsConstraint: true,
		Expressions:  []string{},
		Descending:   []bool{},
		Type:         indexType(false, n.Clustered),
	}
	for _, column := range columns {
		index.Expressions = append(index.Expressions, column)
		index.Descending = append(index.Descending, false)
	}
	if len(columns) == 0 && n.Columns != nil {
		appendIndexColumns(index, n.Columns)
	}
	return index
}

func inlineIndexMetadata(n *ast.InlineIndexDef, source string) *metadata.IndexMetadata {
	index := &metadata.IndexMetadata{
		Name:         n.Name,
		Unique:       n.Unique,
		Type:         indexType(n.Columnstore, n.Clustered),
		Expressions:  []string{},
		Descending:   []bool{},
		IsConstraint: false,
	}
	appendIndexColumns(index, n.Columns)
	index.Definition = indexDefinition(n.IncludeCols, n.WhereClause, source)
	return index
}

func (c *Catalog) appendCheckConstraint(table *metadata.TableMetadata, n *ast.ConstraintDef, source string) {
	name := n.Name
	if name == "" {
		c.indexCounter++
		name = constraintName(constraintPrefix("CK", table.Name), c.indexCounter)
	}
	table.CheckConstraints = append(table.CheckConstraints, &metadata.CheckConstraintMetadata{
		Name:       name,
		Expression: nodeText(n.Expr, source),
	})
}

func (c *Catalog) appendForeignKey(table *metadata.TableMetadata, n *ast.ConstraintDef, columns []string, fallbackSchema string) {
	name := n.Name
	if name == "" {
		c.indexCounter++
		name = constraintName(constraintPrefix("FK", table.Name), c.indexCounter)
	}
	fk := &metadata.ForeignKeyMetadata{
		Name:              name,
		Columns:           columns,
		ReferencedColumns: stringNodeList(n.RefColumns),
		OnDelete:          referentialAction(n.OnDelete),
		OnUpdate:          referentialAction(n.OnUpdate),
	}
	if fk.OnDelete == "" {
		fk.OnDelete = noAction
	}
	if fk.OnUpdate == "" {
		fk.OnUpdate = noAction
	}
	if n.RefTable != nil {
		refSchema, refTable := tableRefSchemaObject(n.RefTable, fallbackSchema)
		fk.ReferencedSchema = refSchema
		fk.ReferencedTable = refTable
		if len(fk.ReferencedColumns) == 0 {
			if parent := c.GetTable(refSchema, refTable); parent != nil {
				for _, index := range parent.Indexes {
					if index.Primary {
						fk.ReferencedColumns = append([]string(nil), index.Expressions...)
						break
					}
				}
			}
		}
	}
	table.ForeignKeys = append(table.ForeignKeys, fk)
}

func (c *Catalog) getOrCreateSchema(name string) *metadata.SchemaMetadata {
	if name == "" {
		name = c.defaultSchema
	}
	if schema := c.schemaByName[name]; schema != nil {
		return schema
	}
	schema := &metadata.SchemaMetadata{
		Name:       name,
		Tables:     []*metadata.TableMetadata{},
		Views:      nil,
		Procedures: nil,
		Functions:  nil,
		Sequences:  nil,
	}
	c.addSchema(schema)
	return schema
}

func (c *Catalog) addSchema(schema *metadata.SchemaMetadata) {
	c.schemas = append(c.schemas, schema)
	c.schemaByName[schema.Name] = schema
	fold := normalize(schema.Name)
	c.schemasByFold[fold] = append(c.schemasByFold[fold], schema)
	for _, table := range schema.Tables {
		c.indexTable(schema.Name, table)
	}
}

func (c *Catalog) indexTable(schemaName string, table *metadata.TableMetadata) {
	c.tableByName[tableKey{schema: schemaName, table: table.Name}] = table
}

func (c *Catalog) resolveOrCreateTable(schemaName, tableName string) *metadata.TableMetadata {
	if table := c.GetTable(schemaName, tableName); table != nil {
		return table
	}
	return c.getOrCreateTable(schemaName, tableName)
}

func (c *Catalog) getOrCreateTable(schemaName, tableName string) *metadata.TableMetadata {
	schemaName = normalizeOrDefault(schemaName, c.defaultSchema)
	key := tableKey{schema: schemaName, table: tableName}
	if table := c.tableByName[key]; table != nil {
		return table
	}
	table := &metadata.TableMetadata{
		Name:             tableName,
		Columns:          []*metadata.ColumnMetadata{},
		Indexes:          []*metadata.IndexMetadata{},
		ForeignKeys:      nil,
		CheckConstraints: nil,
	}
	schema := c.getOrCreateSchema(schemaName)
	schema.Tables = append(schema.Tables, table)
	c.indexTable(schema.Name, table)
	return table
}

// quote renders an identifier as a bracketed name, doubling a closing bracket
// inside it so the name means what the metadata says.
func quote(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}

func normalize(name string) string {
	return strings.ToLower(name)
}

func normalizeOrDefault(name, fallback string) string {
	if name == "" {
		return fallback
	}
	return name
}

func tableRefSchemaObject(ref *ast.TableRef, fallbackSchema string) (string, string) {
	if ref == nil {
		return fallbackSchema, ""
	}
	schemaName := ref.Schema
	if schemaName == "" {
		schemaName = fallbackSchema
	}
	return schemaName, ref.Object
}

func nodeText(node ast.Node, source string) string {
	loc := ast.NodeLoc(node)
	if loc.Start < 0 || loc.End < 0 || loc.Start > loc.End || loc.End > len(source) {
		return ""
	}
	return source[loc.Start:loc.End]
}

func definitionText(node ast.Node, source string) string {
	loc := ast.NodeLoc(node)
	if loc.Start < 0 || loc.End < 0 || loc.Start > loc.End || loc.End > len(source) {
		return ""
	}
	end := loc.End
	for end < len(source) {
		switch source[end] {
		case ' ', '\t', '\r', '\n':
			end++
		default:
			if source[end] == ';' {
				end++
			}
			return strings.TrimSpace(source[loc.Start:end])
		}
	}
	return strings.TrimSpace(source[loc.Start:end])
}

func indexType(columnstore bool, clustered *bool) string {
	if columnstore {
		if clustered != nil && *clustered {
			return "CLUSTERED COLUMNSTORE"
		}
		return "NONCLUSTERED COLUMNSTORE"
	}
	if clustered == nil {
		return ""
	}
	if *clustered {
		return "CLUSTERED"
	}
	return "NONCLUSTERED"
}

func appendIndexColumns(index *metadata.IndexMetadata, columns *ast.List) {
	if columns == nil {
		return
	}
	for _, item := range columns.Items {
		column, ok := item.(*ast.IndexColumn)
		if !ok {
			continue
		}
		index.Expressions = append(index.Expressions, column.Name)
		index.Descending = append(index.Descending, column.SortDir == ast.SortDesc)
	}
}

func indexDefinition(includeCols *ast.List, whereClause ast.ExprNode, source string) string {
	var parts []string
	var include []string
	if includeCols != nil {
		for _, item := range includeCols.Items {
			switch column := item.(type) {
			case *ast.IndexColumn:
				include = append(include, column.Name)
			case *ast.String:
				include = append(include, column.Str)
			}
		}
	}
	if len(include) > 0 {
		parts = append(parts, "INCLUDE:"+strings.Join(include, "\x00"))
	}
	if whereClause != nil {
		parts = append(parts, "WHERE "+strings.TrimSpace(nodeText(whereClause, source)))
	}
	return strings.Join(parts, " ")
}

func indexIncludeColumns(definition string) []string {
	if !strings.Contains(definition, "INCLUDE:") {
		return nil
	}
	after := strings.SplitN(definition, "INCLUDE:", 2)[1]
	if index := strings.Index(after, " WHERE "); index != -1 {
		after = after[:index]
	}
	var columns []string
	for _, column := range strings.Split(after, "\x00") {
		if column = strings.TrimSpace(column); column != "" {
			columns = append(columns, column)
		}
	}
	return columns
}

func indexColumnNames(columns *ast.List) []string {
	var names []string
	if columns == nil {
		return names
	}
	for _, item := range columns.Items {
		column, ok := item.(*ast.IndexColumn)
		if !ok {
			continue
		}
		names = append(names, column.Name)
	}
	return names
}

func stringNodeList(list *ast.List) []string {
	var result []string
	if list == nil {
		return result
	}
	for _, item := range list.Items {
		switch n := item.(type) {
		case *ast.String:
			result = append(result, n.Str)
		case *ast.IndexColumn:
			result = append(result, n.Name)
		}
	}
	return result
}

func referentialAction(action ast.ReferentialAction) string {
	switch action {
	case ast.RefActCascade:
		return "CASCADE"
	case ast.RefActSetNull:
		return "SET NULL"
	case ast.RefActSetDefault:
		return "SET DEFAULT"
	case ast.RefActNoAction:
		return noAction
	default:
		return ""
	}
}

// cloneSchema copies a snapshot through its JSON form. The metadata types are
// protobuf messages, whose mutex rules prevent a struct copy, and the catalog
// must not take a direct protobuf dependency.
func cloneSchema(schema *metadata.SchemaMetadata) *metadata.SchemaMetadata {
	data, err := json.Marshal(schema)
	if err != nil {
		return schema
	}
	out := &metadata.SchemaMetadata{}
	if err := json.Unmarshal(data, out); err != nil {
		return schema
	}
	if out.Tables == nil {
		out.Tables = []*metadata.TableMetadata{}
	}
	return out
}

func cloneTable(table *metadata.TableMetadata) *metadata.TableMetadata {
	data, err := json.Marshal(table)
	if err != nil {
		return table
	}
	out := &metadata.TableMetadata{}
	if err := json.Unmarshal(data, out); err != nil {
		return table
	}
	return out
}
