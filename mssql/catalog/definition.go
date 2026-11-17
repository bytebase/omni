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
	if len(tableFunctions) > 0 {
		if hasPreamble {
			buf.WriteString("GO\n\n")
		}
		writeReferencedFunctions(&buf, to.Schemas, tableFunctions)
	}

	if hasTables(to.Schemas) {
		writeAllTablesInOrder(&buf, to.Schemas)
	}

	viewFunctions := functionsReferencedByViews(to.Schemas)
	for name := range tableFunctions {
		delete(viewFunctions, name)
	}
	if len(viewFunctions) > 0 {
		if hasTables(to.Schemas) || hasPreamble || len(tableFunctions) > 0 {
			buf.WriteString("GO\n\n")
		}
		writeReferencedFunctions(&buf, to.Schemas, viewFunctions)
	}

	if hasViews(to.Schemas) {
		if len(viewFunctions) == 0 && (hasTables(to.Schemas) || hasPreamble) {
			buf.WriteString("GO\n\n")
		}
		writeAllViewsInOrder(&buf, to.Schemas)
	}

	if !hasViews(to.Schemas) && (hasTables(to.Schemas) || hasPreamble) && (hasRoutines(to.Schemas) || hasTriggers(to.Schemas)) {
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

func writeReferencedFunctions(out *strings.Builder, schemas []*metadata.SchemaMetadata, names map[string]bool) {
	for _, schema := range schemas {
		for _, function := range sortFunctions(schema.Name, schema.Functions) {
			if function.GetSkipDump() {
				continue
			}
			if names[objectID(schema.Name, function.Name)] {
				writeFunction(out, function)
			}
		}
	}
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
			if strings.Contains(strings.ToLower(function.Definition), strings.ToLower(other.Name)) {
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

// functionsReferencedByTables names the functions a column default or computed
// expression mentions, so a table can install them before it needs them.
func functionsReferencedByTables(schemas []*metadata.SchemaMetadata) map[string]bool {
	var expressions []string
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			for _, column := range table.Columns {
				expressions = append(expressions, column.Default, column.Type)
			}
			for _, check := range table.CheckConstraints {
				expressions = append(expressions, check.Expression)
			}
		}
	}
	referenced := make(map[string]bool)
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			for _, expression := range expressions {
				if strings.Contains(strings.ToLower(expression), strings.ToLower(function.Name)) {
					referenced[objectID(schema.Name, function.Name)] = true
					break
				}
			}
		}
	}
	return referenced
}

// functionsReferencedByViews names the functions a view body mentions, so they
// are created before the view that resolves them.
func functionsReferencedByViews(schemas []*metadata.SchemaMetadata) map[string]bool {
	var definitions []string
	for _, schema := range schemas {
		for _, view := range schema.Views {
			definitions = append(definitions, view.Definition)
		}
	}
	referenced := make(map[string]bool)
	for _, schema := range schemas {
		for _, function := range schema.Functions {
			for _, definition := range definitions {
				if strings.Contains(strings.ToLower(definition), strings.ToLower(function.Name)) {
					referenced[objectID(schema.Name, function.Name)] = true
					break
				}
			}
		}
	}
	return referenced
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
				fmt.Fprintf(out, "%s\n\nGO\n\n", strings.TrimRight(creationDefinition(trigger.Body), ";"))
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
				fmt.Fprintf(out, "%s\n\nGO\n\n", strings.TrimRight(creationDefinition(trigger.Body), ";"))
			}
		}
	}
}

func writeFunctionsAndProcedures(out *strings.Builder, schema *metadata.SchemaMetadata, early map[string]bool) {
	for _, function := range sortFunctions(schema.Name, schema.Functions) {
		if function.GetSkipDump() || early[objectID(schema.Name, function.Name)] {
			continue
		}
		writeFunction(out, function)
	}
	for _, procedure := range schema.Procedures {
		if procedure.SkipDump {
			continue
		}
		writeProcedure(out, procedure)
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
		writeView(out, vws.view)
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
	fmt.Fprintf(out, "CREATE TABLE %s.%s (\n", quote(schemaName), quote(table.Name))
	for i, column := range table.Columns {
		if i != 0 {
			out.WriteString(",\n")
		}
		writeColumn(out, column)
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
	fmt.Fprintf(out, "CREATE PRIMARY XML INDEX %s ON %s.%s (%s);\n\n",
		quote(index.Name), quote(schemaName), quote(tableName), column)
}

func orderedIndexes(indexes []*metadata.IndexMetadata) []*metadata.IndexMetadata {
	out := make([]*metadata.IndexMetadata, 0, len(indexes))
	for _, index := range indexes {
		if index.Type == "XML" && index.Definition == "PRIMARY" {
			out = append(out, index)
		}
	}
	for _, index := range indexes {
		if index.Type != "XML" || index.Definition != "PRIMARY" {
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

func writeColumn(out *strings.Builder, column *metadata.ColumnMetadata) {
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
	if column.GetDefault() != "" {
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

func writeView(out *strings.Builder, view *metadata.ViewMetadata) {
	// A synced definition keeps the trailing semicolon the CREATE VIEW was run
	// with. Trimming and adding one back makes a dump of a dump identical:
	// appending unconditionally would add a semicolon on every restore.
	fmt.Fprintf(out, "%s;\n\nGO\n\n", strings.TrimRight(strings.TrimSpace(creationDefinition(view.Definition)), ";"))
}

func writeFunction(out *strings.Builder, function *metadata.FunctionMetadata) {
	fmt.Fprintf(out, "%s\n\nGO\n\n", creationDefinition(function.Definition))
}

func writeProcedure(out *strings.Builder, procedure *metadata.ProcedureMetadata) {
	fmt.Fprintf(out, "%s\n\nGO\n\n", creationDefinition(procedure.Definition))
}

// GetViewDefinition renders one view the way DatabaseDefinition does.
func GetViewDefinition(view *metadata.ViewMetadata) (string, error) {
	var buf strings.Builder
	writeView(&buf, view)
	return buf.String(), nil
}

// GetFunctionDefinition renders one function the way DatabaseDefinition does.
func GetFunctionDefinition(function *metadata.FunctionMetadata) (string, error) {
	var buf strings.Builder
	writeFunction(&buf, function)
	return buf.String(), nil
}

// GetProcedureDefinition renders one procedure the way DatabaseDefinition does.
func GetProcedureDefinition(procedure *metadata.ProcedureMetadata) (string, error) {
	var buf strings.Builder
	writeProcedure(&buf, procedure)
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
