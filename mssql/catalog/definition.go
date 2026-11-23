package catalog

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// DefinitionOptions selects the form DatabaseDefinition renders.
type DefinitionOptions struct {
	// PrintHeader is accepted for API parity with the pg and mysql catalogs.
	// SQL Server's dump has no session preamble, so it renders nothing extra.
	PrintHeader bool
	// SDLFormat is not implemented for SQL Server: Bytebase's schema/mssql has
	// no declarative form, and this catalog keeps that boundary.
	SDLFormat bool
}

// DatabaseDefinition renders a synced snapshot as the T-SQL script that
// recreates it. It is Bytebase's schema/mssql rendering, byte for byte: schemas
// first, then tables in foreign-key order, views in dependency order, and
// functions and procedures last.
func DatabaseDefinition(to *metadata.DatabaseSchemaMetadata, opts DefinitionOptions) (string, error) {
	if to == nil {
		return "", nil
	}
	if opts.SDLFormat {
		return "", fmt.Errorf("mssql catalog: SDL format is not supported")
	}

	var buf strings.Builder

	for _, schema := range to.Schemas {
		if schema.Name != defaultSchema {
			fmt.Fprintf(&buf, "CREATE SCHEMA %s;\nGO\n\n", quote(schema.Name))
		}
	}

	writeSynonyms(&buf, to.Schemas)
	writeSequences(&buf, to.Schemas)
	hasPreamble := hasSynonyms(to.Schemas) || hasSequences(to.Schemas)

	tableFunctions := functionsReferencedByTables(to.Schemas)
	lastWasFunction := false
	if len(tableFunctions) > 0 && tableFunctionsDependOnTables(to.Schemas, tableFunctions) {
		lastWasFunction = writeTablesAndFunctionsInOrder(&buf, to.Schemas, tableFunctions, hasPreamble)
	} else {
		if len(tableFunctions) > 0 {
			if hasPreamble {
				buf.WriteString("GO\n\n")
			}
			writeReferencedFunctions(&buf, to.Schemas, tableFunctions)
			lastWasFunction = true
		}

		if hasTables(to.Schemas) {
			writeAllTablesInOrder(&buf, to.Schemas)
			lastWasFunction = false
		}
	}

	viewFunctions := functionsReferencedByViews(to.Schemas)
	for name := range tableFunctions {
		delete(viewFunctions, name)
	}
	if len(viewFunctions) > 0 {
		if !lastWasFunction && (hasTables(to.Schemas) || hasPreamble) {
			buf.WriteString("GO\n\n")
		}
		writeReferencedFunctions(&buf, to.Schemas, viewFunctions)
		lastWasFunction = true
	}

	if hasViews(to.Schemas) {
		if len(viewFunctions) == 0 && !lastWasFunction && (hasTables(to.Schemas) || hasPreamble) {
			buf.WriteString("GO\n\n")
		}
		writeAllViewsInOrder(&buf, to.Schemas)
		lastWasFunction = false
	}

	if !hasViews(to.Schemas) && !lastWasFunction && (hasTables(to.Schemas) || hasPreamble) && (hasRoutines(to.Schemas) || hasTriggers(to.Schemas)) {
		buf.WriteString("GO\n\n")
	}

	writeTriggers(&buf, to.Schemas)

	remaining := make(map[string]bool)
	for name := range tableFunctions {
		remaining[name] = true
	}
	for name := range viewFunctions {
		remaining[name] = true
	}
	for _, schema := range to.Schemas {
		writeFunctionsAndProcedures(&buf, schema, remaining)
	}

	return buf.String(), nil
}

type functionWithSchema struct {
	schema   string
	function *metadata.FunctionMetadata
}

func writeReferencedFunctions(out *strings.Builder, schemas []*metadata.SchemaMetadata, names map[string]bool) {
	for _, fws := range sortFunctionsAcrossSchemas(schemas, names) {
		writeFunction(out, fws.schema, fws.function)
	}
}

// sortFunctionsAcrossSchemas topologically orders functions across all schemas
// so a callee is always created before any function that calls it.
func sortFunctionsAcrossSchemas(schemas []*metadata.SchemaMetadata, filter map[string]bool) []functionWithSchema {
	var all []functionWithSchema
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() {
				continue
			}
			id := objectID(schema.Name, function.Name)
			if filter != nil && !filter[id] {
				continue
			}
			all = append(all, functionWithSchema{schema: schema.Name, function: function})
		}
	}
	if len(all) <= 1 {
		return all
	}
	graph := newGraph()
	byID := make(map[string]functionWithSchema, len(all))
	for _, fws := range all {
		id := objectID(fws.schema, fws.function.Name)
		graph.addNode(id)
		byID[id] = fws
	}
	for _, fws := range all {
		id := objectID(fws.schema, fws.function.Name)
		for _, other := range all {
			otherID := objectID(other.schema, other.function.Name)
			if otherID == id {
				continue
			}
			if definitionReferencesFunction(fws.function.Definition, fws.schema, other.schema, other.function.Name) {
				graph.addEdge(otherID, id)
			}
		}
	}
	order, err := graph.topologicalSort()
	if err != nil {
		return all
	}
	result := make([]functionWithSchema, 0, len(all))
	for _, id := range order {
		if fws, ok := byID[id]; ok {
			result = append(result, fws)
		}
	}
	return result
}

// sortFunctions puts a function after any function its body names.
func sortFunctions(schemaName string, functions []*metadata.FunctionMetadata) []*metadata.FunctionMetadata {
	graph := newGraph()
	byName := make(map[string]*metadata.FunctionMetadata)
	for _, function := range functions {
		if function.GetSkipDump() {
			continue
		}
		id := objectID(schemaName, function.Name)
		graph.addNode(id)
		byName[id] = function
	}
	for _, function := range functions {
		if function.GetSkipDump() {
			continue
		}
		id := objectID(schemaName, function.Name)
		for _, other := range functions {
			if other == function || other.GetSkipDump() {
				continue
			}
			if definitionReferencesFunction(function.Definition, schemaName, schemaName, other.Name) {
				graph.addEdge(objectID(schemaName, other.Name), id)
			}
		}
	}
	order, err := graph.topologicalSort()
	if err != nil {
		return functions
	}
	result := make([]*metadata.FunctionMetadata, 0, len(functions))
	for _, id := range order {
		if function := byName[id]; function != nil {
			result = append(result, function)
		}
	}
	return result
}

// expandFunctionDependencies adds every function transitively called by the
// functions in seeds across all schemas.
func expandFunctionDependencies(schemas []*metadata.SchemaMetadata, seeds map[string]bool) map[string]bool {
	if len(seeds) == 0 {
		return seeds
	}
	var all []functionWithSchema
	byID := make(map[string]functionWithSchema)
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() {
				continue
			}
			id := objectID(schema.Name, function.Name)
			fws := functionWithSchema{schema: schema.Name, function: function}
			all = append(all, fws)
			byID[id] = fws
		}
	}
	result := make(map[string]bool, len(seeds))
	var queue []string
	for id := range seeds {
		result[id] = true
		queue = append(queue, id)
	}
	slices.Sort(queue)
	for len(queue) > 0 {
		curID := queue[0]
		queue = queue[1:]
		cur, ok := byID[curID]
		if !ok {
			continue
		}
		for _, other := range all {
			otherID := objectID(other.schema, other.function.Name)
			if result[otherID] || otherID == curID {
				continue
			}
			if definitionReferencesFunction(cur.function.Definition, cur.schema, other.schema, other.function.Name) {
				result[otherID] = true
				queue = append(queue, otherID)
			}
		}
	}
	return result
}

// functionsReferencedByTables names the functions a column default, computed
// expression, or check constraint mentions, together with their transitive
// function dependencies.
func functionsReferencedByTables(schemas []*metadata.SchemaMetadata) map[string]bool {
	type exprCtx struct {
		schema string
		expr   string
	}
	var expressions []exprCtx
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			if table.GetSkipDump() {
				continue
			}
			for _, column := range table.Columns {
				if column.Default != "" {
					expressions = append(expressions, exprCtx{schema: schema.Name, expr: column.Default})
				}
				if isComputedColumnType(column.Type) {
					expressions = append(expressions, exprCtx{schema: schema.Name, expr: column.Type})
				}
			}
			for _, check := range table.CheckConstraints {
				if check.Expression != "" {
					expressions = append(expressions, exprCtx{schema: schema.Name, expr: check.Expression})
				}
			}
		}
	}
	referenced := make(map[string]bool)
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() {
				continue
			}
			for _, item := range expressions {
				if expressionReferencesFunction(item.expr, item.schema, schema.Name, function.Name) {
					referenced[objectID(schema.Name, function.Name)] = true
					break
				}
			}
		}
	}
	return expandFunctionDependencies(schemas, referenced)
}

// functionsReferencedByViews names the functions a view body mentions, together
// with their transitive function dependencies.
func functionsReferencedByViews(schemas []*metadata.SchemaMetadata) map[string]bool {
	type viewCtx struct {
		schema string
		def    string
	}
	var definitions []viewCtx
	for _, schema := range schemas {
		for _, view := range schema.Views {
			if view.GetSkipDump() {
				continue
			}
			definitions = append(definitions, viewCtx{schema: schema.Name, def: view.Definition})
		}
	}
	referenced := make(map[string]bool)
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() {
				continue
			}
			for _, item := range definitions {
				if expressionReferencesFunction(item.def, item.schema, schema.Name, function.Name) {
					referenced[objectID(schema.Name, function.Name)] = true
					break
				}
			}
		}
	}
	return expandFunctionDependencies(schemas, referenced)
}

// tableFunctionsDependOnTables reports whether any hoisted table function is
// schema-bound to a table in schemas, requiring tables and functions to be
// interleaved in dependency order.
func tableFunctionsDependOnTables(schemas []*metadata.SchemaMetadata, tableFunctions map[string]bool) bool {
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() || !tableFunctions[objectID(schema.Name, function.Name)] {
				continue
			}
			if !isSchemaBound(function.Definition) {
				continue
			}
			for _, ts := range schemas {
				for _, table := range ts.Tables {
					if table.GetSkipDump() {
						continue
					}
					if definitionReferencesTable(function.Definition, schema.Name, ts.Name, table.Name) {
						return true
					}
				}
			}
		}
	}
	return false
}

// writeTablesAndFunctionsInOrder interleaves tables and table-referenced
// functions in dependency order so a schema-bound function follows the table it
// reads while still preceding any table default/check/computed column that
// calls it. Returns true if the last statement emitted was a function (already
// terminated by GO).
func writeTablesAndFunctionsInOrder(out *strings.Builder, schemas []*metadata.SchemaMetadata, tableFunctions map[string]bool, hasPreamble bool) bool {
	var allTables []tableWithSchema
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			if table.GetSkipDump() {
				continue
			}
			allTables = append(allTables, tableWithSchema{schema: schema.Name, table: table})
		}
	}
	var allFuncs []functionWithSchema
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			if function.GetSkipDump() || !tableFunctions[objectID(schema.Name, function.Name)] {
				continue
			}
			allFuncs = append(allFuncs, functionWithSchema{schema: schema.Name, function: function})
		}
	}

	buildGraph := func(includeFKs bool) (*graph, map[string]tableWithSchema, map[string]functionWithSchema) {
		g := newGraph()
		tableMap := make(map[string]tableWithSchema, len(allTables))
		tableIDByObj := make(map[string]string, len(allTables))
		funcMap := make(map[string]functionWithSchema, len(allFuncs))
		for _, tws := range allTables {
			obj := objectID(tws.schema, tws.table.Name)
			node := "T\x00" + obj
			g.addNode(node)
			tableMap[node] = tws
			tableIDByObj[obj] = node
		}
		for _, fws := range allFuncs {
			node := "F\x00" + objectID(fws.schema, fws.function.Name)
			g.addNode(node)
			funcMap[node] = fws
		}
		if includeFKs {
			for _, tws := range allTables {
				tNode := "T\x00" + objectID(tws.schema, tws.table.Name)
				for _, fk := range tws.table.ForeignKeys {
					refObj := objectID(fk.ReferencedSchema, fk.ReferencedTable)
					if refNode, ok := tableIDByObj[refObj]; ok && refNode != tNode {
						g.addEdge(refNode, tNode)
					}
				}
			}
		}
		// Function -> Function edges.
		for _, fws := range allFuncs {
			fNode := "F\x00" + objectID(fws.schema, fws.function.Name)
			for _, other := range allFuncs {
				otherNode := "F\x00" + objectID(other.schema, other.function.Name)
				if otherNode == fNode {
					continue
				}
				if definitionReferencesFunction(fws.function.Definition, fws.schema, other.schema, other.function.Name) {
					g.addEdge(otherNode, fNode)
				}
			}
		}
		// Schema-bound Function reads Table -> Table must precede Function.
		for _, fws := range allFuncs {
			if !isSchemaBound(fws.function.Definition) {
				continue
			}
			fNode := "F\x00" + objectID(fws.schema, fws.function.Name)
			for _, tws := range allTables {
				if definitionReferencesTable(fws.function.Definition, fws.schema, tws.schema, tws.table.Name) {
					g.addEdge("T\x00"+objectID(tws.schema, tws.table.Name), fNode)
				}
			}
		}
		// Table uses Function -> Function must precede Table.
		for _, tws := range allTables {
			tNode := "T\x00" + objectID(tws.schema, tws.table.Name)
			for _, fws := range allFuncs {
				if tableReferencesFunction(tws.schema, tws.table, fws.schema, fws.function.Name) {
					g.addEdge("F\x00"+objectID(fws.schema, fws.function.Name), tNode)
				}
			}
		}
		return g, tableMap, funcMap
	}

	g, tableMap, funcMap := buildGraph(true)
	order, err := g.topologicalSort()
	deferFKs := false
	if err != nil {
		g, tableMap, funcMap = buildGraph(false)
		order, err = g.topologicalSort()
		deferFKs = true
	}

	needGoBeforeFn := hasPreamble
	lastWasFunction := false
	if err == nil {
		for _, node := range order {
			if tws, ok := tableMap[node]; ok {
				if deferFKs {
					writeTableWithoutForeignKeys(out, tws.schema, tws.table)
				} else {
					writeTable(out, tws.schema, tws.table)
				}
				needGoBeforeFn = true
				lastWasFunction = false
			} else if fws, ok := funcMap[node]; ok {
				if needGoBeforeFn {
					out.WriteString("GO\n\n")
				}
				writeFunction(out, fws.schema, fws.function)
				needGoBeforeFn = false
				lastWasFunction = true
			}
		}
		if deferFKs {
			for _, tws := range allTables {
				for _, fk := range tws.table.ForeignKeys {
					writeDeferredForeignKey(out, tws.schema, tws.table.Name, fk)
					lastWasFunction = false
				}
			}
		}
		return lastWasFunction
	}

	// Cycle involves a table's DEFAULT/CHECK constraint and a schema-bound
	// function reading the same table (or cycle of tables). Emit tables without
	// function-referencing defaults/checks (and without FKs if cyclic), emit
	// functions, then attach deferred defaults, checks, and FKs via ALTER TABLE.
	sortedTables, fkCyclic := sortTablesByDependencies(allTables)
	type deferredDefault struct {
		schema string
		table  string
		column *metadata.ColumnMetadata
	}
	type deferredCheck struct {
		schema string
		table  string
		check  *metadata.CheckConstraintMetadata
	}
	var deferredDefaults []deferredDefault
	var deferredChecks []deferredCheck
	for _, tws := range sortedTables {
		deferDefaultCols := make(map[string]bool)
		for _, col := range tws.table.Columns {
			if col.GetDefault() != "" && expressionReferencesAnyFunction(col.GetDefault(), tws.schema, allFuncs) {
				deferDefaultCols[col.Name] = true
				deferredDefaults = append(deferredDefaults, deferredDefault{schema: tws.schema, table: tws.table.Name, column: col})
			}
		}
		deferCheckNames := make(map[string]bool)
		for _, check := range tws.table.CheckConstraints {
			if expressionReferencesAnyFunction(check.Expression, tws.schema, allFuncs) {
				deferCheckNames[check.Name] = true
				deferredChecks = append(deferredChecks, deferredCheck{schema: tws.schema, table: tws.table.Name, check: check})
			}
		}
		writeTableBodyFiltered(out, tws.schema, tws.table, !fkCyclic, deferDefaultCols, deferCheckNames)
		needGoBeforeFn = true
		lastWasFunction = false
	}
	if len(allFuncs) > 0 {
		if needGoBeforeFn {
			out.WriteString("GO\n\n")
		}
		writeReferencedFunctions(out, schemas, tableFunctions)
		lastWasFunction = true
	}
	for _, dd := range deferredDefaults {
		if dd.column.DefaultConstraintName != "" {
			fmt.Fprintf(out, "ALTER TABLE %s.%s ADD CONSTRAINT %s DEFAULT %s FOR %s;\n\n",
				quote(dd.schema), quote(dd.table), quote(dd.column.DefaultConstraintName), dd.column.GetDefault(), quote(dd.column.Name))
		} else {
			fmt.Fprintf(out, "ALTER TABLE %s.%s ADD DEFAULT %s FOR %s;\n\n",
				quote(dd.schema), quote(dd.table), dd.column.GetDefault(), quote(dd.column.Name))
		}
		lastWasFunction = false
	}
	for _, dc := range deferredChecks {
		fmt.Fprintf(out, "ALTER TABLE %s.%s ADD CONSTRAINT %s CHECK %s;\n\n",
			quote(dc.schema), quote(dc.table), quote(dc.check.Name), checkExpression(dc.check.Expression))
		lastWasFunction = false
	}
	if fkCyclic {
		for _, tws := range sortedTables {
			for _, fk := range tws.table.ForeignKeys {
				writeDeferredForeignKey(out, tws.schema, tws.table.Name, fk)
				lastWasFunction = false
			}
		}
	}
	return lastWasFunction
}

func tableReferencesFunction(tableSchema string, table *metadata.TableMetadata, fnSchema, fnName string) bool {
	if table == nil {
		return false
	}
	for _, column := range table.Columns {
		if column.Default != "" && expressionReferencesFunction(column.Default, tableSchema, fnSchema, fnName) {
			return true
		}
		if isComputedColumnType(column.Type) && expressionReferencesFunction(column.Type, tableSchema, fnSchema, fnName) {
			return true
		}
	}
	for _, check := range table.CheckConstraints {
		if check.Expression != "" && expressionReferencesFunction(check.Expression, tableSchema, fnSchema, fnName) {
			return true
		}
	}
	return false
}

func expressionReferencesAnyFunction(expression, defaultSchema string, funcs []functionWithSchema) bool {
	for _, fws := range funcs {
		if expressionReferencesFunction(expression, defaultSchema, fws.schema, fws.function.Name) {
			return true
		}
	}
	return false
}

func hasSynonyms(schemas []*metadata.SchemaMetadata) bool {
	for _, schema := range schemas {
		if len(schema.Synonyms) > 0 {
			return true
		}
	}
	return false
}

func hasSequences(schemas []*metadata.SchemaMetadata) bool {
	for _, schema := range schemas {
		for _, sequence := range schema.Sequences {
			if !sequence.GetSkipDump() {
				return true
			}
		}
	}
	return false
}

func hasTriggers(schemas []*metadata.SchemaMetadata) bool {
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			if len(table.Triggers) > 0 {
				return true
			}
		}
		for _, view := range schema.Views {
			if len(view.Triggers) > 0 {
				return true
			}
		}
	}
	return false
}

func hasRoutines(schemas []*metadata.SchemaMetadata) bool {
	for _, schema := range schemas {
		if len(schema.Functions) > 0 || len(schema.Procedures) > 0 {
			return true
		}
	}
	return false
}

// writeSynonyms renders the aliases the sync reports before objects that may
// use them.
func writeSynonyms(out *strings.Builder, schemas []*metadata.SchemaMetadata) {
	for _, schema := range schemas {
		for _, synonym := range schema.Synonyms {
			fmt.Fprintf(out, "CREATE SYNONYM %s.%s FOR %s;\n\n", quote(schema.Name), quote(synonym.Name), synonymTarget(synonym))
		}
	}
}

// synonymTarget keeps interior empty components, so [db]..[object] does not
// become the different [db].[object].
func synonymTarget(synonym *metadata.SynonymMetadata) string {
	components := []string{synonym.TargetServer, synonym.TargetDatabase, synonym.TargetSchema, synonym.TargetName}
	start := 0
	for i := 0; i < 3; i++ {
		if components[i] != "" {
			start = i
			break
		}
	}
	var parts []string
	for i := start; i < len(components); i++ {
		if components[i] == "" {
			parts = append(parts, "")
		} else {
			parts = append(parts, quote(components[i]))
		}
	}
	return strings.Join(parts, ".")
}

// writeSequences renders CREATE SEQUENCE before tables, since a table may use
// NEXT VALUE FOR in a default.
func writeSequences(out *strings.Builder, schemas []*metadata.SchemaMetadata) {
	for _, schema := range schemas {
		for _, sequence := range schema.Sequences {
			if sequence.GetSkipDump() {
				continue
			}
			fmt.Fprintf(out, "CREATE SEQUENCE %s.%s", quote(schema.Name), quote(sequence.Name))
			if sequence.DataType != "" {
				fmt.Fprintf(out, "\n    AS %s", sequence.DataType)
			}
			if sequence.Start != "" {
				fmt.Fprintf(out, "\n    START WITH %s", sequence.Start)
			}
			if sequence.Increment != "" {
				fmt.Fprintf(out, "\n    INCREMENT BY %s", sequence.Increment)
			}
			if sequence.MinValue != "" {
				fmt.Fprintf(out, "\n    MINVALUE %s", sequence.MinValue)
			}
			if sequence.MaxValue != "" {
				fmt.Fprintf(out, "\n    MAXVALUE %s", sequence.MaxValue)
			}
			if sequence.Cycle {
				out.WriteString("\n    CYCLE")
			}
			if sequence.CacheSize != "" {
				if sequence.CacheSize == "NO CACHE" {
					out.WriteString("\n    NO CACHE")
				} else {
					fmt.Fprintf(out, "\n    CACHE %s", sequence.CacheSize)
				}
			}
			out.WriteString(";\n\n")
			if sequence.Comment != "" {
				fmt.Fprintf(out, "%s;\nGO\n\n", generateSequenceCommentSQL("ADD", schema.Name, sequence.Name, sequence.Comment))
			}
		}
	}
}

// Definition renders the catalog's own snapshot.
func (c *Catalog) Definition(opts DefinitionOptions) (string, error) {
	return DatabaseDefinition(c.Metadata(), opts)
}

// writeTriggers renders the synced trigger definitions after their tables and
// views. The sync stores each as the CREATE TRIGGER text the module holds.
func writeTriggers(out *strings.Builder, schemas []*metadata.SchemaMetadata) {
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			if table.GetSkipDump() {
				continue
			}
			for _, trigger := range table.Triggers {
				if trigger.Body == "" || trigger.SkipDump {
					continue
				}
				body := qualifyTriggerDefinition(creationDefinition(trigger.Body), schema.Name, table.Name, false)
				fmt.Fprintf(out, "%s\n\nGO\n\n", strings.TrimRight(body, ";"))
			}
		}
		for _, view := range schema.Views {
			if view.GetSkipDump() {
				continue
			}
			for _, trigger := range view.Triggers {
				if trigger.Body == "" || trigger.SkipDump {
					continue
				}
				body := qualifyTriggerDefinition(creationDefinition(trigger.Body), schema.Name, view.Name, false)
				fmt.Fprintf(out, "%s\n\nGO\n\n", strings.TrimRight(body, ";"))
			}
		}
	}
}

func writeFunctionsAndProcedures(out *strings.Builder, schema *metadata.SchemaMetadata, early map[string]bool) {
	for _, function := range sortFunctions(schema.Name, schema.Functions) {
		if function.GetSkipDump() || early[objectID(schema.Name, function.Name)] {
			continue
		}
		writeFunction(out, schema.Name, function)
	}
	for _, procedure := range schema.Procedures {
		if procedure.SkipDump {
			continue
		}
		writeProcedure(out, schema.Name, procedure)
	}
}

func hasViews(schemas []*metadata.SchemaMetadata) bool {
	for _, schema := range schemas {
		if len(schema.Views) > 0 {
			return true
		}
	}
	return false
}

func hasTables(schemas []*metadata.SchemaMetadata) bool {
	for _, schema := range schemas {
		if len(schema.Tables) > 0 {
			return true
		}
	}
	return false
}

func writeAllTablesInOrder(out *strings.Builder, schemas []*metadata.SchemaMetadata) {
	var allTables []tableWithSchema
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			if table.GetSkipDump() {
				continue
			}
			allTables = append(allTables, tableWithSchema{schema: schema.Name, table: table})
		}
	}
	if len(allTables) == 0 {
		return
	}
	sorted, cyclic := sortTablesByDependencies(allTables)
	for _, tws := range sorted {
		if cyclic {
			writeTableWithoutForeignKeys(out, tws.schema, tws.table)
		} else {
			writeTable(out, tws.schema, tws.table)
		}
	}
	if cyclic {
		for _, tws := range sorted {
			for _, fk := range tws.table.ForeignKeys {
				writeDeferredForeignKey(out, tws.schema, tws.table.Name, fk)
			}
		}
	}
}

type tableWithSchema struct {
	schema string
	table  *metadata.TableMetadata
}

func sortTablesByDependencies(tables []tableWithSchema) ([]tableWithSchema, bool) {
	if len(tables) <= 1 {
		return tables, false
	}
	tableMap := make(map[string]tableWithSchema)
	for _, tws := range tables {
		tableMap[objectID(tws.schema, tws.table.Name)] = tws
	}
	graph := newGraph()
	for _, tws := range tables {
		graph.addNode(objectID(tws.schema, tws.table.Name))
	}
	for _, tws := range tables {
		tableID := objectID(tws.schema, tws.table.Name)
		for _, fk := range tws.table.ForeignKeys {
			refID := objectID(fk.ReferencedSchema, fk.ReferencedTable)
			if refID == tableID {
				continue
			}
			if _, exists := tableMap[refID]; exists {
				graph.addEdge(refID, tableID)
			}
		}
	}
	sortedIDs, err := graph.topologicalSort()
	if err != nil {
		return tables, true
	}
	result := make([]tableWithSchema, 0, len(sortedIDs))
	for _, id := range sortedIDs {
		if tws, ok := tableMap[id]; ok {
			result = append(result, tws)
		}
	}
	return result, false
}

func writeAllViewsInOrder(out *strings.Builder, schemas []*metadata.SchemaMetadata) {
	var allViews []viewWithSchema
	for _, schema := range schemas {
		for _, view := range schema.Views {
			if view.GetSkipDump() {
				continue
			}
			allViews = append(allViews, viewWithSchema{schema: schema.Name, view: view})
		}
	}
	if len(allViews) == 0 {
		return
	}
	for _, vws := range sortViewsByDependencies(allViews) {
		writeView(out, vws.schema, vws.view)
	}
}

type viewWithSchema struct {
	schema string
	view   *metadata.ViewMetadata
}

func sortViewsByDependencies(views []viewWithSchema) []viewWithSchema {
	if len(views) <= 1 {
		return views
	}
	viewMap := make(map[string]viewWithSchema)
	for _, vws := range views {
		viewMap[objectID(vws.schema, vws.view.Name)] = vws
	}
	graph := newGraph()
	for _, vws := range views {
		graph.addNode(objectID(vws.schema, vws.view.Name))
	}
	for _, vws := range views {
		viewID := objectID(vws.schema, vws.view.Name)
		deps, err := getViewDependencies(vws.view.Definition, vws.schema)
		if err != nil {
			continue
		}
		for _, dep := range deps {
			if _, isView := viewMap[dep]; isView {
				graph.addEdge(dep, viewID)
			}
		}
	}
	sortedIDs, err := graph.topologicalSort()
	if err != nil {
		slices.SortFunc(views, func(a, b viewWithSchema) int {
			if a.schema != b.schema {
				return strings.Compare(a.schema, b.schema)
			}
			return strings.Compare(a.view.Name, b.view.Name)
		})
		return views
	}
	result := make([]viewWithSchema, 0, len(sortedIDs))
	for _, id := range sortedIDs {
		if vws, ok := viewMap[id]; ok {
			result = append(result, vws)
		}
	}
	return result
}

func writeTable(out *strings.Builder, schemaName string, table *metadata.TableMetadata) {
	writeTableBody(out, schemaName, table, true)
}

// writeTableWithoutForeignKeys renders a table whose foreign keys wait for an
// ALTER TABLE after every table exists, which is how a cycle is restored.
func writeTableWithoutForeignKeys(out *strings.Builder, schemaName string, table *metadata.TableMetadata) {
	writeTableBody(out, schemaName, table, false)
}

func writeTableBody(out *strings.Builder, schemaName string, table *metadata.TableMetadata, includeForeignKeys bool) {
	writeTableBodyFiltered(out, schemaName, table, includeForeignKeys, nil, nil)
}

func writeTableBodyFiltered(out *strings.Builder, schemaName string, table *metadata.TableMetadata, includeForeignKeys bool, omitDefaultCols, omitCheckNames map[string]bool) {
	fmt.Fprintf(out, "CREATE TABLE %s.%s (\n", quote(schemaName), quote(table.Name))
	for i, column := range table.Columns {
		if i != 0 {
			out.WriteString(",\n")
		}
		writeColumnFiltered(out, column, !omitDefaultCols[column.Name])
	}
	for _, key := range table.Indexes {
		if !key.IsConstraint {
			continue
		}
		out.WriteString(",\n")
		writeKey(out, key)
	}
	if includeForeignKeys {
		for _, fk := range table.ForeignKeys {
			out.WriteString(",\n")
			writeForeignKey(out, fk)
		}
	}
	for _, check := range table.CheckConstraints {
		if omitCheckNames[check.Name] {
			continue
		}
		out.WriteString(",\n")
		writeCheck(out, check)
	}
	out.WriteString("\n);\n\n")
	for _, index := range orderedIndexes(table.Indexes) {
		if index.IsConstraint {
			continue
		}
		writeIndex(out, schemaName, table.Name, index)
	}
	if table.Comment != "" {
		fmt.Fprintf(out, "%s;\nGO\n\n", generateTableCommentSQL("ADD", schemaName, table.Name, table.Comment))
	}
	for _, column := range table.Columns {
		if column.Comment == "" {
			continue
		}
		fmt.Fprintf(out, "%s;\nGO\n\n", generateColumnCommentSQL("ADD", schemaName, table.Name, column.Name, column.Comment))
	}
}

func writeDeferredForeignKey(out *strings.Builder, schemaName, tableName string, fk *metadata.ForeignKeyMetadata) {
	fmt.Fprintf(out, "ALTER TABLE %s.%s ADD CONSTRAINT %s FOREIGN KEY (", quote(schemaName), quote(tableName), quote(fk.Name))
	for i, column := range fk.Columns {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(quote(column))
	}
	fmt.Fprintf(out, ") REFERENCES %s.%s", quote(fk.ReferencedSchema), quote(fk.ReferencedTable))
	if len(fk.ReferencedColumns) > 0 {
		out.WriteString(" (")
		for i, column := range fk.ReferencedColumns {
			if i != 0 {
				out.WriteString(", ")
			}
			out.WriteString(quote(column))
		}
		out.WriteString(")")
	}
	if fk.OnDelete != "" {
		fmt.Fprintf(out, " ON DELETE %s", fk.OnDelete)
	}
	if fk.OnUpdate != "" {
		fmt.Fprintf(out, " ON UPDATE %s", fk.OnUpdate)
	}
	out.WriteString(";\n\n")
}

func writeXMLIndex(out *strings.Builder, schemaName, tableName string, index *metadata.IndexMetadata) {
	column := ""
	if len(index.Expressions) > 0 {
		column = quote(index.Expressions[0])
	}
	if strings.HasPrefix(index.Definition, "USING XML INDEX ") {
		rest := strings.TrimPrefix(index.Definition, "USING XML INDEX ")
		parent, secondaryType := rest, "PATH"
		if before, after, ok := strings.Cut(rest, " FOR "); ok {
			parent, secondaryType = before, after
		}
		fmt.Fprintf(out, "CREATE XML INDEX %s ON %s.%s (%s)\nUSING XML INDEX %s FOR %s;\n\n",
			quote(index.Name), quote(schemaName), quote(tableName), column, quote(parent), secondaryType)
		return
	}
	if index.Definition == "SELECTIVE" || strings.HasPrefix(index.Definition, "SELECTIVE ") {
		rest := strings.TrimSpace(strings.TrimPrefix(index.Definition, "SELECTIVE"))
		if rest != "" {
			fmt.Fprintf(out, "CREATE SELECTIVE XML INDEX %s ON %s.%s (%s)\n%s;\n\n",
				quote(index.Name), quote(schemaName), quote(tableName), column, rest)
			return
		}
		fmt.Fprintf(out, "CREATE SELECTIVE XML INDEX %s ON %s.%s (%s);\n\n",
			quote(index.Name), quote(schemaName), quote(tableName), column)
		return
	}
	fmt.Fprintf(out, "CREATE PRIMARY XML INDEX %s ON %s.%s (%s);\n\n",
		quote(index.Name), quote(schemaName), quote(tableName), column)
}

func isPrimaryXMLIndex(index *metadata.IndexMetadata) bool {
	return index.Type == "XML" && (index.Definition == "PRIMARY" || index.Definition == "SELECTIVE" || strings.HasPrefix(index.Definition, "SELECTIVE "))
}

func orderedIndexes(indexes []*metadata.IndexMetadata) []*metadata.IndexMetadata {
	out := make([]*metadata.IndexMetadata, 0, len(indexes))
	for _, index := range indexes {
		if isPrimaryXMLIndex(index) {
			out = append(out, index)
		}
	}
	for _, index := range indexes {
		if !isPrimaryXMLIndex(index) {
			out = append(out, index)
		}
	}
	return out
}

func writeClusteredColumnStoreIndex(out *strings.Builder, schemaName, tableName string, index *metadata.IndexMetadata) {
	fmt.Fprintf(out, "CREATE CLUSTERED COLUMNSTORE INDEX %s ON %s.%s;\n\n", quote(index.Name), quote(schemaName), quote(tableName))
}

func writeNonClusteredColumnStoreIndex(out *strings.Builder, schemaName, tableName string, index *metadata.IndexMetadata) {
	fmt.Fprintf(out, "CREATE NONCLUSTERED COLUMNSTORE INDEX %s ON %s.%s (\n", quote(index.Name), quote(schemaName), quote(tableName))
	for i, column := range index.Expressions {
		if i != 0 {
			out.WriteString(",\n")
		}
		fmt.Fprintf(out, "    %s", quote(column))
	}
	out.WriteString("\n)")
	if where := indexWhereClause(index.Definition); where != "" {
		fmt.Fprintf(out, "\nWHERE %s", where)
	}
	out.WriteString(";\n\n")
}

func writeSpatialIndex(out *strings.Builder, schemaName, tableName string, index *metadata.IndexMetadata) {
	out.WriteString(generateSpatialIndexDefinition(index, schemaName, tableName))
	out.WriteString(";\n\n")
}

func writeNormalIndex(out *strings.Builder, schemaName, tableName string, index *metadata.IndexMetadata) {
	out.WriteString("CREATE")
	if index.Unique {
		out.WriteString(" UNIQUE")
	}
	if index.Type != "" {
		fmt.Fprintf(out, " %s", index.Type)
	}
	fmt.Fprintf(out, " INDEX %s ON\n%s.%s (\n", quote(index.Name), quote(schemaName), quote(tableName))
	for i, column := range index.Expressions {
		if i != 0 {
			out.WriteString(",\n")
		}
		fmt.Fprintf(out, "    %s", quote(column))
		if i < len(index.Descending) && index.Descending[i] {
			out.WriteString(" DESC")
		} else {
			out.WriteString(" ASC")
		}
	}
	out.WriteString("\n)")
	if include := indexIncludeColumns(index.Definition); len(include) > 0 {
		out.WriteString(" INCLUDE (")
		for i, column := range include {
			if i != 0 {
				out.WriteString(", ")
			}
			out.WriteString(quote(column))
		}
		out.WriteString(")")
	}
	if where := indexWhereClause(index.Definition); where != "" {
		fmt.Fprintf(out, "\nWHERE %s", where)
	}
	out.WriteString(";\n\n")
}

func writeIndex(out *strings.Builder, schemaName, tableName string, index *metadata.IndexMetadata) {
	switch strings.ToUpper(index.Type) {
	case "XML":
		writeXMLIndex(out, schemaName, tableName, index)
	case "CLUSTERED COLUMNSTORE":
		writeClusteredColumnStoreIndex(out, schemaName, tableName, index)
	case "NONCLUSTERED COLUMNSTORE":
		writeNonClusteredColumnStoreIndex(out, schemaName, tableName, index)
	case "SPATIAL":
		writeSpatialIndex(out, schemaName, tableName, index)
	default:
		writeNormalIndex(out, schemaName, tableName, index)
	}
}

func writeCheck(out *strings.Builder, check *metadata.CheckConstraintMetadata) {
	fmt.Fprintf(out, "    CONSTRAINT %s CHECK %s", quote(check.Name), checkExpression(check.Expression))
}

// checkExpression keeps exactly the parentheses T-SQL requires around a CHECK
// expression, whether the stored text already has them or not.
func checkExpression(expression string) string {
	expression = strings.TrimSpace(expression)
	for strings.HasPrefix(expression, "(") && strings.HasSuffix(expression, ")") && balancedOuterParens(expression) {
		expression = strings.TrimSpace(expression[1 : len(expression)-1])
	}
	if expression == "" {
		return ""
	}
	return "(" + expression + ")"
}

func balancedOuterParens(expression string) bool {
	depth := 0
	inSingle, inDouble, inBracket := false, false, false
	for i := 0; i < len(expression); i++ {
		ch := expression[i]
		if inSingle {
			if ch == '\'' {
				if i+1 < len(expression) && expression[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
			}
			continue
		}
		if inDouble {
			if ch == '"' {
				if i+1 < len(expression) && expression[i+1] == '"' {
					i++
					continue
				}
				inDouble = false
			}
			continue
		}
		if inBracket {
			if ch == ']' {
				if i+1 < len(expression) && expression[i+1] == ']' {
					i++
					continue
				}
				inBracket = false
			}
			continue
		}
		switch ch {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '[':
			inBracket = true
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(expression)-1 {
				return false
			}
		}
	}
	return depth == 0
}

func writeForeignKey(out *strings.Builder, fk *metadata.ForeignKeyMetadata) {
	fmt.Fprintf(out, "    CONSTRAINT %s FOREIGN KEY (", quote(fk.Name))
	for i, column := range fk.Columns {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(quote(column))
	}
	fmt.Fprintf(out, ") REFERENCES %s.%s", quote(fk.ReferencedSchema), quote(fk.ReferencedTable))
	if len(fk.ReferencedColumns) > 0 {
		out.WriteString(" (")
		for i, column := range fk.ReferencedColumns {
			if i != 0 {
				out.WriteString(", ")
			}
			out.WriteString(quote(column))
		}
		out.WriteString(")")
	}
	if fk.OnDelete != "" {
		fmt.Fprintf(out, " ON DELETE %s", fk.OnDelete)
	}
	if fk.OnUpdate != "" {
		fmt.Fprintf(out, " ON UPDATE %s", fk.OnUpdate)
	}
}

func writeKey(out *strings.Builder, key *metadata.IndexMetadata) {
	fmt.Fprintf(out, "    CONSTRAINT %s", quote(key.Name))
	if key.Primary {
		out.WriteString(" PRIMARY KEY")
	} else if key.Unique {
		out.WriteString(" UNIQUE")
	}
	if key.Type != "" {
		fmt.Fprintf(out, " %s", key.Type)
	}
	out.WriteString(" (")
	for i, column := range key.Expressions {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(quote(column))
		if i < len(key.Descending) && key.Descending[i] {
			out.WriteString(" DESC")
		} else {
			out.WriteString(" ASC")
		}
	}
	out.WriteString(")")
}

func writeColumnFiltered(out *strings.Builder, column *metadata.ColumnMetadata, includeDefault bool) {
	fmt.Fprintf(out, "    %s %s", quote(column.Name), column.Type)
	if isComputedColumnType(column.Type) {
		if !column.Nullable && strings.Contains(strings.ToUpper(column.Type), "PERSISTED") {
			out.WriteString(" NOT NULL")
		}
		return
	}
	if column.IsIdentity {
		fmt.Fprintf(out, " IDENTITY(%d,%d)", column.IdentitySeed, column.IdentityIncrement)
	}
	if column.Collation != "" {
		fmt.Fprintf(out, " COLLATE %s", column.Collation)
	}
	if includeDefault && column.GetDefault() != "" {
		if column.DefaultConstraintName != "" {
			fmt.Fprintf(out, " CONSTRAINT %s DEFAULT %s", quote(column.DefaultConstraintName), column.GetDefault())
		} else {
			fmt.Fprintf(out, " DEFAULT %s", column.GetDefault())
		}
	}
	if !column.Nullable {
		out.WriteString(" NOT NULL")
	}
}

// creationDefinition turns the stored module statement into a CREATE, so a
// view, function or procedure last updated with ALTER restores into an empty
// database. Only the leading module header is rewritten, not its body.
func creationDefinition(definition string) string {
	match := creationHeader.FindStringSubmatchIndex(definition)
	if match == nil {
		return definition
	}
	kind := strings.ToUpper(definition[match[2]:match[3]])
	if kind == "PROC" {
		kind = "PROCEDURE"
	}
	return definition[:match[0]] + "CREATE " + kind + definition[match[1]:]
}

var creationHeader = regexp.MustCompile(`(?is)^\s*(?:(?:--[^\n]*\n)|(?:/\*.*?\*/)\s*)*(?:CREATE\s+OR\s+ALTER|ALTER)\s+(VIEW|PROCEDURE|PROC|FUNCTION|TRIGGER)\b`)

// hasRestorableDefinition is false for the placeholder a sync writes where an
// encrypted module's definition is unavailable. Its extended property has no
// object to attach to.
func hasRestorableDefinition(definition string) bool {
	trimmed := strings.TrimSpace(definition)
	return trimmed != "" && !(strings.HasPrefix(trimmed, "/* Definition of ") && strings.HasSuffix(trimmed, " is encrypted. */"))
}

func writeView(out *strings.Builder, schemaName string, view *metadata.ViewMetadata) {
	// A synced definition keeps the trailing semicolon the CREATE VIEW was run
	// with. Trimming and adding one back makes a dump of a dump identical:
	// appending unconditionally would add a semicolon on every restore.
	def := qualifyModuleDefinition(creationDefinition(view.Definition), schemaName, false)
	fmt.Fprintf(out, "%s;\n\nGO\n\n", strings.TrimRight(strings.TrimSpace(def), ";"))
	if schemaName != "" && view.Comment != "" && hasRestorableDefinition(view.Definition) {
		fmt.Fprintf(out, "%s;\nGO\n\n", generateViewCommentSQL("ADD", schemaName, view.Name, view.Comment))
	}
}

func writeFunction(out *strings.Builder, schemaName string, function *metadata.FunctionMetadata) {
	def := qualifyModuleDefinition(creationDefinition(function.Definition), schemaName, false)
	fmt.Fprintf(out, "%s\n\nGO\n\n", def)
	if schemaName != "" && function.Comment != "" && hasRestorableDefinition(function.Definition) {
		fmt.Fprintf(out, "%s;\nGO\n\n", generateFunctionCommentSQL("ADD", schemaName, function.Name, function.Comment))
	}
}

func writeProcedure(out *strings.Builder, schemaName string, procedure *metadata.ProcedureMetadata) {
	def := qualifyModuleDefinition(creationDefinition(procedure.Definition), schemaName, false)
	fmt.Fprintf(out, "%s\n\nGO\n\n", def)
	if schemaName != "" && procedure.Comment != "" && hasRestorableDefinition(procedure.Definition) {
		fmt.Fprintf(out, "%s;\nGO\n\n", generateExtendedPropertySQL("ADD", "PROCEDURE", schemaName, procedure.Name, procedure.Comment))
	}
}

// GetTableDefinition renders one table as the T-SQL that recreates it.
func GetTableDefinition(schemaName string, table *metadata.TableMetadata, _ []*metadata.SequenceMetadata) (string, error) {
	if table == nil || table.GetSkipDump() {
		return "", nil
	}
	if schemaName == "" {
		schemaName = defaultSchema
	}
	var buf strings.Builder
	writeTable(&buf, schemaName, table)
	return buf.String(), nil
}

// GetViewDefinition renders one view's definition without its extended
// properties, which DatabaseDefinition emits after it.
func GetViewDefinition(view *metadata.ViewMetadata) (string, error) {
	var buf strings.Builder
	writeView(&buf, "", view)
	return buf.String(), nil
}

// GetFunctionDefinition renders one function's definition without its extended
// properties, which DatabaseDefinition emits after it.
func GetFunctionDefinition(function *metadata.FunctionMetadata) (string, error) {
	var buf strings.Builder
	writeFunction(&buf, "", function)
	return buf.String(), nil
}

// GetProcedureDefinition renders one procedure's definition without its
// extended properties, which DatabaseDefinition emits after it.
func GetProcedureDefinition(procedure *metadata.ProcedureMetadata) (string, error) {
	var buf strings.Builder
	writeProcedure(&buf, "", procedure)
	return buf.String(), nil
}

// indexWhereClause returns the predicate of a filtered index from either a
// bare predicate or a definition that carries it after WHERE.
func indexWhereClause(definition string) string {
	trimmed := strings.TrimSpace(definition)
	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "WHERE ") {
		return strings.TrimSpace(trimmed[len("WHERE "):])
	}
	if whereIndex := strings.Index(upper, " WHERE "); whereIndex != -1 {
		return strings.TrimSpace(trimmed[whereIndex+len(" WHERE "):])
	}
	return ""
}

func generateSpatialIndexDefinition(index *metadata.IndexMetadata, schemaName, tableName string) string {
	var buf strings.Builder
	buf.WriteString("CREATE SPATIAL INDEX ")
	buf.WriteString(quote(index.Name))
	buf.WriteString(" ON ")
	buf.WriteString(quote(schemaName))
	buf.WriteString(".")
	buf.WriteString(quote(tableName))
	buf.WriteString(" (")
	for i, expr := range index.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(quote(expr))
	}
	buf.WriteString(")")
	if index.SpatialConfig == nil || index.SpatialConfig.Tessellation == nil {
		return buf.String()
	}
	tessellation := index.SpatialConfig.Tessellation
	if tessellation.Scheme != "" {
		buf.WriteString("\nUSING ")
		buf.WriteString(tessellation.Scheme)
	}
	var withParams []string
	withParams = append(withParams, buildTessellationParams(tessellation)...)
	if index.SpatialConfig.Storage != nil {
		withParams = append(withParams, buildStorageParams(index.SpatialConfig.Storage)...)
	}
	if len(withParams) > 0 {
		buf.WriteString("\nWITH (\n    ")
		buf.WriteString(strings.Join(withParams, ",\n    "))
		buf.WriteString("\n)")
	}
	return buf.String()
}

func buildTessellationParams(tessellation *metadata.TessellationConfig) []string {
	var params []string
	if tessellation.BoundingBox != nil {
		bbox := tessellation.BoundingBox
		params = append(params, fmt.Sprintf("BOUNDING_BOX = (%g, %g, %g, %g)", bbox.Xmin, bbox.Ymin, bbox.Xmax, bbox.Ymax))
	}
	if len(tessellation.GridLevels) > 0 {
		var gridParts []string
		for _, level := range tessellation.GridLevels {
			gridParts = append(gridParts, fmt.Sprintf("LEVEL_%d = %s", level.Level, level.Density))
		}
		if len(gridParts) > 0 {
			params = append(params, fmt.Sprintf("GRIDS = (%s)", strings.Join(gridParts, ", ")))
		}
	}
	if tessellation.CellsPerObject > 0 {
		params = append(params, fmt.Sprintf("CELLS_PER_OBJECT = %d", tessellation.CellsPerObject))
	}
	return params
}

func buildStorageParams(storage *metadata.StorageConfig) []string {
	var params []string
	if storage.PadIndex {
		params = append(params, "PAD_INDEX = ON")
	}
	if storage.Fillfactor > 0 {
		params = append(params, fmt.Sprintf("FILLFACTOR = %d", storage.Fillfactor))
	}
	if storage.SortInTempdb != "" {
		params = append(params, fmt.Sprintf("SORT_IN_TEMPDB = %s", storage.SortInTempdb))
	}
	if storage.DropExisting {
		params = append(params, "DROP_EXISTING = ON")
	}
	if storage.Online {
		params = append(params, "ONLINE = ON")
	}
	if !storage.AllowRowLocks {
		params = append(params, "ALLOW_ROW_LOCKS = OFF")
	}
	if !storage.AllowPageLocks {
		params = append(params, "ALLOW_PAGE_LOCKS = OFF")
	}
	if storage.Maxdop > 0 {
		params = append(params, fmt.Sprintf("MAXDOP = %d", storage.Maxdop))
	}
	if storage.DataCompression != "" && storage.DataCompression != "NONE" {
		params = append(params, fmt.Sprintf("DATA_COMPRESSION = %s", storage.DataCompression))
	}
	return params
}
