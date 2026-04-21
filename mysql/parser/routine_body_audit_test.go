package parser

import (
	"strings"
	"testing"
)

// TestRoutineBodyAudit systematically probes parser coverage for stored-
// program bodies after grammar-driven parsing lands. Each probe has an
// expected outcome. A probe "fails the audit" when the actual result does
// not match expected.
//
// Categories:
//   S — Statement dispatch gap (parser should accept an sp_proc_stmt-valid
//       statement); expected=accept; audit fails if reject.
//   E — Expression production gap; expected=accept; audit fails if reject.
//   C — Context-sensitive token handling; expected=accept; audit fails if reject.
//   V — Static validation gap (MySQL CREATE-time rejects what the parser
//       currently accepts); expected=reject; audit fails if accept.
//       These are the items deferred to a future mysql/semantic PR, and the
//       audit locks in "omni still accepts this" to track the backlog.
//
// Test does not fail; we log a structured report so the backlog is visible
// to reviewers and drives the follow-up work.
func TestRoutineBodyAudit(t *testing.T) {
	type probe struct {
		cat    string // S / E / C / V
		expect string // "accept" or "reject"
		name   string
		sql    string
	}
	accept := func(cat, name, sql string) probe { return probe{cat, "accept", name, sql} }
	reject := func(cat, name, sql string) probe { return probe{cat, "reject", name, sql} }

	probes := []probe{
		// --- DECLARE variable types (S) ---
		accept("S", "DECLARE INT", `CREATE PROCEDURE p() BEGIN DECLARE x INT; END`),
		accept("S", "DECLARE INT UNSIGNED", `CREATE PROCEDURE p() BEGIN DECLARE x INT UNSIGNED; END`),
		accept("S", "DECLARE BIGINT UNSIGNED ZEROFILL", `CREATE PROCEDURE p() BEGIN DECLARE x BIGINT UNSIGNED ZEROFILL; END`),
		accept("S", "DECLARE TINYINT(1)", `CREATE PROCEDURE p() BEGIN DECLARE x TINYINT(1); END`),
		accept("S", "DECLARE DECIMAL(10,2)", `CREATE PROCEDURE p() BEGIN DECLARE x DECIMAL(10,2); END`),
		accept("S", "DECLARE NUMERIC(5,2) UNSIGNED", `CREATE PROCEDURE p() BEGIN DECLARE x NUMERIC(5,2) UNSIGNED; END`),
		accept("S", "DECLARE FLOAT(10,2)", `CREATE PROCEDURE p() BEGIN DECLARE x FLOAT(10,2); END`),
		accept("S", "DECLARE DOUBLE", `CREATE PROCEDURE p() BEGIN DECLARE x DOUBLE; END`),
		accept("S", "DECLARE BIT(8)", `CREATE PROCEDURE p() BEGIN DECLARE x BIT(8); END`),
		accept("S", "DECLARE BOOL DEFAULT TRUE", `CREATE PROCEDURE p() BEGIN DECLARE x BOOL DEFAULT TRUE; END`),
		accept("S", "DECLARE DATE", `CREATE PROCEDURE p() BEGIN DECLARE x DATE; END`),
		accept("S", "DECLARE DATETIME(6)", `CREATE PROCEDURE p() BEGIN DECLARE x DATETIME(6); END`),
		accept("S", "DECLARE TIMESTAMP DEFAULT CURRENT_TIMESTAMP", `CREATE PROCEDURE p() BEGIN DECLARE x TIMESTAMP DEFAULT CURRENT_TIMESTAMP; END`),
		accept("S", "DECLARE TIME(3)", `CREATE PROCEDURE p() BEGIN DECLARE x TIME(3); END`),
		accept("S", "DECLARE YEAR", `CREATE PROCEDURE p() BEGIN DECLARE x YEAR; END`),
		accept("S", "DECLARE CHAR(10)", `CREATE PROCEDURE p() BEGIN DECLARE x CHAR(10); END`),
		accept("S", "DECLARE VARCHAR(255) CHARACTER SET utf8mb4", `CREATE PROCEDURE p() BEGIN DECLARE x VARCHAR(255) CHARACTER SET utf8mb4; END`),
		accept("S", "DECLARE VARCHAR(255) COLLATE utf8mb4_bin", `CREATE PROCEDURE p() BEGIN DECLARE x VARCHAR(255) COLLATE utf8mb4_bin; END`),
		accept("S", "DECLARE TEXT", `CREATE PROCEDURE p() BEGIN DECLARE x TEXT; END`),
		accept("S", "DECLARE LONGTEXT", `CREATE PROCEDURE p() BEGIN DECLARE x LONGTEXT; END`),
		accept("S", "DECLARE BLOB", `CREATE PROCEDURE p() BEGIN DECLARE x BLOB; END`),
		accept("S", "DECLARE JSON", `CREATE PROCEDURE p() BEGIN DECLARE x JSON; END`),
		accept("S", "DECLARE ENUM", `CREATE PROCEDURE p() BEGIN DECLARE x ENUM('a','b','c'); END`),
		accept("S", "DECLARE SET", `CREATE PROCEDURE p() BEGIN DECLARE x SET('a','b','c'); END`),
		accept("S", "DECLARE multi comma", `CREATE PROCEDURE p() BEGIN DECLARE a, b, c INT DEFAULT 0; END`),

		// --- SET assignment forms (S) ---
		accept("S", "SET @user_var", `CREATE PROCEDURE p() BEGIN SET @x = 1; END`),
		accept("S", "SET local_var", `CREATE PROCEDURE p() BEGIN DECLARE x INT; SET x = 1; END`),
		accept("S", "SET NEW.col in trigger", `CREATE TRIGGER t BEFORE INSERT ON tbl FOR EACH ROW SET NEW.c = 0`),
		accept("S", "SET OLD.col in trigger", `CREATE TRIGGER t BEFORE UPDATE ON tbl FOR EACH ROW SET OLD.c = 0`),
		accept("S", "SET := (Pascal assignment)", `CREATE PROCEDURE p() BEGIN DECLARE x INT; SET x := 1; END`),
		accept("S", "SET multi comma", `CREATE PROCEDURE p() BEGIN SET @a = 1, @b = 2, @c = 3; END`),
		accept("S", "SET @@session.var", `CREATE PROCEDURE p() BEGIN SET @@session.sql_mode = ''; END`),
		accept("S", "SET @@GLOBAL.var", `CREATE PROCEDURE p() BEGIN SET @@GLOBAL.max_connections = 100; END`),
		accept("S", "SET SESSION scope", `CREATE PROCEDURE p() BEGIN SET SESSION sql_mode = ''; END`),
		accept("S", "SET PERSIST", `CREATE PROCEDURE p() BEGIN SET PERSIST max_connections = 200; END`),
		accept("S", "SET NAMES", `CREATE PROCEDURE p() BEGIN SET NAMES utf8mb4; END`),
		accept("S", "SET CHARACTER SET", `CREATE PROCEDURE p() BEGIN SET CHARACTER SET utf8mb4; END`),
		accept("S", "SET TRANSACTION ISOLATION LEVEL", `CREATE PROCEDURE p() BEGIN SET TRANSACTION ISOLATION LEVEL READ COMMITTED; END`),

		// --- SELECT variants in body (S) ---
		accept("S", "SELECT INTO @user_var", `CREATE PROCEDURE p() BEGIN SELECT 1 INTO @x; END`),
		accept("S", "SELECT INTO bare var", `CREATE PROCEDURE p() BEGIN DECLARE x INT; SELECT 1 INTO x; END`),
		accept("S", "SELECT INTO multiple", `CREATE PROCEDURE p() BEGIN DECLARE a,b INT; SELECT 1, 2 INTO a, b; END`),
		accept("S", "SELECT INTO mixed @ and bare", `CREATE PROCEDURE p() BEGIN DECLARE b INT; SELECT 1, 2 INTO @a, b; END`),
		accept("S", "SELECT FROM with WHERE", `CREATE PROCEDURE p() BEGIN SELECT id FROM t WHERE id = 1; END`),
		accept("S", "SELECT with JOIN", `CREATE PROCEDURE p() BEGIN SELECT a.id FROM t a JOIN u b ON a.id=b.id; END`),
		accept("S", "SELECT with UNION", `CREATE PROCEDURE p() BEGIN SELECT 1 UNION SELECT 2; END`),
		accept("S", "WITH CTE + SELECT", `CREATE PROCEDURE p() BEGIN WITH c AS (SELECT 1 AS v) SELECT v FROM c; END`),
		accept("S", "SELECT FOR UPDATE", `CREATE PROCEDURE p() BEGIN SELECT id FROM t WHERE id = 1 FOR UPDATE; END`),
		accept("S", "SELECT LOCK IN SHARE MODE", `CREATE PROCEDURE p() BEGIN SELECT id FROM t WHERE id = 1 LOCK IN SHARE MODE; END`),

		// --- Cursor operations (S) ---
		accept("S", "DECLARE cursor + OPEN + FETCH + CLOSE", `CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; OPEN c; FETCH c INTO @x; CLOSE c; END`),
		accept("S", "FETCH NEXT FROM cursor", `CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; OPEN c; FETCH NEXT FROM c INTO @x; CLOSE c; END`),
		accept("S", "FETCH FROM cursor (no NEXT)", `CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; OPEN c; FETCH FROM c INTO @x; CLOSE c; END`),
		accept("S", "FETCH into multiple vars", `CREATE PROCEDURE p() BEGIN DECLARE a,b INT; DECLARE c CURSOR FOR SELECT 1,2; OPEN c; FETCH c INTO a, b; CLOSE c; END`),

		// --- Handler condition_value variants (S) ---
		accept("S", "HANDLER FOR SQLSTATE literal", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000' SET @e=1; SELECT 1; END`),
		accept("S", "HANDLER FOR SQLSTATE VALUE literal", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE VALUE '23000' SET @e=1; SELECT 1; END`),
		accept("S", "HANDLER FOR mysql error code int", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1062 SET @e=1; SELECT 1; END`),
		accept("S", "HANDLER FOR named condition", `CREATE PROCEDURE p() BEGIN DECLARE dup CONDITION FOR SQLSTATE '23000'; DECLARE EXIT HANDLER FOR dup SET @e=1; SELECT 1; END`),
		accept("S", "HANDLER FOR SQLWARNING", `CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR SQLWARNING SET @w=1; SELECT 1; END`),
		accept("S", "HANDLER FOR NOT FOUND", `CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR NOT FOUND SET @d=1; SELECT 1; END`),
		accept("S", "HANDLER FOR SQLEXCEPTION", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION SET @e=1; SELECT 1; END`),
		accept("S", "HANDLER FOR multiple conditions", `CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR NOT FOUND, SQLWARNING SET @x=1; SELECT 1; END`),
		accept("S", "CONTINUE handler with compound body", `CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR SQLEXCEPTION BEGIN SET @e=1; ROLLBACK; END; SELECT 1; END`),
		accept("S", "UNDO handler (rare)", `CREATE PROCEDURE p() BEGIN DECLARE UNDO HANDLER FOR SQLEXCEPTION SET @e=1; SELECT 1; END`),

		// --- SIGNAL / RESIGNAL / GET DIAGNOSTICS (S) ---
		accept("S", "SIGNAL SQLSTATE", `CREATE PROCEDURE p() BEGIN SIGNAL SQLSTATE '45000'; END`),
		accept("S", "SIGNAL with SET MESSAGE_TEXT", `CREATE PROCEDURE p() BEGIN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'fail'; END`),
		accept("S", "SIGNAL with multiple items", `CREATE PROCEDURE p() BEGIN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'fail', MYSQL_ERRNO = 9999; END`),
		accept("S", "RESIGNAL bare", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION RESIGNAL; SELECT 1; END`),
		accept("S", "RESIGNAL with new info", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION RESIGNAL SET MESSAGE_TEXT = 'wrapped'; SELECT 1; END`),
		accept("S", "GET DIAGNOSTICS CONDITION", `CREATE PROCEDURE p() BEGIN DECLARE x INT; GET DIAGNOSTICS CONDITION 1 x = MYSQL_ERRNO; END`),
		accept("S", "GET DIAGNOSTICS (no CONDITION)", `CREATE PROCEDURE p() BEGIN DECLARE n INT; GET DIAGNOSTICS n = NUMBER; END`),
		accept("S", "GET CURRENT DIAGNOSTICS", `CREATE PROCEDURE p() BEGIN DECLARE n INT; GET CURRENT DIAGNOSTICS n = NUMBER; END`),
		accept("S", "GET STACKED DIAGNOSTICS", `CREATE PROCEDURE p() BEGIN DECLARE n INT; GET STACKED DIAGNOSTICS n = NUMBER; END`),

		// --- Dynamic SQL (S) ---
		accept("S", "PREPARE from literal", `CREATE PROCEDURE p() BEGIN PREPARE s FROM 'SELECT 1'; EXECUTE s; DEALLOCATE PREPARE s; END`),
		accept("S", "PREPARE from @var", `CREATE PROCEDURE p() BEGIN SET @q = 'SELECT 1'; PREPARE s FROM @q; EXECUTE s; DROP PREPARE s; END`),
		accept("S", "EXECUTE USING", `CREATE PROCEDURE p() BEGIN PREPARE s FROM 'SELECT ?'; SET @a = 1; EXECUTE s USING @a; DEALLOCATE PREPARE s; END`),
		accept("S", "EXECUTE USING multiple", `CREATE PROCEDURE p() BEGIN PREPARE s FROM 'SELECT ?, ?'; SET @a=1; SET @b=2; EXECUTE s USING @a, @b; DEALLOCATE PREPARE s; END`),

		// --- DDL inside body (S) ---
		accept("S", "CREATE TABLE", `CREATE PROCEDURE p() BEGIN CREATE TABLE t (id INT); END`),
		accept("S", "CREATE TABLE IF NOT EXISTS", `CREATE PROCEDURE p() BEGIN CREATE TABLE IF NOT EXISTS t (id INT); END`),
		accept("S", "CREATE TEMPORARY TABLE", `CREATE PROCEDURE p() BEGIN CREATE TEMPORARY TABLE tmp (id INT); END`),
		accept("S", "CREATE TABLE AS SELECT", `CREATE PROCEDURE p() BEGIN CREATE TABLE t AS SELECT 1 AS id; END`),
		accept("S", "DROP TABLE", `CREATE PROCEDURE p() BEGIN DROP TABLE t; END`),
		accept("S", "DROP TABLE IF EXISTS", `CREATE PROCEDURE p() BEGIN DROP TABLE IF EXISTS t; END`),
		accept("S", "TRUNCATE TABLE", `CREATE PROCEDURE p() BEGIN TRUNCATE TABLE t; END`),
		accept("S", "ALTER TABLE", `CREATE PROCEDURE p() BEGIN ALTER TABLE t ADD COLUMN c INT; END`),
		accept("S", "CREATE INDEX", `CREATE PROCEDURE p() BEGIN CREATE INDEX idx ON t (c); END`),
		accept("S", "DROP INDEX", `CREATE PROCEDURE p() BEGIN DROP INDEX idx ON t; END`),

		// --- Transaction (S) ---
		accept("S", "START TRANSACTION", `CREATE PROCEDURE p() BEGIN START TRANSACTION; SELECT 1; COMMIT; END`),
		accept("S", "START TRANSACTION READ ONLY", `CREATE PROCEDURE p() BEGIN START TRANSACTION READ ONLY; COMMIT; END`),
		accept("S", "SAVEPOINT", `CREATE PROCEDURE p() BEGIN SAVEPOINT sp1; ROLLBACK TO SAVEPOINT sp1; RELEASE SAVEPOINT sp1; END`),
		accept("S", "LOCK TABLES", `CREATE PROCEDURE p() BEGIN LOCK TABLES t READ; UNLOCK TABLES; END`),

		// --- Expression production (E) — INTERVAL arithmetic density ---
		accept("E", "INTERVAL literal+unit", `CREATE PROCEDURE p() BEGIN DECLARE d DATE; SET d = CURRENT_DATE + INTERVAL 1 DAY; END`),
		accept("E", "INTERVAL parenthesized expr+unit", `CREATE PROCEDURE p() BEGIN DECLARE d DATE; SET d = CURRENT_DATE + INTERVAL (1+2) DAY; END`),
		accept("E", "INTERVAL DAY(x) DAY (unit happens to be function name)", `CREATE PROCEDURE p() BEGIN DECLARE d DATE; SET d = CURRENT_DATE - INTERVAL DAY(NOW()) DAY; END`),
		accept("E", "INTERVAL DAY(x) - 1 DAY (arithmetic between fn and unit, sakila form)", `CREATE PROCEDURE p() BEGIN DECLARE d DATE; SET d = d - INTERVAL DAY(d) - 1 DAY; END`),
		accept("E", "INTERVAL @var+1 DAY", `CREATE PROCEDURE p() BEGIN SET @d = NOW() + INTERVAL @n + 1 DAY; END`),
		accept("E", "INTERVAL x YEAR_MONTH compound unit", `CREATE PROCEDURE p() BEGIN DECLARE d DATE; SET d = CURRENT_DATE + INTERVAL '1-2' YEAR_MONTH; END`),
		accept("E", "INTERVAL day_hour compound unit", `CREATE PROCEDURE p() BEGIN DECLARE d DATETIME; SET d = NOW() + INTERVAL '1 12' DAY_HOUR; END`),
		accept("E", "DATE_ADD with INTERVAL", `CREATE PROCEDURE p() BEGIN SET @x = DATE_ADD(NOW(), INTERVAL 1 DAY); END`),
		accept("E", "TIMESTAMP + INTERVAL", `CREATE PROCEDURE p() BEGIN SET @x = TIMESTAMP '2024-01-01 00:00:00' + INTERVAL 1 HOUR; END`),

		// --- Expression production (E) — CASE / IF / COALESCE ---
		accept("E", "CASE expr in SET", `CREATE PROCEDURE p() BEGIN SET @x = CASE 1 WHEN 1 THEN 'a' ELSE 'b' END; END`),
		accept("E", "searched CASE expr", `CREATE PROCEDURE p() BEGIN SET @x = CASE WHEN 1>0 THEN 'a' ELSE 'b' END; END`),
		accept("E", "nested CASE", `CREATE PROCEDURE p() BEGIN SET @x = CASE WHEN CASE 1 WHEN 1 THEN 1 END = 1 THEN 'y' END; END`),
		accept("E", "IF() three-arg function", `CREATE PROCEDURE p() BEGIN SET @x = IF(1>0, 'a', 'b'); END`),
		accept("E", "IFNULL", `CREATE PROCEDURE p() BEGIN SET @x = IFNULL(NULL, 'a'); END`),
		accept("E", "NULLIF", `CREATE PROCEDURE p() BEGIN SET @x = NULLIF(1, 2); END`),
		accept("E", "COALESCE", `CREATE PROCEDURE p() BEGIN SET @x = COALESCE(NULL, NULL, 'a'); END`),

		// --- Expression production (E) — JSON / operators ---
		accept("E", "JSON_EXTRACT function", `CREATE PROCEDURE p() BEGIN SET @x = JSON_EXTRACT('{"a":1}', '$.a'); END`),
		accept("E", "-> column path operator", `CREATE PROCEDURE p() BEGIN SET @j = '{"a":1}'; SET @x = @j->'$.a'; END`),
		accept("E", "->> column path unquote", `CREATE PROCEDURE p() BEGIN SET @j = '{"a":1}'; SET @x = @j->>'$.a'; END`),
		accept("E", "JSON_TABLE", `CREATE PROCEDURE p() BEGIN SELECT * FROM JSON_TABLE('[1,2]', '$[*]' COLUMNS (v INT PATH '$')) AS t; END`),

		// --- Expression production (E) — window / aggregates ---
		accept("E", "Window function", `CREATE PROCEDURE p() BEGIN SELECT ROW_NUMBER() OVER (ORDER BY id) FROM t; END`),
		accept("E", "Window with PARTITION BY", `CREATE PROCEDURE p() BEGIN SELECT SUM(v) OVER (PARTITION BY k ORDER BY id) FROM t; END`),
		accept("E", "GROUP_CONCAT with separator", `CREATE PROCEDURE p() BEGIN SELECT GROUP_CONCAT(x SEPARATOR ',') FROM t; END`),
		accept("E", "GROUP_CONCAT DISTINCT ORDER BY", `CREATE PROCEDURE p() BEGIN SELECT GROUP_CONCAT(DISTINCT x ORDER BY x DESC SEPARATOR '|') FROM t; END`),
		accept("E", "COUNT(DISTINCT)", `CREATE PROCEDURE p() BEGIN SELECT COUNT(DISTINCT x) FROM t; END`),

		// --- Expression production (E) — CAST / type ops ---
		accept("E", "CAST AS SIGNED", `CREATE PROCEDURE p() BEGIN SET @x = CAST('1' AS SIGNED); END`),
		accept("E", "CAST AS DECIMAL(10,2)", `CREATE PROCEDURE p() BEGIN SET @x = CAST('1.5' AS DECIMAL(10,2)); END`),
		accept("E", "CAST AS JSON", `CREATE PROCEDURE p() BEGIN SET @x = CAST('[1]' AS JSON); END`),
		accept("E", "CONVERT USING charset", `CREATE PROCEDURE p() BEGIN SET @x = CONVERT('abc' USING utf8mb4); END`),
		accept("E", "BINARY literal prefix", `CREATE PROCEDURE p() BEGIN SET @x = BINARY 'abc'; END`),

		// --- Expression production (E) — predicates / subqueries ---
		accept("E", "EXTRACT YEAR FROM date", `CREATE PROCEDURE p() BEGIN SET @x = EXTRACT(YEAR FROM NOW()); END`),
		accept("E", "STR_TO_DATE", `CREATE PROCEDURE p() BEGIN SET @x = STR_TO_DATE('2024-01-01', '%Y-%m-%d'); END`),
		accept("E", "TIMESTAMPDIFF", `CREATE PROCEDURE p() BEGIN SET @x = TIMESTAMPDIFF(DAY, '2024-01-01', NOW()); END`),
		accept("E", "scalar subquery in expr", `CREATE PROCEDURE p() BEGIN SET @x = (SELECT MAX(id) FROM t); END`),
		accept("E", "EXISTS subquery in IF cond", `CREATE PROCEDURE p() BEGIN IF EXISTS (SELECT 1 FROM t) THEN SELECT 1; END IF; END`),
		accept("E", "IN subquery", `CREATE PROCEDURE p() BEGIN IF 1 IN (SELECT id FROM t) THEN SELECT 1; END IF; END`),
		accept("E", "NOT IN subquery", `CREATE PROCEDURE p() BEGIN IF 1 NOT IN (SELECT id FROM t) THEN SELECT 1; END IF; END`),
		accept("E", "BETWEEN", `CREATE PROCEDURE p() BEGIN IF 5 BETWEEN 1 AND 10 THEN SELECT 1; END IF; END`),
		accept("E", "LIKE with ESCAPE", `CREATE PROCEDURE p() BEGIN IF 'abc' LIKE 'a%' ESCAPE '\\' THEN SELECT 1; END IF; END`),
		accept("E", "REGEXP", `CREATE PROCEDURE p() BEGIN IF 'abc' REGEXP '^a' THEN SELECT 1; END IF; END`),
		accept("E", "row constructor EXISTS", `CREATE PROCEDURE p() BEGIN IF ROW(1,2) = ROW(1,2) THEN SELECT 1; END IF; END`),
		accept("E", "ANY subquery", `CREATE PROCEDURE p() BEGIN IF 1 = ANY (SELECT id FROM t) THEN SELECT 1; END IF; END`),
		accept("E", "ALL subquery", `CREATE PROCEDURE p() BEGIN IF 1 < ALL (SELECT id FROM t) THEN SELECT 1; END IF; END`),

		// --- Context-sensitive (C) ---
		accept("C", "NEW.col in trigger expr", `CREATE TRIGGER t BEFORE INSERT ON tbl FOR EACH ROW BEGIN IF NEW.a > 0 THEN SET NEW.b = 1; END IF; END`),
		accept("C", "OLD.col in trigger expr", `CREATE TRIGGER t BEFORE UPDATE ON tbl FOR EACH ROW BEGIN IF OLD.a <> NEW.a THEN SET @c = 1; END IF; END`),
		accept("C", "NEW.col in INSERT VALUES", `CREATE TRIGGER t BEFORE INSERT ON tbl FOR EACH ROW INSERT INTO audit VALUES (NEW.id, NOW())`),

		// --- Static validation (V) — MySQL rejects, omni currently accepts ---
		reject("V", "LEAVE undeclared label", `CREATE PROCEDURE p() BEGIN LEAVE nonexistent; END`),
		reject("V", "ITERATE undeclared label", `CREATE PROCEDURE p() BEGIN ITERATE nonexistent; END`),
		reject("V", "OPEN undeclared cursor", `CREATE PROCEDURE p() BEGIN OPEN nonexistent; END`),
		reject("V", "FETCH undeclared cursor", `CREATE PROCEDURE p() BEGIN FETCH nonexistent INTO @x; END`),
		reject("V", "CLOSE undeclared cursor", `CREATE PROCEDURE p() BEGIN CLOSE nonexistent; END`),
		reject("V", "HANDLER references undeclared condition", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR nonexistent SET @e=1; SELECT 1; END`),
		reject("V", "assign to undeclared variable", `CREATE PROCEDURE p() BEGIN SET nonexistent_local = 1; END`),
		reject("V", "duplicate DECLARE VAR same scope", `CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE x INT; SELECT 1; END`),
		reject("V", "duplicate DECLARE CURSOR same scope", `CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE c CURSOR FOR SELECT 2; SELECT 1; END`),
		reject("V", "duplicate CONDITION same scope", `CREATE PROCEDURE p() BEGIN DECLARE dk CONDITION FOR SQLSTATE '23000'; DECLARE dk CONDITION FOR SQLSTATE '23001'; SELECT 1; END`),
		reject("V", "duplicate label on same compound", `CREATE PROCEDURE p() lbl: BEGIN lbl: BEGIN SELECT 1; END lbl; END lbl`),
		reject("V", "duplicate HANDLER condition within one DECLARE", `CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23000' SET @e=1; SELECT 1; END`),
		reject("V", "function missing RETURN on else path", `CREATE FUNCTION f() RETURNS INT BEGIN IF @x=1 THEN RETURN 1; END IF; END`),
		reject("V", "function missing RETURN overall", `CREATE FUNCTION f() RETURNS INT BEGIN SELECT 1; END`),

		// --- Real-world shapes (S / mixed) ---
		accept("S", "procedure with handler wrapped body", `CREATE PROCEDURE p() BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN ROLLBACK; SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'failed'; END;
    START TRANSACTION;
    UPDATE t SET v = v + 1;
    COMMIT;
END`),
		accept("S", "procedure with dynamic SQL + variable interpolation (totem-dev style)", `CREATE PROCEDURE p(IN tbl VARCHAR(64), IN quiz INT) BEGIN
    SET @s = CONCAT('SELECT * FROM ', tbl);
    IF quiz THEN
        SET @s = CONCAT(@s, ' WHERE type=?');
    END IF;
    PREPARE st FROM @s;
    EXECUTE st;
    DEALLOCATE PREPARE st;
END`),
		accept("S", "function with cursor + accumulation", `CREATE FUNCTION f() RETURNS INT
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE sum_v INT DEFAULT 0;
    DECLARE v INT;
    DECLARE c CURSOR FOR SELECT id FROM t;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    OPEN c;
    loop1: LOOP
        FETCH c INTO v;
        IF done THEN LEAVE loop1; END IF;
        SET sum_v = sum_v + v;
    END LOOP loop1;
    CLOSE c;
    RETURN sum_v;
END`),
	}

	type outcome struct {
		cat, name, err string
		accepted       bool
		expected       string
	}
	var bad []outcome // probes whose actual differs from expected

	catOK := map[string]int{}
	catBad := map[string]int{}
	for _, pr := range probes {
		_, err := Parse(pr.sql)
		accepted := err == nil
		wantAccept := pr.expect == "accept"
		if accepted == wantAccept {
			catOK[pr.cat]++
			continue
		}
		catBad[pr.cat]++
		o := outcome{cat: pr.cat, name: pr.name, accepted: accepted, expected: pr.expect}
		if err != nil {
			o.err = trimErr(err.Error())
		}
		bad = append(bad, o)
	}

	var sb strings.Builder
	sb.WriteString("\n=== routine-body parser audit ===\n")
	for _, cat := range []string{"S", "E", "C", "V"} {
		sb.WriteString(catBlurb(cat))
		sb.WriteString(": ")
		sb.WriteString(itoaAudit(catOK[cat]))
		sb.WriteString(" ok, ")
		sb.WriteString(itoaAudit(catBad[cat]))
		sb.WriteString(" gap\n")
	}
	if len(bad) > 0 {
		sb.WriteString("\n--- gaps ---\n")
		for _, cat := range []string{"S", "E", "C", "V"} {
			for _, o := range bad {
				if o.cat != cat {
					continue
				}
				sb.WriteString("  [")
				sb.WriteString(o.cat)
				sb.WriteString("] ")
				sb.WriteString(o.name)
				if o.expected == "accept" {
					sb.WriteString(" — parser rejected (should accept):\n     ")
					sb.WriteString(o.err)
				} else {
					sb.WriteString(" — parser accepted (MySQL rejects)")
				}
				sb.WriteString("\n")
			}
		}
	}
	t.Log(sb.String())
}

func trimErr(s string) string {
	before, _, _ := strings.Cut(s, "\nrelated")
	return before
}

func catBlurb(cat string) string {
	switch cat {
	case "S":
		return "S (statement dispatch)"
	case "E":
		return "E (expression production)"
	case "C":
		return "C (context-sensitive)"
	case "V":
		return "V (static validation)"
	}
	return cat
}

func itoaAudit(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
