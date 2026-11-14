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
