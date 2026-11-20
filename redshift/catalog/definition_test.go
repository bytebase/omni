package catalog

import (
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func TestGetObjectDefinition(t *testing.T) {
	tests := []struct {
		name string
		get  func() (string, error)
		want string
	}{
		{
			name: "table with comments",
			get: func() (string, error) {
				return GetTableDefinition("public", &metadatapb.TableMetadata{
					Name:    "orders",
					Comment: "one row per order",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "id", Type: "integer", Comment: "surrogate key"},
						{Name: "total", Type: "numeric(10,2)", Nullable: true},
						{Name: "label", Type: "character varying(50)", Nullable: true, Comment: "display name"},
					},
				})
			},
			want: `CREATE TABLE "public"."orders" (
  "id" integer NOT NULL,
  "total" numeric(10,2),
  "label" character varying(50)
);
COMMENT ON TABLE "public"."orders" IS 'one row per order';
COMMENT ON COLUMN "public"."orders"."id" IS 'surrogate key';
COMMENT ON COLUMN "public"."orders"."label" IS 'display name';
`,
		},
		{
			name: "view with comment",
			get: func() (string, error) {
				return GetViewDefinition("public", &metadatapb.ViewMetadata{
					Name:       "recent_orders",
					Definition: "SELECT id, total FROM orders",
					Comment:    "last 30 days",
				})
			},
			want: `CREATE OR REPLACE VIEW "public"."recent_orders" AS SELECT id, total FROM orders;
COMMENT ON VIEW "public"."recent_orders" IS 'last 30 days';
`,
		},
		{
			name: "table outside the default schema keeps its qualifier",
			get: func() (string, error) {
				return GetTableDefinition("analytics", &metadatapb.TableMetadata{
					Name:    "orders",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				})
			},
			want: `CREATE TABLE "analytics"."orders" (
  "id" integer NOT NULL
);
`,
		},
		{
			name: "comments with apostrophes stay valid",
			get: func() (string, error) {
				return GetTableDefinition("public", &metadatapb.TableMetadata{
					Name:    "orders",
					Comment: "owner's orders",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Comment: "don't drop"}},
				})
			},
			want: `CREATE TABLE "public"."orders" (
  "id" integer NOT NULL
);
COMMENT ON TABLE "public"."orders" IS 'owner''s orders';
COMMENT ON COLUMN "public"."orders"."id" IS 'don''t drop';
`,
		},
		{
			name: "column collation and defaults",
			get: func() (string, error) {
				return GetTableDefinition("public", &metadatapb.TableMetadata{
					Name: "items",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "status", Type: "character varying(20)", Default: "'active'"},
						{Name: "title", Type: "character varying(50)", Collation: "case_insensitive", Default: "'it''s ok'", Nullable: true},
						{Name: "code", Type: "character varying(20)", Collation: "CS", Default: "SYSDATE", Nullable: true},
						{Name: "tag", Type: "character varying(20)", CollationSchema: "pg_catalog", Collation: "default", Default: "owner's", Nullable: true},
						{Name: "qty", Type: "integer", Default: "0"},
						{Name: "active", Type: "boolean", Default: "true"},
					},
				})
			},
			want: `CREATE TABLE "public"."items" (
  "status" character varying(20) NOT NULL DEFAULT 'active',
  "title" character varying(50) COLLATE CASE_INSENSITIVE DEFAULT 'it''s ok',
  "code" character varying(20) COLLATE CS DEFAULT SYSDATE,
  "tag" character varying(20) COLLATE "pg_catalog"."default" DEFAULT 'owner''s',
  "qty" integer NOT NULL DEFAULT 0,
  "active" boolean NOT NULL DEFAULT true
);
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.get()
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestColumnCommentOrderIsStable pins the ordering of the trailing COMMENT ON
// COLUMN statements across runs.
func TestColumnCommentOrderIsStable(t *testing.T) {
	table := &metadatapb.TableMetadata{Name: "wide"}
	for _, name := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		table.Columns = append(table.Columns, &metadatapb.ColumnMetadata{
			Name:     name,
			Type:     "integer",
			Nullable: true,
			Comment:  "note " + name,
		})
	}

	first, err := GetTableDefinition("public", table)
	require.NoError(t, err)
	require.Equal(t, 8, strings.Count(first, "COMMENT ON COLUMN"))

	for i := 0; i < 200; i++ {
		got, err := GetTableDefinition("public", table)
		require.NoError(t, err)
		require.Equalf(t, first, got, "column comment order changed on iteration %d", i)
	}
}

// TestDatabaseDefinitionKeepsSections ensures the whole-database output frames
// each object with its banner and matches the single-object entry points.
func TestDatabaseDefinitionKeepsSections(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "orders",
						Comment: "one row per order",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Comment: "surrogate key"}},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{Name: "recent_orders", Definition: "SELECT id FROM orders"},
				},
			},
		},
	}

	got, err := DatabaseDefinition(metadata)
	require.NoError(t, err)

	require.Equal(t, `--
-- Table structure for `+"`orders`"+`
--
CREATE TABLE "public"."orders" (
  "id" integer NOT NULL
);
COMMENT ON TABLE "public"."orders" IS 'one row per order';
COMMENT ON COLUMN "public"."orders"."id" IS 'surrogate key';

--
-- View structure for `+"`recent_orders`"+`
--
CREATE OR REPLACE VIEW "public"."recent_orders" AS SELECT id FROM orders;
`, got)

	table, err := GetTableDefinition("public", metadata.Schemas[0].Tables[0])
	require.NoError(t, err)
	require.Contains(t, got, table)
}

// TestDatabaseDefinitionCreatesSchemasAndOrdersDependencies pins that
// non-public schemas are created first, FK parent tables precede child tables,
// dependent views follow the views they read, and nil entries do not panic.
func TestDatabaseDefinitionCreatesSchemasAndOrdersDependencies(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			nil,
			{
				Name: "analytics",
				Tables: []*metadatapb.TableMetadata{
					nil,
					{
						Name: "a_child",
						Columns: []*metadatapb.ColumnMetadata{
							nil,
							{Name: "id", Type: "integer"},
							{Name: "parent_id", Type: "integer", Nullable: true},
						},
						ForeignKeys: []*metadatapb.ForeignKeyMetadata{
							nil,
							{
								Name:              "a_child_parent_fk",
								Columns:           []string{"parent_id"},
								ReferencedSchema:  "public",
								ReferencedTable:   "z_parent",
								ReferencedColumns: []string{"id"},
							},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					nil,
					{Name: "a_summary", Definition: `SELECT id FROM "analytics"."z_base"`},
					{Name: "z_base", Definition: `SELECT id FROM "analytics"."a_child"`},
				},
			},
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "z_parent",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{Name: "z_parent_pkey", Primary: true, Expressions: []string{"id"}},
						},
					},
				},
			},
		},
	}

	got, err := DatabaseDefinition(metadata)
	require.NoError(t, err)

	require.Equal(t, `CREATE SCHEMA IF NOT EXISTS "analytics";

--
-- Table structure for `+"`z_parent`"+`
--
CREATE TABLE "public"."z_parent" (
  "id" integer NOT NULL,
  CONSTRAINT "z_parent_pkey" PRIMARY KEY (id)
);

--
-- Table structure for `+"`a_child`"+`
--
CREATE TABLE "analytics"."a_child" (
  "id" integer NOT NULL,
  "parent_id" integer,
  CONSTRAINT "a_child_parent_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."z_parent"("id")
);

--
-- View structure for `+"`z_base`"+`
--
CREATE OR REPLACE VIEW "analytics"."z_base" AS SELECT id FROM "analytics"."a_child";

--
-- View structure for `+"`a_summary`"+`
--
CREATE OR REPLACE VIEW "analytics"."a_summary" AS SELECT id FROM "analytics"."z_base";
`, got)
}

// TestForeignKeyKeepsReferencedSchemaAndCompositeColumns pins the qualifier on
// the REFERENCES target and composite multi-column foreign key syntax.
func TestForeignKeyKeepsReferencedSchemaAndCompositeColumns(t *testing.T) {
	got, err := GetTableDefinition("analytics", &metadatapb.TableMetadata{
		Name: "child",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "id", Type: "integer"},
			{Name: "tenant_id", Type: "integer", Nullable: true},
			{Name: "parent_id", Type: "integer", Nullable: true},
		},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{
			{
				Name:              "child_parent_fk",
				Columns:           []string{"tenant_id", "parent_id"},
				ReferencedSchema:  "warehouse",
				ReferencedTable:   "parent",
				ReferencedColumns: []string{"tenant_id", "id"},
			},
		},
	})
	require.NoError(t, err)
	require.Contains(t, got, `CONSTRAINT "child_parent_fk" FOREIGN KEY ("tenant_id", "parent_id") REFERENCES "warehouse"."parent"("tenant_id", "id")`)
	require.NotContains(t, got, `REFERENCES "parent"`)
}

// TestIdentifiersNeedingQuotesSurvive covers identifiers with spaces, mixed
// case, and embedded double quotes.
func TestIdentifiersNeedingQuotesSurvive(t *testing.T) {
	got, err := GetTableDefinition("Sales Data", &metadatapb.TableMetadata{
		Name: "MixedCase",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: `od"d`, Type: "integer", Comment: "quoted name"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, `CREATE TABLE "Sales Data"."MixedCase" (
  "od""d" integer NOT NULL
);
COMMENT ON COLUMN "Sales Data"."MixedCase"."od""d" IS 'quoted name';
`, got)
}

// TestIndexKeysAreVerbatimAndForeignKeyColumnsAreQuoted pins the split between
// pg_get_indexdef expressions (already quoted) and foreign key columns (bare),
// and checks that explicit constraint names are preserved.
func TestIndexKeysAreVerbatimAndForeignKeyColumnsAreQuoted(t *testing.T) {
	got, err := GetTableDefinition("public", &metadatapb.TableMetadata{
		Name: "t",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "MixedCase", Type: "integer"},
			{Name: "parent_id", Type: "integer", Nullable: true},
		},
		Indexes: []*metadatapb.IndexMetadata{
			{Name: "t_pkey", Primary: true, Expressions: []string{`"MixedCase"`}},
		},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{
			{
				Name:              "t_parent_fk",
				Columns:           []string{"parent_id"},
				ReferencedSchema:  "public",
				ReferencedTable:   "parent",
				ReferencedColumns: []string{"id"},
			},
		},
	})
	require.NoError(t, err)

	require.Contains(t, got, `CONSTRAINT "t_pkey" PRIMARY KEY ("MixedCase")`)
	require.NotContains(t, got, `"""MixedCase"""`)
	require.Contains(t, got, `CONSTRAINT "t_parent_fk" FOREIGN KEY ("parent_id") REFERENCES "public"."parent"("id")`)
}

// TestDatabaseDefinitionDefersCyclicForeignKeys pins that mutually referencing
// tables defer the back-edge FK until after both tables exist.
func TestDatabaseDefinitionDefersCyclicForeignKeys(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "alpha",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
							{Name: "beta_id", Type: "integer", Nullable: true},
						},
						ForeignKeys: []*metadatapb.ForeignKeyMetadata{
							{
								Name:              "alpha_beta_fk",
								Columns:           []string{"beta_id"},
								ReferencedSchema:  "public",
								ReferencedTable:   "beta",
								ReferencedColumns: []string{"id"},
							},
						},
					},
					{
						Name: "beta",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
							{Name: "alpha_id", Type: "integer", Nullable: true},
						},
						ForeignKeys: []*metadatapb.ForeignKeyMetadata{
							{
								Name:              "beta_alpha_fk",
								Columns:           []string{"alpha_id"},
								ReferencedSchema:  "public",
								ReferencedTable:   "alpha",
								ReferencedColumns: []string{"id"},
							},
						},
					},
				},
			},
		},
	}

	got, err := DatabaseDefinition(metadata)
	require.NoError(t, err)
	require.Equal(t, `--
-- Table structure for `+"`alpha`"+`
--
CREATE TABLE "public"."alpha" (
  "id" integer NOT NULL,
  "beta_id" integer
);

--
-- Table structure for `+"`beta`"+`
--
CREATE TABLE "public"."beta" (
  "id" integer NOT NULL,
  "alpha_id" integer,
  CONSTRAINT "beta_alpha_fk" FOREIGN KEY ("alpha_id") REFERENCES "public"."alpha"("id")
);

ALTER TABLE "public"."alpha" ADD CONSTRAINT "alpha_beta_fk" FOREIGN KEY ("beta_id") REFERENCES "public"."beta"("id");
`, got)
}

// TestDatabaseDefinitionScopedCTEShadowing pins that a nested CTE with the
// same name as a real view does not suppress an outer reference to that view.
func TestDatabaseDefinitionScopedCTEShadowing(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Views: []*metadatapb.ViewMetadata{
					{
						Name:       "a_consumer",
						Definition: `SELECT sub.x, v.x AS y FROM (WITH z_base AS (SELECT 1 AS x) SELECT x FROM z_base) sub JOIN z_base v ON true`,
					},
					{
						Name:       "z_base",
						Definition: `SELECT 2 AS x`,
					},
				},
			},
		},
	}

	got, err := DatabaseDefinition(metadata)
	require.NoError(t, err)
	zIdx := strings.Index(got, "View structure for `z_base`")
	aIdx := strings.Index(got, "View structure for `a_consumer`")
	require.Greater(t, aIdx, zIdx, "expected z_base before a_consumer:\n%s", got)
}

// TestDatabaseDefinitionLateBindingViewMaterializedViewsAndSkipDump pins
// late-binding view preservation, materialized view emission, identity columns,
// token-based default expression matching, SkipDump, and mixed-case setting.
func TestDatabaseDefinitionLateBindingViewMaterializedViewsAndSkipDump(t *testing.T) {
	skippedTable, err := GetTableDefinition("public", &metadatapb.TableMetadata{Name: "ignored", SkipDump: true})
	require.NoError(t, err)
	require.Empty(t, skippedTable)

	skippedView, err := GetViewDefinition("public", &metadatapb.ViewMetadata{Name: "ignored", Definition: "SELECT 1", SkipDump: true})
	require.NoError(t, err)
	require.Empty(t, skippedView)

	tableWithIdentityAndStrings, err := GetTableDefinition("public", &metadatapb.TableMetadata{
		Name: "events",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "id", Type: "bigint", IsIdentity: true, IdentitySeed: 10, IdentityIncrement: 5},
			{Name: "weather", Type: "character varying(20)", Default: "snow", Nullable: true},
			{Name: "status", Type: "character varying(20)", Default: "current_status", Nullable: true},
			{Name: "note", Type: "character varying(40)", Default: "timestamp note", Nullable: true},
			{Name: "legacy", Type: "character varying(40)", Default: "pending (old)", Nullable: true},
			{Name: "created_at", Type: "timestamp", Default: "now()"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, `CREATE TABLE "public"."events" (
  "id" bigint IDENTITY(10, 5) NOT NULL,
  "weather" character varying(20) DEFAULT 'snow',
  "status" character varying(20) DEFAULT 'current_status',
  "note" character varying(40) DEFAULT 'timestamp note',
  "legacy" character varying(40) DEFAULT 'pending (old)',
  "created_at" timestamp NOT NULL DEFAULT now()
);
`, tableWithIdentityAndStrings)

	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "hidden_schema", SkipDump: true, Tables: []*metadatapb.TableMetadata{{Name: "t"}}},
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{Name: "SkipMe", SkipDump: true},
					{Name: "Events", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}},
				},
				Views: []*metadatapb.ViewMetadata{
					{Name: "skip_view", Definition: "SELECT 1", SkipDump: true},
					{
						Name:       "late_v",
						Definition: `CREATE VIEW "public"."late_v" AS SELECT id FROM "public"."Events" WITH NO SCHEMA BINDING`,
					},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{Name: "skip_mv", Definition: "SELECT 1", SkipDump: true},
					{Name: "a_mv_top", Definition: `SELECT id FROM "public"."z_mv_base"`},
					{Name: "z_mv_base", Definition: `SELECT id FROM "public"."Events"`, Comment: "matview note"},
				},
			},
		},
	}

	got, err := DatabaseDefinition(metadata)
	require.NoError(t, err)
	require.Equal(t, `SET enable_case_sensitive_identifier TO true;

--
-- Table structure for `+"`Events`"+`
--
CREATE TABLE "public"."Events" (
  "id" integer NOT NULL
);

--
-- View structure for `+"`late_v`"+`
--
CREATE VIEW "public"."late_v" AS SELECT id FROM "public"."Events" WITH NO SCHEMA BINDING;

--
-- Materialized view structure for `+"`z_mv_base`"+`
--
CREATE MATERIALIZED VIEW "public"."z_mv_base" AS SELECT id FROM "public"."Events";
COMMENT ON MATERIALIZED VIEW "public"."z_mv_base" IS 'matview note';

--
-- Materialized view structure for `+"`a_mv_top`"+`
--
CREATE MATERIALIZED VIEW "public"."a_mv_top" AS SELECT id FROM "public"."z_mv_base";
`, got)
}

// TestDatabaseDefinitionCrossKindViewsOmittedFKSchemaAndEmptyView pins that
// materialized views precede ordinary views that read them, omitted FK
// ReferencedSchema defaults to the table's schema, and empty view definitions
// return an error.
func TestDatabaseDefinitionCrossKindViewsOmittedFKSchemaAndEmptyView(t *testing.T) {
	fkTable, err := GetTableDefinition("analytics", &metadatapb.TableMetadata{
		Name: "child",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "parent_id", Type: "integer", Nullable: true},
		},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{
			{
				Name:              "child_parent_fk",
				Columns:           []string{"parent_id"},
				ReferencedTable:   "parent",
				ReferencedColumns: []string{"id"},
			},
		},
	})
	require.NoError(t, err)
	require.Contains(t, fkTable, `REFERENCES "analytics"."parent"("id")`)

	_, err = GetViewDefinition("public", &metadatapb.ViewMetadata{Name: "empty_v", Definition: "   "})
	require.Error(t, err)

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", Views: []*metadatapb.ViewMetadata{{Name: "empty_v"}}},
		},
	})
	require.Error(t, err)

	got, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Views: []*metadatapb.ViewMetadata{
					{Name: "v_uses_mv", Definition: `SELECT id FROM "public"."mv_base"`},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{Name: "mv_base", Definition: `SELECT 1 AS id`},
				},
			},
		},
	})
	require.NoError(t, err)
	mvIdx := strings.Index(got, "Materialized view structure for `mv_base`")
	vIdx := strings.Index(got, "View structure for `v_uses_mv`")
	require.Greater(t, vIdx, mvIdx, "expected mv_base before v_uses_mv:\n%s", got)
}
