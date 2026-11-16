package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// ResolveOptions carries the facts the caller resolves from its metadata store.
// GetDatabaseMetadata returns nil for a database that is not synced.
type ResolveOptions struct {
	DefaultDatabase     string
	DefaultSchema       string
	IgnoreCaseSensitive bool
	ListDatabases       func() ([]string, error)
	GetDatabaseMetadata func(database string) (*metadata.DatabaseSchemaMetadata, error)
}

// Column identifies a physical column resource. Database is a Trino catalog.
type Column struct {
	Database string
	Schema   string
	Table    string
	Column   string
}

// ColumnSet is a set of physical column resources.
type ColumnSet map[Column]bool

// ResolvedQuerySpanResult describes one output column of a resolved span.
type ResolvedQuerySpanResult struct {
	Name          string
	SourceColumns ColumnSet
	IsPlainField  bool
}

// ResolvedQuerySpan is a query span resolved against metadata.
type ResolvedQuerySpan struct {
	Type             QueryType
	SourceColumns    ColumnSet
	PredicateColumns ColumnSet
	Results          []ResolvedQuerySpanResult
	NotFoundError    error
}

// ResourceNotFoundError reports a resource missing from the synced metadata.
type ResourceNotFoundError struct {
	Err      error
	Database *string
	Schema   *string
	Table    *string
	Column   *string
}

func (e *ResourceNotFoundError) Error() string {
	var parts []string
	if e.Database != nil {
		parts = append(parts, fmt.Sprintf("database: %s", *e.Database))
	}
	if e.Schema != nil {
		parts = append(parts, fmt.Sprintf("schema: %s", *e.Schema))
	}
	if e.Table != nil {
		parts = append(parts, fmt.Sprintf("table: %s", *e.Table))
	}
	if e.Column != nil {
		parts = append(parts, fmt.Sprintf("column: %s", *e.Column))
	}
	result := fmt.Sprintf("resource not found: %s", strings.Join(parts, ", "))
	if e.Err != nil {
		result += fmt.Sprintf(", err: %s", e.Err.Error())
	}
	return result
}

func (e *ResourceNotFoundError) Unwrap() error { return e.Err }

// ResolveQuerySpan resolves a single Trino statement against the metadata the
// caller supplies and returns the physical columns it reads, the predicate
// columns, its output columns and lineage, and its statement classification.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedQuerySpan, error) {
	return newResolver(opts.DefaultDatabase, opts.DefaultSchema, opts).getQuerySpan(statement)
}
