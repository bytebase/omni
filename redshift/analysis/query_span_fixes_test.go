package analysis

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/redshift"
	"github.com/stretchr/testify/require"
)

// TestRedshiftQuerySpanFunctionLineageMergesSetOpArms covers defect 6: a
// SQL-function result in a set-operation arm contributes its lineage to the
// merged output column.
func TestRedshiftQuerySpanFunctionLineageMergesSetOpArms(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("SELECT reporting_fn() AS x UNION SELECT analytics.reporting_fn() AS x")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "x",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
			{Database: "db", Schema: "analytics", Table: "events", Column: "event_id"},
		}),
	}}, span.Results)
}

// TestRedshiftQuerySpanStarOverFunctionCTE covers defect 7: `SELECT *` over a
// CTE that projects a SQL function keeps the function's result lineage.
func TestRedshiftQuerySpanStarOverFunctionCTE(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH c AS (SELECT reporting_fn() AS x) SELECT * FROM c")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "x",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
		}),
		IsPlainField: true,
	}}, span.Results)
}

// TestRedshiftQuerySpanStarOverUsingJoinKeepsFunctionAligned covers a star over
// a USING join: its common column is listed once, so a following function still
// lines up with the output projection.
func TestRedshiftQuerySpanStarOverUsingJoinKeepsFunctionAligned(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH a AS (SELECT reporting_fn() AS k), b AS (SELECT 0 AS k) SELECT *, analytics.reporting_fn() AS y FROM a JOIN b USING (k)")
	require.NoError(t, err)
	require.Len(t, span.Results, 2)
	require.Equal(t, "k", span.Results[0].Name)
	require.Equal(t, ColumnSet{
		{Database: "db", Schema: "public", Table: "orders", Column: "id"}: true,
	}, span.Results[0].SourceColumns)
	require.Equal(t, ResolvedQuerySpanResult{
		Name: "y",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "analytics", Table: "events", Column: "event_id"},
		}),
	}, span.Results[1])
}

// TestRedshiftQuerySpanScalarSubqueryFunctionCTE covers a function-producing CTE
// referenced only inside a scalar subquery: the subquery's own FROM binding must
// expose it.
func TestRedshiftQuerySpanScalarSubqueryFunctionCTE(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH c AS (SELECT reporting_fn() AS x) SELECT (SELECT x FROM c) AS y")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "y",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
		}),
	}}, span.Results)
}

// TestRedshiftQuerySpanStarKeepsFollowingFunctionAligned covers a projection
// with a multi-column star followed by a function: the function's output name
// and sources must stay aligned after the star expands.
func TestRedshiftQuerySpanStarKeepsFollowingFunctionAligned(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH a AS (SELECT 0 AS n, reporting_fn() AS x), b AS (SELECT *, analytics.reporting_fn() AS y FROM a) SELECT y FROM b")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "y",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "analytics", Table: "events", Column: "event_id"},
		}),
		IsPlainField: true,
	}}, span.Results)
}

// TestRedshiftQuerySpanAliasedCTEIsNotAmbiguous covers defect 8: an aliased CTE
// is visible only under its alias, so an unqualified column is not ambiguous.
func TestRedshiftQuerySpanAliasedCTEIsNotAmbiguous(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	for _, statement := range []string{
		"WITH c AS (SELECT reporting_fn() AS x) SELECT x FROM c AS cc",
		"WITH c AS (SELECT reporting_fn() AS x) SELECT cc.x FROM c AS cc",
	} {
		t.Run(statement, func(t *testing.T) {
			span, err := q.getQuerySpan(statement)
			require.NoError(t, err)
			require.Equal(t, []ResolvedQuerySpanResult{{
				Name: "x",
				SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
					{Database: "db", Schema: "public", Table: "orders", Column: "id"},
				}),
				IsPlainField: true,
			}}, span.Results)
		})
	}
}

// TestRedshiftQuerySpanCrossDatabaseSourceColumn covers defect 9: a source
// column from a cross-database range variable carries that database, not the
// session default.
func TestRedshiftQuerySpanCrossDatabaseSourceColumn(t *testing.T) {
	other := &metadata.DatabaseSchemaMetadata{
		Name: "other_db",
		Schemas: []*metadata.SchemaMetadata{{
			Name: "public",
			Tables: []*metadata.TableMetadata{{
				Name:    "orders",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "amount", Type: "int"}},
			}},
		}},
	}
	getter, _ := redshiftMockDatabaseMetadataGetter([]*metadata.DatabaseSchemaMetadata{redshiftQuerySpanMetadata(), other})
	q := newQuerySpanResolver("db", []string{"public"}, ResolveOptions{
		DefaultDatabase:     "db",
		SearchPath:          []string{"public"},
		GetDatabaseMetadata: getter,
	})
	span, err := q.getQuerySpan("SELECT id FROM other_db.public.orders")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "id",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "other_db", Schema: "public", Table: "orders", Column: "id"},
		}),
		IsPlainField: true,
	}}, span.Results)
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "other_db", Schema: "public", Table: "orders"},
	}), span.SourceColumns)
}

// TestRedshiftQuerySpanParsesStatementOnce covers defect 10: classification and
// analysis share one parse.
func TestRedshiftQuerySpanParsesStatementOnce(t *testing.T) {
	old := parseRedshiftSpan
	t.Cleanup(func() { parseRedshiftSpan = old })
	calls := 0
	parseRedshiftSpan = func(sql string) ([]redshift.Statement, error) {
		calls++
		return old(sql)
	}
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	_, err := q.getQuerySpan("SELECT id FROM orders")
	require.NoError(t, err)
	require.Equal(t, 1, calls)
}
