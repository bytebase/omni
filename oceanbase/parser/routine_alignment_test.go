package parser

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// routineParserCodes are the errors other than 1064 that OceanBase raises while
// it parses a routine body and that omni's parser raises too. The other routine
// errors (1308 a label with no match, 1313 RETURN outside a function, 1327 an
// undeclared variable, ...) come from OceanBase's resolver and are
// oceanbase/validate's to report; its oracle holds Parse plus Validate to the
// same server, so Parse accepts those bodies.
var routineParserCodes = map[uint16]bool{
	1337: true, // ER_SP_VARCOND_AFTER_CURSHNDLR: a variable or condition declared after a cursor or handler
}

// resolvedSyntaxErrors names the probes OceanBase answers 1064 for although its
// grammar parses them. It reports every failure of an event body as 1064: an
// undeclared variable there answers 1064 too, where a procedure answers 1327
// (measured 2026-11-06), and RETURN is a statement of the PL grammar
// (src/pl/parser/pl_parser_mysql_mode.y, v4.3.5_CE_BP6, line 2053).
var resolvedSyntaxErrors = map[string]bool{
	"RETURN inside event": true,
}

// TestRoutineAlignment holds omni's parser to OceanBase, the oracle of this
// package (oracle_test.go), on a CREATE PROCEDURE / FUNCTION / TRIGGER / EVENT
// corpus. Each probe pins the error 5.7.25-OceanBase_CE-v4.3.5.6 answers when it
// runs the statement, 0 when it creates the object (re-measured on 2026-11-06;
// until then these were MySQL 8.0's answers). Parse must reject exactly the
// probes OceanBase answers 1064 or a code in routineParserCodes, apart from
// resolvedSyntaxErrors, and accept every other one.
//
// The categories group probes by topic; a + marks one OceanBase creates and a -
// one it refuses. A seed creates the table the probes reference, so a body
// that names it does not fail for a missing table.
func TestRoutineAlignment(t *testing.T) {
	o := startParserOracle(t)

	mustExec(t, o, "DROP TABLE IF EXISTS t")
	mustExec(t, o, "CREATE TABLE t (id INT PRIMARY KEY, v INT, c VARCHAR(50))")
	mustExec(t, o, "DROP TABLE IF EXISTS audit_log")
	mustExec(t, o, "CREATE TABLE audit_log (id INT PRIMARY KEY AUTO_INCREMENT, action VARCHAR(50))")
	mustExec(t, o, "INSERT INTO t(id, v, c) VALUES (1,1,'a'),(2,2,'b')")

	type probe struct {
		category string
		name     string
		// dropSQL is run before sql to make the probe rerunnable.
		dropSQL string
		sql     string
		// code is the error OceanBase answers, 0 when it creates the object.
		code uint16
	}

	dropProc := func(name string) string { return "DROP PROCEDURE IF EXISTS " + name }
	dropFunc := func(name string) string { return "DROP FUNCTION IF EXISTS " + name }
	dropTrig := func(name string) string { return "DROP TRIGGER IF EXISTS " + name }
	dropEvt := func(name string) string { return "DROP EVENT IF EXISTS " + name }

	probes := []probe{
		// ---------------- A. Positive baseline ----------------
		{"baseline+", "simple procedure", dropProc("p1"),
			"CREATE PROCEDURE p1() BEGIN SELECT 1; END", 0},
		{"baseline+", "procedure with IN/OUT params", dropProc("p2"),
			"CREATE PROCEDURE p2(IN x INT, OUT y INT) BEGIN SET y = x * 2; END", 0},
		{"baseline+", "function with single RETURN", dropFunc("f1"),
			"CREATE FUNCTION f1(a INT) RETURNS INT DETERMINISTIC RETURN a + 1", 0},
		{"baseline+", "function with BEGIN body and RETURN", dropFunc("f2"),
			`CREATE FUNCTION f2(a INT) RETURNS INT DETERMINISTIC
BEGIN
    DECLARE x INT;
    SET x = a + 1;
    RETURN x;
END`, 0},
		{"baseline+", "trigger BEFORE INSERT", dropTrig("trg1"),
			"CREATE TRIGGER trg1 BEFORE INSERT ON t FOR EACH ROW SET NEW.v = COALESCE(NEW.v, 0)", 0},
		{"baseline+", "event with DO", dropEvt("ev1"),
			"CREATE EVENT ev1 ON SCHEDULE EVERY 1 HOUR DO BEGIN INSERT INTO audit_log(action) VALUES ('tick'); END", 0},

		// ---------------- B. Compound flow ----------------
		{"compound+", "IF / ELSEIF / ELSE", dropProc("p_if"),
			`CREATE PROCEDURE p_if(IN x INT) BEGIN
    IF x = 0 THEN SET @r = 'zero';
    ELSEIF x > 0 THEN SET @r = 'pos';
    ELSE SET @r = 'neg';
    END IF;
END`, 0},
		{"compound+", "WHILE", dropProc("p_while"),
			`CREATE PROCEDURE p_while() BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 3 DO SET i = i + 1; END WHILE;
END`, 0},
		{"compound+", "REPEAT", dropProc("p_repeat"),
			`CREATE PROCEDURE p_repeat() BEGIN
    DECLARE i INT DEFAULT 0;
    REPEAT SET i = i + 1; UNTIL i >= 3 END REPEAT;
END`, 0},
		{"compound+", "LOOP with LEAVE", dropProc("p_loop"),
			`CREATE PROCEDURE p_loop() BEGIN
    DECLARE i INT DEFAULT 0;
    lp: LOOP SET i = i + 1; IF i >= 3 THEN LEAVE lp; END IF; END LOOP lp;
END`, 0},
		{"compound+", "CASE simple", dropProc("p_case_s"),
			`CREATE PROCEDURE p_case_s(IN x INT) BEGIN
    CASE x WHEN 1 THEN SET @r = 'one'; WHEN 2 THEN SET @r = 'two'; ELSE SET @r = 'o'; END CASE;
END`, 0},
		{"compound+", "CASE searched", dropProc("p_case_w"),
			`CREATE PROCEDURE p_case_w(IN x INT) BEGIN
    CASE WHEN x < 0 THEN SET @r = 'n'; WHEN x = 0 THEN SET @r = 'z'; ELSE SET @r = 'p'; END CASE;
END`, 0},
		{"compound+", "labeled BEGIN with end-label", dropProc("p_lbl"),
			`CREATE PROCEDURE p_lbl() my_block: BEGIN SELECT 1; END my_block`, 0},

		// ---------------- B2. Bare (BEGIN-less) compound body ----------------
		// The routine body is a single sp_proc_stmt, which may be a compound
		// statement with its own inner ';' and no surrounding BEGIN ... END. The
		// splitter must not break at that inner ';'.
		{"bare+", "bare IF / ELSEIF / ELSE body", dropProc("p_bare_if"),
			"CREATE PROCEDURE p_bare_if(IN x INT) IF x = 0 THEN SET @r = 'z'; ELSEIF x > 0 THEN SET @r = 'p'; ELSE SET @r = 'n'; END IF", 0},
		{"bare+", "bare IF body, multi-statement THEN", dropProc("p_bare_if2"),
			"CREATE PROCEDURE p_bare_if2(IN x INT) IF x = 0 THEN SET @a = 1; SET @b = 2; END IF", 0},
		{"bare+", "bare WHILE body", dropProc("p_bare_while"),
			"CREATE PROCEDURE p_bare_while(IN x INT) WHILE x > 0 DO SET x = x - 1; END WHILE", 0},
		{"bare+", "bare CASE body", dropProc("p_bare_case"),
			"CREATE PROCEDURE p_bare_case(IN x INT) CASE x WHEN 1 THEN SET @r = 'one'; ELSE SET @r = 'o'; END CASE", 0},

		// ---------------- C. DECLARE / cursor / handler ----------------
		{"decl+", "cursor + continue handler + loop", dropProc("p_cur"),
			`CREATE PROCEDURE p_cur() BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE v INT;
    DECLARE c CURSOR FOR SELECT id FROM t;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    OPEN c;
    rd: LOOP FETCH c INTO v; IF done THEN LEAVE rd; END IF; SET @last = v; END LOOP rd;
    CLOSE c;
END`, 0},
		{"decl+", "named condition + handler", dropProc("p_cnd"),
			`CREATE PROCEDURE p_cnd() BEGIN
    DECLARE dup_key CONDITION FOR SQLSTATE '23000';
    DECLARE EXIT HANDLER FOR dup_key SET @err = 1;
    INSERT INTO t(id) VALUES (1);
END`, 0},
		{"decl-", "variable declared after a handler", dropProc("p_var_after_handler"),
			"CREATE PROCEDURE p_var_after_handler() BEGIN DECLARE CONTINUE HANDLER FOR NOT FOUND SET @d = 1; DECLARE x INT; END", 1337},

		// ---------------- D. Symbol resolution (V #1-#7) ----------------
		{"symbol-", "LEAVE undeclared label", dropProc("p_lv_undecl"),
			"CREATE PROCEDURE p_lv_undecl() BEGIN LEAVE nowhere; END", 1308},
		{"symbol-", "ITERATE undeclared label", dropProc("p_it_undecl"),
			"CREATE PROCEDURE p_it_undecl() BEGIN ITERATE nowhere; END", 1308},
		{"symbol-", "OPEN undeclared cursor", dropProc("p_open_undecl"),
			"CREATE PROCEDURE p_open_undecl() BEGIN OPEN nope; END", 1324},
		{"symbol-", "FETCH undeclared cursor", dropProc("p_fetch_undecl"),
			"CREATE PROCEDURE p_fetch_undecl() BEGIN DECLARE v INT; FETCH nope INTO v; END", 1324},
		// OceanBase's grammar takes no user variable as a FETCH target.
		{"symbol-", "FETCH into a user variable", dropProc("p_fetch_uvar"),
			"CREATE PROCEDURE p_fetch_uvar() BEGIN DECLARE c CURSOR FOR SELECT 1; OPEN c; FETCH c INTO @x; CLOSE c; END", 1064},
		{"symbol-", "CLOSE undeclared cursor", dropProc("p_close_undecl"),
			"CREATE PROCEDURE p_close_undecl() BEGIN CLOSE nope; END", 1324},
		{"symbol-", "HANDLER references undeclared condition", dropProc("p_hcond_undecl"),
			"CREATE PROCEDURE p_hcond_undecl() BEGIN DECLARE EXIT HANDLER FOR no_such SET @e=1; SELECT 1; END", 1319},
		{"symbol-", "SET to undeclared bare var", dropProc("p_set_undecl"),
			"CREATE PROCEDURE p_set_undecl() BEGIN SET nope = 1; END", 1327},

		// ---------------- E. Duplicate detection (V #8-#12) ----------------
		{"dup-", "duplicate DECLARE VAR", dropProc("p_dup_var"),
			"CREATE PROCEDURE p_dup_var() BEGIN DECLARE x INT; DECLARE x INT; SELECT 1; END", 1331},
		{"dup-", "duplicate DECLARE CURSOR", dropProc("p_dup_cur"),
			"CREATE PROCEDURE p_dup_cur() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE c CURSOR FOR SELECT 2; SELECT 1; END", 1333},
		{"dup-", "duplicate DECLARE CONDITION", dropProc("p_dup_cond"),
			"CREATE PROCEDURE p_dup_cond() BEGIN DECLARE dk CONDITION FOR SQLSTATE '23000'; DECLARE dk CONDITION FOR SQLSTATE '23001'; SELECT 1; END", 1332},
		{"dup-", "duplicate label nested", dropProc("p_dup_lbl"),
			`CREATE PROCEDURE p_dup_lbl() lbl: BEGIN lbl: BEGIN SELECT 1; END lbl; END lbl`, 1309},
		{"dup-", "duplicate HANDLER condition value within DECLARE", dropProc("p_dup_hcond"),
			"CREATE PROCEDURE p_dup_hcond() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23000' SET @e=1; SELECT 1; END", 1413},

		// ---------------- F. Function RETURN coverage (V #13-#14) ----------------
		// MySQL only requires a RETURN somewhere in a function body; a path that
		// falls off the end fails when it runs, not when it is created.
		{"return+", "function missing RETURN on else path", dropFunc("f_no_ret_else"),
			`CREATE FUNCTION f_no_ret_else(x INT) RETURNS INT DETERMINISTIC
BEGIN IF x > 0 THEN RETURN 1; END IF; END`, 0},
		// The SELECT is refused (1415) before the missing RETURN is.
		{"return-", "function missing RETURN overall", dropFunc("f_no_ret"),
			`CREATE FUNCTION f_no_ret() RETURNS INT DETERMINISTIC BEGIN SELECT 1; END`, 1415},
		{"return+", "function with IF/ELSE both RETURN", dropFunc("f_ret_both"),
			`CREATE FUNCTION f_ret_both(x INT) RETURNS INT DETERMINISTIC
BEGIN IF x > 0 THEN RETURN 1; ELSE RETURN 0; END IF; END`, 0},
		{"return+", "function CASE with all WHENs + ELSE RETURN", dropFunc("f_case_ret"),
			`CREATE FUNCTION f_case_ret(x INT) RETURNS INT DETERMINISTIC
BEGIN CASE x WHEN 1 THEN RETURN 1; ELSE RETURN 0; END CASE; END`, 0},
		{"return+", "function CASE without ELSE", dropFunc("f_case_no_else"),
			`CREATE FUNCTION f_case_no_else(x INT) RETURNS INT DETERMINISTIC
BEGIN CASE x WHEN 1 THEN RETURN 1; END CASE; END`, 0},
		{"return-", "function with SIGNAL terminal", dropFunc("f_signal"),
			`CREATE FUNCTION f_signal() RETURNS INT DETERMINISTIC
BEGIN SIGNAL SQLSTATE '45000'; END`, 1320},

		// ---------------- G. RETURN outside function ----------------
		{"return-", "RETURN inside procedure", dropProc("p_with_return"),
			"CREATE PROCEDURE p_with_return() BEGIN RETURN 1; END", 1313},
		{"return-", "RETURN inside trigger", dropTrig("trg_ret"),
			"CREATE TRIGGER trg_ret BEFORE INSERT ON t FOR EACH ROW BEGIN RETURN 1; END", 1313},
		{"return-", "RETURN inside event", dropEvt("ev_ret"),
			"CREATE EVENT ev_ret ON SCHEDULE EVERY 1 HOUR DO BEGIN RETURN 1; END", 1064},

		// ---------------- H. Label scope edges ----------------
		{"label+", "LEAVE outer from nested loop", dropProc("p_lv_outer"),
			`CREATE PROCEDURE p_lv_outer() BEGIN
    o: LOOP w: WHILE @i < 3 DO LEAVE o; END WHILE w; END LOOP o;
END`, 0},
		{"label+", "LEAVE BEGIN-kind label", dropProc("p_lv_begin"),
			`CREATE PROCEDURE p_lv_begin() my: BEGIN LEAVE my; END my`, 0},
		{"label-", "ITERATE BEGIN-kind label (loop-only)", dropProc("p_it_begin"),
			`CREATE PROCEDURE p_it_begin() my: BEGIN ITERATE my; END my`, 1308},
		{"label+", "ITERATE loop label", dropProc("p_it_loop"),
			`CREATE PROCEDURE p_it_loop() BEGIN lp: LOOP ITERATE lp; END LOOP lp; END`, 0},
		{"label-", "label barrier — LEAVE outer from HANDLER body", dropProc("p_lv_barrier"),
			`CREATE PROCEDURE p_lv_barrier() my: BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN LEAVE my; END;
    SELECT 1;
END my`, 1308},
		{"label+", "label name reused across handler barrier", dropProc("p_lbl_reuse"),
			`CREATE PROCEDURE p_lbl_reuse() my: BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION my: BEGIN SELECT 1; END my;
    SELECT 2;
END my`, 0},
		{"label-", "end label differs from begin label", dropProc("p_end_label"),
			"CREATE PROCEDURE p_end_label() a: BEGIN SELECT 1; END b", 1308},
		{"label-", "end label without a begin label", dropProc("p_end_only"),
			"CREATE PROCEDURE p_end_only() BEGIN SELECT 1; END b", 1064},

		// ---------------- I. Variable shadowing / SET forms ----------------
		{"var+", "DECLARE shadows parameter in nested BEGIN", dropProc("p_shadow_param"),
			`CREATE PROCEDURE p_shadow_param(IN n INT) BEGIN DECLARE n INT; SET n = 1; END`, 0},
		{"var+", "SET parameter resolves", dropProc("p_set_param"),
			`CREATE PROCEDURE p_set_param(IN n INT) BEGIN SET n = n + 1; END`, 0},
		{"var+", "SET @user_var no scope check", dropProc("p_set_user"),
			`CREATE PROCEDURE p_set_user() BEGIN SET @anything = 1; END`, 0},
		{"var+", "SET SESSION sql_mode", dropProc("p_set_session"),
			`CREATE PROCEDURE p_set_session() BEGIN SET SESSION sql_mode = ''; END`, 0},
		{"var+", "SET NEW.col in trigger", dropTrig("trg_set_new"),
			`CREATE TRIGGER trg_set_new BEFORE INSERT ON t FOR EACH ROW SET NEW.v = 0`, 0},

		// ---------------- J. Sakila-style real procedures ----------------
		{"real+", "sakila-style cursor + handler + loop accumulation", dropFunc("f_sum"),
			`CREATE FUNCTION f_sum() RETURNS INT DETERMINISTIC
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE s INT DEFAULT 0;
    DECLARE v INT;
    DECLARE c CURSOR FOR SELECT v FROM t;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    OPEN c;
    lp: LOOP FETCH c INTO v; IF done THEN LEAVE lp; END IF; SET s = s + COALESCE(v,0); END LOOP lp;
    CLOSE c;
    RETURN s;
END`, 0},
		{"real+", "totem-dev style if(quiz) then dynamic SQL", dropProc("p_dyn"),
			`CREATE PROCEDURE p_dyn(IN quiz INT) BEGIN
    SET @s = 'SELECT * FROM t';
    if(quiz) then SET @s = CONCAT(@s, ' WHERE id > 0'); end if;
    PREPARE st FROM @s; EXECUTE st; DEALLOCATE PREPARE st;
END`, 0},

		// ---------------- K. SIGNAL / RESIGNAL / GET DIAGNOSTICS ----------------
		{"signal+", "SIGNAL with MESSAGE_TEXT", dropProc("p_signal"),
			"CREATE PROCEDURE p_signal() BEGIN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'fail'; END", 0},
		{"signal+", "RESIGNAL inside handler", dropProc("p_resig"),
			"CREATE PROCEDURE p_resig() BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION RESIGNAL; SELECT 1; END", 0},
		{"signal+", "GET DIAGNOSTICS CONDITION", dropProc("p_diag"),
			"CREATE PROCEDURE p_diag() BEGIN DECLARE n INT; GET DIAGNOSTICS CONDITION 1 n = MYSQL_ERRNO; END", 0},

		// ---------------- L. DDL inside body ----------------
		{"ddl+", "DROP TABLE IF EXISTS in body", dropProc("p_drop_if"),
			"CREATE PROCEDURE p_drop_if() BEGIN DROP TABLE IF EXISTS tmp; CREATE TABLE IF NOT EXISTS tmp(id INT); END", 0},

		// ---------------- M. Syntax errors ----------------
		{"syntax-", "statement without its semicolon", dropProc("p_semi"),
			"CREATE PROCEDURE p_semi() BEGIN SELECT 1 END", 1064},
		{"syntax-", "trigger without FOR EACH ROW", dropTrig("trg_no_each"),
			"CREATE TRIGGER trg_no_each BEFORE INSERT ON t SET NEW.v = 0", 1064},
		{"syntax-", "event without DO", dropEvt("ev_no_do"),
			"CREATE EVENT ev_no_do ON SCHEDULE EVERY 1 HOUR", 1064},
		{"syntax-", "IF closed by END", dropProc("p_if_end"),
			"CREATE PROCEDURE p_if_end() BEGIN IF 1 THEN SELECT 1; END; END", 1064},
		{"syntax-", "variable declared after a statement", dropProc("p_decl_late"),
			"CREATE PROCEDURE p_decl_late() BEGIN SELECT 1; DECLARE x INT; END", 1064},
		{"syntax-", "WHILE without DO", dropProc("p_while_no_do"),
			"CREATE PROCEDURE p_while_no_do() BEGIN WHILE 1 SELECT 1; END WHILE; END", 1064},
		{"syntax-", "REPEAT without UNTIL", dropProc("p_repeat_no_until"),
			"CREATE PROCEDURE p_repeat_no_until() BEGIN REPEAT SELECT 1; END REPEAT; END", 1064},
		{"syntax-", "LOOP closed by END", dropProc("p_loop_no_end"),
			"CREATE PROCEDURE p_loop_no_end() BEGIN LOOP SELECT 1; END; END", 1064},
		{"syntax-", "cursor without FOR", dropProc("p_cursor_no_for"),
			"CREATE PROCEDURE p_cursor_no_for() BEGIN DECLARE c CURSOR SELECT 1; END", 1064},
		{"syntax-", "RETURN without a value", dropFunc("f_ret_noexpr"),
			"CREATE FUNCTION f_ret_noexpr() RETURNS INT DETERMINISTIC RETURN", 1064},

		// A DELIMITER script holds several statements, which the driver does not
		// run as one; the splitter's tests cover those.
	}

	for _, pr := range probes {
		t.Run(pr.category+"/"+pr.name, func(t *testing.T) {
			mustExec(t, o, pr.dropSQL)
			code := routineEngineCode(t, o, pr.sql)
			if code != pr.code {
				t.Errorf("OceanBase answered %d, want %d for %s", code, pr.code, pr.sql)
			}
			_, err := Parse(pr.sql)
			wantReject := (code == 1064 && !resolvedSyntaxErrors[pr.name]) || routineParserCodes[code]
			if rejected := err != nil; rejected != wantReject {
				t.Errorf("omni rejects=%v (%v), OceanBase answered %d for %s", rejected, err, code, pr.sql)
			}
		})
	}
}

// routineEngineCode runs sql on the oracle and returns the error the server
// answered, 0 when it ran.
//
// Anything that is not a server error — a dropped connection, a timeout — is no
// verdict at all and fails the test.
func routineEngineCode(t testing.TB, o *parserOracle, sql string) uint16 {
	t.Helper()
	_, err := o.db.ExecContext(o.ctx, sql)
	if err == nil {
		return 0
	}
	var myErr *gomysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number
	}
	t.Fatalf("the oracle gave no verdict for %q: %v", sql, err)
	return 0
}

func mustExec(t testing.TB, o *parserOracle, sql string) {
	t.Helper()
	if _, err := o.db.ExecContext(o.ctx, sql); err != nil {
		t.Fatalf("seed failed: %v\nsql: %s", err, sql)
	}
}
