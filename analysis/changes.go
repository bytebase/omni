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

// DMLStatement is a DML statement whose affected rows the caller estimates.
type DMLStatement struct {
	// Text is the statement as the caller runs it.
	Text string
	// SearchPath is the search path the statement runs under when an earlier
	// statement changed it, and nil when it runs under the initial one.
	SearchPath []string
}

// Changes is what a script changes.
type Changes struct {
	// Resources holds the objects the statements write, in the order they name
	// them, duplicates included.
	Resources []ChangedResource
	// DMLStatements holds the DML statements to estimate, in order.
	DMLStatements []DMLStatement
	// DMLCount counts the DML statements, including those DMLStatements omits.
	DMLCount int
	// InsertCount counts the rows of INSERT ... VALUES statements.
	InsertCount int
}

// ChangeRequest is a script whose changes are extracted.
type ChangeRequest struct {
	DefaultDatabase string
	DefaultSchema   string
	// Database is the default database's synced metadata, or nil when it is
	// not synced.
	Database *Database
	// Script is the whole script the statements were parsed from.
	Script string
}

// ChangeExtractor reports what a script's parsed statements write.
type ChangeExtractor interface {
	ExtractChanges(request ChangeRequest, statements []ParsedStatement) (*Changes, error)
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
