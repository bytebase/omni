package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestConditionalCompilation checks PL/SQL conditional compilation against
// Oracle 23ai: every accepted form compiles there, and every rejected one
// fails with the error noted beside it. Rows ref_264 to ref_271 run them
// against the engine.
func TestConditionalCompilation(t *testing.T) {
	for _, sql := range []string{
		// Selection directives around declarations, statements, and
		// expressions, and in a unit's heading.
		"DECLARE $IF DBMS_DB_VERSION.VERSION >= 12 $THEN a NUMBER; $ELSE a VARCHAR2(10); $END BEGIN NULL; END;",
		"DECLARE x NUMBER; BEGIN x := $IF TRUE $THEN 1 $ELSE 2 $END; END;",
		"DECLARE x NUMBER; BEGIN $IF $$PLSQL_DEBUG $THEN x := 1; $ELSIF FALSE $THEN x := 2; $ELSE x := 3; $END END;",
		"BEGIN $if true $then null; $end end;",
		"BEGIN $IF TRUE $THEN $END NULL; END;",
		"BEGIN $IF (TRUE AND NOT FALSE) OR $$x = 3 $THEN NULL; $END END;",
		"BEGIN $IF $$no_such_flag IS NULL $THEN NULL; $END END;",
		"BEGIN $IF $$a $THEN $IF $$b $THEN NULL; $ELSE NULL; $END $ELSE NULL; $END END;",
		"CREATE OR REPLACE PROCEDURE p IS $IF TRUE $THEN a NUMBER; $END BEGIN NULL; END;",
		"CREATE OR REPLACE PACKAGE k AUTHID $IF TRUE $THEN DEFINER $ELSE CURRENT_USER $END AS c NUMBER; END;",
		"CREATE OR REPLACE PACKAGE BODY k AS $IF $$x $THEN PROCEDURE p IS BEGIN NULL; END; $END END;",
		"CREATE OR REPLACE TRIGGER trg BEFORE INSERT ON t BEGIN $IF $$x $THEN NULL; $END NULL; END;",
		"CREATE OR REPLACE TYPE BODY ty AS MEMBER PROCEDURE p IS BEGIN $IF $$x $THEN NULL; $END NULL; END; END;",
		// Inquiry directives are identifiers, in SQL inside PL/SQL too.
		"BEGIN DBMS_OUTPUT.PUT_LINE($$PLSQL_LINE || $$PLSQL_UNIT); END;",
		"DECLARE x NUMBER := $$PLSQL_LINE; BEGIN NULL; END;",
		"DECLARE x VARCHAR2(100); BEGIN SELECT $$PLSQL_UNIT INTO x FROM dual; END;",
		// Oracle 23ai does not parse the code of a branch it does not
		// select, nor the conditions it does not evaluate: a FALSE, NULL,
		// or TRUE condition settles the selection.
		"BEGIN $IF FALSE $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF NULL $THEN garbage; $ELSIF TRUE $THEN NULL; $END END;",
		"BEGIN $IF TRUE $THEN $IF FALSE $THEN garbage garbage; $END NULL; $END END;",
		"BEGIN $IF TRUE $THEN NULL; $ELSIF garbage garbage $THEN NULL; $END END;",
		"BEGIN $IF TRUE $THEN NULL; $ELSE $IF garbage $THEN x $ELSE y $END $END END;",
		"BEGIN $IF TRUE $THEN NULL; $ELSE NULL $END END;",
		// $ERROR in a branch that may not be selected.
		"BEGIN $IF FALSE $THEN $ERROR 'unsupported' $END $END NULL; END;",
		"BEGIN $IF $$x $THEN $ERROR 'unsupported ' || $$PLSQL_UNIT $END $END NULL; END;",
		// CASE in a selected branch, decided by a scan that starts there.
		"DECLARE x NUMBER; BEGIN $IF $$x $THEN x := CASE WHEN 1 = 1 THEN 1 END; $ELSE x := CASE WHEN 1 = 2 THEN 2 END; $END x := CASE WHEN x > 0 THEN 0 END; END;",
		// A unit's text runs on after its END.
		"BEGIN NULL; END; $IF TRUE $THEN $END",
	} {
		ParseAndCheck(t, sql)
	}

	for _, sql := range []string{
		// Directives out of place (PLS-00176) or without their end
		// (PLS-00177).
		"BEGIN NULL; $END END;",
		"BEGIN $ELSE NULL; END;",
		"BEGIN NULL; $THEN END;",
		"BEGIN $IF TRUE NULL; $END END;",
		"BEGIN $IF TRUE $THEN $THEN NULL; $END END;",
		"BEGIN $IF TRUE $THEN NULL; END;",
		"BEGIN $IF $$x $THEN NULL; $ELSE NULL; $ELSE NULL; $END END;",
		"BEGIN $IF $$x $THEN NULL; $ELSE NULL; $ELSIF TRUE $THEN NULL; $END END;",
		"BEGIN $IF TRUE $THEN NULL; $ELSE $THEN $END END;",
		"BEGIN $IF FALSE $THEN NULL; $ELSE NULL; $ELSE NULL; $END END;",
		"BEGIN $ERROR 'x' END;",
		// A condition is a static expression (PLS-00174, PLS-00103).
		"BEGIN $IF $THEN NULL; $END END;",
		"BEGIN $IF garbage garbage $THEN NULL; $END END;",
		"BEGIN $IF FALSE $THEN NULL; $ELSIF garbage garbage $THEN NULL; $END END;",
		// An $ERROR certainly reached (PLS-00179).
		"BEGIN $ERROR 'x' $END NULL; END;",
		"BEGIN $IF TRUE $THEN $ERROR 'x' $END $END NULL; END;",
		"BEGIN $IF FALSE $THEN NULL; $ELSE $ERROR 'x' || 'y' $END $END NULL; END;",
		// The selected code is parsed: an empty statement after $END
		// (PLS-00103), and an unknown directive (PLS-00181).
		"BEGIN $IF TRUE $THEN NULL; $END; END;",
		"BEGIN DBMS_OUTPUT.PUT_LINE($foo); END;",
		// Only a PL/SQL unit is preprocessed (ORA-00911 elsewhere, PLS-00180
		// in an object type specification).
		"SELECT $IF 1=1 $THEN 1 $ELSE 2 $END FROM dual",
		"SELECT 1 FROM dual WHERE $$x IS NULL",
		"CREATE OR REPLACE VIEW v AS SELECT $IF TRUE $THEN 1 $END a FROM dual",
		"$IF TRUE $THEN BEGIN NULL; END; $END",
		"CREATE OR REPLACE TYPE ty AS OBJECT ($IF TRUE $THEN a NUMBER $ELSE a DATE $END)",
		"SELECT 1 FROM dual; $IF TRUE $THEN $END",
		"SELECT 1 FROM dual; $IF TRUE $THEN $END BEGIN NULL; END;",
	} {
		ParseShouldFail(t, sql)
	}
}

// TestConditionalCompilationTree checks what the parser builds from the
// selected branch, and where it reports a misplaced directive.
func TestConditionalCompilationTree(t *testing.T) {
	block := parseOne[*ast.PLSQLBlock](t, "DECLARE $IF FALSE $THEN a DATE; $ELSE a NUMBER; b NUMBER; $END BEGIN NULL; END;")
	if block.Declarations.Len() != 2 {
		t.Fatalf("declarations = %d, want the 2 of the $ELSE branch", block.Declarations.Len())
	}

	sql := "BEGIN DBMS_OUTPUT.PUT_LINE($$PLSQL_LINE); END;"
	block = parseOne[*ast.PLSQLBlock](t, sql)
	var inquiry *ast.ColumnRef
	ast.Inspect(block, func(n ast.Node) bool {
		if ref, ok := n.(*ast.ColumnRef); ok && ref.Column == "$$PLSQL_LINE" {
			inquiry = ref
		}
		return true
	})
	if inquiry == nil {
		t.Fatalf("no $$PLSQL_LINE reference in %q", sql)
	}
	if got := sql[inquiry.Loc.Start:inquiry.Loc.End]; got != "$$PLSQL_LINE" {
		t.Errorf("inquiry spans %q", got)
	}

	for _, tc := range []struct {
		sql  string
		want string
	}{
		{"BEGIN NULL; $END END;", "$END"},
		{"BEGIN $IF $$x $THEN NULL; $ELSE NULL; $ELSE NULL; $END END;", "$ELSE NULL; $END"},
		{"SELECT 1 FROM dual WHERE $$x IS NULL", "$$x"},
	} {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Fatalf("Parse(%q) = %v, want a ParseError", tc.sql, err)
		}
		if !strings.HasPrefix(tc.sql[pe.Position:], tc.want) {
			t.Errorf("Parse(%q) fails at %d (%q), want at %q", tc.sql, pe.Position, tc.sql[pe.Position:], tc.want)
		}
	}
}

// TestConditionalCompilationSplit checks that Split sees the selected
// branch: an $END is no END, so a unit holding directives stays whole.
func TestConditionalCompilationSplit(t *testing.T) {
	sql := `CREATE OR REPLACE PACKAGE BODY k AS
$IF $$debug $THEN
  PROCEDURE p IS BEGIN NULL; END;
$ELSE
  PROCEDURE p IS BEGIN NULL; END;
$END
  PROCEDURE q IS
  BEGIN
    $if dbms_db_version.version < 12 $then $error 'unsupported' $end $end
    NULL;
  END;
END;
/
SELECT 1 FROM dual;
`
	got := splitTexts(Split(sql))
	if len(got) != 2 || !strings.HasSuffix(got[0], "END;\nEND;") || strings.TrimSpace(got[1]) != "SELECT 1 FROM dual" {
		t.Fatalf("Split = %q", got)
	}
	ParseAndCheck(t, got[0])
}
