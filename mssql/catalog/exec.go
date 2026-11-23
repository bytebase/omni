package catalog

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

// objectKind names a schema-scoped object kind as SQL Server identifies it in
// namespace and DROP/ALTER validation.
type objectKind string

const (
	kindNone       objectKind = ""
	kindTable      objectKind = "table"
	kindView       objectKind = "view"
	kindProcedure  objectKind = "procedure"
	kindFunction   objectKind = "function"
	kindSequence   objectKind = "sequence"
	kindSynonym    objectKind = "synonym"
	kindTrigger    objectKind = "trigger"
	kindConstraint objectKind = "constraint"
)

// resolveTableInSchema resolves name within schema using the same exact-first,
// reject-ambiguous-fold semantics as GetTable.
func (c *Catalog) resolveTableInSchema(schema *metadata.SchemaMetadata, name string) *metadata.TableMetadata {
	if schema == nil || name == "" {
		return nil
	}
	if table := c.tableByName[tableKey{schema: schema.Name, table: name}]; table != nil {
		return table
	}
	var match *metadata.TableMetadata
	for _, table := range schema.Tables {
		if strings.EqualFold(table.Name, name) {
			if match != nil {
				return nil
			}
			match = table
		}
	}
	return match
}

// lookupSchemaObject returns the kind of the object named name in schema,
// checking tables, views, procedures, functions, sequences, synonyms, triggers,
// and named constraints. When exactTableCase is true, table matching requires
// exact case so case-distinct tables in CS_AS collations can coexist.
func (c *Catalog) lookupSchemaObject(schema *metadata.SchemaMetadata, name string, exactTableCase bool) objectKind {
	if schema == nil || name == "" {
		return kindNone
	}
	if exactTableCase {
		for _, t := range schema.Tables {
			if t.Name == name {
				return kindTable
			}
		}
	} else if c.resolveTableInSchema(schema, name) != nil {
		return kindTable
	}
	for _, v := range schema.Views {
		if strings.EqualFold(v.Name, name) {
			return kindView
		}
	}
	for _, p := range schema.Procedures {
		if strings.EqualFold(p.Name, name) {
			return kindProcedure
		}
	}
	for _, fn := range schema.Functions {
		if strings.EqualFold(fn.Name, name) {
			return kindFunction
		}
	}
	for _, sq := range schema.Sequences {
		if strings.EqualFold(sq.Name, name) {
			return kindSequence
		}
	}
	for _, sy := range schema.Synonyms {
		if strings.EqualFold(sy.Name, name) {
			return kindSynonym
		}
	}
	for _, t := range schema.Tables {
		for _, tr := range t.Triggers {
			if strings.EqualFold(tr.Name, name) {
				return kindTrigger
			}
		}
	}
	for _, v := range schema.Views {
		for _, tr := range v.Triggers {
			if strings.EqualFold(tr.Name, name) {
				return kindTrigger
			}
		}
	}
	return kindNone
}

// isValidDataType reports whether dt names a SQL Server built-in data type.
func isValidDataType(dt *ast.DataType, allowCursor bool) bool {
	if dt == nil {
		return false
	}
	// Schema-qualified custom types (e.g. dbo.MyType) are not tracked as
	// built-in types; treat unqualified names against the built-in registry.
	base := strings.Join(strings.Fields(strings.ToLower(dt.Name)), "")
	base = strings.TrimPrefix(strings.TrimSuffix(base, "]"), "[")
	if base == "cursor" {
		return allowCursor
	}
	return isBuiltinSQLServerType(base)
}

func isValidSequenceType(dt *ast.DataType) bool {
	if dt == nil {
		return true
	}
	base := strings.Join(strings.Fields(strings.ToLower(dt.Name)), "")
	base = strings.TrimPrefix(strings.TrimSuffix(base, "]"), "[")
	switch base {
	case "tinyint", "smallint", "int", "integer", "bigint", "decimal", "dec", "numeric":
		return true
	default:
		return false
	}
}

func isValidFilegroup(name string) bool {
	if name == "" {
		return true
	}
	u := strings.ToUpper(strings.Trim(name, "[]\"'"))
	return u == "PRIMARY" || u == "DEFAULT"
}

func hasInboundFK(c *Catalog, targetSchema, targetTable string) (string, bool) {
	return hasInboundFKExcluding(c, targetSchema, targetTable, nil)
}

func hasInboundFKExcluding(c *Catalog, targetSchema, targetTable string, exclude map[tableKey]bool) (string, bool) {
	for _, s := range c.schemas {
		for _, t := range s.Tables {
			if strings.EqualFold(s.Name, targetSchema) && strings.EqualFold(t.Name, targetTable) {
				continue
			}
			if exclude != nil && exclude[tableKey{schema: normalize(s.Name), table: normalize(t.Name)}] {
				continue
			}
			for _, fk := range t.ForeignKeys {
				refSchema := fk.ReferencedSchema
				if refSchema == "" {
					refSchema = c.defaultSchema
				}
				if strings.EqualFold(refSchema, targetSchema) && strings.EqualFold(fk.ReferencedTable, targetTable) {
					return fk.Name, true
				}
			}
		}
	}
	return "", false
}

func tableHasClusteredIndex(table *metadata.TableMetadata) bool {
	return tableHasClusteredIndexExcluding(table, nil)
}

func tableHasClusteredIndexExcluding(table *metadata.TableMetadata, exclude *metadata.IndexMetadata) bool {
	if table == nil {
		return false
	}
	for _, idx := range table.Indexes {
		if idx == exclude {
			continue
		}
		if strings.EqualFold(idx.Type, "CLUSTERED") || strings.EqualFold(idx.Type, "CLUSTERED COLUMNSTORE") {
			return true
		}
		if idx.Type == "" && idx.IsConstraint && idx.Primary {
			return true
		}
	}
	return false
}

func tableHasIdentityColumn(table *metadata.TableMetadata) bool {
	if table == nil {
		return false
	}
	for _, col := range table.Columns {
		if col.IsIdentity {
			return true
		}
	}
	return false
}

func selectStarWithoutFrom(node ast.Node) bool {
	if node == nil {
		return false
	}
	switch s := node.(type) {
	case *ast.SelectStmt:
		if s.FromClause == nil || len(s.FromClause.Items) == 0 {
			if s.TargetList != nil {
				for _, item := range s.TargetList.Items {
					if _, ok := item.(*ast.StarExpr); ok {
						return true
					}
					if res, ok := item.(*ast.ResTarget); ok {
						if _, ok := res.Val.(*ast.StarExpr); ok {
							return true
						}
					}
				}
			}
		}
	case *ast.ReturnStmt:
		return selectStarWithoutFrom(s.Value)
	case *ast.SubqueryExpr:
		return selectStarWithoutFrom(s.Query)
	case *ast.BeginEndStmt:
		if s.Stmts != nil {
			for _, item := range s.Stmts.Items {
				if selectStarWithoutFrom(item) {
					return true
				}
			}
		}
	}
	return false
}

func isPureAlterStmt(node ast.Node, source string) bool {
	text := strings.TrimSpace(nodeText(node, source))
	return len(text) >= 5 && strings.EqualFold(text[:5], "ALTER")
}

// --- CREATE / ALTER SCHEMA ---

func (c *Catalog) execCreateSchema(n *ast.CreateSchemaStmt, source string) error {
	// Legacy `CREATE SCHEMA AUTHORIZATION owner` with no schema name and no
	// elements is accepted as a no-op by SQL Server 2022.
	if n.Name == "" {
		if n.Elements == nil || len(n.Elements.Items) == 0 {
			return nil
		}
	}
	if !c.lenientMetadataLoad && n.Authorization != "" && !strings.EqualFold(n.Authorization, "dbo") {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the user '%s', because it does not exist or you do not have permission.", n.Authorization)
	}
	if c.GetSchema(n.Name) != nil {
		return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
	}
	schema := c.getOrCreateSchema(n.Name)
	if n.Elements != nil {
		savedIndexCounter := c.indexCounter
		previous := c.defaultSchema
		if n.Name != "" {
			c.defaultSchema = n.Name
		}
		for _, item := range n.Elements.Items {
			if err := c.execStatement(item, source); err != nil {
				c.defaultSchema = previous
				c.indexCounter = savedIndexCounter
				c.removeSchema(schema.Name)
				return err
			}
		}
		c.defaultSchema = previous
	}
	return nil
}

func (c *Catalog) removeSchema(name string) {
	for i, s := range c.schemas {
		if strings.EqualFold(s.Name, name) {
			fold := normalize(s.Name)
			for _, t := range s.Tables {
				delete(c.tableByName, tableKey{schema: s.Name, table: t.Name})
				delete(c.systemTimeTables, tableKey{schema: fold, table: t.Name})
			}
			delete(c.schemaByName, s.Name)
			var kept []*metadata.SchemaMetadata
			for _, fs := range c.schemasByFold[fold] {
				if fs != s {
					kept = append(kept, fs)
				}
			}
			if len(kept) == 0 {
				delete(c.schemasByFold, fold)
			} else {
				c.schemasByFold[fold] = kept
			}
			c.schemas = append(c.schemas[:i], c.schemas[i+1:]...)
			return
		}
	}
}

func (c *Catalog) execAlterSchema(n *ast.AlterSchemaStmt) error {
	targetSchema := c.GetSchema(n.Name)
	if targetSchema == nil {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the schema '%s', because it does not exist or you do not have permission.", n.Name)
	}
	if strings.EqualFold(n.TransferType, "TYPE") || strings.EqualFold(n.TransferType, "XML SCHEMA COLLECTION") {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the %s '%s', because it does not exist or you do not have permission.", strings.ToLower(n.TransferType), n.TransferEntity)
	}
	srcSchemaName := c.defaultSchema
	objName := n.TransferObject
	if n.TransferSchema != "" {
		srcSchemaName = n.TransferSchema
	} else if objName == "" {
		objName = n.TransferEntity
		if dot := strings.LastIndex(objName, "."); dot >= 0 {
			srcSchemaName = objName[:dot]
			objName = objName[dot+1:]
		}
	}
	srcSchema := c.GetSchema(srcSchemaName)
	if srcSchema == nil {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
	}
	kind := c.lookupSchemaObject(srcSchema, objName, false)
	if kind == kindNone || kind == kindTrigger {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
	}
	if !strings.EqualFold(srcSchema.Name, targetSchema.Name) {
		if c.lookupSchemaObject(targetSchema, objName, false) != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", objName)
		}
	}
	if strings.EqualFold(srcSchema.Name, targetSchema.Name) {
		return nil
	}
	switch kind {
	case kindTable:
		table := c.resolveTableInSchema(srcSchema, objName)
		if table == nil {
			return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", n.TransferEntity)
		}
		for i, t := range srcSchema.Tables {
			if t == table {
				srcSchema.Tables = append(srcSchema.Tables[:i], srcSchema.Tables[i+1:]...)
				delete(c.tableByName, tableKey{schema: srcSchema.Name, table: t.Name})
				targetSchema.Tables = append(targetSchema.Tables, t)
				c.indexTable(targetSchema.Name, t)
				// Update foreign keys pointing to this table.
				for _, s := range c.schemas {
					for _, tbl := range s.Tables {
						for _, fk := range tbl.ForeignKeys {
							if strings.EqualFold(fk.ReferencedSchema, srcSchema.Name) && strings.EqualFold(fk.ReferencedTable, t.Name) {
								fk.ReferencedSchema = targetSchema.Name
							}
						}
					}
				}
				return nil
			}
		}
	case kindView:
		for i, v := range srcSchema.Views {
			if strings.EqualFold(v.Name, objName) {
				srcSchema.Views = append(srcSchema.Views[:i], srcSchema.Views[i+1:]...)
				targetSchema.Views = append(targetSchema.Views, v)
				return nil
			}
		}
	case kindProcedure:
		for i, p := range srcSchema.Procedures {
			if strings.EqualFold(p.Name, objName) {
				srcSchema.Procedures = append(srcSchema.Procedures[:i], srcSchema.Procedures[i+1:]...)
				targetSchema.Procedures = append(targetSchema.Procedures, p)
				return nil
			}
		}
	case kindFunction:
		for i, fn := range srcSchema.Functions {
			if strings.EqualFold(fn.Name, objName) {
				srcSchema.Functions = append(srcSchema.Functions[:i], srcSchema.Functions[i+1:]...)
				targetSchema.Functions = append(targetSchema.Functions, fn)
				return nil
			}
		}
	case kindSequence:
		for i, sq := range srcSchema.Sequences {
			if strings.EqualFold(sq.Name, objName) {
				srcSchema.Sequences = append(srcSchema.Sequences[:i], srcSchema.Sequences[i+1:]...)
				targetSchema.Sequences = append(targetSchema.Sequences, sq)
				return nil
			}
		}
	case kindSynonym:
		for i, sy := range srcSchema.Synonyms {
			if strings.EqualFold(sy.Name, objName) {
				srcSchema.Synonyms = append(srcSchema.Synonyms[:i], srcSchema.Synonyms[i+1:]...)
				targetSchema.Synonyms = append(targetSchema.Synonyms, sy)
				return nil
			}
		}
	}
	return nil
}

// --- CREATE TABLE / CTAS ---

func (c *Catalog) execCreateTableAsSelect(n *ast.CreateTableAsSelectStmt) error {
	if n.Columns != nil && len(n.Columns.Items) > 0 {
		first := ""
		if s, ok := n.Columns.Items[0].(*ast.String); ok {
			first = s.Str
		}
		return mssqlErr(ErrInvalidTableDef, "The definition for column '%s' must include a data type.", first)
	}
	return mssqlErr(ErrSyntax, "Incorrect syntax near 'HEAP'.")
}

func (c *Catalog) execCreateTable(n *ast.CreateTableStmt, source string) error {
	// Validate column data types first (compile-time check in SQL Server).
	if !c.lenientMetadataLoad && n.Columns != nil {
		for idx, item := range n.Columns.Items {
			col, ok := item.(*ast.ColumnDef)
			if !ok {
				continue
			}
			if col.Computed == nil && !isValidDataType(col.DataType, false) {
				typeName := ""
				if col.DataType != nil {
					typeName = col.DataType.Name
				}
				return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", idx+1, typeName)
			}
		}
	}

	schemaName, tableName := tableRefSchemaObject(n.Name, c.defaultSchema)
	if schema := c.schemaByName[schemaName]; schema != nil {
		if c.lookupSchemaObject(schema, tableName, true) != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", tableName)
		}
	}
	if !c.lenientMetadataLoad && n.IsFileTable {
		return mssqlErr(ErrFileTableNotSupported, "Default FILESTREAM filegroup is not available.")
	}

	table := c.getOrCreateTable(schemaName, tableName)
	if n.PeriodStartCol != "" && n.PeriodEndCol != "" {
		if c.systemTimeTables == nil {
			c.systemTimeTables = make(map[tableKey]bool)
		}
		c.systemTimeTables[tableKey{schema: normalizeOrDefault(schemaName, c.defaultSchema), table: table.Name}] = true
	}

	if n.Columns != nil {
		for _, item := range n.Columns.Items {
			column, ok := item.(*ast.ColumnDef)
			if !ok {
				continue
			}
			c.extractColumn(column, table, schemaName, source)
		}
	}
	if n.Constraints != nil {
		for _, item := range n.Constraints.Items {
			constraint, ok := item.(*ast.ConstraintDef)
			if !ok {
				continue
			}
			c.extractTableConstraint(constraint, table, schemaName, source)
		}
	}
	if n.Indexes != nil {
		for _, item := range n.Indexes.Items {
			index, ok := item.(*ast.InlineIndexDef)
			if !ok {
				continue
			}
			table.Indexes = append(table.Indexes, inlineIndexMetadata(index, source))
		}
	}
	for _, index := range table.Indexes {
		if !index.Primary {
			continue
		}
		for _, name := range index.Expressions {
			if column := columnByName(table, name); column != nil {
				column.Nullable = false
			}
		}
	}
	return nil
}

// --- CREATE / ALTER INDEX ---

func checkDuplicateIndexColumns(cols, includeCols *ast.List) error {
	seen := make(map[string]bool)
	for _, list := range []*ast.List{cols, includeCols} {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			var name string
			switch v := item.(type) {
			case *ast.IndexColumn:
				name = v.Name
			case *ast.String:
				name = v.Str
			}
			if name == "" {
				continue
			}
			low := strings.ToLower(name)
			if seen[low] {
				return mssqlErr(ErrDupIndexColumn, "Cannot use duplicate column names in index. Column name '%s' listed more than once.", name)
			}
			seen[low] = true
		}
	}
	return nil
}

func (c *Catalog) execCreateIndex(n *ast.CreateIndexStmt, source string) error {
	if err := checkDuplicateIndexColumns(n.Columns, n.IncludeCols); err != nil {
		return err
	}
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)

	// Filtered indexes check object kind differently in SQL Server.
	if n.WhereClause != nil && kind != kindTable {
		if kind == kindView {
			if c.skipViewIndexes {
				return nil
			}
			return mssqlErr(ErrFilteredIndexOnView, "Filtered index '%s' cannot be created on object '%s' because it is not a user table.", n.Name, tableName)
		}
		return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", tableName)
	}

	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindView:
		if c.skipViewIndexes {
			return nil
		}
		if n.Columnstore {
			if n.Clustered != nil && *n.Clustered {
				return mssqlErr(ErrClusteredCSOnView, "The statement failed because a clustered columnstore index cannot be created on a view.")
			}
			return mssqlErr(ErrIndexOnViewNotAllowed, "Cannot create index on view '%s'. It does not have a unique clustered index.", tableName)
		}
		view := c.GetView(schemaName, tableName)
		if view != nil && !isSchemaBound(view.Definition) {
			return mssqlErr(ErrViewNotSchemaBound, "Cannot create index on view '%s' because the view is not schema bound.", tableName)
		}
		return fmt.Errorf("indexes on views are not supported: %s.%s", schemaName, tableName)
	case kindTable:
		// proceed
	default:
		return mssqlErr(ErrIndexInvalidObject, "Index cannot be created on object '%s' because the object is not a user table or view.", tableName)
	}

	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	if n.WhereClause != nil {
		if exprReferencesColumn(n.WhereClause, "IDENTITYCOL") && columnByName(table, "IDENTITYCOL") == nil && !tableHasIdentityColumn(table) {
			return mssqlErr(ErrInvalidColumnName, "Invalid column name 'IDENTITYCOL'.")
		}
	}

	dropExisting := hasOptionOn(n.Options, "DROP_EXISTING")
	var existingIdx *metadata.IndexMetadata
	existingPos := -1
	if n.Name != "" {
		for i, idx := range table.Indexes {
			if strings.EqualFold(idx.Name, n.Name) {
				existingIdx = idx
				existingPos = i
				break
			}
		}
	}
	if dropExisting {
		if existingIdx == nil {
			return mssqlErr(ErrIndexNotFound, "Could not find any index named '%s' for table '%s'.", n.Name, tableName)
		}
	} else if existingIdx != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}

	// Validate key, include, and order columns exist on the table.
	for _, list := range []*ast.List{n.Columns, n.IncludeCols, n.OrderCols} {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			var colName string
			switch v := item.(type) {
			case *ast.IndexColumn:
				colName = v.Name
			case *ast.String:
				colName = v.Str
			}
			if colName != "" && columnByName(table, colName) == nil {
				return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", colName)
			}
		}
	}

	if !c.lenientMetadataLoad {
		if n.Columnstore && n.Clustered != nil && *n.Clustered {
			if tableHasClusteredIndexExcluding(table, existingIdx) {
				return mssqlErr(ErrClusteredCSExists, "You cannot create more than one clustered index on table '%s'.", tableName)
			}
		} else if n.Clustered != nil && *n.Clustered {
			if tableHasClusteredIndexExcluding(table, existingIdx) {
				return mssqlErr(ErrMultipleClusteredIndex, "Cannot create more than one clustered index on table '%s'.", tableName)
			}
		}
	}

	if !c.lenientMetadataLoad && !isValidFilegroup(n.OnFileGroup) {
		return mssqlErr(ErrInvalidFilegroup, "Invalid filegroup '%s' specified.", n.OnFileGroup)
	}

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
	if dropExisting && existingPos >= 0 {
		index.Primary = existingIdx.Primary
		index.IsConstraint = existingIdx.IsConstraint
		table.Indexes[existingPos] = index
	} else {
		table.Indexes = append(table.Indexes, index)
	}
	return nil
}

func (c *Catalog) execCreateXMLIndex(n *ast.CreateXmlIndexStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrXmlIndexNotTable, "Could not create the XML or spatial index on object '%s' because that object is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	col := columnByName(table, n.XmlColumn)
	if col == nil {
		return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", n.XmlColumn)
	}
	if !n.Primary {
		parent := indexByName(table, n.UsingIndex)
		if parent == nil || !strings.EqualFold(parent.Type, "XML") {
			return mssqlErr(ErrXmlIndexMissingPrimary, "Could not find SELECTIVE XML index named '%s' on table '%s'.", n.UsingIndex, tableName)
		}
	}
	if n.Name != "" && indexByName(table, n.Name) != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}
	if normalizeType(col.Type) != "xml" {
		return mssqlErr(ErrXmlIndexNonXmlCol, "Column '%s' on table '%s' is not of type XML.", n.XmlColumn, tableName)
	}

	index := &metadata.IndexMetadata{
		Name:         n.Name,
		Type:         "XML",
		Expressions:  []string{n.XmlColumn},
		Descending:   []bool{false},
		IsConstraint: false,
		Definition:   "PRIMARY",
	}
	if !n.Primary {
		index.Definition = "USING XML INDEX " + n.UsingIndex
		if n.SecondaryFor != "" {
			index.Definition += " FOR " + n.SecondaryFor
		}
	}
	table.Indexes = append(table.Indexes, index)
	return nil
}

func (c *Catalog) execCreateSelectiveXMLIndex(n *ast.CreateSelectiveXmlIndexStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrXmlIndexNotTable, "Could not create the XML or spatial index on object '%s' because that object is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	col := columnByName(table, n.XmlColumn)
	if col == nil {
		return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", n.XmlColumn)
	}
	if n.Name != "" && indexByName(table, n.Name) != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}
	if normalizeType(col.Type) != "xml" {
		return mssqlErr(ErrXmlIndexNonXmlCol, "Column '%s' on table '%s' is not of type XML.", n.XmlColumn, tableName)
	}
	table.Indexes = append(table.Indexes, &metadata.IndexMetadata{
		Name:        n.Name,
		Type:        "XML",
		Expressions: []string{n.XmlColumn},
		Descending:  []bool{false},
		Definition:  "SELECTIVE",
	})
	return nil
}

func (c *Catalog) execCreateSpatialIndex(n *ast.CreateSpatialIndexStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrXmlIndexNotTable, "Could not create the XML or spatial index on object '%s' because that object is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	col := columnByName(table, n.SpatialColumn)
	if col == nil {
		return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", n.SpatialColumn)
	}
	colType := normalizeType(col.Type)
	if colType != "geometry" && colType != "geography" {
		return mssqlErr(ErrSpatialIndexNonSpatialCol, "The requested spatial index on column '%s' of table '%s' could not be created because the column type is not geometry or geography.", n.SpatialColumn, tableName)
	}
	if n.Name != "" && indexByName(table, n.Name) != nil {
		return mssqlErr(ErrDupIndex, "The operation failed because an index or statistics with name '%s' already exists on table '%s'.", n.Name, tableName)
	}

	index := &metadata.IndexMetadata{
		Name:         n.Name,
		Type:         "SPATIAL",
		Expressions:  []string{},
		Descending:   []bool{},
		IsConstraint: false,
		SpatialConfig: &metadata.SpatialIndexConfig{
			Method: "SPATIAL",
			Tessellation: &metadata.TessellationConfig{
				Scheme: strings.ToUpper(n.Using),
			},
			Dimensional: &metadata.DimensionalConfig{
				Dimensions: 2,
			},
		},
	}
	if n.SpatialColumn != "" {
		index.Expressions = append(index.Expressions, n.SpatialColumn)
		index.Descending = append(index.Descending, false)
	}
	if strings.Contains(index.SpatialConfig.Tessellation.Scheme, "GEOGRAPHY") {
		index.SpatialConfig.Dimensional.DataType = "GEOGRAPHY"
	} else {
		index.SpatialConfig.Dimensional.DataType = "GEOMETRY"
	}
	parseSpatialOptions(n.Options, index)
	table.Indexes = append(table.Indexes, index)
	return nil
}

func (c *Catalog) execAlterIndex(n *ast.AlterIndexStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable, kindView:
	default:
		return mssqlErr(ErrIndexInvalidObject, "Index cannot be altered on object '%s' because the object is not a user table or view.", tableName)
	}
	var table *metadata.TableMetadata
	if kind == kindTable {
		table = c.resolveTableInSchema(schema, tableName)
		if table == nil {
			return mssqlErr(ErrObjectNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
		}
	}
	if strings.EqualFold(n.IndexName, "ALL") {
		return nil
	}
	var idx *metadata.IndexMetadata
	if table != nil {
		idx = indexByName(table, n.IndexName)
	}
	if idx == nil {
		return mssqlErr(ErrIndexNotFound, "Cannot find index '%s'.", n.IndexName)
	}
	switch strings.ToUpper(n.Action) {
	case "FOR":
		if !strings.EqualFold(idx.Type, "XML") {
			return mssqlErr(ErrXmlIndexMissingPrimary, "Could not find SELECTIVE XML index named '%s' on table '%s'.", n.IndexName, tableName)
		}
	case "PAUSE", "ABORT":
		return mssqlErr(ErrNoResumableIndexOp, "ALTER INDEX '%s' failed. There is no pending resumable index operation for the index '%s' on '%s'.", strings.ToUpper(n.Action), n.IndexName, tableName)
	}
	return nil
}

// --- CREATE / ALTER VIEW, PROCEDURE, FUNCTION, TRIGGER ---

func (c *Catalog) execCreateOrAlterView(n *ast.CreateViewStmt, source string) error {
	if selectStarWithoutFrom(n.Query) {
		return mssqlErr(ErrSelectStarWithoutFrom, "Must specify table to select from.")
	}
	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, name, false)
	pureAlter := isPureAlterStmt(n, source)

	switch {
	case pureAlter:
		if kind == kindNone {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
		if kind != kindView {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	case n.OrAlter:
		if kind != kindNone && kind != kindView {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	default:
		if kind != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}

	schema = c.getOrCreateSchema(schemaName)
	definition := definitionText(n, source)
	if kind == kindView {
		for _, v := range schema.Views {
			if strings.EqualFold(v.Name, name) {
				v.Definition = creationDefinition(definition)
				return nil
			}
		}
	}
	schema.Views = append(schema.Views, &metadata.ViewMetadata{Name: name, Definition: definition})
	return nil
}

func (c *Catalog) execCreateOrAlterProcedure(n *ast.CreateProcedureStmt, source string) error {
	if !c.lenientMetadataLoad && n.Params != nil {
		for idx, item := range n.Params.Items {
			p, ok := item.(*ast.ParamDef)
			if !ok {
				continue
			}
			if !isValidDataType(p.DataType, true) {
				typeName := ""
				if p.DataType != nil {
					typeName = p.DataType.Name
				}
				return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", idx+1, typeName)
			}
		}
	}
	if !c.lenientMetadataLoad && n.ExternalName != "" {
		return mssqlErr(ErrAssemblyNotFound, "Assembly '%s' was not found in the SQL catalog.", n.ExternalName)
	}
	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, name, false)
	pureAlter := isPureAlterStmt(n, source)

	switch {
	case pureAlter:
		if kind == kindNone {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
		if kind != kindProcedure {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	case n.OrAlter:
		if kind != kindNone && kind != kindProcedure {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	default:
		if kind != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}

	if !c.lenientMetadataLoad {
		if err := checkExecuteAsOptions(n.Options); err != nil {
			return err
		}
	}

	schema = c.getOrCreateSchema(schemaName)
	definition := definitionText(n, source)
	if n.OrAlter {
		definition = normalizeProcDefinition(definition)
	}
	if kind == kindProcedure {
		for _, p := range schema.Procedures {
			if strings.EqualFold(p.Name, name) {
				p.Definition = definition
				return nil
			}
		}
	}
	schema.Procedures = append(schema.Procedures, &metadata.ProcedureMetadata{Name: name, Definition: definition})
	return nil
}

// normalizeProcDefinition mirrors SQL Server's OBJECT_DEFINITION storage:
// CREATE OR ALTER PROC ... is stored as CREATE PROC ... and ALTER PROC ... as
// CREATE PROC ... (or keeps the normalized CREATE form for diffing).
func normalizeProcDefinition(def string) string {
	trimmed := strings.TrimSpace(def)
	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "CREATE OR ALTER ") {
		return "CREATE " + strings.TrimSpace(trimmed[len("CREATE OR ALTER "):])
	}
	if strings.HasPrefix(upper, "ALTER ") {
		return "CREATE " + strings.TrimSpace(trimmed[len("ALTER "):])
	}
	return trimmed
}

func (c *Catalog) execCreateOrAlterFunction(n *ast.CreateFunctionStmt, source string) error {
	if !c.lenientMetadataLoad && n.ExternalName != "" {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the assembly '%s', because it does not exist or you do not have permission.", n.ExternalName)
	}
	if !c.lenientMetadataLoad && n.Params != nil {
		for idx, item := range n.Params.Items {
			p, ok := item.(*ast.ParamDef)
			if !ok {
				continue
			}
			if !isValidDataType(p.DataType, false) {
				typeName := ""
				if p.DataType != nil {
					typeName = p.DataType.Name
				}
				return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", idx+1, typeName)
			}
		}
	}
	if !c.lenientMetadataLoad && n.ReturnType != nil && !isValidDataType(n.ReturnType, false) {
		return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #0: Cannot find data type %s.", n.ReturnType.Name)
	}
	if n.ReturnsTable != nil && selectStarWithoutFrom(n.Body) {
		return mssqlErr(ErrSelectStarWithoutFrom, "Must specify table to select from.")
	}

	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, name, false)
	pureAlter := isPureAlterStmt(n, source)

	switch {
	case pureAlter:
		if kind == kindNone {
			return mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", name)
		}
		if kind != kindFunction {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	case n.OrAlter:
		if kind != kindNone && kind != kindFunction {
			return mssqlErr(ErrIncompatibleAlterType, "Cannot perform alter on '%s' because it is an incompatible object type.", name)
		}
	default:
		if kind != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}

	if !c.lenientMetadataLoad {
		if err := checkExecuteAsOptions(n.Options); err != nil {
			return err
		}
	}

	schema = c.getOrCreateSchema(schemaName)
	definition := definitionText(n, source)
	if kind == kindFunction {
		for _, fn := range schema.Functions {
			if strings.EqualFold(fn.Name, name) {
				fn.Definition = creationDefinition(definition)
				return nil
			}
		}
	}
	schema.Functions = append(schema.Functions, &metadata.FunctionMetadata{Name: name, Definition: definition})
	return nil
}

func (c *Catalog) execCreateOrAlterTrigger(n *ast.CreateTriggerStmt, source string) error {
	if n.Table == nil || n.Name == nil {
		return nil
	}
	schemaName, parentName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	parentKind := c.lookupSchemaObject(schema, parentName, false)
	switch parentKind {
	case kindTable:
	case kindView:
		if !strings.EqualFold(n.TriggerType, "INSTEAD OF") {
			return mssqlErr(ErrTriggerTargetInvalid, "The object '%s' does not exist or is invalid for this operation.", parentName)
		}
	default:
		return mssqlErr(ErrTriggerTargetInvalid, "The object '%s' does not exist or is invalid for this operation.", parentName)
	}

	trigSchemaName, trigName := tableRefSchemaObject(n.Name, schemaName)
	trigSchema := c.GetSchema(trigSchemaName)
	existingKind := c.lookupSchemaObject(trigSchema, trigName, false)
	if !n.OrAlter && existingKind != kindNone {
		return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", trigName)
	}

	var events []string
	if n.Events != nil {
		for _, item := range n.Events.Items {
			switch event := item.(type) {
			case *ast.String:
				events = append(events, event.Str)
			case *ast.TriggerEvent:
				events = append(events, event.Name)
			}
		}
	}
	timing := n.TriggerType
	if strings.EqualFold(timing, "FOR") {
		timing = "AFTER"
	}
	trigger := &metadata.TriggerMetadata{
		Name:   trigName,
		Timing: timing,
		Event:  strings.Join(events, ","),
		Body:   definitionText(n, source),
	}
	if parentKind == kindTable {
		table := c.resolveTableInSchema(schema, parentName)
		if table == nil {
			return mssqlErr(ErrTriggerTargetInvalid, "The object '%s' does not exist or is invalid for this operation.", parentName)
		}
		table.Triggers = append(table.Triggers, trigger)
		return nil
	}
	if view := c.GetView(schemaName, parentName); view != nil {
		view.Triggers = append(view.Triggers, trigger)
	}
	return nil
}

// --- CREATE / ALTER SEQUENCE ---

func (c *Catalog) execCreateSequence(n *ast.CreateSequenceStmt, source string) error {
	if !c.lenientMetadataLoad && n.DataType != nil && !isValidSequenceType(n.DataType) {
		return mssqlErr(ErrUnknownType, "Type '%s' is not a defined system type.", n.DataType.Name)
	}
	schemaName, sequenceName := tableRefSchemaObject(n.Name, c.defaultSchema)
	if schema := c.GetSchema(schemaName); schema != nil {
		if c.lookupSchemaObject(schema, sequenceName, false) != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", sequenceName)
		}
	}
	schema := c.getOrCreateSchema(schemaName)
	sequence := &metadata.SequenceMetadata{Name: sequenceName}
	if n.DataType != nil {
		sequence.DataType = nodeText(n.DataType, source)
	}
	if n.Start != nil {
		sequence.Start = strings.TrimSpace(nodeText(n.Start, source))
	}
	if n.Increment != nil {
		sequence.Increment = strings.TrimSpace(nodeText(n.Increment, source))
	}
	if n.MinValue != nil {
		sequence.MinValue = strings.TrimSpace(nodeText(n.MinValue, source))
	}
	if n.MaxValue != nil {
		sequence.MaxValue = strings.TrimSpace(nodeText(n.MaxValue, source))
	}
	if n.Cycle != nil {
		sequence.Cycle = *n.Cycle
	}
	if n.NoCache {
		sequence.CacheSize = "NO CACHE"
	} else if n.Cache != nil {
		sequence.CacheSize = strings.TrimSpace(nodeText(n.Cache, source))
	}
	schema.Sequences = append(schema.Sequences, sequence)
	return nil
}

func (c *Catalog) execAlterSequence(n *ast.AlterSequenceStmt, source string) error {
	schemaName, sequenceName := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	var seq *metadata.SequenceMetadata
	if schema != nil {
		for _, sq := range schema.Sequences {
			if strings.EqualFold(sq.Name, sequenceName) {
				seq = sq
				break
			}
		}
	}
	if seq == nil {
		return mssqlErr(ErrSecurableNotFound, "Cannot find the object '%s', because it does not exist or you do not have permission.", sequenceName)
	}
	if n.Increment != nil {
		seq.Increment = strings.TrimSpace(nodeText(n.Increment, source))
	}
	if n.NoMinVal {
		seq.MinValue = ""
	} else if n.MinValue != nil {
		seq.MinValue = strings.TrimSpace(nodeText(n.MinValue, source))
	}
	if n.NoMaxVal {
		seq.MaxValue = ""
	} else if n.MaxValue != nil {
		seq.MaxValue = strings.TrimSpace(nodeText(n.MaxValue, source))
	}
	if n.Cycle != nil {
		seq.Cycle = *n.Cycle
	}
	if n.NoCache {
		seq.CacheSize = "NO CACHE"
	} else if n.Cache != nil {
		seq.CacheSize = strings.TrimSpace(nodeText(n.Cache, source))
	} else if n.HasCache {
		seq.CacheSize = ""
	}
	return nil
}

// --- CREATE SYNONYM ---

func (c *Catalog) execCreateSynonym(n *ast.CreateSynonymStmt) error {
	if n.Name == nil || n.Target == nil {
		return nil
	}
	schemaName, name := tableRefSchemaObject(n.Name, c.defaultSchema)
	if schema := c.GetSchema(schemaName); schema != nil {
		if c.lookupSchemaObject(schema, name, false) != kindNone {
			return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", name)
		}
	}
	schema := c.getOrCreateSchema(schemaName)
	schema.Synonyms = append(schema.Synonyms, &metadata.SynonymMetadata{
		Name:           name,
		TargetServer:   n.Target.Server,
		TargetDatabase: n.Target.Database,
		TargetSchema:   n.Target.Schema,
		TargetName:     n.Target.Object,
	})
	return nil
}

// --- ALTER TABLE ---

func (c *Catalog) execAlterTable(n *ast.AlterTableStmt, source string) error {
	if n.Actions == nil {
		return nil
	}
	// Pre-lookup compile-time checks SQL Server performs before resolving the
	// table name.
	if err := checkAlterTablePreLookup(n); err != nil {
		return err
	}

	schemaName, tableName := tableRefSchemaObject(n.Name, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrAlterTableNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrAlterTableNotTable, "Cannot alter '%s' because it is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrAlterTableNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}

	// Stage actions against a clone of the table so a failing later action
	// leaves the catalog untouched, matching SQL Server statement atomicity.
	tableIdx := -1
	for i, t := range schema.Tables {
		if t == table {
			tableIdx = i
			break
		}
	}
	staged := cloneTable(table)
	savedIndexCounter := c.indexCounter
	key := tableKey{schema: schema.Name, table: table.Name}
	if tableIdx >= 0 {
		schema.Tables[tableIdx] = staged
		c.tableByName[key] = staged
	}

	for _, item := range n.Actions.Items {
		action, ok := item.(*ast.AlterTableAction)
		if !ok {
			continue
		}
		if err := c.execAlterTableAction(staged, schema.Name, action, source); err != nil {
			if tableIdx >= 0 {
				schema.Tables[tableIdx] = table
				c.tableByName[key] = table
			}
			c.indexCounter = savedIndexCounter
			return err
		}
	}
	return nil
}

func checkAlterTablePreLookup(n *ast.AlterTableStmt) error {
	seenDrop := make(map[string]bool)
	for _, item := range n.Actions.Items {
		action, ok := item.(*ast.AlterTableAction)
		if !ok {
			continue
		}
		switch action.Type {
		case ast.ATAddIndex:
			if action.Index == nil || action.Index.Columns == nil || len(action.Index.Columns.Items) == 0 {
				return mssqlErr(ErrAlterTableAddIndexSyntax, "Table level constraint or index does not specify column list.")
			}
		case ast.ATDropConstraint:
			if action.Constraint != nil && action.Constraint.Name != "" {
				low := strings.ToLower(action.Constraint.Name)
				if seenDrop[low] {
					return mssqlErr(ErrDupTargetInAlter, "Cannot create, drop, enable, or disable more than one constraint, column, index, or trigger named '%s' in this context.", action.Constraint.Name)
				}
				seenDrop[low] = true
			}
		case ast.ATCheckConstraint, ast.ATNocheckConstraint:
			if action.Names != nil {
				seen := make(map[string]bool)
				for _, ni := range action.Names.Items {
					if s, ok := ni.(*ast.String); ok {
						low := strings.ToLower(s.Str)
						if seen[low] {
							return mssqlErr(ErrDupTargetInAlter, "Cannot create, drop, enable, or disable more than one constraint, column, index, or trigger named '%s' in this context.", s.Str)
						}
						seen[low] = true
					}
				}
			}
		}
	}
	return nil
}

func (c *Catalog) execAlterTableAction(table *metadata.TableMetadata, schemaName string, action *ast.AlterTableAction, source string) error {
	switch action.Type {
	case ast.ATAddColumn:
		if action.Column == nil {
			return nil
		}
		if columnByName(table, action.Column.Name) != nil {
			return mssqlErr(ErrDupColumn, "Column names in each table must be unique. Column name '%s' in table '%s' is specified more than once.", action.Column.Name, table.Name)
		}
		if !c.lenientMetadataLoad && action.Column.Computed == nil && !isValidDataType(action.Column.DataType, false) {
			typeName := ""
			if action.Column.DataType != nil {
				typeName = action.Column.DataType.Name
			}
			return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #1: Cannot find data type %s.", typeName)
		}
		c.extractColumn(action.Column, table, schemaName, source)
		return nil

	case ast.ATDropColumn:
		col := columnByName(table, action.ColName)
		if col == nil {
			if action.IfExists {
				return nil
			}
			return mssqlErr(ErrAlterColNotFound, "ALTER TABLE DROP COLUMN %s failed because column '%s' does not exist in table '%s'.", action.ColName, action.ColName, table.Name)
		}
		if dep, ok := columnDependency(c, schemaName, table, col); ok {
			return mssqlErr(ErrColumnDependentObject, "The object '%s' is dependent on column '%s'.", dep, col.Name)
		}
		for i, cur := range table.Columns {
			if cur == col {
				table.Columns = append(table.Columns[:i], table.Columns[i+1:]...)
				break
			}
		}
		for i, cur := range table.Columns {
			cur.Position = int32(i + 1)
		}
		return nil

	case ast.ATAlterColumn:
		col := columnByName(table, action.ColName)
		if col == nil {
			return mssqlErr(ErrAlterColNotFound, "ALTER TABLE ALTER COLUMN %s failed because column '%s' does not exist in table '%s'.", action.ColName, action.ColName, table.Name)
		}
		if !c.lenientMetadataLoad && !isValidDataType(action.DataType, false) {
			typeName := ""
			if action.DataType != nil {
				typeName = action.DataType.Name
			}
			return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #1: Cannot find data type %s.", typeName)
		}
		col.Type = nodeText(action.DataType, source)
		if action.Column != nil && action.Column.Nullable != nil {
			col.Nullable = !action.Column.Nullable.NotNull
		}
		if action.Collation != "" {
			col.Collation = action.Collation
		}
		return nil

	case ast.ATAlterColumnAddDrop:
		return c.execAlterColumnAddDrop(table, action)

	case ast.ATAddConstraint:
		if action.Constraint == nil {
			return nil
		}
		return c.execAlterTableAddConstraint(table, schemaName, action.Constraint, source)

	case ast.ATDropConstraint:
		if action.Constraint == nil {
			if hasOnlineOption(action.Options) {
				return mssqlErr(ErrOnlineDropNonClustered, "Only a clustered index can be dropped online.")
			}
			return nil
		}
		name := action.Constraint.Name
		if hasOnlineOption(action.Options) {
			isClustered := false
			for _, idx := range table.Indexes {
				if idx.IsConstraint && strings.EqualFold(idx.Name, name) &&
					(strings.EqualFold(idx.Type, "CLUSTERED") || strings.EqualFold(idx.Type, "CLUSTERED COLUMNSTORE") || (idx.Type == "" && idx.Primary)) {
					isClustered = true
					break
				}
			}
			if !isClustered {
				return mssqlErr(ErrOnlineDropNonClustered, "Only a clustered index can be dropped online.")
			}
		}
		if c.dropTableConstraint(table, name) {
			return nil
		}
		if action.IfExists {
			return nil
		}
		return mssqlErr(ErrConstraintNotFound, "'%s' is not a constraint.", name)

	case ast.ATAddIndex:
		if action.Index != nil {
			if err := checkDuplicateIndexColumns(action.Index.Columns, action.Index.IncludeCols); err != nil {
				return err
			}
		}
		return mssqlErr(ErrAlterTableIndexDiskTable, "The operation 'ALTER TABLE ADD INDEX' is supported only with memory optimized tables.")

	case ast.ATDropIndex:
		return mssqlErr(ErrAlterTableIndexDiskTable, "The operation 'ALTER TABLE DROP INDEX' is supported only with memory optimized tables.")

	case ast.ATCheckConstraint, ast.ATNocheckConstraint:
		if action.Names == nil {
			return nil
		}
		for _, item := range action.Names.Items {
			s, ok := item.(*ast.String)
			if !ok {
				continue
			}
			if strings.EqualFold(s.Str, "ALL") {
				return nil
			}
			if hasCheckOrFKConstraint(table, s.Str) {
				continue
			}
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, s.Str, false) != kindNone {
				return mssqlErr(ErrCheckConstraintWrongKind, "Object '%s' cannot be disabled or enabled. This action applies only to foreign key and check constraints.", s.Str)
			}
			return mssqlErr(ErrCheckConstraintNotFound, "Constraint '%s' does not exist.", s.Str)
		}
		return nil

	case ast.ATRebuild:
		if action.Partition != nil {
			if _, ok := action.Partition.(*ast.ColumnRef); ok {
				return mssqlErr(ErrInvalidColumnName, "Invalid column name '%s'.", nodeText(action.Partition, source))
			}
			return mssqlErr(ErrRebuildPartitionInvalid, "Cannot specify partition number in the alter index statement as the table '%s' is not partitioned.", table.Name)
		}
		return nil

	case ast.ATSwitchPartition:
		if action.TargetName != nil {
			targetSchema, targetTable := tableRefSchemaObject(action.TargetName, c.defaultSchema)
			if strings.EqualFold(targetSchema, schemaName) && strings.EqualFold(targetTable, table.Name) {
				return mssqlErr(ErrSwitchSameTable, "ALTER TABLE SWITCH statement failed. The source table '%s' and target table '%s' are same.", table.Name, targetTable)
			}
		}
		return nil

	case ast.ATSet:
		return c.execAlterTableSet(schemaName, table, action.Options)

	case ast.ATAddPeriod:
		if action.Names != nil {
			for _, item := range action.Names.Items {
				if s, ok := item.(*ast.String); ok {
					if columnByName(table, s.Str) == nil {
						return mssqlErr(ErrAlterColNotFound, "Column '%s' does not exist in table '%s'.", s.Str, table.Name)
					}
					return mssqlErr(ErrSystemTimeInvalidCols, "Temporal generated always column '%s' has invalid data type.", s.Str)
				}
			}
		}
		return mssqlErr(ErrSystemTimeInvalidCols, "Temporal generated always column has invalid data type.")

	case ast.ATDropPeriod:
		return mssqlErr(ErrDropPeriodNotDefined, "SYSTEM_TIME period cannot be dropped from table '%s', because SYSTEM_TIME period does not exist.", table.Name)
	}
	return nil
}

func (c *Catalog) execAlterColumnAddDrop(table *metadata.TableMetadata, action *ast.AlterTableAction) error {
	var opt *ast.AlterColumnOption
	if action.Options != nil {
		for _, item := range action.Options.Items {
			if o, ok := item.(*ast.AlterColumnOption); ok {
				opt = o
				break
			}
		}
	}
	if opt != nil && strings.EqualFold(opt.Action, "DROP") && strings.EqualFold(opt.Option, "ROWGUIDCOL") {
		return mssqlErr(ErrDropRowguidcolMissing, "ALTER TABLE ALTER COLUMN DROP ROWGUIDCOL failed because a column does not exist in table '%s' with ROWGUIDCOL property.", table.Name)
	}
	col := columnByName(table, action.ColName)
	if col == nil {
		return mssqlErr(ErrAlterColNotFound, "ALTER TABLE ALTER COLUMN %s failed because column '%s' does not exist in table '%s'.", action.ColName, action.ColName, table.Name)
	}
	if opt == nil {
		return nil
	}
	if strings.EqualFold(opt.Action, "ADD") {
		switch strings.ToUpper(opt.Option) {
		case "ROWGUIDCOL":
			if normalizeType(col.Type) != "uniqueidentifier" {
				return mssqlErr(ErrRowguidcolType, "The ROWGUIDCOL property can only be specified on the uniqueidentifier data type.")
			}
		case "NOT FOR REPLICATION":
			if !col.IsIdentity {
				return mssqlErr(ErrAlterColNFRNotIdentity, "Cannot alter NOT FOR REPLICATION attribute on column '%s' in table '%s' because this column is not an identity column.", col.Name, table.Name)
			}
		case "PERSISTED":
			if !isComputedColumnType(col.Type) {
				return mssqlErr(ErrAlterColPersistedNonComp, "PERSISTED attribute cannot be altered on column '%s' because this column is not computed.", col.Name)
			}
		case "MASKED":
			if !isValidMaskingFunction(opt.MaskFunction) {
				return mssqlErr(ErrInvalidMaskingFunction, "Invalid data masking format for function '%s' in column '%s'.", opt.MaskFunction, col.Name)
			}
		case "HIDDEN":
			return mssqlErr(ErrAlterColHiddenNotPeriod, "Cannot alter HIDDEN attribute on column '%s' in table '%s' because this column is not a generated always column.", col.Name, table.Name)
		case "SPARSE":
			return nil
		}
	}
	return nil
}

func isValidMaskingFunction(fn string) bool {
	low := strings.ToLower(strings.TrimSpace(fn))
	return low == "default()" || low == "email()" || strings.HasPrefix(low, "partial(") || strings.HasPrefix(low, "random(") || strings.HasPrefix(low, "datetime(")
}

func (c *Catalog) execAlterTableAddConstraint(table *metadata.TableMetadata, schemaName string, n *ast.ConstraintDef, source string) error {
	if n.Type == ast.ConstraintEdge {
		return mssqlErr(ErrEdgeConstraintNonEdge, "Edge constraint cannot be created on table '%s'. The table is not an edge table.", table.Name)
	}
	switch n.Type {
	case ast.ConstraintPrimaryKey, ast.ConstraintUnique:
		if err := checkDuplicateIndexColumns(n.Columns, nil); err != nil {
			return err
		}
		for _, colName := range stringNodeList(n.Columns) {
			if columnByName(table, colName) == nil {
				return mssqlErr(ErrIndexColumnNotFound, "Column name '%s' does not exist in the target table or view.", colName)
			}
		}
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		c.extractTableConstraint(n, table, schemaName, source)
		return nil

	case ast.ConstraintDefault:
		cols := stringNodeList(n.Columns)
		if len(cols) == 0 {
			return nil
		}
		col := columnByName(table, cols[0])
		if col == nil {
			return mssqlErr(ErrAddDefaultColumnNotFound, "Cannot create default on column '%s' because it does not exist in table '%s'.", cols[0], table.Name)
		}
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		col.Default = nodeText(n.Expr, source)
		if n.Name != "" {
			col.DefaultConstraintName = n.Name
		} else if !c.omitUnnamedDefaultConstraintNames {
			col.DefaultConstraintName = "DF_" + table.Name + "_" + col.Name
		}
		return nil

	case ast.ConstraintForeignKey:
		cols := stringNodeList(n.Columns)
		if len(cols) == 0 {
			refName := ""
			if n.RefTable != nil {
				refName = n.RefTable.Object
			}
			return mssqlErr(ErrFKColumnCountMismatch, "The number of columns in the referencing column list does not match those of the primary key in the referenced table '%s'.", refName)
		}
		for _, colName := range cols {
			if columnByName(table, colName) == nil {
				return mssqlErr(ErrFKColumnNotFound, "Foreign key references invalid column '%s' in referencing table '%s'.", colName, table.Name)
			}
		}
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		c.extractTableConstraint(n, table, schemaName, source)
		return nil

	case ast.ConstraintCheck:
		if n.Name != "" {
			if schema := c.GetSchema(schemaName); schema != nil && c.lookupSchemaObject(schema, n.Name, false) != kindNone {
				return mssqlErr(ErrDupObject, "There is already an object named '%s' in the database.", n.Name)
			}
		}
		c.extractTableConstraint(n, table, schemaName, source)
		return nil
	}
	return nil
}

func (c *Catalog) dropTableConstraint(table *metadata.TableMetadata, name string) bool {
	for i, idx := range table.Indexes {
		if idx.IsConstraint && strings.EqualFold(idx.Name, name) {
			table.Indexes = append(table.Indexes[:i], table.Indexes[i+1:]...)
			return true
		}
	}
	for i, fk := range table.ForeignKeys {
		if strings.EqualFold(fk.Name, name) {
			table.ForeignKeys = append(table.ForeignKeys[:i], table.ForeignKeys[i+1:]...)
			return true
		}
	}
	for i, ck := range table.CheckConstraints {
		if strings.EqualFold(ck.Name, name) {
			table.CheckConstraints = append(table.CheckConstraints[:i], table.CheckConstraints[i+1:]...)
			return true
		}
	}
	for _, col := range table.Columns {
		if col.DefaultConstraintName != "" && strings.EqualFold(col.DefaultConstraintName, name) {
			col.Default = ""
			col.DefaultConstraintName = ""
			return true
		}
	}
	return false
}

func hasCheckOrFKConstraint(table *metadata.TableMetadata, name string) bool {
	for _, fk := range table.ForeignKeys {
		if strings.EqualFold(fk.Name, name) {
			return true
		}
	}
	for _, ck := range table.CheckConstraints {
		if strings.EqualFold(ck.Name, name) {
			return true
		}
	}
	return false
}

func columnDependency(c *Catalog, schemaName string, table *metadata.TableMetadata, col *metadata.ColumnMetadata) (string, bool) {
	if col.Default != "" {
		name := col.DefaultConstraintName
		if name == "" {
			name = "DEFAULT"
		}
		return name, true
	}
	for _, idx := range table.Indexes {
		for _, expr := range idx.Expressions {
			if strings.EqualFold(expr, col.Name) {
				return idx.Name, true
			}
		}
		for _, inc := range indexIncludeColumns(idx.Definition) {
			if strings.EqualFold(inc, col.Name) {
				return idx.Name, true
			}
		}
	}
	for _, fk := range table.ForeignKeys {
		for _, fc := range fk.Columns {
			if strings.EqualFold(fc, col.Name) {
				return fk.Name, true
			}
		}
	}
	for _, ck := range table.CheckConstraints {
		if checkExprReferencesColumn(ck.Expression, table.Name, col.Name) {
			name := ck.Name
			if name == "" {
				name = "CHECK"
			}
			return name, true
		}
	}
	for _, other := range table.Columns {
		if other == col || !isComputedColumnType(other.Type) {
			continue
		}
		if computedExprReferencesColumn(other.Type, table.Name, col.Name) {
			return other.Name, true
		}
	}
	for _, s := range c.schemas {
		for _, t := range s.Tables {
			for _, fk := range t.ForeignKeys {
				refSchema := fk.ReferencedSchema
				if refSchema == "" {
					refSchema = c.defaultSchema
				}
				if strings.EqualFold(refSchema, schemaName) && strings.EqualFold(fk.ReferencedTable, table.Name) {
					for _, rc := range fk.ReferencedColumns {
						if strings.EqualFold(rc, col.Name) {
							return fk.Name, true
						}
					}
				}
			}
		}
	}
	return "", false
}

func exprReferencesColumn(node ast.Node, colName string) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if ref, ok := n.(*ast.ColumnRef); ok && strings.EqualFold(ref.Column, colName) {
			found = true
			return false
		}
		return true
	})
	return found
}

func exprReferencesTableColumn(node ast.Node, tableName, colName string) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		ref, ok := n.(*ast.ColumnRef)
		if !ok {
			return true
		}
		if strings.EqualFold(ref.Column, colName) && (ref.Table == "" || strings.EqualFold(ref.Table, tableName)) {
			found = true
			return false
		}
		return true
	})
	return found
}

func checkExprReferencesColumn(expr, tableName, colName string) bool {
	if strings.TrimSpace(expr) == "" {
		return false
	}
	stmts, err := mssql.Parse("SELECT 1 WHERE (" + expr + ")")
	if err != nil || len(stmts) == 0 {
		return false
	}
	sel, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok || sel.WhereClause == nil {
		return false
	}
	return exprReferencesTableColumn(sel.WhereClause, tableName, colName)
}

func computedExprReferencesColumn(colType, tableName, colName string) bool {
	if !isComputedColumnType(colType) {
		return false
	}
	stmts, err := mssql.Parse("CREATE TABLE _ (_ " + colType + ")")
	if err != nil || len(stmts) == 0 {
		return false
	}
	ct, ok := stmts[0].AST.(*ast.CreateTableStmt)
	if !ok || ct.Columns == nil || len(ct.Columns.Items) == 0 {
		return false
	}
	cd, ok := ct.Columns.Items[0].(*ast.ColumnDef)
	if !ok || cd.Computed == nil || cd.Computed.Expr == nil {
		return false
	}
	return exprReferencesTableColumn(cd.Computed.Expr, tableName, colName)
}

// hasOptionOn reports whether opts contains optionName set to ON.
// Parser option lists use either a single "NAME=VALUE" entry (CREATE/DROP INDEX)
// or paired "NAME", "VALUE" entries (ALTER TABLE DROP CONSTRAINT).
func hasOptionOn(opts *ast.List, optionName string) bool {
	if opts == nil {
		return false
	}
	for i := 0; i < len(opts.Items); i++ {
		s, ok := opts.Items[i].(*ast.String)
		if !ok {
			continue
		}
		raw := strings.TrimSpace(s.Str)
		if k, v, hasEq := strings.Cut(raw, "="); hasEq {
			if strings.EqualFold(strings.TrimSpace(k), optionName) {
				return strings.EqualFold(strings.TrimSpace(v), "ON")
			}
			continue
		}
		if strings.EqualFold(raw, optionName) && i+1 < len(opts.Items) {
			if next, ok := opts.Items[i+1].(*ast.String); ok {
				return strings.EqualFold(strings.TrimSpace(next.Str), "ON")
			}
		}
	}
	return false
}

func hasOnlineOption(opts *ast.List) bool {
	return hasOptionOn(opts, "ONLINE")
}

func (c *Catalog) execAlterTableSet(schemaName string, table *metadata.TableMetadata, opts *ast.List) error {
	if opts == nil {
		return nil
	}
	for i := 0; i < len(opts.Items); i++ {
		s, ok := opts.Items[i].(*ast.String)
		if !ok {
			continue
		}
		name := strings.ToUpper(strings.TrimSpace(s.Str))
		val := ""
		if i+1 < len(opts.Items) {
			if vs, ok := opts.Items[i+1].(*ast.String); ok {
				val = strings.ToUpper(strings.TrimSpace(vs.Str))
			}
		}
		switch name {
		case "SYSTEM_VERSIONING":
			if val == "ON" {
				if !c.systemTimeTables[tableKey{schema: normalizeOrDefault(schemaName, c.defaultSchema), table: table.Name}] {
					return mssqlErr(ErrSystemVersioningNoPeriod, "Cannot set SYSTEM_VERSIONING to ON when SYSTEM_TIME period is not defined.")
				}
				return nil
			}
			if val == "OFF" {
				return mssqlErr(ErrSystemVersioningNotEnabled, "SYSTEM_VERSIONING is not turned ON for table '%s'.", table.Name)
			}
		case "FILESTREAM_ON":
			if !c.lenientMetadataLoad && !isValidFilegroup(val) {
				return mssqlErr(ErrInvalidFilegroup, "Invalid filegroup '%s' specified.", val)
			}
		case "FILETABLE_DIRECTORY":
			return mssqlErr(ErrFileTableDirectoryNonFT, "The option 'FILETABLE_DIRECTORY' is only valid when used on a FileTable.")
		}
	}
	return nil
}

func checkExecuteAsOptions(opts *ast.List) error {
	if opts == nil {
		return nil
	}
	for _, item := range opts.Items {
		s, ok := item.(*ast.String)
		if !ok {
			continue
		}
		upper := strings.ToUpper(strings.TrimSpace(s.Str))
		if strings.HasPrefix(upper, "EXECUTE AS ") {
			target := strings.TrimSpace(s.Str[len("EXECUTE AS "):])
			unquoted := strings.Trim(target, "'\"[]")
			switch strings.ToUpper(unquoted) {
			case "CALLER", "SELF", "OWNER", "DBO":
			default:
				return mssqlErr(ErrSecurableNotFound, "Cannot execute as the database principal because the principal \"%s\" does not exist.", unquoted)
			}
		}
	}
	return nil
}

func isNonNullLiteralExpr(expr ast.ExprNode) bool {
	if expr == nil {
		return false
	}
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return isNonNullLiteralExpr(e.Expr)
	case *ast.Literal:
		return e.Type != ast.LitNull
	}
	return false
}

// --- DROP ---

func (c *Catalog) execDrop(n *ast.DropStmt) error {
	switch n.ObjectType {
	case ast.DropSchema:
		return c.execDropSchema(n)
	case ast.DropIndex:
		return c.execDropIndex(n)
	case ast.DropTable:
		return c.execDropObjects(n, kindTable, "TABLE")
	case ast.DropView:
		return c.execDropObjects(n, kindView, "VIEW")
	case ast.DropProcedure:
		return c.execDropObjects(n, kindProcedure, "PROCEDURE")
	case ast.DropFunction:
		return c.execDropObjects(n, kindFunction, "FUNCTION")
	case ast.DropSequence:
		return c.execDropObjects(n, kindSequence, "SEQUENCE")
	case ast.DropSynonym:
		return c.execDropObjects(n, kindSynonym, "SYNONYM")
	case ast.DropTrigger:
		if n.OnDatabase || n.OnAllServer {
			if n.IfExists {
				return nil
			}
			name := ""
			if n.Names != nil && len(n.Names.Items) > 0 {
				if ref, ok := n.Names.Items[0].(*ast.TableRef); ok {
					name = ref.Object
				}
			}
			return mssqlErr(ErrCannotDropNotFound, "Cannot drop the trigger '%s', because it does not exist or you do not have permission.", name)
		}
		return c.execDropObjects(n, kindTrigger, "TRIGGER")
	}
	return nil
}

func (c *Catalog) execDropSchema(n *ast.DropStmt) error {
	if n.Names == nil {
		return nil
	}
	for _, item := range n.Names.Items {
		ref, ok := item.(*ast.TableRef)
		if !ok {
			continue
		}
		name := ref.Object
		schema := c.GetSchema(name)
		if schema == nil {
			if n.IfExists {
				continue
			}
			return mssqlErr(ErrSecurableNotFound, "Cannot drop the schema '%s', because it does not exist or you do not have permission.", name)
		}
		if len(schema.Tables) > 0 || len(schema.Views) > 0 || len(schema.Procedures) > 0 ||
			len(schema.Functions) > 0 || len(schema.Sequences) > 0 || len(schema.Synonyms) > 0 {
			return mssqlErr(ErrDropSchemaNotEmpty, "Cannot drop schema '%s' because it is being referenced by an object.", name)
		}
		c.removeSchema(schema.Name)
	}
	return nil
}

func (c *Catalog) execDropIndex(n *ast.DropStmt) error {
	if n.Names == nil {
		return nil
	}
	type dropIndexTarget struct {
		table *metadata.TableMetadata
		idx   *metadata.IndexMetadata
	}
	seen := make(map[*metadata.IndexMetadata]bool)
	var targets []dropIndexTarget
	for i, item := range n.Names.Items {
		ref, ok := item.(*ast.TableRef)
		if !ok {
			continue
		}
		var schemaName, tableName, indexName string
		if n.OnTables != nil && i < len(n.OnTables.Items) && n.OnTables.Items[i] != nil {
			onRef, _ := n.OnTables.Items[i].(*ast.TableRef)
			schemaName, tableName = tableRefSchemaObject(onRef, c.defaultSchema)
			indexName = ref.Object
		} else {
			schemaName = ref.Database
			tableName = ref.Schema
			indexName = ref.Object
			if schemaName == "" {
				schemaName = c.defaultSchema
			}
		}

		table := c.GetTable(schemaName, tableName)
		var foundIdx *metadata.IndexMetadata
		if table != nil {
			for _, idx := range table.Indexes {
				if strings.EqualFold(idx.Name, indexName) && !idx.IsConstraint {
					foundIdx = idx
					break
				}
			}
		}
		if foundIdx == nil || seen[foundIdx] {
			if n.IfExists {
				continue
			}
			return mssqlErr(ErrCannotDropNotFound, "Cannot drop the index '%s.%s', because it does not exist or you do not have permission.", tableName, indexName)
		}

		isClustered := strings.EqualFold(foundIdx.Type, "CLUSTERED") || strings.EqualFold(foundIdx.Type, "CLUSTERED COLUMNSTORE")
		if err := checkDropIndexOptions(n.Options, tableName, indexName, isClustered); err != nil {
			return err
		}

		seen[foundIdx] = true
		targets = append(targets, dropIndexTarget{table: table, idx: foundIdx})
	}
	for _, target := range targets {
		target.table.Indexes = slices.DeleteFunc(target.table.Indexes, func(idx *metadata.IndexMetadata) bool {
			return idx == target.idx
		})
	}
	return nil
}

func checkDropIndexOptions(opts *ast.List, tableName, indexName string, isClustered bool) error {
	if opts == nil {
		return nil
	}
	if hasOptionOn(opts, "ONLINE") && !isClustered {
		return mssqlErr(ErrOnlineDropNonClustered, "Only a clustered index can be dropped online.")
	}
	for _, item := range opts.Items {
		s, ok := item.(*ast.String)
		if !ok {
			continue
		}
		upper := strings.ToUpper(strings.TrimSpace(s.Str))
		switch {
		case strings.HasPrefix(upper, "FILESTREAM_ON"):
			return mssqlErr(ErrDropIndexFilestreamOn, "FILESTREAM_ON cannot be specified when a table has no FILESTREAM columns.")
		case strings.HasPrefix(upper, "MAXDOP") || strings.HasPrefix(upper, "MOVE"):
			if !isClustered {
				return mssqlErr(ErrDropIndexClusteredOpt, "Cannot drop non-clustered index '%s.%s' using drop clustered index clause.", tableName, indexName)
			}
		}
	}
	return nil
}

func (c *Catalog) execDropObjects(n *ast.DropStmt, wantKind objectKind, sqlKind string) error {
	if n.Names == nil {
		return nil
	}
	type dropTarget struct {
		schema *metadata.SchemaMetadata
		name   string
	}
	seen := make(map[tableKey]bool)
	var targets []dropTarget
	for _, item := range n.Names.Items {
		ref, ok := item.(*ast.TableRef)
		if !ok {
			continue
		}
		schemaName, objName := tableRefSchemaObject(ref, c.defaultSchema)
		schema := c.GetSchema(schemaName)
		key := tableKey{schema: normalize(schemaName), table: normalize(objName)}
		kind := c.lookupSchemaObject(schema, objName, false)
		if kind == kindNone || seen[key] {
			if n.IfExists {
				continue
			}
			return mssqlErr(ErrCannotDropNotFound, "Cannot drop the %s '%s', because it does not exist or you do not have permission.", strings.ToLower(sqlKind), objName)
		}
		if kind != wantKind {
			return mssqlErr(ErrCannotDropWrongKind, "Cannot use DROP %s with '%s' because '%s' is a %s.", sqlKind, objName, objName, string(kind))
		}
		if wantKind == kindTable {
			if fkName, referenced := hasInboundFKExcluding(c, schema.Name, objName, seen); referenced {
				return mssqlErr(ErrDropTableReferencedByFK, "Could not drop object '%s' because it is referenced by a FOREIGN KEY constraint '%s'.", objName, fkName)
			}
		}
		seen[key] = true
		targets = append(targets, dropTarget{schema: schema, name: objName})
	}
	for _, target := range targets {
		c.removeObject(target.schema, wantKind, target.name)
	}
	return nil
}

func (c *Catalog) removeObject(schema *metadata.SchemaMetadata, kind objectKind, name string) {
	switch kind {
	case kindTable:
		target := c.resolveTableInSchema(schema, name)
		for i, t := range schema.Tables {
			if t == target || (target == nil && strings.EqualFold(t.Name, name)) {
				delete(c.tableByName, tableKey{schema: schema.Name, table: t.Name})
				delete(c.systemTimeTables, tableKey{schema: normalize(schema.Name), table: t.Name})
				schema.Tables = append(schema.Tables[:i], schema.Tables[i+1:]...)
				return
			}
		}
	case kindView:
		schema.Views = slices.DeleteFunc(schema.Views, func(v *metadata.ViewMetadata) bool {
			return strings.EqualFold(v.Name, name)
		})
	case kindProcedure:
		schema.Procedures = slices.DeleteFunc(schema.Procedures, func(p *metadata.ProcedureMetadata) bool {
			return strings.EqualFold(p.Name, name)
		})
	case kindFunction:
		schema.Functions = slices.DeleteFunc(schema.Functions, func(fn *metadata.FunctionMetadata) bool {
			return strings.EqualFold(fn.Name, name)
		})
	case kindSequence:
		schema.Sequences = slices.DeleteFunc(schema.Sequences, func(sq *metadata.SequenceMetadata) bool {
			return strings.EqualFold(sq.Name, name)
		})
	case kindSynonym:
		schema.Synonyms = slices.DeleteFunc(schema.Synonyms, func(sy *metadata.SynonymMetadata) bool {
			return strings.EqualFold(sy.Name, name)
		})
	case kindTrigger:
		for _, t := range schema.Tables {
			t.Triggers = slices.DeleteFunc(t.Triggers, func(tr *metadata.TriggerMetadata) bool {
				return strings.EqualFold(tr.Name, name)
			})
		}
		for _, v := range schema.Views {
			v.Triggers = slices.DeleteFunc(v.Triggers, func(tr *metadata.TriggerMetadata) bool {
				return strings.EqualFold(tr.Name, name)
			})
		}
	}
}

// --- TRUNCATE TABLE ---

func (c *Catalog) execTruncateTable(n *ast.TruncateStmt) error {
	schemaName, tableName := tableRefSchemaObject(n.Table, c.defaultSchema)
	schema := c.GetSchema(schemaName)
	kind := c.lookupSchemaObject(schema, tableName, false)
	switch kind {
	case kindNone:
		return mssqlErr(ErrTruncateNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	case kindTable:
	default:
		return mssqlErr(ErrTruncateNotTable, "Cannot truncate object '%s' because it is not a table.", tableName)
	}
	table := c.resolveTableInSchema(schema, tableName)
	if table == nil {
		return mssqlErr(ErrTruncateNotFound, "Cannot find the object \"%s\" because it does not exist or you do not have permissions.", tableName)
	}
	if fkName, referenced := hasInboundFK(c, schema.Name, table.Name); referenced {
		return mssqlErr(ErrTruncateReferencedByFK, "Cannot truncate table '%s' because it is being referenced by a FOREIGN KEY constraint '%s'.", tableName, fkName)
	}
	return nil
}
