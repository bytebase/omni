// Package analysis holds Bytebase's query analysis for PostgreSQL: query span
// (column lineage, predicate columns and the tables a statement reads) and the
// system object lists it relies on.
package analysis

import (
	"cmp"
	"slices"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

// QueryType is the permission class of a statement. The values mirror
// Bytebase's base.QueryType.
type QueryType = analysis.QueryType

const (
	// QueryTypeUnknown is a statement not deliberately mapped. Bytebase's SQL
	// service denies it before execution.
	QueryTypeUnknown = analysis.QueryTypeUnknown
	// Select is a read-only query.
	Select = analysis.Select
	// Explain is an EXPLAIN statement.
	Explain = analysis.Explain
	// SelectInfoSchema reads system metadata.
	SelectInfoSchema = analysis.SelectInfoSchema
	// DDL changes schema.
	DDL = analysis.DDL
	// DML changes table data.
	DML = analysis.DML
)

// ColumnResource identifies a column, or a whole table when Column is empty.
type ColumnResource = analysis.ColumnResource

// ColumnSet is a set of columns.
type ColumnSet = analysis.ColumnSet

// sortedColumnSet returns the set's resources in order, for walks whose answer
// must not depend on map iteration order.
func sortedColumnSet(s ColumnSet) []ColumnResource {
	resources := make([]ColumnResource, 0, len(s))
	for resource := range s {
		resources = append(resources, resource)
	}
	slices.SortFunc(resources, func(a, b ColumnResource) int {
		return cmp.Or(
			cmp.Compare(a.Database, b.Database),
			cmp.Compare(a.Schema, b.Schema),
			cmp.Compare(a.Table, b.Table),
			cmp.Compare(a.Column, b.Column),
		)
	})
	return resources
}

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult = analysis.QuerySpanResult

// ResolvedSpan is the lineage ResolveQuerySpan returns.
type ResolvedSpan = analysis.ResolvedSpan

// UnresolvedColumnsError identifies relations whose synced metadata has no
// columns.
type UnresolvedColumnsError = analysis.UnresolvedColumnsError

// SyntaxError reports a statement that does not parse. Err is the parser's
// error, a *parser.ParseError when the parser located the failure.
type SyntaxError struct {
	Err error
}

func (e *SyntaxError) Error() string { return e.Err.Error() }

func (e *SyntaxError) Unwrap() error { return e.Err }

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = analysis.ErrMixUserSystemTables

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
