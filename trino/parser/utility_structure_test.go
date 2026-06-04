package parser

import (
	"testing"

	"github.com/bytebase/omni/trino/ast"
)

// This file is the structural (Layer 2) gate for the `parser-utility` node.
// Accept/reject alone does not catch a form that "accepts" but is parsed into
// the wrong node shape (e.g. SHOW TABLES vs SHOW SCHEMAS both accept; the IN/FROM
// scope must land in the right field; SET TIME ZONE LOCAL vs an expression).
// These tests pin the parse-node fields of one representative per alternative.

// parseOneStmt parses sql, requiring exactly one statement and no errors, and
// returns it. Fails the test otherwise.
func parseOneStmt(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q): unexpected errors: %v", sql, errs)
	}
	if file == nil || len(file.Stmts) != 1 {
		got := 0
		if file != nil {
			got = len(file.Stmts)
		}
		t.Fatalf("Parse(%q): got %d statements, want 1", sql, got)
	}
	return file.Stmts[0]
}

func TestUtility_StructureShow(t *testing.T) {
	t.Run("show_tables_from_like_escape", func(t *testing.T) {
		s, ok := parseOneStmt(t, "SHOW TABLES FROM tpch.tiny LIKE 'p%' ESCAPE '#'").(*ShowStmt)
		if !ok {
			t.Fatal("not a *ShowStmt")
		}
		if s.Kind != ShowTables {
			t.Errorf("Kind = %v, want ShowTables", s.Kind)
		}
		if s.In == nil || s.In.Normalize() != "tpch.tiny" {
			t.Errorf("In = %v, want tpch.tiny", s.In)
		}
		if s.InKeyword {
			t.Error("InKeyword = true, want false (FROM)")
		}
		if !s.HasLike || s.Like != "p%" {
			t.Errorf("Like = %q (has=%v), want p%%", s.Like, s.HasLike)
		}
		if !s.HasEscape || s.Escape != "#" {
			t.Errorf("Escape = %q (has=%v), want #", s.Escape, s.HasEscape)
		}
	})

	t.Run("show_tables_in_keyword", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW TABLES IN s").(*ShowStmt)
		if !s.InKeyword {
			t.Error("InKeyword = false, want true (IN)")
		}
	})

	t.Run("show_schemas_single_identifier_scope", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW SCHEMAS FROM tpch").(*ShowStmt)
		if s.Kind != ShowSchemas {
			t.Errorf("Kind = %v, want ShowSchemas", s.Kind)
		}
		if s.In == nil || len(s.In.Parts) != 1 || s.In.Normalize() != "tpch" {
			t.Errorf("In = %v, want single-part tpch", s.In)
		}
	})

	t.Run("show_columns_name_required", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW COLUMNS FROM tpch.sf1.nation").(*ShowStmt)
		if s.Kind != ShowColumns {
			t.Errorf("Kind = %v, want ShowColumns", s.Kind)
		}
		if s.Name == nil || s.Name.Normalize() != "tpch.sf1.nation" {
			t.Errorf("Name = %v, want tpch.sf1.nation", s.Name)
		}
	})

	t.Run("show_create_table", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW CREATE TABLE sf1.orders").(*ShowStmt)
		if s.Kind != ShowCreateTable {
			t.Errorf("Kind = %v, want ShowCreateTable", s.Kind)
		}
		if s.Name == nil || s.Name.Normalize() != "sf1.orders" {
			t.Errorf("Name = %v, want sf1.orders", s.Name)
		}
	})

	t.Run("show_create_materialized_view", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW CREATE MATERIALIZED VIEW mv").(*ShowStmt)
		if s.Kind != ShowCreateMaterializedView {
			t.Errorf("Kind = %v, want ShowCreateMaterializedView", s.Kind)
		}
	})

	t.Run("show_create_function_U1", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW CREATE FUNCTION a.b.c").(*ShowStmt)
		if s.Kind != ShowCreateFunction {
			t.Errorf("Kind = %v, want ShowCreateFunction", s.Kind)
		}
	})

	t.Run("show_grants_on_table", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW GRANTS ON TABLE orders").(*ShowStmt)
		if s.Kind != ShowGrants {
			t.Errorf("Kind = %v, want ShowGrants", s.Kind)
		}
		if !s.OnTable {
			t.Error("OnTable = false, want true")
		}
		if s.Name == nil || s.Name.Normalize() != "orders" {
			t.Errorf("Name = %v, want orders", s.Name)
		}
	})

	t.Run("show_grants_bare", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW GRANTS").(*ShowStmt)
		if s.OnTable || s.Name != nil {
			t.Errorf("bare SHOW GRANTS: OnTable=%v Name=%v, want false/nil", s.OnTable, s.Name)
		}
	})

	t.Run("show_current_roles", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW CURRENT ROLES FROM hive").(*ShowStmt)
		if s.Kind != ShowRoles || !s.Current {
			t.Errorf("Kind=%v Current=%v, want ShowRoles/true", s.Kind, s.Current)
		}
		if s.In == nil || s.In.Normalize() != "hive" {
			t.Errorf("In = %v, want hive", s.In)
		}
	})

	t.Run("show_role_grants", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW ROLE GRANTS").(*ShowStmt)
		if s.Kind != ShowRoleGrants {
			t.Errorf("Kind = %v, want ShowRoleGrants", s.Kind)
		}
	})

	t.Run("show_stats_for_query_raw", func(t *testing.T) {
		s := parseOneStmt(t, "SHOW STATS FOR (SELECT * FROM nation WHERE regionkey = 1)").(*ShowStmt)
		if s.Kind != ShowStatsForQuery {
			t.Errorf("Kind = %v, want ShowStatsForQuery", s.Kind)
		}
		if s.QueryText != "SELECT * FROM nation WHERE regionkey = 1" {
			t.Errorf("QueryText = %q, want the inner query text", s.QueryText)
		}
	})

	t.Run("describe_is_show_columns", func(t *testing.T) {
		s := parseOneStmt(t, "DESCRIBE nation").(*ShowStmt)
		if s.Kind != ShowColumns {
			t.Errorf("Kind = %v, want ShowColumns (DESCRIBE alias)", s.Kind)
		}
		if s.Name == nil || s.Name.Normalize() != "nation" {
			t.Errorf("Name = %v, want nation", s.Name)
		}
	})

	t.Run("describe_input_as_table_name", func(t *testing.T) {
		// `DESCRIBE input` with no trailing name is DESCRIBE-of-table-"input".
		s := parseOneStmt(t, "DESCRIBE input").(*ShowStmt)
		if s.Kind != ShowColumns {
			t.Errorf("Kind = %v, want ShowColumns", s.Kind)
		}
		if s.Name == nil || s.Name.Normalize() != "input" {
			t.Errorf("Name = %v, want input", s.Name)
		}
	})
}

func TestUtility_StructureSession(t *testing.T) {
	t.Run("use_schema_only", func(t *testing.T) {
		s := parseOneStmt(t, "USE information_schema").(*UseStmt)
		if s.Catalog != nil {
			t.Errorf("Catalog = %v, want nil", s.Catalog)
		}
		if s.Schema == nil || s.Schema.Normalize() != "information_schema" {
			t.Errorf("Schema = %v, want information_schema", s.Schema)
		}
	})

	t.Run("use_catalog_schema", func(t *testing.T) {
		s := parseOneStmt(t, "USE hive.finance").(*UseStmt)
		if s.Catalog == nil || s.Catalog.Normalize() != "hive" {
			t.Errorf("Catalog = %v, want hive", s.Catalog)
		}
		if s.Schema == nil || s.Schema.Normalize() != "finance" {
			t.Errorf("Schema = %v, want finance", s.Schema)
		}
	})

	t.Run("set_session_property", func(t *testing.T) {
		s := parseOneStmt(t, "SET SESSION cat.prop = 'v'").(*SetSessionStmt)
		if s.Name == nil || s.Name.Normalize() != "cat.prop" {
			t.Errorf("Name = %v, want cat.prop", s.Name)
		}
		if s.Value == nil {
			t.Error("Value = nil, want an expression")
		}
	})

	t.Run("reset_session_property", func(t *testing.T) {
		s := parseOneStmt(t, "RESET SESSION cat.prop").(*ResetSessionStmt)
		if s.Name == nil || s.Name.Normalize() != "cat.prop" {
			t.Errorf("Name = %v, want cat.prop", s.Name)
		}
	})

	t.Run("set_session_authorization_identifier", func(t *testing.T) {
		s := parseOneStmt(t, "SET SESSION AUTHORIZATION alice").(*SetSessionAuthorizationStmt)
		if s.HasUserString {
			t.Error("HasUserString = true, want false (identifier form)")
		}
		if s.User == nil || s.User.Normalize() != "alice" {
			t.Errorf("User = %v, want alice", s.User)
		}
	})

	t.Run("set_session_authorization_string", func(t *testing.T) {
		s := parseOneStmt(t, "SET SESSION AUTHORIZATION 'alice'").(*SetSessionAuthorizationStmt)
		if !s.HasUserString || s.UserString != "alice" {
			t.Errorf("UserString = %q (has=%v), want alice", s.UserString, s.HasUserString)
		}
	})

	t.Run("reset_session_authorization", func(t *testing.T) {
		if _, ok := parseOneStmt(t, "RESET SESSION AUTHORIZATION").(*ResetSessionAuthorizationStmt); !ok {
			t.Error("not a *ResetSessionAuthorizationStmt")
		}
	})

	t.Run("set_role_named_in_catalog", func(t *testing.T) {
		s := parseOneStmt(t, "SET ROLE analyst IN hive").(*SetRoleStmt)
		if s.Spec != RoleNamed {
			t.Errorf("Spec = %v, want RoleNamed", s.Spec)
		}
		if s.Role == nil || s.Role.Normalize() != "analyst" {
			t.Errorf("Role = %v, want analyst", s.Role)
		}
		if s.Catalog == nil || s.Catalog.Normalize() != "hive" {
			t.Errorf("Catalog = %v, want hive", s.Catalog)
		}
	})

	t.Run("set_role_all", func(t *testing.T) {
		s := parseOneStmt(t, "SET ROLE ALL").(*SetRoleStmt)
		if s.Spec != RoleAll {
			t.Errorf("Spec = %v, want RoleAll", s.Spec)
		}
	})

	t.Run("set_role_none", func(t *testing.T) {
		s := parseOneStmt(t, "SET ROLE NONE").(*SetRoleStmt)
		if s.Spec != RoleNone {
			t.Errorf("Spec = %v, want RoleNone", s.Spec)
		}
	})

	t.Run("set_path_multi", func(t *testing.T) {
		s := parseOneStmt(t, "SET PATH a, b.c, d").(*SetPathStmt)
		if len(s.Elements) != 3 {
			t.Fatalf("Elements = %d, want 3", len(s.Elements))
		}
		if s.Elements[0].Catalog != nil || s.Elements[0].Schema.Normalize() != "a" {
			t.Errorf("element 0 = %+v, want bare schema a", s.Elements[0])
		}
		if s.Elements[1].Catalog == nil || s.Elements[1].Catalog.Normalize() != "b" ||
			s.Elements[1].Schema.Normalize() != "c" {
			t.Errorf("element 1 = %+v, want b.c", s.Elements[1])
		}
	})

	t.Run("set_time_zone_local", func(t *testing.T) {
		s := parseOneStmt(t, "SET TIME ZONE LOCAL").(*SetTimeZoneStmt)
		if !s.Local {
			t.Error("Local = false, want true")
		}
		if s.Value != nil {
			t.Errorf("Value = %v, want nil for LOCAL", s.Value)
		}
	})

	t.Run("set_time_zone_expr", func(t *testing.T) {
		s := parseOneStmt(t, "SET TIME ZONE INTERVAL '10' HOUR").(*SetTimeZoneStmt)
		if s.Local {
			t.Error("Local = true, want false")
		}
		if s.Value == nil {
			t.Error("Value = nil, want the interval expression")
		}
	})
}

func TestUtility_StructureExplainCall(t *testing.T) {
	t.Run("explain_options", func(t *testing.T) {
		// Inner SHOW TABLES is implemented, so this parses cleanly.
		s := parseOneStmt(t, "EXPLAIN (TYPE LOGICAL, FORMAT JSON) SHOW TABLES").(*ExplainStmt)
		if s.Analyze {
			t.Error("Analyze = true, want false")
		}
		if s.Type != ExplainTypeLogical {
			t.Errorf("Type = %v, want ExplainTypeLogical", s.Type)
		}
		if s.Format != ExplainFormatJSON {
			t.Errorf("Format = %v, want ExplainFormatJSON", s.Format)
		}
		if s.Statement == nil {
			t.Error("Statement = nil, want the inner SHOW node")
		}
	})

	t.Run("explain_analyze_verbose", func(t *testing.T) {
		s := parseOneStmt(t, "EXPLAIN ANALYZE VERBOSE SHOW TABLES").(*ExplainStmt)
		if !s.Analyze || !s.Verbose {
			t.Errorf("Analyze=%v Verbose=%v, want both true", s.Analyze, s.Verbose)
		}
	})

	t.Run("explain_io_validate", func(t *testing.T) {
		s := parseOneStmt(t, "EXPLAIN (TYPE IO) SHOW TABLES").(*ExplainStmt)
		if s.Type != ExplainTypeIO {
			t.Errorf("Type = %v, want ExplainTypeIO", s.Type)
		}
	})

	t.Run("call_positional", func(t *testing.T) {
		s := parseOneStmt(t, "CALL test(123, 'apple')").(*CallStmt)
		if s.Name == nil || s.Name.Normalize() != "test" {
			t.Errorf("Name = %v, want test", s.Name)
		}
		if len(s.Arguments) != 2 {
			t.Fatalf("Arguments = %d, want 2", len(s.Arguments))
		}
		for i, a := range s.Arguments {
			if a.Name != nil {
				t.Errorf("argument %d Name = %v, want nil (positional)", i, a.Name)
			}
		}
	})

	t.Run("call_named", func(t *testing.T) {
		s := parseOneStmt(t, "CALL c.s.test(name => 'apple', id => 123)").(*CallStmt)
		if s.Name == nil || s.Name.Normalize() != "c.s.test" {
			t.Errorf("Name = %v, want c.s.test", s.Name)
		}
		if len(s.Arguments) != 2 {
			t.Fatalf("Arguments = %d, want 2", len(s.Arguments))
		}
		if s.Arguments[0].Name == nil || s.Arguments[0].Name.Normalize() != "name" {
			t.Errorf("argument 0 Name = %v, want name", s.Arguments[0].Name)
		}
		if s.Arguments[1].Name == nil || s.Arguments[1].Name.Normalize() != "id" {
			t.Errorf("argument 1 Name = %v, want id", s.Arguments[1].Name)
		}
	})

	t.Run("call_empty_args", func(t *testing.T) {
		s := parseOneStmt(t, "CALL p()").(*CallStmt)
		if len(s.Arguments) != 0 {
			t.Errorf("Arguments = %d, want 0", len(s.Arguments))
		}
	})

	t.Run("call_mixed_named_positional", func(t *testing.T) {
		// Trino 481 accepts named-then-positional order.
		s := parseOneStmt(t, "CALL f(name => 1, 2)").(*CallStmt)
		if len(s.Arguments) != 2 {
			t.Fatalf("Arguments = %d, want 2", len(s.Arguments))
		}
		if s.Arguments[0].Name == nil {
			t.Error("argument 0 should be named")
		}
		if s.Arguments[1].Name != nil {
			t.Error("argument 1 should be positional")
		}
	})
}

// TestUtility_NodeTagsAndLoc verifies every node type returns its declared
// NodeTag and a valid Loc spanning the statement.
func TestUtility_NodeTagsAndLoc(t *testing.T) {
	cases := []struct {
		sql string
		tag ast.NodeTag
	}{
		{"SHOW TABLES", ast.T_ShowStmt},
		{"DESCRIBE t", ast.T_ShowStmt},
		{"USE s", ast.T_UseStmt},
		{"SET SESSION a = 1", ast.T_SetSessionStmt},
		{"RESET SESSION a", ast.T_ResetSessionStmt},
		{"SET SESSION AUTHORIZATION u", ast.T_SetSessionAuthorizationStmt},
		{"RESET SESSION AUTHORIZATION", ast.T_ResetSessionAuthorizationStmt},
		{"SET ROLE ALL", ast.T_SetRoleStmt},
		{"SET PATH a", ast.T_SetPathStmt},
		{"SET TIME ZONE LOCAL", ast.T_SetTimeZoneStmt},
		{"EXPLAIN SHOW TABLES", ast.T_ExplainStmt},
		{"CALL p()", ast.T_CallStmt},
	}
	for _, tc := range cases {
		t.Run(truncateName(tc.sql), func(t *testing.T) {
			stmt := parseOneStmt(t, tc.sql)
			if stmt.Tag() != tc.tag {
				t.Errorf("Tag() = %v, want %v", stmt.Tag(), tc.tag)
			}
			sp, ok := stmt.(interface{ Span() ast.Loc })
			if !ok {
				t.Fatal("statement does not expose Span()")
			}
			loc := sp.Span()
			if !loc.IsValid() || loc.Start != 0 || loc.End != len(tc.sql) {
				t.Errorf("Span() = %+v, want [0,%d)", loc, len(tc.sql))
			}
		})
	}
}
