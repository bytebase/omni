package analysis

import (
	"errors"
	"fmt"
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
	// Explain is a statement whose result is a plan.
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
	// Server is the linked server the statement names; empty for the
	// connected server.
	Server   string
	Database string
	Schema   string
	Table    string
	Column   string
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
}

// NotFoundError reports a database, schema, table or column the catalog does
// not have. A nil field is not part of the report.
type NotFoundError struct {
	Err      error
	Database *string
	Schema   *string
	Table    *string
	Column   *string
}

// Error formats the error as Bytebase's base.ResourceNotFoundError does.
func (e *NotFoundError) Error() string {
	var resourceParts []string
	if e.Database != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("database: %s", *e.Database))
	}
	if e.Schema != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("schema: %s", *e.Schema))
	}
	if e.Table != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("table: %s", *e.Table))
	}
	if e.Column != nil {
		resourceParts = append(resourceParts, fmt.Sprintf("column: %s", *e.Column))
	}
	parts := []string{
		fmt.Sprintf("resource not found: %s", strings.Join(resourceParts, ", ")),
	}
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("err: %s", e.Err.Error()))
	}
	return strings.Join(parts, ", ")
}

// TypeNotSupportedError reports a construct whose lineage is not modeled,
// such as a PIVOT table source.
type TypeNotSupportedError struct {
	Err   error
	Type  string
	Name  string
	Extra string
}

// Error formats the error as Bytebase's base.TypeNotSupportedError does.
func (e *TypeNotSupportedError) Error() string {
	parts := []string{
		fmt.Sprintf("type not supported: %s", e.Type),
	}
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("err: %s", e.Err.Error()))
	}
	if e.Name != "" {
		parts = append(parts, fmt.Sprintf("name: %s", e.Name))
	}
	if e.Extra != "" {
		parts = append(parts, fmt.Sprintf("extra: %s", e.Extra))
	}
	return strings.Join(parts, ", ")
}

func (e *TypeNotSupportedError) Unwrap() error {
	return e.Err
}

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = errors.New("cannot access user and system tables at the same time")

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
