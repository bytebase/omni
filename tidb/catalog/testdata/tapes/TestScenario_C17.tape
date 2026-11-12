-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C17 default#1
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
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT CONCAT(a, b) AS c FROM t
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES ('x', 'y')
< ok 1 0
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v1 LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v1
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT COERCIBILITY(c) FROM testdb.v1 LIMIT 1
< result
< column "COERCIBILITY(c)" "BIGINT"
< row
: int64 2
> query
| SELECT CHARACTER_MAXIMUM_LENGTH FROM information_schema.COLUMNS
| 			WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v1' AND COLUMN_NAME='c'
< result
< column "CHARACTER_MAXIMUM_LENGTH" "BIGINT"
< row
: int64 20
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
| 			b VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_bin
| 		)
< ok 0 0
> exec
| CREATE VIEW v2 AS SELECT CONCAT(a, b) AS c, CONCAT(b, a) AS d FROM t
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES ('x', 'y')
< ok 1 0
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v2 LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v2
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT CHARSET(d), COLLATION(d) FROM testdb.v2 LIMIT 1
< result
< column "CHARSET(d)" "VARCHAR"
< column "COLLATION(d)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v2
: string
| d
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
| 			e VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
| 		)
< ok 0 0
> exec
| CREATE VIEW v3 AS SELECT CONCAT(a, b) AS c, CONCAT(a, e) AS d FROM t WHERE a = e
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES ('x', 'y', 'x')
< ok 1 0
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v3 LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v3
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT COERCIBILITY(c) FROM testdb.v3 LIMIT 1
< result
< column "COERCIBILITY(c)" "BIGINT"
< row
: int64 1
> query
| SELECT CHARSET(d), COLLATION(d) FROM testdb.v3 LIMIT 1
< result
< column "CHARSET(d)" "VARCHAR"
< column "COLLATION(d)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v3
: string
| d
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT COERCIBILITY(d) FROM testdb.v3 LIMIT 1
< result
< column "COERCIBILITY(d)" "BIGINT"
< row
: int64 2
> exec
| CREATE VIEW v3a AS SELECT a FROM t WHERE a = b
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='
> exec
| CREATE VIEW v3b AS SELECT CONCAT(a, b) = 'x' AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,COERCIBLE) for operation '='
> exec
| CREATE VIEW v3c AS SELECT a FROM t WHERE a LIKE b
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation 'like'
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v4 AS SELECT CONCAT_WS(',', a, b, NULL) AS c FROM t
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES (NULL, 'x')
< ok 1 0
> query
| SELECT CONCAT(a,b), CONCAT_WS(',',a,b) FROM testdb.t LIMIT 1
< result
< column "CONCAT(a,b)" "VARCHAR"
< column "CONCAT_WS(',',a,b)" "VARCHAR"
< row
: null
: bytes
| x
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v4 LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v4
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_bin,
| 			s VARCHAR(10) CHARACTER SET ascii
| 		)
< ok 0 0
> exec
| CREATE VIEW v5a AS SELECT CONCAT(a, 'x') AS c FROM t
< ok 0 0
> exec
| CREATE VIEW v5b AS SELECT CONCAT(a, _utf8mb4'x') AS c FROM t
< ok 0 0
> exec
| CREATE VIEW v5c AS SELECT CONCAT('x', s) AS c FROM t
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES ('x', 'y')
< ok 1 0
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v5a LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v5a
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v5b LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v5b
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v5c LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| ascii
: bytes
| ascii_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v5c
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> exec
| CREATE VIEW v5d AS SELECT CONCAT(s, 'é') AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (ascii_bin,IMPLICIT) and (utf8mb4_bin,COERCIBLE) for operation 'concat'
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_bin,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
| 			u VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v6a AS SELECT REPEAT(a, 3) AS c FROM t
< ok 0 0
> exec
| CREATE VIEW v6b AS SELECT LPAD(a, 20, b) AS c FROM t
< ok 0 0
> exec
| CREATE VIEW v6c AS SELECT RPAD(b, 20, a) AS c FROM t
< ok 0 0
> exec
| CREATE VIEW v6d AS SELECT LPAD(b, 20, u) AS c FROM t
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES ('x', 'y', 'z')
< ok 1 0
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v6a LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v6a
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v6b LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v6b
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v6c LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v6c
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v6d LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v6d
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT COERCIBILITY(c) FROM testdb.v6d LIMIT 1
< result
< column "COERCIBILITY(c)" "BIGINT"
< row
: int64 1
> exec
| CREATE VIEW v6e AS SELECT a FROM t WHERE LPAD(b, 20, u) = 'x'
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,COERCIBLE) for operation '='
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v7 AS
| 			SELECT CONCAT(a, CONVERT(b USING utf8mb4)) AS c, CONVERT(a USING latin1) AS d
| 			FROM t WHERE a = CONVERT(b USING utf8mb4)
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES ('x', 'x')
< ok 1 0
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v7 LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v7
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT CHARSET(d), COLLATION(d) FROM testdb.v7 LIMIT 1
< result
< column "CHARSET(d)" "VARCHAR"
< column "COLLATION(d)" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v7
: string
| d
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT COERCIBILITY(d) FROM testdb.v7 LIMIT 1
< result
< column "COERCIBILITY(d)" "BIGINT"
< row
: int64 2
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
| 			l VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_bin
| 		)
< ok 0 0
> exec
| CREATE VIEW v8a AS SELECT CONCAT(a, b COLLATE utf8mb4_bin) AS c FROM t
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES ('x', 'y', 'z')
< ok 1 0
> query
| SELECT CHARSET(c), COLLATION(c) FROM testdb.v8a LIMIT 1
< result
< column "CHARSET(c)" "VARCHAR"
< column "COLLATION(c)" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v8a
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: null
: null
> query
| SELECT COERCIBILITY(c) FROM testdb.v8a LIMIT 1
< result
< column "COERCIBILITY(c)" "BIGINT"
< row
: int64 0
> exec
| CREATE VIEW v8b AS SELECT CONCAT(a COLLATE utf8mb4_general_ci, b COLLATE utf8mb4_bin) AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'
> exec
| CREATE VIEW v8c AS SELECT CONCAT(l COLLATE latin1_bin, a) AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (latin1_bin,EXPLICIT) and (utf8mb4_general_ci,IMPLICIT) for operation 'concat'
> exec
| CREATE VIEW v8d AS SELECT LENGTH(CONCAT(a COLLATE utf8mb4_general_ci, b COLLATE utf8mb4_bin)) AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'
> exec
| SET NAMES utf8mb4 COLLATE utf8mb4_general_ci
< ok 0 0
