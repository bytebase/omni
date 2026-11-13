package v1review

import (
	"slices"

	"github.com/bytebase/omni/metadata"
)

// The codes of the V1 table rules. Each rule function documents which it
// uses and what its findings' fields hold.
const (
	// TableDisallowed: a statement changes a table the rule lists, by DDL or
	// by DML as the rule says. Name is the table as printed.
	TableDisallowed Code = "table-disallowed"
	// TablePartitioned: a statement creates or changes a table's
	// partitioning.
	TablePartitioned Code = "table-partitioned"
	// TableCharsetSet: a statement sets a table's character set.
	TableCharsetSet Code = "table-charset-set"
	// TableTriggerCreated: a statement creates a trigger.
	TableTriggerCreated Code = "table-trigger-created"
	// TableForeignKey: a statement creates a foreign key.
	TableForeignKey Code = "table-foreign-key"
	// TableCharsetMissing: a CREATE TABLE names no character set.
	TableCharsetMissing Code = "table-charset-missing"
	// TableCollationMissing: a CREATE TABLE names no collation.
	TableCollationMissing Code = "table-collation-missing"
	// TableDuplicateIndex: a table would have an index that duplicates
	// another.
	TableDuplicateIndex Code = "table-duplicate-index"
	// TablePrimaryKeyMissing: a table would have no primary key.
	TablePrimaryKeyMissing Code = "table-primary-key-missing"
	// TableTooLarge: a statement changes a table with more rows than the
	// rule allows.
	TableTooLarge Code = "table-too-large"
	// TableTextTooLong: a table's text columns together are longer than the
	// rule allows.
	TableTextTooLong Code = "table-text-too-long"
	// TableCommentMissing: a table has no comment where the rule requires
	// one.
	TableCommentMissing Code = "table-comment-missing"
	// TableCommentTooLong: a table's comment is longer than the rule allows.
	TableCommentTooLong Code = "table-comment-too-long"
)

// TableCommentConfig configures the table comment rule; Bytebase builds it
// from the V1 CommentConventionRulePayload. How MaxLength is read (its unit,
// and which values mean no limit) is each engine's own, as V1's was.
type TableCommentConfig struct {
	// Required requires a comment.
	Required bool
	// MaxLength is the longest allowed comment.
	MaxLength int
}

// HasTable reports whether the schema holds a table of that name or a
// partition of that name.
func (s Schema) HasTable(schemaName, tableName string) bool {
	table, _ := s.table(schemaName, tableName)
	return table != nil
}

// TableColumns returns the columns of the table of that name, in the
// snapshot's order, and whether the table was found. A partition's name
// gives the columns of the table that holds it.
func (s Schema) TableColumns(schemaName, tableName string) ([]*metadata.ColumnMetadata, bool) {
	table, _ := s.table(schemaName, tableName)
	return table.GetColumns(), table != nil
}

// TablePrimaryKey returns the primary key of the table of that name, nil
// when the table is missing or has none. A partition's name finds no
// primary key. Indexes are kept one per name, compared as indexCaseSensitive
// says; where more than one of those is primary, the V1 advisors took
// whichever their map happened to yield, and this takes the first.
func (s Schema) TablePrimaryKey(schemaName, tableName string, indexCaseSensitive bool) *metadata.IndexMetadata {
	table, partition := s.table(schemaName, tableName)
	if table == nil || partition {
		return nil
	}
	for _, index := range lastOfEachName(table.GetIndexes(), indexCaseSensitive, indexNameOf) {
		if index.GetPrimary() {
			return index
		}
	}
	return nil
}

// SchemaNames returns the names of the snapshot's schemas, one per name,
// sorted.
func (s Schema) SchemaNames() []string {
	var names []string
	for _, schema := range lastOfEachName(s.Metadata.GetSchemas(), s.ObjectCaseSensitive, schemaNameOf) {
		names = append(names, schema.GetName())
	}
	slices.Sort(names)
	return names
}

// TableNames returns the names of the schema's tables, one per name, sorted.
// Partitions are not listed.
func (s Schema) TableNames(schemaName string) []string {
	var names []string
	for _, table := range s.tables(schemaName) {
		names = append(names, table.GetName())
	}
	slices.Sort(names)
	return names
}

// table returns the table of that name, else the table holding the
// partition of that name, with partition true. Among partitions of one name
// the last wins, in the order the snapshot lists tables, each partition
// before its subpartitions; tables replaced by a later one of the same name
// still lend their partitions.
func (s Schema) table(schemaName, tableName string) (*metadata.TableMetadata, bool) {
	if table := lastNamed(s.tables(schemaName), tableName, s.ObjectCaseSensitive, tableNameOf); table != nil {
		return table, false
	}
	var owner *metadata.TableMetadata
	id := fold(tableName, s.ObjectCaseSensitive)
	var walk func(*metadata.TableMetadata, []*metadata.TablePartitionMetadata)
	walk = func(table *metadata.TableMetadata, partitions []*metadata.TablePartitionMetadata) {
		for _, partition := range partitions {
			if fold(partition.GetName(), s.ObjectCaseSensitive) == id {
				owner = table
			}
			walk(table, partition.GetSubpartitions())
		}
	}
	for _, table := range s.schema(schemaName).GetTables() {
		walk(table, table.GetPartitions())
	}
	return owner, owner != nil
}
