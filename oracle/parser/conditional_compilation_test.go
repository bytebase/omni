package parser

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

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
		"BEGIN $IF LENGTH('ab') = 2 AND UPPER('a') IS NULL AND ABS($$x) = 1 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH('a') = 2 $THEN garbage garbage; $ELSE NULL; $END END;",
		// Static comparisons of numbers or Booleans, IS NULL of anything.
		"BEGIN $IF TRUE = TRUE AND $$x IS NOT NULL AND 1 IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH(1) = 1 AND LENGTH(UPPER('a')) = 1 AND CHR(65) || 'b' IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF NVL(1, 2) = NVL(3, 4) AND $$x = 1 + $$y AND NULL = 1 AND NOT NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		// Booleans compare, FALSE below TRUE, and a comparison with NULL is
		// NULL.
		"BEGIN $IF TRUE = FALSE $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF TRUE > FALSE AND (1 = 2) = (2 = 3) AND TRUE >= TRUE $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF NULL = TRUE OR NOT (NULL = 1) OR TRUE < FALSE $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF ROUND(1.25, NULL) IS NULL AND 1 IS NOT NULL AND LENGTH(NULL) IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF LENGTH(NULL) IS NOT NULL OR LENGTHB(NULL) = 1 OR ASCII(NULL) = 1 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF ASCII('a') = 97 AND LENGTHB('ab') = 2 AND TO_CHAR('') IS NULL AND NVL('', 'a') IS NOT NULL AND UPPER(NULL) IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
		// NUMBER holds 20 base-100 digits, rounded half away from zero, and
		// a literal rounds as a result does.
		"BEGIN $IF 1 / 3 * 3 = 1 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF 2 / 3 * 3 = 2 AND 1 / 3 = 0.3333333333333333333333333333333333333333 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF 1 / 3 = 0.333333333333333333333333333333333333333 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF 10 / 3 = 3.333333333333333333333333333333333333333 AND 1.00000000000000000000000000000000000000001 = 1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF 99999999999999999999999999999999999999995 = 1e41 AND 1e-130 / 10 = 0 $THEN NULL; $ELSE garbage garbage; $END END;",
		// ROUND and TRUNC to a scale, which rounds to a whole number.
		"BEGIN $IF ROUND(1.24, 1) = 1.3 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF ROUND(1.25, 1) = 1.3 AND ROUND(-1.25, 1) = -1.3 AND ROUND(125, -1) = 130 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF TRUNC(-1.29, 1) = -1.2 AND TRUNC(129, -1) = 120 AND ROUND(1.25, 1.5) = 1.25 AND TRUNC(1.29, 1.6) = 1.29 $THEN NULL; $ELSE garbage garbage; $END END;",
		// IS NULL folds a character value too: a literal, ||, CONCAT, NVL,
		// and the functions that keep their argument's nullness.
		"BEGIN $IF 'x' IS NULL OR 'a' || NULL IS NULL OR UPPER('a') IS NULL OR CHR(65) IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF CONCAT('a', NULL) IS NULL OR NVL($$x, 'a') IS NULL OR $$x || 'a' IS NULL OR 1 || 2 IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF NULL || NULL IS NULL AND TO_CHAR(1) IS NOT NULL AND UPPER(NULL) IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
		// An inquiry directive's kind is the session's: PLSQL_CCFLAGS may
		// set $$PLSQL_LINE to TRUE and $$PLSQL_UNIT to 1.
		"BEGIN $IF $$PLSQL_LINE $THEN NULL; $END END;",
		"BEGIN $IF $$PLSQL_UNIT = 1 $THEN NULL; $END END;",
		// Parentheses around a SELECT INTO or DML RETURNING target.
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO (x) FROM dual; UPDATE t SET a = 1 RETURNING a INTO (x); END;",
		// BITAND (two's complement, toward zero), REMAINDER (to the nearest
		// even quotient), and negative whole powers fold; a NULL operand
		// raises nothing, nor does a condition Oracle does not evaluate.
		"BEGIN $IF BITAND(1, 1) = 0 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF BITAND(6, 3) = 2 AND BITAND(-6, -3) = -8 AND BITAND(2.9, 3) = 2 AND REMAINDER(14, 4) = -2 AND REMAINDER(11, 4) = -1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF POWER(4, -2) = 0.0625 AND 2 ** (-1) = 0.5 AND NULL / 0 IS NULL AND REMAINDER(NULL, 0) IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF TRUE $THEN NULL; $ELSIF 1 / 0 = 0 $THEN NULL; $END END;",
		// The other admitted functions fold where their value is exact:
		// SQRT of an exact square, and EXP, LN, LOG, and the trigonometric
		// functions at 0 or 1.
		"BEGIN $IF SQRT(4) = 3 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF SQRT(2.25) = 1.5 AND EXP(0) = 1 AND LN(1) = 0 AND LOG(3, 1) = 0 AND COS(0) = 1 AND COSH(0) = 1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF SIN(0) = 0 AND TAN(0) = 0 AND ATAN(0) = 0 AND SINH(0) = 0 AND TANH(0) = 0 AND SQRT(NULL) IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF EXP(300) > 0 AND POWER(10, 125) > 0 AND POWER(10, -131) = 0 $THEN NULL; $ELSE NULL; $END END;",
		// SUBSTR, RPAD, and LPAD settle their nullness, the arguments
		// rounded to whole numbers.
		"BEGIN $IF SUBSTR('a', 1) IS NULL OR SUBSTR('abc', 0) IS NULL OR SUBSTR('abc', -1) IS NULL OR SUBSTR(123, 2) IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF SUBSTR('abc', 5) IS NULL AND SUBSTR('abc', -4) IS NULL AND SUBSTR('abc', 2, 0) IS NULL AND SUBSTR('abc', 3.9) IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF SUBSTR('abc', 2, 0.5) IS NULL OR RPAD('a', 0.5) IS NULL OR LPAD('a', 3, 'x') IS NULL OR RPAD(1, 2) IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF RPAD('a', 0) IS NULL AND LPAD('a', 3, NULL) IS NULL AND LPAD(NULL, 3, 'x') IS NULL AND SUBSTR('abc', NULL) IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
		// Character calls fold, nested too, over ASCII values.
		"BEGIN $IF LENGTH(UPPER('a')) = 2 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH('a' || 'bc') = 3 AND ASCII(LOWER('A')) = 97 AND LENGTH(INITCAP('ab cd')) = 5 AND LENGTH(SUBSTR('abc', 2)) = 2 AND LENGTH(RPAD('a', 3)) = 3 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF LENGTH(TO_CHAR(123)) = 3 AND LENGTH(CHR(65) || 'b') = 2 AND LENGTH(NVL(NULL, 'ab')) = 2 AND LENGTH(1 || 2) = 2 AND LENGTH(LPAD('ab', 1)) = 1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF ASCII(SUBSTR(INITCAP('ab-cd'), 4)) = 67 AND ASCII(SUBSTR(INITCAP('aBC dEF'), 2)) = 98 AND ASCII(SUBSTR(INITCAP('a1b'), 3)) = 98 AND ASCII(SUBSTR(LPAD('a', 4, 'xy'), 2)) = 121 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF ASCII(SUBSTR(INITCAP('a_b'), 3)) = 98 $THEN garbage garbage; $ELSE NULL; $END END;",
		// A whole power folds exactly whatever its exponent, and TO_CHAR and
		// || convert a PLS_INTEGER: an integer literal, +, -, and * of
		// them, and ABS, SIGN, BITAND, NVL, LENGTH, and ASCII.
		"BEGIN $IF 1 ** 65 = 2 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF POWER(1, 65) = 1 AND (-1) ** 65 = -1 AND 0 ** 65 = 0 AND 2 ** 100 = 1267650600228229401496703205376 AND POWER(-1, -65) = -1 $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF TO_CHAR(1 + 1) IS NULL OR TO_CHAR(2 * 3) IS NULL OR TO_CHAR(ABS(-1)) IS NULL OR TO_CHAR(LENGTH('ab')) IS NULL OR TO_CHAR(-(1)) IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(1 - 2) IS NULL OR TO_CHAR(NVL(1, 2)) IS NULL OR 1 + 1 || 'a' IS NULL OR TO_CHAR(SIGN(-1)) IS NULL OR TO_CHAR(BITAND(6, 3)) IS NULL OR UPPER(1 + 1) IS NULL $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(NULL + 1) IS NULL AND LENGTH(TO_CHAR(2 * 3 - 1)) = 1 $THEN NULL; $ELSE garbage garbage; $END END;",
		// Arithmetic past NUMBER's range raises nothing.
		"BEGIN $IF 1e125 * 10 > 0 AND 9e125 + 9e125 > 0 AND 1e125 / 1e-5 > 0 $THEN NULL; $ELSE NULL; $END END;",
		// A hint anywhere but after SELECT, INSERT, UPDATE, DELETE, or MERGE
		// is a comment.
		"BEGIN $IF /*+ ignored */ TRUE /*+ x */ AND TRUE $THEN /*+ y */ NULL; $END END;",
		// A constant is a package's, named through it.
		"BEGIN $IF pkg.b AND system.pkg.b AND \"PKG\".\"B\" $THEN NULL; $END END;",
		// An inquiry directive is a value an IN bind or argument takes.
		"BEGIN p(a => $$PLSQL_LINE); EXECUTE IMMEDIATE 'BEGIN NULL; END;' USING IN $$PLSQL_LINE; END;",
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
		// An operation that raises makes no static expression, evaluated
		// or not (PLS-00174): a division or REMAINDER by zero, zero to a
		// negative power, a negative number to a fractional one, and BITAND
		// past -2^127 through 2^127 - 1.
		"BEGIN $IF 1 / 0 = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF FALSE AND 1 / 0 = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TRUE OR 1 / (1 - 1) = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF FALSE $THEN NULL; $ELSIF 1 / 0 = 0 $THEN NULL; $END END;",
		"BEGIN $IF REMAINDER(1, 0) = 1 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF POWER(0, -1) = 1 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF POWER(-8, 1 / 3) < 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF BITAND(170141183460469231731687303715884105728, 1) = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF BITAND(-170141183460469231731687303715884105729, 1) = 0 $THEN NULL; $ELSE NULL; $END END;",
		// So does SQRT, LN, or LOG out of its domain, a power from 1e126 on,
		// CHR of less than 1, and SUBSTR, RPAD, or LPAD arguments past
		// PLS_INTEGER or a pad past 32767.
		"BEGIN $IF SQRT(-1) = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LN(0) = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LOG(1, 2) = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LOG(2, -4) = 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF POWER(10, 126) > 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF POWER(-10, 127) < 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF CHR(0) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF CHR(-0.4) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF SUBSTR('a', 2147483647.5) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF SUBSTR('a', 1, -2147483649) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF RPAD('a', 32768) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LPAD('a', 2147483648) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		// An empty string raises as an operand but of TO_CHAR, TO_NCHAR,
		// and NVL.
		"BEGIN $IF UPPER('') IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH('') IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF ASCII('') IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF 'a' || '' IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF SUBSTR('a', '') IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH(TO_CHAR(q'[]')) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH(NVL((NULL), '')) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		// A character context converts a PLS_INTEGER literal alone.
		"BEGIN $IF TO_CHAR(0.5) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(1.0) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH(TO_CHAR(2147483648)) = 10 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH(TO_CHAR(-2147483649)) = 11 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH(1e3 || 'a') = 5 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF LENGTH(UPPER(0.5)) = 2 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(0.5 + 0) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(4 / 2) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(MOD(5, 2)) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(2 ** 2) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(ROUND(1.5)) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(-2147483648) IS NOT NULL $THEN NULL; $END END;",
		// PLS_INTEGER arithmetic past 32 bits raises.
		"BEGIN $IF 2147483647 + 1 > 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF 2147483647 * 2 > 0 $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF TO_CHAR(2147483647 + 1) IS NULL $THEN NULL; $ELSE NULL; $END END;",
		// ** takes no signed operand (PLS-00103).
		"BEGIN $IF 2 ** -1 = 0.5 $THEN NULL; $ELSE NULL; $END END;",
		"DECLARE x NUMBER; BEGIN x := 2 ** -1; END;",
		// TO_CHAR takes one argument (PLS-00174 for its format overloads).
		"BEGIN $IF TO_CHAR(1, '9', 'NLS_NUMERIC_CHARACTERS = ''.,''') IS NULL $THEN NULL; $ELSE NULL; $END END;",
		// A bind variable after the block's END is the block's too
		// (PLS-00110), when another statement follows.
		"BEGIN NULL; END; $IF :x $THEN $END",
		"BEGIN NULL; END; $IF :x $THEN $END SELECT 1 FROM dual;",
		// A constant is a package's, named through it (PLS-00201), in the
		// package's own body too.
		"BEGIN $IF debug $THEN NULL; $END END;",
		`BEGIN $IF "DEBUG" $THEN NULL; $END END;`,
		"BEGIN $IF TRUE AND debug $THEN NULL; $END END;",
		"BEGIN $IF debug = 1 $THEN NULL; $END END;",
		"CREATE OR REPLACE PACKAGE BODY pkg AS $IF b $THEN x NUMBER; $END END;",
		// A trigger named COMPOUND does not start its body.
		"CREATE TRIGGER compound FOR INSERT ON t $IF TRUE $THEN COMPOUND TRIGGER $END BEFORE STATEMENT IS BEGIN NULL; END BEFORE STATEMENT; END;",
		// The static grammar: no IN, BETWEEN, LIKE, or EXISTS (PLS-00174,
		// PLS-00204), each static function with its arguments (PLS-00306,
		// PLS-00174), and kinds that agree (PLS-00306, PLS-00382).
		"BEGIN $IF 1 IN (1, 2) $THEN NULL; $END END;",
		"BEGIN $IF 1 BETWEEN 0 AND 2 $THEN NULL; $END END;",
		"BEGIN $IF 'a' LIKE 'a' $THEN NULL; $END END;",
		"BEGIN $IF EXISTS (SELECT 1 FROM dual) $THEN NULL; $END END;",
		"BEGIN $IF ABS(TRUE) = 1 $THEN NULL; $END END;",
		"BEGIN $IF ABS(1, 2) = 1 $THEN NULL; $END END;",
		"BEGIN $IF ABS('1') = 1 $THEN NULL; $END END;",
		"BEGIN $IF TO_CHAR(1, '9') IS NULL $THEN NULL; $END END;",
		"BEGIN $IF ROUND(1, 2, 3) = 1 $THEN NULL; $END END;",
		"BEGIN $IF SUBSTR('a') IS NULL $THEN NULL; $END END;",
		"BEGIN $IF NVL(1, TRUE) = 1 $THEN NULL; $END END;",
		"BEGIN $IF 1 + TRUE = 1 $THEN NULL; $END END;",
		"BEGIN $IF NOT 1 $THEN NULL; $END END;",
		"BEGIN $IF TRUE AND 1 $THEN NULL; $END END;",
		"BEGIN $IF 1 < TRUE $THEN NULL; $END END;",
		"BEGIN $IF TO_CHAR(1) = TO_CHAR(1) $THEN NULL; $END END;",
		"BEGIN $IF UPPER('a') = 'A' $THEN NULL; $END END;",
		// An inquiry directive is no cursor (PLS-00103).
		"BEGIN IF $$PLSQL_LINE%FOUND THEN NULL; END IF; END;",
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
		// An inquiry directive is a literal, not a name (PLS-00103): no
		// parameter name (PLS-00306), INTO target (PLS-00103, ORA-32551,
		// PLS-00308), or OUT bind.
		"DECLARE $$PLSQL_LINE NUMBER; BEGIN NULL; END;",
		"BEGIN p($$PLSQL_LINE => 1); END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO $$PLSQL_LINE FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO x, $$X FROM dual; END;",
		"BEGIN EXECUTE IMMEDIATE 'SELECT 1 FROM dual' INTO $$PLSQL_LINE; END;",
		"DECLARE c SYS_REFCURSOR; BEGIN FETCH c INTO $$PLSQL_LINE; END;",
		"BEGIN UPDATE t SET a = 1 RETURNING a INTO $$PLSQL_LINE; END;",
		"BEGIN EXECUTE IMMEDIATE 'BEGIN NULL; END;' USING OUT $$X; END;",
		"BEGIN EXECUTE IMMEDIATE 'BEGIN NULL; END;' USING IN OUT $$X; END;",
		// Parenthesized too.
		"DECLARE x NUMBER; BEGIN x := ($$PLSQL_LINE).y; END;",
		"BEGIN IF ($$PLSQL_LINE)%FOUND THEN NULL; END IF; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO ($$PLSQL_LINE) FROM dual; END;",
		"BEGIN UPDATE t SET a = 1 RETURNING a INTO (($$PLSQL_LINE)); END;",
		// No parentheses around a FETCH, EXECUTE IMMEDIATE, or OUT bind
		// target (PLS-00103).
		"DECLARE x NUMBER; c SYS_REFCURSOR; BEGIN FETCH c INTO (x); END;",
		"DECLARE x NUMBER; BEGIN EXECUTE IMMEDIATE 'SELECT 1 FROM dual' INTO (x); END;",
		"DECLARE x NUMBER; BEGIN EXECUTE IMMEDIATE 'BEGIN NULL; END;' RETURNING INTO (x); END;",
		"DECLARE x NUMBER; BEGIN EXECUTE IMMEDIATE 'BEGIN :x := 1; END;' USING OUT (x); END;",
		"BEGIN $$x := 1; END;",
		"DECLARE x NUMBER; BEGIN x := $$PLSQL_LINE.y; END;",
		"BEGIN $$plsql_unit(1); END;",
		// A directive right after a unit's name, or before a trigger's body
		// (ORA-00911).
		"CREATE OR REPLACE PROCEDURE p $IF TRUE $THEN IS $END BEGIN NULL; END;",
		"CREATE OR REPLACE PROCEDURE $IF TRUE $THEN p IS $END BEGIN NULL; END;",
		"CREATE OR REPLACE PROCEDURE s.$IF TRUE $THEN p IS $END BEGIN NULL; END;",
		"CREATE OR REPLACE FUNCTION $IF TRUE $THEN f $END RETURN NUMBER IS BEGIN RETURN 1; END;",
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

	// A removed directive leaves a hint after its SELECT.
	for _, sql := range []string{
		"DECLARE x NUMBER; BEGIN SELECT $IF TRUE $THEN $END /*+ FULL(dual) */ 1 INTO x FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT $IF FALSE $THEN 2 $END /*+ FULL(dual) */ 1 INTO x FROM dual; END;",
	} {
		hints := 0
		ast.Inspect(parseOne[*ast.PLSQLBlock](t, sql), func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectStmt); ok && sel.Hints != nil {
				hints += sel.Hints.Len()
			}
			return true
		})
		if hints != 1 {
			t.Errorf("%q keeps %d hints, want 1", sql, hints)
		}
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
		// A condition that is no static Boolean fails before its branch.
		{"BEGIN $IF 1 $THEN garbage garbage; $END END;", "1 $THEN"},
		{"BEGIN $IF FALSE $THEN NULL; $ELSIF 1 / 0 = 0 $THEN garbage garbage; $END END;", "1 / 0"},
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
	// To their $END, with the code they select, which Oracle 23ai then
	// rejects (PLS-00103 at SELECT) unless another selection leaves none.
	got = splitTexts(Split("BEGIN NULL; END; $IF TRUE $THEN SELECT 1 FROM dual; $END\nSELECT 2 FROM dual;"))
	if len(got) != 2 || got[0] != "BEGIN NULL; END; $IF TRUE $THEN SELECT 1 FROM dual; $END" {
		t.Fatalf("Split = %q", got)
	}
	ParseShouldFail(t, got[0])
	got = splitTexts(Split("BEGIN NULL; END; $IF $$x $THEN SELECT 1 FROM dual; $ELSE $IF $$y $THEN NULL; $END $END SELECT 2 FROM dual;"))
	if len(got) != 2 || !strings.HasSuffix(got[0], "$END $END") || strings.TrimSpace(got[1]) != "SELECT 2 FROM dual" {
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

	// Nor one whose condition is no static Boolean.
	got = splitTexts(Split("SELECT 1 FROM dual; $IF 1 $THEN SELECT 2 FROM dual; $ELSE SELECT 3 FROM dual; $END SELECT 4 FROM dual;"))
	if len(got) != 4 || !strings.HasSuffix(got[3], "SELECT 4 FROM dual") {
		t.Fatalf("Split = %q", got)
	}

	// Outside a unit, a directive is stray however well formed.
	got = splitTexts(Split("SELECT 1 FROM dual; $IF TRUE $THEN SELECT 2 FROM dual; $ELSE SELECT 3 FROM dual; $END SELECT 4 FROM dual;"))
	if len(got) != 4 || strings.TrimSpace(got[2]) != "$ELSE SELECT 3 FROM dual" || !strings.HasSuffix(got[3], "SELECT 4 FROM dual") {
		t.Fatalf("Split = %q", got)
	}

	// Nor a stray $ERROR.
	got = splitTexts(Split("SELECT 1 FROM dual; $ERROR 'x'; SELECT 2 FROM dual; $END SELECT 3 FROM dual;"))
	if len(got) != 4 || strings.TrimSpace(got[2]) != "SELECT 2 FROM dual" || !strings.HasSuffix(got[3], "SELECT 3 FROM dual") {
		t.Fatalf("Split = %q", got)
	}

	// Tokenize keeps every branch.
	found := false
	for _, tok := range Tokenize("BEGIN $IF $$debug $THEN NULL; $ELSE SELECT * FROM emp; $END END;") {
		found = found || tok.Str == "EMP"
	}
	if !found {
		t.Errorf("Tokenize drops the $ELSE branch")
	}
}

// TestConditionalCompilationCompletion checks that completion reads one
// variant of the code: the branch that holds the caret, and elsewhere the
// branch the preprocessor selects.
func TestConditionalCompilationCompletion(t *testing.T) {
	refs := func(sql string, caret int) []string {
		var names []string
		for _, ref := range CollectCompletion(sql, caret).Scope.LocalReferences {
			names = append(names, ref.Name+" "+ref.Alias)
		}
		return names
	}
	for _, tc := range []struct {
		before, after string
		want          []string
	}{
		{"BEGIN SELECT * INTO x FROM $IF $$x $THEN emp $ELSE dept d WHERE d.", " $END; END;", []string{"DEPT D"}},
		{"BEGIN SELECT * INTO x FROM $IF $$x $THEN emp e WHERE e.", " $ELSE dept d $END; END;", []string{"EMP E"}},
		{"BEGIN SELECT * INTO x FROM $IF FALSE $THEN emp $ELSE dept $END d WHERE d.", "; END;", []string{"DEPT D"}},
		{"BEGIN SELECT * INTO x FROM $IF $$x $THEN emp $ELSE dept $END d WHERE d.", "; END;", []string{"EMP D"}},
		{"BEGIN SELECT * INTO x FROM $IF $$x $THEN emp $ELSIF $$y $THEN dept d WHERE d.", "", []string{"DEPT D"}},
	} {
		sql := tc.before + tc.after
		if got := refs(sql, len(tc.before)); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("CollectCompletion(%q) scope = %q, want %q", sql, got, tc.want)
		}
	}

	// Nested directives are read once: 5000 levels complete at once.
	deep := "BEGIN " + strings.Repeat("$IF TRUE $THEN ", 5000) + "SELECT * INTO x FROM t WHERE t." + strings.Repeat(" $END", 5000) + " END;"
	start := time.Now()
	flat := "BEGIN SELECT * INTO x FROM t WHERE t. END;"
	if got, want := refs(deep, strings.Index(deep, "t.")+2), refs(flat, strings.Index(flat, "t.")+2); !reflect.DeepEqual(got, want) {
		t.Errorf("deep scope = %q, want %q", got, want)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("completion through 5000 nested directives took %v", elapsed)
	}

	// A caret in a condition or an $ERROR message completes it.
	want := CollectCompletion("BEGIN IF pkg.", len("BEGIN IF pkg.")).Intent
	for _, tc := range [][2]string{
		{"BEGIN $IF pkg.", ""},
		{"BEGIN $IF pkg.", " $THEN NULL; $END END;"},
		{"BEGIN $IF TRUE $THEN NULL; $ELSIF pkg.", " $THEN NULL; $END END;"},
		{"BEGIN $IF $$x $THEN $ERROR pkg.", " $END $END END;"},
		{"BEGIN $ERROR 'a' || pkg.", ""},
	} {
		if got := CollectCompletion(tc[0]+tc[1], len(tc[0])).Intent; !reflect.DeepEqual(got, want) {
			t.Errorf("CollectCompletion(%q) intent = %+v, want %+v", tc[0]+tc[1], got, want)
		}
	}
}

// TestConditionalCompilationSelections checks that a condition the text
// does not settle leaves its branches to the environment: a unit parses
// when any selection of them does, and fails when every selection does,
// as Oracle 23ai compiles it under whichever selection its database
// version, PLSQL_CCFLAGS, and package constants settle.
func TestConditionalCompilationSelections(t *testing.T) {
	for _, sql := range []string{
		"BEGIN $IF $$x $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF $$x $THEN NULL; $ELSE garbage garbage; $END END;",
		"BEGIN $IF $$x $THEN garbage garbage; $ELSIF $$y $THEN garbage; $ELSE NULL; $END END;",
		"BEGIN $IF $$x $THEN garbage garbage; $END NULL; END;",
		"BEGIN $IF $$x $THEN $ERROR 'x' $END $END NULL; END;",
		"BEGIN $IF $$x $THEN $ERROR 'x' $END $ELSIF $$y $THEN $ERROR 'y' $END $END NULL; END;",
		"BEGIN $IF $$a $THEN $IF $$b $THEN garbage garbage; $ELSE NULL; $END $ELSE garbage garbage; $END END;",
		"BEGIN $IF $$a $THEN NULL; $END $IF $$b $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF $$a $THEN garbage garbage; $ELSE NULL; $END $IF $$b $THEN NULL; $END END;",
		// Whatever the text cannot settle: an inexact function, arithmetic
		// past NUMBER's range.
		"BEGIN $IF SQRT(2) < 0 $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF 1e125 * 10 < 0 $THEN garbage garbage; $ELSE NULL; $END END;",
		// A selection that leaves the delimiters unbalanced or the unit
		// heading broken is passed over.
		"BEGIN $IF $$x $THEN f( $ELSE NULL; $END END;",
		"CREATE OR REPLACE PROCEDURE p (x NUMBER) $IF $$x $THEN garbage $ELSE IS $END BEGIN NULL; END;",
		"DECLARE $IF $$x $THEN a NUMBER NUMBER; $ELSE a NUMBER; $END BEGIN NULL; END;",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"BEGIN $IF $$x $THEN $ERROR 'x' $END $ELSE $ERROR 'y' $END $END NULL; END;",
		"BEGIN $IF $$x $THEN garbage garbage; $ELSE garbage garbage; $END END;",
		"BEGIN $IF $$x $THEN $ERROR 'x' $END $ELSIF $$y $THEN garbage garbage; $ELSE garbage garbage; $END END;",
		"BEGIN $IF $$a $THEN $IF $$b $THEN garbage garbage; $ELSE garbage garbage; $END $ELSE garbage garbage; $END END;",
		// A settled condition is no selection: FALSE passes its branch.
		"BEGIN $IF TRUE $THEN garbage garbage; $ELSE NULL; $END END;",
		"BEGIN $IF $$x $THEN NULL; $END garbage garbage; END;",
	} {
		ParseShouldFail(t, sql)
	}

	// The tree is the selection that parsed.
	block := parseOne[*ast.PLSQLBlock](t, "DECLARE $IF $$x $THEN a NUMBER NUMBER; $ELSE a NUMBER; b DATE; $END BEGIN NULL; END;")
	if block.Declarations.Len() != 2 {
		t.Errorf("declarations = %d, want the 2 of the $ELSE branch", block.Declarations.Len())
	}

	// When every selection fails, the error is the one that got furthest.
	const sql = "BEGIN $IF $$x $THEN garbage garbage; $ELSE NULL; NULL $END END;"
	_, err := Parse(sql)
	var pe *ParseError
	if !errors.As(err, &pe) || !strings.HasPrefix(sql[pe.Position:], "$END END") && !strings.HasPrefix(sql[pe.Position:], "END;") {
		t.Errorf("Parse(%q) = %v, want the $ELSE branch's error", sql, err)
	}

	// The parentheses balance under the selection that parses: one open
	// outside every branch fails under all of them (PLS-00103 on Oracle
	// 23ai).
	ParseShouldFail(t, "DECLARE x NUMBER; BEGIN $IF $$x $THEN NULL; $END x := (1; END;")
	ParseShouldFail(t, "DECLARE x NUMBER; BEGIN $IF $$x $THEN NULL; $END x := 1); END;")

	// The $IF whose branch holds the failure varies first, so a selection
	// past the first 64 in order is found: here the outer $ELSE, after the
	// 64 selections of six inner $IF directives. An $ERROR is met under
	// every selection that keeps the branches around it.
	six := strings.Repeat("$IF $$i $THEN NULL; $END ", 6)
	ParseAndCheck(t, "BEGIN $IF $$o $THEN "+six+"$ERROR 'e' $END $ELSE NULL; $END END;")
	ParseAndCheck(t, "BEGIN $IF $$o $THEN "+six+"garbage garbage; $ELSE NULL; $END END;")
	ParseShouldFail(t, "BEGIN $IF $$o $THEN "+six+"$ERROR 'e' $END $ELSE "+six+"$ERROR 'f' $END $END END;")

	// The CASE pre-scan reads the selection being parsed, where a column
	// may be named CASE (Oracle 23ai compiles this block).
	ParseAndCheck(t, "DECLARE x NUMBER; BEGIN x := CASE WHEN 1=1 THEN 1 END; $IF $$flag $THEN $ERROR 'bad' $END $ELSE SELECT case + 1 INTO x FROM t; $END END;")

	// A string literal holds 32767 bytes, in code a directive skips too
	// (PLS-00172, ORA-01704), and a static concatenation as many (PLS-00174).
	long := strings.Repeat("a", 32767)
	ParseAndCheck(t, "BEGIN $IF '"+long+"' IS NULL $THEN NULL; $ELSE NULL; $END END;")
	for _, sql := range []string{
		"BEGIN $IF '" + long + "a' IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"DECLARE x CLOB; BEGIN x := '" + long + "a'; END;",
		"BEGIN $IF $$x $THEN NULL; $ELSE x := '" + long + "a'; $END END;",
		"SELECT '" + long + "a' FROM dual",
		"BEGIN $IF '" + long + "' || 'b' IS NULL $THEN NULL; $ELSE NULL; $END END;",
		"BEGIN $IF RPAD('a', 32767) || 'b' IS NULL $THEN NULL; $ELSE NULL; $END END;",
	} {
		ParseShouldFail(t, sql)
	}
	ParseAndCheck(t, "BEGIN NULL; /* "+long+"a */ END;")
	if got := splitTexts(Split("SELECT '" + long + "a' FROM dual; SELECT 1 FROM dual;")); len(got) != 2 {
		t.Errorf("Split past a long literal = %d segments, want 2", len(got))
	}

	// The selections tried are bounded: twenty undecided directives before
	// an error fail at once.
	many := "BEGIN " + strings.Repeat("$IF $$x $THEN NULL; $END ", 20) + "garbage garbage; END;"
	start := time.Now()
	ParseShouldFail(t, many)
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Parse through 20 undecided directives took %v", elapsed)
	}
}

// TestConditionalCompilationLimits checks, against Oracle 23ai, that
// directives nest without bound, that an INTO clause assigns no value, and
// the range of numeric literals in code and in conditions.
func TestConditionalCompilationLimits(t *testing.T) {
	// Oracle 23ai compiles 5000 nested $IF TRUE levels.
	for _, n := range []int{65, 100, 5000} {
		ParseAndCheck(t, "BEGIN "+strings.Repeat("$IF TRUE $THEN ", n)+"NULL; "+strings.Repeat("$END ", n)+"END;")
		ParseAndCheck(t, "BEGIN "+strings.Repeat("$IF $$x $THEN ", n)+"NULL; "+strings.Repeat("$END ", n)+"END;")
	}

	// An INTO target is no value, parenthesized or not (ORA-32551,
	// PLS-00403, PLS-00103).
	for _, sql := range []string{
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO $$PLSQL_LINE + x FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO 1 + x FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO x + 1 FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO -x FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO (x + 1) FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO (1) FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO NULL FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO CASE WHEN 1=1 THEN x END FROM dual; END;",
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO CAST(x AS NUMBER) FROM dual; END;",
		"DECLARE x NUMBER; BEGIN EXECUTE IMMEDIATE 'SELECT 1 FROM dual' INTO 1 + x; END;",
		"DECLARE x NUMBER; c SYS_REFCURSOR; BEGIN FETCH c INTO x + 1; END;",
		"DECLARE x NUMBER; BEGIN DELETE FROM t RETURNING 1 INTO x || 'a'; END;",
	} {
		ParseShouldFail(t, sql)
	}
	for _, sql := range []string{
		"DECLARE x NUMBER; BEGIN SELECT 1 INTO (x) FROM dual; END;",
		"DECLARE TYPE t IS TABLE OF NUMBER INDEX BY PLS_INTEGER; v t; BEGIN SELECT 1 INTO v(1) FROM dual; END;",
		"DECLARE TYPE r IS RECORD (f NUMBER); v r; BEGIN SELECT 1 INTO v.f FROM dual; END;",
		"BEGIN SELECT 1 INTO :b FROM dual; END;",
	} {
		ParseAndCheck(t, sql)
	}

	// A numeric literal past its type's range fails in selected code, SQL,
	// and conditions (PLS-00569, ORA-01426), not in code a directive skips:
	// a NUMBER one that rounds to 1e126 or more, where 40 nines and a 5
	// round up, and a BINARY_FLOAT or BINARY_DOUBLE one that IEEE rounding
	// makes infinite or, from a nonzero value, zero.
	nines := func(n int) string { return strings.Repeat("9", n) }
	for _, sql := range []string{
		"DECLARE x NUMBER; BEGIN x := 1e126; END;",
		"DECLARE x NUMBER; BEGIN x := -1e126; END;",
		"DECLARE x NUMBER; BEGIN x := 9." + nines(39) + "5e125; END;",
		"DECLARE x NUMBER; BEGIN x := " + nines(127) + "; END;",
		"DECLARE x NUMBER; BEGIN x := 0." + strings.Repeat("0", 200) + "1e327; END;",
		"DECLARE x NUMBER; BEGIN $IF 1e126 > 0 $THEN NULL; $ELSE NULL; $END END;",
		"DECLARE x NUMBER; BEGIN $IF -1e126 < 0 $THEN NULL; $ELSE NULL; $END END;",
		"SELECT 1e126 FROM dual",
		"SELECT 1 FROM dual WHERE 1 = 0 AND 9.9999999999999999999999999999999999999999999e125 > 0",
		"CREATE VIEW v AS SELECT 1e126 x FROM dual",
		"DECLARE x BINARY_DOUBLE; BEGIN x := 1.7976931348623159e308d; END;",
		"DECLARE x BINARY_DOUBLE; BEGIN x := 2e-324d; END;",
		"DECLARE x BINARY_FLOAT; BEGIN x := 3.40282357e38f; END;",
		"DECLARE x BINARY_FLOAT; BEGIN x := 1e-46f; END;",
		"DECLARE x NUMBER; BEGIN $IF 1e999d > 0 $THEN NULL; $ELSE NULL; $END END;",
		"SELECT 1e999d FROM dual",
		"SELECT 1 FROM dual WHERE 1 = 0 AND 1e39f > 0",
	} {
		ParseShouldFail(t, sql)
	}
	for _, sql := range []string{
		"DECLARE x NUMBER; BEGIN x := 9.99e125; END;",
		"DECLARE x NUMBER; BEGIN x := 9." + nines(38) + "5e125; END;",
		"DECLARE x NUMBER; BEGIN x := 9." + nines(39) + "4" + nines(10) + "e125; END;",
		"DECLARE x NUMBER; BEGIN x := 1e-131; END;",
		"DECLARE x BINARY_DOUBLE; BEGIN x := 1.7976931348623158e308d; END;",
		"DECLARE x BINARY_DOUBLE; BEGIN x := 5e-324d; END;",
		"DECLARE x BINARY_DOUBLE; BEGIN x := 0.0e-999d; END;",
		"DECLARE x BINARY_FLOAT; BEGIN x := 3.40282356e38f; END;",
		"DECLARE x BINARY_FLOAT; BEGIN x := 1e-40f; END;",
		"DECLARE x NUMBER; BEGIN $IF FALSE $THEN x := 1e999d; $END NULL; END;",
		"DECLARE x NUMBER; BEGIN $IF 9.99e125 > 0 $THEN NULL; $ELSE NULL; $END END;",
		"DECLARE x NUMBER; BEGIN $IF 1e-131 = 0 $THEN NULL; $ELSE NULL; $END END;",
		"DECLARE x NUMBER; BEGIN $IF FALSE $THEN x := 1e126; $END NULL; END;",
		"DECLARE x NUMBER; BEGIN $IF $$x $THEN x := 1e126; $END NULL; END;",
		"SELECT 1e126d FROM dual",
	} {
		ParseAndCheck(t, sql)
	}
}
