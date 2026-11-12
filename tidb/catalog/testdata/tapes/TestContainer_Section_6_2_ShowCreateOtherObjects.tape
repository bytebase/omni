-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_6_2_ShowCreateOtherObjects default#1
> exec
| DROP VIEW IF EXISTS v_sc_basic
< ok 0 0
> exec
| CREATE VIEW v_sc_basic AS SELECT 1 AS col1, 2 AS col2
< ok 0 0
> query
| SHOW CREATE VIEW v_sc_basic
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_sc_basic
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_sc_basic` (`col1`, `col2`) AS SELECT 1 AS `col1`,2 AS `col2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP PROCEDURE IF EXISTS sp_sc_swap
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec
| CREATE PROCEDURE sp_sc_swap(INOUT a INT, INOUT b INT) BEGIN DECLARE tmp INT; SET tmp = a; SET a = b; SET b = tmp; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
