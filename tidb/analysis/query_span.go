package analysis

// ResolveOptions configures ResolveQuerySpan.
type ResolveOptions struct {
	// DefaultDatabase qualifies unqualified table names.
	DefaultDatabase string
	Catalog         Catalog
}

// ResolveQuerySpan returns the column lineage of one TiDB statement,
// resolving the tables and views it names through opts.Catalog, which it
// reads with GetDatabase only.
//
// A statement that reads user and system tables together returns
// ErrMixUserSystemTables. A table or column missing from the catalog is not
// an error: the span carries it in NotFoundError. Statements that are not
// queries return their QueryType and the tables they access, with no Results.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedSpan, error) {
	return newQuerySpanExtractor(opts.DefaultDatabase, opts.Catalog).getQuerySpan(statement)
}
