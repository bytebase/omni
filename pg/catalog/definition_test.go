package catalog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/parser"
)

// TestCheckConstraintNotValidFormatNormalMode tests that CHECK constraints with NOT VALID
// are formatted correctly in normal (non-SDL) mode
func TestCheckConstraintNotValidFormatNormalMode(t *testing.T) {
	metadata := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "namespace_settings",
						Columns: []*metadatapb.ColumnMetadata{
							{
								Name:     "namespace_id",
								Type:     "bigint",
								Nullable: false,
							},
							{
								Name:     "default_branch_protection_defaults",
								Type:     "jsonb",
								Nullable: true,
							},
						},
						CheckConstraints: []*metadatapb.CheckConstraintMetadata{
							{
								Name:       "default_branch_protection_defaults_size_constraint",
								Expression: "(octet_length(default_branch_protection_defaults::text) <= 1024) NOT VALID",
							},
						},
					},
				},
			},
		},
	}

	result, err := DatabaseDefinition(metadata, DefinitionOptions{})
	require.NoError(t, err)

	// Verify the CHECK constraint is output correctly in the CREATE TABLE statement
	assert.Contains(t, result, `CONSTRAINT "default_branch_protection_defaults_size_constraint" CHECK (octet_length(default_branch_protection_defaults::text) <= 1024) NOT VALID`)

	// Additional validation: ensure NOT VALID is outside parentheses (not inside extra parentheses)
	assert.NotContains(t, result, ") NOT VALID)", "NOT VALID should not be inside closing parenthesis")
	assert.Contains(t, result, ") NOT VALID", "NOT VALID should be after CHECK expression parenthesis")

	// Validate that the generated SQL can be parsed without errors using ANTLR parser
	_, err = parser.Parse(result)
	require.NoError(t, err, "Generated SQL should be parseable by PostgreSQL parser")
}

// TestDatabaseDefinition_IndexCommentUsesRealNewline reproduces the A1
// bug where writeIndexComment used a Go raw string literal with `\n\n`,
// producing four literal characters `\`, `n`, `\`, `n` instead of two line
// feeds. The regression surfaced as omni pgparser choking on the backslash
// right after `COMMENT ON INDEX ... IS '...';`.
func TestDatabaseDefinition_IndexCommentUsesRealNewline(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Type: "integer", Nullable: false},
				},
				Indexes: []*metadatapb.IndexMetadata{{
					Name:        "idx_t_id",
					Type:        "btree",
					Expressions: []string{"id"},
					Descending:  []bool{false},
					Definition:  "CREATE INDEX idx_t_id ON public.t USING btree (id);",
					Comment:     "sample comment",
				}},
			}},
		}},
	}

	ddl, err := DatabaseDefinition(meta, DefinitionOptions{})
	if err != nil {
		t.Fatalf("GetDatabaseDefinition: %v", err)
	}

	if strings.Contains(ddl, `\n`) {
		t.Errorf("DDL contains literal backslash-n escape; got:\n%s", ddl)
	}
	if !strings.Contains(ddl, `COMMENT ON INDEX "public"."idx_t_id" IS 'sample comment';`) {
		t.Errorf("expected COMMENT ON INDEX statement, got:\n%s", ddl)
	}

	cat := New()
	if _, execErr := cat.Exec(ddl, &ExecOptions{ContinueOnError: true}); execErr != nil {
		t.Errorf("omni catalog parse error: %v\nDDL:\n%s", execErr, ddl)
	}
}

// TestDatabaseDefinition_TableWithOnlyCheckConstraints reproduces the A2
// bug where a table with zero columns but with CHECK constraints produced
// `CREATE TABLE "..."."..." (,` — an invalid leading comma. We force the
// case by declaring a CHECK constraint on a table with no regular columns.
func TestDatabaseDefinition_TableWithOnlyCheckConstraints(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "cons_only",
				CheckConstraints: []*metadatapb.CheckConstraintMetadata{{
					Name:       "cons_only_dummy",
					Expression: "(1 = 1)",
				}},
			}},
		}},
	}

	ddl, err := DatabaseDefinition(meta, DefinitionOptions{})
	if err != nil {
		t.Fatalf("GetDatabaseDefinition: %v", err)
	}
	if strings.Contains(ddl, "(,") {
		t.Errorf("DDL still has `(,` leading comma; got:\n%s", ddl)
	}
	if !strings.Contains(ddl, `CREATE TABLE "public"."cons_only" (`) {
		t.Errorf("expected CREATE TABLE for cons_only, got:\n%s", ddl)
	}
	// A genuine zero-column table (no indexes or foreign keys) must keep its
	// column-independent CHECK constraint — only privilege-filtered broken
	// metadata goes down the bare-table path.
	if !strings.Contains(ddl, "cons_only_dummy") {
		t.Errorf("expected CHECK constraint cons_only_dummy to be preserved, got:\n%s", ddl)
	}
	// Omni will reject a table with zero columns semantically, but it must
	// at least parse cleanly; the test asserts the parse-level contract.
	if _, parseErr := New().Exec(ddl, &ExecOptions{ContinueOnError: true}); parseErr != nil {
		t.Errorf("omni catalog parse error: %v\nDDL:\n%s", parseErr, ddl)
	}
}

// TestDatabaseDefinition_TableNameContainingDot reproduces A3: a table
// whose name literally contains a period was writing out as `"".".."` —
// because the schema + "." + object object-ID was split back on the first
// dot and produced an empty schema. The fix uses a NUL separator.
func TestDatabaseDefinition_TableNameContainingDot(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "weird.name",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "c", Type: "text", Nullable: true},
				},
			}},
		}},
	}

	ddl, err := DatabaseDefinition(meta, DefinitionOptions{})
	if err != nil {
		t.Fatalf("GetDatabaseDefinition: %v", err)
	}
	if strings.Contains(ddl, `""."`) {
		t.Errorf("DDL contains empty schema `\"\".\"` prefix; got:\n%s", ddl)
	}
	if !strings.Contains(ddl, `CREATE TABLE "public"."weird.name"`) {
		t.Errorf("expected schema to be public, got:\n%s", ddl)
	}
	if _, parseErr := New().Exec(ddl, &ExecOptions{ContinueOnError: true}); parseErr != nil {
		t.Errorf("omni catalog parse error: %v\nDDL:\n%s", parseErr, ddl)
	}
}

func TestParseQualifiedTypeIdent(t *testing.T) {
	testCases := []struct {
		input  string
		schema string
		name   string
		ok     bool
	}{
		{"integer", "", "", false},
		{"character varying(50)", "", "", false},
		{"numeric(9,6)", "", "", false},
		{"text[]", "", "", false},
		{"timestamp(6) with time zone", "", "", false},
		{"public.addr", "public", "addr", true},
		{"public.addr[]", "public", "addr", true},
		{"public.addr(5)", "public", "addr", true},
		{`"My Schema"."Weird Type"`, "My Schema", "Weird Type", true},
		{`"has""quote".t`, `has"quote`, "t", true},
		{`other.geo`, "other", "geo", true},
	}
	for _, tc := range testCases {
		schemaName, typeName, ok := parseQualifiedTypeIdent(tc.input)
		assert.Equal(t, tc.ok, ok, "input %q", tc.input)
		assert.Equal(t, tc.schema, schemaName, "input %q", tc.input)
		assert.Equal(t, tc.name, typeName, "input %q", tc.input)
	}
}

func TestSortCompositeTypesTopologically(t *testing.T) {
	// Chain across schemas: public.a3 -> other.b2 -> public.c1, with names
	// chosen so alphabetical order is the reverse of dependency order.
	c1 := &metadatapb.CompositeTypeMetadata{
		Name:       "c1",
		Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "x", Type: "integer"}},
	}
	b2 := &metadatapb.CompositeTypeMetadata{
		Name:       "b2",
		Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "c", Type: "public.c1"}},
	}
	a3 := &metadatapb.CompositeTypeMetadata{
		Name:       "a3",
		Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "b", Type: "other.b2[]"}},
	}
	ordered := sortCompositeTypesTopologically([]qualifiedCompositeType{
		{Schema: "public", Composite: a3},
		{Schema: "other", Composite: b2},
		{Schema: "public", Composite: c1},
	})
	require.Len(t, ordered, 3)
	assert.Equal(t, "c1", ordered[0].Composite.Name)
	assert.Equal(t, "b2", ordered[1].Composite.Name)
	assert.Equal(t, "a3", ordered[2].Composite.Name)

	// Independent types come out in deterministic (schema, name) order.
	ordered = sortCompositeTypesTopologically([]qualifiedCompositeType{
		{Schema: "public", Composite: &metadatapb.CompositeTypeMetadata{Name: "zz"}},
		{Schema: "public", Composite: &metadatapb.CompositeTypeMetadata{Name: "aa"}},
	})
	require.Len(t, ordered, 2)
	assert.Equal(t, "aa", ordered[0].Composite.Name)
	assert.Equal(t, "zz", ordered[1].Composite.Name)
}

func TestIsDefinitionProcedure(t *testing.T) {
	tests := []struct {
		name        string
		definition  string
		isProcedure bool
	}{
		{
			name: "Simple FUNCTION",
			definition: `CREATE FUNCTION test_function()
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
	NULL;
END;
$$;`,
			isProcedure: false,
		},
		{
			name: "Simple PROCEDURE",
			definition: `CREATE PROCEDURE test_procedure()
LANGUAGE plpgsql
AS $$
BEGIN
	NULL;
END;
$$;`,
			isProcedure: true,
		},
		{
			name: "CREATE OR REPLACE FUNCTION",
			definition: `CREATE OR REPLACE FUNCTION test_function()
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
	NULL;
END;
$$;`,
			isProcedure: false,
		},
		{
			name: "CREATE OR REPLACE PROCEDURE",
			definition: `CREATE OR REPLACE PROCEDURE test_procedure()
LANGUAGE plpgsql
AS $$
BEGIN
	NULL;
END;
$$;`,
			isProcedure: true,
		},
		{
			name: "FUNCTION with word 'procedure' in comment - should NOT be confused",
			definition: `CREATE FUNCTION test_function()
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
	-- This is not a procedure, it's a function
	-- The word PROCEDURE appears in this comment
	RAISE NOTICE 'This function is not a procedure';
END;
$$;`,
			isProcedure: false,
		},
		{
			name: "FUNCTION with word 'procedure' in string literal - should NOT be confused",
			definition: `CREATE FUNCTION test_function()
RETURNS text
LANGUAGE plpgsql
AS $$
BEGIN
	RETURN 'This function is not a PROCEDURE';
END;
$$;`,
			isProcedure: false,
		},
		{
			name: "PROCEDURE with complex body",
			definition: `CREATE OR REPLACE PROCEDURE update_user_data(
	user_id INTEGER,
	new_name VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
	old_name VARCHAR;
BEGIN
	-- Get old name
	SELECT name INTO old_name FROM users WHERE id = user_id;

	-- Update the user name
	UPDATE users SET name = new_name WHERE id = user_id;

	-- Log the change
	INSERT INTO audit_log (message) VALUES (
		'Changed user name from ' || old_name || ' to ' || new_name
	);
END;
$$;`,
			isProcedure: true,
		},
		{
			name: "FUNCTION with RETURNS TABLE",
			definition: `CREATE OR REPLACE FUNCTION get_users()
RETURNS TABLE(id INTEGER, name VARCHAR)
LANGUAGE plpgsql
AS $$
BEGIN
	RETURN QUERY SELECT users.id, users.name FROM users;
END;
$$;`,
			isProcedure: false,
		},
		{
			name:        "Empty definition",
			definition:  "",
			isProcedure: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDefinitionProcedure(tt.definition)
			assert.Equal(t, tt.isProcedure, result,
				"isDefinitionProcedure returned %v, expected %v for:\n%s",
				result, tt.isProcedure, tt.definition)
		})
	}
}

// TestIsDefinitionProcedureRobustness tests edge cases that would fail with string-based detection
func TestIsDefinitionProcedureRobustness(t *testing.T) {
	tests := []struct {
		name        string
		definition  string
		isProcedure bool
		reason      string
	}{
		{
			name: "FUNCTION with 'CREATE PROCEDURE' in string literal",
			definition: `CREATE FUNCTION test_function()
RETURNS text
LANGUAGE sql
AS $outer$
	SELECT 'Example: CREATE PROCEDURE foo() LANGUAGE plpgsql AS $body$ BEGIN NULL; END; $body$;'
$outer$;`,
			isProcedure: false,
			reason:      "String-based detection might incorrectly identify this as a PROCEDURE due to 'CREATE PROCEDURE' in the string literal",
		},
		{
			name: "FUNCTION with 'PROCEDURE' in multi-line comment",
			definition: `CREATE FUNCTION complex_function()
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
	/*
	 * This function performs complex operations
	 * Note: This is NOT a PROCEDURE
	 * PROCEDURE keyword appears in this comment
	 * But it should still be detected as a FUNCTION
	 */
	RAISE NOTICE 'Processing...';
END;
$$;`,
			isProcedure: false,
			reason:      "Multi-line comments with PROCEDURE keyword should not confuse the detector",
		},
		{
			name: "PROCEDURE with unusual formatting",
			definition: `CREATE    OR    REPLACE    PROCEDURE
test_procedure
(
)
LANGUAGE plpgsql
AS
$$
BEGIN
	NULL;
END;
$$;`,
			isProcedure: true,
			reason:      "Should handle unusual whitespace and formatting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDefinitionProcedure(tt.definition)
			assert.Equal(t, tt.isProcedure, result,
				"Test failed: %s\nGot: %v, Expected: %v",
				tt.reason, result, tt.isProcedure)
		})
	}
}

// TestDatabaseDefinitionConstraintComments: a check constraint's and a
// foreign key's comment, which pg/driver carries from COMMENT ON CONSTRAINT,
// are written beside their constraint and read back by the catalog.
func TestDatabaseDefinitionConstraintComments(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{
				{
					Name:    "parent",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Nullable: false}},
					Indexes: []*metadatapb.IndexMetadata{{
						Name: "parent_pkey", Type: "btree", Expressions: []string{`"id"`},
						Unique: true, Primary: true, IsConstraint: true,
					}},
				},
				{
					Name: "child",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "id", Type: "integer", Nullable: false},
						{Name: "parent_id", Type: "integer", Nullable: true},
					},
					CheckConstraints: []*metadatapb.CheckConstraintMetadata{
						{Name: "child_id_check", Expression: "(id > 0)", Comment: "a child has a positive id"},
						{Name: "child_id_cap", Expression: "(id < 1000000)"},
					},
					ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
						Name: "child_parent_fk", Columns: []string{"parent_id"},
						ReferencedSchema: "public", ReferencedTable: "parent", ReferencedColumns: []string{"id"},
						OnDelete: "NO ACTION", OnUpdate: "NO ACTION", Comment: "the parent of a child",
					}},
				},
				{
					Name: "events",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "id", Type: "integer", Nullable: false},
						{Name: "at", Type: "date", Nullable: false},
					},
					Partitions: []*metadatapb.TablePartitionMetadata{{
						Name:       "events_2026",
						Type:       metadatapb.TablePartitionMetadata_RANGE,
						Expression: "RANGE (at)",
						Value:      "FOR VALUES FROM ('2026-01-01') TO ('2027-01-01')",
						CheckConstraints: []*metadatapb.CheckConstraintMetadata{
							{Name: "events_2026_id_check", Expression: "(id > 0)", Comment: "a partition's rows are positive"},
						},
					}},
				},
			},
		}},
	}
	ddl, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	for _, want := range []string{
		`COMMENT ON CONSTRAINT "child_id_check" ON "public"."child" IS 'a child has a positive id';`,
		`COMMENT ON CONSTRAINT "child_parent_fk" ON "public"."child" IS 'the parent of a child';`,
		`COMMENT ON CONSTRAINT "events_2026_id_check" ON "public"."events_2026" IS 'a partition''s rows are positive';`,
	} {
		assert.Contains(t, ddl, want)
	}
	// A constraint the engine states no comment beside writes none.
	assert.NotContains(t, ddl, `COMMENT ON CONSTRAINT "child_id_cap"`)
	// A partition's comment is a statement of its own: written into the
	// CREATE TABLE before its PARTITION BY clause it would not parse.
	assert.Greater(t, strings.Index(ddl, `COMMENT ON CONSTRAINT "events_2026_id_check"`), strings.Index(ddl, `CREATE TABLE "public"."events_2026"`))

	c := New()
	results, err := c.Exec(ddl, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, result := range results {
		require.NoErrorf(t, result.Error, "statement %d (%s) failed", result.Index, result.SQL)
	}
	_, rel, err := c.findRelation("public", "child")
	require.NoError(t, err)
	for name, want := range map[string]string{
		"child_id_check":  "a child has a positive id",
		"child_parent_fk": "the parent of a child",
	} {
		var found bool
		for _, con := range c.consByRel[rel.OID] {
			if con.Name != name {
				continue
			}
			found = true
			if got, ok := c.GetComment('c', con.OID, 0); !ok || got != want {
				t.Errorf("%s's comment is %q (present=%v), want %q", name, got, ok, want)
			}
		}
		if !found {
			t.Errorf("constraint %s is not in the catalog", name)
		}
	}
}

// TestSchemaDefinitionWritesAnExtensionSchemasTypes holds SchemaDefinition to
// the types of the schema it is asked for even when the schema is skipped: the
// caller asked for this schema, and a table of it that names an enum needs the
// enum's declaration.
func TestSchemaDefinitionWritesAnExtensionSchemasTypes(t *testing.T) {
	got, err := SchemaDefinition(&metadatapb.SchemaMetadata{
		Name:     "ext",
		SkipDump: true,
		EnumTypes: []*metadatapb.EnumTypeMetadata{
			{Name: "mood", Values: []string{"happy"}},
		},
		Tables: []*metadatapb.TableMetadata{
			{Name: "t", Columns: []*metadatapb.ColumnMetadata{{Name: "m", Type: "ext.mood"}}},
		},
	})
	require.NoError(t, err)
	assert.Contains(t, got, `CREATE TYPE "ext"."mood" AS ENUM (
    'happy'
)`)
	assert.Contains(t, got, `CREATE TABLE "ext"."t"`)
}

// TestDatabaseDefinitionWritesForeignServerBeforeForeignTable checks that a
// foreign server and the options it and its table state are written, and that
// the server comes before the table that names it: CREATE FOREIGN TABLE fails
// if its server does not exist.
func TestDatabaseDefinitionWritesForeignServerBeforeForeignTable(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name:    "loopback",
			Fdw:     "postgres_fdw",
			Options: []string{"host=localhost", "dbname=postgres"},
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "audit",
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "remote_orders",
				ExternalServerName: "loopback",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Type: "integer", Nullable: false, Options: []string{"column_name=remote_id"}},
					{Name: "note", Type: "text", Nullable: true},
				},
				Options: []string{"schema_name=public", "table_name=orders"},
			}},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	server := `CREATE SERVER "loopback" FOREIGN DATA WRAPPER "postgres_fdw" OPTIONS (host 'localhost', dbname 'postgres');`
	table := `CREATE FOREIGN TABLE "audit"."remote_orders" (
    "id" integer OPTIONS (column_name 'remote_id') NOT NULL,
    "note" text
) SERVER "loopback" OPTIONS (schema_name 'public', table_name 'orders');`
	assert.Contains(t, result, server)
	assert.Contains(t, result, table)
	assert.Less(t, strings.Index(result, server), strings.Index(result, table), "the server must be written before its foreign table")
}

// TestDatabaseDefinitionWritesForeignDataWrapperBeforeServer checks that a
// wrapper the database declares itself, with its handler, validator and
// options, is written before the server that names it.
func TestDatabaseDefinitionWritesForeignDataWrapperBeforeServer(t *testing.T) {
	serverType, serverVersion := "foo", "1"
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:      "custom_fdw",
			Handler:   "custom_handler",
			Validator: "custom_validator",
			Options:   []string{"debug=true", "ApiKey=secret", `"lit"=v`},
		}},
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name:    "custom_server",
			Fdw:     "custom_fdw",
			Type:    &serverType,
			Version: &serverVersion,
		}},
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	wrapper := `CREATE FOREIGN DATA WRAPPER "custom_fdw" HANDLER custom_handler VALIDATOR custom_validator OPTIONS (debug 'true', "ApiKey" 'secret', """lit""" 'v');`
	server := `CREATE SERVER "custom_server" TYPE 'foo' VERSION '1' FOREIGN DATA WRAPPER "custom_fdw";`
	assert.Contains(t, result, wrapper)
	assert.Contains(t, result, server)
	assert.Less(t, strings.Index(result, wrapper), strings.Index(result, server), "the wrapper must be written before its server")
}

// TestDatabaseDefinitionOrdersForeignDataWrapperAfterItsFunctions checks that a
// database-declared wrapper follows a user function its handler or validator
// names.
func TestDatabaseDefinitionOrdersForeignDataWrapperAfterItsFunctions(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "app",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "handler",
				Signature:  "handler()",
				Definition: "CREATE FUNCTION app.handler() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
			}},
		}},
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:    "custom_fdw",
			Handler: "app.handler",
		}},
		ForeignServers: []*metadatapb.ForeignServerMetadata{{Name: "custom_server", Fdw: "custom_fdw"}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	function := `CREATE FUNCTION app.handler()`
	wrapper := `CREATE FOREIGN DATA WRAPPER "custom_fdw" HANDLER app.handler;`
	server := `CREATE SERVER "custom_server" FOREIGN DATA WRAPPER "custom_fdw";`
	assert.Contains(t, result, function)
	assert.Contains(t, result, wrapper)
	assert.Contains(t, result, server)
	assert.Less(t, strings.Index(result, function), strings.Index(result, wrapper), "the handler must be written before its wrapper")
	assert.Less(t, strings.Index(result, wrapper), strings.Index(result, server), "the wrapper must be written before its server")
}

// TestDatabaseDefinitionOrdersForeignDataWrapperAfterQuotedHandler checks that a
// handler regproc writes with quotes still resolves to its function.
func TestDatabaseDefinitionOrdersForeignDataWrapperAfterQuotedHandler(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "app",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "MyHandler",
				Signature:  "MyHandler()",
				Definition: `CREATE FUNCTION app."MyHandler"() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;`,
			}},
		}},
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:    "custom_fdw",
			Handler: `app."MyHandler"`,
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	function := `CREATE FUNCTION app."MyHandler"()`
	wrapper := `CREATE FOREIGN DATA WRAPPER "custom_fdw" HANDLER app."MyHandler";`
	assert.Contains(t, result, function)
	assert.Contains(t, result, wrapper)
	assert.Less(t, strings.Index(result, function), strings.Index(result, wrapper), "the quoted handler must be written before its wrapper")
}

// TestDatabaseDefinitionKeepsForeignTableInExtensionSchema checks that a foreign
// table the extension does not own is written even when its schema is
// extension-owned, and one the extension owns is left to the extension.
func TestDatabaseDefinitionKeepsForeignTableInExtensionSchema(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "loopback",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:     "ext_schema",
			SkipDump: true,
			ExternalTables: []*metadatapb.ExternalTableMetadata{
				{Name: "user_table", ExternalServerName: "loopback"},
				{Name: "extension_table", ExternalServerName: "loopback", SkipDump: true},
			},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	assert.Contains(t, result, `CREATE FOREIGN TABLE "ext_schema"."user_table"`)
	assert.NotContains(t, result, `"ext_schema"."extension_table"`)
}

// TestDatabaseDefinitionScopesQualifiedHandlerToItsSchema checks that a
// qualified handler does not also depend on a same-named function in another
// schema. That false edge can close a cycle with the table the other function
// reads, and the cycle fallback would then put the wrapper before its handler.
func TestDatabaseDefinitionScopesQualifiedHandlerToItsSchema(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:    "custom_fdw",
			Handler: "app.validate",
		}},
		ForeignServers: []*metadatapb.ForeignServerMetadata{{Name: "s", Fdw: "custom_fdw"}},
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "app",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "validate",
					Signature:  "validate()",
					Definition: "CREATE FUNCTION app.validate() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
				}, {
					Name:             "validate",
					Signature:        "validate(public.pt)",
					Definition:       "CREATE FUNCTION app.validate(r public.pt) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
					DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "pt"}},
				}},
			},
			{
				Name: "other",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:             "validate",
					Signature:        "validate(public.pt)",
					Definition:       "CREATE FUNCTION other.validate(r public.pt) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
					DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "pt"}},
				}},
			},
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name:       "pt",
					Partitions: []*metadatapb.TablePartitionMetadata{{Name: "pt_remote", Schema: "public", Value: "FOR VALUES FROM (0) TO (10)"}},
				}},
				ExternalTables: []*metadatapb.ExternalTableMetadata{{Name: "pt_remote", ExternalServerName: "s"}},
			},
		},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	handler := "CREATE FUNCTION app.validate()"
	wrapper := `CREATE FOREIGN DATA WRAPPER "custom_fdw" HANDLER app.validate;`
	assert.Contains(t, result, handler)
	assert.Contains(t, result, wrapper)
	assert.Less(t, strings.Index(result, handler), strings.Index(result, wrapper), "the qualified handler must be written before its wrapper")
}

// TestDatabaseDefinitionRendersDatabaseLevelObjectsWithoutSchemas checks that a
// database with no user schemas still dumps its wrappers and servers.
func TestDatabaseDefinitionRendersDatabaseLevelObjectsWithoutSchemas(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name: "bare_fdw",
		}},
		ForeignServers: []*metadatapb.ForeignServerMetadata{{Name: "s", Fdw: "bare_fdw"}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	assert.Contains(t, result, `CREATE FOREIGN DATA WRAPPER "bare_fdw";`)
	assert.Contains(t, result, `CREATE SERVER "s" FOREIGN DATA WRAPPER "bare_fdw";`)
}

// TestDatabaseDefinitionOrdersViewAfterForeignPartition checks that a view that
// reads a foreign partition waits for the parent that creates the partition.
func TestDatabaseDefinitionOrdersViewAfterForeignPartition(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:    "custom_fdw",
			Handler: "app.handler",
		}},
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "custom_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "app",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "handler",
				Signature:  "handler()",
				Definition: "CREATE FUNCTION app.handler() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
			}},
		}, {
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:       "pt",
				Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				Partitions: []*metadatapb.TablePartitionMetadata{{Name: "pt_remote", Value: "FOR VALUES FROM (0) TO (10)"}},
			}},
			ExternalTables: []*metadatapb.ExternalTableMetadata{{Name: "pt_remote", ExternalServerName: "s"}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:             "read_remote",
				Signature:        "read_remote(public.pt_remote)",
				Definition:       "CREATE FUNCTION public.read_remote(r public.pt_remote) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
				DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "pt_remote"}},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "v",
				Definition:        "SELECT id FROM public.pt_remote",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "pt_remote"}},
			}},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	partition := `CREATE FOREIGN TABLE "public"."pt_remote"`
	function := "CREATE FUNCTION public.read_remote(r public.pt_remote)"
	view := `CREATE VIEW "public"."v"`
	assert.Contains(t, result, partition)
	assert.Contains(t, result, function)
	assert.Contains(t, result, view)
	assert.Less(t, strings.Index(result, partition), strings.Index(result, function), "the foreign partition must exist before the function takes its row type")
	assert.Less(t, strings.Index(result, partition), strings.Index(result, view), "the foreign partition must exist before the view reads it")
}

// TestDatabaseDefinitionWritesColumnlessForeignTableCheck checks that a foreign
// table with no columns but a CHECK constraint does not get a leading comma.
func TestDatabaseDefinitionWritesColumnlessForeignTableCheck(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "ft",
				ExternalServerName: "s",
				CheckConstraints:   []*metadatapb.CheckConstraintMetadata{{Name: "c", Expression: "(true)"}},
			}},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	assert.Contains(t, result, `CREATE FOREIGN TABLE "public"."ft" (
    CONSTRAINT "c" CHECK (true)
) SERVER "s";`)
}

// TestDatabaseDefinitionOrdersForeignDataWrapperAfterNamedValidator checks that
// a validator whose arguments keep their names still matches its signature.
func TestDatabaseDefinitionOrdersForeignDataWrapperAfterNamedValidator(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "app",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "validate",
				Signature:  "validate(\"options,list\" text[], catalog oid)",
				Definition: "CREATE FUNCTION app.validate(\"options,list\" text[], catalog oid) RETURNS void LANGUAGE internal AS 'fdw_validate'",
			}},
		}},
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:      "custom_fdw",
			Validator: "app.validate",
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	function := `CREATE FUNCTION app.validate("options,list" text[], catalog oid)`
	wrapper := `CREATE FOREIGN DATA WRAPPER "custom_fdw" VALIDATOR app.validate;`
	assert.Contains(t, result, function)
	assert.Contains(t, result, wrapper)
	assert.Less(t, strings.Index(result, function), strings.Index(result, wrapper), "the named validator must be written before its wrapper")
}

// TestDatabaseDefinitionOrdersForeignTableAfterCheckFunction checks that a
// function a foreign table's CHECK calls is written before the table.
func TestDatabaseDefinitionOrdersForeignTableAfterCheckFunction(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "zz_type",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:             "validate",
				Signature:        "validate(public.zz_type)",
				Definition:       "CREATE FUNCTION public.validate(r public.zz_type) RETURNS boolean LANGUAGE sql AS $$ SELECT true $$;",
				DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "zz_type"}},
			}},
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "ft",
				ExternalServerName: "s",
				CheckConstraints: []*metadatapb.CheckConstraintMetadata{{
					Name:       "c",
					Expression: "(public.validate(ROW(id)::public.zz_type)) NOT VALID",
				}},
			}},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	function := "CREATE FUNCTION public.validate(r public.zz_type)"
	table := `CREATE FOREIGN TABLE "public"."ft"`
	assert.Contains(t, result, function)
	assert.Contains(t, result, table)
	assert.Less(t, strings.Index(result, function), strings.Index(result, table), "the CHECK's function must be written before the foreign table")
}

// TestDatabaseDefinitionEscapesForeignColumnName checks that a foreign column
// whose name holds a quote is written escaped.
func TestDatabaseDefinitionEscapesForeignColumnName(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "ft",
				ExternalServerName: "s",
				Columns:            []*metadatapb.ColumnMetadata{{Name: `a"b`, Type: "integer"}},
			}},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)
	assert.Contains(t, result, `"a""b" integer`)
}

// TestDatabaseDefinitionOrdersForeignKeyTableAfterRowType checks that a foreign
// table whose column uses another foreign table's row type follows that table,
// even when its name sorts first.
func TestDatabaseDefinitionOrdersForeignKeyTableAfterRowType(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			ExternalTables: []*metadatapb.ExternalTableMetadata{
				{Name: "a", ExternalServerName: "s", Columns: []*metadatapb.ColumnMetadata{{Name: "r", Type: "public.z"}}},
				{Name: "z", ExternalServerName: "s", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}},
			},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	z := `CREATE FOREIGN TABLE "public"."z"`
	a := `CREATE FOREIGN TABLE "public"."a"`
	assert.Contains(t, result, z)
	assert.Contains(t, result, a)
	assert.Less(t, strings.Index(result, z), strings.Index(result, a), "the row type's table must be written first")
}

// TestDatabaseDefinitionOrdersForeignKeyTableAfterArrayRowType checks that an
// array of a foreign table's row type (public._z) also orders the table first.
func TestDatabaseDefinitionOrdersForeignKeyTableAfterArrayRowType(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			ExternalTables: []*metadatapb.ExternalTableMetadata{
				{Name: "a", ExternalServerName: "s", Columns: []*metadatapb.ColumnMetadata{{Name: "rs", Type: "public._z"}}},
				{Name: "z", ExternalServerName: "s", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}},
			},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	z := `CREATE FOREIGN TABLE "public"."z"`
	a := `CREATE FOREIGN TABLE "public"."a"`
	assert.Contains(t, result, z)
	assert.Contains(t, result, a)
	assert.Less(t, strings.Index(result, z), strings.Index(result, a), "the array row type's table must be written first")
}

// TestDatabaseDefinitionChecksForeignPartitionOverloadsAgainstThePartition
// checks that a foreign partition's default function dependencies are matched
// against the partition, not the parent, so an unrelated row-type overload
// does not close a cycle that puts the parent before the default's function.
func TestDatabaseDefinitionChecksForeignPartitionOverloadsAgainstThePartition(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:       "pt",
				Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				Partitions: []*metadatapb.TablePartitionMetadata{{Name: "pt_remote", Schema: "public", Value: "FOR VALUES FROM (0) TO (10)"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "f",
					Signature:  "f(integer)",
					Definition: "CREATE FUNCTION public.f(a integer) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
				},
				{
					Name:             "f",
					Signature:        "f(public.pt_remote)",
					Definition:       "CREATE FUNCTION public.f(r public.pt_remote) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;",
					DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "pt_remote"}},
				},
			},
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "pt_remote",
				ExternalServerName: "s",
				Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Default: "public.f(1)"}},
			}},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	function := "CREATE FUNCTION public.f(a integer)"
	table := `CREATE TABLE "public"."pt"`
	assert.Contains(t, result, function)
	assert.Contains(t, result, table)
	assert.Less(t, strings.Index(result, function), strings.Index(result, table), "the default's function must be written before the parent")
}

// TestDatabaseDefinitionOrdersForeignDataWrapperAfterQuotedParenValidator
// checks that a validator whose name holds a quoted parenthesis still matches
// its argument list.
func TestDatabaseDefinitionOrdersForeignDataWrapperAfterQuotedParenValidator(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "app",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "validate(",
				Signature:  "validate((text[], oid)",
				Definition: `CREATE FUNCTION app."validate("(text[], oid) RETURNS void LANGUAGE internal AS 'fdw_validate'`,
			}},
		}},
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:      "custom_fdw",
			Validator: `app."validate("`,
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	function := `CREATE FUNCTION app."validate("(text[], oid)`
	wrapper := `CREATE FOREIGN DATA WRAPPER "custom_fdw" VALIDATOR app."validate(";`
	assert.Contains(t, result, function)
	assert.Contains(t, result, wrapper)
	assert.Less(t, strings.Index(result, function), strings.Index(result, wrapper), "the quoted-paren validator must be written before its wrapper")
}

// TestDatabaseDefinitionOrdersOrdinaryTableAfterForeignRowType checks that an
// ordinary table using a foreign table's row type follows the foreign table.
func TestDatabaseDefinitionOrdersOrdinaryTableAfterForeignRowType(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		ForeignServers: []*metadatapb.ForeignServerMetadata{{
			Name: "s",
			Fdw:  "postgres_fdw",
		}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "local",
				Columns: []*metadatapb.ColumnMetadata{{Name: "r", Type: "public.remote"}},
			}},
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "remote",
				ExternalServerName: "s",
				Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
			}},
		}},
	}

	result, err := DatabaseDefinition(meta, DefinitionOptions{})
	require.NoError(t, err)

	remote := `CREATE FOREIGN TABLE "public"."remote"`
	local := `CREATE TABLE "public"."local"`
	assert.Contains(t, result, remote)
	assert.Contains(t, result, local)
	assert.Less(t, strings.Index(result, remote), strings.Index(result, local), "the row type's foreign table must be written first")
}

// TestSchemaDefinitionSkipsExternalTableTriggers checks that SchemaDefinition,
// which does not create standalone foreign tables, does not write their
// triggers either.
func TestSchemaDefinitionSkipsExternalTableTriggers(t *testing.T) {
	schema := &metadatapb.SchemaMetadata{
		Name: "public",
		ExternalTables: []*metadatapb.ExternalTableMetadata{{
			Name:               "remote",
			ExternalServerName: "s",
			Triggers: []*metadatapb.TriggerMetadata{{
				Name: "remote_trg",
				Body: "CREATE TRIGGER remote_trg BEFORE INSERT ON public.remote FOR EACH ROW EXECUTE FUNCTION public.audit()",
			}},
		}},
	}

	result, err := SchemaDefinition(schema)
	require.NoError(t, err)
	assert.NotContains(t, result, "CREATE TRIGGER remote_trg")
}
