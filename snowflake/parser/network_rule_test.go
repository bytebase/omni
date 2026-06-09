package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func mustCreateNetworkRule(t *testing.T, input string) *ast.CreateNetworkRuleStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateNetworkRuleStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateNetworkRuleStmt", input, node)
	}
	return stmt
}

func mustAlterNetworkRule(t *testing.T, input string) *ast.AlterNetworkRuleStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterNetworkRuleStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterNetworkRuleStmt", input, node)
	}
	return stmt
}

func mustDropNetworkRule(t *testing.T, input string) *ast.DropStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.DropStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.DropStmt", input, node)
	}
	return stmt
}

// optList returns the parenthesized literal-list values of a CopyOption by
// name (e.g. VALUE_LIST = ('a', 'b')), or nil if the option is absent / not a
// list.
func optList(opts []*ast.CopyOption, name string) []*ast.Literal {
	o := findOption(opts, name)
	if o == nil {
		return nil
	}
	return o.List
}

// ---------------------------------------------------------------------------
// CREATE NETWORK RULE — corpus shapes
// ---------------------------------------------------------------------------

func TestParseCreateNetworkRule_Corpus(t *testing.T) {
	// example_01: TYPE = AWSVPCEID, parenthesized single-element VALUE_LIST,
	// MODE = INTERNAL_STAGE, COMMENT.
	t.Run("example_01", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, "CREATE NETWORK RULE corporate_network\n"+
			"  TYPE = AWSVPCEID\n"+
			"  VALUE_LIST = ('vpce-123abc3420c1931')\n"+
			"  MODE = INTERNAL_STAGE\n"+
			"  COMMENT = 'corporate privatelink endpoint'")
		if stmt.Name.String() != "corporate_network" {
			t.Errorf("Name = %q, want corporate_network", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.IfNotExists {
			t.Errorf("unexpected modifier set: %+v", stmt)
		}
		if w, ok := optWords(stmt.Options, "TYPE"); !ok || w != "AWSVPCEID" {
			t.Errorf("TYPE = %q (ok=%v), want AWSVPCEID", w, ok)
		}
		if w, ok := optWords(stmt.Options, "MODE"); !ok || w != "INTERNAL_STAGE" {
			t.Errorf("MODE = %q (ok=%v), want INTERNAL_STAGE", w, ok)
		}
		if c, ok := optLitString(stmt.Options, "COMMENT"); !ok || c != "corporate privatelink endpoint" {
			t.Errorf("COMMENT = %q (ok=%v)", c, ok)
		}
		vl := optList(stmt.Options, "VALUE_LIST")
		if len(vl) != 1 || vl[0].Value != "vpce-123abc3420c1931" {
			t.Errorf("VALUE_LIST = %+v, want single 'vpce-123abc3420c1931'", vl)
		}
	})

	// example_02: TYPE = IPV4, no MODE, COMMENT.
	t.Run("example_02", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, "CREATE NETWORK RULE cloud_network\n"+
			"  TYPE = IPV4\n"+
			"  VALUE_LIST = ('47.88.25.32/27')\n"+
			"  COMMENT = 'cloud egress ip range'")
		if w, ok := optWords(stmt.Options, "TYPE"); !ok || w != "IPV4" {
			t.Errorf("TYPE = %q (ok=%v), want IPV4", w, ok)
		}
		if findOption(stmt.Options, "MODE") != nil {
			t.Errorf("unexpected MODE option present")
		}
		vl := optList(stmt.Options, "VALUE_LIST")
		if len(vl) != 1 || vl[0].Value != "47.88.25.32/27" {
			t.Errorf("VALUE_LIST = %+v", vl)
		}
	})

	// example_03: property order MODE before VALUE_LIST; no COMMENT.
	t.Run("example_03_free_order", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, "CREATE NETWORK RULE gcp_rule\n"+
			"  TYPE = GCPPSCID\n"+
			"  MODE = INGRESS\n"+
			"  VALUE_LIST = ('31618973889077266')")
		if w, ok := optWords(stmt.Options, "TYPE"); !ok || w != "GCPPSCID" {
			t.Errorf("TYPE = %q (ok=%v), want GCPPSCID", w, ok)
		}
		if w, ok := optWords(stmt.Options, "MODE"); !ok || w != "INGRESS" {
			t.Errorf("MODE = %q (ok=%v), want INGRESS", w, ok)
		}
		if len(stmt.Options) != 3 {
			t.Errorf("expected 3 options, got %d: %+v", len(stmt.Options), stmt.Options)
		}
	})

	// example_04: HOST_PORT with a multi-element VALUE_LIST.
	t.Run("example_04_multi_value_list", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, "CREATE NETWORK RULE external_access_rule\n"+
			"  TYPE = HOST_PORT\n"+
			"  MODE = EGRESS\n"+
			"  VALUE_LIST = ('example.com', 'example.com:443')")
		vl := optList(stmt.Options, "VALUE_LIST")
		if len(vl) != 2 {
			t.Fatalf("VALUE_LIST len = %d, want 2: %+v", len(vl), vl)
		}
		if vl[0].Value != "example.com" || vl[1].Value != "example.com:443" {
			t.Errorf("VALUE_LIST = [%q, %q]", vl[0].Value, vl[1].Value)
		}
	})

	// example_05: CREATE OR REPLACE + 3-part qualified name + free property order.
	t.Run("example_05_or_replace_qualified", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, "CREATE OR REPLACE NETWORK RULE ext_network_access_db.network_rules.azure_sql_private_rule\n"+
			"  MODE = EGRESS\n"+
			"  TYPE = PRIVATE_HOST_PORT\n"+
			"  VALUE_LIST = ('externalaccessdemo.database.windows.net')")
		if !stmt.OrReplace {
			t.Errorf("OrReplace = false, want true")
		}
		if stmt.Name.Database.String() != "ext_network_access_db" ||
			stmt.Name.Schema.String() != "network_rules" ||
			stmt.Name.Name.String() != "azure_sql_private_rule" {
			t.Errorf("qualified name = %q (db=%q schema=%q name=%q)",
				stmt.Name.String(), stmt.Name.Database, stmt.Name.Schema, stmt.Name.Name)
		}
		if w, ok := optWords(stmt.Options, "TYPE"); !ok || w != "PRIVATE_HOST_PORT" {
			t.Errorf("TYPE = %q (ok=%v), want PRIVATE_HOST_PORT", w, ok)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE NETWORK RULE — modifiers, name, IF NOT EXISTS, Loc
// ---------------------------------------------------------------------------

func TestParseCreateNetworkRule_Modifiers(t *testing.T) {
	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, "CREATE NETWORK RULE IF NOT EXISTS r TYPE = IPV4 VALUE_LIST = ('0.0.0.0/0')")
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
		if stmt.OrReplace {
			t.Errorf("OrReplace = true, want false")
		}
		if stmt.Name.String() != "r" {
			t.Errorf("Name = %q, want r", stmt.Name.String())
		}
	})

	t.Run("or replace if not exists", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, "CREATE OR REPLACE NETWORK RULE r TYPE = IPV4 VALUE_LIST = ('0.0.0.0/0')")
		if !stmt.OrReplace {
			t.Errorf("OrReplace = false, want true")
		}
	})

	t.Run("quoted identifier name", func(t *testing.T) {
		stmt := mustCreateNetworkRule(t, `CREATE NETWORK RULE "My Rule" TYPE = IPV4 VALUE_LIST = ('1.2.3.4')`)
		if stmt.Name.Name.String() != `"My Rule"` {
			t.Errorf("Name = %q", stmt.Name.Name.String())
		}
	})
}

func TestParseCreateNetworkRule_Loc(t *testing.T) {
	input := "CREATE NETWORK RULE r TYPE = IPV4 VALUE_LIST = ('1.2.3.4')"
	stmt := mustCreateNetworkRule(t, input)
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if int(stmt.Loc.End) != len(input) {
		t.Errorf("Loc.End = %d, want %d (whole statement)", stmt.Loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// NETWORK RULE vs NETWORK POLICY disambiguation
// ---------------------------------------------------------------------------

func TestParseNetworkRule_NotPolicy(t *testing.T) {
	// CREATE NETWORK POLICY must still route to the policy parser, not the rule
	// parser — the kwNETWORK case must not swallow POLICY as a rule.
	node := mustParseOne(t, "CREATE NETWORK POLICY my_pol ALLOWED_IP_LIST = ('1.2.3.4')")
	if _, ok := node.(*ast.CreateNetworkRuleStmt); ok {
		t.Fatalf("CREATE NETWORK POLICY mis-routed to CreateNetworkRuleStmt")
	}
	if _, ok := node.(*ast.CreatePolicyStmt); !ok {
		t.Errorf("CREATE NETWORK POLICY = %T, want *ast.CreatePolicyStmt", node)
	}
}

// ---------------------------------------------------------------------------
// ALTER NETWORK RULE
// ---------------------------------------------------------------------------

func TestParseAlterNetworkRule_Set(t *testing.T) {
	t.Run("set value_list", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE r SET VALUE_LIST = ('1.2.3.4', '5.6.7.8')")
		if stmt.Action != ast.AlterNetworkRuleSet {
			t.Errorf("Action = %d, want Set", stmt.Action)
		}
		vl := optList(stmt.Options, "VALUE_LIST")
		if len(vl) != 2 {
			t.Errorf("VALUE_LIST len = %d, want 2", len(vl))
		}
	})

	t.Run("set value_list and comment", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE r SET VALUE_LIST = ('1.2.3.4') COMMENT = 'updated'")
		if c, ok := optLitString(stmt.Options, "COMMENT"); !ok || c != "updated" {
			t.Errorf("COMMENT = %q (ok=%v)", c, ok)
		}
	})

	t.Run("set comment only", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE r SET COMMENT = 'just a comment'")
		if len(stmt.Options) != 1 {
			t.Errorf("expected 1 option, got %d", len(stmt.Options))
		}
	})

	t.Run("if exists", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE IF EXISTS r SET VALUE_LIST = ('1.2.3.4')")
		if !stmt.IfExists {
			t.Errorf("IfExists = false, want true")
		}
	})

	t.Run("qualified name", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE db.sch.r SET COMMENT = 'x'")
		if stmt.Name.String() != "db.sch.r" {
			t.Errorf("Name = %q, want db.sch.r", stmt.Name.String())
		}
	})
}

func TestParseAlterNetworkRule_Unset(t *testing.T) {
	t.Run("unset value_list", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE r UNSET VALUE_LIST")
		if stmt.Action != ast.AlterNetworkRuleUnset {
			t.Errorf("Action = %d, want Unset", stmt.Action)
		}
		if len(stmt.UnsetKeys) != 1 || stmt.UnsetKeys[0] != "VALUE_LIST" {
			t.Errorf("UnsetKeys = %v, want [VALUE_LIST]", stmt.UnsetKeys)
		}
	})

	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE r UNSET COMMENT")
		if len(stmt.UnsetKeys) != 1 || stmt.UnsetKeys[0] != "COMMENT" {
			t.Errorf("UnsetKeys = %v, want [COMMENT]", stmt.UnsetKeys)
		}
	})

	t.Run("unset multiple keys", func(t *testing.T) {
		stmt := mustAlterNetworkRule(t, "ALTER NETWORK RULE r UNSET VALUE_LIST, COMMENT")
		if len(stmt.UnsetKeys) != 2 || stmt.UnsetKeys[0] != "VALUE_LIST" || stmt.UnsetKeys[1] != "COMMENT" {
			t.Errorf("UnsetKeys = %v, want [VALUE_LIST COMMENT]", stmt.UnsetKeys)
		}
	})

	t.Run("loc", func(t *testing.T) {
		input := "ALTER NETWORK RULE r UNSET COMMENT"
		stmt := mustAlterNetworkRule(t, input)
		// Loc.Start anchors at the NETWORK keyword (ALTER convention), not ALTER.
		if stmt.Loc.Start != int(len("ALTER ")) {
			t.Errorf("Loc.Start = %d, want %d (NETWORK keyword)", stmt.Loc.Start, len("ALTER "))
		}
		if int(stmt.Loc.End) != len(input) {
			t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
		}
	})
}

// ---------------------------------------------------------------------------
// DROP NETWORK RULE
// ---------------------------------------------------------------------------

func TestParseDropNetworkRule(t *testing.T) {
	t.Run("plain", func(t *testing.T) {
		stmt := mustDropNetworkRule(t, "DROP NETWORK RULE r")
		if stmt.Kind != ast.DropNetworkRule {
			t.Errorf("Kind = %v, want DropNetworkRule", stmt.Kind)
		}
		if stmt.IfExists {
			t.Errorf("IfExists = true, want false")
		}
		if stmt.Name.String() != "r" {
			t.Errorf("Name = %q, want r", stmt.Name.String())
		}
	})

	t.Run("if exists qualified", func(t *testing.T) {
		stmt := mustDropNetworkRule(t, "DROP NETWORK RULE IF EXISTS db.sch.r")
		if !stmt.IfExists {
			t.Errorf("IfExists = false, want true")
		}
		if stmt.Name.String() != "db.sch.r" {
			t.Errorf("Name = %q, want db.sch.r", stmt.Name.String())
		}
	})

	t.Run("kind string", func(t *testing.T) {
		if got := ast.DropNetworkRule.String(); got != "NETWORK RULE" {
			t.Errorf("DropNetworkRule.String() = %q, want %q", got, "NETWORK RULE")
		}
	})
}

// ---------------------------------------------------------------------------
// Negatives
// ---------------------------------------------------------------------------

func TestParseNetworkRule_Negatives(t *testing.T) {
	cases := []string{
		// Missing name after RULE.
		"CREATE NETWORK RULE",
		// ALTER with no action.
		"ALTER NETWORK RULE r",
		// ALTER SET with nothing settable.
		"ALTER NETWORK RULE r SET",
		// ALTER UNSET with no key.
		"ALTER NETWORK RULE r UNSET",
		// IF NOT EXISTS truncated.
		"CREATE NETWORK RULE IF NOT r TYPE = IPV4",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err == nil {
				t.Errorf("Parse(%q) = nil error, want a parse error", input)
			}
		})
	}
}
