package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func mustCreateDocPolicy(t *testing.T, input string) *ast.CreatePolicyStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreatePolicyStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreatePolicyStmt", input, node)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(input) {
		t.Errorf("parse %q: Loc = %+v, want {0 %d}", input, stmt.Loc, len(input))
	}
	return stmt
}

func mustAlterDocPolicy(t *testing.T, input string) *ast.AlterPolicyStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterPolicyStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterPolicyStmt", input, node)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(input) {
		t.Errorf("parse %q: Loc = %+v, want {0 %d}", input, stmt.Loc, len(input))
	}
	return stmt
}

func TestParseCreateBodyPolicies(t *testing.T) {
	for in, want := range map[string]ast.PolicyKind{
		"CREATE AGGREGATION POLICY IF NOT EXISTS db.s.p AS () RETURNS AGGREGATION_CONSTRAINT -> AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5) COMMENT = 'c'": ast.PolicyAggregation,
		"CREATE OR REPLACE JOIN POLICY p AS () RETURNS JOIN_CONSTRAINT -> JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE)":                                          ast.PolicyJoin,
		"CREATE PROJECTION POLICY p AS () RETURNS PROJECTION_CONSTRAINT -> PROJECTION_CONSTRAINT(ALLOW => false)":                                          ast.PolicyProjection,
		"CREATE PRIVACY POLICY p AS () RETURNS PRIVACY_BUDGET -> PRIVACY_BUDGET(BUDGET_NAME => 'b')":                                                       ast.PolicyPrivacy,
	} {
		s := mustCreateDocPolicy(t, in)
		if s.Kind != want || s.Body == nil || s.Returns == nil || s.Args != nil {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	s := mustCreateDocPolicy(t, "CREATE AGGREGATION POLICY IF NOT EXISTS db.s.p AS () RETURNS AGGREGATION_CONSTRAINT -> 1 COMMENT = 'c'")
	if !s.IfNotExists || s.Name.String() != "db.s.p" || s.Returns.Name != "AGGREGATION_CONSTRAINT" || len(s.Options) != 1 || s.Options[0].Name != "COMMENT" {
		t.Errorf("got %+v", s)
	}
}

func TestParseCreatePropertyPolicies(t *testing.T) {
	s := mustCreateDocPolicy(t, "CREATE PACKAGES POLICY p LANGUAGE PYTHON ALLOWLIST = ('numpy', 'pandas==1.2.3') BLOCKLIST = () ADDITIONAL_CREATION_BLOCKLIST = ('x') COMMENT = 'c'")
	if s.Kind != ast.PolicyPackages || s.Language != "PYTHON" || len(s.Options) != 4 {
		t.Errorf("packages: got %+v", s)
	}
	s = mustCreateDocPolicy(t, "CREATE OR REPLACE FEATURE POLICY p BLOCKED_OBJECT_TYPES_FOR_CREATION = (DATABASES, TASK) COMMENT = 'c' AS $$\nblocked_creation_rules: []\n$$")
	if s.Kind != ast.PolicyFeature || !s.OrReplace || len(s.Options) != 2 || s.Definition == nil {
		t.Errorf("feature: got %+v", s)
	}
	s = mustCreateDocPolicy(t, "CREATE FEATURE POLICY p")
	if s.Kind != ast.PolicyFeature || len(s.Options) != 0 || s.Definition != nil {
		t.Errorf("bare feature: got %+v", s)
	}
	s = mustCreateDocPolicy(t, "CREATE MAINTENANCE POLICY IF NOT EXISTS p SCHEDULE = 'USING CRON 0 3 * * * UTC' COMMENT = 'c'")
	if s.Kind != ast.PolicyMaintenance || !s.IfNotExists || len(s.Options) != 2 {
		t.Errorf("maintenance: got %+v", s)
	}
	s = mustCreateDocPolicy(t, "CREATE INGRESS POLICY p COMMENT = 'c' AS 'rules'")
	if s.Kind != ast.PolicyIngress || s.Definition == nil || *s.Definition != "rules" {
		t.Errorf("ingress: got %+v", s)
	}
}

func TestParseAlterDocPolicies(t *testing.T) {
	for in, want := range map[string]ast.AlterPolicyAction{
		"ALTER AGGREGATION POLICY IF EXISTS p RENAME TO q":                                       ast.AlterPolicyRename,
		"ALTER JOIN POLICY p SET BODY -> JOIN_CONSTRAINT(JOIN_REQUIRED => FALSE)":                ast.AlterPolicySetBody,
		"ALTER PROJECTION POLICY p SET TAG t = 's', u = 'v'":                                     ast.AlterPolicySetTag,
		"ALTER PRIVACY POLICY p UNSET TAG t, u":                                                  ast.AlterPolicyUnsetTag,
		"ALTER PRIVACY POLICY IF EXISTS p SET COMMENT = 'c'":                                     ast.AlterPolicySet,
		"ALTER PROJECTION POLICY IF EXISTS p UNSET COMMENT":                                      ast.AlterPolicyUnset,
		"ALTER PACKAGES POLICY p SET ALLOWLIST = () BLOCKLIST = ('a')":                           ast.AlterPolicySet,
		"ALTER PACKAGES POLICY IF EXISTS p UNSET ALLOWLIST, COMMENT":                             ast.AlterPolicyUnset,
		"ALTER FEATURE POLICY p SET BLOCKED_OBJECT_TYPES_FOR_CREATION = () COMMENT = 'c' AS 'y'": ast.AlterPolicySet,
		"ALTER FEATURE POLICY p SET AS 'y'":                                                      ast.AlterPolicySet,
		"ALTER FEATURE POLICY IF EXISTS p AS 'y'":                                                ast.AlterPolicyAs,
		"ALTER FEATURE POLICY p UNSET POLICY_DEFINITION, COMMENT":                                ast.AlterPolicyUnset,
		"ALTER FEATURE POLICY IF EXISTS p SET TAG t = 's'":                                       ast.AlterPolicySetTag,
		"ALTER FEATURE POLICY p RENAME TO q":                                                     ast.AlterPolicyRename,
		"ALTER MAINTENANCE POLICY p SET SCHEDULE = 's' COMMENT = 'c'":                            ast.AlterPolicySet,
		"ALTER MAINTENANCE POLICY p RENAME TO q":                                                 ast.AlterPolicyRename,
		"ALTER INGRESS POLICY p SET COMMENT = 'c' AS 'y'":                                        ast.AlterPolicySet,
		"ALTER INGRESS POLICY IF EXISTS p AS 'y'":                                                ast.AlterPolicyAs,
		"ALTER INGRESS POLICY p UNSET COMMENT":                                                   ast.AlterPolicyUnset,
	} {
		if s := mustAlterDocPolicy(t, in); s.Action != want {
			t.Errorf("%q: action = %v, want %v", in, s.Action, want)
		}
	}
}

func TestParseDropDocPolicies(t *testing.T) {
	for in, want := range map[string]string{
		"DROP AGGREGATION POLICY db.s.p":      "AGGREGATION POLICY",
		"DROP JOIN POLICY p":                  "JOIN POLICY",
		"DROP PROJECTION POLICY p":            "PROJECTION POLICY",
		"DROP PRIVACY POLICY p":               "PRIVACY POLICY",
		"DROP FEATURE POLICY p":               "FEATURE POLICY",
		"DROP PACKAGES POLICY IF EXISTS p":    "PACKAGES POLICY",
		"DROP MAINTENANCE POLICY IF EXISTS p": "MAINTENANCE POLICY",
		"DROP INGRESS POLICY IF EXISTS p":     "INGRESS POLICY",
	} {
		node := mustParseOne(t, in)
		d, ok := node.(*ast.DropStmt)
		if !ok || d.ObjectType() != want || d.Loc.End != len(in) {
			t.Errorf("%q: got %#v", in, node)
		}
	}
}

func TestParseDocPolicyRejects(t *testing.T) {
	for _, in := range []string{
		// The body policies' signature is empty and their RETURNS fixed.
		"CREATE AGGREGATION POLICY p AS (a INT) RETURNS AGGREGATION_CONSTRAINT -> 1",
		"CREATE AGGREGATION POLICY p AS () RETURNS JOIN_CONSTRAINT -> 1",
		"CREATE JOIN POLICY p AS () RETURNS \"JOIN_CONSTRAINT\" -> 1",
		"CREATE PRIVACY POLICY p AS () RETURNS PRIVACY_BUDGET -> 1 OTHER = 1",
		// IF NOT EXISTS excludes OR REPLACE; no CREATE OR ALTER.
		"CREATE OR REPLACE JOIN POLICY IF NOT EXISTS p AS () RETURNS JOIN_CONSTRAINT -> 1",
		"CREATE OR ALTER MAINTENANCE POLICY p SCHEDULE = 's'",
		// Undocumented properties and value shapes.
		"CREATE PACKAGES POLICY p ALLOWLIST = ('a')",
		"CREATE PACKAGES POLICY p LANGUAGE JAVA",
		"CREATE PACKAGES POLICY p LANGUAGE PYTHON ALLOWLIST = (a)",
		"CREATE PACKAGES POLICY p LANGUAGE PYTHON ALLOWLIST = ('a',)",
		"CREATE PACKAGES POLICY p LANGUAGE PYTHON ALLOWLIST = ('a') ALLOWLIST = ('b')",
		"CREATE FEATURE POLICY p BLOCKED_OBJECT_TYPES_FOR_CREATION = ()",
		"CREATE FEATURE POLICY p BLOCKED_OBJECT_TYPES_FOR_CREATION = ('TABLES')",
		"CREATE FEATURE POLICY p BLOCKED_OBJECT_TYPES_FOR_CREATION = (\"TABLES\")",
		"CREATE FEATURE POLICY p BLOCKED_OBJECT_TYPES_FOR_CREATION = (a.b)",
		"CREATE FEATURE POLICY p BLOCKED_OBJECT_TYPES_FOR_CREATION = (TABLE_FORMAT)",
		"ALTER FEATURE POLICY p SET BLOCKED_OBJECT_TYPES_FOR_CREATION = (zzz)",
		"CREATE FEATURE POLICY p AS 'y' COMMENT = 'c'",
		"CREATE MAINTENANCE POLICY p",
		"CREATE MAINTENANCE POLICY p COMMENT = 'c'",
		"CREATE MAINTENANCE POLICY p SCHEDULE = 1",
		"CREATE INGRESS POLICY p COMMENT = 'c'",
		"CREATE INGRESS POLICY p AS rules",
		// ALTER actions a page does not document.
		"ALTER JOIN POLICY IF EXISTS p SET TAG t = 's'",
		"ALTER JOIN POLICY IF EXISTS p UNSET TAG t",
		"ALTER JOIN POLICY p SET OTHER = 1",
		"ALTER JOIN POLICY p UNSET COMMENT, COMMENT",
		"ALTER JOIN POLICY p AS 'y'",
		"ALTER PACKAGES POLICY p RENAME TO q",
		"ALTER PACKAGES POLICY p SET TAG t = 's'",
		"ALTER PACKAGES POLICY p UNSET LANGUAGE",
		"ALTER PACKAGES POLICY p SET",
		"ALTER FEATURE POLICY p SET",
		"ALTER FEATURE POLICY p SET BODY -> 1",
		"ALTER MAINTENANCE POLICY p SET BODY -> 1",
		"ALTER MAINTENANCE POLICY p UNSET SCHEDULE",
		"ALTER INGRESS POLICY p SET AS 'y'",
		"ALTER INGRESS POLICY p RENAME TO q",
		"ALTER INGRESS POLICY p UNSET 'COMMENT'",
		// DROP's IF EXISTS only where the page writes it.
		"DROP AGGREGATION POLICY IF EXISTS p",
		"DROP FEATURE POLICY IF EXISTS p",
	} {
		mustRejectObject(t, in)
	}
}

func TestStartsDocPolicy(t *testing.T) {
	for in, want := range map[string]string{
		"AGGREGATION POLICY p": "AGGREGATION",
		"join policy p":        "JOIN",
		"PACKAGES POLICY":      "PACKAGES",
		"INGRESS POLICY":       "INGRESS",
		"MASKING POLICY p":     "",
		"FEATURE p":            "",
		"\"FEATURE\" POLICY":   "",
		"'JOIN' POLICY":        "",
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		before := p.cur
		got := ""
		if d := p.startsDocPolicy(); d != nil {
			got = d.word
		}
		if got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
		if p.cur != before {
			t.Errorf("%q: lookahead moved the parser", in)
		}
	}
}
