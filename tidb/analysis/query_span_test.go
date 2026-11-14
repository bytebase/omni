package analysis_test

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/analysis"
)

// testDatabase is an in-memory Database; names match exactly.
type testDatabase struct {
	name   string
	tables map[string][]string
	views  map[string]string
}

func (d *testDatabase) Name() string             { return d.name }
func (*testDatabase) HasSchema() bool            { return true }
func (d *testDatabase) ListTableNames() []string { return sortedKeys(d.tables) }
func (d *testDatabase) ListViewNames() []string  { return sortedKeys(d.views) }

func (d *testDatabase) GetTable(name string) *metadata.TableMetadata {
	columns, ok := d.tables[name]
	if !ok {
		return nil
	}
	table := &metadata.TableMetadata{Name: name}
	for _, column := range columns {
		table.Columns = append(table.Columns, &metadata.ColumnMetadata{Name: column})
	}
	return table
}

func (d *testDatabase) GetView(name string) *metadata.ViewMetadata {
	definition, ok := d.views[name]
	if !ok {
		return nil
	}
	return &metadata.ViewMetadata{Name: name, Definition: definition}
}

func sortedKeys[V any](m map[string]V) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

var testDB = &testDatabase{
	name: "db",
	tables: map[string][]string{
		"t1": {"a", "b", "c"},
		"t2": {"a", "d"},
	},
	views: map[string]string{
		"v":      "SELECT a, b + c AS bc FROM t1",
		"cycle1": "SELECT * FROM cycle2",
		"cycle2": "SELECT * FROM cycle1",
		"mixed":  "SELECT * FROM t1, mysql.user",
	},
}

var systemDB = &testDatabase{name: "mysql", tables: map[string][]string{"user": {"host", "user"}}}

func testCatalog() analysis.Catalog {
	return analysis.Catalog{
		GetDatabase: func(name string) (analysis.Database, error) {
			switch name {
			case "db":
				return testDB, nil
			case "mysql":
				return systemDB, nil
			}
			return nil, nil
		},
	}
}

func resolve(t *testing.T, statement string) (*analysis.ResolvedSpan, error) {
	t.Helper()
	return analysis.ResolveQuerySpan(statement, analysis.ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
}

func renderSet(set analysis.ColumnSet) string {
	var out []string
	for c := range set {
		out = append(out, strings.Join([]string{c.Database, c.Table, c.Column}, "."))
	}
	sort.Strings(out)
	return strings.Join(out, ",")
}

func renderResults(span *analysis.ResolvedSpan) []string {
	var out []string
	for _, r := range span.Results {
		s := r.Name + "<-" + renderSet(r.SourceColumns)
		if r.IsPlainField {
			s += " plain"
		}
		out = append(out, s)
	}
	return out
}

func TestResolveQuerySpanLineage(t *testing.T) {
	tests := []struct {
		statement  string
		results    []string
		predicates string
	}{
		{"SELECT * FROM t1", []string{"a<-db.t1.a plain", "b<-db.t1.b plain", "c<-db.t1.c plain"}, ""},
		{"SELECT a + b AS s FROM t1 WHERE c > 0", []string{"s<-db.t1.a,db.t1.b"}, "db.t1.c"},
		{"SELECT * FROM t1 JOIN t2 USING (a)", []string{"a<-db.t1.a plain", "b<-db.t1.b plain", "c<-db.t1.c plain", "d<-db.t2.d plain"}, "db.t1.a,db.t2.a"},
		{"SELECT bc FROM v", []string{"bc<-db.t1.b,db.t1.c"}, ""},
		{"WITH x(k) AS (SELECT d FROM t2) SELECT k FROM x", []string{"k<-db.t2.d"}, ""},
		{"SELECT a FROM t1 UNION SELECT d FROM t2", []string{"a<-db.t1.a,db.t2.d"}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolve(t, tc.statement)
			require.NoError(t, err)
			require.Equal(t, analysis.Select, span.Type)
			require.Equal(t, tc.results, renderResults(span))
			require.Nil(t, span.NotFoundError)
		})
	}
}

func TestResolveQuerySpanNotFound(t *testing.T) {
	span, err := resolve(t, "SELECT nope FROM t1")
	require.NoError(t, err)
	require.Empty(t, span.Results)
	require.Equal(t, "db.t1.", renderSet(span.SourceColumns))
	require.EqualError(t, span.NotFoundError, "resource not found: database: , table: , column: nope, err: cannot find the column ref")

	span, err = resolve(t, "SELECT * FROM missing")
	require.NoError(t, err)
	require.EqualError(t, span.NotFoundError, "resource not found: database: db, schema: , table: missing")
}

func TestResolveQuerySpanErrors(t *testing.T) {
	_, err := resolve(t, "SELECT * FROM t1, mysql.user")
	require.ErrorIs(t, err, analysis.ErrMixUserSystemTables)

	// A system table counts only when the catalog has it, so one the catalog
	// lacks is a missing table instead.
	span, err := resolve(t, "SELECT * FROM t1, mysql.db")
	require.NoError(t, err)
	require.EqualError(t, span.NotFoundError, "resource not found: database: mysql, schema: , table: db")

	_, err = resolve(t, "SELECT * FROM mixed")
	require.ErrorIs(t, err, analysis.ErrMixUserSystemTables)
	require.EqualError(t, err, "failed to extract table source from FROM clause: cannot access user and system tables at the same time")

	_, err = resolve(t, "SELECT * FROM cycle1")
	require.ErrorContains(t, err, `cyclic view reference detected while resolving "cycle1"`)

	boom := errors.New("boom")
	_, err = analysis.ResolveQuerySpan("SELECT * FROM t1", analysis.ResolveOptions{
		DefaultDatabase: "db",
		Catalog:         analysis.Catalog{GetDatabase: func(string) (analysis.Database, error) { return nil, boom }},
	})
	require.ErrorIs(t, err, boom)
	require.EqualError(t, err, `failed to get database metadata for "db": boom`)
}

func TestResolveQuerySpanNonQueries(t *testing.T) {
	tests := []struct {
		statement string
		queryType analysis.QueryType
		analyze   bool
		sources   string
	}{
		{"SELECT * FROM mysql.user", analysis.SelectInfoSchema, false, "mysql.user."},
		{"EXPLAIN SELECT * FROM t1", analysis.Explain, false, "db.t1."},
		{"EXPLAIN ANALYZE SELECT * FROM t1", analysis.Select, true, "db.t1."},
		{"SHOW TABLES", analysis.SelectInfoSchema, false, ""},
		{"UPDATE t1 SET a = 1", analysis.DML, false, "db.t1."},
		{"CREATE TABLE t3 (a INT)", analysis.DDL, false, ""},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolve(t, tc.statement)
			require.NoError(t, err)
			require.Equal(t, tc.queryType, span.Type)
			require.Equal(t, tc.analyze, span.ExplainAnalyze)
			require.Equal(t, tc.sources, renderSet(span.SourceColumns))
		})
	}
}
