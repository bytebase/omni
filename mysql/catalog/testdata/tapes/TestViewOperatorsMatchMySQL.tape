-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestViewOperatorsMatchMySQL default#1
> exec
| CREATE TABLE tt (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE OR REPLACE VIEW v AS SELECT a || b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((0 <> `tt`.`a`) or (0 <> `tt`.`b`)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a && b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((0 <> `tt`.`a`) and (0 <> `tt`.`b`)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a || b && c AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((0 <> `tt`.`a`) or ((0 <> `tt`.`b`) and (0 <> `tt`.`c`))) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a = b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` <> `tt`.`b`) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a < b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` >= `tt`.`b`) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a = b is null AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((`tt`.`a` = `tt`.`b`) is null) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a is null = b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((`tt`.`a` is null) = `tt`.`b`) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a between 1 and 2 | 3 AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` between 1 and (2 | 3)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a = any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` in (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a = some (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` in (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a = all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` = all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a <> any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` <> any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a <> some (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` <> any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a <> all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (not(`tt`.`a` in (select 1))) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a > any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` > any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a > some (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` > any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a > all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` > all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT a <= all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` <= all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a = any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (not(`tt`.`a` in (select 1))) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a > all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` <= any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a > any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` <= all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a = all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` <> any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a <> any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` = all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a <> all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` in (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a >= any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` < all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a >= all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` < any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a < any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` >= all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a < all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` >= any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a <= any (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` > all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a <= all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` > any (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT not a in (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (not(`tt`.`a` in (select 1))) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a > all (select 1)) = 1 AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` > all (select 1) = 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a > any (select 1)) and b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` > any (select 1) and (0 <> `tt`.`b`)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT b or (a > all (select 1)) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((0 <> `tt`.`b`) or `tt`.`a` > all (select 1)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a = any (select 1)) and b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` in (select 1) and (0 <> `tt`.`b`)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a <> all (select 1)) and b AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((not(`tt`.`a` in (select 1))) and (0 <> `tt`.`b`)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a > all (select 1)) is null AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` > all (select 1) is null) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a = any (select 1)) is null AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (`tt`.`a` in (select 1) is null) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 = (a = any (select 1)) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (1 = `tt`.`a` in (select 1)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT 1 = (a <> all (select 1)) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (1 = (not(`tt`.`a` in (select 1)))) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a <> all (select 1)) + 1 AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select ((not(`tt`.`a` in (select 1))) + 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a > all (select 1)) > all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `tt`.`a` > all (select 1) > all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT (a <> all (select 1)) > all (select 1) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (not(`tt`.`a` in (select 1))) > all (select 1) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT -(a > all (select 1)) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select -(`tt`.`a` > all (select 1)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT abs(a > all (select 1)) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select abs(`tt`.`a` > all (select 1)) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT coalesce(a = any (select 1), a <> all (select 1)) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select coalesce(`tt`.`a` in (select 1),(not(`tt`.`a` in (select 1)))) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT cast(a > all (select 1) as signed) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select cast(`tt`.`a` > all (select 1) as signed) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT if(a > all (select 1), 1, 2) AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select if(`tt`.`a` > all (select 1),1,2) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT case a > all (select 1) when 1 then 2 end AS x FROM tt
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select (case `tt`.`a` > all (select 1) when 1 then 2 end) AS `x` from `tt`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
