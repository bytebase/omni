package catalog

import (
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle"
	"github.com/bytebase/omni/oracle/ast"
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
			dependencyID := getObjectID(dependencyOwner(dependency, databaseName), dependency.Table)
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
			dependencyID := getObjectID(dependencyOwner(dependency, databaseName), dependency.Table)
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
			dependencyID := getObjectID(tableDependencyOwner(dependency, databaseName), dependency.Table)
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
		// A view that calls a function needs it before the view, as a virtual
		// column does.
		checkView := func(definition string) error {
			for _, function := range schema.Functions {
				if !prewrittenFunctions[function.Name] && expressionReferencesFunction(definition, function.Name) {
					if err := writeFunction(&buf, function); err != nil {
						return err
					}
					prewrittenFunctions[function.Name] = true
				}
			}
			return nil
		}
		for _, view := range schema.Views {
			if err := checkView(view.Definition); err != nil {
				return "", err
			}
		}
		for _, definition := range materializedViewDefinitions(schema) {
			if err := checkView(definition); err != nil {
				return "", err
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
		for _, pkg := range schema.Packages {
			if err := writePackage(&buf, pkg); err != nil {
				return "", err
			}
		}
		if err := writeDeferredFunctionIndexes(&buf, tables); err != nil {
			return "", err
		}
		if err := writeDeferredMaterializedViewFunctionIndexes(&buf, schema.MaterializedViews); err != nil {
			return "", err
		}
		if err := writeSynonyms(&buf, schema.Synonyms, to.PublicSynonyms, databaseName); err != nil {
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
	for _, pkg := range schema.Packages {
		if err := writePackage(&buf, pkg); err != nil {
			return "", err
		}
	}
	if err := writeDeferredFunctionIndexes(&buf, tables); err != nil {
		return "", err
	}
	if err := writeDeferredMaterializedViewFunctionIndexes(&buf, schema.MaterializedViews); err != nil {
		return "", err
	}
	if err := writeSynonyms(&buf, schema.Synonyms, to.PublicSynonyms, databaseName); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// dependencyOwner names the schema a dependency belongs to: the one the
// snapshot gives, or, when it gives none, the schema being defined. A synced
// snapshot leaves the owner of an unqualified name empty, so connecting the
// raw value would attach the edge to a ghost node and lose the order.
func dependencyOwner(dependency *metadatapb.DependencyColumn, schema string) string {
	if dependency.GetSchema() != "" {
		return dependency.GetSchema()
	}
	return schema
}

// tableDependencyOwner is dependencyOwner for a function's dependency, which
// carries the same owner rule in another message.
func tableDependencyOwner(dependency *metadatapb.DependencyTable, schema string) string {
	if dependency.GetSchema() != "" {
		return dependency.GetSchema()
	}
	return schema
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
	addEdges := func(key, owner string, dependencies []*metadatapb.DependencyColumn) {
		for _, dependency := range dependencies {
			// A dependency the snapshot leaves unowned belongs to the view's
			// schema: without that, a synced view whose text names its
			// relations unqualified connects to a ghost node and keeps a
			// dependent view ahead of the view it reads.
			if owner != "" && dependency.GetSchema() != "" && dependency.GetSchema() != owner {
				continue
			}
			for _, candidate := range []string{"v:" + dependency.Table, "m:" + dependency.Table} {
				if views[candidate] != nil || matViews[candidate] != nil {
					graph.AddEdge(candidate, key)
				}
			}
		}
	}
	for key, view := range views {
		addEdges(key, schema.Name, view.DependencyColumns)
	}
	for key, view := range matViews {
		addEdges(key, schema.Name, view.DependencyColumns)
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

// materializedViewDefinitions returns the query text of each materialized view,
// for the passes that ask whether a definition calls a function.
func materializedViewDefinitions(schema *metadatapb.SchemaMetadata) []string {
	var definitions []string
	for _, view := range schema.MaterializedViews {
		definitions = append(definitions, view.GetDefinition())
	}
	return definitions
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

func GetPackageDefinition(_ string, pkg *metadatapb.PackageMetadata) (string, error) {
	var buf strings.Builder
	if err := writePackage(&buf, pkg); err != nil {
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
		if !constraint.IsConstraint || needsReusedIndexClause(constraint) {
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

	var deferredForeignKeys []*metadatapb.ForeignKeyMetadata
	if includeForeignKeys {
		for _, fk := range table.ForeignKeys {
			// A self reference to a key a reused index backs waits for the
			// constraint, which is added after the table.
			if selfReferenceTargetsReusedConstraint(table, fk) {
				deferredForeignKeys = append(deferredForeignKeys, fk)
				continue
			}
			if err := writeSeparator(); err != nil {
				return err
			}
			if err := writeForeignKey(buf, fk); err != nil {
				return err
			}
		}
	}
	if _, err := buf.WriteString("\n)"); err != nil {
		return err
	}
	if err := writeTablePartitions(buf, table.Partitions); err != nil {
		return err
	}
	if _, err := buf.WriteString(";\n\n"); err != nil {
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
	if err := writeReusedIndexConstraints(buf, table, deferredForeignKeys); err != nil {
		return err
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
		if _, err := buf.WriteString(keyColumn(column)); err != nil {
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
		if _, err := buf.WriteString(keyColumn(column)); err != nil {
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

// needsReusedIndexClause reports a constraint the snapshot says an existing
// index backs, where that index is not the constraint's own: it is not unique,
// or its keys are not the constraint's. The index must exist before the
// constraint, which then names it with USING INDEX. The driver's sync states
// this shape for ALTER TABLE ... ADD CONSTRAINT ... USING INDEX.
func needsReusedIndexClause(index *metadatapb.IndexMetadata) bool {
	constraints := index.GetConstraints()
	if len(constraints) == 0 {
		return false
	}
	if !index.GetUnique() {
		return true
	}
	for _, constraint := range constraints {
		// A constraint that reuses the index keeps its own name.
		if constraint.GetName() != index.GetName() {
			return true
		}
		columns := constraint.GetColumns()
		if len(columns) == 0 {
			columns = index.GetExpressions()
		}
		if len(columns) != len(index.GetExpressions()) {
			return true
		}
		for i, column := range columns {
			if unquoteIdentifier(column) != unquoteIdentifier(index.GetExpressions()[i]) {
				return true
			}
		}
	}
	return false
}

// writeReusedIndexConstraints writes a table's indexes whose constraint needs
// the index clause: the index first, then the constraint naming it.
func selfReferenceTargetsReusedConstraint(table *metadatapb.TableMetadata, fk *metadatapb.ForeignKeyMetadata) bool {
	if fk.GetReferencedTable() != table.GetName() {
		return false
	}
	for _, index := range table.GetIndexes() {
		if !index.GetIsConstraint() || !needsReusedIndexClause(index) {
			continue
		}
		for _, constraint := range index.GetConstraints() {
			columns := constraint.GetColumns()
			if len(columns) == 0 {
				columns = index.GetExpressions()
			}
			if len(fk.GetReferencedColumns()) == 0 {
				if constraint.GetPrimary() {
					return true
				}
				continue
			}
			if len(fk.GetReferencedColumns()) != len(columns) {
				continue
			}
			matches := true
			for i, column := range fk.GetReferencedColumns() {
				if unquoteIdentifier(column) != unquoteIdentifier(columns[i]) {
					matches = false
				}
			}
			if matches {
				return true
			}
		}
	}
	return false
}

func writeReusedIndexConstraints(buf *strings.Builder, table *metadatapb.TableMetadata, deferredForeignKeys []*metadatapb.ForeignKeyMetadata) error {
	for _, index := range table.Indexes {
		if !index.GetIsConstraint() || !needsReusedIndexClause(index) {
			continue
		}
		if err := writeIndex(buf, table.Name, index); err != nil {
			return err
		}
		constraints := index.GetConstraints()
		if len(constraints) == 0 {
			constraints = []*metadatapb.IndexConstraintMetadata{{Name: index.GetName(), Columns: index.GetExpressions(), Primary: index.GetPrimary()}}
		}
		for _, constraint := range constraints {
			if _, err := buf.WriteString(`ALTER TABLE "` + table.Name + `" ADD CONSTRAINT "` + constraint.GetName() + `"`); err != nil {
				return err
			}
			if constraint.GetPrimary() {
				if _, err := buf.WriteString(` PRIMARY KEY (`); err != nil {
					return err
				}
			} else {
				if _, err := buf.WriteString(` UNIQUE (`); err != nil {
					return err
				}
			}
			columns := constraint.GetColumns()
			if len(columns) == 0 {
				columns = index.GetExpressions()
			}
			for i, column := range columns {
				if i != 0 {
					if _, err := buf.WriteString(", "); err != nil {
						return err
					}
				}
				if _, err := buf.WriteString(keyColumn(column)); err != nil {
					return err
				}
			}
			if _, err := buf.WriteString(`) USING INDEX "` + index.GetName() + `";

`); err != nil {
				return err
			}
		}
	}
	for _, fk := range deferredForeignKeys {
		if err := writeAddForeignKey(buf, table.Name, fk); err != nil {
			return err
		}
	}
	return nil
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
			if _, err := buf.WriteString(keyColumn(column)); err != nil {
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
			if _, err := buf.WriteString(keyColumn(column)); err != nil {
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
	if collation := explicitOracleCollation(column.Collation); collation != "" {
		if _, err := buf.WriteString(` COLLATE `); err != nil {
			return err
		}
		if _, err := buf.WriteString(collation); err != nil {
			return err
		}
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
	if _, err := buf.WriteString(`"`); err != nil {
		return err
	}
	if viewNeedsColumnList(view) {
		if err := writeViewColumnList(buf, view.Columns); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(` AS `); err != nil {
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
	if err := writeTriggers(buf, view.Triggers); err != nil {
		return err
	}
	return writeDefinitionComment(buf, "TABLE", view.Name, view.Comment)
}

// explicitOracleCollation returns collation when it names an explicit
// non-default Oracle column collation. USING_NLS_COMP is Oracle's default and
// fails to re-apply on MAX_STRING_SIZE=STANDARD instances.
func explicitOracleCollation(collation string) string {
	trimmed := strings.TrimSpace(collation)
	if trimmed == "" || strings.EqualFold(trimmed, "USING_NLS_COMP") {
		return ""
	}
	return trimmed
}

func viewNeedsColumnList(view *metadatapb.ViewMetadata) bool {
	if view == nil || len(view.Columns) == 0 {
		return false
	}
	for _, column := range view.Columns {
		if column.GetIsInvisible() {
			return true
		}
	}
	inferred, ok := inferredViewQueryColumnNames(view.GetDefinition())
	if !ok || len(inferred) != len(view.Columns) {
		return true
	}
	for i, column := range view.Columns {
		if inferred[i] != column.GetName() {
			return true
		}
	}
	return false
}

func inferredViewQueryColumnNames(definition string) ([]string, bool) {
	return inferredViewQueryColumnNamesForSchema(definition, nil)
}

type viewFromSource struct {
	alias  string
	schema string
	name   string
}

func inferredViewQueryColumnNamesForSchema(definition string, schema *metadatapb.SchemaMetadata) ([]string, bool) {
	query := strings.TrimSpace(definition)
	if query == "" {
		return nil, false
	}
	stmts, err := oracle.Parse("CREATE VIEW OMNI_PROBE AS " + query)
	if err != nil || len(stmts) == 0 {
		return nil, false
	}
	cv, ok := stmts[0].AST.(*ast.CreateViewStmt)
	if !ok || cv.Query == nil {
		return nil, false
	}
	sel, ok := cv.Query.(*ast.SelectStmt)
	if !ok {
		return nil, false
	}
	for sel.Larg != nil {
		sel = sel.Larg
	}
	if sel.TargetList == nil || len(sel.TargetList.Items) == 0 {
		return nil, false
	}
	names := make([]string, 0, len(sel.TargetList.Items))
	for _, item := range sel.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok {
			return nil, false
		}
		if target.Name != "" {
			names = append(names, target.Name)
			continue
		}
		if _, isStar := target.Expr.(*ast.Star); isStar {
			expanded, ok := expandViewWildcardColumns(sel, "", schema)
			if !ok {
				return nil, false
			}
			names = append(names, expanded...)
			continue
		}
		colRef, ok := target.Expr.(*ast.ColumnRef)
		if !ok || colRef.Column == "" {
			return nil, false
		}
		if colRef.Column == "*" {
			expanded, ok := expandViewWildcardColumns(sel, colRef.Table, schema)
			if !ok {
				return nil, false
			}
			names = append(names, expanded...)
			continue
		}
		names = append(names, colRef.Column)
	}
	return names, true
}

func expandViewWildcardColumns(sel *ast.SelectStmt, qualifier string, schema *metadatapb.SchemaMetadata) ([]string, bool) {
	if schema == nil || sel.FromClause == nil || len(sel.FromClause.Items) == 0 {
		return nil, false
	}
	var sources []viewFromSource
	if !collectSelectFromSources(sel.FromClause.Items, &sources) || len(sources) == 0 {
		return nil, false
	}
	lookupColumns := func(src viewFromSource) ([]string, bool) {
		if src.schema != "" && schema.GetName() != "" && src.schema != schema.GetName() {
			return nil, false
		}
		for _, table := range schema.GetTables() {
			if table.GetName() == src.name {
				if len(table.GetColumns()) == 0 {
					return nil, false
				}
				var cols []string
				for _, col := range table.GetColumns() {
					if !col.GetIsInvisible() {
						cols = append(cols, col.GetName())
					}
				}
				return cols, len(cols) > 0
			}
		}
		for _, view := range schema.GetViews() {
			if view.GetName() == src.name {
				if len(view.GetColumns()) == 0 {
					return nil, false
				}
				var cols []string
				for _, col := range view.GetColumns() {
					if !col.GetIsInvisible() {
						cols = append(cols, col.GetName())
					}
				}
				return cols, len(cols) > 0
			}
		}
		return nil, false
	}
	if qualifier != "" {
		var matched []string
		found := false
		for _, src := range sources {
			if src.alias == qualifier {
				if found {
					return nil, false
				}
				cols, ok := lookupColumns(src)
				if !ok {
					return nil, false
				}
				matched = cols
				found = true
			}
		}
		return matched, found
	}
	var out []string
	for _, src := range sources {
		cols, ok := lookupColumns(src)
		if !ok {
			return nil, false
		}
		out = append(out, cols...)
	}
	return out, len(out) > 0
}

func collectSelectFromSources(items []ast.Node, out *[]viewFromSource) bool {
	for _, item := range items {
		switch node := item.(type) {
		case *ast.TableRef:
			if node.Name == nil || node.Dblink != "" || node.Name.DBLink != "" {
				return false
			}
			name := objectName(node.Name)
			if name == "" {
				return false
			}
			alias := name
			if node.Alias != nil && node.Alias.Name != "" {
				alias = node.Alias.Name
			}
			*out = append(*out, viewFromSource{
				alias:  alias,
				schema: node.Name.Schema,
				name:   name,
			})
		case *ast.JoinClause:
			if (node.Using != nil && node.Using.Len() > 0) ||
				node.Type == ast.JOIN_NATURAL_INNER ||
				node.Type == ast.JOIN_NATURAL_LEFT ||
				node.Type == ast.JOIN_NATURAL_RIGHT ||
				node.Type == ast.JOIN_NATURAL_FULL {
				return false
			}
			if !collectSelectFromSources([]ast.Node{node.Left, node.Right}, out) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func writeViewColumnList(buf *strings.Builder, columns []*metadatapb.ColumnMetadata) error {
	if _, err := buf.WriteString(` (`); err != nil {
		return err
	}
	for i, column := range columns {
		if i > 0 {
			if _, err := buf.WriteString(`, `); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(`"` + column.GetName() + `"`); err != nil {
			return err
		}
		if column.GetIsInvisible() {
			if _, err := buf.WriteString(` INVISIBLE`); err != nil {
				return err
			}
		}
	}
	_, err := buf.WriteString(`)`)
	return err
}

// keyColumn renders a key column the snapshot states. The driver's sync keeps a
// case-sensitive name quoted, and an extracted one raw, so a quoted name is
// written as it is and a raw one is quoted.
func keyColumn(column string) string {
	if strings.HasPrefix(column, `"`) && strings.HasSuffix(column, `"`) && len(column) > 1 {
		return column
	}
	return `"` + column + `"`
}

// unquoteIdentifier strips the quotes from a bare quoted identifier, which the
// sync states but an extracted snapshot does not. An expression is left alone.
func unquoteIdentifier(identifier string) string {
	if strings.HasPrefix(identifier, `"`) && strings.HasSuffix(identifier, `"`) && len(identifier) > 1 && !strings.Contains(identifier[1:len(identifier)-1], `"`) {
		return identifier[1 : len(identifier)-1]
	}
	return identifier
}

// writeDefinitionComment writes the COMMENT statement a synced object carries.
// Oracle comments a view as a table and a materialized view by its own kind.
func writeDefinitionComment(buf *strings.Builder, object, name, comment string) error {
	if comment == "" {
		return nil
	}
	if _, err := buf.WriteString("COMMENT ON " + object + ` "` + name + `" IS '`); err != nil {
		return err
	}
	if _, err := buf.WriteString(strings.ReplaceAll(comment, "'", "''")); err != nil {
		return err
	}
	_, err := buf.WriteString("';\n\n")
	return err
}

// virtualColumnClause writes the clause that follows a virtual column's type
// and visibility.
func virtualColumnClause(generation *metadatapb.GenerationMetadata) string {
	return " GENERATED ALWAYS AS (" + generation.Expression + ") VIRTUAL"
}

// writeMaterializedView writes a CREATE MATERIALIZED VIEW statement
func materializedViewDefinitionHasClause(definition string) bool {
	upper := strings.ToUpper(strings.TrimSpace(definition))
	return strings.HasPrefix(upper, "BUILD ") ||
		strings.HasPrefix(upper, "REFRESH ") ||
		strings.HasPrefix(upper, "NEVER ") ||
		strings.HasPrefix(upper, "ON ")
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
	if err := writeDefinitionComment(buf, "MATERIALIZED VIEW", view.Name, view.Comment); err != nil {
		return err
	}
	for _, index := range materializedViewIndexes(view.Indexes) {
		// A function-based index waits for the function it calls; the deferred
		// pass after the functions writes it.
		if strings.Contains(index.GetType(), "FUNCTION-BASED") {
			continue
		}
		if err := writeIndex(buf, view.Name, index); err != nil {
			return err
		}
	}
	return nil
}

// writeDeferredMaterializedViewFunctionIndexes writes the function-based
// indexes a materialized view's definition holds, after the functions they
// call exist.
func writeDeferredMaterializedViewFunctionIndexes(buf *strings.Builder, views []*metadatapb.MaterializedViewMetadata) error {
	for _, view := range views {
		for _, index := range materializedViewIndexes(view.Indexes) {
			if !strings.Contains(index.GetType(), "FUNCTION-BASED") {
				continue
			}
			if err := writeIndex(buf, view.Name, index); err != nil {
				return err
			}
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

// writePackage writes a CREATE OR REPLACE PACKAGE [BODY] statement (or both
// specification and body when stored together).
func writePackage(buf *strings.Builder, pkg *metadatapb.PackageMetadata) error {
	definition := strings.TrimSpace(pkg.GetDefinition())
	if definition == "" {
		return nil
	}
	if !strings.HasPrefix(strings.ToUpper(definition), "CREATE") {
		if _, err := buf.WriteString("CREATE OR REPLACE "); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(definition); err != nil {
		return err
	}
	if !strings.HasSuffix(definition, ";") {
		if _, err := buf.WriteString(";"); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("\n\n"); err != nil {
		return err
	}
	return nil
}

func writeSynonyms(buf *strings.Builder, synonyms, publicSynonyms []*metadatapb.SynonymMetadata, databaseName string) error {
	for _, synonym := range synonyms {
		if err := writeSynonym(buf, synonym, databaseName, false); err != nil {
			return err
		}
	}
	for _, synonym := range dumpablePublicSynonyms(publicSynonyms, databaseName) {
		if err := writeSynonym(buf, synonym, databaseName, true); err != nil {
			return err
		}
	}
	return nil
}

// dumpablePublicSynonyms returns the public synonyms that belong to the schema
// being defined. Oracle's ALL_SYNONYMS WHERE OWNER = 'PUBLIC' includes
// thousands of built-in SYS/XDB/MDSYS synonyms that must not be emitted into a
// user schema's DDL.
func dumpablePublicSynonyms(publicSynonyms []*metadatapb.SynonymMetadata, databaseName string) []*metadatapb.SynonymMetadata {
	var out []*metadatapb.SynonymMetadata
	for _, synonym := range publicSynonyms {
		if synonym.GetName() == "" || synonym.GetTargetName() == "" {
			continue
		}
		targetSchema := synonym.GetTargetSchema()
		if databaseName != "" {
			if targetSchema != "" && targetSchema != databaseName {
				continue
			}
		} else if isOracleSystemSchema(targetSchema) {
			continue
		}
		out = append(out, synonym)
	}
	return out
}

func isOracleSystemSchema(schema string) bool {
	switch strings.ToUpper(schema) {
	case "ANONYMOUS", "APPQOSSYS", "AUDSYS", "CTXSYS", "DBSFWUSER", "DBSNMP",
		"DGPDB_INT", "DIP", "DVF", "DVSYS", "EXFSYS", "FLOWS_FILES", "GGSYS",
		"GSMADMIN_INTERNAL", "GSMCATUSER", "GSMROOTUSER", "GSMUSER", "LBACSYS",
		"MDDATA", "MDSYS", "MGMT_VIEW", "OJVMSYS", "OLAPSYS", "OPS$ORACLE",
		"ORACLE_OCM", "ORDDATA", "ORDPLUGINS", "ORDSYS", "OUTLN", "OWBSYS",
		"OWBSYS_AUDIT", "PUBLIC", "REMOTE_SCHEDULER_AGENT", "SI_INFORMTN_SCHEMA",
		"SPATIAL_CSW_ADMIN_USR", "SPATIAL_WFS_ADMIN_USR", "SYS", "SYS$UMF",
		"SYSBACKUP", "SYSDG", "SYSKM", "SYSMAN", "SYSRAC", "SYSTEM", "WMSYS",
		"XDB", "XS$NULL", "XS$$NULL":
		return true
	}
	return strings.HasPrefix(strings.ToUpper(schema), "APEX_")
}

func writeSynonym(buf *strings.Builder, synonym *metadatapb.SynonymMetadata, databaseName string, public bool) error {
	if synonym.GetName() == "" || synonym.GetTargetName() == "" {
		return nil
	}
	if public {
		if _, err := buf.WriteString(`CREATE OR REPLACE PUBLIC SYNONYM "`); err != nil {
			return err
		}
	} else {
		if _, err := buf.WriteString(`CREATE OR REPLACE SYNONYM "`); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(synonym.GetName()); err != nil {
		return err
	}
	if _, err := buf.WriteString(`" FOR `); err != nil {
		return err
	}
	if targetSchema := synonym.GetTargetSchema(); targetSchema != "" && (synonym.GetTargetDbLink() != "" || targetSchema != databaseName) {
		if _, err := buf.WriteString(`"` + targetSchema + `".`); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(`"` + synonym.GetTargetName() + `"`); err != nil {
		return err
	}
	if dbLink := synonym.GetTargetDbLink(); dbLink != "" {
		if _, err := buf.WriteString(`@` + formatDBLinkName(dbLink)); err != nil {
			return err
		}
	}
	_, err := buf.WriteString(";\n\n")
	return err
}

func formatDBLinkName(dbLink string) string {
	// Oracle database link names often include a domain suffix separated by dots
	// (e.g., REMOTE.REGRESS.RDBMS.DEV.US.ORACLE.COM). Simple dotted identifiers
	// stay unquoted so Oracle parses the domain cleanly.
	for _, part := range strings.Split(dbLink, ".") {
		if !isSimpleIdentifier(part) {
			return `"` + dbLink + `"`
		}
	}
	return dbLink
}

func partitionTypeKeyword(t metadatapb.TablePartitionMetadata_Type) string {
	switch t {
	case metadatapb.TablePartitionMetadata_RANGE, metadatapb.TablePartitionMetadata_RANGE_COLUMNS:
		return "RANGE"
	case metadatapb.TablePartitionMetadata_LIST, metadatapb.TablePartitionMetadata_LIST_COLUMNS:
		return "LIST"
	case metadatapb.TablePartitionMetadata_HASH, metadatapb.TablePartitionMetadata_LINEAR_HASH:
		return "HASH"
	default:
		return ""
	}
}

func formatPartitionKeyExpression(expr string) string {
	parts := strings.Split(expr, ",")
	quoted := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		quoted = append(quoted, keyColumn(trimmed))
	}
	return strings.Join(quoted, ", ")
}

func writeTablePartitions(buf *strings.Builder, partitions []*metadatapb.TablePartitionMetadata) error {
	if len(partitions) == 0 {
		return nil
	}
	first := partitions[0]
	kind := partitionTypeKeyword(first.GetType())
	if kind == "" || strings.TrimSpace(first.GetExpression()) == "" {
		return nil
	}
	if _, err := buf.WriteString(" PARTITION BY " + kind + " (" + formatPartitionKeyExpression(first.GetExpression()) + ")"); err != nil {
		return err
	}
	if interval := strings.TrimSpace(first.GetIntervalClause()); interval != "" {
		if _, err := buf.WriteString(" INTERVAL (" + interval + ")"); err != nil {
			return err
		}
	}
	var subKind string
	for _, p := range partitions {
		if len(p.GetSubpartitions()) > 0 {
			sub := p.GetSubpartitions()[0]
			subKind = partitionTypeKeyword(sub.GetType())
			if subKind != "" && strings.TrimSpace(sub.GetExpression()) != "" {
				if _, err := buf.WriteString(" SUBPARTITION BY " + subKind + " (" + formatPartitionKeyExpression(sub.GetExpression()) + ")"); err != nil {
					return err
				}
			}
			break
		}
	}
	if _, err := buf.WriteString(" (\n"); err != nil {
		return err
	}
	for i, part := range partitions {
		if i > 0 {
			if _, err := buf.WriteString(",\n"); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(`  PARTITION "`); err != nil {
			return err
		}
		if _, err := buf.WriteString(part.GetName()); err != nil {
			return err
		}
		if _, err := buf.WriteString(`"`); err != nil {
			return err
		}
		switch kind {
		case "RANGE":
			if _, err := buf.WriteString(" VALUES LESS THAN (" + strings.TrimSpace(part.GetValue()) + ")"); err != nil {
				return err
			}
		case "LIST":
			if _, err := buf.WriteString(" VALUES (" + strings.TrimSpace(part.GetValue()) + ")"); err != nil {
				return err
			}
		}
		if subKind != "" && len(part.GetSubpartitions()) > 0 {
			if _, err := buf.WriteString(" (\n"); err != nil {
				return err
			}
			for j, sub := range part.GetSubpartitions() {
				if j > 0 {
					if _, err := buf.WriteString(",\n"); err != nil {
						return err
					}
				}
				if _, err := buf.WriteString(`    SUBPARTITION "` + sub.GetName() + `"`); err != nil {
					return err
				}
				switch subKind {
				case "RANGE":
					if _, err := buf.WriteString(" VALUES LESS THAN (" + strings.TrimSpace(sub.GetValue()) + ")"); err != nil {
						return err
					}
				case "LIST":
					if _, err := buf.WriteString(" VALUES (" + strings.TrimSpace(sub.GetValue()) + ")"); err != nil {
						return err
					}
				}
			}
			if _, err := buf.WriteString("\n  )"); err != nil {
				return err
			}
		}
	}
	_, err := buf.WriteString("\n)")
	return err
}
