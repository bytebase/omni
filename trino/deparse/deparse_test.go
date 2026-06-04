package deparse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetDatabaseDefinition is the differential-vs-legacy gate: the omni output
// must be byte-identical to the legacy bytebase plugin/schema/trino
// GetDatabaseDefinition output. The first three cases are lifted verbatim
// (inputs and expected strings) from
// bytebase/backend/plugin/schema/trino/get_database_definition_test.go so any
// drift from the legacy format fails here. The remaining cases extend coverage
// to inputs the legacy suite never exercised (multiple tables, multiple
// schemas, ordering, and nil/empty handling).
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
			expected: `CREATE TABLE IF NOT EXISTS "testschema"."testtable" (
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
			expected: `CREATE TABLE IF NOT EXISTS "testschema"."empty_table" (

);

`,
		},
		{
			// Verbatim from the legacy suite: special characters in identifiers
			// are emitted as-is inside the double quotes (no escaping), exactly
			// like the legacy source.
			name: "special characters in identifiers (legacy golden)",
			metadata: oneTable("test-catalog", "test_schema", &TableMetadata{
				Name: "test.table",
				Columns: []*ColumnMetadata{
					{Name: "id-field", Type: "bigint", Nullable: false},
				},
			}),
			expected: `CREATE TABLE IF NOT EXISTS "test_schema"."test.table" (
    "id-field" bigint NOT NULL
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
			expected: `CREATE TABLE IF NOT EXISTS "s"."first" (
    "a" integer NOT NULL
);

CREATE TABLE IF NOT EXISTS "s"."second" (
    "b" varchar
);

`,
		},
		{
			// Tables from multiple schemas render schema-by-schema in snapshot
			// order, qualified by their own schema name.
			name: "multiple schemas",
			metadata: &DatabaseSchemaMetadata{
				Name: "cat",
				Schemas: []*SchemaMetadata{
					{Name: "s1", Tables: []*TableMetadata{{Name: "t1", Columns: []*ColumnMetadata{{Name: "c1", Type: "bigint", Nullable: false}}}}},
					{Name: "s2", Tables: []*TableMetadata{{Name: "t2", Columns: []*ColumnMetadata{{Name: "c2", Type: "double", Nullable: true}}}}},
				},
			},
			expected: `CREATE TABLE IF NOT EXISTS "s1"."t1" (
    "c1" bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "s2"."t2" (
    "c2" double
);

`,
		},
		{
			// A schema with no tables contributes nothing.
			name:     "schema with no tables",
			metadata: &DatabaseSchemaMetadata{Name: "cat", Schemas: []*SchemaMetadata{{Name: "empty"}}},
			expected: "",
		},
		{
			// No schemas at all yields the empty string.
			name:     "no schemas",
			metadata: &DatabaseSchemaMetadata{Name: "cat"},
			expected: "",
		},
		{
			// A nil metadata is handled gracefully (the legacy source would
			// panic on nil; the omni version returns empty instead).
			name:     "nil metadata",
			metadata: nil,
			expected: "",
		},
		{
			// nil entries in the slices are skipped rather than panicking.
			name: "nil schema, table, and column entries are skipped",
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
					},
				},
			},
			expected: `CREATE TABLE IF NOT EXISTS "s"."t" (
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

// TestWriteCreateTable pins the single-statement renderer, mirroring the legacy
// TestWriteCreateTable (same input, same expected — no trailing blank line, the
// caller adds that).
func TestWriteCreateTable(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		table    *TableMetadata
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
