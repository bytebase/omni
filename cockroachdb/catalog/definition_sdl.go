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

// databaseSDL renders a snapshot as SDL: one statement per object, with a
// table's keys, checks and foreign keys inside its CREATE TABLE.
func databaseSDL(metadata *metadatapb.DatabaseSchemaMetadata) string {
	r := newRenderer(metadata.Name, metadata.Schemas)
	var buf strings.Builder
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		writeSchema(&buf, schema)
		if schema.Comment != "" {
			fmt.Fprintf(&buf, "COMMENT ON SCHEMA %s IS '%s';\n\n", quoteIdent(schema.Name), escapeSingleQuote(schema.Comment))
		}
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
			if !sequence.SkipDump {
				writeSequenceSDL(&buf, schema.Name, sequence)
			}
		}
	}
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		r.writeSchemaObjectsSDL(&buf, schema)
	}
	return buf.String()
}

// writeSchemaObjectsSDL writes a schema's tables, views, materialized views
// and functions as SDL.
func (r renderer) writeSchemaObjectsSDL(buf *strings.Builder, schema *metadatapb.SchemaMetadata) {
	for _, table := range schema.Tables {
		if !table.SkipDump {
			r.writeTableSDL(buf, schema.Name, table)
		}
	}
	for _, sequence := range schema.Sequences {
		if !sequence.SkipDump && sequence.OwnerTable != "" && sequence.OwnerColumn != "" {
			writeAlterSequenceOwnedBy(buf, schema.Name, sequence)
		}
	}
	for _, view := range schema.Views {
		if !view.SkipDump {
			r.writeView(buf, schema.Name, view)
		}
	}
	for _, view := range schema.MaterializedViews {
		if !view.SkipDump {
			r.writeMaterializedViewSDL(buf, schema.Name, view)
		}
	}
	for _, function := range schema.Functions {
		if !function.SkipDump {
			r.writeFunction(buf, schema.Name, function)
		}
	}
}

func (r renderer) writeTableSDL(buf *strings.Builder, schema string, table *metadatapb.TableMetadata) {
	columns, primary, checks := tableShape(table)
	fmt.Fprintf(buf, "CREATE TABLE %s.%s (", quoteIdent(schema), quoteIdent(table.Name))
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
	for _, fk := range table.ForeignKeys {
		sep()
		fmt.Fprintf(buf, "CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s.%s (%s)",
			quoteIdent(fk.Name), quoteList(fk.Columns), quoteIdent(fk.ReferencedSchema), quoteIdent(fk.ReferencedTable), quoteList(fk.ReferencedColumns))
		if fk.MatchType == "FULL" {
			buf.WriteString(" MATCH FULL")
		}
		if fk.OnDelete != "" && fk.OnDelete != "NO ACTION" {
			buf.WriteString(" ON DELETE " + fk.OnDelete)
		}
		if fk.OnUpdate != "" && fk.OnUpdate != "NO ACTION" {
			buf.WriteString(" ON UPDATE " + fk.OnUpdate)
		}
	}
	buf.WriteString("\n);\n\n")
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
	for _, fk := range table.ForeignKeys {
		if fk.Comment != "" {
			writeConstraintComment(buf, schema, table.Name, fk.Name, fk.Comment)
		}
	}
	r.writeIndexes(buf, schema, table.Name, table.Indexes)
}

func (r renderer) writeMaterializedViewSDL(buf *strings.Builder, schema string, view *metadatapb.MaterializedViewMetadata) {
	fmt.Fprintf(buf, "CREATE MATERIALIZED VIEW %s.%s AS %s;\n\n", quoteIdent(schema), quoteIdent(view.Name), r.query(view.Definition))
	r.writeIndexes(buf, schema, view.Name, materializedViewIndexes(view))
}

func writeSequenceSDL(buf *strings.Builder, schema string, sequence *metadatapb.SequenceMetadata) {
	fmt.Fprintf(buf, "CREATE SEQUENCE %s.%s", quoteIdent(schema), quoteIdent(sequence.Name))
	for _, clause := range [][2]string{
		{"AS", sequence.DataType},
		{"START WITH", sequence.Start},
		{"INCREMENT BY", sequence.Increment},
		{"MINVALUE", sequence.MinValue},
		{"MAXVALUE", sequence.MaxValue},
	} {
		if clause[1] != "" {
			fmt.Fprintf(buf, " %s %s", clause[0], clause[1])
		}
	}
	if sequence.Cycle {
		buf.WriteString(" CYCLE")
	} else {
		buf.WriteString(" NO CYCLE")
	}
	if sequence.CacheSize != "" {
		fmt.Fprintf(buf, " CACHE %s", sequence.CacheSize)
	}
	buf.WriteString(";\n\n")
}

// MultiFileDefinition renders a synced database snapshot as SDL split into one
// file per object under schemas/<schema>/. It returns an empty, non-nil slice
// for a snapshot without schemas, and nil when the schemas render no files.
func MultiFileDefinition(metadata *metadatapb.DatabaseSchemaMetadata) ([]DefinitionFile, error) {
	if len(metadata.GetSchemas()) == 0 {
		return []DefinitionFile{}, nil
	}
	r := newRenderer(metadata.Name, metadata.Schemas)
	var files []DefinitionFile
	add := func(name string, write func(*strings.Builder)) {
		var buf strings.Builder
		write(&buf)
		if buf.Len() > 0 {
			files = append(files, DefinitionFile{Name: name, Content: buf.String()})
		}
	}
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		add(fmt.Sprintf("schemas/%s/schema.sql", schema.Name), func(buf *strings.Builder) {
			writeSchema(buf, schema)
			if schema.Comment != "" {
				fmt.Fprintf(buf, "COMMENT ON SCHEMA %s IS '%s';\n\n", quoteIdent(schema.Name), escapeSingleQuote(schema.Comment))
			}
		})
	}
	add("extensions.sql", func(buf *strings.Builder) {
		for _, extension := range metadata.Extensions {
			writeExtension(buf, extension)
		}
	})
	for _, schema := range metadata.Schemas {
		if schema.SkipDump {
			continue
		}
		add(fmt.Sprintf("schemas/%s/types.sql", schema.Name), func(buf *strings.Builder) {
			r.writeTypes(buf, collectTypes([]*metadatapb.SchemaMetadata{schema}, false))
		})
		for _, table := range schema.Tables {
			if table.SkipDump {
				continue
			}
			add(fmt.Sprintf("schemas/%s/tables/%s.sql", schema.Name, table.Name), func(buf *strings.Builder) {
				r.writeTableSDL(buf, schema.Name, table)
				for _, sequence := range schema.Sequences {
					if !sequence.SkipDump && sequence.OwnerTable == table.Name && sequence.OwnerColumn != "" {
						writeSequenceSDL(buf, schema.Name, sequence)
						writeAlterSequenceOwnedBy(buf, schema.Name, sequence)
					}
				}
			})
		}
		for _, view := range schema.Views {
			if view.SkipDump {
				continue
			}
			add(fmt.Sprintf("schemas/%s/views/%s.sql", schema.Name, view.Name), func(buf *strings.Builder) {
				r.writeView(buf, schema.Name, view)
			})
		}
		for _, view := range schema.MaterializedViews {
			if view.SkipDump {
				continue
			}
			add(fmt.Sprintf("schemas/%s/materialized_views/%s.sql", schema.Name, view.Name), func(buf *strings.Builder) {
				r.writeMaterializedViewSDL(buf, schema.Name, view)
			})
		}
		// Overloads of one name share its file.
		var routineFiles []string
		routines := make(map[string][]*metadatapb.FunctionMetadata)
		for _, function := range schema.Functions {
			if function.SkipDump {
				continue
			}
			folder := "functions"
			if isProcedure(function.Definition) {
				folder = "procedures"
			}
			name := fmt.Sprintf("schemas/%s/%s/%s.sql", schema.Name, folder, function.Name)
			if routines[name] == nil {
				routineFiles = append(routineFiles, name)
			}
			routines[name] = append(routines[name], function)
		}
		for _, name := range routineFiles {
			add(name, func(buf *strings.Builder) {
				for _, function := range routines[name] {
					r.writeFunction(buf, schema.Name, function)
				}
			})
		}
		add(fmt.Sprintf("schemas/%s/sequences.sql", schema.Name), func(buf *strings.Builder) {
			for _, sequence := range schema.Sequences {
				if !sequence.SkipDump && (sequence.OwnerTable == "" || sequence.OwnerColumn == "") {
					writeSequenceSDL(buf, schema.Name, sequence)
				}
			}
		})
	}
	return files, nil
}
