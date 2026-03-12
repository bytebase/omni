package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// =============================================================================
// Phase 1: Dispatch Completeness Tests
// =============================================================================

// parseOne parses a single SQL statement and returns the first node.
func parseOne(t *testing.T, sql string) nodes.Node {
	t.Helper()
	list, err := pgparser.Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v\nSQL: %s", err, sql)
	}
	if list == nil || len(list.Items) == 0 {
		t.Fatalf("parse returned no statements\nSQL: %s", sql)
	}
	item := list.Items[0]
	if raw, ok := item.(*nodes.RawStmt); ok {
		return raw.Stmt
	}
	return item
}

// TestPhase1_NewNoOpDispatch tests that all newly added no-op statement types
// are accepted by ProcessUtility without error.
func TestPhase1_NewNoOpDispatch(t *testing.T) {
	noOpSQL := []struct {
		name string
		sql  string
	}{
		// Foreign data wrappers
		{"CreateFdwStmt", "CREATE FOREIGN DATA WRAPPER test_fdw"},
		{"AlterFdwStmt", "ALTER FOREIGN DATA WRAPPER test_fdw OPTIONS (ADD host 'localhost')"},
		{"CreateForeignServerStmt", "CREATE SERVER test_srv FOREIGN DATA WRAPPER test_fdw"},
		{"AlterForeignServerStmt", "ALTER SERVER test_srv OPTIONS (ADD host 'localhost')"},
		{"CreateUserMappingStmt", "CREATE USER MAPPING FOR CURRENT_USER SERVER test_srv"},
		{"AlterUserMappingStmt", "ALTER USER MAPPING FOR CURRENT_USER SERVER test_srv OPTIONS (ADD password 'secret')"},
		{"DropUserMappingStmt", "DROP USER MAPPING FOR CURRENT_USER SERVER test_srv"},
		{"ImportForeignSchemaStmt", "IMPORT FOREIGN SCHEMA remote_schema FROM SERVER test_srv INTO public"},

		// Operators/Collations
		{"AlterCollationStmt", "ALTER COLLATION \"C\" REFRESH VERSION"},

		// Tablespaces
		{"CreateTableSpaceStmt", "CREATE TABLESPACE test_ts LOCATION '/tmp/ts'"},
		{"DropTableSpaceStmt", "DROP TABLESPACE test_ts"},
		{"AlterTableSpaceOptionsStmt", "ALTER TABLESPACE pg_default SET (seq_page_cost = 1.0)"},

		// Roles/Databases
		{"CreateRoleStmt", "CREATE ROLE test_role"},
		{"AlterRoleStmt", "ALTER ROLE test_role SUPERUSER"},
		{"AlterRoleSetStmt", "ALTER ROLE test_role SET search_path = 'public'"},
		{"DropRoleStmt", "DROP ROLE test_role"},
		{"GrantRoleStmt", "GRANT test_role TO CURRENT_USER"},
		{"CreatedbStmt", "CREATE DATABASE test_db"},
		{"AlterDatabaseStmt", "ALTER DATABASE test_db ALLOW_CONNECTIONS true"},
		{"AlterDatabaseSetStmt", "ALTER DATABASE test_db SET search_path = 'public'"},
		{"DropdbStmt", "DROP DATABASE test_db"},

		// Misc DDL
		{"AlterObjectDependsStmt", "ALTER FUNCTION pg_catalog.now() DEPENDS ON EXTENSION test_ext"},
		{"CreatePLangStmt", "CREATE LANGUAGE plpgsql"},
		{"AlterSystemStmt", "ALTER SYSTEM SET max_connections = 200"},
		{"DropOwnedStmt", "DROP OWNED BY test_role"},
		{"ReassignOwnedStmt", "REASSIGN OWNED BY test_role TO CURRENT_USER"},
		{"SecLabelStmt", "SECURITY LABEL FOR test_provider ON TABLE pg_class IS 'public'"},

		// Non-DDL utility statements
		{"DoStmt", "DO $$ BEGIN NULL; END $$"},
		{"ExplainStmt", "EXPLAIN SELECT 1"},
		{"VacuumStmt", "VACUUM"},
		{"ClusterStmt", "CLUSTER"},
		{"CheckPointStmt", "CHECKPOINT"},
		{"DiscardStmt", "DISCARD ALL"},
		{"ListenStmt", "LISTEN test_channel"},
		{"UnlistenStmt", "UNLISTEN test_channel"},
		{"NotifyStmt", "NOTIFY test_channel"},
		{"LoadStmt", "LOAD 'test_lib'"},
		{"ConstraintsSetStmt", "SET CONSTRAINTS ALL DEFERRED"},
		{"VariableShowStmt", "SHOW search_path"},
		{"CallStmt", "CALL test_proc()"},
		{"PrepareStmt", "PREPARE test_plan AS SELECT 1"},
		{"ExecuteStmt", "EXECUTE test_plan"},
		{"DeallocateStmt", "DEALLOCATE test_plan"},
	}

	c := New()
	for _, tc := range noOpSQL {
		t.Run(tc.name, func(t *testing.T) {
			stmt := parseOne(t, tc.sql)
			err := c.ProcessUtility(stmt)
			if err != nil {
				t.Errorf("expected no error for %s, got: %v", tc.name, err)
			}
		})
	}
}

// TestPhase1_DropForeignTable tests DROP FOREIGN TABLE end-to-end.
func TestPhase1_DropForeignTable(t *testing.T) {
	c := New()
	// Setup: create a foreign table via ProcessUtility.
	for _, sql := range []string{
		"CREATE FOREIGN TABLE ft (id int, name text) SERVER test_srv",
	} {
		stmt := parseOne(t, sql)
		if err := c.ProcessUtility(stmt); err != nil {
			t.Fatalf("setup error: %v", err)
		}
	}

	// Verify it exists.
	_, rel, err := c.findRelation("", "ft")
	if err != nil {
		t.Fatalf("foreign table not found: %v", err)
	}
	if rel.RelKind != 'f' {
		t.Fatalf("expected relkind 'f', got '%c'", rel.RelKind)
	}

	// Drop it.
	stmt := parseOne(t, "DROP FOREIGN TABLE ft")
	if err := c.ProcessUtility(stmt); err != nil {
		t.Fatalf("DROP FOREIGN TABLE failed: %v", err)
	}

	// Verify it's gone.
	_, _, err = c.findRelation("", "ft")
	if err == nil {
		t.Fatal("expected foreign table to be gone after DROP")
	}
}

// TestPhase1_DropForeignTableIfExists tests DROP FOREIGN TABLE IF EXISTS.
func TestPhase1_DropForeignTableIfExists(t *testing.T) {
	c := New()
	stmt := parseOne(t, "DROP FOREIGN TABLE IF EXISTS nonexistent")
	if err := c.ProcessUtility(stmt); err != nil {
		t.Fatalf("DROP FOREIGN TABLE IF EXISTS should not error: %v", err)
	}
}

// TestPhase1_DropForeignTableWrongKind tests DROP FOREIGN TABLE on a regular table.
func TestPhase1_DropForeignTableWrongKind(t *testing.T) {
	c := New()
	stmt := parseOne(t, "CREATE TABLE t (id int)")
	if err := c.ProcessUtility(stmt); err != nil {
		t.Fatalf("setup: %v", err)
	}
	stmt = parseOne(t, "DROP FOREIGN TABLE t")
	err := c.ProcessUtility(stmt)
	assertErrorCode(t, err, CodeWrongObjectType)
}

// TestPhase1_CommentOnProcedure tests COMMENT ON PROCEDURE.
func TestPhase1_CommentOnProcedure(t *testing.T) {
	c := New()
	// Create a procedure.
	stmt := parseOne(t, "CREATE PROCEDURE test_proc() LANGUAGE plpgsql AS $$ BEGIN NULL; END $$")
	if err := c.ProcessUtility(stmt); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Comment on it.
	stmt = parseOne(t, "COMMENT ON PROCEDURE test_proc() IS 'test comment'")
	if err := c.ProcessUtility(stmt); err != nil {
		t.Fatalf("COMMENT ON PROCEDURE failed: %v", err)
	}
}

// TestPhase1_DropNoOpObjectTypes tests that DROP for no-op object types returns nil.
func TestPhase1_DropNoOpObjectTypes(t *testing.T) {
	c := New()

	// DROP COLLATION, OPERATOR CLASS, etc. should all succeed as no-ops.
	dropSQL := []struct {
		name string
		sql  string
	}{
		{"DropCollation", "DROP COLLATION IF EXISTS test_coll"},
		{"DropConversion", "DROP CONVERSION IF EXISTS test_conv"},
		{"DropOperator", "DROP OPERATOR IF EXISTS +(int4, int4)"},
		{"DropOpClass", "DROP OPERATOR CLASS IF EXISTS test_opc USING btree"},
		{"DropOpFamily", "DROP OPERATOR FAMILY IF EXISTS test_opf USING btree"},
		{"DropLanguage", "DROP LANGUAGE IF EXISTS test_lang"},
		{"DropFDW", "DROP FOREIGN DATA WRAPPER IF EXISTS test_fdw"},
		{"DropServer", "DROP SERVER IF EXISTS test_srv"},
		{"DropExtension", "DROP EXTENSION IF EXISTS test_ext"},
		{"DropRule", "DROP RULE IF EXISTS test_rule ON test_tbl"},
		{"DropCast", "DROP CAST IF EXISTS (int4 AS text)"},
		{"DropPublication", "DROP PUBLICATION IF EXISTS test_pub"},
		{"DropStatistics", "DROP STATISTICS IF EXISTS test_stat"},
		{"DropTextSearchConfig", "DROP TEXT SEARCH CONFIGURATION IF EXISTS test_tsc"},
		{"DropTextSearchDict", "DROP TEXT SEARCH DICTIONARY IF EXISTS test_tsd"},
		{"DropTextSearchParser", "DROP TEXT SEARCH PARSER IF EXISTS test_tsp"},
		{"DropTextSearchTemplate", "DROP TEXT SEARCH TEMPLATE IF EXISTS test_tst"},
		{"DropAccessMethod", "DROP ACCESS METHOD IF EXISTS test_am"},
		{"DropEventTrigger", "DROP EVENT TRIGGER IF EXISTS test_evt"},
	}

	for _, tc := range dropSQL {
		t.Run(tc.name, func(t *testing.T) {
			stmt := parseOne(t, tc.sql)
			err := c.ProcessUtility(stmt)
			if err != nil {
				t.Errorf("expected no error for %s, got: %v", tc.name, err)
			}
		})
	}
}

// TestPhase1_CommentNoOpTypes tests COMMENT ON for no-op object types.
func TestPhase1_CommentNoOpTypes(t *testing.T) {
	c := New()
	noOpComments := []struct {
		name string
		sql  string
	}{
		{"Collation", "COMMENT ON COLLATION \"C\" IS 'test'"},
		{"Conversion", "COMMENT ON CONVERSION test_conv IS 'test'"},
		{"FDW", "COMMENT ON FOREIGN DATA WRAPPER test_fdw IS 'test'"},
		{"Server", "COMMENT ON SERVER test_srv IS 'test'"},
		{"Database", "COMMENT ON DATABASE test_db IS 'test'"},
		{"Tablespace", "COMMENT ON TABLESPACE pg_default IS 'test'"},
		{"Role", "COMMENT ON ROLE test_role IS 'test'"},
		{"TextSearchConfig", "COMMENT ON TEXT SEARCH CONFIGURATION english IS 'test'"},
		{"TextSearchDict", "COMMENT ON TEXT SEARCH DICTIONARY english_stem IS 'test'"},
		{"TextSearchTemplate", "COMMENT ON TEXT SEARCH TEMPLATE simple IS 'test'"},
	}

	for _, tc := range noOpComments {
		t.Run(tc.name, func(t *testing.T) {
			stmt := parseOne(t, tc.sql)
			err := c.ProcessUtility(stmt)
			if err != nil {
				t.Errorf("expected no error for %s, got: %v", tc.name, err)
			}
		})
	}
}
