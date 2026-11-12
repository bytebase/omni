-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerViewCollationScopes default#1
> exec
| DROP DATABASE IF EXISTS c17src
< ok 0 0
> exec
| DROP DATABASE IF EXISTS c17dst
< ok 0 0
> exec
| CREATE DATABASE c17src; CREATE DATABASE c17dst;
| CREATE TABLE c17src.t (g VARCHAR(10) COLLATE utf8mb4_general_ci, u VARCHAR(10) COLLATE utf8mb4_general_ci);
| CREATE TABLE c17dst.t (g VARCHAR(10) COLLATE utf8mb4_general_ci, u VARCHAR(10) COLLATE utf8mb4_unicode_ci);
| USE c17src;
< ok 0 0
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT g FROM t WHERE g = u
< ok 0 0
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT g FROM c17dst.t WHERE g = u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS WITH t AS (SELECT 'x' AS g, 'x' AS u) SELECT g FROM c17dst.t AS t2 WHERE g = 'x' UNION ALL SELECT g FROM t WHERE g = u
< ok 0 0
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS WITH t AS (SELECT 'x' AS g, 'x' AS u) (SELECT g FROM t WHERE g = u)
< ok 0 0
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT g FROM c17dst.t WHERE g = u UNION ALL SELECT 'x'
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT c17src.t.g FROM c17src.t JOIN c17dst.t ON c17src.t.g = c17src.t.u
< ok 0 0
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT c17src.t.g FROM c17src.t JOIN c17dst.t ON c17dst.t.g = c17dst.t.u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='
> exec
| INSERT INTO c17src.t VALUES ('x', 'x')
< ok 1 0
> exec
| INSERT INTO c17dst.t VALUES ('x', 'x')
< ok 1 0
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT u FROM t
< ok 0 0
> query
| SELECT COLLATION(u) FROM c17dst.v
< result
< column "COLLATION(u)" "VARCHAR"
< row
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT c17dst.t.u FROM c17src.t JOIN c17dst.t
< ok 0 0
> query
| SELECT COLLATION(u) FROM c17dst.v
< result
< column "COLLATION(u)" "VARCHAR"
< row
: bytes
| utf8mb4_unicode_ci
> exec
| DROP VIEW IF EXISTS c17dst.v
< ok 0 0
> exec
| CREATE VIEW c17dst.v AS SELECT c17src.t.u FROM c17src.t JOIN c17dst.t
< ok 0 0
> query
| SELECT COLLATION(u) FROM c17dst.v
< result
< column "COLLATION(u)" "VARCHAR"
< row
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS c17src
< ok 0 0
> exec
| DROP DATABASE IF EXISTS c17dst
< ok 0 0
