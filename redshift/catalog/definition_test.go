package catalog

import (
	"os"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	pgparser "github.com/bytebase/omni/redshift/parser"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
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
			want: `SET search_path TO "public";
CREATE OR REPLACE VIEW "public"."recent_orders" AS SELECT id, total FROM orders;
RESET search_path;
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
SET search_path TO "public";
CREATE OR REPLACE VIEW "public"."recent_orders" AS SELECT id FROM orders;
RESET search_path;
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
SET search_path TO "analytics", public;
CREATE OR REPLACE VIEW "analytics"."z_base" AS SELECT id FROM "analytics"."a_child";
RESET search_path;

--
-- View structure for `+"`a_summary`"+`
--
SET search_path TO "analytics", public;
CREATE OR REPLACE VIEW "analytics"."a_summary" AS SELECT id FROM "analytics"."z_base";
RESET search_path;
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
	require.Equal(t, `SET enable_case_sensitive_identifier TO true;

CREATE TABLE "Sales Data"."MixedCase" (
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

	require.True(t, strings.HasPrefix(got, "SET enable_case_sensitive_identifier TO true;\n\n"))
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
			{Name: "ext_id", Type: "bigint", IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, IdentitySeed: 1, IdentityIncrement: 1},
			{Name: "weather", Type: "character varying(20)", Default: "snow", Nullable: true},
			{Name: "status", Type: "character varying(20)", Default: "current_status", Nullable: true},
			{Name: "note", Type: "character varying(40)", Default: "timestamp note", Nullable: true},
			{Name: "legacy", Type: "character varying(40)", Default: "pending (old)", Nullable: true},
			{Name: "created_at", Type: "timestamp", Default: "now()"},
			{Name: "max_score", Type: "integer", Default: "GREATEST(1, 2)"},
			{Name: "big_num", Type: "double precision", Default: "2.0e19"},
			{Name: "opt", Type: "integer", Default: "null", Nullable: true},
			{Name: "esc_str", Type: "character varying(20)", Default: `E'line\n'`, Nullable: true},
			{Name: "dollar_str", Type: "character varying(20)", Default: `$$pending$$`, Nullable: true},
			{Name: "win_path", Type: "character varying(64)", Default: `C:\temp\a'b`, Comment: `\'; DROP TABLE victim; --`, Nullable: true},
		},
	})
	require.NoError(t, err)
	require.Equal(t, `CREATE TABLE "public"."events" (
  "id" bigint IDENTITY(10, 5) NOT NULL,
  "ext_id" bigint GENERATED BY DEFAULT AS IDENTITY(1, 1) NOT NULL,
  "weather" character varying(20) DEFAULT 'snow',
  "status" character varying(20) DEFAULT 'current_status',
  "note" character varying(40) DEFAULT 'timestamp note',
  "legacy" character varying(40) DEFAULT 'pending (old)',
  "created_at" timestamp NOT NULL DEFAULT now(),
  "max_score" integer NOT NULL DEFAULT GREATEST(1, 2),
  "big_num" double precision NOT NULL DEFAULT 2.0e19,
  "opt" integer DEFAULT NULL,
  "esc_str" character varying(20) DEFAULT E'line\n',
  "dollar_str" character varying(20) DEFAULT $$pending$$,
  "win_path" character varying(64) DEFAULT 'C:\\temp\\a''b'
);
COMMENT ON COLUMN "public"."events"."win_path" IS '\\''; DROP TABLE victim; --';
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
						Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{Name: "skip_mv", Definition: "SELECT 1", SkipDump: true},
					{Name: "a_mv_top", Definition: `SELECT id FROM "public"."z_mv_base"`},
					{Name: "z_mv_base", Definition: `SELECT id FROM "public"."Events"`},
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
SET search_path TO "public";
CREATE VIEW "public"."late_v" AS SELECT id FROM "public"."Events" WITH NO SCHEMA BINDING;
RESET search_path;

--
-- Materialized view structure for `+"`z_mv_base`"+`
--
SET search_path TO "public";
CREATE MATERIALIZED VIEW "public"."z_mv_base" AS SELECT id FROM "public"."Events";
RESET search_path;

--
-- Materialized view structure for `+"`a_mv_top`"+`
--
SET search_path TO "public";
CREATE MATERIALIZED VIEW "public"."a_mv_top" AS SELECT id FROM "public"."z_mv_base";
RESET search_path;
`, got)

	mixedView, err := GetViewDefinition("public", &metadatapb.ViewMetadata{Name: "MixedView", Definition: "SELECT 1"})
	require.NoError(t, err)
	require.Equal(t, "SET enable_case_sensitive_identifier TO true;\n\nSET search_path TO \"public\";\nCREATE OR REPLACE VIEW \"public\".\"MixedView\" AS SELECT 1;\nRESET search_path;\n", mixedView)

	mixedColView, err := GetViewDefinition("public", &metadatapb.ViewMetadata{
		Name:       "v",
		Columns:    []*metadatapb.ColumnMetadata{{Name: "MixedAlias"}},
		Definition: `SELECT 1 AS "MixedAlias"`,
	})
	require.NoError(t, err)
	require.Equal(t, "SET enable_case_sensitive_identifier TO true;\n\nSET search_path TO \"public\";\nCREATE OR REPLACE VIEW \"public\".\"v\" AS SELECT 1 AS \"MixedAlias\";\nRESET search_path;\n", mixedColView)

	newlineNameDump, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{Name: "safe\nDROP TABLE victim; --", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Contains(t, newlineNameDump, "-- Table structure for `safe DROP TABLE victim; --`\n--\nCREATE TABLE")
}

// TestDatabaseDefinitionRoutinesAndViewColumnComments pins UDF/procedure
// emission order and view column COMMENT ON COLUMN statements.
func TestDatabaseDefinitionRoutinesAndViewColumnComments(t *testing.T) {
	viewDef, err := GetViewDefinition("public", &metadatapb.ViewMetadata{
		Name:       "v",
		Definition: "SELECT 1 AS id",
		Comment:    "view note",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "id", Type: "integer", Comment: "view col note"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, `SET search_path TO "public";
CREATE OR REPLACE VIEW "public"."v" AS SELECT 1 AS id;
RESET search_path;
COMMENT ON VIEW "public"."v" IS 'view note';
COMMENT ON COLUMN "public"."v"."id" IS 'view col note';
`, viewDef)

	got, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "skip_fn", SkipDump: true},
					{
						Name:       "a_wrapper",
						Definition: "CREATE FUNCTION public.a_wrapper(x int) RETURNS int IMMUTABLE AS $$ SELECT public.f_inc(x) * 2 $$ LANGUAGE SQL",
					},
					{
						Name:       "f_inc",
						Definition: "CREATE FUNCTION public.f_inc(x int) RETURNS int IMMUTABLE AS $$ SELECT x + 1 $$ LANGUAGE SQL",
					},
				},
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "t",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Default: "0"}},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:       "v",
						Definition: "SELECT public.f_inc(id) AS next_id FROM public.t",
						Columns:    []*metadatapb.ColumnMetadata{{Name: "next_id", Type: "integer", Comment: "computed"}},
					},
				},
				Procedures: []*metadatapb.ProcedureMetadata{
					{Name: "skip_proc", SkipDump: true},
					{
						Name:       "p_cleanup",
						Definition: "CREATE PROCEDURE public.p_cleanup() AS $$ BEGIN DELETE FROM public.t; END; $$ LANGUAGE plpgsql",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	fnIdx := strings.Index(got, "Function structure for `f_inc`")
	wrapperIdx := strings.Index(got, "Function structure for `a_wrapper`")
	tblIdx := strings.Index(got, "Table structure for `t`")
	viewIdx := strings.Index(got, "View structure for `v`")
	procIdx := strings.Index(got, "Procedure structure for `p_cleanup`")
	require.GreaterOrEqual(t, fnIdx, 0)
	require.Greater(t, wrapperIdx, fnIdx)
	require.Greater(t, tblIdx, wrapperIdx)
	require.Greater(t, viewIdx, tblIdx)
	require.Greater(t, procIdx, viewIdx)
	require.Contains(t, got, `COMMENT ON COLUMN "public"."v"."next_id" IS 'computed';`)
	require.NotContains(t, got, "skip_fn")
	require.NotContains(t, got, "skip_proc")

	// A table with no columns returns an error instead of invalid SQL.
	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{Name: "empty_t"})
	require.Error(t, err)
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", Tables: []*metadatapb.TableMetadata{{Name: "empty_t"}}},
		},
	})
	require.Error(t, err)

	// A caller without DependencyTables that calls a table-dependent UDF moves
	// to the post-table pass after its callee.
	crossPass, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:       "a_caller",
						Definition: "CREATE FUNCTION public.a_caller(x int) RETURNS int AS $$ SELECT public.z_table_fn(x) $$ LANGUAGE SQL",
					},
					{
						Name:             "z_table_fn",
						Definition:       "CREATE FUNCTION public.z_table_fn(x int) RETURNS int AS $$ SELECT x $$ LANGUAGE SQL",
						DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "t"}},
					},
				},
				Tables: []*metadatapb.TableMetadata{
					{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}},
				},
			},
		},
	})
	require.NoError(t, err)
	tPos := strings.Index(crossPass, "Table structure for `t`")
	zPos := strings.Index(crossPass, "Function structure for `z_table_fn`")
	aPos := strings.Index(crossPass, "Function structure for `a_caller`")
	require.Greater(t, zPos, tPos, "expected z_table_fn after table t:\n%s", crossPass)
	require.Greater(t, aPos, zPos, "expected a_caller after z_table_fn:\n%s", crossPass)

	// Overloaded UDFs: f(int) calls wrapper(int) which calls f(int, int).
	// Name-only matching would create a false cycle between f and wrapper.
	overloaded, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:    "analytics",
				Owner:   "app_owner",
				Comment: "analytics schema",
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:       "f",
						Definition: "CREATE FUNCTION analytics.f(a int) RETURNS int IMMUTABLE AS $$ SELECT analytics.wrapper(a) $$ LANGUAGE SQL",
					},
					{
						Name:       "wrapper",
						Definition: "CREATE FUNCTION analytics.wrapper(x int) RETURNS int IMMUTABLE AS $$ SELECT analytics.f(x, 0) $$ LANGUAGE SQL",
					},
					{
						Name:       "f",
						Definition: "CREATE FUNCTION analytics.f(a int, b int) RETURNS int IMMUTABLE AS $$ SELECT a + b $$ LANGUAGE SQL",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Contains(t, overloaded, "ALTER SCHEMA \"analytics\" OWNER TO \"app_owner\";\n")
	require.Contains(t, overloaded, "COMMENT ON SCHEMA \"analytics\" IS 'analytics schema';\n")
	fTwoPos := strings.Index(overloaded, "CREATE FUNCTION analytics.f(a int, b int)")
	wrapPos := strings.Index(overloaded, "CREATE FUNCTION analytics.wrapper(x int)")
	fOnePos := strings.Index(overloaded, "CREATE FUNCTION analytics.f(a int)")
	require.Greater(t, wrapPos, fTwoPos, "expected f(int,int) before wrapper(int):\n%s", overloaded)
	require.Greater(t, fOnePos, wrapPos, "expected wrapper(int) before f(int):\n%s", overloaded)

	// Mismatched FK column counts return an error instead of truncating.
	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name: "bad_fk",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "a", Type: "integer"},
			{Name: "b", Type: "integer"},
		},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{
			{
				Name:              "bad_fk",
				Columns:           []string{"a", "b"},
				ReferencedSchema:  "public",
				ReferencedTable:   "parent",
				ReferencedColumns: []string{"x"},
			},
		},
	})
	require.Error(t, err)
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

func TestDatabaseDefinitionCrossKindOrderingExternalTablesAndValidation(t *testing.T) {
	// 1 & 2. Unified cross-kind ordering: view-dependent functions follow views,
	// table-dependent functions follow tables, and SkipDump external tables are ignored cleanly.
	unified, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:             "v_summary_fn",
						Definition:       "CREATE FUNCTION public.v_summary_fn() RETURNS int AS $$ SELECT count(*)::int FROM public.v_events $$ LANGUAGE SQL",
						DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "v_events"}},
					},
					{
						Name:             "next_code_fn",
						Definition:       "CREATE FUNCTION public.next_code_fn() RETURNS int AS $$ SELECT coalesce(max(id), 0) + 1 FROM public.base_tbl $$ LANGUAGE SQL",
						DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "base_tbl"}},
					},
				},
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "base_tbl",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					},
				},
				ExternalTables: []*metadatapb.ExternalTableMetadata{
					{Name: "skip_ext", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}, SkipDump: true},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:       "v_events",
						Definition: `SELECT id FROM "public"."base_tbl"`,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.NotContains(t, unified, "skip_ext")

	basePos := strings.Index(unified, "Table structure for `base_tbl`")
	nextFnPos := strings.Index(unified, "Function structure for `next_code_fn`")
	viewPos := strings.Index(unified, "View structure for `v_events`")
	vFnPos := strings.Index(unified, "Function structure for `v_summary_fn`")

	require.Greater(t, nextFnPos, basePos, "expected base_tbl before next_code_fn:\n%s", unified)
	require.Greater(t, viewPos, basePos, "expected base_tbl before v_events:\n%s", unified)
	require.Greater(t, vFnPos, viewPos, "expected v_events before v_summary_fn:\n%s", unified)

	// Any table column DEFAULT referencing a UDF (even without a cycle) is rejected.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:       "pure_fn",
						Definition: "CREATE FUNCTION public.pure_fn() RETURNS int IMMUTABLE AS $$ SELECT 42 $$ LANGUAGE SQL",
					},
				},
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "consumer_tbl",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "code", Type: "integer", Default: "public.pure_fn()"},
						},
					},
				},
			},
		},
	})
	require.ErrorContains(t, err, "default_expr")

	// Obsolete Python UDFs (LANGUAGE plpythonu) are rejected explicitly.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:       "py_add",
						Definition: "CREATE FUNCTION public.py_add(a int, b int) RETURNS int IMMUTABLE AS $$ return a + b $$ LANGUAGE plpythonu",
					},
				},
			},
		},
	})
	require.ErrorContains(t, err, "plpythonu")

	// Non-skipped external table returns an explicit error because SchemaMetadata
	// does not capture CREATE EXTERNAL SCHEMA configuration.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "spectrum",
				ExternalTables: []*metadatapb.ExternalTableMetadata{
					{
						Name:    "ext_events",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					},
				},
			},
		},
	})
	require.ErrorContains(t, err, "external schema")

	// 4. Same-arity overload resolution via nested builtin / expression return type:
	// fmt_val(int) -> caller_wrap(text) -> fmt_val(text), where caller_wrap calls fmt_val(length(x)).
	sameArity, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:       "fmt_val",
						Definition: "CREATE FUNCTION public.fmt_val(s text) RETURNS text IMMUTABLE AS $$ SELECT public.caller_wrap(s) $$ LANGUAGE SQL",
					},
					{
						Name:       "caller_wrap",
						Definition: "CREATE FUNCTION public.caller_wrap(s text) RETURNS text IMMUTABLE AS $$ SELECT public.fmt_val(length(s)) $$ LANGUAGE SQL",
					},
					{
						Name:       "fmt_val",
						Definition: "CREATE FUNCTION public.fmt_val(n int) RETURNS text IMMUTABLE AS $$ SELECT n::text $$ LANGUAGE SQL",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	intFnPos := strings.Index(sameArity, "CREATE FUNCTION public.fmt_val(n int)")
	wrapFnPos := strings.Index(sameArity, "CREATE FUNCTION public.caller_wrap(s text)")
	textFnPos := strings.Index(sameArity, "CREATE FUNCTION public.fmt_val(s text)")
	require.Greater(t, wrapFnPos, intFnPos, "expected fmt_val(int) before caller_wrap(text):\n%s", sameArity)
	require.Greater(t, textFnPos, wrapFnPos, "expected caller_wrap(text) before fmt_val(text):\n%s", sameArity)

	// 5. Multiple unnamed foreign keys are both emitted.
	unnamedFKs, err := GetTableDefinition("public", &metadatapb.TableMetadata{
		Name: "order_items",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "order_id", Type: "integer"},
			{Name: "product_id", Type: "integer"},
		},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{
			{
				Name:              "",
				Columns:           []string{"order_id"},
				ReferencedSchema:  "public",
				ReferencedTable:   "orders",
				ReferencedColumns: []string{"id"},
			},
			{
				Name:              "",
				Columns:           []string{"product_id"},
				ReferencedSchema:  "public",
				ReferencedTable:   "products",
				ReferencedColumns: []string{"id"},
			},
		},
	})
	require.NoError(t, err)
	require.Contains(t, unnamedFKs, `FOREIGN KEY ("order_id") REFERENCES "public"."orders"("id")`)
	require.Contains(t, unnamedFKs, `FOREIGN KEY ("product_id") REFERENCES "public"."products"("id")`)

	// 6. Primary/unique key constraints without columns fail explicitly.
	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:    "empty_pk",
		Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
		Indexes: []*metadatapb.IndexMetadata{{Name: "empty_pk_pkey", Primary: true}},
	})
	require.Error(t, err)
	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:    "empty_uk",
		Columns: []*metadatapb.ColumnMetadata{{Name: "code", Type: "text"}},
		Indexes: []*metadatapb.IndexMetadata{{Name: "empty_uk_key", Unique: true}},
	})
	require.Error(t, err)

	// 7. Table owner and key/FK constraint comments:
	// - In GetTableDefinition: OWNER TO follows table, column, and constraint comments.
	// - In DatabaseDefinition: schema/table OWNER TO statements run at the end after objects & comments.
	ownedTbl, err := GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:    "owned_tbl",
		Owner:   "AppOwner",
		Comment: "table comment",
		Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Comment: "id comment"}},
		Indexes: []*metadatapb.IndexMetadata{{Name: "owned_tbl_pkey", Primary: true, Expressions: []string{"id"}, Comment: "pk comment"}},
	})
	require.NoError(t, err)
	require.Equal(t, `SET enable_case_sensitive_identifier TO true;

CREATE TABLE "public"."owned_tbl" (
  "id" integer NOT NULL,
  CONSTRAINT "owned_tbl_pkey" PRIMARY KEY (id)
);
COMMENT ON TABLE "public"."owned_tbl" IS 'table comment';
COMMENT ON COLUMN "public"."owned_tbl"."id" IS 'id comment';
COMMENT ON CONSTRAINT "owned_tbl_pkey" ON "public"."owned_tbl" IS 'pk comment';
ALTER TABLE "public"."owned_tbl" OWNER TO "AppOwner";
`, ownedTbl)

	deferredOwnersAndComments, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:    "reporting",
				Owner:   "reporting_owner",
				Comment: "reporting schema",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "t_child",
						Owner:   "table_owner",
						Comment: "child table",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
						Indexes: []*metadatapb.IndexMetadata{{Name: "t_child_pkey", Primary: true, Expressions: []string{"id"}, Comment: "pk note"}},
						ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
							Name:              "fk_self",
							Columns:           []string{"id"},
							ReferencedSchema:  "reporting",
							ReferencedTable:   "t_child",
							ReferencedColumns: []string{"id"},
							Comment:           "self fk note",
						}},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:       "v_line_comment",
						Definition: "SELECT id FROM reporting.t_child -- end of view",
					},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{
						Name:       "mv_line_comment",
						Definition: "SELECT id FROM reporting.t_child -- end of mv",
					},
				},
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:       "f_line_comment",
						Definition: "CREATE FUNCTION reporting.f_line_comment(a int) RETURNS int IMMUTABLE AS $$ SELECT a -- inside body\n$$ LANGUAGE SQL -- trailing comment",
					},
				},
				Procedures: []*metadatapb.ProcedureMetadata{
					{
						Name:       "p_line_comment",
						Definition: "CREATE PROCEDURE reporting.p_line_comment(x varchar(20)) AS $$ BEGIN NULL; END; $$ LANGUAGE plpgsql -- trailing proc comment",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Contains(t, deferredOwnersAndComments, "COMMENT ON CONSTRAINT \"t_child_pkey\" ON \"reporting\".\"t_child\" IS 'pk note';\n")
	require.Contains(t, deferredOwnersAndComments, "COMMENT ON CONSTRAINT \"fk_self\" ON \"reporting\".\"t_child\" IS 'self fk note';\n")

	publicOwned, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:  "public",
			Owner: "PublicOwner",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, `SET enable_case_sensitive_identifier TO true;
ALTER SCHEMA "public" OWNER TO "PublicOwner";
`, publicOwned)

	// Cross-database view references (otherdb.s.b) must not match local views
	// and introduce false dependency cycles.
	crossDBViews, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Name: "localdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{
				{Name: "a", Definition: "SELECT x FROM otherdb.s.b WITH NO SCHEMA BINDING"},
				{Name: "b", Definition: "SELECT x FROM s.a WITH NO SCHEMA BINDING"},
			},
		}},
	})
	require.NoError(t, err)
	require.Less(t, strings.Index(crossDBViews, `CREATE OR REPLACE VIEW "s"."a"`), strings.Index(crossDBViews, `CREATE OR REPLACE VIEW "s"."b"`))

	// Redshift numeric overload resolution respects implicit cast direction
	// (int8 argument cannot implicitly coerce to int4 overload).
	narrowCastFuncs, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Name: "localdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "f_num",
					Definition: "CREATE FUNCTION s.f_num(x int4) RETURNS int4 IMMUTABLE AS $$ SELECT s.caller_big(CAST(1 AS int8)) $$ LANGUAGE SQL",
				},
				{
					Name:       "caller_big",
					Definition: "CREATE FUNCTION s.caller_big(v int8) RETURNS int8 IMMUTABLE AS $$ SELECT s.f_num(v) $$ LANGUAGE SQL",
				},
				{
					Name:       "f_num",
					Definition: "CREATE FUNCTION s.f_num(x int8) RETURNS int8 IMMUTABLE AS $$ SELECT x $$ LANGUAGE SQL",
				},
			},
		}},
	})
	require.NoError(t, err)
	require.Less(t, strings.Index(narrowCastFuncs, "f_num(x int8)"), strings.Index(narrowCastFuncs, "caller_big("))
	require.Less(t, strings.Index(narrowCastFuncs, "caller_big("), strings.Index(narrowCastFuncs, "f_num(x int4)"))

	// Database-level Owner emits ALTER DATABASE ... OWNER TO ... (and enables
	// case-sensitive identifiers when needed), or rejects an empty database name.
	dbOwned, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Name:  "analytics",
		Owner: "DbAdmin",
	})
	require.NoError(t, err)
	require.Equal(t, `SET enable_case_sensitive_identifier TO true;
ALTER DATABASE "analytics" OWNER TO "DbAdmin";
`, dbOwned)

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Owner: "DbAdmin",
	})
	require.ErrorContains(t, err, "database name is empty")

	// Late-binding views (WITH NO SCHEMA BINDING) must not add reverse relation
	// dependencies onto views that later depend on them.
	lbvOrder, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{
				{Name: "regular_v", Definition: "SELECT x FROM s.late_v"},
				{Name: "late_v", Definition: "CREATE VIEW s.late_v AS SELECT x FROM s.regular_v WITH NO SCHEMA BINDING"},
			},
		}},
	})
	require.NoError(t, err)
	require.Less(t, strings.Index(lbvOrder, `CREATE VIEW s.late_v AS`), strings.Index(lbvOrder, `CREATE OR REPLACE VIEW "s"."regular_v"`))

	// Unrepresented table/column options fail explicitly rather than dropping attributes.
	_, err = GetTableDefinition("s", &metadatapb.TableMetadata{
		Name:          "t",
		CreateOptions: "DISTSTYLE EVEN",
		Columns:       []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
	})
	require.ErrorContains(t, err, "create options")

	_, err = GetTableDefinition("s", &metadatapb.TableMetadata{
		Name:    "t",
		Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Options: []string{"encode=az64"}}},
	})
	require.ErrorContains(t, err, "options")

	// Non-skipped PostgreSQL sequences fail explicitly; SkipDump sequences are ignored.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			Sequences: []*metadatapb.SequenceMetadata{{Name: "ord_seq"}},
		}},
	})
	require.ErrorContains(t, err, "sequence")

	// Non-public schema UDFs set search_path so unqualified sibling calls resolve.
	require.Contains(t, narrowCastFuncs, "SET search_path TO \"s\", public;\nCREATE FUNCTION s.caller_big")
	require.Contains(t, narrowCastFuncs, "RESET search_path;\n")
	require.Contains(t, deferredOwnersAndComments, "SELECT id FROM reporting.t_child -- end of view\n;\n")
	require.Contains(t, deferredOwnersAndComments, "SELECT id FROM reporting.t_child -- end of mv\n;\n")
	require.Contains(t, deferredOwnersAndComments, "$$ LANGUAGE SQL -- trailing comment\n;\n")
	require.Contains(t, deferredOwnersAndComments, "$$ LANGUAGE plpgsql -- trailing proc comment\n;\n")

	createTablePos := strings.Index(deferredOwnersAndComments, `CREATE TABLE "reporting"."t_child"`)
	createViewPos := strings.Index(deferredOwnersAndComments, `CREATE OR REPLACE VIEW "reporting"."v_line_comment"`)
	createProcPos := strings.Index(deferredOwnersAndComments, `CREATE PROCEDURE reporting.p_line_comment(x varchar(20))`)
	tblOwnerPos := strings.Index(deferredOwnersAndComments, `ALTER TABLE "reporting"."t_child" OWNER TO "table_owner";`)
	schemaOwnerPos := strings.Index(deferredOwnersAndComments, `ALTER SCHEMA "reporting" OWNER TO "reporting_owner";`)
	require.Greater(t, tblOwnerPos, createTablePos)
	require.Greater(t, tblOwnerPos, createViewPos)
	require.Greater(t, tblOwnerPos, createProcPos)
	require.Greater(t, schemaOwnerPos, tblOwnerPos)

	// 8. Unsupported Redshift metadata returns explicit errors instead of emitting
	// invalid DDL, and apostrophes in SQL comments do not suppress quoted
	// uppercase identifier detection.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Name:      "shared_db",
		Datashare: true,
	})
	require.ErrorContains(t, err, "datashare")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "f_comment",
				Definition: "CREATE FUNCTION public.f_comment() RETURNS int IMMUTABLE AS $$ SELECT 1 $$ LANGUAGE SQL",
				Comment:    "unsupported",
			}},
		}},
	})
	require.ErrorContains(t, err, "COMMENT ON FUNCTION")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "p_comment",
				Definition: "CREATE PROCEDURE public.p_comment() AS $$ BEGIN NULL; END; $$ LANGUAGE plpgsql",
				Comment:    "unsupported",
			}},
		}},
	})
	require.ErrorContains(t, err, "COMMENT ON PROCEDURE")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_comment",
				Definition: "SELECT 1",
				Comment:    "unsupported",
			}},
		}},
	})
	require.ErrorContains(t, err, "COMMENT ON MATERIALIZED VIEW")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:             "cyc_next",
				Definition:       "CREATE FUNCTION public.cyc_next() RETURNS int AS $$ SELECT coalesce(max(id), 0) + 1 FROM public.cyc_tbl $$ LANGUAGE SQL",
				DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "cyc_tbl"}},
			}},
			Tables: []*metadatapb.TableMetadata{{
				Name:    "cyc_tbl",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Default: "public.cyc_next()"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "default_expr")

	commentQuoteIdent, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_comment_quote",
				Definition: "SELECT\n-- don't treat this comment apostrophe as a string\n1 AS \"MixedAlias\"",
			}},
		}},
	})
	require.NoError(t, err)
	require.Contains(t, commentQuoteIdent, "SET enable_case_sensitive_identifier TO true;\n")

	// 9. CASE expression result types resolve overloaded UDFs without false cycles.
	caseOverload, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "f_case",
					Definition: "CREATE FUNCTION public.f_case(s text) RETURNS text IMMUTABLE AS $$ SELECT public.case_wrap(s) $$ LANGUAGE SQL",
				},
				{
					Name:       "case_wrap",
					Definition: "CREATE FUNCTION public.case_wrap(s text) RETURNS text IMMUTABLE AS $$ SELECT public.f_case(CASE WHEN s IS NOT NULL THEN length(s) ELSE 0 END) $$ LANGUAGE SQL",
				},
				{
					Name:       "f_case",
					Definition: "CREATE FUNCTION public.f_case(n int) RETURNS text IMMUTABLE AS $$ SELECT n::text $$ LANGUAGE SQL",
				},
			},
		}},
	})
	require.NoError(t, err)
	caseIntPos := strings.Index(caseOverload, "CREATE FUNCTION public.f_case(n int)")
	caseWrapPos := strings.Index(caseOverload, "CREATE FUNCTION public.case_wrap(s text)")
	caseTextPos := strings.Index(caseOverload, "CREATE FUNCTION public.f_case(s text)")
	require.Greater(t, caseWrapPos, caseIntPos)
	require.Greater(t, caseTextPos, caseWrapPos)

	// 10. Unsupported FK actions (anything other than empty or NO ACTION) are rejected.
	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:    "fk_cascade",
		Columns: []*metadatapb.ColumnMetadata{{Name: "parent_id", Type: "integer"}},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
			Name:              "fk_parent",
			Columns:           []string{"parent_id"},
			ReferencedSchema:  "public",
			ReferencedTable:   "parent",
			ReferencedColumns: []string{"id"},
			OnDelete:          "CASCADE",
		}},
	})
	require.ErrorContains(t, err, "ON DELETE")

	noActionFK, err := GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:    "fk_no_action",
		Columns: []*metadatapb.ColumnMetadata{{Name: "parent_id", Type: "integer"}},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
			Name:              "fk_parent",
			Columns:           []string{"parent_id"},
			ReferencedSchema:  "public",
			ReferencedTable:   "parent",
			ReferencedColumns: []string{"id"},
			OnDelete:          "NO ACTION",
			OnUpdate:          "NO ACTION",
		}},
	})
	require.NoError(t, err)
	require.Contains(t, noActionFK, `CONSTRAINT "fk_parent" FOREIGN KEY ("parent_id") REFERENCES "public"."parent"("id")`)

	// 11. Check constraints and duplicate column names are rejected explicitly.
	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:             "chk_tbl",
		Columns:          []*metadatapb.ColumnMetadata{{Name: "age", Type: "integer"}},
		CheckConstraints: []*metadatapb.CheckConstraintMetadata{{Name: "chk_age", Expression: "age > 0"}},
	})
	require.ErrorContains(t, err, "CHECK")

	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name: "dup_tbl",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "id", Type: "integer"},
			{Name: "id", Type: "bigint"},
		},
	})
	require.ErrorContains(t, err, "duplicate column")

	// 12. Colliding relation names across tables, views, and materialized views fail explicitly.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{{Name: "events", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}}},
			Views:  []*metadatapb.ViewMetadata{{Name: "events", Definition: "SELECT 1"}},
		}},
	})
	require.ErrorContains(t, err, "duplicate relation")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:              "public",
			Views:             []*metadatapb.ViewMetadata{{Name: "mv_dup", Definition: "SELECT 1"}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "mv_dup", Definition: "SELECT 1"}},
		}},
	})
	require.ErrorContains(t, err, "duplicate relation")

	// Case-distinct relation names ("t" and "T") are preserved when enable_case_sensitive_identifier is enabled.
	caseDistinctRels, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}}},
			Views:  []*metadatapb.ViewMetadata{{Name: "T", Definition: "SELECT 1"}},
		}},
	})
	require.NoError(t, err)
	require.Contains(t, caseDistinctRels, `CREATE TABLE "public"."t"`)
	require.Contains(t, caseDistinctRels, `CREATE OR REPLACE VIEW "public"."T"`)

	// Duplicate schema names fail explicitly.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "app"},
			{Name: "app"},
		},
	})
	require.ErrorContains(t, err, "duplicate schema")

	// Redshift function identity mismatch is rejected explicitly.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "f",
				Definition: "CREATE FUNCTION s.g() RETURNS int IMMUTABLE AS $$ SELECT 1 $$ LANGUAGE SQL",
			}},
		}},
	})
	require.ErrorContains(t, err, "identity")

	// Single-candidate UDF with different arity than a built-in call (e.g. length(text) vs s.length(int, int))
	// must not form a false dependency edge.
	aritySingle, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "length",
					Definition: "CREATE FUNCTION s.length(a int, b int) RETURNS int IMMUTABLE AS $$ SELECT s.txt_caller('hi') $$ LANGUAGE SQL",
				},
				{
					Name:       "txt_caller",
					Definition: "CREATE FUNCTION s.txt_caller(v text) RETURNS int IMMUTABLE AS $$ SELECT length(v) $$ LANGUAGE SQL",
				},
			},
		}},
	})
	require.NoError(t, err)
	require.Less(t, strings.Index(aritySingle, "CREATE FUNCTION s.txt_caller"), strings.Index(aritySingle, "CREATE FUNCTION s.length"))

	// Unqualified non-public procedures set and reset search_path around creation.
	procOut, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "reporting",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "p_unqual",
				Definition: "CREATE PROCEDURE p_unqual() AS $$ BEGIN NULL; END; $$ LANGUAGE plpgsql",
			}},
		}},
	})
	require.NoError(t, err)
	require.Contains(t, procOut, "SET search_path TO \"reporting\", public;\nCREATE PROCEDURE p_unqual()")
	require.Contains(t, procOut, "RESET search_path;\n")

	// Procedure identity / kind / signature mismatches and complete view target mismatches fail explicitly.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "p",
				Definition: "CREATE PROCEDURE s.q() AS $$ BEGIN NULL; END; $$ LANGUAGE plpgsql",
			}},
		}},
	})
	require.ErrorContains(t, err, "identity")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "p",
				Definition: "CREATE FUNCTION s.p() RETURNS int IMMUTABLE AS $$ SELECT 1 $$ LANGUAGE SQL",
			}},
		}},
	})
	require.ErrorContains(t, err, "CREATE PROCEDURE")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "p",
				Signature:  "(int)",
				Definition: "CREATE PROCEDURE s.p(v varchar(20)) AS $$ BEGIN NULL; END; $$ LANGUAGE plpgsql",
			}},
		}},
	})
	require.ErrorContains(t, err, "signature")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "CREATE VIEW s.w AS SELECT 1 WITH NO SCHEMA BINDING",
			}},
		}},
	})
	require.ErrorContains(t, err, "identity")

	// Unsupported secondary indexes (neither primary nor unique) fail explicitly.
	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:    "t_sec_idx",
		Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
		Indexes: []*metadatapb.IndexMetadata{{Name: "idx_t_id", Expressions: []string{"id"}}},
	})
	require.ErrorContains(t, err, "secondary index")

	// pg_catalog precedes local UDFs for same-signature unqualified calls.
	builtinPrec, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "abs",
					Definition: "CREATE FUNCTION s.abs(x int4) RETURNS int4 IMMUTABLE AS $$ SELECT s.builtin_caller(x) $$ LANGUAGE SQL",
				},
				{
					Name:       "builtin_caller",
					Definition: "CREATE FUNCTION s.builtin_caller(n int4) RETURNS int4 IMMUTABLE AS $$ SELECT abs(n) $$ LANGUAGE SQL",
				},
			},
		}},
	})
	require.NoError(t, err)
	require.Less(t, strings.Index(builtinPrec, "CREATE FUNCTION s.builtin_caller("), strings.Index(builtinPrec, "CREATE FUNCTION s.abs("))

	// public schema functions and procedures set/reset search_path too.
	require.Contains(t, builtinPrec, "SET search_path TO \"s\", public;\nCREATE FUNCTION s.builtin_caller")
	pubProcOut, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "p_pub",
				Definition: "CREATE PROCEDURE p_pub() AS $$ BEGIN NULL; END; $$ LANGUAGE plpgsql",
			}},
		}},
	})
	require.NoError(t, err)
	require.Contains(t, pubProcOut, "SET search_path TO \"public\";\nCREATE PROCEDURE p_pub()")
	require.Contains(t, pubProcOut, "RESET search_path;\n")

	// CHAR (bpchar) stays distinct from VARCHAR so VARCHAR args do not match CHAR overloads.
	charVsVarchar, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "str_fn",
					Definition: "CREATE FUNCTION s.str_fn(c char(4)) RETURNS varchar IMMUTABLE AS $$ SELECT s.str_wrap('ok'::varchar) $$ LANGUAGE SQL",
				},
				{
					Name:       "str_wrap",
					Definition: "CREATE FUNCTION s.str_wrap(v varchar) RETURNS varchar IMMUTABLE AS $$ SELECT s.str_fn(v) $$ LANGUAGE SQL",
				},
				{
					Name:       "str_fn",
					Definition: "CREATE FUNCTION s.str_fn(v varchar) RETURNS varchar IMMUTABLE AS $$ SELECT v $$ LANGUAGE SQL",
				},
			},
		}},
	})
	require.NoError(t, err)
	require.Less(t, strings.Index(charVsVarchar, "CREATE FUNCTION s.str_fn(v varchar)"), strings.Index(charVsVarchar, "CREATE FUNCTION s.str_wrap("))
	require.Less(t, strings.Index(charVsVarchar, "CREATE FUNCTION s.str_wrap("), strings.Index(charVsVarchar, "CREATE FUNCTION s.str_fn(c char(4))"))

	// Non-public views wrap creation in schema search_path so unqualified relations resolve.
	require.Contains(t, deferredOwnersAndComments, "SET search_path TO \"reporting\", public;\nCREATE OR REPLACE VIEW \"reporting\".\"v_line_comment\"")

	// Unrendered view triggers and rules fail explicitly.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_trig",
				Definition: "SELECT 1",
				Triggers:   []*metadatapb.TriggerMetadata{{Name: "tr_v", Body: "CREATE TRIGGER tr_v ..."}},
			}},
		}},
	})
	require.ErrorContains(t, err, "triggers")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_rule",
				Definition: "SELECT 1",
				Rules:      []*metadatapb.RuleMetadata{{Name: "rl_v"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "rules")

	// Public views (including query-only and complete CREATE VIEW public.v) and public materialized views set search_path.
	pubViews, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Views: []*metadatapb.ViewMetadata{
				{Name: "v_qual", Definition: "CREATE VIEW public.v_qual AS SELECT id FROM t"},
				{Name: "v_query", Definition: "SELECT id FROM t"},
			},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{
				{Name: "mv_pub", Definition: "SELECT id FROM t"},
			},
		}},
	})
	require.NoError(t, err)
	require.Contains(t, pubViews, "SET search_path TO \"public\";\nCREATE VIEW public.v_qual AS SELECT id FROM t;\nRESET search_path;\n")
	require.Contains(t, pubViews, "SET search_path TO \"public\";\nCREATE OR REPLACE VIEW \"public\".\"v_query\" AS SELECT id FROM t;\nRESET search_path;\n")
	require.Contains(t, pubViews, "SET search_path TO \"public\";\nCREATE MATERIALIZED VIEW \"public\".\"mv_pub\" AS SELECT id FROM t;\nRESET search_path;\n")

	// Unrendered extensions, late-binding view column comments, and FKs missing ReferencedTable fail explicitly.
	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Name:       "dev",
		Extensions: []*metadatapb.ExtensionMetadata{{Name: "plpgsql"}},
	})
	require.ErrorContains(t, err, "unsupported")

	// Multiword routine parameter types (timestamp with time zone) match Signature metadata.
	multiwordSig, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "ts_fn",
				Signature:  "(timestamp with time zone)",
				Definition: "CREATE FUNCTION s.ts_fn(ts timestamptz) RETURNS int IMMUTABLE AS $$ SELECT 1 $$ LANGUAGE SQL",
			}},
		}},
	})
	require.NoError(t, err)
	require.Contains(t, multiwordSig, "CREATE FUNCTION s.ts_fn(ts timestamptz)")

	_, err = DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "late_cmt",
				Definition: `CREATE VIEW "public"."late_cmt" AS SELECT 1 AS id WITH NO SCHEMA BINDING`,
				Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Comment: "unsupported"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "late-binding")

	_, err = GetTableDefinition("public", &metadatapb.TableMetadata{
		Name:    "fk_missing_tbl",
		Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
			Name:              "fk_empty_ref",
			Columns:           []string{"id"},
			ReferencedColumns: []string{"id"},
		}},
	})
	require.ErrorContains(t, err, "referenced table")
}

// TestDatabaseDefinitionDevSnapshotReplay verifies DatabaseDefinition on the
// recorded Redshift 1.0.539507 cluster sync snapshot (../driver/testdata/sync/dev.json),
// parsing every statement and executing it against the Redshift catalog simulator.
func TestDatabaseDefinitionDevSnapshotReplay(t *testing.T) {
	body, err := os.ReadFile("../driver/testdata/sync/dev.json")
	require.NoError(t, err)
	db := &metadatapb.DatabaseSchemaMetadata{}
	require.NoError(t, protojson.Unmarshal(body, db))

	ddl, err := DatabaseDefinition(db)
	require.NoError(t, err)
	require.NotEmpty(t, ddl)

	lookupPos := strings.Index(ddl, `CREATE TABLE "public"."lookup"`)
	factsPos := strings.Index(ddl, `CREATE TABLE "app"."facts"`)
	require.Greater(t, lookupPos, -1)
	require.Greater(t, factsPos, lookupPos)

	_, err = pgparser.Parse(ddl)
	require.NoError(t, err)

	results, err := New().Exec(ddl, nil)
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error)
	}
}
