package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// The documented forms of CREATE PROCEDURE, ALTER RESOURCE MONITOR and DROP
// DATABASE ROLE omni missed.

func TestParseProcedureScopedTable(t *testing.T) {
	in := "CREATE OR REPLACE PROCEDURE SCOPED TEMPORARY TABLE db.s.t (a INT, \"B\" NUMBER(3, 1))"
	s, ok := mustParseOne(t, in).(*ast.CreateTableStmt)
	if !ok || !s.ProcedureScoped || !s.Temporary || !s.OrReplace || s.Name.String() != "db.s.t" || len(s.Columns) != 2 {
		t.Fatalf("got %+v", s)
	}
	if s.Columns[1].Name.Name != "B" || !s.Columns[1].Name.Quoted || s.Columns[1].DataType.Name != "NUMBER" {
		t.Errorf("column = %+v", s.Columns[1])
	}
	if s.Loc.Start != 0 || s.Loc.End != len(in) {
		t.Errorf("Loc = %+v", s.Loc)
	}
	mustParseOne(t, "CREATE PROCEDURE SCOPED TEMP TABLE t (a INT)")
	for _, bad := range []string{
		"CREATE PROCEDURE SCOPED TEMP TABLE t ()",               // at least one column
		"CREATE PROCEDURE SCOPED TEMP TABLE t (a INT,)",         // trailing comma
		"CREATE PROCEDURE SCOPED TEMP TABLE t (a INT NOT NULL)", // a column is a name and a type
		"CREATE PROCEDURE SCOPED TEMP TABLE t (a)",              // a type is required
		"CREATE PROCEDURE SCOPED TEMP TABLE t",                  // columns are required
		"CREATE PROCEDURE SCOPED TABLE t (a INT)",               // TEMP or TEMPORARY is required
		"CREATE PROCEDURE SCOPED TRANSIENT TABLE t (a INT)",     // TEMP or TEMPORARY
		"CREATE TEMP PROCEDURE SCOPED TEMP TABLE t (a INT)",     // no modifier before PROCEDURE
		"CREATE SECURE PROCEDURE SCOPED TEMP TABLE t (a INT)",   // no modifier before PROCEDURE
		"CREATE OR ALTER PROCEDURE SCOPED TEMP TABLE t (a INT)", // no OR ALTER
		"CREATE PROCEDURE SCOPED TEMP TABLE IF NOT EXISTS t (a INT)",
		"CREATE PROCEDURE SCOPED TEMP TABLE t (a INT) COMMENT = 'x'",
	} {
		mustRejectObject(t, bad)
	}
}

func TestStartsProcedureScopedTable(t *testing.T) {
	for in, want := range map[string]bool{
		"SCOPED TEMP TABLE":      true,
		"SCOPED TEMPORARY TABLE": true,
		"SCOPED TABLE":           false,
		"\"SCOPED\" TEMP TABLE":  false,
		"scoped(a INT)":          false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsProcedureScopedTable(); got != want || p.cur.Loc.Start != 0 {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
	}
}

func TestParseProcedureArgModesCopyGrantsRestrictedCaller(t *testing.T) {
	s, ok := mustParseOne(t, "CREATE PROCEDURE p(a IN INT, b INPUT INT, c OUT VARCHAR, d OUTPUT INT DEFAULT 1, e INT) "+
		"COPY GRANTS RETURNS INT LANGUAGE SQL EXECUTE AS RESTRICTED CALLER AS 'BEGIN RETURN 1; END'").(*ast.CreateRoutineStmt)
	if !ok || !s.CopyGrants || s.ExecuteAs != ast.ExecuteAsRestrictedCaller || len(s.Args) != 5 {
		t.Fatalf("got %+v", s)
	}
	for i, want := range []string{"IN", "INPUT", "OUT", "OUTPUT", ""} {
		if s.Args[i].Mode != want {
			t.Errorf("arg %d mode = %q, want %q", i, s.Args[i].Mode, want)
		}
	}
	if f, ok := mustParseOne(t, "CREATE FUNCTION f() COPY GRANTS RETURNS INT AS '1'").(*ast.CreateRoutineStmt); !ok || !f.CopyGrants {
		t.Errorf("CREATE FUNCTION COPY GRANTS: got %+v", f)
	}
	if a, ok := mustParseOne(t, "ALTER PROCEDURE p() EXECUTE AS RESTRICTED CALLER").(*ast.AlterRoutineStmt); !ok || a.ExecuteAs != ast.ExecuteAsRestrictedCaller {
		t.Errorf("ALTER PROCEDURE: got %+v", a)
	}
	for _, bad := range []string{
		"CREATE PROCEDURE p(a OUT INT) RETURNS INT LANGUAGE JAVASCRIPT AS 's'",            // modes are LANGUAGE SQL's
		"CREATE PROCEDURE p(a OUT INT) RETURNS INT LANGUAGE PYTHON HANDLER = 'h'",         // modes are LANGUAGE SQL's
		"CREATE FUNCTION f(a IN INT) RETURNS INT AS '1'",                                  // modes are a procedure's
		"CREATE PROCEDURE p(a IN OUT INT) RETURNS INT LANGUAGE SQL AS 's'",                // one mode
		"CREATE PROCEDURE p(a INOUT INT) RETURNS INT LANGUAGE SQL AS 's'",                 // not a mode
		"CREATE PROCEDURE p() COPY RETURNS INT LANGUAGE SQL AS 's'",                       // COPY GRANTS is whole
		"CREATE EXTERNAL FUNCTION f() COPY GRANTS RETURNS INT API_INTEGRATION = i AS 'u'", // not documented
		"CREATE PROCEDURE p() RETURNS INT LANGUAGE SQL EXECUTE AS RESTRICTED AS 's'",      // RESTRICTED CALLER is whole
		"ALTER PROCEDURE p() EXECUTE AS RESTRICTED OWNER",                                 // RESTRICTED CALLER is whole
	} {
		mustRejectObject(t, bad)
	}
}

func TestParseAlterResourceMonitorTriggers(t *testing.T) {
	s, ok := mustParseOne(t, "ALTER RESOURCE MONITOR IF EXISTS limiter TRIGGERS ON 80 PERCENT DO NOTIFY ON 100 PERCENT DO SUSPEND_IMMEDIATE").(*ast.AlterIntegrationStmt)
	if !ok || s.Kind != ast.ResourceMonitor || !s.IfExists || len(s.Triggers) != 2 || len(s.Options) != 0 {
		t.Fatalf("got %+v", s)
	}
	for _, bad := range []string{
		"ALTER RESOURCE MONITOR m TRIGGERS",                          // at least one trigger
		"ALTER RESOURCE MONITOR m TRIGGERS ON 1 PERCENT",             // DO is required
		"ALTER RESOURCE MONITOR m TRIGGERS ON 1 PERCENT DO NOTIFY x", // trailing tokens
		"ALTER SECRET s TRIGGERS ON 1 PERCENT DO NOTIFY",             // a resource monitor's alone
		"ALTER INTEGRATION i TRIGGERS ON 1 PERCENT DO NOTIFY",        // a resource monitor's alone
	} {
		mustRejectObject(t, bad)
	}
}

func TestParseDropDatabaseRole(t *testing.T) {
	d, ok := mustParseOne(t, "DROP DATABASE ROLE IF EXISTS d1.dr1").(*ast.DropStmt)
	if !ok || d.ObjectType() != "DATABASE ROLE" || !d.IfExists || d.Name.String() != "d1.dr1" {
		t.Fatalf("got %+v", d)
	}
	// ROLE before the statement's end, CASCADE or RESTRICT names a database.
	for _, in := range []string{"DROP DATABASE ROLE", "DROP DATABASE ROLE CASCADE", "DROP DATABASE IF EXISTS ROLE RESTRICT"} {
		if _, ok := mustParseOne(t, in).(*ast.DropDatabaseStmt); !ok {
			t.Errorf("%q: not a DROP DATABASE", in)
		}
	}
	for _, bad := range []string{
		"DROP DATABASE ROLE a.b.c",     // at most two parts
		"DROP DATABASE ROLE r CASCADE", // not documented
		"DROP DATABASE ROLE IF EXISTS", // a name is required
	} {
		mustRejectObject(t, bad)
	}
}

func TestStartsDropDatabaseRole(t *testing.T) {
	for in, want := range map[string]bool{
		"DATABASE ROLE r":       true,
		"DATABASE ROLE IF":      true,
		"DATABASE ROLE":         false,
		"DATABASE ROLE;":        false,
		"DATABASE ROLE CASCADE": false,
		"DATABASE d":            false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsDropDatabaseRole(); got != want || p.cur.Loc.Start != 0 {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
	}
}
