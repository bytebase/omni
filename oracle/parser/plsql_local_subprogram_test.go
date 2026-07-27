package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// TestPLSQLLocalSubprograms pins local procedure/function declarations in
// PL/SQL declaration sections. All accepted shapes are engine-verified on
// Oracle 23ai. Assertions check AST shape, not just parse success: before
// the fix, a parameterless local procedure was wrongly accepted as a run of
// PLSQLVarDecl fragments (fail-open), so err == nil alone proves nothing.
func TestPLSQLLocalSubprograms(t *testing.T) {
	declsOf := func(t *testing.T, sql string) []nodes.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		block := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.PLSQLBlock)
		if block.Declarations == nil {
			t.Fatal("no declarations")
		}
		return block.Declarations.Items
	}

	t.Run("local procedure with parameters", func(t *testing.T) {
		decls := declsOf(t, "DECLARE PROCEDURE p1(x IN VARCHAR2) IS BEGIN NULL; END; BEGIN p1('a'); END;")
		if len(decls) != 1 {
			t.Fatalf("expected 1 declaration, got %d", len(decls))
		}
		ps, ok := decls[0].(*nodes.CreateProcedureStmt)
		if !ok {
			t.Fatalf("expected CreateProcedureStmt, got %T", decls[0])
		}
		if ps.Body == nil || ps.Parameters == nil || ps.Parameters.Len() != 1 {
			t.Fatalf("procedure shape wrong: body=%v params=%v", ps.Body, ps.Parameters)
		}
	})

	t.Run("local function with parameters", func(t *testing.T) {
		decls := declsOf(t, "DECLARE FUNCTION g1(x NUMBER) RETURN NUMBER IS BEGIN RETURN x + 1; END; BEGIN NULL; END;")
		fs, ok := decls[0].(*nodes.CreateFunctionStmt)
		if !ok {
			t.Fatalf("expected CreateFunctionStmt, got %T", decls[0])
		}
		if fs.Body == nil || fs.ReturnType == nil {
			t.Fatal("function missing body or return type")
		}
	})

	t.Run("parameterless forms are single correct nodes", func(t *testing.T) {
		// The fail-open regression guard: these must NOT parse into
		// PLSQLVarDecl fragments.
		decls := declsOf(t, "DECLARE v NUMBER := 0; PROCEDURE p2 IS BEGIN NULL; END; FUNCTION g2 RETURN NUMBER IS BEGIN RETURN 2; END; BEGIN p2; END;")
		if len(decls) != 3 {
			t.Fatalf("expected 3 declarations, got %d: %#v", len(decls), decls)
		}
		if _, ok := decls[0].(*nodes.PLSQLVarDecl); !ok {
			t.Fatalf("decl[0] = %T, want *PLSQLVarDecl", decls[0])
		}
		if _, ok := decls[1].(*nodes.CreateProcedureStmt); !ok {
			t.Fatalf("decl[1] = %T, want *CreateProcedureStmt", decls[1])
		}
		if _, ok := decls[2].(*nodes.CreateFunctionStmt); !ok {
			t.Fatalf("decl[2] = %T, want *CreateFunctionStmt", decls[2])
		}
	})

	t.Run("nested local subprogram two levels", func(t *testing.T) {
		decls := declsOf(t, "DECLARE FUNCTION outer_f(p_h NUMBER) RETURN VARCHAR2 IS FUNCTION hx(v NUMBER) RETURN VARCHAR2 IS BEGIN RETURN TO_CHAR(v); END; BEGIN RETURN hx(p_h); END; BEGIN NULL; END;")
		outer, ok := decls[0].(*nodes.CreateFunctionStmt)
		if !ok {
			t.Fatalf("expected CreateFunctionStmt, got %T", decls[0])
		}
		body, ok := outer.Body.(*nodes.PLSQLBlock)
		if !ok || body.Declarations == nil || body.Declarations.Len() != 1 {
			t.Fatalf("outer body missing nested declaration: %T", outer.Body)
		}
		if _, ok := body.Declarations.Items[0].(*nodes.CreateFunctionStmt); !ok {
			t.Fatalf("nested decl = %T, want *CreateFunctionStmt", body.Declarations.Items[0])
		}
	})

	t.Run("forward declaration then definition", func(t *testing.T) {
		decls := declsOf(t, "DECLARE PROCEDURE fwd; PROCEDURE fwd IS BEGIN NULL; END; BEGIN fwd; END;")
		if len(decls) != 2 {
			t.Fatalf("expected 2 declarations, got %d", len(decls))
		}
		spec := decls[0].(*nodes.CreateProcedureStmt)
		def := decls[1].(*nodes.CreateProcedureStmt)
		if spec.Body != nil {
			t.Fatal("forward declaration must have no body")
		}
		if def.Body == nil {
			t.Fatal("definition must have a body")
		}
	})

	t.Run("procedure declare section supports local subprograms", func(t *testing.T) {
		// The customer shape: local subprogram inside a standalone
		// function/procedure body's declaration section.
		list, err := Parse("CREATE OR REPLACE FUNCTION f_outer(p NUMBER) RETURN VARCHAR2 IS FUNCTION hsl(x NUMBER) RETURN VARCHAR2 IS BEGIN RETURN TO_CHAR(x); END; BEGIN RETURN hsl(p); END;")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		fn := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateFunctionStmt)
		body, ok := fn.Body.(*nodes.PLSQLBlock)
		if !ok || body.Declarations == nil || body.Declarations.Len() != 1 {
			t.Fatalf("outer function body missing local subprogram declaration: %T", fn.Body)
		}
	})
}
