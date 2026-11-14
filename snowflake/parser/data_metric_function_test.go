package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreateDataMetricFunction(t *testing.T) {
	in := "CREATE OR REPLACE SECURE DATA METRIC FUNCTION db.s.referential_check(" +
		"arg_t1 TABLE (arg_c1 INT), arg_t2 TABLE (arg_c2 INT, arg_c3 VARCHAR)) " +
		"RETURNS NUMBER NOT NULL LANGUAGE SQL COMMENT = 'c' " +
		"AS 'SELECT COUNT(*) FROM arg_t1 WHERE arg_c1 NOT IN (SELECT arg_c2 FROM arg_t2)'"
	node := mustParseOne(t, in)
	s, ok := node.(*ast.CreateRoutineStmt)
	if !ok {
		t.Fatalf("got %T", node)
	}
	if s.Kind != ast.RoutineDataMetricFunction || !s.OrReplace || !s.Secure || !s.ReturnNotNull || s.Loc.End != len(in) {
		t.Fatalf("got %+v", s)
	}
	if len(s.Args) != 2 || len(s.Args[0].Table) != 1 || len(s.Args[1].Table) != 2 || s.Args[0].Type != nil {
		t.Fatalf("args = %+v", s.Args)
	}
	if len(s.Options) != 2 || s.Options[0].Name != "LANGUAGE" || s.Options[0].Words != "SQL" || s.Options[1].Name != "COMMENT" {
		t.Fatalf("options = %+v", s.Options)
	}

	s = mustParseOne(t, "CREATE DATA METRIC FUNCTION IF NOT EXISTS f(t TABLE(c INT)) RETURNS NUMBER AS $$SELECT 1$$").(*ast.CreateRoutineStmt)
	if !s.IfNotExists || s.Secure || !s.BodyDollar || len(s.Args) != 1 {
		t.Fatalf("got %+v", s)
	}
}

func TestParseAlterDataMetricFunction(t *testing.T) {
	for in, want := range map[string]ast.AlterRoutineAction{
		"ALTER FUNCTION IF EXISTS f(TABLE(INT)) RENAME TO g":           ast.AlterRoutineRename,
		"ALTER FUNCTION f(TABLE(INT, VARCHAR), TABLE(INT)) SET SECURE": ast.AlterRoutineSet,
		"ALTER FUNCTION f(TABLE(INT)) UNSET SECURE":                    ast.AlterRoutineUnset,
		"ALTER FUNCTION f(TABLE(INT)) SET COMMENT = 'c'":               ast.AlterRoutineSet,
		"ALTER FUNCTION f(TABLE(INT)) UNSET COMMENT":                   ast.AlterRoutineUnset,
		"ALTER FUNCTION f(TABLE(INT)) SET TAG t = 's', u = 'v'":        ast.AlterRoutineSetTag,
		"ALTER FUNCTION f(TABLE(INT)) UNSET TAG t, u":                  ast.AlterRoutineUnsetTag,
		"ALTER FUNCTION f(INT) SET TAG t = 's', u = 'v'":               ast.AlterRoutineSetTag,
		"ALTER PROCEDURE p() UNSET TAG t, u":                           ast.AlterRoutineUnsetTag,
	} {
		node := mustParseOne(t, in)
		s, ok := node.(*ast.AlterRoutineStmt)
		if !ok || s.Action != want || s.Loc.End != len(in) {
			t.Errorf("%q: got %#v", in, node)
		}
	}
	s := mustParseOne(t, "ALTER FUNCTION f(TABLE(INT, VARCHAR), TABLE(INT)) SET SECURE").(*ast.AlterRoutineStmt)
	if len(s.TableArgs) != 2 || len(s.TableArgs[0].Types) != 2 || s.ArgTypes != nil {
		t.Errorf("table args = %+v", s.TableArgs)
	}
}

func TestParseDataMetricFunctionRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE DATA METRIC FUNCTION f() RETURNS NUMBER AS 's'",
		"CREATE DATA METRIC FUNCTION f(t INT) RETURNS NUMBER AS 's'",
		"CREATE DATA METRIC FUNCTION f(t TABLE()) RETURNS NUMBER AS 's'",
		"CREATE DATA METRIC FUNCTION f(a TABLE(c INT), b TABLE(c INT), d TABLE(c INT)) RETURNS NUMBER AS 's'",
		"CREATE DATA METRIC FUNCTION f(t TABLE(c INT)) RETURNS INT AS 's'",
		"CREATE DATA METRIC FUNCTION f(t TABLE(c INT)) RETURNS NUMBER(38, 0) AS 's'",
		"CREATE DATA METRIC FUNCTION f(t TABLE(c INT)) RETURNS NUMBER",
		"CREATE DATA METRIC FUNCTION f(t TABLE(c INT)) RETURNS NUMBER LANGUAGE PYTHON AS 's'",
		"CREATE DATA METRIC FUNCTION f(t TABLE(c INT)) RETURNS NUMBER COMMENT = 1 AS 's'",
		"CREATE DATA METRIC FUNCTION f(t TABLE(c INT)) RETURNS NUMBER IMMUTABLE AS 's'",
		"CREATE OR REPLACE DATA METRIC FUNCTION IF NOT EXISTS f(t TABLE(c INT)) RETURNS NUMBER AS 's'",
		"CREATE OR ALTER DATA METRIC FUNCTION f(t TABLE(c INT)) RETURNS NUMBER AS 's'",
		"ALTER FUNCTION f(TABLE()) SET SECURE",
		"ALTER FUNCTION f(TABLE(INT), TABLE(INT), TABLE(INT)) SET SECURE",
		"ALTER FUNCTION f(TABLE(INT), INT) SET SECURE",
		"ALTER FUNCTION f(TABLE(INT)) SET LOG_LEVEL = 'INFO'",
		"ALTER FUNCTION f(TABLE(INT)) SET SECURE COMMENT = 'c'",
		"ALTER FUNCTION f(TABLE(INT)) UNSET LOG_LEVEL",
		"ALTER FUNCTION f(TABLE(INT)) UNSET SECURE, COMMENT",
		"ALTER PROCEDURE p(TABLE(INT)) SET SECURE",
	} {
		mustRejectObject(t, in)
	}
}

func TestStartsCreateDataMetricFunction(t *testing.T) {
	for in, want := range map[string]bool{
		"DATA METRIC FUNCTION f":        true,
		"SECURE DATA METRIC FUNCTION f": true,
		"DATA METRIC f":                 false,
		"SECURE FUNCTION f":             false,
		"DATABASE d":                    false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		before := p.cur
		if got := p.startsCreateDataMetricFunction(); got != want {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
		if p.cur != before {
			t.Errorf("%q: lookahead moved the parser", in)
		}
	}
}
