-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_4_2_StoredRoutines default#1
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
| 1305 42000 procedure/function does not exist
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
| 1305 42000 procedure/function does not exist
> exec
| DROP FUNCTION IF EXISTS fn_nonexistent_xyz
< ok 0 0
> exec
| DROP FUNCTION fn_totally_nonexistent
< error
| 1305 42000 FUNCTION test.fn_totally_nonexistent does not exist
