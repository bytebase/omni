-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_2_10_CreateDropView default#1
> exec
| DROP VIEW IF EXISTS v_basic
< ok 0 0
> exec
| CREATE VIEW v_basic AS SELECT 1 AS a
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_replace
< ok 0 0
> exec
| CREATE VIEW v_replace AS SELECT 1 AS a
< ok 0 0
> exec
| CREATE OR REPLACE VIEW v_replace AS SELECT 2 AS b
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_cols
< ok 0 0
> exec
| CREATE VIEW v_cols (x, y) AS SELECT 1, 2
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_opts
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW v_opts AS SELECT 1 AS a
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_chk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_chk_view
< ok 0 0
> exec
| CREATE TABLE t_chk_view (id INT, val INT)
< ok 0 0
> exec
| CREATE VIEW v_chk AS SELECT * FROM t_chk_view WHERE val > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_drop1
< ok 0 0
> exec
| CREATE VIEW v_drop1 AS SELECT 1 AS a
< ok 0 0
> exec
| DROP VIEW v_drop1
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_noexist
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_noexist
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_m1
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_m2
< ok 0 0
> exec
| CREATE VIEW v_m1 AS SELECT 1 AS a
< ok 0 0
> exec
| CREATE VIEW v_m2 AS SELECT 2 AS b
< ok 0 0
> exec
| DROP VIEW v_m1, v_m2
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_dup
< ok 0 0
> exec
| CREATE VIEW v_dup AS SELECT 1 AS a
< ok 0 0
> exec
| CREATE VIEW v_dup AS SELECT 2 AS b
< error
| 1050 42S01 Table 'test.v_dup' already exists
> exec
| DROP VIEW IF EXISTS v_tbl_conflict
< ok 0 0
> exec
| DROP TABLE IF EXISTS v_tbl_conflict
< ok 0 0
> exec
| CREATE TABLE v_tbl_conflict (id INT)
< ok 0 0
> exec
| CREATE VIEW v_tbl_conflict AS SELECT 1 AS a
< error
| 1050 42S01 Table 'test.v_tbl_conflict' already exists
> exec
| DROP VIEW IF EXISTS v_nonexist_err
< ok 0 0
> exec
| DROP VIEW v_nonexist_err
< error
| 1051 42S02 Unknown table 'test.v_nonexist_err'
