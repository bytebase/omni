-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_RoutineMultiObjectRoundTrip default#1
> exec
| DROP DATABASE IF EXISTS rtn_multi
< ok 0 0
> exec
| CREATE DATABASE rtn_multi DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE rtn_multi
< ok 0 0
> exec
| CREATE FUNCTION add2(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b
< ok 0 0
> exec
| CREATE FUNCTION greet(s VARCHAR(20)) RETURNS VARCHAR(40) DETERMINISTIC READS SQL DATA COMMENT 'g' RETURN CONCAT('hi ', s)
< ok 0 0
> exec
| CREATE PROCEDURE bump(INOUT n INT) BEGIN SET n = n + 1; END
< ok 0 0
> exec
| CREATE PROCEDURE noop() MODIFIES SQL DATA BEGIN SET @v = 1; END
< ok 0 0
> query
| SHOW CREATE FUNCTION rtn_multi.add2
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| add2
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `add2`(a INT, b INT) RETURNS int(11)
|     DETERMINISTIC
| RETURN a + b
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION rtn_multi.greet
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| greet
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `greet`(s VARCHAR(20)) RETURNS varchar(40)
|     READS SQL DATA
|     DETERMINISTIC
|     COMMENT 'g'
| RETURN CONCAT('hi ', s)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE PROCEDURE rtn_multi.bump
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| bump
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `bump`(INOUT n INT)
| BEGIN SET n = n + 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE PROCEDURE rtn_multi.noop
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| noop
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` PROCEDURE `noop`()
|     MODIFIES SQL DATA
| BEGIN SET @v = 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_general_ci
