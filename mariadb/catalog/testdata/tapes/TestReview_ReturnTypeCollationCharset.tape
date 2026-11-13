-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ReturnTypeCollationCharset default#1
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
| CREATE FUNCTION rf0() RETURNS VARCHAR(20) COLLATE latin1_bin DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf0
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf0
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf0`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf0
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf0`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf1() RETURNS VARCHAR(20) COLLATE utf8mb4_bin DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf1
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf1
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf1`() RETURNS varchar(20) CHARSET utf8mb4 COLLATE utf8mb4_bin
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf1
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf1`() RETURNS varchar(20) CHARSET utf8mb4 COLLATE utf8mb4_bin
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf2() RETURNS VARCHAR(20) COLLATE utf8mb3_bin DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf2
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf2
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf2`() RETURNS varchar(20) CHARSET utf8mb3 COLLATE utf8mb3_bin
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf2
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf2`() RETURNS varchar(20) CHARSET utf8mb3 COLLATE utf8mb3_bin
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf3() RETURNS VARCHAR(20) COLLATE latin1_swedish_ci DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf3
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf3
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf3`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_swedish_ci
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf3
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf3`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_swedish_ci
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf4() RETURNS VARCHAR(20) CHARACTER SET latin1 DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf4
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf4
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf4`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_swedish_ci
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf4
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf4`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_swedish_ci
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf5() RETURNS VARCHAR(20) CHARACTER SET latin1 COLLATE latin1_bin DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf5
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf5
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf5`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf5
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf5`() RETURNS varchar(20) CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf6() RETURNS VARCHAR(20) DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf6
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf6
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf6`() RETURNS varchar(20) CHARSET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf6
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf6`() RETURNS varchar(20) CHARSET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf7() RETURNS TEXT COLLATE latin1_bin DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf7
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf7
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf7`() RETURNS text CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf7
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf7`() RETURNS text CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
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
| CREATE FUNCTION rf8() RETURNS CHAR(5) COLLATE latin1_bin DETERMINISTIC RETURN 'a'
< ok 0 0
> query
| SHOW CREATE FUNCTION rf8
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| rf8
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `rf8`() RETURNS char(5) CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP FUNCTION rf8
< ok 0 0
> exec
| CREATE DEFINER=`root`@`%` FUNCTION `rf8`() RETURNS char(5) CHARSET latin1 COLLATE latin1_bin
|     DETERMINISTIC
| RETURN 'a'
< ok 0 0
