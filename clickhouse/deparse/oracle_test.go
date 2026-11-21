package deparse

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
	metadatapb "github.com/bytebase/omni/metadata"
)

func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// TestDatabaseDefinitionExecutesOnClickHouse renders synced-shape ClickHouse
// metadata — including quoted literal defaults, non-nullable columns carrying
// sync's "NULL" sentinel, Nullable columns, special identifiers, escaped
// comments, MergeTree tables with and without ORDER BY, a full CREATE VIEW
// definition, and a bare SELECT view definition — and executes every emitted
// statement on the pinned ClickHouse 26.8.38.9 server.
func TestDatabaseDefinitionExecutesOnClickHouse(t *testing.T) {
	ctx := context.Background()
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	database := testcontainer.Database(t)

	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: database,
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{
				{
					Name:        "events",
					Engine:      "MergeTree",
					SortingKeys: []string{"id", "name"},
					Comment:     "events' log\\v1",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "id", Type: "UInt64", Default: "NULL"},
						{Name: "name", Type: "String", Default: "NULL", Comment: "the 'event'\\name"},
						{Name: "at", Type: "DateTime64(3, 'UTC')", Default: "NULL"},
						{Name: "amount", Type: "Nullable(Decimal(12, 2))", Nullable: true, Default: "NULL"},
						{Name: "tags", Type: "Array(String)", Default: "NULL"},
						{Name: "attrs", Type: "Map(String, String)", Default: "NULL"},
						{Name: "level", Type: "Enum8('info' = 1, 'warn' = 2)", Default: "'info'"},
						{Name: "seen", Type: "Date", Default: "'2026-01-02'"},
						{Name: "retries", Type: "UInt32", Default: "3"},
						{Name: "created_at", Type: "DateTime", Default: "now()"},
					},
					Indexes: []*metadatapb.IndexMetadata{
						{Name: "name_bloom", Type: "bloom_filter", Granularity: 4, Expressions: []string{"name"}},
						{Primary: true, Expressions: []string{"id"}},
					},
				},
				{
					Name:   "scratch",
					Engine: "ReplacingMergeTree",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "k", Type: "Int32", Default: "NULL"},
					},
				},
				{
					Name:    "odd-table name",
					Engine:  "Memory",
					Comment: "line1\nline2",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "select", Type: "Int32", Default: "0"},
						{Name: "col with space", Type: "Nullable(String)", Nullable: true, Default: "NULL"},
					},
				},
			},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "a_derived",
					Definition: "SELECT id FROM `query view`",
				},
				{
					Name:       "synced_view",
					Definition: "CREATE VIEW " + database + ".synced_view (`id` UInt64, `name` String) AS SELECT id, name FROM " + database + ".events",
				},
				{
					Name:       "query view",
					Definition: "SELECT id, seen FROM events WHERE retries > 0",
					Comment:    "filtered 'events'",
				},
			},
		}},
	}

	ddl, err := DatabaseDefinition(meta)
	require.NoError(t, err)

	conn, err := db.Conn(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	// Execute each generated DDL statement (one per table/view banner block).
	for _, block := range splitDeparsedBlocks(ddl) {
		_, err := conn.ExecContext(ctx, block)
		require.NoError(t, err, "statement failed:\n%s", block)
	}

	// Verify ClickHouse stored all 3 tables and 3 views with expected comments and defaults.
	var count uint64
	require.NoError(t, conn.QueryRowContext(ctx, "SELECT count() FROM system.tables WHERE database = $1", database).Scan(&count))
	require.Equal(t, uint64(6), count)

	var seenDefault, eventComment string
	require.NoError(t, conn.QueryRowContext(ctx,
		"SELECT column_default FROM information_schema.columns WHERE table_schema = $1 AND table_name = 'events' AND column_name = 'seen'",
		database,
	).Scan(&seenDefault))
	require.Equal(t, "'2026-01-02'", seenDefault)

	require.NoError(t, conn.QueryRowContext(ctx,
		"SELECT comment FROM system.tables WHERE database = $1 AND name = 'events'",
		database,
	).Scan(&eventComment))
	require.Equal(t, "events' log\\v1", eventComment)

	var viewComment string
	require.NoError(t, conn.QueryRowContext(ctx,
		"SELECT comment FROM system.tables WHERE database = $1 AND name = 'query view'",
		database,
	).Scan(&viewComment))
	require.Equal(t, "filtered 'events'", viewComment)
}

func splitDeparsedBlocks(ddl string) []string {
	parts := strings.Split(ddl, ";\n")
	var stmts []string
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			stmts = append(stmts, s)
		}
	}
	return stmts
}
