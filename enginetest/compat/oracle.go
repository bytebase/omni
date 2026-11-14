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
// PL/SQL compiler's errors, each PLS-nnnnn, or ORA-nnnnn for an error of a
// SQL statement inside the PL/SQL ("PL/SQL: ORA-00942", whose code
// USER_ERRORS keeps only in its text): for ORA-24344, a unit created with
// compilation errors, those USER_ERRORS lists (PLS-00000 is "Statement
// ignored", which follows another), and for ORA-06550, an
// anonymous block that did not compile, those of its error stack.
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
    SELECT LISTAGG(CASE WHEN text LIKE 'PL/SQL: ORA-%' THEN REGEXP_SUBSTR(text, '^PL/SQL: (ORA-[0-9]{5})', 1, 1, NULL, 1)
                        ELSE 'PLS-' || TO_CHAR(message_number, 'FM00000') END, ',')
             WITHIN GROUP (ORDER BY name, type, sequence)
      INTO compile FROM user_errors WHERE attribute = 'ERROR';
  ELSIF code = 6550 THEN
    FOR i IN 1 .. REGEXP_COUNT(stack, '^(PLS|PL/SQL: ORA)-[0-9]{5}', 1, 'm') LOOP
      compile := compile || CASE WHEN i > 1 THEN ',' END ||
        REPLACE(REGEXP_SUBSTR(stack, '^(PLS|PL/SQL: ORA)-[0-9]{5}', 1, i, 'm'), 'PL/SQL: ', '');
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

// oracleTail is where Oracle reports an error at the probe's tail: the ~,
// two bytes after the statement. It is ORA-00911, or the error of what the
// parser wanted there, as ORA-14006 for a partition's name.
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
// missing expression, ...), the lexer's and the parser's of 01700 to 01799
// (oracleLexer), 02000 to 02099 (missing keyword and option errors) and
// 03048 to 03050 (23ai's errors naming the word that is not valid where it
// stands), less the codes of those ranges raised after the parse, in
// oracleResolved. It decides an error at the end of the statement, where
// the probe cannot, one whose probe says nothing, and a SQL statement's
// error inside PL/SQL.
func oracleParserCode(code int) bool {
	return (code >= 900 && code <= 999 || oracleLexer[code] ||
		code >= 2000 && code <= 2099 || code >= 3048 && code <= 3050) && !oracleResolved[code]
}

// oracleLexer are the codes of 01700 to 01799 the lexer and the parser
// raise: an invalid ALTER TABLE option (01735), a missing double quote
// (01740), a zero-length identifier (01741), a comment not terminated
// (01742), an invalid bind variable name (01745), an invalid column
// specification (01747), a quoted string not terminated (01756). The rest
// of the range is checked after the parse, as 01790, a set operation's
// columns of different types.
var oracleLexer = map[int]bool{1735: true, 1740: true, 1741: true, 1742: true, 1745: true, 1747: true, 1756: true}

// oracleResolved are codes among the parser's that Oracle raises after the
// parse, resolving names and types: an invalid identifier (00904, a column
// that is not there), a column ambiguously defined (00918), inconsistent
// datatypes (00932), group functions where none may stand (00934, 00937,
// 00979), a missing table or view (00942), too few values (00947), a name
// already used (00955), a column named twice (00957).
var oracleResolved = map[int]bool{904: true, 918: true, 932: true, 934: true, 937: true, 942: true, 947: true, 955: true, 957: true, 979: true}

// oracleAnswer reads what OracleParse raised.
var oracleAnswer = regexp.MustCompile(`ORA-20999: omni-parse code=(\d+) position=(-?\d+) compile=([A-Z0-9,-]*)`)

type oracleError struct {
	code, position int
	compile        []string
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
		if n != "" {
			e.compile = append(e.compile, n)
		}
	}
	return e, true
}

// oracleFailure matches Oracle failing in its own code while it parses: an
// internal error (ORA-00600), one that leaves the session failing the next
// query it runs with ORA-00000, "normal, successful completion", as ALTER
// MATERIALIZED VIEW LOG ON a SPLIT PARTITION a AT(1) INCLUDING ROWS WHERE
// 1 = 1 does, a process that died (ORA-07445), or the session the server
// dropped with it (end-of-file on the connection).
var oracleFailure = regexp.MustCompile(`ORA-00600|ORA-00000|ORA-07445|ORA-03113|ORA-03114|^EOF$`)

// OracleVerdict reads Oracle's answer to OracleParse of a statement:
// whether Oracle takes its syntax. parse runs OracleParse on a text and
// returns its error. It is called again with OracleProbeTail after the
// statement where the first answer leaves the verdict open.
//
// DBMS_SQL.PARSE resolves names, checks privileges and, for DDL, runs the
// statement, and an error from any of those comes with a position too. The
// probe tells them from a syntax error, as PrepareVerdict's does on MySQL:
// Oracle reads all of a statement before it looks anything up, so a
// statement refused after its parse is answered at the probe's tail, and
// one its parser refused is answered with the same error at the same byte,
// which comes before the tail. A rejection carries that byte. A probe that
// stops anywhere else has changed how the parser reads the text before the
// tail, as ALTER INDEX a UNUSABLE followed by more is refused at UNUSABLE,
// and says nothing: the first answer's code decides, a rejection when it
// is the parser's (oracleParserCode) and an acceptance otherwise.
//
// Some answers need no probe. ORA-24344 is a PL/SQL unit Oracle created
// with compilation errors, and ORA-06550 an anonymous block the PL/SQL
// compiler refused: either is a rejection when one of the compiler's errors
// is a syntax error (oracleCompileSyntax): the PL/SQL parser's, PLS-00103
// "Encountered the symbol" among them, or the SQL parser's on a statement
// inside the PL/SQL. A stored unit's errors are at a line and column of the
// source Oracle keeps, so its rejection carries no byte; an anonymous
// block's carries LAST_ERROR_POSITION. A probe could not tell: its tail is
// in the block. The errors in oracleAfterParse are acceptances. An error at
// the very end of the statement, where the probe's tail is met either way,
// is decided by its code. And Oracle failing in its own code on the
// statement (oracleFailure), on two sessions, is a rejection with no byte:
// no statement Oracle runs, though not one its parser names an error in.
//
// An answer that is no verdict comes back as an error, so the caller fails
// rather than records agreement: an error OracleParse did not raise and
// that is no failure of Oracle's (a timeout), and a probe the server took.
func OracleVerdict(statement string, parse func(text string) error) (string, error) {
	err := parse(statement)
	if err != nil && oracleFailure.MatchString(firstLine(err)) {
		// Again, on the new session database/sql opens for a dropped one.
		if err = parse(statement); err != nil && oracleFailure.MatchString(firstLine(err)) {
			return Reject, nil
		}
	}
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
			if oracleCompileSyntax(c) {
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
	byCode := func() string {
		if oracleParserCode(e.code) {
			return reject
		}
		return Accept
	}
	if e.position == len(statement) {
		return byCode(), nil
	}
	probeErr := parse(statement + OracleProbeTail)
	if probeErr == nil {
		return "", fmt.Errorf("Oracle took the probe's tail after answering %w", err)
	}
	p, ok := readOracle(probeErr)
	switch {
	case !ok && oracleFailure.MatchString(firstLine(probeErr)):
		return byCode(), nil
	case !ok:
		return "", probeErr
	case p.code == e.code && p.position == e.position:
		return reject, nil
	case p.position == len(statement)+oracleTail:
		return Accept, nil
	}
	return byCode(), nil
}

// oracleCompileSyntax reports whether a PL/SQL compilation error is a
// syntax error: the PL/SQL parser's (PLS-00100 to 00199), or the SQL
// parser's on a statement inside the PL/SQL (an ORA-nnnnn of
// oracleParserCode).
func oracleCompileSyntax(e string) bool {
	kind, num, ok := strings.Cut(e, "-")
	n, err := strconv.Atoi(num)
	if !ok || err != nil {
		return false
	}
	switch kind {
	case "PLS":
		return n >= 100 && n <= 199
	case "ORA":
		return oracleParserCode(n)
	}
	return false
}

// firstLine is an error's text to its first line break: go-ora puts the
// error stack and the position after it.
func firstLine(err error) string {
	line, _, _ := strings.Cut(err.Error(), "\n")
	return strings.TrimSpace(line)
}
