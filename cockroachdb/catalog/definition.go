// Package catalog renders CockroachDB schema definitions from a synced
// metadata snapshot and generates the migration between two snapshots.
//
// It began as a copy of the parts of pg/catalog that render definitions, and
// differs where CockroachDB does. A primary key is part of CREATE TABLE,
// because CockroachDB gives a table created without one a hidden rowid key. A
// secondary index renders from its synced definition through CockroachDB's own
// formatter, which keeps STORING, hash sharding, inverted and partial indexes.
// The columns CockroachDB adds itself (the hidden rowid, a hash-sharded
// index's shard column) are left out. CockroachDB stores the relations of a
// view, an index and a SQL routine body under their database's name, which the
// rendering drops so a definition replays on a database of any name.
// CockroachDB holds no comment on a view, materialized view, sequence or
// routine, so none is written.
//
// CockroachDB's syntax is parsed with its own parser, never pg's.
package catalog

import (
	"fmt"

	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

type metadataIndex = metadatapb.IndexMetadata

const (
	header = `
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

`

	setDefaultTableSpace = "SET default_tablespace = '';\n\n"
)

// DefinitionOptions selects the form DatabaseDefinition renders.
type DefinitionOptions struct {
	// PrintHeader prefixes the restore script with session settings and sets
	// the default tablespace before the first table.
	PrintHeader bool
	// SDLFormat renders the declarative form instead of the restore script.
	// It ignores PrintHeader.
	SDLFormat bool
}

// DatabaseDefinition renders a synced database snapshot as a restore script,
// or with opts.SDLFormat as SDL.
func DatabaseDefinition(metadata *metadatapb.DatabaseSchemaMetadata, opts DefinitionOptions) (string, error) {
	if len(metadata.GetSchemas()) == 0 {
		return "", nil
	}
	if opts.SDLFormat {
		return databaseSDL(metadata), nil
	}
	r := newRenderer(metadata.Name, metadata.Schemas)
	var buf strings.Builder
	if opts.PrintHeader {
		buf.WriteString(header)
	}
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		writeSchema(&buf, schema)
		writeSchemaOwnerAndComment(&buf, schema)
	}
	for _, extension := range metadata.Extensions {
		writeExtension(&buf, extension)
	}
	r.writeTypes(&buf, collectTypes(metadata.Schemas, true))
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		for _, sequence := range schema.Sequences {
			if sequence.SkipDump {
				continue
			}
			writeCreateSequence(&buf, schema.Name, sequence)
		}
	}
	if opts.PrintHeader {
		buf.WriteString(setDefaultTableSpace)
	}
	var schemas []*metadatapb.SchemaMetadata
	for _, schema := range metadata.Schemas {
		if !schema.SkipDump {
			schemas = append(schemas, schema)
		}
	}
	r.writeObjects(&buf, schemas)
	for _, schema := range schemas {
		for _, sequence := range schema.Sequences {
			if !sequence.SkipDump && sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
				writeAlterSequenceOwnedBy(&buf, schema.Name, sequence)
			}
		}
	}
	for _, schema := range schemas {
		writeSchemaForeignKeys(&buf, schema)
	}
	return buf.String(), nil
}

// SchemaDefinition renders one schema and the objects in it.
func SchemaDefinition(schema *metadatapb.SchemaMetadata) (string, error) {
	var buf strings.Builder
	writeSchema(&buf, schema)
	writeSchemaOwnerAndComment(&buf, schema)
	r := newRenderer("", []*metadatapb.SchemaMetadata{schema})
	r.writeTypes(&buf, collectTypes([]*metadatapb.SchemaMetadata{schema}, false))
	for _, sequence := range schema.Sequences {
		if !sequence.SkipDump {
			writeCreateSequence(&buf, schema.Name, sequence)
		}
	}
	r.writeObjects(&buf, []*metadatapb.SchemaMetadata{schema})
	for _, sequence := range schema.Sequences {
		if !sequence.SkipDump && sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
			writeAlterSequenceOwnedBy(&buf, schema.Name, sequence)
		}
	}
	writeSchemaForeignKeys(&buf, schema)
	return buf.String(), nil
}

// TableDefinition renders one table, its indexes and its foreign keys.
func TableDefinition(schema string, table *metadatapb.TableMetadata, sequences []*metadatapb.SequenceMetadata) (string, error) {
	var buf strings.Builder
	renderer{}.writeTable(&buf, schema, table)
	for _, sequence := range sequences {
		if sequence.OwnerTable == table.Name && sequence.OwnerColumn != "" {
			writeAlterSequenceOwnedBy(&buf, schema, sequence)
		}
	}
	for _, fk := range table.ForeignKeys {
		writeForeignKey(&buf, schema, table.Name, fk)
	}
	return buf.String(), nil
}

// ViewDefinition renders one view.
func ViewDefinition(schema string, view *metadatapb.ViewMetadata) (string, error) {
	var buf strings.Builder
	renderer{}.writeView(&buf, schema, view)
	return buf.String(), nil
}

// MaterializedViewDefinition renders one materialized view and its indexes.
func MaterializedViewDefinition(schema string, view *metadatapb.MaterializedViewMetadata) (string, error) {
	var buf strings.Builder
	renderer{}.writeMaterializedView(&buf, schema, view)
	return buf.String(), nil
}

// FunctionDefinition renders one function or procedure.
func FunctionDefinition(schema string, function *metadatapb.FunctionMetadata) (string, error) {
	var buf strings.Builder
	renderer{}.writeFunction(&buf, schema, function)
	return buf.String(), nil
}

// SequenceDefinition renders one sequence.
func SequenceDefinition(schema string, sequence *metadatapb.SequenceMetadata) (string, error) {
	var buf strings.Builder
	writeCreateSequence(&buf, schema, sequence)
	if sequence.OwnerColumn != "" && sequence.OwnerTable != "" {
		writeAlterSequenceOwnedBy(&buf, schema, sequence)
	}
	return buf.String(), nil
}

// renderer writes the objects of one database. database is the database's
// own name, dropped from the relations a definition names; it is empty when
// one object is rendered alone, which then keeps its names as stored.
type renderer struct {
	database string
	// userTypes maps an enum or composite type as the sync states a column
	// of it, schema and name unquoted, to the name a definition writes.
	userTypes map[string]string
}

func newRenderer(database string, schemas []*metadatapb.SchemaMetadata) renderer {
	r := renderer{database: database, userTypes: make(map[string]string)}
	for _, schema := range schemas {
		for _, t := range collectTypes([]*metadatapb.SchemaMetadata{schema}, false) {
			r.userTypes[schema.Name+"."+t.name()] = quoteIdent(schema.Name) + "." + quoteIdent(t.name())
		}
	}
	return r
}

// typeName returns a synced column type as a definition writes it. The sync
// states a user-defined type as schema.name, unquoted, so a name that needs
// quoting is quoted from the snapshot's own types.
func (r renderer) typeName(t string) string {
	if quoted, ok := r.userTypes[t]; ok {
		return quoted
	}
	if element, ok := strings.CutSuffix(t, "[]"); ok {
		if quoted, ok := r.userTypes[element]; ok {
			return quoted + "[]"
		}
	}
	return columnType(t)
}

// writeObjects writes the functions, tables, views and materialized views of
// schemas, each after the objects it depends on.
func (r renderer) writeObjects(buf *strings.Builder, schemas []*metadatapb.SchemaMetadata) {
	for _, object := range orderObjects(schemas) {
		switch {
		case object.function != nil:
			r.writeFunction(buf, object.schema, object.function)
		case object.table != nil:
			r.writeTable(buf, object.schema, object.table)
		case object.view != nil:
			r.writeView(buf, object.schema, object.view)
		case object.materializedView != nil:
			r.writeMaterializedView(buf, object.schema, object.materializedView)
		}
	}
}

// objectRef is one function, table, view or materialized view of a schema.
type objectRef struct {
	schema           string
	function         *metadatapb.FunctionMetadata
	table            *metadatapb.TableMetadata
	view             *metadatapb.ViewMetadata
	materializedView *metadatapb.MaterializedViewMetadata
}

// orderObjects returns the functions, tables, views and materialized views of
// schemas, each after the objects it depends on.
func orderObjects(schemas []*metadatapb.SchemaMetadata) []objectRef {
	graph := newDefinitionGraph()
	functions := make(map[string]*metadatapb.FunctionMetadata)
	tables := make(map[string]*metadatapb.TableMetadata)
	views := make(map[string]*metadatapb.ViewMetadata)
	materializedViews := make(map[string]*metadatapb.MaterializedViewMetadata)
	functionsByName := make(map[string][]string)
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			if function.SkipDump {
				continue
			}
			id := getObjectID(schema.Name, metadataFuncIdentity(function))
			functions[id] = function
			graph.AddNode(id)
			functionsByName[schema.Name+"."+function.Name] = append(functionsByName[schema.Name+"."+function.Name], id)
			functionsByName[function.Name] = append(functionsByName[function.Name], id)
		}
	}
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			if table.SkipDump {
				continue
			}
			id := getObjectID(schema.Name, table.Name)
			tables[id] = table
			graph.AddNode(id)
			// A column default or a check that calls a function follows it.
			var exprs []string
			for _, column := range table.Columns {
				exprs = append(exprs, column.Default)
			}
			for _, check := range table.CheckConstraints {
				exprs = append(exprs, check.Expression)
			}
			for _, expr := range exprs {
				for _, name := range calledFunctions(expr) {
					for _, functionID := range functionsByName[name] {
						graph.AddEdge(functionID, id)
					}
				}
			}
		}
		for _, view := range schema.Views {
			if view.SkipDump {
				continue
			}
			id := getObjectID(schema.Name, view.Name)
			views[id] = view
			graph.AddNode(id)
			for _, dependency := range view.DependencyColumns {
				graph.AddEdge(getObjectID(dependency.Schema, dependency.Table), id)
			}
			// The sync states a dependency per column read, so a relation the
			// query reads without naming a column comes from the query.
			for _, relation := range queryRelations(view.Definition) {
				graph.AddEdge(relation, id)
			}
		}
		for _, view := range schema.MaterializedViews {
			if view.SkipDump {
				continue
			}
			id := getObjectID(schema.Name, view.Name)
			materializedViews[id] = view
			graph.AddNode(id)
			for _, dependency := range view.DependencyColumns {
				graph.AddEdge(getObjectID(dependency.Schema, dependency.Table), id)
			}
			// The sync states a dependency per column read, so a relation the
			// query reads without naming a column comes from the query.
			for _, relation := range queryRelations(view.Definition) {
				graph.AddEdge(relation, id)
			}
		}
	}
	// CockroachDB resolves what a routine body names when it creates the
	// routine, and the sync records no routine dependencies, so a routine
	// follows the relations its body reads and the routines it calls, read on
	// CockroachDB's tree.
	for _, functionID := range sortedMapIDs(functions) {
		function := functions[functionID]
		for _, relation := range routineRelations(function.Definition) {
			graph.AddEdge(relation, functionID)
		}
		for _, call := range routineCalls(function.Definition) {
			for _, calledID := range functionsByName[call] {
				graph.AddEdge(calledID, functionID)
			}
		}
	}
	ordered, err := graph.TopologicalSort()
	if err != nil {
		// A cycle falls back to tables, views, materialized views, functions.
		ordered = nil
		ordered = append(ordered, sortedMapIDs(tables)...)
		ordered = append(ordered, sortedMapIDs(views)...)
		ordered = append(ordered, sortedMapIDs(materializedViews)...)
		ordered = append(ordered, sortedMapIDs(functions)...)
	}
	var objects []objectRef
	for _, id := range ordered {
		object := objectRef{
			schema:           getSchemaNameFromID(id),
			function:         functions[id],
			table:            tables[id],
			view:             views[id],
			materializedView: materializedViews[id],
		}
		if object.function != nil || object.table != nil || object.view != nil || object.materializedView != nil {
			objects = append(objects, object)
		}
	}
	return objects
}

// writeTable writes CREATE TABLE with its primary key and checks, the table's
// comments, and its secondary indexes.
func (r renderer) writeTable(buf *strings.Builder, schema string, table *metadatapb.TableMetadata) {
	columns, primary, checks := tableShape(table)
	r.writeCreateTable(buf, schema, table.Name, columns, primary, checks)
	buf.WriteString(";\n\n")
	if table.Owner != "" {
		fmt.Fprintf(buf, "ALTER TABLE %s.%s OWNER TO %s;\n\n", quoteIdent(schema), quoteIdent(table.Name), quoteIdent(table.Owner))
	}
	if table.Comment != "" {
		writeRelationComment(buf, "TABLE", schema, table.Name, table.Comment)
	}
	for _, column := range columns {
		if column.Comment != "" {
			writeColumnComment(buf, schema, table.Name, column)
		}
	}
	for _, check := range checks {
		if check.Comment != "" {
			writeConstraintComment(buf, schema, table.Name, check.Name, check.Comment)
		}
	}
	r.writeIndexes(buf, schema, table.Name, table.Indexes)
}

// writeIndexes writes the secondary indexes of a table or materialized view
// and every index comment, the primary key's included.
func (r renderer) writeIndexes(buf *strings.Builder, schema, table string, indexes []*metadatapb.IndexMetadata) {
	for _, index := range indexes {
		if !index.Primary {
			buf.WriteString(indexStatement(schema, table, index))
			buf.WriteString(";\n\n")
		}
		if index.Comment != "" {
			writeIndexComment(buf, schema, table, index)
		}
	}
}

// tableShape returns the columns, primary key and checks CREATE TABLE
// declares, leaving out what CockroachDB adds itself: the hidden rowid key of a
// table created without one, and the hidden shard column of a hash-sharded
// index with the check that bounds it, which the index creates again.
func tableShape(table *metadatapb.TableMetadata) ([]*metadatapb.ColumnMetadata, *metadatapb.IndexMetadata, []*metadatapb.CheckConstraintMetadata) {
	var primary *metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if index.Primary {
			primary = index
			break
		}
	}
	implicit := ""
	if primary != nil && len(primary.Expressions) == 1 {
		for _, column := range table.Columns {
			if column.Name == primary.Expressions[0] && isImplicitRowID(column) {
				implicit = column.Name
				primary = nil
				break
			}
		}
	}
	columns := make([]*metadatapb.ColumnMetadata, 0, len(table.Columns))
	shardChecks := make(map[string]bool)
	for _, column := range table.Columns {
		if isShardColumn(column) {
			shardChecks["check_"+column.Name] = true
			continue
		}
		if column.Name != implicit {
			columns = append(columns, column)
		}
	}
	checks := make([]*metadatapb.CheckConstraintMetadata, 0, len(table.CheckConstraints))
	for _, check := range table.CheckConstraints {
		if !shardChecks[check.Name] {
			checks = append(checks, check)
		}
	}
	return columns, primary, checks
}

// isImplicitRowID reports whether column is the hidden key CockroachDB gives a
// table created without a primary key.
func isImplicitRowID(column *metadatapb.ColumnMetadata) bool {
	return column.IsInvisible && strings.HasPrefix(column.Name, "rowid") && column.Default == "unique_rowid()"
}

// isShardColumn reports whether column is the hidden column of a hash-sharded
// index or primary key.
func isShardColumn(column *metadatapb.ColumnMetadata) bool {
	return column.IsInvisible && strings.HasPrefix(column.Name, "crdb_internal_")
}

func (r renderer) writeCreateTable(buf *strings.Builder, schema, table string, columns []*metadatapb.ColumnMetadata, primary *metadatapb.IndexMetadata, checks []*metadatapb.CheckConstraintMetadata) {
	fmt.Fprintf(buf, "CREATE TABLE %s.%s (", quoteIdent(schema), quoteIdent(table))
	first := true
	sep := func() {
		if first {
			first = false
			buf.WriteString("\n    ")
			return
		}
		buf.WriteString(",\n    ")
	}
	for _, column := range columns {
		sep()
		r.writeColumnDefinition(buf, column)
	}
	if primary != nil {
		sep()
		buf.WriteString(primaryKeyClause(primary))
	}
	for _, check := range checks {
		sep()
		fmt.Fprintf(buf, "CONSTRAINT %s CHECK %s", quoteIdent(check.Name), check.Expression)
	}
	buf.WriteString("\n)")
}

func (r renderer) writeColumnDefinition(buf *strings.Builder, column *metadatapb.ColumnMetadata) {
	fmt.Fprintf(buf, "%s %s", quoteIdent(column.Name), r.typeName(column.Type))
	if column.Collation != "" {
		fmt.Fprintf(buf, " COLLATE %s", quoteIdent(column.Collation))
	}
	if column.IsInvisible {
		buf.WriteString(" NOT VISIBLE")
	}
	if column.Generation != nil {
		storage := "STORED"
		if column.Generation.Type == metadatapb.GenerationMetadata_TYPE_VIRTUAL {
			storage = "VIRTUAL"
		}
		fmt.Fprintf(buf, " AS (%s) %s", column.Generation.Expression, storage)
	} else if column.Default != "" {
		buf.WriteString(" DEFAULT ")
		buf.WriteString(defaultExpression(column.Default))
	}
	if !column.Nullable {
		buf.WriteString(" NOT NULL")
	}
}

func writeSchemaForeignKeys(buf *strings.Builder, schema *metadatapb.SchemaMetadata) {
	for _, table := range schema.Tables {
		if table.SkipDump {
			continue
		}
		for _, fk := range table.ForeignKeys {
			writeForeignKey(buf, schema.Name, table.Name, fk)
			if fk.Comment != "" {
				writeConstraintComment(buf, schema.Name, table.Name, fk.Name, fk.Comment)
			}
		}
	}
}

func writeForeignKey(buf *strings.Builder, schema, table string, fk *metadatapb.ForeignKeyMetadata) {
	fmt.Fprintf(buf, "ALTER TABLE %s.%s\n    ADD CONSTRAINT %s FOREIGN KEY (%s)\n    REFERENCES %s.%s (%s)",
		quoteIdent(schema), quoteIdent(table), quoteIdent(fk.Name), quoteList(fk.Columns),
		quoteIdent(fk.ReferencedSchema), quoteIdent(fk.ReferencedTable), quoteList(fk.ReferencedColumns))
	// MATCH SIMPLE is the default and is omitted.
	if fk.MatchType == "FULL" {
		buf.WriteString("\n    MATCH FULL")
	}
	if fk.OnDelete != "" && fk.OnDelete != "NO ACTION" {
		buf.WriteString("\n    ON DELETE " + fk.OnDelete)
	}
	if fk.OnUpdate != "" && fk.OnUpdate != "NO ACTION" {
		buf.WriteString("\n    ON UPDATE " + fk.OnUpdate)
	}
	buf.WriteString(";\n\n")
}

func quoteList(names []string) string {
	quoted := make([]string, len(names))
	for i, name := range names {
		quoted[i] = quoteIdent(name)
	}
	return strings.Join(quoted, ", ")
}

func (r renderer) writeView(buf *strings.Builder, schema string, view *metadatapb.ViewMetadata) {
	// CockroachDB keeps a view's column list out of its stored query, so the
	// synced output names are written.
	columns := ""
	if len(view.Columns) > 0 {
		names := make([]string, len(view.Columns))
		for i, column := range view.Columns {
			names[i] = column.Name
		}
		columns = " (" + quoteList(names) + ")"
	}
	fmt.Fprintf(buf, "CREATE VIEW %s.%s%s AS \n%s;\n\n", quoteIdent(schema), quoteIdent(view.Name), columns, r.query(view.Definition))
}

func (r renderer) writeMaterializedView(buf *strings.Builder, schema string, view *metadatapb.MaterializedViewMetadata) {
	// Not WITH NO DATA: CockroachDB indexes only a populated materialized view.
	fmt.Fprintf(buf, "CREATE MATERIALIZED VIEW %s.%s AS \n%s;\n\n", quoteIdent(schema), quoteIdent(view.Name), r.query(view.Definition))
	r.writeIndexes(buf, schema, view.Name, materializedViewIndexes(view))
}

// materializedViewIndexes leaves out the hidden rowid key CockroachDB gives
// every materialized view.
func materializedViewIndexes(view *metadatapb.MaterializedViewMetadata) []*metadatapb.IndexMetadata {
	var indexes []*metadatapb.IndexMetadata
	for _, index := range view.Indexes {
		if index.Primary && len(index.Expressions) == 1 && strings.HasPrefix(index.Expressions[0], "rowid") {
			continue
		}
		indexes = append(indexes, index)
	}
	return indexes
}

func (r renderer) query(definition string) string {
	if r.database == "" {
		return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(definition), ";"))
	}
	return relationQuery(definition, r.database)
}

func (r renderer) writeFunction(buf *strings.Builder, _ string, function *metadatapb.FunctionMetadata) {
	definition := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(function.Definition), ";"))
	if r.database != "" {
		definition = routineDefinition(definition, r.database)
	}
	buf.WriteString(definition)
	buf.WriteString(";\n\n")
}

func writeSchema(buf *strings.Builder, schema *metadatapb.SchemaMetadata) {
	if schema.Name == "public" {
		return
	}
	fmt.Fprintf(buf, "CREATE SCHEMA IF NOT EXISTS %s;\n\n", quoteIdent(schema.Name))
}

// writeSchemaOwnerAndComment restores a schema's owner and comment. The public
// schema already exists and keeps its owner, so only its comment is written.
func writeSchemaOwnerAndComment(buf *strings.Builder, schema *metadatapb.SchemaMetadata) {
	if schema.Name != "public" && schema.Owner != "" {
		fmt.Fprintf(buf, "ALTER SCHEMA %s OWNER TO %s;\n\n", quoteIdent(schema.Name), quoteIdent(schema.Owner))
	}
	if schema.Comment != "" {
		fmt.Fprintf(buf, "COMMENT ON SCHEMA %s IS '%s';\n\n", quoteIdent(schema.Name), escapeSingleQuote(schema.Comment))
	}
}

func writeExtension(buf *strings.Builder, extension *metadatapb.ExtensionMetadata) {
	fmt.Fprintf(buf, "CREATE EXTENSION IF NOT EXISTS %s", quoteIdent(extension.Name))
	if extension.Schema != "" {
		fmt.Fprintf(buf, " WITH SCHEMA %s", quoteIdent(extension.Schema))
	}
	if extension.Version != "" {
		fmt.Fprintf(buf, " VERSION '%s'", extension.Version)
	}
	buf.WriteString(";\n\n")
	if extension.Description != "" {
		fmt.Fprintf(buf, "COMMENT ON EXTENSION %s IS '%s';\n\n", quoteIdent(extension.Name), escapeSingleQuote(extension.Description))
	}
}

func writeCreateSequence(buf *strings.Builder, schema string, sequence *metadatapb.SequenceMetadata) {
	fmt.Fprintf(buf, "CREATE SEQUENCE %s.%s\n    AS %s\n\tSTART WITH %s\n\tINCREMENT BY %s\n\tMINVALUE %s\n\tMAXVALUE %s",
		quoteIdent(schema), quoteIdent(sequence.Name), sequence.DataType, sequence.Start, sequence.Increment, sequence.MinValue, sequence.MaxValue)
	if sequence.CacheSize != "" {
		fmt.Fprintf(buf, "\n\tCACHE %s", sequence.CacheSize)
	}
	if sequence.Cycle {
		buf.WriteString("\n\tCYCLE")
	} else {
		buf.WriteString("\n\tNO CYCLE")
	}
	buf.WriteString(";\n\n")
}

func writeAlterSequenceOwnedBy(buf *strings.Builder, schema string, sequence *metadatapb.SequenceMetadata) {
	fmt.Fprintf(buf, "ALTER SEQUENCE %s.%s OWNED BY %s.%s.%s;\n\n", quoteIdent(schema), quoteIdent(sequence.Name), quoteIdent(schema), quoteIdent(sequence.OwnerTable), quoteIdent(sequence.OwnerColumn))
}

// qualifiedType is an enum or composite type of a schema.
type qualifiedType struct {
	schema    string
	enum      *metadatapb.EnumTypeMetadata
	composite *metadatapb.CompositeTypeMetadata
}

func (t qualifiedType) name() string {
	if t.enum != nil {
		return t.enum.Name
	}
	return t.composite.Name
}

func collectTypes(schemas []*metadatapb.SchemaMetadata, skipDumped bool) []qualifiedType {
	var types []qualifiedType
	for _, schema := range schemas {
		if skipDumped && schema.SkipDump {
			continue
		}
		for _, enum := range schema.EnumTypes {
			if !enum.SkipDump {
				types = append(types, qualifiedType{schema: schema.Name, enum: enum})
			}
		}
		for _, composite := range schema.CompositeTypes {
			if !composite.SkipDump {
				types = append(types, qualifiedType{schema: schema.Name, composite: composite})
			}
		}
	}
	return orderTypes(types)
}

// orderTypes puts each composite type after the types its attributes use.
func orderTypes(types []qualifiedType) []qualifiedType {
	graph := newDefinitionGraph()
	byID := make(map[string]qualifiedType)
	byName := make(map[string]string)
	for _, t := range types {
		id := getObjectID(t.schema, t.name())
		byID[id] = t
		graph.AddNode(id)
		byName[t.schema+"."+t.name()] = id
	}
	for _, t := range types {
		if t.composite == nil {
			continue
		}
		id := getObjectID(t.schema, t.name())
		for _, attribute := range t.composite.Attributes {
			name := strings.TrimSuffix(strings.ReplaceAll(attribute.Type, `"`, ""), "[]")
			if dep, ok := byName[name]; ok && dep != id {
				graph.AddEdge(dep, id)
			} else if dep, ok := byName[t.schema+"."+name]; ok && dep != id {
				graph.AddEdge(dep, id)
			}
		}
	}
	ordered, err := graph.TopologicalSort()
	if err != nil {
		return types
	}
	out := make([]qualifiedType, 0, len(types))
	for _, id := range ordered {
		out = append(out, byID[id])
	}
	return out
}

func (r renderer) writeTypes(buf *strings.Builder, types []qualifiedType) {
	for _, t := range types {
		if t.enum != nil {
			writeEnum(buf, t.schema, t.enum)
			buf.WriteString(";\n\n")
			if t.enum.Comment != "" {
				writeRelationComment(buf, "TYPE", t.schema, t.enum.Name, t.enum.Comment)
			}
			continue
		}
		r.writeCompositeType(buf, t.schema, t.composite)
		buf.WriteString(";\n\n")
		if t.composite.Comment != "" {
			writeRelationComment(buf, "TYPE", t.schema, t.composite.Name, t.composite.Comment)
		}
	}
}

func writeEnum(buf *strings.Builder, schema string, enum *metadatapb.EnumTypeMetadata) {
	fmt.Fprintf(buf, "CREATE TYPE %s.%s AS ENUM (\n", quoteIdent(schema), quoteIdent(enum.Name))
	for i, value := range enum.Values {
		if i > 0 {
			buf.WriteString(",\n")
		}
		fmt.Fprintf(buf, "    '%s'", escapeSingleQuote(value))
	}
	buf.WriteString("\n)")
}

func (r renderer) writeCompositeType(buf *strings.Builder, schema string, composite *metadatapb.CompositeTypeMetadata) {
	fmt.Fprintf(buf, "CREATE TYPE %s.%s AS (\n", quoteIdent(schema), quoteIdent(composite.Name))
	for i, attribute := range composite.Attributes {
		if i > 0 {
			buf.WriteString(",\n")
		}
		fmt.Fprintf(buf, "    %s %s", quoteIdent(attribute.Name), r.typeName(attribute.Type))
		if attribute.Collation != "" {
			fmt.Fprintf(buf, " COLLATE %s", attribute.Collation)
		}
	}
	buf.WriteString("\n)")
}

func writeRelationComment(buf *strings.Builder, kind, schema, name, comment string) {
	fmt.Fprintf(buf, "COMMENT ON %s %s.%s IS '%s';\n\n", kind, quoteIdent(schema), quoteIdent(name), escapeSingleQuote(comment))
}

func writeColumnComment(buf *strings.Builder, schema, table string, column *metadatapb.ColumnMetadata) {
	fmt.Fprintf(buf, "COMMENT ON COLUMN %s.%s.%s IS '%s';\n\n", quoteIdent(schema), quoteIdent(table), quoteIdent(column.Name), escapeSingleQuote(column.Comment))
}

func writeConstraintComment(buf *strings.Builder, schema, table, name, comment string) {
	fmt.Fprintf(buf, "COMMENT ON CONSTRAINT %s ON %s.%s IS '%s';\n\n", quoteIdent(name), quoteIdent(schema), quoteIdent(table), escapeSingleQuote(comment))
}

// writeIndexComment names the index through its table, the form CockroachDB
// reads in every database.
func writeIndexComment(buf *strings.Builder, schema, table string, index *metadatapb.IndexMetadata) {
	fmt.Fprintf(buf, "COMMENT ON INDEX %s.%s@%s IS '%s';\n\n", quoteIdent(schema), quoteIdent(table), quoteIdent(index.Name), escapeSingleQuote(index.Comment))
}

func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// objectIDSeparator joins schema and object into a graph key. NUL cannot
// appear in an identifier, so a name holding a dot stays unambiguous.
const objectIDSeparator = "\x00"

func getObjectID(schema, object string) string {
	return schema + objectIDSeparator + object
}

func getSchemaNameFromID(id string) string {
	schema, _, _ := strings.Cut(id, objectIDSeparator)
	return schema
}

// metadataFuncIdentity returns a function's identity, its signature where the
// sync states one, so overloads stay apart.
func metadataFuncIdentity(f *metadatapb.FunctionMetadata) string {
	if f.Signature != "" {
		return f.Signature
	}
	return f.Name
}

func sortedMapIDs[V any](m map[string]V) []string {
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids
}
