package deparse

import (
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func TestDatabaseDefinitionNilAndEmpty(t *testing.T) {
	got, err := DatabaseDefinition(nil)
	require.NoError(t, err)
	require.Equal(t, "", got)

	got, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{Name: "analytics"})
	require.NoError(t, err)
	require.Equal(t, "", got)

	got, err = GetTableDefinition("", nil)
	require.NoError(t, err)
	require.Equal(t, "", got)

	got, err = GetViewDefinition("", nil)
	require.NoError(t, err)
	require.Equal(t, "", got)
}

func TestDatabaseDefinitionTablesAndViews(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "analytics",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:        "events",
						Engine:      "MergeTree",
						SortingKeys: []string{"tenant_id", "created_at"},
						Comment:     "raw 'events'\\log",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "UUID", Default: "generateUUIDv4()", Comment: "event's id"},
							{Name: "tenant_id", Type: "UInt64", Default: "NULL"},
							{Name: "payload", Type: "Nullable(String)", Nullable: true, Default: "NULL"},
							{Name: "status", Type: "String", Default: "'pending'"},
							{Name: "seen", Type: "Date", Default: "'2026-01-02'"},
							{Name: "created_at", Type: "DateTime", Default: "now()"},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{Name: "status_bloom", Primary: false, Type: "bloom_filter", Granularity: 4, Expressions: []string{"status"}},
							{Name: "primary", Primary: true, Expressions: []string{"tenant_id"}},
						},
					},
					{
						Name:   "scratch",
						Engine: "ReplacingMergeTree",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "k", Type: "Int32"},
						},
					},
					{
						Name:   "odd-name `with` space",
						Engine: "Memory",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "select", Type: "Int32"},
							{Name: "v", Type: "Nullable(String)", Nullable: true},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:       "active_events",
						Definition: "  SELECT id, tenant_id FROM events WHERE status = 'active' ; \n",
						Comment:    "only 'active' events",
					},
					{
						Name:       "synced_view",
						Definition: "CREATE VIEW analytics.synced_view (`id` UUID, `tenant_id` UInt64) AS SELECT id, tenant_id FROM analytics.events",
					},
				},
			},
		},
	}

	got, err := DatabaseDefinition(meta)
	require.NoError(t, err)

	want := `--
-- Table structure for ` + "`events`" + `
--
CREATE TABLE ` + "`events`" + ` (
  ` + "`id`" + ` UUID DEFAULT generateUUIDv4() COMMENT 'event\'s id',
  ` + "`tenant_id`" + ` UInt64,
  ` + "`payload`" + ` Nullable(String) DEFAULT NULL,
  ` + "`status`" + ` String DEFAULT 'pending',
  ` + "`seen`" + ` Date DEFAULT '2026-01-02',
  ` + "`created_at`" + ` DateTime DEFAULT now(),
  INDEX ` + "`status_bloom`" + ` status TYPE bloom_filter GRANULARITY 4
)
ENGINE = MergeTree
ORDER BY (tenant_id, created_at)
PRIMARY KEY tenant_id
COMMENT 'raw \'events\'\\log';

--
-- Table structure for ` + "`scratch`" + `
--
CREATE TABLE ` + "`scratch`" + ` (
  ` + "`k`" + ` Int32
)
ENGINE = ReplacingMergeTree
ORDER BY tuple();

--
-- Table structure for ` + "`odd-name \\`with\\` space`" + `
--
CREATE TABLE ` + "`odd-name \\`with\\` space`" + ` (
  ` + "`select`" + ` Int32,
  ` + "`v`" + ` Nullable(String)
)
ENGINE = Memory;

--
-- View structure for ` + "`active_events`" + `
--
CREATE OR REPLACE VIEW ` + "`active_events`" + ` AS SELECT id, tenant_id FROM events WHERE status = 'active' COMMENT 'only \'active\' events';

--
-- View structure for ` + "`synced_view`" + `
--
CREATE VIEW analytics.synced_view (` + "`id`" + ` UUID, ` + "`tenant_id`" + ` UInt64) AS SELECT id, tenant_id FROM analytics.events;
`
	require.Equal(t, want, got)
}

func TestGetTableAndViewDefinition(t *testing.T) {
	tbl := &metadatapb.TableMetadata{
		Name:   "metrics",
		Engine: "MergeTree",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "ts", Type: "DateTime", Default: "now()"},
			{Name: "day", Type: "Date", Default: "today()"},
			{Name: "tag", Type: "String", Default: "'default_tag'"},
		},
	}

	gotTable, err := GetTableDefinition("ignored_schema", tbl)
	require.NoError(t, err)
	wantTable := `
--
-- Table structure for ` + "`metrics`" + `
--
CREATE TABLE ` + "`metrics`" + ` (
  ` + "`ts`" + ` DateTime DEFAULT now(),
  ` + "`day`" + ` Date DEFAULT today(),
  ` + "`tag`" + ` String DEFAULT 'default_tag'
)
ENGINE = MergeTree
ORDER BY tuple();
`
	require.Equal(t, wantTable, gotTable)

	view := &metadatapb.ViewMetadata{
		Name:       "recent_metrics",
		Definition: "\n  SELECT ts, tag FROM metrics;  ",
		Comment:    "recent\nmetrics",
	}

	gotView, err := GetViewDefinition("ignored_schema", view)
	require.NoError(t, err)
	wantView := `
--
-- View structure for ` + "`recent_metrics`" + `
--
CREATE OR REPLACE VIEW ` + "`recent_metrics`" + ` AS SELECT ts, tag FROM metrics COMMENT 'recent\nmetrics';
`
	require.Equal(t, wantView, gotView)
}

func TestDatabaseDefinitionPreservesSchemaAndObjectOrder(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "s1",
				Tables: []*metadatapb.TableMetadata{
					{Name: "z_tbl", Columns: []*metadatapb.ColumnMetadata{{Name: "a", Type: "Int32"}}},
					{Name: "a_tbl", Columns: []*metadatapb.ColumnMetadata{{Name: "b", Type: "Int32"}}},
				},
				Views: []*metadatapb.ViewMetadata{
					{Name: "v1", Definition: "SELECT 1"},
				},
			},
			{
				Name: "s2",
				Tables: []*metadatapb.TableMetadata{
					{Name: "m_tbl", Columns: []*metadatapb.ColumnMetadata{{Name: "c", Type: "Int32"}}},
				},
				Views: []*metadatapb.ViewMetadata{
					{Name: "v2", Definition: "SELECT 2"},
				},
			},
		},
	}

	got, err := DatabaseDefinition(meta)
	require.NoError(t, err)

	want := `--
-- Table structure for ` + "`z_tbl`" + `
--
CREATE TABLE ` + "`z_tbl`" + ` (
  ` + "`a`" + ` Int32
);

--
-- Table structure for ` + "`a_tbl`" + `
--
CREATE TABLE ` + "`a_tbl`" + ` (
  ` + "`b`" + ` Int32
);

--
-- Table structure for ` + "`m_tbl`" + `
--
CREATE TABLE ` + "`m_tbl`" + ` (
  ` + "`c`" + ` Int32
);

--
-- View structure for ` + "`v1`" + `
--
CREATE OR REPLACE VIEW ` + "`v1`" + ` AS SELECT 1;

--
-- View structure for ` + "`v2`" + `
--
CREATE OR REPLACE VIEW ` + "`v2`" + ` AS SELECT 2;
`
	require.Equal(t, want, got)
}
