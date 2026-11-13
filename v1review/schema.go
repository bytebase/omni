package v1review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
)

// Schema is the schema as it stood before the change, for the rules that
// must look something up in it: the columns of an index a statement only
// renames, say. Bytebase passes the synced snapshot it already holds.
//
// The lookups reproduce how the V1 advisors read that snapshot, name
// comparison included: a name is compared as it is when its kind is case
// sensitive, lower-cased otherwise. Schema and table names follow
// ObjectCaseSensitive. Whether index names are case sensitive is the
// engine's own fact, so each lookup takes it from the engine's rule
// function.
type Schema struct {
	// Metadata is the snapshot. nil holds nothing.
	Metadata *metadata.DatabaseSchemaMetadata
	// ObjectCaseSensitive is whether schema and table names are case
	// sensitive on the instance the snapshot came from.
	ObjectCaseSensitive bool
}

// TableIndex returns the index of that name on the table of that name, nil
// when the schema, the table, or the index is missing. A partition's name
// does not find its table's indexes.
func (s Schema) TableIndex(schemaName, tableName, indexName string, indexCaseSensitive bool) *metadata.IndexMetadata {
	table := lastNamed(s.tables(schemaName), tableName, s.ObjectCaseSensitive, tableNameOf)
	if table == nil {
		return nil
	}
	return lastNamed(table.GetIndexes(), indexName, indexCaseSensitive, indexNameOf)
}

// SchemaIndex returns the first index of that name anywhere in the schema
// and the name of the table or materialized view that holds it: tables
// first, in the snapshot's order, then materialized views. A materialized
// view's index name is compared as an object name, as the V1 advisors did.
//
// Where two tables hold an index of the same name, the V1 advisors took
// whichever their map happened to yield; this takes the first.
func (s Schema) SchemaIndex(schemaName, indexName string, indexCaseSensitive bool) (string, *metadata.IndexMetadata) {
	schema := s.schema(schemaName)
	if schema == nil {
		return "", nil
	}
	for _, table := range s.tables(schemaName) {
		if index := lastNamed(table.GetIndexes(), indexName, indexCaseSensitive, indexNameOf); index != nil {
			return table.GetName(), index
		}
	}
	for _, view := range lastOfEachName(schema.GetMaterializedViews(), s.ObjectCaseSensitive, viewNameOf) {
		for _, index := range view.GetIndexes() {
			if fold(index.GetName(), s.ObjectCaseSensitive) == fold(indexName, s.ObjectCaseSensitive) {
				return view.GetName(), index
			}
		}
	}
	return "", nil
}

// schema returns the schema of that name; the last one wins among equals.
func (s Schema) schema(name string) *metadata.SchemaMetadata {
	return lastNamed(s.Metadata.GetSchemas(), name, s.ObjectCaseSensitive, schemaNameOf)
}

// tables returns the schema's tables, one per name: a later table replaces
// an earlier one of the same name.
func (s Schema) tables(schemaName string) []*metadata.TableMetadata {
	return lastOfEachName(s.schema(schemaName).GetTables(), s.ObjectCaseSensitive, tableNameOf)
}

func schemaNameOf(s *metadata.SchemaMetadata) string         { return s.GetName() }
func tableNameOf(t *metadata.TableMetadata) string           { return t.GetName() }
func indexNameOf(i *metadata.IndexMetadata) string           { return i.GetName() }
func viewNameOf(v *metadata.MaterializedViewMetadata) string { return v.GetName() }

// fold returns a name as its kind compares it.
func fold(name string, caseSensitive bool) string {
	if caseSensitive {
		return name
	}
	return strings.ToLower(name)
}

// lastNamed returns the last item of that name, the zero value when none
// has it.
func lastNamed[T any](items []T, name string, caseSensitive bool, nameOf func(T) string) T {
	var found T
	id := fold(name, caseSensitive)
	for _, item := range items {
		if fold(nameOf(item), caseSensitive) == id {
			found = item
		}
	}
	return found
}

// lastOfEachName returns the items with each name kept once, the last item
// of a name standing where its first stood.
func lastOfEachName[T any](items []T, caseSensitive bool, nameOf func(T) string) []T {
	at := make(map[string]int, len(items))
	var kept []T
	for _, item := range items {
		id := fold(nameOf(item), caseSensitive)
		if i, ok := at[id]; ok {
			kept[i] = item
			continue
		}
		at[id] = len(kept)
		kept = append(kept, item)
	}
	return kept
}
