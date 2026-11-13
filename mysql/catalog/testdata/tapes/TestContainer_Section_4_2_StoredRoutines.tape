-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_4_2_StoredRoutines default#1
> exec
| DROP FUNCTION IF EXISTS fn_add
< ok 0 0
> exec
| CREATE FUNCTION fn_add(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b
< ok 0 0
> query
| SHOW CREATE FUNCTION fn_add
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_add
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_add`(a INT, b INT) RETURNS int
|     DETERMINISTIC
| RETURN a + b
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP PROCEDURE IF EXISTS sp_greet
< ok 0 0
> exec
| CREATE PROCEDURE sp_greet(IN name VARCHAR(100)) BEGIN SELECT CONCAT('Hello, ', name); END
< ok 0 0
> query
| SHOW CREATE PROCEDURE sp_greet
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| sp_greet
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_greet`(IN name VARCHAR(100))
| BEGIN SELECT CONCAT('Hello, ', name); END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP FUNCTION IF EXISTS fn_drop_test
< ok 0 0
> exec
| CREATE FUNCTION fn_drop_test() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| DROP FUNCTION fn_drop_test
< ok 0 0
> query
| SHOW CREATE FUNCTION fn_drop_test
< error
| 1305 42000 FUNCTION fn_drop_test does not exist
> exec
| DROP PROCEDURE IF EXISTS sp_drop_test
< ok 0 0
> exec
| CREATE PROCEDURE sp_drop_test() BEGIN SELECT 1; END
< ok 0 0
> exec
| DROP PROCEDURE sp_drop_test
< ok 0 0
> query
| SHOW CREATE PROCEDURE sp_drop_test
< error
| 1305 42000 PROCEDURE sp_drop_test does not exist
> exec
| DROP FUNCTION IF EXISTS fn_nonexistent_xyz
< ok 0 0
> exec
| DROP FUNCTION fn_totally_nonexistent
< error
| 1305 42000 FUNCTION test.fn_totally_nonexistent does not exist
> exec
| DROP FUNCTION IF EXISTS fn_alter_test
< ok 0 0
> exec
| CREATE FUNCTION fn_alter_test() RETURNS INT DETERMINISTIC RETURN 42
< ok 0 0
> exec
| ALTER FUNCTION fn_alter_test COMMENT 'test comment'
< ok 0 0
> query
| SHOW CREATE FUNCTION fn_alter_test
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_alter_test
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_alter_test`() RETURNS int
|     DETERMINISTIC
|     COMMENT 'test comment'
| RETURN 42
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP PROCEDURE IF EXISTS sp_alter_test
< ok 0 0
> exec
| CREATE PROCEDURE sp_alter_test() BEGIN SELECT 1; END
< ok 0 0
> exec
| ALTER PROCEDURE sp_alter_test SQL SECURITY INVOKER
< ok 0 0
> query
| SHOW CREATE PROCEDURE sp_alter_test
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| sp_alter_test
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_alter_test`()
|     SQL SECURITY INVOKER
| BEGIN SELECT 1; END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP FUNCTION IF EXISTS fn_show_test
< ok 0 0
> exec
| CREATE FUNCTION fn_show_test(x INT) RETURNS VARCHAR(100) DETERMINISTIC RETURN CONCAT('val=', x)
< ok 0 0
> query
| SHOW CREATE FUNCTION fn_show_test
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| fn_show_test
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `fn_show_test`(x INT) RETURNS varchar(100) CHARSET utf8mb4
|     DETERMINISTIC
| RETURN CONCAT('val=', x)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP PROCEDURE IF EXISTS sp_show_test
< ok 0 0
> exec
| CREATE PROCEDURE sp_show_test(IN id INT, OUT result VARCHAR(100)) BEGIN SET result = CONCAT('id=', id); END
< ok 0 0
> query
| SHOW CREATE PROCEDURE sp_show_test
< result
< column "Procedure" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| sp_show_test
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `sp_show_test`(IN id INT, OUT result VARCHAR(100))
| BEGIN SET result = CONCAT('id=', id); END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
