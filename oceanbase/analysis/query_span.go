package analysis

// ResolveOptions configures ResolveQuerySpan.
type ResolveOptions struct {
	// DefaultDatabase qualifies unqualified table names.
	DefaultDatabase string
	// IgnoreCaseSensitive matches database, table and view names
	// case-insensitively, as on an instance with lower_case_table_names set.
	IgnoreCaseSensitive bool
	Catalog             Catalog
}

// ResolveQuerySpan returns the column lineage of one OceanBase statement,
// resolving the tables and views it names through opts.Catalog.
//
// A parse error is the parser's error, unwrapped. A statement that reads user
// and system tables together returns ErrMixUserSystemTables. A table or column
// missing from the catalog is not an error: the span carries it in
// NotFoundError. Statements that are not queries return their QueryType and
// the tables they access, with no Results.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedSpan, error) {
	return newQuerySpanExtractor(opts.DefaultDatabase, opts.Catalog, opts.IgnoreCaseSensitive).getQuerySpan(statement)
}
