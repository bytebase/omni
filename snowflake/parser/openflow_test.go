package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreateOpenflowDeployment(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OPENFLOW DEPLOYMENT IF NOT EXISTS d DEPLOYMENT_TYPE = BYOC VPC_TYPE = 'MANAGED' EVENT_TABLE = NONE")
	if s.Verb != "CREATE" || s.Kind != "OPENFLOW DEPLOYMENT" || !s.IfNotExists || len(s.Clauses) != 1 || len(s.Clauses[0].Options) != 3 {
		t.Fatalf("got %+v", s)
	}
	if s := mustObjectStmt(t, "CREATE OPENFLOW DEPLOYMENT d"); len(s.Clauses) != 0 {
		t.Fatalf("clauses = %q", clauseKeywords(s))
	}
}

func TestParseCreateOpenflowRuntime(t *testing.T) {
	in := "CREATE OPENFLOW RUNTIME IF NOT EXISTS db.sch.r IN DEPLOYMENT d NODE_TYPE = SMALL NODE_TYPE_TIER = 'S1' " +
		"MIN_NODES = 1 MAX_NODES = 2 EXECUTE_AS_ROLE = r1 EXTERNAL_ACCESS_INTEGRATIONS = (e1, e2) DISPLAY_NAME = 'n'"
	s := mustObjectStmt(t, in)
	if s.Kind != "OPENFLOW RUNTIME" || s.Name.String() != "db.sch.r" {
		t.Fatalf("got %+v", s)
	}
	if got := clauseKeywords(s); len(got) != 2 || got[0] != "IN DEPLOYMENT" || got[1] != "" {
		t.Fatalf("clauses = %q", got)
	}
	if n := s.Clauses[0].Names; len(n) != 1 || n[0].String() != "d" {
		t.Errorf("IN DEPLOYMENT names = %v", n)
	}
	if len(s.Clauses[1].Options) != 7 {
		t.Errorf("options = %d", len(s.Clauses[1].Options))
	}
}

func TestParseAlterOpenflow(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER OPENFLOW DEPLOYMENT IF EXISTS d UPGRADE":                         "UPGRADE",
		"ALTER OPENFLOW DEPLOYMENT d TERMINATE":                                 "TERMINATE",
		"ALTER OPENFLOW DEPLOYMENT d SET DISPLAY_NAME = 'n' EVENT_TABLE = NONE": "SET",
		"ALTER OPENFLOW DEPLOYMENT d UNSET DISPLAY_NAME, EVENT_TABLE":           "UNSET",
		"ALTER OPENFLOW DEPLOYMENT d RENAME TO e":                               "RENAME TO",
		"ALTER OPENFLOW RUNTIME r UPGRADE":                                      "UPGRADE",
		"ALTER OPENFLOW RUNTIME r UPGRADE RECOVERY":                             "UPGRADE RECOVERY",
		"ALTER OPENFLOW RUNTIME r UPGRADE FORCE":                                "UPGRADE FORCE",
		"ALTER OPENFLOW RUNTIME r SUSPEND":                                      "SUSPEND",
		"ALTER OPENFLOW RUNTIME r RESUME RECOVERY":                              "RESUME RECOVERY",
		"ALTER OPENFLOW RUNTIME r RESTART":                                      "RESTART",
		"ALTER OPENFLOW RUNTIME r TERMINATE CASCADE":                            "TERMINATE CASCADE",
		"ALTER OPENFLOW RUNTIME r ADD EXTERNAL_ACCESS_INTEGRATIONS = (e1, e2)":  "ADD",
		"ALTER OPENFLOW RUNTIME r REMOVE EXTERNAL_ACCESS_INTEGRATIONS = (e1)":   "REMOVE",
		"ALTER OPENFLOW RUNTIME r SET MIN_NODES = 1 NODE_TYPE_TIER = 'M4'":      "SET",
		"ALTER OPENFLOW RUNTIME r UNSET EXTERNAL_ACCESS_INTEGRATIONS":           "UNSET",
		"ALTER OPENFLOW RUNTIME r RENAME TO s":                                  "RENAME TO",
		"ALTER OPENFLOW DATA PLANE INTEGRATION i SET EVENT_TABLE = 'db.s.t'":    "SET",
	} {
		s := mustObjectStmt(t, in)
		if s.Verb != "ALTER" || len(s.Clauses) != 1 || s.Clauses[0].Keywords != want {
			t.Errorf("%q: clauses = %q, want [%q]", in, clauseKeywords(s), want)
		}
	}
	s := mustObjectStmt(t, "ALTER OPENFLOW DATA PLANE INTEGRATION i SET EVENT_TABLE = 'db.s.t'")
	if s.Kind != "OPENFLOW DATA PLANE INTEGRATION" {
		t.Errorf("Kind = %q", s.Kind)
	}
	s = mustObjectStmt(t, "ALTER OPENFLOW RUNTIME r ADD EXTERNAL_ACCESS_INTEGRATIONS = (e1, e2)")
	if o := s.Clauses[0].Options; len(o) != 1 || o[0].Name != "EXTERNAL_ACCESS_INTEGRATIONS" || len(o[0].Group) != 2 {
		t.Errorf("ADD options = %+v", o)
	}
}

func TestParseDropOpenflow(t *testing.T) {
	for in, want := range map[string]*ast.DropStmt{
		"DROP OPENFLOW DEPLOYMENT IF EXISTS d": {Object: "OPENFLOW DEPLOYMENT", IfExists: true},
		"DROP OPENFLOW RUNTIME r":              {Object: "OPENFLOW RUNTIME"},
		"DROP OPENFLOW RUNTIME r CASCADE":      {Object: "OPENFLOW RUNTIME", Cascade: true},
	} {
		node := mustParseOne(t, in)
		d, ok := node.(*ast.DropStmt)
		if !ok || d.Kind != ast.DropObject || d.Object != want.Object || d.IfExists != want.IfExists || d.Cascade != want.Cascade {
			t.Errorf("%q: got %#v", in, node)
			continue
		}
		if d.Loc.Start != 0 || d.Loc.End != len(in) {
			t.Errorf("%q: Loc = %+v", in, d.Loc)
		}
	}
}

func TestParseOpenflowRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE OR REPLACE OPENFLOW DEPLOYMENT d", // not documented
		"CREATE OPENFLOW DEPLOYMENT d UPGRADE",
		// IN DEPLOYMENT and EXECUTE_AS_ROLE are required.
		"CREATE OPENFLOW RUNTIME r NODE_TYPE = SMALL MIN_NODES = 1 MAX_NODES = 1 EXECUTE_AS_ROLE = x",
		"CREATE OPENFLOW RUNTIME r IN DEPLOYMENT d NODE_TYPE = SMALL MIN_NODES = 1 MAX_NODES = 1",
		// A runtime's name has one part or three; a deployment's has one.
		"CREATE OPENFLOW RUNTIME s.r IN DEPLOYMENT d NODE_TYPE = SMALL MIN_NODES = 1 MAX_NODES = 1 EXECUTE_AS_ROLE = x",
		"CREATE OPENFLOW RUNTIME r IN DEPLOYMENT s.d NODE_TYPE = SMALL MIN_NODES = 1 MAX_NODES = 1 EXECUTE_AS_ROLE = x",
		"ALTER OPENFLOW DEPLOYMENT d",
		"ALTER OPENFLOW DEPLOYMENT d SUSPEND",
		"ALTER OPENFLOW DEPLOYMENT d UPGRADE FORCE",
		"ALTER OPENFLOW RUNTIME r SUSPEND RECOVERY",
		"ALTER OPENFLOW RUNTIME r RESUME FORCE",
		"ALTER OPENFLOW RUNTIME r TERMINATE RECOVERY",
		"ALTER OPENFLOW RUNTIME r ADD",
		"ALTER OPENFLOW RUNTIME r ADD COMMENT = 'c'",
		"ALTER OPENFLOW RUNTIME r REMOVE EXTERNAL_ACCESS_INTEGRATIONS",
		"ALTER OPENFLOW DATA PLANE INTEGRATION i",
		"ALTER OPENFLOW DATA PLANE INTEGRATION i UNSET EVENT_TABLE",
		"CREATE OPENFLOW DATA PLANE INTEGRATION i",
		"DROP OPENFLOW DEPLOYMENT d CASCADE",
		"DROP OPENFLOW RUNTIME r RESTRICT",
		"DROP OPENFLOW DATA PLANE INTEGRATION i",
	} {
		mustRejectObject(t, in)
	}
}
