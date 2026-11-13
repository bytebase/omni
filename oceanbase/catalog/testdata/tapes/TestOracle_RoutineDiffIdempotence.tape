-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_RoutineDiffIdempotence default#1
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_simple
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_simple
< ok 0 0
> exec
| CREATE FUNCTION f(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_simple.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT, b INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + b
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_string_return
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_string_return DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_string_return
< ok 0 0
> exec
| CREATE FUNCTION f(s VARCHAR(50)) RETURNS VARCHAR(100) DETERMINISTIC RETURN CONCAT(s, s)
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_string_return.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(s VARCHAR(50)) RETURNS varchar(100)
|     DETERMINISTIC
| RETURN CONCAT(s, s)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_all_characteristics
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_all_characteristics DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_all_characteristics
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA SQL SECURITY INVOKER COMMENT 'hi' RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_all_characteristics.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS int(11)
|     READS SQL DATA
|     DETERMINISTIC
|     SQL SECURITY INVOKER
|     COMMENT 'hi'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_no_sql
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_no_sql DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_no_sql
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS INT NO SQL RETURN 42
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_no_sql.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`() RETURNS int(11)
|     NO SQL
| RETURN 42
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_modifies
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_modifies DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_modifies
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC MODIFIES SQL DATA RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_modifies.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS int(11)
|     MODIFIES SQL DATA
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_explicit_defaults
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_explicit_defaults DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_explicit_defaults
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_explicit_defaults.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_language_sql
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_language_sql DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_language_sql
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS INT LANGUAGE SQL DETERMINISTIC RETURN 7
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_language_sql.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`() RETURNS int(11)
|     DETERMINISTIC
| RETURN 7
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_begin_body
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_begin_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_begin_body
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC
| BEGIN
|   DECLARE r INT;
|   SET r = a * 2;
|   RETURN r;
| END
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_begin_body.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS int(11)
|     DETERMINISTIC
| BEGIN
|   DECLARE r INT;
|   SET r = a * 2;
|   RETURN r;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_comment_quote
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_comment_quote DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_comment_quote
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'it''s here' RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_comment_quote.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'it's here'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_comment_backslash
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_comment_backslash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_comment_backslash
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'a\\b' RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_comment_backslash.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'a\b'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_decimal_return
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_decimal_return DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_decimal_return
< ok 0 0
> exec
| CREATE FUNCTION f(a DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_decimal_return.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a DECIMAL(10,2)) RETURNS decimal(10,2)
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_fn_enum_return
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_fn_enum_return DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_fn_enum_return
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS ENUM('a','b','c') DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_idem_fn_enum_return.f
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`() RETURNS enum('a','b','c')
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_enum_set_params
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_enum_set_params DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_enum_set_params
< ok 0 0
> exec
| CREATE PROCEDURE p(IN s SET('x','y'), IN e ENUM('p','q')) BEGIN SET @v = s; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_enum_set_params.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN s SET('x','y'), IN e ENUM('p','q'))
| BEGIN SET @v = s; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_reserved_param
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_reserved_param DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_reserved_param
< ok 0 0
> exec
| CREATE PROCEDURE p(IN `select` INT, IN `from` INT) BEGIN SET @v = `select`; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_reserved_param.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN `select` INT, IN `from` INT)
| BEGIN SET @v = `select`; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_in_out_inout
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_in_out_inout DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_in_out_inout
< ok 0 0
> exec
| CREATE PROCEDURE p(IN x INT, OUT y INT, INOUT z INT) BEGIN SET y = x + z; SET z = z + 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_in_out_inout.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN x INT, OUT y INT, INOUT z INT)
| BEGIN SET y = x + z; SET z = z + 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_no_params
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_no_params DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_no_params
< ok 0 0
> exec
| CREATE PROCEDURE p() BEGIN DECLARE n INT; SET n = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_no_params.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`()
| BEGIN DECLARE n INT; SET n = 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_modifies
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_modifies DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_modifies
< ok 0 0
> exec
| CREATE PROCEDURE p() MODIFIES SQL DATA BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_modifies.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`()
|     MODIFIES SQL DATA
| BEGIN SET @v = 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_default_in_direction
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_default_in_direction DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_default_in_direction
< ok 0 0
> exec
| CREATE PROCEDURE p(x INT, y INT) BEGIN SET @v = x + y; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_default_in_direction.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(x INT, y INT)
| BEGIN SET @v = x + y; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_all_characteristics
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_all_characteristics DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_all_characteristics
< ok 0 0
> exec
| CREATE PROCEDURE p(IN a INT) DETERMINISTIC READS SQL DATA SQL SECURITY INVOKER COMMENT 'pc' BEGIN SET @v = a; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_all_characteristics.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN a INT)
|     READS SQL DATA
|     DETERMINISTIC
|     SQL SECURITY INVOKER
|     COMMENT 'pc'
| BEGIN SET @v = a; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_explicit_defaults
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_explicit_defaults DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_explicit_defaults
< ok 0 0
> exec
| CREATE PROCEDURE p(IN a INT) NOT DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER BEGIN SET @v = a; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_explicit_defaults.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN a INT)
| BEGIN SET @v = a; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_idem_proc_adjacent_literals
< ok 0 0
> exec
| CREATE DATABASE rtn_idem_proc_adjacent_literals DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_idem_proc_adjacent_literals
< ok 0 0
> exec
| CREATE PROCEDURE p()
| BEGIN
|     SELECT CONCAT('tmp disk tables: ', 3, '\n'
|                   'select scan: ', 4, '\n');
| END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_idem_proc_adjacent_literals.p
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| p
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `p`()
| BEGIN
|     SELECT CONCAT('tmp disk tables: ', 3, '\n'
|                   'select scan: ', 4, '\n');
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
