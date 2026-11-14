package catalog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/parser"
)

func TestMultiFileDefinition_WithComments(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:    "app_schema",
				Comment: "Application schema",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "products",
						Comment: "Product catalog",
						Columns: []*metadatapb.ColumnMetadata{
							{
								Name:     "id",
								Type:     "INTEGER",
								Nullable: false,
								Comment:  "Product ID",
							},
							{
								Name:     "name",
								Type:     "TEXT",
								Nullable: false,
								Comment:  "Product name",
							},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{
								Name:        "idx_products_name",
								Expressions: []string{"name"},
								Unique:      false,
								Primary:     false,
								Comment:     "Index for product search",
							},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:       "active_products",
						Definition: "SELECT id, name FROM app_schema.products WHERE active = true",
						Comment:    "View of active products",
					},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{
						Name:       "product_summary_mv",
						Definition: "SELECT id, name FROM app_schema.products",
						Comment:    "Materialized view of product summary",
					},
				},
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:       "count_products",
						Signature:  "count_products()",
						Definition: "CREATE FUNCTION app_schema.count_products() RETURNS INTEGER AS $$ BEGIN RETURN (SELECT COUNT(*) FROM app_schema.products); END; $$ LANGUAGE plpgsql",
						Comment:    "Returns product count",
					},
				},
				Sequences: []*metadatapb.SequenceMetadata{
					{
						Name:      "order_seq",
						DataType:  "bigint",
						Start:     "1",
						Increment: "1",
						MinValue:  "1",
						MaxValue:  "9223372036854775807",
						Cycle:     false,
						CacheSize: "1",
						Comment:   "Sequence for orders",
					},
				},
			},
		},
	}

	result, err := MultiFileDefinition(metadata)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Build a map for easier testing
	fileMap := make(map[string]string)
	for _, file := range result {
		fileMap[file.Name] = file.Content
	}

	// Verify schema-level file with schema comment
	schemaFile, ok := fileMap["schemas/app_schema/schema.sql"]
	require.True(t, ok, "schema.sql file should exist for non-public schema")
	assert.Contains(t, schemaFile, `CREATE SCHEMA IF NOT EXISTS "app_schema";`)
	assert.Contains(t, schemaFile, `COMMENT ON SCHEMA "app_schema" IS 'Application schema';`)

	// Verify consolidated sequences file with comment (independent sequences go in sequences.sql)
	sequenceFile, ok := fileMap["schemas/app_schema/sequences.sql"]
	require.True(t, ok, "sequences.sql file should exist for independent sequences")
	assert.Contains(t, sequenceFile, `CREATE SEQUENCE "app_schema"."order_seq"`)
	assert.Contains(t, sequenceFile, `COMMENT ON SEQUENCE "app_schema"."order_seq" IS 'Sequence for orders';`)

	// Verify table file with comments
	tableFile, ok := fileMap["schemas/app_schema/tables/products.sql"]
	require.True(t, ok, "table file should exist")
	assert.Contains(t, tableFile, `CREATE TABLE "app_schema"."products"`)
	assert.Contains(t, tableFile, `COMMENT ON TABLE "app_schema"."products" IS 'Product catalog';`)
	assert.Contains(t, tableFile, `COMMENT ON COLUMN "app_schema"."products"."id" IS 'Product ID';`)
	assert.Contains(t, tableFile, `COMMENT ON COLUMN "app_schema"."products"."name" IS 'Product name';`)
	assert.Contains(t, tableFile, `COMMENT ON INDEX "app_schema"."idx_products_name" IS 'Index for product search';`)

	// Verify view file with comment
	viewFile, ok := fileMap["schemas/app_schema/views/active_products.sql"]
	require.True(t, ok, "view file should exist")
	assert.Contains(t, viewFile, `CREATE VIEW "app_schema"."active_products"`)
	assert.Contains(t, viewFile, `COMMENT ON VIEW "app_schema"."active_products" IS 'View of active products';`)

	// Verify materialized view file with comment
	materializedViewFile, ok := fileMap["schemas/app_schema/materialized_views/product_summary_mv.sql"]
	require.True(t, ok, "materialized view file should exist")
	assert.Contains(t, materializedViewFile, `CREATE MATERIALIZED VIEW "app_schema"."product_summary_mv"`)
	assert.Contains(t, materializedViewFile, `COMMENT ON MATERIALIZED VIEW "app_schema"."product_summary_mv" IS 'Materialized view of product summary';`)

	// Verify function file with comment
	functionFile, ok := fileMap["schemas/app_schema/functions/count_products.sql"]
	require.True(t, ok, "function file should exist")
	assert.Contains(t, functionFile, `CREATE FUNCTION app_schema.count_products()`)
	assert.Contains(t, functionFile, `COMMENT ON FUNCTION "app_schema".count_products() IS 'Returns product count';`)

	// Validate that each file's SQL can be parsed
	for fileName, content := range fileMap {
		_, err := parser.Parse(content)
		require.NoError(t, err, "SQL in file %s should be parseable by PostgreSQL parser", fileName)
	}
}

func TestMultiFileDefinition_SchemaFileMakesCombinedSDLLoadable(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "metric_helpers",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "metrics",
						Columns: []*metadatapb.ColumnMetadata{
							{
								Name:     "id",
								Type:     "INTEGER",
								Nullable: false,
							},
						},
					},
				},
			},
		},
	}

	result, err := MultiFileDefinition(metadata)
	require.NoError(t, err)

	fileMap := make(map[string]string)
	var combined strings.Builder
	for _, file := range result {
		fileMap[file.Name] = file.Content
		combined.WriteString(file.Content)
		combined.WriteString("\n")
	}

	schemaFile, ok := fileMap["schemas/metric_helpers/schema.sql"]
	assert.True(t, ok, "schema.sql file should exist for non-public schema")
	assert.Contains(t, schemaFile, `CREATE SCHEMA IF NOT EXISTS "metric_helpers";`)

	_, err = LoadSDL(combined.String())
	require.NoError(t, err)
}

func TestMultiFileDefinition_SerialColumnWithSequence(t *testing.T) {
	// This test verifies that multi-file format correctly handles serial columns and their sequences
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Sequences: []*metadatapb.SequenceMetadata{
					{
						Name:        "users_id_seq",
						DataType:    "integer",
						Start:       "1",
						Increment:   "1",
						MinValue:    "1",
						MaxValue:    "2147483647",
						Cycle:       false,
						OwnerTable:  "users",
						OwnerColumn: "id",
						Comment:     "Sequence for users id - should NOT appear as separate file",
					},
					{
						Name:       "independent_seq",
						DataType:   "bigint",
						Start:      "1",
						Increment:  "1",
						MinValue:   "1",
						MaxValue:   "9223372036854775807",
						Cycle:      false,
						OwnerTable: "", // Independent sequence - not owned by any column
						Comment:    "Independent sequence - should appear as separate file",
					},
					{
						Name:       "custom_seq",
						DataType:   "bigint",
						Start:      "100",
						Increment:  "5",
						MinValue:   "100",
						MaxValue:   "9223372036854775807",
						Cycle:      false,
						OwnerTable: "", // Not owned by a column (independent sequence)
						Comment:    "Custom sequence for orders - should appear as separate file",
					},
				},
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "users",
						Columns: []*metadatapb.ColumnMetadata{
							{
								Name:     "id",
								Type:     "integer",
								Nullable: false,
								Default:  "nextval('users_id_seq'::regclass)",
							},
							{
								Name:     "name",
								Type:     "text",
								Nullable: false,
							},
						},
					},
					{
						Name: "orders",
						Columns: []*metadatapb.ColumnMetadata{
							{
								Name:     "id",
								Type:     "integer",
								Nullable: false,
							},
							{
								Name:     "order_number",
								Type:     "bigint",
								Nullable: false,
								Default:  "nextval('custom_seq'::regclass)",
								Comment:  "Order number using custom sequence",
							},
						},
					},
				},
			},
		},
	}

	result, err := MultiFileDefinition(metadata)
	require.NoError(t, err)

	// Build a map of files for easy lookup
	fileMap := make(map[string]string)
	for _, file := range result {
		fileMap[file.Name] = file.Content
	}

	// Verify users table file exists and contains serial column
	usersFile, ok := fileMap["schemas/public/tables/users.sql"]
	require.True(t, ok, "users table file should exist")
	require.Contains(t, usersFile, `"id" serial`, "users table should use serial type")
	require.Contains(t, usersFile, `"name" text NOT NULL`)

	// Verify orders table file exists
	ordersFile, ok := fileMap["schemas/public/tables/orders.sql"]
	require.True(t, ok, "orders table file should exist")
	require.Contains(t, ordersFile, `"id" integer NOT NULL`)
	require.Contains(t, ordersFile, `"order_number" bigint DEFAULT nextval('custom_seq'::regclass) NOT NULL`,
		"order_number should use custom_seq via DEFAULT nextval()")

	// Verify the consolidated sequences file exists with both independent sequences
	sequencesFile, ok := fileMap["schemas/public/sequences.sql"]
	require.True(t, ok, "independent_seq file should exist")
	require.Contains(t, sequencesFile, `CREATE SEQUENCE "public"."independent_seq"`)
	require.Contains(t, sequencesFile, `COMMENT ON SEQUENCE "public"."independent_seq" IS 'Independent sequence - should appear as separate file'`)
	require.Contains(t, sequencesFile, `CREATE SEQUENCE "public"."custom_seq"`)
	require.Contains(t, sequencesFile, `COMMENT ON SEQUENCE "public"."custom_seq" IS 'Custom sequence for orders - should appear as separate file'`)
	require.NotContains(t, sequencesFile, `ALTER SEQUENCE`, "Independent sequences should not have ALTER SEQUENCE OWNED BY")

	// Verify individual sequence files do NOT exist (they should be in the consolidated file)
	_, ok = fileMap["schemas/public/sequences/independent_seq.sql"]
	require.False(t, ok, "independent_seq should NOT have individual file (should be in sequences.sql)")
	_, ok = fileMap["schemas/public/sequences/custom_seq.sql"]
	require.False(t, ok, "custom_seq should NOT have individual file (should be in sequences.sql)")

	// Verify users_id_seq does NOT appear in the sequences file (it's a serial sequence)
	require.NotContains(t, sequencesFile, "users_id_seq", "users_id_seq should NOT appear in sequences.sql because it's owned by a serial column")

	// Total files should be: 2 tables + 1 sequences file = 3 files
	require.Equal(t, 3, len(result), "Should have exactly 3 files (2 tables + 1 consolidated sequences file)")
}

func TestMultiFileDefinition_FunctionAndProcedureSeparation(t *testing.T) {
	// This test verifies that functions and procedures are separated into different folders
	// in multi-file mode
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name: "get_user_count",
						Definition: `CREATE FUNCTION "public"."get_user_count"() RETURNS integer
    LANGUAGE sql
    AS $$
    SELECT COUNT(*)::integer FROM users;
$$`,
					},
					{
						Name: "get_user_by_id",
						Definition: `CREATE FUNCTION "public"."get_user_by_id"(user_id integer) RETURNS TABLE(id integer, name character varying)
    LANGUAGE sql
    AS $$
    SELECT u.id, u.name FROM users u WHERE u.id = user_id;
$$`,
					},
					{
						Name: "update_user_name",
						Definition: `CREATE PROCEDURE "public"."update_user_name"(IN user_id integer, IN new_name character varying)
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE users SET name = new_name WHERE id = user_id;
END;
$$`,
					},
					{
						Name: "delete_old_users",
						Definition: `CREATE PROCEDURE "public"."delete_old_users"(IN days_old integer)
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM users WHERE created_at < NOW() - INTERVAL '1 day' * days_old;
END;
$$`,
					},
				},
			},
		},
	}

	result, err := MultiFileDefinition(metadata)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Build a map of files for easy lookup
	fileMap := make(map[string]string)
	for _, file := range result {
		fileMap[file.Name] = file.Content
	}

	// Verify functions are in functions folder
	getUserCountFile, ok := fileMap["schemas/public/functions/get_user_count.sql"]
	require.True(t, ok, "get_user_count function should be in functions folder")
	assert.Contains(t, getUserCountFile, `CREATE FUNCTION "public"."get_user_count"()`)

	getUserByIDFile, ok := fileMap["schemas/public/functions/get_user_by_id.sql"]
	require.True(t, ok, "get_user_by_id function should be in functions folder")
	assert.Contains(t, getUserByIDFile, `CREATE FUNCTION "public"."get_user_by_id"`)

	// Verify procedures are in procedures folder (not functions folder)
	updateUserNameFile, ok := fileMap["schemas/public/procedures/update_user_name.sql"]
	require.True(t, ok, "update_user_name procedure should be in procedures folder")
	assert.Contains(t, updateUserNameFile, `CREATE PROCEDURE "public"."update_user_name"`)

	deleteOldUsersFile, ok := fileMap["schemas/public/procedures/delete_old_users.sql"]
	require.True(t, ok, "delete_old_users procedure should be in procedures folder")
	assert.Contains(t, deleteOldUsersFile, `CREATE PROCEDURE "public"."delete_old_users"`)

	// Verify procedures are NOT in functions folder
	_, ok = fileMap["schemas/public/functions/update_user_name.sql"]
	require.False(t, ok, "update_user_name should NOT be in functions folder")

	_, ok = fileMap["schemas/public/functions/delete_old_users.sql"]
	require.False(t, ok, "delete_old_users should NOT be in functions folder")

	// Total files should be: 2 functions + 2 procedures = 4 files
	require.Equal(t, 4, len(result), "Should have exactly 4 files (2 functions + 2 procedures)")

	// Validate that each file's SQL can be parsed
	for fileName, content := range fileMap {
		_, err := parser.Parse(content)
		require.NoError(t, err, "SQL in file %s should be parseable by PostgreSQL parser", fileName)
	}
}

// TestDatabaseDefinition_TableWithConstraintsButNoColumns reproduces a
// dump replay failure observed on a real workspace export: the sync user
// lacked column privileges on some tables, so information_schema.columns
// returned nothing while pg_catalog still exposed the tables' constraints.
// The dump then emitted `CREATE TABLE "s1"."t1" ();` followed by
// `ALTER TABLE ... ADD CONSTRAINT ... PRIMARY KEY (c1, c2);`, which
// PostgreSQL rejects because the columns do not exist. The dump must skip
// column-dependent DDL for such tables instead of producing non-replayable
// statements.
func TestDatabaseDefinition_TableWithConstraintsButNoColumns(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "s1",
			Sequences: []*metadatapb.SequenceMetadata{
				{
					// Owned by a column missing from t1's metadata: the sequence
					// itself must dump, its OWNED BY clause must not.
					Name:        "t1_c1_seq",
					DataType:    "bigint",
					Start:       "1",
					Increment:   "1",
					MinValue:    "1",
					MaxValue:    "9223372036854775807",
					OwnerTable:  "t1",
					OwnerColumn: "c1",
				},
				{
					// Owned by a column of t4, a zero-column table with no
					// broken-metadata signal (no indexes/FKs/partitions):
					// OWNED BY must still be skipped — a sequence cannot be
					// owned by a column of a zero-column table.
					Name:        "t4_id_seq",
					DataType:    "bigint",
					Start:       "1",
					Increment:   "1",
					MinValue:    "1",
					MaxValue:    "9223372036854775807",
					OwnerTable:  "t4",
					OwnerColumn: "id",
				},
			},
			Tables: []*metadatapb.TableMetadata{
				{
					// Broken sync: constraints and indexes, but no columns.
					Name: "t1",
					Indexes: []*metadatapb.IndexMetadata{
						{
							Name:         "pk_t1",
							Primary:      true,
							Unique:       true,
							IsConstraint: true,
							Expressions:  []string{"c1", "c2"},
						},
						{
							Name:         "uk_t1_c1",
							Unique:       true,
							IsConstraint: true,
							Expressions:  []string{"c1"},
						},
						{
							Name:        "idx_t1_c2",
							Type:        "btree",
							Expressions: []string{"c2"},
							Definition:  `CREATE INDEX idx_t1_c2 ON s1.t1 USING btree (c2);`,
							Comment:     "comment on a never-created index",
						},
					},
					CheckConstraints: []*metadatapb.CheckConstraintMetadata{{
						Name:       "chk_t1_c1",
						Expression: "(c1 > 0)",
					}},
					ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
						Name:              "fk_t1_t2",
						Columns:           []string{"c1"},
						ReferencedSchema:  "s1",
						ReferencedTable:   "t2",
						ReferencedColumns: []string{"id"},
					}},
					Triggers: []*metadatapb.TriggerMetadata{{
						Name:    "trg_t1_c1",
						Body:    `CREATE TRIGGER trg_t1_c1 BEFORE UPDATE OF c1 ON s1.t1 FOR EACH ROW EXECUTE FUNCTION s1.f()`,
						Comment: "comment on a never-created trigger",
					}},
					Rules: []*metadatapb.RuleMetadata{{
						Name:       "rule_t1",
						Event:      "INSERT",
						Definition: `CREATE RULE rule_t1 AS ON INSERT TO s1.t1 WHERE new.c1 > 0 DO INSTEAD NOTHING;`,
					}},
				},
				{
					// Healthy table; its own constraints must still dump, but
					// its foreign key into the broken table must not.
					Name: "t2",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "id", Type: "integer", Nullable: false},
						{Name: "t1_c1", Type: "integer", Nullable: true},
					},
					Indexes: []*metadatapb.IndexMetadata{{
						Name:         "pk_t2",
						Primary:      true,
						Unique:       true,
						IsConstraint: true,
						Expressions:  []string{"id"},
					}},
					ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
						Name:              "fk_t2_t1",
						Columns:           []string{"t1_c1"},
						ReferencedSchema:  "s1",
						ReferencedTable:   "t1",
						ReferencedColumns: []string{"c1"},
					}},
				},
				{
					// Broken sync on a partitioned table with no indexes: the
					// partition key references missing columns, so the table
					// must dump bare, without the PARTITION BY clause.
					Name: "t3",
					Partitions: []*metadatapb.TablePartitionMetadata{{
						Name:       "t3_p1",
						Expression: "RANGE (created_at)",
					}},
				},
				{
					// Genuine zero-column table: no indexes, foreign keys, or
					// partitions. Its column-independent CHECK must be kept.
					Name: "t4",
					CheckConstraints: []*metadatapb.CheckConstraintMetadata{{
						Name:       "chk_t4_ok",
						Expression: "(1 = 1)",
					}},
				},
			},
		}},
	}

	banned := []string{"pk_t1", "uk_t1_c1", "idx_t1_c2", "chk_t1_c1", "fk_t1_t2", "fk_t2_t1", "trg_t1_c1", "rule_t1", "t3_p1", "PARTITION BY", "OWNED BY", "COMMENT ON INDEX"}

	for name, ctx := range map[string]DefinitionOptions{
		"dump": {},
		"sdl":  {SDLFormat: true},
	} {
		ddl, err := DatabaseDefinition(meta, ctx)
		if err != nil {
			t.Fatalf("[%s] GetDatabaseDefinition: %v", name, err)
		}

		if !strings.Contains(ddl, `CREATE TABLE "s1"."t1" (`) {
			t.Errorf("[%s] expected bare CREATE TABLE for t1, got:\n%s", name, ddl)
		}
		for _, b := range banned {
			if strings.Contains(ddl, b) {
				t.Errorf("[%s] DDL references %q, which depends on columns missing from t1's metadata; got:\n%s", name, b, ddl)
			}
		}
		if !strings.Contains(ddl, "pk_t2") {
			t.Errorf("[%s] expected pk_t2 on the healthy table, got:\n%s", name, ddl)
		}
		if !strings.Contains(ddl, `CREATE SEQUENCE "s1"."t1_c1_seq"`) {
			t.Errorf("[%s] expected CREATE SEQUENCE for t1_c1_seq, got:\n%s", name, ddl)
		}
		if !strings.Contains(ddl, `CREATE TABLE "s1"."t3" (`) {
			t.Errorf("[%s] expected bare CREATE TABLE for t3, got:\n%s", name, ddl)
		}
		if !strings.Contains(ddl, "chk_t4_ok") {
			t.Errorf("[%s] expected genuine zero-column table t4 to keep its CHECK constraint, got:\n%s", name, ddl)
		}

		if _, parseErr := New().Exec(ddl, &ExecOptions{ContinueOnError: true}); parseErr != nil {
			t.Errorf("[%s] omni catalog parse error: %v\nDDL:\n%s", name, parseErr, ddl)
		}
	}

	multiFile, err := MultiFileDefinition(meta)
	if err != nil {
		t.Fatalf("GetMultiFileDatabaseDefinition: %v", err)
	}
	var combined strings.Builder
	for _, file := range multiFile {
		combined.WriteString(file.Content)
	}
	for _, b := range banned {
		if strings.Contains(combined.String(), b) {
			t.Errorf("[multi-file] DDL references %q, which depends on columns missing from t1's metadata; got:\n%s", b, combined.String())
		}
	}
}

func TestMultiFileDefinition_CompositeTypes(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				EnumTypes: []*metadatapb.EnumTypeMetadata{
					{Name: "mood", Values: []string{"happy", "sad"}},
				},
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{
					{
						Name: "wrapper",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "b", Type: "public.base"},
						},
					},
					{
						Name: "base",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "x", Type: "integer"},
						},
						Comment: "base type",
					},
				},
			},
		},
	}

	result, err := MultiFileDefinition(metadata)
	require.NoError(t, err)

	var typesFile string
	var combined strings.Builder
	for _, file := range result {
		if file.Name == "schemas/public/types.sql" {
			typesFile = file.Content
		}
		combined.WriteString(file.Content)
		combined.WriteString("\n")
	}
	require.NotEmpty(t, typesFile, "types.sql should exist")

	assert.Contains(t, typesFile, `CREATE TYPE "public"."mood" AS ENUM (`)
	assert.Contains(t, typesFile, `CREATE TYPE "public"."base" AS (`)
	assert.Contains(t, typesFile, `CREATE TYPE "public"."wrapper" AS (`)
	assert.Contains(t, typesFile, `COMMENT ON TYPE "public"."base" IS 'base type';`)
	assert.Less(t,
		strings.Index(typesFile, `CREATE TYPE "public"."base" AS (`),
		strings.Index(typesFile, `CREATE TYPE "public"."wrapper" AS (`),
		"referenced composite must be emitted before its dependent")

	_, err = LoadSDL(combined.String())
	require.NoError(t, err, "combined multi-file SDL should load")
}

func TestEventTriggerSDLMultiFileOutput(t *testing.T) {
	result, err := MultiFileDefinition(eventTriggerSDLMetadata())
	require.NoError(t, err)
	require.NotNil(t, result)

	fileMap := make(map[string]string)
	for _, file := range result {
		fileMap[file.Name] = file.Content
	}

	eventTriggerFile, ok := fileMap["event_triggers.sql"]
	require.True(t, ok, "event_triggers.sql file should exist")
	assert.Contains(t, eventTriggerFile, `CREATE EVENT TRIGGER "audit_ddl_start" ON ddl_command_start`)
	assert.Contains(t, eventTriggerFile, `WHEN TAG IN ('CREATE TABLE')`)
	assert.Contains(t, eventTriggerFile, `EXECUTE FUNCTION "public"."audit_ddl"();`)
	assert.Contains(t, eventTriggerFile, `COMMENT ON EVENT TRIGGER "audit_ddl_start" IS 'Audit DDL start';`)
}

func TestTriggerSDLMultiFileOutput(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Name: "testdb",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "users",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
						},
						Triggers: []*metadatapb.TriggerMetadata{
							{
								Name: "audit_trigger",
								Body: "CREATE TRIGGER audit_trigger AFTER INSERT ON public.users FOR EACH ROW EXECUTE FUNCTION audit_log()",
							},
						},
					},
				},
			},
		},
	}

	result, err := MultiFileDefinition(metadata)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Find the users table file
	var usersFileContent string
	for _, file := range result {
		if file.Name == "schemas/public/tables/users.sql" {
			usersFileContent = file.Content
			break
		}
	}

	require.NotEmpty(t, usersFileContent, "Should have users table file")
	assert.Contains(t, usersFileContent, "CREATE TRIGGER audit_trigger")
	assert.Contains(t, usersFileContent, "AFTER INSERT ON public.users")
}
