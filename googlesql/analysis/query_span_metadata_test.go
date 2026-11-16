package analysis

import (
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

// These pins cover masking-leak shapes that CANNOT live in the legacy-recorded
// corpora under testdata/query-span, because the legacy ANTLR resolvers either
// error on them (fail-closed) or produce lineage that is itself a known leak.
// For each, the resolved span must carry CORRECT lineage — the structural rule
// for shapes legacy cannot resolve: correct lineage or a closed failure, never
// silently empty or misaligned attribution (Bytebase's masker is positional and
// fails open).

func resolvePin(t *testing.T, opts ResolveOptions, statement, defaultDatabase string, schemas []*metadata.SchemaMetadata) *ResolvedQuerySpan {
	t.Helper()
	dbMetadata := &metadata.DatabaseSchemaMetadata{Name: defaultDatabase, Schemas: schemas}
	caseOpts := opts
	caseOpts.DefaultDatabase = defaultDatabase
	caseOpts.GetDatabaseMetadata = corpusGetter(dbMetadata)
	span, err := ResolveQuerySpan(statement, caseOpts)
	require.NoError(t, err, "statement: %s", statement)
	return span
}

func pinTables(tables ...*metadata.TableMetadata) []*metadata.SchemaMetadata {
	return []*metadata.SchemaMetadata{{Name: "", Tables: tables}}
}

// sourcesOf renders a result's lineage as sorted "[schema.]table.column" names.
func sourcesOf(r ResolvedQuerySpanResult) []string {
	out := make([]string, 0, len(r.SourceColumns))
	for c := range r.SourceColumns {
		name := c.Table + "." + c.Column
		if c.Schema != "" {
			name = c.Schema + "." + name
		}
		out = append(out, name)
	}
	slices.Sort(out)
	return out
}

// TestLeakPin_LowercaseUsingCoalesce: `SELECT * FROM a JOIN b USING (k)` with a
// lowercase key. Real BigQuery coalesces USING keys case-insensitively (3 output
// columns), but the LEGACY resolver only coalesced upper-case-written keys — its
// lowercase output was 4 uncoalesced columns, shifting every later position
// against the executed result. The coalesced key reads BOTH sides.
func TestLeakPin_LowercaseUsingCoalesce(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT * FROM a JOIN b USING (k);", "ds1", pinTables(
		&metadata.TableMetadata{Name: "a", Columns: []*metadata.ColumnMetadata{{Name: "k"}, {Name: "x"}}},
		&metadata.TableMetadata{Name: "b", Columns: []*metadata.ColumnMetadata{{Name: "k"}, {Name: "y"}}},
	))
	require.Len(t, span.Results, 3, "USING key must be coalesced: real output is [k, x, y]")
	require.Equal(t, []string{"a.k", "b.k"}, sourcesOf(span.Results[0]), "coalesced key reads both sides")
	require.Equal(t, []string{"a.x"}, sourcesOf(span.Results[1]))
	require.Equal(t, []string{"b.y"}, sourcesOf(span.Results[2]))
}

// TestLeakPin_SpannerLowercaseUsingCoalesce: legacy spanner only coalesced a
// fully-upper-case-written USING token; the lowercase output kept BOTH key
// columns — a positional masking leak. The key keeps the field's metadata case.
func TestLeakPin_SpannerLowercaseUsingCoalesce(t *testing.T) {
	span := resolvePin(t, spannerResolveOptions(), "SELECT * FROM users JOIN logins USING (id);", "db", pinTables(
		&metadata.TableMetadata{Name: "users", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "email"}}},
		&metadata.TableMetadata{Name: "logins", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "secret"}}},
	))
	require.Len(t, span.Results, 3, "USING key must be coalesced: real output is [id, email, secret]")
	require.Equal(t, "id", span.Results[0].Name, "the coalesced key keeps the field's metadata case")
	require.Equal(t, []string{"logins.id", "users.id"}, sourcesOf(span.Results[0]), "coalesced key reads both sides")
	require.Equal(t, []string{"users.email"}, sourcesOf(span.Results[1]))
	require.Equal(t, []string{"logins.secret"}, sourcesOf(span.Results[2]))
}

// TestLeakPin_UnnestLineage: the UNNEST output column's lineage is the unnested
// array column. The legacy resolver errored on UNNEST in FROM (fail-closed); a
// silently empty lineage would return the sensitive array elements unmasked.
func TestLeakPin_UnnestLineage(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT elem FROM victim, UNNEST(victim.secret_tokens) AS elem", "ds1", pinTables(
		&metadata.TableMetadata{Name: "victim", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "secret_tokens"}}},
	))
	require.Len(t, span.Results, 1)
	require.Equal(t, []string{"victim.secret_tokens"}, sourcesOf(span.Results[0]),
		"UNNEST element column must carry the array column's lineage")
}

// TestLeakPin_SpannerUnnestLineage: same shape under the spanner dialect.
func TestLeakPin_SpannerUnnestLineage(t *testing.T) {
	span := resolvePin(t, spannerResolveOptions(), "SELECT elem FROM victim, UNNEST(victim.secret_tokens) AS elem", "db", pinTables(
		&metadata.TableMetadata{Name: "victim", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "secret_tokens"}}},
	))
	require.Len(t, span.Results, 1)
	require.Equal(t, []string{"victim.secret_tokens"}, sourcesOf(span.Results[0]),
		"UNNEST element column must carry the array column's lineage")
}

// TestLeakPin_UnnestWithOffsetNoLineage: the WITH OFFSET companion column is
// positional metadata, not data — empty lineage is correct for it, while the
// element column keeps the array lineage.
func TestLeakPin_UnnestWithOffsetNoLineage(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(),
		"SELECT elem, pos FROM victim, UNNEST(victim.secret_tokens) AS elem WITH OFFSET AS pos", "ds1", pinTables(
			&metadata.TableMetadata{Name: "victim", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "secret_tokens"}}},
		))
	require.Len(t, span.Results, 2)
	require.Equal(t, []string{"victim.secret_tokens"}, sourcesOf(span.Results[0]))
	require.Empty(t, sourcesOf(span.Results[1]), "OFFSET is positional, carries no data lineage")
}

// TestLeakPin_ByNameMerge: `UNION ALL BY NAME` aligns set-op arms by column
// NAME, not ordinal. The legacy parser errors on BY NAME (fail-closed); an
// ordinal merge would attribute the right arm's sensitive column to the wrong
// output position.
func TestLeakPin_ByNameMerge(t *testing.T) {
	for name, opts := range map[string]ResolveOptions{"bigquery": bigqueryResolveOptions(), "spanner": spannerResolveOptions()} {
		t.Run(name, func(t *testing.T) {
			span := resolvePin(t, opts, "SELECT id, label FROM lt UNION ALL BY NAME SELECT b_secret AS label, id FROM rt", "ds1", pinTables(
				&metadata.TableMetadata{Name: "lt", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "label"}}},
				&metadata.TableMetadata{Name: "rt", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "b_secret"}}},
			))
			require.Len(t, span.Results, 2)
			require.Equal(t, []string{"lt.id", "rt.id"}, sourcesOf(span.Results[0]),
				"id merges both arms' id lineage by NAME")
			require.Equal(t, []string{"lt.label", "rt.b_secret"}, sourcesOf(span.Results[1]),
				"label carries the right arm's b_secret lineage (aliased AS label)")
		})
	}
}

// TestLeakPin_ByNameStarArmReversedOrder: a BY NAME set-op whose left arm is a
// star takes the deferred SetOpMerge path (metadata expansion first), which must
// name-merge too — an ordinal merge mis-attributes whenever the arms' column
// orders differ.
func TestLeakPin_ByNameStarArmReversedOrder(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT * FROM pub UNION ALL BY NAME SELECT secret AS label, id FROM priv", "ds1", pinTables(
		&metadata.TableMetadata{Name: "pub", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "label"}}},
		&metadata.TableMetadata{Name: "priv", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "secret"}}},
	))
	require.Len(t, span.Results, 2)
	require.Equal(t, []string{"priv.id", "pub.id"}, sourcesOf(span.Results[0]),
		"id merges by NAME despite the right arm listing it second")
	require.Equal(t, []string{"priv.secret", "pub.label"}, sourcesOf(span.Results[1]),
		"label carries priv.secret (aliased AS label) despite it being first in the right arm")
}

// TestLeakPin_ByNameOnMatchColumns: `BY NAME ON (cols)` outputs ONLY the listed
// columns, in list order — emitting the full ordinal merge would shift the
// positional masker off every real output column.
func TestLeakPin_ByNameOnMatchColumns(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT * FROM pub UNION ALL BY NAME ON (label) SELECT id, secret AS label FROM priv", "ds1", pinTables(
		&metadata.TableMetadata{Name: "pub", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "label"}}},
		&metadata.TableMetadata{Name: "priv", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "secret"}}},
	))
	require.Len(t, span.Results, 1, "ON (label) restricts the output to exactly the listed column")
	require.Equal(t, "label", span.Results[0].Name)
	require.Equal(t, []string{"priv.secret", "pub.label"}, sourcesOf(span.Results[0]),
		"the single output column carries both arms' label lineage")
}

// TestLeakPin_SchemaQualifiedStar: a schema-qualified dot-star expands the named
// schema's table with schema-qualified lineage (the legacy resolver errored with
// "resource not found: column: analytics").
func TestLeakPin_SchemaQualifiedStar(t *testing.T) {
	span := resolvePin(t, spannerResolveOptions(), "SELECT analytics.events.* FROM analytics.events", "db", []*metadata.SchemaMetadata{
		{Name: ""},
		{Name: "analytics", Tables: []*metadata.TableMetadata{
			{Name: "events", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "payload"}}},
		}},
	})
	require.Len(t, span.Results, 2, "the schema-qualified star expands events' two columns")
	require.Equal(t, []string{"analytics.events.id"}, sourcesOf(span.Results[0]))
	require.Equal(t, []string{"analytics.events.payload"}, sourcesOf(span.Results[1]))
}

// TestLeakPin_MixedUserSystemRejected: mixing user and system tables keeps
// failing closed, and a system-only query early-returns an EMPTY
// SelectInfoSchema span under the spanner dialect (SPANNER_SYS included).
func TestLeakPin_MixedUserSystemRejected(t *testing.T) {
	opts := spannerResolveOptions()
	opts.DefaultDatabase = "db"
	opts.GetDatabaseMetadata = corpusGetter(&metadata.DatabaseSchemaMetadata{
		Name:    "db",
		Schemas: pinTables(&metadata.TableMetadata{Name: "users", Columns: []*metadata.ColumnMetadata{{Name: "id"}}}),
	})
	_, err := ResolveQuerySpan("SELECT * FROM users JOIN INFORMATION_SCHEMA.TABLES ON TRUE", opts)
	require.ErrorIs(t, err, ErrMixUserSystemTables)

	span, err := ResolveQuerySpan("SELECT * FROM INFORMATION_SCHEMA.TABLES", opts)
	require.NoError(t, err)
	require.Equal(t, SelectInfoSchema, span.Type)
	require.Empty(t, span.Results, "system-only query early-returns an empty span (legacy behavior)")
	require.Empty(t, span.SourceColumns)

	span, err = ResolveQuerySpan("SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE", opts)
	require.NoError(t, err)
	require.Equal(t, SelectInfoSchema, span.Type, "SPANNER_SYS is a system schema for the Spanner dialect")
	require.Empty(t, span.Results)
}

// TestQuerySpanSingleStatement: the span covers exactly one statement, and an
// empty text yields an empty span rather than an error.
func TestQuerySpanSingleStatement(t *testing.T) {
	empty, err := ResolveQuerySpan("", bigqueryResolveOptions())
	require.NoError(t, err)
	require.Empty(t, empty.Results)
	require.Empty(t, empty.SourceColumns)

	_, err = ResolveQuerySpan("SELECT 1; SELECT 2", bigqueryResolveOptions())
	require.EqualError(t, err, "expecting only one statement to get query span, but got 2")
}

// TestQuerySpanMissingMetadata: a table missing from the synced metadata keeps
// its table-level resource and reports a non-fatal not-found error, which the
// masking consumer treats as "must not mask on this span".
func TestQuerySpanMissingMetadata(t *testing.T) {
	opts := bigqueryResolveOptions()
	opts.DefaultDatabase = "ds1"
	opts.GetDatabaseMetadata = corpusGetter(&metadata.DatabaseSchemaMetadata{Name: "ds1"})
	span, err := ResolveQuerySpan("SELECT id FROM missing", opts)
	require.NoError(t, err)
	require.Equal(t, Select, span.Type)
	require.NotNil(t, span.NotFoundError)
	require.Equal(t, ColumnSet{{Database: "ds1", Table: "missing"}: true}, span.SourceColumns)
}

// TestQuerySpanMetadataErrorIsFatal: a metadata lookup failure that is not a
// not-found error aborts the span instead of leaving result lineage silently
// empty (the masker fails open, so an empty lineage would leak).
func TestQuerySpanMetadataErrorIsFatal(t *testing.T) {
	opts := bigqueryResolveOptions()
	opts.DefaultDatabase = "ds1"
	opts.GetDatabaseMetadata = func(string) (*metadata.DatabaseSchemaMetadata, error) {
		return nil, errors.New("store unavailable")
	}
	_, err := ResolveQuerySpan("SELECT id FROM t", opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database metadata for database: ds1")
}

// TestPredicateColumns_NarrowedByDatabase: a predicate reference resolved from a
// table written without a database names the DEFAULT database's column only, so
// another database's table of the same name and column is not reported. It cannot
// live in the recorded corpora, which serve one database per case.
func TestPredicateColumns_NarrowedByDatabase(t *testing.T) {
	column := func(name string) *metadata.ColumnMetadata { return &metadata.ColumnMetadata{Name: name} }
	databases := map[string]*metadata.DatabaseSchemaMetadata{
		"ds1": {Name: "ds1", Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{
			{Name: "t", Columns: []*metadata.ColumnMetadata{column("id"), column("secret")}},
		}}}},
		"other": {Name: "other", Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{
			{Name: "t", Columns: []*metadata.ColumnMetadata{column("id"), column("secret")}},
		}}}},
	}
	opts := bigqueryResolveOptions()
	opts.DefaultDatabase = "ds1"
	opts.GetDatabaseMetadata = func(name string) (*metadata.DatabaseSchemaMetadata, error) {
		return databases[name], nil
	}
	span, err := ResolveQuerySpan("SELECT 1 FROM t CROSS JOIN other.t u WHERE t.secret = 1", opts)
	require.NoError(t, err)
	require.Equal(t, []string{"ds1.t.secret"}, predicateSourcesOf(span))
}

// predicateSourcesOf renders the predicate set as sorted "database.table.column"
// names.
func predicateSourcesOf(span *ResolvedQuerySpan) []string {
	out := make([]string, 0, len(span.PredicateColumns))
	for c := range span.PredicateColumns {
		out = append(out, c.Database+"."+c.Table+"."+c.Column)
	}
	slices.Sort(out)
	return out
}

// TestPredicateColumnsFailClosed pins the contract the predicate columns are
// reported under: a statement whose row conditions are ALL resolved reports the
// columns they read, and a statement with any row condition the walk cannot
// resolve reports NO predicate columns, with the reason kept for the caller's
// tests. An empty set is honest; a set that silently misses a read is not,
// because the masking consumer acts on what it contains.
func TestPredicateColumnsFailClosed(t *testing.T) {
	column := func(name string) *metadata.ColumnMetadata { return &metadata.ColumnMetadata{Name: name} }
	table := func(name string, columns ...string) *metadata.TableMetadata {
		out := &metadata.TableMetadata{Name: name}
		for _, c := range columns {
			out.Columns = append(out.Columns, column(c))
		}
		return out
	}
	schemas := []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{
		table("t", "id", "secret"),
		table("u", "id", "secret"),
		table("v", "secret"),
	}}}
	opts := bigqueryResolveOptions()

	for _, tc := range []struct {
		name   string
		sql    string
		reason string
	}{
		{
			// The derived relation's `k` is a select-list scalar subquery, which the
			// projection layer does not carry lineage for.
			name:   "a derived relation projecting a scalar subquery",
			sql:    "SELECT 1 FROM (SELECT (SELECT secret FROM t) AS k) c WHERE c.k > 5",
			reason: "resolves to no lineage",
		},
		{
			// A BY NAME merge materializes its star arm, so the position arrives with no
			// lineage.
			name:   "a BY NAME merge over a star arm",
			sql:    "SELECT 1 FROM t WHERE 5 IN (SELECT * FROM t UNION ALL BY NAME ON (id) (SELECT * FROM u UNION ALL SELECT 0 AS id, secret FROM v))",
			reason: "",
		},
		{
			name:   "a condition on a column the metadata does not hold",
			sql:    "SELECT 1 FROM t WHERE missing = 1",
			reason: "is not in the synced metadata",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			span := resolvePin(t, opts, tc.sql, "ds1", schemas)
			if tc.reason == "" {
				// This shape resolves conservatively today; what matters is that a set is
				// only reported together with no unresolved reason.
				if span.predicateColumnsUnresolved != "" && len(span.PredicateColumns) > 0 {
					t.Errorf("reported %d predicate columns while unresolved: %s", len(span.PredicateColumns), span.predicateColumnsUnresolved)
				}
				return
			}
			require.Empty(t, span.PredicateColumns, "a statement with an unresolved row condition reports no predicate columns")
			require.Contains(t, span.predicateColumnsUnresolved, tc.reason)
		})
	}

	t.Run("resolved", func(t *testing.T) {
		span := resolvePin(t, opts, "WITH c AS (SELECT secret AS k FROM t) SELECT 1 FROM c WHERE c.k = 1", "ds1", schemas)
		require.Equal(t, []string{"ds1.t.secret"}, predicateSourcesOf(span))
		require.Empty(t, span.predicateColumnsUnresolved)
	})

	t.Run("no row condition", func(t *testing.T) {
		span := resolvePin(t, opts, "SELECT id FROM t", "ds1", schemas)
		require.Empty(t, span.PredicateColumns)
		require.Empty(t, span.predicateColumnsUnresolved, "a statement without a row condition is not unresolved")
	})
}
