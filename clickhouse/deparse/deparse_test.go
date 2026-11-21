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
						Comment:     "raw events",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "UUID", Default: "generateUUIDv4()", Comment: "event id"},
							{Name: "tenant_id", Type: "UInt64"},
							{Name: "payload", Type: "Nullable(String)", Nullable: true, Default: "NULL"},
							{Name: "status", Type: "String", Default: "pending"},
							{Name: "created_at", Type: "DateTime", Default: "now()"},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{Name: "secondary", Primary: false, Expressions: []string{"status"}},
							{Name: "primary", Primary: true, Expressions: []string{"tenant_id"}},
						},
					},
					{
						Name:   "scratch",
						Engine: "MergeTree",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "k", Type: "Int32"},
						},
					},
					{
						Name:   "mem_cache",
						Engine: "Memory",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "k", Type: "Int32"},
							{Name: "v", Type: "String", Nullable: true},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:       "active_events",
						Definition: "  SELECT id, tenant_id FROM events WHERE status = 'active' ; \n",
						Comment:    "only active events",
					},
					{
						Name:       "plain_view",
						Definition: "SELECT k FROM scratch",
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
CREATE TABLE events (
  id UUID NOT NULL DEFAULT generateUUIDv4() COMMENT 'event id',
  tenant_id UInt64 NOT NULL,
  payload Nullable(String) DEFAULT NULL,
  status String NOT NULL DEFAULT 'pending',
  created_at DateTime NOT NULL DEFAULT now()
)
ENGINE = MergeTree
ORDER BY (tenant_id, created_at)
PRIMARY KEY (tenant_id)
COMMENT 'raw events';

--
-- Table structure for ` + "`scratch`" + `
--
CREATE TABLE scratch (
  k Int32 NOT NULL
)
ENGINE = MergeTree
ORDER BY tuple();

--
-- Table structure for ` + "`mem_cache`" + `
--
CREATE TABLE mem_cache (
  k Int32 NOT NULL,
  v String
)
ENGINE = Memory;

--
-- View structure for ` + "`active_events`" + `
--
CREATE OR REPLACE VIEW active_events AS (SELECT id, tenant_id FROM events WHERE status = 'active') COMMENT 'only active events';

--
-- View structure for ` + "`plain_view`" + `
--
CREATE OR REPLACE VIEW plain_view AS (SELECT k FROM scratch);
`
	require.Equal(t, want, got)
}

func TestGetTableAndViewDefinition(t *testing.T) {
	tbl := &metadatapb.TableMetadata{
		Name:   "metrics",
		Engine: "MergeTree",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "ts", Type: "DateTime", Default: "CURRENT_TIMESTAMP"},
			{Name: "day", Type: "Date", Default: "today"},
			{Name: "tag", Type: "String", Default: "default_tag"},
		},
	}

	gotTable, err := GetTableDefinition("ignored_schema", tbl)
	require.NoError(t, err)
	wantTable := `
--
-- Table structure for ` + "`metrics`" + `
--
CREATE TABLE metrics (
  ts DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  day Date NOT NULL DEFAULT today,
  tag String NOT NULL DEFAULT 'default_tag'
)
ENGINE = MergeTree
ORDER BY tuple();
`
	require.Equal(t, wantTable, gotTable)

	view := &metadatapb.ViewMetadata{
		Name:       "recent_metrics",
		Definition: "\n  SELECT ts, tag FROM metrics;  ",
		Comment:    "recent",
	}

	gotView, err := GetViewDefinition("ignored_schema", view)
	require.NoError(t, err)
	wantView := `
--
-- View structure for ` + "`recent_metrics`" + `
--
CREATE OR REPLACE VIEW recent_metrics AS (SELECT ts, tag FROM metrics) COMMENT 'recent';
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
CREATE TABLE z_tbl (
  a Int32 NOT NULL
);

--
-- Table structure for ` + "`a_tbl`" + `
--
CREATE TABLE a_tbl (
  b Int32 NOT NULL
);

--
-- Table structure for ` + "`m_tbl`" + `
--
CREATE TABLE m_tbl (
  c Int32 NOT NULL
);

--
-- View structure for ` + "`v1`" + `
--
CREATE OR REPLACE VIEW v1 AS (SELECT 1);

--
-- View structure for ` + "`v2`" + `
--
CREATE OR REPLACE VIEW v2 AS (SELECT 2);
`
	require.Equal(t, want, got)
}
