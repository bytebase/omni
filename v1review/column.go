package v1review

import "github.com/bytebase/omni/metadata"

// The codes the column rules report. The V1 payloads these rules read are
// plain values, so each rule function takes them as parameters: a number
// payload as an int, a string list payload as a []string, and a comment
// convention payload as whether a comment is required and its longest
// allowed length.
const (
	// ColumnMissing: a table lacks a column the rule requires. Name is the
	// missing column.
	ColumnMissing Code = "column-missing"
	// ColumnNoDefault: a column has no DEFAULT.
	ColumnNoDefault Code = "column-no-default"
	// ColumnNotNullNoDefault: a NOT NULL column has no DEFAULT.
	ColumnNotNullNoDefault Code = "column-not-null-no-default"
	// ColumnDefaultVolatile: a column's DEFAULT is volatile.
	ColumnDefaultVolatile Code = "column-default-volatile"
	// ColumnNullable: a column can hold NULL.
	ColumnNullable Code = "column-nullable"
	// ColumnTypeDisallowed: a column has a type the rule disallows.
	ColumnTypeDisallowed Code = "column-type-disallowed"
	// ColumnCommentMissing: a column has no comment.
	ColumnCommentMissing Code = "column-comment-missing"
	// ColumnCommentTooLong: a column's comment exceeds the allowed length.
	ColumnCommentTooLong Code = "column-comment-too-long"
	// ColumnCharTooLong: a CHAR column exceeds the allowed length.
	ColumnCharTooLong Code = "column-char-too-long"
	// ColumnVarcharTooLong: a VARCHAR column exceeds the allowed length.
	ColumnVarcharTooLong Code = "column-varchar-too-long"
	// ColumnAutoIncrementInitialValue: a table's initial auto-increment
	// value is not the configured one.
	ColumnAutoIncrementInitialValue Code = "column-auto-increment-initial-value"
	// ColumnAutoIncrementNotInteger: an auto-increment column is not of an
	// integer type.
	ColumnAutoIncrementNotInteger Code = "column-auto-increment-not-integer"
	// ColumnAutoIncrementNotUnsigned: an auto-increment column is not
	// UNSIGNED.
	ColumnAutoIncrementNotUnsigned Code = "column-auto-increment-not-unsigned"
	// ColumnDefaultCurrentTimeCount: a table has too many columns that
	// default to the current time.
	ColumnDefaultCurrentTimeCount Code = "column-default-current-time-count"
	// ColumnOnUpdateCurrentTimeCount: a table has too many columns that
	// update to the current time.
	ColumnOnUpdateCurrentTimeCount Code = "column-on-update-current-time-count"
	// ColumnChanged: a statement uses CHANGE COLUMN.
	ColumnChanged Code = "column-changed"
	// ColumnOrderChanged: a statement changes a column's position.
	ColumnOrderChanged Code = "column-order-changed"
	// ColumnTypeChanged: a statement changes a column's type.
	ColumnTypeChanged Code = "column-type-changed"
	// ColumnCharsetSet: a statement sets a column's character set.
	ColumnCharsetSet Code = "column-charset-set"
	// ColumnNoCharset: a column does not name its character set.
	ColumnNoCharset Code = "column-no-charset"
	// ColumnNoCollation: a column does not name its collation.
	ColumnNoCollation Code = "column-no-collation"
	// ColumnDropped: a statement drops a column.
	ColumnDropped Code = "column-dropped"
	// ColumnInIndexDropped: a statement drops a column an index uses.
	ColumnInIndexDropped Code = "column-in-index-dropped"
)

// TableColumn returns the column of that name on the table of that name, nil
// when the schema, the table, or the column is missing. A name no table has
// finds the partition of that name, whose columns are its table's; among
// partitions of the same name across the schema, subpartitions included, the
// last wins. Whether column names are case sensitive is the engine's own
// fact, so the engine's rule function passes it.
func (s Schema) TableColumn(schemaName, tableName, columnName string, columnCaseSensitive bool) *metadata.ColumnMetadata {
	table, _ := s.tableOrPartition(schemaName, tableName)
	if table == nil {
		return nil
	}
	return lastNamed(table.GetColumns(), columnName, columnCaseSensitive, columnNameOf)
}

// TableIndexes returns the indexes on the table of that name, one per name, a
// later index replacing an earlier one of the same name where the first
// stood. A partition's name finds the partition, which has no indexes of its
// own, so it returns nil, as it does for a missing schema or table.
func (s Schema) TableIndexes(schemaName, tableName string, indexCaseSensitive bool) []*metadata.IndexMetadata {
	table, partition := s.tableOrPartition(schemaName, tableName)
	if table == nil || partition {
		return nil
	}
	return lastOfEachName(table.GetIndexes(), indexCaseSensitive, indexNameOf)
}

// tableOrPartition returns the table of that name, else the table owning the
// last partition of that name and true.
func (s Schema) tableOrPartition(schemaName, tableName string) (*metadata.TableMetadata, bool) {
	if table := lastNamed(s.tables(schemaName), tableName, s.ObjectCaseSensitive, tableNameOf); table != nil {
		return table, false
	}
	var owner *metadata.TableMetadata
	id := fold(tableName, s.ObjectCaseSensitive)
	var walk func(table *metadata.TableMetadata, partitions []*metadata.TablePartitionMetadata)
	walk = func(table *metadata.TableMetadata, partitions []*metadata.TablePartitionMetadata) {
		for _, partition := range partitions {
			if fold(partition.GetName(), s.ObjectCaseSensitive) == id {
				owner = table
			}
			walk(table, partition.GetSubpartitions())
		}
	}
	for _, table := range s.schema(schemaName).GetTables() {
		if table != nil {
			walk(table, table.GetPartitions())
		}
	}
	return owner, owner != nil
}

func columnNameOf(c *metadata.ColumnMetadata) string { return c.GetName() }
