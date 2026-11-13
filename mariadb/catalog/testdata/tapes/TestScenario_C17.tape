-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C17 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT CONCAT(a, b) AS c FROM t
< ok 0 0
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME
| 			   FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v1
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT CHARACTER_MAXIMUM_LENGTH FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v1
: string
| c
< result
< column "CHARACTER_MAXIMUM_LENGTH" "UNSIGNED BIGINT"
< row
: uint64 20
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 			b VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_swedish_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v2 AS SELECT CONCAT(a, b) AS c FROM t
< ok 0 0
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME
| 			   FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v2
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
| 		)
< ok 0 0
> exec
| INSERT INTO t VALUES ('x','y')
< ok 1 0
> exec
| SELECT 1 FROM t WHERE a = b
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_0900_ai_ci,IMPLICIT) and (utf8mb4_0900_as_cs,IMPLICIT) for operation '='
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v4 AS SELECT CONCAT_WS(',', a, b, NULL) AS c FROM t
< ok 0 0
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME
| 			   FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v4
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> exec
| INSERT INTO t VALUES (NULL, 'x')
< ok 1 0
> query
| SELECT CONCAT(a,b), CONCAT_WS(',',a,b) FROM t LIMIT 1
< result
< column "CONCAT(a,b)" "VARCHAR"
< column "CONCAT_WS(',',a,b)" "VARCHAR"
< row
: null
: bytes
| x
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci)
< ok 0 0
> exec
| CREATE VIEW v5a AS SELECT CONCAT(a, 'x') AS c FROM t
< ok 0 0
> exec
| CREATE VIEW v5b AS SELECT CONCAT(a, _utf8mb4'x') AS c FROM t
< ok 0 0
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME
| 			   FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v5a
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_swedish_ci
> query ^
: string
| v5b
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_swedish_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 3 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_swedish_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
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
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME
| 			   FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v6a
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| latin1
: bytes
| latin1_swedish_ci
> query ^
: string
| v6b
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> query ^
: string
| v6c
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 4 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
| 		)
< ok 0 0
> exec
| CREATE VIEW v7 AS SELECT CONCAT(a, CONVERT(b USING utf8mb4)) AS c FROM t
< ok 0 0
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME
| 			   FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v7
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_as_cs
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
| 		)
< ok 0 0
> exec
| CREATE VIEW v8a AS SELECT CONCAT(a, b COLLATE utf8mb4_bin) AS c FROM t
< ok 0 0
> query
| SELECT CHARACTER_SET_NAME, COLLATION_NAME
| 			   FROM information_schema.COLUMNS
| 			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| v8a
: string
| c
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
| 			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
| 		)
< ok 0 0
> exec
| CREATE VIEW v8b AS SELECT CONCAT(a COLLATE utf8mb4_0900_ai_ci, b COLLATE utf8mb4_bin) AS c FROM t
< error
| 1267 HY000 Illegal mix of collations (utf8mb4_0900_ai_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'
