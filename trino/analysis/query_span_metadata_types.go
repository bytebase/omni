package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

// ResolveOptions carries the facts the caller resolves from its metadata store.
// GetDatabaseMetadata returns nil for a database that is not synced.
type ResolveOptions struct {
	DefaultDatabase string
	DefaultSchema   string
	// IgnoreCaseSensitive is read by no Trino resolution: Trino matches every
	// catalog, schema, table and column name case-insensitively, quoted or
	// not, and so does ResolveQuerySpan.
	IgnoreCaseSensitive bool
	ListDatabases       func() ([]string, error)
	GetDatabaseMetadata func(database string) (*metadata.DatabaseSchemaMetadata, error)
}

// Column identifies a physical column resource. Database is a Trino catalog.
type Column = analysis.ColumnResource

// ColumnSet is a set of physical column resources.
type ColumnSet = analysis.ColumnSet

// ResolvedQuerySpanResult describes one output column of a resolved span.
type ResolvedQuerySpanResult = analysis.QuerySpanResult

// ResolvedQuerySpan is a query span resolved against metadata.
type ResolvedQuerySpan struct {
	Type             QueryType
	SourceColumns    ColumnSet
	PredicateColumns ColumnSet
	Results          []ResolvedQuerySpanResult
	NotFoundError    error
}

// ResourceNotFoundError reports a resource missing from the synced metadata.
type ResourceNotFoundError = analysis.ResourceNotFoundError

// ResolveQuerySpan resolves a single Trino statement against the metadata the
// caller supplies and returns the physical columns it reads, the predicate
// columns, its output columns and lineage, and its statement classification.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedQuerySpan, error) {
	return newResolver(opts.DefaultDatabase, opts.DefaultSchema, opts).getQuerySpan(statement)
}
