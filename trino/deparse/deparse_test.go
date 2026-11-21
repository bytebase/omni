package deparse

import (
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	trinoparser "github.com/bytebase/omni/trino/parser"
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
				Name: "s",
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
					{Name: "mv_base", Definition: `SELECT id FROM "s"."t" -- trailing mv comment`},
				},
			},
		},
	}
	got, err := DatabaseDefinitionFromProto(protoMeta)
	require.NoError(t, err)
	assert.Equal(t, `CREATE SCHEMA IF NOT EXISTS "s";

CREATE FUNCTION "cat".s.inc(x bigint) RETURNS bigint COMMENT 'inc''s doc' RETURN x + 1 -- udf note
;

CREATE FUNCTION "cat".s.caller(x bigint) RETURNS bigint RETURN "cat".s.inc(x) * 2;

CREATE TABLE IF NOT EXISTS "s"."t" (
    "id" bigint DEFAULT 0 NOT NULL COMMENT 'pk col',
    "label" varchar DEFAULT 'pending'
) COMMENT 'events table';

CREATE MATERIALIZED VIEW IF NOT EXISTS "s"."mv_base" AS
SELECT id FROM "s"."t" -- trailing mv comment
;

CREATE OR REPLACE VIEW "s"."v_base" COMMENT 'ordinary note' AS
SELECT id FROM "s"."mv_base";

CREATE MATERIALIZED VIEW IF NOT EXISTS "s"."mv_top_unparsed" AS
SELECT /* syntax unknown to GetQuerySpan */ * FROM TABLE(exclude_columns(TABLE(s.mv_base), DESCRIPTOR(x)));

CREATE OR REPLACE VIEW "s"."v_top" AS
SELECT "cat".s.caller(id) FROM "s"."v_base" -- trailing view comment
;
COMMENT ON COLUMN "s"."v_top"."id" IS 'computed id';

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
	fBigintPos := strings.Index(overloaded, "CREATE FUNCTION \"cat\".s.f(n bigint)")
	wrapPos := strings.Index(overloaded, "CREATE FUNCTION \"cat\".s.wrap(n bigint)")
	fVarcharPos := strings.Index(overloaded, "CREATE FUNCTION \"cat\".s.f(v varchar)")
	require.Greater(t, wrapPos, fBigintPos, "expected f(bigint) before wrap(bigint):\n%s", overloaded)
	require.Greater(t, fVarcharPos, wrapPos, "expected wrap(bigint) before f(varchar):\n%s", overloaded)

	// Small integer literal (1) selects f(integer) over f(bigint), and CASE
	// widens arms (1, 1.5) to decimal over bigint: snapshot order puts callers
	// before their callees to ensure cycles aren't formed and callees lead.
	litAndCase, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "num_fn",
					Definition: "CREATE FUNCTION s.num_fn(n bigint) RETURNS bigint RETURN s.lit_caller()",
				},
				{
					Name:       "lit_caller",
					Definition: "CREATE FUNCTION s.lit_caller() RETURNS bigint RETURN s.num_fn(1)",
				},
				{
					Name:       "num_fn",
					Definition: "CREATE FUNCTION s.num_fn(n integer) RETURNS bigint RETURN CAST(n AS bigint)",
				},
				{
					Name:       "wide_fn",
					Definition: "CREATE FUNCTION s.wide_fn(n bigint) RETURNS decimal RETURN s.case_caller(n)",
				},
				{
					Name:       "case_caller",
					Definition: "CREATE FUNCTION s.case_caller(n bigint) RETURNS decimal RETURN s.wide_fn(CASE WHEN n > 0 THEN 1 ELSE 1.5 END)",
				},
				{
					Name:       "wide_fn",
					Definition: "CREATE FUNCTION s.wide_fn(d decimal) RETURNS decimal RETURN d",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.num_fn(n integer)"), strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.lit_caller()"))
	assert.Less(t, strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.lit_caller()"), strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.num_fn(n bigint)"))
	assert.Less(t, strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.wide_fn(d decimal)"), strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.case_caller(n bigint)"))
	assert.Less(t, strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.case_caller(n bigint)"), strings.Index(litAndCase, "CREATE FUNCTION \"cat\".s.wide_fn(n bigint)"))

	// Case-insensitive duplicate columns fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name: "dup_cols",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Type: "bigint"},
					{Name: "ID", Type: "bigint"},
				},
			}},
		}},
	})
	require.ErrorContains(t, err, "duplicate column")

	// Case-insensitive relation lookup in orderedViews.
	caseViews, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "Cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "S",
			Views: []*metadatapb.ViewMetadata{
				{Name: "TopView", Definition: `SELECT id FROM "cat"."s"."baseview"`},
				{Name: "BaseView", Definition: "SELECT 1 AS id"},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(caseViews, `CREATE OR REPLACE VIEW "S"."BaseView"`), strings.Index(caseViews, `CREATE OR REPLACE VIEW "S"."TopView"`))

	// Callees inside SQL predicate expressions (=, IS NULL, BETWEEN, IN, LIKE)
	// must also precede their callers even when listed later in snapshot order.
	predFuncs, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "pred_caller",
					Definition: "CREATE FUNCTION s.pred_caller(x bigint, p varchar) RETURNS boolean BEGIN IF s.cmp_callee(x) = 1 AND s.null_callee(x) IS NOT NULL AND s.between_callee(x) BETWEEN 1 AND 10 AND s.in_callee(x) IN (1, 2) AND s.like_callee(p) LIKE 'a%' THEN RETURN true; END IF; RETURN false; END",
				},
				{
					Name:       "cmp_callee",
					Definition: "CREATE FUNCTION s.cmp_callee(x bigint) RETURNS bigint RETURN x",
				},
				{
					Name:       "null_callee",
					Definition: "CREATE FUNCTION s.null_callee(x bigint) RETURNS bigint RETURN x",
				},
				{
					Name:       "between_callee",
					Definition: "CREATE FUNCTION s.between_callee(x bigint) RETURNS bigint RETURN x",
				},
				{
					Name:       "in_callee",
					Definition: "CREATE FUNCTION s.in_callee(x bigint) RETURNS bigint RETURN x",
				},
				{
					Name:       "like_callee",
					Definition: "CREATE FUNCTION s.like_callee(p varchar) RETURNS varchar RETURN p",
				},
			},
		}},
	})
	require.NoError(t, err)
	callerPos := strings.Index(predFuncs, "CREATE FUNCTION \"cat\".s.pred_caller(")
	for _, callee := range []string{"cmp_callee", "null_callee", "between_callee", "in_callee", "like_callee"} {
		assert.Less(t, strings.Index(predFuncs, "CREATE FUNCTION \"cat\".s."+callee+"("), callerPos, "expected %s before pred_caller:\n%s", callee, predFuncs)
	}

	// Parameterized overloads (array(integer) vs array(varchar)) disambiguate
	// via element type instead of forming cycles through unrelated overloads.
	arrayOverloads, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "arr_fn",
					Definition: "CREATE FUNCTION s.arr_fn(a array(varchar)) RETURNS varchar RETURN s.arr_caller()",
				},
				{
					Name:       "arr_caller",
					Definition: "CREATE FUNCTION s.arr_caller() RETURNS bigint RETURN CAST(cardinality(s.arr_fn(ARRAY[1])) AS bigint)",
				},
				{
					Name:       "arr_fn",
					Definition: "CREATE FUNCTION s.arr_fn(a array(integer)) RETURNS array(integer) RETURN a",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(arrayOverloads, "CREATE FUNCTION \"cat\".s.arr_fn(a array(integer))"), strings.Index(arrayOverloads, "CREATE FUNCTION \"cat\".s.arr_caller()"))
	assert.Less(t, strings.Index(arrayOverloads, "CREATE FUNCTION \"cat\".s.arr_caller()"), strings.Index(arrayOverloads, "CREATE FUNCTION \"cat\".s.arr_fn(a array(varchar))"))

	// Non-skipped Trino procedures fail explicitly; SkipDump procedures are ignored.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Procedures: []*metadatapb.ProcedureMetadata{
				{Name: "sync_partitions"},
			},
		}},
	})
	require.ErrorContains(t, err, "procedure")

	skippedProcOut, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Procedures: []*metadatapb.ProcedureMetadata{
				{Name: "sync_partitions", SkipDump: true},
			},
		}},
	})
	require.NoError(t, err)
	assert.Equal(t, "CREATE SCHEMA IF NOT EXISTS \"s\";\n\n", skippedProcOut)

	// Cross-catalog calls (other_catalog.s.f()) do not create false dependencies
	// against local s.f.
	crossCatFuncs, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "a",
					Definition: "CREATE FUNCTION s.a() RETURNS bigint RETURN other_catalog.s.b()",
				},
				{
					Name:       "b",
					Definition: "CREATE FUNCTION s.b() RETURNS bigint RETURN s.a()",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(crossCatFuncs, "CREATE FUNCTION \"cat\".s.a()"), strings.Index(crossCatFuncs, "CREATE FUNCTION \"cat\".s.b()"))

	// Duplicate tables (case-insensitive) fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{
				{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}},
				{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}},
			},
		}},
	})
	require.ErrorContains(t, err, "duplicate table")

	// Array subscript propagates element type for overload disambiguation.
	subscriptOverloads, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "elem_fn",
					Definition: "CREATE FUNCTION s.elem_fn(v varchar) RETURNS varchar RETURN CAST(s.elem_caller(ARRAY[CAST(v AS integer)]) AS varchar)",
				},
				{
					Name:       "elem_caller",
					Definition: "CREATE FUNCTION s.elem_caller(vals array(integer)) RETURNS integer RETURN s.elem_fn(vals[1])",
				},
				{
					Name:       "elem_fn",
					Definition: "CREATE FUNCTION s.elem_fn(n integer) RETURNS integer RETURN n",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(subscriptOverloads, "CREATE FUNCTION \"cat\".s.elem_fn(n integer)"), strings.Index(subscriptOverloads, "CREATE FUNCTION \"cat\".s.elem_caller("))
	assert.Less(t, strings.Index(subscriptOverloads, "CREATE FUNCTION \"cat\".s.elem_caller("), strings.Index(subscriptOverloads, "CREATE FUNCTION \"cat\".s.elem_fn(v varchar)"))

	// Named ROW field dereference propagates field type into overload disambiguation.
	rowOverloads, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "row_fn",
					Definition: "CREATE FUNCTION s.row_fn(v varchar) RETURNS varchar RETURN CAST(s.row_caller(ROW(CAST(v AS integer))) AS varchar)",
				},
				{
					Name:       "row_caller",
					Definition: "CREATE FUNCTION s.row_caller(r ROW(a integer)) RETURNS integer RETURN s.row_fn(r.a)",
				},
				{
					Name:       "row_fn",
					Definition: "CREATE FUNCTION s.row_fn(n integer) RETURNS integer RETURN n",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(rowOverloads, "CREATE FUNCTION \"cat\".s.row_fn(n integer)"), strings.Index(rowOverloads, "CREATE FUNCTION \"cat\".s.row_caller("))
	assert.Less(t, strings.Index(rowOverloads, "CREATE FUNCTION \"cat\".s.row_caller("), strings.Index(rowOverloads, "CREATE FUNCTION \"cat\".s.row_fn(v varchar)"))

	// Duplicate view identifiers (case-insensitive) fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{
				{Name: "v", Definition: "SELECT 1"},
				{Name: "V", Definition: "SELECT 2"},
			},
		}},
	})
	require.ErrorContains(t, err, "duplicate view")

	// Special SQL functions (CURRENT_DATE, CURRENT_USER) infer result types for overload disambiguation.
	specialFuncs, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "dt_fn",
					Definition: "CREATE FUNCTION s.dt_fn(v varchar) RETURNS varchar RETURN s.dt_caller()",
				},
				{
					Name:       "dt_caller",
					Definition: "CREATE FUNCTION s.dt_caller() RETURNS date RETURN s.dt_fn(CURRENT_DATE)",
				},
				{
					Name:       "dt_fn",
					Definition: "CREATE FUNCTION s.dt_fn(d date) RETURNS date RETURN d",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(specialFuncs, "CREATE FUNCTION \"cat\".s.dt_fn(d date)"), strings.Index(specialFuncs, "CREATE FUNCTION \"cat\".s.dt_caller()"))
	assert.Less(t, strings.Index(specialFuncs, "CREATE FUNCTION \"cat\".s.dt_caller()"), strings.Index(specialFuncs, "CREATE FUNCTION \"cat\".s.dt_fn(v varchar)"))

	// Unrendered table constraints/indexes and function identity mismatches fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t_idx",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}},
				Indexes: []*metadatapb.IndexMetadata{{Name: "t_idx_pkey", Primary: true, Expressions: []string{"id"}}},
			}},
		}},
	})
	require.ErrorContains(t, err, "constraints")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "expected_name",
				Definition: "CREATE FUNCTION s.other_name() RETURNS bigint RETURN 1",
			}},
		}},
	})
	require.ErrorContains(t, err, "identity")

	// Single local candidate with different arity than a built-in call (length(x) vs length(varchar, varchar))
	// does not add a false edge.
	aritySingle, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "length",
					Definition: "CREATE FUNCTION s.length(a varchar, b varchar) RETURNS bigint RETURN s.str_caller(a)",
				},
				{
					Name:       "str_caller",
					Definition: "CREATE FUNCTION s.str_caller(x varchar) RETURNS bigint RETURN length(x)",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(aritySingle, "CREATE FUNCTION \"cat\".s.str_caller("), strings.Index(aritySingle, "CREATE FUNCTION \"cat\".s.length("))

	// Named arguments map by parameter name rather than call-site position.
	namedArgs, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "named_fn",
					Definition: "CREATE FUNCTION s.named_fn(a integer, b varchar) RETURNS varchar RETURN s.named_caller()",
				},
				{
					Name:       "named_caller",
					Definition: "CREATE FUNCTION s.named_caller() RETURNS varchar RETURN s.named_fn(b => 1, a => 'x')",
				},
				{
					Name:       "named_fn",
					Definition: "CREATE FUNCTION s.named_fn(a varchar, b integer) RETURNS varchar RETURN a",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(namedArgs, "CREATE FUNCTION \"cat\".s.named_fn(a varchar, b integer)"), strings.Index(namedArgs, "CREATE FUNCTION \"cat\".s.named_caller()"))
	assert.Less(t, strings.Index(namedArgs, "CREATE FUNCTION \"cat\".s.named_caller()"), strings.Index(namedArgs, "CREATE FUNCTION \"cat\".s.named_fn(a integer, b varchar)"))

	// Unqualified Trino function definitions and sibling calls are qualified with
	// catalog and schema, and MV comments/indexes on Trino 482 fail explicitly.
	unqualFn, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "analytics",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "f_caller",
					Definition: "CREATE FUNCTION f_caller() RETURNS bigint RETURN f_callee() + 1",
				},
				{
					Name:       "f_callee",
					Definition: "CREATE FUNCTION f_callee() RETURNS bigint RETURN 1",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, unqualFn, `CREATE FUNCTION "cat"."analytics".f_callee()`)
	assert.Contains(t, unqualFn, `CREATE FUNCTION "cat"."analytics".f_caller() RETURNS bigint RETURN "cat"."analytics".f_callee() + 1;`)
	assert.Less(t, strings.Index(unqualFn, `CREATE FUNCTION "cat"."analytics".f_callee()`), strings.Index(unqualFn, `CREATE FUNCTION "cat"."analytics".f_caller()`))

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:              "analytics",
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "mv_cmt", Definition: "SELECT 1", Comment: "unsupported on 482"}},
		}},
	})
	require.ErrorContains(t, err, "482")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "analytics",
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_idx",
				Definition: "SELECT 1",
				Indexes:    []*metadatapb.IndexMetadata{{Name: "idx_1", Expressions: []string{"id"}}},
			}},
		}},
	})
	require.ErrorContains(t, err, "indexes")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "analytics",
			Sequences: []*metadatapb.SequenceMetadata{{Name: "seq_1"}},
		}},
	})
	require.ErrorContains(t, err, "sequence")

	// Non-skipped external tables fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:           "analytics",
			ExternalTables: []*metadatapb.ExternalTableMetadata{{Name: "ext_t"}},
		}},
	})
	require.ErrorContains(t, err, "external table")

	// AT TIME ZONE produces a zoned datetime type and CHAR->VARCHAR coercion is one-way.
	tzAndChar, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "tz_fn",
					Definition: "CREATE FUNCTION s.tz_fn(t timestamp) RETURNS varchar RETURN s.tz_caller()",
				},
				{
					Name:       "tz_caller",
					Definition: "CREATE FUNCTION s.tz_caller() RETURNS varchar RETURN s.tz_fn(CAST('2026-01-01 00:00:00' AS timestamp) AT TIME ZONE 'UTC')",
				},
				{
					Name:       "tz_fn",
					Definition: "CREATE FUNCTION s.tz_fn(t timestamp with time zone) RETURNS varchar RETURN 'ok'",
				},
				{
					Name:       "chr_fn",
					Definition: "CREATE FUNCTION s.chr_fn(c char) RETURNS varchar RETURN s.chr_caller('hi')",
				},
				{
					Name:       "chr_caller",
					Definition: "CREATE FUNCTION s.chr_caller(v varchar) RETURNS varchar RETURN s.chr_fn(v)",
				},
				{
					Name:       "chr_fn",
					Definition: "CREATE FUNCTION s.chr_fn(v varchar) RETURNS varchar RETURN v",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Less(t, strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.tz_fn(t timestamp with time zone)"), strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.tz_caller()"))
	assert.Less(t, strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.tz_caller()"), strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.tz_fn(t timestamp)"))
	assert.Less(t, strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.chr_fn(v varchar)"), strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.chr_caller("))
	assert.Less(t, strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.chr_caller("), strings.Index(tzAndChar, "CREATE FUNCTION \"cat\".s.chr_fn(c char)"))

	// FunctionMetadata.Signature parameter type mismatches fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "sig_fn",
				Signature:  "(bigint)",
				Definition: "CREATE FUNCTION s.sig_fn(v varchar) RETURNS varchar RETURN v",
			}},
		}},
	})
	require.ErrorContains(t, err, "signature")

	// View triggers and rules cannot be restored in Trino and fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_trig",
				Definition: "SELECT 1",
				Triggers:   []*metadatapb.TriggerMetadata{{Name: "tr_v", Body: "CREATE TRIGGER tr_v ..."}},
			}},
		}},
	})
	require.ErrorContains(t, err, "triggers or rules")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_rule",
				Definition: "SELECT 1",
				Rules:      []*metadatapb.RuleMetadata{{Name: "rl_v"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "triggers or rules")

	// Python UDFs with AS $$...$$ bodies preserve and replace metadata comments.
	pyOut, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "py_wrap",
					Definition: "CREATE FUNCTION py_wrap(x bigint) RETURNS bigint RETURN py_add(x, 1)",
				},
				{
					Name:       "py_add",
					Comment:    "adds one's arg",
					Definition: "CREATE FUNCTION py_add(x bigint, y bigint) RETURNS bigint LANGUAGE PYTHON WITH (handler = 'add') AS $$def add(x, y): return x + y$$",
				},
				{
					Name:       "py_replace",
					Comment:    "updated comment",
					Definition: "CREATE FUNCTION s.py_replace(x bigint) RETURNS bigint LANGUAGE PYTHON COMMENT 'stale' WITH (handler = 'id') AS $$def id(x): return x$$",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, pyOut, `CREATE FUNCTION "cat"."s".py_add(x bigint, y bigint) RETURNS bigint LANGUAGE PYTHON WITH (handler = 'add') COMMENT 'adds one''s arg' AS $$def add(x, y): return x + y$$;`)
	assert.Contains(t, pyOut, `CREATE FUNCTION "cat".s.py_replace(x bigint) RETURNS bigint LANGUAGE PYTHON COMMENT 'updated comment' WITH (handler = 'id') AS $$def id(x): return x$$;`)
	assert.Less(t, strings.Index(pyOut, `CREATE FUNCTION "cat"."s".py_add(`), strings.Index(pyOut, `CREATE FUNCTION "cat"."s".py_wrap(`))

	// Unrendered schema collations and range types fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:       "s",
			Collations: []*metadatapb.CollationMetadata{{Name: "custom_coll"}},
		}},
	})
	require.ErrorContains(t, err, "unsupported")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:       "s",
			RangeTypes: []*metadatapb.RangeTypeMetadata{{Name: "floatrange"}},
		}},
	})
	require.ErrorContains(t, err, "unsupported")

	// Unrendered database-level objects and generated columns fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name:          "cat",
		EventTriggers: []*metadatapb.EventTriggerMetadata{{Name: "evt"}},
	})
	require.ErrorContains(t, err, "unsupported")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t_gen",
				Columns: []*metadatapb.ColumnMetadata{{
					Name:       "c",
					Type:       "bigint",
					Generation: &metadatapb.GenerationMetadata{Expression: "a + 1"},
				}},
			}},
		}},
	})
	require.ErrorContains(t, err, "generated expression")

	// Unqualified built-in calls (length(x)) stay unqualified even when a same-signature
	// catalog UDF exists in the caller's schema, whereas explicitly qualified calls resolve locally.
	builtinShadow, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "length",
					Definition: "CREATE FUNCTION s.length(x varchar) RETURNS bigint RETURN CAST(x AS bigint)",
				},
				{
					Name:       "use_builtin",
					Definition: "CREATE FUNCTION s.use_builtin(x varchar) RETURNS bigint RETURN length(x)",
				},
				{
					Name:       "use_qual",
					Definition: "CREATE FUNCTION s.use_qual(x varchar) RETURNS bigint RETURN s.length(x)",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, builtinShadow, `CREATE FUNCTION "cat".s.use_builtin(x varchar) RETURNS bigint RETURN length(x);`)
	assert.Contains(t, builtinShadow, `CREATE FUNCTION "cat".s.use_qual(x varchar) RETURNS bigint RETURN "cat".s.length(x);`)
	assert.Less(t, strings.Index(builtinShadow, `CREATE FUNCTION "cat".s.length(`), strings.Index(builtinShadow, `CREATE FUNCTION "cat".s.use_qual(`))

	// SQL UDFs containing queries or subqueries fail explicitly (Trino disallows them),
	// whereas scalar subqueries inside views qualify nested UDF calls.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "sub_caller",
				Definition: "CREATE FUNCTION s.sub_caller(x bigint) RETURNS bigint RETURN (SELECT x)",
			}},
		}},
	})
	require.ErrorContains(t, err, "cannot process")

	// Unqualified relations in Trino view/MV bodies are qualified with catalog and schema
	// while CTE references stay unqualified.
	unqualViewOut, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_unqual",
				Definition: "WITH cte AS (SELECT id FROM t) SELECT id FROM cte",
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_unqual",
				Definition: "SELECT id FROM t",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, unqualViewOut, `WITH cte AS (SELECT id FROM "cat"."s".t) SELECT id FROM cte;`)
	assert.Contains(t, unqualViewOut, "CREATE MATERIALIZED VIEW IF NOT EXISTS \"s\".\"mv_unqual\" AS\nSELECT id FROM \"cat\".\"s\".t;")

	// Ambiguous schema/table owners (user vs role) and table columns without a type fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name:    "cat",
		Schemas: []*metadatapb.SchemaMetadata{{Name: "s", Owner: "role_a"}},
	})
	require.ErrorContains(t, err, "owner")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t_empty_type",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "no data type")

	// View/MV bodies qualify unqualified relations inside JSON expressions and
	// qualify UDF calls while preserving unqualified built-ins (length).
	viewExprAndUDF, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "my_fn",
				Definition: "CREATE FUNCTION s.my_fn(x bigint) RETURNS bigint RETURN x + 1",
			}},
			Tables: []*metadatapb.TableMetadata{{
				Name: "t",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Type: "bigint"},
					{Name: "name", Type: "varchar"},
				},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_json_udf",
				Definition: "SELECT JSON_OBJECT('id' VALUE (SELECT my_fn(id) FROM t)), length(name) FROM t",
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_udf",
				Definition: "SELECT s.my_fn(id) FROM t",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, viewExprAndUDF, `SELECT JSON_OBJECT('id' VALUE (SELECT "cat"."s".my_fn(id) FROM "cat"."s".t)), length(name) FROM "cat"."s".t;`)
	assert.Contains(t, viewExprAndUDF, "CREATE MATERIALIZED VIEW IF NOT EXISTS \"s\".\"mv_udf\" AS\nSELECT \"cat\".s.my_fn(id) FROM \"cat\".\"s\".t;")

	// Non-CREATE FUNCTION definitions, cyclic UDFs, duplicate function signatures,
	// and view column defaults fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "not_fn",
				Definition: "CREATE TABLE s.not_fn(id bigint)",
			}},
		}},
	})
	require.ErrorContains(t, err, "CREATE FUNCTION")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "c_a", Definition: "CREATE FUNCTION s.c_a(x bigint) RETURNS bigint RETURN s.c_b(x)"},
				{Name: "c_b", Definition: "CREATE FUNCTION s.c_b(x bigint) RETURNS bigint RETURN s.c_a(x)"},
			},
		}},
	})
	require.ErrorContains(t, err, "cyclic")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "dup_f", Definition: "CREATE FUNCTION s.dup_f(x int) RETURNS int RETURN x"},
				{Name: "dup_f", Definition: "CREATE FUNCTION s.dup_f(y integer) RETURNS integer RETURN y + 1"},
			},
		}},
	})
	require.ErrorContains(t, err, "duplicate function")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_def_col",
				Definition: "SELECT 1 AS id",
				Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Default: "1"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "default cannot be restored")

	// JSON_TABLE relation inputs/defaults and interval-qualified UDF overloads.
	jsonTableAndIntervals, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "make_json", Definition: "CREATE FUNCTION s.make_json() RETURNS varchar RETURN '{\"a\":1,\"b\":[{\"n\":2}]}'"},
				{Name: "fallback_fn", Definition: "CREATE FUNCTION s.fallback_fn() RETURNS integer RETURN 0"},
				{Name: "nested_fn", Definition: "CREATE FUNCTION s.nested_fn() RETURNS integer RETURN -1"},
				{Name: "fmt_iv", Definition: "CREATE FUNCTION s.fmt_iv(x INTERVAL YEAR TO MONTH) RETURNS varchar RETURN 'ym'"},
				{Name: "fmt_iv", Definition: "CREATE FUNCTION s.fmt_iv(x INTERVAL DAY TO SECOND) RETURNS varchar RETURN s.tag_ds()"},
				{Name: "tag_ds", Definition: "CREATE FUNCTION s.tag_ds() RETURNS varchar RETURN 'ds'"},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_jt_iv",
				Definition: "SELECT jt.v, jt.n, s.fmt_iv(INTERVAL '1-2' YEAR TO MONTH) AS ym, s.fmt_iv(INTERVAL '1 02:03:04' DAY TO SECOND) AS ds FROM JSON_TABLE(s.make_json(), 'lax $' COLUMNS (v INTEGER PATH '$.a' DEFAULT s.fallback_fn() ON EMPTY, NESTED PATH '$.b[*]' COLUMNS (n INTEGER PATH '$.n' DEFAULT s.nested_fn() ON ERROR))) AS jt",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, jsonTableAndIntervals, `JSON_TABLE("cat".s.make_json(), 'lax $'`)
	assert.Contains(t, jsonTableAndIntervals, `DEFAULT "cat".s.fallback_fn() ON EMPTY`)
	assert.Contains(t, jsonTableAndIntervals, `DEFAULT "cat".s.nested_fn() ON ERROR`)
	assert.Contains(t, jsonTableAndIntervals, `"cat".s.fmt_iv(INTERVAL '1-2' YEAR TO MONTH)`)
	assert.Less(t, strings.Index(jsonTableAndIntervals, `CREATE FUNCTION "cat".s.tag_ds(`), strings.Index(jsonTableAndIntervals, `CREATE FUNCTION "cat".s.fmt_iv(x INTERVAL DAY TO SECOND)`))

	// MATCH_RECOGNIZE relation qualification, variadic map_concat/concat_ws built-ins,
	// and cyclic view rejection.
	mrAndVariadic, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "events",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}, {Name: "ts", Type: "bigint"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "score_fn", Definition: "CREATE FUNCTION s.score_fn(v bigint) RETURNS bigint RETURN v + 1"},
				{Name: "map_concat", Definition: "CREATE FUNCTION s.map_concat(a map(varchar, integer)) RETURNS map(varchar, integer) RETURN a"},
				{Name: "concat_ws", Definition: "CREATE FUNCTION s.concat_ws(sep varchar, a varchar) RETURNS varchar RETURN a"},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_mr_var",
				Definition: "SELECT m.id, m.sc, map_concat(MAP(ARRAY['a'], ARRAY[1]), MAP(ARRAY['b'], ARRAY[2])), concat_ws(',', 'x', 'y', 'z') FROM events MATCH_RECOGNIZE (ORDER BY ts MEASURES A.id AS id, s.score_fn(A.id) AS sc ONE ROW PER MATCH PATTERN (A) DEFINE A AS A.id > 0 AND s.score_fn(A.id) > 0) AS m",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, mrAndVariadic, `FROM "cat"."s".events MATCH_RECOGNIZE`)
	assert.Contains(t, mrAndVariadic, `"cat".s.score_fn(A.id) AS sc`)
	assert.Contains(t, mrAndVariadic, `DEFINE A AS A.id > 0 AND "cat".s.score_fn(A.id) > 0`)
	assert.Contains(t, mrAndVariadic, `, map_concat(MAP(`)
	assert.Contains(t, mrAndVariadic, `, concat_ws(',', 'x', 'y', 'z')`)

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{
				{Name: "v1", Definition: "SELECT * FROM v2"},
				{Name: "v2", Definition: "SELECT * FROM v1"},
			},
		}},
	})
	require.ErrorContains(t, err, "cyclic view")

	// Inline WITH FUNCTION shadows unqualified calls while qualifying inner/schema calls.
	inlineUDFView, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "helper", Definition: "CREATE FUNCTION s.helper(x bigint) RETURNS bigint RETURN x + 100"},
				{Name: "cat_fn", Definition: "CREATE FUNCTION s.cat_fn(x bigint) RETURNS bigint RETURN x * 2"},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_inline",
				Definition: "WITH FUNCTION helper(x bigint) RETURNS bigint RETURN cat_fn(x) + 1 SELECT helper(5), s.helper(5)",
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_inline",
				Definition: "WITH FUNCTION helper(x bigint) RETURNS bigint RETURN cat_fn(x) + 2 SELECT helper(6), s.helper(6)",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, inlineUDFView, `WITH FUNCTION helper(x bigint) RETURNS bigint RETURN "cat"."s".cat_fn(x) + 1 SELECT helper(5), "cat".s.helper(5)`)
	assert.Contains(t, inlineUDFView, `WITH FUNCTION helper(x bigint) RETURNS bigint RETURN "cat"."s".cat_fn(x) + 2 SELECT helper(6), "cat".s.helper(6)`)

	// Self-recursive Trino SQL UDFs fail explicitly.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "rec_f",
				Definition: "CREATE FUNCTION s.rec_f(x bigint) RETURNS bigint RETURN CASE WHEN x <= 0 THEN 0 ELSE s.rec_f(x - 1) END",
			}},
		}},
	})
	require.ErrorContains(t, err, "self-recursive")

	// Lambda-taking built-ins (transform) take precedence over colliding UDFs,
	// unparsed CREATE FUNCTION declarations are not flagged as self-recursive,
	// and GetDatabaseDefinition qualifies direct-model function declarations.
	lambdaAndUnparsed, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "transform", Definition: "CREATE FUNCTION s.transform(a array(bigint), b bigint) RETURNS array(bigint) RETURN a"},
				{Name: "ext_fn", Definition: "CREATE FUNCTION s.ext_fn(x bigint) RETURNS bigint LANGUAGE PYTHON WITH (handler = 'mod.fn') AS 'body'"},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_transform",
				Definition: "SELECT transform(ARRAY[1, 2], x -> x + 1)",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, lambdaAndUnparsed, "SELECT transform(ARRAY[1, 2], x -> x + 1);")
	assert.Contains(t, lambdaAndUnparsed, `CREATE FUNCTION "cat".s.ext_fn(x bigint)`)

	directQualified, err := GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			Functions: []*FunctionMetadata{{
				Name:       "f",
				Definition: "CREATE FUNCTION f() RETURNS bigint RETURN 1",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, directQualified, `CREATE FUNCTION "cat"."s".f()`)

	// Unqualified Trino special forms (coalesce, if, nullif, try) are not
	// catalog-qualified even when a schema defines identically named UDFs.
	specialForms, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "coalesce", Definition: "CREATE FUNCTION s.coalesce(a bigint, b bigint) RETURNS bigint RETURN a"},
				{Name: "if", Definition: "CREATE FUNCTION s.if(c boolean, a bigint, b bigint) RETURNS bigint RETURN a"},
				{Name: "nullif", Definition: "CREATE FUNCTION s.nullif(a bigint, b bigint) RETURNS bigint RETURN a"},
				{Name: "try", Definition: "CREATE FUNCTION s.try(a bigint) RETURNS bigint RETURN a"},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_special",
				Definition: "SELECT coalesce(1, 0), if(true, 1, 2), nullif(1, 2), try(1), s.coalesce(1, 0)",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, specialForms, `SELECT coalesce(1, 0), if(true, 1, 2), nullif(1, 2), try(1), "cat".s.coalesce(1, 0);`)

	// Table DEFAULT expressions qualify catalog UDF calls consistently.
	defaultUDF, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "next_id",
				Definition: "CREATE FUNCTION s.next_id() RETURNS bigint RETURN 42",
			}},
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t_def_udf",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint", Default: "next_id()"}},
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, defaultUDF, `"id" bigint DEFAULT "cat"."s".next_id() NOT NULL`)

	// Non-varchar format() calls resolve to colliding schema UDFs, and duplicate
	// view columns fail explicitly.
	formatCollision, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "format",
				Definition: "CREATE FUNCTION s.format(x bigint) RETURNS varchar RETURN CAST(x AS varchar)",
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_fmt",
				Definition: "SELECT format(1), format('%s', 1)",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, formatCollision, `SELECT "cat"."s".format(1), format('%s', 1);`)

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_dup_col",
				Definition: "SELECT 1 AS a, 2 AS b",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "c", Comment: "first"},
					{Name: "C", Comment: "second"},
				},
			}},
		}},
	})
	require.ErrorContains(t, err, "duplicate column")

	// Unqualified view relations resolve from DependencyColumns when their actual
	// schema differs from the view's own schema.
	crossSchemaDep, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name:    "t",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}},
				}},
			},
			{
				Name: "app",
				Views: []*metadatapb.ViewMetadata{{
					Name:       "v",
					Definition: "SELECT id FROM t",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public",
						Table:  "t",
					}},
				}},
			},
		},
	})
	require.NoError(t, err)
	assert.Contains(t, crossSchemaDep, `SELECT id FROM "cat"."public".t;`)

	// Direct GetDatabaseDefinition rejects duplicate function signatures while
	// distinguishing interval overloads.
	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			Functions: []*FunctionMetadata{
				{Name: "dup_f", Definition: "CREATE FUNCTION s.dup_f(x int) RETURNS int RETURN x"},
				{Name: "dup_f", Definition: "CREATE FUNCTION s.dup_f(y integer) RETURNS integer RETURN y + 1"},
			},
		}},
	})
	require.ErrorContains(t, err, "duplicate function")

	directIntervalOK, err := GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			Functions: []*FunctionMetadata{
				{Name: "iv_f", Definition: "CREATE FUNCTION iv_f(x INTERVAL YEAR TO MONTH) RETURNS varchar RETURN 'ym'"},
				{Name: "iv_f", Definition: "CREATE FUNCTION iv_f(x INTERVAL DAY TO SECOND) RETURNS varchar RETURN 'ds'"},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, directIntervalOK, `CREATE FUNCTION "cat"."s".iv_f(x INTERVAL YEAR TO MONTH)`)
	assert.Contains(t, directIntervalOK, `CREATE FUNCTION "cat"."s".iv_f(x INTERVAL DAY TO SECOND)`)

	// Row-pattern MEASURES/DEFINE inside both inline OVER (...) and named WINDOW w AS (...) qualify UDF calls.
	windowPattern, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "score_fn",
				Definition: "CREATE FUNCTION s.score_fn(x bigint) RETURNS bigint RETURN x + 1",
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_win_pat",
				Definition: "SELECT FIRST_VALUE(m) OVER (ORDER BY id MEASURES score_fn(id) AS m ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING PATTERN (A) DEFINE A AS score_fn(id) > 0), LAST_VALUE(n) OVER w FROM t WINDOW w AS (ORDER BY id MEASURES score_fn(id) AS n ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING PATTERN (B) DEFINE B AS score_fn(id) > 1)",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, windowPattern, `OVER (ORDER BY id MEASURES "cat"."s".score_fn(id) AS m`)
	assert.Contains(t, windowPattern, `DEFINE A AS "cat"."s".score_fn(id) > 0`)
	assert.Contains(t, windowPattern, `WINDOW w AS (ORDER BY id MEASURES "cat"."s".score_fn(id) AS n`)
	assert.Contains(t, windowPattern, `DEFINE B AS "cat"."s".score_fn(id) > 1`)

	// Built-in collisions with column arguments resolve column types rather than
	// suppressing UDFs, and GetDatabaseDefinition enforces duplicate relations,
	// MV comment rejection, and view qualification on the direct model.
	absCollision, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "txt", Type: "varchar"},
					{Name: "num", Type: "bigint"},
				},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "abs",
				Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
			}},
			Views: []*metadatapb.ViewMetadata{
				{Name: "v_known_col", Definition: "SELECT abs(txt), abs(num) FROM t"},
				{Name: "v_ext_col", Definition: "SELECT abs(unknown_col) FROM ext_tbl"},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, absCollision, `SELECT "cat"."s".abs(txt), abs(num) FROM "cat"."s".t;`)
	assert.Contains(t, absCollision, `SELECT "cat"."s".abs(unknown_col) FROM "cat"."s".ext_tbl;`)

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			Tables: []*TableMetadata{
				{Name: "dup_t", Columns: []*ColumnMetadata{{Name: "id", Type: "bigint"}}},
				{Name: "DUP_T", Columns: []*ColumnMetadata{{Name: "id", Type: "bigint"}}},
			},
		}},
	})
	require.ErrorContains(t, err, "duplicate table")

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name:              "s",
			MaterializedViews: []*ViewMetadata{{Name: "mv", Definition: "SELECT 1", Comment: "unsupported"}},
		}},
	})
	require.ErrorContains(t, err, "comment cannot be restored")

	directViewQual, err := GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name:              "s",
			Tables:            []*TableMetadata{{Name: "t", Columns: []*ColumnMetadata{{Name: "id", Type: "bigint"}}}},
			Functions:         []*FunctionMetadata{{Name: "f", Definition: "CREATE FUNCTION f() RETURNS bigint RETURN 1"}},
			Views:             []*ViewMetadata{{Name: "v", Definition: "SELECT id, f() FROM t"}},
			MaterializedViews: []*ViewMetadata{{Name: "mv", Definition: "SELECT id, f() FROM t"}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, directViewQual, `SELECT id, "cat"."s".f() FROM "cat"."s".t;`)
	assert.Contains(t, directViewQual, "CREATE MATERIALIZED VIEW IF NOT EXISTS \"s\".\"mv\" AS\nSELECT id, \"cat\".\"s\".f() FROM \"cat\".\"s\".t;")

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{
			{Name: "s", Tables: []*TableMetadata{{Name: "t", Columns: []*ColumnMetadata{{Name: "id", Type: "bigint"}}}}},
			{Name: "S", Tables: []*TableMetadata{{Name: "t", Columns: []*ColumnMetadata{{Name: "name", Type: "varchar"}}}}},
		},
	})
	require.ErrorContains(t, err, "duplicate schema")

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name:    "cat",
		Schemas: []*SchemaMetadata{{Name: "s", Owner: "role_a"}},
	})
	require.ErrorContains(t, err, "owner")

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name:   "s",
			Tables: []*TableMetadata{{Name: "t", Owner: "role_a", Columns: []*ColumnMetadata{{Name: "id", Type: "bigint"}}}},
		}},
	})
	require.ErrorContains(t, err, "owner")

	// Incompatible greatest/least args resolve to schema UDFs, skipped unsupported
	// schema objects are omitted cleanly, and self-referential views fail explicitly.
	greatestAndSkip, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:           "s",
			EnumTypes:      []*metadatapb.EnumTypeMetadata{{Name: "e", Values: []string{"a"}, SkipDump: true}},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{Name: "c", SkipDump: true}},
			Collations:     []*metadatapb.CollationMetadata{{Name: "col", SkipDump: true}},
			RangeTypes:     []*metadatapb.RangeTypeMetadata{{Name: "r", SkipDump: true}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "greatest",
				Definition: "CREATE FUNCTION s.greatest(a boolean, b bigint) RETURNS bigint RETURN b",
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_gr",
				Definition: "SELECT greatest(true, 1), greatest(1, 2)",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, greatestAndSkip, `SELECT "cat"."s".greatest(true, 1), greatest(1, 2);`)

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:  "s",
			Views: []*metadatapb.ViewMetadata{{Name: "self_v", Definition: "SELECT * FROM s.self_v"}},
		}},
	})
	require.ErrorContains(t, err, "cyclic view")

	// CTE output types propagate to outer SELECT built-in matching, and direct-model
	// table columns with empty types fail explicitly.
	cteBuiltin, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "abs",
				Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_cte",
				Definition: "WITH x(n) AS (SELECT 1) SELECT abs(n) FROM x",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, cteBuiltin, `WITH x(n) AS (SELECT 1) SELECT abs(n) FROM x;`)

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name:   "s",
			Tables: []*TableMetadata{{Name: "t_empty_coltype", Columns: []*ColumnMetadata{{Name: "c", Type: " "}}}},
		}},
	})
	require.ErrorContains(t, err, "no data type")

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			Tables: []*TableMetadata{{
				Name: "t_dup_col",
				Columns: []*ColumnMetadata{
					{Name: "id", Type: "bigint"},
					{Name: "ID", Type: "varchar"},
				},
			}},
		}},
	})
	require.ErrorContains(t, err, "duplicate column")

	// Derived-table output types propagate to outer SELECT built-in matching, skipped
	// event triggers are ignored while non-skipped ones fail, and skipped UDFs remain
	// in the resolution catalog so callers still qualify calls to them.
	derivedAndSkipped, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name:          "cat",
		EventTriggers: []*metadatapb.EventTriggerMetadata{{Name: "skip_evt", SkipDump: true}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "abs",
					Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
				},
				{
					Name:       "ext_fn",
					Definition: "CREATE FUNCTION cat.s.ext_fn(x bigint) RETURNS bigint RETURN x + 1",
					SkipDump:   true,
				},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_derived_and_ext",
				Definition: "SELECT abs(n), abs(x.n), ext_fn(n) FROM (SELECT 1 AS n) x",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, derivedAndSkipped, `SELECT abs(n), abs(x.n), "cat"."s".ext_fn(n) FROM (SELECT 1 AS n) x;`)
	assert.Contains(t, derivedAndSkipped, `CREATE FUNCTION "cat".s.abs`)
	assert.NotContains(t, derivedAndSkipped, "ext_fn(x bigint)")
	assert.NotContains(t, derivedAndSkipped, "CREATE FUNCTION \"cat\".s.ext_fn")
	assert.NotContains(t, derivedAndSkipped, "CREATE FUNCTION \"cat\".\"s\".ext_fn")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name:          "cat",
		EventTriggers: []*metadatapb.EventTriggerMetadata{{Name: "evt"}},
	})
	require.ErrorContains(t, err, "unsupported")

	// Direct-model views with duplicate columns and table defaults containing
	// subqueries fail explicitly rather than emitting invalid DDL.
	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			Views: []*ViewMetadata{{
				Name:       "v_dup_col",
				Definition: "SELECT 1 AS a, 2 AS b",
				Columns: []*ColumnMetadata{
					{Name: "id", Comment: "first"},
					{Name: "ID", Comment: "second"},
				},
			}},
		}},
	})
	require.ErrorContains(t, err, "duplicate column")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t_subquery_def",
				Columns: []*metadatapb.ColumnMetadata{{
					Name:    "id",
					Type:    "bigint",
					Default: "(SELECT 1)",
				}},
			}},
		}},
	})
	require.ErrorContains(t, err, "subquery")

	// Nested column references in table defaults fail explicitly, and skipped
	// tables still supply column types for view built-in vs UDF resolution.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t_nested_colref_def",
				Columns: []*metadatapb.ColumnMetadata{{
					Name:    "id",
					Type:    "bigint",
					Default: "abs(other_col)",
				}},
			}},
		}},
	})
	require.ErrorContains(t, err, "column")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t_nested_cast_def",
				Columns: []*metadatapb.ColumnMetadata{{
					Name:    "id",
					Type:    "bigint",
					Default: "CAST(other_col AS bigint)",
				}},
			}},
		}},
	})
	require.ErrorContains(t, err, "column")

	skipTableRes, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "abs",
				Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
			}},
			Tables: []*metadatapb.TableMetadata{{
				Name:     "ext_t",
				SkipDump: true,
				Columns:  []*metadatapb.ColumnMetadata{{Name: "n", Type: "bigint"}},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_skip_tbl",
				Definition: "SELECT abs(n) FROM ext_t",
			}},
		}},
	})
	require.NoError(t, err)
	assert.NotContains(t, skipTableRes, `CREATE TABLE IF NOT EXISTS "s"."ext_t"`)
	assert.Contains(t, skipTableRes, `SELECT abs(n) FROM "cat"."s".ext_t;`)

	// Unparseable views still qualify explicit cross-schema dependencies, and
	// direct-model views reject column defaults.
	unparseableDef := "SELECT id FROM t AS OF TIMESTAMP '2024-01-01 00:00:00'"
	_, parseErr := trinoparser.Parse(unparseableDef)
	require.Error(t, parseErr)
	downgraded, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:   "shared",
				Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}}},
			},
			{
				Name: "app",
				Views: []*metadatapb.ViewMetadata{{
					Name:              "v_downgraded",
					Definition:        unparseableDef,
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "shared", Table: "t"}},
				}},
			},
		},
	})
	require.NoError(t, err)
	assert.Contains(t, downgraded, `FROM "cat"."shared".t AS OF TIMESTAMP`)

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			Views: []*ViewMetadata{{
				Name:       "v_direct_col_def",
				Definition: "SELECT 1 AS id",
				Columns:    []*ColumnMetadata{{Name: "id", Default: "1"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "default cannot be restored")

	// View columns (including skipped views) supply types for downstream call resolution,
	// and definition-less skipped functions do not fabricate a return type that silently
	// mis-binds ambiguous downstream overload resolution.
	viewOnView, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "abs",
				Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
			}},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "src_v",
					Definition: "SELECT 1 AS n",
					Columns:    []*metadatapb.ColumnMetadata{{Name: "n", Type: "bigint"}},
				},
				{
					Name:       "skip_v",
					Definition: "SELECT 2 AS m",
					Columns:    []*metadatapb.ColumnMetadata{{Name: "m", Type: "bigint"}},
					SkipDump:   true,
				},
				{
					Name:       "dst_v",
					Definition: "SELECT abs(n), abs(m) FROM src_v, skip_v",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, viewOnView, `SELECT abs(n), abs(m) FROM "cat"."s".src_v, "cat"."s".skip_v;`)
	assert.NotContains(t, viewOnView, `CREATE OR REPLACE VIEW "s"."skip_v"`)

	// Unambiguous call to definition-less skipped function succeeds and qualifies;
	// ambiguous call abs(ext_fn()) against built-in abs vs s.abs(varchar) fails explicitly.
	extOnly, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:      "ext_fn",
				Signature: "()",
				SkipDump:  true,
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_ext",
				Definition: "SELECT ext_fn()",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, extOnly, `SELECT "cat"."s".ext_fn();`)

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "abs",
					Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
				},
				{
					Name:      "ext_fn",
					Signature: "()",
					SkipDump:  true,
				},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_ambig",
				Definition: "SELECT abs(ext_fn())",
			}},
		}},
	})
	require.ErrorContains(t, err, "unambiguously")

	// Lexical fallback preserves CTE references, nested built-in return types
	// inform outer overload matching, and base-table column aliases remap types.
	cteDowngraded, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:   "shared",
				Tables: []*metadatapb.TableMetadata{{Name: "x", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}}},
			},
			{
				Name: "app",
				Views: []*metadatapb.ViewMetadata{{
					Name:              "v_cte_downgraded",
					Definition:        "WITH x AS (SELECT id FROM shared.x) SELECT id FROM x AS OF TIMESTAMP '2024-01-01 00:00:00'",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "shared", Table: "x"}},
				}},
			},
		},
	})
	require.NoError(t, err)
	assert.Contains(t, cteDowngraded, "WITH x AS (SELECT id FROM shared.x) SELECT id FROM x AS OF TIMESTAMP")

	nestedAndAlias, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{
				{
					Name:    "t",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}, {Name: "n", Type: "bigint"}},
				},
				{
					Name:    "u",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}, {Name: "m", Type: "bigint"}},
				},
			},
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "abs",
					Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
				},
				{
					Name:       "sqrt",
					Definition: "CREATE FUNCTION s.sqrt(x varchar) RETURNS varchar RETURN x",
				},
			},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "v_nested_alias",
					Definition: "SELECT sqrt(abs(1)), abs(x) FROM t AS q(k, x)",
				},
				{
					Name:       "v_join_alias",
					Definition: "SELECT abs(j.n), abs(j.m) FROM (t JOIN u ON t.id = u.id) j",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, nestedAndAlias, `SELECT sqrt(abs(1)), abs(x) FROM "cat"."s".t AS q(k, x);`)
	assert.Contains(t, nestedAndAlias, `SELECT abs(j.n), abs(j.m) FROM ("cat"."s".t JOIN "cat"."s".u ON t.id = u.id) j;`)

	// SkipDump triggers on tables, views, and materialized views are ignored,
	// and direct-model materialized view column comments fail explicitly.
	skipTrigs, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:     "t_trig",
				Columns:  []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}},
				Triggers: []*metadatapb.TriggerMetadata{{Name: "tr_t", SkipDump: true}},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_trig",
				Definition: "SELECT id FROM t_trig",
				Triggers:   []*metadatapb.TriggerMetadata{{Name: "tr_v", SkipDump: true}},
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_trig",
				Definition: "SELECT id FROM t_trig",
				Triggers:   []*metadatapb.TriggerMetadata{{Name: "tr_mv", SkipDump: true}},
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, skipTrigs, `CREATE TABLE IF NOT EXISTS "s"."t_trig"`)

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name: "s",
			MaterializedViews: []*ViewMetadata{{
				Name:       "mv_col_cmt",
				Definition: "SELECT 1 AS id",
				Columns:    []*ColumnMetadata{{Name: "id", Comment: "unsupported"}},
			}},
		}},
	})
	require.ErrorContains(t, err, "comment cannot be restored")

	// Wildcard derived projections propagate column types, and empty schema names fail explicitly.
	wildcardProj, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "source",
				Columns: []*metadatapb.ColumnMetadata{{Name: "n", Type: "bigint"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "abs",
				Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
			}},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "v_star",
					Definition: "SELECT abs(q.n) FROM (SELECT * FROM source) q",
				},
				{
					Name:       "v_star_from",
					Definition: "SELECT abs(q.n) FROM (SELECT src.* FROM source src) q",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, wildcardProj, `SELECT abs(q.n) FROM (SELECT * FROM "cat"."s".source) q;`)
	assert.Contains(t, wildcardProj, `SELECT abs(q.n) FROM (SELECT src.* FROM "cat"."s".source src) q;`)

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name:    "cat",
		Schemas: []*SchemaMetadata{{Name: ""}},
	})
	require.ErrorContains(t, err, "empty name")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name:    "cat",
		Schemas: []*metadatapb.SchemaMetadata{{Name: ""}},
	})
	require.ErrorContains(t, err, "empty name")

	// Lexical fallback function scanning excludes return type modifiers such as decimal(10,2).
	fallbackRetType, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "decimal",
				Signature:  "(integer, integer)",
				Definition: "CREATE FUNCTION s.decimal(a integer, b integer) RETURNS decimal(10,2) COMMENT 'unparseable LANGUAGE clause' LANGUAGE PYTHON RETURN 0",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, fallbackRetType, "RETURNS decimal(10,2)")

	// VALUES projection types propagate to enclosing aliases, and downgraded
	// function definitions qualify unambiguous UDF calls in their bodies.
	valuesAndFallback, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "abs",
					Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
				},
				{
					Name:       "callee",
					Definition: "CREATE FUNCTION s.callee() RETURNS bigint RETURN 1",
				},
				{
					Name:       "caller_fb",
					Signature:  "()",
					Definition: "CREATE FUNCTION s.caller_fb() RETURNS bigint COMMENT 'unparseable LANGUAGE clause' LANGUAGE PYTHON RETURN callee()",
				},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_vals",
				Definition: "SELECT abs(q.n) FROM (VALUES BIGINT '1') AS q(n)",
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, valuesAndFallback, `SELECT abs(q.n) FROM (VALUES BIGINT '1') AS q(n);`)
	assert.Contains(t, valuesAndFallback, `RETURN "cat"."s".callee();`)

	// Same-schema explicit dependencies are qualified even when view parsing downgrades.
	downgradedSameSchema, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "s",
			Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}}},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "v_same_down",
				Definition:        "SELECT id FROM t AS OF TIMESTAMP '2024-01-01 00:00:00'",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "s", Table: "t"}},
			}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, downgradedSameSchema, `FROM "cat"."s".t AS OF TIMESTAMP`)

	// Empty relation names fail explicitly on both public entry points.
	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name:   "s",
			Tables: []*TableMetadata{{Name: "", Columns: []*ColumnMetadata{{Name: "id", Type: "bigint"}}}},
		}},
	})
	require.ErrorContains(t, err, "empty name")

	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:  "s",
			Views: []*metadatapb.ViewMetadata{{Name: "", Definition: "SELECT 1"}},
		}},
	})
	require.ErrorContains(t, err, "empty name")

	// Empty column names fail explicitly on both public entry points.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "s",
			Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "", Type: "bigint"}}}},
		}},
	})
	require.ErrorContains(t, err, "column with empty name")

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name:   "s",
			Tables: []*TableMetadata{{Name: "t", Columns: []*ColumnMetadata{{Name: "", Type: "bigint"}}}},
		}},
	})
	require.ErrorContains(t, err, "column with empty name")

	_, err = GetDatabaseDefinition(&DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*SchemaMetadata{{
			Name:  "s",
			Views: []*ViewMetadata{{Name: "v", Definition: "SELECT 1 AS x", Columns: []*ColumnMetadata{{Name: ""}}}},
		}},
	})
	require.ErrorContains(t, err, "column with empty name")

	// MATCH_RECOGNIZE returns its partition columns and typed measures to the
	// enclosing query under its alias, so a numeric measure keeps Trino's
	// built-in abs instead of the same-schema abs(varchar) UDF.
	matchRecognize, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}, {Name: "k", Type: "bigint"}, {Name: "n", Type: "bigint"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "abs",
				Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x",
			}},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "v_mr_measure",
					Definition: "SELECT abs(m.n), abs(m.k) FROM t MATCH_RECOGNIZE (PARTITION BY k ORDER BY id MEASURES A.n AS n PATTERN (A) DEFINE A AS true) AS m",
				},
				{
					Name:       "v_mr_all_rows",
					Definition: "SELECT abs(r.id), abs(r.q) FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES B.n AS m ALL ROWS PER MATCH PATTERN (A B) DEFINE B AS B.n > 0) AS r(id, m2, k2, q)",
				},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, matchRecognize, `SELECT abs(m.n), abs(m.k) FROM "cat"."s".t MATCH_RECOGNIZE`)
	assert.Contains(t, matchRecognize, `SELECT abs(r.id), abs(r.q) FROM "cat"."s".t MATCH_RECOGNIZE`)

	// JSON_TABLE columns, ORDER BY output aliases, and array concatenation
	// keep their types, so calls stay on Trino's built-ins instead of the
	// same-schema varchar UDFs.
	typedBuiltins, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Columns: []*metadatapb.ColumnMetadata{{Name: "n", Type: "bigint"}, {Name: "j", Type: "varchar"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "abs", Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x"},
				{Name: "cardinality", Definition: "CREATE FUNCTION s.cardinality(x varchar) RETURNS bigint RETURN 1"},
			},
			Views: []*metadatapb.ViewMetadata{
				{Name: "v_json_table", Definition: "SELECT abs(jt.n), abs(jt.o) FROM t, JSON_TABLE(t.j, 'lax $' COLUMNS (o FOR ORDINALITY, n BIGINT PATH 'lax $.n')) AS jt"},
				{Name: "v_order_alias", Definition: "SELECT n AS x FROM t ORDER BY abs(x)"},
				{Name: "v_order_source", Definition: "SELECT n AS x FROM t ORDER BY abs(n), abs(t.n)"},
				{Name: "v_array_concat", Definition: "SELECT cardinality(ARRAY[1] || ARRAY[2]) AS c"},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, typedBuiltins, `SELECT abs(jt.n), abs(jt.o) FROM "cat"."s".t, JSON_TABLE(`)
	assert.Contains(t, typedBuiltins, `ORDER BY abs(x);`)
	assert.Contains(t, typedBuiltins, `ORDER BY abs(n), abs(t.n);`)
	assert.Contains(t, typedBuiltins, `SELECT cardinality(ARRAY[1] || ARRAY[2]) AS c;`)

	// JOIN ... USING exposes one coalesced column and UNNEST returns its
	// element types, so both keep Trino's numeric abs.
	usingUnnest, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Tables: []*metadatapb.TableMetadata{
				{Name: "a", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}},
				{Name: "b", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}},
			},
			Functions: []*metadatapb.FunctionMetadata{{Name: "abs", Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x"}},
			Views: []*metadatapb.ViewMetadata{
				{Name: "v_using", Definition: "SELECT abs(id) FROM a JOIN b USING (id)"},
				{Name: "v_unnest", Definition: "SELECT abs(q.n), abs(q.o) FROM UNNEST(ARRAY[BIGINT '1']) WITH ORDINALITY AS q(n, o)"},
			},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, usingUnnest, `SELECT abs(id) FROM "cat"."s".a JOIN "cat"."s".b USING (id);`)
	assert.Contains(t, usingUnnest, `SELECT abs(q.n), abs(q.o) FROM UNNEST(`)

	// A definition the parser cannot read still has its parameter list
	// checked against the metadata signature.
	_, err = DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "fb",
				Signature:  "(bigint)",
				Definition: "CREATE FUNCTION s.fb(x varchar) RETURNS bigint NOT_A_CLAUSE RETURN 0",
			}},
		}},
	})
	require.ErrorContains(t, err, "does not match signature")

	// Lambda parameters take the element type of the array a higher-order
	// built-in passes them, so abs stays Trino's numeric built-in.
	lambdaTyped, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "s",
			Functions: []*metadatapb.FunctionMetadata{{Name: "abs", Definition: "CREATE FUNCTION s.abs(x varchar) RETURNS varchar RETURN x"}},
			Views:     []*metadatapb.ViewMetadata{{Name: "v_lambda", Definition: "SELECT transform(ARRAY[1, 2], x -> abs(x)) AS t"}},
		}},
	})
	require.NoError(t, err)
	assert.Contains(t, lambdaTyped, `SELECT transform(ARRAY[1, 2], x -> abs(x)) AS t;`)

	// An unqualified relation whose recorded dependency is in another schema
	// is ordered against that schema, not as a self-reference.
	unqualifiedCrossDep, err := DatabaseDefinitionFromProto(&metadatapb.DatabaseSchemaMetadata{
		Name: "cat",
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "shared", Tables: []*metadatapb.TableMetadata{{Name: "v", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}}}}},
			{Name: "app", Views: []*metadatapb.ViewMetadata{{
				Name:              "v",
				Definition:        "SELECT id FROM v",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "shared", Table: "v"}},
			}}},
		},
	})
	require.NoError(t, err)
	assert.Contains(t, unqualifiedCrossDep, `SELECT id FROM "cat"."shared".v;`)
}
