package analysis

// ChangedResource is a database object a statement writes.
type ChangedResource struct {
	Database string
	Schema   string
	// Table is the changed table, empty for a database-only change or when
	// TableID identifies the table instead.
	Table string
	// DatabaseOnly marks a non-table object (such as a view or routine) whose
	// name carries a database qualifier and has no table to scope.
	DatabaseOnly bool
	// AffectsData marks a change that touches existing rows or the table's
	// definition, such as ALTER, DROP, or TRUNCATE.
	AffectsData bool
	// TableID is an engine-internal table identifier (for example, Doris
	// DORIS_INTERNAL_TABLE_ID) when the statement does not name a table.
	TableID int64
}

// StatementChanges summarizes the objects a statement or script writes and its
// DML counts.
type StatementChanges struct {
	// Resources holds each written resource in statement order, duplicates
	// included.
	Resources []ChangedResource
	// DMLStatements holds the DML statements eligible for affected-row
	// estimation.
	DMLStatements []string
	// DMLCount counts all DML statements, including those omitted from
	// DMLStatements.
	DMLCount int
	// InsertCount counts rows written by INSERT ... VALUES statements.
	InsertCount int
}

// ByteRange is a half-open byte range [Start, End) in a SQL script.
type ByteRange struct {
	Start int
	End   int
}

// Statement is one segment of a split SQL script.
type Statement struct {
	// Text is the segment's source text, typically script[Start:End].
	Text string
	// Empty reports whether the segment holds no SQL statement (whitespace,
	// comments, or a lone delimiter).
	Empty bool
	// Start and End are the half-open byte offsets [Start, End) in the script.
	Start int
	End   int
}
