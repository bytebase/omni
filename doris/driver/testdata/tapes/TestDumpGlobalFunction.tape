-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestDumpGlobalFunction default#1
> exec
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(INT) WITH PARAMETER(x) AS abs(x);
< ok 0 0
> exec
|
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(INT, INT) WITH PARAMETER(x, y) AS greatest(x, y);
< ok 0 0
> exec
|
| CREATE GLOBAL ALIAS FUNCTION dump_global_unused(INT) WITH PARAMETER(x) AS abs(x);
< ok 0 0
> exec
|
| CREATE ALIAS FUNCTION dump_local(INT) WITH PARAMETER(x) AS abs(x);
< ok 0 0
> exec
|
| CREATE VIEW calls_functions AS SELECT dump_global_called(-1) AS g, dump_local(-2) AS l
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_global_function
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_global_function
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_global_function") ORDER BY Name
< result
< column "Name" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_global_function
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< row
: bytes
| calls_functions
: bytes
| VIEW
> query
| SHOW FULL FUNCTIONS FROM `omni_dump_global_function`
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| dump_local(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_global_function
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE VIEW `omni_dump_global_function`.`calls_functions`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| calls_functions
: bytes
| CREATE VIEW `calls_functions`
| (g,l)
|  AS SELECT `omni_dump_global_function`.`dump_global_called`(-1) AS `g`, `omni_dump_global_function`.`dump_local`(-2) AS `l`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE FUNCTION dump_local(int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_local(int)
: bytes
| CREATE ALIAS FUNCTION dump_local(int) WITH PARAMETER(x) AS abs(`x`);
> query
| SHOW GLOBAL FULL FUNCTIONS
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| dump_global_called(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
< row
: bytes
| dump_global_called(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"greatest(`x`, `y`)"}
< row
: bytes
| dump_global_unused(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
> query
| SHOW CREATE GLOBAL FUNCTION dump_global_called(int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_global_called(int)
: bytes
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(int) WITH PARAMETER(x) AS abs(`x`);
> query
| SHOW CREATE GLOBAL FUNCTION dump_global_called(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_global_called(int, int)
: bytes
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(int, int) WITH PARAMETER(x, y) AS greatest(`x`, `y`);
> exec
| DROP DATABASE `omni_dump_global_function`
< ok 0 0
> exec
| CREATE DATABASE `omni_dump_global_function`
< ok 0 0
> exec
| USE `omni_dump_global_function`
< ok 0 0
> exec
| --
| -- Global function structure for dump_global_called(int)
| --
| CREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS dump_global_called(int) WITH PARAMETER(x) AS abs(`x`);
< ok 0 0
> exec
|
|
| --
| -- Global function structure for dump_global_called(int, int)
| --
| CREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS dump_global_called(int, int) WITH PARAMETER(x, y) AS greatest(`x`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for dump_local(int)
| --
| CREATE ALIAS FUNCTION dump_local(int) WITH PARAMETER(x) AS abs(`x`);
< ok 0 0
> exec
|
|
| --
| -- View structure for `calls_functions`
| --
| CREATE VIEW `calls_functions`
| (g,l)
|  AS SELECT `omni_dump_global_function`.`dump_global_called`(-1) AS `g`, `omni_dump_global_function`.`dump_local`(-2) AS `l`;
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_global_function
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_global_function
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_global_function") ORDER BY Name
< result
< column "Name" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_global_function
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< row
: bytes
| calls_functions
: bytes
| VIEW
> query
| SHOW FULL FUNCTIONS FROM `omni_dump_global_function`
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| dump_local(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_global_function
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE VIEW `omni_dump_global_function`.`calls_functions`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| calls_functions
: bytes
| CREATE VIEW `calls_functions`
| (g,l)
|  AS SELECT `omni_dump_global_function`.`dump_global_called`(-1) AS `g`, `omni_dump_global_function`.`dump_local`(-2) AS `l`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE FUNCTION dump_local(int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_local(int)
: bytes
| CREATE ALIAS FUNCTION dump_local(int) WITH PARAMETER(x) AS abs(`x`);
> query
| SHOW GLOBAL FULL FUNCTIONS
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| dump_global_called(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
< row
: bytes
| dump_global_called(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"greatest(`x`, `y`)"}
< row
: bytes
| dump_global_unused(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
> query
| SHOW CREATE GLOBAL FUNCTION dump_global_called(int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_global_called(int)
: bytes
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(int) WITH PARAMETER(x) AS abs(`x`);
> query
| SHOW CREATE GLOBAL FUNCTION dump_global_called(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_global_called(int, int)
: bytes
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(int, int) WITH PARAMETER(x, y) AS greatest(`x`, `y`);
> exec
| DROP DATABASE `omni_dump_global_function`
< ok 0 0
> exec
| CREATE DATABASE `omni_dump_global_function`
< ok 0 0
> exec
| USE `omni_dump_global_function`
< ok 0 0
> exec
| --
| -- Global function structure for dump_global_called(int)
| --
| CREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS dump_global_called(int) WITH PARAMETER(x) AS abs(`x`);
< ok 0 0
> exec
|
|
| --
| -- Global function structure for dump_global_called(int, int)
| --
| CREATE GLOBAL ALIAS FUNCTION IF NOT EXISTS dump_global_called(int, int) WITH PARAMETER(x, y) AS greatest(`x`, `y`);
< ok 0 0
> exec
|
|
| --
| -- Function structure for dump_local(int)
| --
| CREATE ALIAS FUNCTION dump_local(int) WITH PARAMETER(x) AS abs(`x`);
< ok 0 0
> exec
|
|
| --
| -- View structure for `calls_functions`
| --
| CREATE VIEW `calls_functions`
| (g,l)
|  AS SELECT `omni_dump_global_function`.`dump_global_called`(-1) AS `g`, `omni_dump_global_function`.`dump_local`(-2) AS `l`;
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_global_function
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_global_function
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_global_function") ORDER BY Name
< result
< column "Name" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_global_function
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< row
: bytes
| calls_functions
: bytes
| VIEW
> query
| SHOW FULL FUNCTIONS FROM `omni_dump_global_function`
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| dump_local(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_global_function
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE VIEW `omni_dump_global_function`.`calls_functions`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| calls_functions
: bytes
| CREATE VIEW `calls_functions`
| (g,l)
|  AS SELECT `omni_dump_global_function`.`dump_global_called`(-1) AS `g`, `omni_dump_global_function`.`dump_local`(-2) AS `l`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE FUNCTION dump_local(int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_local(int)
: bytes
| CREATE ALIAS FUNCTION dump_local(int) WITH PARAMETER(x) AS abs(`x`);
> query
| SHOW GLOBAL FULL FUNCTIONS
< result
< column "Signature" "CHAR"
< column "Return Type" "CHAR"
< column "Function Type" "CHAR"
< column "Intermediate Type" "CHAR"
< column "Properties" "CHAR"
< row
: bytes
| dump_global_called(int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x","origin_function":"abs(`x`)"}
< row
: bytes
| dump_global_called(int, int)
: bytes
| VARCHAR
: bytes
| Alias
: bytes
| NULL
: bytes
| {"parameter":"x, y","origin_function":"greatest(`x`, `y`)"}
> query
| SHOW CREATE GLOBAL FUNCTION dump_global_called(int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_global_called(int)
: bytes
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(int) WITH PARAMETER(x) AS abs(`x`);
> query
| SHOW CREATE GLOBAL FUNCTION dump_global_called(int, int)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_global_called(int, int)
: bytes
| CREATE GLOBAL ALIAS FUNCTION dump_global_called(int, int) WITH PARAMETER(x, y) AS greatest(`x`, `y`);
