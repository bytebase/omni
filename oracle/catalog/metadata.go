package catalog

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle"
	"github.com/bytebase/omni/oracle/ast"
)

// metadataFromSQL parses Oracle DDL and returns the schema it defines.
func metadataFromSQL(schemaText string) (*metadatapb.DatabaseSchemaMetadata, error) {
	stmts, err := oracle.Parse(schemaText)
	if err != nil {
		return nil, fmt.Errorf("parse Oracle schema: %w", err)
	}
	if len(stmts) == 0 {
		return nil, fmt.Errorf("no parse results")
	}

	extractor := newMetadataExtractor(schemaText)
	for _, stmt := range stmts {
		node, ok := stmt.AST.(ast.StmtNode)
		if !ok || node == nil {
			continue
		}
		extractor.extractStatement(node)
	}

	if extractor.err != nil {
		return nil, extractor.err
	}
	return extractor.databaseMetadata(), nil
}

func newMetadataExtractor(schemaText string) *oracleMetadataExtractor {
	return &oracleMetadataExtractor{
		tables:            make(map[string]*metadatapb.TableMetadata),
		views:             make(map[string]*metadatapb.ViewMetadata),
		materializedViews: make(map[string]*metadatapb.MaterializedViewMetadata),
		functions:         make(map[string]*metadatapb.FunctionMetadata),
		procedures:        make(map[string]*metadatapb.ProcedureMetadata),
		triggers:          make(map[string]*metadatapb.TriggerMetadata),
		sequences:         make(map[string]*metadatapb.SequenceMetadata),
		packages:          make(map[string]*metadatapb.PackageMetadata),
		checkNames:        make(map[string]map[string]bool),
		checkNameNext:     make(map[string]map[string]int),
		schemaText:        schemaText,
	}
}

type oracleMetadataExtractor struct {
	err               error
	currentDatabase   string
	currentSchema     string
	tables            map[string]*metadatapb.TableMetadata
	views             map[string]*metadatapb.ViewMetadata
	materializedViews map[string]*metadatapb.MaterializedViewMetadata
	functions         map[string]*metadatapb.FunctionMetadata
	procedures        map[string]*metadatapb.ProcedureMetadata
	triggers          map[string]*metadatapb.TriggerMetadata
	sequences         map[string]*metadatapb.SequenceMetadata
	packages          map[string]*metadatapb.PackageMetadata
	checkNames        map[string]map[string]bool
	checkNameNext     map[string]map[string]int
	schemaText        string
}

func (e *oracleMetadataExtractor) databaseMetadata() *metadatapb.DatabaseSchemaMetadata {
	e.resolveForeignKeyReferencedColumns()

	schemaMetadata := &metadatapb.DatabaseSchemaMetadata{
		Name:    e.currentDatabase,
		Schemas: []*metadatapb.SchemaMetadata{},
	}
	schema := &metadatapb.SchemaMetadata{
		Name:              e.currentSchema,
		Tables:            []*metadatapb.TableMetadata{},
		Views:             []*metadatapb.ViewMetadata{},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{},
		Procedures:        []*metadatapb.ProcedureMetadata{},
		Functions:         []*metadatapb.FunctionMetadata{},
		Sequences:         []*metadatapb.SequenceMetadata{},
		Packages:          []*metadatapb.PackageMetadata{},
	}

	var tableNames []string
	for name := range e.tables {
		if _, ok := e.materializedViews[name]; !ok {
			tableNames = append(tableNames, name)
		}
	}
	slices.Sort(tableNames)
	for _, name := range tableNames {
		schema.Tables = append(schema.Tables, e.tables[name])
	}

	var viewNames []string
	for name := range e.views {
		viewNames = append(viewNames, name)
	}
	slices.Sort(viewNames)
	for _, name := range viewNames {
		schema.Views = append(schema.Views, e.views[name])
	}

	var materializedViewNames []string
	for name := range e.materializedViews {
		materializedViewNames = append(materializedViewNames, name)
	}
	slices.Sort(materializedViewNames)
	for _, name := range materializedViewNames {
		schema.MaterializedViews = append(schema.MaterializedViews, e.materializedViews[name])
	}

	var functionNames []string
	for name := range e.functions {
		functionNames = append(functionNames, name)
	}
	slices.Sort(functionNames)
	for _, name := range functionNames {
		schema.Functions = append(schema.Functions, e.functions[name])
	}

	var procedureNames []string
	for name := range e.procedures {
		procedureNames = append(procedureNames, name)
	}
	slices.Sort(procedureNames)
	for _, name := range procedureNames {
		schema.Procedures = append(schema.Procedures, e.procedures[name])
	}

	var sequenceNames []string
	for name := range e.sequences {
		sequenceNames = append(sequenceNames, name)
	}
	slices.Sort(sequenceNames)
	for _, name := range sequenceNames {
		schema.Sequences = append(schema.Sequences, e.sequences[name])
	}

	var packageNames []string
	for name := range e.packages {
		packageNames = append(packageNames, name)
	}
	slices.Sort(packageNames)
	for _, name := range packageNames {
		schema.Packages = append(schema.Packages, e.packages[name])
	}

	schemaMetadata.Schemas = append(schemaMetadata.Schemas, schema)
	return schemaMetadata
}

func (e *oracleMetadataExtractor) getOrCreateTable(tableName string) *metadatapb.TableMetadata {
	if table, ok := e.tables[tableName]; ok {
		return table
	}
	table := &metadatapb.TableMetadata{
		Name:             tableName,
		Columns:          []*metadatapb.ColumnMetadata{},
		Indexes:          []*metadatapb.IndexMetadata{},
		ForeignKeys:      []*metadatapb.ForeignKeyMetadata{},
		CheckConstraints: []*metadatapb.CheckConstraintMetadata{},
		Triggers:         []*metadatapb.TriggerMetadata{},
		Partitions:       []*metadatapb.TablePartitionMetadata{},
	}
	e.tables[tableName] = table
	return table
}

// useSchema keeps one schema in the extractor. A statement that names another
// owner would move earlier objects, so it is refused rather than merged.
func (e *oracleMetadataExtractor) useSchema(schema string) bool {
	if schema == "" {
		return true
	}
	if e.currentSchema == "" {
		e.currentSchema = schema
		return true
	}
	if schema == e.currentSchema {
		return true
	}
	if e.err == nil {
		e.err = fmt.Errorf("statement names schema %q, but the catalog is loading schema %q", schema, e.currentSchema)
	}
	return false
}

func (e *oracleMetadataExtractor) extractStatement(stmt ast.StmtNode) {
	switch n := stmt.(type) {
	case *ast.CreateTableStmt:
		e.extractCreateTable(n)
	case *ast.CreateIndexStmt:
		e.extractCreateIndex(n)
	case *ast.CreateViewStmt:
		e.extractCreateView(n)
	case *ast.CreateSchemaStmt:
		e.extractCreateSchema(n)
	case *ast.CreateSequenceStmt:
		e.extractCreateSequence(n)
	case *ast.CreateProcedureStmt:
		if name := objectName(n.Name); name != "" {
			e.procedures[name] = &metadatapb.ProcedureMetadata{Name: name, Definition: e.definitionText(n)}
		}
	case *ast.CreateFunctionStmt:
		if name := objectName(n.Name); name != "" {
			e.functions[name] = &metadatapb.FunctionMetadata{Name: name, Definition: e.functionDefinitionText(n)}
		}
	case *ast.CreatePackageStmt:
		if name := objectName(n.Name); name != "" {
			e.extractCreatePackage(name, n)
		}
	case *ast.CreateTriggerStmt:
		e.extractCreateTrigger(n)
	case *ast.AlterTableStmt:
		e.extractAlterTable(n)
	case *ast.CommentStmt:
		e.extractComment(n)
	case *ast.AlterSequenceStmt:
		e.extractAlterSequence(n)
	case *ast.DropStmt:
		e.extractDrop(n)
	default:
	}
}

func (e *oracleMetadataExtractor) extractCreateSchema(n *ast.CreateSchemaStmt) {
	if n.SchemaName != "" {
		e.currentSchema = n.SchemaName
	}
	for _, item := range listItems(n.Stmts) {
		stmt, ok := item.(ast.StmtNode)
		if !ok {
			continue
		}
		e.extractStatement(stmt)
	}
}

func (e *oracleMetadataExtractor) extractCreateTable(n *ast.CreateTableStmt) {
	tableName := objectName(n.Name)
	if tableName == "" {
		return
	}
	if !e.useSchema(n.Name.Schema) {
		return
	}

	table := e.getOrCreateTable(tableName)
	for _, item := range listItems(n.Columns) {
		column, ok := item.(*ast.ColumnDef)
		if !ok {
			continue
		}
		e.extractColumn(column, table)
	}
	for _, item := range listItems(n.Constraints) {
		constraint, ok := item.(*ast.TableConstraint)
		if !ok {
			continue
		}
		e.extractTableConstraint(constraint, table)
	}
}

func (e *oracleMetadataExtractor) extractColumn(n *ast.ColumnDef, table *metadatapb.TableMetadata) {
	column := &metadatapb.ColumnMetadata{
		Name:     n.Name,
		Nullable: true,
		Position: int32(len(table.Columns) + 1),
	}
	if n.TypeName != nil {
		column.Type = normalizeDataTypeText(e.nodeText(n.TypeName))
	} else if n.Domain != nil {
		column.Type = objectName(n.Domain)
	}
	if n.Default != nil {
		column.Default = e.nodeText(n.Default)
	}
	if n.Virtual != nil {
		column.Generation = &metadatapb.GenerationMetadata{
			Type:       metadatapb.GenerationMetadata_TYPE_VIRTUAL,
			Expression: e.exprText(n.Virtual),
		}
	}
	column.IsInvisible = n.Invisible
	column.DefaultOnNull = n.DefaultOnNull
	if n.Collation != "" {
		column.Collation = n.Collation
	}
	if n.Identity != nil {
		applyIdentity(column, n.Identity, e.nodeText(n.Identity))
		column.Nullable = false
	}
	if n.NotNull {
		column.Nullable = false
	}
	if n.Null {
		column.Nullable = true
	}

	for _, item := range listItems(n.Constraints) {
		constraint, ok := item.(*ast.ColumnConstraint)
		if !ok {
			continue
		}
		e.extractColumnConstraint(constraint, table, column)
		if constraint.Type == ast.CONSTRAINT_NOT_NULL && constraint.Name != "" {
			e.addNamedNotNull(table, constraint.Name, column.Name)
		}
	}

	table.Columns = append(table.Columns, column)
}

// addNamedNotNull records a named NOT NULL constraint the way sync reports it:
// as a check constraint on the column.
func (e *oracleMetadataExtractor) addNamedNotNull(table *metadatapb.TableMetadata, name, column string) {
	for _, check := range table.CheckConstraints {
		if check.Name == name {
			return
		}
	}
	table.CheckConstraints = append(table.CheckConstraints, &metadatapb.CheckConstraintMetadata{
		Name:       name,
		Expression: `"` + column + `" IS NOT NULL`,
	})
}

func (e *oracleMetadataExtractor) removeNamedNotNull(table *metadatapb.TableMetadata, column string) {
	want := `"` + column + `" IS NOT NULL`
	var checks []*metadatapb.CheckConstraintMetadata
	for _, check := range table.CheckConstraints {
		if !strings.EqualFold(strings.Join(strings.Fields(check.Expression), " "), want) {
			checks = append(checks, check)
		}
	}
	table.CheckConstraints = checks
}

func (e *oracleMetadataExtractor) clearNamedNotNull(table *metadatapb.TableMetadata, check *metadatapb.CheckConstraintMetadata) {
	expression := strings.TrimSpace(check.Expression)
	suffix := " IS NOT NULL"
	if !strings.HasSuffix(strings.ToUpper(expression), suffix) {
		return
	}
	columnName := strings.TrimSpace(expression[:len(expression)-len(suffix)])
	columnName = strings.Trim(columnName, `"`)
	for _, other := range table.CheckConstraints {
		if other.Name == check.Name {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(other.Expression), suffix)), `"`+columnName+`"`) {
			return
		}
	}
	for _, column := range table.Columns {
		if column.Name == columnName {
			column.Nullable = true
			return
		}
	}
}

func (e *oracleMetadataExtractor) extractColumnConstraint(n *ast.ColumnConstraint, table *metadatapb.TableMetadata, column *metadatapb.ColumnMetadata) {
	switch n.Type {
	case ast.CONSTRAINT_NOT_NULL:
		column.Nullable = false
	case ast.CONSTRAINT_NULL:
		column.Nullable = true
	case ast.CONSTRAINT_DEFAULT:
		if n.Expr != nil {
			column.Default = e.nodeText(n.Expr)
		}
	case ast.CONSTRAINT_PRIMARY:
		column.Nullable = false
		if !hasPrimaryIndex(table) {
			table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
				Name:         fallbackName(n.Name, fmt.Sprintf("PK_%s", table.Name)),
				Primary:      true,
				Unique:       true,
				Type:         "NORMAL",
				Expressions:  []string{column.Name},
				Visible:      true,
				IsConstraint: true,
			})
		}
	case ast.CONSTRAINT_UNIQUE:
		name := fallbackName(n.Name, fmt.Sprintf("UK_%s_%s", table.Name, column.Name))
		if !hasUniqueIndex(table, []string{column.Name}) {
			table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
				Name:         name,
				Unique:       true,
				Type:         "NORMAL",
				Expressions:  []string{column.Name},
				Visible:      true,
				IsConstraint: true,
			})
		}
	case ast.CONSTRAINT_CHECK:
		e.appendCheckConstraint(table, e.checkConstraintName(table, n.Name, fmt.Sprintf("CHK_%s_%s", table.Name, column.Name)), n.Expr)
	case ast.CONSTRAINT_FOREIGN:
		e.appendForeignKey(table, fallbackName(n.Name, fmt.Sprintf("FK_%s_%s", table.Name, column.Name)), []string{column.Name}, n.RefTable, n.RefColumns, n.OnDelete)
	default:
	}
}

func (e *oracleMetadataExtractor) extractTableConstraint(n *ast.TableConstraint, table *metadatapb.TableMetadata) {
	columns := stringList(n.Columns)
	switch n.Type {
	case ast.CONSTRAINT_PRIMARY:
		if len(columns) == 0 {
			return
		}
		table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
			Name:         fallbackName(n.Name, fmt.Sprintf("PK_%s", table.Name)),
			Primary:      true,
			Unique:       true,
			Type:         "NORMAL",
			Expressions:  columns,
			Visible:      true,
			IsConstraint: true,
		})
		markColumnsNotNull(table, columns)
	case ast.CONSTRAINT_UNIQUE:
		if len(columns) == 0 {
			return
		}
		table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
			Name:         fallbackName(n.Name, fmt.Sprintf("UK_%s_%d", table.Name, len(table.Indexes)+1)),
			Unique:       true,
			Type:         "NORMAL",
			Expressions:  columns,
			Visible:      false,
			IsConstraint: true,
		})
	case ast.CONSTRAINT_CHECK:
		e.appendCheckConstraint(table, e.checkConstraintName(table, n.Name, fmt.Sprintf("CHK_%s_%d", table.Name, len(table.CheckConstraints)+1)), n.Expr)
	case ast.CONSTRAINT_FOREIGN:
		e.appendForeignKey(table, fallbackName(n.Name, fmt.Sprintf("FK_%s_%d", table.Name, len(table.ForeignKeys)+1)), columns, n.RefTable, n.RefColumns, n.OnDelete)
	default:
	}
}

func (e *oracleMetadataExtractor) extractAlterTable(n *ast.AlterTableStmt) {
	tableName := objectName(n.Name)
	if tableName == "" {
		return
	}
	table := e.getOrCreateTable(tableName)
	for _, item := range listItems(n.Actions) {
		action, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch action.Action {
		case ast.AT_ADD_COLUMN:
			if action.ColumnDef != nil {
				e.extractColumn(action.ColumnDef, table)
			}
			for _, item := range listItems(action.ColumnDefs) {
				column, ok := item.(*ast.ColumnDef)
				if !ok {
					continue
				}
				e.extractColumn(column, table)
			}
		case ast.AT_MODIFY_COLUMN:
			if action.ColumnDef != nil {
				e.modifyColumn(action.ColumnDef, table)
			}
			for _, item := range listItems(action.ColumnDefs) {
				column, ok := item.(*ast.ColumnDef)
				if !ok {
					continue
				}
				e.modifyColumn(column, table)
			}
		case ast.AT_DROP_COLUMN:
			e.dropColumn(table, action.ColumnName)
		case ast.AT_RENAME_COLUMN:
			e.renameColumn(table, action.ColumnName, action.NewName)
		case ast.AT_ADD_CONSTRAINT:
			if action.Constraint != nil {
				e.extractTableConstraint(action.Constraint, table)
			}
		case ast.AT_DROP_CONSTRAINT:
			if action.Constraint != nil {
				e.dropConstraint(table, action.Constraint.Name)
			}
		case ast.AT_RENAME_CONSTRAINT:
			e.renameConstraint(table, action.ColumnName, action.NewName)
		case ast.AT_RENAME:
			if action.NewName != "" {
				e.renameTable(tableName, action.NewName)
			}
		default:
		}
	}
}

// modifyColumn applies an ALTER TABLE MODIFY column definition to the column
// the table already holds. Oracle keeps the attributes the clause does not
// name, including a type or NOT NULL when the clause only changes visibility.
func (e *oracleMetadataExtractor) modifyColumn(def *ast.ColumnDef, table *metadatapb.TableMetadata) {
	for _, column := range table.Columns {
		if column.Name != def.Name {
			continue
		}
		if def.TypeName != nil {
			column.Type = normalizeDataTypeText(e.nodeText(def.TypeName))
		} else if def.Domain != nil {
			column.Type = objectName(def.Domain)
		}
		if def.Default != nil {
			column.Default = e.nodeText(def.Default)
			column.DefaultOnNull = def.DefaultOnNull
		}
		switch {
		case def.NotNull:
			column.Nullable = false
		case def.Null:
			column.Nullable = true
			e.removeNamedNotNull(table, column.Name)
		case def.Identity != nil:
			column.Nullable = false
		}
		if def.DropIdentity {
			column.IsIdentity = false
			column.IdentityGeneration = metadatapb.ColumnMetadata_IDENTITY_GENERATION_UNSPECIFIED
			column.IdentitySeed = 0
			column.IdentityIncrement = 0
		} else if def.Identity != nil {
			applyIdentity(column, def.Identity, e.nodeText(def.Identity))
		}
		if def.Collation != "" {
			column.Collation = def.Collation
		}
		if def.Virtual != nil {
			column.Generation = &metadatapb.GenerationMetadata{
				Type:       metadatapb.GenerationMetadata_TYPE_VIRTUAL,
				Expression: e.exprText(def.Virtual),
			}
		}
		if def.Visible {
			column.IsInvisible = false
		} else if def.Invisible {
			column.IsInvisible = true
		}
		for _, item := range listItems(def.Constraints) {
			constraint, ok := item.(*ast.ColumnConstraint)
			if !ok || constraint.Type != ast.CONSTRAINT_NOT_NULL || constraint.Name == "" {
				continue
			}
			e.addNamedNotNull(table, constraint.Name, column.Name)
		}
		return
	}
}

func (e *oracleMetadataExtractor) dropColumn(table *metadatapb.TableMetadata, name string) {
	var columns []*metadatapb.ColumnMetadata
	for _, column := range table.Columns {
		if column.Name != name {
			columns = append(columns, column)
		}
	}
	table.Columns = columns

	var indexes []*metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if !slices.ContainsFunc(index.Expressions, func(expression string) bool { return expressionReferencesColumn(expression, name) }) {
			indexes = append(indexes, index)
		}
	}
	table.Indexes = indexes

	var checks []*metadatapb.CheckConstraintMetadata
	for _, check := range table.CheckConstraints {
		if !expressionReferencesColumn(check.Expression, name) {
			checks = append(checks, check)
		}
	}
	table.CheckConstraints = checks

	var foreignKeys []*metadatapb.ForeignKeyMetadata
	for _, fk := range table.ForeignKeys {
		onColumn := slices.Contains(fk.Columns, name)
		referencesColumn := fk.ReferencedTable == table.Name && slices.Contains(fk.ReferencedColumns, name)
		if !onColumn && !referencesColumn {
			foreignKeys = append(foreignKeys, fk)
		}
	}
	table.ForeignKeys = foreignKeys
	for _, other := range e.tables {
		if other == table {
			continue
		}
		var remaining []*metadatapb.ForeignKeyMetadata
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable == table.Name && slices.Contains(fk.ReferencedColumns, name) {
				continue
			}
			remaining = append(remaining, fk)
		}
		other.ForeignKeys = remaining
	}
}

func (e *oracleMetadataExtractor) renameColumn(table *metadatapb.TableMetadata, oldName, newName string) {
	if oldName == "" || newName == "" {
		return
	}
	for _, column := range table.Columns {
		if column.Name == oldName {
			column.Name = newName
			break
		}
	}
	for _, index := range table.Indexes {
		for i, expression := range index.Expressions {
			index.Expressions[i] = renameIdentifier(expression, oldName, newName)
		}
	}
	for _, column := range table.Columns {
		if column.GetGeneration() != nil {
			column.Generation.Expression = renameIdentifier(column.Generation.GetExpression(), oldName, newName)
		}
	}
	for _, check := range table.CheckConstraints {
		check.Expression = renameIdentifier(check.Expression, oldName, newName)
	}
	for _, fk := range table.ForeignKeys {
		fk.Columns = renameColumnNames(fk.Columns, oldName, newName)
	}
	for _, other := range e.tables {
		if other == table {
			continue
		}
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable == table.Name {
				fk.ReferencedColumns = renameColumnNames(fk.ReferencedColumns, oldName, newName)
			}
		}
	}
}

func renameColumnNames(columns []string, oldName, newName string) []string {
	out := make([]string, len(columns))
	for i, column := range columns {
		if column == oldName {
			out[i] = newName
		} else {
			out[i] = column
		}
	}
	return out
}

// renameIdentifier rewrites identifier tokens that name oldName and leaves
// string literals and other identifiers alone.
func renameIdentifier(expression, oldName, newName string) string {
	var out strings.Builder
	for i := 0; i < len(expression); {
		switch c := expression[i]; {
		case c == '\'':
			end := i + 1
			for end < len(expression) {
				if expression[end] == '\'' {
					if end+1 < len(expression) && expression[end+1] == '\'' {
						end += 2
						continue
					}
					end++
					break
				}
				end++
			}
			out.WriteString(expression[i:end])
			i = end
		case c == '"':
			end := i + 1
			for end < len(expression) && expression[end] != '"' {
				end++
			}
			if end < len(expression) {
				end++
			}
			identifier := expression[i:end]
			if len(identifier) >= 2 && identifier[1:len(identifier)-1] == oldName {
				out.WriteString(`"`)
				out.WriteString(newName)
				out.WriteString(`"`)
			} else {
				out.WriteString(identifier)
			}
			i = end
		case isIdentifierByte(c):
			end := i + 1
			for end < len(expression) && isIdentifierByte(expression[end]) {
				end++
			}
			token := expression[i:end]
			if strings.EqualFold(token, oldName) {
				out.WriteString(newName)
			} else {
				out.WriteString(token)
			}
			i = end
		default:
			out.WriteByte(c)
			i++
		}
	}
	return out.String()
}

func isIdentifierByte(c byte) bool {
	return c == '_' || c == '$' || c == '#' ||
		(c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func (e *oracleMetadataExtractor) dropConstraint(table *metadatapb.TableMetadata, name string) {
	if name == "" {
		return
	}
	var indexes []*metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if index.Name != name {
			indexes = append(indexes, index)
		}
	}
	table.Indexes = indexes
	var checks []*metadatapb.CheckConstraintMetadata
	for _, check := range table.CheckConstraints {
		if check.Name == name {
			e.clearNamedNotNull(table, check)
			continue
		}
		checks = append(checks, check)
	}
	table.CheckConstraints = checks
	var foreignKeys []*metadatapb.ForeignKeyMetadata
	for _, fk := range table.ForeignKeys {
		if fk.Name != name {
			foreignKeys = append(foreignKeys, fk)
		}
	}
	table.ForeignKeys = foreignKeys
}

func (e *oracleMetadataExtractor) renameConstraint(table *metadatapb.TableMetadata, oldName, newName string) {
	if oldName == "" || newName == "" {
		return
	}
	for _, index := range table.Indexes {
		if index.Name == oldName {
			index.Name = newName
		}
	}
	for _, check := range table.CheckConstraints {
		if check.Name == oldName {
			check.Name = newName
		}
	}
	for _, fk := range table.ForeignKeys {
		if fk.Name == oldName {
			fk.Name = newName
		}
	}
}

func (e *oracleMetadataExtractor) renameTable(oldName, newName string) {
	table := e.tables[oldName]
	if table == nil {
		return
	}
	delete(e.tables, oldName)
	table.Name = newName
	for _, fk := range table.ForeignKeys {
		if fk.ReferencedTable == oldName {
			fk.ReferencedTable = newName
		}
	}
	for _, other := range e.tables {
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable == oldName {
				fk.ReferencedTable = newName
			}
		}
	}
	e.tables[newName] = table
}

// extractDrop removes the objects a DROP statement names from the extractor.
func (e *oracleMetadataExtractor) extractDrop(n *ast.DropStmt) {
	for _, item := range listItems(n.Names) {
		nameNode, ok := item.(*ast.ObjectName)
		if !ok {
			continue
		}
		name := objectName(nameNode)
		switch n.ObjectType {
		case ast.OBJECT_TABLE:
			delete(e.tables, name)
		case ast.OBJECT_VIEW:
			delete(e.views, name)
		case ast.OBJECT_MATERIALIZED_VIEW:
			delete(e.materializedViews, name)
		case ast.OBJECT_SEQUENCE:
			delete(e.sequences, name)
		case ast.OBJECT_FUNCTION:
			delete(e.functions, name)
		case ast.OBJECT_PROCEDURE:
			delete(e.procedures, name)
		case ast.OBJECT_PACKAGE, ast.OBJECT_PACKAGE_BODY:
			delete(e.packages, name)
		case ast.OBJECT_INDEX:
			e.dropIndex(name)
		case ast.OBJECT_TRIGGER:
			e.dropTrigger(name)
		default:
		}
	}
}

func (e *oracleMetadataExtractor) dropIndex(name string) {
	for _, table := range e.tables {
		var indexes []*metadatapb.IndexMetadata
		for _, index := range table.Indexes {
			if index.Name != name {
				indexes = append(indexes, index)
			}
		}
		table.Indexes = indexes
	}
	for _, view := range e.materializedViews {
		var indexes []*metadatapb.IndexMetadata
		for _, index := range view.Indexes {
			if index.Name != name {
				indexes = append(indexes, index)
			}
		}
		view.Indexes = indexes
	}
}

func (e *oracleMetadataExtractor) dropTrigger(name string) {
	for _, table := range e.tables {
		var triggers []*metadatapb.TriggerMetadata
		for _, trigger := range table.Triggers {
			if trigger.Name != name {
				triggers = append(triggers, trigger)
			}
		}
		table.Triggers = triggers
	}
	for _, view := range e.views {
		var triggers []*metadatapb.TriggerMetadata
		for _, trigger := range view.Triggers {
			if trigger.Name != name {
				triggers = append(triggers, trigger)
			}
		}
		view.Triggers = triggers
	}
	for _, view := range e.materializedViews {
		var triggers []*metadatapb.TriggerMetadata
		for _, trigger := range view.Triggers {
			if trigger.Name != name {
				triggers = append(triggers, trigger)
			}
		}
		view.Triggers = triggers
	}
}

func (e *oracleMetadataExtractor) appendCheckConstraint(table *metadatapb.TableMetadata, name string, expr ast.ExprNode) {
	if expr == nil {
		return
	}
	e.reserveCheckConstraintName(table, name)
	table.CheckConstraints = append(table.CheckConstraints, &metadatapb.CheckConstraintMetadata{
		Name:       name,
		Expression: e.exprText(expr),
	})
}

func (e *oracleMetadataExtractor) checkConstraintName(table *metadatapb.TableMetadata, name string, fallback string) string {
	if name != "" {
		return name
	}
	return e.uniqueCheckConstraintName(table, fallback)
}

func (e *oracleMetadataExtractor) uniqueCheckConstraintName(table *metadatapb.TableMetadata, fallback string) string {
	used := e.checkConstraintNameSet(table)
	if !used[fallback] {
		return fallback
	}
	nextByFallback := e.checkConstraintNameNextSet(table)
	next := nextByFallback[fallback]
	if next < 2 {
		next = 2
	}
	for {
		candidate := fmt.Sprintf("%s_%d", fallback, next)
		next++
		if !used[candidate] {
			nextByFallback[fallback] = next
			return candidate
		}
	}
}

func (e *oracleMetadataExtractor) reserveCheckConstraintName(table *metadatapb.TableMetadata, name string) {
	e.checkConstraintNameSet(table)[name] = true
}

func (e *oracleMetadataExtractor) checkConstraintNameSet(table *metadatapb.TableMetadata) map[string]bool {
	used := e.checkNames[table.Name]
	if used != nil {
		return used
	}
	used = make(map[string]bool)
	for _, check := range table.CheckConstraints {
		used[check.Name] = true
	}
	e.checkNames[table.Name] = used
	return used
}

func (e *oracleMetadataExtractor) checkConstraintNameNextSet(table *metadatapb.TableMetadata) map[string]int {
	next := e.checkNameNext[table.Name]
	if next != nil {
		return next
	}
	next = make(map[string]int)
	e.checkNameNext[table.Name] = next
	return next
}

func (e *oracleMetadataExtractor) appendForeignKey(table *metadatapb.TableMetadata, name string, columns []string, refTable *ast.ObjectName, refColumns *ast.List, onDelete string) {
	referencedTable := objectName(refTable)
	referencedColumns := stringList(refColumns)
	if len(columns) == 0 || referencedTable == "" {
		return
	}
	referencedSchema := ""
	if refTable != nil && refTable.Schema != "" && refTable.Schema != e.currentSchema && refTable.Schema != e.currentDatabase {
		referencedSchema = refTable.Schema
	}
	foreignKey := &metadatapb.ForeignKeyMetadata{
		Name:              name,
		Columns:           columns,
		ReferencedSchema:  referencedSchema,
		ReferencedTable:   referencedTable,
		ReferencedColumns: referencedColumns,
		OnDelete:          onDelete,
	}
	table.ForeignKeys = append(table.ForeignKeys, foreignKey)
}

func (e *oracleMetadataExtractor) resolveForeignKeyReferencedColumns() {
	for _, table := range e.tables {
		for _, foreignKey := range table.ForeignKeys {
			if len(foreignKey.ReferencedColumns) != 0 {
				continue
			}
			// An owner outside this snapshot keeps its implicit reference: the
			// definition and migration render it without a column list.
			if foreignKey.ReferencedSchema != "" && foreignKey.ReferencedSchema != e.currentSchema && foreignKey.ReferencedSchema != e.currentDatabase {
				continue
			}
			referencedTable := e.tables[foreignKey.ReferencedTable]
			foreignKey.ReferencedColumns = primaryKeyColumns(referencedTable)
		}
	}
}

func (e *oracleMetadataExtractor) extractCreateIndex(n *ast.CreateIndexStmt) {
	tableName := objectName(n.Table)
	indexName := objectName(n.Name)
	if tableName == "" || indexName == "" {
		return
	}
	if !e.useSchema(n.Name.Schema) {
		return
	}

	index := &metadatapb.IndexMetadata{
		Name:        indexName,
		Unique:      n.Unique,
		Type:        "NORMAL",
		Expressions: []string{},
		Visible:     !n.Invisible,
	}
	if n.Bitmap {
		index.Type = "BITMAP"
	}

	isFunctionBased := false
	for _, item := range listItems(n.Columns) {
		column, ok := item.(*ast.IndexColumn)
		if !ok {
			continue
		}
		expression := e.indexExpression(column.Expr)
		if expression == "" {
			continue
		}
		if _, ok := column.Expr.(*ast.ColumnRef); !ok {
			isFunctionBased = true
		}
		index.Expressions = append(index.Expressions, expression)
		index.Descending = append(index.Descending, column.Dir == ast.SORTBY_DESC)
	}
	for _, desc := range index.Descending {
		if desc && index.Type == "NORMAL" {
			index.Type = "FUNCTION-BASED NORMAL"
			break
		}
	}
	if isFunctionBased && index.Type == "NORMAL" {
		index.Type = "FUNCTION-BASED NORMAL"
	} else if isFunctionBased && index.Type == "BITMAP" {
		index.Type = "FUNCTION-BASED BITMAP"
	}

	if materializedView := e.materializedViews[tableName]; materializedView != nil {
		materializedView.Indexes = append(materializedView.Indexes, index)
		return
	}
	table := e.getOrCreateTable(tableName)
	table.Indexes = append(table.Indexes, index)
}

func (e *oracleMetadataExtractor) extractCreateView(n *ast.CreateViewStmt) {
	viewName := objectName(n.Name)
	if viewName == "" {
		return
	}
	if !e.useSchema(n.Name.Schema) {
		return
	}

	definition := e.nodeText(n.Query)
	if n.Materialized {
		materializedView := &metadatapb.MaterializedViewMetadata{
			Name:              viewName,
			Definition:        definition,
			DependencyColumns: viewDependencyColumns(n.Query, e.currentSchema),
		}
		if table := e.tables[viewName]; table != nil {
			materializedView.Triggers = append(materializedView.Triggers, table.Triggers...)
			materializedView.Indexes = append(materializedView.Indexes, table.Indexes...)
		}
		if definition != "" && !strings.HasSuffix(definition, "\n") {
			materializedView.Definition += "\n"
		}
		e.materializedViews[viewName] = materializedView
		delete(e.tables, viewName)
		return
	}

	view := &metadatapb.ViewMetadata{
		Name:              viewName,
		Definition:        definition,
		DependencyColumns: viewDependencyColumns(n.Query, e.currentSchema),
	}
	if table := e.tables[viewName]; table != nil {
		view.Triggers = append(view.Triggers, table.Triggers...)
		delete(e.tables, viewName)
	}
	e.views[viewName] = view
}

func (e *oracleMetadataExtractor) extractCreateSequence(n *ast.CreateSequenceStmt) {
	sequenceName := objectName(n.Name)
	if sequenceName == "" {
		return
	}
	sequence := &metadatapb.SequenceMetadata{Name: sequenceName}
	if start := e.nodeText(n.StartWith); start != "" {
		sequence.Start = start
	}
	if increment := e.nodeText(n.IncrementBy); increment != "" && increment != "1" {
		sequence.Increment = increment
	}
	sequence.MinValue = e.nodeText(n.MinValue)
	sequence.MaxValue = e.nodeText(n.MaxValue)
	sequence.CacheSize = e.nodeText(n.Cache)
	if n.NoCache {
		sequence.CacheSize = "0"
	}
	sequence.Cycle = n.Cycle
	e.sequences[sequenceName] = sequence
}

func (e *oracleMetadataExtractor) extractAlterSequence(n *ast.AlterSequenceStmt) {
	name := objectName(n.Name)
	if name == "" {
		return
	}
	sequence := e.sequences[name]
	if sequence == nil {
		if n.IfExists {
			return
		}
		sequence = &metadatapb.SequenceMetadata{Name: name}
		e.sequences[name] = sequence
	}
	if n.IncrementBy != nil {
		sequence.Increment = e.nodeText(n.IncrementBy)
	}
	if n.MaxValue != nil {
		sequence.MaxValue = e.nodeText(n.MaxValue)
	} else if n.NoMaxValue {
		sequence.MaxValue = ""
	}
	if n.MinValue != nil {
		sequence.MinValue = e.nodeText(n.MinValue)
	} else if n.NoMinValue {
		sequence.MinValue = ""
	}
	if n.Cache != nil {
		sequence.CacheSize = e.nodeText(n.Cache)
	} else if n.NoCache {
		sequence.CacheSize = "0"
	}
	switch {
	case n.Cycle:
		sequence.Cycle = true
	case n.NoCycle:
		sequence.Cycle = false
	}
	if n.RestartWith != nil {
		sequence.Start = e.nodeText(n.RestartWith)
	}
}

func (e *oracleMetadataExtractor) extractCreatePackage(name string, n *ast.CreatePackageStmt) {
	definition := e.definitionText(n)
	if definition == "" {
		return
	}
	pkg := e.packages[name]
	if pkg == nil {
		e.packages[name] = &metadatapb.PackageMetadata{Name: name, Definition: definition}
		return
	}
	if pkg.Definition == "" {
		pkg.Definition = definition
		return
	}
	pkg.Definition += "\n\n" + definition
}

func (e *oracleMetadataExtractor) extractCreateTrigger(n *ast.CreateTriggerStmt) {
	triggerName := objectName(n.Name)
	tableName := objectName(n.Table)
	if triggerName == "" || tableName == "" {
		return
	}
	trigger := &metadatapb.TriggerMetadata{
		Name: triggerName,
		Body: e.definitionText(n),
	}
	e.triggers[triggerName] = trigger
	if view := e.views[tableName]; view != nil {
		view.Triggers = append(view.Triggers, trigger)
		return
	}
	if materializedView := e.materializedViews[tableName]; materializedView != nil {
		materializedView.Triggers = append(materializedView.Triggers, trigger)
		return
	}
	table := e.getOrCreateTable(tableName)
	table.Triggers = append(table.Triggers, trigger)
}

func (e *oracleMetadataExtractor) extractComment(n *ast.CommentStmt) {
	if n.Column != "" {
		table := e.tables[objectName(n.Object)]
		if table == nil {
			return
		}
		for _, column := range table.Columns {
			if column.Name == n.Column {
				column.Comment = n.Comment
				return
			}
		}
		return
	}

	switch n.ObjectType {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW:
		name := objectName(n.Object)
		if table := e.tables[name]; table != nil {
			table.Comment = n.Comment
			return
		}
		if view := e.views[name]; view != nil {
			view.Comment = n.Comment
		}
	case ast.OBJECT_MATERIALIZED_VIEW:
		if materializedView := e.materializedViews[objectName(n.Object)]; materializedView != nil {
			materializedView.Comment = n.Comment
		}
	default:
	}
}

func (e *oracleMetadataExtractor) indexExpression(expr ast.ExprNode) string {
	switch n := expr.(type) {
	case *ast.ColumnRef:
		if text := strings.TrimSpace(e.nodeText(expr)); strings.HasPrefix(text, `"`) && strings.HasSuffix(text, `"`) {
			return text
		}
		return n.Column
	default:
		return e.nodeText(expr)
	}
}

func (e *oracleMetadataExtractor) exprText(expr ast.ExprNode) string {
	if text := e.nodeText(expr); text != "" {
		return text
	}
	switch n := expr.(type) {
	case *ast.BinaryExpr:
		return e.exprText(n.Left) + n.Op + e.exprText(n.Right)
	case *ast.BoolExpr:
		var parts []string
		for _, item := range listItems(n.Args) {
			expr, ok := item.(ast.ExprNode)
			if !ok {
				continue
			}
			parts = append(parts, e.exprText(expr))
		}
		switch n.Boolop {
		case ast.BOOL_AND:
			return strings.Join(parts, "AND")
		case ast.BOOL_OR:
			return strings.Join(parts, "OR")
		case ast.BOOL_NOT:
			return "NOT" + strings.Join(parts, "")
		default:
			return strings.Join(parts, "")
		}
	case *ast.ColumnRef:
		if n.Table != "" {
			return n.Table + "." + n.Column
		}
		return n.Column
	case *ast.InExpr:
		operator := "IN"
		if n.Not {
			operator = "NOTIN"
		}
		return e.exprText(n.Expr) + operator + "(" + strings.Join(e.exprListText(n.List), ",") + ")"
	case *ast.BetweenExpr:
		operator := "BETWEEN"
		if n.Not {
			operator = "NOTBETWEEN"
		}
		return e.exprText(n.Expr) + operator + e.exprText(n.Low) + "AND" + e.exprText(n.High)
	case *ast.LikeExpr:
		operator := "LIKE"
		if n.Not {
			operator = "NOTLIKE"
		}
		result := e.exprText(n.Expr) + operator + e.exprText(n.Pattern)
		if n.Escape != nil {
			result += "ESCAPE" + e.exprText(n.Escape)
		}
		return result
	case *ast.IsExpr:
		operator := "IS"
		if n.Not {
			operator = "ISNOT"
		}
		return e.exprText(n.Expr) + operator + n.Test
	case *ast.ParenExpr:
		return "(" + e.exprText(n.Expr) + ")"
	case *ast.StringLiteral:
		prefix := ""
		if n.IsNChar {
			prefix = "N"
		}
		return prefix + "'" + strings.ReplaceAll(n.Val, "'", "''") + "'"
	case *ast.NumberLiteral:
		return n.Val
	case *ast.NullLiteral:
		return "NULL"
	case *ast.FuncCallExpr:
		return objectName(n.FuncName) + "(" + strings.Join(e.exprListText(n.Args), ",") + ")"
	default:
		return strings.ReplaceAll(e.nodeText(expr), " ", "")
	}
}

func (e *oracleMetadataExtractor) exprListText(list *ast.List) []string {
	var result []string
	for _, item := range listItems(list) {
		expr, ok := item.(ast.ExprNode)
		if !ok {
			continue
		}
		result = append(result, e.exprText(expr))
	}
	return result
}

func (e *oracleMetadataExtractor) functionDefinitionText(node ast.Node) string {
	definition := e.definitionText(node)
	upper := strings.ToUpper(definition)
	for _, prefix := range []string{"CREATE OR REPLACE", "CREATE"} {
		if strings.HasPrefix(upper, prefix) {
			definition = strings.TrimSpace(definition[len(prefix):])
			break
		}
	}
	if definition != "" && !strings.HasPrefix(strings.ToUpper(definition), "FUNCTION") {
		definition = "FUNCTION " + definition
	}
	return definition
}

func (e *oracleMetadataExtractor) definitionText(node ast.Node) string {
	loc := ast.NodeLoc(node)
	if loc.Start < 0 || loc.End < 0 || loc.Start > loc.End || loc.End > len(e.schemaText) {
		return ""
	}
	end := loc.End
	for end < len(e.schemaText) {
		switch e.schemaText[end] {
		case ' ', '\t', '\r', '\n':
			end++
		default:
			if e.schemaText[end] == ';' {
				end++
			}
			return strings.TrimSpace(e.schemaText[loc.Start:end])
		}
	}
	return strings.TrimSpace(e.schemaText[loc.Start:end])
}

func (e *oracleMetadataExtractor) nodeText(node ast.Node) string {
	loc := ast.NodeLoc(node)
	if loc.Start < 0 || loc.End < 0 || loc.Start > loc.End || loc.End > len(e.schemaText) {
		return ""
	}
	return e.schemaText[loc.Start:loc.End]
}

func objectName(name *ast.ObjectName) string {
	if name == nil {
		return ""
	}
	return name.Name
}

func listItems(list *ast.List) []ast.Node {
	if list == nil {
		return nil
	}
	return list.Items
}

func stringList(list *ast.List) []string {
	var result []string
	for _, item := range listItems(list) {
		switch n := item.(type) {
		case *ast.String:
			result = append(result, n.Str)
		case *ast.ColumnRef:
			result = append(result, n.Column)
		default:
		}
	}
	return result
}

func fallbackName(name, fallback string) string {
	if name != "" {
		return name
	}
	return fallback
}

func hasPrimaryIndex(table *metadatapb.TableMetadata) bool {
	for _, index := range table.Indexes {
		if index.Primary {
			return true
		}
	}
	return false
}

func primaryKeyColumns(table *metadatapb.TableMetadata) []string {
	if table == nil {
		return nil
	}
	for _, index := range table.Indexes {
		if index.Primary {
			return append([]string(nil), index.Expressions...)
		}
	}
	return nil
}

func hasUniqueIndex(table *metadatapb.TableMetadata, columns []string) bool {
	for _, index := range table.Indexes {
		if !index.Unique || index.Primary || len(index.Expressions) != len(columns) {
			continue
		}
		matches := true
		for i, column := range columns {
			if index.Expressions[i] != column {
				matches = false
				break
			}
		}
		if matches {
			return true
		}
	}
	return false
}

func markColumnsNotNull(table *metadatapb.TableMetadata, columns []string) {
	for _, column := range table.Columns {
		if slices.Contains(columns, column.Name) {
			column.Nullable = false
		}
	}
}

// newMetadataExtractorFromMetadata seeds an extractor with an existing snapshot
// and the text of the statement about to be applied.
func newMetadataExtractorFromMetadata(meta *metadatapb.DatabaseSchemaMetadata, schemaText string) *oracleMetadataExtractor {
	e := newMetadataExtractor(schemaText)
	e.currentDatabase = meta.GetName()
	if len(meta.GetSchemas()) == 0 {
		return e
	}
	schema := meta.GetSchemas()[0]
	e.currentSchema = schema.GetName()
	for _, table := range schema.GetTables() {
		e.tables[table.GetName()] = cloneTable(table)
	}
	for _, view := range schema.GetViews() {
		e.views[view.GetName()] = cloneView(view)
	}
	for _, view := range schema.GetMaterializedViews() {
		e.materializedViews[view.GetName()] = cloneMaterializedView(view)
	}
	for _, fn := range schema.GetFunctions() {
		e.functions[fn.GetName()] = cloneFunction(fn)
	}
	for _, proc := range schema.GetProcedures() {
		e.procedures[proc.GetName()] = cloneProcedure(proc)
	}
	for _, seq := range schema.GetSequences() {
		e.sequences[seq.GetName()] = cloneSequence(seq)
	}
	for _, pkg := range schema.GetPackages() {
		e.packages[pkg.GetName()] = clonePackage(pkg)
	}
	return e
}

var (
	identityStartPattern     = regexp.MustCompile(`(?i)START\s+WITH\s+([-+]?[0-9]+)`)
	identityIncrementPattern = regexp.MustCompile(`(?i)INCREMENT\s+BY\s+([-+]?[0-9]+)`)
)

func applyIdentity(column *metadatapb.ColumnMetadata, identity *ast.IdentityClause, text string) {
	column.IsIdentity = true
	if identity != nil && !identity.Always {
		column.IdentityGeneration = metadatapb.ColumnMetadata_BY_DEFAULT
	} else {
		column.IdentityGeneration = metadatapb.ColumnMetadata_ALWAYS
	}
	if match := identityStartPattern.FindStringSubmatch(text); len(match) == 2 {
		if value, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			column.IdentitySeed = value
		}
	}
	if match := identityIncrementPattern.FindStringSubmatch(text); len(match) == 2 {
		if value, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			column.IdentityIncrement = value
		}
	}
}

type viewDependencyVisitor struct {
	schema string
	seen   map[string]bool
	result *[]*metadatapb.DependencyColumn
}

func (v *viewDependencyVisitor) Visit(node ast.Node) ast.Visitor {
	if ref, ok := node.(*ast.TableRef); ok && ref.Name != nil && ref.Name.Name != "" {
		schema := ref.Name.Schema
		if schema == "" {
			schema = v.schema
		}
		key := schema + "\x00" + ref.Name.Name
		if !v.seen[key] {
			v.seen[key] = true
			*v.result = append(*v.result, &metadatapb.DependencyColumn{Schema: schema, Table: ref.Name.Name, Column: "*"})
		}
	}
	return v
}

func viewDependencyColumns(query ast.StmtNode, schema string) []*metadatapb.DependencyColumn {
	if query == nil {
		return nil
	}
	var result []*metadatapb.DependencyColumn
	ast.Walk(&viewDependencyVisitor{schema: schema, seen: map[string]bool{}, result: &result}, query)
	return result
}
