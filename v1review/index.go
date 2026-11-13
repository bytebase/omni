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
	// IndexTypeNotAllowed: an index's type, held in Value, is not in the
	// configured list.
	IndexTypeNotAllowed Code = "index-type-not-allowed"
	// IndexCreatedNotConcurrently: an index is created in a way that
	// blocks writes to its table.
	IndexCreatedNotConcurrently Code = "index-created-not-concurrently"
	// IndexDroppedNotConcurrently: an index is dropped in a way that
	// blocks writes to its table.
	IndexDroppedNotConcurrently Code = "index-dropped-not-concurrently"
	// PrimaryKeyTypeNotAllowed: a primary key column, held in Name, has a
	// type, held in Value, the rule does not allow.
	PrimaryKeyTypeNotAllowed Code = "primary-key-type-not-allowed"
	// IndexColumnTypeNotAllowed: an indexed column, held in Name, has a
	// type, held in Value, that the rule keeps out of indexes.
	IndexColumnTypeNotAllowed Code = "index-column-type-not-allowed"
	// IndexCountExceeded: a table has more indexes, counted in Value, than
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

// The kinds of key a finding's Value names, in the words the V1 advisors
// printed. A rule function sets one where its advice says what kind of key
// broke the rule, or where its caller must tell the kinds apart.
const (
	KeyPrimary     = "PRIMARY KEY"
	KeyUnique      = "UNIQUE KEY"
	KeyForeign     = "FOREIGN KEY"
	KeyIndex       = "INDEX"
	KeyUniqueIndex = "UNIQUE INDEX"
)

// TableIndexes returns the indexes of the table of that name, in the
// snapshot's order, and whether the table was found. A partition's name
// gives the indexes of the table that holds it.
func (s Schema) TableIndexes(schemaName, tableName string) ([]*metadata.IndexMetadata, bool) {
	table, _ := s.table(schemaName, tableName)
	return table.GetIndexes(), table != nil
}

// TableColumn returns the column of that name on the table TableColumns
// finds, nil when either is missing; the last column wins among equals.
// Whether column names are case sensitive is the engine's own fact.
func (s Schema) TableColumn(schemaName, tableName, columnName string, columnCaseSensitive bool) *metadata.ColumnMetadata {
	columns, _ := s.TableColumns(schemaName, tableName)
	return lastNamed(columns, columnName, columnCaseSensitive, columnNameOf)
}

func columnNameOf(c *metadata.ColumnMetadata) string { return c.GetName() }
