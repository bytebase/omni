package parser

import (
	"context"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestRoutineAlignment compares omni's accept/reject decision against
// MySQL 8.0 (testcontainers oracle) on a curated CREATE PROCEDURE /
// FUNCTION / TRIGGER / EVENT corpus. For each probe we record:
//
//   - omni:  Parse(sql) returned no error?
//   - mysql: container exec(sql) returned no error of any kind?
//
// Alignment is exact when both decisions agree. The test is gated by
// -short (the existing oracle convention in this package). It does not
// fail on disagreement; it logs a structured matrix so reviewers can
// see exactly where omni and MySQL diverge.
//
// Container is started once via existing startParserOracle helper.
// A pre-test cleanup creates the schema objects each probe references
// (no-op tables / functions) so semantic refs that aren't part of the
// alignment scope don't dominate the failure column.
func TestRoutineAlignment(t *testing.T) {
	o := startParserOracle(t)

	// Seed a referent schema so probes that DROP / SELECT FROM existing
	// names don't fail for non-existence reasons.
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
	}

	dropProc := func(name string) string { return "DROP PROCEDURE IF EXISTS " + name }
	dropFunc := func(name string) string { return "DROP FUNCTION IF EXISTS " + name }
	dropTrig := func(name string) string { return "DROP TRIGGER IF EXISTS " + name }
	dropEvt := func(name string) string { return "DROP EVENT IF EXISTS " + name }

	probes := []probe{
		// ---------------- A. Positive baseline ----------------
		{"baseline+", "simple procedure", dropProc("p1"),
			"CREATE PROCEDURE p1() BEGIN SELECT 1; END"},
		{"baseline+", "procedure with IN/OUT params", dropProc("p2"),
			"CREATE PROCEDURE p2(IN x INT, OUT y INT) BEGIN SET y = x * 2; END"},
		{"baseline+", "function with single RETURN", dropFunc("f1"),
			"CREATE FUNCTION f1(a INT) RETURNS INT DETERMINISTIC RETURN a + 1"},
		{"baseline+", "function with BEGIN body and RETURN", dropFunc("f2"),
			`CREATE FUNCTION f2(a INT) RETURNS INT DETERMINISTIC
BEGIN
    DECLARE x INT;
    SET x = a + 1;
    RETURN x;
END`},
		{"baseline+", "trigger BEFORE INSERT", dropTrig("trg1"),
			"CREATE TRIGGER trg1 BEFORE INSERT ON t FOR EACH ROW SET NEW.v = COALESCE(NEW.v, 0)"},
		{"baseline+", "event with DO", dropEvt("ev1"),
			"CREATE EVENT ev1 ON SCHEDULE EVERY 1 HOUR DO BEGIN INSERT INTO audit_log(action) VALUES ('tick'); END"},

		// ---------------- B. Compound flow ----------------
		{"compound+", "IF / ELSEIF / ELSE", dropProc("p_if"),
			`CREATE PROCEDURE p_if(IN x INT) BEGIN
    IF x = 0 THEN SET @r = 'zero';
    ELSEIF x > 0 THEN SET @r = 'pos';
    ELSE SET @r = 'neg';
    END IF;
END`},
		{"compound+", "WHILE", dropProc("p_while"),
			`CREATE PROCEDURE p_while() BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 3 DO SET i = i + 1; END WHILE;
END`},
		{"compound+", "REPEAT", dropProc("p_repeat"),
			`CREATE PROCEDURE p_repeat() BEGIN
    DECLARE i INT DEFAULT 0;
    REPEAT SET i = i + 1; UNTIL i >= 3 END REPEAT;
END`},
		{"compound+", "LOOP with LEAVE", dropProc("p_loop"),
			`CREATE PROCEDURE p_loop() BEGIN
    DECLARE i INT DEFAULT 0;
    lp: LOOP SET i = i + 1; IF i >= 3 THEN LEAVE lp; END IF; END LOOP lp;
END`},
		{"compound+", "CASE simple", dropProc("p_case_s"),
			`CREATE PROCEDURE p_case_s(IN x INT) BEGIN
    CASE x WHEN 1 THEN SET @r = 'one'; WHEN 2 THEN SET @r = 'two'; ELSE SET @r = 'o'; END CASE;
END`},
		{"compound+", "CASE searched", dropProc("p_case_w"),
			`CREATE PROCEDURE p_case_w(IN x INT) BEGIN
    CASE WHEN x < 0 THEN SET @r = 'n'; WHEN x = 0 THEN SET @r = 'z'; ELSE SET @r = 'p'; END CASE;
END`},
		{"compound+", "labeled BEGIN with end-label", dropProc("p_lbl"),
			`CREATE PROCEDURE p_lbl() my_block: BEGIN SELECT 1; END my_block`},

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
END`},
		{"decl+", "named condition + handler", dropProc("p_cnd"),
			`CREATE PROCEDURE p_cnd() BEGIN
    DECLARE dup_key CONDITION FOR SQLSTATE '23000';
    DECLARE EXIT HANDLER FOR dup_key SET @err = 1;
    INSERT INTO t(id) VALUES (1);
END`},

		// ---------------- D. Symbol resolution (V #1-#7) ----------------
		{"symbol-", "LEAVE undeclared label", dropProc("p_lv_undecl"),
			"CREATE PROCEDURE p_lv_undecl() BEGIN LEAVE nowhere; END"},
		{"symbol-", "ITERATE undeclared label", dropProc("p_it_undecl"),
			"CREATE PROCEDURE p_it_undecl() BEGIN ITERATE nowhere; END"},
		{"symbol-", "OPEN undeclared cursor", dropProc("p_open_undecl"),
			"CREATE PROCEDURE p_open_undecl() BEGIN OPEN nope; END"},
		{"symbol-", "FETCH undeclared cursor", dropProc("p_fetch_undecl"),
			"CREATE PROCEDURE p_fetch_undecl() BEGIN FETCH nope INTO @x; END"},
		{"symbol-", "CLOSE undeclared cursor", dropProc("p_close_undecl"),
			"CREATE PROCEDURE p_close_undecl() BEGIN CLOSE nope; END"},
		{"symbol-", "HANDLER references undeclared condition", dropProc("p_hcond_undecl"),
			"CREATE PROCEDURE p_hcond_undecl() BEGIN DECLARE EXIT HANDLER FOR no_such SET @e=1; SELECT 1; END"},
		{"symbol-", "SET to undeclared bare var", dropProc("p_set_undecl"),
			"CREATE PROCEDURE p_set_undecl() BEGIN SET nope = 1; END"},

		// ---------------- E. Duplicate detection (V #8-#12) ----------------
		{"dup-", "duplicate DECLARE VAR", dropProc("p_dup_var"),
			"CREATE PROCEDURE p_dup_var() BEGIN DECLARE x INT; DECLARE x INT; SELECT 1; END"},
		{"dup-", "duplicate DECLARE CURSOR", dropProc("p_dup_cur"),
			"CREATE PROCEDURE p_dup_cur() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE c CURSOR FOR SELECT 2; SELECT 1; END"},
		{"dup-", "duplicate DECLARE CONDITION", dropProc("p_dup_cond"),
			"CREATE PROCEDURE p_dup_cond() BEGIN DECLARE dk CONDITION FOR SQLSTATE '23000'; DECLARE dk CONDITION FOR SQLSTATE '23001'; SELECT 1; END"},
		{"dup-", "duplicate label nested", dropProc("p_dup_lbl"),
			`CREATE PROCEDURE p_dup_lbl() lbl: BEGIN lbl: BEGIN SELECT 1; END lbl; END lbl`},
		{"dup-", "duplicate HANDLER condition value within DECLARE", dropProc("p_dup_hcond"),
			"CREATE PROCEDURE p_dup_hcond() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23000' SET @e=1; SELECT 1; END"},

		// ---------------- F. Function RETURN coverage (V #13-#14) ----------------
		{"return-", "function missing RETURN on else path", dropFunc("f_no_ret_else"),
			`CREATE FUNCTION f_no_ret_else(x INT) RETURNS INT DETERMINISTIC
BEGIN IF x > 0 THEN RETURN 1; END IF; END`},
		{"return-", "function missing RETURN overall", dropFunc("f_no_ret"),
			`CREATE FUNCTION f_no_ret() RETURNS INT DETERMINISTIC BEGIN SELECT 1; END`},
		{"return+", "function with IF/ELSE both RETURN", dropFunc("f_ret_both"),
			`CREATE FUNCTION f_ret_both(x INT) RETURNS INT DETERMINISTIC
BEGIN IF x > 0 THEN RETURN 1; ELSE RETURN 0; END IF; END`},
		{"return+", "function CASE with all WHENs + ELSE RETURN", dropFunc("f_case_ret"),
			`CREATE FUNCTION f_case_ret(x INT) RETURNS INT DETERMINISTIC
BEGIN CASE x WHEN 1 THEN RETURN 1; ELSE RETURN 0; END CASE; END`},
		{"return-", "function CASE without ELSE", dropFunc("f_case_no_else"),
			`CREATE FUNCTION f_case_no_else(x INT) RETURNS INT DETERMINISTIC
BEGIN CASE x WHEN 1 THEN RETURN 1; END CASE; END`},
		{"return+", "function with SIGNAL terminal", dropFunc("f_signal"),
			`CREATE FUNCTION f_signal() RETURNS INT DETERMINISTIC
BEGIN SIGNAL SQLSTATE '45000'; END`},

		// ---------------- G. RETURN outside function ----------------
		{"return-", "RETURN inside procedure", dropProc("p_with_return"),
			"CREATE PROCEDURE p_with_return() BEGIN RETURN 1; END"},
		{"return-", "RETURN inside trigger", dropTrig("trg_ret"),
			"CREATE TRIGGER trg_ret BEFORE INSERT ON t FOR EACH ROW BEGIN RETURN 1; END"},
		{"return-", "RETURN inside event", dropEvt("ev_ret"),
			"CREATE EVENT ev_ret ON SCHEDULE EVERY 1 HOUR DO BEGIN RETURN 1; END"},

		// ---------------- H. Label scope edges ----------------
		{"label+", "LEAVE outer from nested loop", dropProc("p_lv_outer"),
			`CREATE PROCEDURE p_lv_outer() BEGIN
    o: LOOP w: WHILE @i < 3 DO LEAVE o; END WHILE w; END LOOP o;
END`},
		{"label+", "LEAVE BEGIN-kind label", dropProc("p_lv_begin"),
			`CREATE PROCEDURE p_lv_begin() my: BEGIN LEAVE my; END my`},
		{"label-", "ITERATE BEGIN-kind label (loop-only)", dropProc("p_it_begin"),
			`CREATE PROCEDURE p_it_begin() my: BEGIN ITERATE my; END my`},
		{"label+", "ITERATE loop label", dropProc("p_it_loop"),
			`CREATE PROCEDURE p_it_loop() BEGIN lp: LOOP ITERATE lp; END LOOP lp; END`},
		{"label-", "label barrier — LEAVE outer from HANDLER body", dropProc("p_lv_barrier"),
			`CREATE PROCEDURE p_lv_barrier() my: BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN LEAVE my; END;
    SELECT 1;
END my`},
		{"label+", "label name reused across handler barrier", dropProc("p_lbl_reuse"),
			`CREATE PROCEDURE p_lbl_reuse() my: BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION my: BEGIN SELECT 1; END my;
    SELECT 2;
END my`},

		// ---------------- I. Variable shadowing / SET forms ----------------
		{"var+", "DECLARE shadows parameter in nested BEGIN", dropProc("p_shadow_param"),
			`CREATE PROCEDURE p_shadow_param(IN n INT) BEGIN DECLARE n INT; SET n = 1; END`},
		{"var+", "SET parameter resolves", dropProc("p_set_param"),
			`CREATE PROCEDURE p_set_param(IN n INT) BEGIN SET n = n + 1; END`},
		{"var+", "SET @user_var no scope check", dropProc("p_set_user"),
			`CREATE PROCEDURE p_set_user() BEGIN SET @anything = 1; END`},
		{"var+", "SET SESSION sql_mode", dropProc("p_set_session"),
			`CREATE PROCEDURE p_set_session() BEGIN SET SESSION sql_mode = ''; END`},
		{"var+", "SET NEW.col in trigger", dropTrig("trg_set_new"),
			`CREATE TRIGGER trg_set_new BEFORE INSERT ON t FOR EACH ROW SET NEW.v = 0`},

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
END`},
		{"real+", "totem-dev style if(quiz) then dynamic SQL", dropProc("p_dyn"),
			`CREATE PROCEDURE p_dyn(IN quiz INT) BEGIN
    SET @s = 'SELECT * FROM t';
    if(quiz) then SET @s = CONCAT(@s, ' WHERE id > 0'); end if;
    PREPARE st FROM @s; EXECUTE st; DEALLOCATE PREPARE st;
END`},

		// ---------------- K. SIGNAL / RESIGNAL / GET DIAGNOSTICS ----------------
		{"signal+", "SIGNAL with MESSAGE_TEXT", dropProc("p_signal"),
			"CREATE PROCEDURE p_signal() BEGIN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'fail'; END"},
		{"signal+", "RESIGNAL inside handler", dropProc("p_resig"),
			"CREATE PROCEDURE p_resig() BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION RESIGNAL; SELECT 1; END"},
		{"signal+", "GET DIAGNOSTICS CONDITION", dropProc("p_diag"),
			"CREATE PROCEDURE p_diag() BEGIN DECLARE n INT; GET DIAGNOSTICS CONDITION 1 n = MYSQL_ERRNO; END"},

		// ---------------- L. DDL inside body ----------------
		{"ddl+", "DROP TABLE IF EXISTS in body", dropProc("p_drop_if"),
			"CREATE PROCEDURE p_drop_if() BEGIN DROP TABLE IF EXISTS tmp; CREATE TABLE IF NOT EXISTS tmp(id INT); END"},

		// ---------------- M. DELIMITER + multi-procedure (Parse-only; oracle skips) ----------------
		// Multi-stmt scripts can't be exec'd as one shot via the driver
		// in the same way (they're handled at split time). Skip oracle.
	}

	type result struct {
		category    string
		name        string
		omniAccept  bool
		mysqlAccept bool
		omniErr     string
		mysqlErr    string
	}
	var results []result
	matchCounts := map[string]int{"agree": 0, "omni-strict": 0, "omni-loose": 0}

	for _, pr := range probes {
		// Drop any prior version in MySQL.
		_, _ = o.db.ExecContext(o.ctx, pr.dropSQL)

		// omni
		_, omniErr := Parse(pr.sql)
		omniAccept := omniErr == nil

		// mysql
		_, mysqlExecErr := o.db.ExecContext(o.ctx, pr.sql)
		mysqlAccept := mysqlExecErr == nil

		r := result{
			category:    pr.category,
			name:        pr.name,
			omniAccept:  omniAccept,
			mysqlAccept: mysqlAccept,
		}
		if omniErr != nil {
			r.omniErr = trimErrLineCol(omniErr.Error())
		}
		if mysqlExecErr != nil {
			r.mysqlErr = mysqlErrText(mysqlExecErr)
		}
		results = append(results, r)

		switch {
		case omniAccept == mysqlAccept:
			matchCounts["agree"]++
		case !omniAccept && mysqlAccept:
			matchCounts["omni-strict"]++
		case omniAccept && !mysqlAccept:
			matchCounts["omni-loose"]++
		}
	}

	// Print structured report.
	var sb strings.Builder
	sb.WriteString("\n=== routine accept/reject alignment vs MySQL 8.0 ===\n")
	sb.WriteString("agree: " + itoaAudit(matchCounts["agree"]))
	sb.WriteString("  omni-strict (omni rejects, mysql accepts): " + itoaAudit(matchCounts["omni-strict"]))
	sb.WriteString("  omni-loose  (omni accepts, mysql rejects): " + itoaAudit(matchCounts["omni-loose"]))
	sb.WriteString("  total: " + itoaAudit(len(results)) + "\n\n")

	// Disagreements first.
	sb.WriteString("--- disagreements ---\n")
	disagreements := 0
	for _, r := range results {
		if r.omniAccept == r.mysqlAccept {
			continue
		}
		disagreements++
		marker := "OMNI-STRICT"
		if r.omniAccept {
			marker = "OMNI-LOOSE "
		}
		sb.WriteString("  [" + marker + "] [" + r.category + "] " + r.name + "\n")
		sb.WriteString("     omni:  ")
		if r.omniAccept {
			sb.WriteString("ACCEPT\n")
		} else {
			sb.WriteString("REJECT (" + r.omniErr + ")\n")
		}
		sb.WriteString("     mysql: ")
		if r.mysqlAccept {
			sb.WriteString("ACCEPT\n")
		} else {
			sb.WriteString("REJECT (" + r.mysqlErr + ")\n")
		}
	}
	if disagreements == 0 {
		sb.WriteString("  (none)\n")
	}

	// Per-category summary.
	sb.WriteString("\n--- per-category summary ---\n")
	type catTotal struct{ agree, total int }
	catSum := map[string]*catTotal{}
	for _, r := range results {
		if catSum[r.category] == nil {
			catSum[r.category] = &catTotal{}
		}
		catSum[r.category].total++
		if r.omniAccept == r.mysqlAccept {
			catSum[r.category].agree++
		}
	}
	cats := []string{
		"baseline+", "compound+", "decl+",
		"symbol-", "dup-", "return-", "return+",
		"label+", "label-", "var+",
		"signal+", "ddl+", "real+",
	}
	for _, c := range cats {
		if v, ok := catSum[c]; ok {
			sb.WriteString("  " + padRight(c, 12) + ": " + itoaAudit(v.agree) + "/" + itoaAudit(v.total) + "\n")
		}
	}

	t.Log(sb.String())
}

func mustExec(t *testing.T, o *parserOracle, sql string) {
	t.Helper()
	if _, err := o.db.ExecContext(context.Background(), sql); err != nil {
		t.Fatalf("seed failed: %v\nsql: %s", err, sql)
	}
}

func mysqlErrText(err error) string {
	if my, ok := err.(*gomysql.MySQLError); ok {
		msg := my.Message
		if i := strings.Index(msg, "\n"); i >= 0 {
			msg = msg[:i]
		}
		return "ERR " + itoaAudit(int(my.Number)) + ": " + msg
	}
	return err.Error()
}

func trimErrLineCol(s string) string {
	if i := strings.Index(s, " (line"); i >= 0 {
		s = s[:i]
	}
	return s
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
