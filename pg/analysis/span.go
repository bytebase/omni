// Package analysis holds Bytebase's query analysis for PostgreSQL: query span
// (column lineage, predicate columns and the tables a statement reads) and the
// system object lists it relies on.
package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// QueryType is the permission class of a statement. The values mirror
// Bytebase's base.QueryType.
type QueryType int

const (
	// QueryTypeUnknown is a statement not deliberately mapped. Bytebase's SQL
	// service denies it before execution.
	QueryTypeUnknown QueryType = iota
	// Select is a read-only query.
	Select
	// Explain is an EXPLAIN statement.
	Explain
	// SelectInfoSchema reads system metadata.
	SelectInfoSchema
	// DDL changes schema.
	DDL
	// DML changes table data.
	DML
)

// ColumnResource identifies a column, or a whole table when Column is empty.
type ColumnResource struct {
	Database string
	Schema   string
	Table    string
	Column   string
}

// String formats the resource as Bytebase's base.ColumnResource does.
func (c ColumnResource) String() string {
	var parts []string
	for _, part := range []string{c.Database, c.Schema, c.Table, c.Column} {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}

// ColumnSet is a set of columns.
type ColumnSet map[ColumnResource]bool

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult struct {
	Name          string
	SourceColumns ColumnSet
	// IsPlainField is true when the result is a bare column reference rather
	// than an expression.
	IsPlainField bool
}

// ResolvedSpan is the lineage ResolveQuerySpan returns.
type ResolvedSpan struct {
	Type QueryType
	// ExplainAnalyze marks an EXPLAIN ANALYZE, which runs the statement it
	// wraps; Type is that statement's.
	ExplainAnalyze bool
	// SourceColumns holds the tables the statement accesses, as table-level
	// resources (Column empty), including those read inside the bodies of the
	// functions it calls.
	SourceColumns ColumnSet
	Results       []QuerySpanResult
	// PredicateColumns holds the columns that decide which rows return.
	PredicateColumns ColumnSet
	// UnresolvedColumnsError is set when a table the statement reads has no
	// columns in the synced metadata, so its lineage cannot be trusted.
	UnresolvedColumnsError *UnresolvedColumnsError
}

// UnresolvedColumnsError identifies relations whose synced metadata has no
// columns.
type UnresolvedColumnsError struct {
	// Relations each carry an empty Column field.
	Relations []ColumnResource
}

// Error formats the error as Bytebase's base.UnresolvedColumnsError does.
func (e *UnresolvedColumnsError) Error() string {
	names := make([]string, 0, len(e.Relations))
	for _, r := range e.Relations {
		names = append(names, r.String())
	}
	slices.Sort(names)
	return fmt.Sprintf("the synced schema describes no columns for %s", strings.Join(names, ", "))
}

// SyntaxError reports a statement that does not parse. Err is the parser's
// error, a *parser.ParseError when the parser located the failure.
type SyntaxError struct {
	Err error
}

func (e *SyntaxError) Error() string { return e.Err.Error() }

func (e *SyntaxError) Unwrap() error { return e.Err }

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = errors.New("cannot access user and system tables at the same time")

// Catalog is the synced metadata of the instance a statement runs on.
type Catalog struct {
	// GetDatabase returns a database's metadata, or a nil Database when the
	// database has not been synced.
	GetDatabase func(name string) (Database, error)
}

// Database is one database's synced metadata. The lookups follow the
// metadata's own name matching, so the caller decides case sensitivity.
type Database interface {
	// Metadata is the whole snapshot, which loads the analysis catalog.
	Metadata() *metadata.DatabaseSchemaMetadata
	// GetTable returns the table named table in schema, or the table a
	// partition of that name belongs to, or nil.
	GetTable(schema, table string) *metadata.TableMetadata
	// SearchRelation returns the schema and name of the first table, view,
	// materialized view, external table or sequence named name in searchPath's
	// schemas, or two empty strings.
	SearchRelation(searchPath []string, name string) (schema string, relation string)
	// SearchFunctions returns the functions named name in searchPath's schemas,
	// in search path order.
	SearchFunctions(searchPath []string, name string) []*metadata.FunctionMetadata
}
