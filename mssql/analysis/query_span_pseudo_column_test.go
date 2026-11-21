package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuerySpanPseudoColumns(t *testing.T) {
	identSrc := ColumnResource{Database: "db", Schema: "dbo", Table: "ident_t", Column: "id"}

	t.Run("identity_bare", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT $IDENTITY FROM ident_t")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.Equal(t, "id", span.Results[0].Name)
		require.Equal(t, ColumnSet{identSrc: true}, span.Results[0].SourceColumns)
	})

	t.Run("identitycol_keyword", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT IDENTITYCOL FROM ident_t")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.Equal(t, ColumnSet{identSrc: true}, span.Results[0].SourceColumns)
	})

	t.Run("identity_qualified_alias", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT x.$IDENTITY FROM ident_t AS x")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.Equal(t, ColumnSet{identSrc: true}, span.Results[0].SourceColumns)
	})

	t.Run("identity_qualified_disambiguates", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT ident_t.$IDENTITY FROM ident_t, ident_t2")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.Equal(t, ColumnSet{identSrc: true}, span.Results[0].SourceColumns)
	})

	t.Run("identity_in_predicate", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT payload FROM ident_t WHERE $IDENTITY = 5")
		require.NoError(t, err)
		require.Equal(t, ColumnSet{identSrc: true}, span.PredicateColumns)
	})

	t.Run("identity_ambiguous_errors", func(t *testing.T) {
		_, err := newExtractorSession("db").resolve("SELECT $IDENTITY FROM ident_t, ident_t2")
		require.ErrorContains(t, err, "ambiguous")
	})

	t.Run("identity_no_identity_column_errors", func(t *testing.T) {
		_, err := newExtractorSession("db").resolve("SELECT $IDENTITY FROM t")
		require.ErrorContains(t, err, "no identity column")
	})

	t.Run("graph_table_level_lineage", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT $node_id FROM t")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.Equal(t, ColumnSet{
			{Database: "db", Schema: "dbo", Table: "t", Column: "a"}: true,
			{Database: "db", Schema: "dbo", Table: "t", Column: "b"}: true,
			{Database: "db", Schema: "dbo", Table: "t", Column: "c"}: true,
		}, span.Results[0].SourceColumns)
	})

	t.Run("graph_qualified_two_tables", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT t2.$from_id FROM t, t2")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.Equal(t, ColumnSet{
			{Database: "db", Schema: "dbo", Table: "t2", Column: "a"}: true,
			{Database: "db", Schema: "dbo", Table: "t2", Column: "b"}: true,
		}, span.Results[0].SourceColumns)
	})

	t.Run("graph_bare_ambiguous_errors", func(t *testing.T) {
		_, err := newExtractorSession("db").resolve("SELECT $node_id FROM t, t2")
		require.ErrorContains(t, err, "ambiguous")
	})

	t.Run("graph_empty_table_fails_closed", func(t *testing.T) {
		_, err := newExtractorSession("db").resolve("SELECT $edge_id FROM bare_edge")
		require.ErrorContains(t, err, "no columns in the catalog")
	})

	t.Run("delimited_identifiers_are_real_columns", func(t *testing.T) {
		for _, sql := range []string{
			"SELECT [$node_id] FROM weird",
			"SELECT weird.[$node_id] FROM weird",
		} {
			span, err := newExtractorSession("db").resolve(sql)
			require.NoError(t, err, sql)
			require.Len(t, span.Results, 1)
			require.Equal(t, ColumnSet{
				{Database: "db", Schema: "dbo", Table: "weird", Column: "$node_id"}: true,
			}, span.Results[0].SourceColumns, sql)
		}
		span, err := newExtractorSession("db").resolve("SELECT [IDENTITYCOL] FROM weird")
		require.NoError(t, err)
		require.Equal(t, ColumnSet{
			{Database: "db", Schema: "dbo", Table: "weird", Column: "IDENTITYCOL"}: true,
		}, span.Results[0].SourceColumns)
	})

	t.Run("delimited_table_bare_column_stays_pseudo", func(t *testing.T) {
		span, err := newExtractorSession("db").resolve("SELECT [weird].$node_id FROM weird")
		require.NoError(t, err)
		require.Len(t, span.Results, 1)
		require.Equal(t, ColumnSet{
			{Database: "db", Schema: "dbo", Table: "weird", Column: "$node_id"}:    true,
			{Database: "db", Schema: "dbo", Table: "weird", Column: "IDENTITYCOL"}: true,
		}, span.Results[0].SourceColumns)
	})

	t.Run("rowguid_stays_fail_closed", func(t *testing.T) {
		_, err := newExtractorSession("db").resolve("SELECT $ROWGUID FROM t")
		require.Error(t, err)
	})
}
