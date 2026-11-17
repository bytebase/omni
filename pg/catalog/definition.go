package catalog

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

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
	// PrintHeader prefixes the restore script with pg_dump's session settings
	// and sets the default tablespace before the first table.
	PrintHeader bool
	// SDLFormat renders the declarative form instead of the restore script.
	// It ignores PrintHeader.
	SDLFormat bool
}

// DatabaseDefinition renders a synced database snapshot as a pg_dump-style
// restore script, or with opts.SDLFormat as SDL.
func DatabaseDefinition(metadata *metadatapb.DatabaseSchemaMetadata, opts DefinitionOptions) (string, error) {
	if len(metadata.Schemas) == 0 && len(metadata.ForeignDataWrappers) == 0 && len(metadata.ForeignServers) == 0 {
		return "", nil
	}

	var buf strings.Builder

	if opts.SDLFormat {
		return databaseSDL(metadata)
	}

	if opts.PrintHeader {
		if _, err := buf.WriteString(header); err != nil {
			return "", err
		}
	}

	// Construct schemas.
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		if err := writeSchema(&buf, schema); err != nil {
			return "", err
		}
		if err := writeSchemaOwnerAndComment(&buf, schema); err != nil {
			return "", err
		}
	}

	// Construct extensions.
	for _, extension := range metadata.Extensions {
		if err := writeExtension(&buf, extension); err != nil {
			return "", err
		}
	}

	// Construct user-defined collations, before the columns that sort with
	// them.
	for _, schema := range metadata.Schemas {
		for _, collation := range schema.Collations {
			if collation.SkipDump {
				continue
			}
			if err := writeCollation(&buf, schema.Name, collation); err != nil {
				return "", err
			}
		}
	}

	// Construct enums, composite types and range types, each after the types it
	// is declared over.
	if err := writeTypesSection(&buf, collectTypes(metadata.Schemas)); err != nil {
		return "", err
	}

	// Build the graph for topological sort.
	graph := newDefinitionGraph()
	functionMap := make(map[string]*metadatapb.FunctionMetadata)
	tableMap := make(map[string]*metadatapb.TableMetadata)
	viewMap := make(map[string]*metadatapb.ViewMetadata)
	materializedViewMap := make(map[string]*metadatapb.MaterializedViewMetadata)
	foreignTableMap := make(map[string]*metadatapb.ExternalTableMetadata)
	foreignDataWrapperMap := make(map[string]*metadatapb.ForeignDataWrapperMetadata)
	foreignServerMap := make(map[string]*metadatapb.ForeignServerMetadata)
	externalTables := collectExternalTables(metadata.Schemas)
	foreignPartitions := collectForeignPartitionIDs(metadata.Schemas, externalTables)
	foreignPartitionParents := collectForeignPartitionParents(metadata.Schemas, externalTables)

	// Construct functions.
	functionIDsByName := make(map[string][]string)
	for _, schema := range metadata.Schemas {
		for _, function := range schema.Functions {
			if function.SkipDump {
				continue
			}
			funcID := getObjectID(schema.Name, metadataFuncIdentity(function))
			functionMap[funcID] = function
			qualified := schema.Name + "." + function.Name
			functionIDsByName[qualified] = append(functionIDsByName[qualified], funcID)
			functionIDsByName[function.Name] = append(functionIDsByName[function.Name], funcID)
			graph.AddNode(funcID)
			for _, dependency := range function.DependencyTables {
				dependencyID := getObjectID(dependency.Schema, dependency.Table)
				if parentID, ok := foreignPartitionParents[dependencyID]; ok {
					dependencyID = parentID
				}
				graph.AddEdge(dependencyID, funcID)
			}
		}
	}

	// A wrapper the database declares itself follows the functions its handler
	// and validator name; an extension-provided wrapper is already written.
	wrapperIDsByName := make(map[string]string)
	for _, wrapper := range metadata.ForeignDataWrappers {
		wrapperID := getObjectID("\x00fdw", wrapper.Name)
		foreignDataWrapperMap[wrapperID] = wrapper
		wrapperIDsByName[wrapper.Name] = wrapperID
		graph.AddNode(wrapperID)
		addFunctionNameEdges(graph, functionMap, functionIDsByName, wrapper.Handler, nil, wrapperID)
		addFunctionNameEdges(graph, functionMap, functionIDsByName, wrapper.Validator, []string{"text[]", "oid"}, wrapperID)
	}
	serverIDsByName := make(map[string]string)
	for _, server := range metadata.ForeignServers {
		serverID := getObjectID("\x00server", server.Name)
		foreignServerMap[serverID] = server
		serverIDsByName[server.Name] = serverID
		graph.AddNode(serverID)
		if wrapperID := wrapperIDsByName[server.Fdw]; wrapperID != "" {
			graph.AddEdge(wrapperID, serverID)
		}
	}

	// Create non-identity sequences before tables to prevent errors when tables reference
	// sequences owned by other tables. Identity sequences are created via writeColumnIdentityGeneration.
	// Build identity column map once for O(m*c) instead of O(s*m*c) when checking each sequence.
	identityColumnMap := buildIdentityColumnMap(metadata)

	sequenceOwnershipMap := make(map[string][]*metadatapb.SequenceMetadata)
	for _, schema := range metadata.Schemas {
		for _, sequence := range schema.Sequences {
			if sequence.SkipDump {
				continue
			}
			// Check if this sequence belongs to an identity column using the prebuilt map
			isIdentity := false
			if sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
				key := schema.Name + "." + sequence.OwnerTable + "." + sequence.OwnerColumn
				isIdentity = identityColumnMap[key]
			}

			if isIdentity {
				tableID := getObjectID(schema.Name, sequence.OwnerTable)
				sequenceOwnershipMap[tableID] = append(sequenceOwnershipMap[tableID], sequence)
				continue
			}
			if err := writeCreateSequence(&buf, schema.Name, sequence); err != nil {
				return "", err
			}
			if sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
				tableID := getObjectID(schema.Name, sequence.OwnerTable)
				sequenceOwnershipMap[tableID] = append(sequenceOwnershipMap[tableID], sequence)
			}
		}
	}

	if opts.PrintHeader {
		if _, err := buf.WriteString(setDefaultTableSpace); err != nil {
			return "", err
		}
	}

	// Construct tables, views and materialized views.
	for _, schema := range metadata.Schemas {
		for _, table := range schema.Tables {
			if table.SkipDump {
				continue
			}
			tableID := getObjectID(schema.Name, table.Name)
			tableMap[tableID] = table
			graph.AddNode(tableID)
			addColumnFunctionEdges(graph, functionMap, functionIDsByName, schema.Name, table)
			addForeignPartitionEdges(graph, functionMap, functionIDsByName, schema.Name, table, externalTables, serverIDsByName)
		}

		for _, externalTable := range schema.ExternalTables {
			if externalTable.SkipDump {
				continue
			}
			externalTableID := getObjectID(schema.Name, externalTable.Name)
			if foreignPartitions[externalTableID] {
				// A foreign partition is created with its parent, so it is not a
				// node of its own.
				continue
			}
			foreignTableMap[externalTableID] = externalTable
			graph.AddNode(externalTableID)
			addColumnFunctionEdges(graph, functionMap, functionIDsByName, schema.Name, &metadatapb.TableMetadata{
				Name:    externalTable.Name,
				Columns: externalTable.Columns,
			})
			addCheckFunctionEdges(graph, functionMap, functionIDsByName, schema.Name, externalTable.Name, externalTable.CheckConstraints)
			if serverID := serverIDsByName[externalTable.ExternalServerName]; serverID != "" {
				graph.AddEdge(serverID, externalTableID)
			}
		}

		for _, view := range schema.Views {
			if view.SkipDump {
				continue
			}
			viewID := getObjectID(schema.Name, view.Name)
			viewMap[viewID] = view
			graph.AddNode(viewID)
			for _, dependency := range view.DependencyColumns {
				dependencyID := getObjectID(dependency.Schema, dependency.Table)
				if parentID, ok := foreignPartitionParents[dependencyID]; ok {
					dependencyID = parentID
				}
				graph.AddEdge(dependencyID, viewID)
			}
		}

		for _, view := range schema.MaterializedViews {
			if view.SkipDump {
				continue
			}
			viewID := getObjectID(schema.Name, view.Name)
			materializedViewMap[viewID] = view
			graph.AddNode(viewID)
			for _, dependency := range view.DependencyColumns {
				dependencyID := getObjectID(dependency.Schema, dependency.Table)
				if parentID, ok := foreignPartitionParents[dependencyID]; ok {
					dependencyID = parentID
				}
				graph.AddEdge(dependencyID, viewID)
			}
		}
	}

	addRelationTypeEdges(graph, foreignTableMap, tableMap, foreignPartitionParents)

	orderedList, err := graph.TopologicalSort()
	if err != nil {
		// If topological sort fails (e.g., due to circular dependencies),
		// fall back to a safe dependency order: tables -> views -> materialized views -> functions
		orderedList = make([]string, 0, len(functionMap)+len(tableMap)+len(viewMap)+len(materializedViewMap)+len(foreignTableMap)+len(foreignDataWrapperMap)+len(foreignServerMap))

		// Add wrappers and servers first, then tables.
		var wrapperIDs []string
		for id := range foreignDataWrapperMap {
			wrapperIDs = append(wrapperIDs, id)
		}
		slices.Sort(wrapperIDs)
		orderedList = append(orderedList, wrapperIDs...)
		var serverIDs []string
		for id := range foreignServerMap {
			serverIDs = append(serverIDs, id)
		}
		slices.Sort(serverIDs)
		orderedList = append(orderedList, serverIDs...)

		// Add tables first (sorted alphabetically within each category)
		var tableIDs []string
		for id := range tableMap {
			tableIDs = append(tableIDs, id)
		}
		slices.Sort(tableIDs)
		orderedList = append(orderedList, tableIDs...)

		// Add foreign tables, whose defaults may call functions, before views.
		var foreignTableIDs []string
		for id := range foreignTableMap {
			foreignTableIDs = append(foreignTableIDs, id)
		}
		slices.Sort(foreignTableIDs)
		orderedList = append(orderedList, foreignTableIDs...)

		// Add views second
		var viewIDs []string
		for id := range viewMap {
			viewIDs = append(viewIDs, id)
		}
		slices.Sort(viewIDs)
		orderedList = append(orderedList, viewIDs...)

		// Add materialized views third
		var materializedViewIDs []string
		for id := range materializedViewMap {
			materializedViewIDs = append(materializedViewIDs, id)
		}
		slices.Sort(materializedViewIDs)
		orderedList = append(orderedList, materializedViewIDs...)

		// Add functions last
		var functionIDs []string
		for id := range functionMap {
			functionIDs = append(functionIDs, id)
		}
		slices.Sort(functionIDs)
		orderedList = append(orderedList, functionIDs...)
	}

	// Construct functions, wrappers, servers, tables, views and materialized
	// views in order.
	for _, objectID := range orderedList {
		if wrapper, ok := foreignDataWrapperMap[objectID]; ok {
			if err := writeForeignDataWrapper(&buf, wrapper); err != nil {
				return "", err
			}
			continue
		}
		if server, ok := foreignServerMap[objectID]; ok {
			if err := writeForeignServer(&buf, server); err != nil {
				return "", err
			}
			continue
		}
		if function, ok := functionMap[objectID]; ok {
			if err := writeFunction(&buf, getSchemaNameFromID(objectID), function); err != nil {
				return "", err
			}
			continue
		}
		if table, ok := tableMap[objectID]; ok {
			if err := writeTable(&buf, getSchemaNameFromID(objectID), table, sequenceOwnershipMap[objectID], externalTables); err != nil {
				return "", err
			}
		}
		if externalTable, ok := foreignTableMap[objectID]; ok {
			if err := writeForeignTable(&buf, getSchemaNameFromID(objectID), externalTable, externalTable.CheckConstraints); err != nil {
				return "", err
			}
			for _, sequence := range sequenceOwnershipMap[objectID] {
				if err := writeAlterSequenceOwnedBy(&buf, getSchemaNameFromID(objectID), sequence); err != nil {
					return "", err
				}
			}
			foreignSchema := getSchemaNameFromID(objectID)
			for _, column := range externalTable.Columns {
				if len(column.Comment) > 0 {
					if err := writeColumnComment(&buf, foreignSchema, externalTable.Name, column); err != nil {
						return "", err
					}
				}
			}
			for _, check := range externalTable.CheckConstraints {
				if len(check.Comment) > 0 {
					if err := writeConstraintComment(&buf, foreignSchema, externalTable.Name, check.Name, check.Comment); err != nil {
						return "", err
					}
				}
			}
		}
		if view, ok := viewMap[objectID]; ok {
			if err := writeView(&buf, getSchemaNameFromID(objectID), view); err != nil {
				return "", err
			}
			continue
		}
		if view, ok := materializedViewMap[objectID]; ok {
			if err := writeMaterializedView(&buf, getSchemaNameFromID(objectID), view); err != nil {
				return "", err
			}
		}
	}

	// Construct triggers.
	for _, schema := range metadata.Schemas {
		if err := writeSchemaTriggers(&buf, schema, true); err != nil {
			return "", err
		}
	}

	// Construct rules.
	for _, schema := range metadata.Schemas {
		if err := writeSchemaRules(&buf, schema); err != nil {
			return "", err
		}
	}

	// Construct foreign keys.
	tablesMissingColumns := collectTables(metadata.Schemas, tableMissingColumnMetadata)
	for _, schema := range metadata.Schemas {
		if err := writeSchemaForeignKeys(&buf, schema, tablesMissingColumns); err != nil {
			return "", err
		}
	}

	// Construct event triggers (last, as they depend on functions).
	for _, eventTrigger := range metadata.EventTriggers {
		if eventTrigger.SkipDump {
			continue
		}
		if err := writeEventTrigger(&buf, eventTrigger); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func SchemaDefinition(schema *metadatapb.SchemaMetadata) (string, error) {
	var buf strings.Builder
	if err := writeSchema(&buf, schema); err != nil {
		return "", err
	}
	if err := writeSchemaOwnerAndComment(&buf, schema); err != nil {
		return "", err
	}

	// Construct user-defined collations, before the columns that sort with them.
	for _, collation := range schema.Collations {
		if collation.SkipDump {
			continue
		}
		if err := writeCollation(&buf, schema.Name, collation); err != nil {
			return "", err
		}
	}

	// Construct enums, composite types and range types, each after the types it
	// is declared over. The schema-level SkipDump filter does not apply here:
	// the caller asked for this specific schema's definition.
	if err := writeTypesSection(&buf, collectSchemaTypes(schema)); err != nil {
		return "", err
	}

	// Build the graph for topological sort.
	graph := newDefinitionGraph()
	functionMap := make(map[string]*metadatapb.FunctionMetadata)
	tableMap := make(map[string]*metadatapb.TableMetadata)
	viewMap := make(map[string]*metadatapb.ViewMetadata)
	materializedViewMap := make(map[string]*metadatapb.MaterializedViewMetadata)

	// Construct functions.
	for _, function := range schema.Functions {
		if function.SkipDump {
			continue
		}
		funcID := getObjectID(schema.Name, metadataFuncIdentity(function))
		functionMap[funcID] = function
		graph.AddNode(funcID)
		for _, dependency := range function.DependencyTables {
			dependencyID := getObjectID(dependency.Schema, dependency.Table)
			graph.AddEdge(dependencyID, funcID)
		}
	}

	// Create non-identity sequences before tables to prevent errors when tables reference
	// sequences owned by other tables. Identity sequences are created via writeColumnIdentityGeneration.
	// Build identity column map once for O(m*c) instead of O(s*m*c) when checking each sequence.
	identityColumnMap := buildIdentityColumnMapForSchema(schema)

	sequenceOwnershipMap := make(map[string][]*metadatapb.SequenceMetadata)
	for _, sequence := range schema.Sequences {
		if sequence.SkipDump {
			continue
		}
		// Check if this sequence belongs to an identity column using the prebuilt map
		isIdentity := false
		if sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
			key := sequence.OwnerTable + "." + sequence.OwnerColumn
			isIdentity = identityColumnMap[key]
		}

		if isIdentity {
			tableID := getObjectID(schema.Name, sequence.OwnerTable)
			sequenceOwnershipMap[tableID] = append(sequenceOwnershipMap[tableID], sequence)
			continue
		}
		if err := writeCreateSequence(&buf, schema.Name, sequence); err != nil {
			return "", err
		}
		if sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
			tableID := getObjectID(schema.Name, sequence.OwnerTable)
			sequenceOwnershipMap[tableID] = append(sequenceOwnershipMap[tableID], sequence)
		}
	}

	// Construct tables, views and materialized views.
	for _, table := range schema.Tables {
		if table.SkipDump {
			continue
		}
		tableID := getObjectID(schema.Name, table.Name)
		tableMap[tableID] = table
		graph.AddNode(tableID)
	}

	for _, view := range schema.Views {
		if view.SkipDump {
			continue
		}
		viewID := getObjectID(schema.Name, view.Name)
		viewMap[viewID] = view
		graph.AddNode(viewID)
		for _, dependency := range view.DependencyColumns {
			dependencyID := getObjectID(dependency.Schema, dependency.Table)
			graph.AddEdge(dependencyID, viewID)
		}
	}

	for _, view := range schema.MaterializedViews {
		if view.SkipDump {
			continue
		}
		viewID := getObjectID(schema.Name, view.Name)
		materializedViewMap[viewID] = view
		graph.AddNode(viewID)
		for _, dependency := range view.DependencyColumns {
			dependencyID := getObjectID(dependency.Schema, dependency.Table)
			graph.AddEdge(dependencyID, viewID)
		}
	}

	orderedList, err := graph.TopologicalSort()
	if err != nil {
		// If topological sort fails (e.g., due to circular dependencies),
		// fall back to a safe dependency order: tables -> views -> materialized views -> functions
		orderedList = make([]string, 0, len(functionMap)+len(tableMap)+len(viewMap)+len(materializedViewMap))

		// Add tables first (sorted alphabetically within each category)
		var tableIDs []string
		for id := range tableMap {
			tableIDs = append(tableIDs, id)
		}
		slices.Sort(tableIDs)
		orderedList = append(orderedList, tableIDs...)

		// Add views second
		var viewIDs []string
		for id := range viewMap {
			viewIDs = append(viewIDs, id)
		}
		slices.Sort(viewIDs)
		orderedList = append(orderedList, viewIDs...)

		// Add materialized views third
		var materializedViewIDs []string
		for id := range materializedViewMap {
			materializedViewIDs = append(materializedViewIDs, id)
		}
		slices.Sort(materializedViewIDs)
		orderedList = append(orderedList, materializedViewIDs...)

		// Add functions last
		var functionIDs []string
		for id := range functionMap {
			functionIDs = append(functionIDs, id)
		}
		slices.Sort(functionIDs)
		orderedList = append(orderedList, functionIDs...)
	}

	// Construct functions, tables, views and materialized views in order.
	for _, objectID := range orderedList {
		if function, ok := functionMap[objectID]; ok {
			if err := writeFunction(&buf, getSchemaNameFromID(objectID), function); err != nil {
				return "", err
			}
			continue
		}
		if table, ok := tableMap[objectID]; ok {
			if err := writeTable(&buf, getSchemaNameFromID(objectID), table, sequenceOwnershipMap[objectID], collectExternalTables([]*metadatapb.SchemaMetadata{schema})); err != nil {
				return "", err
			}
		}
		if view, ok := viewMap[objectID]; ok {
			if err := writeView(&buf, getSchemaNameFromID(objectID), view); err != nil {
				return "", err
			}
			continue
		}
		if view, ok := materializedViewMap[objectID]; ok {
			if err := writeMaterializedView(&buf, getSchemaNameFromID(objectID), view); err != nil {
				return "", err
			}
		}
	}

	// Construct triggers.
	if err := writeSchemaTriggers(&buf, schema, false); err != nil {
		return "", err
	}

	// Construct rules.
	if err := writeSchemaRules(&buf, schema); err != nil {
		return "", err
	}

	// Construct foreign keys.
	if err := writeSchemaForeignKeys(&buf, schema, collectTables([]*metadatapb.SchemaMetadata{schema}, tableMissingColumnMetadata)); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func TableDefinition(schema string, table *metadatapb.TableMetadata, sequences []*metadatapb.SequenceMetadata) (string, error) {
	var buf strings.Builder
	if err := writeTable(&buf, schema, table, sequences, nil); err != nil {
		return "", err
	}
	// Skip triggers, rules, and foreign keys for a table without column
	// metadata; their definitions can reference the missing columns.
	if tableMissingColumnMetadata(table) {
		return buf.String(), nil
	}
	// Construct triggers.
	for _, trigger := range table.Triggers {
		if trigger.SkipDump {
			continue
		}
		if err := writeTrigger(&buf, schema, table.Name, trigger); err != nil {
			return "", err
		}
	}
	// Construct rules.
	if len(table.Rules) > 0 {
		if err := writeRules(&buf, schema, table.Name, table.Rules); err != nil {
			return "", err
		}
	}
	// Construct foreign keys.
	for _, fk := range table.ForeignKeys {
		if err := writeForeignKey(&buf, schema, table.Name, fk); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func ViewDefinition(schema string, view *metadatapb.ViewMetadata) (string, error) {
	var buf strings.Builder
	if err := writeView(&buf, schema, view); err != nil {
		return "", err
	}
	// Construct triggers.
	for _, trigger := range view.Triggers {
		if trigger.SkipDump {
			continue
		}
		if err := writeTrigger(&buf, schema, view.Name, trigger); err != nil {
			return "", err
		}
	}
	// Construct rules (non-SELECT rules only, as SELECT rules are part of the view definition).
	var nonSelectRules []*metadatapb.RuleMetadata
	for _, rule := range view.Rules {
		if rule.Event != "SELECT" {
			nonSelectRules = append(nonSelectRules, rule)
		}
	}
	if len(nonSelectRules) > 0 {
		if err := writeRules(&buf, schema, view.Name, nonSelectRules); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func MaterializedViewDefinition(schema string, view *metadatapb.MaterializedViewMetadata) (string, error) {
	var buf strings.Builder
	if err := writeMaterializedView(&buf, schema, view); err != nil {
		return "", err
	}
	// Construct triggers.
	for _, trigger := range view.Triggers {
		if trigger.SkipDump {
			continue
		}
		if err := writeTrigger(&buf, schema, view.Name, trigger); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func FunctionDefinition(schema string, function *metadatapb.FunctionMetadata) (string, error) {
	var buf strings.Builder
	if err := writeFunction(&buf, schema, function); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SequenceDefinition(schema string, sequence *metadatapb.SequenceMetadata) (string, error) {
	var buf strings.Builder
	if err := writeCreateSequence(&buf, schema, sequence); err != nil {
		return "", err
	}
	if sequence.OwnerColumn != "" && sequence.OwnerTable != "" {
		if err := writeAlterSequenceOwnedBy(&buf, schema, sequence); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func writeTrigger(out io.Writer, schema string, table string, trigger *metadatapb.TriggerMetadata) error {
	if _, err := io.WriteString(out, trigger.Body); err != nil {
		return err
	}
	if _, err := io.WriteString(out, ";\n\n"); err != nil {
		return err
	}

	if len(trigger.Comment) > 0 {
		if err := writeTriggerComment(out, schema, table, trigger); err != nil {
			return err
		}
	}

	return nil
}

func writeTriggerComment(out io.Writer, schema string, table string, trigger *metadatapb.TriggerMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON TRIGGER "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, trigger.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" ON "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, table); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(trigger.Comment)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeEventTrigger(out io.Writer, eventTrigger *metadatapb.EventTriggerMetadata) error {
	// Use the stored definition if available
	if eventTrigger.Definition != "" {
		if _, err := io.WriteString(out, eventTrigger.Definition); err != nil {
			return err
		}
		if _, err := io.WriteString(out, ";\n\n"); err != nil {
			return err
		}
	} else {
		// Build the CREATE EVENT TRIGGER statement
		if _, err := io.WriteString(out, `CREATE EVENT TRIGGER "`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, eventTrigger.Name); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `" ON `); err != nil {
			return err
		}
		if _, err := io.WriteString(out, eventTrigger.Event); err != nil {
			return err
		}

		// Add WHEN TAG IN clause if tags are specified
		if len(eventTrigger.Tags) > 0 {
			if _, err := io.WriteString(out, "\n  WHEN TAG IN ("); err != nil {
				return err
			}
			for i, tag := range eventTrigger.Tags {
				if i > 0 {
					if _, err := io.WriteString(out, ", "); err != nil {
						return err
					}
				}
				if _, err := io.WriteString(out, "'"); err != nil {
					return err
				}
				if _, err := io.WriteString(out, tag); err != nil {
					return err
				}
				if _, err := io.WriteString(out, "'"); err != nil {
					return err
				}
			}
			if _, err := io.WriteString(out, ")"); err != nil {
				return err
			}
		}

		// Add EXECUTE FUNCTION clause
		if _, err := io.WriteString(out, "\n  EXECUTE FUNCTION "); err != nil {
			return err
		}
		if eventTrigger.FunctionSchema != "" {
			if _, err := io.WriteString(out, `"`); err != nil {
				return err
			}
			if _, err := io.WriteString(out, eventTrigger.FunctionSchema); err != nil {
				return err
			}
			if _, err := io.WriteString(out, `".`); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, eventTrigger.FunctionName); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `"();\n\n`); err != nil {
			return err
		}
	}

	// A disabled event trigger is created enabled; restore its disabled state.
	if !eventTrigger.Enabled {
		if _, err := fmt.Fprintf(out, "ALTER EVENT TRIGGER \"%s\" DISABLE;\n\n", eventTrigger.Name); err != nil {
			return err
		}
	}

	// Write event trigger comment if present
	if eventTrigger.Comment != "" {
		if err := writeEventTriggerComment(out, eventTrigger); err != nil {
			return err
		}
	}

	return nil
}

func writeEventTriggerComment(out io.Writer, eventTrigger *metadatapb.EventTriggerMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON EVENT TRIGGER "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, eventTrigger.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(eventTrigger.Comment)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeEnum(out io.Writer, schema string, enum *metadatapb.EnumTypeMetadata) error {
	if _, err := io.WriteString(out, `CREATE TYPE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, enum.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" AS ENUM (\n"); err != nil {
		return err
	}
	for i, value := range enum.Values {
		if i > 0 {
			if _, err := io.WriteString(out, ",\n"); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, `    '`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, escapeSingleQuote(value)); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `'`); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, "\n)"); err != nil {
		return err
	}

	return nil
}

func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// escapeIdentifier doubles the double quotes in a name so it is safe inside a
// quoted identifier, as a role name with a double quote would otherwise be.
func escapeIdentifier(s string) string {
	return strings.ReplaceAll(s, `"`, `""`)
}

func writeCompositeType(out io.Writer, schema string, composite *metadatapb.CompositeTypeMetadata) error {
	if _, err := io.WriteString(out, `CREATE TYPE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, composite.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" AS (\n"); err != nil {
		return err
	}
	for i, attribute := range composite.Attributes {
		if i > 0 {
			if _, err := io.WriteString(out, ",\n"); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, `    "`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, attribute.Name); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `" `); err != nil {
			return err
		}
		if _, err := io.WriteString(out, attribute.Type); err != nil {
			return err
		}
		// The collation is stored as an emit-ready identifier reference
		// (quoted as needed, schema-qualified outside pg_catalog).
		if attribute.Collation != "" {
			if _, err := io.WriteString(out, ` COLLATE `); err != nil {
				return err
			}
			if _, err := io.WriteString(out, attribute.Collation); err != nil {
				return err
			}
		}
	}

	if _, err := io.WriteString(out, "\n)"); err != nil {
		return err
	}

	return nil
}

func compositeTypeHasComments(composite *metadatapb.CompositeTypeMetadata) bool {
	if composite.Comment != "" {
		return true
	}
	for _, attribute := range composite.Attributes {
		if attribute.Comment != "" {
			return true
		}
	}
	return false
}

func writeCompositeTypeComments(out io.Writer, schema string, composite *metadatapb.CompositeTypeMetadata) error {
	if composite.Comment != "" {
		if _, err := fmt.Fprintf(out, "COMMENT ON TYPE \"%s\".\"%s\" IS '%s';\n\n", schema, composite.Name, escapeSingleQuote(composite.Comment)); err != nil {
			return err
		}
	}
	for _, attribute := range composite.Attributes {
		if attribute.Comment == "" {
			continue
		}
		if _, err := fmt.Fprintf(out, "COMMENT ON COLUMN \"%s\".\"%s\".\"%s\" IS '%s';\n\n", schema, composite.Name, attribute.Name, escapeSingleQuote(attribute.Comment)); err != nil {
			return err
		}
	}
	return nil
}

// qualifiedCompositeType pairs a composite type with its schema for
// cross-schema ordering.
type qualifiedCompositeType struct {
	Schema    string
	Composite *metadatapb.CompositeTypeMetadata
}

// collectCompositeTypes gathers the dumpable composite types of the given
// schemas. Schemas with SkipDump never get a CREATE SCHEMA statement, so
// their types must not be emitted either.
// qualifiedType is one enum, composite type or range type of the dumpable
// snapshot, with the schema it lives in so the types it names can be resolved.
type qualifiedType struct {
	Schema    string
	Enum      *metadatapb.EnumTypeMetadata
	Composite *metadatapb.CompositeTypeMetadata
	RangeType *metadatapb.RangeTypeMetadata
}

const (
	typeKindEnum = iota
	typeKindComposite
	typeKindRange
)

func (t qualifiedType) kind() int {
	switch {
	case t.Enum != nil:
		return typeKindEnum
	case t.Composite != nil:
		return typeKindComposite
	default:
		return typeKindRange
	}
}

// name is the type's own name in its schema.
func (t qualifiedType) name() string {
	switch {
	case t.Enum != nil:
		return t.Enum.Name
	case t.Composite != nil:
		return t.Composite.Name
	default:
		return t.RangeType.Name
	}
}

// skipDump reports whether creating the type is another statement's work: an
// extension's, or a schema's the dump never creates.
func (t qualifiedType) skipDump() bool {
	switch {
	case t.Enum != nil:
		return t.Enum.SkipDump
	case t.Composite != nil:
		return t.Composite.SkipDump
	default:
		return t.RangeType.SkipDump
	}
}

// collectTypes gathers the dumpable enums, composite types and range types of
// the given schemas. Schemas with SkipDump never get a CREATE SCHEMA statement,
// so their types must not be emitted either.
func collectTypes(schemas []*metadatapb.SchemaMetadata) []qualifiedType {
	var kept []qualifiedType
	for _, schema := range schemas {
		if schema.SkipDump {
			continue
		}
		kept = append(kept, collectSchemaTypes(schema)...)
	}
	return kept
}

// collectSchemaTypes gathers one schema's dumpable types whether or not the
// schema itself is dumped: a caller that asked for this schema's definition gets
// it, though a type an extension owns is still left out.
func collectSchemaTypes(schema *metadatapb.SchemaMetadata) []qualifiedType {
	var types []qualifiedType
	for _, enum := range schema.EnumTypes {
		types = append(types, qualifiedType{Schema: schema.Name, Enum: enum})
	}
	for _, composite := range schema.CompositeTypes {
		types = append(types, qualifiedType{Schema: schema.Name, Composite: composite})
	}
	for _, rangeType := range schema.RangeTypes {
		types = append(types, qualifiedType{Schema: schema.Name, RangeType: rangeType})
	}
	kept := types[:0]
	for _, t := range types {
		if !t.skipDump() {
			kept = append(kept, t)
		}
	}
	return kept
}

// orderTypes orders the types so that each comes after the types it is declared
// over: a composite attribute's type, and a range's subtype, which may be
// another range's multirange type. A type the list does not hold adds no edge,
// since a pg_catalog type needs no statement. Ties go to enums, then composite
// types, then range types, in name order, so the section reads the way it did
// before the types could depend on each other. A cycle cannot reach the server;
// the rest of the list follows in that order if one did.
func orderTypes(types []qualifiedType) []qualifiedType {
	key := func(schema, name string) string { return schema + "\x00" + name }
	nodes := make(map[string]qualifiedType, len(types))
	for _, t := range types {
		nodes[key(t.Schema, t.name())] = t
	}
	// A range's multirange type is made by the range's own CREATE, so a subtype
	// that names it depends on that range.
	multiranges := make(map[string]string, len(nodes))
	for k, t := range nodes {
		if t.RangeType == nil || t.RangeType.MultirangeName == "" {
			continue
		}
		schema, name, ok := parseQualifiedTypeIdent(t.RangeType.MultirangeName)
		if !ok {
			schema, name = t.Schema, t.RangeType.MultirangeName
		}
		multiranges[key(schema, name)] = k
	}
	inDegree := make(map[string]int, len(nodes))
	dependents := make(map[string][]string, len(nodes))
	for k := range nodes {
		inDegree[k] = 0
	}
	addEdges := func(k string, references []string) {
		for _, reference := range references {
			depSchema, depName, ok := parseQualifiedTypeIdent(reference)
			if !ok {
				continue
			}
			depKey := key(depSchema, depName)
			if owner, isMultirange := multiranges[depKey]; isMultirange {
				depKey = owner
			}
			if depKey == k {
				continue
			}
			if _, exists := nodes[depKey]; !exists {
				continue
			}
			dependents[depKey] = append(dependents[depKey], k)
			inDegree[k]++
		}
	}
	for k, t := range nodes {
		switch {
		case t.Composite != nil:
			var attributes []string
			for _, attribute := range t.Composite.Attributes {
				attributes = append(attributes, attribute.Type)
			}
			addEdges(k, attributes)
		case t.RangeType != nil:
			addEdges(k, []string{t.RangeType.Subtype})
		}
	}

	ready := make([]string, 0, len(nodes))
	for k, d := range inDegree {
		if d == 0 {
			ready = append(ready, k)
		}
	}
	sortTypeKeys(ready, nodes)

	ordered := make([]qualifiedType, 0, len(nodes))
	for len(ready) > 0 {
		k := ready[0]
		ready = ready[1:]
		ordered = append(ordered, nodes[k])
		delete(inDegree, k)
		for _, dependent := range dependents[k] {
			if _, remaining := inDegree[dependent]; !remaining {
				continue
			}
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		sortTypeKeys(ready, nodes)
	}

	if len(ordered) < len(nodes) {
		var rest []string
		for k := range inDegree {
			rest = append(rest, k)
		}
		sortTypeKeys(rest, nodes)
		for _, k := range rest {
			ordered = append(ordered, nodes[k])
		}
	}

	return ordered
}

// sortTypeKeys orders type keys by kind, then name, so independent types keep a
// stable order across runs.
func sortTypeKeys(keys []string, nodes map[string]qualifiedType) {
	slices.SortFunc(keys, func(a, b string) int {
		if c := cmp.Compare(nodes[a].kind(), nodes[b].kind()); c != 0 {
			return c
		}
		return strings.Compare(a, b)
	})
}

// writeTypesSection emits the enums, composite types and range types in
// dependency order, each with its comments.
func writeTypesSection(out io.Writer, types []qualifiedType) error {
	for _, t := range orderTypes(types) {
		switch {
		case t.Enum != nil:
			if err := writeEnum(out, t.Schema, t.Enum); err != nil {
				return err
			}
			if _, err := io.WriteString(out, ";\n\n"); err != nil {
				return err
			}
			if t.Enum.Comment != "" {
				if err := writeEnumComment(out, t.Schema, t.Enum); err != nil {
					return err
				}
			}
		case t.Composite != nil:
			if err := writeCompositeType(out, t.Schema, t.Composite); err != nil {
				return err
			}
			if _, err := io.WriteString(out, ";\n\n"); err != nil {
				return err
			}
			if compositeTypeHasComments(t.Composite) {
				if err := writeCompositeTypeComments(out, t.Schema, t.Composite); err != nil {
					return err
				}
			}
		default:
			if err := writeRangeType(out, t.Schema, t.RangeType); err != nil {
				return err
			}
		}
	}
	return nil
}

func collectCompositeTypes(schemas []*metadatapb.SchemaMetadata) []qualifiedCompositeType {
	var compositeTypes []qualifiedCompositeType
	for _, schema := range schemas {
		if schema.SkipDump {
			continue
		}
		for _, composite := range schema.CompositeTypes {
			if composite.SkipDump {
				continue
			}
			compositeTypes = append(compositeTypes, qualifiedCompositeType{Schema: schema.Name, Composite: composite})
		}
	}
	return compositeTypes
}

// writeCompositeTypesSection emits the composite types with their comments,
// ordered so referenced composites come first.
func writeCompositeTypesSection(buf *strings.Builder, compositeTypes []qualifiedCompositeType) error {
	for _, t := range sortCompositeTypesTopologically(compositeTypes) {
		if err := writeCompositeType(buf, t.Schema, t.Composite); err != nil {
			return err
		}
		if _, err := buf.WriteString(";\n\n"); err != nil {
			return err
		}
		if compositeTypeHasComments(t.Composite) {
			if err := writeCompositeTypeComments(buf, t.Schema, t.Composite); err != nil {
				return err
			}
		}
	}
	return nil
}

// parseQualifiedTypeIdent splits a type reference produced by the sync layer
// (format_type output, schema-qualified for user-defined types) into
// (schema, name), ignoring any typmod or array suffix. It returns ok=false
// for unqualified references — built-in types are never schema-qualified.
func parseQualifiedTypeIdent(s string) (string, string, bool) {
	var parts []string
	for i := 0; i < len(s); {
		if s[i] == '"' {
			var b strings.Builder
			i++
			for i < len(s) {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						b.WriteByte('"')
						i += 2
						continue
					}
					i++
					break
				}
				b.WriteByte(s[i])
				i++
			}
			parts = append(parts, b.String())
		} else {
			start := i
			for i < len(s) && s[i] != '.' && s[i] != '(' && s[i] != '[' {
				i++
			}
			parts = append(parts, strings.TrimSpace(s[start:i]))
		}
		if i < len(s) && s[i] == '.' {
			i++
			continue
		}
		break
	}
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], true
	}
	return "", "", false
}

// sortCompositeTypesTopologically orders composite types so that a composite
// referenced by another composite's attribute is emitted first. PostgreSQL
// rejects self-containment, so the graph is acyclic; any leftover on an
// unexpected cycle falls back to (schema, name) order. A composite attribute
// typed with a table or view row type is not reordered here — that requires
// catalog-grade dependency data (see pg_dump's pg_depend-driven sort).
func sortCompositeTypesTopologically(types []qualifiedCompositeType) []qualifiedCompositeType {
	key := func(schema, name string) string { return schema + "\x00" + name }
	nodes := make(map[string]qualifiedCompositeType, len(types))
	for _, t := range types {
		nodes[key(t.Schema, t.Composite.Name)] = t
	}
	inDegree := make(map[string]int, len(nodes))
	dependents := make(map[string][]string, len(nodes))
	for k := range nodes {
		inDegree[k] = 0
	}
	for k, t := range nodes {
		for _, attribute := range t.Composite.Attributes {
			depSchema, depName, ok := parseQualifiedTypeIdent(attribute.Type)
			if !ok {
				continue
			}
			depKey := key(depSchema, depName)
			if depKey == k {
				continue
			}
			if _, exists := nodes[depKey]; !exists {
				continue
			}
			dependents[depKey] = append(dependents[depKey], k)
			inDegree[k]++
		}
	}

	var ready []string
	for k, d := range inDegree {
		if d == 0 {
			ready = append(ready, k)
		}
	}
	slices.Sort(ready)

	ordered := make([]qualifiedCompositeType, 0, len(nodes))
	for len(ready) > 0 {
		k := ready[0]
		ready = ready[1:]
		ordered = append(ordered, nodes[k])
		delete(inDegree, k)
		for _, dependent := range dependents[k] {
			if _, remaining := inDegree[dependent]; !remaining {
				continue
			}
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		slices.Sort(ready)
	}

	if len(ordered) < len(nodes) {
		var rest []string
		for k := range inDegree {
			rest = append(rest, k)
		}
		slices.Sort(rest)
		for _, k := range rest {
			ordered = append(ordered, nodes[k])
		}
	}

	return ordered
}

// signatureTypesOnly strips parameter names from a function/procedure signature,
// keeping only types. PostgreSQL's COMMENT ON FUNCTION/PROCEDURE requires
// type-only signatures: "func(integer, text)" not "func(id integer, name text)".
//
// It works by wrapping the signature in a CREATE FUNCTION statement, parsing it
// with omni to discover parameter names, then removing those names from the
// original signature text to preserve the original type spelling.
func signatureTypesOnly(signature string) string {
	openParen := strings.Index(signature, "(")
	if openParen < 0 {
		return signature
	}
	closeParen := strings.LastIndex(signature, ")")
	if closeParen <= openParen {
		return signature
	}
	funcName := signature[:openParen]
	paramStr := signature[openParen+1 : closeParen]
	if strings.TrimSpace(paramStr) == "" {
		return signature
	}

	// Parse the signature as a CREATE FUNCTION to extract param names from AST.
	// The signature comes from a synced function, so a parse failure keeps it
	// unchanged rather than failing the dump.
	sql := fmt.Sprintf("CREATE FUNCTION %s RETURNS void LANGUAGE sql AS 'SELECT 1'", signature)
	stmts, err := parseDefinition(sql)
	if err != nil {
		return signature
	}
	fn, ok := stmts[0].(*ast.CreateFunctionStmt)
	if !ok || fn.Parameters == nil {
		return signature
	}

	// Collect parameter names from AST.
	var paramNames []string
	for _, item := range fn.Parameters.Items {
		if fp, ok := item.(*ast.FunctionParameter); ok && fp.Name != "" {
			paramNames = append(paramNames, fp.Name)
		}
	}

	// Strip parameter names from the original text, preserving original type spelling.
	rawParams := strings.Split(paramStr, ",")
	var typesOnly []string
	nameIdx := 0
	for _, raw := range rawParams {
		param := strings.TrimSpace(raw)
		words := strings.Fields(param)
		if len(words) == 0 {
			continue
		}
		// Skip mode keywords (IN, OUT, INOUT, VARIADIC).
		startIdx := 0
		for startIdx < len(words) {
			w := strings.ToLower(words[startIdx])
			if w == "in" || w == "out" || w == "inout" || w == "variadic" {
				startIdx++
			} else {
				break
			}
		}
		// Check if first non-mode word matches a known param name from AST.
		if startIdx < len(words) && nameIdx < len(paramNames) && words[startIdx] == paramNames[nameIdx] {
			typesOnly = append(typesOnly, strings.Join(words[startIdx+1:], " "))
			nameIdx++
		} else {
			typesOnly = append(typesOnly, strings.Join(words[startIdx:], " "))
		}
	}

	return funcName + "(" + strings.Join(typesOnly, ", ") + ")"
}

// writeCollation writes a user-defined collation and its comment. The provider
// decides which locale fields PostgreSQL takes: ICU and builtin take a locale,
// libc takes the pair it sorts and classifies with.
// writeForeignDataWrapper writes CREATE FOREIGN DATA WRAPPER and the handler,
// validator and options it states. It is written before the servers that use it.
func writeForeignDataWrapper(out io.Writer, wrapper *metadatapb.ForeignDataWrapperMetadata) error {
	if _, err := fmt.Fprintf(out, "CREATE FOREIGN DATA WRAPPER \"%s\"", escapeIdentifier(wrapper.Name)); err != nil {
		return err
	}
	if wrapper.Handler != "" {
		if _, err := fmt.Fprintf(out, " HANDLER %s", wrapper.Handler); err != nil {
			return err
		}
	}
	if wrapper.Validator != "" {
		if _, err := fmt.Fprintf(out, " VALIDATOR %s", wrapper.Validator); err != nil {
			return err
		}
	}
	if err := writeForeignOptions(out, wrapper.Options); err != nil {
		return err
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

// It is written after extensions, which provide the wrapper, and before every
// foreign table that names it.
func writeForeignServer(out io.Writer, server *metadatapb.ForeignServerMetadata) error {
	if _, err := fmt.Fprintf(out, "CREATE SERVER \"%s\"", escapeIdentifier(server.Name)); err != nil {
		return err
	}
	if server.Type != nil {
		if _, err := fmt.Fprintf(out, " TYPE '%s'", escapeSingleQuote(*server.Type)); err != nil {
			return err
		}
	}
	if server.Version != nil {
		if _, err := fmt.Fprintf(out, " VERSION '%s'", escapeSingleQuote(*server.Version)); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(out, " FOREIGN DATA WRAPPER \"%s\"", escapeIdentifier(server.Fdw)); err != nil {
		return err
	}
	if err := writeForeignOptions(out, server.Options); err != nil {
		return err
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

// writeForeignTable writes CREATE FOREIGN TABLE with the columns the sync read,
// any check constraints a partition carries, and the options the table states
// itself. The server must already exist.
func writeForeignTable(out io.Writer, schema string, table *metadatapb.ExternalTableMetadata, checks []*metadatapb.CheckConstraintMetadata) error {
	if _, err := fmt.Fprintf(out, "CREATE FOREIGN TABLE \"%s\".\"%s\" (",
		escapeIdentifier(schema), escapeIdentifier(table.Name)); err != nil {
		return err
	}
	first := true
	writeSeparator := func() error {
		separator := ",\n    "
		if first {
			separator = "\n    "
			first = false
		}
		_, err := io.WriteString(out, separator)
		return err
	}
	for _, column := range table.Columns {
		if err := writeSeparator(); err != nil {
			return err
		}
		if err := writeColumnDefinition(out, column); err != nil {
			return err
		}
	}
	for _, check := range checks {
		if err := writeSeparator(); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "CONSTRAINT \"%s\" CHECK %s", escapeIdentifier(check.Name), check.Expression); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, "\n)"); err != nil {
		return err
	}
	if table.ExternalServerName != "" {
		if _, err := fmt.Fprintf(out, " SERVER \"%s\"", escapeIdentifier(table.ExternalServerName)); err != nil {
			return err
		}
	}
	if err := writeForeignOptions(out, table.Options); err != nil {
		return err
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

// writeForeignOptions writes an OPTIONS list from the entries pg_catalog stores
// as "key=value". A value is a string; a key that is not a simple identifier is
// quoted.
func writeForeignOptions(out io.Writer, options []string) error {
	if len(options) == 0 {
		return nil
	}
	if _, err := io.WriteString(out, " OPTIONS ("); err != nil {
		return err
	}
	for i, option := range options {
		if i != 0 {
			if _, err := io.WriteString(out, ", "); err != nil {
				return err
			}
		}
		key, value, ok := strings.Cut(option, "=")
		if !ok {
			key, value = option, ""
		}
		if err := writeForeignOptionName(out, key); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, " '%s'", escapeSingleQuote(value)); err != nil {
			return err
		}
	}
	_, err := io.WriteString(out, ")")
	return err
}

func writeForeignOptionName(out io.Writer, name string) error {
	// An unquoted PostgreSQL identifier folds to lower case, so an upper-case
	// character in an option name has to be quoted to keep its spelling.
	simple := name != ""
	for i, r := range name {
		if r == '_' || r >= 'a' && r <= 'z' || i > 0 && r >= '0' && r <= '9' {
			continue
		}
		simple = false
		break
	}
	if simple {
		_, err := io.WriteString(out, name)
		return err
	}
	_, err := fmt.Fprintf(out, "\"%s\"", escapeIdentifier(name))
	return err
}

func writeCollation(out io.Writer, schema string, collation *metadatapb.CollationMetadata) error {
	var provider, locale string
	switch collation.Provider {
	case "i":
		provider, locale = "icu", fmt.Sprintf("locale = '%s'", escapeSingleQuote(collation.Locale))
	case "c":
		provider = "libc"
		locale = fmt.Sprintf("lc_collate = '%s', lc_ctype = '%s'", escapeSingleQuote(collation.LcCollate), escapeSingleQuote(collation.LcCtype))
	case "b":
		provider, locale = "builtin", fmt.Sprintf("locale = '%s'", escapeSingleQuote(collation.Locale))
	default:
		// A collation with no provider to state is one the server derives, which
		// a restore cannot create.
		return nil
	}
	if _, err := fmt.Fprintf(out, "CREATE COLLATION \"%s\".\"%s\" (provider = %s, %s",
		escapeIdentifier(schema), escapeIdentifier(collation.Name), provider, locale); err != nil {
		return err
	}
	if collation.Rules != "" {
		if _, err := fmt.Fprintf(out, ", rules = '%s'", escapeSingleQuote(collation.Rules)); err != nil {
			return err
		}
	}
	if collation.Deterministic != nil {
		if _, err := fmt.Fprintf(out, ", deterministic = %t", *collation.Deterministic); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, ");\n\n"); err != nil {
		return err
	}
	if collation.Comment != "" {
		if _, err := fmt.Fprintf(out, "COMMENT ON COLLATION \"%s\".\"%s\" IS '%s';\n\n",
			escapeIdentifier(schema), escapeIdentifier(collation.Name), escapeSingleQuote(collation.Comment)); err != nil {
			return err
		}
	}
	return nil
}

// writeRangeType writes CREATE TYPE ... AS RANGE and its comment. The subtype
// is written as the server states it; the clauses it left out are the ones the
// server derives, so writing only what pg_range filled keeps the restore to
// what the source declared. canonical and subtype_diff are not written: a range
// that declares one is created from a shell type, then the function, then the
// full type, and that function's definition is a C symbol in a LANGUAGE
// internal function, which a dump cannot write. The metadata carries both, so a
// range that declares one shows the loss in the restored snapshot rather than
// restoring a type that compares differently.
func writeRangeType(out io.Writer, schema string, rangeType *metadatapb.RangeTypeMetadata) error {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TYPE \"%s\".\"%s\" AS RANGE (\n    subtype = %s",
		escapeIdentifier(schema), escapeIdentifier(rangeType.Name), rangeType.Subtype)
	for _, clause := range []struct{ key, value string }{
		{"subtype_opclass", rangeType.SubtypeOpclass},
		{"collation", rangeType.Collation},
		{"multirange_type_name", rangeType.MultirangeName},
	} {
		if clause.value == "" {
			continue
		}
		fmt.Fprintf(&b, ",\n    %s = %s", clause.key, clause.value)
	}
	b.WriteString("\n);\n\n")
	if _, err := io.WriteString(out, b.String()); err != nil {
		return err
	}
	if rangeType.Comment != "" {
		if _, err := fmt.Fprintf(out, "COMMENT ON TYPE \"%s\".\"%s\" IS '%s';\n\n",
			escapeIdentifier(schema), escapeIdentifier(rangeType.Name), escapeSingleQuote(rangeType.Comment)); err != nil {
			return err
		}
	}
	return nil
}

func writeEnumComment(out io.Writer, schema string, enum *metadatapb.EnumTypeMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON TYPE "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, enum.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(enum.Comment)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeMaterializedView(out io.Writer, schema string, view *metadatapb.MaterializedViewMetadata) error {
	if _, err := io.WriteString(out, `CREATE MATERIALIZED VIEW "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" AS \n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, strings.TrimRight(view.Definition, ";")); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n  WITH NO DATA;\n\n"); err != nil {
		return err
	}

	for _, index := range view.Indexes {
		if err := writeIndex(out, schema, view.Name, index, false); err != nil {
			return err
		}
	}

	if len(view.Comment) > 0 {
		if err := writeMaterializedViewComment(out, schema, view); err != nil {
			return err
		}
	}

	return nil
}

func writeMaterializedViewComment(out io.Writer, schema string, view *metadatapb.MaterializedViewMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON MATERIALIZED VIEW "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(view.Comment)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeView(out io.Writer, schema string, view *metadatapb.ViewMetadata) error {
	if _, err := io.WriteString(out, `CREATE VIEW "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" AS \n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Definition); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n\n"); err != nil {
		return err
	}

	if len(view.Comment) > 0 {
		if err := writeViewComment(out, schema, view); err != nil {
			return err
		}
	}

	return nil
}

func writeViewComment(out io.Writer, schema string, view *metadatapb.ViewMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON VIEW "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(view.Comment)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeRules(out io.Writer, _ string, _ string, rules []*metadatapb.RuleMetadata) error {
	for _, rule := range rules {
		// Write the full rule definition
		if _, err := io.WriteString(out, rule.Definition); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "\n\n"); err != nil {
			return err
		}
	}
	return nil
}

// objectIDSeparator joins schema + object into a map/graph key. NUL cannot
// legally appear in a PostgreSQL identifier, so this avoids ambiguity when
// either half itself contains a dot (e.g. a table literally named "a.b").
const objectIDSeparator = "\x00"

func getSchemaNameFromID(id string) string {
	if schema, _, ok := strings.Cut(id, objectIDSeparator); ok {
		return schema
	}
	return ""
}

func writeColumnIdentityGeneration(out io.Writer, schema string, generationTypes metadatapb.ColumnMetadata_IdentityGeneration, sequence *metadatapb.SequenceMetadata) error {
	if _, err := io.WriteString(out, `ALTER TABLE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.OwnerTable); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" ALTER COLUMN "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.OwnerColumn); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" ADD GENERATED `); err != nil {
		return err
	}
	if generationTypes == metadatapb.ColumnMetadata_ALWAYS {
		if _, err := io.WriteString(out, "ALWAYS "); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(out, "BY DEFAULT "); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, "AS IDENTITY (\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "    SEQUENCE NAME \""); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "    START WITH "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.Start); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n    INCREMENT BY "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.Increment); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n    MINVALUE "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.MinValue); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n    MAXVALUE "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.MaxValue); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n    CACHE "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.CacheSize); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\n);\n\n")
	return err
}

func writeCreateSequence(out io.Writer, schema string, sequence *metadatapb.SequenceMetadata) error {
	if _, err := io.WriteString(out, `CREATE SEQUENCE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\"\n    "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "AS "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.DataType); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n	START WITH "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.Start); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n	INCREMENT BY "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.Increment); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n	MINVALUE "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.MinValue); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\n	MAXVALUE "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.MaxValue); err != nil {
		return err
	}
	if sequence.CacheSize != "" {
		if _, err := io.WriteString(out, "\n	CACHE "); err != nil {
			return err
		}
		if _, err := io.WriteString(out, sequence.CacheSize); err != nil {
			return err
		}
	}
	if sequence.Cycle {
		if _, err := io.WriteString(out, "\n	CYCLE"); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(out, "\n	NO CYCLE"); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, ";\n\n"); err != nil {
		return err
	}

	if len(sequence.Comment) > 0 {
		if err := writeSequenceComment(out, schema, sequence); err != nil {
			return err
		}
	}

	return nil
}

func writeSequenceComment(out io.Writer, schema string, sequence *metadatapb.SequenceMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON SEQUENCE "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, sequence.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(sequence.Comment)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeAlterSequenceOwnedBy(out io.Writer, schema string, sequence *metadatapb.SequenceMetadata) error {
	if _, err := io.WriteString(out, `ALTER SEQUENCE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" OWNED BY \""); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.OwnerTable); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, sequence.OwnerColumn); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\";\n\n")
	return err
}

// metadataFuncIdentity returns the unique identity for a function, preferring
// the full signature (name + parameter types) to support overloaded functions.
func metadataFuncIdentity(f *metadatapb.FunctionMetadata) string {
	if f.Signature != "" {
		return f.Signature
	}
	return f.Name
}

// addColumnFunctionEdges adds an edge from each user function a column's
// default or generation expression calls to the column's table, so the
// function is created first. A STORED generated column's expression and a
// column default are validated against their functions at CREATE TABLE time,
// and TopologicalSort otherwise orders a table before a function it needs
// whenever the table's id sorts first. Built-in functions and unparsable
// expressions add no edge.
//
// The call's exact overload is not resolved, so an edge is added for every
// function sharing the referenced name — except an overload that itself
// depends on this table (for example one taking the table's row type), which
// the referenced overload never is and which would make a cycle that drops
// the sort to its all-tables-then-all-functions fallback.
func addColumnFunctionEdges(graph *definitionGraph, functionMap map[string]*metadatapb.FunctionMetadata, functionIDsByName map[string][]string, schema string, table *metadatapb.TableMetadata) {
	addColumnFunctionEdgesTo(graph, functionMap, functionIDsByName, schema, table.Name, getObjectID(schema, table.Name), table.Columns)
}

// addColumnFunctionEdgesTo adds the edges for columns of a relation named
// relation that is written as targetID. A foreign partition is written by its
// parent, so its overload dependencies are checked against the partition while
// the edge points at the parent.
func addColumnFunctionEdgesTo(graph *definitionGraph, functionMap map[string]*metadatapb.FunctionMetadata, functionIDsByName map[string][]string, schema, relation, targetID string, columns []*metadatapb.ColumnMetadata) {
	for _, column := range columns {
		var exprs []string
		if column.Default != "" {
			exprs = append(exprs, column.Default)
		}
		if column.Generation != nil && column.Generation.Expression != "" {
			exprs = append(exprs, column.Generation.Expression)
		}
		for _, expr := range exprs {
			node, err := parseScalar("SELECT " + expr)
			if err != nil {
				continue
			}
			funcRefs, _, _ := collectExprDeps(node)
			for _, ref := range funcRefs {
				for _, funcID := range functionIDsByName[ref] {
					if functionDependsOnTable(functionMap[funcID], schema, relation) {
						continue
					}
					graph.AddEdge(funcID, targetID)
				}
			}
		}
	}
}

// checkExpressionForDependencies drops the constraint suffixes that are not
// part of the expression, so the rest parses as a SELECT.
func checkExpressionForDependencies(expression string) string {
	for {
		trimmed := strings.TrimSpace(expression)
		upper := strings.ToUpper(trimmed)
		switch {
		case strings.HasSuffix(upper, " NO INHERIT"):
			expression = trimmed[:len(trimmed)-len(" NO INHERIT")]
		case strings.HasSuffix(upper, " NOT VALID"):
			expression = trimmed[:len(trimmed)-len(" NOT VALID")]
		default:
			return trimmed
		}
	}
}

// addRelationTypeEdges makes a relation follow the relation whose row type one
// of its columns uses. PostgreSQL row types are not composite types, so the
// type section does not create them; a foreign table can name another foreign
// table's row type, and the table that defines it must exist first.
func addRelationTypeEdges(graph *definitionGraph, foreignTableMap map[string]*metadatapb.ExternalTableMetadata, tableMap map[string]*metadatapb.TableMetadata, partitionParents map[string]string) {
	for id, table := range foreignTableMap {
		addRelationTypeEdgesFor(graph, id, table.Columns, table.CheckConstraints, tableMap, foreignTableMap, partitionParents)
	}
	for id, table := range tableMap {
		addRelationTypeEdgesFor(graph, id, table.Columns, table.CheckConstraints, tableMap, foreignTableMap, partitionParents)
	}
}

func addRelationTypeEdgesFor(graph *definitionGraph, id string, columns []*metadatapb.ColumnMetadata, checks []*metadatapb.CheckConstraintMetadata, tableMap map[string]*metadatapb.TableMetadata, foreignTableMap map[string]*metadatapb.ExternalTableMetadata, partitionParents map[string]string) {
	addTypeRef := func(typeName string) {
		schema, name, ok := strings.Cut(typeName, ".")
		if !ok {
			return
		}
		// An array of a row type is named _<relation>, as in public._z.
		names := []string{name}
		if strings.HasPrefix(name, "_") {
			names = append(names, name[1:])
		}
		for _, relationName := range names {
			dependencyID := getObjectID(schema, relationName)
			if parentID, ok := partitionParents[dependencyID]; ok {
				dependencyID = parentID
			}
			if dependencyID == id {
				return
			}
			if _, ok := tableMap[dependencyID]; ok {
				graph.AddEdge(dependencyID, id)
				return
			}
			if _, ok := foreignTableMap[dependencyID]; ok {
				graph.AddEdge(dependencyID, id)
				return
			}
		}
	}
	for _, column := range columns {
		addTypeRef(column.Type)
	}
	for _, check := range checks {
		node, err := parseScalar("SELECT " + checkExpressionForDependencies(check.Expression))
		if err != nil {
			continue
		}
		_, _, typeRefs := collectExprDeps(node)
		for _, typeRef := range typeRefs {
			addTypeRef(typeRef)
		}
	}
}

// addCheckFunctionEdges adds an edge from each user function a check expression
// calls to the relation that states the check, so the function exists first.
func addCheckFunctionEdges(graph *definitionGraph, functionMap map[string]*metadatapb.FunctionMetadata, functionIDsByName map[string][]string, schema, table string, checks []*metadatapb.CheckConstraintMetadata) {
	addCheckFunctionEdgesTo(graph, functionMap, functionIDsByName, schema, table, getObjectID(schema, table), checks)
}

func addCheckFunctionEdgesTo(graph *definitionGraph, functionMap map[string]*metadatapb.FunctionMetadata, functionIDsByName map[string][]string, schema, relation, targetID string, checks []*metadatapb.CheckConstraintMetadata) {
	for _, check := range checks {
		node, err := parseScalar("SELECT " + checkExpressionForDependencies(check.Expression))
		if err != nil {
			continue
		}
		funcRefs, _, _ := collectExprDeps(node)
		for _, ref := range funcRefs {
			for _, funcID := range functionIDsByName[ref] {
				if functionDependsOnTable(functionMap[funcID], schema, relation) {
					continue
				}
				graph.AddEdge(funcID, targetID)
			}
		}
	}
}

// functionDependsOnTable reports whether function lists (schema, table) among
// its dependency tables.
func functionDependsOnTable(function *metadatapb.FunctionMetadata, schema, table string) bool {
	if function == nil {
		return false
	}
	for _, dependency := range function.DependencyTables {
		if dependency.Schema == schema && dependency.Table == table {
			return true
		}
	}
	return false
}

// collectExternalTables indexes every foreign table by object ID.
func collectExternalTables(schemas []*metadatapb.SchemaMetadata) map[string]*metadatapb.ExternalTableMetadata {
	externalTables := make(map[string]*metadatapb.ExternalTableMetadata)
	for _, schema := range schemas {
		for _, externalTable := range schema.ExternalTables {
			externalTables[getObjectID(schema.Name, externalTable.Name)] = externalTable
		}
	}
	return externalTables
}

// collectForeignPartitionIDs names the foreign tables that a table's partition
// list also holds, so they are rendered with their parent rather than twice.
func collectForeignPartitionIDs(schemas []*metadatapb.SchemaMetadata, externalTables map[string]*metadatapb.ExternalTableMetadata) map[string]bool {
	ids := make(map[string]bool)
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			walkPartitions(schema.Name, table.Partitions, func(partitionSchema string, partition *metadatapb.TablePartitionMetadata) {
				id := getObjectID(partitionSchema, partition.Name)
				if externalTables[id] != nil {
					ids[id] = true
				}
			})
		}
	}
	return ids
}

// collectForeignPartitionParents maps a foreign partition's ID to the table
// that creates and attaches it, so a view that reads the partition waits for
// the parent rather than for a node of its own.
func collectForeignPartitionParents(schemas []*metadatapb.SchemaMetadata, externalTables map[string]*metadatapb.ExternalTableMetadata) map[string]string {
	parents := make(map[string]string)
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			parentID := getObjectID(schema.Name, table.Name)
			walkPartitions(schema.Name, table.Partitions, func(partitionSchema string, partition *metadatapb.TablePartitionMetadata) {
				id := getObjectID(partitionSchema, partition.Name)
				if externalTables[id] != nil {
					parents[id] = parentID
				}
			})
		}
	}
	return parents
}

func walkPartitions(schema string, partitions []*metadatapb.TablePartitionMetadata, visit func(string, *metadatapb.TablePartitionMetadata)) {
	for _, partition := range partitions {
		partitionSchema := partitionSchemaName(schema, partition)
		visit(partitionSchema, partition)
		walkPartitions(partitionSchema, partition.Subpartitions, visit)
	}
}

// addFunctionNameEdges makes every function a name resolves to precede node.
// regproc writes a name that needs quoting quoted, so the map's unquoted keys
// are tried too.
func addFunctionNameEdges(graph *definitionGraph, functionMap map[string]*metadatapb.FunctionMetadata, functionIDsByName map[string][]string, name string, expectedArgs []string, node string) {
	if name == "" {
		return
	}
	seen := make(map[string]bool)
	for _, key := range functionNameKeys(name) {
		for _, functionID := range functionIDsByName[key] {
			if seen[functionID] {
				continue
			}
			if !functionSignatureMatches(functionMap[functionID], expectedArgs) {
				continue
			}
			seen[functionID] = true
			graph.AddEdge(functionID, node)
		}
	}
}

// functionSignatureMatches reports whether a function's identity arguments are
// the ones a handler or validator has. A function map from a snapshot without
// a signature matches nothing, so only the function the catalog can name is
// ordered.
func functionSignatureMatches(function *metadatapb.FunctionMetadata, expected []string) bool {
	if function == nil || function.Signature == "" {
		return false
	}
	args, ok := signatureArguments(function.Signature, function.Name)
	if !ok {
		return false
	}
	actual := splitSignatureArguments(args)
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if signatureArgumentType(actual[i]) != normalizeSignatureType(expected[i]) {
			return false
		}
	}
	return true
}

// signatureArgumentType drops an argument's mode and name, leaving its type.
// pg_get_function_identity_arguments may keep both.
func signatureArgumentType(arg string) string {
	tokens := splitSignatureTokens(arg)
	for len(tokens) > 0 && isSignatureMode(tokens[0]) {
		tokens = tokens[1:]
	}
	if len(tokens) > 1 && !isMultiWordTypeStart(tokens[0]) {
		tokens = tokens[1:]
	}
	return normalizeSignatureType(strings.Join(tokens, " "))
}

func isSignatureMode(token string) bool {
	switch strings.ToUpper(token) {
	case "IN", "OUT", "INOUT", "VARIADIC":
		return true
	}
	return false
}

func isMultiWordTypeStart(token string) bool {
	switch strings.ToLower(token) {
	case "double", "character", "bit", "national", "timestamp", "time", "interval":
		return true
	}
	return false
}

func splitSignatureTokens(arg string) []string {
	var tokens []string
	var b strings.Builder
	inQuote := false
	for i := 0; i < len(arg); i++ {
		c := arg[i]
		if inQuote {
			b.WriteByte(c)
			if c != '"' {
				continue
			}
			if i+1 < len(arg) && arg[i+1] == '"' {
				b.WriteByte(arg[i+1])
				i++
				continue
			}
			inQuote = false
			continue
		}
		switch c {
		case '"':
			inQuote = true
			b.WriteByte(c)
		case ' ', '\t':
			if b.Len() > 0 {
				tokens = append(tokens, b.String())
				b.Reset()
			}
		default:
			b.WriteByte(c)
		}
	}
	if b.Len() > 0 {
		tokens = append(tokens, b.String())
	}
	return tokens
}

// signatureArguments returns the text between a signature's argument-list
// parentheses. The function name may itself contain a quoted parenthesis.
func signatureArguments(signature, name string) (string, bool) {
	// getFunctions builds the signature as name(arguments), so a name that
	// itself ends in a parenthesis still has its argument list right after it.
	if name != "" && strings.HasPrefix(signature, name) {
		rest := signature[len(name):]
		if strings.HasPrefix(rest, "(") && strings.HasSuffix(rest, ")") {
			return rest[1 : len(rest)-1], true
		}
	}
	inQuote := false
	for i := 0; i < len(signature); i++ {
		switch signature[i] {
		case '"':
			if inQuote && i+1 < len(signature) && signature[i+1] == '"' {
				i++
				continue
			}
			inQuote = !inQuote
		case '(':
			if inQuote {
				continue
			}
			close := strings.LastIndex(signature, ")")
			if close < i {
				return "", false
			}
			return signature[i+1 : close], true
		}
	}
	return "", false
}

func splitSignatureArguments(args string) []string {
	args = strings.TrimSpace(args)
	if args == "" {
		return nil
	}
	var parts []string
	depth := 0
	start := 0
	inQuote := false
	for i := 0; i < len(args); i++ {
		if inQuote {
			if args[i] != '"' {
				continue
			}
			if i+1 < len(args) && args[i+1] == '"' {
				i++
				continue
			}
			inQuote = false
			continue
		}
		switch args[i] {
		case '"':
			inQuote = true
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, args[start:i])
				start = i + 1
			}
		}
	}
	return append(parts, args[start:])
}

func normalizeSignatureType(t string) string {
	t = strings.TrimSpace(t)
	t = strings.ReplaceAll(t, "pg_catalog.", "")
	return strings.ReplaceAll(t, " ", "")
}

// functionNameKeys is a SQL-qualified name's unquoted forms: its qualified name
// and its bare name. regproc writes app."MyHandler" where the function map holds
// app.MyHandler.
func functionNameKeys(name string) []string {
	parts, ok := splitQualifiedIdentifier(name)
	if !ok {
		return []string{name}
	}
	if len(parts) == 1 {
		return parts
	}
	return []string{strings.Join(parts, ".")}
}

func splitQualifiedIdentifier(name string) ([]string, bool) {
	var parts []string
	var b strings.Builder
	inQuote := false
	for i := 0; i < len(name); i++ {
		c := name[i]
		if inQuote {
			if c != '"' {
				b.WriteByte(c)
				continue
			}
			if i+1 < len(name) && name[i+1] == '"' {
				b.WriteByte('"')
				i++
				continue
			}
			inQuote = false
			continue
		}
		switch c {
		case '"':
			inQuote = true
		case '.':
			parts = append(parts, b.String())
			b.Reset()
		case ' ', '\t':
		default:
			b.WriteByte(c)
		}
	}
	if inQuote {
		return nil, false
	}
	return append(parts, b.String()), true
}

// addForeignPartitionEdges makes the functions a foreign partition's columns
// call, and its server, precede the parent table that creates the partition.
func addForeignPartitionEdges(graph *definitionGraph, functionMap map[string]*metadatapb.FunctionMetadata, functionIDsByName map[string][]string, schema string, table *metadatapb.TableMetadata, externalTables map[string]*metadatapb.ExternalTableMetadata, serverIDsByName map[string]string) {
	parentID := getObjectID(schema, table.Name)
	walkPartitions(schema, table.Partitions, func(partitionSchema string, partition *metadatapb.TablePartitionMetadata) {
		externalTable := externalTables[getObjectID(partitionSchema, partition.Name)]
		if externalTable == nil {
			return
		}
		addColumnFunctionEdgesTo(graph, functionMap, functionIDsByName, partitionSchema, partition.Name, parentID, externalTable.Columns)
		addCheckFunctionEdgesTo(graph, functionMap, functionIDsByName, partitionSchema, partition.Name, parentID, externalTable.CheckConstraints)
		if serverID := serverIDsByName[externalTable.ExternalServerName]; serverID != "" {
			graph.AddEdge(serverID, parentID)
		}
	})
}

func getObjectID(schema string, object string) string {
	var buf strings.Builder
	_, _ = buf.WriteString(schema)
	_, _ = buf.WriteString(objectIDSeparator)
	_, _ = buf.WriteString(object)
	return buf.String()
}

// writeColumnDefinition writes one column of a relation, after the separator
// that begins its line: its name, type, collation, generation or default, and
// nullability.
func writeColumnDefinition(out io.Writer, column *metadatapb.ColumnMetadata) error {
	if _, err := io.WriteString(out, `"`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeIdentifier(column.Name)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" `); err != nil {
		return err
	}

	if _, err := io.WriteString(out, column.Type); err != nil {
		return err
	}

	// A foreign table's column options map its name for the wrapper.
	if err := writeForeignOptions(out, column.Options); err != nil {
		return err
	}

	// A user-defined collation is written with its schema: the dump restores
	// under an empty search_path, so its bare name reaches only the collations
	// pg_catalog holds.
	if column.Collation != "" {
		if column.CollationSchema != "" {
			if _, err := fmt.Fprintf(out, ` COLLATE "%s"."%s"`, escapeIdentifier(column.CollationSchema), escapeIdentifier(column.Collation)); err != nil {
				return err
			}
		} else if _, err := fmt.Fprintf(out, ` COLLATE "%s"`, escapeIdentifier(column.Collation)); err != nil {
			return err
		}
	}

	// A generated column carries its expression in place of a default.
	if column.Generation != nil {
		storage := "STORED"
		if column.Generation.Type == metadatapb.GenerationMetadata_TYPE_VIRTUAL {
			storage = "VIRTUAL"
		}
		if _, err := fmt.Fprintf(out, " GENERATED ALWAYS AS (%s) %s", column.Generation.Expression, storage); err != nil {
			return err
		}
	} else if column.Default != "" {
		if _, err := io.WriteString(out, ` DEFAULT `); err != nil {
			return err
		}
		if _, err := io.WriteString(out, column.Default); err != nil {
			return err
		}
	}

	if !column.Nullable {
		if _, err := io.WriteString(out, ` NOT NULL`); err != nil {
			return err
		}
	}
	return nil
}

func writeCreateTable(out io.Writer, schema string, tableName string, columns []*metadatapb.ColumnMetadata, checks []*metadatapb.CheckConstraintMetadata, excludes []*metadatapb.ExcludeConstraintMetadata) error {
	if _, err := io.WriteString(out, `CREATE TABLE "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, tableName); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" (`); err != nil {
		return err
	}

	// first tracks whether we've written any element inside the parens yet so
	// the leading comma on constraints is suppressed when there are no columns
	// (e.g. tables that only carry CHECK constraints).
	first := true
	writeItemSep := func() error {
		if first {
			first = false
			_, err := io.WriteString(out, "\n    ")
			return err
		}
		_, err := io.WriteString(out, ",\n    ")
		return err
	}

	for _, column := range columns {
		if err := writeItemSep(); err != nil {
			return err
		}
		if err := writeColumnDefinition(out, column); err != nil {
			return err
		}
	}

	for _, check := range checks {
		if err := writeItemSep(); err != nil {
			return err
		}
		_, _ = io.WriteString(out, `CONSTRAINT "`)
		_, _ = io.WriteString(out, check.Name)
		_, _ = io.WriteString(out, `" CHECK `)
		_, _ = io.WriteString(out, check.Expression)
	}

	for _, exclude := range excludes {
		if err := writeItemSep(); err != nil {
			return err
		}
		_, _ = io.WriteString(out, `CONSTRAINT "`)
		_, _ = io.WriteString(out, exclude.Name)
		_, _ = io.WriteString(out, `" `)
		_, _ = io.WriteString(out, exclude.Expression) // Already includes "EXCLUDE USING ..."
	}

	_, err := io.WriteString(out, "\n)")
	return err
}

// isIdentityColumn checks if a column is an identity column.
func isIdentityColumn(column *metadatapb.ColumnMetadata) bool {
	return column.IdentityGeneration == metadatapb.ColumnMetadata_ALWAYS ||
		column.IdentityGeneration == metadatapb.ColumnMetadata_BY_DEFAULT
}

// buildIdentityColumnMap builds a map of identity columns for the entire database.
// Key format: "schemaName.tableName.columnName" -> true
// This is O(m*c) which is done once, vs O(s*m*c) when checking each sequence individually.
func buildIdentityColumnMap(metadata *metadatapb.DatabaseSchemaMetadata) map[string]bool {
	identityMap := make(map[string]bool)
	for _, schema := range metadata.Schemas {
		for _, table := range schema.Tables {
			for _, column := range table.Columns {
				if isIdentityColumn(column) {
					key := schema.Name + "." + table.Name + "." + column.Name
					identityMap[key] = true
				}
			}
		}
	}
	return identityMap
}

// buildIdentityColumnMapForSchema builds a map of identity columns for a single schema.
// Key format: "tableName.columnName" -> true
func buildIdentityColumnMapForSchema(schema *metadatapb.SchemaMetadata) map[string]bool {
	identityMap := make(map[string]bool)
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			if isIdentityColumn(column) {
				key := table.Name + "." + column.Name
				identityMap[key] = true
			}
		}
	}
	return identityMap
}

func splitSequencesByIdentityOrNot(table *metadatapb.TableMetadata, sequences []*metadatapb.SequenceMetadata) ([]metadatapb.ColumnMetadata_IdentityGeneration, []*metadatapb.SequenceMetadata, []*metadatapb.SequenceMetadata) {
	columnMap := make(map[string]*metadatapb.ColumnMetadata)
	for _, column := range table.Columns {
		columnMap[column.Name] = column
	}
	var generationType []metadatapb.ColumnMetadata_IdentityGeneration
	var identitySequences []*metadatapb.SequenceMetadata
	var nonIdentitySequences []*metadatapb.SequenceMetadata
	for _, sequence := range sequences {
		if column, ok := columnMap[sequence.OwnerColumn]; ok {
			if isIdentityColumn(column) {
				generationType = append(generationType, column.IdentityGeneration)
				identitySequences = append(identitySequences, sequence)
				continue
			}
		}
		nonIdentitySequences = append(nonIdentitySequences, sequence)
	}
	return generationType, identitySequences, nonIdentitySequences
}

// tableMissingColumnMetadata reports whether a synced table's column list
// was filtered away. information_schema.columns is privilege-filtered, so a
// sync user without column privileges on a table sees zero columns while the
// pg_catalog-based queries still return the table's constraints, indexes,
// triggers, and rules. DDL referencing the missing columns would make the
// dump non-replayable, so writers skip column-dependent statements for such
// tables.
//
// A genuine zero-column table (CREATE TABLE t ()) is legal in PostgreSQL and
// may carry column-independent objects such as CHECK (1 = 1), but it cannot
// carry primary keys, unique constraints, foreign keys, or a partition key
// (PostgreSQL rejects constant partition-key expressions); indexes on it are
// limited to constant-expression corner cases. Treat zero columns as missing
// metadata only when such objects are present, so genuine zero-column tables
// keep their objects.
func tableMissingColumnMetadata(table *metadatapb.TableMetadata) bool {
	return len(table.Columns) == 0 &&
		(len(table.Indexes) > 0 || len(table.ForeignKeys) > 0 || len(table.Partitions) > 0)
}

// tableHasNoColumns reports an empty synced column list, whether genuinely
// zero-column or privilege-filtered. A sequence cannot be owned by a column
// of such a table, so ALTER SEQUENCE ... OWNED BY is skipped either way.
func tableHasNoColumns(table *metadatapb.TableMetadata) bool {
	return len(table.Columns) == 0
}

// collectTables returns the tables matching include, keyed by
// getObjectID(schema, table).
func collectTables(schemas []*metadatapb.SchemaMetadata, include func(*metadatapb.TableMetadata) bool) map[string]bool {
	set := make(map[string]bool)
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			if include(table) {
				set[getObjectID(schema.Name, table.Name)] = true
			}
		}
	}
	return set
}

// writeSchemaTriggers writes the triggers of every table, view, and
// materialized view in the schema. Triggers on tables without column
// metadata are skipped; their definitions (UPDATE OF, WHEN clauses) can
// reference the missing columns.
func writeSchemaTriggers(out io.Writer, schema *metadatapb.SchemaMetadata, includeForeignTables bool) error {
	for _, table := range schema.Tables {
		if tableMissingColumnMetadata(table) {
			continue
		}
		for _, trigger := range table.Triggers {
			if trigger.SkipDump {
				continue
			}
			if err := writeTrigger(out, schema.Name, table.Name, trigger); err != nil {
				return err
			}
		}
	}

	for _, view := range schema.Views {
		for _, trigger := range view.Triggers {
			if trigger.SkipDump {
				continue
			}
			if err := writeTrigger(out, schema.Name, view.Name, trigger); err != nil {
				return err
			}
		}
	}

	for _, materializedView := range schema.MaterializedViews {
		for _, trigger := range materializedView.Triggers {
			if trigger.SkipDump {
				continue
			}
			if err := writeTrigger(out, schema.Name, materializedView.Name, trigger); err != nil {
				return err
			}
		}
	}

	if includeForeignTables {
		for _, externalTable := range schema.ExternalTables {
			if externalTable.SkipDump {
				continue
			}
			for _, trigger := range externalTable.Triggers {
				if trigger.SkipDump {
					continue
				}
				if err := writeTrigger(out, schema.Name, externalTable.Name, trigger); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// writeSchemaRules writes the rules of every table and the non-SELECT rules
// of every view in the schema (SELECT rules are part of the view
// definition). Rules on tables without column metadata are skipped; their
// definitions can reference the missing columns.
func writeSchemaRules(out io.Writer, schema *metadatapb.SchemaMetadata) error {
	for _, table := range schema.Tables {
		if len(table.Rules) > 0 && !tableMissingColumnMetadata(table) {
			if err := writeRules(out, schema.Name, table.Name, table.Rules); err != nil {
				return err
			}
		}
	}

	for _, view := range schema.Views {
		var nonSelectRules []*metadatapb.RuleMetadata
		for _, rule := range view.Rules {
			if rule.Event != "SELECT" {
				nonSelectRules = append(nonSelectRules, rule)
			}
		}
		if len(nonSelectRules) > 0 {
			if err := writeRules(out, schema.Name, view.Name, nonSelectRules); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeSchemaForeignKeys writes the foreign keys of every table in the
// schema, skipping foreign keys touching a table without column metadata on
// either side; they would reference columns absent from the dump.
func writeSchemaForeignKeys(out io.Writer, schema *metadatapb.SchemaMetadata, tablesMissingColumns map[string]bool) error {
	for _, table := range schema.Tables {
		if table.SkipDump || tablesMissingColumns[getObjectID(schema.Name, table.Name)] {
			continue
		}
		for _, fk := range table.ForeignKeys {
			if tablesMissingColumns[getObjectID(fk.ReferencedSchema, fk.ReferencedTable)] {
				continue
			}
			if err := writeForeignKey(out, schema.Name, table.Name, fk); err != nil {
				return err
			}
			if len(fk.Comment) > 0 {
				if err := writeConstraintComment(out, schema.Name, table.Name, fk.Name, fk.Comment); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func writeTable(out io.Writer, schema string, table *metadatapb.TableMetadata, sequences []*metadatapb.SequenceMetadata, externalTables map[string]*metadatapb.ExternalTableMetadata) error {
	if tableMissingColumnMetadata(table) {
		// Emit a bare CREATE TABLE so the table still exists on replay, and
		// skip constraints, indexes, partitions, and sequence ownership —
		// they all reference columns absent from the dump.
		if err := writeCreateTable(out, schema, table.Name, nil, nil, nil); err != nil {
			return err
		}
		if _, err := io.WriteString(out, ";\n\n"); err != nil {
			return err
		}
		if len(table.Comment) > 0 {
			return writeTableComment(out, schema, table)
		}
		return nil
	}

	generationTypes, identitySequences, nonIdentitySequences := splitSequencesByIdentityOrNot(table, sequences)

	if err := writeCreateTable(out, schema, table.Name, table.Columns, table.CheckConstraints, table.ExcludeConstraints); err != nil {
		return err
	}

	if len(table.Partitions) > 0 {
		if err := writePartitionClause(out, table.Partitions[0]); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, ";\n\n"); err != nil {
		return err
	}

	// Skip OWNED BY for a table with no columns; the clause would reference
	// a column absent from the dump.
	if !tableHasNoColumns(table) {
		for _, sequence := range nonIdentitySequences {
			if err := writeAlterSequenceOwnedBy(out, schema, sequence); err != nil {
				return err
			}
		}
	}

	// Set the owner after linking OWNED BY sequences, so the change cascades
	// to them: PostgreSQL requires an owned sequence to share its table's
	// owner. Identity sequences are created below, under this owner.
	if table.Owner != "" {
		if _, err := fmt.Fprintf(out, "ALTER TABLE \"%s\".\"%s\" OWNER TO \"%s\";\n\n", escapeIdentifier(schema), escapeIdentifier(table.Name), escapeIdentifier(table.Owner)); err != nil {
			return err
		}
	}

	// Construct comments.
	if len(table.Comment) > 0 {
		if err := writeTableComment(out, schema, table); err != nil {
			return err
		}
	}

	for _, column := range table.Columns {
		if len(column.Comment) > 0 {
			if err := writeColumnComment(out, schema, table.Name, column); err != nil {
				return err
			}
		}
	}

	for _, check := range table.CheckConstraints {
		if len(check.Comment) > 0 {
			if err := writeConstraintComment(out, schema, table.Name, check.Name, check.Comment); err != nil {
				return err
			}
		}
	}

	for i, sequence := range identitySequences {
		if err := writeColumnIdentityGeneration(out, schema, generationTypes[i], sequence); err != nil {
			return err
		}
	}

	// Construct partition tables. A foreign partition is created as a foreign
	// table, and the attach below attaches it to its parent.
	for _, partition := range table.Partitions {
		partitionSchema := partitionSchemaName(schema, partition)
		if externalTable, ok := externalTables[getObjectID(partitionSchema, partition.Name)]; ok {
			if err := writeForeignTable(out, partitionSchema, externalTable, partition.CheckConstraints); err != nil {
				return err
			}
			if err := writeForeignPartitionComments(out, partitionSchema, partition); err != nil {
				return err
			}
			continue
		}
		if err := writePartitionTable(out, schema, table.Columns, partition, externalTables); err != nil {
			return err
		}
	}

	for _, partition := range table.Partitions {
		if err := writeAttachPartition(out, schema, table.Name, partition); err != nil {
			return err
		}
	}

	// Construct Primary Key.
	for _, index := range table.Indexes {
		if index.Primary {
			if err := writePrimaryKey(out, schema, table.Name, index); err != nil {
				return err
			}
		}
	}

	// Construct Partition table primary key.
	for _, partition := range table.Partitions {
		if err := writePartitionPrimaryKey(out, schema, partition); err != nil {
			return err
		}
	}

	// Construct Unique Key.
	for _, index := range table.Indexes {
		if index.Unique && !index.Primary && index.IsConstraint {
			if err := writeUniqueKey(out, schema, table.Name, index); err != nil {
				return err
			}
		}
	}

	// Construct Partition table unique key.
	for _, partition := range table.Partitions {
		if err := writePartitionUniqueKey(out, schema, partition); err != nil {
			return err
		}
	}

	// Construct Index.
	for _, index := range table.Indexes {
		if !index.Primary && !index.IsConstraint {
			if err := writeIndex(out, schema, table.Name, index, len(table.Partitions) > 0); err != nil {
				return err
			}
		}
	}

	// Construct Partition table index.
	for _, partition := range table.Partitions {
		if err := writePartitionIndex(out, schema, partition); err != nil {
			return err
		}
	}

	// Construct index attach partition.
	for _, partition := range table.Partitions {
		if err := writeAttachPartitionIndex(out, schema, partition); err != nil {
			return err
		}
	}

	return nil
}

func writeForeignKey(out io.Writer, schema string, table string, fk *metadatapb.ForeignKeyMetadata) error {
	if _, err := io.WriteString(out, `ALTER TABLE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\"\n    ADD CONSTRAINT \""); err != nil {
		return err
	}
	if _, err := io.WriteString(out, fk.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" FOREIGN KEY ("); err != nil {
		return err
	}
	for i, column := range fk.Columns {
		if i > 0 {
			if _, err := io.WriteString(out, ", "); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, column); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, ")\n    REFERENCES \""); err != nil {
		return err
	}
	if _, err := io.WriteString(out, fk.ReferencedSchema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, fk.ReferencedTable); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" ("); err != nil {
		return err
	}
	for i, column := range fk.ReferencedColumns {
		if i > 0 {
			if _, err := io.WriteString(out, ", "); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, column); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, ")"); err != nil {
		return err
	}

	// MATCH SIMPLE is the default and PostgreSQL omits it; FULL and PARTIAL
	// are written.
	switch fk.MatchType {
	case "FULL":
		if _, err := io.WriteString(out, "\n    MATCH FULL"); err != nil {
			return err
		}
	case "PARTIAL":
		if _, err := io.WriteString(out, "\n    MATCH PARTIAL"); err != nil {
			return err
		}
	}

	// Add ON DELETE clause if specified
	if fk.OnDelete != "" && fk.OnDelete != "NO ACTION" {
		if _, err := io.WriteString(out, "\n    ON DELETE "); err != nil {
			return err
		}
		if _, err := io.WriteString(out, fk.OnDelete); err != nil {
			return err
		}
	}

	// Add ON UPDATE clause if specified
	if fk.OnUpdate != "" && fk.OnUpdate != "NO ACTION" {
		if _, err := io.WriteString(out, "\n    ON UPDATE "); err != nil {
			return err
		}
		if _, err := io.WriteString(out, fk.OnUpdate); err != nil {
			return err
		}
	}

	_, err := io.WriteString(out, ";\n\n")
	return err
}

func writeAttachPartitionIndex(out io.Writer, parentSchema string, partition *metadatapb.TablePartitionMetadata) error {
	schema := partitionSchemaName(parentSchema, partition)
	for _, index := range partition.Indexes {
		if err := writeAttachIndex(out, schema, index); err != nil {
			return err
		}
	}

	for _, subpartition := range partition.Subpartitions {
		if err := writeAttachPartitionIndex(out, schema, subpartition); err != nil {
			return err
		}
	}
	return nil
}

func writeAttachIndex(out io.Writer, schema string, index *metadatapb.IndexMetadata) error {
	if len(index.ParentIndexName) == 0 || len(index.ParentIndexSchema) == 0 {
		return nil
	}

	if _, err := io.WriteString(out, `ALTER INDEX "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, index.ParentIndexSchema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, index.ParentIndexName); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" ATTACH PARTITION "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, index.Name); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\";\n\n")
	return err
}

func writePartitionIndex(out io.Writer, parentSchema string, partition *metadatapb.TablePartitionMetadata) error {
	schema := partitionSchemaName(parentSchema, partition)
	for _, index := range partition.Indexes {
		if !index.IsConstraint && !index.Primary {
			if err := writeIndex(out, schema, partition.Name, index, len(partition.Subpartitions) > 0); err != nil {
				return err
			}
		}
	}

	for _, subpartition := range partition.Subpartitions {
		if err := writePartitionIndex(out, schema, subpartition); err != nil {
			return err
		}
	}
	return nil
}

func writeIndex(out io.Writer, schema string, table string, index *metadatapb.IndexMetadata, useOnlyClause bool) error {
	return writeIndexInternal(out, schema, table, index, useOnlyClause, true)
}

func writeIndexKeyList(out io.Writer, index *metadatapb.IndexMetadata) error {
	if _, err := io.WriteString(out, `(`); err != nil {
		return err
	}

	// NULLS ordering is not stored on IndexMetadata; it is read back from the
	// index definition, which pg_get_indexdef writes only when non-default.
	nulls := indexKeyNulls(index.Definition)

	for i, expression := range index.Expressions {
		if i > 0 {
			if _, err := io.WriteString(out, ", "); err != nil {
				return err
			}
		}

		// expression is expected to be in canonical pg_get_indexdef form —
		// see IndexMetadata.expressions. Callers repair historical rows that
		// pre-date the contract before rendering.
		if _, err := io.WriteString(out, expression); err != nil {
			return err
		}

		// Add opclass if it's not the default
		if i < len(index.OpclassNames) && i < len(index.OpclassDefaults) {
			if !index.OpclassDefaults[i] && index.OpclassNames[i] != "" {
				if _, err := io.WriteString(out, " "); err != nil {
					return err
				}
				if _, err := io.WriteString(out, index.OpclassNames[i]); err != nil {
					return err
				}
			}
		}

		// Add DESC if this column is marked as descending
		if i < len(index.Descending) && index.Descending[i] {
			if _, err := io.WriteString(out, " DESC"); err != nil {
				return err
			}
		}

		if i < len(nulls) {
			switch nulls[i] {
			case ast.SORTBY_NULLS_FIRST:
				if _, err := io.WriteString(out, " NULLS FIRST"); err != nil {
					return err
				}
			case ast.SORTBY_NULLS_LAST:
				if _, err := io.WriteString(out, " NULLS LAST"); err != nil {
					return err
				}
			}
		}
	}

	_, err := io.WriteString(out, ")")
	return err
}

// indexKeyNulls returns the per-key NULLS ordering of a CREATE INDEX
// definition, aligned with its key expressions. A definition that does not
// parse yields no orderings, so the key list renders without them.
func indexKeyNulls(definition string) []ast.SortByNulls {
	if definition == "" {
		return nil
	}
	stmts, err := parseDefinition(definition)
	if err != nil || len(stmts) == 0 {
		return nil
	}
	idx, ok := stmts[0].(*ast.IndexStmt)
	if !ok || idx == nil || idx.IndexParams == nil {
		return nil
	}
	var out []ast.SortByNulls
	for _, item := range idx.IndexParams.Items {
		if elem, ok := item.(*ast.IndexElem); ok {
			out = append(out, elem.NullsOrdering)
		}
	}
	return out
}

func writeIndexComment(out io.Writer, schema string, index *metadatapb.IndexMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON INDEX "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, index.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(index.Comment)); err != nil {
		return err
	}

	_, err := io.WriteString(out, "';\n\n")
	return err
}

func writePartitionUniqueKey(out io.Writer, parentSchema string, partition *metadatapb.TablePartitionMetadata) error {
	schema := partitionSchemaName(parentSchema, partition)
	for _, index := range partition.Indexes {
		if index.Unique && !index.Primary && index.IsConstraint {
			if err := writeUniqueKey(out, schema, partition.Name, index); err != nil {
				return err
			}
		}
	}

	for _, subpartition := range partition.Subpartitions {
		if err := writePartitionUniqueKey(out, schema, subpartition); err != nil {
			return err
		}
	}
	return nil
}

func writeUniqueKey(out io.Writer, schema string, table string, index *metadatapb.IndexMetadata) error {
	if _, err := io.WriteString(out, `ALTER TABLE ONLY "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" ADD CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, index.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" UNIQUE (`); err != nil {
		return err
	}
	for i, expression := range index.Expressions {
		if i > 0 {
			if _, err := io.WriteString(out, ", "); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, expression); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, ");\n\n"); err != nil {
		return err
	}

	if len(index.Comment) > 0 {
		if err := writeConstraintComment(out, schema, table, index.Name, index.Comment); err != nil {
			return err
		}
	}

	return nil
}

func writePartitionPrimaryKey(out io.Writer, parentSchema string, partition *metadatapb.TablePartitionMetadata) error {
	schema := partitionSchemaName(parentSchema, partition)
	for _, index := range partition.Indexes {
		if index.Primary {
			if err := writePrimaryKey(out, schema, partition.Name, index); err != nil {
				return err
			}
		}
	}

	for _, subpartition := range partition.Subpartitions {
		if err := writePartitionPrimaryKey(out, schema, subpartition); err != nil {
			return err
		}
	}
	return nil
}

func writePrimaryKey(out io.Writer, schema string, table string, index *metadatapb.IndexMetadata) error {
	if _, err := io.WriteString(out, `ALTER TABLE ONLY "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" ADD CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, index.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" PRIMARY KEY (`); err != nil {
		return err
	}
	for i, expression := range index.Expressions {
		if i > 0 {
			if _, err := io.WriteString(out, ", "); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, expression); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, ");\n\n"); err != nil {
		return err
	}

	if len(index.Comment) > 0 {
		if err := writeConstraintComment(out, schema, table, index.Name, index.Comment); err != nil {
			return err
		}
	}

	return nil
}

// writeConstraintComment writes COMMENT ON CONSTRAINT for a constraint the
// engine states a comment beside: a key constraint (carried as an index), a
// check constraint or a foreign key.
func writeConstraintComment(out io.Writer, schema string, table string, name string, comment string) error {
	if _, err := io.WriteString(out, `COMMENT ON CONSTRAINT "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" ON "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, table); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(comment)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeColumnComment(out io.Writer, schema string, table string, column *metadatapb.ColumnMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON COLUMN "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, column.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, escapeSingleQuote(column.Comment)); err != nil {
		return err
	}
	_, err := io.WriteString(out, "';\n\n")
	return err
}

func writeTableComment(out io.Writer, schema string, table *metadatapb.TableMetadata) error {
	return writeRelationComment(out, schema, table.Name, table.Comment)
}

// writeForeignPartitionComments writes the comments a foreign partition and its
// check constraints carry, which the ordinary partition path writes after its
// CREATE TABLE.
func writeForeignPartitionComments(out io.Writer, schema string, partition *metadatapb.TablePartitionMetadata) error {
	for _, check := range partition.CheckConstraints {
		if len(check.Comment) > 0 {
			if err := writeConstraintComment(out, schema, partition.Name, check.Name, check.Comment); err != nil {
				return err
			}
		}
	}
	if len(partition.Comment) > 0 {
		return writeForeignPartitionComment(out, schema, partition)
	}
	return nil
}

func writeForeignPartitionComment(out io.Writer, schema string, partition *metadatapb.TablePartitionMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON FOREIGN TABLE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, partition.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, escapeSingleQuote(partition.Comment)); err != nil {
		return err
	}
	_, err := io.WriteString(out, "';\n\n")
	return err
}

func writePartitionComment(out io.Writer, schema string, partition *metadatapb.TablePartitionMetadata) error {
	return writeRelationComment(out, schema, partition.Name, partition.Comment)
}

func writeRelationComment(out io.Writer, schema string, table string, text string) error {
	if _, err := io.WriteString(out, `COMMENT ON TABLE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, escapeSingleQuote(text)); err != nil {
		return err
	}
	_, err := io.WriteString(out, "';\n\n")
	return err
}

func writePartitionClause(out io.Writer, partition *metadatapb.TablePartitionMetadata) error {
	if _, err := io.WriteString(out, " PARTITION BY "); err != nil {
		return err
	}
	_, err := io.WriteString(out, partition.Expression)
	return err
}

func writeAttachPartition(out io.Writer, parentSchema string, tableName string, partition *metadatapb.TablePartitionMetadata) error {
	if _, err := io.WriteString(out, `ALTER TABLE ONLY "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, parentSchema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tableName); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" ATTACH PARTITION "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, partitionSchemaName(parentSchema, partition)); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, partition.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\" "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, partition.Value); err != nil {
		return err
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

func writePartitionTable(out io.Writer, parentSchema string, columns []*metadatapb.ColumnMetadata, partition *metadatapb.TablePartitionMetadata, externalTables map[string]*metadatapb.ExternalTableMetadata) error {
	schema := partitionSchemaName(parentSchema, partition)
	if err := writeCreateTable(out, schema, partition.Name, columns, partition.CheckConstraints, partition.ExcludeConstraints); err != nil {
		return err
	}

	if len(partition.Subpartitions) > 0 {
		if err := writePartitionClause(out, partition.Subpartitions[0]); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, ";\n\n"); err != nil {
		return err
	}

	// A partition's check constraint comment is a statement of its own: written
	// into the CREATE TABLE before its optional PARTITION BY clause, it would
	// not parse.
	for _, check := range partition.CheckConstraints {
		if len(check.Comment) > 0 {
			if err := writeConstraintComment(out, schema, partition.Name, check.Name, check.Comment); err != nil {
				return err
			}
		}
	}

	if len(partition.Comment) > 0 {
		if err := writePartitionComment(out, schema, partition); err != nil {
			return err
		}
	}

	// Construct subpartition tables. A foreign subpartition is created as a
	// foreign table, and the attach below attaches it.
	for _, subpartition := range partition.Subpartitions {
		subpartitionSchema := partitionSchemaName(schema, subpartition)
		if externalTable, ok := externalTables[getObjectID(subpartitionSchema, subpartition.Name)]; ok {
			if err := writeForeignTable(out, subpartitionSchema, externalTable, subpartition.CheckConstraints); err != nil {
				return err
			}
			if err := writeForeignPartitionComments(out, subpartitionSchema, subpartition); err != nil {
				return err
			}
			continue
		}
		if err := writePartitionTable(out, schema, columns, subpartition, externalTables); err != nil {
			return err
		}
	}

	for _, subpartition := range partition.Subpartitions {
		if err := writeAttachPartition(out, schema, partition.Name, subpartition); err != nil {
			return err
		}
	}

	return nil
}

func writeFunction(out io.Writer, schema string, function *metadatapb.FunctionMetadata) error {
	if _, err := io.WriteString(out, function.Definition); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ";\n\n"); err != nil {
		return err
	}

	if len(function.Comment) > 0 {
		if err := writeFunctionComment(out, schema, function); err != nil {
			return err
		}
	}

	return nil
}

// isDefinitionProcedure checks if the definition string represents a PROCEDURE (not a FUNCTION).
// Uses omni parser: procedures have a DefElem{Defname: "isProcedure"} in Options.
func isDefinitionProcedure(definition string) bool {
	if definition == "" {
		return false
	}
	// The definition comes from a synced function; one that does not parse is
	// treated as a function.
	stmts, err := parseDefinition(definition)
	if err != nil || len(stmts) == 0 {
		return false
	}
	fn, ok := stmts[0].(*ast.CreateFunctionStmt)
	if !ok || fn == nil || fn.Options == nil {
		return false
	}
	for _, item := range fn.Options.Items {
		if def, ok := item.(*ast.DefElem); ok && def.Defname == "isProcedure" {
			return true
		}
	}
	return false
}

func writeFunctionComment(out io.Writer, schema string, function *metadatapb.FunctionMetadata) error {
	// Determine if this is a PROCEDURE or FUNCTION by checking the definition
	objectType := "FUNCTION"
	if isDefinitionProcedure(function.Definition) {
		objectType = "PROCEDURE"
	}

	if _, err := io.WriteString(out, "COMMENT ON "+objectType+" \""); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `".`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, signatureTypesOnly(function.Signature)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ` IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(function.Comment)); err != nil {
		return err
	}

	_, err := io.WriteString(out, "';\n\n")
	return err
}

func writeExtension(out io.Writer, extension *metadatapb.ExtensionMetadata) error {
	if _, err := io.WriteString(out, `CREATE EXTENSION IF NOT EXISTS "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, extension.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"`); err != nil {
		return err
	}

	// Add WITH SCHEMA clause if schema is specified and not empty
	if extension.Schema != "" {
		if _, err := io.WriteString(out, ` WITH SCHEMA "`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, extension.Schema); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
	}

	// Add VERSION clause if version is specified
	if extension.Version != "" {
		if _, err := io.WriteString(out, ` VERSION '`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, extension.Version); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `'`); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, `;`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, "\n"); err != nil {
		return err
	}

	// Write extension description (comment) if present
	if extension.Description != "" {
		if _, err := io.WriteString(out, "\n"); err != nil {
			return err
		}
		if err := writeExtensionComment(out, extension); err != nil {
			return err
		}
	} else {
		// Add blank line even if no description
		if _, err := io.WriteString(out, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func writeExtensionComment(out io.Writer, extension *metadatapb.ExtensionMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON EXTENSION "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, extension.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, escapeSingleQuote(extension.Description)); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

// writeIndexInternal is the core index writing function with options for different modes
func writeIndexInternal(out io.Writer, schema string, table string, index *metadatapb.IndexMetadata, useOnlyClause bool, includeTerminatorAndComment bool) error {
	if index.Unique {
		if _, err := io.WriteString(out, `CREATE UNIQUE INDEX "`); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(out, `CREATE INDEX "`); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, index.Name); err != nil {
		return err
	}

	// TODO: useOnlyClause assumes partition-structured metadata. If legacy
	// PostgreSQL table inheritance support is ever needed, revisit with
	// definition-driven ONLY extraction.
	if useOnlyClause {
		if _, err := io.WriteString(out, `" ON ONLY "`); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(out, `" ON "`); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"`); err != nil {
		return err
	}

	// Add USING clause for non-btree indexes
	if index.Type != "" && index.Type != "btree" {
		if _, err := fmt.Fprintf(out, " USING %s", strings.ToUpper(index.Type)); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, ` `); err != nil {
		return err
	}
	if err := writeIndexKeyList(out, index); err != nil {
		return err
	}

	// Write INCLUDE, WITH and WHERE clauses from the full index definition
	// (from pg_get_indexdef), in that order.
	if index.Definition != "" {
		includeClause, whereClause := extractIndexClauses(index.Definition)
		if includeClause != "" {
			if _, err := fmt.Fprintf(out, " %s", includeClause); err != nil {
				return err
			}
		}
		if withClause := extractIndexWith(index.Definition); withClause != "" {
			if _, err := fmt.Fprintf(out, " %s", withClause); err != nil {
				return err
			}
		}
		if whereClause != "" {
			if _, err := fmt.Fprintf(out, " WHERE %s", whereClause); err != nil {
				return err
			}
		}
	}

	if includeTerminatorAndComment {
		if _, err := io.WriteString(out, ";\n\n"); err != nil {
			return err
		}

		if len(index.Comment) > 0 {
			if err := writeIndexComment(out, schema, index); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeSchema(out io.Writer, schema *metadatapb.SchemaMetadata) error {
	if schema.Name == "public" {
		return nil
	}

	if _, err := io.WriteString(out, `CREATE SCHEMA IF NOT EXISTS "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schema.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `";`); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n\n")
	return err
}

// writeSchemaOwnerAndComment restores a schema's owner and comment in a dump.
// The public schema already exists and keeps its owner, so only its comment
// is written.
func writeSchemaOwnerAndComment(out io.Writer, schema *metadatapb.SchemaMetadata) error {
	if schema.Name != "public" && schema.Owner != "" {
		if _, err := fmt.Fprintf(out, "ALTER SCHEMA \"%s\" OWNER TO \"%s\";\n\n", escapeIdentifier(schema.Name), escapeIdentifier(schema.Owner)); err != nil {
			return err
		}
	}
	if schema.Comment != "" {
		if _, err := fmt.Fprintf(out, "COMMENT ON SCHEMA \"%s\" IS '%s';\n\n", schema.Name, escapeSingleQuote(schema.Comment)); err != nil {
			return err
		}
	}
	return nil
}

// extractIndexClauses parses a CREATE INDEX statement using the omni parser and
// extracts the INCLUDE clause and WHERE clause as strings.
func extractIndexClauses(definition string) (includeClause, whereClause string) {
	// The definition comes from a synced index; one that does not parse renders
	// without INCLUDE and WHERE clauses.
	stmts, err := parseDefinition(definition)
	if err != nil || len(stmts) == 0 {
		return "", ""
	}
	idx, ok := stmts[0].(*ast.IndexStmt)
	if !ok || idx == nil {
		return "", ""
	}

	// Extract INCLUDE columns using Loc to preserve original quoting.
	if idx.IndexIncludingParams != nil {
		var cols []string
		for _, item := range idx.IndexIncludingParams.Items {
			if elem, ok := item.(*ast.IndexElem); ok && elem.Loc.Start >= 0 && elem.Loc.End > elem.Loc.Start {
				cols = append(cols, definition[elem.Loc.Start:elem.Loc.End])
			}
		}
		if len(cols) > 0 {
			includeClause = "INCLUDE (" + strings.Join(cols, ", ") + ")"
		}
	}

	// Extract WHERE clause using AST node Loc for precise extraction.
	if idx.WhereClause != nil {
		loc := ast.NodeLoc(idx.WhereClause)
		if loc.Start >= 0 && loc.End > loc.Start && loc.End <= len(definition) {
			raw := definition[loc.Start:loc.End]
			// Strip outer parentheses if the entire predicate is wrapped.
			if strings.HasPrefix(raw, "(") && strings.HasSuffix(raw, ")") {
				inner := raw[1 : len(raw)-1]
				// Only strip if parentheses are balanced.
				if isBalancedParens(inner) {
					raw = inner
				}
			}
			whereClause = strings.TrimSpace(raw)
		}
	}

	return includeClause, whereClause
}

// extractIndexWith returns the storage-parameter clause of a CREATE INDEX
// definition, such as `WITH (pages_per_range='16')`, or "" when the index has
// none. It reads the clause verbatim from the definition, which
// pg_get_indexdef writes between the key list and any WHERE clause. The clause
// is the `WITH (` keyword at bracket depth zero and outside any quotes, so a
// `WITH (` inside the quoted index name, a key expression, or an INCLUDE
// column is not mistaken for it.
func extractIndexWith(definition string) string {
	stmts, err := parseDefinition(definition)
	if err != nil || len(stmts) == 0 {
		return ""
	}
	idx, ok := stmts[0].(*ast.IndexStmt)
	if !ok || idx == nil || idx.Options == nil || len(idx.Options.Items) == 0 {
		return ""
	}
	inSingle, inDouble, depth := false, false, 0
	for i := 0; i < len(definition); i++ {
		c := definition[i]
		switch {
		case inSingle:
			if c == '\'' {
				inSingle = false
			}
		case inDouble:
			if c == '"' {
				inDouble = false
			}
		case c == '\'':
			inSingle = true
		case c == '"':
			inDouble = true
		case c == '(':
			depth++
		case c == ')':
			if depth > 0 {
				depth--
			}
		case depth == 0 && (c == 'W' || c == 'w') && (i == 0 || definition[i-1] == ' ') && strings.HasPrefix(definition[i:], "WITH ("):
			return indexWithClause(definition, i)
		}
	}
	return ""
}

// indexWithClause returns the `WITH (...)` clause of definition that begins at
// start, with its parentheses balanced and single-quoted values skipped.
func indexWithClause(definition string, start int) string {
	open := start + len("WITH ") // the '(' of the clause
	depth, inSingle := 0, false
	for j := open; j < len(definition); j++ {
		c := definition[j]
		switch {
		case inSingle:
			if c == '\'' {
				inSingle = false
			}
		case c == '\'':
			inSingle = true
		case c == '(':
			depth++
		case c == ')':
			depth--
			if depth == 0 {
				return definition[start : j+1]
			}
		}
	}
	return ""
}

// isBalancedParens checks if a string has balanced parentheses.
func isBalancedParens(s string) bool {
	depth := 0
	for _, c := range s {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				return false
			}
		default:
		}
	}
	return depth == 0
}

// parseDefinition parses sql and returns each statement's node in order.
func parseDefinition(sql string) ([]ast.Node, error) {
	list, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	var stmts []ast.Node
	for _, item := range list.Items {
		if raw, ok := item.(*ast.RawStmt); ok {
			stmts = append(stmts, raw.Stmt)
		}
	}
	return stmts, nil
}

// definitionGraph orders object IDs so that each object follows the objects it
// depends on. Ties break by ID.
type definitionGraph struct {
	nodes map[string]bool
	edges [][2]string
}

func newDefinitionGraph() *definitionGraph {
	return &definitionGraph{nodes: make(map[string]bool)}
}

func (g *definitionGraph) AddNode(id string) {
	g.nodes[id] = true
}

func (g *definitionGraph) AddEdge(start, end string) {
	g.nodes[start] = true
	g.nodes[end] = true
	g.edges = append(g.edges, [2]string{start, end})
}

// TopologicalSort returns the IDs in dependency order, or an error when the
// graph has a cycle.
func (g *definitionGraph) TopologicalSort() ([]string, error) {
	inDegree := make(map[string]int)
	outEdges := make(map[string][]string)
	for _, edge := range g.edges {
		inDegree[edge[1]]++
		outEdges[edge[0]] = append(outEdges[edge[0]], edge[1])
	}

	var queue []string
	for id := range g.nodes {
		if inDegree[id] == 0 {
			queue = append(queue, id)
		}
	}
	slices.Sort(queue)

	var result []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, end := range outEdges[node] {
			inDegree[end]--
			if inDegree[end] == 0 {
				queue = append(queue, end)
			}
		}
		result = append(result, node)
	}
	if len(result) != len(g.nodes) {
		return nil, errors.New("graph has cycle")
	}
	return result, nil
}
