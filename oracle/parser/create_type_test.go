package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateTypeObject(t *testing.T) {
	p := newTestParser("TYPE address_t AS OBJECT (street VARCHAR2(100), city VARCHAR2(50), zip NUMBER(5))")
	stmt, parseErr1 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if stmt == nil {
		t.Fatal("expected CreateTypeStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "ADDRESS_T" {
		t.Errorf("expected type name ADDRESS_T, got %v", stmt.Name)
	}
	if stmt.Attributes == nil || stmt.Attributes.Len() != 3 {
		t.Fatalf("expected 3 attributes, got %d", stmt.Attributes.Len())
	}
	attr0 := stmt.Attributes.Items[0].(*ast.ColumnDef)
	if attr0.Name != "STREET" {
		t.Errorf("expected attribute name STREET, got %q", attr0.Name)
	}
	if attr0.TypeName == nil || attr0.TypeName.Names.Len() == 0 {
		t.Error("expected non-nil TypeName for first attribute")
	}
}

func TestParseCreateTypeOrReplace(t *testing.T) {
	p := newTestParser("TYPE my_type AS OBJECT (id NUMBER)")
	stmt, parseErr2 := p.parseCreateTypeStmt(0, true, false, false, false)
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
}

func TestParseCreateTypeTableOf(t *testing.T) {
	p := newTestParser("TYPE num_table AS TABLE OF NUMBER")
	stmt, parseErr3 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if stmt.AsTable == nil {
		t.Fatal("expected non-nil AsTable")
	}
	if stmt.AsTable.Names == nil || stmt.AsTable.Names.Len() == 0 {
		t.Error("expected AsTable type names")
	}
	nameStr := stmt.AsTable.Names.Items[0].(*ast.String)
	if nameStr.Str != "NUMBER" {
		t.Errorf("expected NUMBER, got %q", nameStr.Str)
	}
}

func TestParseCreateTypeVarray(t *testing.T) {
	p := newTestParser("TYPE phone_list AS VARRAY (10) OF VARCHAR2(20)")
	stmt, parseErr4 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if stmt.AsVarray == nil {
		t.Fatal("expected non-nil AsVarray")
	}
	if stmt.VarraySize == nil {
		t.Fatal("expected non-nil VarraySize")
	}
	sizeNum := stmt.VarraySize.(*ast.NumberLiteral)
	if sizeNum.Ival != 10 {
		t.Errorf("expected varray size 10, got %d", sizeNum.Ival)
	}
}

func TestParseCreateTypeBody(t *testing.T) {
	// A type body holds subprogram definitions and ends with END; Oracle
	// rejects a body without it (PLS-00103), and so does omni.
	p := newTestParser("TYPE BODY my_type IS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END my_type;")
	stmt, parseErr5 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if !stmt.IsBody {
		t.Error("expected IsBody to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "MY_TYPE" {
		t.Errorf("expected type name MY_TYPE, got %v", stmt.Name)
	}
}

func TestParseCreateTypeSchemaQualified(t *testing.T) {
	p := newTestParser("TYPE hr.address_t AS OBJECT (id NUMBER)")
	stmt, parseErr6 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if stmt.Name == nil || stmt.Name.Schema != "HR" || stmt.Name.Name != "ADDRESS_T" {
		t.Errorf("expected HR.ADDRESS_T, got %v", stmt.Name)
	}
}

func TestParseCreateTypeLoc(t *testing.T) {
	p := newTestParser("TYPE t AS OBJECT (id NUMBER)")
	stmt, parseErr7 := p.parseCreateTypeStmt(0, false, false, false, false)
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}

func TestP2CreateTypeBodyMemberStructure(t *testing.T) {
	result := ParseAndCheck(t, `CREATE TYPE BODY person_type AS
  MEMBER FUNCTION get_name RETURN VARCHAR2 IS
  BEGIN
    RETURN first_name;
  END get_name;
  STATIC PROCEDURE set_count(p_count NUMBER) IS
  BEGIN
    NULL;
  END set_count;
END;`)
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateTypeStmt)
	if !ok {
		t.Fatalf("expected CreateTypeStmt, got %T", raw.Stmt)
	}
	if !stmt.IsBody {
		t.Fatal("expected type body")
	}
	if stmt.Body == nil || stmt.Body.Len() != 2 {
		t.Fatalf("expected 2 body members, got %#v", stmt.Body)
	}

	memberFn := stmt.Body.Items[0].(*ast.TypeBodyMember)
	if memberFn.Kind != ast.TYPE_BODY_MEMBER {
		t.Fatalf("member 1 kind=%d, want MEMBER", memberFn.Kind)
	}
	fn, ok := memberFn.Subprog.(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("member 1 subprogram=%T, want function", memberFn.Subprog)
	}
	fnBlock, ok := fn.Body.(*ast.PLSQLBlock)
	if !ok {
		t.Fatalf("function body=%T, want PLSQLBlock", fn.Body)
	}
	if fn.Name == nil || fn.Name.Name != "GET_NAME" || fnBlock.Statements == nil || fnBlock.Statements.Len() != 1 {
		t.Fatalf("function body was not structured: %#v", fn)
	}
	if _, ok := fnBlock.Statements.Items[0].(*ast.PLSQLReturn); !ok {
		t.Fatalf("function statement=%T, want PLSQLReturn", fnBlock.Statements.Items[0])
	}

	memberProc := stmt.Body.Items[1].(*ast.TypeBodyMember)
	if memberProc.Kind != ast.TYPE_BODY_STATIC {
		t.Fatalf("member 2 kind=%d, want STATIC", memberProc.Kind)
	}
	proc, ok := memberProc.Subprog.(*ast.CreateProcedureStmt)
	if !ok {
		t.Fatalf("member 2 subprogram=%T, want procedure", memberProc.Subprog)
	}
	procBlock, ok := proc.Body.(*ast.PLSQLBlock)
	if !ok {
		t.Fatalf("procedure body=%T, want PLSQLBlock", proc.Body)
	}
	if proc.Name == nil || proc.Name.Name != "SET_COUNT" || procBlock.Statements == nil || procBlock.Statements.Len() != 1 {
		t.Fatalf("procedure body was not structured: %#v", proc)
	}
	if _, ok := procBlock.Statements.Items[0].(*ast.PLSQLNull); !ok {
		t.Fatalf("procedure statement=%T, want PLSQLNull", procBlock.Statements.Items[0])
	}
}

func TestP2CreateTypeBodyRequiresPLSQLBlock(t *testing.T) {
	ParseShouldFail(t, "CREATE TYPE BODY person_type AS MEMBER PROCEDURE p IS END p; END;")
}

func TestP2CreateTypeBodyLoc(t *testing.T) {
	stmt := ParseAndCheck(t, `CREATE TYPE BODY person_type AS
  MEMBER PROCEDURE p IS
  BEGIN
    NULL;
  END p;
END;`).Items[0].(*ast.RawStmt).Stmt.(*ast.CreateTypeStmt)
	if stmt.Loc.Start != 0 || stmt.Loc.End <= stmt.Loc.Start {
		t.Fatalf("invalid type body Loc=%+v", stmt.Loc)
	}
	member := stmt.Body.Items[0].(*ast.TypeBodyMember)
	if member.Loc.Start <= stmt.Loc.Start || member.Loc.End > stmt.Loc.End {
		t.Fatalf("member Loc=%+v not inside stmt Loc=%+v", member.Loc, stmt.Loc)
	}
}

// TestParseObjectTypeSpecMethods covers method specifications in an object
// type specification, with and without call specs. Oracle 23ai compiles each
// spec VALID.
func TestParseObjectTypeSpecMethods(t *testing.T) {
	sql := "CREATE OR REPLACE TYPE emp_t AS OBJECT (" +
		"id NUMBER, name VARCHAR2(30), " +
		"MEMBER FUNCTION wages RETURN NUMBER AS LANGUAGE JAVA NAME 'Paymaster.wages() return java.math.BigDecimal', " +
		"STATIC PROCEDURE log_it(msg VARCHAR2) AS LANGUAGE C LIBRARY lib WITH CONTEXT PARAMETERS (CONTEXT, msg STRING), " +
		"NOT FINAL MEMBER PROCEDURE raise_pay(pct NUMBER), " +
		"FINAL INSTANTIABLE CONSTRUCTOR FUNCTION emp_t(SELF IN OUT NOCOPY emp_t, id NUMBER) RETURN SELF AS RESULT, " +
		"MAP MEMBER FUNCTION sort_key RETURN NUMBER DETERMINISTIC, " +
		"PRAGMA RESTRICT_REFERENCES (wages, WNDS, RNPS)" +
		") NOT FINAL NOT INSTANTIABLE"
	result := ParseAndCheck(t, sql)
	stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateTypeStmt)
	if stmt.Attributes.Len() != 2 {
		t.Fatalf("Attributes = %d, want 2", stmt.Attributes.Len())
	}
	if stmt.Methods == nil || stmt.Methods.Len() != 6 {
		t.Fatalf("Methods = %v, want 6 items", stmt.Methods)
	}
	wages := stmt.Methods.Items[0].(*ast.TypeBodyMember).Subprog.(*ast.CreateFunctionStmt)
	if wages.CallSpec == nil || wages.CallSpec.Language != "JAVA" || wages.Body != nil {
		t.Fatalf("wages = %s, want a Java call spec", ast.NodeToString(wages))
	}
	logIt := stmt.Methods.Items[1].(*ast.TypeBodyMember)
	if logIt.Kind != ast.TYPE_BODY_STATIC || logIt.Subprog.(*ast.CreateProcedureStmt).CallSpec == nil {
		t.Fatalf("log_it = %s, want a STATIC C call spec", ast.NodeToString(logIt))
	}
	raise := stmt.Methods.Items[2].(*ast.TypeBodyMember)
	if len(raise.Modifiers) != 1 || raise.Modifiers[0] != "NOT FINAL" || raise.Subprog.(*ast.CreateProcedureStmt).CallSpec != nil {
		t.Fatalf("raise_pay = %s", ast.NodeToString(raise))
	}
	ctor := stmt.Methods.Items[3].(*ast.TypeBodyMember)
	if ctor.Kind != ast.TYPE_BODY_CONSTRUCTOR || len(ctor.Modifiers) != 2 {
		t.Fatalf("constructor = %s", ast.NodeToString(ctor))
	}
	if pragma, ok := stmt.Methods.Items[5].(*ast.PLSQLPragma); !ok || pragma.Name != "RESTRICT_REFERENCES" || pragma.Args.Len() != 3 {
		t.Fatalf("pragma = %s", ast.NodeToString(stmt.Methods.Items[5]))
	}
	if len(stmt.Modifiers) != 2 || stmt.Modifiers[0] != "NOT FINAL" || stmt.Modifiers[1] != "NOT INSTANTIABLE" {
		t.Fatalf("Modifiers = %q", stmt.Modifiers)
	}
	if violations := CheckLocations(t, sql); len(violations) > 0 {
		t.Fatalf("Loc violations: %v", violations)
	}
}

// TestParseObjectTypeSpecMLECallSpecs checks the MLE call specs Oracle 23ai
// accepts on a method spec: the MODULE form reaches name resolution
// (PLS-00201 for an absent module), the inline form the MLE resolver.
func TestParseObjectTypeSpecMLECallSpecs(t *testing.T) {
	for _, sql := range []string{
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS MLE MODULE m SIGNATURE 'f()')",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; }}, MEMBER PROCEDURE p)",
	} {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateTypeStmt)
			f := stmt.Methods.Items[0].(*ast.TypeBodyMember).Subprog.(*ast.CreateFunctionStmt)
			if f.CallSpec == nil || !f.CallSpec.MLE {
				t.Fatalf("f = %s, want an MLE call spec", ast.NodeToString(f))
			}
			if violations := CheckLocations(t, sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
}

// TestParseObjectTypeSpecMethodRules pins the method spec rules Oracle 23ai
// checks from the statement text, each case beside one it accepts.
func TestParseObjectTypeSpecMethodRules(t *testing.T) {
	accept := []string{
		// A constructor bears its type's name, quoted or not.
		`CREATE TYPE t AS OBJECT (a NUMBER, CONSTRUCTOR FUNCTION "T" RETURN SELF AS RESULT, CONSTRUCTOR FUNCTION t(x NUMBER) RETURN SELF AS RESULT)`,
		// MAP takes no parameter but SELF, ORDER one besides SELF. Their
		// return and parameter types are not checked: that needs the types
		// resolved (PLS-00522 through PLS-00524).
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER FUNCTION m RETURN VARCHAR2)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER FUNCTION m(SELF t) RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, ORDER MEMBER FUNCTION o(x t) RETURN INTEGER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, ORDER MEMBER FUNCTION o(SELF t, x t) RETURN NUMBER)",
		// [ OVERRIDING ] { [ NOT ] FINAL | [ NOT ] INSTANTIABLE }...
		"CREATE TYPE t AS OBJECT (a NUMBER, OVERRIDING MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, OVERRIDING NOT FINAL MEMBER FUNCTION f RETURN NUMBER) NOT FINAL",
		"CREATE TYPE t AS OBJECT (a NUMBER, INSTANTIABLE FINAL MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT INSTANTIABLE NOT FINAL MEMBER FUNCTION f RETURN NUMBER) NOT FINAL NOT INSTANTIABLE",
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT INSTANTIABLE MAP MEMBER FUNCTION m RETURN NUMBER) NOT FINAL NOT INSTANTIABLE",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
	reject := []string{
		// PLS-00658: a constructor named otherwise.
		"CREATE TYPE t AS OBJECT (a NUMBER, CONSTRUCTOR FUNCTION wrong RETURN SELF AS RESULT)",
		// PLS-00520 and PLS-00521: MAP and ORDER parameter counts.
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER FUNCTION m(x NUMBER) RETURN VARCHAR2)",
		"CREATE TYPE t AS OBJECT (a NUMBER, ORDER MEMBER FUNCTION o RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, ORDER MEMBER FUNCTION o(x t, y t) RETURN NUMBER)",
		// PLS-00154: one MAP or one ORDER method per type.
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER FUNCTION m RETURN NUMBER, MAP MEMBER FUNCTION n RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER FUNCTION m RETURN NUMBER, ORDER MEMBER FUNCTION o(x t) RETURN NUMBER)",
		// PLS-00169: conflicting clauses.
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT INSTANTIABLE FINAL MEMBER FUNCTION f RETURN NUMBER) NOT FINAL NOT INSTANTIABLE",
		"CREATE TYPE t AS OBJECT (a NUMBER, FINAL NOT INSTANTIABLE MEMBER FUNCTION f RETURN NUMBER) NOT FINAL NOT INSTANTIABLE",
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT INSTANTIABLE STATIC FUNCTION f RETURN NUMBER) NOT FINAL NOT INSTANTIABLE",
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT INSTANTIABLE CONSTRUCTOR FUNCTION t RETURN SELF AS RESULT) NOT FINAL NOT INSTANTIABLE",
		"CREATE TYPE t AS OBJECT (a NUMBER, OVERRIDING STATIC FUNCTION g RETURN NUMBER)",
		// PLS-00103: OVERRIDING takes no NOT and comes first, once.
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT OVERRIDING MEMBER FUNCTION f RETURN NUMBER) NOT FINAL",
		"CREATE TYPE t AS OBJECT (a NUMBER, FINAL OVERRIDING MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, OVERRIDING OVERRIDING MEMBER FUNCTION f RETURN NUMBER)",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}

// TestParseObjectTypeSpecRejects lists object type specifications Oracle 23ai
// rejects: a call spec ended by ';', the EXTERNAL call spec, or a PL/SQL body
// in the spec (PLS-00103), a method without a name, a function without
// RETURN, or a constructor without all of RETURN SELF AS RESULT (PLS-00103,
// PLS-00659 for another return type), an attribute named like a method
// keyword unless quoted (PLS-00103), an attribute after a method or no
// attribute at all (PLS-00589), a repeated modifier (PLS-00168), PERSISTABLE
// before a method (PLS-00771), a MAP or ORDER procedure (PLS-00155), and a
// pragma other than RESTRICT_REFERENCES (PLS-00127).
func TestParseObjectTypeSpecRejects(t *testing.T) {
	for _, sql := range []string{
		"CREATE TYPE t AS OBJECT (a NUMBER, CONSTRUCTOR FUNCTION t RETURN SELF)",
		"CREATE TYPE t AS OBJECT (a NUMBER, CONSTRUCTOR FUNCTION t RETURN SELF AS)",
		"CREATE TYPE t AS OBJECT (a NUMBER, CONSTRUCTOR FUNCTION t RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, PERSISTABLE MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT PERSISTABLE MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, FINAL PERSISTABLE MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER PROCEDURE p)",
		"CREATE TYPE t AS OBJECT (a NUMBER, ORDER MEMBER PROCEDURE p)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS MLE MODULE m SIGNATURE 'f()';)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS EXTERNAL LIBRARY lib)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f, MEMBER PROCEDURE p)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER PROCEDURE)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END)",
		"CREATE TYPE t AS OBJECT (map NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, member NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, final NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (MEMBER FUNCTION f RETURN NUMBER, a NUMBER)",
		"CREATE TYPE t AS OBJECT (MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT ()",
		"CREATE TYPE t AS OBJECT (a NUMBER) FINAL FINAL",
		"CREATE TYPE t AS OBJECT (a NUMBER) NOT FINAL FINAL",
		"CREATE TYPE t AS OBJECT (a NUMBER, FINAL FINAL MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER, PRAGMA INLINE (f, 'YES'))",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER, PRAGMA RESTRICT_REFERENCES (f))",
		"CREATE TYPE t AS OBJECT (a NUMBER, CONSTRUCTOR PROCEDURE p)",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
	// Quoted, the method keywords name ordinary attributes.
	ParseAndCheck(t, `CREATE TYPE t AS OBJECT ("MAP" NUMBER, "ORDER" NUMBER, a NUMBER)`)
}

// TestParseTypeBodyMethodRejects lists type body methods Oracle 23ai rejects
// with PLS-00103: a function without RETURN, a method without a name, a
// PL/SQL block or call spec without IS|AS before it, and a constructor
// without all of RETURN SELF AS RESULT.
func TestParseTypeBodyMethodRejects(t *testing.T) {
	for _, sql := range []string{
		"CREATE TYPE BODY t AS MEMBER FUNCTION f IS BEGIN RETURN 1; END; END;",
		"CREATE TYPE BODY t AS MEMBER PROCEDURE IS BEGIN NULL; END; END;",
		"CREATE TYPE BODY t AS MEMBER FUNCTION f RETURN NUMBER LANGUAGE JAVA NAME 'X.f() return int'; END;",
		"CREATE TYPE BODY t AS MEMBER FUNCTION f RETURN NUMBER BEGIN RETURN 1; END; END;",
		"CREATE TYPE BODY t AS MEMBER PROCEDURE p BEGIN NULL; END; END;",
		"CREATE TYPE BODY t AS CONSTRUCTOR FUNCTION t(x NUMBER) RETURN SELF IS BEGIN RETURN; END; END;",
		"CREATE TYPE BODY t AS CONSTRUCTOR FUNCTION t(x NUMBER) RETURN SELF AS IS BEGIN RETURN; END; END;",
	} {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
	// A body may define a MAP procedure: Oracle 23ai rejects it only for not
	// matching a spec, which cannot declare one (PLS-00538).
	ParseAndCheck(t, "CREATE TYPE BODY t AS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; "+
		"MAP MEMBER PROCEDURE p IS BEGIN NULL; END; "+
		"CONSTRUCTOR FUNCTION t(x NUMBER) RETURN SELF AS RESULT IS BEGIN RETURN; END; END;")
}
