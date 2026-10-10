package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestPLSQLTerminatorsRequired checks that every PL/SQL statement,
// declaration, package item, and END carries its ';', and that a stray ';'
// is an error, as on Oracle 23ai (PLS-00103; reference rows ref_108 to
// ref_118). at is the text from the error onward, "" for the end of input.
func TestPLSQLTerminatorsRequired(t *testing.T) {
	cases := []struct {
		sql string
		at  string
	}{
		{"BEGIN NULL END;", "END;"},
		{"BEGIN x := 1 y := 2; END;", "y := 2; END;"},
		{"BEGIN IF 1 = 1 THEN NULL; END IF END;", "END;"},
		{"BEGIN LOOP NULL; END LOOP END;", "END;"},
		{"BEGIN CASE 1 WHEN 1 THEN NULL; END CASE END;", "END;"},
		{"BEGIN BEGIN NULL; END END;", "END;"},
		{"BEGIN BEGIN NULL; END IF; END;", "IF; END;"},
		{"BEGIN <<l>> BEGIN NULL; END l; END", ""},
		{"BEGIN NULL; END", ""},
		{"BEGIN NULL; END lbl", ""},
		{"BEGIN NULL;; END;", "; END;"},
		{"BEGIN COMMIT END;", "END;"},
		{"BEGIN PRAGMA AUTONOMOUS_TRANSACTION; NULL; END;", "AUTONOMOUS_TRANSACTION; NULL; END;"},
		{"BEGIN SET ROLE ALL; END;", "SET ROLE ALL; END;"},
		{"DECLARE x NUMBER y NUMBER; BEGIN NULL; END;", "y NUMBER; BEGIN NULL; END;"},
		{"DECLARE x NUMBER;; BEGIN NULL; END;", "; BEGIN NULL; END;"},
		{"CREATE PROCEDURE p AS BEGIN NULL; END", ""},
		{"CREATE PACKAGE k AS PROCEDURE q; END", ""},
		{"CREATE PACKAGE k AS PROCEDURE q END;", "END;"},
		{"CREATE PACKAGE k AS c CONSTANT PLS_INTEGER := 1;; END;", "; END;"},
		{"CREATE PACKAGE k AS ; c CONSTANT PLS_INTEGER := 1; END;", "; c CONSTANT PLS_INTEGER := 1; END;"},
		{"CREATE PACKAGE k AS c CONSTANT PLS_INTEGER := 1; 42; END;", "42; END;"},
		{"CREATE PACKAGE BODY k AS PROCEDURE q IS BEGIN NULL; END; END", ""},
		{"CREATE TYPE BODY ty AS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END", ""},
		{"CREATE TRIGGER trg FOR INSERT ON t COMPOUND TRIGGER BEFORE STATEMENT IS BEGIN NULL; END BEFORE STATEMENT; END trg", ""},
		{"BEGIN CASE WHEN 1 = 1 THEN NULL x := 1; END CASE; END;", "x := 1; END CASE; END;"},
		{"BEGIN CASE 1 WHEN 1 THEN NULL WHEN 2 THEN NULL; END CASE; END;", "WHEN 2 THEN NULL; END CASE; END;"},
		{"BEGIN CASE 1 WHEN 1 THEN NULL; END; END;", "; END;"},
		{"BEGIN CASE 1 ELSE NULL; END CASE; END;", "ELSE NULL; END CASE; END;"},
		{"BEGIN IF 1 = 1 THEN NULL; END; END;", "; END;"},
		{"BEGIN IF 1 = 1 THEN NULL; END IF IF; END;", "IF; END;"},
		{"BEGIN END;", "END;"},
		{"BEGIN IF 1 = 1 THEN END IF; END;", "END IF; END;"},
		{"BEGIN IF 1 = 1 THEN NULL; ELSE END IF; END;", "END IF; END;"},
		{"BEGIN CASE 1 WHEN 1 THEN WHEN 2 THEN NULL; END CASE; END;", "WHEN 2 THEN NULL; END CASE; END;"},
		{"BEGIN CASE 1 WHEN 1 THEN NULL; ELSE END CASE; END;", "END CASE; END;"},
		{"BEGIN WHILE 1 = 0 LOOP END LOOP; END;", "END LOOP; END;"},
		{"BEGIN NULL; EXCEPTION WHEN OTHERS THEN END;", "END;"},
		{"BEGIN NULL; EXCEPTION WHEN OTHERS NULL; END;", "NULL; END;"},
		{"DECLARE SUBTYPE s IS PLS_INTEGER RANGE 1 .. ; BEGIN NULL; END;", "; BEGIN NULL; END;"},
		{"DECLARE SUBTYPE s IS PLS_INTEGER RANGE .. 3; BEGIN NULL; END;", ".. 3; BEGIN NULL; END;"},
		{"BEGIN LOCK t IN EXCLUSIVE MODE; END;", "t IN EXCLUSIVE MODE; END;"},
		{"LOCK t IN EXCLUSIVE MODE", "t IN EXCLUSIVE MODE"},
	}
	for _, tc := range cases {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a ParseError", tc.sql, err)
			continue
		}
		if want := len(tc.sql) - len(tc.at); pe.Position != want {
			t.Errorf("Parse(%q) error at %d (%v), want %d, at %q", tc.sql, pe.Position, err, want, tc.at)
		}
	}
}

// TestPLSQLTransactionStatements checks the SQL statements a PL/SQL body runs
// as they are, which the parser used to skip as unknown words (ref_105 to
// ref_107).
func TestPLSQLTransactionStatements(t *testing.T) {
	sql := "BEGIN SAVEPOINT s; ROLLBACK TO SAVEPOINT s; ROLLBACK WORK; COMMIT WORK; SET TRANSACTION READ ONLY; " +
		"LOCK TABLE t IN EXCLUSIVE MODE NOWAIT; PRAGMA INLINE (pr, 'YES'); PRAGMA COVERAGE ('NOT_FEASIBLE'); NULL; END;"
	result := ParseAndCheck(t, sql)
	block := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock)
	var got []string
	for _, item := range block.Statements.Items {
		switch n := item.(type) {
		case *ast.SavepointStmt:
			got = append(got, "SAVEPOINT")
		case *ast.RollbackStmt:
			got = append(got, "ROLLBACK")
		case *ast.CommitStmt:
			got = append(got, "COMMIT")
		case *ast.SetTransactionStmt:
			got = append(got, "SET TRANSACTION")
		case *ast.LockTableStmt:
			got = append(got, "LOCK TABLE")
		case *ast.PLSQLPragma:
			got = append(got, "PRAGMA "+n.Name)
		case *ast.PLSQLNull:
			got = append(got, "NULL")
		default:
			got = append(got, "?")
		}
	}
	want := []string{"SAVEPOINT", "ROLLBACK", "ROLLBACK", "COMMIT", "SET TRANSACTION", "LOCK TABLE",
		"PRAGMA INLINE", "PRAGMA COVERAGE", "NULL"}
	if len(got) != len(want) {
		t.Fatalf("statements = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("statements = %v, want %v", got, want)
		}
	}
	if violations := CheckLocations(t, sql); len(violations) > 0 {
		t.Errorf("Loc violations: %v", violations)
	}
}

// TestPLSQLCaseAndIfStatementEnds checks that the arms of a CASE statement
// keep every statement, and that END CASE and END IF take a label, as on
// Oracle 23ai (ref_135).
func TestPLSQLCaseAndIfStatementEnds(t *testing.T) {
	sql := "DECLARE x NUMBER; BEGIN <<l>> CASE WHEN 1 = 1 THEN NULL; x := 1; ELSE x := 2; NULL; END CASE l; " +
		"<<m>> IF x = 1 THEN NULL; END IF m; IF x = 2 THEN NULL; END IF foo; END;"
	result := ParseAndCheck(t, sql)
	block := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock)
	labeled, ok := block.Statements.Items[0].(*ast.PLSQLLabeledStatement)
	if !ok {
		t.Fatalf("first statement = %T, want the labeled CASE", block.Statements.Items[0])
	}
	cs, ok := labeled.Statement.(*ast.PLSQLCase)
	if !ok {
		t.Fatalf("labeled statement = %T, want PLSQLCase", labeled.Statement)
	}
	if len(cs.Whens) != 1 || len(cs.Whens[0].Stmts) != 2 || len(cs.Else) != 2 {
		t.Errorf("CASE arms = %d WHEN with %v, ELSE %v; want one WHEN and ELSE with 2 statements each", len(cs.Whens), cs.Whens, cs.Else)
	}
	if got := sql[cs.Loc.Start:cs.Loc.End]; got[len(got)-len("END CASE l;"):] != "END CASE l;" {
		t.Errorf("CASE Loc covers %q, want it to end at END CASE l;", got)
	}
	if violations := CheckLocations(t, sql); len(violations) > 0 {
		t.Errorf("Loc violations: %v", violations)
	}
}

// TestPLSQLSubtypeDecl checks SUBTYPE declarations (ref_119 to ref_121).
func TestPLSQLSubtypeDecl(t *testing.T) {
	sql := "DECLARE SUBTYPE s IS PLS_INTEGER RANGE 0 .. 3 NOT NULL; SUBTYPE v IS VARCHAR2(10); SUBTYPE r IS t.c%TYPE; " +
		"SUBTYPE loop IS NUMBER; SUBTYPE \"SELECT\" IS NUMBER; BEGIN NULL; END;"
	result := ParseAndCheck(t, sql)
	decls := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock).Declarations
	if decls.Len() != 5 {
		t.Fatalf("declarations = %d, want 5", decls.Len())
	}
	s := decls.Items[0].(*ast.PLSQLSubtypeDecl)
	if s.Name != "S" || s.RangeLow == nil || s.RangeHigh == nil || !s.NotNull {
		t.Errorf("SUBTYPE s = %+v, want RANGE 0 .. 3 NOT NULL", s)
	}
	if got := sql[s.Loc.Start:s.Loc.End]; got != "SUBTYPE s IS PLS_INTEGER RANGE 0 .. 3 NOT NULL;" {
		t.Errorf("SUBTYPE s Loc covers %q", got)
	}
	v := decls.Items[1].(*ast.PLSQLSubtypeDecl)
	if v.RangeLow != nil || v.NotNull || v.BaseType.TypeMods.Len() != 1 {
		t.Errorf("SUBTYPE v = %+v, want VARCHAR2(10)", v)
	}
	if r := decls.Items[2].(*ast.PLSQLSubtypeDecl); !r.BaseType.IsPercType {
		t.Errorf("SUBTYPE r base type = %+v, want t.c%%TYPE", r.BaseType)
	}
	if violations := CheckLocations(t, sql); len(violations) > 0 {
		t.Errorf("Loc violations: %v", violations)
	}

	for _, bad := range []string{
		"DECLARE SUBTYPE s NUMBER; BEGIN NULL; END;",
		"DECLARE SUBTYPE s IS NUMBER 8, 2; BEGIN NULL; END;",
		"DECLARE SUBTYPE s IS PLS_INTEGER RANGE 1; BEGIN NULL; END;",
		"DECLARE SUBTYPE s IS; BEGIN NULL; END;",
		"DECLARE SUBTYPE IS IS NUMBER; BEGIN NULL; END;",
		"DECLARE SUBTYPE prior IS NUMBER; BEGIN NULL; END;",
	} {
		ParseShouldFail(t, bad)
	}
}

// TestPLSQLTypeLengthExpressions checks that a PL/SQL declaration takes an
// expression for a length, precision, or scale, while SQL keeps to integer
// literals (ref_122, ref_123).
func TestPLSQLTypeLengthExpressions(t *testing.T) {
	sql := "DECLARE v VARCHAR2(ORA_MAX_NAME_LEN + 2); w NUMBER(k.p, 2); c CHAR(k.n BYTE); z NUMBER(*, 2); BEGIN NULL; END;"
	result := ParseAndCheck(t, sql)
	decls := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock).Declarations
	mods := func(i int) []ast.Node {
		return decls.Items[i].(*ast.PLSQLVarDecl).TypeName.TypeMods.Items
	}
	if m := mods(0); len(m) != 1 || isIntegerMod(m[0]) {
		t.Errorf("VARCHAR2(ORA_MAX_NAME_LEN + 2) mods = %#v, want one expression", m)
	}
	if m := mods(1); len(m) != 2 || isIntegerMod(m[0]) || !isIntegerMod(m[1]) {
		t.Errorf("NUMBER(k.p, 2) mods = %#v, want an expression and an integer", m)
	}
	if m := mods(2); len(m) != 2 || isIntegerMod(m[0]) {
		t.Errorf("CHAR(k.n BYTE) mods = %#v, want an expression and BYTE", m)
	} else if s, ok := m[1].(*ast.String); !ok || s.Str != "BYTE" {
		t.Errorf("CHAR(k.n BYTE) length semantics = %#v, want BYTE", m[1])
	}
	if m := mods(3); len(m) != 2 || !isIntegerMod(m[1]) {
		t.Errorf("NUMBER(*, 2) mods = %#v, want * and an integer", m)
	}
	if violations := CheckLocations(t, sql); len(violations) > 0 {
		t.Errorf("Loc violations: %v", violations)
	}

	for _, bad := range []string{
		"DECLARE v VARCHAR2(10 x); BEGIN NULL; END;",
		"DECLARE v VARCHAR2(); BEGIN NULL; END;",
		"CREATE TABLE t (c VARCHAR2(n))",
		"CREATE TABLE t (c VARCHAR2(10 + 2))",
		// A datatype nested in a modifier follows the SQL rules (ref_138).
		"DECLARE v VARCHAR2(CAST(1 AS NUMBER(foo))); BEGIN NULL; END;",
	} {
		ParseShouldFail(t, bad)
	}
}

func isIntegerMod(n ast.Node) bool {
	_, ok := n.(*ast.Integer)
	return ok
}

// TestPLSQLCollectionElementNotNull checks NOT NULL on the element type of a
// nested table, associative array, and varray (ref_124).
func TestPLSQLCollectionElementNotNull(t *testing.T) {
	sql := "DECLARE TYPE a IS TABLE OF NUMBER NOT NULL INDEX BY PLS_INTEGER; TYPE b IS VARRAY(3) OF BLOB NOT NULL; " +
		"TYPE c IS TABLE OF VARCHAR2(ORA_MAX_NAME_LEN); BEGIN NULL; END;"
	result := ParseAndCheck(t, sql)
	decls := result.Items[0].(*ast.RawStmt).Stmt.(*ast.PLSQLBlock).Declarations
	for i, want := range []bool{true, true, false} {
		decl := decls.Items[i].(*ast.PLSQLTypeDecl)
		if decl.ElementNotNull != want {
			t.Errorf("TYPE %s ElementNotNull = %v, want %v", decl.Name, decl.ElementNotNull, want)
		}
	}
	if decls.Items[0].(*ast.PLSQLTypeDecl).IndexBy == nil {
		t.Error("TYPE a lost its INDEX BY after NOT NULL")
	}
	if violations := CheckLocations(t, sql); len(violations) > 0 {
		t.Errorf("Loc violations: %v", violations)
	}
}
