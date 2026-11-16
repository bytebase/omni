package analysis

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest/cputime"
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

// TestRedshiftQuerySpanStarOverCrossJoin covers an ordinary joined star, whose
// join has no USING clause to read.
func TestRedshiftQuerySpanStarOverCrossJoin(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("SELECT * FROM (SELECT reporting_fn() AS x) a CROSS JOIN (SELECT 0 AS y) b")
	require.NoError(t, err)
	require.Len(t, span.Results, 2)
	require.Equal(t, ResolvedQuerySpanResult{
		Name: "x",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
		}),
		IsPlainField: true,
	}, span.Results[0])
	require.Equal(t, "y", span.Results[1].Name)
}

// TestRedshiftQuerySpanScalarSubqueryShadowsOuterBinding covers an inner FROM
// alias over a CTE that is also bound outside: the inner scope wins instead of
// making the unqualified column ambiguous.
func TestRedshiftQuerySpanScalarSubqueryShadowsOuterBinding(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH c AS (SELECT reporting_fn() AS x) SELECT (SELECT x FROM c AS cc) AS y FROM c")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "y",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
		}),
	}}, span.Results)
}

// TestRedshiftQuerySpanUsingJoinUnqualifiedColumn covers an unqualified common
// column: the join exposes it once, with both inputs' sources.
func TestRedshiftQuerySpanUsingJoinUnqualifiedColumn(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH a AS (SELECT reporting_fn() AS k), b AS (SELECT 0 AS k) SELECT k FROM a JOIN b USING (k)")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "k",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
		}),
		IsPlainField: true,
	}}, span.Results)
}

// TestRedshiftQuerySpanSubqueryPredicateFunction covers a subquery whose result
// depends on a SQL function only in its predicate.
func TestRedshiftQuerySpanSubqueryPredicateFunction(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	for _, statement := range []string{
		"SELECT EXISTS (SELECT 1 WHERE reporting_fn() > 0) AS x",
		"SELECT (SELECT 1 WHERE reporting_fn() > 0) AS x",
	} {
		t.Run(statement, func(t *testing.T) {
			span, err := q.getQuerySpan(statement)
			require.NoError(t, err)
			require.Equal(t, []ResolvedQuerySpanResult{{
				Name: "x",
				SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
					{Database: "db", Schema: "public", Table: "orders", Column: "id"},
				}),
			}}, span.Results)
		})
	}
}

// TestRedshiftQuerySpanOpaqueStarKeepsFollowingFunction covers a CTE whose
// projection mixes an expanded star with an unresolved base-table star; the
// unknown width must stop expansion so the following function keeps its source.
func TestRedshiftQuerySpanOpaqueStarKeepsFollowingFunction(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH a AS (SELECT reporting_fn() AS x, 0 AS n), b AS (SELECT a.*, o.* FROM a CROSS JOIN orders o) SELECT b.*, analytics.reporting_fn() AS y FROM b")
	require.NoError(t, err)
	last := span.Results[len(span.Results)-1]
	require.Equal(t, "y", last.Name)
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "db", Schema: "analytics", Table: "events", Column: "event_id"},
	}), last.SourceColumns)
}

// TestRedshiftQuerySpanSubquerySetOpFunction covers a set operation inside a
// scalar subquery: each arm resolves a function-producing CTE in its own scope.
func TestRedshiftQuerySpanSubquerySetOpFunction(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH c AS (SELECT reporting_fn() AS x) SELECT (SELECT x FROM c UNION ALL SELECT 0 LIMIT 1) AS y")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "y",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
		}),
	}}, span.Results)
}

// TestRedshiftQuerySpanSetOpOpaqueStarKeepsFollowingFunction covers a UNION CTE
// with an unresolved star: its width stays unknown, so a later function keeps
// its own source instead of taking the first arm's.
func TestRedshiftQuerySpanSetOpOpaqueStarKeepsFollowingFunction(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH b AS (SELECT reporting_fn() AS x, o.* FROM orders o UNION ALL SELECT reporting_fn(), o.* FROM orders o) SELECT b.*, analytics.reporting_fn() AS y FROM b")
	require.NoError(t, err)
	last := span.Results[len(span.Results)-1]
	require.Equal(t, "y", last.Name)
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "db", Schema: "analytics", Table: "events", Column: "event_id"},
	}), last.SourceColumns)
}

// TestRedshiftQuerySpanMixedStarKeepsFollowingFunction covers a final SELECT
// that mixes an expanded CTE star with an unresolved base-table star: the
// later function is matched by name, not given the first empty result.
func TestRedshiftQuerySpanMixedStarKeepsFollowingFunction(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH a AS (SELECT reporting_fn() AS x, 0 AS n) SELECT a.*, o.*, analytics.reporting_fn() AS y FROM a CROSS JOIN orders o")
	require.NoError(t, err)
	last := span.Results[len(span.Results)-1]
	require.Equal(t, "y", last.Name)
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "db", Schema: "analytics", Table: "events", Column: "event_id"},
	}), last.SourceColumns)
}

// TestRedshiftQuerySpanBaseStarKeepsUnnamedFunction covers a base-table star
// followed by an unaliased function expression: the star now expands to its
// metadata width, so the expression's sources stay on the last result.
func TestRedshiftQuerySpanBaseStarKeepsUnnamedFunction(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("SELECT o.*, reporting_fn() + 1 FROM orders o")
	require.NoError(t, err)
	last := span.Results[len(span.Results)-1]
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "db", Schema: "public", Table: "orders", Column: "id"},
	}), last.SourceColumns)
}

// TestRedshiftQuerySpanPartialColumnAliasKeepsWidth covers a table alias that
// renames only a prefix: the remaining metadata columns still count toward the
// star's width, so the trailing function stays aligned.
func TestRedshiftQuerySpanPartialColumnAliasKeepsWidth(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("SELECT o.*, reporting_fn() + 1 FROM orders o(first_id)")
	require.NoError(t, err)
	require.Len(t, span.Results, 6)
	last := span.Results[len(span.Results)-1]
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "db", Schema: "public", Table: "orders", Column: "id"},
	}), last.SourceColumns)
}

// TestRedshiftQuerySpanPartialCTEAliasKeepsWidth covers a partial alias on a
// CTE reference: the remaining projection names must survive so a later
// function stays aligned.
func TestRedshiftQuerySpanPartialCTEAliasKeepsWidth(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH c AS (SELECT 0 AS n, reporting_fn() AS x), b AS (SELECT c.*, analytics.reporting_fn() AS y FROM c AS c(nn)) SELECT y FROM b")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "y",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "analytics", Table: "events", Column: "event_id"},
		}),
		IsPlainField: true,
	}}, span.Results)
}

// TestRedshiftQuerySpanCorrelatedReferenceUsesNearestScope covers a correlated
// reference that skips an enclosing query with no matching column: the nearest
// scope that has the column wins over older same-named bindings.
func TestRedshiftQuerySpanCorrelatedReferenceUsesNearestScope(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH a AS (SELECT 0 AS x), b AS (SELECT reporting_fn() AS x) SELECT (SELECT (SELECT (SELECT x) FROM orders) FROM b) AS y FROM a")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "y",
		SourceColumns: redshiftSourceColumnSetFromList([]ColumnResource{
			{Database: "db", Schema: "public", Table: "orders", Column: "id"},
		}),
	}}, span.Results)
}

// TestRedshiftQuerySpanSchemaQualifiedTableShadowsCTE covers a schema-qualified
// relation whose name matches a CTE: the physical table, not the CTE, supplies
// the star width.
func TestRedshiftQuerySpanSchemaQualifiedTableShadowsCTE(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH orders AS (SELECT 0 AS n) SELECT o.*, reporting_fn() + 1 FROM public.orders o")
	require.NoError(t, err)
	require.Len(t, span.Results, 6)
	last := span.Results[len(span.Results)-1]
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "db", Schema: "public", Table: "orders", Column: "id"},
	}), last.SourceColumns)
}

// TestRedshiftQuerySpanMaterializedViewStarKeepsFunction covers a materialized
// view star whose width comes from its definition, since the metadata has no
// column list.
func TestRedshiftQuerySpanMaterializedViewStarKeepsFunction(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("SELECT m.*, reporting_fn() + 1 FROM order_totals_mv m")
	require.NoError(t, err)
	require.Len(t, span.Results, 3)
	last := span.Results[len(span.Results)-1]
	require.Equal(t, redshiftSourceColumnSetFromList([]ColumnResource{
		{Database: "db", Schema: "public", Table: "orders", Column: "id"},
	}), last.SourceColumns)
}

// TestRedshiftQuerySpanScalarSubqueryStarThroughDerived covers a scalar
// subquery selecting a star over a derived relation that renames a function
// column.
func TestRedshiftQuerySpanScalarSubqueryStarThroughDerived(t *testing.T) {
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	span, err := q.getQuerySpan("WITH c AS (SELECT reporting_fn() AS x) SELECT (SELECT * FROM (SELECT x AS y FROM c) d) AS z")
	require.NoError(t, err)
	require.Equal(t, []ResolvedQuerySpanResult{{
		Name: "z",
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

// TestRedshiftFunctionSourceManyResolvedJoinsStayLinear: the function-source
// projection is built only when a star or unqualified column needs it, and an
// ordinary join chain is flattened once, so a generated query of resolved CTE
// joins does not copy an ever-growing prefix at every node.
func TestRedshiftFunctionSourceManyResolvedJoinsStayLinear(t *testing.T) {
	const joins = 10000
	parts := make([]string, joins)
	for i := range parts {
		parts[i] = fmt.Sprintf("c a%d", i)
	}
	sql := "WITH c AS (SELECT 1 AS x) SELECT 1 FROM " + strings.Join(parts, " CROSS JOIN ")
	stmts, err := parseRedshiftSpan(sql)
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	cputime.Within(t, 2*time.Second, fmt.Sprintf("function sources over %d resolved CTE joins", joins), func() {
		_, _, _, err := q.collectFunctionBodySources(stmts[0].AST, map[string]bool{}, make(map[string]functionBodySources))
		require.NoError(t, err)
	})
}

// TestRedshiftFunctionSourceNestedScalarSubqueriesStayLinear: each nested
// scalar subquery is analyzed once, not once through its target and again
// through the whole SELECT, which made depth exponential.
func TestRedshiftFunctionSourceNestedScalarSubqueriesStayLinear(t *testing.T) {
	const depth = 1000
	sql := "SELECT " + strings.Repeat("(", depth) + "1" + strings.Repeat(")", depth)
	stmts, err := parseRedshiftSpan(sql)
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	q := newQuerySpanResolver("db", []string{"public"}, redshiftQuerySpanContext(t))
	cputime.Within(t, 2*time.Second, fmt.Sprintf("function sources through %d nested scalar subqueries", depth), func() {
		_, _, _, err := q.collectFunctionBodySources(stmts[0].AST, map[string]bool{}, make(map[string]functionBodySources))
		require.NoError(t, err)
	})
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
