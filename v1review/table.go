package v1review

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
