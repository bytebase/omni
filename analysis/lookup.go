package analysis

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// Lookup finds the objects of one synced database by name, the way the
// caller's own metadata model does, so an engine resolves a name to the same
// object the caller would. A name of a kind that compares case-insensitively is
// folded to lower case; where two objects share a folded name, the later one
// wins. Object case sensitivity covers schemas, tables, views, materialized
// views and external tables; detail case sensitivity covers sequences and
// functions, and is the engine's own fact.
type Lookup struct {
	metadata            *metadata.DatabaseSchemaMetadata
	objectCaseSensitive bool
	detailCaseSensitive bool
	schemas             map[string]*SchemaLookup
}

// SchemaLookup finds the objects of one schema by name.
type SchemaLookup struct {
	proto               *metadata.SchemaMetadata
	objectCaseSensitive bool
	detailCaseSensitive bool
	tables              map[string]*metadata.TableMetadata
	// partitions maps a partition's name, at any depth, to its root table.
	partitions        map[string]*metadata.TableMetadata
	views             map[string]*metadata.ViewMetadata
	materializedViews map[string]*metadata.MaterializedViewMetadata
	externalTables    map[string]*metadata.ExternalTableMetadata
	sequences         map[string]*metadata.SequenceMetadata
}

// NewLookup indexes a database's metadata. A nil metadata finds nothing.
func NewLookup(database *metadata.DatabaseSchemaMetadata, objectCaseSensitive, detailCaseSensitive bool) *Lookup {
	l := &Lookup{metadata: database, objectCaseSensitive: objectCaseSensitive, detailCaseSensitive: detailCaseSensitive, schemas: map[string]*SchemaLookup{}}
	for _, schema := range database.GetSchemas() {
		s := &SchemaLookup{
			proto:               schema,
			objectCaseSensitive: objectCaseSensitive,
			detailCaseSensitive: detailCaseSensitive,
			tables:              map[string]*metadata.TableMetadata{},
			partitions:          map[string]*metadata.TableMetadata{},
			views:               map[string]*metadata.ViewMetadata{},
			materializedViews:   map[string]*metadata.MaterializedViewMetadata{},
			externalTables:      map[string]*metadata.ExternalTableMetadata{},
			sequences:           map[string]*metadata.SequenceMetadata{},
		}
		for _, table := range schema.GetTables() {
			s.tables[fold(table.GetName(), objectCaseSensitive)] = table
			s.addPartitions(table, table.GetPartitions())
		}
		for _, table := range schema.GetExternalTables() {
			s.externalTables[fold(table.GetName(), objectCaseSensitive)] = table
		}
		for _, view := range schema.GetViews() {
			s.views[fold(view.GetName(), objectCaseSensitive)] = view
		}
		for _, view := range schema.GetMaterializedViews() {
			s.materializedViews[fold(view.GetName(), objectCaseSensitive)] = view
		}
		for _, sequence := range schema.GetSequences() {
			s.sequences[fold(sequence.GetName(), detailCaseSensitive)] = sequence
		}
		l.schemas[fold(schema.GetName(), objectCaseSensitive)] = s
	}
	return l
}

func (s *SchemaLookup) addPartitions(root *metadata.TableMetadata, partitions []*metadata.TablePartitionMetadata) {
	for _, partition := range partitions {
		s.partitions[fold(partition.GetName(), s.objectCaseSensitive)] = root
		s.addPartitions(root, partition.GetSubpartitions())
	}
}

// Metadata returns the database's metadata.
func (l *Lookup) Metadata() *metadata.DatabaseSchemaMetadata { return l.metadata }

// Schema returns the schema of that name, or nil.
func (l *Lookup) Schema(name string) *SchemaLookup {
	if l == nil {
		return nil
	}
	s, ok := l.schemas[fold(name, l.objectCaseSensitive)]
	if !ok {
		return nil
	}
	return s
}

// ListSchemaNames returns each schema's name once, sorted.
func (l *Lookup) ListSchemaNames() []string {
	names := make([]string, 0, len(l.schemas))
	for _, s := range l.schemas {
		names = append(names, s.proto.GetName())
	}
	slices.Sort(names)
	return names
}

// SearchRelation returns the schema of the first relation of that name along
// searchPath, and the name; empty strings when none is found. A sequence
// counts, being in a relation namespace; a routine does not.
func (l *Lookup) SearchRelation(searchPath []string, name string) (string, string) {
	for _, schemaName := range searchPath {
		s := l.Schema(schemaName)
		if s == nil {
			continue
		}
		if s.Table(name) != nil || s.View(name) != nil || s.MaterializedView(name) != nil || s.ExternalTable(name) != nil || s.Sequence(name) != nil {
			return s.proto.GetName(), name
		}
	}
	return "", ""
}

// SearchFunctions returns every function of that name along searchPath, with
// the schema each was found in.
func (l *Lookup) SearchFunctions(searchPath []string, name string) ([]string, []*metadata.FunctionMetadata) {
	var schemas []string
	var functions []*metadata.FunctionMetadata
	for _, schemaName := range searchPath {
		s := l.Schema(schemaName)
		if s == nil {
			continue
		}
		for _, function := range s.proto.GetFunctions() {
			if sameName(function.GetName(), name, l.detailCaseSensitive) {
				schemas = append(schemas, s.proto.GetName())
				functions = append(functions, function)
			}
		}
	}
	return schemas, functions
}

// SearchIndex returns the schema and holding table of the first index of that
// name along searchPath, as SchemaLookup.IndexTable finds it.
func (l *Lookup) SearchIndex(searchPath []string, name string) (string, string, bool) {
	for _, schemaName := range searchPath {
		s := l.Schema(schemaName)
		if table, ok := s.IndexTable(name); ok {
			return s.proto.GetName(), table, true
		}
	}
	return "", "", false
}

// Proto returns the schema's metadata.
func (s *SchemaLookup) Proto() *metadata.SchemaMetadata { return s.proto }

// Table returns the table of that name, or the root table of the partition of
// that name, or nil.
func (s *SchemaLookup) Table(name string) *metadata.TableMetadata {
	if s == nil {
		return nil
	}
	key := fold(name, s.objectCaseSensitive)
	if table, ok := s.tables[key]; ok {
		return table
	}
	return s.partitions[key]
}

// View returns the view of that name, or nil.
func (s *SchemaLookup) View(name string) *metadata.ViewMetadata {
	if s == nil {
		return nil
	}
	return s.views[fold(name, s.objectCaseSensitive)]
}

// MaterializedView returns the materialized view of that name, or nil.
func (s *SchemaLookup) MaterializedView(name string) *metadata.MaterializedViewMetadata {
	if s == nil {
		return nil
	}
	return s.materializedViews[fold(name, s.objectCaseSensitive)]
}

// ExternalTable returns the external table of that name, or nil.
func (s *SchemaLookup) ExternalTable(name string) *metadata.ExternalTableMetadata {
	if s == nil {
		return nil
	}
	return s.externalTables[fold(name, s.objectCaseSensitive)]
}

// Sequence returns the sequence of that name, or nil.
func (s *SchemaLookup) Sequence(name string) *metadata.SequenceMetadata {
	if s == nil {
		return nil
	}
	return s.sequences[fold(name, s.detailCaseSensitive)]
}

// Function returns the first function of that name, or nil.
func (s *SchemaLookup) Function(name string) *metadata.FunctionMetadata {
	if s == nil {
		return nil
	}
	for _, function := range s.proto.GetFunctions() {
		if sameName(function.GetName(), name, s.detailCaseSensitive) {
			return function
		}
	}
	return nil
}

// IndexTable returns the name of the table, or else the materialized view,
// that holds the index of that name. A table's index names follow detail case
// sensitivity and a materialized view's follow object case sensitivity, as the
// caller's model has it. Where two such indexes share a name, the first in
// metadata order wins.
func (s *SchemaLookup) IndexTable(name string) (string, bool) {
	if s == nil {
		return "", false
	}
	for _, table := range s.proto.GetTables() {
		if s.tables[fold(table.GetName(), s.objectCaseSensitive)] != table {
			continue
		}
		for _, index := range table.GetIndexes() {
			if sameName(index.GetName(), name, s.detailCaseSensitive) {
				return table.GetName(), true
			}
		}
	}
	for _, view := range s.proto.GetMaterializedViews() {
		if s.materializedViews[fold(view.GetName(), s.objectCaseSensitive)] != view {
			continue
		}
		for _, index := range view.GetIndexes() {
			if sameName(index.GetName(), name, s.objectCaseSensitive) {
				return view.GetName(), true
			}
		}
	}
	return "", false
}

// ListTableNames returns the names of the schema's tables, partitions
// excluded, sorted.
func (s *SchemaLookup) ListTableNames() []string {
	return sortedNames(s.tables, (*metadata.TableMetadata).GetName)
}

// ListViewNames returns the names of the schema's views, sorted.
func (s *SchemaLookup) ListViewNames() []string {
	return sortedNames(s.views, (*metadata.ViewMetadata).GetName)
}

// ListMaterializedViewNames returns the names of the schema's materialized
// views, sorted.
func (s *SchemaLookup) ListMaterializedViewNames() []string {
	return sortedNames(s.materializedViews, (*metadata.MaterializedViewMetadata).GetName)
}

// ListExternalTableNames returns the names of the schema's external tables,
// sorted.
func (s *SchemaLookup) ListExternalTableNames() []string {
	return sortedNames(s.externalTables, (*metadata.ExternalTableMetadata).GetName)
}

func sortedNames[T any](items map[string]T, nameOf func(T) string) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, nameOf(item))
	}
	slices.Sort(names)
	return names
}

func fold(name string, caseSensitive bool) string {
	if caseSensitive {
		return name
	}
	return strings.ToLower(name)
}

func sameName(a, b string, caseSensitive bool) bool {
	if caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}
