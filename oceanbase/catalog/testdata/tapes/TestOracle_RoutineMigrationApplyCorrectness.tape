-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_RoutineMigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_fn_simple
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_simple
< ok 0 0
> exec
| CREATE FUNCTION f(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_simple.f
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
| DROP DATABASE IF EXISTS rtn_apply_create_fn_simple
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_simple
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_create_fn_simple`.`f`(`a` INT, `b` INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + b
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_simple.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT, `b` INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + b
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_fn_string
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_string DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_string
< ok 0 0
> exec
| CREATE FUNCTION f(s VARCHAR(50)) RETURNS VARCHAR(100) DETERMINISTIC RETURN CONCAT(s, s)
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_string.f
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
| DROP DATABASE IF EXISTS rtn_apply_create_fn_string
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_string DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_string
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_create_fn_string`.`f`(`s` VARCHAR(50)) RETURNS varchar(100) CHARSET utf8mb4
|     DETERMINISTIC
| RETURN CONCAT(s, s)
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_string.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`s` VARCHAR(50)) RETURNS varchar(100)
|     DETERMINISTIC
| RETURN CONCAT(s, s)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_fn_all_chars
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_all_chars DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_all_chars
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA SQL SECURITY INVOKER COMMENT 'c' RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_all_chars.f
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
|     COMMENT 'c'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_fn_all_chars
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_all_chars DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_all_chars
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_create_fn_all_chars`.`f`(`a` INT) RETURNS int(11)
|     DETERMINISTIC
|     READS SQL DATA
|     SQL SECURITY INVOKER
|     COMMENT 'c'
| RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_all_chars.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT) RETURNS int(11)
|     READS SQL DATA
|     DETERMINISTIC
|     SQL SECURITY INVOKER
|     COMMENT 'c'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_fn_begin
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_begin DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_begin
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
| SHOW CREATE FUNCTION rtn_apply_create_fn_begin.f
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
| DROP DATABASE IF EXISTS rtn_apply_create_fn_begin
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_begin DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_begin
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_create_fn_begin`.`f`(`a` INT) RETURNS int(11)
|     DETERMINISTIC
| BEGIN
|   DECLARE r INT;
|   SET r = a * 2;
|   RETURN r;
| END
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_begin.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT) RETURNS int(11)
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_inout
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_inout DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_inout
< ok 0 0
> exec
| CREATE PROCEDURE p(IN x INT, OUT y INT, INOUT z INT) BEGIN SET y = x + z; SET z = z + 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_inout.p
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_inout
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_inout DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_inout
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_create_proc_inout`.`p`(IN `x` INT, OUT `y` INT, INOUT `z` INT)
| BEGIN SET y = x + z; SET z = z + 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_inout.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN `x` INT, OUT `y` INT, INOUT `z` INT)
| BEGIN SET y = x + z; SET z = z + 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_proc_modifies
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_modifies DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_modifies
< ok 0 0
> exec
| CREATE PROCEDURE p() MODIFIES SQL DATA BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_modifies.p
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_modifies
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_modifies DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_modifies
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_create_proc_modifies`.`p`()
|     MODIFIES SQL DATA
| BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_modifies.p
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_comment
< ok 0 0
> exec
| CREATE PROCEDURE p(IN a INT) SQL SECURITY INVOKER COMMENT 'pc' BEGIN SET @v = a; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_comment.p
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_comment
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_create_proc_comment`.`p`(IN `a` INT)
|     SQL SECURITY INVOKER
|     COMMENT 'pc'
| BEGIN SET @v = a; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_comment.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN `a` INT)
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
| DROP DATABASE IF EXISTS rtn_apply_create_fn_decimal
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_decimal DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_decimal
< ok 0 0
> exec
| CREATE FUNCTION f(a DECIMAL(10,2)) RETURNS DECIMAL(12,4) DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_decimal.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a DECIMAL(10,2)) RETURNS decimal(12,4)
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_fn_decimal
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_decimal DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_decimal
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_create_fn_decimal`.`f`(`a` DECIMAL(10,2)) RETURNS decimal(12,4)
|     DETERMINISTIC
| RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_decimal.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` DECIMAL(10,2)) RETURNS decimal(12,4)
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_fn_enum_return
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_enum_return DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_enum_return
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS ENUM('a','b','c') DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_enum_return.f
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
| DROP DATABASE IF EXISTS rtn_apply_create_fn_enum_return
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_fn_enum_return DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_fn_enum_return
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_create_fn_enum_return`.`f`() RETURNS enum('a','b','c') CHARSET utf8mb4
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_create_fn_enum_return.f
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_enum_set_params
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_enum_set_params DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_enum_set_params
< ok 0 0
> exec
| CREATE PROCEDURE p(IN s SET('x','y'), IN e ENUM('p','q')) BEGIN SET @v = s; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_enum_set_params.p
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_enum_set_params
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_enum_set_params DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_enum_set_params
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_create_proc_enum_set_params`.`p`(IN `s` SET('x','y'), IN `e` ENUM('p','q'))
| BEGIN SET @v = s; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_enum_set_params.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN `s` SET('x','y'), IN `e` ENUM('p','q'))
| BEGIN SET @v = s; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_create_proc_reserved_param
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_reserved_param DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_reserved_param
< ok 0 0
> exec
| CREATE PROCEDURE p(IN `select` INT, IN `from` INT) BEGIN SET @v = `select`; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_reserved_param.p
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
| DROP DATABASE IF EXISTS rtn_apply_create_proc_reserved_param
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_create_proc_reserved_param DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_create_proc_reserved_param
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_create_proc_reserved_param`.`p`(IN `select` INT, IN `from` INT)
| BEGIN SET @v = `select`; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_create_proc_reserved_param.p
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
| DROP DATABASE IF EXISTS rtn_apply_drop_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_drop_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_drop_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_drop_fn.f
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
| DROP DATABASE IF EXISTS rtn_apply_drop_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_drop_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_drop_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> exec
| DROP FUNCTION `rtn_apply_drop_fn`.`f`
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_drop_fn.f
< error
| 1305 42000 procedure/function does not exist
> exec
| DROP DATABASE IF EXISTS rtn_apply_drop_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_drop_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_drop_proc
< ok 0 0
> exec
| CREATE PROCEDURE p() BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_drop_proc.p
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
| BEGIN SET @v = 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_drop_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_drop_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_drop_proc
< ok 0 0
> exec
| CREATE PROCEDURE p() BEGIN SET @v = 1; END
< ok 0 0
> exec
| DROP PROCEDURE `rtn_apply_drop_proc`.`p`
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_drop_proc.p
< error
| 1305 42000 procedure/function does not exist
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_comment
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'new' RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_comment.f
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
|     COMMENT 'new'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_comment
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_comment.f
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
|     COMMENT 'old'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_comment
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a
< ok 0 0
> exec
| ALTER FUNCTION `rtn_apply_alter_fn_comment`.`f` CONTAINS SQL SQL SECURITY DEFINER COMMENT 'new'
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_comment.f
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
|     COMMENT 'new'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_security
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_security
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC SQL SECURITY INVOKER RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_security.f
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
|     SQL SECURITY INVOKER
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_security
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_security
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC SQL SECURITY DEFINER RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_security.f
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
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_security
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_security
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC SQL SECURITY DEFINER RETURN a
< ok 0 0
> exec
| ALTER FUNCTION `rtn_apply_alter_fn_security`.`f` CONTAINS SQL SQL SECURITY INVOKER COMMENT ''
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_security.f
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
|     SQL SECURITY INVOKER
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_dataaccess
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_dataaccess DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_dataaccess
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_dataaccess.f
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
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_dataaccess
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_dataaccess DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_dataaccess
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_dataaccess.f
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
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_dataaccess
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_dataaccess DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_dataaccess
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> exec
| ALTER FUNCTION `rtn_apply_alter_fn_dataaccess`.`f` READS SQL DATA SQL SECURITY DEFINER COMMENT ''
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_dataaccess.f
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
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_add_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_add_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_add_comment
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'added' RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_add_comment.f
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
|     COMMENT 'added'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_add_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_add_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_add_comment
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_add_comment.f
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
| DROP DATABASE IF EXISTS rtn_apply_alter_fn_add_comment
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_fn_add_comment DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_fn_add_comment
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> exec
| ALTER FUNCTION `rtn_apply_alter_fn_add_comment`.`f` CONTAINS SQL SQL SECURITY DEFINER COMMENT 'added'
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_alter_fn_add_comment.f
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
|     COMMENT 'added'
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_proc_security
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_proc_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_proc_security
< ok 0 0
> exec
| CREATE PROCEDURE p(IN a INT) SQL SECURITY INVOKER BEGIN SET @v = a; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_alter_proc_security.p
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
|     SQL SECURITY INVOKER
| BEGIN SET @v = a; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_proc_security
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_proc_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_proc_security
< ok 0 0
> exec
| CREATE PROCEDURE p(IN a INT) SQL SECURITY DEFINER BEGIN SET @v = a; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_alter_proc_security.p
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
| DROP DATABASE IF EXISTS rtn_apply_alter_proc_security
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_proc_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_proc_security
< ok 0 0
> exec
| CREATE PROCEDURE p(IN a INT) SQL SECURITY DEFINER BEGIN SET @v = a; END
< ok 0 0
> exec
| ALTER PROCEDURE `rtn_apply_alter_proc_security`.`p` CONTAINS SQL SQL SECURITY INVOKER COMMENT ''
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_alter_proc_security.p
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
|     SQL SECURITY INVOKER
| BEGIN SET @v = a; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_proc_dataaccess
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_proc_dataaccess DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_proc_dataaccess
< ok 0 0
> exec
| CREATE PROCEDURE p() MODIFIES SQL DATA BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_alter_proc_dataaccess.p
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
| DROP DATABASE IF EXISTS rtn_apply_alter_proc_dataaccess
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_proc_dataaccess DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_proc_dataaccess
< ok 0 0
> exec
| CREATE PROCEDURE p() CONTAINS SQL BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_alter_proc_dataaccess.p
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
| BEGIN SET @v = 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_alter_proc_dataaccess
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_alter_proc_dataaccess DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_alter_proc_dataaccess
< ok 0 0
> exec
| CREATE PROCEDURE p() CONTAINS SQL BEGIN SET @v = 1; END
< ok 0 0
> exec
| ALTER PROCEDURE `rtn_apply_alter_proc_dataaccess`.`p` MODIFIES SQL DATA SQL SECURITY DEFINER COMMENT ''
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_alter_proc_dataaccess.p
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
| DROP DATABASE IF EXISTS rtn_apply_body_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 2
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_body_change_fn.f
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
| RETURN a + 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_body_change_fn.f
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
| RETURN a + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1
< ok 0 0
> exec
| DROP FUNCTION `rtn_apply_body_change_fn`.`f`
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_body_change_fn`.`f`(`a` INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + 2
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_body_change_fn.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_param_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_param_change_fn.f
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
| DROP DATABASE IF EXISTS rtn_apply_param_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_param_change_fn.f
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
| DROP DATABASE IF EXISTS rtn_apply_param_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> exec
| DROP FUNCTION `rtn_apply_param_change_fn`.`f`
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_param_change_fn`.`f`(`a` INT, `b` INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + b
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_param_change_fn.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT, `b` INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + b
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_returns_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_returns_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_returns_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS BIGINT DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_returns_change_fn.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS bigint(20)
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_returns_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_returns_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_returns_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_returns_change_fn.f
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
| DROP DATABASE IF EXISTS rtn_apply_returns_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_returns_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_returns_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a
< ok 0 0
> exec
| DROP FUNCTION `rtn_apply_returns_change_fn`.`f`
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_returns_change_fn`.`f`(`a` INT) RETURNS bigint(20)
|     DETERMINISTIC
| RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_returns_change_fn.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT) RETURNS bigint(20)
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_determ_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_determ_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_determ_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_determ_change_fn.f
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
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_determ_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_determ_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_determ_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT NOT DETERMINISTIC READS SQL DATA RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_determ_change_fn.f
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
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_determ_change_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_determ_change_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_determ_change_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT NOT DETERMINISTIC READS SQL DATA RETURN a
< ok 0 0
> exec
| DROP FUNCTION `rtn_apply_determ_change_fn`.`f`
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_determ_change_fn`.`f`(`a` INT) RETURNS int(11)
|     DETERMINISTIC
|     READS SQL DATA
| RETURN a
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_determ_change_fn.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT) RETURNS int(11)
|     READS SQL DATA
|     DETERMINISTIC
| RETURN a
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p() BEGIN SET @v = 2; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_body_change_proc.p
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
| BEGIN SET @v = 2; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p() BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_body_change_proc.p
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
| BEGIN SET @v = 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p() BEGIN SET @v = 1; END
< ok 0 0
> exec
| DROP PROCEDURE `rtn_apply_body_change_proc`.`p`
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_body_change_proc`.`p`()
| BEGIN SET @v = 2; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_body_change_proc.p
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
| BEGIN SET @v = 2; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_param_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p(IN x INT, OUT y INT) BEGIN SET y = x; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_param_change_proc.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN x INT, OUT y INT)
| BEGIN SET y = x; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_param_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p(IN x INT) BEGIN SET @v = x; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_param_change_proc.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN x INT)
| BEGIN SET @v = x; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_param_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p(IN x INT) BEGIN SET @v = x; END
< ok 0 0
> exec
| DROP PROCEDURE `rtn_apply_param_change_proc`.`p`
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_param_change_proc`.`p`(IN `x` INT, OUT `y` INT)
| BEGIN SET y = x; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_param_change_proc.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN `x` INT, OUT `y` INT)
| BEGIN SET y = x; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_param_enum_case_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_enum_case_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_enum_case_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p(IN e ENUM('A','b')) BEGIN SET @v = e; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_param_enum_case_change_proc.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN e ENUM('A','b'))
| BEGIN SET @v = e; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_param_enum_case_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_enum_case_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_enum_case_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p(IN e ENUM('a','b')) BEGIN SET @v = e; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_param_enum_case_change_proc.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN e ENUM('a','b'))
| BEGIN SET @v = e; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_param_enum_case_change_proc
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_param_enum_case_change_proc DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_param_enum_case_change_proc
< ok 0 0
> exec
| CREATE PROCEDURE p(IN e ENUM('a','b')) BEGIN SET @v = e; END
< ok 0 0
> exec
| DROP PROCEDURE `rtn_apply_param_enum_case_change_proc`.`p`
< ok 0 0
> exec
| CREATE PROCEDURE `rtn_apply_param_enum_case_change_proc`.`p`(IN `e` ENUM('A','b'))
| BEGIN SET @v = e; END
< ok 0 0
> query
| SHOW CREATE PROCEDURE rtn_apply_param_enum_case_change_proc.p
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
| CREATE DEFINER = `root`@`%` PROCEDURE `p`(IN `e` ENUM('A','b'))
| BEGIN SET @v = e; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_and_comment_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_and_comment_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_and_comment_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'new' RETURN a + 2
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_body_and_comment_fn.f
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
|     COMMENT 'new'
| RETURN a + 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_and_comment_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_and_comment_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_and_comment_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a + 1
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_body_and_comment_fn.f
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
|     COMMENT 'old'
| RETURN a + 1
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS rtn_apply_body_and_comment_fn
< ok 0 0
> exec
| CREATE DATABASE rtn_apply_body_and_comment_fn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_apply_body_and_comment_fn
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a + 1
< ok 0 0
> exec
| DROP FUNCTION `rtn_apply_body_and_comment_fn`.`f`
< ok 0 0
> exec
| CREATE FUNCTION `rtn_apply_body_and_comment_fn`.`f`(`a` INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'new'
| RETURN a + 2
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_apply_body_and_comment_fn.f
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
| CREATE DEFINER = `root`@`%` FUNCTION `f`(`a` INT) RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'new'
| RETURN a + 2
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
