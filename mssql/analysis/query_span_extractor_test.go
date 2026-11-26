package analysis

import (
	"cmp"
	"errors"
	"slices"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

var extractorTestMetadata = []*metadatapb.DatabaseSchemaMetadata{
	{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "dbo",
			Tables: []*metadatapb.TableMetadata{
				{Name: "t", Columns: columns("a", "b", "c")},
				{Name: "t1", Columns: columns("a", "b", "c")},
				{Name: "t2", Columns: columns("a", "b")},
				{Name: "ident_t", Columns: []*metadatapb.ColumnMetadata{{Name: "id", IsIdentity: true}, {Name: "payload"}}},
				{Name: "ident_t2", Columns: []*metadatapb.ColumnMetadata{{Name: "seq", IsIdentity: true}, {Name: "v"}}},
				{Name: "bare_edge"},
				{Name: "weird", Columns: columns("$node_id", "IDENTITYCOL")},
			},
			Views: []*metadatapb.ViewMetadata{
				{Name: "vw", Definition: "CREATE VIEW [dbo].[vw] AS SELECT a, b FROM t"},
			},
		}},
	},
	{
		Name: "db2",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "dbo",
			Tables: []*metadatapb.TableMetadata{
				{Name: "t2", Columns: columns("x", "y")},
			},
		}},
	},
}

type extractorSession struct {
	defaultDatabase string
	catalog         Catalog
	tempTables      map[string]*TempTable
}

func newExtractorSession(defaultDatabase string) *extractorSession {
	return &extractorSession{
		defaultDatabase: defaultDatabase,
		catalog:         catalogFromMetadata(extractorTestMetadata),
		tempTables:      make(map[string]*TempTable),
	}
}

func (s *extractorSession) resolve(statement string) (*ResolvedSpan, error) {
	return ResolveQuerySpan(statement, ResolveOptions{
		DefaultDatabase:     s.defaultDatabase,
		DefaultSchema:       "dbo",
		IgnoreCaseSensitive: true,
		Catalog:             s.catalog,
		TempTables:          s.tempTables,
	})
}

func sortedResources(set ColumnSet) []ColumnResource {
	out := make([]ColumnResource, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	slices.SortFunc(out, func(a, b ColumnResource) int { return cmp.Compare(a.String(), b.String()) })
	return out
}

func TestQuerySpanExtractorSupportedShapes(t *testing.T) {
	type expectedColumn struct {
		name    string
		sources []ColumnResource
	}
	src := func(db, schema, table, col string) ColumnResource {
		return ColumnResource{Database: db, Schema: schema, Table: table, Column: col}
	}
	access := func(db, schema, table string) ColumnResource {
		return ColumnResource{Database: db, Schema: schema, Table: table}
	}

	cases := []struct {
		name             string
		sql              string
		wantResults      []expectedColumn
		wantPredicate    []ColumnResource
		wantAccessTables []ColumnResource
	}{
		{
			name:             "bare_column",
			sql:              "SELECT a FROM t",
			wantResults:      []expectedColumn{{"a", []ColumnResource{src("db", "dbo", "t", "a")}}},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name:             "qualified_column",
			sql:              "SELECT t.a FROM t",
			wantResults:      []expectedColumn{{"a", []ColumnResource{src("db", "dbo", "t", "a")}}},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name:             "alias_as",
			sql:              "SELECT a AS x FROM t",
			wantResults:      []expectedColumn{{"x", []ColumnResource{src("db", "dbo", "t", "a")}}},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name: "multi_columns",
			sql:  "SELECT a, b FROM t",
			wantResults: []expectedColumn{
				{"a", []ColumnResource{src("db", "dbo", "t", "a")}},
				{"b", []ColumnResource{src("db", "dbo", "t", "b")}},
			},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name: "star_unqualified",
			sql:  "SELECT * FROM t",
			wantResults: []expectedColumn{
				{"a", []ColumnResource{src("db", "dbo", "t", "a")}},
				{"b", []ColumnResource{src("db", "dbo", "t", "b")}},
				{"c", []ColumnResource{src("db", "dbo", "t", "c")}},
			},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name:             "table_alias_qualified",
			sql:              "SELECT x.a FROM t AS x",
			wantResults:      []expectedColumn{{"a", []ColumnResource{src("db", "dbo", "t", "a")}}},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name: "table_alias_star",
			sql:  "SELECT x.* FROM t AS x",
			wantResults: []expectedColumn{
				{"a", []ColumnResource{src("db", "dbo", "t", "a")}},
				{"b", []ColumnResource{src("db", "dbo", "t", "b")}},
				{"c", []ColumnResource{src("db", "dbo", "t", "c")}},
			},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name:             "where_predicate_columns",
			sql:              "SELECT a FROM t WHERE b > 1 AND c = 2",
			wantResults:      []expectedColumn{{"a", []ColumnResource{src("db", "dbo", "t", "a")}}},
			wantPredicate:    []ColumnResource{src("db", "dbo", "t", "b"), src("db", "dbo", "t", "c")},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name:             "cross_database_table",
			sql:              "SELECT x FROM db2.dbo.t2",
			wantResults:      []expectedColumn{{"x", []ColumnResource{src("db2", "dbo", "t2", "x")}}},
			wantAccessTables: []ColumnResource{access("db2", "dbo", "t2")},
		},
		{
			name:             "explicit_schema",
			sql:              "SELECT a FROM dbo.t",
			wantResults:      []expectedColumn{{"a", []ColumnResource{src("db", "dbo", "t", "a")}}},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name: "star_with_alias",
			sql:  "SELECT y.* FROM t2 AS y",
			wantResults: []expectedColumn{
				{"a", []ColumnResource{src("db", "dbo", "t2", "a")}},
				{"b", []ColumnResource{src("db", "dbo", "t2", "b")}},
			},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t2")},
		},
		{
			name: "pivot",
			sql:  "SELECT * FROM t PIVOT (SUM(a) FOR b IN ([1], [2])) AS p",
			wantResults: []expectedColumn{
				{"c", []ColumnResource{src("db", "dbo", "t", "c")}},
				{"1", []ColumnResource{src("db", "dbo", "t", "a"), src("db", "dbo", "t", "b")}},
				{"2", []ColumnResource{src("db", "dbo", "t", "a"), src("db", "dbo", "t", "b")}},
			},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
		{
			name: "pivot_over_join",
			sql:  "SELECT p.a, p.[1] FROM t JOIN db2.dbo.t2 AS u ON t.a = u.x PIVOT (MAX(y) FOR b IN ([1])) AS p",
			wantResults: []expectedColumn{
				{"a", []ColumnResource{src("db", "dbo", "t", "a")}},
				{"1", []ColumnResource{src("db", "dbo", "t", "b"), src("db2", "dbo", "t2", "y")}},
			},
			wantPredicate:    []ColumnResource{src("db", "dbo", "t", "a"), src("db2", "dbo", "t2", "x")},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t"), access("db2", "dbo", "t2")},
		},
		{
			name: "unpivot",
			sql:  "SELECT * FROM t UNPIVOT (v FOR k IN (a, b)) AS u",
			wantResults: []expectedColumn{
				{"c", []ColumnResource{src("db", "dbo", "t", "c")}},
				{"v", []ColumnResource{src("db", "dbo", "t", "a"), src("db", "dbo", "t", "b")}},
				{"k", nil},
			},
			wantPredicate:    []ColumnResource{src("db", "dbo", "t", "a"), src("db", "dbo", "t", "b")},
			wantAccessTables: []ColumnResource{access("db", "dbo", "t")},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			span, err := newExtractorSession("db").resolve(tc.sql)
			require.NoError(t, err)
			require.Equal(t, Select, span.Type)
			require.Len(t, span.Results, len(tc.wantResults))
			for i, want := range tc.wantResults {
				require.Equal(t, want.name, span.Results[i].Name)
				require.ElementsMatch(t, want.sources, sortedResources(span.Results[i].SourceColumns))
			}
			require.ElementsMatch(t, tc.wantPredicate, sortedResources(span.PredicateColumns))
			require.ElementsMatch(t, tc.wantAccessTables, sortedResources(span.SourceColumns))
		})
	}
}

func TestQuerySpanExtractorUnsupportedTableSources(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
	}{
		{"cross_apply_tvf", "SELECT * FROM t1 CROSS APPLY fn(t1.a) AS x(v)"},
		{"aliased_tvf", "SELECT * FROM fn(t1.a) AS x(v)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newExtractorSession("db").resolve(tc.sql)
			require.Error(t, err)
			var typeNotSupported *TypeNotSupportedError
			require.True(t, errors.As(err, &typeNotSupported))
		})
	}
}

func TestQuerySpanExtractorSystemTableValuedFunctions(t *testing.T) {
	wantT1A := []ColumnResource{{Database: "db", Schema: "dbo", Table: "t1", Column: "a"}}
	for _, tc := range []struct {
		name     string
		sql      string
		wantName string
	}{
		{"string_split_value_flows_from_input", "SELECT s.value FROM t1 CROSS APPLY STRING_SPLIT(t1.a, ',') AS s", "value"},
		{"openjson_value_flows_from_input", "SELECT j.value FROM t1 CROSS APPLY OPENJSON(t1.a) AS j", "value"},
		{"tvf_column_alias_list", "SELECT s.v FROM t1 CROSS APPLY STRING_SPLIT(t1.a, ',', 1) AS s(v, ordinal)", "v"},
		{"string_split_third_arg_zero_has_no_ordinal", "SELECT s.value FROM t1 CROSS APPLY STRING_SPLIT(t1.a, ',', 0) AS s(value)", "value"},
		{"string_split_third_arg_null_has_no_ordinal", "SELECT s.value FROM t1 CROSS APPLY STRING_SPLIT(t1.a, ',', NULL) AS s(value)", "value"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			span, err := newExtractorSession("db").resolve(tc.sql)
			require.NoError(t, err)
			require.Len(t, span.Results, 1)
			require.Equal(t, tc.wantName, span.Results[0].Name)
			require.ElementsMatch(t, wantT1A, sortedResources(span.Results[0].SourceColumns))
		})
	}
}

func TestQuerySpanExtractorTempTableRoundTrips(t *testing.T) {
	t.Run("create temp table case insensitive", func(t *testing.T) {
		s := newExtractorSession("db")
		span1, err := s.resolve("CREATE TABLE #Tmp(id INT, name NVARCHAR(50))")
		require.NoError(t, err)
		require.Equal(t, DDL, span1.Type)
		span2, err := s.resolve("SELECT id, name FROM #tmp")
		require.NoError(t, err)
		require.Len(t, span2.Results, 2)
		require.Equal(t, "id", span2.Results[0].Name)
		require.Equal(t, "name", span2.Results[1].Name)
	})

	t.Run("select into temp table registers columns", func(t *testing.T) {
		s := newExtractorSession("db")
		span1, err := s.resolve("SELECT a AS id, b INTO #selected FROM t")
		require.NoError(t, err)
		require.Equal(t, Select, span1.Type)
		require.Equal(t, []string{"id", "b"}, s.tempTables["#selected"].Columns)
		span2, err := s.resolve("SELECT * FROM #selected")
		require.NoError(t, err)
		require.Len(t, span2.Results, 2)
		require.Equal(t, "id", span2.Results[0].Name)
		require.Equal(t, "b", span2.Results[1].Name)
	})
}

func TestQuerySpanExtractorRegressions(t *testing.T) {
	t.Run("cte visible in set op arms", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("WITH cte AS (SELECT a FROM t) SELECT a FROM cte UNION SELECT a FROM cte")
		require.NoError(t, err)
		require.Nil(t, span.NotFoundError)
		require.Len(t, span.Results, 1)
		require.Equal(t, "a", span.Results[0].Name)
	})

	t.Run("not found propagates from subqueries and recursive CTE arm", func(t *testing.T) {
		for _, sql := range []string{
			"SELECT a FROM t WHERE EXISTS (SELECT 1 FROM no_such_table)",
			"WITH rec AS (SELECT a FROM t UNION ALL SELECT a FROM rec JOIN no_such_table nx ON rec.a = nx.a) SELECT * FROM rec",
			"SELECT a FROM t WHERE a > ANY (SELECT x FROM no_such_table)",
		} {
			span, err := newExtractorSession("db").resolve(sql)
			require.NoError(t, err, sql)
			require.NotNil(t, span.NotFoundError, sql)
			require.Empty(t, span.Results, sql)
		}
	})

	t.Run("apply values clause merges rows", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT x FROM t1 CROSS APPLY (VALUES (t1.a), (t1.b)) AS v(x)")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.ElementsMatch(t, []ColumnResource{
			{Database: "db", Schema: "dbo", Table: "t1", Column: "a"},
			{Database: "db", Schema: "dbo", Table: "t1", Column: "b"},
		}, sortedResources(span.Results[0].SourceColumns))
	})

	t.Run("window and within group order by sources", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT SUM(a) OVER (ORDER BY b) FROM t")
		require.NoError(t, err)
		require.ElementsMatch(t, []ColumnResource{
			{Database: "db", Schema: "dbo", Table: "t", Column: "a"},
			{Database: "db", Schema: "dbo", Table: "t", Column: "b"},
		}, sortedResources(span.Results[0].SourceColumns))

		span, err = newExtractorSession("db").resolve("SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY b) FROM t")
		require.NoError(t, err)
		require.ElementsMatch(t, []ColumnResource{
			{Database: "db", Schema: "dbo", Table: "t", Column: "b"},
		}, sortedResources(span.Results[0].SourceColumns))
	})

	t.Run("in list access tables", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT a FROM t WHERE a IN (ABS((SELECT x FROM db2.dbo.t2)))")
		require.NoError(t, err)
		require.ElementsMatch(t, []ColumnResource{
			{Database: "db", Schema: "dbo", Table: "t"},
			{Database: "db2", Schema: "dbo", Table: "t2"},
		}, sortedResources(span.SourceColumns))
	})

	t.Run("correlated exists with dateadd datepart", func(t *testing.T) {
		sql := `
SELECT outer_t.a
FROM t AS outer_t
WHERE outer_t.c > DATEADD(HOUR, 2, GETUTCDATE())
  AND NOT EXISTS (
    SELECT 1
    FROM t1 AS inner_t
    WHERE inner_t.a = outer_t.a
      AND inner_t.b > outer_t.b
  )
ORDER BY outer_t.a
`
		span, err := newExtractorSession("db").resolve(sql)
		require.NoError(t, err)
		require.ElementsMatch(t, []ColumnResource{
			{Database: "db", Schema: "dbo", Table: "t", Column: "a"},
			{Database: "db", Schema: "dbo", Table: "t", Column: "b"},
			{Database: "db", Schema: "dbo", Table: "t", Column: "c"},
			{Database: "db", Schema: "dbo", Table: "t1", Column: "a"},
			{Database: "db", Schema: "dbo", Table: "t1", Column: "b"},
		}, sortedResources(span.PredicateColumns))
	})

	t.Run("unresolved column errors", func(t *testing.T) {
		for _, sql := range []string{
			"SELECT no_such_col FROM t",
			"SELECT t.no_such_col FROM t",
			"SELECT a FROM t WHERE no_such_col = 1",
			"SELECT t.a FROM t JOIN t1 ON t1.no_such_col = t.a",
			"SELECT SUM(a) OVER (PARTITION BY no_such_col) FROM t",
			"SELECT SUM(a) OVER (ORDER BY no_such_col) FROM t",
			"SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY no_such_col) FROM t",
			"SELECT CASE WHEN CONTAINS(no_such_col, 'x') THEN 1 ELSE 0 END FROM t",
		} {
			_, err := newExtractorSession("db").resolve(sql)
			require.Errorf(t, err, "unresolved column %q should error", sql)
		}
	})

	t.Run("empty and GO", func(t *testing.T) {
		for _, sql := range []string{"", "GO", ";"} {
			span, err := newExtractorSession("db").resolve(sql)
			require.NoError(t, err, sql)
			require.Equal(t, Select, span.Type, sql)
			require.Empty(t, span.Results, sql)
		}
	})
}
