package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseDbtProject(t *testing.T) {
	s := mustObjectStmt(t, "CREATE DBT PROJECT IF NOT EXISTS db.sch.p FROM '@db.sch.git/branches/main' "+
		"DBT_VERSION = '1.9.4' DEFAULT_TARGET = 'prod' DEFAULT_ENVIRONMENT = 'dev' "+
		"EXTERNAL_ACCESS_INTEGRATIONS = (i1, i2) AUTO_COMPILE = FALSE DEFAULT_WRITEBACK = TRUE COMMENT = 'c'")
	if s.Verb != "CREATE" || s.Kind != "DBT PROJECT" || !s.IfNotExists || s.Name.String() != "db.sch.p" {
		t.Fatalf("got %+v", s)
	}
	if got := clauseKeywords(s); len(got) != 2 || got[0] != "FROM" || got[1] != "" {
		t.Fatalf("clauses = %q", got)
	}
	if n := len(s.Clauses[1].Options); n != 7 {
		t.Fatalf("options = %d", n)
	}
	mustObjectStmt(t, "CREATE OR REPLACE DBT PROJECT p FROM 'snow://dbt/db.sch.q/versions/live'")

	for in, want := range map[string]string{
		"ALTER DBT PROJECT IF EXISTS p RENAME TO db.sch.q":                   "RENAME TO",
		"ALTER DBT PROJECT p DEPLOY FROM '@s/p'":                             "DEPLOY FROM",
		"ALTER DBT PROJECT p SET DBT_VERSION = '1.9.4' AUTO_COMPILE = TRUE":  "SET",
		"ALTER DBT PROJECT p SET EXTERNAL_ACCESS_INTEGRATIONS = (i)":         "SET",
		"ALTER DBT PROJECT IF EXISTS p UNSET DBT_VERSION, DEFAULT_WRITEBACK": "UNSET",
	} {
		s := mustObjectStmt(t, in)
		if s.Verb != "ALTER" || len(s.Clauses) != 1 || s.Clauses[0].Keywords != want {
			t.Errorf("%q: clauses = %q, want [%q]", in, clauseKeywords(s), want)
		}
	}
	if d, ok := mustParseOne(t, "DROP DBT PROJECT IF EXISTS db.sch.p").(*ast.DropStmt); !ok || d.Object != "DBT PROJECT" || !d.IfExists {
		t.Errorf("DROP DBT PROJECT: got %+v", d)
	}

	for _, in := range []string{
		"CREATE DBT PROJECT p",                                   // FROM is required
		"CREATE DBT PROJECT p FROM @s/p",                         // FROM takes a string
		"CREATE OR REPLACE DBT PROJECT IF NOT EXISTS p FROM 's'", // the two exclude each other
		"CREATE OR ALTER DBT PROJECT p FROM 's'",                 // not documented
		"CREATE DBT PROJECT p FROM 's' DBT_VERSION = 1",          // a string
		"CREATE DBT PROJECT p FROM 's' AUTO_COMPILE = 'TRUE'",    // TRUE or FALSE
		"CREATE DBT PROJECT p FROM 's' EXTERNAL_ACCESS_INTEGRATIONS = (db.i)",
		"CREATE DBT PROJECT p FROM 's' EXTERNAL_ACCESS_INTEGRATIONS = ()",
		"CREATE DBT PROJECT p FROM 's' COMMENT = 'a' COMMENT = 'b'",
		"CREATE DBT PROJECT p FROM 's' WITH TAG (t = 'v')",
		"CREATE DBT PROJECT p FROM 's' NOT_A_PROPERTY = 1",
		"ALTER DBT PROJECT IF EXISTS p DEPLOY FROM 's'", // the page writes no IF EXISTS here
		"ALTER DBT PROJECT p SET",
		"ALTER DBT PROJECT p SET TAG t = 'v'",
		"ALTER DBT PROJECT p UNSET NOT_A_PROPERTY",
		"ALTER DBT PROJECT p SUSPEND",
		"ALTER DBT PROJECT p DEPLOY FROM",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseDcmProject(t *testing.T) {
	s := mustObjectStmt(t, "CREATE DCM PROJECT IF NOT EXISTS db.sch.p LOG_LEVEL = WARN COMMENT = 'c'")
	if s.Kind != "DCM PROJECT" || len(s.Clauses) != 1 || len(s.Clauses[0].Options) != 2 {
		t.Fatalf("got %+v", s)
	}
	if s := mustObjectStmt(t, "CREATE OR REPLACE DCM PROJECT p"); len(s.Clauses) != 0 {
		t.Errorf("clauses = %q", clauseKeywords(s))
	}
	for in, want := range map[string]string{
		"ALTER DCM PROJECT IF EXISTS p SET LOG_LEVEL = TRACE COMMENT = 'c'": "SET",
		"ALTER DCM PROJECT p UNSET LOG_LEVEL, COMMENT":                      "UNSET",
		"ALTER DCM PROJECT p DROP DEPLOYMENT DEPLOYMENT$1":                  "DROP DEPLOYMENT",
	} {
		s := mustObjectStmt(t, in)
		if len(s.Clauses) != 1 || s.Clauses[0].Keywords != want {
			t.Errorf("%q: clauses = %q, want [%q]", in, clauseKeywords(s), want)
		}
	}
	mustParseOne(t, "DROP DCM PROJECT p")

	for _, in := range []string{
		"CREATE DCM PROJECT p LOG_LEVEL = TRACE", // CREATE lists four levels
		"CREATE DCM PROJECT p LOG_LEVEL = 'INFO'",
		"CREATE DCM PROJECT p WITH TAG (t = 'v')",
		"ALTER DCM PROJECT p SET LOG_LEVEL = a",
		"ALTER DCM PROJECT IF EXISTS p DROP DEPLOYMENT d",
		"ALTER DCM PROJECT p DROP DEPLOYMENT db.d",
		"ALTER DCM PROJECT p UNSET LOG_LEVEL = INFO",
		"ALTER DCM PROJECT p RENAME TO q",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseCodeBundle(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OR REPLACE CODE BUNDLE db.sch.b FROM 'snow://workspace/user$.public.ws/versions/live' COMMENT = 'c'")
	if s.Kind != "CODE BUNDLE" || !s.OrReplace || len(s.Clauses) != 2 {
		t.Fatalf("got %+v", s)
	}
	mustObjectStmt(t, "CREATE CODE BUNDLE IF NOT EXISTS db.sch.b FROM '@db.sch.st'")
	for _, in := range []string{
		"CREATE CODE BUNDLE b FROM 's'",     // the name is fully qualified
		"CREATE CODE BUNDLE sch.b FROM 's'", // likewise
		"CREATE CODE BUNDLE db.sch.b",
		"CREATE CODE BUNDLE db.sch.b FROM 's' NOT_A_PROPERTY = 's'",
		"ALTER CODE BUNDLE db.sch.b SET COMMENT = 'c'", // not documented
		"DROP CODE BUNDLE db.sch.b",                    // not documented
	} {
		mustRejectObject(t, in)
	}
}

func mustExecuteStmt(t *testing.T, input, kind string) *ast.ObjectStmt {
	t.Helper()
	s := mustObjectStmt(t, input)
	if s.Verb != "EXECUTE" || s.Kind != kind {
		t.Fatalf("parse %q: Verb %q Kind %q", input, s.Verb, s.Kind)
	}
	return s
}

func TestParseExecuteDbtProject(t *testing.T) {
	s := mustExecuteStmt(t, "EXECUTE DBT PROJECT IF EXISTS db.sch.p ARGS = 'run --select m' DBT_VERSION = '1.9.4' "+
		"EXTERNAL_ACCESS_INTEGRATIONS = (i) ENVIRONMENT = 'prod' ENV_VARS = ('K' = 'v', 'L' = 'w') "+
		"IMPORTS = ('@db.sch.st/target/' AS 'state', SYSTEM$DBT_GET_LAST_SUCCESSFUL_RUN_TARGET('db.sch.p') AS 'prev', '@s/x') "+
		"WRITEBACK = FALSE", "DBT PROJECT")
	if !s.IfExists || s.Name.String() != "db.sch.p" || len(s.Clauses) != 1 || len(s.Clauses[0].Options) != 7 {
		t.Fatalf("got %+v", s)
	}
	imports := findOption(s.Clauses[0].Options, "IMPORTS")
	if len(imports.Group) != 3 || imports.Group[0].Alias == nil || imports.Group[1].Expr == nil || imports.Group[2].Alias != nil {
		t.Fatalf("IMPORTS = %+v", imports.Group)
	}

	for in, clauses := range map[string]int{
		"EXECUTE DBT PROJECT":                0,
		"EXECUTE DBT PROJECT IF EXISTS":      0,
		"EXECUTE DBT PROJECT ARGS = 'build'": 1,
		"EXECUTE DBT PROJECT FROM WORKSPACE \"USER$\".PUBLIC.\"ws\" PROJECT_ROOT = 'p'": 2,
		"EXECUTE DBT PROJECT p":  0,
		"EXECUTE DBT PROJECT if": 0, // a project named IF
	} {
		if s := mustExecuteStmt(t, in, "DBT PROJECT"); len(s.Clauses) != clauses {
			t.Errorf("%q: clauses = %q", in, clauseKeywords(s))
		}
	}

	for _, in := range []string{
		"EXECUTE DBT PROJECT p PROJECT_ROOT = 'x'",              // the workspace form's
		"EXECUTE DBT PROJECT FROM WORKSPACE w IMPORTS = ('s')",  // the named form's
		"EXECUTE DBT PROJECT FROM WORKSPACE w WRITEBACK = TRUE", // likewise
		"EXECUTE DBT PROJECT p ARGS = 'a' ARGS = 'b'",
		"EXECUTE DBT PROJECT p ARGS = run",
		"EXECUTE DBT PROJECT p IMPORTS = (a)",
		"EXECUTE DBT PROJECT p IMPORTS = (SYSTEM$TYPEOF('x'))",
		"EXECUTE DBT PROJECT p IMPORTS = (SYSTEM$LOCATE_DBT_ARCHIVE('q') || 'x')",
		"EXECUTE DBT PROJECT p IMPORTS = ('s' AS a)",
		"EXECUTE DBT PROJECT p IMPORTS = ()",
		"EXECUTE DBT PROJECT p ENV_VARS = (K = 'v')",
		"EXECUTE DBT PROJECT p ENV_VARS = ('K' = 1)",
		"EXECUTE DBT PROJECT p ENV_VARS = ('K' = 'v',)",
		"EXECUTE DBT PROJECT p EXTERNAL_ACCESS_INTEGRATIONS = (db.i)",
		"EXECUTE DBT PROJECT p WRITEBACK = 'TRUE'",
		"EXECUTE DBT PROJECT p NOT_A_PROPERTY = 's'",
		"EXECUTE DBT p",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseExecuteDcmProject(t *testing.T) {
	for in, want := range map[string][]string{
		"EXECUTE DCM PROJECT p PLAN FROM '@s/defs'":                                       {"PLAN", "FROM"},
		"EXECUTE DCM PROJECT p PLAN DELTA FROM DEPLOYMENT LAST OUTPUT_PATH '@s/out'":      {"PLAN DELTA", "FROM DEPLOYMENT", "OUTPUT_PATH"},
		"EXECUTE DCM PROJECT p PLAN USING (env => 'dev', regions => ['a', 'b']) FROM 's'": {"PLAN", "USING", "FROM"},
		"EXECUTE DCM PROJECT p DEPLOY AS \"my update\" USING CONFIGURATION PROD (desc => 'x') ENVIRONMENT (:env_json) FROM 's'": {
			"DEPLOY AS", "USING CONFIGURATION", "ENVIRONMENT", "FROM"},
		"EXECUTE DCM PROJECT p DEPLOY FROM DEPLOYMENT DEPLOYMENT$3":                                       {"DEPLOY", "FROM DEPLOYMENT"},
		"EXECUTE DCM PROJECT p DEPLOY ENVIRONMENT ($env) FROM 's'":                                        {"DEPLOY", "ENVIRONMENT", "FROM"},
		"EXECUTE DCM PROJECT p TEST ALL":                                                                  {"TEST ALL"},
		"EXECUTE DCM PROJECT p PREVIEW db.sch.t USING CONFIGURATION c ENVIRONMENT (:e) FROM 's' LIMIT 10": {"PREVIEW", "USING CONFIGURATION", "ENVIRONMENT", "FROM", "LIMIT"},
		"EXECUTE DCM PROJECT p PURGE":                                                                     {"PURGE"},
		"EXECUTE DCM PROJECT p PURGE AS \"sandbox clean\"":                                                {"PURGE AS"},
	} {
		s := mustExecuteStmt(t, in, "DCM PROJECT")
		got := clauseKeywords(s)
		if len(got) != len(want) {
			t.Errorf("%q: clauses = %q, want %q", in, got, want)
			continue
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("%q: clauses = %q, want %q", in, got, want)
				break
			}
		}
	}
	s := mustExecuteStmt(t, "EXECUTE DCM PROJECT p PLAN USING CONFIGURATION c (a => 1, b => 'x') FROM 's'", "DCM PROJECT")
	if using := s.Clauses[1]; len(using.Args) != 2 || using.Args[0].Name.Name != "a" || len(using.Names) != 1 {
		t.Fatalf("USING = %+v", using)
	}

	for _, in := range []string{
		"EXECUTE DCM PROJECT p",
		"EXECUTE DCM PROJECT p PLAN",
		"EXECUTE DCM PROJECT p PLAN USING FROM 's'", // USING with nothing
		"EXECUTE DCM PROJECT p PLAN USING (1) FROM 's'",
		"EXECUTE DCM PROJECT p PLAN USING () FROM 's'",
		"EXECUTE DCM PROJECT p PLAN USING CONFIGURATION c FROM DEPLOYMENT LAST",
		"EXECUTE DCM PROJECT p PLAN ENVIRONMENT (e) FROM 's'",
		"EXECUTE DCM PROJECT p PLAN ENVIRONMENT (: e) FROM 's'",
		"EXECUTE DCM PROJECT p PLAN ENVIRONMENT (:a, :b) FROM 's'",
		"EXECUTE DCM PROJECT p PLAN FROM @s/defs",
		"EXECUTE DCM PROJECT p PLAN FROM 's' OUTPUT_PATH = 'o'",
		"EXECUTE DCM PROJECT p DEPLOY FROM 's' OUTPUT_PATH 'o'",
		"EXECUTE DCM PROJECT p DEPLOY AS 'alias' FROM 's'",
		"EXECUTE DCM PROJECT p TEST",
		"EXECUTE DCM PROJECT p PREVIEW t FROM 's'",
		"EXECUTE DCM PROJECT p PREVIEW t USING CONFIGURATION c FROM 's' LIMIT x",
		"EXECUTE DCM PROJECT IF EXISTS p TEST ALL",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseExecuteCodeBundle(t *testing.T) {
	s := mustExecuteStmt(t, "EXECUTE CODE BUNDLE db.sch.b ENTRYPOINT = 'jobs/main.py' ARGUMENTS = ('--table', :source_table) "+
		"WITH SPECIFICATION $$\nbundle:\n  type: custom\n$$", "CODE BUNDLE")
	if got := clauseKeywords(s); len(got) != 2 || got[0] != "" || got[1] != "WITH SPECIFICATION" {
		t.Fatalf("clauses = %q", got)
	}
	args := findOption(s.Clauses[0].Options, "ARGUMENTS")
	if len(args.Group) != 2 || args.Group[1].Name != ":source_table" || !args.Group[1].Bare {
		t.Fatalf("ARGUMENTS = %+v", args.Group)
	}
	mustExecuteStmt(t, "EXECUTE CODE BUNDLE b ENTRYPOINT = 'main.py'", "CODE BUNDLE")

	for _, in := range []string{
		"EXECUTE CODE BUNDLE b",
		"EXECUTE CODE BUNDLE b ARGUMENTS = ('a') ENTRYPOINT = 'm'", // the page orders them
		"EXECUTE CODE BUNDLE b ENTRYPOINT = m",
		"EXECUTE CODE BUNDLE b ENTRYPOINT = 'm' ARGUMENTS = (a)",
		"EXECUTE CODE BUNDLE b ENTRYPOINT = 'm' ARGUMENTS = (1)",
		"EXECUTE CODE BUNDLE b ENTRYPOINT = 'm' ARGUMENTS = ()",
		"EXECUTE CODE BUNDLE b ENTRYPOINT = 'm' WITH SPECIFICATION",
		"EXECUTE CODE BUNDLE b ENTRYPOINT = 'm' COMMENT = 'c'",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseExecuteImmediateFrom(t *testing.T) {
	for in, want := range map[string]*ast.ExecuteImmediateStmt{
		"EXECUTE IMMEDIATE FROM @db.sch.st/scripts/setup.sql":                {File: "@db.sch.st/scripts/setup.sql"},
		"EXECUTE IMMEDIATE FROM './insert.sql'":                              {File: "'./insert.sql'"},
		"EXECUTE IMMEDIATE FROM @st/t.sql USING (env => 'dev', n => 0)":      {File: "@st/t.sql"},
		"EXECUTE IMMEDIATE FROM @st/t.sql DRY_RUN = TRUE":                    {File: "@st/t.sql", DryRun: "TRUE"},
		"EXECUTE IMMEDIATE FROM '../t.sql' USING (e => 'p') DRY_RUN = FALSE": {File: "'../t.sql'", DryRun: "FALSE"},
	} {
		s, ok := mustParseOne(t, in).(*ast.ExecuteImmediateStmt)
		if !ok || s.Source != ast.ExecImmFrom || s.File != want.File || s.DryRun != want.DryRun {
			t.Errorf("%q: got %+v", in, s)
			continue
		}
		if s.Loc.Start != 0 || s.Loc.End != len(in) {
			t.Errorf("%q: Loc = %+v", in, s.Loc)
		}
	}
	s := mustParseOne(t, "EXECUTE IMMEDIATE FROM @st/t.sql USING (env => 'dev', n => 0)").(*ast.ExecuteImmediateStmt)
	if len(s.Vars) != 2 || s.Vars[0].Name.Name != "env" {
		t.Fatalf("Vars = %+v", s.Vars)
	}
	// A stored procedure passes its arguments as :<name>.
	s = mustParseOne(t, "EXECUTE IMMEDIATE FROM @st/t.sql USING (env => :env, n => :\"N\")").(*ast.ExecuteImmediateStmt)
	if v, ok := s.Vars[1].Value.(*ast.BindVariable); !ok || v.Name.Name != "N" || !v.Name.Quoted || v.Loc.End != s.Vars[1].Loc.End {
		t.Fatalf("Vars[1] = %+v", s.Vars[1].Value)
	}
	// The page's example runs a hyphenated file name.
	s = mustParseOne(t, "EXECUTE IMMEDIATE FROM @my_stage/scripts/create-inventory.sql").(*ast.ExecuteImmediateStmt)
	if s.File != "@my_stage/scripts/create-inventory.sql" {
		t.Fatalf("File = %q", s.File)
	}

	for _, in := range []string{
		"EXECUTE IMMEDIATE FROM",
		"EXECUTE IMMEDIATE FROM st/t.sql",
		"EXECUTE IMMEDIATE FROM @",
		"EXECUTE IMMEDIATE FROM @ DRY_RUN = TRUE",
		"EXECUTE IMMEDIATE FROM @st/t.sql USING (env)", // template values are named
		"EXECUTE IMMEDIATE FROM @st/t.sql USING ()",
		"EXECUTE IMMEDIATE FROM @st/t.sql USING (a => : v)",
		"EXECUTE IMMEDIATE FROM @st/t.sql USING (a => :)",
		"EXECUTE IMMEDIATE FROM @st/a - b.sql",
		"EXECUTE IMMEDIATE FROM @st/t.sql DRY_RUN = 'TRUE'",
		"EXECUTE IMMEDIATE FROM @st/t.sql DRY_RUN = TRUE USING (a => 1)", // the page orders them
	} {
		mustRejectObject(t, in)
	}
}

func TestParseShowInProject(t *testing.T) {
	for _, in := range []string{
		"SHOW VERSIONS IN DBT PROJECT db.sch.p LIMIT 5",
		"SHOW DEPLOYMENTS IN DCM PROJECT p",
		"SHOW ENTITIES LIKE '%t%' IN DCM PROJECT p STARTS WITH 'T' LIMIT 3 FROM 'T1'",
		"SHOW GRANTS IN DCM PROJECT p LIMIT 1",
		"SHOW FUTURE GRANTS IN DCM PROJECT p",
	} {
		s, ok := mustParseOne(t, in).(*ast.ShowStmt)
		if !ok || s.Scope != ast.ShowScopeOther || s.ScopeName == nil || s.ScopeName.Name.Name == "" {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	for _, in := range []string{
		"SHOW VERSIONS IN DBT PROJECT p STARTS WITH 'a'",
		"SHOW DEPLOYMENTS LIKE 'a' IN DCM PROJECT p",
		"SHOW TERSE DEPLOYMENTS IN DCM PROJECT p",
		"SHOW TABLES IN DCM PROJECT p",
		"SHOW DEPLOYMENTS IN DBT PROJECT p",
		"SHOW GRANTS IN DBT PROJECT p",
		"SHOW GRANTS IN DCM PROJECT",
	} {
		mustRejectObject(t, in)
	}
}

func TestStartsExecuteProject(t *testing.T) {
	for in, want := range map[string]bool{
		"DBT PROJECT p":   true,
		"DCM PROJECT p":   true,
		"CODE BUNDLE b":   true,
		"\"DBT\" PROJECT": false,
		"DBT p":           false,
		"TASK t":          false,
		"IMMEDIATE 'x'":   false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsExecuteProject(); got != want {
			t.Errorf("startsExecuteProject(%q) = %v, want %v", in, got, want)
		}
	}
}
