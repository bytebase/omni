package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreateStreamlit(t *testing.T) {
	in := "CREATE OR REPLACE STREAMLIT db.s.app FROM '@db.s.stage/src' MAIN_FILE = 'main.py' " +
		"QUERY_WAREHOUSE = wh COMPUTE_POOL = pool TITLE = 't' IMPORTS = ('@s/a.py', '@s/b.py') " +
		"EXTERNAL_ACCESS_INTEGRATIONS = (i1, i2) SECRETS = ('k' = db.s.secret)"
	s := mustObjectStmt(t, in)
	if s.Verb != "CREATE" || s.Kind != "STREAMLIT" || !s.OrReplace || s.Name.String() != "db.s.app" {
		t.Fatalf("got %+v", s)
	}
	if got := clauseKeywords(s); len(got) != 2 || got[0] != "FROM" || got[1] != "" {
		t.Fatalf("clauses = %q", got)
	}
	if src := s.Clauses[0].Options; len(src) != 1 || src[0].Name != "" || src[0].Lit.Value != "@db.s.stage/src" {
		t.Errorf("FROM = %+v", src)
	}
	if n := len(s.Clauses[1].Options); n != 7 {
		t.Errorf("options = %d", n)
	}

	for _, in := range []string{
		"CREATE STREAMLIT a",
		"CREATE STREAMLIT IF NOT EXISTS a TITLE = 't'",
		"CREATE STREAMLIT a ROOT_LOCATION = '@s/app' MAIN_FILE = 'main.py' QUERY_WAREHOUSE = wh",
	} {
		mustObjectStmt(t, in)
	}
}

func TestParseAlterStreamlit(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER STREAMLIT IF EXISTS a SET MAIN_FILE = 'm.py' SECRETS = ('k' = s, 'l' = t)": "SET",
		"ALTER STREAMLIT a UNSET SECRETS, TITLE":                                          "UNSET",
		"ALTER STREAMLIT IF EXISTS a RENAME TO b":                                         "RENAME TO",
		"ALTER STREAMLIT a COMMIT":                                                        "COMMIT",
		"ALTER STREAMLIT a ABORT":                                                         "ABORT",
		"ALTER STREAMLIT a PULL":                                                          "PULL",
		"ALTER STREAMLIT a ADD LIVE VERSION FROM LAST":                                    "ADD LIVE VERSION FROM LAST",
		"ALTER STREAMLIT a PUSH":                                                          "PUSH",
	} {
		s := mustObjectStmt(t, in)
		if s.Verb != "ALTER" || len(s.Clauses) != 1 || s.Clauses[0].Keywords != want {
			t.Errorf("%q: clauses = %q, want [%q]", in, clauseKeywords(s), want)
		}
	}
	s := mustObjectStmt(t, "ALTER STREAMLIT a PUSH TO 'https://g/r/main' USERNAME = u PASSWORD = 'p' NAME = n EMAIL = e COMMENT = c")
	if opts := s.Clauses[0].Options; len(opts) != 6 || opts[0].Name != "TO" || opts[0].Lit.Value != "https://g/r/main" ||
		opts[5].Name != "COMMENT" {
		t.Errorf("PUSH options = %+v", opts)
	}
}

func TestParseDropStreamlit(t *testing.T) {
	node := mustParseOne(t, "DROP STREAMLIT IF EXISTS a")
	d, ok := node.(*ast.DropStmt)
	if !ok || d.Kind != ast.DropObject || d.Object != "STREAMLIT" || !d.IfExists || d.Name.String() != "a" {
		t.Fatalf("got %#v", node)
	}
}

func TestParseStreamlitRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE STREAMLIT a FROM @s",                                        // the source is a string
		"CREATE STREAMLIT a FROM",                                           //
		"CREATE STREAMLIT a ROOT_LOCATION = '@s/app'",                       // MAIN_FILE is required
		"CREATE STREAMLIT a FROM '@s' ROOT_LOCATION = '@s' MAIN_FILE = 'm'", // two forms mixed
		"CREATE STREAMLIT a MAIN_FILE",
		"CREATE OR ALTER STREAMLIT a",
		"ALTER STREAMLIT a",
		"ALTER STREAMLIT a SET",
		"ALTER STREAMLIT IF EXISTS a COMMIT", // IF EXISTS is not written before COMMIT
		"ALTER STREAMLIT IF EXISTS a PUSH",
		"ALTER STREAMLIT a ADD LIVE VERSION",
		"ALTER STREAMLIT a COMMIT ABORT",
		"ALTER STREAMLIT a PUSH TO",
		"ALTER STREAMLIT a PUSH TO (x)",
		"ALTER STREAMLIT a PUSH NAME",
		"DROP STREAMLIT a CASCADE",
		// Properties the page does not document, or of another shape.
		"CREATE STREAMLIT a NOT_A_PROPERTY = 1",
		"CREATE STREAMLIT a MAIN_FILE = main",
		"CREATE STREAMLIT a IMPORTS = ()",
		"CREATE STREAMLIT a EXTERNAL_ACCESS_INTEGRATIONS = ('i')",
		"CREATE STREAMLIT a SECRETS = (s)",
		"ALTER STREAMLIT a SET TITLE = 1",
		"ALTER STREAMLIT a UNSET MAIN_FILE",
		"ALTER STREAMLIT a PUSH BRANCH = b",
		"ALTER STREAMLIT a PUSH PASSWORD = p",
		"ALTER STREAMLIT a PUSH GIT_CREDENTIALS = 's'",
	} {
		mustRejectObject(t, in)
	}
}
