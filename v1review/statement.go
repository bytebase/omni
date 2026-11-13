package v1review

import "github.com/bytebase/omni/metadata"

// The codes of the V1 statement rules (STATEMENT_*). A statement rule judges
// a statement's form, not a name, so its findings usually carry only a code
// and an anchor; Name, Table, Detail and Count are set where the V1 advisor
// printed a fact the statement held, and each code says which.
const (
	// Commit: a COMMIT statement. STATEMENT_DISALLOW_COMMIT.
	Commit Code = "commit"
	// SelectAll: a select list holds *. STATEMENT_SELECT_NO_SELECT_ALL.
	SelectAll Code = "select-all"
	// WhereMissing: a statement has no WHERE clause.
	// STATEMENT_WHERE_REQUIRE_SELECT and STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
	// Name is the statement's kind, "SELECT", "UPDATE" or "DELETE", where the
	// V1 advisor told kinds apart.
	WhereMissing Code = "where-missing"
	// DMLLimit: an INSERT, UPDATE or DELETE has a LIMIT clause.
	// STATEMENT_DISALLOW_LIMIT. Name is the statement's kind, "INSERT",
	// "UPDATE" or "DELETE".
	DMLLimit Code = "dml-limit"
	// DMLOrderBy: an UPDATE or DELETE has an ORDER BY clause.
	// STATEMENT_DISALLOW_ORDER_BY. Name is the statement's kind, "UPDATE" or
	// "DELETE".
	DMLOrderBy Code = "dml-order-by"
	// LeadingWildcardLike: a LIKE pattern starts with %.
	// STATEMENT_WHERE_NO_LEADING_WILDCARD_LIKE.
	LeadingWildcardLike Code = "leading-wildcard-like"
	// InsertColumnsUnspecified: an INSERT names no column list.
	// STATEMENT_INSERT_MUST_SPECIFY_COLUMN.
	InsertColumnsUnspecified Code = "insert-columns-unspecified"
	// InsertOrderByRand: an INSERT ... SELECT orders by a random function.
	// STATEMENT_INSERT_DISALLOW_ORDER_BY_RAND.
	InsertOrderByRand Code = "insert-order-by-rand"
	// Truncate: a TRUNCATE TABLE. STATEMENT_DISALLOW_TRUNCATE. Table is the
	// truncated table as the V1 advisor printed it. Where the V1 advisor
	// printed the statement's options, Detail lists the ones written,
	// comma-separated in a fixed order ("CASCADE,RESTART IDENTITY").
	Truncate Code = "truncate"
	// TruncatePartition: an ALTER TABLE that truncates partitions.
	// STATEMENT_DISALLOW_TRUNCATE. Table is the table and Name the partition
	// as printed; Detail is the partition keyword where the V1 advisor
	// printed one.
	TruncatePartition Code = "truncate-partition"
	// RedundantAlterTable: several statements alter one table.
	// STATEMENT_MERGE_ALTER_TABLE. Table is the table as printed, Count the
	// number of statements, and Statement the last of them.
	RedundantAlterTable Code = "redundant-alter-table"
	// OptionMissing: an ALTER TABLE lacks a required option.
	// STATEMENT_REQUIRE_ALGORITHM_OPTION and STATEMENT_REQUIRE_LOCK_OPTION.
	// Name is the option, "ALGORITHM" or "LOCK".
	OptionMissing Code = "option-missing"
	// MaxExecutionTimeUnset: a statement runs before the session sets the
	// execution time limit. STATEMENT_MAX_EXECUTION_TIME. Name is the
	// variable as printed.
	MaxExecutionTimeUnset Code = "max-execution-time-unset"
	// ColumnPositioned: an added column is given a position.
	// STATEMENT_ADD_COLUMN_WITHOUT_POSITION. Name is the position keyword,
	// "FIRST" or "AFTER".
	ColumnPositioned Code = "column-positioned"
	// CheckValidated: a CHECK constraint is added without NOT VALID.
	// STATEMENT_ADD_CHECK_NOT_VALID.
	CheckValidated Code = "check-validated"
	// ForeignKeyValidated: a foreign key is added without NOT VALID.
	// STATEMENT_ADD_FOREIGN_KEY_NOT_VALID.
	ForeignKeyValidated Code = "foreign-key-validated"
	// ColumnWithDefault: an added column has a DEFAULT.
	// STATEMENT_DISALLOW_ADD_COLUMN_WITH_DEFAULT.
	ColumnWithDefault Code = "column-with-default"
	// NotNullSet: a column is set NOT NULL. STATEMENT_DISALLOW_ADD_NOT_NULL.
	// Name is the column as printed.
	NotNullSet Code = "not-null-set"
	// OnDeleteCascade: a foreign key deletes ON DELETE CASCADE.
	// STATEMENT_DISALLOW_ON_DEL_CASCADE.
	OnDeleteCascade Code = "on-delete-cascade"
	// DropCascade: a table is dropped or truncated with CASCADE.
	// STATEMENT_DISALLOW_RM_TBL_CASCADE.
	DropCascade Code = "drop-cascade"
	// NonTransactional: a statement cannot run in a transaction.
	// STATEMENT_NON_TRANSACTIONAL.
	NonTransactional Code = "non-transactional"
	// SchemaUnspecified: an object is created without its schema.
	// STATEMENT_CREATE_SPECIFY_SCHEMA. Name is the object's kind as printed.
	SchemaUnspecified Code = "schema-unspecified"
	// SetRoleMissing: no statement sets the role before the first statement
	// that is not a SET. STATEMENT_CHECK_SET_ROLE_VARIABLE.
	SetRoleMissing Code = "set-role-missing"
	// CrossDatabaseQuery: a statement names a table in another database.
	// STATEMENT_DISALLOW_CROSS_DB_QUERIES. Name is the database the table is
	// in, and Detail the current database, as printed.
	CrossDatabaseQuery Code = "cross-database-query"
	// LogicalOperatorsExceeded: a WHERE clause holds more operands of one
	// logical operator than allowed.
	// STATEMENT_WHERE_MAXIMUM_LOGICAL_OPERATOR_COUNT. Name is the operator,
	// "IN" or "OR", and Count the operands.
	LogicalOperatorsExceeded Code = "logical-operators-exceeded"
	// EqualNull: a WHERE clause compares with NULL using = or !=.
	// STATEMENT_WHERE_NO_EQUAL_NULL.
	EqualNull Code = "equal-null"
	// FunctionOnIndexedColumn: a WHERE clause applies a function to an
	// indexed column. STATEMENT_WHERE_DISALLOW_FUNCTIONS_AND_CALCULATIONS.
	// Name is the column and Detail the function, as printed.
	FunctionOnIndexedColumn Code = "function-on-indexed-column"
	// CalculationOnIndexedColumn: a WHERE clause computes with an indexed
	// column. STATEMENT_WHERE_DISALLOW_FUNCTIONS_AND_CALCULATIONS. Name is the
	// column as printed.
	CalculationOnIndexedColumn Code = "calculation-on-indexed-column"
)

// The codes of the statement rules that judge what Bytebase read from the
// database over a live connection. The rule functions never see the
// connection: each engine's rule says which statements it needs probed, and
// Bytebase runs the probe and passes each result in as a Probe.
const (
	// DryRunFailed: probing a statement on the database failed. Detail is
	// the error text. STATEMENT_DML_DRY_RUN, and STATEMENT_INSERT_ROW_LIMIT
	// and STATEMENT_AFFECTED_ROW_LIMIT when their EXPLAIN fails.
	DryRunFailed Code = "dry-run-failed"
	// RowCountFailed: the probe ran but its row count could not be read.
	// Detail is the error text. STATEMENT_INSERT_ROW_LIMIT and
	// STATEMENT_AFFECTED_ROW_LIMIT.
	RowCountFailed Code = "row-count-failed"
	// RowLimitExceeded: a statement inserts or affects more rows than
	// allowed. Count is the rows. STATEMENT_INSERT_ROW_LIMIT and
	// STATEMENT_AFFECTED_ROW_LIMIT.
	RowLimitExceeded Code = "row-limit-exceeded"
	// OfflineDDL: a statement runs an operation the engine performs offline.
	// STATEMENT_DISALLOW_OFFLINE_DDL. Name is the operation as printed.
	OfflineDDL Code = "offline-ddl"
	// OwnerMismatch: a statement changes an object another role owns.
	// STATEMENT_OBJECT_OWNER_CHECK. Name is the object's kind as printed,
	// Table the object's name, and Detail its owner.
	OwnerMismatch Code = "owner-mismatch"
)

// Probe is what Bytebase read from the database for one statement a rule
// asked to have probed (an EXPLAIN, a dry run, a row estimate).
type Probe struct {
	// Error is the probe's error text, empty when the probe succeeded.
	Error string
	// CountError is the error reading the probe's row count, empty when the
	// count was read or the rule reads none.
	CountError string
	// Rows is the row count the probe reported.
	Rows int64
}

// ColumnRef names a column a rule needs a fact about from the database.
type ColumnRef struct {
	Database, Table, Column string
}

// Table returns the table of that name in the schema of that name, nil when
// either is missing. A name no table has finds the table that owns a
// partition or subpartition of that name, as Bytebase's metadata lookup the
// V1 advisors used did; where several tables or partitions share a name, the
// last one wins. Schema, table and partition names follow ObjectCaseSensitive.
// STATEMENT_WHERE_DISALLOW_FUNCTIONS_AND_CALCULATIONS reads it.
func (s Schema) Table(schemaName, tableName string) *metadata.TableMetadata {
	id := fold(tableName, s.ObjectCaseSensitive)
	var table, partitionOwner *metadata.TableMetadata
	for _, t := range s.schema(schemaName).GetTables() {
		if t == nil {
			continue
		}
		if fold(t.GetName(), s.ObjectCaseSensitive) == id {
			table = t
		}
		if partitionNamed(t.GetPartitions(), id, s.ObjectCaseSensitive) {
			partitionOwner = t
		}
	}
	if table != nil {
		return table
	}
	return partitionOwner
}

// partitionNamed reports whether a partition or subpartition has the folded
// name id.
func partitionNamed(partitions []*metadata.TablePartitionMetadata, id string, caseSensitive bool) bool {
	for _, p := range partitions {
		if fold(p.GetName(), caseSensitive) == id || partitionNamed(p.GetSubpartitions(), id, caseSensitive) {
			return true
		}
	}
	return false
}

// SchemaNames returns the name of each schema in the snapshot, in the
// snapshot's order, a name that several schemas share (compared as
// ObjectCaseSensitive says) given once as the last of them spells it, as
// Bytebase's metadata listing the V1 advisors used did; that listing had no
// order, so a rule sorts what it needs sorted.
// STATEMENT_WHERE_DISALLOW_FUNCTIONS_AND_CALCULATIONS reads it.
func (s Schema) SchemaNames() []string {
	var names []string
	for _, schema := range lastOfEachName(s.Metadata.GetSchemas(), s.ObjectCaseSensitive, schemaNameOf) {
		names = append(names, schema.GetName())
	}
	return names
}

// SchemaByName returns the schema of that name, nil when the snapshot has
// none; the last of several schemas of one name (compared as
// ObjectCaseSensitive says) wins, as in Bytebase's metadata lookup the V1
// advisors used. STATEMENT_OBJECT_OWNER_CHECK reads it.
func (s Schema) SchemaByName(name string) *metadata.SchemaMetadata {
	return s.schema(name)
}
