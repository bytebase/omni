package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestConditionalCompilation checks PL/SQL conditional compilation against
// Oracle 23ai: every accepted form compiles there, and every rejected one
// fails with the error noted beside it. Rows ref_264 to ref_274 and ref_281
// to ref_283 run them against the engine.
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
		// Oracle 23ai folds static Boolean expressions and literal
		// comparisons, short-circuiting past what the text does not settle.
		"BEGIN $IF NOT TRUE $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF TRUE AND FALSE $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF FALSE OR (TRUE AND NOT FALSE) $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF NOT NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF TRUE OR $$x $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF FALSE AND $$x $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF 1 = 1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF 1 = 2 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF 1 + 1 = 3 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF 1 - 1 = 0 AND 10 / 2 > 4 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF NVL(NULL, TRUE) $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF XOR(TRUE, TRUE) $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF ABS(-1) = 1 AND pkg.c $THEN NULL; $END END;",
		// Decimal arithmetic and whole powers, as NUMBER computes them.
		"BEGIN $IF 0.1 + 0.2 = 0.3 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF 2 ** 3 = 8 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF 2 ** 3 = 9 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF NVL(1, 2) = 1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF $$PLSQL_UNIT_OWNER IS NULL $THEN NULL; $END END;",
		"BEGIN $IF $$PLSQL_LINE = 1 $THEN NULL; $END END;",
		// The static functions fold too, MOD with the dividend's sign.
		"BEGIN $IF ABS(-1) = 2 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF MOD(-11, 4) = -3 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF ROUND(-1.5) = -2 AND CEIL(1.2) = 2 AND FLOOR(-1.2) = -2 AND TRUNC(-1.5) = -1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF SIGN(-5) = -1 AND POWER(2, 10) = 1024 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF LENGTH('ab') = 2 AND UPPER('a') IS NULL AND ABS($$x) = 1 $THEN NULL; $END END;",
		// A quoted column named $$X is no inquiry directive.
		`DECLARE x NUMBER; BEGIN SELECT 1 INTO x FROM dual t WHERE t."$$X"(+) = 1; END;`,
		// A label before the block's BEGIN does not count as its first token.
		"<<l>> $IF TRUE $THEN BEGIN $END NULL; END;",
		// A unit's first directive follows a token of its own; a trigger's
		// follows the start of its body.
		"CREATE OR REPLACE PROCEDURE p (x NUMBER) $IF TRUE $THEN IS $ELSE AS $END BEGIN NULL; END;",
		"CREATE OR REPLACE FUNCTION f RETURN NUMBER $IF TRUE $THEN IS $END BEGIN RETURN 1; END;",
		"CREATE OR REPLACE TRIGGER trg BEFORE INSERT ON t DECLARE $IF TRUE $THEN x NUMBER; $END BEGIN NULL; END;",
		"CREATE OR REPLACE TRIGGER trg FOR INSERT ON t COMPOUND TRIGGER $IF TRUE $THEN BEFORE STATEMENT IS BEGIN NULL; END BEFORE STATEMENT; $END END;",
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
		// A condition is Boolean (PLS-00382), after the unit's END too.
		"BEGIN $IF 1 $THEN NULL; $END END;",
		"BEGIN $IF 'x' $THEN NULL; $END END;",
		"BEGIN $IF 1 + 2 $THEN NULL; $END END;",
		"BEGIN $IF $$x + 1 $THEN NULL; $END END;",
		"BEGIN NULL; END; $IF 1 + $THEN $END",
		// A call is no static Boolean but STANDARD's NVL and XOR (PLS-00382,
		// PLS-00174, PLS-00306), nor are CASE and DECODE.
		"BEGIN $IF ABS(1) $THEN NULL; $END END;",
		"BEGIN $IF TO_CHAR(1) $THEN NULL; $END END;",
		"BEGIN $IF NOT ABS(1) $THEN NULL; $END END;",
		"BEGIN $IF f(1) $THEN NULL; $END END;",
		"BEGIN $IF NVL2(TRUE, TRUE, FALSE) $THEN NULL; $END END;",
		"BEGIN $IF CASE WHEN 1 = 1 THEN TRUE ELSE FALSE END $THEN NULL; $END END;",
		"BEGIN $IF DECODE(1, 1, TRUE, FALSE) $THEN NULL; $END END;",
		"BEGIN $IF NVL(1, 2) $THEN NULL; $END END;",
		"BEGIN $IF NVL(TRUE, 1) $THEN NULL; $END END;",
		"BEGIN $IF pkg.NVL(TRUE, FALSE) $THEN NULL; $END END;",
		// A static expression calls the static functions alone (PLS-00174,
		// PLS-00201): no INSTR, GREATEST, TO_NUMBER, function of the user's,
		// pseudo-column, DECODE, CASE, or CAST.
		"BEGIN $IF f() = 1 $THEN NULL; $END END;",
		"BEGIN $IF INSTR('ab', 'b') = 2 $THEN NULL; $END END;",
		"BEGIN $IF GREATEST(1, 2) = 2 $THEN NULL; $END END;",
		"BEGIN $IF TO_NUMBER('1') = 1 $THEN NULL; $END END;",
		"BEGIN $IF SYSDATE IS NULL $THEN NULL; $END END;",
		"BEGIN $IF DECODE(1, 1, 1, 2) = 1 $THEN NULL; $END END;",
		"BEGIN $IF CASE WHEN 1 = 1 THEN 1 END = 1 $THEN NULL; $END END;",
		"BEGIN $IF CAST(1 AS NUMBER) = 1 $THEN NULL; $END END;",
		// A bind variable after the block's END is the block's too
		// (PLS-00110).
		"BEGIN NULL; END; $IF :x $THEN $END",
		// A trigger named COMPOUND does not start its body.
		"CREATE TRIGGER compound FOR INSERT ON t $IF TRUE $THEN COMPOUND TRIGGER $END BEFORE STATEMENT IS BEGIN NULL; END BEFORE STATEMENT; END;",
		// A static comparison compares no character values (PLS-00174).
		"BEGIN $IF $$PLSQL_UNIT = 'P' $THEN NULL; $END END;",
		"BEGIN $IF $$PLSQL_UNIT_TYPE = 'X' $THEN NULL; $END END;",
		"BEGIN $IF 'a' = 'a' $THEN NULL; $END END;",
		"BEGIN $IF $$x = 'a' $THEN NULL; $END END;",
		// No outer join on an inquiry directive (ORA-03048).
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO x FROM dual t WHERE $$PLSQL_LINE(+) = 1; END;",
		// The selected code is parsed: here, an empty block (PLS-00103).
		"BEGIN $IF NOT NVL(TRUE, FALSE) $THEN NULL; $END END;",
		// An anonymous block takes no bind variable inside a directive
		// (PLS-00110).
		"BEGIN $IF TRUE $THEN :n := 1; $END END;",
		"BEGIN $IF FALSE $THEN :n := 1; $ELSE NULL; $END END;",
		"BEGIN $IF :x $THEN NULL; $END END;",
		// An inquiry directive is a literal, not a name (PLS-00103).
		"DECLARE $$PLSQL_LINE NUMBER; BEGIN NULL; END;",
		"BEGIN $$x := 1; END;",
		"DECLARE x NUMBER; BEGIN x := $$PLSQL_LINE.y; END;",
		"BEGIN $$plsql_unit(1); END;",
		// A directive right after a unit's name, or before a trigger's body
		// (ORA-00911).
		"CREATE OR REPLACE PROCEDURE p $IF TRUE $THEN IS $END BEGIN NULL; END;",
		"CREATE OR REPLACE PROCEDURE p $IF TRUE $THEN (x NUMBER) $END IS BEGIN NULL; END;",
		"CREATE OR REPLACE PACKAGE k $IF TRUE $THEN AS $END c NUMBER; END;",
		"CREATE OR REPLACE PACKAGE BODY k $IF TRUE $THEN AS $END PROCEDURE p IS BEGIN NULL; END; END;",
		"CREATE OR REPLACE TRIGGER trg BEFORE INSERT ON t $IF TRUE $THEN FOR EACH ROW $END BEGIN NULL; END;",
		"CREATE OR REPLACE TRIGGER trg BEFORE INSERT ON t $IF TRUE $THEN BEGIN NULL; END; $END",
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

	// Directives after a unit's END are the unit's.
	got = splitTexts(Split("BEGIN NULL; END; $IF TRUE $THEN $END\nSELECT 1 FROM dual;"))
	if len(got) != 2 || got[0] != "BEGIN NULL; END; $IF TRUE $THEN $END" {
		t.Fatalf("Split = %q", got)
	}
	ParseAndCheck(t, got[0])

	// A malformed directive ends no statement early: Split lexes on past it,
	// past one whose condition is no expression too.
	got = splitTexts(Split("SELECT 1 FROM dual; $IF broken\nSELECT 2 FROM dual;\nSELECT 3 FROM dual;"))
	if len(got) != 3 || strings.TrimSpace(got[2]) != "SELECT 3 FROM dual" {
		t.Fatalf("Split = %q", got)
	}
	got = splitTexts(Split("SELECT 1 FROM dual; $IF broken; SELECT 2 FROM dual; $THEN SELECT 3 FROM dual; $END SELECT 4 FROM dual;"))
	if len(got) != 5 || !strings.HasSuffix(got[4], "SELECT 4 FROM dual") {
		t.Fatalf("Split = %q", got)
	}

	// Completion sees every branch.
	found := false
	for _, tok := range Tokenize("BEGIN $IF $$debug $THEN NULL; $ELSE SELECT * FROM emp; $END END;") {
		found = found || tok.Str == "EMP"
	}
	if !found {
		t.Errorf("Tokenize drops the $ELSE branch")
	}
}
