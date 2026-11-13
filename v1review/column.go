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

// TableIndexesByName returns the indexes on the table of that name, one per
// name, a later index replacing an earlier one of the same name where the
// first stood: what Bytebase's GetTable(t).ListIndexes() holds. A
// partition's name finds the partition, which has no indexes of its own, so
// it returns nil, as it does for a missing schema or table; TableIndexes
// gives a partition its table's indexes instead.
func (s Schema) TableIndexesByName(schemaName, tableName string, indexCaseSensitive bool) []*metadata.IndexMetadata {
	table, partition := s.table(schemaName, tableName)
	if table == nil || partition {
		return nil
	}
	return lastOfEachName(table.GetIndexes(), indexCaseSensitive, indexNameOf)
}
