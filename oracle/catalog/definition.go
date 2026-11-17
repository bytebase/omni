package catalog

import (
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

func DatabaseDefinition(to *metadatapb.DatabaseSchemaMetadata) (string, error) {
	if len(to.Schemas) == 0 {
		return "", nil
	}

	var buf strings.Builder
	schema := to.Schemas[0]
	// For Oracle, the database name (user/schema name) is stored in to.Name;
	// a snapshot loaded from SQL may carry it only on the schema.
	databaseName := to.Name
	if databaseName == "" {
		databaseName = schema.Name
	}

	// First generate sequences (they need to exist before tables that reference them)
	for _, sequence := range schema.Sequences {
		// Skip system-generated sequences
		if !strings.HasPrefix(sequence.Name, "ISEQ$$_") {
			if err := writeSequence(&buf, sequence); err != nil {
				return "", err
			}
		}
	}

	// Build dependency graph for topological sorting
	graph := newGraph()
	tableMap := make(map[string]*metadatapb.TableMetadata)
	viewMap := make(map[string]*metadatapb.ViewMetadata)
	materializedViewMap := make(map[string]*metadatapb.MaterializedViewMetadata)
	functionMap := make(map[string]*metadatapb.FunctionMetadata)
	procedureMap := make(map[string]*metadatapb.ProcedureMetadata)

	// SkipDump marks object and XMLType tables: CREATE TABLE ... OF needs the
	// object type, which the metadata does not carry.
	var tables []*metadatapb.TableMetadata
	for _, table := range schema.Tables {
		if !table.SkipDump {
			tables = append(tables, table)
		}
	}

	// Add tables to graph
	for _, table := range tables {
		tableID := getObjectID(databaseName, table.Name)
		tableMap[tableID] = table
		graph.AddNode(tableID)
	}

	// Add views to graph
	for _, view := range schema.Views {
		viewID := getObjectID(databaseName, view.Name)
		viewMap[viewID] = view
		graph.AddNode(viewID)
		// Add dependencies from view to tables/views it references
		for _, dependency := range view.DependencyColumns {
			dependencyID := getObjectID(dependency.Schema, dependency.Table)
			graph.AddEdge(dependencyID, viewID)
		}
	}

	// Add materialized views to graph
	for _, view := range schema.MaterializedViews {
		viewID := getObjectID(databaseName, view.Name)
		materializedViewMap[viewID] = view
		graph.AddNode(viewID)
		// Add dependencies from materialized view to tables/views it references
		for _, dependency := range view.DependencyColumns {
			dependencyID := getObjectID(dependency.Schema, dependency.Table)
			graph.AddEdge(dependencyID, viewID)
		}
	}

	// Add functions to graph
	for _, function := range schema.Functions {
		functionID := getObjectID(databaseName, function.Name)
		functionMap[functionID] = function
		graph.AddNode(functionID)
		// Add dependencies from function to tables it references
		for _, dependency := range function.DependencyTables {
			dependencyID := getObjectID(dependency.Schema, dependency.Table)
			graph.AddEdge(dependencyID, functionID)
		}
	}

	// Add procedures to graph
	for _, procedure := range schema.Procedures {
		procedureID := getObjectID(databaseName, procedure.Name)
		procedureMap[procedureID] = procedure
		graph.AddNode(procedureID)
		// Note: ProcedureMetadata doesn't have DependencyTables field
		// Procedures will be created after tables by default order
	}

	// Views and virtual columns can call functions, so the function must exist
	// first.
	for viewID, view := range viewMap {
		for functionID, function := range functionMap {
			if expressionReferencesFunction(view.Definition, function.Name) {
				graph.AddEdge(functionID, viewID)
			}
		}
	}
	for viewID, view := range materializedViewMap {
		for functionID, function := range functionMap {
			if expressionReferencesFunction(view.Definition, function.Name) {
				graph.AddEdge(functionID, viewID)
			}
		}
	}
	for _, table := range tables {
		tableID := getObjectID(databaseName, table.Name)
		for _, column := range table.Columns {
			if column.GetGeneration() == nil {
				continue
			}
			for functionID, function := range functionMap {
				if expressionReferencesFunction(column.GetGeneration().GetExpression(), function.Name) {
					graph.AddEdge(functionID, tableID)
				}
			}
		}
	}

	// Add foreign key dependencies between tables
	for _, table := range tables {
		tableID := getObjectID(databaseName, table.Name)
		for _, fk := range table.ForeignKeys {
			if fk.ReferencedTable != table.Name { // Avoid self-references
				referencedTableID := getObjectID(databaseName, fk.ReferencedTable)
				graph.AddEdge(referencedTableID, tableID)
			}
		}
	}

	// Perform topological sort
	orderedList, err := graph.TopologicalSort()
	if err != nil {
		// If there are cycles, fall back to original order. Functions a virtual
		// column calls must exist before their table; the foreign keys are added
		// once all tables exist.
		prewrittenFunctions := map[string]bool{}
		for _, table := range tables {
			for _, column := range table.Columns {
				if column.GetGeneration() == nil {
					continue
				}
				for _, function := range schema.Functions {
					if !prewrittenFunctions[function.Name] && expressionReferencesFunction(column.GetGeneration().GetExpression(), function.Name) {
						if err := writeFunction(&buf, function); err != nil {
							return "", err
						}
						prewrittenFunctions[function.Name] = true
					}
				}
			}
		}
		for _, table := range tables {
			if err := writeTableWithOptions(&buf, table, false, false); err != nil {
				return "", err
			}
		}
		for _, table := range tables {
			for _, fk := range table.ForeignKeys {
				if err := writeAddForeignKey(&buf, table.Name, fk); err != nil {
					return "", err
				}
			}
		}
		if err := writeViewsInDependencyOrder(&buf, schema); err != nil {
			return "", err
		}
		for _, function := range schema.Functions {
			if prewrittenFunctions[function.Name] {
				continue
			}
			if err := writeFunction(&buf, function); err != nil {
				return "", err
			}
		}
		for _, procedure := range schema.Procedures {
			if err := writeProcedure(&buf, procedure); err != nil {
				return "", err
			}
		}
		if err := writeDeferredFunctionIndexes(&buf, tables); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// Generate objects in dependency order
	for _, objectID := range orderedList {
		if table, ok := tableMap[objectID]; ok {
			if err := writeTableWithoutFunctionIndexes(&buf, table); err != nil {
				return "", err
			}
		}
		if view, ok := viewMap[objectID]; ok {
			if err := writeView(&buf, view); err != nil {
				return "", err
			}
		}
		if view, ok := materializedViewMap[objectID]; ok {
			if err := writeMaterializedView(&buf, view); err != nil {
				return "", err
			}
		}
		if function, ok := functionMap[objectID]; ok {
			if err := writeFunction(&buf, function); err != nil {
				return "", err
			}
		}
		if procedure, ok := procedureMap[objectID]; ok {
			if err := writeProcedure(&buf, procedure); err != nil {
				return "", err
			}
		}
	}
	if err := writeDeferredFunctionIndexes(&buf, tables); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func writeViewsInDependencyOrder(buf *strings.Builder, schema *metadatapb.SchemaMetadata) error {
	graph := newGraph()
	views := map[string]*metadatapb.ViewMetadata{}
	matViews := map[string]*metadatapb.MaterializedViewMetadata{}
	for _, view := range schema.Views {
		key := "v:" + view.Name
		views[key] = view
		graph.AddNode(key)
	}
	for _, view := range schema.MaterializedViews {
		key := "m:" + view.Name
		matViews[key] = view
		graph.AddNode(key)
	}
	addEdges := func(key string, dependencies []*metadatapb.DependencyColumn) {
		for _, dependency := range dependencies {
			for _, candidate := range []string{"v:" + dependency.Table, "m:" + dependency.Table} {
				if views[candidate] != nil || matViews[candidate] != nil {
					graph.AddEdge(candidate, key)
				}
			}
		}
	}
	for key, view := range views {
		addEdges(key, view.DependencyColumns)
	}
	for key, view := range matViews {
		addEdges(key, view.DependencyColumns)
	}
	ordered, err := graph.TopologicalSort()
	if err != nil {
		ordered = nil
		for key := range views {
			ordered = append(ordered, key)
		}
		for key := range matViews {
			ordered = append(ordered, key)
		}
		slices.Sort(ordered)
	}
	for _, key := range ordered {
		if view := views[key]; view != nil {
			if err := writeView(buf, view); err != nil {
				return err
			}
		}
		if view := matViews[key]; view != nil {
			if err := writeMaterializedView(buf, view); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeDeferredFunctionIndexes(buf *strings.Builder, tables []*metadatapb.TableMetadata) error {
	for _, table := range tables {
		for _, index := range table.Indexes {
			if index.IsConstraint || !strings.Contains(index.Type, "FUNCTION-BASED") {
				continue
			}
			if err := writeIndex(buf, table.Name, index); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetTableDefinition(_ string, table *metadatapb.TableMetadata, _ []*metadatapb.SequenceMetadata) (string, error) {
	if table.SkipDump {
		return "", nil
	}
	var buf strings.Builder
	if err := writeTable(&buf, table); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// The single-object definitions below all ignore the schema argument: Oracle has
// one schema per database and the writers emit unqualified names, matching
// GetTableDefinition and the whole-database output.

func GetViewDefinition(_ string, view *metadatapb.ViewMetadata) (string, error) {
	var buf strings.Builder
	if err := writeView(&buf, view); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetMaterializedViewDefinition(_ string, view *metadatapb.MaterializedViewMetadata) (string, error) {
	var buf strings.Builder
	if err := writeMaterializedView(&buf, view); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetFunctionDefinition(_ string, function *metadatapb.FunctionMetadata) (string, error) {
	var buf strings.Builder
	if err := writeFunction(&buf, function); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetProcedureDefinition(_ string, procedure *metadatapb.ProcedureMetadata) (string, error) {
	var buf strings.Builder
	if err := writeProcedure(&buf, procedure); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func writeTable(buf *strings.Builder, table *metadatapb.TableMetadata) error {
	return writeTableWithOptions(buf, table, true, true)
}

func writeTableWithoutFunctionIndexes(buf *strings.Builder, table *metadatapb.TableMetadata) error {
	return writeTableWithOptions(buf, table, true, false)
}

func writeTableWithForeignKeys(buf *strings.Builder, table *metadatapb.TableMetadata, includeForeignKeys bool) error {
	return writeTableWithOptions(buf, table, includeForeignKeys, true)
}

func writeTableWithOptions(buf *strings.Builder, table *metadatapb.TableMetadata, includeForeignKeys, includeFunctionIndexes bool) error {
	if _, err := buf.WriteString(`CREATE TABLE "`); err != nil {
		return err
	}
	if _, err := buf.WriteString(table.Name); err != nil {
		return err
	}
	if _, err := buf.WriteString("\" (\n"); err != nil {
		return err
	}
	itemCount := 0
	writeSeparator := func() error {
		if itemCount > 0 {
			if _, err := buf.WriteString(",\n"); err != nil {
				return err
			}
		}
		itemCount++
		_, err := buf.WriteString(`  `)
		return err
	}

	for _, column := range table.Columns {
		if err := writeSeparator(); err != nil {
			return err
		}
		if err := writeColumn(buf, column); err != nil {
			return err
		}
	}

	for _, constraint := range table.Indexes {
		if !constraint.IsConstraint {
			continue
		}
		if err := writeSeparator(); err != nil {
			return err
		}
		if err := writeConstraint(buf, constraint); err != nil {
			return err
		}
	}

	for _, check := range table.CheckConstraints {
		if err := writeSeparator(); err != nil {
			return err
		}
		if err := writeCheckConstraint(buf, check); err != nil {
			return err
		}
	}

	if includeForeignKeys {
		for _, fk := range table.ForeignKeys {
			if err := writeSeparator(); err != nil {
				return err
			}
			if err := writeForeignKey(buf, fk); err != nil {
				return err
			}
		}
	}
	if _, err := buf.WriteString("\n);\n\n"); err != nil {
		return err
	}

	for _, index := range table.Indexes {
		if index.IsConstraint {
			continue
		}
		if !includeFunctionIndexes && strings.Contains(index.Type, "FUNCTION-BASED") {
			continue
		}
		if err := writeIndex(buf, table.Name, index); err != nil {
			return err
		}
	}

	if err := writeTriggers(buf, table.Triggers); err != nil {
		return err
	}
	if table.Comment != "" {
		writeTableComment(buf, "", table.Name, table.Comment)
	}
	for _, column := range table.Columns {
		if column.Comment != "" {
			writeColumnComment(buf, "", table.Name, column.Name, column.Comment)
		}
	}
	return nil
}

func writeIndex(buf *strings.Builder, table string, index *metadatapb.IndexMetadata) error {
	if _, err := buf.WriteString(`CREATE`); err != nil {
		return err
	}

	switch index.Type {
	case "BITMAP", "FUNCTION-BASED BITMAP":
		if _, err := buf.WriteString(` BITMAP`); err != nil {
			return err
		}
	default:
		// Other index types don't need special syntax
	}

	if index.Unique {
		if _, err := buf.WriteString(` UNIQUE`); err != nil {
			return err
		}
	}

	if _, err := buf.WriteString(` INDEX "`); err != nil {
		return err
	}

	if _, err := buf.WriteString(index.Name); err != nil {
		return err
	}

	if _, err := buf.WriteString(`" ON "`); err != nil {
		return err
	}

	if _, err := buf.WriteString(table); err != nil {
		return err
	}

	if _, err := buf.WriteString(`" (`); err != nil {
		return err
	}

	if strings.Contains(index.Type, "FUNCTION-BASED") {
		for i, expression := range index.Expressions {
			if i != 0 {
				if _, err := buf.WriteString(", "); err != nil {
					return err
				}
			}
			if _, err := buf.WriteString(expression); err != nil {
				return err
			}

			if i < len(index.Descending) && index.Descending[i] {
				if _, err := buf.WriteString(` DESC`); err != nil {
					return err
				}
			} else {
				if _, err := buf.WriteString(` ASC`); err != nil {
					return err
				}
			}
		}
	} else {
		for i, column := range index.Expressions {
			if i != 0 {
				if _, err := buf.WriteString(", "); err != nil {
					return err
				}
			}

			// Check if column is already quoted to avoid double quoting
			if strings.HasPrefix(column, `"`) && strings.HasSuffix(column, `"`) {
				// Column is already quoted, use as-is
				if _, err := buf.WriteString(column); err != nil {
					return err
				}
			} else {
				// Column is not quoted, add quotes
				if _, err := buf.WriteString(`"`); err != nil {
					return err
				}
				if _, err := buf.WriteString(column); err != nil {
					return err
				}
				if _, err := buf.WriteString(`"`); err != nil {
					return err
				}
			}

			if i < len(index.Descending) && index.Descending[i] {
				if _, err := buf.WriteString(` DESC`); err != nil {
					return err
				}
			} else {
				if _, err := buf.WriteString(` ASC`); err != nil {
					return err
				}
			}
		}
	}

	if _, err := buf.WriteString(`)`); err != nil {
		return err
	}
	if !index.Visible {
		if _, err := buf.WriteString(` INVISIBLE`); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(`;`); err != nil {
		return err
	}
	if _, err := buf.WriteString("\n\n"); err != nil {
		return err
	}
	return nil
}

func writeForeignKey(buf *strings.Builder, fk *metadatapb.ForeignKeyMetadata) error {
	if _, err := buf.WriteString(`CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := buf.WriteString(fk.Name); err != nil {
		return err
	}
	if _, err := buf.WriteString(`"`); err != nil {
		return err
	}
	if _, err := buf.WriteString(` FOREIGN KEY (`); err != nil {
		return err
	}
	for i, column := range fk.Columns {
		if i != 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(`"`); err != nil {
			return err
		}
		if _, err := buf.WriteString(column); err != nil {
			return err
		}
		if _, err := buf.WriteString(`"`); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(`) REFERENCES "`); err != nil {
		return err
	}
	// Only include schema prefix for cross-schema foreign key references
	// Same-schema references don't need prefix for portability
	if fk.ReferencedSchema != "" {
		if _, err := buf.WriteString(fk.ReferencedSchema); err != nil {
			return err
		}
		if _, err := buf.WriteString(`"."`); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(fk.ReferencedTable); err != nil {
		return err
	}
	if _, err := buf.WriteString(`"`); err != nil {
		return err
	}
	if len(fk.ReferencedColumns) > 0 {
		if _, err := buf.WriteString(` (`); err != nil {
			return err
		}
	}
	for i, column := range fk.ReferencedColumns {
		if i != 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(`"`); err != nil {
			return err
		}
		if _, err := buf.WriteString(column); err != nil {
			return err
		}
		if _, err := buf.WriteString(`"`); err != nil {
			return err
		}
	}
	if len(fk.ReferencedColumns) > 0 {
		if _, err := buf.WriteString(`)`); err != nil {
			return err
		}
	}
	if fk.OnDelete != "" {
		if _, err := buf.WriteString(` ON DELETE `); err != nil {
			return err
		}
		if _, err := buf.WriteString(fk.OnDelete); err != nil {
			return err
		}
	}
	return nil
}

func writeAddForeignKey(buf *strings.Builder, table string, fk *metadatapb.ForeignKeyMetadata) error {
	if _, err := buf.WriteString(`ALTER TABLE "`); err != nil {
		return err
	}
	if _, err := buf.WriteString(table); err != nil {
		return err
	}
	if _, err := buf.WriteString(`" ADD `); err != nil {
		return err
	}
	if err := writeForeignKey(buf, fk); err != nil {
		return err
	}
	if _, err := buf.WriteString(";\n\n"); err != nil {
		return err
	}
	return nil
}

func writeCheckConstraint(buf *strings.Builder, check *metadatapb.CheckConstraintMetadata) error {
	if _, err := buf.WriteString(`CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := buf.WriteString(check.Name); err != nil {
		return err
	}
	if _, err := buf.WriteString(`"`); err != nil {
		return err
	}
	if _, err := buf.WriteString(` CHECK (`); err != nil {
		return err
	}
	if _, err := buf.WriteString(check.Expression); err != nil {
		return err
	}
	_, err := buf.WriteString(`)`)
	return err
}

func writeConstraint(buf *strings.Builder, constraint *metadatapb.IndexMetadata) error {
	if _, err := buf.WriteString(`CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := buf.WriteString(constraint.Name); err != nil {
		return err
	}
	if _, err := buf.WriteString(`"`); err != nil {
		return err
	}

	switch {
	case constraint.Primary:
		if _, err := buf.WriteString(` PRIMARY KEY (`); err != nil {
			return err
		}
		for i, column := range constraint.Expressions {
			if i != 0 {
				if _, err := buf.WriteString(", "); err != nil {
					return err
				}
			}
			if _, err := buf.WriteString(`"`); err != nil {
				return err
			}
			if _, err := buf.WriteString(column); err != nil {
				return err
			}
			if _, err := buf.WriteString(`"`); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(`)`); err != nil {
			return err
		}
	case !constraint.Primary && constraint.Unique:
		if _, err := buf.WriteString(` UNIQUE (`); err != nil {
			return err
		}
		for i, column := range constraint.Expressions {
			if i != 0 {
				if _, err := buf.WriteString(", "); err != nil {
					return err
				}
			}
			if _, err := buf.WriteString(`"`); err != nil {
				return err
			}
			if _, err := buf.WriteString(column); err != nil {
				return err
			}
			if _, err := buf.WriteString(`"`); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(`)`); err != nil {
			return err
		}
	default:
		// Other constraint types (e.g. CHECK, FOREIGN KEY)
	}
	return nil
}

func writeColumn(buf *strings.Builder, column *metadatapb.ColumnMetadata) error {
	if _, err := buf.WriteString(`"`); err != nil {
		return err
	}
	if _, err := buf.WriteString(column.Name); err != nil {
		return err
	}
	if _, err := buf.WriteString(`" `); err != nil {
		return err
	}
	if _, err := buf.WriteString(column.Type); err != nil {
		return err
	}
	if column.IsInvisible {
		if _, err := buf.WriteString(` INVISIBLE`); err != nil {
			return err
		}
	}
	if column.IsIdentity {
		if _, err := buf.WriteString(columnIdentityClause(column)); err != nil {
			return err
		}
		if !column.Nullable {
			if _, err := buf.WriteString(` NOT NULL`); err != nil {
				return err
			}
		}
		return nil
	}
	if column.Generation != nil {
		if _, err := buf.WriteString(virtualColumnClause(column.Generation)); err != nil {
			return err
		}
		if !column.Nullable {
			if _, err := buf.WriteString(` NOT NULL`); err != nil {
				return err
			}
		}
		return nil
	}
	// Skip collation for Oracle as it causes issues with standard string size settings
	// if column.Collation != "" {
	//	if _, err := buf.WriteString(` COLLATE "`); err != nil {
	//		return err
	//	}
	//	if _, err := buf.WriteString(column.Collation); err != nil {
	//		return err
	//	}
	//	if _, err := buf.WriteString(`"`); err != nil {
	//		return err
	//	}
	// }
	// Handle default values
	hasDefault := column.Default != ""
	if hasDefault {
		defaultExpr := column.Default

		// Skip system-generated sequence references as they can't be manually
		// created. A quoted literal that merely contains the prefix is a value.
		systemSequence := strings.Contains(defaultExpr, "ISEQ$$_") && !strings.HasPrefix(strings.TrimSpace(defaultExpr), "'")
		if defaultExpr != "" && !systemSequence {
			if _, err := buf.WriteString(` DEFAULT `); err != nil {
				return err
			}
			if column.DefaultOnNull {
				if _, err := buf.WriteString(`ON NULL `); err != nil {
					return err
				}
			}
			if _, err := buf.WriteString(defaultExpr); err != nil {
				return err
			}
		}
	}
	if !column.Nullable {
		if _, err := buf.WriteString(` NOT NULL`); err != nil {
			return err
		}
	}
	return nil
}

func writeSequence(buf *strings.Builder, sequence *metadatapb.SequenceMetadata) error {
	if _, err := buf.WriteString("CREATE SEQUENCE \""); err != nil {
		return err
	}
	if _, err := buf.WriteString(sequence.Name); err != nil {
		return err
	}
	if _, err := buf.WriteString("\""); err != nil {
		return err
	}

	if sequence.Start != "" {
		if _, err := buf.WriteString(" START WITH "); err != nil {
			return err
		}
		if _, err := buf.WriteString(sequence.Start); err != nil {
			return err
		}
	}

	if sequence.Increment != "" && sequence.Increment != "0" && sequence.Increment != "1" {
		if _, err := buf.WriteString(" INCREMENT BY "); err != nil {
			return err
		}
		if _, err := buf.WriteString(sequence.Increment); err != nil {
			return err
		}
	}

	if sequence.MaxValue != "" {
		if _, err := buf.WriteString(" MAXVALUE "); err != nil {
			return err
		}
		if _, err := buf.WriteString(sequence.MaxValue); err != nil {
			return err
		}
	}

	if sequence.MinValue != "" {
		if _, err := buf.WriteString(" MINVALUE "); err != nil {
			return err
		}
		if _, err := buf.WriteString(sequence.MinValue); err != nil {
			return err
		}
	}

	if sequence.CacheSize == "0" {
		if _, err := buf.WriteString(" NOCACHE"); err != nil {
			return err
		}
	} else if sequence.CacheSize != "" {
		if _, err := buf.WriteString(" CACHE "); err != nil {
			return err
		}
		if _, err := buf.WriteString(sequence.CacheSize); err != nil {
			return err
		}
	}

	if sequence.Cycle {
		if _, err := buf.WriteString(" CYCLE"); err != nil {
			return err
		}
	}

	if _, err := buf.WriteString(";\n\n"); err != nil {
		return err
	}

	return nil
}

// getObjectID returns a unique identifier for database objects
func getObjectID(schema, objectName string) string {
	if schema == "" {
		return objectName
	}
	return schema + "." + objectName
}

// writeView writes a CREATE VIEW statement
func writeView(buf *strings.Builder, view *metadatapb.ViewMetadata) error {
	if _, err := buf.WriteString(`CREATE VIEW "`); err != nil {
		return err
	}
	if _, err := buf.WriteString(view.Name); err != nil {
		return err
	}
	// The synced column names are deliberately not written back as an explicit
	// list. They would restore the aliases of CREATE VIEW V (RENAMED) AS SELECT
	// BASE ..., which ALL_VIEWS.TEXT does not carry, but getTableColumns drops
	// invisible columns (COLUMN_ID IS NULL, sync.go:485) while the query keeps
	// their expressions -- so a view with an invisible column would get fewer
	// names than output columns and fail with ORA-01730. Losing an alias beats
	// emitting DDL that cannot run.
	if _, err := buf.WriteString(`" AS `); err != nil {
		return err
	}
	if _, err := buf.WriteString(view.Definition); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.TrimSpace(view.Definition), ";") {
		if _, err := buf.WriteString(`;`); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("\n\n"); err != nil {
		return err
	}

	// An INSTEAD OF trigger carries the view's DML behavior, so it belongs with
	// the view the way writeTable emits a table's triggers.
	return writeTriggers(buf, view.Triggers)
}

// virtualColumnClause writes the clause that follows a virtual column's type
// and visibility.
func virtualColumnClause(generation *metadatapb.GenerationMetadata) string {
	return " GENERATED ALWAYS AS (" + generation.Expression + ") VIRTUAL"
}

// writeMaterializedView writes a CREATE MATERIALIZED VIEW statement
func materializedViewDefinitionHasClause(definition string) bool {
	upper := strings.ToUpper(strings.TrimSpace(definition))
	return strings.HasPrefix(upper, "BUILD ") || strings.HasPrefix(upper, "REFRESH ")
}

func writeMaterializedView(buf *strings.Builder, view *metadatapb.MaterializedViewMetadata) error {
	if _, err := buf.WriteString(`CREATE MATERIALIZED VIEW "`); err != nil {
		return err
	}
	if _, err := buf.WriteString(view.Name); err != nil {
		return err
	}
	if _, err := buf.WriteString(`"`); err != nil {
		return err
	}
	if materializedViewDefinitionHasClause(view.Definition) {
		if _, err := buf.WriteString(` `); err != nil {
			return err
		}
	} else if _, err := buf.WriteString(` AS `); err != nil {
		return err
	}
	if _, err := buf.WriteString(view.Definition); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.TrimSpace(view.Definition), ";") {
		if _, err := buf.WriteString(`;`); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("\n\n"); err != nil {
		return err
	}
	if err := writeTriggers(buf, view.Triggers); err != nil {
		return err
	}
	for _, index := range materializedViewIndexes(view.Indexes) {
		if err := writeIndex(buf, view.Name, index); err != nil {
			return err
		}
	}
	return nil
}

// materializedViewIndexes drops the indexes Oracle creates on a materialized
// view's container. Their name is SYS_C_SNAP$_<object id>, and the id differs
// between snapshots, so neither a definition nor a diff may carry them: the
// view makes Oracle create one again.
func materializedViewIndexes(indexes []*metadatapb.IndexMetadata) []*metadatapb.IndexMetadata {
	kept := make([]*metadatapb.IndexMetadata, 0, len(indexes))
	for _, index := range indexes {
		if strings.HasPrefix(index.GetName(), "SYS_C_SNAP$") {
			continue
		}
		kept = append(kept, index)
	}
	return kept
}

// writeFunction writes a CREATE FUNCTION statement
func writeFunction(buf *strings.Builder, function *metadatapb.FunctionMetadata) error {
	definition := function.Definition
	// If the definition doesn't start with CREATE, add the CREATE OR REPLACE prefix
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(definition)), "CREATE") {
		if _, err := buf.WriteString("CREATE OR REPLACE "); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(definition); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.TrimSpace(definition), ";") {
		if _, err := buf.WriteString(";"); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("\n\n"); err != nil {
		return err
	}
	return nil
}

// writeProcedure writes a CREATE PROCEDURE statement
func writeProcedure(buf *strings.Builder, procedure *metadatapb.ProcedureMetadata) error {
	definition := procedure.Definition
	// If the definition doesn't start with CREATE, add the CREATE OR REPLACE prefix
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(definition)), "CREATE") {
		if _, err := buf.WriteString("CREATE OR REPLACE "); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(definition); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.TrimSpace(definition), ";") {
		if _, err := buf.WriteString(";"); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("\n\n"); err != nil {
		return err
	}
	return nil
}

// writeTriggers writes the CREATE TRIGGER statements owned by a table or view.
func writeTriggers(buf *strings.Builder, triggers []*metadatapb.TriggerMetadata) error {
	for _, trigger := range triggers {
		if err := writeTrigger(buf, trigger); err != nil {
			return err
		}
	}
	return nil
}

// writeTrigger writes a CREATE TRIGGER statement
func writeTrigger(buf *strings.Builder, trigger *metadatapb.TriggerMetadata) error {
	// The trigger body should already contain the full CREATE TRIGGER statement
	if _, err := buf.WriteString(trigger.Body); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.TrimSpace(trigger.Body), ";") {
		if _, err := buf.WriteString(";"); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("\n\n"); err != nil {
		return err
	}
	return nil
}
