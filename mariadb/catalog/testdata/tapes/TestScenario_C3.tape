-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C3 default#1
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
| SET SESSION explicit_defaults_for_timestamp=0
< ok 0 0
> exec
| SET SESSION sql_mode=''
< ok 0 0
> exec
| CREATE TABLE t (
|     ts1 TIMESTAMP NOT NULL,
|     ts2 TIMESTAMP NOT NULL
| )
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
|   `ts1` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
|   `ts2` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET SESSION explicit_defaults_for_timestamp=1
< ok 0 0
> exec
| SET SESSION sql_mode=DEFAULT
< ok 0 0
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
| CREATE TABLE t (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NO
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
| CREATE TABLE t (a INT AUTO_INCREMENT, KEY (a))
< ok 0 0
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NO
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
| CREATE TABLE t (a INT NULL PRIMARY KEY)
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
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
| CREATE TABLE t (a INT UNIQUE)
< ok 0 0
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
> query
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| a
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
|     a INT NULL,
|     b INT GENERATED ALWAYS AS (a+1) VIRTUAL
| )
< ok 0 0
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
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
| CREATE TABLE t (id INT NULL AUTO_INCREMENT, KEY(id))
< ok 0 0
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='id'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NO
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
| SET SESSION explicit_defaults_for_timestamp=0
< ok 0 0
> exec
| SET SESSION sql_mode=''
< ok 0 0
> exec
| CREATE TABLE t (ts1 TIMESTAMP, ts2 TIMESTAMP)
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
|   `ts1` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
|   `ts2` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET SESSION explicit_defaults_for_timestamp=1
< ok 0 0
> exec
| SET SESSION sql_mode=DEFAULT
< ok 0 0
