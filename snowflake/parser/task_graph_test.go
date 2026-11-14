package parser

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreateTaskGraph(t *testing.T) {
	in := "CREATE OR REPLACE TASK GRAPH db.s.g WITH TAG (t = 'v') WITH CONTACT (STEWARD = c1, SUPPORT = c2) " +
		"SCHEDULE = 'USING CRON 0 * * * * UTC' ROOTS = (t_root) FINALIZER = db.s.t_fin OVERLAP_POLICY = NO_OVERLAP " +
		"CONFIG = $${\"a\": 1}$$ COMMENT = 'c' TASK_DEFAULTS (WAREHOUSE = wh, LOG_LEVEL = 'INFO') WHEN x > 1"
	s := mustObjectStmt(t, in)
	if s.Verb != "CREATE" || s.Kind != "TASK GRAPH" || !s.OrReplace || s.Name.String() != "db.s.g" {
		t.Fatalf("got %+v", s)
	}
	got := clauseKeywords(s)
	want := []string{"", "WITH CONTACT", "", "TASK_DEFAULTS", "WHEN"}
	if len(got) != len(want) {
		t.Fatalf("clauses = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("clauses = %q, want %q", got, want)
		}
	}
	if len(s.Clauses[1].Options) != 1 || len(s.Clauses[1].Options[0].Group) != 2 {
		t.Errorf("contacts = %+v", s.Clauses[1].Options)
	}
	if _, ok := s.Clauses[4].Query.(*ast.BinaryExpr); !ok {
		t.Errorf("WHEN = %T", s.Clauses[4].Query)
	}

	for _, in := range []string{
		"CREATE TASK GRAPH g",
		"CREATE TASK GRAPH IF NOT EXISTS g SUSPEND_TASK_AFTER_NUM_FAILURES = 3 TASK_AUTO_RETRY_ATTEMPTS = 2",
		"CREATE OR ALTER TASK GRAPH g SCHEDULE = '5 MINUTES' ROOTS = (t) ERROR_INTEGRATION = i SUCCESS_INTEGRATION = j",
		"CREATE TASK GRAPH g TAG (t = 'v') OVERLAP_POLICY = ALLOW_ALL_OVERLAP",
	} {
		mustObjectStmt(t, in)
	}
}

func TestParseCreateTaskGraphRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE TASK GRAPH g ROOTS = (a, b)",              // exactly one root
		"CREATE TASK GRAPH g ROOTS = a",                   // a parenthesized list
		"CREATE TASK GRAPH g COMMENT = 'x' COMMENT = 'y'", // a property once
		"CREATE TASK GRAPH g COMMENT = 'x' WITH CONTACT (s = c) COMMENT = 'y'",
		"CREATE TASK GRAPH g NOT_A_PROPERTY = 1",                // unknown property
		"CREATE TASK GRAPH g OVERLAP_POLICY = SOMETIMES",        // value outside the choices
		"CREATE TASK GRAPH g SCHEDULE = 5",                      // a string
		"CREATE TASK GRAPH g ERROR_INTEGRATION = db.i",          // one-part name
		"CREATE TASK GRAPH g TASK_DEFAULTS (NOT_A_DEFAULT = x)", // listed defaults only
		"CREATE TASK GRAPH g TASK_DEFAULTS (WAREHOUSE = 'w')",   // a name
		"CREATE TASK GRAPH g TASK_DEFAULTS (LOG_LEVEL = 'a', LOG_LEVEL = 'b')",
		"CREATE TASK GRAPH g TASK_DEFAULTS ()",
		"CREATE TASK GRAPH g TASK_DEFAULTS (WAREHOUSE = w,)",
		"CREATE TASK GRAPH g TASK_DEFAULTS = (WAREHOUSE = w)",
		"CREATE TASK GRAPH g WITH CONTACT ()",
		"CREATE TASK GRAPH g WITH CONTACT (s = 'c')",
		"CREATE TASK GRAPH g WITH CONTACT (s)",
		"CREATE TASK GRAPH g WHEN x > 1 COMMENT = 'c'", // WHEN is last
		"CREATE TASK GRAPH g TAG (t = 'v') TAG (u = 'v')",
		"CREATE OR ALTER TASK GRAPH g WITH CONTACT (s = c)", // not in CREATE OR ALTER
		"CREATE OR ALTER TASK GRAPH g TAG (t = 'v')",
		"CREATE OR ALTER TASK GRAPH IF NOT EXISTS g",
		"CREATE OR REPLACE TASK GRAPH IF NOT EXISTS g",
		"CREATE TASK GRAPH 'COMMENT' = 'x'",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAlterTaskGraph(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER TASK GRAPH g COMMIT":                                     "COMMIT",
		"ALTER TASK GRAPH IF EXISTS g REVERT":                           "REVERT",
		"ALTER TASK GRAPH g RESUME":                                     "RESUME",
		"ALTER TASK GRAPH g SUSPEND":                                    "SUSPEND",
		"ALTER TASK GRAPH g MODIFY WHEN SYSTEM$STREAM_HAS_DATA('s')":    "MODIFY WHEN",
		"ALTER TASK GRAPH g REMOVE WHEN":                                "REMOVE WHEN",
		"ALTER TASK GRAPH g SET SCHEDULE = '30 SECONDS', COMMENT = 'c'": "SET",
		"ALTER TASK GRAPH g UNSET CONFIG, FINALIZER":                    "UNSET",
		"ALTER TASK GRAPH g SET TASK_DEFAULTS (WAREHOUSE = new_wh)":     "SET TASK_DEFAULTS",
		"ALTER TASK GRAPH g UNSET TASK_DEFAULTS (LOG_LEVEL, COMMENT)":   "UNSET TASK_DEFAULTS",
		"ALTER TASK GRAPH g SET TAG t = 'v', u = 'w'":                   "SET TAG",
		"ALTER TASK GRAPH g UNSET TAG t, u":                             "UNSET TAG",
		"ALTER TASK GRAPH g SET CONTACT STEWARD = c, SUPPORT = d":       "SET CONTACT",
		"ALTER TASK GRAPH g UNSET CONTACT STEWARD, SUPPORT":             "UNSET CONTACT",
	} {
		s := mustObjectStmt(t, in)
		if s.Verb != "ALTER" || s.Kind != "TASK GRAPH" {
			t.Errorf("parse %q: got %+v", in, s)
		}
		if got := clauseKeywords(s); len(got) != 1 || got[0] != want {
			t.Errorf("parse %q: clauses = %q, want [%q]", in, got, want)
		}
	}
	for _, in := range []string{
		"ALTER TASK GRAPH g SET SCHEDULE = '30 SECONDS' COMMENT = 'c'", // the page writes commas
		"ALTER TASK GRAPH g SET SCHEDULE = 's',",
		"ALTER TASK GRAPH g SET NOT_A_PROPERTY = 1",
		"ALTER TASK GRAPH g SET COMMENT = 'a', COMMENT = 'b'",
		"ALTER TASK GRAPH g UNSET NOT_A_PROPERTY",
		"ALTER TASK GRAPH g UNSET TASK_DEFAULTS",
		"ALTER TASK GRAPH g UNSET TASK_DEFAULTS (WAREHOUSE = w)",
		"ALTER TASK GRAPH g SET CONTACT = c",
		"ALTER TASK GRAPH g UNSET CONTACT",
		"ALTER TASK GRAPH g COMMIT REVERT",
		"ALTER TASK GRAPH g MODIFY WHEN",
	} {
		mustRejectObject(t, in)
	}
}

// A task may be named GRAPH: those statements still read as a task's.
func TestParseTaskNamedGraph(t *testing.T) {
	for in, want := range map[string]string{
		"CREATE TASK GRAPH AS SELECT 1":                "*ast.CreateTaskStmt",
		"CREATE TASK GRAPH WAREHOUSE = w AS SELECT 1":  "*ast.CreateTaskStmt",
		"ALTER TASK GRAPH REMOVE WHEN":                 "*ast.AlterTaskStmt",
		"ALTER TASK GRAPH SUSPEND":                     "*ast.AlterTaskStmt",
		"DROP TASK GRAPH":                              "*ast.DropStmt",
		"EXECUTE TASK GRAPH RETRY LAST":                "*ast.ExecuteTaskStmt",
		"EXECUTE TASK GRAPH RETRY GRAPH RUN GROUP 'g'": "*ast.ExecuteTaskStmt",
		"EXECUTE TASK GRAPH AWAIT":                     "*ast.ExecuteTaskStmt",
		"CREATE TASK GRAPH g":                          "*ast.ObjectStmt",
		"ALTER TASK GRAPH g SUSPEND":                   "*ast.ObjectStmt",
	} {
		node := mustParseOne(t, in)
		if got := fmt.Sprintf("%T", node); got != want {
			t.Errorf("parse %q: got %s, want %s", in, got, want)
		}
	}
	if s := mustParseOne(t, "EXECUTE TASK GRAPH RETRY LAST").(*ast.ExecuteTaskStmt); s.Graph || s.Name.String() != "GRAPH" || !s.RetryLast {
		t.Errorf("got %+v", s)
	}
	if s := mustParseOne(t, "DROP TASK GRAPH").(*ast.DropStmt); s.Kind != ast.DropTask {
		t.Errorf("got %+v", s)
	}
}

func TestParseDropExecuteTaskGraph(t *testing.T) {
	d := mustParseOne(t, "DROP TASK GRAPH IF EXISTS db.s.g").(*ast.DropStmt)
	if d.Kind != ast.DropObject || d.ObjectType() != "TASK GRAPH" || !d.IfExists || d.Name.String() != "db.s.g" {
		t.Errorf("got %+v", d)
	}
	for in, check := range map[string]func(*ast.ExecuteTaskStmt) bool{
		"EXECUTE TASK GRAPH g":                            func(s *ast.ExecuteTaskStmt) bool { return s.Graph },
		"EXECUTE TASK GRAPH g USING CONFIG = '{}'":        func(s *ast.ExecuteTaskStmt) bool { return s.Graph && s.UsingConfig == "'{}'" },
		"EXECUTE TASK GRAPH g RETRY LAST":                 func(s *ast.ExecuteTaskStmt) bool { return s.Graph && s.RetryLast },
		"EXECUTE TASK GRAPH g RETRY GRAPH RUN GROUP 'r1'": func(s *ast.ExecuteTaskStmt) bool { return s.Graph && s.RetryGroup == "r1" },
		"EXECUTE TASK t RETRY GRAPH RUN GROUP 'r1'":       func(s *ast.ExecuteTaskStmt) bool { return !s.Graph && s.RetryGroup == "r1" },
		"EXECUTE TASK t AWAIT":                            func(s *ast.ExecuteTaskStmt) bool { return s.Await },
		"EXECUTE TASK t USING CONFIG = '{}' AWAIT":        func(s *ast.ExecuteTaskStmt) bool { return s.Await && s.UsingConfig != "" },
	} {
		s := mustExecTask(t, in)
		if !check(s) || s.Loc.Start != 0 || s.Loc.End != len(in) {
			t.Errorf("parse %q: got %+v", in, s)
		}
	}
	for _, in := range []string{
		"EXECUTE TASK GRAPH g AWAIT",
		"EXECUTE TASK t RETRY LAST AWAIT",
		"EXECUTE TASK t RETRY GRAPH RUN GROUP r1",
		"EXECUTE TASK GRAPH g RETRY",
		"DROP TASK GRAPH g h",
	} {
		mustRejectObject(t, in)
	}
}
