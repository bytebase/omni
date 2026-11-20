package deparse

import (
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetDatabaseDefinition pins the DDL the deparser emits. The first three
// table cases are lifted verbatim (inputs and table text) from
// bytebase/backend/plugin/schema/trino/get_database_definition_test.go; the
// CREATE SCHEMA and CREATE VIEW statements are the fidelity additions that make
// a dump restore into a catalog that does not already hold them. The remaining
// cases extend coverage to multiple tables, multiple schemas, ordering, and
// nil/empty handling.
func TestGetDatabaseDefinition(t *testing.T) {
	tests := []struct {
		name     string
		metadata *DatabaseSchemaMetadata
		expected string
	}{
		{
			// Verbatim from the legacy suite: a non-nullable bigint and a
			// nullable varchar, in declaration order.
			name: "simple table (legacy golden)",
			metadata: oneTable("testcatalog", "testschema", &TableMetadata{
				Name: "testtable",
				Columns: []*ColumnMetadata{
					{Name: "id", Type: "bigint", Nullable: false},
					{Name: "name", Type: "varchar", Nullable: true},
				},
			}),
			expected: `CREATE SCHEMA IF NOT EXISTS "testschema";

CREATE TABLE IF NOT EXISTS "testschema"."testtable" (
    "id" bigint NOT NULL,
    "name" varchar
);

`,
		},
		{
			// Verbatim from the legacy suite: a table with no columns still
			// emits the parens with the blank line the legacy format leaves.
			name: "empty columns (legacy golden)",
			metadata: oneTable("testcatalog", "testschema", &TableMetadata{
				Name:    "empty_table",
				Columns: []*ColumnMetadata{},
			}),
			expected: `CREATE SCHEMA IF NOT EXISTS "testschema";

CREATE TABLE IF NOT EXISTS "testschema"."empty_table" (

);

`,
		},
		{
			// Verbatim from the legacy suite: special characters in identifiers
			// are emitted as-is inside the double quotes.
			name: "special characters in identifiers (legacy golden)",
			metadata: oneTable("test-catalog", "test_schema", &TableMetadata{
				Name: "test.table",
				Columns: []*ColumnMetadata{
					{Name: "id-field", Type: "bigint", Nullable: false},
				},
			}),
			expected: `CREATE SCHEMA IF NOT EXISTS "test_schema";

CREATE TABLE IF NOT EXISTS "test_schema"."test.table" (
    "id-field" bigint NOT NULL
);

`,
		},
		{
			// A double quote in any name is written as two double quotes, which
			// Trino reads back as one inside the identifier.
			name: "embedded double quotes",
			metadata: oneTable("cat", `s"1`, &TableMetadata{
				Name: `t"2`,
				Columns: []*ColumnMetadata{
					{Name: `c"3`, Type: "varchar", Nullable: true},
				},
			}),
			expected: `CREATE SCHEMA IF NOT EXISTS "s""1";

CREATE TABLE IF NOT EXISTS "s""1"."t""2" (
    "c""3" varchar
);

`,
		},
		{
			// Two tables in one schema render in snapshot order, each as its own
			// statement separated by the trailing blank line.
			name: "multiple tables preserve order",
			metadata: &DatabaseSchemaMetadata{
				Name: "cat",
				Schemas: []*SchemaMetadata{
					{
						Name: "s",
						Tables: []*TableMetadata{
							{Name: "first", Columns: []*ColumnMetadata{{Name: "a", Type: "integer", Nullable: false}}},
							{Name: "second", Columns: []*ColumnMetadata{{Name: "b", Type: "varchar", Nullable: true}}},
						},
					},
				},
			},
			expected: `CREATE SCHEMA IF NOT EXISTS "s";

CREATE TABLE IF NOT EXISTS "s"."first" (
    "a" integer NOT NULL
);

CREATE TABLE IF NOT EXISTS "s"."second" (
    "b" varchar
);

`,
		},
		{
			// Every CREATE SCHEMA comes before every table, so a table can
			// reference one in a schema the dump has not written yet.
			name: "multiple schemas",
			metadata: &DatabaseSchemaMetadata{
				Name: "cat",
				Schemas: []*SchemaMetadata{
					{Name: "s1", Tables: []*TableMetadata{{Name: "t1", Columns: []*ColumnMetadata{{Name: "c1", Type: "bigint", Nullable: false}}}}},
					{Name: "s2", Tables: []*TableMetadata{{Name: "t2", Columns: []*ColumnMetadata{{Name: "c2", Type: "double", Nullable: true}}}}},
				},
			},
			expected: `CREATE SCHEMA IF NOT EXISTS "s1";

CREATE SCHEMA IF NOT EXISTS "s2";

CREATE TABLE IF NOT EXISTS "s1"."t1" (
    "c1" bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "s2"."t2" (
    "c2" double
);

`,
		},
		{
			// A schema with no relations still restores as a schema.
			name:     "schema with no tables",
			metadata: &DatabaseSchemaMetadata{Name: "cat", Schemas: []*SchemaMetadata{{Name: "empty"}}},
			expected: `CREATE SCHEMA IF NOT EXISTS "empty";

`,
		},
		{
			// Every table is written before any view, because a view may read a
			// table in another schema.
			name: "view",
			metadata: &DatabaseSchemaMetadata{
				Name: "cat",
				Schemas: []*SchemaMetadata{
					{
						Name:   "s",
						Tables: []*TableMetadata{{Name: "t", Columns: []*ColumnMetadata{{Name: "id", Type: "integer", Nullable: true}}}},
						Views:  []*ViewMetadata{{Name: "v", Definition: "SELECT id\nFROM\n  \"s\".\"t\"\n"}},
					},
				},
			},
			expected: `CREATE SCHEMA IF NOT EXISTS "s";

CREATE TABLE IF NOT EXISTS "s"."t" (
    "id" integer
);

CREATE OR REPLACE VIEW "s"."v" AS
SELECT id
FROM
  "s"."t";

`,
		},
		{
			// A view that reads another view follows it, even when its name
			// sorts first.
			name: "views follow their dependencies",
			metadata: &DatabaseSchemaMetadata{
				Name: "cat",
				Schemas: []*SchemaMetadata{
					{
						Name:   "s",
						Tables: []*TableMetadata{{Name: "t", Columns: []*ColumnMetadata{{Name: "id", Type: "integer", Nullable: true}}}},
						Views: []*ViewMetadata{
							{Name: "a_summary", Definition: `SELECT id FROM "s"."z_base"`},
							{Name: "z_base", Definition: `SELECT id FROM "s"."t"`},
						},
					},
				},
			},
			expected: `CREATE SCHEMA IF NOT EXISTS "s";

CREATE TABLE IF NOT EXISTS "s"."t" (
    "id" integer
);

CREATE OR REPLACE VIEW "s"."z_base" AS
SELECT id FROM "s"."t";

CREATE OR REPLACE VIEW "s"."a_summary" AS
SELECT id FROM "s"."z_base";

`,
		},
		{
			// No schemas at all yields the empty string.
			name:     "no schemas",
			metadata: &DatabaseSchemaMetadata{Name: "cat"},
			expected: "",
		},
		{
			// A nil metadata is handled gracefully.
			name:     "nil metadata",
			metadata: nil,
			expected: "",
		},
		{
			// nil entries in the slices are skipped rather than panicking.
			name: "nil schema, table, column, and view entries are skipped",
			metadata: &DatabaseSchemaMetadata{
				Name: "cat",
				Schemas: []*SchemaMetadata{
					nil,
					{
						Name: "s",
						Tables: []*TableMetadata{
							nil,
							{Name: "t", Columns: []*ColumnMetadata{nil, {Name: "c", Type: "integer", Nullable: false}}},
						},
						Views: []*ViewMetadata{nil},
					},
				},
			},
			expected: `CREATE SCHEMA IF NOT EXISTS "s";

CREATE TABLE IF NOT EXISTS "s"."t" (
    "c" integer NOT NULL
);

`,
		},
	}

	a := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			result, err := GetDatabaseDefinition(tt.metadata)
			a.NoError(err)
			a.Equal(tt.expected, result, tt.name)
		})
	}
}

// TestGetDatabaseDefinitionRefusesAViewWithoutDefinition pins that a view the
// dump cannot recreate is an error, not a statement that silently drops it.
func TestGetDatabaseDefinitionRefusesAViewWithoutDefinition(t *testing.T) {
	_, err := GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name:    "cat",
		Schemas: []*SchemaMetadata{{Name: "s", Views: []*ViewMetadata{{Name: "v"}}}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `view s.v has no definition`)
}

// TestWriteCreateTable pins the single-statement renderer, mirroring the legacy
// TestWriteCreateTable (same input, same expected — no trailing blank line, the
// caller adds that).
func TestWriteCreateTable(t *testing.T) {
	tests := []struct {
		name   string
		schema string
		table  *TableMetadata
		// expected is the single table statement, without the trailing blank
		// line.
		expected string
	}{
		{
			name:   "simple table (legacy golden)",
			schema: "schema",
			table: &TableMetadata{
				Name:    "table",
				Columns: []*ColumnMetadata{{Name: "col1", Type: "integer", Nullable: false}},
			},
			expected: `CREATE TABLE IF NOT EXISTS "schema"."table" (
    "col1" integer NOT NULL
);`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			writeCreateTable(&buf, tt.schema, tt.table)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

// oneTable builds a DatabaseSchemaMetadata holding a single schema with a single
// table, matching how the legacy test assembled its fixtures.
func oneTable(catalog, schema string, table *TableMetadata) *DatabaseSchemaMetadata {
	return &DatabaseSchemaMetadata{
		Name: catalog,
		Schemas: []*SchemaMetadata{
			{Name: schema, Tables: []*TableMetadata{table}},
		},
	}
}

func TestDatabaseDefinitionFromProto(t *testing.T) {
	gotNil, err := DatabaseDefinitionFromProto(nil)
	require.NoError(t, err)
	assert.Equal(t, "", gotNil)

	protoMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{
			nil,
			{Name: "skipped_schema", SkipDump: true, Tables: []*metadatapb.TableMetadata{{Name: "t_hidden"}}},
			{
				Name:  "s",
				Owner: "schema_admin",
				Functions: []*metadatapb.FunctionMetadata{
					nil,
					{Name: "fn_skip", Definition: "CREATE FUNCTION s.fn_skip() RETURNS int RETURN 1", SkipDump: true},
					{Name: "caller", Definition: "CREATE FUNCTION s.caller(x bigint) RETURNS bigint RETURN s.inc(x) * 2"},
					{Name: "inc", Definition: "CREATE FUNCTION s.inc(x bigint) RETURNS bigint RETURN x + 1 -- udf note", Comment: "inc's doc"},
				},
				Tables: []*metadatapb.TableMetadata{
					nil,
					{Name: "t_skip", SkipDump: true},
					{
						Name:    "t",
						Owner:   "tbl_owner",
						Comment: "events table",
						Columns: []*metadatapb.ColumnMetadata{
							nil,
							{Name: "id", Type: "bigint", Nullable: false, Default: "0", Comment: "pk col"},
							{Name: "label", Type: "varchar", Nullable: true, Default: "pending"},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					nil,
					{Name: "v_skip", Definition: "SELECT 1", SkipDump: true},
					{
						Name:       "v_top",
						Definition: `SELECT s.caller(id) FROM "s"."v_base" -- trailing view comment`,
						Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Comment: "computed id"}},
					},
					{Name: "v_base", Definition: `SELECT id FROM "s"."mv_base"`, Comment: "ordinary note"},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					nil,
					{Name: "mv_skip", Definition: "SELECT 1", SkipDump: true},
					{
						Name:              "mv_top_unparsed",
						Definition:        "SELECT /* syntax unknown to GetQuerySpan */ * FROM TABLE(exclude_columns(TABLE(s.mv_base), DESCRIPTOR(x)))",
						DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "s", Table: "mv_base"}},
					},
					{Name: "mv_base", Definition: `SELECT id FROM "s"."t" -- trailing mv comment`, Comment: "mv's note"},
				},
			},
		},
	}
	got, err := DatabaseDefinitionFromProto(protoMeta)
	require.NoError(t, err)
	assert.Equal(t, `CREATE SCHEMA IF NOT EXISTS "s";

CREATE FUNCTION s.inc(x bigint) RETURNS bigint COMMENT 'inc''s doc' RETURN x + 1 -- udf note
;

CREATE FUNCTION s.caller(x bigint) RETURNS bigint RETURN s.inc(x) * 2;

CREATE TABLE IF NOT EXISTS "s"."t" (
    "id" bigint DEFAULT 0 NOT NULL COMMENT 'pk col',
    "label" varchar DEFAULT 'pending'
) COMMENT 'events table';

CREATE MATERIALIZED VIEW IF NOT EXISTS "s"."mv_base" COMMENT 'mv''s note' AS
SELECT id FROM "s"."t" -- trailing mv comment
;

CREATE OR REPLACE VIEW "s"."v_base" COMMENT 'ordinary note' AS
SELECT id FROM "s"."mv_base";

CREATE MATERIALIZED VIEW IF NOT EXISTS "s"."mv_top_unparsed" AS
SELECT /* syntax unknown to GetQuerySpan */ * FROM TABLE(exclude_columns(TABLE(s.mv_base), DESCRIPTOR(x)));

CREATE OR REPLACE VIEW "s"."v_top" AS
SELECT s.caller(id) FROM "s"."v_base" -- trailing view comment
;
COMMENT ON COLUMN "s"."v_top"."id" IS 'computed id';

ALTER TABLE "s"."t" SET AUTHORIZATION "tbl_owner";

ALTER SCHEMA "s" SET AUTHORIZATION "schema_admin";

`, got)

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "s",
			Tables: []*metadatapb.TableMetadata{{Name: "empty_t"}},
		}},
	})
	require.ErrorContains(t, err, "no columns")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "s",
			Functions: []*metadatapb.FunctionMetadata{{Name: "empty_fn"}},
		}},
	})
	require.ErrorContains(t, err, "no definition")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:    "s",
			Comment: "schema note",
		}},
	})
	require.ErrorContains(t, err, "COMMENT ON SCHEMA")

	// Same-arity overloaded UDFs resolved by argument type:
	// [f(varchar) calls wrap(bigint), wrap(bigint) calls f(bigint), f(bigint) independent].
	overloaded, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "f",
					Definition: "CREATE FUNCTION s.f(v varchar) RETURNS bigint RETURN s.wrap(length(v))",
				},
				{
					Name:       "wrap",
					Definition: "CREATE FUNCTION s.wrap(n bigint) RETURNS bigint RETURN s.f(n) + 1",
				},
				{
					Name:       "f",
					Definition: "CREATE FUNCTION s.f(n bigint) RETURNS bigint RETURN n * 2",
				},
			},
		}},
	})
	require.NoError(t, err)
	fBigintPos := strings.Index(overloaded, "CREATE FUNCTION s.f(n bigint)")
	wrapPos := strings.Index(overloaded, "CREATE FUNCTION s.wrap(n bigint)")
	fVarcharPos := strings.Index(overloaded, "CREATE FUNCTION s.f(v varchar)")
	require.Greater(t, wrapPos, fBigintPos, "expected f(bigint) before wrap(bigint):\n%s", overloaded)
	require.Greater(t, fVarcharPos, wrapPos, "expected wrap(bigint) before f(varchar):\n%s", overloaded)
}
