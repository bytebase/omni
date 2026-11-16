package v1review

// The condition the online migration advice reports (ADVICE_ONLINE_MIGRATION).
const (
	// OnlineMigrationNeeded: a statement changes a table whose estimated row
	// count reaches the configured limit, so the change may hold a lock or
	// copy the table. Table holds the table as "database.table", the text the
	// V1 advisor printed, and Value the row count it read.
	OnlineMigrationNeeded Code = "online-migration-needed"
)
