// Package analysis is the contract between Bytebase and the engine analysis
// packages (mysql/analysis, pg/analysis, ...). It holds the shared types and
// sentinel errors that cross that boundary and nothing else: no Bytebase
// imports and no engine imports.
//
// Each engine analysis package owns the AST walks and engine-specific options
// for query span resolution, statement classification, changed resources,
// statement splitting, and byte ranges, and re-exports the types here through
// Go type aliases so callers across engines share one canonical vocabulary.
package analysis

import "strings"

// QueryType is the permission class of a statement. Its numeric values mirror
// Bytebase's base.QueryType so callers and corpus fixtures map them 1:1.
type QueryType int

const (
	// QueryTypeUnknown is a statement not deliberately classified.
	QueryTypeUnknown QueryType = iota
	// Select is a read-only user query.
	Select
	// Explain is a query plan statement.
	Explain
	// SelectInfoSchema is a read-only query against system or information
	// schema objects, or a SHOW / DESCRIBE statement.
	SelectInfoSchema
	// DDL changes schema.
	DDL
	// DML changes table data.
	DML
)

// String returns the canonical name of the query type.
func (q QueryType) String() string {
	switch q {
	case Select:
		return "SELECT"
	case Explain:
		return "EXPLAIN"
	case SelectInfoSchema:
		return "SELECT_INFO_SCHEMA"
	case DDL:
		return "DDL"
	case DML:
		return "DML"
	default:
		return "UNKNOWN"
	}
}

// ColumnResource identifies a physical column, or a whole table when Column is
// empty.
type ColumnResource struct {
	// Instance names the remote instance a link resolved to (for example, an
	// Oracle database link); empty means the connected instance.
	Instance string
	// Server is the normalized remote server or database-link name as written
	// in the statement; empty means the connected server.
	Server   string
	Database string
	Schema   string
	Table    string
	Column   string
}

// String renders the resource as a dot-separated qualified name, omitting empty
// components. Instance is not part of the SQL name and is excluded.
func (c ColumnResource) String() string {
	parts := make([]string, 0, 5)
	for _, part := range []string{c.Server, c.Database, c.Schema, c.Table, c.Column} {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}

// ColumnSet is a set of column resources.
type ColumnSet map[ColumnResource]bool

// QuerySpanResult describes one output column of a resolved query span and the
// physical columns it reads.
type QuerySpanResult struct {
	Name          string
	SourceColumns ColumnSet
	// IsPlainField is true when the result is a bare column reference rather
	// than an expression.
	IsPlainField bool
	// SelectAsterisk marks a `SELECT *` that could not be expanded against
	// metadata.
	SelectAsterisk bool
	// Invisible marks a table column that `SELECT *` omits (such as an Oracle
	// INVISIBLE column), while a reference by name still resolves it.
	Invisible bool
}

// ResolvedSpan is the lineage and classification returned by metadata-backed
// query span resolution.
type ResolvedSpan struct {
	Type QueryType
	// ExplainAnalyze marks a statement that executes and also returns its plan
	// (such as EXPLAIN ANALYZE or SET STATISTICS XML).
	ExplainAnalyze bool
	// SourceColumns holds the tables the statement accesses, as table-level
	// resources (Column empty).
	SourceColumns ColumnSet
	Results       []QuerySpanResult
	// PredicateColumns holds the columns that decide which rows return.
	PredicateColumns ColumnSet
	// NotFoundError is set when a resource the query names is missing from the
	// catalog; Results is then empty.
	NotFoundError *ResourceNotFoundError
	// UnresolvedColumnsError is set when a relation the query reads has no
	// columns in the synced metadata, so its lineage cannot be trusted.
	UnresolvedColumnsError *UnresolvedColumnsError
}
