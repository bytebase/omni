-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerGenerationExpressionsMatchTiDB default#1
> exec
| DROP DATABASE IF EXISTS omni_genexpr
< ok 0 0
> exec
| CREATE DATABASE omni_genexpr
< ok 0 0
> exec
| USE omni_genexpr
< ok 0 0
> exec
| CREATE TABLE gen_exprs (a int, b int, s varchar(20), j json,
|  g1 int AS (a * (b + 1)),
|  g2 int AS ((a + b) * 2),
|  g3 int AS (a + b * 2),
|  g4 varchar(30) AS (concat('it''s', s)),
|  g5 varchar(30) AS (CASE WHEN a < 5 THEN 'x' WHEN a IS NULL THEN 'n' ELSE 'y' END) STORED,
|  g6 int AS (a IN (1, 2, 3)),
|  g7 int AS (s LIKE 'a%'),
|  g8 int AS (a BETWEEN 1 AND 5),
|  g9 char(10) AS (CAST(a AS char)),
|  g10 int AS (IF(a > 0, a, -a)),
|  g11 varchar(20) AS (concat('a\\b', s)),
|  g12 int AS (NOT (a > 1)),
|  g13 int AS (a IS NOT NULL),
|  g14 varchar(40) AS (j->>'$.x'),
|  g15 int AS (-a + 1),
|  g16 int AS (a DIV 2 MOD 3),
|  g17 int AS (coalesce(a, b, 0)));
< ok 0 0
> query
| SHOW CREATE TABLE `gen_exprs`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| gen_exprs
: bytes
| CREATE TABLE `gen_exprs` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   `j` json DEFAULT NULL,
|   `g1` int GENERATED ALWAYS AS (`a` * (`b` + 1)) VIRTUAL,
|   `g2` int GENERATED ALWAYS AS ((`a` + `b`) * 2) VIRTUAL,
|   `g3` int GENERATED ALWAYS AS (`a` + `b` * 2) VIRTUAL,
|   `g4` varchar(30) GENERATED ALWAYS AS (concat(_utf8mb4'it''s', `s`)) VIRTUAL,
|   `g5` varchar(30) GENERATED ALWAYS AS (case when `a` < 5 then _utf8mb4'x' when `a` is null then _utf8mb4'n' else _utf8mb4'y' end) STORED,
|   `g6` int GENERATED ALWAYS AS (`a` in (1,2,3)) VIRTUAL,
|   `g7` int GENERATED ALWAYS AS (`s` like _utf8mb4'a%') VIRTUAL,
|   `g8` int GENERATED ALWAYS AS (`a` between 1 and 5) VIRTUAL,
|   `g9` char(10) GENERATED ALWAYS AS (cast(`a` as char)) VIRTUAL,
|   `g10` int GENERATED ALWAYS AS (if(`a` > 0, `a`, -`a`)) VIRTUAL,
|   `g11` varchar(20) GENERATED ALWAYS AS (concat(_utf8mb4'a\\b', `s`)) VIRTUAL,
|   `g12` int GENERATED ALWAYS AS (not (`a` > 1)) VIRTUAL,
|   `g13` int GENERATED ALWAYS AS (`a` is not null) VIRTUAL,
|   `g14` varchar(40) GENERATED ALWAYS AS (json_unquote(json_extract(`j`, _utf8mb4'$.x'))) VIRTUAL,
|   `g15` int GENERATED ALWAYS AS (-`a` + 1) VIRTUAL,
|   `g16` int GENERATED ALWAYS AS (`a` div 2 % 3) VIRTUAL,
|   `g17` int GENERATED ALWAYS AS (coalesce(`a`, `b`, 0)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| USE test
< ok 0 0
