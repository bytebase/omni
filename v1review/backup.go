package v1review

// The conditions a prior-backup check reports (BUILTIN_PRIOR_BACKUP_CHECK).
// The judgment reads two facts Bytebase resolves: the size of the reviewed
// SQL against the size a backup may take, and whether the database a backup
// writes to exists. Both arrive in a PriorBackupFacts, so no connection
// crosses the contract. Oracle's LONG/LONG RAW warning reads the pre-change
// snapshot instead, through the v1review.Schema every rule that looks
// something up already receives. A condition an engine judges from SQL alone
// gets a code here; none of them needs another fact.
const (
	// PriorBackupScriptTooLarge: the reviewed SQL is larger than a prior
	// backup can take. Value holds the limit in bytes, as the V1 advisor
	// printed it.
	PriorBackupScriptTooLarge Code = "prior-backup-script-too-large"
	// PriorBackupMixedDDLAndDML: a statement in the reviewed batch is DDL,
	// which a prior backup cannot take. The V1 advisor reported every such
	// statement, whether or not the batch changes rows as well. Statement is
	// the index of the DDL statement; the V1 advisor anchored each at the
	// statement's start.
	PriorBackupMixedDDLAndDML Code = "prior-backup-mixed-ddl-and-dml"
	// PriorBackupDDLInBatch: the reviewed batch holds a DDL statement, which a
	// prior backup cannot take. The Oracle and MSSQL V1 advisors reported this
	// once for the batch wherever the DDL is, and anchored it at the script's
	// start rather than the statement's, so the finding carries no statement
	// index and no range; MySQL's per-statement code is
	// PriorBackupMixedDDLAndDML.
	PriorBackupDDLInBatch Code = "prior-backup-ddl-in-batch"
	// PriorBackupDatabaseMissing: the database a prior backup writes to does
	// not exist. Name holds it as the V1 advisor printed it.
	PriorBackupDatabaseMissing Code = "prior-backup-database-missing"
	// PriorBackupMixedDMLOnTable: the reviewed statements change one table
	// with more than one kind of DML. Table holds the table as the V1 advisor
	// printed it, qualified by every name it resolved: "database.table" for
	// MySQL, "schema.table" for Oracle and PostgreSQL, "database.schema.table"
	// for MSSQL. TiDB printed MySQL's qualifiers lower-cased.
	PriorBackupMixedDMLOnTable Code = "prior-backup-mixed-dml-on-table"
	// PriorBackupUncopyableColumn: a table the reviewed statements change has
	// a column a backup cannot copy: Oracle's LONG and LONG RAW, which CREATE
	// TABLE AS SELECT rejects (ORA-00997). Schema and Table hold the table as
	// the V1 advisor printed them, and Value the column names as it printed
	// them, each double-quoted and separated by ", ".
	PriorBackupUncopyableColumn Code = "prior-backup-uncopyable-column"
	// PriorBackupDataModifyingCTE: a statement holds a WITH query that changes
	// data, so a prior backup cannot save the rows it changes. Value holds the
	// CTE names as the PostgreSQL V1 advisor printed them, each double-quoted
	// and separated by ", ". The V1 advisor worded the advice differently for
	// an UPDATE or DELETE, which the backup task rejects, and for any other
	// statement, whose rows it merely does not cover; the caller reads that
	// from the statement.
	PriorBackupDataModifyingCTE Code = "prior-backup-data-modifying-cte"
	// PriorBackupTooManyDMLStatements: the reviewed batch changes rows in more
	// statements than a prior backup groups into one backup. Value holds the
	// limit, as the V1 advisor printed it. TiDB's backup transformer routes
	// that many statements into a single-table backup, so a batch above the
	// limit that spans more than one table is rejected at run time.
	PriorBackupTooManyDMLStatements Code = "prior-backup-too-many-dml-statements"
)

// PriorBackupFacts is what Bytebase resolves for a prior-backup check: what
// the reviewed SQL weighs and whether a backup can be written. A rule function
// compares the first pair and reports the second; it reads no connection.
type PriorBackupFacts struct {
	// TotalSize is the size of the reviewed SQL in bytes.
	TotalSize int
	// MaxSize is the largest reviewed SQL a prior backup takes, in bytes.
	MaxSize int
	// Database is the database a prior backup writes to.
	Database string
	// DatabaseExists is whether Database exists on the instance.
	DatabaseExists bool
}
