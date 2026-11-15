package analysis

import (
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// Conditions CheckSDLIntegrity reports, besides SDLSyntaxError. Schemas are
// "public" where the SQL names none. Start is the later definition or the
// offending object; First is the earlier definition.
const (
	// SDLDuplicateTable: a file creates a table twice. Schema, Table,
	// First.Start.
	SDLDuplicateTable SDLCode = "DUPLICATE_TABLE"
	// SDLDuplicateColumn: a table defines a column twice. Schema, Table,
	// Column.
	SDLDuplicateColumn SDLCode = "DUPLICATE_COLUMN"
	// SDLDuplicateConstraint: a file reuses a PRIMARY KEY or UNIQUE name within
	// a schema, or a CHECK or FOREIGN KEY name within a table. Name,
	// ConstraintType, Schema, Table, First.ConstraintType, First.Table,
	// First.Start.
	SDLDuplicateConstraint SDLCode = "DUPLICATE_CONSTRAINT"
	// SDLDuplicateIndex: a file creates an index name twice within a schema.
	// Schema, Name, First.Start.
	SDLDuplicateIndex SDLCode = "DUPLICATE_INDEX"
	// SDLDuplicateTableAcrossFiles: an earlier file creates the same table.
	// Schema, Table, First.Path, First.Start.
	SDLDuplicateTableAcrossFiles SDLCode = "DUPLICATE_TABLE_ACROSS_FILES"
	// SDLDuplicateIndexAcrossFiles: an earlier file creates the same index.
	// Schema, Name, First.Path, First.Start.
	SDLDuplicateIndexAcrossFiles SDLCode = "DUPLICATE_INDEX_ACROSS_FILES"
	// SDLDuplicateViewAcrossFiles: an earlier file creates the same view.
	// Schema, Name, First.Path, First.Start.
	SDLDuplicateViewAcrossFiles SDLCode = "DUPLICATE_VIEW_ACROSS_FILES"
	// SDLDuplicateConstraintAcrossFiles: an earlier file uses the same PRIMARY
	// KEY or UNIQUE name in the schema. Name, Schema, ConstraintType, Table,
	// First.Path, First.ConstraintType, First.Table, First.Start.
	SDLDuplicateConstraintAcrossFiles SDLCode = "DUPLICATE_CONSTRAINT_ACROSS_FILES"
	// SDLMultiplePrimaryKeys: a table has more than one table-level PRIMARY
	// KEY. Schema, Table, PrimaryKeys (Name, Start); Start is the second one's.
	SDLMultiplePrimaryKeys SDLCode = "MULTIPLE_PRIMARY_KEYS"
	// SDLForeignKeyTableNotFound: no file creates the referenced table. Name
	// (the constraint, empty when unnamed), Schema, Table, RefSchema, RefTable,
	// Columns (the key's own).
	SDLForeignKeyTableNotFound SDLCode = "FOREIGN_KEY_TABLE_NOT_FOUND"
	// SDLForeignKeyColumnNotFound: the referenced table lacks a referenced
	// column. Name, Schema, Table, Column (the key's own), RefSchema, RefTable,
	// RefColumn, Columns (the referenced table's, in declaration order).
	SDLForeignKeyColumnNotFound SDLCode = "FOREIGN_KEY_COLUMN_NOT_FOUND"
	// SDLViewDependencyNotFound: a view reads a relation no file creates.
	// Schema, Name (the view), RefSchema, RefTable.
	SDLViewDependencyNotFound SDLCode = "VIEW_DEPENDENCY_NOT_FOUND"
)

// SDLFile is one file of a declarative schema project.
type SDLFile struct {
	Path string
	Text string
}

// CheckSDLIntegrity checks the files of one declarative schema project
// together and returns each file's findings, keyed by Path. Paths must be
// distinct.
//
// It judges top-level CREATE TABLE, CREATE INDEX and CREATE VIEW. A file's
// findings are its cross-file duplicates, then its own duplicates in statement
// order, then multiple primary keys, unresolved foreign keys, and unresolved
// view dependencies. "Earlier" means earlier in files.
//
// What is judged, as plain facts:
//   - Only table-level constraints (TableElts and Constraints) are tracked;
//     column-level PRIMARY KEY, UNIQUE and REFERENCES are not.
//   - A duplicate table's second definition is dropped whole. A view defined
//     twice in one file keeps its last definition and is not reported.
//   - A foreign key resolves against the first file's definition of its table.
//     Keys whose column counts differ, and own columns the table lacks, are
//     not checked further.
//   - A view's dependencies are the RangeVars anywhere in its query, minus the
//     names of the top-level WITH's CTEs, resolved against the view's schema
//     when unqualified; system schemas, and system table and view names in
//     any schema, are skipped.
//
// If a file fails parseSDL, the result holds only that file's SDLSyntaxError.
func CheckSDLIntegrity(files []SDLFile) map[string][]SDLFinding {
	if len(files) == 0 {
		return make(map[string][]SDLFinding)
	}

	symbols := make([]*sdlSymbols, 0, len(files))
	for _, file := range files {
		stmts, syntaxError := parseSDL(file.Text)
		if syntaxError != nil {
			return map[string][]SDLFinding{file.Path: {*syntaxError}}
		}
		symbols = append(symbols, buildSDLSymbols(file, stmts))
	}

	merged := mergeSDLSymbols(symbols)
	crossFile := sdlCrossFileDuplicates(symbols)

	results := make(map[string][]SDLFinding, len(files))
	for i, st := range symbols {
		findings := append([]SDLFinding{}, crossFile[i]...)
		findings = append(findings, st.findings...)
		findings = st.appendMultiplePrimaryKeys(findings)
		findings = st.appendForeignKeys(findings, merged)
		findings = st.appendViewDependencies(findings, merged)
		results[st.path] = findings
	}
	return results
}

type sdlSymbols struct {
	path string

	tables     map[string]*sdlTable
	tableOrder []*sdlTable
	indexes    map[string]*sdlIndex
	indexOrder []*sdlIndex
	views      map[string]*sdlView
	viewOrder  []string
	// schemaConstraints holds named PRIMARY KEY and UNIQUE constraints by
	// schema, unique within it; tableConstraints holds named CHECK and
	// FOREIGN KEY constraints by qualified table, unique within it.
	schemaConstraints     map[string]map[string]*sdlConstraint
	schemaConstraintOrder []*sdlConstraint
	tableConstraints      map[string]map[string]*sdlConstraint

	findings []SDLFinding
}

type sdlTable struct {
	schema      string
	name        string
	columns     map[string]bool
	columnOrder []string
	primaryKeys []*sdlConstraint
	constraints []*sdlConstraint
	start       int
}

type sdlIndex struct {
	schema string
	name   string
	start  int
}

type sdlView struct {
	schema string
	name   string
	query  ast.Node
	start  int
}

type sdlConstraint struct {
	name           string
	constraintType string
	schema         string
	table          string
	start          int
	refSchema      string
	refTable       string
	columns        []string
	refColumns     []string
}

type sdlMerged struct {
	tables map[string]*sdlTable
	views  map[string]bool
}

func sdlQualifiedName(schema, name string) string {
	if schema == "" {
		schema = "public"
	}
	return schema + "." + name
}

func sdlSchemaOrPublic(schema string) string {
	if schema == "" {
		return "public"
	}
	return schema
}

func buildSDLSymbols(file SDLFile, stmts []omnipg.Statement) *sdlSymbols {
	st := &sdlSymbols{
		path:              file.Path,
		tables:            make(map[string]*sdlTable),
		indexes:           make(map[string]*sdlIndex),
		views:             make(map[string]*sdlView),
		schemaConstraints: make(map[string]map[string]*sdlConstraint),
		tableConstraints:  make(map[string]map[string]*sdlConstraint),
	}
	for _, stmt := range stmts {
		switch n := stmt.AST.(type) {
		case *ast.CreateStmt:
			st.addCreateStmt(n)
		case *ast.IndexStmt:
			st.addIndexStmt(n)
		case *ast.ViewStmt:
			st.addViewStmt(n)
		default:
		}
	}
	return st
}

func (st *sdlSymbols) addCreateStmt(n *ast.CreateStmt) {
	if n.Relation == nil {
		return
	}
	tableName := n.Relation.Relname
	schemaName := sdlSchemaOrPublic(n.Relation.Schemaname)
	start := n.Loc.Start

	key := sdlQualifiedName(schemaName, tableName)
	if existing := st.tables[key]; existing != nil {
		st.findings = append(st.findings, SDLFinding{
			Code:   SDLDuplicateTable,
			Start:  start,
			Schema: schemaName,
			Table:  tableName,
			First:  &SDLDefinition{Start: existing.start},
		})
		return
	}
	table := &sdlTable{
		schema:  schemaName,
		name:    tableName,
		columns: make(map[string]bool),
		start:   start,
	}
	st.tables[key] = table
	st.tableOrder = append(st.tableOrder, table)

	if n.TableElts != nil {
		for _, item := range n.TableElts.Items {
			switch elem := item.(type) {
			case *ast.ColumnDef:
				st.addColumn(table, elem)
			case *ast.Constraint:
				st.addConstraint(table, elem)
			default:
			}
		}
	}
	if n.Constraints != nil {
		for _, item := range n.Constraints.Items {
			if c, ok := item.(*ast.Constraint); ok {
				st.addConstraint(table, c)
			}
		}
	}
}

func (st *sdlSymbols) addColumn(table *sdlTable, col *ast.ColumnDef) {
	if table.columns[col.Colname] {
		st.findings = append(st.findings, SDLFinding{
			Code:   SDLDuplicateColumn,
			Start:  col.Loc.Start,
			Schema: table.schema,
			Table:  table.name,
			Column: col.Colname,
		})
		return
	}
	table.columns[col.Colname] = true
	table.columnOrder = append(table.columnOrder, col.Colname)
}

func (st *sdlSymbols) addConstraint(table *sdlTable, c *ast.Constraint) {
	constraint := &sdlConstraint{
		name:   c.Conname,
		schema: table.schema,
		table:  table.name,
		start:  c.Loc.Start,
	}
	switch c.Contype {
	case ast.CONSTR_PRIMARY:
		constraint.constraintType = "PRIMARY KEY"
		table.primaryKeys = append(table.primaryKeys, constraint)
	case ast.CONSTR_UNIQUE:
		constraint.constraintType = "UNIQUE"
	case ast.CONSTR_CHECK:
		constraint.constraintType = "CHECK"
	case ast.CONSTR_FOREIGN:
		constraint.constraintType = "FOREIGN KEY"
		constraint.columns = sdlListStrings(c.FkAttrs)
		if c.Pktable != nil {
			constraint.refTable = c.Pktable.Relname
			constraint.refSchema = sdlSchemaOrPublic(c.Pktable.Schemaname)
		}
		constraint.refColumns = sdlListStrings(c.PkAttrs)
	default:
		return
	}

	if constraint.name != "" {
		var names map[string]*sdlConstraint
		if constraint.isSchemaScoped() {
			if st.schemaConstraints[table.schema] == nil {
				st.schemaConstraints[table.schema] = make(map[string]*sdlConstraint)
			}
			names = st.schemaConstraints[table.schema]
		} else {
			tableKey := sdlQualifiedName(table.schema, table.name)
			if st.tableConstraints[tableKey] == nil {
				st.tableConstraints[tableKey] = make(map[string]*sdlConstraint)
			}
			names = st.tableConstraints[tableKey]
		}
		if existing, ok := names[constraint.name]; ok {
			st.findings = append(st.findings, SDLFinding{
				Code:           SDLDuplicateConstraint,
				Start:          constraint.start,
				Schema:         table.schema,
				Table:          table.name,
				Name:           constraint.name,
				ConstraintType: constraint.constraintType,
				First: &SDLDefinition{
					Table:          existing.table,
					ConstraintType: existing.constraintType,
					Start:          existing.start,
				},
			})
		} else {
			names[constraint.name] = constraint
			if constraint.isSchemaScoped() {
				st.schemaConstraintOrder = append(st.schemaConstraintOrder, constraint)
			}
		}
	}

	table.constraints = append(table.constraints, constraint)
}

// isSchemaScoped reports whether PostgreSQL requires the constraint's name to
// be unique within its schema rather than its table.
func (c *sdlConstraint) isSchemaScoped() bool {
	return c.constraintType == "PRIMARY KEY" || c.constraintType == "UNIQUE"
}

func (st *sdlSymbols) addIndexStmt(n *ast.IndexStmt) {
	var schemaName string
	if n.Relation != nil {
		schemaName = n.Relation.Schemaname
	}
	schemaName = sdlSchemaOrPublic(schemaName)
	if n.Idxname == "" {
		return
	}
	start := n.Loc.Start

	key := sdlQualifiedName(schemaName, n.Idxname)
	if existing, ok := st.indexes[key]; ok {
		st.findings = append(st.findings, SDLFinding{
			Code:   SDLDuplicateIndex,
			Start:  start,
			Schema: schemaName,
			Name:   n.Idxname,
			First:  &SDLDefinition{Start: existing.start},
		})
		return
	}
	index := &sdlIndex{schema: schemaName, name: n.Idxname, start: start}
	st.indexes[key] = index
	st.indexOrder = append(st.indexOrder, index)
}

func (st *sdlSymbols) addViewStmt(n *ast.ViewStmt) {
	if n.View == nil {
		return
	}
	schemaName := sdlSchemaOrPublic(n.View.Schemaname)
	key := sdlQualifiedName(schemaName, n.View.Relname)
	if _, ok := st.views[key]; !ok {
		st.viewOrder = append(st.viewOrder, key)
	}
	st.views[key] = &sdlView{
		schema: schemaName,
		name:   n.View.Relname,
		query:  n.Query,
		start:  n.Loc.Start,
	}
}

func mergeSDLSymbols(symbols []*sdlSymbols) *sdlMerged {
	merged := &sdlMerged{
		tables: make(map[string]*sdlTable),
		views:  make(map[string]bool),
	}
	for _, st := range symbols {
		for key, table := range st.tables {
			if _, ok := merged.tables[key]; !ok {
				merged.tables[key] = table
			}
		}
		for key := range st.views {
			merged.views[key] = true
		}
	}
	return merged
}

// sdlCrossFileDuplicates returns, per file index, the tables, indexes, views
// and schema-scoped constraint names an earlier file already defines.
func sdlCrossFileDuplicates(symbols []*sdlSymbols) map[int][]SDLFinding {
	findings := make(map[int][]SDLFinding)
	seenTables := make(map[string]SDLDefinition)
	seenIndexes := make(map[string]SDLDefinition)
	seenViews := make(map[string]SDLDefinition)
	seenConstraints := make(map[string]SDLDefinition)

	for i, st := range symbols {
		for _, table := range st.tableOrder {
			key := sdlQualifiedName(table.schema, table.name)
			if first, ok := seenTables[key]; ok {
				findings[i] = append(findings[i], SDLFinding{
					Code:   SDLDuplicateTableAcrossFiles,
					Start:  table.start,
					Schema: table.schema,
					Table:  table.name,
					First:  &first,
				})
			} else {
				seenTables[key] = SDLDefinition{Path: st.path, Start: table.start}
			}
		}
		for _, index := range st.indexOrder {
			key := sdlQualifiedName(index.schema, index.name)
			if first, ok := seenIndexes[key]; ok {
				findings[i] = append(findings[i], SDLFinding{
					Code:   SDLDuplicateIndexAcrossFiles,
					Start:  index.start,
					Schema: index.schema,
					Name:   index.name,
					First:  &first,
				})
			} else {
				seenIndexes[key] = SDLDefinition{Path: st.path, Start: index.start}
			}
		}
		for _, key := range st.viewOrder {
			view := st.views[key]
			if first, ok := seenViews[key]; ok {
				findings[i] = append(findings[i], SDLFinding{
					Code:   SDLDuplicateViewAcrossFiles,
					Start:  view.start,
					Schema: view.schema,
					Name:   view.name,
					First:  &first,
				})
			} else {
				seenViews[key] = SDLDefinition{Path: st.path, Start: view.start}
			}
		}
		for _, constraint := range st.schemaConstraintOrder {
			key := constraint.schema + "\x00" + constraint.name
			if first, ok := seenConstraints[key]; ok {
				findings[i] = append(findings[i], SDLFinding{
					Code:           SDLDuplicateConstraintAcrossFiles,
					Start:          constraint.start,
					Schema:         constraint.schema,
					Table:          constraint.table,
					Name:           constraint.name,
					ConstraintType: constraint.constraintType,
					First:          &first,
				})
			} else {
				seenConstraints[key] = SDLDefinition{
					Path:           st.path,
					Table:          constraint.table,
					ConstraintType: constraint.constraintType,
					Start:          constraint.start,
				}
			}
		}
	}
	return findings
}

func (st *sdlSymbols) appendMultiplePrimaryKeys(findings []SDLFinding) []SDLFinding {
	for _, table := range st.tableOrder {
		if len(table.primaryKeys) <= 1 {
			continue
		}
		keys := make([]SDLDefinition, 0, len(table.primaryKeys))
		for _, pk := range table.primaryKeys {
			keys = append(keys, SDLDefinition{Name: pk.name, Start: pk.start})
		}
		findings = append(findings, SDLFinding{
			Code:        SDLMultiplePrimaryKeys,
			Start:       table.primaryKeys[1].start,
			Schema:      table.schema,
			Table:       table.name,
			PrimaryKeys: keys,
		})
	}
	return findings
}

func (st *sdlSymbols) appendForeignKeys(findings []SDLFinding, merged *sdlMerged) []SDLFinding {
	for _, table := range st.tableOrder {
		for _, constraint := range table.constraints {
			if constraint.constraintType != "FOREIGN KEY" {
				continue
			}
			refTable := merged.tables[sdlQualifiedName(constraint.refSchema, constraint.refTable)]
			if refTable == nil {
				findings = append(findings, SDLFinding{
					Code:      SDLForeignKeyTableNotFound,
					Start:     constraint.start,
					Schema:    table.schema,
					Table:     table.name,
					Name:      constraint.name,
					RefSchema: constraint.refSchema,
					RefTable:  constraint.refTable,
					Columns:   constraint.columns,
				})
				continue
			}
			if len(constraint.columns) != len(constraint.refColumns) {
				continue
			}
			for i, column := range constraint.columns {
				refColumn := constraint.refColumns[i]
				if !table.columns[column] || refTable.columns[refColumn] {
					continue
				}
				findings = append(findings, SDLFinding{
					Code:      SDLForeignKeyColumnNotFound,
					Start:     constraint.start,
					Schema:    table.schema,
					Table:     table.name,
					Name:      constraint.name,
					Column:    column,
					RefSchema: constraint.refSchema,
					RefTable:  constraint.refTable,
					RefColumn: refColumn,
					Columns:   refTable.columnOrder,
				})
			}
		}
	}
	return findings
}

func (st *sdlSymbols) appendViewDependencies(findings []SDLFinding, merged *sdlMerged) []SDLFinding {
	for _, key := range st.viewOrder {
		view := st.views[key]
		if view.query == nil {
			continue
		}
		for _, dep := range sdlViewDependencies(view.query, view.schema) {
			if IsSystemSchema(dep.schema) || IsSystemTable(dep.table) || IsSystemView(dep.table) {
				continue
			}
			refKey := sdlQualifiedName(dep.schema, dep.table)
			if merged.tables[refKey] != nil || merged.views[refKey] {
				continue
			}
			findings = append(findings, SDLFinding{
				Code:      SDLViewDependencyNotFound,
				Start:     view.start,
				Schema:    view.schema,
				Name:      view.name,
				RefSchema: dep.schema,
				RefTable:  dep.table,
			})
		}
	}
	return findings
}

type sdlRelation struct {
	schema string
	table  string
}

// sdlViewDependencies returns the distinct RangeVars in query, in walk order,
// except those naming a CTE of query's own WITH clause.
func sdlViewDependencies(query ast.Node, defaultSchema string) []sdlRelation {
	cteNames := make(map[string]bool)
	if sel, ok := query.(*ast.SelectStmt); ok && sel != nil && sel.WithClause != nil && sel.WithClause.Ctes != nil {
		for _, item := range sel.WithClause.Ctes.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok {
				cteNames[cte.Ctename] = true
			}
		}
	}

	var refs []sdlRelation
	seen := make(map[string]bool)
	ast.Inspect(query, func(n ast.Node) bool {
		rv, ok := n.(*ast.RangeVar)
		if !ok || rv.Relname == "" || cteNames[rv.Relname] {
			return true
		}
		schema := rv.Schemaname
		if schema == "" {
			schema = defaultSchema
		}
		key := schema + "." + rv.Relname
		if !seen[key] {
			seen[key] = true
			refs = append(refs, sdlRelation{schema: schema, table: rv.Relname})
		}
		return true
	})
	return refs
}
