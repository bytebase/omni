package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// ResolveOptions carries the facts the caller resolves from its metadata store,
// plus the dialect deltas that decide how the resolved span is shaped.
//
// The metadata model flags are the BigQuery vs Spanner seams the two legacy
// extractors differed on; testdata/query-span and the dialect tests beside this
// file pin them.
type ResolveOptions struct {
	// Dialect selects the GoogleSQL dialect (system schemas, name-part
	// bucketing, default join type) for the structural walk.
	Dialect Dialect

	// DefaultDatabase is the database/dataset a reference without one resolves
	// against, and the database the metadata lookup falls back to.
	DefaultDatabase string

	// IgnoreWrittenDatabase pins every resource's Database to DefaultDatabase.
	// Spanner: the legacy resolver silently ignored the db part of a 3-part
	// db.schema.table and recorded the session database. BigQuery keeps the
	// written dataset with default fill-in.
	IgnoreWrittenDatabase bool

	// SchemaScopedMetadata looks tables up under the metadata schema named by
	// the reference (case-SENSITIVE, "" = default schema) within the default
	// database — the Spanner named-schema model. BigQuery instead resolves the
	// written dataset as the DATABASE and uses its single unnamed schema.
	SchemaScopedMetadata bool

	// SystemOnlyEmptySpan early-returns an EMPTY SelectInfoSchema span (no
	// source columns, no results, no metadata access) for a query reading ONLY
	// system tables — the legacy spanner extractor's behavior. BigQuery instead
	// keeps the system tables as table-level resources.
	SystemOnlyEmptySpan bool

	// SetStatementIsSelect classifies a SET statement as Select — both legacy
	// query-type listeners treated SET as a select, and the access check would
	// otherwise reject an allowed statement as an unknown query type.
	SetStatementIsSelect bool

	// UppercaseStarMergeName upper-cases a set-operation star-merge output name
	// (the legacy BigQuery extractor's rendering). The legacy spanner extractor
	// passed a LEFT-star-derived name through in the star column's metadata
	// case instead.
	UppercaseStarMergeName bool

	// CanonicalBaseFieldName renders a base-table FIELD passthrough (the JOIN
	// ... USING coalesced key, marked BaseFieldName by the span walker) in the
	// field's METADATA case — the legacy spanner resolver named it after the
	// left PhysicalTable's field, not the written USING token.
	CanonicalBaseFieldName bool

	// GetDatabaseMetadata returns the synced metadata for a database (BigQuery
	// dataset, or the one Spanner database). It returns nil for a database that
	// is not synced; an error other than a *ResourceNotFoundError aborts the
	// span.
	GetDatabaseMetadata func(database string) (*metadata.DatabaseSchemaMetadata, error)
}

// ColumnResource identifies a physical column resource. Database is a BigQuery
// dataset or a Spanner database; Schema is empty under the BigQuery model.
type ColumnResource struct {
	Database string
	Schema   string
	Table    string
	Column   string
}

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
	Type          QueryType
	SourceColumns ColumnSet
	// PredicateColumns holds the columns the statement's row conditions read —
	// WHERE, JOIN ... ON / USING, and HAVING — as physical resources resolved the
	// same way result lineage is, so they carry the dialect's name model. Empty
	// for a statement that produces no result set (DML, DDL, Explain, Unknown),
	// which keeps no lineage either.
	PredicateColumns ColumnSet
	Results          []ResolvedQuerySpanResult
	NotFoundError    error

	// predicateColumnsUnresolved is why PredicateColumns is empty when one of the
	// statement's row conditions could not be resolved: an empty set is honest, a set
	// missing a read is not, so such a statement reports none at all.
	predicateColumnsUnresolved string
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

// ErrMixUserSystemTables reports a statement that reads user and system
// (INFORMATION_SCHEMA / SPANNER_SYS) tables at the same time.
var ErrMixUserSystemTables = fmt.Errorf("cannot access user and system tables at the same time")

// ResolveQuerySpan analyzes a single GoogleSQL statement in the dialect the
// options name and returns the physical columns it reads, the per-output-column
// lineage (expanded against catalog metadata), its statement classification,
// and the first relation whose metadata was missing.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedQuerySpan, error) {
	return newQuerySpanResolver(opts).getQuerySpan(statement)
}
