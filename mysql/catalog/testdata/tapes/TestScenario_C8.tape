-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C8 default#1
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
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SELECT IFNULL(ENGINE,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(ENGINE,'')" "VARCHAR"
< row
: bytes
| InnoDB
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
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SELECT IFNULL(ROW_FORMAT,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(ROW_FORMAT,'')" "VARCHAR"
< row
: bytes
| Dynamic
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
| CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY)
< ok 0 0
> query
| SELECT IFNULL(AUTO_INCREMENT,0) FROM information_schema.TABLES
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(AUTO_INCREMENT,0)" "DECIMAL"
< row
: bytes
| 0
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
|   `id` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb CHARACTER SET latin1
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10))
< ok 0 0
> query
| SELECT IFNULL(TABLE_COLLATION,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(TABLE_COLLATION,'')" "VARCHAR"
< row
: bytes
| latin1_swedish_ci
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| latin1
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
| CREATE TABLE t (a INT) COLLATE=latin1_german2_ci
< ok 0 0
> query
| SELECT IFNULL(TABLE_COLLATION,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(TABLE_COLLATION,'')" "VARCHAR"
< row
: bytes
| latin1_german2_ci
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
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SELECT IFNULL(CREATE_OPTIONS,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(CREATE_OPTIONS,'')" "VARCHAR"
< row
: bytes
|
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
|   `a` int DEFAULT NULL
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
| CREATE TABLE t (a INT)
< ok 0 0
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
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT IFNULL(CREATE_OPTIONS,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(CREATE_OPTIONS,'')" "VARCHAR"
< row
: bytes
|
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
| CREATE TABLE t (a INT)
< ok 0 0
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
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT IFNULL(CREATE_OPTIONS,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(CREATE_OPTIONS,'')" "VARCHAR"
< row
: bytes
|
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
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SELECT IFNULL(CREATE_OPTIONS,'') FROM information_schema.TABLES WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "IFNULL(CREATE_OPTIONS,'')" "VARCHAR"
< row
: bytes
|
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
|   `a` int DEFAULT NULL
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
| CREATE TABLE t (a INT)
< ok 0 0
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
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT IFNULL(SPACE,0) FROM information_schema.INNODB_TABLES
|             WHERE NAME='testdb/t'
< result
< column "IFNULL(SPACE,0)" "BIGINT"
< row
: int64 2969
