-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_4_2_StoredRoutines default#1
> exec
| CREATE PROCEDURE sp_greet(IN name VARCHAR(100)) BEGIN SELECT CONCAT('Hello, ', name); END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec
| CREATE PROCEDURE sp_noop() BEGIN SELECT 1; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec
| DROP PROCEDURE sp_greet
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec
| DROP PROCEDURE IF EXISTS sp_absent
< error
| 8108 HY000 Unsupported type *resolve.NodeW
