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
	for name, opts := range map[string]ResolveOptions{"bigquery": bigqueryResolveOptions()} {
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

// TestNameResults: results are named as BigQuery names its fields
// (testdata/span.tsv), and a wide select list is named in linear time.
func TestNameResults(t *testing.T) {
	named := func(names ...string) []string {
		var in []ResolvedQuerySpanResult
		for _, n := range names {
			in = append(in, ResolvedQuerySpanResult{Name: n})
		}
		var out []string
		for _, r := range nameResults(in) {
			out = append(out, r.Name)
		}
		return out
	}
	// Each case is the field names a BigQuery dry run states.
	require.Equal(t, []string{"f1_", "f0_", "f2_"}, named("", "f0_", ""))
	require.Equal(t, []string{"x", "X_1"}, named("x", "X"))
	require.Equal(t, []string{"x", "x_2", "x_1"}, named("x", "x", "x_1"))
	require.Equal(t, []string{"x_1", "x", "x_2"}, named("x_1", "x", "x"))
	require.Equal(t, []string{"f1_", "F0_"}, named("", "F0_"))

	wide := make([]string, 100000)
	for i := range wide {
		if i%2 == 1 {
			wide[i] = "x"
		}
	}
	out := named(wide...)
	require.Equal(t, "f49999_", out[len(out)-2])
	require.Equal(t, "x_49999", out[len(out)-1])
}
