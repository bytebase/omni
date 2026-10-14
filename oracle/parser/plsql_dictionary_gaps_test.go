package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// The tests below cover constructs Oracle's own PL/SQL uses, checked against
// Oracle 23ai: every accepted form compiles there (or fails only on names it
// cannot resolve, or at run time), and every rejected one fails with the
// error noted beside it. Rows ref_209 to ref_215 and ref_219 to ref_222
// run them against the engine where a standalone statement can.

// TestPackageInitialization checks the initialization section of a package
// body.
func TestPackageInitialization(t *testing.T) {
	for _, sql := range []string{
		"CREATE PACKAGE BODY k AS PROCEDURE p IS BEGIN NULL; END; BEGIN NULL; END;",
		"CREATE PACKAGE BODY k AS PROCEDURE p IS BEGIN NULL; END; BEGIN NULL; EXCEPTION WHEN OTHERS THEN NULL; END k;",
		"CREATE PACKAGE BODY k AS c NUMBER; BEGIN c := 1; DECLARE d NUMBER; BEGIN d := c; END; END;",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		// The section holds a statement (PLS-00103), the package ends with
		// ';' (PLS-00103), and a specification has none (PLS-00103).
		"CREATE PACKAGE BODY k AS PROCEDURE p IS BEGIN NULL; END; BEGIN END;",
		"CREATE PACKAGE BODY k AS PROCEDURE p IS BEGIN NULL; END; BEGIN NULL; EXCEPTION END;",
		"CREATE PACKAGE BODY k AS PROCEDURE p IS BEGIN NULL; END; BEGIN NULL; END k",
		"CREATE PACKAGE k AS c NUMBER; BEGIN NULL; END;",
	} {
		ParseShouldFail(t, sql)
	}

	pkg := parseOne[*ast.CreatePackageStmt](t, "CREATE PACKAGE BODY k AS c NUMBER; BEGIN c := 1; c := 2; EXCEPTION WHEN OTHERS THEN NULL; END k;")
	if pkg.Body.Len() != 1 || pkg.InitBlock == nil || pkg.InitBlock.Statements.Len() != 2 || pkg.InitBlock.Exceptions.Len() != 1 {
		t.Fatalf("package body = %d items, initialization %+v", pkg.Body.Len(), pkg.InitBlock)
	}
}

// TestTriggerCallBody checks a trigger whose body is a CALL.
func TestTriggerCallBody(t *testing.T) {
	for _, sql := range []string{
		"CREATE TRIGGER t AFTER UPDATE OF a ON tab FOR EACH ROW CALL k.p(:old.a, :new.b)",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW WHEN (new.status = 'D') CALL k.p(:old.a, :new.status)",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(1, 'x')",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(x => 1, y => 'x')",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(:old.a + 1, UPPER(:new.b))",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(1, y => 'a')",
		// A packaged routine with its owner (Oracle 23ai creates both).
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL app.audit_pkg.record_change(1)",
		"CREATE TRIGGER t AFTER UPDATE ON tab CALL app.audit_pkg.record_change(1)",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		// An argument list holds an argument (ORA-00936), the call takes no
		// INTO (ORA-00911), the routine has a name (ORA-03050), and no
		// positional argument follows a named one (PLS-00312).
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p()",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL (1)",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(x => 1, 'a')",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(1,)",
		"CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(:old.a, :new.b) INTO :x",
		// Nor a database link (ORA-00911), nor a dangling component.
		"CREATE TRIGGER t AFTER UPDATE ON tab CALL app.audit_pkg.record_change@lnk(1)",
		"CREATE TRIGGER t AFTER UPDATE ON tab CALL k.p@lnk(1)",
		"CREATE TRIGGER t AFTER UPDATE ON tab CALL app.audit_pkg.(1)",
	} {
		ParseShouldFail(t, sql)
	}

	// The middle component of schema.package.routine is the package, in a
	// trigger and a CALL statement alike.
	full := parseOne[*ast.CreateTriggerStmt](t, "CREATE TRIGGER t AFTER UPDATE ON tab CALL app.audit_pkg.record_change(1)")
	if call, ok := full.Body.(*ast.CallStmt); !ok || call.Name.Schema != "APP" || call.Package != "AUDIT_PKG" || call.Name.Name != "RECORD_CHANGE" {
		t.Errorf("trigger body = %#v", full.Body)
	}
	if call := parseOne[*ast.CallStmt](t, "CALL app.audit_pkg.record_change(1)"); call.Name.Schema != "APP" || call.Package != "AUDIT_PKG" || call.Name.Name != "RECORD_CHANGE" {
		t.Errorf("CALL = %#v", call)
	}
}

// TestCallStmtArguments checks the argument list of a CALL statement
// against Oracle 23ai: it is in parentheses even when empty (ORA-06576 for
// CALL k.p), in positional, named, or mixed notation, with no positional
// argument after a named one (PLS-00312) and none left out (ORA-00936).
func TestCallStmtArguments(t *testing.T) {
	for _, sql := range []string{
		"CALL app.audit_pkg.record_change()",
		"CALL app.audit_pkg.record_change(1, 2)",
		"CALL app.audit_pkg.record_change(x => 1)",
		"CALL app.audit_pkg.record_change(1, y => 2)",
		"CALL k.p(y => 2)",
		"CALL f(1) INTO :x",
		// A database link's name may hold dots (ORA-04054 only for a
		// link that does not exist).
		"CALL app.pkg.proc@remote.example.com()",
		"CALL pkg.proc@remote.example.com()",
		"CALL app.pkg.f() INTO :v",
		`CALL app.pkg."SELECT"()`,
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"CALL app.audit_pkg.record_change",
		"CALL k.p",
		"CALL app.audit_pkg.record_change(1,)",
		"CALL app.audit_pkg.record_change(,1)",
		"CALL app.audit_pkg.record_change(1,,2)",
		"CALL app.audit_pkg.record_change(x => 1, 2)",
		"CALL k.p(1",
		// The INTO target is a host variable (ORA-00921, ORA-01745), no
		// component a reserved word unquoted (ORA-01747), and a link
		// follows every @ (ORA-01729).
		"CALL app.pkg.f() INTO",
		"CALL app.pkg.f() INTO 1",
		"CALL app.pkg.SELECT()",
		"CALL app.pkg@.proc()",
		"CREATE TRIGGER t AFTER UPDATE ON tab CALL app.pkg@.proc",
		"CREATE TRIGGER t AFTER UPDATE ON tab CALL app.pkg.SELECT(1)",
		// A name follows every dot (ORA-03050).
		"CALL app. .proc()",
		"CALL app.(1)",
		"CALL app.pkg.proc@remote.()",
		"CREATE TRIGGER t AFTER UPDATE ON tab CALL app. /*x*/ .proc",
	} {
		ParseShouldFail(t, sql)
	}
	call := parseOne[*ast.CallStmt](t, "CALL k.p(1, y => 2)")
	if named, ok := call.Args.Items[1].(*ast.NamedArgExpr); !ok || named.Name != "Y" {
		t.Errorf("second argument = %#v, want named argument Y", call.Args.Items[1])
	}

	trigger := parseOne[*ast.CreateTriggerStmt](t, "CREATE TRIGGER t AFTER UPDATE ON tab FOR EACH ROW CALL k.p(:old.a, y => 'x')")
	call, ok := trigger.Body.(*ast.CallStmt)
	if !ok || call.Name.Schema != "K" || call.Name.Name != "P" || call.Args.Len() != 2 {
		t.Fatalf("trigger body = %#v", trigger.Body)
	}
	if named, ok := call.Args.Items[1].(*ast.NamedArgExpr); !ok || named.Name != "Y" {
		t.Errorf("second argument = %#v, want named argument Y", call.Args.Items[1])
	}
}

// TestTypeBodyMethodModifiers checks inheritance clauses before the
// methods of a type body.
func TestTypeBodyMethodModifiers(t *testing.T) {
	for _, sql := range []string{
		"CREATE TYPE BODY t AS FINAL INSTANTIABLE CONSTRUCTOR FUNCTION t (x NUMBER) RETURN SELF AS RESULT IS BEGIN RETURN; END; NOT FINAL MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;",
		// OVERRIDING parses; it fails to match a spec without it (PLS-00538).
		"CREATE TYPE BODY t AS OVERRIDING MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		// A modifier repeated (PLS-00168), conflicting with the method kind
		// (PLS-00169), without a method, or NOT OVERRIDING (PLS-00103).
		"CREATE TYPE BODY t AS FINAL FINAL MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;",
		"CREATE TYPE BODY t AS NOT INSTANTIABLE CONSTRUCTOR FUNCTION t (x NUMBER) RETURN SELF AS RESULT IS BEGIN RETURN; END; END;",
		"CREATE TYPE BODY t AS FINAL END;",
		"CREATE TYPE BODY t AS FINAL MEMBER END;",
		"CREATE TYPE BODY t AS MEMBER END;",
		"CREATE TYPE BODY t AS OVERRIDING CONSTRUCTOR END;",
		"CREATE TYPE BODY t AS NOT OVERRIDING MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;",
	} {
		ParseShouldFail(t, sql)
	}
}

// TestBulkCursorAttributes checks the collection attributes of a cursor.
func TestBulkCursorAttributes(t *testing.T) {
	for _, sql := range []string{
		"DECLARE n NUMBER; BEGIN FOR i IN 1..SQL%BULK_EXCEPTIONS.COUNT LOOP n := SQL%BULK_EXCEPTIONS(i).ERROR_CODE; END LOOP; END;",
		"DECLARE n NUMBER; BEGIN n := SQL%BULK_EXCEPTIONS(1).ERROR_INDEX + SQL%BULK_ROWCOUNT(1); END;",
		// A named cursor's bulk attribute compiles; it fails at run time. A
		// bare one parses; assigning it to a number fails on its type
		// (PLS-00382).
		"DECLARE CURSOR c IS SELECT 1 FROM dual; n NUMBER; BEGIN n := c%BULK_ROWCOUNT(1); END;",
		"DECLARE n NUMBER; BEGIN n := SQL%BULK_ROWCOUNT; END;",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		// One index (PLS-00884 for none, PLS-00431 for two).
		"DECLARE n NUMBER; BEGIN n := SQL%BULK_ROWCOUNT(); END;",
		"DECLARE n NUMBER; BEGIN n := SQL%BULK_EXCEPTIONS(1, 2).ERROR_CODE; END;",
	} {
		ParseShouldFail(t, sql)
	}

	block := parseOne[*ast.PLSQLBlock](t, "DECLARE n NUMBER; BEGIN n := SQL%BULK_ROWCOUNT(3); END;")
	assign := block.Statements.Items[0].(*ast.PLSQLAssign)
	attr, ok := assign.Value.(*ast.CursorAttributeExpr)
	if !ok || attr.Attribute != "BULK_ROWCOUNT" || attr.Index == nil {
		t.Fatalf("assigned value = %#v", assign.Value)
	}
}

// TestCallArgumentsAndHandlers checks two rules of every call and block:
// no positional argument follows a named one (PLS-00312, in SQL and PL/SQL
// calls alike), and EXCEPTION takes at least one handler (PLS-00103).
func TestCallArgumentsAndHandlers(t *testing.T) {
	for _, sql := range []string{
		"BEGIN p(1, y => 'a'); END;",
		"SELECT f(1, y => 2) FROM dual",
		"BEGIN NULL; EXCEPTION WHEN OTHERS THEN NULL; END;",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"BEGIN p(x => 1, 'a'); END;",
		"SELECT f(x => 1, 2) FROM dual",
		"DECLARE n NUMBER; BEGIN n := f(x => 1, 2); END;",
		"BEGIN NULL; EXCEPTION END;",
	} {
		ParseShouldFail(t, sql)
	}
}
