package catalog

import (
	"errors"
	"fmt"
	"io"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

// writePrimaryKeyConstraintSDL writes a single primary key constraint SDL
func writePrimaryKeyConstraintSDL(out io.Writer, index *metadatapb.IndexMetadata) error {
	if index == nil || !index.Primary {
		return errors.New("invalid primary key constraint")
	}

	if _, err := io.WriteString(out, `CONSTRAINT "`); err != nil {
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
	if _, err := io.WriteString(out, ")"); err != nil {
		return err
	}
	return nil
}

// writeUniqueKeyConstraintSDL writes a single unique key constraint SDL
func writeUniqueKeyConstraintSDL(out io.Writer, index *metadatapb.IndexMetadata) error {
	if index == nil || !index.Unique || index.Primary || !index.IsConstraint {
		return errors.New("invalid unique key constraint")
	}

	if _, err := io.WriteString(out, `CONSTRAINT "`); err != nil {
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
	if _, err := io.WriteString(out, ")"); err != nil {
		return err
	}
	return nil
}

func databaseSDL(metadata *metadatapb.DatabaseSchemaMetadata) (string, error) {
	var buf strings.Builder

	tablesMissingColumns := collectTables(metadata.Schemas, tableMissingColumnMetadata)
	tablesWithoutColumns := collectTables(metadata.Schemas, tableHasNoColumns)

	// Write CREATE SCHEMA statements first for non-public schemas
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		if err := writeSchema(&buf, schema); err != nil {
			return "", err
		}

		// Write schema comment if present
		if len(schema.Comment) > 0 {
			if err := writeSchemaCommentSDL(&buf, schema); err != nil {
				return "", err
			}
		}
	}

	// Write extensions (before enum types and tables as they might provide types used in definitions)
	for _, extension := range metadata.Extensions {
		if err := writeExtension(&buf, extension); err != nil {
			return "", err
		}
	}

	// Build a map of serial and identity sequences that should be skipped
	// Serial sequences are identified by checking if their owner columns match the serial pattern
	// Identity sequences are identified by checking if their owner columns have IdentityGeneration set
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

					// Find the sequence that belongs to this serial column
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
					// Find the sequence that belongs to this identity column
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

	// Write all enum types before sequences and tables
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		for _, enumType := range schema.EnumTypes {
			if enumType.SkipDump {
				continue
			}
			if err := writeEnum(&buf, schema.Name, enumType); err != nil {
				return "", err
			}
			if _, err := buf.WriteString(";\n\n"); err != nil {
				return "", err
			}

			// Write enum type comment if present
			if len(enumType.Comment) > 0 {
				if err := writeEnumComment(&buf, schema.Name, enumType); err != nil {
					return "", err
				}
			}
		}
	}

	// Write all composite types before sequences and tables.
	if err := writeCompositeTypesSection(&buf, collectCompositeTypes(metadata.Schemas)); err != nil {
		return "", err
	}

	// Write all sequences before tables to ensure they exist before any table references them.
	// Skip sequences that belong to serial or identity columns as they will be implicitly created.
	sequenceOwnershipMap := make(map[string][]*metadatapb.SequenceMetadata)

	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}

		for _, sequence := range schema.Sequences {
			if sequence.SkipDump {
				continue
			}
			// Skip sequences that belong to serial or identity columns
			sequenceKey := schema.Name + "." + sequence.Name
			if skipSequences[sequenceKey] {
				continue
			}
			if err := writeSequenceSDL(&buf, schema.Name, sequence); err != nil {
				return "", err
			}
			if _, err := buf.WriteString(";\n\n"); err != nil {
				return "", err
			}

			// Write sequence comment if present
			if len(sequence.Comment) > 0 {
				if err := writeSequenceCommentSDL(&buf, schema.Name, sequence); err != nil {
					return "", err
				}
			}

			// Track sequences with owners for later ownership establishment
			if sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
				tableID := getObjectID(schema.Name, sequence.OwnerTable)
				sequenceOwnershipMap[tableID] = append(sequenceOwnershipMap[tableID], sequence)
			}
		}
	}

	// Build a map of sequences by table for easy lookup during table creation
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

	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}

		// Write tables (sequences are already written above)
		for _, table := range schema.Tables {
			if table.SkipDump {
				continue
			}

			// Get sequences for this table
			tableKey := schema.Name + "." + table.Name
			tableSequences := tableSequencesMap[tableKey]

			// Write CREATE TABLE statement
			if err := writeCreateTableSDL(&buf, schema.Name, table, tableSequences, tablesMissingColumns); err != nil {
				return "", err
			}

			if _, err := buf.WriteString(";\n\n"); err != nil {
				return "", err
			}

			// Write table comment if present
			if len(table.Comment) > 0 {
				if err := writeTableCommentSDL(&buf, schema.Name, table); err != nil {
					return "", err
				}
			}

			// Write column comments if present
			for _, column := range table.Columns {
				if len(column.Comment) > 0 {
					if err := writeColumnCommentSDL(&buf, schema.Name, table.Name, column); err != nil {
						return "", err
					}
				}
			}

			// Write CREATE INDEX statements for non-constraint indexes
			if err := writeIndexesSDL(&buf, schema.Name, table); err != nil {
				return "", err
			}

			// Write index comments if present. Skip when the table has no
			// column metadata — its indexes were not emitted, so a COMMENT ON
			// INDEX would target a nonexistent index.
			if !tableMissingColumnMetadata(table) {
				for _, index := range table.Indexes {
					// Only write comment for standalone indexes (not primary key or unique constraint indexes)
					if !index.Primary && !index.Unique {
						if len(index.Comment) > 0 {
							if err := writeIndexCommentSDL(&buf, schema.Name, index); err != nil {
								return "", err
							}
						}
					}
				}
			}

			// Write CREATE TRIGGER statements
			if err := writeTriggersSDL(&buf, schema.Name, table); err != nil {
				return "", err
			}

			// Write trigger comments if present. Skip when the table has no
			// column metadata — its triggers were not emitted.
			if !tableMissingColumnMetadata(table) {
				for _, trigger := range table.Triggers {
					if len(trigger.Comment) > 0 {
						if err := writeTriggerCommentSDL(&buf, schema.Name, table.Name, trigger); err != nil {
							return "", err
						}
					}
				}
			}
		}

		// Write ALTER SEQUENCE OWNED BY statements for sequences that have owners
		// but are not serial or identity sequences
		for _, sequence := range schema.Sequences {
			if sequence.SkipDump {
				continue
			}
			// Skip sequences that belong to serial or identity columns
			sequenceKey := schema.Name + "." + sequence.Name
			if skipSequences[sequenceKey] {
				continue
			}
			// Write ALTER SEQUENCE OWNED BY for sequences with owners. Skip
			// when the owning table has no columns; the OWNED BY clause
			// would reference a column absent from the dump.
			if sequence.OwnerTable != "" && sequence.OwnerColumn != "" &&
				!tablesWithoutColumns[getObjectID(schema.Name, sequence.OwnerTable)] {
				if err := writeAlterSequenceOwnedBy(&buf, schema.Name, sequence); err != nil {
					return "", err
				}
			}
		}

		// Write views after tables
		for _, view := range schema.Views {
			if view.SkipDump {
				continue
			}

			if err := writeViewSDL(&buf, schema.Name, view); err != nil {
				return "", err
			}

			if _, err := buf.WriteString(";\n\n"); err != nil {
				return "", err
			}

			// Write view comment if present
			if len(view.Comment) > 0 {
				if err := writeViewCommentSDL(&buf, schema.Name, view); err != nil {
					return "", err
				}
			}
		}

		// Write materialized views after views
		for _, materializedView := range schema.MaterializedViews {
			if materializedView.SkipDump {
				continue
			}

			if err := writeMaterializedViewSDL(&buf, schema.Name, materializedView); err != nil {
				return "", err
			}

			if _, err := buf.WriteString(";\n\n"); err != nil {
				return "", err
			}

			// Write materialized view comment if present
			if len(materializedView.Comment) > 0 {
				if err := writeMaterializedViewCommentSDL(&buf, schema.Name, materializedView); err != nil {
					return "", err
				}
			}

			// Write indexes on materialized view
			for _, index := range materializedView.Indexes {
				// Skip constraint-based indexes as they are part of table definition
				if index.Primary || index.IsConstraint {
					continue
				}

				if err := writeIndexSDL(&buf, schema.Name, materializedView.Name, index, false); err != nil {
					return "", err
				}
				if _, err := buf.WriteString(";\n\n"); err != nil {
					return "", err
				}

				// Write index comment if present
				if len(index.Comment) > 0 {
					if err := writeIndexCommentSDL(&buf, schema.Name, index); err != nil {
						return "", err
					}
				}
			}
		}

		// Write functions and procedures after materialized views
		for _, function := range schema.Functions {
			if function.SkipDump {
				continue
			}

			if err := writeFunctionSDL(&buf, schema.Name, function); err != nil {
				return "", err
			}

			if _, err := buf.WriteString(";\n\n"); err != nil {
				return "", err
			}

			// Write function comment if present
			if len(function.Comment) > 0 {
				if err := writeFunctionCommentSDL(&buf, schema.Name, function); err != nil {
					return "", err
				}
			}
		}
	}

	// Write event triggers (last, as they depend on functions)
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

// isSerialColumn checks if a column is a serial/bigserial/smallserial type based on its properties
// Serial pattern: integer/bigint/smallint + NOT NULL + nextval(sequence) default
// extractSequenceNameFromNextval extracts the sequence name from a nextval() default value.
// It handles various PostgreSQL formats including:
// - nextval('sequence_name')
// - nextval('schema.sequence_name')
// - nextval('"Schema"."Sequence"')
// - nextval('public."MixedCase"')
// - nextval('sequence_name'::regclass)
// Returns the unquoted sequence name without schema qualification.
func extractSequenceNameFromNextval(defaultValue string) string {
	if defaultValue == "" {
		return ""
	}

	// Find nextval( - case insensitive
	defaultLower := strings.ToLower(defaultValue)
	nextvalIdx := strings.Index(defaultLower, "nextval(")
	if nextvalIdx == -1 {
		return ""
	}

	// Find the opening quote after nextval(
	startIdx := nextvalIdx + len("nextval(")
	if startIdx >= len(defaultValue) {
		return ""
	}

	// Determine quote type (' or ")
	quoteChar := byte(0)
	for startIdx < len(defaultValue) && (defaultValue[startIdx] == ' ' || defaultValue[startIdx] == '\t') {
		startIdx++
	}
	if startIdx >= len(defaultValue) {
		return ""
	}

	switch defaultValue[startIdx] {
	case '\'':
		quoteChar = '\''
	case '"':
		quoteChar = '"'
	default:
		return ""
	}

	startIdx++ // Skip opening quote

	// Find the closing quote
	endIdx := startIdx
	for endIdx < len(defaultValue) && defaultValue[endIdx] != quoteChar {
		endIdx++
	}
	if endIdx >= len(defaultValue) {
		return ""
	}

	// Extract the full sequence reference (may include schema and ::regclass)
	sequenceRef := defaultValue[startIdx:endIdx]

	// Remove ::regclass suffix if present
	if idx := strings.Index(sequenceRef, "::"); idx != -1 {
		sequenceRef = sequenceRef[:idx]
	}

	// Parse schema-qualified name: schema.sequence or "schema"."sequence"
	// We need to extract just the sequence name, handling quoted identifiers
	sequenceName := extractIdentifierFromQualifiedName(sequenceRef)

	// Remove surrounding quotes from the sequence name if present
	sequenceName = strings.TrimSpace(sequenceName)
	if len(sequenceName) >= 2 && sequenceName[0] == '"' && sequenceName[len(sequenceName)-1] == '"' {
		sequenceName = sequenceName[1 : len(sequenceName)-1]
	}

	return sequenceName
}

// extractIdentifierFromQualifiedName extracts the identifier (sequence name) from a qualified name.
// Handles cases like:
// - sequence_name -> sequence_name
// - schema.sequence_name -> sequence_name
// - "Schema"."Sequence" -> "Sequence"
// - public."MixedCase" -> "MixedCase"
func extractIdentifierFromQualifiedName(qualifiedName string) string {
	if qualifiedName == "" {
		return ""
	}

	// Look for the last dot that's not inside quotes
	inQuotes := false
	lastDotIdx := -1

	for i := 0; i < len(qualifiedName); i++ {
		if qualifiedName[i] == '"' {
			inQuotes = !inQuotes
		} else if qualifiedName[i] == '.' && !inQuotes {
			lastDotIdx = i
		}
	}

	// If we found a dot outside quotes, take everything after it
	if lastDotIdx != -1 {
		return qualifiedName[lastDotIdx+1:]
	}

	// No schema qualification found
	return qualifiedName
}

func isSerialColumn(column *metadatapb.ColumnMetadata, tableName string, sequences []*metadatapb.SequenceMetadata) (isSerial bool, serialType string) {
	// Serial columns must be NOT NULL
	if column.Nullable {
		return false, ""
	}

	// Serial columns must have a nextval() default
	if !strings.Contains(strings.ToLower(column.Default), "nextval(") {
		return false, ""
	}

	// Check type and map to serial type
	var expectedSerialType string
	switch strings.ToLower(column.Type) {
	case "integer":
		expectedSerialType = "serial"
	case "bigint":
		expectedSerialType = "bigserial"
	case "smallint":
		expectedSerialType = "smallserial"
	default:
		return false, ""
	}

	// IMPORTANT: Only treat as serial if the sequence is owned by this column
	// This prevents converting columns that use independent sequences to serial type
	sequenceName := extractSequenceNameFromNextval(column.Default)
	if sequenceName == "" {
		return false, ""
	}

	// Check if any sequence with this name is owned by this column
	for _, seq := range sequences {
		if seq.Name == sequenceName && seq.OwnerTable == tableName && seq.OwnerColumn == column.Name {
			return true, expectedSerialType
		}
	}

	// Sequence is not owned by this column, so don't convert to serial
	return false, ""
}

// writeColumnSDL writes a single column definition to the output writer
// This function is extracted from writeCreateTableSDL to enable code reuse
// sequences parameter is optional and used to find identity sequences for the column
// tableName is used to verify sequence ownership for serial column detection
func writeColumnSDL(out io.Writer, column *metadatapb.ColumnMetadata, tableName string, sequences []*metadatapb.SequenceMetadata) error {
	if _, err := io.WriteString(out, `"`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, column.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" `); err != nil {
		return err
	}

	// Check if this is an identity column
	if isIdentityColumn(column) {
		// Find the sequence for this identity column
		var identitySequence *metadatapb.SequenceMetadata
		for _, seq := range sequences {
			if seq.OwnerColumn == column.Name {
				identitySequence = seq
				break
			}
		}

		// Write type
		if _, err := io.WriteString(out, column.Type); err != nil {
			return err
		}

		// Write GENERATED ... AS IDENTITY
		if _, err := io.WriteString(out, " GENERATED "); err != nil {
			return err
		}
		if column.IdentityGeneration == metadatapb.ColumnMetadata_ALWAYS {
			if _, err := io.WriteString(out, "ALWAYS"); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(out, "BY DEFAULT"); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, " AS IDENTITY"); err != nil {
			return err
		}

		// Write identity options if we have the sequence
		if identitySequence != nil {
			hasOptions := false
			// Check if we need to write options (non-default values)
			if identitySequence.Start != "" && identitySequence.Start != "1" {
				hasOptions = true
			}
			if identitySequence.Increment != "" && identitySequence.Increment != "1" {
				hasOptions = true
			}

			if hasOptions {
				if _, err := io.WriteString(out, " ("); err != nil {
					return err
				}

				first := true
				if identitySequence.Start != "" && identitySequence.Start != "1" {
					if _, err := fmt.Fprintf(out, "START WITH %s", identitySequence.Start); err != nil {
						return err
					}
					first = false
				}
				if identitySequence.Increment != "" && identitySequence.Increment != "1" {
					if !first {
						if _, err := io.WriteString(out, " "); err != nil {
							return err
						}
					}
					if _, err := fmt.Fprintf(out, "INCREMENT BY %s", identitySequence.Increment); err != nil {
						return err
					}
				}

				if _, err := io.WriteString(out, ")"); err != nil {
					return err
				}
			}
		}

		// Identity columns are implicitly NOT NULL, skip that
		return nil
	}

	// Check if this is a serial column and convert it back to serial type
	isSerial, serialType := isSerialColumn(column, tableName, sequences)
	if isSerial {
		// Use serial type instead of integer + nextval + NOT NULL
		if _, err := io.WriteString(out, serialType); err != nil {
			return err
		}
		// Skip DEFAULT and NOT NULL as they are implicit in serial
	} else {
		// Normal column handling
		if _, err := io.WriteString(out, column.Type); err != nil {
			return err
		}

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

	if column.Collation != "" {
		if _, err := io.WriteString(out, ` COLLATE "`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, column.Collation); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
	}

	return nil
}

func writeCreateTableSDL(out io.Writer, schemaName string, table *metadatapb.TableMetadata, sequences []*metadatapb.SequenceMetadata, tablesMissingColumns map[string]bool) error {
	if _, err := io.WriteString(out, `CREATE TABLE "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schemaName); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, table.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" (`); err != nil {
		return err
	}

	// first tracks whether any element has been written inside the parens so
	// the separator on constraints is suppressed when there are no columns.
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

	// Write columns
	for _, column := range table.Columns {
		if err := writeItemSep(); err != nil {
			return err
		}
		if err := writeColumnSDL(out, column, table.Name, sequences); err != nil {
			return err
		}
	}

	// Write table constraints
	if err := writeTableConstraintsSDL(out, table, writeItemSep, tablesMissingColumns); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n)")
	return err
}

func writeTableConstraintsSDL(out io.Writer, table *metadatapb.TableMetadata, writeSep func() error, tablesMissingColumns map[string]bool) error {
	if tableMissingColumnMetadata(table) {
		// Constraints would reference columns absent from the dump.
		return nil
	}

	// Write primary key constraint
	for _, index := range table.Indexes {
		if index.Primary {
			if err := writeSep(); err != nil {
				return err
			}
			if err := writePrimaryKeyConstraintSDL(out, index); err != nil {
				return err
			}
		}
	}

	// Write unique constraints
	for _, index := range table.Indexes {
		if index.Unique && !index.Primary && index.IsConstraint {
			if err := writeSep(); err != nil {
				return err
			}
			if err := writeUniqueKeyConstraintSDL(out, index); err != nil {
				return err
			}
		}
	}

	// Write check constraints
	for _, check := range table.CheckConstraints {
		if err := writeSep(); err != nil {
			return err
		}
		if err := writeCheckConstraintSDL(out, check); err != nil {
			return err
		}
	}

	// Write EXCLUDE constraints
	for _, exclude := range table.ExcludeConstraints {
		if err := writeSep(); err != nil {
			return err
		}
		if err := writeExcludeConstraintSDL(out, exclude); err != nil {
			return err
		}
	}

	// Write foreign key constraints. Skip foreign keys whose referenced
	// table has no column metadata; the REFERENCES clause would name
	// columns absent from the dump.
	for _, fk := range table.ForeignKeys {
		if tablesMissingColumns[getObjectID(fk.ReferencedSchema, fk.ReferencedTable)] {
			continue
		}
		if err := writeSep(); err != nil {
			return err
		}
		if err := writeForeignKeyConstraintSDL(out, fk); err != nil {
			return err
		}
	}

	return nil
}

func writeIndexesSDL(out io.Writer, schemaName string, table *metadatapb.TableMetadata) error {
	if tableMissingColumnMetadata(table) {
		// Indexes would reference columns absent from the dump.
		return nil
	}
	for _, index := range table.Indexes {
		// Skip indexes that are constraints (primary key, unique constraints)
		// These are already handled in the CREATE TABLE statement
		if index.Primary || index.IsConstraint {
			continue
		}

		if err := writeIndexSDL(out, schemaName, table.Name, index, len(table.Partitions) > 0); err != nil {
			return err
		}

		if _, err := io.WriteString(out, ";\n\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeIndexSDL(out io.Writer, schemaName string, tableName string, index *metadatapb.IndexMetadata, useOnlyClause bool) error {
	return writeIndexInternal(out, schemaName, tableName, index, useOnlyClause, false)
}

func writeViewSDL(out io.Writer, schemaName string, view *metadatapb.ViewMetadata) error {
	if _, err := io.WriteString(out, `CREATE VIEW "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schemaName); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" AS `); err != nil {
		return err
	}

	// The view definition should already include the SELECT statement
	definition := strings.TrimSpace(view.Definition)
	// Remove trailing semicolon if present
	definition = strings.TrimSuffix(definition, ";")

	_, err := io.WriteString(out, definition)
	return err
}

func writeMaterializedViewSDL(out io.Writer, schemaName string, view *metadatapb.MaterializedViewMetadata) error {
	if _, err := io.WriteString(out, `CREATE MATERIALIZED VIEW "`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, schemaName); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}

	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}

	if _, err := io.WriteString(out, `" AS `); err != nil {
		return err
	}

	// The materialized view definition should already include the SELECT statement
	definition := strings.TrimSpace(view.Definition)
	// Remove trailing semicolon if present
	definition = strings.TrimSuffix(definition, ";")

	_, err := io.WriteString(out, definition)
	return err
}

func writeFunctionSDL(out io.Writer, _ string, function *metadatapb.FunctionMetadata) error {
	// The function definition should already include the complete CREATE FUNCTION statement
	definition := strings.TrimSpace(function.Definition)
	// Remove trailing semicolon if present
	definition = strings.TrimSuffix(definition, ";")

	_, err := io.WriteString(out, definition)
	return err
}

func writeSequenceSDL(out io.Writer, schemaName string, sequence *metadatapb.SequenceMetadata) error {
	// Write CREATE SEQUENCE statement with schema and name
	if _, err := fmt.Fprintf(out, "CREATE SEQUENCE \"%s\".\"%s\"", schemaName, sequence.Name); err != nil {
		return err
	}

	// Add data type (always output for consistency with non-SDL format)
	if sequence.DataType != "" {
		if _, err := fmt.Fprintf(out, " AS %s", sequence.DataType); err != nil {
			return err
		}
	}

	// Add START WITH (always output for consistency with non-SDL format)
	if sequence.Start != "" {
		if _, err := fmt.Fprintf(out, " START WITH %s", sequence.Start); err != nil {
			return err
		}
	}

	// Add INCREMENT BY (always output for consistency with non-SDL format)
	if sequence.Increment != "" {
		if _, err := fmt.Fprintf(out, " INCREMENT BY %s", sequence.Increment); err != nil {
			return err
		}
	}

	// Add MINVALUE (always output for consistency with non-SDL format)
	if sequence.MinValue != "" {
		if _, err := fmt.Fprintf(out, " MINVALUE %s", sequence.MinValue); err != nil {
			return err
		}
	}

	// Add MAXVALUE (always output for consistency with non-SDL format)
	if sequence.MaxValue != "" {
		if _, err := fmt.Fprintf(out, " MAXVALUE %s", sequence.MaxValue); err != nil {
			return err
		}
	}

	// Add CYCLE/NO CYCLE (always output for consistency with non-SDL format)
	if sequence.Cycle {
		if _, err := io.WriteString(out, " CYCLE"); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(out, " NO CYCLE"); err != nil {
			return err
		}
	}

	// Add CACHE (new field to match non-SDL format)
	if sequence.CacheSize != "" {
		if _, err := fmt.Fprintf(out, " CACHE %s", sequence.CacheSize); err != nil {
			return err
		}
	}

	return nil
}

// writeCheckConstraintSDL writes a single check constraint SDL
func writeCheckConstraintSDL(out io.Writer, check *metadatapb.CheckConstraintMetadata) error {
	if _, err := io.WriteString(out, `CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, check.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" CHECK `); err != nil {
		return err
	}
	if _, err := io.WriteString(out, check.Expression); err != nil {
		return err
	}
	return nil
}

func writeExcludeConstraintSDL(out io.Writer, exclude *metadatapb.ExcludeConstraintMetadata) error {
	if _, err := io.WriteString(out, `CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, exclude.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" `); err != nil {
		return err
	}
	// Expression already includes "EXCLUDE USING ..." prefix
	if _, err := io.WriteString(out, exclude.Expression); err != nil {
		return err
	}
	return nil
}

// writeForeignKeyConstraintSDL writes a single foreign key constraint SDL
func writeForeignKeyConstraintSDL(out io.Writer, fk *metadatapb.ForeignKeyMetadata) error {
	if _, err := io.WriteString(out, `CONSTRAINT "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, fk.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" FOREIGN KEY (`); err != nil {
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
	if _, err := io.WriteString(out, ") REFERENCES "); err != nil {
		return err
	}

	// Always add schema qualifier according to project principle
	if fk.ReferencedSchema != "" {
		if _, err := io.WriteString(out, `"`); err != nil {
			return err
		}
		if _, err := io.WriteString(out, fk.ReferencedSchema); err != nil {
			return err
		}
		if _, err := io.WriteString(out, `".`); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, `"`); err != nil {
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
		if _, err := io.WriteString(out, " ON DELETE "); err != nil {
			return err
		}
		if _, err := io.WriteString(out, fk.OnDelete); err != nil {
			return err
		}
	}

	// Add ON UPDATE clause if specified
	if fk.OnUpdate != "" && fk.OnUpdate != "NO ACTION" {
		if _, err := io.WriteString(out, " ON UPDATE "); err != nil {
			return err
		}
		if _, err := io.WriteString(out, fk.OnUpdate); err != nil {
			return err
		}
	}
	return nil
}

func writeMaterializedViewCommentSDL(out io.Writer, schemaName string, view *metadatapb.MaterializedViewMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON MATERIALIZED VIEW "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
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

func writeSchemaCommentSDL(out io.Writer, schema *metadatapb.SchemaMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON SCHEMA "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schema.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" IS '`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, escapeSingleQuote(schema.Comment)); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeTableCommentSDL(out io.Writer, schemaName string, table *metadatapb.TableMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON TABLE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
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
	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeColumnCommentSDL(out io.Writer, schemaName, tableName string, column *metadatapb.ColumnMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON COLUMN "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tableName); err != nil {
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
	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeViewCommentSDL(out io.Writer, schemaName string, view *metadatapb.ViewMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON VIEW "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
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

func writeFunctionCommentSDL(out io.Writer, schemaName string, function *metadatapb.FunctionMetadata) error {
	// Determine if this is a PROCEDURE or FUNCTION by checking the definition
	objectType := "FUNCTION"
	if isDefinitionProcedure(function.Definition) {
		objectType = "PROCEDURE"
	}

	if _, err := io.WriteString(out, "COMMENT ON "+objectType+" \""); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
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
	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeSequenceCommentSDL(out io.Writer, schemaName string, sequence *metadatapb.SequenceMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON SEQUENCE "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
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

func writeIndexCommentSDL(out io.Writer, schemaName string, index *metadatapb.IndexMetadata) error {
	if _, err := io.WriteString(out, `COMMENT ON INDEX "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
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
	if _, err := io.WriteString(out, `';`); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\n\n")
	return err
}

func writeTriggersSDL(out io.Writer, schemaName string, table *metadatapb.TableMetadata) error {
	if tableMissingColumnMetadata(table) {
		// Trigger definitions (UPDATE OF, WHEN clauses) can reference
		// columns absent from the dump.
		return nil
	}
	for _, trigger := range table.Triggers {
		if trigger.SkipDump {
			continue
		}

		if err := writeTriggerSDL(out, schemaName, table.Name, trigger); err != nil {
			return err
		}

		if _, err := io.WriteString(out, ";\n\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeTriggerSDL(out io.Writer, _ /* schemaName */, _ /* tableName */ string, trigger *metadatapb.TriggerMetadata) error {
	// For PostgreSQL, trigger.Body holds the complete CREATE TRIGGER statement.
	_, err := io.WriteString(out, trigger.Body)
	return err
}

func writeTriggerCommentSDL(out io.Writer, schemaName, tableName string, trigger *metadatapb.TriggerMetadata) error {
	if len(trigger.Comment) == 0 {
		return nil
	}

	if _, err := io.WriteString(out, `COMMENT ON TRIGGER "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `" ON "`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, schemaName); err != nil {
		return err
	}
	if _, err := io.WriteString(out, `"."`); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tableName); err != nil {
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
