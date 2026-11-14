package compat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// OracleParse is the PL/SQL block an Oracle judge runs for each statement,
// as the user that owns its scratch schema. :1 is the statement, and :2 a
// DDL statement that builds the schema's fixture, or NULL for none.
//
// DBMS_SQL.PARSE parses a statement and runs it only when it is DDL
// (https://docs.oracle.com/en/database/oracle/oracle-database/23/arpls/DBMS_SQL.html,
// "PARSE Procedures"). So after a statement it ran, the block drops
// everything the user owns, purges the recycle bin, and builds the fixture
// again: the next statement finds the schema as the first did. It tells
// what ran by the statement's SQL function code
// (DBMS_SQL.LAST_SQL_FUNCTION_CODE): queries, DML, transaction and session
// control, and the other statements in oracleParsedOnly are parsed and not
// run.
//
// The block answers nothing when the parse succeeds. Otherwise it raises
// ORA-20999 with what OracleVerdict reads: the error's code, its position
// (DBMS_SQL.LAST_ERROR_POSITION, a byte offset into the statement), and the
// PL/SQL compiler's errors: for ORA-24344, a unit created with compilation
// errors, the USER_ERRORS message numbers (PLS-nnnnn; 0 for "Statement
// ignored", which follows another), and for ORA-06550, an anonymous block
// that did not compile, the PLS-nnnnn numbers of its error stack.
const OracleParse = `DECLARE
  c INTEGER := DBMS_SQL.OPEN_CURSOR;
  code PLS_INTEGER := 0;
  pos PLS_INTEGER := -1;
  fn PLS_INTEGER := 0;
  compile VARCHAR2(4000);
  stack VARCHAR2(4000);
  left_over VARCHAR2(4000);
BEGIN
  BEGIN
    DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
  EXCEPTION
    WHEN OTHERS THEN
      code := -SQLCODE;
      pos := DBMS_SQL.LAST_ERROR_POSITION;
      stack := DBMS_UTILITY.FORMAT_ERROR_STACK;
  END;
  fn := DBMS_SQL.LAST_SQL_FUNCTION_CODE;
  DBMS_SQL.CLOSE_CURSOR(c);
  IF code = 24344 THEN
    SELECT LISTAGG(message_number, ',') WITHIN GROUP (ORDER BY name, type, sequence)
      INTO compile FROM user_errors WHERE attribute = 'ERROR';
  ELSIF code = 6550 THEN
    FOR i IN 1 .. REGEXP_COUNT(stack, 'PLS-[0-9]{5}') LOOP
      compile := compile || CASE WHEN i > 1 THEN ',' END || TO_NUMBER(SUBSTR(REGEXP_SUBSTR(stack, 'PLS-[0-9]{5}', 1, i), 5));
    END LOOP;
  END IF;
  IF code IN (0, 24344) AND fn NOT IN (` + oracleParsedOnly + `) THEN
    FOR pass IN 1 .. 3 LOOP
      FOR o IN (SELECT object_name, object_type FROM user_objects
                 WHERE object_type NOT IN ('LOB', 'LOB PARTITION', 'LOB SUBPARTITION', 'TABLE PARTITION',
                                           'TABLE SUBPARTITION', 'INDEX PARTITION', 'INDEX SUBPARTITION')
                 ORDER BY DECODE(object_type, 'MATERIALIZED VIEW', 1, 'TRIGGER', 2, 'VIEW', 3, 'TABLE', 5, 4)) LOOP
        BEGIN
          EXECUTE IMMEDIATE 'DROP ' || o.object_type || ' "' || o.object_name || '"' ||
            CASE o.object_type
              WHEN 'TABLE' THEN ' CASCADE CONSTRAINTS PURGE'
              WHEN 'TYPE' THEN ' FORCE'
              WHEN 'OPERATOR' THEN ' FORCE'
              WHEN 'INDEXTYPE' THEN ' FORCE'
              WHEN 'CLUSTER' THEN ' INCLUDING TABLES CASCADE CONSTRAINTS'
            END;
        EXCEPTION
          WHEN OTHERS THEN NULL;
        END;
      END LOOP;
    END LOOP;
    EXECUTE IMMEDIATE 'PURGE RECYCLEBIN';
    SELECT LISTAGG(object_type || ' ' || object_name, ', ') WITHIN GROUP (ORDER BY object_type, object_name)
      INTO left_over FROM user_objects;
    IF left_over IS NOT NULL THEN
      RAISE_APPLICATION_ERROR(-20998, 'omni-parse: the scratch schema still holds ' || left_over);
    END IF;
    IF :2 IS NOT NULL THEN
      EXECUTE IMMEDIATE :2;
    END IF;
  END IF;
  IF code <> 0 THEN
    RAISE_APPLICATION_ERROR(-20999, 'omni-parse code=' || code || ' position=' || pos || ' compile=' || compile);
  END IF;
END;`

// oracleParsedOnly lists the OCI SQL function codes of the statements
// DBMS_SQL.PARSE parses and does not run: SET ROLE (2), INSERT (3), SELECT
// (4), UPDATE (5), DELETE (9), EXPLAIN (27), ALTER SYSTEM (32), SET
// TRANSACTION (33), LOCK TABLE (35), ALTER SESSION (52), COMMIT (54),
// ROLLBACK (55), SAVEPOINT (56), SET CONSTRAINTS (90), CALL (170) and MERGE
// (189), as Oracle Free 23.26.3 reports them. None of them changes the
// schema when parsed.
const oracleParsedOnly = "2, 3, 4, 5, 9, 27, 32, 33, 35, 52, 54, 55, 56, 90, 170, 189"

// OracleProbeTail is what OracleVerdict's probe appends to a statement: a
// character Oracle's lexer refuses (ORA-00911), on a line of its own so a
// trailing -- comment cannot take it in.
const OracleProbeTail = "\n ~"

// oracleTail is where Oracle reports the probe's tail: the ~, two bytes
// after the statement.
const oracleTail = len(OracleProbeTail) - 1

// oracleAfterParse are errors Oracle raises before it has read all of a
// statement, about the database and not the text: a privilege the user
// lacks (01031) and a table or view that is not there (00942). The SQL
// layer checks both on a PL/SQL unit's header before it passes the body to
// the compiler, as on CREATE TRIGGER ... ON DATABASE, so the probe, whose
// tail goes into the body, is answered the same and cannot tell.
var oracleAfterParse = map[int]bool{
	942:  true,
	1031: true,
}

// oracleParserCode reports whether a code is one Oracle's SQL parser
// raises: ORA-00900 to 00999 (invalid SQL statement, missing keyword,
// missing expression, ...), 01700 to 01799 (the lexer's: an unterminated
// string or comment, a zero-length identifier), 02000 to 02099 (missing
// keyword and option errors) and 03048 to 03050 (23ai's errors naming the
// word that is not valid where it stands). It decides an error at the end
// of the statement, where the probe cannot.
func oracleParserCode(code int) bool {
	return code >= 900 && code <= 999 || code >= 1700 && code <= 1799 ||
		code >= 2000 && code <= 2099 || code >= 3048 && code <= 3050
}

// oracleAnswer reads what OracleParse raised.
var oracleAnswer = regexp.MustCompile(`ORA-20999: omni-parse code=(\d+) position=(-?\d+) compile=([\d,]*)`)

type oracleError struct {
	code, position int
	compile        []int
}

func readOracle(err error) (oracleError, bool) {
	m := oracleAnswer.FindStringSubmatch(err.Error())
	if m == nil {
		return oracleError{}, false
	}
	e := oracleError{}
	e.code, _ = strconv.Atoi(m[1])
	e.position, _ = strconv.Atoi(m[2])
	for _, n := range strings.Split(m[3], ",") {
		if n == "" {
			continue
		}
		v, _ := strconv.Atoi(n)
		e.compile = append(e.compile, v)
	}
	return e, true
}

// OracleVerdict reads Oracle's answer to OracleParse of a statement:
// whether Oracle takes its syntax. parse runs OracleParse on a text and
// returns its error. It is called again with OracleProbeTail after the
// statement where the first answer leaves the verdict open.
//
// DBMS_SQL.PARSE resolves names, checks privileges and, for DDL, runs the
// statement, and an error from any of those comes with a position too. The
// probe tells them from a syntax error, as PrepareVerdict's does on MySQL:
// Oracle reads all of a statement before it looks anything up, so a
// statement refused after its parse is answered with the tail's ORA-00911,
// and one its parser refused is answered with the same error at the same
// byte, which comes before the tail. A rejection carries that byte.
//
// Three answers need no probe. ORA-24344 is a PL/SQL unit Oracle created
// with compilation errors, and ORA-06550 an anonymous block the PL/SQL
// compiler refused: either is a rejection when one of the compiler's errors
// is its parser's (PLS-00100 to 00199, PLS-00103 "Encountered the symbol"
// among them). A stored unit's errors are at a line and column of the
// source Oracle keeps, so its rejection carries no byte; an anonymous
// block's carries LAST_ERROR_POSITION. A probe could not tell: its tail is
// in the block. The errors in oracleAfterParse are
// acceptances. And an error at the very end of the statement, where the
// probe's tail is met either way, is a rejection when its code is the
// parser's (oracleParserCode), and an acceptance otherwise: CREATE
// PLUGGABLE DATABASE a FROM b is answered ORA-65016 there, for a clause the
// statement leaves out.
//
// An answer that is no verdict comes back as an error, so the caller fails
// rather than records agreement: an error OracleParse did not raise (a
// dropped connection, a timeout), a probe the server took, and a probe that
// stopped neither where the statement did nor at its tail.
func OracleVerdict(statement string, parse func(text string) error) (string, error) {
	err := parse(statement)
	if err == nil {
		return Accept, nil
	}
	e, ok := readOracle(err)
	if !ok {
		return "", err
	}
	if e.code == 24344 || e.code == 6550 {
		if len(e.compile) == 0 {
			return "", fmt.Errorf("Oracle answered ORA-%05d and lists no PL/SQL error: %w", e.code, err)
		}
		for _, c := range e.compile {
			if c >= 100 && c <= 199 {
				if e.code == 6550 && e.position >= 0 && e.position <= len(statement) {
					return RejectAt(e.position), nil
				}
				return Reject, nil
			}
		}
		return Accept, nil
	}
	if oracleAfterParse[e.code] {
		return Accept, nil
	}
	reject := Reject
	if e.position >= 0 && e.position <= len(statement) {
		reject = RejectAt(e.position)
	}
	if e.position == len(statement) && oracleParserCode(e.code) {
		return reject, nil
	}
	probeErr := parse(statement + OracleProbeTail)
	if probeErr == nil {
		return "", fmt.Errorf("Oracle took the probe's tail after answering %w", err)
	}
	p, ok := readOracle(probeErr)
	if !ok {
		return "", probeErr
	}
	switch {
	case p.code == e.code && p.position == e.position:
		return reject, nil
	case p.code == 911 && p.position == len(statement)+oracleTail:
		return Accept, nil
	}
	return "", fmt.Errorf("Oracle answered ORA-%05d at byte %d, and the probe ORA-%05d at byte %d, neither where the statement stopped nor at the tail",
		e.code, e.position, p.code, p.position)
}
