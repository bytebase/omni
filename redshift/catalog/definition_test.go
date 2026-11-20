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

// TestForeignKeyKeepsReferencedSchema pins the qualifier on the REFERENCES
// target, since Redshift sync strips it off ReferencedTable and keeps it in
// ReferencedSchema.
func TestForeignKeyKeepsReferencedSchema(t *testing.T) {
	got, err := GetTableDefinition("analytics", &metadatapb.TableMetadata{
		Name: "child",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "id", Type: "integer"},
			{Name: "parent_id", Type: "integer", Nullable: true},
		},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{
			{
				Name:              "child_parent_fk",
				Columns:           []string{"parent_id"},
				ReferencedSchema:  "warehouse",
				ReferencedTable:   "parent",
				ReferencedColumns: []string{"id"},
			},
		},
	})
	require.NoError(t, err)
	require.Contains(t, got, `REFERENCES "warehouse"."parent"("id")`)
	require.NotContains(t, got, `REFERENCES "parent"("id")`)
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
// pg_get_indexdef expressions (already quoted) and foreign key columns (bare).
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

	require.Contains(t, got, `PRIMARY KEY ("MixedCase")`)
	require.NotContains(t, got, `"""MixedCase"""`)
	require.Contains(t, got, `FOREIGN KEY ("parent_id") REFERENCES "public"."parent"("id")`)
}
