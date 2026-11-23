// Package analysis holds Bytebase's query analysis for Oracle: splitting and
// parsing scripts, and query span (column lineage and the tables a statement
// reads).
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
	// Explain is an EXPLAIN PLAN statement.
	Explain = analysis.Explain
	// SelectInfoSchema reads only system schemas.
	SelectInfoSchema = analysis.SelectInfoSchema
	// DDL changes schema.
	DDL = analysis.DDL
	// DML changes table data.
	DML = analysis.DML
)

// ColumnResource identifies a column, or a whole table when Column is empty.
// In Oracle, Database is the schema.
type ColumnResource = analysis.ColumnResource

// ColumnSet is a set of columns.
type ColumnSet = analysis.ColumnSet

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult = analysis.QuerySpanResult

// ResolvedSpan is the lineage ResolveQuerySpan returns.
type ResolvedSpan = analysis.ResolvedSpan

// NotFoundError reports a resource missing from the catalog. Its fields and
// text are those of Bytebase's base.ResourceNotFoundError.
type NotFoundError = analysis.ResourceNotFoundError

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = analysis.ErrMixUserSystemTables

// Catalog is the synced metadata a statement resolves against. An instance is
// the name the catalog gives it: ResolveOptions.Instance for the connected
// instance, and what GetLinkedDatabase returns for a linked one.
type Catalog struct {
	// GetDatabase returns the metadata of database (an Oracle schema) on
	// instance, or a nil Database when it has not been synced. An error
	// matching *NotFoundError under errors.As is reported as a missing
	// resource.
	GetDatabase func(instance, database string) (Database, error)
	// GetLinkedDatabase resolves the database link link from instance, for a
	// table in schema, to the instance and database it reaches. A nil
	// Database means the link resolves to nothing; a nil GetLinkedDatabase
	// resolves every link to nothing.
	GetLinkedDatabase func(instance, link, schema string) (linkedInstance string, database string, db Database, err error)
}

// Database is one database's synced metadata. The lookups follow the
// metadata's own name matching.
type Database interface {
	// HasSchema reports whether the database has its unnamed schema, where
	// Oracle metadata keeps its objects.
	HasSchema() bool
	GetTable(name string) *metadata.TableMetadata
	GetView(name string) *metadata.ViewMetadata
	GetMaterializedView(name string) *metadata.MaterializedViewMetadata
	GetExternalTable(name string) *metadata.ExternalTableMetadata
}

// SynonymDatabase is a Database that also knows its synonyms (CREATE
// SYNONYM). A name that matches no table or view of a database resolves
// through the database's synonym of that name, and an unqualified one then
// through the instance's public synonym. A Database that does not implement
// it resolves no synonym, and a name only a synonym gives is not found.
type SynonymDatabase interface {
	// GetSynonym returns the database's synonym name, or nil.
	GetSynonym(name string) *metadata.SynonymMetadata
	// GetPublicSynonym returns the public synonym name of the database's
	// instance, or nil.
	GetPublicSynonym(name string) *metadata.SynonymMetadata
}

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

// spanTable is a relation a query reads: a table, a view, a CTE or a derived
// table.
type spanTable interface {
	GetTableName() string
	GetSchemaName() string
	GetDatabaseName() string
	GetServerName() string
	// GetQuerySpanResult returns the relation's columns; callers copy before
	// changing them.
	GetQuerySpanResult() []QuerySpanResult
}

type pseudoTable struct {
	Name    string
	Columns []QuerySpanResult
}

func (p *pseudoTable) GetTableName() string                  { return p.Name }
func (*pseudoTable) GetSchemaName() string                   { return "" }
func (*pseudoTable) GetDatabaseName() string                 { return "" }
func (*pseudoTable) GetServerName() string                   { return "" }
func (p *pseudoTable) GetQuerySpanResult() []QuerySpanResult { return p.Columns }

type physicalTable struct {
	Server           string
	Database         string
	Schema           string
	Name             string
	Columns          []string
	InvisibleColumns map[string]bool
}

func (p *physicalTable) GetTableName() string    { return p.Name }
func (p *physicalTable) GetSchemaName() string   { return p.Schema }
func (p *physicalTable) GetDatabaseName() string { return p.Database }
func (p *physicalTable) GetServerName() string   { return p.Server }

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
			Invisible:     p.InvisibleColumns[column],
		})
	}
	return result
}
