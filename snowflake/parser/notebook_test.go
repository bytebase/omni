package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreateNotebook(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OR REPLACE NOTEBOOK db.s.nb FROM '@db.s.stage' MAIN_FILE = 'nb.ipynb' "+
		"QUERY_WAREHOUSE = wh IDLE_AUTO_SHUTDOWN_TIME_SECONDS = 3600 COMPUTE_POOL = 'pool'")
	if s.Kind != "NOTEBOOK" || !s.OrReplace || s.Name.String() != "db.s.nb" {
		t.Fatalf("got %+v", s)
	}
	if got := clauseKeywords(s); len(got) != 2 || got[0] != "FROM" || len(s.Clauses[1].Options) != 4 {
		t.Fatalf("clauses = %q", got)
	}

	// A notebook may be named PROJECT.
	for _, in := range []string{
		"CREATE NOTEBOOK project",
		"CREATE NOTEBOOK PROJECT FROM '@s'",
		"CREATE NOTEBOOK PROJECT COMMENT = 'c'",
		"CREATE NOTEBOOK PROJECT.s.n FROM '@s'",
		"CREATE NOTEBOOK IF NOT EXISTS project",
		`CREATE NOTEBOOK "PROJECT"`,
	} {
		if s := mustObjectStmt(t, in); s.Kind != "NOTEBOOK" {
			t.Errorf("%q: Kind = %q", in, s.Kind)
		}
	}
}

func TestParseCreateNotebookProject(t *testing.T) {
	for in, ifNotExists := range map[string]bool{
		"CREATE NOTEBOOK PROJECT db.s.p FROM 'snow://workspace/w' COMMENT = 'c'": false,
		"CREATE NOTEBOOK PROJECT IF NOT EXISTS db.s.p FROM '@db.s.stage'":        true,
	} {
		s := mustObjectStmt(t, in)
		if s.Kind != "NOTEBOOK PROJECT" || s.IfNotExists != ifNotExists || s.Name.String() != "db.s.p" {
			t.Errorf("%q: got %+v", in, s)
		}
		if got := clauseKeywords(s); got[0] != "FROM" {
			t.Errorf("%q: clauses = %q", in, got)
		}
	}
}

func TestParseDropNotebook(t *testing.T) {
	d, ok := mustParseOne(t, "DROP NOTEBOOK db.s.nb").(*ast.DropStmt)
	if !ok || d.Object != "NOTEBOOK" || d.IfExists || d.Name.String() != "db.s.nb" {
		t.Fatalf("got %#v", d)
	}
}

func TestParseAlterNotebook(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER NOTEBOOK IF EXISTS a RENAME TO b":                          "RENAME TO",
		"ALTER NOTEBOOK a SET COMMENT = 'c' QUERY_WAREHOUSE = wh":         "SET",
		"ALTER NOTEBOOK a SET SECRETS = ('k' = db.s.secret, 'l' = other)": "SET",
		"ALTER NOTEBOOK a UNSET QUERY_WAREHOUSE, COMMENT":                 "UNSET",
		"ALTER NOTEBOOK a UNSET TAG t, u":                                 "UNSET TAG",
	} {
		s := mustObjectStmt(t, in)
		if len(s.Clauses) != 1 || s.Clauses[0].Keywords != want {
			t.Errorf("%q: clauses = %q, want [%q]", in, clauseKeywords(s), want)
		}
	}
}

func TestParseExecuteNotebook(t *testing.T) {
	for in, n := range map[string]int{
		"EXECUTE NOTEBOOK db.s.nb()":               0,
		"EXECUTE NOTEBOOK nb('a', 'b c')":          2,
		"EXECUTE NOTEBOOK nb(1, -2.5, TRUE, NULL)": 4,
		"EXECUTE NOTEBOOK project('x')":            1,
		"EXECUTE NOTEBOOK PROJECT.s.n()":           0,
		`EXECUTE NOTEBOOK "PROJECT" ()`:            0,
		"EXECUTE NOTEBOOK IDENTIFIER('nb')('x')":   1,
	} {
		s := mustObjectStmt(t, in)
		if s.Verb != "EXECUTE" || s.Kind != "NOTEBOOK" || len(s.Clauses) != 1 || len(s.Clauses[0].Options[0].List) != n {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	s := mustObjectStmt(t, "EXECUTE NOTEBOOK PROJECT p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' ARGUMENTS = :nb_args")
	if o := s.Clauses[0].Options[4]; o.Name != "ARGUMENTS" || o.Words != ":NB_ARGS" {
		t.Errorf("ARGUMENTS = %+v", o)
	}
	s = mustObjectStmt(t, "EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'nb.ipynb' COMPUTE_POOL = 'cp' "+
		"QUERY_WAREHOUSE = 'wh' RUNTIME = '2.2-CPU-PY3.11' ARGUMENTS = '-x 1' SECRETS = (db.s.k1, db.s.k2)")
	if s.Kind != "NOTEBOOK PROJECT" || s.Name.String() != "db.s.p" || len(s.Clauses) != 1 || len(s.Clauses[0].Options) != 6 {
		t.Fatalf("got %+v", s)
	}
}

func TestParseNotebookRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE NOTEBOOK a FROM @s",
		"CREATE NOTEBOOK a SET COMMENT = 'c'",
		"CREATE OR ALTER NOTEBOOK a",
		"CREATE NOTEBOOK PROJECT db.s.p",                     // FROM is required
		"CREATE NOTEBOOK PROJECT db.s.p FROM 's' MAIN_FILE",  //
		"CREATE OR REPLACE NOTEBOOK PROJECT db.s.p FROM 's'", // not documented
		"ALTER NOTEBOOK a",
		"ALTER NOTEBOOK a UNSET MAIN_FILE", // only QUERY_WAREHOUSE and COMMENT
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' ARGUMENTS = : v",
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' ARGUMENTS = v",
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' ARGUMENTS = :v ARGUMENTS = 'x'",
		"ALTER NOTEBOOK a SET SECRETS = ('k' = 1)",
		"ALTER NOTEBOOK a SET SECRETS = ('k' = a.b.c.d)",
		"ALTER NOTEBOOK a SET",
		"EXECUTE NOTEBOOK nb",
		"EXECUTE NOTEBOOK nb(x)",
		"EXECUTE NOTEBOOK nb(- 'a')",
		"EXECUTE NOTEBOOK nb(1 + 1)",
		"EXECUTE NOTEBOOK nb('a',)",
		"EXECUTE NOTEBOOK nb('a' 'b')",
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w'", // RUNTIME is required
		"EXECUTE NOTEBOOK PROJECT db.s.p",
		"EXECUTE NOTEBOOK",
		"DROP NOTEBOOK IF EXISTS a", // the page writes no IF EXISTS
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' RUNTIME = 'r'",
		// Properties the page does not document, or of another shape.
		"CREATE NOTEBOOK a SECRETS = ('k' = s)",                   // ALTER's alone
		"CREATE NOTEBOOK a IDLE_AUTO_SHUTDOWN_TIME_SECONDS = 'x'", //
		"ALTER NOTEBOOK a SET MAIN_FILE = 'm'",                    // CREATE's alone
		"CREATE NOTEBOOK PROJECT db.s.p FROM 's' MAIN_FILE = 'm'",
		"CREATE NOTEBOOK PROJECT db.s.p FROM 's' TAG (t = 'v')",
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' X = 1",
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = c QUERY_WAREHOUSE = 'w' RUNTIME = 'r'",
		"EXECUTE NOTEBOOK PROJECT db.s.p MAIN_FILE = 'm' COMPUTE_POOL = 'c' QUERY_WAREHOUSE = 'w' RUNTIME = 'r' SECRETS = ('s')",
	} {
		mustRejectObject(t, in)
	}
}
