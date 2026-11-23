package catalog

import (
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle"
	"github.com/bytebase/omni/oracle/analysis"
	"github.com/bytebase/omni/oracle/ast"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

// GetDatabaseMetadata parses Oracle DDL and returns the schema it defines.
// Unlike LoadSQL, it skips SQL*Plus commands, allows constraints backed by an
// existing index (USING INDEX), and keeps materialized view definitions as the
// query text synced snapshots store.
func GetDatabaseMetadata(schemaText string) (*metadatapb.DatabaseSchemaMetadata, error) {
	list, err := analysis.Parse(schemaText)
	if err != nil {
		return nil, fmt.Errorf("parse Oracle schema: %w", err)
	}
	if list == nil || len(list.Items) == 0 {
		return nil, fmt.Errorf("no parse results")
	}

	extractor := newMetadataExtractor(schemaText)
	extractor.allowUsingExistingIndex = true
	extractor.omitMaterializedViewClauses = true
	for _, item := range list.Items {
		raw, ok := item.(*ast.RawStmt)
		if !ok || raw.Stmt == nil {
			continue
		}
		node, ok := raw.Stmt.(ast.StmtNode)
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
		namedNotNull:         make(map[string]bool),
		primaryKeyNotNull:    make(map[string]bool),
		syntheticConstraints: make(map[string]bool),
		viewConstraints:      make(map[string][]viewConstraint),
		editioningViews:      make(map[string]bool),
		encryptedColumns:     make(map[string]bool),
		domains:              make(map[string]bool),
		editions:             map[string]bool{"ORA$BASE": true},
		sequenceCurrent:      make(map[string]*big.Int),
		tables:               make(map[string]*metadatapb.TableMetadata),
		views:                make(map[string]*metadatapb.ViewMetadata),
		materializedViews:    make(map[string]*metadatapb.MaterializedViewMetadata),
		functions:            make(map[string]*metadatapb.FunctionMetadata),
		procedures:           make(map[string]*metadatapb.ProcedureMetadata),
		triggers:             make(map[string]*metadatapb.TriggerMetadata),
		sequences:            make(map[string]*metadatapb.SequenceMetadata),
		packages:             make(map[string]*metadatapb.PackageMetadata),
		synonyms:             make(map[string]*metadatapb.SynonymMetadata),
		publicSynonyms:       make(map[string]*metadatapb.SynonymMetadata),
		checkNames:           make(map[string]map[string]bool),
		checkNameNext:        make(map[string]map[string]int),
		schemaText:           schemaText,
	}
}

type oracleMetadataExtractor struct {
	err                         error
	serverConfig                ServerConfig
	allowUsingExistingIndex     bool
	omitMaterializedViewClauses bool
	namedNotNull                map[string]bool
	primaryKeyNotNull           map[string]bool
	syntheticConstraints        map[string]bool
	viewConstraints             map[string][]viewConstraint
	editioningViews             map[string]bool
	encryptedColumns            map[string]bool
	domains                     map[string]bool
	editions                    map[string]bool
	sequenceCurrent             map[string]*big.Int
	explicitSchema              bool
	owner                       string
	currentDatabase             string
	currentSchema               string
	tables                      map[string]*metadatapb.TableMetadata
	views                       map[string]*metadatapb.ViewMetadata
	materializedViews           map[string]*metadatapb.MaterializedViewMetadata
	functions                   map[string]*metadatapb.FunctionMetadata
	procedures                  map[string]*metadatapb.ProcedureMetadata
	triggers                    map[string]*metadatapb.TriggerMetadata
	sequences                   map[string]*metadatapb.SequenceMetadata
	packages                    map[string]*metadatapb.PackageMetadata
	synonyms                    map[string]*metadatapb.SynonymMetadata
	publicSynonyms              map[string]*metadatapb.SynonymMetadata
	checkNames                  map[string]map[string]bool
	checkNameNext               map[string]map[string]int
	schemaText                  string
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

	if len(e.synonyms) > 0 {
		var synonymNames []string
		for name := range e.synonyms {
			synonymNames = append(synonymNames, name)
		}
		slices.Sort(synonymNames)
		for _, name := range synonymNames {
			schema.Synonyms = append(schema.Synonyms, e.synonyms[name])
		}
	}

	if len(e.publicSynonyms) > 0 {
		var publicSynonymNames []string
		for name := range e.publicSynonyms {
			publicSynonymNames = append(publicSynonymNames, name)
		}
		slices.Sort(publicSynonymNames)
		for _, name := range publicSynonymNames {
			schemaMetadata.PublicSynonyms = append(schemaMetadata.PublicSynonyms, e.publicSynonyms[name])
		}
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
		if e.owner != "" && schema != e.owner {
			if e.err == nil {
				e.err = fmt.Errorf("statement names schema %q, but the catalog is loading schema %q", schema, e.owner)
			}
			return false
		}
		e.currentSchema = schema
		e.explicitSchema = true
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
		if name := objectName(n.Name); name != "" && e.useSchema(n.Name.Schema) {
			e.procedures[name] = &metadatapb.ProcedureMetadata{Name: name, Definition: e.definitionText(n)}
		}
	case *ast.CreateFunctionStmt:
		if name := objectName(n.Name); name != "" && e.useSchema(n.Name.Schema) {
			e.functions[name] = &metadatapb.FunctionMetadata{Name: name, Definition: e.functionDefinitionText(n)}
		}
	case *ast.CreatePackageStmt:
		if name := objectName(n.Name); name != "" && e.useSchema(n.Name.Schema) {
			e.extractCreatePackage(name, n)
		}
	case *ast.CreateSynonymStmt:
		e.extractCreateSynonym(n)
	case *ast.CreateDomainStmt:
		e.extractCreateDomain(n)
	case *ast.CreateTriggerStmt:
		e.extractCreateTrigger(n)
	case *ast.AlterTableStmt:
		e.extractAlterTable(n)
	case *ast.AlterViewStmt:
		e.extractAlterView(n)
	case *ast.TruncateStmt:
		e.extractTruncate(n)
	case *ast.CommentStmt:
		e.extractComment(n)
	case *ast.AlterSequenceStmt:
		e.extractAlterSequence(n)
	case *ast.AlterTriggerStmt:
		e.extractAlterTrigger(n)
	case *ast.RenameStmt:
		e.extractRename(n)
	case *ast.AlterIndexStmt:
		e.extractAlterIndex(n)
	case *ast.DropStmt:
		e.extractDrop(n)
	case *ast.AdminDDLStmt:
		e.extractAdminDDL(n)
	case *ast.AlterSessionStmt:
		e.extractAlterSession(n)
	default:
	}
}

func (e *oracleMetadataExtractor) extractCreateSchema(n *ast.CreateSchemaStmt) {
	if n.SchemaName != "" && !e.useSchema(n.SchemaName) {
		return
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
	if e.serverConfig.NonApplicationRoot && isApplicationSharing(n.Sharing) {
		e.setError(oraError("65021", "illegal use of SHARING clause"))
		return
	}
	if e.isStandardMaxStringSize() {
		if n.Collation != "" {
			e.setError(oraError("43929", "Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set"))
			return
		}
		for _, item := range listItems(n.Columns) {
			if col, ok := item.(*ast.ColumnDef); ok && col.Collation != "" {
				e.setError(oraError("43929", "Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set"))
				return
			}
		}
	}
	if !e.useSchema(n.Name.Schema) {
		return
	}

	if n.AsQuery != nil && queryHasUnaliasedExpression(n.AsQuery) {
		e.setError(oraError("00998", "must name this expression with a column alias"))
		return
	}

	// Validate column definitions before checking table existence, matching Oracle's order:
	// duplicate column names first, then missing datatypes / domains, then all-virtual tables.
	seenCols := map[string]bool{}
	var colDefs []*ast.ColumnDef
	for _, item := range listItems(n.Columns) {
		column, ok := item.(*ast.ColumnDef)
		if !ok {
			continue
		}
		if seenCols[column.Name] {
			e.setError(oraError("00957", "duplicate column name %q", column.Name))
			return
		}
		seenCols[column.Name] = true
		colDefs = append(colDefs, column)
	}
	realColCount := 0
	for _, column := range colDefs {
		if column.Domain != nil {
			if e.isForeignSchema(column.Domain.Schema) {
				e.setError(oraError("01917", "user or role %q does not exist", column.Domain.Schema))
				return
			}
			domName := objectName(column.Domain)
			if !e.domains[domName] {
				e.setError(oraError("11504", "domain %q does not exist", domName))
				return
			}
		}
		if column.TypeName == nil && column.Domain == nil && column.Virtual == nil && !colHasInlineReferences(column) {
			e.setError(oraError("02263", "need to specify the datatype for this column"))
			return
		}
		if column.Virtual == nil {
			realColCount++
		}
	}
	if len(colDefs) > 0 && realColCount == 0 {
		e.setError(oraError("54037", "table must have at least 1 column that is not virtual"))
		return
	}
	if len(colDefs) == 0 && e.hasPeriodWithExplicitCols(n) {
		e.setError(oraError("54037", "table must have at least 1 column that is not virtual"))
		return
	}

	if n.IfNotExists && e.tables[tableName] != nil {
		return
	}
	if n.AsQuery != nil {
		if e.err == nil {
			e.err = fmt.Errorf("CREATE TABLE AS SELECT is not supported by oracle/catalog")
		}
		return
	}
	if n.Global || n.Private {
		if e.err == nil {
			e.err = fmt.Errorf("temporary tables are not supported by oracle/catalog")
		}
		return
	}
	if n.OfType != nil {
		if e.err == nil {
			e.err = fmt.Errorf("object and XMLType tables are not supported by oracle/catalog")
		}
		return
	}
	if !n.IfNotExists && e.objectExistsInNamespace(tableName) {
		e.setError(oraError("00955", "name %q is already used by an existing object", tableName))
		return
	}
	for _, column := range colDefs {
		if !e.validateVirtualColumnEditions(column) {
			return
		}
		if isColumnEncrypt(column.Encrypt) {
			if colHasInlineReferences(column) {
				e.setError(oraError("28335", "referenced or referencing foreign key constraint column %q cannot be encrypted", column.Name))
				return
			}
			if e.serverConfig.KeystoreClosed {
				e.setError(oraError("28365", "wallet is not open"))
				return
			}
		}
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
	for _, item := range listItems(n.Options) {
		opt, ok := item.(*ast.DDLOption)
		if !ok || !strings.EqualFold(opt.Key, "PERIOD FOR") {
			continue
		}
		if period, startCol, endCol, ok := parsePeriodForClause(opt.Value, false); ok {
			e.addPeriodConstraint(table, period, startCol, endCol)
		}
	}
	// Validate inline REFERENCES constraints: untyped columns must resolve
	// against the target table, and external schema references are refused
	// when ValidateExternalSchemaRefs is enabled.
	for _, item := range listItems(n.Columns) {
		column, ok := item.(*ast.ColumnDef)
		if !ok || column.Virtual != nil || column.Domain != nil {
			continue
		}
		untyped := column.TypeName == nil
		for _, citem := range listItems(column.Constraints) {
			cc, ok := citem.(*ast.ColumnConstraint)
			if !ok || cc.Type != ast.CONSTRAINT_FOREIGN || cc.RefTable == nil {
				continue
			}
			if e.serverConfig.ValidateExternalSchemaRefs && e.isUnresolvedExternalSchema(cc.RefTable.Schema) {
				e.setError(oraError("00942", "table or view %q.%q does not exist", cc.RefTable.Schema, objectName(cc.RefTable)))
				e.dropTableProvenance(tableName)
				delete(e.tables, tableName)
				return
			}
			if e.isForeignSchema(cc.RefTable.Schema) {
				continue
			}
			if !untyped {
				continue
			}
			refName := objectName(cc.RefTable)
			target := e.tables[refName]
			if target == nil {
				e.setError(oraError("00942", "table or view %q does not exist", refName))
				e.dropTableProvenance(tableName)
				delete(e.tables, tableName)
				return
			}
			refCols := stringList(cc.RefColumns)
			if len(refCols) == 0 {
				if !hasPrimaryIndex(target) {
					e.setError(oraError("02268", "referenced table %q does not have a primary key", refName))
					e.dropTableProvenance(tableName)
					delete(e.tables, tableName)
					return
				}
			} else if !tableHasUniqueKey(target, refCols) {
				e.setError(oraError("02270", "no matching unique or primary key for this column-list on %q", refName))
				e.dropTableProvenance(tableName)
				delete(e.tables, tableName)
				return
			}
		}
	}
	if e.serverConfig.ValidateExternalSchemaRefs {
		for _, item := range listItems(n.Constraints) {
			tc, ok := item.(*ast.TableConstraint)
			if !ok || tc.Type != ast.CONSTRAINT_FOREIGN || tc.RefTable == nil {
				continue
			}
			if e.isUnresolvedExternalSchema(tc.RefTable.Schema) {
				e.setError(oraError("00942", "table or view %q.%q does not exist", tc.RefTable.Schema, objectName(tc.RefTable)))
				e.dropTableProvenance(tableName)
				delete(e.tables, tableName)
				return
			}
		}
	}
	if n.Partition != nil {
		if partitions := e.extractTablePartitions(n.Partition); len(partitions) > 0 {
			table.Partitions = partitions
		}
	}
}

func (e *oracleMetadataExtractor) hasPeriodWithExplicitCols(n *ast.CreateTableStmt) bool {
	text := strings.ToUpper(e.nodeText(n))
	idx := strings.Index(text, "PERIOD FOR ")
	if idx < 0 {
		return false
	}
	rest := text[idx+len("PERIOD FOR "):]
	return strings.Contains(rest, "(")
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
	} else if n.Virtual != nil {
		column.Type = inferVirtualColumnType(n.Virtual)
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
	if n.DefaultOnNull {
		column.Nullable = false
	}
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
	if n.Encrypt != "" && !strings.EqualFold(n.Encrypt, "DECRYPT") {
		e.encryptedColumns[table.Name+"\x00"+column.Name] = true
	}

	explicitNotNull := n.NotNull
	for _, item := range listItems(n.Constraints) {
		constraint, ok := item.(*ast.ColumnConstraint)
		if !ok {
			continue
		}
		if constraint.Type == ast.CONSTRAINT_NOT_NULL {
			explicitNotNull = true
		}
		e.extractColumnConstraint(constraint, table, column)
		if constraint.Type == ast.CONSTRAINT_NOT_NULL && constraint.Name != "" {
			e.addNamedNotNull(table, constraint.Name, column.Name)
		}
	}
	if explicitNotNull {
		delete(e.primaryKeyNotNull, table.Name+"\x00"+column.Name)
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
	e.namedNotNull[table.Name+"\x00"+name] = true
}

func (e *oracleMetadataExtractor) removeNamedNotNull(table *metadatapb.TableMetadata, column string) {
	want := `"` + column + `" IS NOT NULL`
	var checks []*metadatapb.CheckConstraintMetadata
	for _, check := range table.CheckConstraints {
		namedNotNull := e.namedNotNull[table.Name+"\x00"+check.Name]
		if !namedNotNull || !strings.EqualFold(strings.Join(strings.Fields(check.Expression), " "), want) {
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
	if e.constraintDisabled(n.Loc) {
		return
	}
	if !e.allowUsingExistingIndex && e.constraintUsesExistingIndex(n.UsingIndex) {
		return
	}
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
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable constraints are not supported by oracle/catalog")
			}
			return
		}
		pkMadeNotNull := column.Nullable
		column.Nullable = false
		if pkMadeNotNull {
			e.primaryKeyNotNull[table.Name+"\x00"+column.Name] = true
		}
		constraintName := e.resolveConstraintName(table, n.Name, fmt.Sprintf("PK_%s", table.Name))
		if e.allowUsingExistingIndex && n.UsingIndex != "" {
			if !e.attachConstraintToExistingIndex(table, n.UsingIndex, constraintName, []string{column.Name}, true) {
				_ = e.constraintUsesExistingIndex(n.UsingIndex)
			}
			return
		}
		if !hasPrimaryIndex(table) {
			table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
				Name:         constraintName,
				Primary:      true,
				Unique:       true,
				Type:         "NORMAL",
				Expressions:  []string{column.Name},
				Visible:      true,
				IsConstraint: true,
			})
		}
	case ast.CONSTRAINT_UNIQUE:
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable constraints are not supported by oracle/catalog")
			}
			return
		}
		name := e.resolveConstraintName(table, n.Name, fmt.Sprintf("UK_%s_%s", table.Name, column.Name))
		if e.allowUsingExistingIndex && n.UsingIndex != "" {
			if !e.attachConstraintToExistingIndex(table, n.UsingIndex, name, []string{column.Name}, false) {
				_ = e.constraintUsesExistingIndex(n.UsingIndex)
			}
			return
		}
		if reusable := findStandaloneUniqueIndex(table, []string{column.Name}); reusable != nil {
			e.attachConstraintToExistingIndex(table, reusable.Name, name, []string{column.Name}, false)
		} else if !hasUniqueIndex(table, []string{column.Name}) {
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
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable constraints are not supported by oracle/catalog")
			}
			return
		}
		if n.Name != "" {
			e.disambiguateSyntheticConstraint(table, n.Name)
		}
		checkName := e.checkConstraintName(table, n.Name, fmt.Sprintf("CHK_%s_%s", table.Name, column.Name))
		if n.Name == "" {
			e.syntheticConstraints[table.Name+"\x00"+checkName] = true
		}
		e.appendCheckConstraint(table, checkName, n.Expr)
	case ast.CONSTRAINT_FOREIGN:
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable foreign keys are not supported by oracle/catalog")
			}
			return
		}
		fkName := e.resolveConstraintName(table, n.Name, fmt.Sprintf("FK_%s_%s", table.Name, column.Name))
		e.appendForeignKey(table, fkName, []string{column.Name}, n.RefTable, n.RefColumns, n.OnDelete)
	default:
	}
}

func (e *oracleMetadataExtractor) extractTableConstraint(n *ast.TableConstraint, table *metadatapb.TableMetadata) {
	if e.constraintDisabled(n.Loc) {
		return
	}
	if !e.allowUsingExistingIndex && e.constraintUsesExistingIndex(n.UsingIndex) {
		return
	}
	columns := stringList(n.Columns)
	switch n.Type {
	case ast.CONSTRAINT_PRIMARY:
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable constraints are not supported by oracle/catalog")
			}
			return
		}
		if len(columns) == 0 {
			return
		}
		for _, columnName := range columns {
			for _, column := range table.Columns {
				if column.Name == columnName && column.Nullable {
					e.primaryKeyNotNull[table.Name+"\x00"+columnName] = true
				}
			}
		}
		constraintName := e.resolveConstraintName(table, n.Name, fmt.Sprintf("PK_%s", table.Name))
		if e.allowUsingExistingIndex && n.UsingIndex != "" {
			if e.attachConstraintToExistingIndex(table, n.UsingIndex, constraintName, columns, true) {
				markColumnsNotNull(table, columns)
			} else {
				_ = e.constraintUsesExistingIndex(n.UsingIndex)
			}
			return
		}
		table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
			Name:         constraintName,
			Primary:      true,
			Unique:       true,
			Type:         "NORMAL",
			Expressions:  columns,
			Visible:      true,
			IsConstraint: true,
		})
		markColumnsNotNull(table, columns)
	case ast.CONSTRAINT_UNIQUE:
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable constraints are not supported by oracle/catalog")
			}
			return
		}
		if len(columns) == 0 {
			return
		}
		constraintName := e.resolveConstraintName(table, n.Name, fmt.Sprintf("UK_%s_%d", table.Name, len(table.Indexes)+1))
		if e.allowUsingExistingIndex && n.UsingIndex != "" {
			if !e.attachConstraintToExistingIndex(table, n.UsingIndex, constraintName, columns, false) {
				_ = e.constraintUsesExistingIndex(n.UsingIndex)
			}
			return
		}
		table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
			Name:         constraintName,
			Unique:       true,
			Type:         "NORMAL",
			Expressions:  columns,
			Visible:      false,
			IsConstraint: true,
		})
	case ast.CONSTRAINT_CHECK:
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable constraints are not supported by oracle/catalog")
			}
			return
		}
		if n.Name != "" {
			e.disambiguateSyntheticConstraint(table, n.Name)
		}
		checkName := e.checkConstraintName(table, n.Name, fmt.Sprintf("CHK_%s_%d", table.Name, len(table.CheckConstraints)+1))
		if n.Name == "" {
			e.syntheticConstraints[table.Name+"\x00"+checkName] = true
		}
		e.appendCheckConstraint(table, checkName, n.Expr)
	case ast.CONSTRAINT_FOREIGN:
		if n.Deferrable || n.Initially != "" {
			if e.err == nil {
				e.err = fmt.Errorf("deferrable foreign keys are not supported by oracle/catalog")
			}
			return
		}
		fkName := e.resolveConstraintName(table, n.Name, fmt.Sprintf("FK_%s_%d", table.Name, len(table.ForeignKeys)+1))
		e.appendForeignKey(table, fkName, columns, n.RefTable, n.RefColumns, n.OnDelete)
	default:
	}
}

func (e *oracleMetadataExtractor) extractRename(n *ast.RenameStmt) {
	oldName := objectName(n.OldName)
	newName := objectName(n.NewName)
	if oldName == "" || newName == "" {
		return
	}
	if !e.useSchema(n.OldName.Schema) || !e.useSchema(n.NewName.Schema) {
		return
	}
	if e.sequences[oldName] == nil && e.views[oldName] == nil && e.materializedViews[oldName] == nil && e.tables[oldName] == nil && e.synonyms[oldName] == nil {
		e.setError(oraError("04043", "object %q does not exist", oldName))
		return
	}
	if e.objectExistsInNamespace(newName) {
		e.setError(oraError("00955", "name %q is already used by an existing object", newName))
		return
	}
	if e.renameSequence(oldName, newName) {
		return
	}
	if e.renameView(oldName, newName) {
		return
	}
	if e.renameMaterializedView(oldName, newName) {
		return
	}
	if e.renameSynonym(oldName, newName) {
		return
	}
	if e.tables[oldName] != nil {
		e.renameTable(oldName, newName)
		return
	}
}

func (e *oracleMetadataExtractor) renameSequence(oldName, newName string) bool {
	sequence := e.sequences[oldName]
	if sequence == nil {
		return false
	}
	delete(e.sequences, oldName)
	if curr, ok := e.sequenceCurrent[oldName]; ok {
		delete(e.sequenceCurrent, oldName)
		e.sequenceCurrent[newName] = curr
	}
	sequence.Name = newName
	e.sequences[newName] = sequence
	return true
}

func (e *oracleMetadataExtractor) renameView(oldName, newName string) bool {
	view := e.views[oldName]
	if view == nil {
		return false
	}
	delete(e.views, oldName)
	if vcs, ok := e.viewConstraints[oldName]; ok {
		delete(e.viewConstraints, oldName)
		e.viewConstraints[newName] = vcs
	}
	if e.editioningViews[oldName] {
		delete(e.editioningViews, oldName)
		e.editioningViews[newName] = true
	}
	view.Name = newName
	for _, trigger := range view.Triggers {
		trigger.Body = renameTriggerTarget(trigger.Body, oldName, newName)
	}
	e.views[newName] = view
	return true
}

func (e *oracleMetadataExtractor) renameMaterializedView(oldName, newName string) bool {
	view := e.materializedViews[oldName]
	if view == nil {
		return false
	}
	delete(e.materializedViews, oldName)
	view.Name = newName
	for _, trigger := range view.Triggers {
		trigger.Body = renameTriggerTarget(trigger.Body, oldName, newName)
	}
	e.materializedViews[newName] = view
	return true
}

func (e *oracleMetadataExtractor) extractAlterTable(n *ast.AlterTableStmt) {
	tableName := objectName(n.Name)
	if tableName == "" {
		return
	}
	if e.isStandardMaxStringSize() && alterTableHasCollationClause(n.Actions) {
		e.setError(oraError("43929", "Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set"))
		return
	}
	// Check intra-statement duplicate column lists before table lookup, matching Oracle.
	for _, item := range listItems(n.Actions) {
		action, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch action.Action {
		case ast.AT_ADD_COLUMN, ast.AT_MODIFY_COLUMN:
			seen := map[string]bool{}
			if action.ColumnDef != nil && action.ColumnDef.Name != "" {
				seen[action.ColumnDef.Name] = true
			}
			for _, citem := range listItems(action.ColumnDefs) {
				col, ok := citem.(*ast.ColumnDef)
				if !ok || col.Name == "" {
					continue
				}
				if seen[col.Name] {
					e.setError(oraError("00957", "duplicate column name %q", col.Name))
					return
				}
				seen[col.Name] = true
			}
		case ast.AT_DROP_CONSTRAINT:
			if action.Constraint != nil && action.Constraint.Name == "" && action.Subtype == "UNIQUE" {
				seen := map[string]bool{}
				for _, col := range stringList(action.Constraint.Columns) {
					if seen[col] {
						e.setError(oraError("00957", "duplicate column name %q", col))
						return
					}
					seen[col] = true
				}
			}
		}
	}
	if !e.checkSchema(n.Name.Schema, "00942") {
		return
	}
	targetName := tableName
	table := e.tables[targetName]
	if table == nil && e.materializedViews[targetName] == nil && alterActionsFollowSynonym(n.Actions) {
		resolved := e.resolveSynonymTarget(tableName, n.Name.Schema)
		if resolved != tableName {
			targetName = resolved
			table = e.tables[targetName]
		}
	}
	if table == nil {
		if e.materializedViews[targetName] != nil {
			for _, item := range listItems(n.Actions) {
				action, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				switch action.Action {
				case ast.AT_RENAME:
					e.setError(oraError("32318", "cannot rename a materialized view %q", tableName))
					return
				case ast.AT_ENABLE_DISABLE, ast.AT_MODIFY_CONSTRAINT:
					e.setError(fmt.Errorf("constraint enable/disable is not supported by oracle/catalog"))
					return
				case ast.AT_ENABLE_DISABLE_TRIGGERS:
					e.setError(fmt.Errorf("table-wide trigger enable/disable is not supported by oracle/catalog"))
					return
				case ast.AT_SET_UNUSED:
					e.setError(fmt.Errorf("SET UNUSED is not supported by oracle/catalog"))
					return
				}
			}
			return
		}
		e.setError(oraError("00942", "table or view %q does not exist", tableName))
		return
	}
	type pendingAlterFK struct {
		refTable   *ast.ObjectName
		refColumns *ast.List
	}
	var pendingFKs []pendingAlterFK
	flushPendingFKs := func() bool {
		for _, fk := range pendingFKs {
			if !e.validateAlterForeignKey(table, fk.refTable, fk.refColumns) {
				return false
			}
		}
		pendingFKs = pendingFKs[:0]
		return true
	}
	for _, item := range listItems(n.Actions) {
		action, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch action.Action {
		case ast.AT_ADD_COLUMN:
			var cols []*ast.ColumnDef
			if action.ColumnDef != nil {
				cols = append(cols, action.ColumnDef)
			}
			for _, citem := range listItems(action.ColumnDefs) {
				if col, ok := citem.(*ast.ColumnDef); ok {
					cols = append(cols, col)
				}
			}
			for _, col := range cols {
				if tableFindColumn(table, col.Name) != nil {
					e.setError(oraError("01430", "column %q being added already exists in table %q", col.Name, table.Name))
					return
				}
				if !e.validateVirtualColumnEditions(col) {
					return
				}
				if isColumnEncrypt(col.Encrypt) {
					if colHasInlineReferences(col) {
						e.setError(oraError("28335", "referenced or referencing foreign key constraint column %q cannot be encrypted", col.Name))
						return
					}
					if e.serverConfig.KeystoreClosed {
						e.setError(oraError("28365", "wallet is not open"))
						return
					}
				}
				for _, citem := range listItems(col.Constraints) {
					cc, ok := citem.(*ast.ColumnConstraint)
					if !ok {
						continue
					}
					if cc.Name != "" && e.constraintNameExists(cc.Name) {
						e.setError(oraError("02264", "name %q already used by an existing constraint", cc.Name))
						return
					}
					switch cc.Type {
					case ast.CONSTRAINT_SCOPE, ast.CONSTRAINT_WITH_ROWID:
						if col.TypeName == nil || !col.TypeName.IsRef {
							e.setError(oraError("22893", "constraint may only be specified on a REF column"))
							return
						}
					case ast.CONSTRAINT_PRIMARY:
						if hasPrimaryIndex(table) {
							e.setError(oraError("02260", "table %q can have only one primary key", table.Name))
							return
						}
					case ast.CONSTRAINT_FOREIGN:
						pendingFKs = append(pendingFKs, pendingAlterFK{refTable: cc.RefTable, refColumns: cc.RefColumns})
					}
				}
				e.extractColumn(col, table)
			}
		case ast.AT_ADD_PERIOD:
			if action.ColumnName != "" && tableFindColumn(table, action.ColumnName) != nil {
				e.setError(oraError("01430", "column %q being added already exists in table %q", action.ColumnName, table.Name))
				return
			}
			if period, startCol, endCol, ok := parsePeriodForClause(e.nodeText(action), true); ok {
				e.addPeriodConstraint(table, period, startCol, endCol)
			} else if action.ColumnName != "" {
				e.addPeriodConstraint(table, action.ColumnName, "", "")
			}
		case ast.AT_DROP_PERIOD:
			if action.ColumnName != "" {
				e.dropPeriodConstraint(table, action.ColumnName)
			}
		case ast.AT_MODIFY_COLUMN:
			if action.ColumnDef != nil {
				for _, fk := range e.modifyColumn(action.ColumnDef, table) {
					pendingFKs = append(pendingFKs, pendingAlterFK{refTable: fk.RefTable, refColumns: fk.RefColumns})
				}
			}
			for _, citem := range listItems(action.ColumnDefs) {
				column, ok := citem.(*ast.ColumnDef)
				if !ok {
					continue
				}
				for _, fk := range e.modifyColumn(column, table) {
					pendingFKs = append(pendingFKs, pendingAlterFK{refTable: fk.RefTable, refColumns: fk.RefColumns})
				}
			}
		case ast.AT_DROP_COLUMN:
			if !flushPendingFKs() {
				return
			}
			columns := []string{action.ColumnName}
			if action.Options != nil && len(action.Options.Items) > 0 {
				columns = stringList(action.Options)
			}
			for _, column := range columns {
				if tableFindColumn(table, column) == nil {
					e.setError(oraError("00904", "%q: invalid identifier", column))
					return
				}
				e.dropColumn(table, column)
			}
		case ast.AT_RENAME_COLUMN:
			if tableFindColumn(table, action.ColumnName) == nil {
				e.setError(oraError("00904", "%q: invalid identifier", action.ColumnName))
				return
			}
			if tableFindColumn(table, action.NewName) != nil {
				e.setError(oraError("00957", "duplicate column name %q", action.NewName))
				return
			}
			e.renameColumn(table, action.ColumnName, action.NewName)
		case ast.AT_ADD_CONSTRAINT:
			if action.Constraint != nil {
				tc := action.Constraint
				for _, colName := range stringList(tc.Columns) {
					if tableFindColumnOrPath(table, colName) == nil {
						e.setError(oraError("00904", "%q: invalid identifier", colName))
						return
					}
				}
				if tc.Name != "" && e.constraintNameExists(tc.Name) {
					e.setError(oraError("02264", "name %q already used by an existing constraint", tc.Name))
					return
				}
				switch tc.Type {
				case ast.CONSTRAINT_PRIMARY:
					if hasPrimaryIndex(table) {
						e.setError(oraError("02260", "table %q can have only one primary key", table.Name))
						return
					}
				case ast.CONSTRAINT_FOREIGN:
					pendingFKs = append(pendingFKs, pendingAlterFK{refTable: tc.RefTable, refColumns: tc.RefColumns})
				}
				e.extractTableConstraint(tc, table)
			}
		case ast.AT_DROP_CONSTRAINT:
			if !flushPendingFKs() {
				return
			}
			var droppedColumns []string
			opts := stringList(action.Options)
			keepIndex := slices.Contains(opts, "KEEP INDEX")
			dropIndex := slices.Contains(opts, "DROP INDEX")
			cascade := slices.Contains(opts, "CASCADE")
			online := slices.Contains(opts, "ONLINE") || keywordOutsideQuotes(e.nodeText(action), "ONLINE")
			if action.Constraint != nil && action.Constraint.Name != "" {
				if !tableHasConstraint(table, action.Constraint.Name) {
					e.setError(oraError("02443", "cannot drop constraint %q - nonexistent constraint", action.Constraint.Name))
					return
				}
				droppedColumns = constraintKeyColumns(table, action.Constraint.Name)
				if len(droppedColumns) > 0 && !cascade && e.columnsHaveIncomingForeignKeys(table.Name, droppedColumns) {
					if online {
						e.setError(oraError("02297", "cannot disable/drop constraint - dependencies exist"))
					} else {
						e.setError(oraError("02273", "this unique/primary key is referenced by some foreign keys"))
					}
					return
				}
				e.dropConstraint(table, action.Constraint.Name, keepIndex, dropIndex)
			} else {
				var columns []string
				if action.Constraint != nil {
					columns = stringList(action.Constraint.Columns)
				}
				switch action.Subtype {
				case "PRIMARY KEY":
					if !hasPrimaryIndex(table) {
						e.setError(oraError("02441", "cannot drop nonexistent primary key"))
						return
					}
					columns = primaryKeyColumns(table)
				case "UNIQUE":
					for _, col := range columns {
						if tableFindColumnOrPath(table, col) == nil {
							e.setError(oraError("00904", "%q: invalid identifier", col))
							return
						}
					}
					if !tableHasUniqueConstraint(table, columns) {
						e.setError(oraError("02442", "cannot drop nonexistent unique key"))
						return
					}
				}
				if !cascade && e.columnsHaveIncomingForeignKeys(table.Name, columns) {
					if online {
						e.setError(oraError("02297", "cannot disable/drop constraint - dependencies exist"))
					} else {
						e.setError(oraError("02273", "this unique/primary key is referenced by some foreign keys"))
					}
					return
				}
				droppedColumns = columns
				e.dropConstraintSubtype(table, action.Subtype, columns, keepIndex, dropIndex)
			}
			if cascade && len(droppedColumns) > 0 {
				e.dropIncomingForeignKeysForColumns(table.Name, droppedColumns)
			}
		case ast.AT_ENABLE_DISABLE, ast.AT_MODIFY_CONSTRAINT:
			if e.err == nil {
				e.err = fmt.Errorf("constraint enable/disable is not supported by oracle/catalog")
			}
			return
		case ast.AT_ENABLE_DISABLE_TRIGGERS:
			if e.err == nil {
				e.err = fmt.Errorf("table-wide trigger enable/disable is not supported by oracle/catalog")
			}
			return
		case ast.AT_SET_UNUSED:
			columns, cascade := extractSetUnusedSpec(e.nodeText(action), action.ColumnName)
			seen := make(map[string]bool, len(columns))
			for _, column := range columns {
				if column == "" {
					continue
				}
				if seen[column] {
					e.setError(oraError("00957", "duplicate column name %q", column))
					return
				}
				seen[column] = true
				if tableFindColumn(table, column) == nil {
					e.setError(oraError("00904", "%q: invalid identifier", column))
					return
				}
			}
			if !e.serverConfig.ApplySetUnused {
				if e.err == nil {
					e.err = fmt.Errorf("SET UNUSED is not supported by oracle/catalog")
				}
				return
			}
			if len(seen) > 0 && len(seen) >= len(table.Columns) {
				e.setError(oraError("12983", "cannot drop all columns in a table"))
				return
			}
			if !cascade {
				if col := setUnusedMultiColumnConstraintColumn(table, seen); col != "" {
					e.setError(oraError("12991", "column %q is referenced in a multi-column constraint", col))
					return
				}
				for _, column := range columns {
					if e.columnReferencedByIncomingForeignKey(table, column) {
						e.setError(oraError("12992", "cannot drop parent key column %q", column))
						return
					}
				}
			}
			for _, column := range columns {
				e.dropColumn(table, column)
			}
		case ast.AT_RENAME_CONSTRAINT:
			if !tableHasConstraint(table, action.ColumnName) {
				e.setError(oraError("23292", "constraint %q does not exist", action.ColumnName))
				return
			}
			if e.constraintNameExists(action.NewName) {
				e.setError(oraError("02264", "name %q already used by an existing constraint", action.NewName))
				return
			}
			e.disambiguateSyntheticConstraint(table, action.NewName)
			e.renameConstraint(table, action.ColumnName, action.NewName)
		case ast.AT_RENAME:
			if action.NewName != "" {
				if e.objectExistsInNamespace(action.NewName) {
					e.setError(oraError("00955", "name %q is already used by an existing object", action.NewName))
					return
				}
				e.renameTable(tableName, action.NewName)
			}
		default:
		}
		if e.err != nil {
			return
		}
	}
	flushPendingFKs()
}

// modifyColumn applies an ALTER TABLE MODIFY column definition to the column
// the table already holds. Oracle keeps the attributes the clause does not
// name, including a type or NOT NULL when the clause only changes visibility.
// It returns any inline foreign-key constraints so the caller can validate
// them once all columns/keys in the ALTER TABLE statement have been installed.
func (e *oracleMetadataExtractor) modifyColumn(def *ast.ColumnDef, table *metadatapb.TableMetadata) []*ast.ColumnConstraint {
	column := tableFindColumn(table, def.Name)
	if column == nil {
		e.setError(oraError("00904", "%q: invalid identifier", def.Name))
		return nil
	}
	if def.Virtual != nil && column.GetGeneration() == nil {
		e.setError(oraError("54026", "real column %q cannot have an expression", def.Name))
		return nil
	}
	if !e.validateVirtualColumnEditions(def) {
		return nil
	}
	if (def.Identity != nil || def.DropIdentity) && !column.IsIdentity {
		e.setError(oraError("30673", "column %q to be modified is not an identity column", def.Name))
		return nil
	}
	if strings.EqualFold(def.Encrypt, "DECRYPT") {
		if !e.encryptedColumns[table.Name+"\x00"+def.Name] {
			e.setError(oraError("28333", "cannot decrypt a column that is not encrypted"))
			return nil
		}
	}
	if isColumnEncrypt(def.Encrypt) {
		if e.columnInForeignKey(table, def.Name) || colHasInlineReferences(def) {
			e.setError(oraError("28335", "referenced or referencing foreign key constraint column %q cannot be encrypted", def.Name))
			return nil
		}
		if e.serverConfig.KeystoreClosed {
			e.setError(oraError("28365", "wallet is not open"))
			return nil
		}
	}
	if def.Null && column.Nullable {
		e.setError(oraError("01451", "column %q to be modified to NULL cannot be modified to NULL", def.Name))
		return nil
	}
	var pendingFKs []*ast.ColumnConstraint
	for _, item := range listItems(def.Constraints) {
		cc, ok := item.(*ast.ColumnConstraint)
		if !ok {
			continue
		}
		if cc.Type == ast.CONSTRAINT_NULL && column.Nullable {
			e.setError(oraError("01451", "column %q to be modified to NULL cannot be modified to NULL", def.Name))
			return nil
		}
		if cc.Name != "" && e.constraintNameExists(cc.Name) {
			e.setError(oraError("02264", "name %q already used by an existing constraint", cc.Name))
			return nil
		}
		switch cc.Type {
		case ast.CONSTRAINT_PRIMARY:
			if hasPrimaryIndex(table) {
				e.setError(oraError("02260", "table %q can have only one primary key", table.Name))
				return nil
			}
		case ast.CONSTRAINT_FOREIGN:
			pendingFKs = append(pendingFKs, cc)
		}
	}

	if def.TypeName != nil {
		column.Type = normalizeDataTypeText(e.nodeText(def.TypeName))
	} else if def.Domain != nil {
		column.Type = objectName(def.Domain)
	} else if def.Virtual != nil {
		column.Type = inferVirtualColumnType(def.Virtual)
	}
	if def.Default != nil {
		column.Default = e.nodeText(def.Default)
		column.DefaultOnNull = def.DefaultOnNull
		if def.DefaultOnNull {
			column.Nullable = false
		}
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
		column.Options = nil
	} else if def.Identity != nil {
		applyIdentity(column, def.Identity, e.nodeText(def.Identity))
	}
	if def.Collation != "" {
		column.Collation = def.Collation
	} else if def.TypeName != nil && !oracleTypeSupportsCollation(column.Type) {
		column.Collation = ""
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
	if strings.EqualFold(def.Encrypt, "DECRYPT") {
		delete(e.encryptedColumns, table.Name+"\x00"+column.Name)
	} else if def.Encrypt != "" {
		e.encryptedColumns[table.Name+"\x00"+column.Name] = true
	}
	for _, item := range listItems(def.Constraints) {
		constraint, ok := item.(*ast.ColumnConstraint)
		if !ok {
			continue
		}
		if constraint.Type == ast.CONSTRAINT_NOT_NULL {
			column.Nullable = false
			if constraint.Name != "" {
				e.addNamedNotNull(table, constraint.Name, column.Name)
			}
			continue
		}
		e.extractColumnConstraint(constraint, table, column)
	}
	return pendingFKs
}

func (e *oracleMetadataExtractor) dropColumn(table *metadatapb.TableMetadata, name string) {
	delete(e.primaryKeyNotNull, table.Name+"\x00"+name)
	delete(e.encryptedColumns, table.Name+"\x00"+name)
	for _, check := range table.CheckConstraints {
		if e.namedNotNull[table.Name+"\x00"+check.Name] && expressionReferencesColumn(check.Expression, name) {
			delete(e.namedNotNull, table.Name+"\x00"+check.Name)
		}
	}
	var columns []*metadatapb.ColumnMetadata
	for _, column := range table.Columns {
		if column.Name != name {
			columns = append(columns, column)
		}
	}
	table.Columns = columns

	var indexes []*metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		matches := func(expression string) bool { return expressionReferencesColumn(expression, name) }
		if index.IsConstraint {
			matches = func(expression string) bool { return strings.Trim(expression, `"`) == name }
		}
		if !slices.ContainsFunc(index.Expressions, matches) {
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
		referencesColumn := fk.ReferencedTable == table.Name && slices.Contains(fk.ReferencedColumns, name) &&
			(fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema)
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
			if fk.ReferencedTable == table.Name && slices.Contains(fk.ReferencedColumns, name) &&
				(fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema) {
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
	oldKey := table.Name + "\x00" + oldName
	newKey := table.Name + "\x00" + newName
	if e.primaryKeyNotNull[oldKey] {
		delete(e.primaryKeyNotNull, oldKey)
		e.primaryKeyNotNull[newKey] = true
	}
	if e.encryptedColumns[oldKey] {
		delete(e.encryptedColumns, oldKey)
		e.encryptedColumns[newKey] = true
	}
	for _, column := range table.Columns {
		if column.Name == oldName {
			column.Name = newName
			break
		}
	}
	for _, index := range table.Indexes {
		for i, expression := range index.Expressions {
			if index.IsConstraint {
				if strings.Trim(expression, `"`) == oldName {
					index.Expressions[i] = newName
				}
				continue
			}
			index.Expressions[i] = renameIdentifier(expression, oldName, newName)
		}
		for _, constraint := range index.Constraints {
			for i, col := range constraint.Columns {
				if strings.Trim(col, `"`) == oldName {
					if strings.HasPrefix(col, `"`) && strings.HasSuffix(col, `"`) {
						constraint.Columns[i] = `"` + newName + `"`
					} else {
						constraint.Columns[i] = newName
					}
				}
			}
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
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable == table.Name && (fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema) {
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
		if end, ok := alternativeQuoteEnd(expression, i); ok {
			out.WriteString(expression[i:end])
			i = end
			continue
		}
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
				writeRenamedIdentifier(&out, newName, true)
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
			next := end
			for next < len(expression) && (expression[next] == ' ' || expression[next] == '\t' || expression[next] == '\n' || expression[next] == '\r') {
				next++
			}
			isFunction := next < len(expression) && expression[next] == '('
			if !isFunction && strings.ToUpper(token) == oldName {
				writeRenamedIdentifier(&out, newName, false)
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

func renderRenamedIdentifier(name string) string {
	if needsConstraintQuotes(name) {
		return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
	}
	return name
}

func writeRenamedIdentifier(out *strings.Builder, name string, forceQuote bool) {
	quote := forceQuote || needsConstraintQuotes(name)
	if !quote {
		out.WriteString(name)
		return
	}
	out.WriteString(`"`)
	out.WriteString(strings.ReplaceAll(name, `"`, `""`))
	out.WriteString(`"`)
}

func expressionReferencesFunction(expression, function string) bool {
	for i := 0; i < len(expression); {
		if end, ok := alternativeQuoteEnd(expression, i); ok {
			i = end
			continue
		}
		if expression[i] == '\'' {
			if end := strings.IndexByte(expression[i+1:], '\''); end >= 0 {
				i += end + 2
			} else {
				return false
			}
			continue
		}
		if expression[i] == '"' {
			end := i + 1
			for end < len(expression) && expression[end] != '"' {
				end++
			}
			token := expression[i+1 : end]
			if end < len(expression) {
				end++
			}
			next := end
			for next < len(expression) && (expression[next] == ' ' || expression[next] == '\t' || expression[next] == '\n' || expression[next] == '\r') {
				next++
			}
			if token == function && next < len(expression) && expression[next] == '(' {
				return true
			}
			i = end
			continue
		}
		if !isIdentifierByte(expression[i]) {
			i++
			continue
		}
		end := i + 1
		for end < len(expression) && isIdentifierByte(expression[end]) {
			end++
		}
		next := end
		for next < len(expression) && (expression[next] == ' ' || expression[next] == '\t' || expression[next] == '\n' || expression[next] == '\r') {
			next++
		}
		if strings.EqualFold(expression[i:end], function) && next < len(expression) && expression[next] == '(' {
			return true
		}
		i = end
	}
	return false
}

func expressionReferencesPackage(expression, pkg string) bool {
	if pkg == "" {
		return false
	}
	for i := 0; i < len(expression); {
		if end, ok := alternativeQuoteEnd(expression, i); ok {
			i = end
			continue
		}
		if i+1 < len(expression) && expression[i] == '-' && expression[i+1] == '-' {
			if end := strings.IndexByte(expression[i:], '\n'); end >= 0 {
				i += end + 1
			} else {
				return false
			}
			continue
		}
		if i+1 < len(expression) && expression[i] == '/' && expression[i+1] == '*' {
			if end := strings.Index(expression[i+2:], "*/"); end >= 0 {
				i += end + 4
			} else {
				return false
			}
			continue
		}
		if expression[i] == '\'' {
			if end := strings.IndexByte(expression[i+1:], '\''); end >= 0 {
				i += end + 2
			} else {
				return false
			}
			continue
		}
		matched := false
		end := i
		if expression[i] == '"' {
			end = i + 1
			for end < len(expression) && expression[end] != '"' {
				end++
			}
			matched = expression[i+1:end] == pkg
			if end < len(expression) {
				end++
			}
		} else if isIdentifierByte(expression[i]) {
			end = i + 1
			for end < len(expression) && isIdentifierByte(expression[end]) {
				end++
			}
			matched = strings.EqualFold(expression[i:end], pkg)
		} else {
			i++
			continue
		}
		if matched {
			next := end
			for next < len(expression) && (expression[next] == ' ' || expression[next] == '\t' || expression[next] == '\n' || expression[next] == '\r') {
				next++
			}
			if next < len(expression) && expression[next] == '.' {
				member := next + 1
				for member < len(expression) && (expression[member] == ' ' || expression[member] == '\t' || expression[member] == '\n' || expression[member] == '\r') {
					member++
				}
				if member < len(expression) && (expression[member] == '"' || isIdentifierByte(expression[member])) {
					return true
				}
			}
		}
		i = end
	}
	return false
}

func viewReferencesSynonym(definition string, deps []*metadatapb.DependencyColumn, synName string) bool {
	if synName == "" {
		return false
	}
	for _, dep := range deps {
		if dep.GetTable() == synName {
			return true
		}
	}
	for i := 0; i < len(definition); {
		if end, ok := alternativeQuoteEnd(definition, i); ok {
			i = end
			continue
		}
		if i+1 < len(definition) && definition[i] == '-' && definition[i+1] == '-' {
			if end := strings.IndexByte(definition[i:], '\n'); end >= 0 {
				i += end + 1
			} else {
				return false
			}
			continue
		}
		if i+1 < len(definition) && definition[i] == '/' && definition[i+1] == '*' {
			if end := strings.Index(definition[i+2:], "*/"); end >= 0 {
				i += end + 4
			} else {
				return false
			}
			continue
		}
		if definition[i] == '\'' {
			if end := strings.IndexByte(definition[i+1:], '\''); end >= 0 {
				i += end + 2
			} else {
				return false
			}
			continue
		}
		if definition[i] == '"' {
			end := i + 1
			for end < len(definition) && definition[end] != '"' {
				end++
			}
			if definition[i+1:end] == synName {
				return true
			}
			if end < len(definition) {
				end++
			}
			i = end
			continue
		}
		if !isIdentifierByte(definition[i]) {
			i++
			continue
		}
		end := i + 1
		for end < len(definition) && isIdentifierByte(definition[end]) {
			end++
		}
		if strings.EqualFold(definition[i:end], synName) {
			return true
		}
		i = end
	}
	return false
}

func isIdentifierByte(c byte) bool {
	return c == '_' || c == '$' || c == '#' ||
		(c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func (e *oracleMetadataExtractor) dropConstraint(table *metadatapb.TableMetadata, name string, keepIndex, dropIndex bool) {
	if name == "" {
		return
	}
	delete(e.syntheticConstraints, table.Name+"\x00"+name)
	var indexes []*metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if len(index.Constraints) > 0 {
			reused := needsReusedIndexClause(index)
			var remaining []*metadatapb.IndexConstraintMetadata
			matched := false
			for _, c := range index.Constraints {
				if c.Name == name {
					matched = true
					if c.Primary {
						cols := c.Columns
						if len(cols) == 0 {
							cols = index.Expressions
						}
						e.clearPrimaryKeyNotNull(table, cols)
					}
					continue
				}
				remaining = append(remaining, c)
			}
			if !matched {
				indexes = append(indexes, index)
				continue
			}
			if len(remaining) > 0 {
				index.Constraints = remaining
				index.IsConstraint = true
				index.Primary = slices.ContainsFunc(remaining, func(c *metadatapb.IndexConstraintMetadata) bool { return c.Primary })
				indexes = append(indexes, index)
				continue
			}
			if keepIndex || (reused && !dropIndex) {
				index.Constraints = nil
				index.IsConstraint = false
				index.Primary = false
				indexes = append(indexes, index)
			}
			continue
		}
		if index.Name == name && index.IsConstraint {
			if index.Primary {
				e.clearPrimaryKeyNotNull(table, index.Expressions)
			}
			if keepIndex {
				index.IsConstraint = false
				index.Primary = false
				indexes = append(indexes, index)
			}
			continue
		}
		indexes = append(indexes, index)
	}
	table.Indexes = indexes
	var checks []*metadatapb.CheckConstraintMetadata
	for _, check := range table.CheckConstraints {
		if check.Name == name {
			if e.namedNotNull[table.Name+"\x00"+check.Name] {
				e.clearNamedNotNull(table, check)
			}
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

func (e *oracleMetadataExtractor) dropConstraintSubtype(table *metadatapb.TableMetadata, subtype string, columns []string, keepIndex, dropIndex bool) {
	var wantSorted []string
	if subtype == "UNIQUE" && len(columns) > 0 {
		wantSorted = make([]string, len(columns))
		for i, c := range columns {
			wantSorted[i] = strings.Trim(c, `"`)
		}
		slices.Sort(wantSorted)
	}
	var indexes []*metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if len(index.Constraints) > 0 {
			reused := needsReusedIndexClause(index)
			var remaining []*metadatapb.IndexConstraintMetadata
			matched := false
			for _, c := range index.Constraints {
				cCols := c.Columns
				if len(cCols) == 0 {
					cCols = index.Expressions
				}
				drop := false
				if subtype == "PRIMARY KEY" && c.Primary {
					drop = true
				} else if subtype == "UNIQUE" && !c.Primary {
					drop = len(wantSorted) == 0 || columnsMatchSet(cCols, wantSorted)
				}
				if drop {
					matched = true
					if c.Primary {
						e.clearPrimaryKeyNotNull(table, cCols)
					}
					delete(e.syntheticConstraints, table.Name+"\x00"+c.Name)
					continue
				}
				remaining = append(remaining, c)
			}
			if !matched {
				indexes = append(indexes, index)
				continue
			}
			if len(remaining) > 0 {
				index.Constraints = remaining
				index.IsConstraint = true
				index.Primary = slices.ContainsFunc(remaining, func(c *metadatapb.IndexConstraintMetadata) bool { return c.Primary })
				indexes = append(indexes, index)
				continue
			}
			if keepIndex || (reused && !dropIndex) {
				index.Constraints = nil
				index.IsConstraint = false
				index.Primary = false
				indexes = append(indexes, index)
			}
			continue
		}
		drop := subtype == "PRIMARY KEY" && index.Primary
		if subtype == "UNIQUE" && index.Unique && !index.Primary {
			drop = len(columns) == 0 || expressionsMatchColumns(index.Expressions, columns)
		}
		if !drop {
			indexes = append(indexes, index)
			continue
		}
		if index.Primary {
			e.clearPrimaryKeyNotNull(table, index.Expressions)
		}
		delete(e.syntheticConstraints, table.Name+"\x00"+index.Name)
		if keepIndex {
			index.IsConstraint = false
			index.Primary = false
			indexes = append(indexes, index)
		}
	}
	table.Indexes = indexes
}

func (e *oracleMetadataExtractor) clearPrimaryKeyNotNull(table *metadatapb.TableMetadata, columns []string) {
	for _, columnName := range columns {
		columnName = strings.Trim(columnName, `"`)
		key := table.Name + "\x00" + columnName
		if !e.primaryKeyNotNull[key] {
			continue
		}
		delete(e.primaryKeyNotNull, key)
		for _, column := range table.Columns {
			if column.Name == columnName {
				column.Nullable = true
			}
		}
	}
}

func expressionsMatchColumns(expressions, columns []string) bool {
	if len(expressions) != len(columns) {
		return false
	}
	for i, expression := range expressions {
		expression = strings.Trim(expression, `"`)
		column := strings.Trim(columns[i], `"`)
		if expression != column {
			return false
		}
	}
	return true
}

func (e *oracleMetadataExtractor) renameConstraint(table *metadatapb.TableMetadata, oldName, newName string) {
	if oldName == "" || newName == "" {
		return
	}
	for _, index := range table.Indexes {
		if len(index.Constraints) > 0 {
			for _, constraint := range index.Constraints {
				if constraint.Name == oldName {
					constraint.Name = newName
				}
			}
		} else if index.Name == oldName && index.IsConstraint {
			index.Name = newName
		}
	}
	oldKey := table.Name + "\x00" + oldName
	newKey := table.Name + "\x00" + newName
	delete(e.syntheticConstraints, oldKey)
	if e.namedNotNull[oldKey] {
		delete(e.namedNotNull, oldKey)
		e.namedNotNull[newKey] = true
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

// renameTriggerTarget rewrites the ON relation in a trigger declaration and
// leaves the rest of its body alone.
func renameTriggerTarget(body, oldName, newName string) string {
	upper := strings.ToUpper(body)
	for offset := 0; ; {
		match := triggerOnClause.FindStringIndex(upper[offset:])
		if match == nil {
			return body
		}
		start := offset + match[1]
		first, firstStart, end := scanTriggerIdentifier(body, start)
		relation, relationStart, relationEnd := first, firstStart, end
		if end < len(body) && body[end] == '.' {
			relation, relationStart, relationEnd = scanTriggerIdentifier(body, end+1)
		}
		if strings.EqualFold(relation, oldName) {
			return body[:relationStart] + renderRenamedIdentifier(newName) + body[relationEnd:]
		}
		offset = end
	}
}

var triggerOnClause = regexp.MustCompile(`(?i)\bON\s+`)

func scanTriggerIdentifier(body string, start int) (string, int, int) {
	if start >= len(body) {
		return "", start, start
	}
	if body[start] == '"' {
		end := start + 1
		for end < len(body) && body[end] != '"' {
			end++
		}
		token := body[start+1 : end]
		if end < len(body) {
			end++
		}
		return token, start, end
	}
	end := start
	for end < len(body) && isIdentifierByte(body[end]) {
		end++
	}
	return body[start:end], start, end
}

func rekeyProvenance(provenance map[string]bool, oldTable, newTable string) {
	prefix := oldTable + "\x00"
	for key := range provenance {
		if strings.HasPrefix(key, prefix) {
			value := provenance[key]
			delete(provenance, key)
			provenance[newTable+"\x00"+strings.TrimPrefix(key, prefix)] = value
		}
	}
}

func (e *oracleMetadataExtractor) renameTable(oldName, newName string) {
	table := e.tables[oldName]
	if table == nil {
		return
	}
	delete(e.tables, oldName)
	rekeyProvenance(e.primaryKeyNotNull, oldName, newName)
	rekeyProvenance(e.namedNotNull, oldName, newName)
	rekeyProvenance(e.syntheticConstraints, oldName, newName)
	rekeyProvenance(e.encryptedColumns, oldName, newName)
	table.Name = newName
	for _, trigger := range table.Triggers {
		trigger.Body = renameTriggerTarget(trigger.Body, oldName, newName)
	}
	for _, fk := range table.ForeignKeys {
		if fk.ReferencedTable == oldName && (fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema) {
			fk.ReferencedTable = newName
		}
	}
	for _, other := range e.tables {
		for _, fk := range other.ForeignKeys {
			if fk.ReferencedTable == oldName && (fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema) {
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
			if !e.checkSchema(nameNode.Schema, "00942") {
				return
			}
			if e.tables[name] == nil {
				if n.IfExists {
					continue
				}
				e.setError(oraError("00942", "table or view %q does not exist", name))
				return
			}
			if !n.Cascade && e.tableHasIncomingForeignKeys(name) {
				e.setError(oraError("02449", "unique/primary keys in table %q referenced by foreign keys", name))
				return
			}
			if n.Cascade {
				e.dropIncomingForeignKeys(name)
			}
			e.dropTableProvenance(name)
			delete(e.tables, name)
		case ast.OBJECT_VIEW:
			if !e.checkSchema(nameNode.Schema, "00942") {
				return
			}
			if e.views[name] == nil {
				if n.IfExists {
					continue
				}
				e.setError(oraError("00942", "table or view %q does not exist", name))
				return
			}
			delete(e.views, name)
			delete(e.viewConstraints, name)
			delete(e.editioningViews, name)
		case ast.OBJECT_MATERIALIZED_VIEW:
			if !e.checkSchema(nameNode.Schema, "12003") {
				return
			}
			delete(e.materializedViews, name)
		case ast.OBJECT_SEQUENCE:
			if !e.checkSchema(nameNode.Schema, "02289") {
				return
			}
			if e.sequences[name] == nil {
				if n.IfExists {
					continue
				}
				e.setError(oraError("02289", "sequence %q does not exist", name))
				return
			}
			delete(e.sequences, name)
			delete(e.sequenceCurrent, name)
		case ast.OBJECT_FUNCTION:
			if !e.useSchema(nameNode.Schema) {
				continue
			}
			delete(e.functions, name)
		case ast.OBJECT_PROCEDURE:
			if !e.useSchema(nameNode.Schema) {
				continue
			}
			delete(e.procedures, name)
		case ast.OBJECT_PACKAGE:
			if !e.useSchema(nameNode.Schema) {
				continue
			}
			delete(e.packages, name)
		case ast.OBJECT_PACKAGE_BODY:
			if !e.useSchema(nameNode.Schema) {
				continue
			}
			if pkg := e.packages[name]; pkg != nil {
				spec, _ := splitPackageUnits(pkg.Definition)
				if spec == "" {
					delete(e.packages, name)
				} else {
					pkg.Definition = spec
				}
			}
		case ast.OBJECT_SYNONYM:
			if isDropPublicSynonym(e.nodeText(n)) {
				if e.publicSynonyms[name] == nil {
					if n.IfExists {
						continue
					}
					e.setError(oraError("01432", "public synonym %q to be dropped does not exist", name))
					return
				}
				delete(e.publicSynonyms, name)
			} else {
				if !e.checkSchema(nameNode.Schema, "01434") {
					return
				}
				if e.synonyms[name] == nil {
					if n.IfExists {
						continue
					}
					e.setError(oraError("01434", "private synonym %q to be dropped does not exist", name))
					return
				}
				delete(e.synonyms, name)
			}
		case ast.OBJECT_INDEX:
			if !e.checkSchema(nameNode.Schema, "01418") {
				return
			}
			idx := e.findIndex(name)
			if idx == nil {
				if n.IfExists {
					continue
				}
				e.setError(oraError("01418", "specified index %q does not exist", name))
				return
			}
			if n.Force && !strings.EqualFold(idx.Type, "DOMAIN") {
				e.setError(oraError("29862", "cannot specify FORCE option for dropping a non-domain index"))
				return
			}
			e.dropIndex(name)
		case ast.OBJECT_TRIGGER:
			if !e.useSchema(nameNode.Schema) {
				continue
			}
			e.dropTrigger(name)
		case ast.OBJECT_DOMAIN:
			if !e.checkSchema(nameNode.Schema, "01031") {
				return
			}
			if !e.domains[name] {
				if n.IfExists {
					continue
				}
				e.setError(oraError("11504", "domain %q does not exist", name))
				return
			}
			delete(e.domains, name)
		case ast.OBJECT_EDITION:
			if !e.editions[name] {
				if n.IfExists {
					continue
				}
				e.setError(oraError("38802", "edition %q does not exist", name))
				return
			}
			delete(e.editions, name)
		default:
			if !e.useSchema(nameNode.Schema) {
				continue
			}
		}
	}
}

func primaryKeyColumns(table *metadatapb.TableMetadata) []string {
	for _, index := range table.Indexes {
		for _, c := range index.Constraints {
			if c.Primary {
				if len(c.Columns) > 0 {
					return append([]string(nil), c.Columns...)
				}
				return append([]string(nil), index.Expressions...)
			}
		}
		if index.Primary {
			return append([]string(nil), index.Expressions...)
		}
	}
	return nil
}

func (e *oracleMetadataExtractor) dropIncomingForeignKeysForColumns(tableName string, columns []string) {
	if len(columns) == 0 {
		e.dropIncomingForeignKeys(tableName)
		return
	}
	columnSet := map[string]bool{}
	for _, column := range columns {
		columnSet[strings.Trim(column, `"`)] = true
	}
	for _, table := range e.tables {
		var foreignKeys []*metadatapb.ForeignKeyMetadata
		for _, fk := range table.ForeignKeys {
			matches := fk.ReferencedTable == tableName
			if matches && len(fk.ReferencedColumns) > 0 {
				matches = len(fk.ReferencedColumns) == len(columnSet)
				if matches {
					for _, column := range fk.ReferencedColumns {
						if !columnSet[strings.Trim(column, `"`)] {
							matches = false
							break
						}
					}
				}
			}
			if matches && (fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema) {
				continue
			}
			foreignKeys = append(foreignKeys, fk)
		}
		table.ForeignKeys = foreignKeys
	}
}

func (e *oracleMetadataExtractor) dropIncomingForeignKeys(tableName string) {
	for _, table := range e.tables {
		if table.Name == tableName {
			continue
		}
		var foreignKeys []*metadatapb.ForeignKeyMetadata
		for _, fk := range table.ForeignKeys {
			if fk.ReferencedTable == tableName && (fk.ReferencedSchema == "" || fk.ReferencedSchema == e.currentSchema) {
				continue
			}
			foreignKeys = append(foreignKeys, fk)
		}
		table.ForeignKeys = foreignKeys
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

func (e *oracleMetadataExtractor) attachConstraintToExistingIndex(table *metadatapb.TableMetadata, indexName, constraintName string, columns []string, primary bool) bool {
	for _, index := range table.Indexes {
		if index.Name != indexName {
			continue
		}
		index.IsConstraint = true
		if primary {
			index.Primary = true
		}
		for _, existing := range index.Constraints {
			if existing.Name == constraintName {
				return true
			}
		}
		constraint := &metadatapb.IndexConstraintMetadata{
			Name:    constraintName,
			Columns: append([]string(nil), columns...),
			Primary: primary,
		}
		if primary {
			index.Constraints = append([]*metadatapb.IndexConstraintMetadata{constraint}, index.Constraints...)
		} else {
			index.Constraints = append(index.Constraints, constraint)
		}
		return true
	}
	return false
}

// constraintUsesExistingIndex reports a constraint that Oracle backs with an
// index the statement names (USING INDEX [schema.]index), or that it asks
// Oracle to create inline (USING INDEX (CREATE INDEX ...)). The snapshot holds
// the constraint's index, not the separate index it reuses, so the definition
// would create both. The physical properties form (USING INDEX TABLESPACE ...,
// USING INDEX LOCAL) names no index and stays a documented gap: the catalog
// ignores it like every other storage clause.
func (e *oracleMetadataExtractor) constraintUsesExistingIndex(indexName string) bool {
	if indexName == "" {
		return false
	}
	if e.err == nil {
		e.err = fmt.Errorf("a constraint backed by the existing index %q is not supported by oracle/catalog", indexName)
	}
	return true
}

func (e *oracleMetadataExtractor) constraintDisabled(loc ast.Loc) bool {
	if loc.Start < 0 || loc.End < loc.Start || loc.End > len(e.schemaText) {
		return false
	}
	if keywordOutsideQuotes(e.schemaText[loc.Start:loc.End], "DISABLE") {
		if e.err == nil {
			e.err = fmt.Errorf("disabled constraints are not supported by oracle/catalog")
		}
		return true
	}
	return false
}

func keywordOutsideQuotes(text, keyword string) bool {
	var quote byte
	for i := 0; i < len(text); i++ {
		c := text[i]
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			quote = c
			continue
		}
		if i > 0 && isIdentifierByte(text[i-1]) {
			continue
		}
		if len(text)-i >= len(keyword) && strings.EqualFold(text[i:i+len(keyword)], keyword) &&
			(i+len(keyword) == len(text) || !isIdentifierByte(text[i+len(keyword)])) {
			return true
		}
	}
	return false
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

func (e *oracleMetadataExtractor) extractAlterIndex(n *ast.AlterIndexStmt) {
	name := objectName(n.Name)
	if name == "" || !e.useSchema(n.Name.Schema) {
		return
	}
	var index *metadatapb.IndexMetadata
	for _, table := range e.tables {
		for _, candidate := range table.Indexes {
			if candidate.Name == name {
				index = candidate
			}
		}
	}
	for _, view := range e.materializedViews {
		for _, candidate := range view.Indexes {
			if candidate.Name == name {
				index = candidate
			}
		}
	}
	if index == nil {
		return
	}
	switch strings.ToUpper(n.Action) {
	case "RENAME":
		if n.NewName != "" {
			index.Name = n.NewName
		}
	case "VISIBLE":
		index.Visible = true
	case "INVISIBLE":
		index.Visible = false
	default:
		if e.err == nil {
			e.err = fmt.Errorf("ALTER INDEX action %q is not supported by oracle/catalog", n.Action)
		}
	}
}

func (e *oracleMetadataExtractor) extractCreateIndex(n *ast.CreateIndexStmt) {
	tableName := objectName(n.Table)
	indexName := objectName(n.Name)
	if tableName == "" || indexName == "" {
		return
	}
	// Check duplicate column references before table lookup, matching Oracle.
	seenCols := map[string]bool{}
	for _, item := range listItems(n.Columns) {
		column, ok := item.(*ast.IndexColumn)
		if !ok {
			continue
		}
		if colRef, ok := column.Expr.(*ast.ColumnRef); ok && colRef.Column != "" {
			key := indexColumnRefKey(colRef)
			if seenCols[key] {
				e.setError(oraError("00957", "duplicate column name %q", key))
				return
			}
			seenCols[key] = true
		}
	}
	if !e.checkSchema(n.Name.Schema, "00942") || !e.checkSchema(n.Table.Schema, "00942") {
		return
	}

	resolved := e.resolveSynonymTarget(tableName, n.Table.Schema)
	if e.views[resolved] != nil {
		e.setError(oraError("01702", "a view is not appropriate here"))
		return
	}
	if e.sequences[resolved] != nil {
		e.setError(oraError("02201", "sequence not allowed here"))
		return
	}
	if e.tables[resolved] == nil && e.materializedViews[resolved] == nil {
		e.setError(oraError("00942", "table or view %q does not exist", tableName))
		return
	}
	if e.explicitIndexExists(indexName) {
		if n.IfNotExists {
			return
		}
		e.setError(oraError("00955", "name %q is already used by an existing object", indexName))
		return
	}
	e.disambiguateSyntheticIndex(indexName)
	if table := e.tables[resolved]; table != nil {
		for _, item := range listItems(n.Columns) {
			column, ok := item.(*ast.IndexColumn)
			if !ok {
				continue
			}
			if colRef, ok := column.Expr.(*ast.ColumnRef); ok && colRef.Column != "" {
				if !indexColumnRefExists(table, n.Alias, colRef) {
					e.setError(oraError("00904", "%q: invalid identifier", indexColumnRefKey(colRef)))
					return
				}
			}
		}
	}

	index := &metadatapb.IndexMetadata{
		Name:        indexName,
		Unique:      n.Unique,
		Type:        "NORMAL",
		Expressions: []string{},
		Visible:     !n.Invisible,
	}
	if n.IndexType != nil {
		index.Type = "DOMAIN"
	} else if n.Bitmap {
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

	if materializedView := e.materializedViews[resolved]; materializedView != nil {
		materializedView.Indexes = append(materializedView.Indexes, index)
		return
	}
	table := e.tables[resolved]
	table.Indexes = append(table.Indexes, index)
}

func (e *oracleMetadataExtractor) extractCreateView(n *ast.CreateViewStmt) {
	viewName := objectName(n.Name)
	if viewName == "" {
		return
	}
	if e.serverConfig.NonApplicationRoot && isApplicationSharing(n.Sharing) {
		e.setError(oraError("65021", "illegal use of SHARING clause"))
		return
	}
	if e.isStandardMaxStringSize() && n.DefaultCollation != "" {
		e.setError(oraError("43929", "Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set"))
		return
	}
	// Check duplicate explicit column names before schema/existence checks, matching Oracle.
	if !n.Materialized && n.Columns != nil && n.Columns.Len() > 0 {
		seen := map[string]bool{}
		for _, item := range n.Columns.Items {
			colName := ""
			switch v := item.(type) {
			case *ast.String:
				colName = v.Str
			default:
				colName = strings.TrimSpace(e.nodeText(item))
			}
			if colName == "" {
				continue
			}
			if seen[colName] {
				e.setError(oraError("00957", "duplicate column name %q", colName))
				return
			}
			seen[colName] = true
		}
	}
	if !e.checkSchema(n.Name.Schema, "01917") {
		return
	}
	if !n.Materialized {
		if e.nonViewObjectExists(viewName) {
			if n.IfNotExists {
				return
			}
			e.setError(oraError("00955", "name %q is already used by an existing object", viewName))
			return
		}
		if !n.Force {
			if n.Columns != nil && n.Columns.Len() > 0 {
				invisible := viewInvisibleAliasPositions(e.schemaText, n)
				visibleCount := 0
				for i := range n.Columns.Items {
					if !invisible[i] {
						visibleCount++
					}
				}
				if visibleCount == 0 {
					e.setError(oraError("54039", "view must have at least one column that is not invisible"))
					return
				}
			} else if queryHasUnaliasedExpression(n.Query) {
				e.setError(oraError("00998", "must name this expression with a column alias"))
				return
			}
		}
		if e.views[viewName] != nil && !n.OrReplace {
			if n.IfNotExists {
				return
			}
			e.setError(oraError("00955", "name %q is already used by an existing object", viewName))
			return
		}
	}

	definition := e.nodeText(n.Query)
	if n.WithReadOnly {
		definition = strings.TrimSpace(definition) + " WITH READ ONLY"
	} else if n.WithCheckOpt {
		restriction := " WITH CHECK OPTION"
		if n.ConstraintName != "" {
			restriction += ` CONSTRAINT "` + n.ConstraintName + `"`
		}
		definition = strings.TrimSpace(definition) + restriction
	}
	if n.Materialized {
		if prefix := e.materializedViewPrefix(n); !e.omitMaterializedViewClauses && prefix != "" {
			definition = prefix + " AS " + strings.TrimSpace(definition)
		}
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
	if n.Columns != nil && n.Columns.Len() > 0 {
		invisible := viewInvisibleAliasPositions(e.schemaText, n)
		for i, item := range n.Columns.Items {
			colName := ""
			switch v := item.(type) {
			case *ast.String:
				colName = v.Str
			default:
				colName = strings.TrimSpace(e.nodeText(item))
			}
			if colName == "" {
				continue
			}
			view.Columns = append(view.Columns, &metadatapb.ColumnMetadata{
				Name:        colName,
				Position:    int32(i + 1),
				Nullable:    true,
				IsInvisible: invisible[i],
			})
		}
	}
	if table := e.tables[viewName]; table != nil {
		view.Triggers = append(view.Triggers, table.Triggers...)
		delete(e.tables, viewName)
	}
	e.views[viewName] = view
	var vcs []viewConstraint
	for _, item := range listItems(n.Constraints) {
		tc, ok := item.(*ast.TableConstraint)
		if !ok {
			continue
		}
		vcs = append(vcs, viewConstraint{
			name:    tc.Name,
			kind:    tc.Type,
			columns: stringList(tc.Columns),
		})
	}
	vcs = append(vcs, extractViewInlineConstraints(e.schemaText, n)...)
	if (n.WithReadOnly || n.WithCheckOpt) && n.ConstraintName != "" {
		vcs = append(vcs, viewConstraint{
			name: n.ConstraintName,
			kind: ast.CONSTRAINT_CHECK,
		})
	}
	if len(vcs) > 0 {
		e.viewConstraints[viewName] = vcs
	} else {
		delete(e.viewConstraints, viewName)
	}
	if n.Editioning {
		e.editioningViews[viewName] = true
	} else {
		delete(e.editioningViews, viewName)
	}
}

func (e *oracleMetadataExtractor) materializedViewPrefix(n *ast.CreateViewStmt) string {
	var parts []string
	if n.OnPrebuilt {
		parts = append(parts, "ON PREBUILT TABLE")
	} else if n.BuildMode != "" {
		parts = append(parts, "BUILD "+strings.ToUpper(n.BuildMode))
	}
	if n.NeverRefresh {
		parts = append(parts, "NEVER REFRESH")
	} else if n.RefreshMethod != "" || n.RefreshMode != "" || n.RefreshOnStmt || n.StartWith != nil || n.Next != nil {
		refresh := "REFRESH"
		if n.RefreshMethod != "" {
			refresh += " " + strings.ToUpper(n.RefreshMethod)
		}
		if n.RefreshOnStmt {
			refresh += " ON STATEMENT"
		} else if n.RefreshMode != "" {
			refresh += " " + strings.ToUpper(n.RefreshMode)
		}
		if n.StartWith != nil {
			refresh += " START WITH " + e.nodeText(n.StartWith)
		}
		if n.Next != nil {
			refresh += " NEXT " + e.nodeText(n.Next)
		}
		parts = append(parts, refresh)
	}
	return strings.Join(parts, " ")
}

func (e *oracleMetadataExtractor) extractCreateSequence(n *ast.CreateSequenceStmt) {
	sequenceName := objectName(n.Name)
	if sequenceName == "" {
		return
	}
	if e.serverConfig.NonApplicationRoot && isApplicationSharing(n.Sharing) {
		e.setError(oraError("65021", "illegal use of SHARING clause"))
		return
	}
	if e.serverConfig.DisableShardDDL && (n.Shard || n.ShardExtend || n.ShardNoExtend || n.NoShard) {
		e.setError(oraError("02511", "Shard DDL is not enabled"))
		return
	}
	if !e.checkSchema(n.Name.Schema, "01031") {
		return
	}
	if n.IfNotExists && e.sequences[sequenceName] != nil {
		return
	}
	if !n.IfNotExists && e.objectExistsInNamespace(sequenceName) {
		e.setError(oraError("00955", "name %q is already used by an existing object", sequenceName))
		return
	}
	inc := big.NewInt(1)
	if v, ok := parseSignedBigInt(e.nodeText(n.IncrementBy)); ok {
		inc = v
	}
	if cacheStr := e.nodeText(n.Cache); cacheStr != "" {
		if v, ok := parseSignedBigInt(cacheStr); ok && v.Cmp(big.NewInt(2)) < 0 {
			e.setError(oraError("04010", "the number of values to CACHE must be greater than 1"))
			return
		}
	}
	if n.Cycle && inc.Sign() > 0 && n.MaxValue == nil {
		e.setError(oraError("04015", "ascending sequences that CYCLE must specify MAXVALUE"))
		return
	}
	minText := e.nodeText(n.MinValue)
	minVal, hasMin := parseSignedBigInt(minText)
	if strings.TrimSpace(minText) == "" && inc.Sign() > 0 {
		minVal, hasMin = big.NewInt(1), true
	}
	maxVal, hasMax := parseSignedBigInt(e.nodeText(n.MaxValue))
	if hasMax && hasMin && minVal.Cmp(maxVal) >= 0 {
		e.setError(oraError("04004", "MINVALUE must be less than MAXVALUE"))
		return
	}
	startVal, hasStart := parseSignedBigInt(e.nodeText(n.StartWith))

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
	e.sequenceCurrent[sequenceName] = e.sequenceInitialCurrent(inc, startVal, minVal, maxVal, hasStart, hasMin, hasMax)
}

func (e *oracleMetadataExtractor) extractAlterSequence(n *ast.AlterSequenceStmt) {
	name := objectName(n.Name)
	if name == "" {
		return
	}
	if e.serverConfig.DisableShardDDL && (n.Shard || n.ShardExtend || n.ShardNoExtend || n.NoShard) {
		e.setError(oraError("02511", "Shard DDL is not enabled"))
		return
	}
	if !e.checkSchema(n.Name.Schema, "02289") {
		return
	}
	sequence := e.sequences[name]
	if sequence == nil {
		if n.IfExists {
			return
		}
		e.setError(oraError("02289", "sequence %q does not exist", name))
		return
	}
	if cacheStr := e.nodeText(n.Cache); cacheStr != "" {
		if v, ok := parseSignedBigInt(cacheStr); ok && v.Cmp(big.NewInt(2)) < 0 {
			e.setError(oraError("04010", "the number of values to CACHE must be greater than 1"))
			return
		}
	}
	inc := big.NewInt(1)
	if v, ok := parseSignedBigInt(sequence.Increment); ok {
		inc = v
	}
	if v, ok := parseSignedBigInt(e.nodeText(n.IncrementBy)); ok {
		inc = v
	}
	minText := sequence.MinValue
	if n.MinValue != nil {
		minText = e.nodeText(n.MinValue)
	} else if n.NoMinValue {
		minText = ""
	}
	minVal, hasMin := parseSignedBigInt(minText)
	if strings.TrimSpace(minText) == "" && inc.Sign() > 0 {
		minVal, hasMin = big.NewInt(1), true
	}
	maxText := sequence.MaxValue
	if n.MaxValue != nil {
		maxText = e.nodeText(n.MaxValue)
	} else if n.NoMaxValue {
		maxText = ""
	}
	maxVal, hasMax := parseSignedBigInt(maxText)
	if strings.TrimSpace(maxText) == "" && inc.Sign() < 0 {
		maxVal, hasMax = big.NewInt(-1), true
	}
	if hasMax && hasMin && minVal.Cmp(maxVal) >= 0 {
		e.setError(oraError("04004", "MINVALUE must be less than MAXVALUE"))
		return
	}
	if currVal, hasCurr := e.sequenceCurrent[name]; hasCurr && currVal != nil && !n.Restart && n.RestartWith == nil {
		if (n.MinValue != nil || n.NoMinValue) && hasMin && minVal.Cmp(currVal) > 0 {
			e.setError(oraError("04007", "MINVALUE cannot be made to exceed the current value"))
			return
		}
		if (n.MaxValue != nil || n.NoMaxValue) && hasMax && maxVal.Cmp(currVal) < 0 {
			e.setError(oraError("04009", "MAXVALUE cannot be made to be less than the current value"))
			return
		}
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
		if restartVal, ok := parseSignedBigInt(sequence.Start); ok {
			e.sequenceCurrent[name] = new(big.Int).Sub(restartVal, inc)
		}
	} else if n.Restart {
		e.sequenceCurrent[name] = e.sequenceInitialCurrent(inc, nil, minVal, maxVal, false, hasMin, hasMax)
	}
}

func (e *oracleMetadataExtractor) extractCreatePackage(name string, n *ast.CreatePackageStmt) {
	definition := e.definitionText(n)
	if definition == "" {
		return
	}
	pkg := e.packages[name]
	if pkg == nil || pkg.Definition == "" {
		e.packages[name] = &metadatapb.PackageMetadata{Name: name, Definition: definition}
		return
	}
	spec, body := splitPackageUnits(pkg.Definition)
	if n.IsBody {
		body = definition
	} else {
		spec = definition
	}
	switch {
	case spec != "" && body != "":
		pkg.Definition = spec + "\n\n" + body
	case spec != "":
		pkg.Definition = spec
	default:
		pkg.Definition = body
	}
}

func splitPackageUnits(definition string) (spec, body string) {
	existing := normalizePLSQLPackageDefinition(definition)
	if existing == "" {
		return "", ""
	}
	for _, seg := range oracleparser.Split(existing) {
		if seg.Empty() {
			continue
		}
		text := strings.TrimSpace(seg.Text)
		if isPackageBodyUnit(text) {
			body = text
		} else {
			spec = text
		}
	}
	return spec, body
}

func normalizePLSQLPackageDefinition(definition string) string {
	trimmed := strings.TrimSpace(definition)
	if trimmed == "" {
		return ""
	}
	if words, _ := leadingWords(trimmed, 1); len(words) == 0 || words[0] != "CREATE" {
		trimmed = "CREATE OR REPLACE " + trimmed
	}
	if !strings.HasSuffix(trimmed, ";") {
		trimmed += ";"
	}
	return trimmed
}

func isPackageBodyUnit(unit string) bool {
	words, _ := leadingWords(strings.TrimSpace(unit), 6)
	for i := 0; i+1 < len(words); i++ {
		if words[i] == "PACKAGE" && words[i+1] == "BODY" {
			return true
		}
	}
	return false
}

func removeTriggerDisableClause(body string) string {
	headerEnd := len(body)
	upper := strings.ToUpper(body)
	for _, marker := range []string{"\nBEGIN", "\nDECLARE"} {
		if index := strings.Index(upper, marker); index >= 0 && index < headerEnd {
			headerEnd = index
		}
	}
	for index := 0; index+len("DISABLE") <= headerEnd; index++ {
		if upper[index:index+len("DISABLE")] != "DISABLE" {
			continue
		}
		if index > 0 && isIdentifierByte(body[index-1]) {
			continue
		}
		if index+len("DISABLE") < len(upper) && isIdentifierByte(body[index+len("DISABLE")]) {
			continue
		}
		// The clause follows SQL whitespace, which may be a newline: take the
		// run before the keyword with it and leave what follows in place.
		start := index
		for start > 0 && (body[start-1] == ' ' || body[start-1] == '\t' || body[start-1] == '\n' || body[start-1] == '\r') {
			start--
		}
		if start == index {
			continue
		}
		return body[:start] + body[index+len("DISABLE"):]
	}
	return body
}

func triggerAttributes(n *ast.CreateTriggerStmt) (string, string) {
	var events []string
	for _, item := range listItems(n.Events) {
		switch event := item.(type) {
		case *ast.Integer:
			switch ast.TriggerEvent(event.Ival) {
			case ast.TRIGGER_INSERT:
				events = append(events, "INSERT")
			case ast.TRIGGER_UPDATE:
				events = append(events, "UPDATE")
			case ast.TRIGGER_DELETE:
				events = append(events, "DELETE")
			}
		case *ast.ColumnRef:
			events = append(events, strings.ToUpper(event.Column))
		}
	}
	var timing string
	switch n.Timing {
	case ast.TRIGGER_BEFORE:
		timing = "BEFORE"
	case ast.TRIGGER_AFTER:
		timing = "AFTER"
	case ast.TRIGGER_INSTEAD_OF:
		timing = "INSTEAD OF"
	}
	if n.ForEachRow {
		timing += " EACH ROW"
	} else {
		timing += " STATEMENT"
	}
	return strings.Join(events, " OR "), timing
}

func (e *oracleMetadataExtractor) dropTableProvenance(tableName string) {
	prefix := tableName + "\x00"
	for key := range e.primaryKeyNotNull {
		if strings.HasPrefix(key, prefix) {
			delete(e.primaryKeyNotNull, key)
		}
	}
	for key := range e.namedNotNull {
		if strings.HasPrefix(key, prefix) {
			delete(e.namedNotNull, key)
		}
	}
	for key := range e.syntheticConstraints {
		if strings.HasPrefix(key, prefix) {
			delete(e.syntheticConstraints, key)
		}
	}
	for key := range e.encryptedColumns {
		if strings.HasPrefix(key, prefix) {
			delete(e.encryptedColumns, key)
		}
	}
}

func renameTriggerName(body, oldName, newName string) string {
	upper := strings.ToUpper(body)
	index := strings.Index(upper, "TRIGGER ")
	if index < 0 {
		return body
	}
	start := index + len("TRIGGER ")
	token, tokenStart, tokenEnd := scanTriggerIdentifier(body, start)
	if tokenEnd < len(body) && body[tokenEnd] == '.' {
		token, tokenStart, tokenEnd = scanTriggerIdentifier(body, tokenEnd+1)
	}
	if !strings.EqualFold(token, oldName) {
		return body
	}
	body = body[:tokenStart] + renderRenamedIdentifier(newName) + body[tokenEnd:]
	body = strings.ReplaceAll(body, "ALTER TRIGGER \""+oldName+"\" DISABLE;", "ALTER TRIGGER \""+newName+"\" DISABLE;")
	return body
}

func (e *oracleMetadataExtractor) extractAlterTrigger(n *ast.AlterTriggerStmt) {
	name := objectName(n.Name)
	if name == "" || !e.useSchema(n.Name.Schema) {
		return
	}
	trigger := e.triggers[name]
	if trigger == nil {
		return
	}
	switch strings.ToUpper(n.Action) {
	case "RENAME":
		if n.NewName == "" {
			return
		}
		delete(e.triggers, name)
		trigger.Name = n.NewName
		trigger.Body = renameTriggerName(trigger.Body, name, n.NewName)
		e.triggers[n.NewName] = trigger
	case "DISABLE":
		marker := "ALTER TRIGGER \"" + name + "\" DISABLE;"
		if !strings.Contains(strings.ToUpper(trigger.Body), strings.ToUpper(marker)) {
			trigger.Body = strings.TrimSpace(trigger.Body) + "\n" + marker
		}
	case "ENABLE":
		marker := strings.ToUpper("ALTER TRIGGER \"" + name + "\" DISABLE;")
		var lines []string
		for _, line := range strings.Split(trigger.Body, "\n") {
			if strings.TrimSpace(strings.ToUpper(line)) == marker {
				continue
			}
			lines = append(lines, line)
		}
		trigger.Body = strings.Join(lines, "\n")
	}
}

func (e *oracleMetadataExtractor) extractCreateTrigger(n *ast.CreateTriggerStmt) {
	triggerName := objectName(n.Name)
	tableName := objectName(n.Table)
	if triggerName == "" || tableName == "" {
		return
	}
	if !e.useSchema(n.Name.Schema) || !e.useSchema(n.Table.Schema) {
		return
	}
	if strings.EqualFold(tableName, "DATABASE") || strings.EqualFold(tableName, "SCHEMA") {
		if e.err == nil {
			e.err = fmt.Errorf("non-table triggers are not supported by oracle/catalog")
		}
		return
	}
	if n.OrReplace {
		e.dropTrigger(triggerName)
	}
	event, timing := triggerAttributes(n)
	body := e.definitionText(n)
	if !n.Enable {
		body = removeTriggerDisableClause(body)
		body = strings.TrimSpace(body) + "\nALTER TRIGGER \"" + triggerName + "\" DISABLE;"
	}
	trigger := &metadatapb.TriggerMetadata{
		Name:   triggerName,
		Event:  event,
		Timing: timing,
		Body:   body,
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
	if !e.checkSchema(n.Object.Schema, "00942") {
		return
	}
	rawName := objectName(n.Object)
	name := e.resolveSynonymTarget(rawName, n.Object.Schema)
	if e.sequences[name] != nil {
		e.setError(oraError("02201", "sequence not allowed here"))
		return
	}
	if n.Column != "" {
		if table := e.tables[name]; table != nil {
			col := tableFindColumn(table, n.Column)
			if col == nil {
				e.setError(oraError("00904", "%q: invalid identifier", n.Column))
				return
			}
			col.Comment = n.Comment
			return
		}
		if view := e.views[name]; view != nil {
			if !e.viewHasColumn(view, n.Column) {
				e.setError(oraError("00904", "%q: invalid identifier", n.Column))
				return
			}
			for _, col := range view.Columns {
				if col.Name == n.Column {
					col.Comment = n.Comment
					return
				}
			}
			return
		}
		if e.materializedViews[name] != nil {
			return
		}
		e.setError(oraError("00942", "table or view %q does not exist", rawName))
		return
	}

	switch n.ObjectType {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW:
		if table := e.tables[name]; table != nil {
			table.Comment = n.Comment
			return
		}
		if view := e.views[name]; view != nil {
			view.Comment = n.Comment
			return
		}
		if e.materializedViews[name] != nil {
			return
		}
		e.setError(oraError("00942", "table or view %q does not exist", rawName))
	case ast.OBJECT_MATERIALIZED_VIEW:
		if materializedView := e.materializedViews[name]; materializedView != nil {
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
	if match := functionCreatePrefix.FindStringIndex(definition); match != nil {
		definition = strings.TrimSpace(definition[match[1]:])
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

func hasPrimaryIndex(table *metadatapb.TableMetadata) bool {
	for _, index := range table.Indexes {
		if index.Primary {
			return true
		}
		for _, c := range index.Constraints {
			if c.Primary {
				return true
			}
		}
	}
	return false
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
	e.owner = e.currentSchema
	if e.owner == "" {
		e.owner = meta.GetName()
	}
	for _, table := range schema.GetTables() {
		cloned := cloneTable(table)
		e.tables[table.GetName()] = cloned
		for _, trigger := range cloned.GetTriggers() {
			e.triggers[trigger.GetName()] = trigger
		}
	}
	for _, view := range schema.GetViews() {
		cloned := cloneView(view)
		e.views[view.GetName()] = cloned
		for _, trigger := range cloned.GetTriggers() {
			e.triggers[trigger.GetName()] = trigger
		}
	}
	for _, view := range schema.GetMaterializedViews() {
		cloned := cloneMaterializedView(view)
		e.materializedViews[view.GetName()] = cloned
		for _, trigger := range cloned.GetTriggers() {
			e.triggers[trigger.GetName()] = trigger
		}
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
	for _, syn := range schema.GetSynonyms() {
		e.synonyms[syn.GetName()] = cloneSynonym(syn)
	}
	for _, syn := range meta.GetPublicSynonyms() {
		e.publicSynonyms[syn.GetName()] = cloneSynonym(syn)
	}
	return e
}

var (
	functionCreatePrefix = regexp.MustCompile(`(?is)^\s*CREATE(?:\s+OR\s+REPLACE)?(?:\s+(?:NON)?EDITIONABLE)?\s+`)
)

func applyIdentity(column *metadatapb.ColumnMetadata, identity *ast.IdentityClause, text string) {
	column.IsIdentity = true
	if identity != nil && !identity.Always {
		column.IdentityGeneration = metadatapb.ColumnMetadata_BY_DEFAULT
		column.DefaultOnNull = identityTextHasOnNull(text)
	} else {
		column.IdentityGeneration = metadatapb.ColumnMetadata_ALWAYS
	}
	applyIdentitySequenceOptionsFromText(column, text)
}

func identityTextHasOnNull(text string) bool {
	lexer := oracleparser.NewLexer(text)
	var prevQuoted bool
	var prevStr string
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			return false
		}
		quoted := tok.Loc >= 0 && tok.Loc < len(text) && text[tok.Loc] == '"'
		if !quoted && strings.EqualFold(tok.Str, "IDENTITY") {
			return false
		}
		if !prevQuoted && !quoted && strings.EqualFold(prevStr, "ON") && strings.EqualFold(tok.Str, "NULL") {
			return true
		}
		prevQuoted = quoted
		prevStr = tok.Str
	}
}

func (e *oracleMetadataExtractor) extractCreateSynonym(n *ast.CreateSynonymStmt) {
	name := objectName(n.Name)
	if name == "" || n.Target == nil || objectName(n.Target) == "" {
		return
	}
	if e.serverConfig.NonApplicationRoot && isApplicationSharing(n.Sharing) {
		e.setError(oraError("65021", "illegal use of SHARING clause"))
		return
	}
	if !n.Public && !e.checkSchema(n.Name.Schema, "01031") {
		return
	}
	if !n.Public {
		targetName := objectName(n.Target)
		sameSchema := n.Target.Schema == "" || n.Target.Schema == e.currentSchema || (e.currentSchema == "" && n.Target.Schema == e.owner)
		if n.Target.DBLink == "" && sameSchema && targetName == name {
			e.setError(oraError("01471", "cannot create a synonym %q with same name as object", name))
			return
		}
		if n.IfNotExists && e.synonyms[name] != nil {
			return
		}
		if n.OrReplace {
			if e.nonSynonymObjectExists(name) {
				e.setError(oraError("00955", "name %q is already used by an existing object", name))
				return
			}
		} else if e.objectExistsInNamespace(name) {
			e.setError(oraError("00955", "name %q is already used by an existing object", name))
			return
		}
	}
	targetSchema := n.Target.Schema
	if targetSchema == "" && n.Target.DBLink == "" && e.currentSchema != "" {
		targetSchema = e.currentSchema
	}
	syn := &metadatapb.SynonymMetadata{
		Name:         name,
		TargetSchema: targetSchema,
		TargetName:   objectName(n.Target),
		TargetDbLink: n.Target.DBLink,
	}
	if n.Public {
		if n.IfNotExists && e.publicSynonyms[name] != nil {
			return
		}
		e.publicSynonyms[name] = syn
	} else {
		e.synonyms[name] = syn
	}
}

func isDropPublicSynonym(stmtText string) bool {
	lexer := oracleparser.NewLexer(stmtText)
	tok1 := lexer.NextToken()
	tok2 := lexer.NextToken()
	return strings.EqualFold(tok1.Str, "DROP") && strings.EqualFold(tok2.Str, "PUBLIC")
}

func viewInvisibleAliasPositions(schemaText string, n *ast.CreateViewStmt) map[int]bool {
	out := map[int]bool{}
	if n == nil || n.Name == nil || n.Query == nil {
		return out
	}
	start := n.Name.Loc.End
	end := n.Loc.End
	if qloc := ast.NodeLoc(n.Query); qloc.Start > start && qloc.Start <= len(schemaText) {
		end = qloc.Start
	}
	if start < 0 || end <= start || end > len(schemaText) {
		return out
	}
	header := schemaText[start:end]
	openIdx := strings.IndexByte(header, '(')
	if openIdx < 0 {
		return out
	}
	aliasHeader := header[openIdx:]
	lexer := oracleparser.NewLexer(aliasHeader)
	depth := 0
	itemIdx := 0
	itemTokIdx := 0
	prevWasConstraint := false
	started := false
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			break
		}
		if tok.Type == '(' {
			depth++
			started = true
			continue
		}
		if tok.Type == ')' {
			depth--
			if started && depth == 0 {
				break
			}
			continue
		}
		if depth == 1 {
			if tok.Type == ',' {
				itemIdx++
				itemTokIdx = 0
				prevWasConstraint = false
				continue
			}
			isQuoted := tok.Loc >= 0 && tok.Loc < len(aliasHeader) && aliasHeader[tok.Loc] == '"'
			if itemTokIdx == 1 && !isQuoted && !prevWasConstraint && strings.EqualFold(tok.Str, "INVISIBLE") {
				out[itemIdx] = true
			}
			prevWasConstraint = !isQuoted && strings.EqualFold(tok.Str, "CONSTRAINT")
			itemTokIdx++
		}
	}
	return out
}

func formatPartitionKeyIdentifier(name string) string {
	if isSimpleIdentifier(name) {
		return name
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func partitionTypeFromAST(t ast.PartitionType) metadatapb.TablePartitionMetadata_Type {
	switch t {
	case ast.PARTITION_RANGE:
		return metadatapb.TablePartitionMetadata_RANGE
	case ast.PARTITION_LIST:
		return metadatapb.TablePartitionMetadata_LIST
	case ast.PARTITION_HASH:
		return metadatapb.TablePartitionMetadata_HASH
	default:
		return metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED
	}
}

func (e *oracleMetadataExtractor) extractPartitionKeyExpr(cols *ast.List) string {
	var keyCols []string
	for _, item := range listItems(cols) {
		if colRef, ok := item.(*ast.ColumnRef); ok && colRef.Column != "" {
			keyCols = append(keyCols, formatPartitionKeyIdentifier(colRef.Column))
		} else if txt := strings.TrimSpace(e.nodeText(item)); txt != "" {
			keyCols = append(keyCols, txt)
		}
	}
	return strings.Join(keyCols, ", ")
}

func (e *oracleMetadataExtractor) extractPartitionValues(values *ast.List) string {
	var vals []string
	for _, val := range listItems(values) {
		if txt := strings.TrimSpace(e.nodeText(val)); txt != "" {
			vals = append(vals, txt)
		}
	}
	return strings.Join(vals, ", ")
}

func (e *oracleMetadataExtractor) extractTablePartitions(clause *ast.PartitionClause) []*metadatapb.TablePartitionMetadata {
	if clause == nil {
		return nil
	}
	partType := partitionTypeFromAST(clause.Type)
	if partType == metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED {
		return nil
	}
	keyExpr := e.extractPartitionKeyExpr(clause.Columns)
	intervalExpr := strings.TrimSpace(e.nodeText(clause.Interval))
	var subType metadatapb.TablePartitionMetadata_Type
	var subKeyExpr string
	var templateSubs []ast.Node
	if clause.Subpartition != nil {
		subType = partitionTypeFromAST(clause.Subpartition.Type)
		subKeyExpr = e.extractPartitionKeyExpr(clause.Subpartition.Columns)
		templateSubs = listItems(clause.Subpartition.Partitions)
	}
	var partitions []*metadatapb.TablePartitionMetadata
	for _, item := range listItems(clause.Partitions) {
		def, ok := item.(*ast.PartitionDef)
		if !ok {
			continue
		}
		pm := &metadatapb.TablePartitionMetadata{
			Name:           def.Name,
			Type:           partType,
			Expression:     keyExpr,
			Value:          e.extractPartitionValues(def.Values),
			IntervalClause: intervalExpr,
		}
		if subType != metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED {
			subItems := listItems(def.Subpartitions)
			fromTemplate := false
			if len(subItems) == 0 && len(templateSubs) > 0 {
				subItems = templateSubs
				fromTemplate = true
			}
			for _, subItem := range subItems {
				subDef, ok := subItem.(*ast.PartitionDef)
				if !ok {
					continue
				}
				subName := subDef.Name
				if fromTemplate && def.Name != "" && subName != "" {
					subName = def.Name + "_" + subName
				}
				pm.Subpartitions = append(pm.Subpartitions, &metadatapb.TablePartitionMetadata{
					Name:       subName,
					Type:       subType,
					Expression: subKeyExpr,
					Value:      e.extractPartitionValues(subDef.Values),
				})
			}
		}
		partitions = append(partitions, pm)
	}
	return partitions
}

type viewDependencyVisitor struct {
	schema  string
	seen    map[string]bool
	scopes  []map[string]bool
	markers []bool
	result  *[]*metadatapb.DependencyColumn
}

func (v *viewDependencyVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		if len(v.markers) > 0 {
			marker := v.markers[len(v.markers)-1]
			v.markers = v.markers[:len(v.markers)-1]
			if marker && len(v.scopes) > 0 {
				v.scopes = v.scopes[:len(v.scopes)-1]
			}
		}
		return v
	}
	marker := false
	if selectStmt, ok := node.(*ast.SelectStmt); ok {
		scope := map[string]bool{}
		if len(v.scopes) > 0 {
			for name := range v.scopes[len(v.scopes)-1] {
				scope[name] = true
			}
		}
		if selectStmt.WithClause != nil {
			for _, item := range listItems(selectStmt.WithClause.CTEs) {
				if cte, ok := item.(*ast.CTE); ok && cte.Name != "" {
					scope[cte.Name] = true
				}
			}
		}
		v.scopes = append(v.scopes, scope)
		marker = true
	}
	v.markers = append(v.markers, marker)
	// A derived analytic view names a persistent analytic view, in a FROM
	// clause and in a WITH element: the change depends on it.
	if clause, ok := node.(*ast.SubAVClause); ok && clause.Using != nil && clause.Using.Name != "" {
		schema := clause.Using.Schema
		if schema == "" {
			schema = v.schema
		}
		key := schema + "\x00" + clause.Using.Name
		if !v.seen[key] {
			v.seen[key] = true
			*v.result = append(*v.result, &metadatapb.DependencyColumn{Schema: schema, Table: clause.Using.Name, Column: "*"})
		}
	}
	if ref, ok := node.(*ast.TableRef); ok && ref.Name != nil && ref.Name.Name != "" {
		inScope := len(v.scopes) > 0 && v.scopes[len(v.scopes)-1][ref.Name.Name]
		if ref.Name.Schema == "" && inScope {
			return v
		}
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
