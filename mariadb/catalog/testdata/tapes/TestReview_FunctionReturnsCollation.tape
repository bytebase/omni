-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_FunctionReturnsCollation default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| DROP DATABASE IF EXISTS fdb
< ok 0 0
> exec
| CREATE DATABASE fdb CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
< ok 1 0
> exec
| USE fdb
< ok 0 0
> exec
| CREATE FUNCTION f_plain() RETURNS VARCHAR(20) RETURN 'a'
< ok 0 0
> exec
| CREATE FUNCTION f_cs_other() RETURNS VARCHAR(20) CHARACTER SET utf8mb3 RETURN 'a'
< ok 0 0
> exec
| CREATE FUNCTION f_cs_same() RETURNS VARCHAR(20) CHARACTER SET utf8mb4 RETURN 'a'
< ok 0 0
> exec
| CREATE FUNCTION f_text() RETURNS TEXT RETURN 'a'
< ok 0 0
> exec
| CREATE FUNCTION f_collate() RETURNS VARCHAR(20) COLLATE utf8mb4_general_ci RETURN 'a'
< ok 0 0
> exec
| CREATE FUNCTION f_int() RETURNS INT RETURN 1
< ok 0 0
> query
| SHOW CREATE FUNCTION f_plain
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f_plain
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f_plain`() RETURNS varchar(20) CHARSET utf8mb4 COLLATE utf8mb4_bin
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_bin
> query
| SHOW CREATE FUNCTION f_cs_other
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f_cs_other
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f_cs_other`() RETURNS varchar(20) CHARSET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_bin
> query
| SHOW CREATE FUNCTION f_cs_same
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f_cs_same
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f_cs_same`() RETURNS varchar(20) CHARSET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_bin
> query
| SHOW CREATE FUNCTION f_text
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f_text
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f_text`() RETURNS text CHARSET utf8mb4 COLLATE utf8mb4_bin
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_bin
> query
| SHOW CREATE FUNCTION f_collate
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f_collate
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f_collate`() RETURNS varchar(20) CHARSET utf8mb4 COLLATE utf8mb4_general_ci
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_bin
> query
| SHOW CREATE FUNCTION f_int
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f_int
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f_int`() RETURNS int(11)
| RETURN 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_bin
