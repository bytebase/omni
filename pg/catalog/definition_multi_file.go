package catalog

import (
	"fmt"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

// DefinitionFile is one file of MultiFileDefinition's output.
type DefinitionFile struct {
	// Name is the file's path, such as "schemas/public/tables/t.sql".
	Name string
	// Content is the file's SDL.
	Content string
}

// MultiFileDefinition renders a synced database snapshot as SDL split into one
// file per object under schemas/<schema>/. It returns an empty, non-nil slice
// for a snapshot without schemas, and nil when the schemas render no files.
func MultiFileDefinition(metadata *metadatapb.DatabaseSchemaMetadata) ([]DefinitionFile, error) {
	if len(metadata.Schemas) == 0 {
		return []DefinitionFile{}, nil
	}

	var files []DefinitionFile

	// Build skip sequences map (for serial and identity columns)
	skipSequences := buildSkipSequencesMap(metadata)

	// Build table sequences map for each table
	tableSequencesMap := buildTableSequencesMap(metadata)

	tablesMissingColumns := collectTables(metadata.Schemas, tableMissingColumnMetadata)

	// Generate every schema.sql first: extensions may install into these
	// schemas and all other files live inside them.
	for _, schemaMetadata := range metadata.Schemas {
		if schemaMetadata.SkipDump {
			continue
		}
		schemaName := schemaMetadata.Name
		if schemaName == "" {
			schemaName = "public"
		}
		schemaContent, err := getMultiFileSchemaSDL(schemaName, schemaMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to generate schema SDL for %s: %w", schemaName, err)
		}
		if schemaContent != "" {
			files = append(files, DefinitionFile{
				Name:    fmt.Sprintf("schemas/%s/schema.sql", schemaName),
				Content: schemaContent,
			})
		}
	}

	// Generate extensions.sql next: composite type attributes and table
	// columns may use extension-provided types, and the extensions may
	// install into the schemas created above.
	if len(metadata.Extensions) > 0 {
		var buf strings.Builder
		for i, extension := range metadata.Extensions {
			if i > 0 {
				buf.WriteString("\n")
			}

			if err := writeExtension(&buf, extension); err != nil {
				return nil, fmt.Errorf("failed to generate extension SDL for %s: %w", extension.Name, err)
			}
		}

		files = append(files, DefinitionFile{
			Name:    "extensions.sql",
			Content: buf.String(),
		})
	}

	// Generate files for each schema
	for _, schemaMetadata := range metadata.Schemas {
		if schemaMetadata.SkipDump {
			continue
		}

		schemaName := schemaMetadata.Name
		if schemaName == "" {
			schemaName = "public"
		}

		// Collect independent sequences (no owner) for this schema
		var independentSequences []*metadatapb.SequenceMetadata
		for _, sequence := range schemaMetadata.Sequences {
			if sequence.SkipDump {
				continue
			}
			// Skip sequences that belong to serial or identity columns
			sequenceKey := schemaName + "." + sequence.Name
			if skipSequences[sequenceKey] {
				continue
			}

			// Collect independent sequences (no owner)
			if sequence.OwnerTable == "" || sequence.OwnerColumn == "" {
				independentSequences = append(independentSequences, sequence)
			}
		}

		// Generate a single file for all enum and composite types in this
		// schema, listed before table files: table columns may use these types.
		if len(schemaMetadata.EnumTypes) > 0 || len(schemaMetadata.CompositeTypes) > 0 {
			var buf strings.Builder
			hasTypes := false
			for i, enumType := range schemaMetadata.EnumTypes {
				if enumType.SkipDump {
					continue
				}

				if hasTypes {
					buf.WriteString("\n")
				}
				hasTypes = true

				if i > 0 {
					buf.WriteString("\n")
				}

				if err := writeEnum(&buf, schemaName, enumType); err != nil {
					return nil, fmt.Errorf("failed to generate enum type SDL for %s.%s: %w", schemaName, enumType.Name, err)
				}
				buf.WriteString(";\n")

				// Add enum type comment if present
				if len(enumType.Comment) > 0 {
					buf.WriteString("\n")
					if err := writeEnumComment(&buf, schemaName, enumType); err != nil {
						return nil, fmt.Errorf("failed to generate enum type comment for %s.%s: %w", schemaName, enumType.Name, err)
					}
				}
			}

			// Composite types follow enums, ordered so referenced composites
			// come first.
			var compositeTypes []qualifiedCompositeType
			for _, composite := range schemaMetadata.CompositeTypes {
				if composite.SkipDump {
					continue
				}
				compositeTypes = append(compositeTypes, qualifiedCompositeType{Schema: schemaName, Composite: composite})
			}
			for _, t := range sortCompositeTypesTopologically(compositeTypes) {
				if hasTypes {
					buf.WriteString("\n")
				}
				hasTypes = true

				if err := writeCompositeType(&buf, schemaName, t.Composite); err != nil {
					return nil, fmt.Errorf("failed to generate composite type SDL for %s.%s: %w", schemaName, t.Composite.Name, err)
				}
				buf.WriteString(";\n")

				if compositeTypeHasComments(t.Composite) {
					buf.WriteString("\n")
					if err := writeCompositeTypeComments(&buf, schemaName, t.Composite); err != nil {
						return nil, fmt.Errorf("failed to generate composite type comment for %s.%s: %w", schemaName, t.Composite.Name, err)
					}
				}
			}

			if hasTypes {
				files = append(files, DefinitionFile{
					Name:    fmt.Sprintf("schemas/%s/types.sql", schemaName),
					Content: buf.String(),
				})
			}
		}

		// Generate table files
		for _, table := range schemaMetadata.Tables {
			if table.SkipDump {
				continue
			}

			// Get sequences for this table
			tableKey := schemaName + "." + table.Name
			tableSequences := tableSequencesMap[tableKey]

			var buf strings.Builder
			// Write CREATE TABLE statement in SDL format
			if err := writeCreateTableSDL(&buf, schemaName, table, tableSequences, tablesMissingColumns); err != nil {
				return nil, fmt.Errorf("failed to generate table SDL for %s.%s: %w", schemaName, table.Name, err)
			}
			buf.WriteString(";\n\n")

			// Write table comment if present
			if len(table.Comment) > 0 {
				if err := writeTableCommentSDL(&buf, schemaName, table); err != nil {
					return nil, fmt.Errorf("failed to generate table comment for %s.%s: %w", schemaName, table.Name, err)
				}
			}

			// Write column comments if present
			for _, column := range table.Columns {
				if len(column.Comment) > 0 {
					if err := writeColumnCommentSDL(&buf, schemaName, table.Name, column); err != nil {
						return nil, fmt.Errorf("failed to generate column comment for %s.%s.%s: %w", schemaName, table.Name, column.Name, err)
					}
				}
			}

			// Write CREATE INDEX statements for non-constraint indexes
			if err := writeIndexesSDL(&buf, schemaName, table); err != nil {
				return nil, fmt.Errorf("failed to generate indexes SDL for %s.%s: %w", schemaName, table.Name, err)
			}

			// Write index comments if present. Skip when the table has no
			// column metadata — its indexes were not emitted, so a COMMENT ON
			// INDEX would target a nonexistent index.
			if !tableMissingColumnMetadata(table) {
				for _, index := range table.Indexes {
					// Only write comment for standalone indexes (not primary key or unique constraint indexes)
					if !index.Primary && !index.Unique {
						if len(index.Comment) > 0 {
							if err := writeIndexCommentSDL(&buf, schemaName, index); err != nil {
								return nil, fmt.Errorf("failed to generate index comment for %s.%s: %w", schemaName, index.Name, err)
							}
						}
					}
				}
			}

			// Write CREATE TRIGGER statements
			if err := writeTriggersSDL(&buf, schemaName, table); err != nil {
				return nil, fmt.Errorf("failed to generate triggers SDL for %s.%s: %w", schemaName, table.Name, err)
			}

			// Write trigger comments if present. Skip when the table has no
			// column metadata — its triggers were not emitted.
			if !tableMissingColumnMetadata(table) {
				for _, trigger := range table.Triggers {
					if len(trigger.Comment) > 0 {
						if err := writeTriggerCommentSDL(&buf, schemaName, table.Name, trigger); err != nil {
							return nil, fmt.Errorf("failed to generate trigger comment for %s.%s.%s: %w", schemaName, table.Name, trigger.Name, err)
						}
					}
				}
			}

			// Write owned sequences (non-serial/non-identity) for this table
			for _, sequence := range tableSequences {
				// Skip sequences that belong to serial or identity columns
				sequenceKey := schemaName + "." + sequence.Name
				if skipSequences[sequenceKey] {
					continue
				}

				// Only include sequences with owners (non-independent sequences)
				if sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
					buf.WriteString("\n")
					if err := writeSequenceSDL(&buf, schemaName, sequence); err != nil {
						return nil, fmt.Errorf("failed to generate sequence SDL for %s.%s: %w", schemaName, sequence.Name, err)
					}
					buf.WriteString(";\n")

					// Add sequence comment if present
					if len(sequence.Comment) > 0 {
						buf.WriteString("\n")
						if err := writeSequenceCommentSDL(&buf, schemaName, sequence); err != nil {
							return nil, fmt.Errorf("failed to generate sequence comment for %s.%s: %w", schemaName, sequence.Name, err)
						}
					}

					// Add ALTER SEQUENCE OWNED BY. Skip when the owning table
					// has no columns; the OWNED BY clause would reference a
					// column absent from the dump.
					if !tableHasNoColumns(table) {
						buf.WriteString("\n")
						if err := writeAlterSequenceOwnedBy(&buf, schemaName, sequence); err != nil {
							return nil, fmt.Errorf("failed to generate ALTER SEQUENCE OWNED BY for %s.%s: %w", schemaName, sequence.Name, err)
						}
					}
				}
			}

			files = append(files, DefinitionFile{
				Name:    fmt.Sprintf("schemas/%s/tables/%s.sql", schemaName, table.Name),
				Content: buf.String(),
			})
		}

		// Generate view files
		for _, view := range schemaMetadata.Views {
			if view.SkipDump {
				continue
			}

			var buf strings.Builder
			if err := writeViewSDL(&buf, schemaName, view); err != nil {
				return nil, fmt.Errorf("failed to generate view SDL for %s.%s: %w", schemaName, view.Name, err)
			}
			buf.WriteString(";\n")

			// Write view comment if present
			if len(view.Comment) > 0 {
				buf.WriteString("\n")
				if err := writeViewCommentSDL(&buf, schemaName, view); err != nil {
					return nil, fmt.Errorf("failed to generate view comment for %s.%s: %w", schemaName, view.Name, err)
				}
			}

			files = append(files, DefinitionFile{
				Name:    fmt.Sprintf("schemas/%s/views/%s.sql", schemaName, view.Name),
				Content: buf.String(),
			})
		}

		// Generate materialized view files
		for _, materializedView := range schemaMetadata.MaterializedViews {
			if materializedView.SkipDump {
				continue
			}

			var buf strings.Builder
			if err := writeMaterializedViewSDL(&buf, schemaName, materializedView); err != nil {
				return nil, fmt.Errorf("failed to generate materialized view SDL for %s.%s: %w", schemaName, materializedView.Name, err)
			}
			buf.WriteString(";\n")

			// Write materialized view comment if present
			if len(materializedView.Comment) > 0 {
				buf.WriteString("\n")
				if err := writeMaterializedViewCommentSDL(&buf, schemaName, materializedView); err != nil {
					return nil, fmt.Errorf("failed to generate materialized view comment for %s.%s: %w", schemaName, materializedView.Name, err)
				}
			}

			// Write indexes on materialized view
			for _, index := range materializedView.Indexes {
				// Skip constraint-based indexes
				if index.Primary || index.IsConstraint {
					continue
				}

				buf.WriteString("\n")
				if err := writeIndexSDL(&buf, schemaName, materializedView.Name, index, false); err != nil {
					return nil, fmt.Errorf("failed to generate index SDL for %s.%s: %w", schemaName, index.Name, err)
				}
				buf.WriteString(";\n")

				// Write index comment if present
				if len(index.Comment) > 0 {
					buf.WriteString("\n")
					if err := writeIndexCommentSDL(&buf, schemaName, index); err != nil {
						return nil, fmt.Errorf("failed to generate index comment for %s.%s: %w", schemaName, index.Name, err)
					}
				}
			}

			files = append(files, DefinitionFile{
				Name:    fmt.Sprintf("schemas/%s/materialized_views/%s.sql", schemaName, materializedView.Name),
				Content: buf.String(),
			})
		}

		// Generate function files
		for _, function := range schemaMetadata.Functions {
			if function.SkipDump {
				continue
			}

			var buf strings.Builder
			if err := writeFunctionSDL(&buf, schemaName, function); err != nil {
				return nil, fmt.Errorf("failed to generate function SDL for %s.%s: %w", schemaName, function.Name, err)
			}
			buf.WriteString(";\n")

			// Write function comment if present
			if len(function.Comment) > 0 {
				buf.WriteString("\n")
				if err := writeFunctionCommentSDL(&buf, schemaName, function); err != nil {
					return nil, fmt.Errorf("failed to generate function comment for %s.%s: %w", schemaName, function.Name, err)
				}
			}

			// Determine if this is a PROCEDURE or FUNCTION to use different folder
			folderName := "functions"
			if isDefinitionProcedure(function.Definition) {
				folderName = "procedures"
			}

			files = append(files, DefinitionFile{
				Name:    fmt.Sprintf("schemas/%s/%s/%s.sql", schemaName, folderName, function.Name),
				Content: buf.String(),
			})
		}

		// Generate a single file for all independent sequences (no owner) in this schema
		if len(independentSequences) > 0 {
			var buf strings.Builder
			for i, sequence := range independentSequences {
				if i > 0 {
					buf.WriteString("\n")
				}

				if err := writeSequenceSDL(&buf, schemaName, sequence); err != nil {
					return nil, fmt.Errorf("failed to generate sequence SDL for %s.%s: %w", schemaName, sequence.Name, err)
				}
				buf.WriteString(";\n")

				// Add sequence comment if present
				if len(sequence.Comment) > 0 {
					buf.WriteString("\n")
					if err := writeSequenceCommentSDL(&buf, schemaName, sequence); err != nil {
						return nil, fmt.Errorf("failed to generate sequence comment for %s.%s: %w", schemaName, sequence.Name, err)
					}
				}
			}

			files = append(files, DefinitionFile{
				Name:    fmt.Sprintf("schemas/%s/sequences.sql", schemaName),
				Content: buf.String(),
			})
		}
	}

	// Generate event_triggers.sql file if there are any event triggers
	if len(metadata.EventTriggers) > 0 {
		var buf strings.Builder
		for i, eventTrigger := range metadata.EventTriggers {
			if eventTrigger.SkipDump {
				continue
			}
			if i > 0 {
				buf.WriteString("\n")
			}

			if err := writeEventTrigger(&buf, eventTrigger); err != nil {
				return nil, fmt.Errorf("failed to generate event trigger SDL for %s: %w", eventTrigger.Name, err)
			}
		}

		if buf.Len() > 0 {
			files = append(files, DefinitionFile{
				Name:    "event_triggers.sql",
				Content: buf.String(),
			})
		}
	}

	return files, nil
}

func getMultiFileSchemaSDL(schemaName string, schemaMetadata *metadatapb.SchemaMetadata) (string, error) {
	var buf strings.Builder

	if schemaName != "public" {
		if err := writeSchema(&buf, &metadatapb.SchemaMetadata{Name: schemaName}); err != nil {
			return "", err
		}
	}

	if schemaMetadata.GetComment() != "" {
		if err := writeSchemaCommentSDL(&buf, &metadatapb.SchemaMetadata{
			Name:    schemaName,
			Comment: schemaMetadata.GetComment(),
		}); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

// buildSkipSequencesMap builds a map of sequences that should be skipped (serial and identity sequences).
func buildSkipSequencesMap(metadata *metadatapb.DatabaseSchemaMetadata) map[string]bool {
	skipSequences := make(map[string]bool)
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		for _, table := range schema.Tables {
			if table.SkipDump {
				continue
			}
			for _, column := range table.Columns {
				// Check for serial columns
				isSerial, _ := isSerialColumn(column, table.Name, schema.Sequences)
				if isSerial {
					// Extract the sequence name from the DEFAULT clause to match the exact sequence
					// This ensures we skip the correct sequence, especially when multiple sequences
					// claim ownership of the same column.
					sequenceName := extractSequenceNameFromNextval(column.Default)

					for _, sequence := range schema.Sequences {
						// Match by sequence name AND ownership to ensure we skip the exact sequence
						// referenced in the DEFAULT clause
						if sequence.Name == sequenceName && sequence.OwnerTable == table.Name && sequence.OwnerColumn == column.Name {
							sequenceKey := schema.Name + "." + sequence.Name
							skipSequences[sequenceKey] = true
							break
						}
					}
				}
				// Check for identity columns
				if isIdentityColumn(column) {
					for _, sequence := range schema.Sequences {
						if sequence.OwnerTable == table.Name && sequence.OwnerColumn == column.Name {
							sequenceKey := schema.Name + "." + sequence.Name
							skipSequences[sequenceKey] = true
							break
						}
					}
				}
			}
		}
	}
	return skipSequences
}

// buildTableSequencesMap builds a map of sequences by table for easy lookup during table creation.
func buildTableSequencesMap(metadata *metadatapb.DatabaseSchemaMetadata) map[string][]*metadatapb.SequenceMetadata {
	tableSequencesMap := make(map[string][]*metadatapb.SequenceMetadata)
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		for _, sequence := range schema.Sequences {
			if sequence.SkipDump {
				continue
			}
			if sequence.OwnerTable != "" {
				tableKey := schema.Name + "." + sequence.OwnerTable
				tableSequencesMap[tableKey] = append(tableSequencesMap[tableKey], sequence)
			}
		}
	}
	return tableSequencesMap
}
