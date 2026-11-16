package analysis

import (
	"errors"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
	"github.com/stretchr/testify/require"
)

func fixSpanMetadata() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Name: "catalog1",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{
				{Name: "a", Columns: []*metadatapb.ColumnMetadata{{Name: "k"}, {Name: "p"}}},
				{Name: "b", Columns: []*metadatapb.ColumnMetadata{{Name: "k"}, {Name: "q"}}},
				{Name: "users", Columns: []*metadatapb.ColumnMetadata{{Name: "id"}, {Name: "email"}}},
			},
		}},
	}
}

func fixSpanResolver(t *testing.T) *resolver {
	t.Helper()
	getter, lister := buildMockDatabaseMetadataGetter([]*metadatapb.DatabaseSchemaMetadata{fixSpanMetadata()})
	return newResolver("catalog1", "public", ResolveOptions{
		GetDatabaseMetadata: getter,
		ListDatabases:       lister,
	})
}

// TestQuerySpan_UsingJoinStarUsesCoalescedProjection covers defect 1: SELECT *
// over a USING join lists each join column once at the front, not the flat
// concatenation of both relations' columns.
func TestQuerySpan_UsingJoinStarUsesCoalescedProjection(t *testing.T) {
	span, err := fixSpanResolver(t).getQuerySpan("SELECT * FROM a JOIN b USING (k)")
	require.NoError(t, err)
	require.Len(t, span.Results, 3)
	require.Equal(t, []string{"k", "p", "q"}, []string{span.Results[0].Name, span.Results[1].Name, span.Results[2].Name})
	require.Equal(t, ColumnSet{
		{Database: "catalog1", Schema: "public", Table: "a", Column: "k"}: true,
		{Database: "catalog1", Schema: "public", Table: "b", Column: "k"}: true,
	}, span.Results[0].SourceColumns)
}

// TestQuerySpan_ColumnAliasedStarIsNotWildcard covers defect 2: a column
// aliased to "*" is one ordinary output column, not a star to expand.
func TestQuerySpan_ColumnAliasedStarIsNotWildcard(t *testing.T) {
	span, err := fixSpanResolver(t).getQuerySpan(`SELECT email AS "*" FROM users`)
	require.NoError(t, err)
	require.Len(t, span.Results, 1)
	require.Equal(t, "*", span.Results[0].Name)
	require.Equal(t, ColumnSet{
		{Database: "catalog1", Schema: "public", Table: "users", Column: "email"}: true,
	}, span.Results[0].SourceColumns)
}

// TestQuerySpan_ParseErrorSurfaces covers defect 3: a malformed statement
// reports the parser's ParseError instead of only its text.
func TestQuerySpan_ParseErrorSurfaces(t *testing.T) {
	_, err := fixSpanResolver(t).getQuerySpan("SELECT id FROM users )))")
	require.Error(t, err)
	var parseErr *parser.ParseError
	require.ErrorAs(t, err, &parseErr)
	require.Positive(t, parseErr.Position)
}

// TestQuerySpan_MetadataFailureSurfaces covers defect 4: a metadata lookup
// failure that is not a missing resource fails the span instead of falling back
// to a table-level source.
func TestQuerySpan_MetadataFailureSurfaces(t *testing.T) {
	getter := func(string) (*metadatapb.DatabaseSchemaMetadata, error) {
		return nil, errors.New("metadata backend unavailable")
	}
	lister := func() ([]string, error) { return []string{"catalog1"}, nil }
	extractor := newResolver("catalog1", "public", ResolveOptions{
		GetDatabaseMetadata: getter,
		ListDatabases:       lister,
	})
	_, err := extractor.getQuerySpan("SELECT id FROM users")
	require.ErrorContains(t, err, "metadata backend unavailable")
}

// TestQuerySpan_CreateAsSelectReadsDefiningQuery covers defect 5: CREATE VIEW /
// CREATE TABLE … AS SELECT stays DDL but reports the tables its query reads.
func TestQuerySpan_CreateAsSelectReadsDefiningQuery(t *testing.T) {
	for _, statement := range []string{
		"CREATE VIEW v AS SELECT id FROM users",
		"CREATE TABLE copied AS SELECT id FROM users",
		"CREATE MATERIALIZED VIEW mv AS SELECT id FROM users",
	} {
		t.Run(statement, func(t *testing.T) {
			span, err := fixSpanResolver(t).getQuerySpan(statement)
			require.NoError(t, err)
			require.Equal(t, DDL, span.Type)
			require.Empty(t, span.Results)
			require.Equal(t, ColumnSet{
				{Database: "catalog1", Schema: "public", Table: "users"}: true,
			}, span.SourceColumns)
		})
	}
}

// TestQuerySpan_ParsesStatementOnce covers defect 10: classification and
// analysis share one parse of the statement. The query-bearing DDL unwrap also
// reuses the parsed AST.
func TestQuerySpan_ParsesStatementOnce(t *testing.T) {
	old := parseTrinoSpan
	t.Cleanup(func() { parseTrinoSpan = old })
	for _, statement := range []string{
		"SELECT id FROM users",
		"EXPLAIN SELECT id FROM users",
		"CREATE VIEW v AS SELECT id FROM users",
	} {
		calls := 0
		parseTrinoSpan = func(s string) (*ast.File, error) {
			calls++
			return old(s)
		}
		_, err := fixSpanResolver(t).getQuerySpan(statement)
		require.NoError(t, err, statement)
		require.Equal(t, 1, calls, "%s must parse once", statement)
	}
}
