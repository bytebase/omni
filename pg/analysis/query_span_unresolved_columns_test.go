package analysis

import (
	"context"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func TestUnresolvedColumnsSignalFiresOnDegradedSnapshot(t *testing.T) {
	statements := []string{
		"SELECT * FROM public.t",
		"SELECT email FROM public.t",
		"SELECT t.email FROM public.t",
		"SELECT * FROM public.t WHERE email = 'x'",
		"EXPLAIN ANALYZE SELECT * FROM public.t WHERE email = 'x'",
		"SELECT id, email FROM public.t ORDER BY ssn",
		"WITH c AS (SELECT * FROM public.t) SELECT * FROM c",
		// A non-recursive CTE cannot reference itself; the inner t is the physical table.
		"WITH t AS (SELECT count(*) FILTER (WHERE EXISTS (SELECT 1 FROM t)) AS c FROM (SELECT 1) x) SELECT * FROM t",
		"WITH t AS (SELECT string_agg('x', ',' ORDER BY (SELECT count(*) FROM t)) AS c FROM (SELECT 1) x) SELECT * FROM t",
		`WITH "T" AS (SELECT 1 AS n) SELECT count(*) FILTER (WHERE EXISTS (SELECT 1 FROM t)) FROM "T"`,
		"SELECT * FROM generate_series(1, (SELECT count(*)::int FROM public.t))",
		"SELECT * FROM (VALUES ((SELECT count(*) FROM public.t))) v(x)",
	}
	for _, statement := range statements {
		t.Run(statement, func(t *testing.T) {
			span := spanFor(t, statement, degradedSchema())
			require.NotNil(t, span.UnresolvedColumnsError,
				"a table with no synced columns must mark the span as unresolvable")
			require.Contains(t, span.UnresolvedColumnsError.Error(), "public.t")
			require.Equal(t, []string{"db"}, span.UnresolvedColumnsError.Databases(),
				"the re-sync must target the database holding the unresolved relation")
		})
	}
}

func TestUnresolvedColumnsSignalQuietOnHealthySnapshot(t *testing.T) {
	statements := []string{
		"SELECT * FROM public.t",
		"SELECT email FROM public.t",
		"SELECT * FROM public.t NATURAL JOIN public.o",
		"SELECT * FROM public.t JOIN public.o USING (id)",
		"SELECT * FROM public.t LEFT JOIN public.o ON t.id = o.id",
		"SELECT *, md5(email) FROM public.t",
		"SELECT * FROM public.v",
		"SELECT * FROM public.mv",
		"WITH c AS (SELECT * FROM public.t) SELECT * FROM c",
		"SELECT * FROM public.t UNION ALL SELECT * FROM public.t",
		"SELECT DISTINCT ON (id) * FROM public.t",
		"SELECT count(*) FROM public.t",
		"SELECT 1",
		"SELECT * FROM generate_series(1, 3)",
		"EXPLAIN SELECT * FROM public.t",
		"EXPLAIN ANALYZE SELECT * FROM public.t",
		"-- plan\nEXPLAIN ANALYZE SELECT * FROM public.t",
		"SHOW search_path",
		"SET search_path TO public",
	}
	for _, statement := range statements {
		t.Run(statement, func(t *testing.T) {
			span := spanFor(t, statement, healthySchema())
			require.Nil(t, span.UnresolvedColumnsError,
				"a fully synced snapshot must never mark a span unresolvable")
		})
	}
}

func TestUnresolvedColumnsSignalScope(t *testing.T) {
	t.Run("partial column loss is out of scope", func(t *testing.T) {
		span := spanFor(t, "SELECT * FROM public.t", partialSchema())
		require.Nil(t, span.UnresolvedColumnsError)
	})

	t.Run("a view with no synced columns is not refused", func(t *testing.T) {
		span := spanFor(t, "SELECT * FROM public.v", degradedViewSchema())
		require.Nil(t, span.UnresolvedColumnsError)
	})

	t.Run("a foreign table with no synced columns is not refused", func(t *testing.T) {
		metadata := &metadatapb.DatabaseSchemaMetadata{
			Name: "db",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:           "public",
				ExternalTables: []*metadatapb.ExternalTableMetadata{{Name: "ft"}},
			}},
		}
		span := spanFor(t, "SELECT * FROM public.ft", metadata)
		require.Nil(t, span.UnresolvedColumnsError)
	})

	t.Run("a CTE named like a degraded relation reads the CTE", func(t *testing.T) {
		for _, statement := range []string{
			"WITH t AS (SELECT 1 AS n) SELECT * FROM t",
			"WITH outer_q AS (WITH t AS (SELECT 1 AS n) SELECT * FROM t) SELECT * FROM outer_q",
			"SELECT * FROM (WITH t AS (SELECT 1 AS n) SELECT * FROM t) q",
			"WITH RECURSIVE t AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM t WHERE n < 3) SELECT * FROM t",
		} {
			span := spanFor(t, statement, degradedSchema())
			require.Nil(t, span.UnresolvedColumnsError, statement)
		}
	})

	t.Run("a CTE out of scope leaves the relation of its name read", func(t *testing.T) {
		for _, statement := range []string{
			"SELECT * FROM t, (WITH t AS (SELECT 1 AS n) SELECT * FROM t) q",
			"WITH a AS (SELECT * FROM t), t AS (SELECT 1 AS n) SELECT * FROM a",
		} {
			span := spanFor(t, statement, degradedSchema())
			require.NotNil(t, span.UnresolvedColumnsError, statement)
		}
	})

	t.Run("a non-recursive CTE reading its own name reads the table", func(t *testing.T) {
		span := spanFor(t, "WITH t AS (SELECT * FROM t) SELECT * FROM t", degradedSchema())
		require.NotNil(t, span.UnresolvedColumnsError)
		require.Contains(t, span.UnresolvedColumnsError.Error(), "public.t")
	})

	t.Run("a qualified read is still checked when a CTE shares the name", func(t *testing.T) {
		span := spanFor(t, "WITH t AS (SELECT 1 AS n) SELECT * FROM public.t", degradedSchema())
		require.NotNil(t, span.UnresolvedColumnsError)
		require.Contains(t, span.UnresolvedColumnsError.Error(), "public.t")
	})

	t.Run("a statement reading both the CTE and the qualified table is checked", func(t *testing.T) {
		span := spanFor(t, "WITH t AS (SELECT 1 AS n) SELECT * FROM t UNION ALL SELECT * FROM public.t", degradedSchema())
		require.NotNil(t, span.UnresolvedColumnsError)
	})

	t.Run("materialized view without columns is not a degraded table", func(t *testing.T) {
		span := spanFor(t, "SELECT * FROM public.mv", healthySchema())
		require.Nil(t, span.UnresolvedColumnsError)
	})

	t.Run("error names every unresolved relation", func(t *testing.T) {
		metadata := degradedSchema()
		metadata.Schemas[0].Tables = append(metadata.Schemas[0].Tables, &metadatapb.TableMetadata{Name: "o"})
		span := spanFor(t, "SELECT * FROM public.t, public.o", metadata)
		require.NotNil(t, span.UnresolvedColumnsError)
		require.Contains(t, span.UnresolvedColumnsError.Error(), "public.o")
		require.Contains(t, span.UnresolvedColumnsError.Error(), "public.t")
	})
}

func TestUnresolvedColumnsSignalResolvesRelationsNotRoutines(t *testing.T) {
	// Routines do not shadow relations in PostgreSQL search_path resolution.
	shadowed := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "a", Functions: []*metadatapb.FunctionMetadata{{Name: "t", Signature: "t()"}}},
			{Name: "b", Tables: []*metadatapb.TableMetadata{{Name: "t"}}},
		},
	}
	span, err := ResolveQuerySpan(context.Background(), "SELECT * FROM t", ResolveOptions{
		DefaultDatabase: "db",
		SearchPath:      []string{"a", "b"},
		Catalog:         testCatalog(shadowed),
	})
	require.NoError(t, err)
	require.NotNil(t, span.UnresolvedColumnsError,
		"a routine named like the table must not shadow it, or the column-less b.t goes unchecked")
	require.Contains(t, span.UnresolvedColumnsError.Error(), "b.t")

	// Sequences share the relation namespace and do shadow later tables.
	sequenceFirst := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "a", Sequences: []*metadatapb.SequenceMetadata{{Name: "t"}}},
			{Name: "b", Tables: []*metadatapb.TableMetadata{{Name: "t"}}},
		},
	}
	span, err = ResolveQuerySpan(context.Background(), "SELECT * FROM t", ResolveOptions{
		DefaultDatabase: "db",
		SearchPath:      []string{"a", "b"},
		Catalog:         testCatalog(sequenceFirst),
	})
	require.NoError(t, err)
	require.Nil(t, span.UnresolvedColumnsError,
		"the sequence in a is the relation this query reads, and a sequence carries no column list to judge")
}
