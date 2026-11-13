-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_ViewMigrationApplyCorrectness default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_simple
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_simple
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_simple.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_simple.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_simple`.`t`.`a` AS `a` from `vgen_to_create_simple`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_simple
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_simple
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_simple.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 0 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_with_options
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_with_options DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_with_options
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_with_options.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_with_options.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_with_options`.`t`.`a` AS `a` from `vgen_to_create_with_options`.`t` where (`vgen_to_create_with_options`.`t`.`a` > 0) WITH CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_with_options
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_with_options DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_with_options
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_with_options.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t` where (`vt`.`t`.`a` > 0) WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t` where (`vt`.`t`.`a` > 0) WITH CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_stored_forms
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_stored_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_stored_forms
< ok 0 0
> exec
| CREATE TABLE t (a INT, s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v_cte AS WITH c AS (SELECT a FROM t) SELECT a FROM c
< ok 0 0
> exec
| CREATE VIEW v_union (p) AS SELECT a FROM t UNION SELECT a + 1 FROM t UNION ALL SELECT 1
< ok 0 0
> exec
| CREATE VIEW v_if AS SELECT IF(a > 0, 'YES', 'NO') AS x, TRIM('x' FROM s) AS y FROM t
< ok 0 0
> exec
| CREATE VIEW v_in AS SELECT a FROM t WHERE a IN (SELECT a FROM t) AND a NOT IN (SELECT a + 1 FROM t)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_stored_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_stored_forms.v_cte
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte
: bytes
| CREATE VIEW `v_cte` AS WITH RECURSIVE `c` as (select `vgen_to_create_stored_forms`.`t`.`a` AS `a` from `vgen_to_create_stored_forms`.`t`) select `c`.`a` AS `a` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_stored_forms.v_union
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_union
: bytes
| CREATE VIEW `v_union` AS ((select `vgen_to_create_stored_forms`.`t`.`a` AS `p` from `vgen_to_create_stored_forms`.`t`) union (select (`vgen_to_create_stored_forms`.`t`.`a` + 1) AS `p` from `vgen_to_create_stored_forms`.`t`)) union  all (select 1 AS `p`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_stored_forms.v_if
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_if
: bytes
| CREATE VIEW `v_if` AS select (case when (`vgen_to_create_stored_forms`.`t`.`a` > 0) then 'YES' else 'NO' end) AS `x`,trim(both 'x' from `vgen_to_create_stored_forms`.`t`.`s`) AS `y` from `vgen_to_create_stored_forms`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_stored_forms.v_in
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_in
: bytes
| CREATE VIEW `v_in` AS select `vgen_to_create_stored_forms`.`t`.`a` AS `a` from `vgen_to_create_stored_forms`.`t` where (`vgen_to_create_stored_forms`.`t`.`a` = any (select `vgen_to_create_stored_forms`.`t`.`a` from `vgen_to_create_stored_forms`.`t`)) and (`vgen_to_create_stored_forms`.`t`.`a` <> all (select (`vgen_to_create_stored_forms`.`t`.`a` + 1) from `vgen_to_create_stored_forms`.`t`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_stored_forms
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_stored_forms DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_stored_forms
< ok 0 0
> exec
| CREATE TABLE t (a INT, s VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_stored_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, s VARCHAR(20))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v_cte` AS with recursive `c` as (select `vt`.`t`.`a` AS `a` from `vt`.`t`) select `c`.`a` AS `a` from `c`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v_if` AS select (case when (`vt`.`t`.`a` > 0) then 'YES' else 'NO' end) AS `x`,trim(both 'x' from `vt`.`t`.`s`) AS `y` from `vt`.`t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v_in` AS select `vt`.`t`.`a` AS `a` from `vt`.`t` where ((`vt`.`t`.`a` = any (select `vt`.`t`.`a` from `vt`.`t`)) and (`vt`.`t`.`a` <> all (select (`vt`.`t`.`a` + 1) from `vt`.`t`)))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v_union` AS ((select `vt`.`t`.`a` AS `p` from `vt`.`t`) union (select (`vt`.`t`.`a` + 1) AS `p` from `vt`.`t`)) union all (select 1 AS `p`)
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v_cte
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cte
: bytes
| CREATE VIEW `v_cte` AS WITH RECURSIVE `c` as (select `vt`.`t`.`a` AS `a` from `vt`.`t`) select `c`.`a` AS `a` from `c`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v_union
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_union
: bytes
| CREATE VIEW `v_union` AS ((select `vt`.`t`.`a` AS `p` from `vt`.`t`) union (select (`vt`.`t`.`a` + 1) AS `p` from `vt`.`t`)) union  all (select 1 AS `p`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v_if
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_if
: bytes
| CREATE VIEW `v_if` AS select (case when (`vt`.`t`.`a` > 0) then 'YES' else 'NO' end) AS `x`,trim(both 'x' from `vt`.`t`.`s`) AS `y` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v_in
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_in
: bytes
| CREATE VIEW `v_in` AS select `vt`.`t`.`a` AS `a` from `vt`.`t` where (`vt`.`t`.`a` = any (select `vt`.`t`.`a` from `vt`.`t`)) and (`vt`.`t`.`a` <> all (select (`vt`.`t`.`a` + 1) from `vt`.`t`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_and_table
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_and_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_and_table
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_and_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_view_and_table.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_view_and_table`.`t`.`a` AS `a` from `vgen_to_create_view_and_table`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 5 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE `vt`.`t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a, b FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_view_on_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE VIEW `base` AS select `vgen_to_create_view_on_view`.`t`.`a` AS `a`,`vgen_to_create_view_on_view`.`t`.`b` AS `b` from `vgen_to_create_view_on_view`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_on_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_view_on_view`.`base`.`a` AS `a` from `vgen_to_create_view_on_view`.`base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`base` AS select `vt`.`t`.`a` AS `a`,`vt`.`t`.`b` AS `b` from `vt`.`t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`base`.`a` AS `a` from `vt`.`base`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.base
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE VIEW `base` AS select `vt`.`t`.`a` AS `a`,`vt`.`t`.`b` AS `b` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`base`.`a` AS `a` from `vt`.`base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_chain
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_chain DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_chain
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v2 AS SELECT a FROM v1
< ok 0 0
> exec
| CREATE VIEW v3 AS SELECT a FROM v2
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_chain.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_view_chain.v1
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE VIEW `v1` AS select `vgen_to_create_view_chain`.`t`.`a` AS `a` from `vgen_to_create_view_chain`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_chain.v2
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE VIEW `v2` AS select `vgen_to_create_view_chain`.`v1`.`a` AS `a` from `vgen_to_create_view_chain`.`v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_chain.v3
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v3
: bytes
| CREATE VIEW `v3` AS select `vgen_to_create_view_chain`.`v2`.`a` AS `a` from `vgen_to_create_view_chain`.`v2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_view_chain
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_view_chain DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_view_chain
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_view_chain.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 3 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v1` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v2` AS select `vt`.`v1`.`a` AS `a` from `vt`.`v1`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v3` AS select `vt`.`v2`.`a` AS `a` from `vt`.`v2`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v1
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE VIEW `v1` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v2
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE VIEW `v2` AS select `vt`.`v1`.`a` AS `a` from `vt`.`v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v3
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v3
: bytes
| CREATE VIEW `v3` AS select `vt`.`v2`.`a` AS `a` from `vt`.`v2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_alias_name_clash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| CREATE VIEW a AS SELECT id AS b FROM t
< ok 0 0
> exec
| CREATE VIEW b AS SELECT b FROM a
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_alias_name_clash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_view_alias_name_clash.a
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `vgen_to_create_view_alias_name_clash`.`t`.`id` AS `b` from `vgen_to_create_view_alias_name_clash`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_alias_name_clash.b
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE VIEW `b` AS select `vgen_to_create_view_alias_name_clash`.`a`.`b` AS `b` from `vgen_to_create_view_alias_name_clash`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_view_alias_name_clash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_view_alias_name_clash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`a` AS select `vt`.`t`.`id` AS `b` from `vt`.`t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`b` AS select `vt`.`a`.`b` AS `b` from `vt`.`a`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.a
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `vt`.`t`.`id` AS `b` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.b
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE VIEW `b` AS select `vt`.`a`.`b` AS `b` from `vt`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y, w.n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_create_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_create_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_multi_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_multi_join`.`t`.`a` AS `a`,`vgen_to_create_multi_join`.`u`.`y` AS `y`,`vgen_to_create_multi_join`.`w`.`n` AS `n` from ((`vgen_to_create_multi_join`.`t` join `vgen_to_create_multi_join`.`u` on (`vgen_to_create_multi_join`.`t`.`c` = `vgen_to_create_multi_join`.`u`.`x`)) join `vgen_to_create_multi_join`.`w` on (`vgen_to_create_multi_join`.`u`.`y` = `vgen_to_create_multi_join`.`w`.`m`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_create_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_create_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 3 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`a` AS `a`,`vt`.`u`.`y` AS `y`,`vt`.`w`.`n` AS `n` from ((`vt`.`t` join `vt`.`u` on((`vt`.`t`.`c` = `vt`.`u`.`x`))) join `vt`.`w` on((`vt`.`u`.`y` = `vt`.`w`.`m`)))
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a`,`vt`.`u`.`y` AS `y`,`vt`.`w`.`n` AS `n` from ((`vt`.`t` join `vt`.`u` on (`vt`.`t`.`c` = `vt`.`u`.`x`)) join `vt`.`w` on (`vt`.`u`.`y` = `vt`.`w`.`m`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_sakila_film_list
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_sakila_film_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_sakila_film_list
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT film.film_id AS FID, category.name AS category, actor.name AS actor FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id JOIN film_actor ON film.film_id = film_actor.film_id JOIN actor ON film_actor.actor_id = actor.actor_id GROUP BY film.film_id
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int(11) DEFAULT NULL,
|   `category_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int(11) DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `film_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_sakila_film_list.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_sakila_film_list`.`film`.`film_id` AS `FID`,`vgen_to_create_sakila_film_list`.`category`.`name` AS `category`,`vgen_to_create_sakila_film_list`.`actor`.`name` AS `actor` from ((((`vgen_to_create_sakila_film_list`.`category` left join `vgen_to_create_sakila_film_list`.`film_category` on (`vgen_to_create_sakila_film_list`.`category`.`category_id` = `vgen_to_create_sakila_film_list`.`film_category`.`category_id`)) left join `vgen_to_create_sakila_film_list`.`film` on (`vgen_to_create_sakila_film_list`.`film_category`.`film_id` = `vgen_to_create_sakila_film_list`.`film`.`film_id`)) join `vgen_to_create_sakila_film_list`.`film_actor` on (`vgen_to_create_sakila_film_list`.`film`.`film_id` = `vgen_to_create_sakila_film_list`.`film_actor`.`film_id`)) join `vgen_to_create_sakila_film_list`.`actor` on (`vgen_to_create_sakila_film_list`.`film_actor`.`actor_id` = `vgen_to_create_sakila_film_list`.`actor`.`actor_id`)) group by `vgen_to_create_sakila_film_list`.`film`.`film_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_sakila_film_list
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_sakila_film_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_sakila_film_list
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int(11) DEFAULT NULL,
|   `category_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int(11) DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `film_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`film`.`film_id` AS `FID`,`vt`.`category`.`name` AS `category`,`vt`.`actor`.`name` AS `actor` from ((((`vt`.`category` left join `vt`.`film_category` on((`vt`.`category`.`category_id` = `vt`.`film_category`.`category_id`))) left join `vt`.`film` on((`vt`.`film_category`.`film_id` = `vt`.`film`.`film_id`))) join `vt`.`film_actor` on((`vt`.`film`.`film_id` = `vt`.`film_actor`.`film_id`))) join `vt`.`actor` on((`vt`.`film_actor`.`actor_id` = `vt`.`actor`.`actor_id`))) group by `vt`.`film`.`film_id`
< ok 0 0
> query
| SHOW CREATE TABLE vt.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int(11) DEFAULT NULL,
|   `category_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int(11) DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `film_id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int(11) DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`film`.`film_id` AS `FID`,`vt`.`category`.`name` AS `category`,`vt`.`actor`.`name` AS `actor` from ((((`vt`.`category` left join `vt`.`film_category` on (`vt`.`category`.`category_id` = `vt`.`film_category`.`category_id`)) left join `vt`.`film` on (`vt`.`film_category`.`film_id` = `vt`.`film`.`film_id`)) join `vt`.`film_actor` on (`vt`.`film`.`film_id` = `vt`.`film_actor`.`film_id`)) join `vt`.`actor` on (`vt`.`film_actor`.`actor_id` = `vt`.`actor`.`actor_id`)) group by `vt`.`film`.`film_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_employees_current_dept_emp DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> exec
| CREATE VIEW current_dept_emp AS SELECT l.emp_no, d.dept_no, l.from_date, l.to_date FROM dept_emp d JOIN dept_emp_latest_date l ON d.emp_no = l.emp_no AND d.from_date = l.from_date JOIN departments dp ON d.dept_no = dp.dept_no
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_employees_current_dept_emp.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int(11) NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`, `dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_create_employees_current_dept_emp.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_employees_current_dept_emp.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE VIEW `dept_emp_latest_date` AS select `vgen_to_create_employees_current_dept_emp`.`dept_emp`.`emp_no` AS `emp_no`,max(`vgen_to_create_employees_current_dept_emp`.`dept_emp`.`from_date`) AS `from_date`,max(`vgen_to_create_employees_current_dept_emp`.`dept_emp`.`to_date`) AS `to_date` from `vgen_to_create_employees_current_dept_emp`.`dept_emp` group by `vgen_to_create_employees_current_dept_emp`.`dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_employees_current_dept_emp.current_dept_emp
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| current_dept_emp
: bytes
| CREATE VIEW `current_dept_emp` AS select `vgen_to_create_employees_current_dept_emp`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`vgen_to_create_employees_current_dept_emp`.`l`.`from_date` AS `from_date`,`vgen_to_create_employees_current_dept_emp`.`l`.`to_date` AS `to_date` from ((`vgen_to_create_employees_current_dept_emp`.`dept_emp` `d` join `vgen_to_create_employees_current_dept_emp`.`dept_emp_latest_date` `l` on (`d`.`emp_no` = `vgen_to_create_employees_current_dept_emp`.`l`.`emp_no`) and (`d`.`from_date` = `vgen_to_create_employees_current_dept_emp`.`l`.`from_date`)) join `vgen_to_create_employees_current_dept_emp`.`departments` `dp` on (`d`.`dept_no` = `dp`.`dept_no`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_employees_current_dept_emp DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_employees_current_dept_emp.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int(11) NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`, `dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_create_employees_current_dept_emp.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_create_employees_current_dept_emp.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE VIEW `dept_emp_latest_date` AS select `vgen_from_create_employees_current_dept_emp`.`dept_emp`.`emp_no` AS `emp_no`,max(`vgen_from_create_employees_current_dept_emp`.`dept_emp`.`from_date`) AS `from_date`,max(`vgen_from_create_employees_current_dept_emp`.`dept_emp`.`to_date`) AS `to_date` from `vgen_from_create_employees_current_dept_emp`.`dept_emp` group by `vgen_from_create_employees_current_dept_emp`.`dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 6 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`current_dept_emp` AS select `vt`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`vt`.`l`.`from_date` AS `from_date`,`vt`.`l`.`to_date` AS `to_date` from ((`vt`.`dept_emp` `d` join `vt`.`dept_emp_latest_date` `l` on(((`d`.`emp_no` = `vt`.`l`.`emp_no`) and (`d`.`from_date` = `vt`.`l`.`from_date`)))) join `vt`.`departments` `dp` on((`d`.`dept_no` = `dp`.`dept_no`)))
< ok 0 0
> query
| SHOW CREATE TABLE vt.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int(11) NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`, `dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE VIEW `dept_emp_latest_date` AS select `vt`.`dept_emp`.`emp_no` AS `emp_no`,max(`vt`.`dept_emp`.`from_date`) AS `from_date`,max(`vt`.`dept_emp`.`to_date`) AS `to_date` from `vt`.`dept_emp` group by `vt`.`dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.current_dept_emp
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| current_dept_emp
: bytes
| CREATE VIEW `current_dept_emp` AS select `vt`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`vt`.`l`.`from_date` AS `from_date`,`vt`.`l`.`to_date` AS `to_date` from ((`vt`.`dept_emp` `d` join `vt`.`dept_emp_latest_date` `l` on (`d`.`emp_no` = `vt`.`l`.`emp_no`) and (`d`.`from_date` = `vt`.`l`.`from_date`)) join `vt`.`departments` `dp` on (`d`.`dept_no` = `dp`.`dept_no`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_into_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_into_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_into_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, w.n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_into_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_replace_into_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_to_replace_into_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_replace_into_multi_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_replace_into_multi_join`.`t`.`a` AS `a`,`vgen_to_replace_into_multi_join`.`w`.`n` AS `n` from ((`vgen_to_replace_into_multi_join`.`t` join `vgen_to_replace_into_multi_join`.`u` on (`vgen_to_replace_into_multi_join`.`t`.`c` = `vgen_to_replace_into_multi_join`.`u`.`x`)) join `vgen_to_replace_into_multi_join`.`w` on (`vgen_to_replace_into_multi_join`.`u`.`y` = `vgen_to_replace_into_multi_join`.`w`.`m`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_into_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_into_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_into_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_into_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_replace_into_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vgen_from_replace_into_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_replace_into_multi_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_replace_into_multi_join`.`t`.`a` AS `a` from `vgen_from_replace_into_multi_join`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`a` AS `a`,`vt`.`w`.`n` AS `n` from ((`vt`.`t` join `vt`.`u` on((`vt`.`t`.`c` = `vt`.`u`.`x`))) join `vt`.`w` on((`vt`.`u`.`y` = `vt`.`w`.`m`)))
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int(11) DEFAULT NULL,
|   `y` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE vt.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int(11) DEFAULT NULL,
|   `n` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a`,`vt`.`w`.`n` AS `n` from ((`vt`.`t` join `vt`.`u` on (`vt`.`t`.`c` = `vt`.`u`.`x`)) join `vt`.`w` on (`vt`.`u`.`y` = `vt`.`w`.`m`))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_paren_subquery_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT if(((SELECT count(0) FROM t) = 0),'NO','YES') AS x
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_paren_subquery_compare.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_paren_subquery_compare.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (case when ((select count(0) from `vgen_to_create_paren_subquery_compare`.`t`) = 0) then 'NO' else 'YES' end) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_paren_subquery_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_paren_subquery_compare.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select (case when ((select count(0) from `vt`.`t`) = 0) then 'NO' else 'YES' end) AS `x`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select (case when ((select count(0) from `vt`.`t`) = 0) then 'NO' else 'YES' end) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_body
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_body
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_body.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_replace_body.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_replace_body`.`t`.`b` AS `b` from `vgen_to_replace_body`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_body
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_body
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_body.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_replace_body.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_replace_body`.`t`.`a` AS `a` from `vgen_from_replace_body`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`b` AS `b` from `vt`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`b` AS `b` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_algorithm
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_algorithm DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_algorithm
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_algorithm.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_replace_algorithm.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_replace_algorithm`.`t`.`a` AS `a` from `vgen_to_replace_algorithm`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_algorithm
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_algorithm DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_algorithm
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_algorithm.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_replace_algorithm.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_replace_algorithm`.`t`.`a` AS `a` from `vgen_from_replace_algorithm`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_sql_security
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_sql_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_sql_security
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_sql_security.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_replace_sql_security.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_replace_sql_security`.`t`.`a` AS `a` from `vgen_to_replace_sql_security`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_sql_security
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_sql_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_sql_security
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_sql_security.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_replace_sql_security.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_replace_sql_security`.`t`.`a` AS `a` from `vgen_from_replace_sql_security`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_add_check_option
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_add_check_option DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_add_check_option
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_add_check_option.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_replace_add_check_option.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_replace_add_check_option`.`t`.`a` AS `a` from `vgen_to_replace_add_check_option`.`t` where (`vgen_to_replace_add_check_option`.`t`.`a` > 0) WITH CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_add_check_option
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_add_check_option DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_add_check_option
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_add_check_option.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_replace_add_check_option.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_replace_add_check_option`.`t`.`a` AS `a` from `vgen_from_replace_add_check_option`.`t` where (`vgen_from_replace_add_check_option`.`t`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t` where (`vt`.`t`.`a` > 0) WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `a` from `vt`.`t` where (`vt`.`t`.`a` > 0) WITH CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_drop_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_drop_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_drop_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_drop_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vgen_from_drop_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_drop_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_drop_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_drop_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_drop_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_drop_view`.`t`.`a` AS `a` from `vgen_from_drop_view`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| DROP VIEW IF EXISTS `vt`.`v`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< error
| 1146 42S02 Table 'vt.v' doesn't exist
> exec
| DROP DATABASE IF EXISTS vgen_to_drop_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_drop_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_drop_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_drop_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vgen_from_drop_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_drop_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_drop_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_drop_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_drop_view_on_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE VIEW `base` AS select `vgen_from_drop_view_on_view`.`t`.`a` AS `a` from `vgen_from_drop_view_on_view`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_from_drop_view_on_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_drop_view_on_view`.`base`.`a` AS `a` from `vgen_from_drop_view_on_view`.`base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 1 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> exec
| DROP VIEW IF EXISTS `vt`.`base`
< ok 0 0
> exec
| DROP VIEW IF EXISTS `vt`.`v`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.base
< error
| 1146 42S02 Table 'vt.base' doesn't exist
> query
| SHOW CREATE VIEW vt.v
< error
| 1146 42S02 Table 'vt.v' doesn't exist
> exec
| DROP DATABASE IF EXISTS vgen_to_create_explicit_columns
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_explicit_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_explicit_columns
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v (p, q) AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_explicit_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_explicit_columns.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_explicit_columns`.`t`.`a` AS `p`,`vgen_to_create_explicit_columns`.`t`.`b` AS `q` from `vgen_to_create_explicit_columns`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_explicit_columns
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_explicit_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_explicit_columns
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_explicit_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 1 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`a` AS `p`,`vt`.`t`.`b` AS `q` from `vt`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`a` AS `p`,`vt`.`t`.`b` AS `q` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_create_tsdiff_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_tsdiff_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_tsdiff_view
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id, TIMESTAMPDIFF(SECOND, created, NOW()) AS age_secs FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_tsdiff_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_create_tsdiff_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_to_create_tsdiff_view`.`t`.`id` AS `id`,timestampdiff(second,`vgen_to_create_tsdiff_view`.`t`.`created`,now()) AS `age_secs` from `vgen_to_create_tsdiff_view`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_create_tsdiff_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_tsdiff_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_tsdiff_view
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_tsdiff_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `vt`.`t`.`id` AS `id`,timestampdiff(SECOND,`vt`.`t`.`created`,now()) AS `age_secs` from `vt`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vt`.`t`.`id` AS `id`,timestampdiff(second,`vt`.`t`.`created`,now()) AS `age_secs` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_into_temporal_medley DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TIMESTAMPDIFF(DAY, a, NOW()) AS c1, TIMESTAMPADD(HOUR, 1, a) AS c2, EXTRACT(DAY_HOUR FROM a) AS c3, GET_FORMAT(DATE, 'USA') AS c4, a + INTERVAL '1:1' DAY_HOUR AS c5, TRIM('x' FROM s) AS c6 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_into_temporal_medley.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_to_replace_into_temporal_medley.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select timestampdiff(day,`vgen_to_replace_into_temporal_medley`.`t`.`a`,now()) AS `c1`,timestampadd(HOUR ,1,`vgen_to_replace_into_temporal_medley`.`t`.`a`) AS `c2`,extract(day_hour from `vgen_to_replace_into_temporal_medley`.`t`.`a`) AS `c3`,get_format(date, 'USA') AS `c4`,(`vgen_to_replace_into_temporal_medley`.`t`.`a` + interval '1:1' day_hour) AS `c5`,trim(both 'x' from `vgen_to_replace_into_temporal_medley`.`t`.`s`) AS `c6` from `vgen_to_replace_into_temporal_medley`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_into_temporal_medley DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_into_temporal_medley.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vgen_from_replace_into_temporal_medley.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select `vgen_from_replace_into_temporal_medley`.`t`.`a` AS `a` from `vgen_from_replace_into_temporal_medley`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select timestampdiff(DAY,`vt`.`t`.`a`,now()) AS `c1`,(`vt`.`t`.`a` + interval 1 hour) AS `c2`,extract(day_hour from `vt`.`t`.`a`) AS `c3`,get_format(DATE, 'USA') AS `c4`,(`vt`.`t`.`a` + interval '1:1' day_hour) AS `c5`,trim(both 'x' from `vt`.`t`.`s`) AS `c6` from `vt`.`t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select timestampdiff(day,`vt`.`t`.`a`,now()) AS `c1`,(`vt`.`t`.`a` + interval 1 hour) AS `c2`,extract(day_hour from `vt`.`t`.`a`) AS `c3`,get_format(date, 'USA') AS `c4`,(`vt`.`t`.`a` + interval '1:1' day_hour) AS `c5`,trim(both 'x' from `vt`.`t`.`s`) AS `c6` from `vt`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
