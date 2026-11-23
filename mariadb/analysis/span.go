package analysis

import (
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
// Schema is always empty here.
type ColumnResource = analysis.ColumnResource

// ColumnSet is a set of columns.
type ColumnSet = analysis.ColumnSet

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult = analysis.QuerySpanResult

// ResolvedSpan is the lineage ResolveQuerySpan returns.
type ResolvedSpan = analysis.ResolvedSpan

// NotFoundError reports a database, schema, table or column the catalog does
// not have. A nil field is not part of the report.
type NotFoundError = analysis.ResourceNotFoundError

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = analysis.ErrMixUserSystemTables

// Catalog is the synced metadata of the instance a statement runs on.
type Catalog struct {
	// ListDatabaseNames lists the instance's databases.
	ListDatabaseNames func() ([]string, error)
	// GetDatabase returns a database's metadata, or nil when the database
	// has not been synced. A nil GetDatabase means no metadata is available.
	GetDatabase func(name string) (Database, error)
}

// Database is one database's synced metadata. MariaDB keeps its tables and
// views in the database's single unnamed schema, which the list and get
// methods read.
type Database interface {
	// Name is the database's name as the metadata spells it.
	Name() string
	// HasSchema reports whether the metadata has the unnamed schema.
	HasSchema() bool
	ListTableNames() []string
	// GetTable returns the named table, or nil.
	GetTable(name string) *metadata.TableMetadata
	ListViewNames() []string
	// GetView returns the named view, or nil.
	GetView(name string) *metadata.ViewMetadata
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
	GetDatabaseName() string
	// GetQuerySpanResult returns the relation's columns. A pseudoTable and a
	// physicalView return their own slice, which callers must copy before
	// modifying.
	GetQuerySpanResult() []QuerySpanResult
}

// pseudoTable is a derived relation: a CTE, subquery, alias or join.
type pseudoTable struct {
	// Database is the database of a table or view read through an alias,
	// which MySQL still matches in db.alias.column; empty otherwise.
	Database string
	Name     string
	Columns  []QuerySpanResult
}

func (p *pseudoTable) GetTableName() string                  { return p.Name }
func (p *pseudoTable) GetDatabaseName() string               { return p.Database }
func (p *pseudoTable) GetQuerySpanResult() []QuerySpanResult { return p.Columns }

// physicalView is a catalog view with its resolved columns.
type physicalView struct {
	Database string
	Schema   string
	Name     string
	Columns  []QuerySpanResult
}

func (p *physicalView) GetTableName() string                  { return p.Name }
func (p *physicalView) GetDatabaseName() string               { return p.Database }
func (p *physicalView) GetQuerySpanResult() []QuerySpanResult { return p.Columns }

// physicalTable is a catalog table; each column is its own plain source.
type physicalTable struct {
	Database string
	Schema   string
	Name     string
	Columns  []string
}

func (p *physicalTable) GetTableName() string    { return p.Name }
func (p *physicalTable) GetDatabaseName() string { return p.Database }

func (p *physicalTable) GetQuerySpanResult() []QuerySpanResult {
	result := make([]QuerySpanResult, 0, len(p.Columns))
	for _, column := range p.Columns {
		sourceColumnSet := make(ColumnSet, 1)
		sourceColumnSet[ColumnResource{
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
