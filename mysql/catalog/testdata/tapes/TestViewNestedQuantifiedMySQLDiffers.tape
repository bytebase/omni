-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestViewNestedQuantifiedMySQLDiffers default#1
> exec
| CREATE TABLE tt (a INT, b INT, c INT)
< ok 0 0
> exec
| INSERT INTO tt VALUES (2, 0, 0)
< ok 1 0
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a > all (select 1)) + 1 AS x FROM tt
< ok 0 0
> query
| SHOW CREATE VIEW v
< error
| 1064 42000 …near '+ 1) AS `x` from `test`.`tt`' at line 1
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 = (a > all (select 1)) AS x FROM tt
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (1 = `tt`.`a`) > all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
