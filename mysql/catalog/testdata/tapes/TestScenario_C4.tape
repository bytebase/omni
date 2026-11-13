-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C4 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb CHARACTER SET latin1 COLLATE latin1_swedish_ci
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c VARCHAR(10))
< ok 0 0
> query
| SELECT IFNULL(TABLE_COLLATION,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(TABLE_COLLATION,'')" "VARCHAR"
< row
: bytes
| latin1_swedish_ci
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
| CREATE TABLE t (c VARCHAR(10)) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
< ok 0 0
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `c` varchar(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| CREATE TABLE t (c VARCHAR(10) COLLATE utf8mb4_unicode_ci) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4_unicode_ci
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
| CREATE TABLE t (c VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| latin1
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| latin1_swedish_ci
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
| CREATE TABLE t_bad (c VARCHAR(10)) CHARACTER SET latin1 COLLATE utf8mb4_0900_ai_ci
< error
| 1253 42000 COLLATION 'utf8mb4_0900_ai_ci' is not valid for CHARACTER SET 'latin1'
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
| CREATE TABLE t_ok (c VARCHAR(10)) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
< ok 0 0
> query
| SELECT IFNULL(TABLE_COLLATION,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_ok'
< result
< column "IFNULL(TABLE_COLLATION,'')" "VARCHAR"
< row
: bytes
| utf8mb4_0900_ai_ci
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
|   a CHAR(10) BINARY,
|   b VARCHAR(10) CHARACTER SET latin1 BINARY
| ) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4_bin
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| latin1_bin
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
|   a BINARY(10),
|   b CHAR(10) CHARACTER SET binary,
|   c VARBINARY(10)
| )
< ok 0 0
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
|
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
|
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
|
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| binary
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| binary
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| varbinary
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
| CREATE TABLE t (c VARCHAR(10) CHARACTER SET utf8)
< ok 0 0
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3_general_ci
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
|   a NCHAR(10),
|   b NATIONAL CHARACTER(10),
|   cc NATIONAL VARCHAR(10),
|   d NCHAR VARYING(10)
| )
< ok 0 0
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3_general_ci
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3_general_ci
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3_general_ci
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3_general_ci
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
|   a ENUM('x','y'),
|   b ENUM('x','y') CHARACTER SET latin1,
|   cc SET('p','q') COLLATE utf8mb4_unicode_ci
| ) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| latin1
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| latin1_swedish_ci
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4_unicode_ci
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
| CREATE TABLE t1 (c VARCHAR(10) CHARACTER SET latin1, KEY k (c(5)))
< ok 0 0
> exec
| CREATE TABLE t2 (c VARCHAR(10) CHARACTER SET utf8mb4, KEY k (c(5)))
< ok 0 0
> exec
| CREATE TABLE t3 (c VARCHAR(200) CHARACTER SET utf8mb4, KEY k (c(768)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> query
| SELECT SUB_PART FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND INDEX_NAME='k'
< result
< column "SUB_PART" "BIGINT"
< row
: int64 5
> query
| SELECT SUB_PART FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2' AND INDEX_NAME='k'
< result
< column "SUB_PART" "BIGINT"
< row
: int64 5
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
| CREATE TABLE t (c VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci)
< ok 0 0
> query
| SELECT IFNULL(COLLATION_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IFNULL(COLLATION_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4_0900_ai_ci
> exec
| SELECT c = 'abc' FROM t
< ok 0 0
> exec
| SELECT c = _utf8mb4'abc' COLLATE utf8mb4_bin FROM t
< ok 0 0
> exec
| SELECT c COLLATE utf8mb4_bin = _latin1'abc' FROM t
< ok 0 0
> exec
| SELECT CAST('abc' AS CHAR CHARACTER SET latin1) = c FROM t
< ok 0 0
