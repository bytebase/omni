package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// QueryType classifies a Redshift statement. The values mirror Bytebase's
// base.QueryType.
type QueryType int

const (
	QueryTypeUnknown QueryType = iota
	Select
	Explain
	SelectInfoSchema
	DDL
	DML
)

// ColumnSet is a set of physical column resources.
type ColumnSet map[ColumnResource]bool

// ResolvedQuerySpanResult describes one output column of a resolved span.
type ResolvedQuerySpanResult struct {
	Name          string
	SourceColumns ColumnSet
	IsPlainField  bool
}

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

// FunctionNotSupportedError reports a user function whose body cannot be
// analyzed for lineage.
type FunctionNotSupportedError struct {
	Err        error
	Function   string
	Definition string
}

func (e *FunctionNotSupportedError) Error() string {
	parts := []string{fmt.Sprintf("function is not supported with data masking: %s", e.Function)}
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("err: %s", e.Err.Error()))
	}
	if e.Definition != "" {
		parts = append(parts, fmt.Sprintf("definition: %s", e.Definition))
	}
	return strings.Join(parts, ", ")
}

func (e *FunctionNotSupportedError) Unwrap() error { return e.Err }

// ErrMixUserSystemTables reports a statement that reads user and system tables
// together.
var ErrMixUserSystemTables = fmt.Errorf("cannot access user and system tables at the same time")

// ResolveQuerySpan resolves a single Redshift statement against the metadata the
// caller supplies and returns the physical columns it reads, the predicate
// columns, its output columns and lineage, and its statement classification.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedQuerySpan, error) {
	return newQuerySpanResolver(opts.DefaultDatabase, opts.SearchPath, opts).getQuerySpan(statement)
}
