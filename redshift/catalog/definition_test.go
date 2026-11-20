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
						Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Comment: "ignored on late-binding view"}},
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

--
-- Materialized view structure for `+"`a_mv_top`"+`
--
CREATE MATERIALIZED VIEW "public"."a_mv_top" AS SELECT id FROM "public"."z_mv_base";
`, got)

	mixedView, err := GetViewDefinition("public", &metadatapb.ViewMetadata{Name: "MixedView", Definition: "SELECT 1"})
	require.NoError(t, err)
	require.Equal(t, "SET enable_case_sensitive_identifier TO true;\n\nCREATE OR REPLACE VIEW \"public\".\"MixedView\" AS SELECT 1;\n", mixedView)

	mixedColView, err := GetViewDefinition("public", &metadatapb.ViewMetadata{
		Name:       "v",
		Columns:    []*metadatapb.ColumnMetadata{{Name: "MixedAlias"}},
		Definition: `SELECT 1 AS "MixedAlias"`,
	})
	require.NoError(t, err)
	require.Equal(t, "SET enable_case_sensitive_identifier TO true;\n\nCREATE OR REPLACE VIEW \"public\".\"v\" AS SELECT 1 AS \"MixedAlias\";\n", mixedColView)

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
	require.Equal(t, `CREATE OR REPLACE VIEW "public"."v" AS SELECT 1 AS id;
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
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Default: "public.f_inc(0)"}},
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
	// table-dependent functions precede tables whose DEFAULT calls them, and
	// SkipDump external tables are ignored cleanly.
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
						Name: "consumer_tbl",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "code", Type: "integer", Default: "public.next_code_fn()"},
						},
					},
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
	consumerPos := strings.Index(unified, "Table structure for `consumer_tbl`")
	viewPos := strings.Index(unified, "View structure for `v_events`")
	vFnPos := strings.Index(unified, "Function structure for `v_summary_fn`")

	require.Greater(t, nextFnPos, basePos, "expected base_tbl before next_code_fn:\n%s", unified)
	require.Greater(t, consumerPos, nextFnPos, "expected next_code_fn before consumer_tbl:\n%s", unified)
	require.Greater(t, viewPos, basePos, "expected base_tbl before v_events:\n%s", unified)
	require.Greater(t, vFnPos, viewPos, "expected v_events before v_summary_fn:\n%s", unified)

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
}
