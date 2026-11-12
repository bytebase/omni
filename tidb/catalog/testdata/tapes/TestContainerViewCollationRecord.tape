-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerViewCollationRecord default#1
> query
| SELECT @@collation_connection
< result
< column "@@collation_connection" "VARCHAR"
< row
: bytes
| utf8mb4_general_ci
> exec
| SET NAMES utf8mb4 COLLATE utf8mb4_bin
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_view_collation
< ok 0 0
> exec
| CREATE DATABASE omni_view_collation
< ok 0 0
> exec
| USE omni_view_collation
< ok 0 0
> exec
| CREATE TABLE t (
| 	a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
| 	g VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
| 	u VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
| 	z VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 	l VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_bin,
| 	s VARCHAR(10) CHARACTER SET ascii,
| 	vb VARBINARY(10),
| 	i INT,
| 	d DATE
| )
< ok 0 0
> exec
| INSERT INTO t VALUES ('x', 'x', 'x', 'x', 'x', 'x', 'x', 1, '2020-01-01')
< ok 1 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(a, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(g, u) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 1
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT g AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(z, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 1
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(g, u, a) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 1
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(CONCAT(g, u), l) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 1
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT REPEAT(CONCAT(g, u), 2) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 1
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT IF(1, g, u) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 1
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(a, l) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(l, a) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT IFNULL(g, l) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(s, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(s, l) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(l, 'é') AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT('é', l) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(l, _utf8mb4'x') AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(_latin1'x', g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(l, CONCAT('é', 'x')) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT LPAD(g, 5, 'x') AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(s, 'é') AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (ascii_bin,IMPLICIT) and (utf8mb4_bin,COERCIBLE) for operation 'concat'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(1, 2) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 4
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(g, 1) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(i, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(d, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(NULL, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(vb, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| binary
: bytes
| binary
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(g, vb) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| binary
: bytes
| binary
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(g = l, s) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| ascii
: bytes
| ascii_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT LOWER(g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TRIM(BOTH 'x' FROM g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT LPAD(l, 20, g) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT LPAD(g, 20, u) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 1
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(g, u COLLATE utf8mb4_general_ci) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: int64 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(l COLLATE latin1_bin, 'é') AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONVERT(g USING latin1) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(l, CONVERT(g USING latin1)) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| latin1
: bytes
| latin1_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(u, CONVERT(g USING utf8mb4)) AS c FROM t
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(g COLLATE utf8mb4_general_ci, u COLLATE utf8mb4_bin) AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT LENGTH(CONCAT(g COLLATE utf8mb4_general_ci, u COLLATE utf8mb4_bin)) AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT CONCAT(l COLLATE latin1_bin, g) AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (latin1_bin,EXPLICIT) and (utf8mb4_general_ci,IMPLICIT) for operation 'concat'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE a = g
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE l = g
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE a = 'x' OR i = g OR d = g OR vb = g
< ok 0 0
> query
| SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< column "COERCIBILITY(c)" "BIGINT"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: int64 2
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE g = u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE z = g
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_0900_ai_ci,IMPLICIT) and (utf8mb4_general_ci,IMPLICIT) for operation '='
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE g <> u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '!='
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE g <=> u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '<=>'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE g < u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '<'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE g NOT LIKE u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation 'like'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE STRCMP(g, u) = 0
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation 'strcmp'
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE CONCAT(g, u) = CONCAT(g, u)
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,NONE) for operation '='
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a AS c FROM t WHERE COALESCE(g, u, z) = 'x'
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,COERCIBLE) for operation '='
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a AS c FROM t JOIN t AS t2 ON t.g = t2.u
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='
> exec
| DROP DATABASE IF EXISTS omni_view_collation
< ok 0 0
> exec
| SET NAMES utf8mb4 COLLATE utf8mb4_general_ci
< ok 0 0
