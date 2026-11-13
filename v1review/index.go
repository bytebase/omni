package v1review

import (
	"errors"

	"github.com/bytebase/omni/metadata"
)

// The conditions the index rules report.
const (
	// IndexKeyCountExceeded: an index or key has more key parts than the
	// configured limit.
	IndexKeyCountExceeded Code = "index-key-count-exceeded"
	// IndexColumnDuplicated: an index or key lists a column twice. Other
	// holds the column.
	IndexColumnDuplicated Code = "index-column-duplicated"
	// IndexTypeNotAllowed: an index's type, held in Type, is not in the
	// configured list.
	IndexTypeNotAllowed Code = "index-type-not-allowed"
	// IndexCreatedNotConcurrently: an index is created in a way that
	// blocks writes to its table.
	IndexCreatedNotConcurrently Code = "index-created-not-concurrently"
	// IndexDroppedNotConcurrently: an index is dropped in a way that
	// blocks writes to its table.
	IndexDroppedNotConcurrently Code = "index-dropped-not-concurrently"
	// PrimaryKeyTypeNotAllowed: a primary key column, held in Name, has a
	// type, held in Type, the rule does not allow.
	PrimaryKeyTypeNotAllowed Code = "primary-key-type-not-allowed"
	// IndexColumnTypeNotAllowed: an indexed column, held in Name, has a
	// type, held in Type, that the rule keeps out of indexes.
	IndexColumnTypeNotAllowed Code = "index-column-type-not-allowed"
	// IndexCountExceeded: a table has more indexes, counted in Count, than
	// the configured limit.
	IndexCountExceeded Code = "index-count-exceeded"
	// IndexRedundant: a new index, held in Name, repeats what an existing
	// index, held in Other, already covers.
	IndexRedundant Code = "index-redundant"
)

// ErrSchemaMissing is the error of a rule function that had to read a Schema
// whose Metadata is nil. The V1 advisors read their snapshot unguarded and
// failed at that point, so a rule function returns this, and no findings,
// exactly where its advisor would have failed; a script that never sends the
// rule to the schema is judged as usual.
var ErrSchemaMissing = errors.New("v1review: the rule read a schema that holds no snapshot")

// Table returns the table of that name, nil when the schema or the table is
// missing. Where no table has the name, the name of a partition or
// subpartition finds the table it divides, the last such table in the
// snapshot's order.
func (s Schema) Table(schemaName, tableName string) *metadata.TableMetadata {
	tables := s.schema(schemaName).GetTables()
	if table := lastNamed(tables, tableName, s.ObjectCaseSensitive, tableNameOf); table != nil {
		return table
	}
	var found *metadata.TableMetadata
	for _, table := range tables {
		if hasPartition(table.GetPartitions(), fold(tableName, s.ObjectCaseSensitive), s.ObjectCaseSensitive) {
			found = table
		}
	}
	return found
}

// TableColumn returns the column of that name on the table Table finds, nil
// when either is missing; the last column wins among equals. Whether column
// names are case sensitive is the engine's own fact.
func (s Schema) TableColumn(schemaName, tableName, columnName string, columnCaseSensitive bool) *metadata.ColumnMetadata {
	return lastNamed(s.Table(schemaName, tableName).GetColumns(), columnName, columnCaseSensitive, columnNameOf)
}

func columnNameOf(c *metadata.ColumnMetadata) string { return c.GetName() }

// hasPartition reports whether a partition or subpartition has that name,
// already folded.
func hasPartition(partitions []*metadata.TablePartitionMetadata, id string, caseSensitive bool) bool {
	for _, partition := range partitions {
		if fold(partition.GetName(), caseSensitive) == id || hasPartition(partition.GetSubpartitions(), id, caseSensitive) {
			return true
		}
	}
	return false
}

// The kinds of key a finding's Type names, in the words the V1 advisors
// printed. A rule function sets one where its advice says what kind of key
// broke the rule, or where its caller must tell the kinds apart.
const (
	KeyPrimary     = "PRIMARY KEY"
	KeyUnique      = "UNIQUE KEY"
	KeyForeign     = "FOREIGN KEY"
	KeyIndex       = "INDEX"
	KeyUniqueIndex = "UNIQUE INDEX"
)
