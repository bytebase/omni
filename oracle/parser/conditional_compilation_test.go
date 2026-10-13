package parser

import (
	"errors"
	"reflect"
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
		"BEGIN $IF ROUND(1.25, NULL) IS NULL AND 1 IS NOT NULL AND LENGTH('') IS NULL $THEN NULL; $ELSE garbage garbage; $END END;",
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
