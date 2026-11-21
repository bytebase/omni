package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

// QueryType classifies a Redshift statement. The values mirror Bytebase's
// base.QueryType.
type QueryType = analysis.QueryType

const (
	QueryTypeUnknown = analysis.QueryTypeUnknown
	Select           = analysis.Select
	Explain          = analysis.Explain
	SelectInfoSchema = analysis.SelectInfoSchema
	DDL              = analysis.DDL
	DML              = analysis.DML
)

// ColumnSet is a set of physical column resources.
type ColumnSet = analysis.ColumnSet

// ResolvedQuerySpanResult describes one output column of a resolved span.
type ResolvedQuerySpanResult = analysis.QuerySpanResult

// ResolvedQuerySpan is a query span resolved against metadata.
type ResolvedQuerySpan struct {
	Type                      QueryType
	SourceColumns             ColumnSet
	PredicateColumns          ColumnSet
	Results                   []ResolvedQuerySpanResult
	NotFoundError             error
	FunctionNotSupportedError error
}

// ResolveOptions carries the facts the caller resolves from its metadata store.
// SearchPath is the session search path, already resolved by the caller; an
// empty path means "public". GetDatabaseMetadata returns nil for a database
// that is not synced.
type ResolveOptions struct {
	DefaultDatabase     string
	SearchPath          []string
	GetDatabaseMetadata func(database string) (*metadata.DatabaseSchemaMetadata, error)
}

// ResourceNotFoundError reports a resource missing from the synced metadata.
type ResourceNotFoundError = analysis.ResourceNotFoundError

// FunctionNotSupportedError reports a user function whose body cannot be
// analyzed for lineage.
type FunctionNotSupportedError = analysis.FunctionNotSupportedError

// ErrMixUserSystemTables reports a statement that reads user and system tables
// together.
var ErrMixUserSystemTables = analysis.ErrMixUserSystemTables

// ResolveQuerySpan resolves a single Redshift statement against the metadata the
// caller supplies and returns the physical columns it reads, the predicate
// columns, its output columns and lineage, and its statement classification.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedQuerySpan, error) {
	return newQuerySpanResolver(opts.DefaultDatabase, opts.SearchPath, opts).getQuerySpan(statement)
}
