-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSimulation default#1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a
< error
| 1051 42S02 Unknown table 'omni_t_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a
< error
| 1051 42S02 Unknown table 'omni_t_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX a ON a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a
< error
| 1017 HY000 File not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCESS_INFO)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCOUNT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(LIKE a)
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `a` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `a` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `a` tinytext DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(FULLTEXT(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(VECTOR KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE RESTRICT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON DELETE RESTRICT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) CHECK(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE CASCADE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET NULL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE NO ACTION
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH SIMPLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH FULL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH PARTIAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT AS(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a) ENFORCED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a) ENFORCED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT GENERATED ALWAYS AS(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NULL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NOT NULL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE KEY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) KEY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) COMMENT 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) CHECK(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) SRID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) INVISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VIRTUAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) STORED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT, a TINYINT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT, a TINYINT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARCHAR(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR VARCHAR(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NVARCHAR(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL VARCHAR(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER VARYING(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER VARYING(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARBINARY(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LINESTRING
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POLYGON
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTILINESTRING
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOLYGON
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRYCOLLECTION
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY(TINYINT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < TINYINT >
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < ARRAY < TINYINT >>
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT [ ]
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT [ ]
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VECTOR(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MAP(TINYINT, TINYINT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SPARSEVECTOR
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ROARINGBITMAP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SMALLINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INTEGER
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL PRECISION
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARCHAR
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a STRING
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARBINARY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1.5)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1, 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER(1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNSIGNED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SIGNED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ZEROFILL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER BINARY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL NULL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NOT NULL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ORIG_DEFAULT CURRENT_TIMESTAMP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AUTO_INCREMENT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE KEY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT KEY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COMMENT 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHECK(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SRID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLLATE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX()
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT FIXED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DYNAMIC
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DISK
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE MEMORY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT VISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT INVISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT + 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT - 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT PRIMARY KEY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `a` tinyint(4) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL NOPARALLEL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL, NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_MODE 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_SCOPE 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXPIRE_INFO(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PROGRESSIVE_MERGE_NUM 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a BLOCK_SIZE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REPLICA_NUM 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_FORMAT_VERSION 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a USE_BLOOM_FILTER TRUE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLEGROUP a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a READ ONLY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLET_SIZE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PCTFREE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_USED_PART_ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_MODE 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_CACHE_SIZE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_EXTENDED_ROWID TRUE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION_TYPE USER_SPECIFIED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCATION 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PROPERTIES(TYPE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PATTERN 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KV_ATTRIBUTES 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DEFAULT_LOB_INROW_THRESHOLD 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOB_INROW_THRESHOLD 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MICRO_INDEX_CLUSTERED TRUE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_MACRO_BLOCK_BLOOM_FILTER TRUE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH OFF
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH IMMEDIATE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH INTERVAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_READ_CONSISTENCY 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENCRYPTION 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION()
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION INDEX
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE PARTIAL_UPDATE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE ENABLE = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(PRECREATE_TIME = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(EXPIRE_TIME = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(BIGINT_PRECISION = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a, BOUNDARY_COLUMN = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN_UNIT = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 YEAR)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MONTH)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 WEEK)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 DAY)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 HOUR)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MINUTE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 SECOND)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE DELETE_INSERT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARALLEL 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND, a + INTERVAL 1 SECOND)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MINUTE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 HOUR)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 DAY)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MONTH)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 YEAR)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY HASH(a) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY RANGE() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(EACH COLUMN) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(a(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS, ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE AUTO PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE COLUMNS(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)(PARTITION a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST COLUMNS(a)(PARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)(PARTITION a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE()
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS()
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE() SIZE('s')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1 SUBPARTITION BY RANGE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) PARTITIONS 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a, PARTITION a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE, PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a), PARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES LESS THAN MAXVALUE))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES IN(a)))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a, SUBPARTITION a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ENGINE = INNODB)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE, SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a), SUBPARTITION a VALUES IN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT_AUTH))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a, a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ID 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a, a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1.5
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION HEAP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT CONDENSED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's' TYPE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_DELIMITER = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(LINE_DELIMITER = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ESCAPE = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_OPTIONALLY_ENCLOSED_BY = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_ENCLOSED_BY = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ENCODING = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_HEADER = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_BLANK_LINES = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TRIM_SPACE = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(NULL_IF =(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(EMPTY_FIELD_AS_NULL = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_GROUP_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_INDEX_TYPE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(STRIPE_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION_BLOCK_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_INDEX_STRIDE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_BLOOM_FILTER =(1))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FILE_EXTENSION = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(PARSE_HEADER = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(IGNORE_LAST_EMPTY_COLUMN = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CASCADED CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH LOCAL CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a(a) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a TABLE_ID = 1 AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VECTOR INDEX a ON a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a NOT ENFORCED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a ENFORCED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT zzz UNIQUE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a ID 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL GLOBAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL GLOBAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1 MIN_TOKEN_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_TOKEN_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(NGRAM_TOKEN_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(STOPWORD_TABLE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(DICT_TABLE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(QUANTIFIER_TABLE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(IK_MODE = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_NGRAM_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_NGRAM_SIZE = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) LOCAL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) BLOCK_SIZE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORING(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH ROWID
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) DATA_TABLE_ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) INDEX_TABLE_ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) MAX_USED_PART_ID 1
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) NOPARALLEL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY HASH(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a, a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a zzz
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a zzz
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a zzz
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO zzz
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCOUNT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a
< error
| 1017 HY000 File not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a, NOPARALLEL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP TABLEGROUP
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REFRESH
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE TTL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TTL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INPLACE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM COPY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INSTANT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK NONE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK SHARED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK EXCLUSIVE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a zzz
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER PARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER SUBPARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a zzz
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a zzz
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION a WITH TABLE a WITHOUT VALIDATION
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT PRIMARY KEY(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a VISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a NOPARALLEL
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CHECK a ENFORCED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(KEY(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(UNIQUE(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT UNIQUE(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT PRIMARY KEY(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(SPATIAL(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(PRIMARY KEY(a))
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT FOREIGN KEY(a) REFERENCES a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a INVISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS) DELAYED
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN GROUP(ALL COLUMNS)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN a TINYINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN a DROP DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE COLUMN a a TINYINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a TINYINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY COLUMN a TINYINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a TINYINT
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT FIRST
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT BEFORE a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AFTER a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE
< error
| 1051 42S02 Unknown table 'omni_t_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT
< error
| 1051 42S02 Unknown table 'omni_t_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(zzz)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1, a = 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUNCATE)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1.5)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = a)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(MIN_MAX)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(, MIN_MAX)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(SUM)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 's'
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1 CHUNK 1)
< error
| 1146 42S02 Table 'omni_t_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a
< error
| 1347 HY000 'omni_t_simulation.a' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX a ON a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCESS_INFO)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCOUNT)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(LIKE a)
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(FULLTEXT(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(VECTOR KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON DELETE RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) CHECK(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE NO ACTION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH SIMPLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH FULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH PARTIAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a)
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT AS(a)
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a) ENFORCED
< error
| 3820 HY000 Check constraint 'a_OBCHECK_0000000000000001' refers to non-existing column 'a'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a) ENFORCED
< error
| 3820 HY000 Check constraint 'a_OBCHECK_0000000000000002' refers to non-existing column 'a'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a)
< error
| 3820 HY000 Check constraint 'a_OBCHECK_0000000000000003' refers to non-existing column 'a'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a)
< error
| 3820 HY000 Check constraint 'a_OBCHECK_0000000000000004' refers to non-existing column 'a'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT GENERATED ALWAYS AS(a)
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NULL
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NOT NULL
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE KEY
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) KEY
< error
| 3106 HY000 'Defining a generated column as primary key' is not supported for generated columns.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) COMMENT 's'
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) CHECK(a)
< error
| 1235 0A000 Add column-level check constraint while altering table not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) SRID 1
< error
| 1221 HY000 Incorrect usage of srid.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VISIBLE
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) INVISIBLE
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VIRTUAL
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) STORED
< error
| 1054 42S22 Unknown column 'a' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinytext DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` float DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` datetime DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` date DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR VARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NVARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL VARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER VARYING(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER VARYING(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyblob DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinytext DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARBINARY(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` varbinary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a 's'
< error
| 5874 HY000 invalid datatype
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` bit(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s')
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` enum('s') DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s')
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` set('s') DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` json DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` geometry DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` point DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LINESTRING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` linestring DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POLYGON
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` polygon DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` multipoint DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTILINESTRING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` multilinestring DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOLYGON
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` multipolygon DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRYCOLLECTION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` geomcollection DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY(TINYINT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` ARRAY(TINYINT) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < TINYINT >
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` ARRAY(TINYINT) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < ARRAY < TINYINT >>
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` ARRAY(ARRAY(TINYINT)) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT [ ]
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` ARRAY(TINYINT) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT [ ]
< error
| 1235 0A000 Not supported feature or function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VECTOR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` VECTOR(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MAP(TINYINT, TINYINT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` MAP(TINYINT,TINYINT) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SPARSEVECTOR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` SPARSEVECTOR DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ROARINGBITMAP
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` roaringbitmap DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SMALLINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` smallint(6) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` mediumint(9) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INTEGER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL PRECISION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` timestamp NULL DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` time DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` year(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` text DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` longtext DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARCHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a STRING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` string DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` blob DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` mediumblob DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` longblob DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARBINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` mediumblob DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` bit(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` float(1,1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` float DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1.5)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` float DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1, 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` decimal(1,1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` decimal(1,0) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` year(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` datetime(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNSIGNED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(3) unsigned DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SIGNED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ZEROFILL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(3) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER BINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NOT NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) NOT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ORIG_DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AUTO_INCREMENT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE KEY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COMMENT 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL COMMENT 's',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHECK(a)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000005'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SRID 1
< error
| 1221 HY000 Incorrect usage of srid.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLLATE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 1
< error
| 1235 0A000 chunk size column attribute only supported for json, others not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT FIXED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DISK
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE MEMORY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT INVISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) INVISIBLE DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT(a)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT + 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT '1',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT '1',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT - 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT '-1',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT PRIMARY KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) NOPARALLEL
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL, NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_MODE 's'
< error
| 1235 0A000 Table mode s is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_SCOPE 's'
< error
| 1210 HY000 Incorrect arguments to duplicate_scope
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXPIRE_INFO(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PROGRESSIVE_MERGE_NUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 PROGRESSIVE_MERGE_NUM = 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REPLICA_NUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION 's'
< error
| 1210 HY000 Incorrect arguments to s
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = REDUNDANT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_FORMAT_VERSION 1
< error
| 4179 HY000 target storage format version cannot be smaller than current version not allowed
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a USE_BLOOM_FILTER TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = TRUE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE a
< error
| 1273 HY000 Unknown collation: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLEGROUP a
< error
| 5151 42P01 tablegroup not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a READ ONLY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 READ ONLY
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLET_SIZE 1
< error
| 4147 HY000 Invalid config
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PCTFREE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_USED_PART_ID 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a
< error
| 3510 HY000 Tablespace does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_MODE 's'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_CACHE_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_EXTENDED_ROWID TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION_TYPE USER_SPECIFIED
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCATION 's'
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PROPERTIES(TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PATTERN 's'
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND)
< error
| 10502 HY000 TTL column 'a' not exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KV_ATTRIBUTES 's'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DEFAULT_LOB_INROW_THRESHOLD 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 LOB_INROW_THRESHOLD=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOB_INROW_THRESHOLD 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 LOB_INROW_THRESHOLD=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MICRO_INDEX_CLUSTERED TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_MACRO_BLOCK_BLOOM_FILTER TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH OFF
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH IMMEDIATE
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH INTERVAL
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_READ_CONSISTENCY 's'
< error
| 1210 HY000 Incorrect arguments to duplicate_read_consistency
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENCRYPTION 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION INDEX
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 's'
< error
| 1210 HY000 Incorrect arguments to semistruct encoding type, should be none or encoding
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE PARTIAL_UPDATE
< error
| 1235 0A000 alter merge engine table option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE)
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE ENABLE = TRUE)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(PRECREATE_TIME = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(EXPIRE_TIME = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = 's')
< error
| 1298 HY000 Unknown or incorrect time zone
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(BIGINT_PRECISION = 's')
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a, BOUNDARY_COLUMN = a)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN_UNIT = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 YEAR)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MONTH)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 WEEK)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 DAY)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 HOUR)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MINUTE)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 SECOND)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE DELETE_INSERT
< error
| 1235 0A000 alter merge engine table option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARALLEL 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND, a + INTERVAL 1 SECOND)
< error
| 10502 HY000 TTL column 'a' not exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MINUTE)
< error
| 10502 HY000 TTL column 'a' not exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 HOUR)
< error
| 10502 HY000 TTL column 'a' not exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 DAY)
< error
| 10502 HY000 TTL column 'a' not exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MONTH)
< error
| 10502 HY000 TTL column 'a' not exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 YEAR)
< error
| 10502 HY000 TTL column 'a' not exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY HASH(a) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY RANGE() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(EACH COLUMN) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(a(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS, ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE AUTO PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE COLUMNS(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)(PARTITION a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST COLUMNS(a)(PARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)(PARTITION a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `p0`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE() SIZE('s')
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1 SUBPARTITION BY RANGE(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) PARTITIONS 1
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `p0`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a)
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a, PARTITION a)
< error
| 1517 HY000 Duplicate partition name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE, PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a), PARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES LESS THAN MAXVALUE))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES IN(a)))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a, SUBPARTITION a))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ENGINE = INNODB)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE, SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a), SUBPARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT_AUTH))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a, a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ID 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a, a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1.5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION HEAP
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT CONDENSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = CONDENSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's' TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_DELIMITER = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(LINE_DELIMITER = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ESCAPE = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_OPTIONALLY_ENCLOSED_BY = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_ENCLOSED_BY = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ENCODING = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_HEADER = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_BLANK_LINES = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TRIM_SPACE = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(NULL_IF =(a))
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(EMPTY_FIELD_AS_NULL = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_GROUP_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_INDEX_TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(STRIPE_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION_BLOCK_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_INDEX_STRIDE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_BLOOM_FILTER =(1))
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FILE_EXTENSION = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(PARSE_HEADER = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(IGNORE_LAST_EMPTY_COLUMN = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CASCADED CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH LOCAL CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a(a) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a TABLE_ID = 1 AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VECTOR INDEX a ON a(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a NOT ENFORCED
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a ENFORCED
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT zzz UNIQUE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK a
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((a))
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a ID 1)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1))
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL GLOBAL
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL GLOBAL
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1 MIN_TOKEN_SIZE = 1)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_TOKEN_SIZE = 1)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(NGRAM_TOKEN_SIZE = 1)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(STOPWORD_TABLE = 's')
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(DICT_TABLE = 's')
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(QUANTIFIER_TABLE = 's')
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(IK_MODE = 's')
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_NGRAM_SIZE = 1)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_NGRAM_SIZE = 1)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) LOCAL
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's'
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORING(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH ROWID
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) DATA_TABLE_ID 1
< error
| 1235 0A000 specify data table id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) INDEX_TABLE_ID 1
< error
| 1235 0A000 specify index table id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID 1
< error
| 1235 0A000 specify virtual column id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) MAX_USED_PART_ID 1
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) NOPARALLEL
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY HASH(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a, a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCOUNT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a, NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP TABLEGROUP
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REFRESH
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ttl
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ttl`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| ttl
: bytes
| CREATE TABLE `ttl` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INPLACE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM COPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INSTANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK NONE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK SHARED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK EXCLUSIVE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT(a)
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK(a)
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION
< error
| 1091 42000 Can't DROP 'SUBPARTITION'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER PARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER SUBPARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION a WITH TABLE a WITHOUT VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION(a)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT PRIMARY KEY(a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a VISIBLE
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a NOPARALLEL
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CHECK a ENFORCED
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(KEY(a))
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(UNIQUE(a))
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT UNIQUE(a))
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT PRIMARY KEY(a))
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(SPATIAL(a))
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(PRIMARY KEY(a))
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT FOREIGN KEY(a) REFERENCES a(a)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a INVISIBLE
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN GROUP(ALL COLUMNS)
< error
| 4185 HY000 Column group 'all column group' not found
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN a DROP DEFAULT
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE COLUMN a a TINYINT
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a TINYINT
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY COLUMN a TINYINT
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a TINYINT
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT FIRST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `a` tinyint(4) DEFAULT NULL,
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT BEFORE a
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AFTER a
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT 's'
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT(a)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(zzz)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1, a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUE)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUNCATE)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1.5)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = a)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(MIN_MAX)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL SKIP_INDEX(MIN_MAX),
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(, MIN_MAX)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL SKIP_INDEX(MIN_MAX),
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(SUM)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` tinyint(4) DEFAULT NULL SKIP_INDEX(SUM),
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 's'
< error
| 1235 0A000 chunk size column attribute only supported for json, others not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1 CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a
< error
| 1347 HY000 'omni_t_simulation.a' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX a ON a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCESS_INFO)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCOUNT)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(LIKE a)
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(FULLTEXT(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(VECTOR KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE RESTRICT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON DELETE RESTRICT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) CHECK(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE CASCADE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE NO ACTION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET DEFAULT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH SIMPLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH FULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH PARTIAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a) ENFORCED
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000006'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a) ENFORCED
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000007'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000008'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000009'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT GENERATED ALWAYS AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NOT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE KEY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) KEY
< error
| 3106 HY000 'Defining a generated column as primary key' is not supported for generated columns.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) COMMENT 's'
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) CHECK(a)
< error
| 1235 0A000 Add column-level check constraint while altering table not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) SRID 1
< error
| 1221 HY000 Incorrect usage of srid.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) INVISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VIRTUAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) STORED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NVARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARBINARY(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a 's'
< error
| 5874 HY000 invalid datatype
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s')
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s')
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LINESTRING
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POLYGON
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTILINESTRING
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOLYGON
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRYCOLLECTION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY(TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < TINYINT >
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < ARRAY < TINYINT >>
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT [ ]
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT [ ]
< error
| 1235 0A000 Not supported feature or function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VECTOR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MAP(TINYINT, TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SPARSEVECTOR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ROARINGBITMAP
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SMALLINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INTEGER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL PRECISION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARCHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a STRING
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARBINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1.5)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1, 1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNSIGNED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SIGNED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ZEROFILL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER BINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NOT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ORIG_DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AUTO_INCREMENT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE KEY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COMMENT 's'
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHECK(a)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000010'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SRID 1
< error
| 1221 HY000 Incorrect usage of srid.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLLATE a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX()
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 1
< error
| 1235 0A000 chunk size column attribute only supported for json, others not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DEFAULT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT FIXED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DYNAMIC
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DEFAULT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DISK
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE MEMORY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT VISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT INVISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT(a)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT + 1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT - 1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT PRIMARY KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) NOPARALLEL
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL, NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_MODE 's'
< error
| 1235 0A000 Table mode s is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_SCOPE 's'
< error
| 1210 HY000 Incorrect arguments to duplicate_scope
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXPIRE_INFO(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PROGRESSIVE_MERGE_NUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 PROGRESSIVE_MERGE_NUM = 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REPLICA_NUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION 's'
< error
| 1210 HY000 Incorrect arguments to s
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = REDUNDANT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_FORMAT_VERSION 1
< error
| 4179 HY000 target storage format version cannot be smaller than current version not allowed
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a USE_BLOOM_FILTER TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = TRUE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE a
< error
| 1273 HY000 Unknown collation: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLEGROUP a
< error
| 5151 42P01 tablegroup not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a READ ONLY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 READ ONLY
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLET_SIZE 1
< error
| 4147 HY000 Invalid config
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PCTFREE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_USED_PART_ID 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a
< error
| 3510 HY000 Tablespace does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_MODE 's'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_CACHE_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_EXTENDED_ROWID TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION_TYPE USER_SPECIFIED
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCATION 's'
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PROPERTIES(TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PATTERN 's'
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KV_ATTRIBUTES 's'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DEFAULT_LOB_INROW_THRESHOLD 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 LOB_INROW_THRESHOLD=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOB_INROW_THRESHOLD 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 LOB_INROW_THRESHOLD=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MICRO_INDEX_CLUSTERED TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_MACRO_BLOCK_BLOOM_FILTER TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH OFF
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH IMMEDIATE
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH INTERVAL
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_READ_CONSISTENCY 's'
< error
| 1210 HY000 Incorrect arguments to duplicate_read_consistency
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENCRYPTION 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION INDEX
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 's'
< error
| 1210 HY000 Incorrect arguments to semistruct encoding type, should be none or encoding
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE PARTIAL_UPDATE
< error
| 1235 0A000 alter merge engine table option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE)
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE ENABLE = TRUE)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(PRECREATE_TIME = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(EXPIRE_TIME = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = 's')
< error
| 1298 HY000 Unknown or incorrect time zone
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(BIGINT_PRECISION = 's')
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a, BOUNDARY_COLUMN = a)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN_UNIT = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 YEAR)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MONTH)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 WEEK)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 DAY)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 HOUR)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MINUTE)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 SECOND)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE DELETE_INSERT
< error
| 1235 0A000 alter merge engine table option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARALLEL 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND, a + INTERVAL 1 SECOND)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MINUTE)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 HOUR)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 DAY)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MONTH)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 YEAR)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY HASH(a) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY RANGE() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(EACH COLUMN) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(a(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS, ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE AUTO PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE COLUMNS(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)(PARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST COLUMNS(a)(PARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)(PARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `p0`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS(a)
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE() SIZE('s')
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1 SUBPARTITION BY RANGE(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) PARTITIONS 1
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `p0`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a, PARTITION a)
< error
| 1517 HY000 Duplicate partition name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE, PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a), PARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES LESS THAN MAXVALUE))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES IN(a)))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a, SUBPARTITION a))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ENGINE = INNODB)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE, SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a), SUBPARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT_AUTH))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a, a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ID 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a, a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1.5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION HEAP
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT CONDENSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = CONDENSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's' TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_DELIMITER = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(LINE_DELIMITER = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ESCAPE = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_OPTIONALLY_ENCLOSED_BY = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_ENCLOSED_BY = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ENCODING = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_HEADER = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_BLANK_LINES = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TRIM_SPACE = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(NULL_IF =(a))
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(EMPTY_FIELD_AS_NULL = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_GROUP_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_INDEX_TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(STRIPE_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION_BLOCK_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_INDEX_STRIDE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_BLOOM_FILTER =(1))
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FILE_EXTENSION = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(PARSE_HEADER = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(IGNORE_LAST_EMPTY_COLUMN = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CASCADED CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH LOCAL CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a(a) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a TABLE_ID = 1 AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL WITH COLUMN GROUP(all columns)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VECTOR INDEX a ON a(a)
< error
| 7601 HY000 Column 'a' cannot be part of VECTOR index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a NOT ENFORCED
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a ENFORCED
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT zzz UNIQUE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK a
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((a))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a ID 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 GLOBAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL GLOBAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 GLOBAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL GLOBAL
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1 MIN_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(NGRAM_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(STOPWORD_TABLE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(DICT_TABLE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(QUANTIFIER_TABLE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(IK_MODE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_NGRAM_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_NGRAM_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) LOCAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 's'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORING(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH ROWID
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) DATA_TABLE_ID 1
< error
| 1235 0A000 specify data table id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) INDEX_TABLE_ID 1
< error
| 1235 0A000 specify index table id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID 1
< error
| 1235 0A000 specify virtual column id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) MAX_USED_PART_ID 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `zzz` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `BTREE` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY HASH(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `HASH` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a, a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCOUNT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a, NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP TABLEGROUP
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REFRESH
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ttl
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ttl`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| ttl
: bytes
| CREATE TABLE `ttl` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INPLACE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM COPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INSTANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK NONE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK SHARED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK EXCLUSIVE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT(a)
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK(a)
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION
< error
| 1091 42000 Can't DROP 'SUBPARTITION'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER PARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER SUBPARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION a WITH TABLE a WITHOUT VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION(a)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT PRIMARY KEY(a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a VISIBLE
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a NOPARALLEL
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CHECK a ENFORCED
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(KEY(a))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL WITH COLUMN GROUP(all columns)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(UNIQUE(a))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL WITH COLUMN GROUP(all columns),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT UNIQUE(a))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL WITH COLUMN GROUP(all columns),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT PRIMARY KEY(a))
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(SPATIAL(a))
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(PRIMARY KEY(a))
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT FOREIGN KEY(a) REFERENCES a(a)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a INVISIBLE
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN GROUP(ALL COLUMNS)
< error
| 4185 HY000 Column group 'all column group' not found
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN a TINYINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN a DROP DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE COLUMN a a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY COLUMN a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT FIRST
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT BEFORE a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AFTER a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT(a)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(zzz)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1, a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUE)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUNCATE)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1.5)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = a)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(MIN_MAX)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(, MIN_MAX)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(SUM)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 's'
< error
| 1235 0A000 chunk size column attribute only supported for json, others not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1 CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP VIEW a
< error
| 1347 HY000 'omni_t_simulation.a' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a)
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP INDEX a ON a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| RENAME TABLE a TO a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| TRUNCATE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCESS_INFO)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(BINARY_FORMAT = ACCOUNT)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(LIKE a)
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY KEY() WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS) AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL IGNORE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT)
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYTEXT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(PRIMARY KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(FULLTEXT(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(VECTOR KEY(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(CONSTRAINT FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(FOREIGN KEY(a) REFERENCES a(a))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE RESTRICT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON DELETE RESTRICT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) CHECK(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE CASCADE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE NO ACTION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) ON UPDATE SET DEFAULT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH SIMPLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH FULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT REFERENCES a(a) MATCH PARTIAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a) ENFORCED
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000011'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a) ENFORCED
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000012'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT CHECK(a)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000013'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000014'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT GENERATED ALWAYS AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) NOT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE KEY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) KEY
< error
| 3106 HY000 'Defining a generated column as primary key' is not supported for generated columns.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) COMMENT 's'
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) CHECK(a)
< error
| 1235 0A000 Add column-level check constraint while altering table not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) SRID 1
< error
| 1221 HY000 Incorrect usage of srid.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) INVISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) VIRTUAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AS(a) STORED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYTEXT, a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NVARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VARBINARY(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a 's'
< error
| 5874 HY000 invalid datatype
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s')
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s')
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LINESTRING
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a POLYGON
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTILINESTRING
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MULTIPOLYGON
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a GEOMETRYCOLLECTION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY(TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < TINYINT >
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ARRAY < ARRAY < TINYINT >>
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT [ ]
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT [ ]
< error
| 1235 0A000 Not supported feature or function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a VECTOR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MAP(TINYINT, TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SPARSEVECTOR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ROARINGBITMAP
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SMALLINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INTEGER
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL PRECISION
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARCHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a STRING
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT VARBINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1.5)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1, 1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMBER(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNSIGNED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SIGNED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ZEROFILL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHARACTER BINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NULL NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT NOT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ORIG_DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AUTO_INCREMENT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE KEY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COMMENT 's'
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHECK(a)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000015'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SRID 1
< error
| 1221 HY000 Incorrect usage of srid.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLLATE a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX()
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 1
< error
| 1235 0A000 chunk size column attribute only supported for json, others not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DEFAULT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT FIXED
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT COLUMN_FORMAT DYNAMIC
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DEFAULT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE DISK
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT STORAGE MEMORY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT VISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT INVISIBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT(a)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT + 1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT 1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT DEFAULT - 1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT PRIMARY KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a TINYINT) NOPARALLEL
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a NOPARALLEL NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL, NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL NOPARALLEL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_MODE 's'
< error
| 1235 0A000 Table mode s is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_SCOPE 's'
< error
| 1210 HY000 Incorrect arguments to duplicate_scope
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a EXPIRE_INFO(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PROGRESSIVE_MERGE_NUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 PROGRESSIVE_MERGE_NUM = 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_ID 1
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REPLICA_NUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION 's'
< error
| 1210 HY000 Incorrect arguments to s
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = REDUNDANT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_FORMAT_VERSION 1
< error
| 4179 HY000 target storage format version cannot be smaller than current version not allowed
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a USE_BLOOM_FILTER TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = TRUE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE a
< error
| 1273 HY000 Unknown collation: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TABLEGROUP a
< error
| 5151 42P01 tablegroup not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a READ ONLY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 READ ONLY
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TABLET_SIZE 1
< error
| 4147 HY000 Invalid config
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PCTFREE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_USED_PART_ID 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a
< error
| 3510 HY000 Tablespace does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_MODE 's'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_CACHE_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_EXTENDED_ROWID TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION_TYPE USER_SPECIFIED
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a LOCATION 's'
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PROPERTIES(TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PATTERN 's'
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a KV_ATTRIBUTES 's'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DEFAULT_LOB_INROW_THRESHOLD 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 LOB_INROW_THRESHOLD=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a LOB_INROW_THRESHOLD 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 LOB_INROW_THRESHOLD=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MICRO_INDEX_CLUSTERED TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE_MACRO_BLOCK_BLOOM_FILTER TRUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH OFF
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH IMMEDIATE
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_REFRESH INTERVAL
< error
| 1235 0A000 Using ALTER AUTO REFRESH as a ALTER TABLE option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_READ_CONSISTENCY 's'
< error
| 1210 HY000 Incorrect arguments to duplicate_read_consistency
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ENCRYPTION 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a UNION()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION INDEX
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 's'
< error
| 1210 HY000 Incorrect arguments to semistruct encoding type, should be none or encoding
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE PARTIAL_UPDATE
< error
| 1235 0A000 alter merge engine table option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE)
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(ENABLE = TRUE ENABLE = TRUE)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(PRECREATE_TIME = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(EXPIRE_TIME = 's')
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = 's')
< error
| 1298 HY000 Unknown or incorrect time zone
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(BIGINT_PRECISION = 's')
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN = a, BOUNDARY_COLUMN = a)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(BOUNDARY_COLUMN_UNIT = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 YEAR)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MONTH)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 WEEK)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 DAY)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 HOUR)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 MINUTE)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE_CACHE_POLICY(HOT_RETENTION = 1 SECOND)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MERGE_ENGINE DELETE_INSERT
< error
| 1235 0A000 alter merge engine table option not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARALLEL 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 SECOND, a + INTERVAL 1 SECOND)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MINUTE)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 HOUR)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 DAY)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 MONTH)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TTL(a + INTERVAL 1 YEAR)
< error
| 10503 HY000 Column type of 'a' is not supported for TTL definition
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY HASH(a) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY RANGE() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY KEY() SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(EACH COLUMN) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(a(a)) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a WITH COLUMN GROUP(ALL COLUMNS, ALL COLUMNS) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE() PARTITION SIZE AUTO PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a NOPARALLEL PARTITION BY RANGE COLUMNS(a) PARTITION SIZE 's' PARTITIONS AUTO SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(a)(PARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST(a)(PARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LIST COLUMNS(a)(PARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(a)(PARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `p0`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS(a)
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE() SIZE('s')
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1064 42000 Parse error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a) SIZE('s')
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1 SUBPARTITION BY RANGE(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) PARTITIONS 1
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() PARTITIONS 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `p0`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a) SUBPARTITION TEMPLATE(SUBPARTITION a VALUES IN(a))
< error
| 1054 42S22 Unknown column 'a' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a) SUBPARTITION TEMPLATE(SUBPARTITION a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY RANGE COLUMNS(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY LIST COLUMNS(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY KEY(a)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a, PARTITION a)
< error
| 1517 HY000 Duplicate partition name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN MAXVALUE, PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a), PARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES LESS THAN MAXVALUE))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a VALUES IN(a)))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a(SUBPARTITION a, SUBPARTITION a))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ENGINE = INNODB)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES LESS THAN MAXVALUE, SUBPARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY PARTITION a ADD SUBPARTITION(SUBPARTITION a VALUES IN(a), SUBPARTITION a VALUES IN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(DEFAULT_AUTH))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES IN(a, a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a ID 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(PARTITION a VALUES LESS THAN(a, a))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1.5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ORGANIZATION HEAP
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT CONDENSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = CONDENSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TYPE = 's' TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_DELIMITER = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(LINE_DELIMITER = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ESCAPE = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_OPTIONALLY_ENCLOSED_BY = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FIELD_ENCLOSED_BY = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ENCODING = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_HEADER = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(SKIP_BLANK_LINES = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(TRIM_SPACE = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(NULL_IF =(a))
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(EMPTY_FIELD_AS_NULL = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_GROUP_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION = a)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_INDEX_TYPE = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(STRIPE_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COMPRESSION_BLOCK_SIZE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(ROW_INDEX_STRIDE = 1)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(COLUMN_BLOOM_FILTER =(1))
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(FILE_EXTENSION = 's')
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(PARSE_HEADER = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORMAT(IGNORE_LAST_EMPTY_COLUMN = TRUE)
< error
| 4016 HY000 Internal error
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| TRUNCATE TABLE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH CASCADED CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a WITH LOCAL CHECK OPTION
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a(a) AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a TABLE_ID = 1 AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VECTOR INDEX a ON a(a)
< error
| 7601 HY000 Column 'a' cannot be part of VECTOR index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a)
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a NOT ENFORCED
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT a ENFORCED
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT zzz UNIQUE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK a
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((a))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a ID 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 GLOBAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) GLOBAL GLOBAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 GLOBAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) GLOBAL GLOBAL
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_TOKEN_SIZE = 1 MIN_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(NGRAM_TOKEN_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(STOPWORD_TABLE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(DICT_TABLE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(QUANTIFIER_TABLE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(IK_MODE = 's')
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MIN_NGRAM_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) PARSER_PROPERTIES(MAX_NGRAM_SIZE = 1)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) LOCAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL COMMENT 's'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORING(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH ROWID
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) DATA_TABLE_ID 1
< error
| 1235 0A000 specify data table id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) INDEX_TABLE_ID 1
< error
| 1235 0A000 specify index table id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VIRTUAL_COLUMN_ID 1
< error
| 1235 0A000 specify virtual column id not for restore purpose in create index not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) MAX_USED_PART_ID 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `zzz` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `BTREE` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY HASH(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `HASH` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a, a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCESS_INFO zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION ACCOUNT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 WITH COLUMN GROUP(all columns)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a, NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP TABLEGROUP
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a
< error
| 1091 42000 Can't DROP 'a'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 3940 HY000 Constraint 'a' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REFRESH
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ttl
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ttl`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| ttl
: bytes
| CREATE TABLE `ttl` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INPLACE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM COPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM INSTANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK NONE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK SHARED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK EXCLUSIVE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT(a)
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK(a)
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION
< error
| 1091 42000 Can't DROP 'SUBPARTITION'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SUBPARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER PARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER SUBPARTITION a STORAGE_CACHE_POLICY = 's'
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(PARTITION a VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE SUBPARTITION a zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION a WITH TABLE a WITHOUT VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION(a)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT PRIMARY KEY(a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a STORAGE_CACHE_POLICY(GLOBAL = 's')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a NOPARALLEL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CHECK a ENFORCED
< error
| 3821 HY000 Check constraint 'a' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(KEY(a))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL WITH COLUMN GROUP(all columns)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(UNIQUE(a))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL WITH COLUMN GROUP(all columns),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT UNIQUE(a))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT UNIQUE(a) WITH COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a_2` (`a`) BLOCK_SIZE 16384 LOCAL WITH COLUMN GROUP(all columns),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(CONSTRAINT PRIMARY KEY(a))
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(SPATIAL(a))
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) WITH COLUMN GROUP(ALL COLUMNS)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(PRIMARY KEY(a))
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT FOREIGN KEY(a) REFERENCES a(a)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX a INVISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL /*!80000 INVISIBLE */
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN GROUP(ALL COLUMNS) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN GROUP(ALL COLUMNS)
< error
| 4185 HY000 Column group 'all column group' not found
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN a TINYINT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a TINYINT) JSON(a) STORE AS(CHUNK 1)
< error
| 1210 HY000 Invalid argument
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN a DROP DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE COLUMN a a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY COLUMN a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` tinyint(4) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT FIRST
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT BEFORE a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT AFTER a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT(a)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(zzz)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a UNION(a)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1, a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUE)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = TRUNCATE)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = 1.5)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) WITH(a = a)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(MIN_MAX)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(, MIN_MAX)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT SKIP_INDEX(SUM)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT CHUNK 's'
< error
| 1235 0A000 chunk size column attribute only supported for json, others not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a JSON(a) STORE AS(CHUNK 1 CHUNK 1)
< error
| 1210 HY000 Invalid argument
