package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// ALTER SESSION { SET | UNSET } <session_parameter> ... (gap-alter-family)
// ---------------------------------------------------------------------------

func mustAlterSession(t *testing.T, input string) *ast.AlterSessionStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterSessionStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterSessionStmt", input, node)
	}
	return stmt
}

func TestAlterSession_Set(t *testing.T) {
	stmt := mustAlterSession(t, "ALTER SESSION SET LOCK_TIMEOUT = 3600")
	if stmt.Action != ast.AlterSessionSet {
		t.Fatalf("action = %v, want AlterSessionSet", stmt.Action)
	}
	if len(stmt.Options) != 1 {
		t.Fatalf("Options len = %d, want 1", len(stmt.Options))
	}
	opt := stmt.Options[0]
	if opt.Name != "LOCK_TIMEOUT" {
		t.Errorf("option name = %q, want LOCK_TIMEOUT", opt.Name)
	}
	if opt.Lit == nil || opt.Lit.Ival != 3600 {
		t.Errorf("option value = %+v, want int 3600", opt.Lit)
	}
}

func TestAlterSession_SetStringValue(t *testing.T) {
	stmt := mustAlterSession(t, "ALTER SESSION SET DEFAULT_NULL_ORDERING = 'LAST'")
	if stmt.Action != ast.AlterSessionSet {
		t.Fatalf("action = %v, want AlterSessionSet", stmt.Action)
	}
	opt := stmt.Options[0]
	if opt.Name != "DEFAULT_NULL_ORDERING" {
		t.Errorf("option name = %q, want DEFAULT_NULL_ORDERING", opt.Name)
	}
	if opt.Lit == nil || opt.Lit.Kind != ast.LitString || opt.Lit.Value != "LAST" {
		t.Errorf("option value = %+v, want string 'LAST'", opt.Lit)
	}
}

func TestAlterSession_SetNoSpaceWord(t *testing.T) {
	// ERROR_ON_NONDETERMINISTIC_UPDATE=TRUE (no surrounding spaces, word value).
	stmt := mustAlterSession(t, "ALTER SESSION SET ERROR_ON_NONDETERMINISTIC_UPDATE=TRUE")
	opt := stmt.Options[0]
	if opt.Name != "ERROR_ON_NONDETERMINISTIC_UPDATE" {
		t.Errorf("option name = %q", opt.Name)
	}
	if opt.Words != "TRUE" {
		t.Errorf("option words = %q, want TRUE", opt.Words)
	}
}

func TestAlterSession_SetMultiple(t *testing.T) {
	stmt := mustAlterSession(t, "ALTER SESSION SET AUTOCOMMIT = TRUE LOCK_TIMEOUT = 3600")
	if len(stmt.Options) != 2 {
		t.Fatalf("Options len = %d, want 2", len(stmt.Options))
	}
	if stmt.Options[0].Name != "AUTOCOMMIT" || stmt.Options[1].Name != "LOCK_TIMEOUT" {
		t.Errorf("option names = [%q %q], want [AUTOCOMMIT LOCK_TIMEOUT]",
			stmt.Options[0].Name, stmt.Options[1].Name)
	}
}

func TestAlterSession_Unset(t *testing.T) {
	stmt := mustAlterSession(t, "ALTER SESSION UNSET LOCK_TIMEOUT")
	if stmt.Action != ast.AlterSessionUnset {
		t.Fatalf("action = %v, want AlterSessionUnset", stmt.Action)
	}
	if len(stmt.UnsetKeys) != 1 || stmt.UnsetKeys[0] != "LOCK_TIMEOUT" {
		t.Errorf("UnsetKeys = %v, want [LOCK_TIMEOUT]", stmt.UnsetKeys)
	}
}

func TestAlterSession_UnsetMultiple(t *testing.T) {
	stmt := mustAlterSession(t, "ALTER SESSION UNSET LOCK_TIMEOUT, AUTOCOMMIT")
	if len(stmt.UnsetKeys) != 2 || stmt.UnsetKeys[0] != "LOCK_TIMEOUT" || stmt.UnsetKeys[1] != "AUTOCOMMIT" {
		t.Errorf("UnsetKeys = %v, want [LOCK_TIMEOUT AUTOCOMMIT]", stmt.UnsetKeys)
	}
}

func TestAlterSession_Loc(t *testing.T) {
	input := "ALTER SESSION SET LOCK_TIMEOUT = 3600"
	stmt := mustAlterSession(t, input)
	// Loc.Start anchors at the SESSION keyword (ALTER convention), not ALTER.
	if stmt.Loc.Start != len("ALTER ") {
		t.Errorf("Loc.Start = %d, want %d (SESSION keyword)", stmt.Loc.Start, len("ALTER "))
	}
	if int(stmt.Loc.End) != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

func TestAlterSession_Negatives(t *testing.T) {
	for _, in := range []string{
		"ALTER SESSION SET",          // SET with nothing settable
		"ALTER SESSION UNSET",        // UNSET with no keys
		"ALTER SESSION SET = 3600",   // missing parameter name
		"ALTER SESSION FROBNICATE x", // unknown action keyword
	} {
		if _, errs := parseSingle(in, 0); len(errs) == 0 {
			t.Errorf("expected parse error for %q, got none", in)
		}
	}
}

// TestAlterSessionPolicy_NotRegressed confirms that ALTER SESSION POLICY (a
// policy object) is still routed to the policy parser and NOT to the new
// session-parameter parser.
func TestAlterSessionPolicy_NotRegressed(t *testing.T) {
	node := mustParseOne(t, "ALTER SESSION POLICY sp RENAME TO sp2")
	if _, ok := node.(*ast.AlterSessionStmt); ok {
		t.Fatalf("ALTER SESSION POLICY mis-routed to *ast.AlterSessionStmt")
	}
	stmt, ok := node.(*ast.AlterPolicyStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.AlterPolicyStmt", node)
	}
	_ = stmt
}

// ---------------------------------------------------------------------------
// ALTER TABLE ... SEARCH OPTIMIZATION ON <method>(<cols>) [, ...]
// ---------------------------------------------------------------------------

func TestAlterTable_SearchOptMultiColumn(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION ON EQUALITY(c1, c2, c3)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableAddSearchOpt {
		t.Fatalf("kind = %v, want AlterTableAddSearchOpt", a.Kind)
	}
	if len(a.SearchOptOn) != 1 || a.SearchOptOn[0] != "EQUALITY(c1, c2, c3)" {
		t.Errorf("SearchOptOn = %v, want [EQUALITY(c1, c2, c3)]", a.SearchOptOn)
	}
}

func TestAlterTable_SearchOptMethodList(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t1 ADD SEARCH OPTIMIZATION ON EQUALITY(c1), EQUALITY(c2, c3)")
	a := requireAction(t, stmt, 0)
	if len(a.SearchOptOn) != 2 {
		t.Fatalf("SearchOptOn len = %d, want 2 (%v)", len(a.SearchOptOn), a.SearchOptOn)
	}
	if a.SearchOptOn[0] != "EQUALITY(c1)" || a.SearchOptOn[1] != "EQUALITY(c2, c3)" {
		t.Errorf("SearchOptOn = %v, want [EQUALITY(c1) EQUALITY(c2, c3)]", a.SearchOptOn)
	}
}

func TestAlterTable_SearchOptDropMethodList(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t1 DROP SEARCH OPTIMIZATION ON EQUALITY(c1, c2)")
	a := requireAction(t, stmt, 0)
	if a.Kind != ast.AlterTableDropSearchOpt {
		t.Fatalf("kind = %v, want AlterTableDropSearchOpt", a.Kind)
	}
	if len(a.SearchOptOn) != 1 || a.SearchOptOn[0] != "EQUALITY(c1, c2)" {
		t.Errorf("SearchOptOn = %v, want [EQUALITY(c1, c2)]", a.SearchOptOn)
	}
}

func TestAlterTable_SearchOptStar(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION ON EQUALITY(*)")
	a := requireAction(t, stmt, 0)
	if len(a.SearchOptOn) != 1 || a.SearchOptOn[0] != "EQUALITY(*)" {
		t.Errorf("SearchOptOn = %v, want [EQUALITY(*)]", a.SearchOptOn)
	}
}

func TestAlterTable_SearchOptNegatives(t *testing.T) {
	for _, in := range []string{
		"ALTER TABLE t ADD SEARCH OPTIMIZATION ON EQUALITY(",      // unterminated
		"ALTER TABLE t ADD SEARCH OPTIMIZATION ON EQUALITY(c1,",   // trailing comma, no col
		"ALTER TABLE t ADD SEARCH OPTIMIZATION ON EQUALITY(c1) ,", // trailing comma, no method
	} {
		if _, errs := parseSingle(in, 0); len(errs) == 0 {
			t.Errorf("expected parse error for %q, got none", in)
		}
	}
}

// ---------------------------------------------------------------------------
// ALTER VIEW ... SET/UNSET { JOIN | AGGREGATION } POLICY
// ---------------------------------------------------------------------------

func TestAlterView_SetJoinPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW join_view SET JOIN POLICY jp1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewSetJoinPolicy {
		t.Errorf("action = %v, want AlterViewSetJoinPolicy", stmt.Action)
	}
	if stmt.PolicyName == nil || stmt.PolicyName.Normalize() != "JP1" {
		t.Errorf("policy name = %v, want JP1", stmt.PolicyName)
	}
	if stmt.PolicyKeyKind != "" || len(stmt.PolicyKeyCols) != 0 {
		t.Errorf("unexpected key clause: kind=%q cols=%v", stmt.PolicyKeyKind, stmt.PolicyKeyCols)
	}
}

func TestAlterView_SetJoinPolicyAllowedKeys(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v SET JOIN POLICY jp1 ALLOWED JOIN KEYS (a, b)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.PolicyKeyKind != "ALLOWED" {
		t.Errorf("key kind = %q, want ALLOWED", stmt.PolicyKeyKind)
	}
	if len(stmt.PolicyKeyCols) != 2 {
		t.Errorf("key cols = %v, want 2", stmt.PolicyKeyCols)
	}
}

func TestAlterView_SetJoinPolicyEnforcedKeys(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v SET JOIN POLICY jp1 ENFORCED JOIN KEYS (a)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.PolicyKeyKind != "ENFORCED" {
		t.Errorf("key kind = %q, want ENFORCED", stmt.PolicyKeyKind)
	}
}

func TestAlterView_UnsetJoinPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v UNSET JOIN POLICY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewUnsetJoinPolicy {
		t.Errorf("action = %v, want AlterViewUnsetJoinPolicy", stmt.Action)
	}
}

func TestAlterView_SetAggregationPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v SET AGGREGATION POLICY ap1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewSetAggregationPolicy {
		t.Errorf("action = %v, want AlterViewSetAggregationPolicy", stmt.Action)
	}
	if stmt.PolicyName == nil || stmt.PolicyName.Normalize() != "AP1" {
		t.Errorf("policy name = %v, want AP1", stmt.PolicyName)
	}
}

func TestAlterView_SetAggregationPolicyEntityKeyForce(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v SET AGGREGATION POLICY ap1 ENTITY KEY (a, b) FORCE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.PolicyKeyCols) != 2 {
		t.Errorf("entity key cols = %v, want 2", stmt.PolicyKeyCols)
	}
	if !stmt.PolicyForce {
		t.Error("expected PolicyForce=true")
	}
}

func TestAlterView_UnsetAggregationPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v UNSET AGGREGATION POLICY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewUnsetAggregationPolicy {
		t.Errorf("action = %v, want AlterViewUnsetAggregationPolicy", stmt.Action)
	}
}

func TestAlterView_PolicyNegatives(t *testing.T) {
	for _, in := range []string{
		"ALTER VIEW v SET JOIN jp1",                // missing POLICY
		"ALTER VIEW v SET JOIN POLICY",             // missing name
		"ALTER VIEW v SET AGGREGATION jp1",         // missing POLICY
		"ALTER VIEW v SET JOIN POLICY jp1 ALLOWED", // ALLOWED with no JOIN KEYS (...)
		"ALTER VIEW v UNSET JOIN",                  // missing POLICY
	} {
		if _, errs := parseSingle(in, 0); len(errs) == 0 {
			t.Errorf("expected parse error for %q, got none", in)
		}
	}
}

// TestAlterView_ExistingFormsNotRegressed confirms the previously-supported
// ALTER VIEW forms still parse unchanged after adding the policy SET branches.
func TestAlterView_ExistingFormsNotRegressed(t *testing.T) {
	cases := []struct {
		in   string
		want ast.AlterViewAction
	}{
		{"ALTER VIEW v RENAME TO v2", ast.AlterViewRename},
		{"ALTER VIEW v SET SECURE", ast.AlterViewSetSecure},
		{"ALTER VIEW v UNSET SECURE", ast.AlterViewUnsetSecure},
		{"ALTER VIEW v SET COMMENT = 'x'", ast.AlterViewSetComment},
		{"ALTER VIEW v SET TAG (t = 'v')", ast.AlterViewSetTag},
		{"ALTER VIEW v ADD ROW ACCESS POLICY rap ON (c1)", ast.AlterViewAddRowAccessPolicy},
		{"ALTER VIEW v DROP ROW ACCESS POLICY rap", ast.AlterViewDropRowAccessPolicy},
	}
	for _, c := range cases {
		stmt, errs := testParseAlterView(c.in)
		if len(errs) > 0 {
			t.Errorf("unexpected errors for %q: %v", c.in, errs)
			continue
		}
		if stmt.Action != c.want {
			t.Errorf("%q: action = %v, want %v", c.in, stmt.Action, c.want)
		}
	}
}
