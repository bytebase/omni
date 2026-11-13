-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestRoutineAlignment default#1
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, v INT, c VARCHAR(50))
< ok 0 0
> exec
| DROP TABLE IF EXISTS audit_log
< ok 0 0
> exec
| CREATE TABLE audit_log (id INT PRIMARY KEY AUTO_INCREMENT, action VARCHAR(50))
< ok 0 0
> exec
| INSERT INTO t(id, v, c) VALUES (1,1,'a'),(2,2,'b')
< ok 2 0
> exec
| DROP PROCEDURE IF EXISTS p1
< ok 0 0
> exec
| CREATE PROCEDURE p1() BEGIN SELECT 1; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p2
< ok 0 0
> exec
| CREATE PROCEDURE p2(IN x INT, OUT y INT) BEGIN SET y = x * 2; END
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f1
< ok 0 0
> exec
| CREATE FUNCTION f1(a INT) RETURNS INT DETERMINISTIC RETURN a + 1
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f2
< ok 0 0
> exec
| CREATE FUNCTION f2(a INT) RETURNS INT DETERMINISTIC
| BEGIN
|     DECLARE x INT;
|     SET x = a + 1;
|     RETURN x;
| END
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS trg1
< ok 0 0
> exec
| CREATE TRIGGER trg1 BEFORE INSERT ON t FOR EACH ROW SET NEW.v = COALESCE(NEW.v, 0)
< ok 0 0
> exec
| DROP EVENT IF EXISTS ev1
< ok 0 0
> exec
| CREATE EVENT ev1 ON SCHEDULE EVERY 1 HOUR DO BEGIN INSERT INTO audit_log(action) VALUES ('tick'); END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_if
< ok 0 0
> exec
| CREATE PROCEDURE p_if(IN x INT) BEGIN
|     IF x = 0 THEN SET @r = 'zero';
|     ELSEIF x > 0 THEN SET @r = 'pos';
|     ELSE SET @r = 'neg';
|     END IF;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_while
< ok 0 0
> exec
| CREATE PROCEDURE p_while() BEGIN
|     DECLARE i INT DEFAULT 0;
|     WHILE i < 3 DO SET i = i + 1; END WHILE;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_repeat
< ok 0 0
> exec
| CREATE PROCEDURE p_repeat() BEGIN
|     DECLARE i INT DEFAULT 0;
|     REPEAT SET i = i + 1; UNTIL i >= 3 END REPEAT;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_loop
< ok 0 0
> exec
| CREATE PROCEDURE p_loop() BEGIN
|     DECLARE i INT DEFAULT 0;
|     lp: LOOP SET i = i + 1; IF i >= 3 THEN LEAVE lp; END IF; END LOOP lp;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_case_s
< ok 0 0
> exec
| CREATE PROCEDURE p_case_s(IN x INT) BEGIN
|     CASE x WHEN 1 THEN SET @r = 'one'; WHEN 2 THEN SET @r = 'two'; ELSE SET @r = 'o'; END CASE;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_case_w
< ok 0 0
> exec
| CREATE PROCEDURE p_case_w(IN x INT) BEGIN
|     CASE WHEN x < 0 THEN SET @r = 'n'; WHEN x = 0 THEN SET @r = 'z'; ELSE SET @r = 'p'; END CASE;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_lbl
< ok 0 0
> exec
| CREATE PROCEDURE p_lbl() my_block: BEGIN SELECT 1; END my_block
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_bare_if
< ok 0 0
> exec
| CREATE PROCEDURE p_bare_if(IN x INT) IF x = 0 THEN SET @r = 'z'; ELSEIF x > 0 THEN SET @r = 'p'; ELSE SET @r = 'n'; END IF
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_bare_if2
< ok 0 0
> exec
| CREATE PROCEDURE p_bare_if2(IN x INT) IF x = 0 THEN SET @a = 1; SET @b = 2; END IF
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_bare_while
< ok 0 0
> exec
| CREATE PROCEDURE p_bare_while(IN x INT) WHILE x > 0 DO SET x = x - 1; END WHILE
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_bare_case
< ok 0 0
> exec
| CREATE PROCEDURE p_bare_case(IN x INT) CASE x WHEN 1 THEN SET @r = 'one'; ELSE SET @r = 'o'; END CASE
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_cur
< ok 0 0
> exec
| CREATE PROCEDURE p_cur() BEGIN
|     DECLARE done INT DEFAULT 0;
|     DECLARE v INT;
|     DECLARE c CURSOR FOR SELECT id FROM t;
|     DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
|     OPEN c;
|     rd: LOOP FETCH c INTO v; IF done THEN LEAVE rd; END IF; SET @last = v; END LOOP rd;
|     CLOSE c;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_cnd
< ok 0 0
> exec
| CREATE PROCEDURE p_cnd() BEGIN
|     DECLARE dup_key CONDITION FOR SQLSTATE '23000';
|     DECLARE EXIT HANDLER FOR dup_key SET @err = 1;
|     INSERT INTO t(id) VALUES (1);
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_var_after_handler
< ok 0 0
> exec
| CREATE PROCEDURE p_var_after_handler() BEGIN DECLARE CONTINUE HANDLER FOR NOT FOUND SET @d = 1; DECLARE x INT; END
< error
| 1337 42000 Variable or condition declaration after cursor or handler declaration
> exec
| DROP PROCEDURE IF EXISTS p_lv_undecl
< ok 0 0
> exec
| CREATE PROCEDURE p_lv_undecl() BEGIN LEAVE nowhere; END
< error
| 1308 42000 LEAVE with no matching label: nowhere
> exec
| DROP PROCEDURE IF EXISTS p_it_undecl
< ok 0 0
> exec
| CREATE PROCEDURE p_it_undecl() BEGIN ITERATE nowhere; END
< error
| 1308 42000 ITERATE with no matching label: nowhere
> exec
| DROP PROCEDURE IF EXISTS p_open_undecl
< ok 0 0
> exec
| CREATE PROCEDURE p_open_undecl() BEGIN OPEN nope; END
< error
| 1324 42000 Undefined CURSOR: nope
> exec
| DROP PROCEDURE IF EXISTS p_fetch_undecl
< ok 0 0
> exec
| CREATE PROCEDURE p_fetch_undecl() BEGIN FETCH nope INTO @x; END
< error
| 1324 42000 Undefined CURSOR: nope
> exec
| DROP PROCEDURE IF EXISTS p_close_undecl
< ok 0 0
> exec
| CREATE PROCEDURE p_close_undecl() BEGIN CLOSE nope; END
< error
| 1324 42000 Undefined CURSOR: nope
> exec
| DROP PROCEDURE IF EXISTS p_hcond_undecl
< ok 0 0
> exec
| CREATE PROCEDURE p_hcond_undecl() BEGIN DECLARE EXIT HANDLER FOR no_such SET @e=1; SELECT 1; END
< error
| 1319 42000 Undefined CONDITION: no_such
> exec
| DROP PROCEDURE IF EXISTS p_set_undecl
< ok 0 0
> exec
| CREATE PROCEDURE p_set_undecl() BEGIN SET nope = 1; END
< error
| 1193 HY000 Unknown system variable 'nope'
> exec
| DROP PROCEDURE IF EXISTS p_dup_var
< ok 0 0
> exec
| CREATE PROCEDURE p_dup_var() BEGIN DECLARE x INT; DECLARE x INT; SELECT 1; END
< error
| 1331 42000 Duplicate variable: x
> exec
| DROP PROCEDURE IF EXISTS p_dup_cur
< ok 0 0
> exec
| CREATE PROCEDURE p_dup_cur() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE c CURSOR FOR SELECT 2; SELECT 1; END
< error
| 1333 42000 Duplicate cursor: c
> exec
| DROP PROCEDURE IF EXISTS p_dup_cond
< ok 0 0
> exec
| CREATE PROCEDURE p_dup_cond() BEGIN DECLARE dk CONDITION FOR SQLSTATE '23000'; DECLARE dk CONDITION FOR SQLSTATE '23001'; SELECT 1; END
< error
| 1332 42000 Duplicate condition: dk
> exec
| DROP PROCEDURE IF EXISTS p_dup_lbl
< ok 0 0
> exec
| CREATE PROCEDURE p_dup_lbl() lbl: BEGIN lbl: BEGIN SELECT 1; END lbl; END lbl
< error
| 1309 42000 Redefining label lbl
> exec
| DROP PROCEDURE IF EXISTS p_dup_hcond
< ok 0 0
> exec
| CREATE PROCEDURE p_dup_hcond() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23000' SET @e=1; SELECT 1; END
< error
| 1413 42000 Duplicate handler declared in the same block
> exec
| DROP FUNCTION IF EXISTS f_no_ret_else
< ok 0 0
> exec
| CREATE FUNCTION f_no_ret_else(x INT) RETURNS INT DETERMINISTIC
| BEGIN IF x > 0 THEN RETURN 1; END IF; END
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f_no_ret
< ok 0 0
> exec
| CREATE FUNCTION f_no_ret() RETURNS INT DETERMINISTIC BEGIN SELECT 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec
| DROP FUNCTION IF EXISTS f_ret_both
< ok 0 0
> exec
| CREATE FUNCTION f_ret_both(x INT) RETURNS INT DETERMINISTIC
| BEGIN IF x > 0 THEN RETURN 1; ELSE RETURN 0; END IF; END
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f_case_ret
< ok 0 0
> exec
| CREATE FUNCTION f_case_ret(x INT) RETURNS INT DETERMINISTIC
| BEGIN CASE x WHEN 1 THEN RETURN 1; ELSE RETURN 0; END CASE; END
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f_case_no_else
< ok 0 0
> exec
| CREATE FUNCTION f_case_no_else(x INT) RETURNS INT DETERMINISTIC
| BEGIN CASE x WHEN 1 THEN RETURN 1; END CASE; END
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f_signal
< ok 0 0
> exec
| CREATE FUNCTION f_signal() RETURNS INT DETERMINISTIC
| BEGIN SIGNAL SQLSTATE '45000'; END
< error
| 1320 42000 No RETURN found in FUNCTION test.f_signal
> exec
| DROP PROCEDURE IF EXISTS p_with_return
< ok 0 0
> exec
| CREATE PROCEDURE p_with_return() BEGIN RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec
| DROP TRIGGER IF EXISTS trg_ret
< ok 0 0
> exec
| CREATE TRIGGER trg_ret BEFORE INSERT ON t FOR EACH ROW BEGIN RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec
| DROP EVENT IF EXISTS ev_ret
< ok 0 0
> exec
| CREATE EVENT ev_ret ON SCHEDULE EVERY 1 HOUR DO BEGIN RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec
| DROP PROCEDURE IF EXISTS p_lv_outer
< ok 0 0
> exec
| CREATE PROCEDURE p_lv_outer() BEGIN
|     o: LOOP w: WHILE @i < 3 DO LEAVE o; END WHILE w; END LOOP o;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_lv_begin
< ok 0 0
> exec
| CREATE PROCEDURE p_lv_begin() my: BEGIN LEAVE my; END my
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_it_begin
< ok 0 0
> exec
| CREATE PROCEDURE p_it_begin() my: BEGIN ITERATE my; END my
< error
| 1308 42000 ITERATE with no matching label: my
> exec
| DROP PROCEDURE IF EXISTS p_it_loop
< ok 0 0
> exec
| CREATE PROCEDURE p_it_loop() BEGIN lp: LOOP ITERATE lp; END LOOP lp; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_lv_barrier
< ok 0 0
> exec
| CREATE PROCEDURE p_lv_barrier() my: BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN LEAVE my; END;
|     SELECT 1;
| END my
< error
| 1308 42000 LEAVE with no matching label: my
> exec
| DROP PROCEDURE IF EXISTS p_lbl_reuse
< ok 0 0
> exec
| CREATE PROCEDURE p_lbl_reuse() my: BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION my: BEGIN SELECT 1; END my;
|     SELECT 2;
| END my
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_end_label
< ok 0 0
> exec
| CREATE PROCEDURE p_end_label() a: BEGIN SELECT 1; END b
< error
| 1310 42000 End-label b without match
> exec
| DROP PROCEDURE IF EXISTS p_shadow_param
< ok 0 0
> exec
| CREATE PROCEDURE p_shadow_param(IN n INT) BEGIN DECLARE n INT; SET n = 1; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_set_param
< ok 0 0
> exec
| CREATE PROCEDURE p_set_param(IN n INT) BEGIN SET n = n + 1; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_set_user
< ok 0 0
> exec
| CREATE PROCEDURE p_set_user() BEGIN SET @anything = 1; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_set_session
< ok 0 0
> exec
| CREATE PROCEDURE p_set_session() BEGIN SET SESSION sql_mode = ''; END
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS trg_set_new
< ok 0 0
> exec
| CREATE TRIGGER trg_set_new BEFORE INSERT ON t FOR EACH ROW SET NEW.v = 0
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f_sum
< ok 0 0
> exec
| CREATE FUNCTION f_sum() RETURNS INT DETERMINISTIC
| BEGIN
|     DECLARE done INT DEFAULT 0;
|     DECLARE s INT DEFAULT 0;
|     DECLARE v INT;
|     DECLARE c CURSOR FOR SELECT v FROM t;
|     DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
|     OPEN c;
|     lp: LOOP FETCH c INTO v; IF done THEN LEAVE lp; END IF; SET s = s + COALESCE(v,0); END LOOP lp;
|     CLOSE c;
|     RETURN s;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_dyn
< ok 0 0
> exec
| CREATE PROCEDURE p_dyn(IN quiz INT) BEGIN
|     SET @s = 'SELECT * FROM t';
|     if(quiz) then SET @s = CONCAT(@s, ' WHERE id > 0'); end if;
|     PREPARE st FROM @s; EXECUTE st; DEALLOCATE PREPARE st;
| END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_signal
< ok 0 0
> exec
| CREATE PROCEDURE p_signal() BEGIN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'fail'; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_resig
< ok 0 0
> exec
| CREATE PROCEDURE p_resig() BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION RESIGNAL; SELECT 1; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_diag
< ok 0 0
> exec
| CREATE PROCEDURE p_diag() BEGIN DECLARE n INT; GET DIAGNOSTICS CONDITION 1 n = MYSQL_ERRNO; END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_drop_if
< ok 0 0
> exec
| CREATE PROCEDURE p_drop_if() BEGIN DROP TABLE IF EXISTS tmp; CREATE TABLE IF NOT EXISTS tmp(id INT); END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_semi
< ok 0 0
> exec
| CREATE PROCEDURE p_semi() BEGIN SELECT 1 END
< error
| 1064 42000 …near '' at line 1
> exec
| DROP TRIGGER IF EXISTS trg_no_each
< ok 0 0
> exec
| CREATE TRIGGER trg_no_each BEFORE INSERT ON t SET NEW.v = 0
< error
| 1064 42000 …near 'SET NEW.v = 0' at line 1
> exec
| DROP EVENT IF EXISTS ev_no_do
< ok 0 0
> exec
| CREATE EVENT ev_no_do ON SCHEDULE EVERY 1 HOUR
< error
| 1064 42000 …near '' at line 1
> exec
| DROP PROCEDURE IF EXISTS p_if_end
< ok 0 0
> exec
| CREATE PROCEDURE p_if_end() BEGIN IF 1 THEN SELECT 1; END; END
< error
| 1064 42000 …near '; END' at line 1
> exec
| DROP PROCEDURE IF EXISTS p_decl_late
< ok 0 0
> exec
| CREATE PROCEDURE p_decl_late() BEGIN SELECT 1; DECLARE x INT; END
< error
| 1064 42000 …near 'DECLARE x INT; END' at line 1
> exec
| DROP PROCEDURE IF EXISTS p_while_no_do
< ok 0 0
> exec
| CREATE PROCEDURE p_while_no_do() BEGIN WHILE 1 SELECT 1; END WHILE; END
< error
| 1064 42000 …near 'SELECT 1; END WHILE; END' at line 1
> exec
| DROP PROCEDURE IF EXISTS p_repeat_no_until
< ok 0 0
> exec
| CREATE PROCEDURE p_repeat_no_until() BEGIN REPEAT SELECT 1; END REPEAT; END
< error
| 1064 42000 …near 'END REPEAT; END' at line 1
> exec
| DROP PROCEDURE IF EXISTS p_loop_no_end
< ok 0 0
> exec
| CREATE PROCEDURE p_loop_no_end() BEGIN LOOP SELECT 1; END; END
< error
| 1064 42000 …near '; END' at line 1
> exec
| DROP PROCEDURE IF EXISTS p_cursor_no_for
< ok 0 0
> exec
| CREATE PROCEDURE p_cursor_no_for() BEGIN DECLARE c CURSOR SELECT 1; END
< error
| 1064 42000 …near 'SELECT 1; END' at line 1
> exec
| DROP FUNCTION IF EXISTS f_ret_noexpr
< ok 0 0
> exec
| CREATE FUNCTION f_ret_noexpr() RETURNS INT DETERMINISTIC RETURN
< error
| 1064 42000 …near '' at line 1
