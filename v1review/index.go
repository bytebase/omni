package v1review

// The conditions the index rules report.
const (
	// IndexKeyCountExceeded: an index or key has more key parts than the
	// configured limit.
	IndexKeyCountExceeded Code = "index-key-count-exceeded"
	// IndexColumnDuplicated: an index or key lists a column twice. Other
	// holds the column.
	IndexColumnDuplicated Code = "index-column-duplicated"
	// IndexTypeNotAllowed: an index's type, held in Type, is not in the
	// configured list.
	IndexTypeNotAllowed Code = "index-type-not-allowed"
	// IndexCreatedNotConcurrently: an index is created in a way that
	// blocks writes to its table.
	IndexCreatedNotConcurrently Code = "index-created-not-concurrently"
	// IndexDroppedNotConcurrently: an index is dropped in a way that
	// blocks writes to its table.
	IndexDroppedNotConcurrently Code = "index-dropped-not-concurrently"
	// PrimaryKeyTypeNotAllowed: a primary key column, held in Name, has a
	// type, held in Type, the rule does not allow.
	PrimaryKeyTypeNotAllowed Code = "primary-key-type-not-allowed"
)

// The kinds of key a finding's Type names, in the words the V1 advisors
// printed. A rule function sets one where its advice says what kind of key
// broke the rule, or where its caller must tell the kinds apart.
const (
	KeyPrimary     = "PRIMARY KEY"
	KeyUnique      = "UNIQUE KEY"
	KeyForeign     = "FOREIGN KEY"
	KeyIndex       = "INDEX"
	KeyUniqueIndex = "UNIQUE INDEX"
)
