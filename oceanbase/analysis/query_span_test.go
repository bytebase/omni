package analysis_test

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/analysis"
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

func testCatalog() analysis.Catalog {
	return analysis.Catalog{
		ListDatabaseNames: func() ([]string, error) { return []string{"db"}, nil },
		GetDatabase: func(name string) (analysis.Database, error) {
			if name == "db" {
				return testDB, nil
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
		{"SELECT a + b AS s FROM t1 WHERE c > 0", []string{"s<-db.t1.a,db.t1.b"}, ""},
		{"SELECT * FROM t1 JOIN t2 USING (a)", []string{"a<-db.t1.a,db.t2.a plain", "b<-db.t1.b plain", "c<-db.t1.c plain", "d<-db.t2.d plain"}, ""},
		// MySQL returns the merged columns first, in the first table's order; RIGHT JOIN's first table is the right one.
		{"SELECT * FROM (SELECT b, a FROM t1) AS x JOIN t2 USING (a)", []string{"a<-db.t1.a,db.t2.a plain", "b<-db.t1.b plain", "d<-db.t2.d plain"}, ""},
		{"SELECT * FROM (SELECT b, a FROM t1) AS x RIGHT JOIN t2 USING (a)", []string{"a<-db.t1.a,db.t2.a plain", "d<-db.t2.d plain", "b<-db.t1.b plain"}, ""},
		{"SELECT * FROM (SELECT b, a FROM t1) AS x NATURAL RIGHT JOIN t2", []string{"a<-db.t1.a,db.t2.a plain", "d<-db.t2.d plain", "b<-db.t1.b plain"}, ""},
		{"SELECT * FROM (SELECT b, a FROM t1) AS x NATURAL LEFT JOIN t2", []string{"a<-db.t1.a,db.t2.a plain", "b<-db.t1.b plain", "d<-db.t2.d plain"}, ""},
		// A subquery's own tables come before the enclosing query's.
		{"SELECT (SELECT a FROM t2 LIMIT 1) AS x FROM t1", []string{"x<-db.t2.a plain"}, ""},
		{"SELECT (SELECT y.a FROM t2 AS y LIMIT 1) AS x FROM t1 AS y", []string{"x<-db.t2.a plain"}, ""},
		{"SELECT (SELECT d FROM t2 WHERE t2.a = t1.b LIMIT 1) AS x FROM t1", []string{"x<-db.t2.d plain"}, ""},
		{"SELECT (SELECT b FROM t2 LIMIT 1) AS x FROM t1", []string{"x<-db.t1.b plain"}, ""},
		{"SELECT bc FROM v", []string{"bc<-db.t1.b,db.t1.c"}, ""},
		{"SELECT a, v.bc FROM v", []string{"a<-db.t1.a plain", "bc<-db.t1.b,db.t1.c"}, ""},
		{"SELECT s FROM (SELECT a + b AS s, c FROM t1) AS d", []string{"s<-db.t1.a,db.t1.b"}, ""},
		{"SELECT d.c FROM (SELECT a + b AS s, c FROM t1) AS d", []string{"c<-db.t1.c plain"}, ""},
		{"WITH x AS (SELECT a + b AS s FROM t1) SELECT s FROM x", []string{"s<-db.t1.a,db.t1.b"}, ""},
		{"WITH x(k) AS (SELECT d FROM t2) SELECT k FROM x", []string{"k<-db.t2.d plain"}, ""},
		{"SELECT a FROM t1 UNION SELECT d FROM t2", []string{"a<-db.t1.a,db.t2.d"}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolve(t, tc.statement)
			require.NoError(t, err)
			require.Equal(t, analysis.Select, span.Type)
			require.Equal(t, tc.results, renderResults(span))
			require.Equal(t, tc.predicates, renderSet(span.PredicateColumns))
			require.Nil(t, span.NotFoundError)
		})
	}
}

func TestResolveQuerySpanNotFound(t *testing.T) {
	span, err := resolve(t, "SELECT nope FROM t1")
	require.NoError(t, err)
	require.Empty(t, span.Results)
	require.Equal(t, "db.t1.", renderSet(span.SourceColumns))
	require.EqualError(t, span.NotFoundError, "resource not found: database: , table: , column: nope")

	span, err = resolve(t, "SELECT * FROM missing")
	require.NoError(t, err)
	require.EqualError(t, span.NotFoundError, "resource not found: database: db, schema: , table: missing")
}

func TestResolveQuerySpanErrors(t *testing.T) {
	_, err := resolve(t, "SELECT * FROM t1, mysql.user")
	require.ErrorIs(t, err, analysis.ErrMixUserSystemTables)

	_, err = resolve(t, "SELECT * FROM mixed")
	require.ErrorIs(t, err, analysis.ErrMixUserSystemTables)
	require.EqualError(t, err, `failed to get columns for view "mixed": failed to get query span for view: cannot access user and system tables at the same time`)

	_, err = resolve(t, "SELECT * FROM cycle1")
	require.ErrorContains(t, err, `cyclic view reference detected while resolving "cycle1"`)

	boom := errors.New("boom")
	_, err = analysis.ResolveQuerySpan("SELECT * FROM t1", analysis.ResolveOptions{
		DefaultDatabase: "db",
		Catalog:         analysis.Catalog{ListDatabaseNames: func() ([]string, error) { return nil, boom }},
	})
	require.ErrorIs(t, err, boom)
	require.EqualError(t, err, "failed to list databases: boom")
}

func TestResolveQuerySpanNonQueries(t *testing.T) {
	tests := []struct {
		statement string
		queryType analysis.QueryType
		analyze   bool
		sources   string
	}{
		{"SELECT * FROM information_schema.tables", analysis.SelectInfoSchema, false, ""},
		{"EXPLAIN SELECT * FROM t1", analysis.Explain, false, "db.t1."},
		{"SHOW TABLES", analysis.SelectInfoSchema, false, ""},
		{"UPDATE t1 SET a = 1", analysis.DML, false, "db.t1."},
		{"CREATE TABLE t3 (a INT)", analysis.DDL, false, "db.t3."},
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

// TestResolveQuerySpanTableFunctions: the table sources OceanBase's grammar adds
// — UNNEST, RB_ITERATE and a table read from a location — are analyzed rather
// than refused as an unsupported source.
func TestResolveQuerySpanTableFunctions(t *testing.T) {
	for _, statement := range []string{
		"SELECT * FROM UNNEST(t1.a, t1.b) AS u(c, d)",
		"SELECT c FROM UNNEST(t1.a) AS u(c)",
		"SELECT * FROM RB_ITERATE(t1.a) AS x(v)",
		"SELECT * FROM SOURCE(TYPE = 's')",
		"SELECT * FROM 's'(FORMAT(TYPE = 'csv')) AS e",
		// A column of a source whose file schema is unknown resolves rather
		// than reading as a missing column.
		"SELECT a FROM SOURCE(TYPE = 's')",
		"SELECT e.a FROM 's'(FORMAT(TYPE = 'csv')) AS e",
		// A snapshot may read a table too.
		"SELECT * FROM t1 AS OF SNAPSHOT (SELECT max(a) FROM t1)",
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO TABLE t1",
		"LOAD DATA FROM(SELECT a FROM SOURCE(TYPE = 's')) INTO TABLE t1",
		// Without a column alias list the outputs are the arguments.
		"SELECT * FROM t1, UNNEST(t1.a)",
		"SELECT * FROM RB_ITERATE(t1.a)",
		// A format value may read a table.
		"SELECT a FROM 's'(FORMAT(FIELD_DELIMITER = (SELECT a FROM t1)))",
		// Remote DUAL over a database link is a remote relation, not local 0-column DUAL.
		"SELECT r.x FROM `dual`@remote_link AS r",
		"FOR COLUMNS",
	} {
		if _, err := resolve(t, statement); err != nil {
			t.Errorf("ResolveQuerySpan(%q): %v", statement, err)
		}
	}
}

// TestResolveQuerySpanMethodOptClassified pins that a method_opt statement has a
// permission class: Bytebase's SQL service denies a statement left unknown, and
// this one is now parsed.
func TestResolveQuerySpanMethodOptClassified(t *testing.T) {
	span, err := resolve(t, "FOR COLUMNS")
	if err != nil {
		t.Fatalf("ResolveQuerySpan: %v", err)
	}
	if span.Type == analysis.QueryTypeUnknown {
		t.Errorf("FOR COLUMNS is QueryTypeUnknown, want a class")
	}
}

// TestFlashbackAccessResources pins that FLASHBACK's restored table, and the
// name it is restored under, are access resources of the statement.
func TestFlashbackAccessResources(t *testing.T) {
	for _, tc := range []struct {
		statement string
		sources   string
	}{
		{"FLASHBACK TABLE db.t1 TO BEFORE DROP", "db.t1."},
		{"FLASHBACK TABLE db.t1 TO BEFORE DROP RENAME TO db.t2", "db.t1.,db.t2."},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolve(t, tc.statement)
			require.NoError(t, err)
			require.Equal(t, tc.sources, renderSet(span.SourceColumns))
		})
	}
}

// TestRemoteCatalogColumnDoesNotBindLocally pins that a four-part reference to
// another catalog never resolves to a same-named local table that is also in
// FROM.
func TestRemoteCatalogColumnDoesNotBindLocally(t *testing.T) {
	span, err := resolve(t, "SELECT c1.db.t1.a FROM c1.db.t1, db.t1")
	require.NoError(t, err)
	// The remote relation's columns are known only when the query runs, so the
	// reference reads no local column.
	require.Equal(t, []string{"a<-"}, renderResults(span))

	// The local reference, without a catalog, keeps the local table's column.
	span, err = resolve(t, "SELECT t1.a FROM c1.db.t1, db.t1")
	require.NoError(t, err)
	require.Equal(t, []string{"a<-db.t1.a plain"}, renderResults(span))
}

// TestRemoteAliasAndWildcardStayRemote pins that an alias-qualified reference to
// a remote relation resolves (its columns are unknown, not missing), and that a
// catalog-qualified wildcard never reports a same-named local table's columns.
func TestRemoteAliasAndWildcardStayRemote(t *testing.T) {
	span, err := resolve(t, "SELECT x.a FROM c1.db.t1 AS x")
	require.NoError(t, err)
	require.Equal(t, []string{"a<-"}, renderResults(span))

	span, err = resolve(t, "SELECT c1.db.t1.* FROM c1.db.t1, db.t1")
	require.NoError(t, err)
	require.Empty(t, span.Results)
	for resource := range span.SourceColumns {
		if resource.Table == "t1" && resource.Column != "" {
			t.Errorf("remote wildcard reported a local column: %s", renderSet(span.SourceColumns))
		}
	}
}

func TestRemoteUnknownDoesNotSwallowLocalLineage(t *testing.T) {
	// Unqualified and qualified references over catalog tables resolve.
	for _, sql := range []string{
		"SELECT a FROM c1.db.t1",
		"SELECT a FROM c1.db.t1 AS x",
		"SELECT (SELECT a) FROM c1.db.t1",
		"SELECT POSITION('b' IN a), ? FROM db.t1",
	} {
		span, err := resolve(t, sql)
		require.NoError(t, err, sql)
		require.Nil(t, span.NotFoundError, sql)
	}

	// Unknown remote/external relations must not wipe out local column lineage.
	for _, sql := range []string{
		"SELECT a FROM db.t1, other@remote_link",
		"SELECT (SELECT a FROM other@remote_link LIMIT 1) AS a FROM db.t1",
		"SELECT db.t1.a FROM db.t1, db.t1@remote_link",
	} {
		span, err := resolve(t, sql)
		require.NoError(t, err, sql)
		require.Nil(t, span.NotFoundError, sql)
		require.Equal(t, []string{"a<-db.t1.a plain"}, renderResults(span), sql)
	}

	// Qualified alias references (`db.x.a`, `c2.db.x.a`, `c1.db.t1.*` on aliased table) must not bind to the alias.
	for _, sql := range []string{
		"SELECT c2.db.x.a FROM c1.db.t1 AS x",
		"SELECT db.x.a FROM c1.db.t1 AS x",
		"SELECT c1.db.t1.* FROM c1.db.t1 AS x",
	} {
		span, err := resolve(t, sql)
		require.NoError(t, err, sql)
		require.NotNil(t, span.NotFoundError, sql)
	}
}
