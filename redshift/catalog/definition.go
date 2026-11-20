package catalog

import (
	"fmt"
	"io"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/redshift/ast"
	pgparser "github.com/bytebase/omni/redshift/parser"
)

// DatabaseDefinition renders a synced Redshift database snapshot as DDL.
func DatabaseDefinition(to *metadatapb.DatabaseSchemaMetadata) (string, error) {
	if to == nil {
		return "", nil
	}
	toState := convertToDatabaseState(to)

	var sb strings.Builder

	if needsCaseSensitiveIdentifiers(to) {
		if _, err := sb.WriteString("SET enable_case_sensitive_identifier TO true;\n"); err != nil {
			return "", err
		}
	}
	if err := writeSchemas(&sb, to); err != nil {
		return "", err
	}
	if err := writeTables(&sb, to, toState); err != nil {
		return "", err
	}
	if err := writeAllViews(&sb, to, toState); err != nil {
		return "", err
	}

	s := sb.String()
	// Trim the leading newline so the output starts directly with the first statement or banner.
	s = strings.TrimLeft(s, "\n")
	return s, nil
}

// GetTableDefinition renders one table as DDL, matching DatabaseDefinition
// without the section banner.
func GetTableDefinition(schemaName string, table *metadatapb.TableMetadata) (string, error) {
	if table == nil || table.SkipDump {
		return "", nil
	}
	var sb strings.Builder
	if err := writeTable(&sb, convertToTableState(0, schemaName, table)); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// GetViewDefinition renders one view as DDL, matching DatabaseDefinition
// without the section banner.
func GetViewDefinition(schemaName string, view *metadatapb.ViewMetadata) (string, error) {
	if view == nil || view.SkipDump {
		return "", nil
	}
	var sb strings.Builder
	if err := writeView(&sb, convertToViewState(0, schemaName, view)); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func needsCaseSensitiveIdentifiers(to *metadatapb.DatabaseSchemaMetadata) bool {
	hasUpper := func(s string) bool {
		for i := 0; i < len(s); i++ {
			if s[i] >= 'A' && s[i] <= 'Z' {
				return true
			}
		}
		return false
	}
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		if hasUpper(schema.Name) {
			return true
		}
		for _, table := range schema.Tables {
			if table == nil || table.SkipDump {
				continue
			}
			if hasUpper(table.Name) {
				return true
			}
			for _, col := range table.Columns {
				if col != nil && hasUpper(col.Name) {
					return true
				}
			}
			for _, idx := range table.Indexes {
				if idx != nil && (idx.Primary || idx.Unique) && hasUpper(idx.Name) {
					return true
				}
			}
			for _, fk := range table.ForeignKeys {
				if fk == nil {
					continue
				}
				if hasUpper(fk.Name) || hasUpper(fk.ReferencedSchema) || hasUpper(fk.ReferencedTable) {
					return true
				}
				for _, c := range fk.Columns {
					if hasUpper(c) {
						return true
					}
				}
				for _, c := range fk.ReferencedColumns {
					if hasUpper(c) {
						return true
					}
				}
			}
		}
		for _, view := range schema.Views {
			if view != nil && !view.SkipDump && hasUpper(view.Name) {
				return true
			}
		}
		for _, mv := range schema.MaterializedViews {
			if mv != nil && !mv.SkipDump && hasUpper(mv.Name) {
				return true
			}
		}
	}
	return false
}

func writeSchemas(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata) error {
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump || schema.Name == "" || schema.Name == "public" {
			continue
		}
		if _, err := w.WriteString(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;\n", quoteDefinitionIdentifier(schema.Name))); err != nil {
			return err
		}
	}
	return nil
}

func writeTables(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState) error {
	tables := orderedTables(to, state)
	known := make(map[string]bool, len(tables))
	for _, table := range tables {
		known[objectKey(table.schema, table.name)] = true
	}

	type deferredFK struct {
		table *tableState
		fk    *foreignKeyState
	}
	var deferred []deferredFK
	emitted := make(map[string]bool, len(tables))

	for _, table := range tables {
		if _, err := w.WriteString(getTableAnnouncement(table.name)); err != nil {
			return err
		}

		tableKey := objectKey(table.schema, table.name)
		var inlineFKs []*foreignKeyState
		for _, fk := range sortedForeignKeys(table.foreignKeys) {
			refSchema := fk.referencedSchema
			if refSchema == "" {
				refSchema = table.schema
			}
			refKey := objectKey(refSchema, fk.referencedTable)
			// Defer any FK whose target table exists in this dump and has not been
			// created yet (cyclic references) so every CREATE TABLE stays valid.
			if refKey != tableKey && known[refKey] && !emitted[refKey] {
				deferred = append(deferred, deferredFK{table: table, fk: fk})
				continue
			}
			inlineFKs = append(inlineFKs, fk)
		}

		if err := writeTableWithFKs(w, table, inlineFKs); err != nil {
			return err
		}
		emitted[tableKey] = true
	}

	if len(deferred) > 0 {
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
		for _, item := range deferred {
			var buf strings.Builder
			if err := item.fk.toString(&buf); err != nil {
				return err
			}
			if _, err := w.WriteString(fmt.Sprintf("ALTER TABLE %s ADD %s;\n", item.table.qualifiedName(), buf.String())); err != nil {
				return err
			}
		}
	}
	return nil
}

type orderedViewEntry struct {
	state        *viewState
	meta         *metadatapb.ViewMetadata
	depCols      []*metadatapb.DependencyColumn
	materialized bool
}

func writeAllViews(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState) error {
	for _, entry := range orderedAllViews(to, state) {
		if entry.materialized {
			if _, err := w.WriteString(getMaterializedViewAnnouncement(entry.state.name)); err != nil {
				return err
			}
			if err := writeMaterializedView(w, entry.state); err != nil {
				return err
			}
			continue
		}
		if _, err := w.WriteString(getViewAnnouncement(entry.state.name)); err != nil {
			return err
		}
		if err := writeView(w, entry.state); err != nil {
			return err
		}
	}
	return nil
}

func objectKey(schemaName, objectName string) string {
	return schemaName + "\x00" + objectName
}

// orderedTables returns tables in snapshot order, topologically sorted so
// foreign-key referenced parent tables are created before child tables.
func orderedTables(to *metadatapb.DatabaseSchemaMetadata, state *databaseState) []*tableState {
	var tables []*tableState
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		for _, table := range schema.Tables {
			if table == nil || table.SkipDump {
				continue
			}
			ts, ok := schemaState.tables[table.Name]
			if !ok {
				continue
			}
			tables = append(tables, ts)
			delete(schemaState.tables, table.Name)
		}
	}
	if len(tables) < 2 {
		return tables
	}

	index := make(map[string]int, len(tables))
	for i, t := range tables {
		index[objectKey(t.schema, t.name)] = i
	}
	dependencies := make([]map[int]struct{}, len(tables))
	for i, t := range tables {
		for _, fk := range t.foreignKeys {
			refSchema := fk.referencedSchema
			if refSchema == "" {
				refSchema = t.schema
			}
			dep, ok := index[objectKey(refSchema, fk.referencedTable)]
			if !ok || dep == i {
				continue
			}
			if dependencies[i] == nil {
				dependencies[i] = map[int]struct{}{}
			}
			dependencies[i][dep] = struct{}{}
		}
	}

	order := topoOrder(len(tables), dependencies)
	ordered := make([]*tableState, len(order))
	for i, idx := range order {
		ordered[i] = tables[idx]
	}
	return ordered
}

// orderedAllViews returns ordinary and materialized views in one topological
// order so cross-kind dependencies (e.g., an ordinary view reading a
// materialized view) emit the dependency first.
func orderedAllViews(to *metadatapb.DatabaseSchemaMetadata, state *databaseState) []orderedViewEntry {
	var entries []orderedViewEntry
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		for _, view := range schema.Views {
			if view == nil || view.SkipDump {
				continue
			}
			vs, ok := schemaState.views[view.Name]
			if !ok {
				continue
			}
			entries = append(entries, orderedViewEntry{state: vs, meta: view})
			delete(schemaState.views, view.Name)
		}
	}
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		for _, mv := range schema.MaterializedViews {
			if mv == nil || mv.SkipDump {
				continue
			}
			vs, ok := schemaState.materializedViews[mv.Name]
			if !ok {
				continue
			}
			entries = append(entries, orderedViewEntry{state: vs, depCols: mv.DependencyColumns, materialized: true})
			delete(schemaState.materializedViews, mv.Name)
		}
	}
	return sortViewEntries(entries)
}

func sortViewEntries(entries []orderedViewEntry) []orderedViewEntry {
	if len(entries) < 2 {
		return entries
	}

	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[objectKey(e.state.schema, e.state.name)] = i
	}
	dependencies := make([]map[int]struct{}, len(entries))
	addDep := func(i, dep int) {
		if dep == i {
			return
		}
		if dependencies[i] == nil {
			dependencies[i] = map[int]struct{}{}
		}
		dependencies[i][dep] = struct{}{}
	}
	for i, e := range entries {
		depCols := e.depCols
		if e.meta != nil {
			depCols = e.meta.DependencyColumns
		}
		for _, dep := range depCols {
			if dep == nil {
				continue
			}
			depSchema := dep.Schema
			if depSchema == "" {
				depSchema = e.state.schema
			}
			if depIdx, ok := index[objectKey(depSchema, dep.Table)]; ok {
				addDep(i, depIdx)
			}
		}
		for _, ref := range collectViewRelationRefs(e.state.definition) {
			if ref.schema != "" {
				if depIdx, ok := index[objectKey(ref.schema, ref.name)]; ok {
					addDep(i, depIdx)
				}
				continue
			}
			if depIdx, ok := index[objectKey(e.state.schema, ref.name)]; ok {
				addDep(i, depIdx)
			} else if depIdx, ok := index[objectKey("public", ref.name)]; ok {
				addDep(i, depIdx)
			}
		}
	}

	order := topoOrder(len(entries), dependencies)
	ordered := make([]orderedViewEntry, len(order))
	for i, idx := range order {
		ordered[i] = entries[idx]
	}
	return ordered
}

type relationRef struct {
	schema string
	name   string
}

func collectViewRelationRefs(definition string) []relationRef {
	if strings.TrimSpace(definition) == "" {
		return nil
	}
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil {
		return nil
	}
	var refs []relationRef
	nodes.Walk(&scopedRelVisitor{refs: &refs}, stmts)
	return refs
}

type scopedRelVisitor struct {
	cteScope map[string]bool
	refs     *[]relationRef
}

func (v *scopedRelVisitor) Visit(node nodes.Node) nodes.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *nodes.RangeVar:
		if n.Relname != "" && (n.Schemaname != "" || !v.cteScope[n.Relname]) {
			*v.refs = append(*v.refs, relationRef{schema: n.Schemaname, name: n.Relname})
		}
		return v
	case *nodes.SelectStmt:
		if n.WithClause == nil {
			return v
		}
		childScope := make(map[string]bool, len(v.cteScope))
		for k, val := range v.cteScope {
			childScope[k] = val
		}
		if n.WithClause.Recursive && n.WithClause.Ctes != nil {
			for _, item := range n.WithClause.Ctes.Items {
				if cte, ok := item.(*nodes.CommonTableExpr); ok && cte.Ctename != "" {
					childScope[cte.Ctename] = true
				}
			}
		}
		if n.WithClause.Ctes != nil {
			for _, item := range n.WithClause.Ctes.Items {
				cte, ok := item.(*nodes.CommonTableExpr)
				if !ok {
					continue
				}
				if cte.Ctequery != nil {
					cteVisitor := v
					if n.WithClause.Recursive || len(childScope) > len(v.cteScope) {
						cteVisitor = &scopedRelVisitor{cteScope: childScope, refs: v.refs}
					}
					nodes.Walk(cteVisitor, cte.Ctequery)
				}
				if !n.WithClause.Recursive && cte.Ctename != "" {
					childScope[cte.Ctename] = true
				}
			}
		}
		child := &scopedRelVisitor{cteScope: childScope, refs: v.refs}
		walkNodeList(child, n.DistinctClause)
		if n.IntoClause != nil {
			nodes.Walk(child, n.IntoClause)
		}
		walkNodeList(child, n.TargetList)
		walkNodeList(child, n.FromClause)
		nodes.Walk(child, n.WhereClause)
		walkNodeList(child, n.GroupClause)
		nodes.Walk(child, n.HavingClause)
		nodes.Walk(child, n.StartWith)
		nodes.Walk(child, n.ConnectBy)
		nodes.Walk(child, n.QualifyClause)
		walkNodeList(child, n.WindowClause)
		walkNodeList(child, n.ExcludeList)
		walkNodeList(child, n.ValuesLists)
		walkNodeList(child, n.SortClause)
		nodes.Walk(child, n.LimitOffset)
		nodes.Walk(child, n.LimitCount)
		walkNodeList(child, n.LockingClause)
		if n.Larg != nil {
			nodes.Walk(child, n.Larg)
		}
		if n.Rarg != nil {
			nodes.Walk(child, n.Rarg)
		}
		return nil
	default:
		return v
	}
}

func walkNodeList(v nodes.Visitor, list *nodes.List) {
	if list == nil {
		return
	}
	for _, item := range list.Items {
		nodes.Walk(v, item)
	}
}

// topoOrder returns indices 0..n-1 ordered so each item follows its
// dependencies, preserving input order among independent items and breaking
// cycles in input order.
func topoOrder(n int, dependencies []map[int]struct{}) []int {
	remaining := make([]int, n)
	for i, deps := range dependencies {
		remaining[i] = len(deps)
	}
	ordered := make([]int, 0, n)
	done := make([]bool, n)
	for len(ordered) < n {
		progressed := false
		for i := 0; i < n; i++ {
			if done[i] || remaining[i] != 0 {
				continue
			}
			done[i] = true
			ordered = append(ordered, i)
			progressed = true
			for j, deps := range dependencies {
				if _, ok := deps[i]; ok && !done[j] {
					remaining[j]--
				}
			}
		}
		if !progressed {
			for i := 0; i < n; i++ {
				if !done[i] {
					done[i] = true
					ordered = append(ordered, i)
				}
			}
		}
	}
	return ordered
}

// writeTable writes one table and the comments that belong to it. Column
// comments follow the column order the CREATE TABLE body uses so output is
// deterministic across runs.
func writeTable(w io.StringWriter, table *tableState) error {
	return writeTableWithFKs(w, table, sortedForeignKeys(table.foreignKeys))
}

func writeTableWithFKs(w io.StringWriter, table *tableState, foreignKeys []*foreignKeyState) error {
	buf := &strings.Builder{}
	if err := table.toStringWithFKs(buf, foreignKeys); err != nil {
		return err
	}
	if table.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON TABLE %s IS '%s';\n", table.qualifiedName(), escapeSingleQuote(table.comment)); err != nil {
			return err
		}
	}
	for _, column := range sortedColumns(table.columns) {
		if column.comment == "" {
			continue
		}
		if _, err := fmt.Fprintf(buf, "COMMENT ON COLUMN %s.%s IS '%s';\n", table.qualifiedName(), quoteDefinitionIdentifier(column.name), escapeSingleQuote(column.comment)); err != nil {
			return err
		}
	}
	_, err := w.WriteString(buf.String())
	return err
}

// writeView writes one view and its comment.
func writeView(w io.StringWriter, view *viewState) error {
	buf := &strings.Builder{}
	if err := view.toString(buf); err != nil {
		return err
	}
	if view.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON VIEW %s IS '%s';\n", view.qualifiedName(), escapeSingleQuote(view.comment)); err != nil {
			return err
		}
	}
	_, err := w.WriteString(buf.String())
	return err
}

func writeMaterializedView(w io.StringWriter, mv *viewState) error {
	buf := &strings.Builder{}
	if err := mv.toMaterializedString(buf); err != nil {
		return err
	}
	if mv.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON MATERIALIZED VIEW %s IS '%s';\n", mv.qualifiedName(), escapeSingleQuote(mv.comment)); err != nil {
			return err
		}
	}
	_, err := w.WriteString(buf.String())
	return err
}

func getTableAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Table structure for `%s`\n--\n", name)
}

func getViewAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- View structure for `%s`\n--\n", name)
}

func getMaterializedViewAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Materialized view structure for `%s`\n--\n", name)
}

type databaseState struct {
	name    string
	schemas map[string]*schemaState
}

func newDatabaseState() *databaseState {
	return &databaseState{
		schemas: make(map[string]*schemaState),
	}
}

func convertToDatabaseState(database *metadatapb.DatabaseSchemaMetadata) *databaseState {
	state := newDatabaseState()
	state.name = database.Name
	for _, schema := range database.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		state.schemas[schema.Name] = convertToSchemaState(schema)
	}
	return state
}

type schemaState struct {
	name              string
	tables            map[string]*tableState
	views             map[string]*viewState
	materializedViews map[string]*viewState
}

// qualifyName prefixes the object with its schema, so DDL for analytics.orders
// cannot land in public. The prefix is dropped when no schema is known, which is
// how the single-object entry points behave when the caller passes none.
func qualifyName(schemaName, objectName string) string {
	if schemaName == "" {
		return quoteDefinitionIdentifier(objectName)
	}
	return quoteDefinitionIdentifier(schemaName) + "." + quoteDefinitionIdentifier(objectName)
}

// quoteDefinitionIdentifier renders a name as a quoted SQL identifier. Redshift
// folds an unquoted name to lower case and rejects one containing a space or a
// reserved word outright, so a name the catalog reports verbatim has to be
// quoted to come back as itself. An embedded double quote is doubled, as in
// PostgreSQL.
func quoteDefinitionIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func (t *tableState) qualifiedName() string {
	return qualifyName(t.schema, t.name)
}

func (v *viewState) qualifiedName() string {
	return qualifyName(v.schema, v.name)
}

// escapeSingleQuote encodes a comment or literal for a SQL string literal.
func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func newSchemaState() *schemaState {
	return &schemaState{
		tables:            make(map[string]*tableState),
		views:             make(map[string]*viewState),
		materializedViews: make(map[string]*viewState),
	}
}

func convertToSchemaState(schema *metadatapb.SchemaMetadata) *schemaState {
	state := newSchemaState()
	state.name = schema.Name
	for i, table := range schema.Tables {
		if table == nil || table.SkipDump {
			continue
		}
		state.tables[table.Name] = convertToTableState(i, schema.Name, table)
	}
	for i, view := range schema.Views {
		if view == nil || view.SkipDump {
			continue
		}
		state.views[view.Name] = convertToViewState(i, schema.Name, view)
	}
	for i, mv := range schema.MaterializedViews {
		if mv == nil || mv.SkipDump {
			continue
		}
		state.materializedViews[mv.Name] = convertToMaterializedViewState(i, schema.Name, mv)
	}
	return state
}

type tableState struct {
	id          int
	schema      string
	name        string
	columns     map[string]*columnState
	indexes     map[string]*indexState
	foreignKeys map[string]*foreignKeyState
	comment     string
}

func sortedForeignKeys(fks map[string]*foreignKeyState) []*foreignKeyState {
	foreignKeys := make([]*foreignKeyState, 0, len(fks))
	for _, fk := range fks {
		if len(fk.columns) > 0 && len(fk.referencedColumns) > 0 {
			foreignKeys = append(foreignKeys, fk)
		}
	}
	slices.SortFunc(foreignKeys, func(a, b *foreignKeyState) int {
		return strings.Compare(a.name, b.name)
	})
	return foreignKeys
}

func (t *tableState) toStringWithFKs(buf *strings.Builder, foreignKeys []*foreignKeyState) error {
	if _, err := fmt.Fprintf(buf, "CREATE TABLE %s (\n  ", t.qualifiedName()); err != nil {
		return err
	}
	columns := sortedColumns(t.columns)
	for i, column := range columns {
		if i > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := column.toString(buf); err != nil {
			return err
		}
	}

	indexes := make([]*indexState, 0, len(t.indexes))
	for _, index := range t.indexes {
		if index.primary || index.unique {
			indexes = append(indexes, index)
		}
	}
	slices.SortFunc(indexes, func(a, b *indexState) int {
		if a.primary && !b.primary {
			return -1
		}
		if !a.primary && b.primary {
			return 1
		}
		return strings.Compare(a.name, b.name)
	})

	for i, index := range indexes {
		if i+len(columns) > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := index.toString(buf); err != nil {
			return err
		}
	}

	for i, fk := range foreignKeys {
		if i+len(columns)+len(indexes) > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := fk.toString(buf); err != nil {
			return err
		}
	}

	if _, err := buf.WriteString("\n);\n"); err != nil {
		return err
	}
	return nil
}

func newTableState(id int, name string) *tableState {
	return &tableState{
		id:          id,
		name:        name,
		columns:     make(map[string]*columnState),
		indexes:     make(map[string]*indexState),
		foreignKeys: make(map[string]*foreignKeyState),
	}
}

// sortedColumns returns a table's columns in declaration order. Both the
// CREATE TABLE body and the trailing COMMENT ON COLUMN statements go through
// it, so the two agree and neither depends on map iteration order.
func sortedColumns(columns map[string]*columnState) []*columnState {
	sorted := make([]*columnState, 0, len(columns))
	for _, column := range columns {
		sorted = append(sorted, column)
	}
	slices.SortFunc(sorted, func(a, b *columnState) int {
		if a.id < b.id {
			return -1
		}
		if a.id > b.id {
			return 1
		}
		return 0
	})
	return sorted
}

func convertToTableState(id int, schemaName string, table *metadatapb.TableMetadata) *tableState {
	state := newTableState(id, table.Name)
	state.schema = schemaName
	state.comment = table.Comment
	for i, column := range table.Columns {
		if column == nil {
			continue
		}
		state.columns[column.Name] = convertToColumnState(i, column)
	}
	for i, index := range table.Indexes {
		if index == nil {
			continue
		}
		state.indexes[index.Name] = convertToIndexState(i, index)
	}
	for i, fk := range table.ForeignKeys {
		if fk == nil {
			continue
		}
		state.foreignKeys[fk.Name] = convertToForeignKeyState(i, schemaName, fk)
	}
	return state
}

type foreignKeyState struct {
	id                int
	name              string
	columns           []string
	referencedSchema  string
	referencedTable   string
	referencedColumns []string
}

func convertToForeignKeyState(id int, tableSchema string, foreignKey *metadatapb.ForeignKeyMetadata) *foreignKeyState {
	refSchema := foreignKey.ReferencedSchema
	if refSchema == "" {
		refSchema = tableSchema
	}
	return &foreignKeyState{
		id:                id,
		name:              foreignKey.Name,
		columns:           foreignKey.Columns,
		referencedSchema:  refSchema,
		referencedTable:   foreignKey.ReferencedTable,
		referencedColumns: foreignKey.ReferencedColumns,
	}
}

func (f *foreignKeyState) toString(buf *strings.Builder) error {
	n := min(len(f.columns), len(f.referencedColumns))
	if n == 0 {
		return nil
	}
	if f.name != "" {
		if _, err := fmt.Fprintf(buf, "CONSTRAINT %s ", quoteDefinitionIdentifier(f.name)); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("FOREIGN KEY ("); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(quoteDefinitionIdentifier(f.columns[i])); err != nil {
			return err
		}
	}
	// The sync strips the qualifier off ReferencedTable and keeps it in
	// ReferencedSchema, so an unqualified target here would bind to whatever
	// the search path resolves rather than the table the constraint names.
	referenced := qualifyName(f.referencedSchema, f.referencedTable)
	if _, err := fmt.Fprintf(buf, ") REFERENCES %s(", referenced); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(quoteDefinitionIdentifier(f.referencedColumns[i])); err != nil {
			return err
		}
	}
	_, err := buf.WriteString(")")
	return err
}

type indexState struct {
	id      int
	name    string
	keys    []string
	lengths []int64
	primary bool
	unique  bool
	tp      string
}

func convertToIndexState(id int, index *metadatapb.IndexMetadata) *indexState {
	return &indexState{
		id:      id,
		name:    index.Name,
		keys:    index.Expressions,
		lengths: index.KeyLength,
		primary: index.Primary,
		unique:  index.Unique,
		tp:      index.Type,
	}
}

func (i *indexState) toString(buf *strings.Builder) error {
	var keyword string
	switch {
	case i.primary:
		keyword = "PRIMARY KEY ("
	case i.unique:
		keyword = "UNIQUE ("
	default:
		return nil
	}
	if i.name != "" {
		if _, err := fmt.Fprintf(buf, "CONSTRAINT %s ", quoteDefinitionIdentifier(i.name)); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(keyword); err != nil {
		return err
	}
	for idx, key := range i.keys {
		if idx > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		// Written verbatim: these come from pg_get_indexdef(.., true), which
		// already quotes what needs it and can be an expression rather than
		// a column, so quoting again would name something else.
		if _, err := buf.WriteString(key); err != nil {
			return err
		}
	}
	_, err := buf.WriteString(")")
	return err
}

type defaultValue interface {
	toString() string
}

type defaultValueNull struct{}

func (*defaultValueNull) toString() string {
	return "NULL"
}

type defaultValueString struct {
	value string
}

func (d *defaultValueString) toString() string {
	return fmt.Sprintf("'%s'", escapeSingleQuote(d.value))
}

type defaultValueExpression struct {
	value string
}

func (d *defaultValueExpression) toString() string {
	return d.value
}

type columnState struct {
	id                int
	name              string
	tp                string
	collationSchema   string
	collation         string
	defaultValue      defaultValue
	comment           string
	nullable          bool
	isIdentity        bool
	identitySeed      int64
	identityIncrement int64
}

func (c *columnState) toString(buf *strings.Builder) error {
	if _, err := fmt.Fprintf(buf, "%s ", quoteDefinitionIdentifier(c.name)); err != nil {
		return err
	}
	if _, err := buf.WriteString(c.tp); err != nil {
		return err
	}
	if c.collation != "" {
		if _, err := fmt.Fprintf(buf, " COLLATE %s", formatRedshiftCollation(c.collationSchema, c.collation)); err != nil {
			return err
		}
	}
	if c.isIdentity {
		seed, step := c.identitySeed, c.identityIncrement
		if step == 0 {
			seed, step = 1, 1
		}
		if _, err := fmt.Fprintf(buf, " IDENTITY(%d, %d)", seed, step); err != nil {
			return err
		}
	}
	if !c.nullable {
		if _, err := buf.WriteString(" NOT NULL"); err != nil {
			return err
		}
	}
	if !c.isIdentity && c.defaultValue != nil {
		if _, err := fmt.Fprintf(buf, " DEFAULT %s", c.defaultValue.toString()); err != nil {
			return err
		}
	}
	return nil
}

func convertToColumnState(id int, column *metadatapb.ColumnMetadata) *columnState {
	result := &columnState{
		id:                id,
		name:              column.Name,
		tp:                column.Type,
		collationSchema:   column.CollationSchema,
		collation:         column.Collation,
		nullable:          column.Nullable,
		comment:           column.Comment,
		isIdentity:        column.IsIdentity,
		identitySeed:      column.IdentitySeed,
		identityIncrement: column.IdentityIncrement,
	}
	if column.Default == "NULL" {
		result.defaultValue = &defaultValueNull{}
	} else if column.Default != "" {
		if isExpression(column.Default) || isQuotedSQLLiteral(column.Default) || isNumericOrBoolLiteral(column.Default) {
			result.defaultValue = &defaultValueExpression{value: column.Default}
		} else {
			result.defaultValue = &defaultValueString{value: column.Default}
		}
	}
	return result
}

// formatRedshiftCollation emits Redshift's unquoted collation keywords
// (CASE_SENSITIVE, CS, CASE_INSENSITIVE, CI) verbatim and quotes other
// collation names as identifiers.
func formatRedshiftCollation(schema, name string) string {
	if schema == "" {
		switch strings.ToUpper(name) {
		case "CASE_SENSITIVE", "CS", "CASE_INSENSITIVE", "CI":
			return strings.ToUpper(name)
		}
	}
	return qualifyName(schema, name)
}

func isExpression(value string) bool {
	trimmed := strings.TrimSpace(value)
	switch strings.ToUpper(trimmed) {
	case "SYSDATE", "CURRENT_DATE", "CURRENT_TIME", "CURRENT_TIMESTAMP",
		"CURRENT_USER", "CURRENT_SCHEMA", "LOCALTIME", "LOCALTIMESTAMP",
		"SESSION_USER", "USER":
		return true
	}
	stmts, err := pgparser.Parse("SELECT " + trimmed)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	sel, ok := stmt.(*nodes.SelectStmt)
	if !ok || sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return false
	}
	rt, ok := sel.TargetList.Items[0].(*nodes.ResTarget)
	if !ok || rt.Name != "" {
		return false
	}
	switch rt.Val.(type) {
	case *nodes.FuncCall, *nodes.TypeCast, *nodes.A_Expr, *nodes.SQLValueFunction, *nodes.CoalesceExpr, *nodes.CaseExpr:
		// Column defaults cannot reference columns; rejecting ColumnRef keeps
		// raw strings like "pending (old)" from matching as FuncCall(ColumnRef).
		hasColRef := false
		nodes.Inspect(rt.Val, func(n nodes.Node) bool {
			if _, ok := n.(*nodes.ColumnRef); ok {
				hasColRef = true
				return false
			}
			return true
		})
		return !hasColRef
	default:
		return false
	}
}

// isQuotedSQLLiteral reports whether s is already a single-quoted SQL string
// literal (such as 'active' or 'it”s ok'), so column defaults copied from
// information_schema.columns.column_default are preserved verbatim.
func isQuotedSQLLiteral(s string) bool {
	if len(s) < 2 || s[0] != '\'' || s[len(s)-1] != '\'' {
		return false
	}
	for i := 1; i < len(s)-1; i++ {
		if s[i] == '\'' {
			if i+1 >= len(s)-1 || s[i+1] != '\'' {
				return false
			}
			i++
		}
	}
	return true
}

func isNumericOrBoolLiteral(s string) bool {
	if strings.EqualFold(s, "true") || strings.EqualFold(s, "false") {
		return true
	}
	if s == "" {
		return false
	}
	i := 0
	if s[i] == '+' || s[i] == '-' {
		i++
	}
	digits := false
	dot := false
	for ; i < len(s); i++ {
		switch {
		case s[i] >= '0' && s[i] <= '9':
			digits = true
		case s[i] == '.' && !dot:
			dot = true
		default:
			return false
		}
	}
	return digits
}

type viewState struct {
	id         int
	schema     string
	name       string
	definition string
	comment    string
}

func convertToViewState(id int, schemaName string, view *metadatapb.ViewMetadata) *viewState {
	return &viewState{
		id:         id,
		schema:     schemaName,
		name:       view.Name,
		definition: view.Definition,
		comment:    view.Comment,
	}
}

func convertToMaterializedViewState(id int, schemaName string, mv *metadatapb.MaterializedViewMetadata) *viewState {
	return &viewState{
		id:         id,
		schema:     schemaName,
		name:       mv.Name,
		definition: mv.Definition,
		comment:    mv.Comment,
	}
}

func (v *viewState) toString(buf io.StringWriter) error {
	trimmed := strings.TrimSpace(v.definition)
	if trimmed == "" {
		return fmt.Errorf("view %s has no definition to restore", v.qualifiedName())
	}
	var stmt string
	if isCompleteCreateView(trimmed) {
		stmt = trimmed
	} else {
		stmt = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", v.qualifiedName(), v.definition)
	}
	if !strings.HasSuffix(stmt, ";") {
		stmt += ";"
	}
	stmt += "\n"
	_, err := buf.WriteString(stmt)
	return err
}

func (v *viewState) toMaterializedString(buf io.StringWriter) error {
	trimmed := strings.TrimSpace(v.definition)
	if trimmed == "" {
		return fmt.Errorf("materialized view %s has no definition to restore", v.qualifiedName())
	}
	var stmt string
	if isCompleteCreateMaterializedView(trimmed) {
		stmt = trimmed
	} else {
		stmt = fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS %s", v.qualifiedName(), v.definition)
	}
	if !strings.HasSuffix(stmt, ";") {
		stmt += ";"
	}
	stmt += "\n"
	_, err := buf.WriteString(stmt)
	return err
}

// isCompleteCreateView reports whether definition already starts with a full
// CREATE [OR REPLACE] VIEW statement (as SHOW VIEW returns for late-binding
// views with WITH NO SCHEMA BINDING).
func isCompleteCreateView(definition string) bool {
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	_, ok := stmt.(*nodes.ViewStmt)
	return ok
}

func isCompleteCreateMaterializedView(definition string) bool {
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	ctas, ok := stmt.(*nodes.CreateTableAsStmt)
	return ok && ctas.Objtype == nodes.OBJECT_MATVIEW
}
