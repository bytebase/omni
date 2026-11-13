-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestRoutineValidationOracleAlignment default#1
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
| DROP PROCEDURE IF EXISTS p_ok
< ok 0 0
> exec
| CREATE PROCEDURE p_ok() BEGIN SELECT 1; END
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f_ok
< ok 0 0
> exec
| CREATE FUNCTION f_ok(a INT) RETURNS INT DETERMINISTIC RETURN a + 1
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS trg_ok
< ok 0 0
> exec
| CREATE TRIGGER trg_ok BEFORE INSERT ON t FOR EACH ROW SET NEW.v = COALESCE(NEW.v, 0)
< ok 0 0
> exec
| DROP EVENT IF EXISTS ev_ok
< ok 0 0
> exec
| CREATE EVENT ev_ok ON SCHEDULE EVERY 1 HOUR DO BEGIN INSERT INTO audit_log(action) VALUES ('tick'); END
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS p_leave_undeclared
< ok 0 0
> exec
| CREATE PROCEDURE p_leave_undeclared() BEGIN LEAVE nowhere; END
< error
| 1308 42000 LEAVE with no matching label: nowhere
> exec
| DROP PROCEDURE IF EXISTS p_iterate_undeclared
< ok 0 0
> exec
| CREATE PROCEDURE p_iterate_undeclared() BEGIN ITERATE nowhere; END
< error
| 1308 42000 ITERATE with no matching label: nowhere
> exec
| DROP PROCEDURE IF EXISTS p_iterate_begin_label
< ok 0 0
> exec
| CREATE PROCEDURE p_iterate_begin_label() my_block: BEGIN ITERATE my_block; END my_block
< error
| 1308 42000 ITERATE with no matching label: my_block
> exec
| DROP PROCEDURE IF EXISTS p_open_undeclared
< ok 0 0
> exec
| CREATE PROCEDURE p_open_undeclared() BEGIN OPEN nope; END
< error
| 1324 42000 Undefined CURSOR: nope
> exec
| DROP PROCEDURE IF EXISTS p_fetch_undeclared
< ok 0 0
> exec
| CREATE PROCEDURE p_fetch_undeclared() BEGIN FETCH nope INTO @x; END
< error
| 1324 42000 Undefined CURSOR: nope
> exec
| DROP PROCEDURE IF EXISTS p_close_undeclared
< ok 0 0
> exec
| CREATE PROCEDURE p_close_undeclared() BEGIN CLOSE nope; END
< error
| 1324 42000 Undefined CURSOR: nope
> exec
| DROP PROCEDURE IF EXISTS p_handler_condition_undeclared
< ok 0 0
> exec
| CREATE PROCEDURE p_handler_condition_undeclared() BEGIN DECLARE EXIT HANDLER FOR no_such SET @e = 1; SELECT 1; END
< error
| 1319 42000 Undefined CONDITION: no_such
> exec
| DROP PROCEDURE IF EXISTS p_set_undeclared
< ok 0 0
> exec
| CREATE PROCEDURE p_set_undeclared() BEGIN SET nope = 1; END
< error
| 1193 HY000 Unknown system variable 'nope'
> exec
| DROP PROCEDURE IF EXISTS p_duplicate_var
< ok 0 0
> exec
| CREATE PROCEDURE p_duplicate_var() BEGIN DECLARE x INT; DECLARE x INT; SELECT 1; END
< error
| 1331 42000 Duplicate variable: x
> exec
| DROP PROCEDURE IF EXISTS p_duplicate_cursor
< ok 0 0
> exec
| CREATE PROCEDURE p_duplicate_cursor() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE c CURSOR FOR SELECT 2; SELECT 1; END
< error
| 1333 42000 Duplicate cursor: c
> exec
| DROP PROCEDURE IF EXISTS p_duplicate_condition
< ok 0 0
> exec
| CREATE PROCEDURE p_duplicate_condition() BEGIN DECLARE dk CONDITION FOR SQLSTATE '23000'; DECLARE dk CONDITION FOR SQLSTATE '23001'; SELECT 1; END
< error
| 1332 42000 Duplicate condition: dk
> exec
| DROP PROCEDURE IF EXISTS p_duplicate_label
< ok 0 0
> exec
| CREATE PROCEDURE p_duplicate_label() lbl: BEGIN lbl: BEGIN SELECT 1; END lbl; END lbl
< error
| 1309 42000 Redefining label lbl
> exec
| DROP PROCEDURE IF EXISTS p_duplicate_handler_condition
< ok 0 0
> exec
| CREATE PROCEDURE p_duplicate_handler_condition() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23000' SET @e = 1; SELECT 1; END
< error
| 1413 42000 Duplicate handler declared in the same block
> exec
| DROP FUNCTION IF EXISTS f_no_return
< ok 0 0
> exec
| CREATE FUNCTION f_no_return() RETURNS INT DETERMINISTIC BEGIN SELECT 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec
| DROP FUNCTION IF EXISTS f_signal_only
< ok 0 0
> exec
| CREATE FUNCTION f_signal_only() RETURNS INT DETERMINISTIC BEGIN SIGNAL SQLSTATE '45000'; END
< error
| 1320 42000 No RETURN found in FUNCTION omni_t_mysql_validate.f_signal_only
> exec
| DROP PROCEDURE IF EXISTS p_return
< ok 0 0
> exec
| CREATE PROCEDURE p_return() BEGIN RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec
| DROP TRIGGER IF EXISTS trg_return
< ok 0 0
> exec
| CREATE TRIGGER trg_return BEFORE INSERT ON t FOR EACH ROW BEGIN RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec
| DROP EVENT IF EXISTS ev_return
< ok 0 0
> exec
| CREATE EVENT ev_return ON SCHEDULE EVERY 1 HOUR DO BEGIN RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec
| DROP PROCEDURE IF EXISTS p_handler_leave_outer
< ok 0 0
> exec
| CREATE PROCEDURE p_handler_leave_outer() my_block: BEGIN
|     DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN LEAVE my_block; END;
|     SELECT 1;
| END my_block
< error
| 1308 42000 LEAVE with no matching label: my_block
> exec
| DROP FUNCTION IF EXISTS f_return_one_branch
< ok 0 0
> exec
| CREATE FUNCTION f_return_one_branch(x INT) RETURNS INT DETERMINISTIC BEGIN IF x > 0 THEN RETURN 1; END IF; END
< ok 0 0
