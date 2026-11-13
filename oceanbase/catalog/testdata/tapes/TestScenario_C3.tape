-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `ts1` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `ts2` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT NULL PRIMARY KEY)
< error
| 1171 42000 All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead
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
< ok 2 0
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
< ok 2 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `ts1` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
|   `ts2` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| SET SESSION explicit_defaults_for_timestamp=1
< ok 0 0
> exec
| SET SESSION sql_mode=DEFAULT
< ok 0 0
