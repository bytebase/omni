package catalog

import (
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
}

// DatabaseDefinition renders a synced database snapshot as a pg_dump-style
// restore script.
func DatabaseDefinition(metadata *metadatapb.DatabaseSchemaMetadata, opts DefinitionOptions) (string, error) {
	if len(metadata.Schemas) == 0 {
		return "", nil
	}

	var buf strings.Builder

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
	}

	// Construct extensions.
	for _, extension := range metadata.Extensions {
		if err := writeExtension(&buf, extension); err != nil {
			return "", err
		}
	}

	// Construct enums.
	for _, schema := range metadata.Schemas {
		for _, enum := range schema.EnumTypes {
			if enum.SkipDump {
				continue
			}
			if err := writeEnum(&buf, schema.Name, enum); err != nil {
				return "", err
			}
			if _, err := buf.WriteString(";\n\n"); err != nil {
				return "", err
			}
			// Write enum comment if present
			if len(enum.Comment) > 0 {
				if err := writeEnumComment(&buf, schema.Name, enum); err != nil {
					return "", err
				}
			}
		}
	}

	// Construct composite types.
	if err := writeCompositeTypesSection(&buf, collectCompositeTypes(metadata.Schemas)); err != nil {
		return "", err
	}

	// Build the graph for topological sort.
	graph := newDefinitionGraph()
	functionMap := make(map[string]*metadatapb.FunctionMetadata)
	tableMap := make(map[string]*metadatapb.TableMetadata)
	viewMap := make(map[string]*metadatapb.ViewMetadata)
	materializedViewMap := make(map[string]*metadatapb.MaterializedViewMetadata)

	// Construct functions.
	for _, schema := range metadata.Schemas {
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
			if err := writeTable(&buf, getSchemaNameFromID(objectID), table, sequenceOwnershipMap[objectID]); err != nil {
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
	for _, schema := range metadata.Schemas {
		if err := writeSchemaTriggers(&buf, schema); err != nil {
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

	// Construct enums.
	for _, enum := range schema.EnumTypes {
		if enum.SkipDump {
			continue
		}
		if err := writeEnum(&buf, schema.Name, enum); err != nil {
			return "", err
		}
		if _, err := buf.WriteString(";\n\n"); err != nil {
			return "", err
		}
		// Write enum comment if present
		if len(enum.Comment) > 0 {
			if err := writeEnumComment(&buf, schema.Name, enum); err != nil {
				return "", err
			}
		}
	}

	// Construct composite types. The schema-level SkipDump filter does not
	// apply here: the caller asked for this specific schema's definition.
	var compositeTypes []qualifiedCompositeType
	for _, composite := range schema.CompositeTypes {
		if composite.SkipDump {
			continue
		}
		compositeTypes = append(compositeTypes, qualifiedCompositeType{Schema: schema.Name, Composite: composite})
	}
	if err := writeCompositeTypesSection(&buf, compositeTypes); err != nil {
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
			if err := writeTable(&buf, getSchemaNameFromID(objectID), table, sequenceOwnershipMap[objectID]); err != nil {
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
	if err := writeSchemaTriggers(&buf, schema); err != nil {
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
	if err := writeTable(&buf, schema, table, sequences); err != nil {
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

func getObjectID(schema string, object string) string {
	var buf strings.Builder
	_, _ = buf.WriteString(schema)
	_, _ = buf.WriteString(objectIDSeparator)
	_, _ = buf.WriteString(object)
	return buf.String()
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

		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}

		if _, err := io.WriteString(out, column.Name); err != nil {
			return err
		}

		if _, err := io.WriteString(out, `" `); err != nil {
			return err
		}

		if _, err := io.WriteString(out, column.Type); err != nil {
			return err
		}

		// Handle default values
		if column.Default != "" {
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
func writeSchemaTriggers(out io.Writer, schema *metadatapb.SchemaMetadata) error {
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
		}
	}
	return nil
}

func writeTable(out io.Writer, schema string, table *metadatapb.TableMetadata, sequences []*metadatapb.SequenceMetadata) error {
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

	for i, sequence := range identitySequences {
		if err := writeColumnIdentityGeneration(out, schema, generationTypes[i], sequence); err != nil {
			return err
		}
	}

	// Construct partition tables.
	for _, partition := range table.Partitions {
		if err := writePartitionTable(out, schema, table.Columns, partition); err != nil {
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

func writeAttachPartitionIndex(out io.Writer, schema string, partition *metadatapb.TablePartitionMetadata) error {
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

func writePartitionIndex(out io.Writer, schema string, partition *metadatapb.TablePartitionMetadata) error {
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
		// Note: NULLS ordering information is not available in IndexMetadata
		// so we omit it in the generated DDL
	}

	_, err := io.WriteString(out, ")")
	return err
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

func writePartitionUniqueKey(out io.Writer, schema string, partition *metadatapb.TablePartitionMetadata) error {
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
		if err := writeConstraintComment(out, schema, table, index); err != nil {
			return err
		}
	}

	return nil
}

func writePartitionPrimaryKey(out io.Writer, schema string, partition *metadatapb.TablePartitionMetadata) error {
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
		if err := writeConstraintComment(out, schema, table, index); err != nil {
			return err
		}
	}

	return nil
}

func writeConstraintComment(out io.Writer, schema string, table string, index *metadatapb.IndexMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON CONSTRAINT "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, index.Name); err != nil {
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

	if _, err := io.WriteString(out, escapeSingleQuote(index.Comment)); err != nil {
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
	if _, err := io.WriteString(out, `COMMENT ON TABLE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, escapeSingleQuote(table.Comment)); err != nil {
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

func writeAttachPartition(out io.Writer, schema string, tableName string, partition *metadatapb.TablePartitionMetadata) error {
	if _, err := io.WriteString(out, `ALTER TABLE ONLY "`); err != nil {
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
	if _, err := io.WriteString(out, `" ATTACH PARTITION "`); err != nil {
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
	if _, err := io.WriteString(out, "\" "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, partition.Value); err != nil {
		return err
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

func writePartitionTable(out io.Writer, schema string, columns []*metadatapb.ColumnMetadata, partition *metadatapb.TablePartitionMetadata) error {
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

	// Construct subpartition tables.
	for _, subpartition := range partition.Subpartitions {
		if err := writePartitionTable(out, schema, columns, subpartition); err != nil {
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

	// Write INCLUDE and WHERE clauses from the full index definition (from pg_get_indexdef).
	if index.Definition != "" {
		includeClause, whereClause := extractIndexClauses(index.Definition)
		if includeClause != "" {
			if _, err := fmt.Fprintf(out, " %s", includeClause); err != nil {
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
