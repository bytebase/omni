package v1review

// The conditions a prior-backup check reports (BUILTIN_PRIOR_BACKUP_CHECK).
// The judgment reads two facts Bytebase resolves: the size of the reviewed
// SQL against the size a backup may take, and whether the database a backup
// writes to exists. Both arrive in a PriorBackupFacts, so no connection
// crosses the contract.
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
	// PriorBackupDatabaseMissing: the database a prior backup writes to does
	// not exist. Name holds it as the V1 advisor printed it.
	PriorBackupDatabaseMissing Code = "prior-backup-database-missing"
	// PriorBackupMixedDMLOnTable: the reviewed statements change one table
	// with more than one kind of DML. Table holds the table as
	// "database.table", the text the V1 advisor printed.
	PriorBackupMixedDMLOnTable Code = "prior-backup-mixed-dml-on-table"
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
