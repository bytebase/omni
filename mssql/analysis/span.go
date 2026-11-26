package analysis

import (
	"fmt"
	"slices"
	"strings"

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
	// Explain is a statement whose result is a plan.
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

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult = analysis.QuerySpanResult

// ResolvedSpan is the lineage ResolveQuerySpan returns.
type ResolvedSpan struct {
	Type QueryType
	// ExplainAnalyze marks a statement that runs and also returns its plan, as
	// under SET STATISTICS XML.
	ExplainAnalyze bool
	// SourceColumns holds the tables the statement accesses, as table-level
	// resources (Column empty).
	SourceColumns ColumnSet
	Results       []QuerySpanResult
	// PredicateColumns holds the columns that decide which rows return.
	PredicateColumns ColumnSet
	// NotFoundError is set when a table the query names is not in the catalog.
	// Results is then empty.
	NotFoundError *NotFoundError
	// Untraced is set when the statement reads a value whose columns this
	// package cannot trace. The statement is legal and its rows are right; a
	// masking caller must withhold them, because it cannot tell which column a
	// policy applies to.
	Untraced *UntracedError
}

// UntracedError reports the values a statement reads that this package cannot
// trace to a column of a table the synced metadata describes, so a masking
// caller must withhold the rows the statement returns. ResolveQuerySpan sees
// one statement and the batch's temp tables as the caller carries them, which
// is why the batch's own writes and variable assignments are out of reach; see
// the package doc.
type UntracedError struct {
	// Relations are the temp tables and table variables the statement reads,
	// each with an empty Column field. The batch's own statements wrote their
	// rows.
	Relations []ColumnResource
	// Variables are the local variables the statement reads, with their
	// leading @. A statement of the same batch assigned each value, from a
	// column this package did not see.
	Variables []string
}

// Error names what the statement reads that has no column lineage.
func (e *UntracedError) Error() string {
	parts := make([]string, 0, 2)
	if len(e.Relations) > 0 {
		names := make([]string, 0, len(e.Relations))
		for _, relation := range e.Relations {
			names = append(names, relation.String())
		}
		slices.Sort(names)
		parts = append(parts, fmt.Sprintf("%s (whose rows the batch's own statements wrote)", strings.Join(names, ", ")))
	}
	if len(e.Variables) > 0 {
		variables := slices.Clone(e.Variables)
		slices.Sort(variables)
		parts = append(parts, fmt.Sprintf("%s (whose value a statement of the batch assigned)", strings.Join(variables, ", ")))
	}
	return "the statement's rows cannot be traced to columns: it reads " + strings.Join(parts, " and ")
}

// LinkedServerError reports a statement that reads a relation on a linked
// server, or through an ad-hoc provider, whose columns the connected
// instance's metadata does not describe. It is returned as an error rather
// than carried on the span: the server, not the connected instance, resolves
// the relation, so a span without its columns would read as a statement that
// touches nothing sensitive and a masking caller would return the rows.
type LinkedServerError struct {
	// Server is the linked server as written; empty for an ad-hoc provider
	// (OPENROWSET, OPENDATASOURCE).
	Server string
	// Table is the relation as written after the server, when the statement
	// names one: a four-part name's database, schema and table.
	Table string
	// PassThrough marks a call whose query text the remote server runs
	// (OPENQUERY), so the statement spells none of the columns even when it
	// names no server.
	PassThrough bool
}

// Error names the server the read goes through and why it has no lineage.
func (e *LinkedServerError) Error() string {
	switch {
	case e.Table != "":
		return fmt.Sprintf("linked server %q: the connected instance's metadata does not describe %q, so the columns it returns cannot be traced", e.Server, e.Table)
	case e.Server != "":
		return fmt.Sprintf("linked server %q: a pass-through query's columns are the remote server's, so they cannot be traced", e.Server)
	case e.PassThrough:
		return "a pass-through query's columns are the remote server's, so they cannot be traced"
	default:
		return "an ad-hoc provider read's columns are the provider's, not the connected instance's, so they cannot be traced"
	}
}

// NotFoundError reports a database, schema, table or column the catalog does
// not have. A nil field is not part of the report.
type NotFoundError = analysis.ResourceNotFoundError

// TypeNotSupportedError reports a construct whose lineage is not modeled,
// such as a table source the resolver does not model.
type TypeNotSupportedError = analysis.TypeNotSupportedError

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = analysis.ErrMixUserSystemTables

// Catalog is the synced metadata of the instance a statement runs on.
type Catalog struct {
	// ListDatabaseNames lists the instance's databases.
	ListDatabaseNames func() ([]string, error)
	// GetDatabase returns a database's metadata, or nil when the database
	// has not been synced.
	GetDatabase func(name string) (Database, error)
}

// Database is one database's synced metadata.
type Database interface {
	ListSchemaNames() []string
	// GetSchema returns the named schema, or nil.
	GetSchema(name string) Schema
}

// Schema is one schema's synced metadata.
type Schema interface {
	ListTableNames() []string
	// GetTable returns the named table, or nil.
	GetTable(name string) *metadata.TableMetadata
	ListViewNames() []string
	// GetView returns the named view, or nil.
	GetView(name string) *metadata.ViewMetadata
}

// TempTable is a temp table or table variable a batch declared, keyed in
// ResolveOptions.TempTables by its name as written, such as "#t" or "@t".
type TempTable struct {
	Server   string
	Database string
	Schema   string
	Name     string
	Columns  []string
}

// mergeColumnSet returns the union of m and n, and whether it is larger than m.
func mergeColumnSet(m, n ColumnSet) (ColumnSet, bool) {
	r := make(ColumnSet)
	for k := range m {
		r[k] = true
	}
	for k := range n {
		if _, ok := r[k]; !ok {
			r[k] = true
		}
	}
	return r, len(r) != len(m)
}

// tableSource is a relation in scope.
type tableSource interface {
	GetTableName() string
	GetSchemaName() string
	GetDatabaseName() string
	// GetQuerySpanResult returns the relation's columns. A pseudoTable and a
	// physicalView return their own slice, which callers must copy before
	// modifying.
	GetQuerySpanResult() []QuerySpanResult
}

// pseudoTable is a derived relation: a CTE, subquery, alias or join.
type pseudoTable struct {
	Name    string
	Columns []QuerySpanResult
}

func (p *pseudoTable) GetTableName() string                  { return p.Name }
func (*pseudoTable) GetSchemaName() string                   { return "" }
func (*pseudoTable) GetDatabaseName() string                 { return "" }
func (p *pseudoTable) GetQuerySpanResult() []QuerySpanResult { return p.Columns }

// physicalView is a catalog view with its resolved columns.
type physicalView struct {
	Database string
	Schema   string
	Name     string
	Columns  []QuerySpanResult
}

func (p *physicalView) GetTableName() string                  { return p.Name }
func (p *physicalView) GetSchemaName() string                 { return p.Schema }
func (p *physicalView) GetDatabaseName() string               { return p.Database }
func (p *physicalView) GetQuerySpanResult() []QuerySpanResult { return p.Columns }

// physicalTable is a catalog or temp table; each column is its own plain
// source.
type physicalTable struct {
	Server   string
	Database string
	Schema   string
	Name     string
	Columns  []string
}

func (p *physicalTable) GetTableName() string    { return p.Name }
func (p *physicalTable) GetSchemaName() string   { return p.Schema }
func (p *physicalTable) GetDatabaseName() string { return p.Database }

func (p *physicalTable) GetQuerySpanResult() []QuerySpanResult {
	result := make([]QuerySpanResult, 0, len(p.Columns))
	for _, column := range p.Columns {
		sourceColumnSet := make(ColumnSet, 1)
		sourceColumnSet[ColumnResource{
			Server:   p.Server,
			Database: p.Database,
			Schema:   p.Schema,
			Table:    p.Name,
			Column:   column,
		}] = true
		result = append(result, QuerySpanResult{
			Name:          column,
			SourceColumns: sourceColumnSet,
			IsPlainField:  true,
		})
	}
	return result
}
