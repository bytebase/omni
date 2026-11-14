package parser

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreatePostgresInstance(t *testing.T) {
	in := "CREATE POSTGRES INSTANCE IF NOT EXISTS pg1 COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 10 " +
		"AUTHENTICATION_AUTHORITY = EXTERNAL_OAUTH AUTHENTICATION_PROPERTIES = (ISSUER = 'i', AUDIENCE = 'a', " +
		"ACCEPTED_TOKENS = ((MAPPING_CLAIM = 'm', SCOPES_CLAIM = 's', SCOPES = ('x', 'y')), " +
		"(MAPPING_CLAIM = 'm', SCOPES_CLAIM = 's', SCOPES = 'z'))) POSTGRES_VERSION = 17 " +
		"HIGH_AVAILABILITY = TRUE COMMENT = 'c' WITH TAG (t = 'v')"
	s := mustObjectStmt(t, in)
	if s.Verb != "CREATE" || s.Kind != "POSTGRES INSTANCE" || !s.IfNotExists || s.Name.String() != "pg1" {
		t.Fatalf("got %+v", s)
	}
	if got := clauseKeywords(s); len(got) != 1 || got[0] != "" {
		t.Fatalf("clauses = %q", got)
	}
	props := s.Clauses[0]
	if len(props.Options) != 7 || len(props.Tags) != 1 {
		t.Fatalf("options = %d, tags = %d", len(props.Options), len(props.Tags))
	}
	auth := findOption(props.Options, "AUTHENTICATION_PROPERTIES")
	if auth == nil || len(auth.Group) != 3 || len(auth.Group[2].Group) != 2 || len(auth.Group[2].Group[0].Group[2].List) != 2 {
		t.Fatalf("AUTHENTICATION_PROPERTIES = %+v", auth)
	}
	if w, ok := optWords(props.Options, "AUTHENTICATION_AUTHORITY"); !ok || w != "EXTERNAL_OAUTH" {
		t.Errorf("AUTHENTICATION_AUTHORITY = %q", w)
	}
}

func TestParseCreatePostgresInstanceFork(t *testing.T) {
	for in, want := range map[string][]string{
		"CREATE POSTGRES INSTANCE a FORK b":                                              {"FORK"},
		"CREATE POSTGRES INSTANCE a FORK b AT (TIMESTAMP => '2025-01-01') COMMENT = 'c'": {"FORK", "AT", ""},
		"CREATE POSTGRES INSTANCE a FORK b BEFORE (OFFSET => -60) TAG (t = 'v')":         {"FORK", "BEFORE", ""},
	} {
		s := mustObjectStmt(t, in)
		if got := clauseKeywords(s); !slices.Equal(got, want) {
			t.Fatalf("%q: clauses = %q, want %q", in, got, want)
		}
		if n := s.Clauses[0].Names; len(n) != 1 || n[0].String() != "b" {
			t.Errorf("%q: FORK names = %v", in, n)
		}
	}
	s := mustObjectStmt(t, "CREATE POSTGRES INSTANCE a FORK b BEFORE (OFFSET => -60)")
	if o := s.Clauses[1].Options; len(o) != 1 || o[0].Name != "OFFSET" || o[0].Lit == nil || o[0].Lit.Value != "-60" {
		t.Errorf("BEFORE options = %+v", o)
	}
}

func TestParseAlterPostgresInstance(t *testing.T) {
	for in, want := range map[string][]string{
		"ALTER POSTGRES INSTANCE IF EXISTS a SUSPEND":                                  {"SUSPEND"},
		"ALTER POSTGRES INSTANCE a RESUME":                                             {"RESUME"},
		"ALTER POSTGRES INSTANCE a RENAME TO b":                                        {"RENAME TO"},
		"ALTER POSTGRES INSTANCE a RESET ACCESS FOR 'snowflake_admin'":                 {"RESET ACCESS FOR"},
		"ALTER POSTGRES INSTANCE a SET COMMENT = 'c'":                                  {"SET"},
		"ALTER POSTGRES INSTANCE a SET POSTGRES_VERSION = 18 APPLY IMMEDIATELY":        {"SET", "APPLY IMMEDIATELY"},
		"ALTER POSTGRES INSTANCE a SET POSTGRES_SETTINGS = '{}' APPLY ON '2025-01-01'": {"SET", "APPLY ON"},
		"ALTER POSTGRES INSTANCE a SET APPLY IMMEDIATELY":                              {"SET", "APPLY IMMEDIATELY"},
		"ALTER POSTGRES INSTANCE a UNSET COMMENT, NETWORK_POLICY":                      {"UNSET"},
		"ALTER POSTGRES INSTANCE a SET TAG t = 'v', u = 'w'":                           {"SET TAG"},
		"ALTER POSTGRES INSTANCE a UNSET TAG t, u":                                     {"UNSET TAG"},
	} {
		s := mustObjectStmt(t, in)
		if got := clauseKeywords(s); s.Verb != "ALTER" || !slices.Equal(got, want) {
			t.Fatalf("%q: clauses = %q, want %q", in, got, want)
		}
	}
	s := mustObjectStmt(t, "ALTER POSTGRES INSTANCE a RESET ACCESS FOR 'application'")
	if o := s.Clauses[0].Options; len(o) != 1 || o[0].Lit == nil || o[0].Lit.Value != "application" {
		t.Errorf("RESET ACCESS FOR options = %+v", o)
	}
	s = mustObjectStmt(t, "ALTER POSTGRES INSTANCE a SET APPLY ON '2025-01-01'")
	if o := s.Clauses[1].Options; len(o) != 1 || o[0].Lit == nil || o[0].Lit.Value != "2025-01-01" {
		t.Errorf("APPLY ON options = %+v", o)
	}
}

func TestParseDropPostgresInstance(t *testing.T) {
	node := mustParseOne(t, "DROP POSTGRES INSTANCE IF EXISTS a")
	d, ok := node.(*ast.DropStmt)
	if !ok || d.Object != "POSTGRES INSTANCE" || !d.IfExists || d.Name.String() != "a" {
		t.Fatalf("got %#v", node)
	}
}

func TestParsePostgresInstanceRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE POSTGRES INSTANCE a",
		"CREATE POSTGRES INSTANCE a COMPUTE_FAMILY = 's' STORAGE_SIZE_GB = 1", // AUTHENTICATION_AUTHORITY is required
		"CREATE OR REPLACE POSTGRES INSTANCE a FORK b",                        // not documented
		"CREATE POSTGRES INSTANCE a FORK",
		"CREATE POSTGRES INSTANCE a FORK b.c",
		"CREATE POSTGRES INSTANCE a FORK b AT (TIMESTAMP => 's'",
		"CREATE POSTGRES INSTANCE a FORK b AT (STATEMENT => 's')",
		"CREATE POSTGRES INSTANCE a FORK b AT (OFFSET => x)",
		"CREATE POSTGRES INSTANCE a FORK b AT OFFSET => 1",
		"ALTER POSTGRES INSTANCE a",
		"ALTER POSTGRES INSTANCE a SET",
		"ALTER POSTGRES INSTANCE a SET APPLY",
		"ALTER POSTGRES INSTANCE a SET APPLY ON",
		"ALTER POSTGRES INSTANCE a SET APPLY ON x",
		"ALTER POSTGRES INSTANCE a SET COMMENT = 'c' APPLY LATER",
		"ALTER POSTGRES INSTANCE a RESET ACCESS FOR snowflake_admin",
		"ALTER POSTGRES INSTANCE a RESET ACCESS",
		"ALTER POSTGRES INSTANCE a SUSPEND RESUME",
		"DROP POSTGRES INSTANCE a CASCADE",
	} {
		mustRejectObject(t, in)
	}
}
