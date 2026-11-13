-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSimulation default#1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a;
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a;
< error
| 1051 42S02 Unknown table 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
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
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY SYSTEM_TIME;
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY SYSTEM_TIME;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD IF NOT EXISTS a INT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS a;
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE a;
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a, ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_CHECKSUM 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION();
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET utf8mb4;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE utf8mb4_bin;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE DISK;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE MEMORY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRANSACTIONAL 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SEQUENCE 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WITH SYSTEM VERSIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFINER;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a READ_ONLY = 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WRAPPER = 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a WITH SYSTEM VERSIONING;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a VERSIONING;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT FIXED;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT PAGE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a INT);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT, a INT);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a));
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a));
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(PERIOD FOR SYSTEM_TIME(a, a));
< error
| 4125 HY000 Wrong parameters for `a`: missing 'WITH SYSTEM VERSIONING'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a TYPE BTREE(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE a TYPE BTREE(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(a, a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CHECK(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT PRIMARY KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT REFERENCES a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT zzz REFERENCES a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW END;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT GENERATED ALWAYS AS ROW START;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) VIRTUAL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a)
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) PERSISTENT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) STORED;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE UNIQUE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COMMENT 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT INVISIBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT + 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s');
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMERIC;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARBINARY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARCHAR;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s');
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL(1, 1);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1.5);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL NULL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AUTO_INCREMENT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT SERIAL DEFAULT VALUE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT VALUE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COLLATE utf8mb4_bin;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL ENABLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT a = 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITH SYSTEM VERSIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITHOUT SYSTEM VERSIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR CHARSET utf8mb4;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR UNICODE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BYTE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII BINARY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY ASCII;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE utf8mb4_bin;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY zzz;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE BINARY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR BINARY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT PRIMARY KEY;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a, a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH FULL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH PARTIAL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH SIMPLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT ON DELETE RESTRICT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT ON UPDATE RESTRICT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE CASCADE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET NULL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE NO ACTION;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY zzz(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE zzz(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT zzz(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY zzz(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL zzz(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR zzz(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE VISIBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE VISIBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE VISIBLE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a USING BTREE ON a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a TYPE BTREE ON a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) TYPE BTREE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) IGNORED;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) a = 's';
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING RTREE(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a IGNORED;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a NOT IGNORED;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a, a WITHOUT OVERLAPS);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1));
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s', 's');
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE REMOVE PARTITIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME PARTITIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION ALL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a OPTIMIZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECK PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REPAIR PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COALESCE PARTITION 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION ALL;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT PARTITION a TO TABLE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION PARTITIONS 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE, FORCE;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS a CHECK(a);
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a INT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a INT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISABLE KEYS;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE KEYS;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4 COLLATE utf8mb4_bin;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SYSTEM VERSIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SYSTEM VERSIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP VERSIONING;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR SYSTEM_TIME;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT LOCK DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE IGNORE TABLE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ONLINE TABLE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT;
< error
| 1051 42S02 Unknown table 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a
< error
| 1051 42S02 Unknown table 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE;
< error
| 1051 42S02 Unknown table 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AFTER a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT FIRST;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME = a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a, a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 9223372036854775808;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 2147483648;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1.5;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1e0;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 9223372036854775808;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 0x41;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 2147483648;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM DEFAULT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX a ON a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a;
< error
| 4092 42S02 Unknown VIEW: 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP IF EXISTS a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE IF EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT - 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a.a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a WAIT 1;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a NOWAIT;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW a AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW a AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW a AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW a AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a(a) AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a(a, a) AS SELECT a;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CHECK OPTION;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CASCADED CHECK OPTION;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH LOCAL CHECK OPTION;
< error
| 1146 42S02 Table 'omni_simulation.a' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT);
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY SYSTEM_TIME;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY SYSTEM_TIME;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD IF NOT EXISTS a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE a;
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a, ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MIN_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci AVG_ROW_LENGTH=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='s'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci PACK_KEYS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_AUTO_RECALC=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_PERSISTENT=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_SAMPLE_PAGES=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci PAGE_CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci DELAY_KEY_WRITE=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION();
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE DISK;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE MEMORY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CONNECTION='s'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci KEY_BLOCK_SIZE=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRANSACTIONAL 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci TRANSACTIONAL=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 's';
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SEQUENCE 1;
< error
| 4086 HY000 Sequence 'omni_simulation.a' table structure is invalid (Wrong number of columns)
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WITH SYSTEM VERSIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = a;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 1;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFAULT;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFINER;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a READ_ONLY = 's';
< error
| 1911 HY000 Unknown option 'READ_ONLY'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WRAPPER = 's';
< error
| 1911 HY000 Unknown option 'WRAPPER'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a WITH SYSTEM VERSIONING;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a VERSIONING;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT FIXED;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=COMPRESSED
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=REDUNDANT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=COMPACT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT PAGE;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a INT);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT, a INT);
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a));
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a));
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(PERIOD FOR SYSTEM_TIME(a, a));
< error
| 4125 HY000 Wrong parameters for `a`: missing 'WITH SYSTEM VERSIONING'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a TYPE BTREE(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE a TYPE BTREE(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a);
< error
| 1054 42S22 Unknown column 'a' in 'CHECK'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(a, a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CHECK(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL CHECK (`a`),
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT REFERENCES a;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT zzz REFERENCES a;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY;
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY KEY;
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY KEY;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a);
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW END;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT GENERATED ALWAYS AS ROW START;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) VIRTUAL;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a)
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) PERSISTENT;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) STORED;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE UNIQUE;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE KEY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL COMMENT 's',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT INVISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL INVISIBLE,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT(a);
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT a;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT + 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT 1,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT 1,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s');
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMERIC;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARBINARY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARCHAR;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`a`)),
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s');
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL(1, 1);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` double(1,1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1.5);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL NULL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AUTO_INCREMENT;
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT SERIAL DEFAULT VALUE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT VALUE;
< error
| 1054 42S22 Unknown column 'VALUE' in 'DEFAULT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT KEY;
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL ENABLE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) NOT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT a = 's';
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITH SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITHOUT SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR UNICODE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET ucs2 COLLATE ucs2_uca1400_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BYTE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII BINARY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY ASCII;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY zzz;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE BINARY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR BINARY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT PRIMARY KEY;
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a);
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a, a);
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH FULL;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH PARTIAL;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH SIMPLE;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT ON DELETE RESTRICT;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT ON UPDATE RESTRICT;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE CASCADE;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET NULL;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE NO ACTION;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET DEFAULT;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY zzz(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE zzz(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a;
< error
| 1091 42000 Can't DROP INDEX `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a;
< error
| 1091 42000 Can't DROP INDEX `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT zzz(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY zzz(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL zzz(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR zzz(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE VISIBLE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE VISIBLE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE VISIBLE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a USING BTREE ON a(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a TYPE BTREE ON a(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) TYPE BTREE;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's';
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) IGNORED;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) a = 's';
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING RTREE(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a IGNORED;
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a NOT IGNORED;
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a, a WITHOUT OVERLAPS);
< error
| 4156 HY000 Period `a` is not found in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1));
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a(a);
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s', 's');
< error
| 1291 HY000 Column 'a' has duplicated value 's' in ENUM
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE;
< error
| 1813 HY000 Tablespace for table 'omni_simulation/a' exists. Please DISCARD the tablespace before IMPORT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME PARTITIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a OPTIMIZE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECK PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REPAIR PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COALESCE PARTITION 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT PARTITION a TO TABLE a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(a);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION PARTITIONS 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(a);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS a CHECK(a);
< error
| 1054 42S22 Unknown column 'a' in 'CHECK'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a INT;
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a INT;
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a;
< error
| 1091 42000 Can't DROP COLUMN `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a;
< error
| 1091 42000 Can't DROP CONSTRAINT `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a;
< error
| 1091 42000 Can't DROP FOREIGN KEY `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT a;
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT;
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a;
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a;
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4 COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SYSTEM VERSIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP VERSIONING;
< error
| 1091 42000 Can't DROP COLUMN `VERSIONING`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR SYSTEM_TIME;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR a;
< error
| 1091 42000 Can't DROP PERIOD `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT ALGORITHM DEFAULT;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT LOCK DEFAULT;
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a;
< error
| 1091 42000 Can't DROP COLUMN `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE IGNORE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ONLINE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AFTER a;
< error
| 1054 42S22 Unknown column 'a' in 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT FIRST;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME = a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a, a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 9223372036854775808;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 2147483648;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=2147483648
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1.5;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1e0;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 9223372036854775808;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 0x41;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 2147483648;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX a ON a;
< error
| 1091 42000 Can't DROP INDEX `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a;
< error
| 4092 42S02 Unknown VIEW: 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP IF EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE IF EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT - 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT -1,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a.a;
< error
| 1091 42000 Can't DROP COLUMN `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
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
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a WAIT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a NOWAIT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a(a) AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a(a, a) AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CASCADED CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH LOCAL CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a);
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT);
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY SYSTEM_TIME;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY SYSTEM_TIME;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD IF NOT EXISTS a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE a;
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a, ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MIN_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci AVG_ROW_LENGTH=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='s'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci PACK_KEYS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_AUTO_RECALC=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_PERSISTENT=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_SAMPLE_PAGES=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci PAGE_CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci DELAY_KEY_WRITE=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a UNION();
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE DISK;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE MEMORY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CONNECTION='s'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci KEY_BLOCK_SIZE=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRANSACTIONAL 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci TRANSACTIONAL=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 's';
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SEQUENCE 1;
< error
| 4086 HY000 Sequence 'omni_simulation.a' table structure is invalid (Wrong number of columns)
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WITH SYSTEM VERSIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = a;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 1;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFAULT;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFINER;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a READ_ONLY = 's';
< error
| 1911 HY000 Unknown option 'READ_ONLY'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WRAPPER = 's';
< error
| 1911 HY000 Unknown option 'WRAPPER'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a WITH SYSTEM VERSIONING;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a VERSIONING;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT FIXED;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=COMPRESSED
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=REDUNDANT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=COMPACT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT PAGE;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a INT);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT, a INT);
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a));
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a));
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a(PERIOD FOR SYSTEM_TIME(a, a));
< error
| 4125 HY000 Wrong parameters for `a`: missing 'WITH SYSTEM VERSIONING'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a TYPE BTREE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a);
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE a TYPE BTREE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) USING BTREE,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(a, a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CHECK(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT REFERENCES a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT zzz REFERENCES a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY KEY;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW END;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT GENERATED ALWAYS AS ROW START;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) VIRTUAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) PERSISTENT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) STORED;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE UNIQUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COMMENT 's';
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT INVISIBLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT + 1;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT 1;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s');
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMERIC;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARBINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARCHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s');
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL(1, 1);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1.5);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AUTO_INCREMENT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT SERIAL DEFAULT VALUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT VALUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COLLATE utf8mb4_bin;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL ENABLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT a = 's';
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITH SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITHOUT SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR CHARSET utf8mb4;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR UNICODE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BYTE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY ASCII;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE DEFAULT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE utf8mb4_bin;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE DEFAULT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY zzz;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT PRIMARY KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a, a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH FULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH PARTIAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH SIMPLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE CASCADE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE NO ACTION;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET DEFAULT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY zzz(a);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE zzz(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a;
< error
| 1091 42000 Can't DROP INDEX `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a;
< error
| 1091 42000 Can't DROP INDEX `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT zzz(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY zzz(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL zzz(a);
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR(a);
< error
| 1210 HY000 Incorrect arguments to VECTOR INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR zzz(a);
< error
| 1210 HY000 Incorrect arguments to VECTOR INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE;
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE;
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE VISIBLE;
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE VISIBLE;
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `zzz` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `BTREE` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a USING BTREE ON a(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a TYPE BTREE ON a(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) TYPE BTREE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) KEY_BLOCK_SIZE=1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) COMMENT 's'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) IGNORED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) a = 's';
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING RTREE(a);
< error
| 1210 HY000 Incorrect arguments to RTREE INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a IGNORED;
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a NOT IGNORED;
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a, a WITHOUT OVERLAPS);
< error
| 4156 HY000 Period `a` is not found in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1));
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s', 's');
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE;
< error
| 1813 HY000 Tablespace for table 'omni_simulation/a' exists. Please DISCARD the tablespace before IMPORT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME PARTITIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a OPTIMIZE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECK PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REPAIR PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COALESCE PARTITION 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT PARTITION a TO TABLE a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(a);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION PARTITIONS 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(a);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS a CHECK(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   CONSTRAINT `a` CHECK (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a;
< error
| 1091 42000 Can't DROP CONSTRAINT `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a;
< error
| 1091 42000 Can't DROP FOREIGN KEY `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT a;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a;
< error
| 1176 42000 Key 'a' doesn't exist in table 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4 COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SYSTEM VERSIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP VERSIONING;
< error
| 1091 42000 Can't DROP COLUMN `VERSIONING`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR SYSTEM_TIME;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR a;
< error
| 1091 42000 Can't DROP PERIOD `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE IGNORE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ONLINE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AFTER a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT FIRST;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME = a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a, a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 9223372036854775808;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 2147483648;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=2147483648
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1.5;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1e0;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 9223372036854775808;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 0x41;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a a = 2147483648;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX a ON a;
< error
| 1091 42000 Can't DROP INDEX `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a;
< error
| 4092 42S02 Unknown VIEW: 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP IF EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE IF EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT - 1;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a.a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a WAIT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a NOWAIT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a(a) AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a(a, a) AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CASCADED CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH LOCAL CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| RENAME TABLE a TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| TRUNCATE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a);
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX a ON a(a);
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX a ON a(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX a ON a(a);
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT);
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a zzz;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a PARTITION BY SYSTEM_TIME;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY SYSTEM_TIME;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a AS SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD IF NOT EXISTS a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE OR REPLACE TABLE a;
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a ENGINE a, ENGINE a;
< error
| 1286 42000 Unknown storage engine 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MIN_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MIN_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AVG_ROW_LENGTH 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci AVG_ROW_LENGTH=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PASSWORD 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='s'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci PACK_KEYS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PACK_KEYS DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_AUTO_RECALC=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_AUTO_RECALC DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_PERSISTENT=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_PERSISTENT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci STATS_SAMPLE_PAGES=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STATS_SAMPLE_PAGES DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TABLE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci PAGE_CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DELAY_KEY_WRITE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci DELAY_KEY_WRITE=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a UNION();
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD NO;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DATA DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INDEX DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE DISK;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a STORAGE MEMORY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CONNECTION 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CONNECTION='s'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci KEY_BLOCK_SIZE=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TRANSACTIONAL 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci TRANSACTIONAL=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = 's';
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a SEQUENCE 1;
< error
| 4086 HY000 Sequence 'omni_simulation.a' table structure is invalid (Wrong number of columns)
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a WITH SYSTEM VERSIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = a;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = 1;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFAULT;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = DEFINER;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a READ_ONLY = 's';
< error
| 1911 HY000 Unknown option 'READ_ONLY'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a WRAPPER = 's';
< error
| 1911 HY000 Unknown option 'WRAPPER'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a WITH SYSTEM VERSIONING;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT a VERSIONING;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT FIXED;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT DYNAMIC;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPRESSED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=COMPRESSED
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT REDUNDANT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=REDUNDANT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT COMPACT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=COMPACT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT PAGE;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD FIRST;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a INSERT_METHOD LAST;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD(a INT);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(a INT, a INT);
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(KEY(a));
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(CHECK(a));
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a(PERIOD FOR SYSTEM_TIME(a, a));
< error
| 4125 HY000 Wrong parameters for `a`: missing 'WITH SYSTEM VERSIONING'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a TYPE BTREE(a);
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a);
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a_2` (`a`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE a TYPE BTREE(a);
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a;
< error
| 1005 HY000 Can't create table `omni_simulation`.`a` (errno: 150 "Foreign key constraint is incorrectly formed")
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(a, a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CHECK(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT REFERENCES a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT CONSTRAINT zzz REFERENCES a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL KEY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY KEY;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW START KEY
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS ROW END;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT GENERATED ALWAYS AS ROW START;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) VIRTUAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) PERSISTENT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) STORED;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE UNIQUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AS(a) UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT UNIQUE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COMMENT 's';
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT INVISIBLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT + 1;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT 1;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a YEAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYBLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYINT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s');
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BOOLEAN;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DECIMAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NUMERIC;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FIXED;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIME;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TIMESTAMP;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DATETIME;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMBLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMINT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGBLOB;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARBINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG VARCHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TINYTEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a TEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a MEDIUMTEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONGTEXT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a LONG;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a JSON;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SET('s');
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NATIONAL CHAR
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE PRECISION;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a DOUBLE
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a FLOAT(1, 1);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a REAL(1, 1);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1.5);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT(1);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AUTO_INCREMENT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT SERIAL DEFAULT VALUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT VALUE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT COLLATE utf8mb4_bin;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL ENABLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT NOT NULL
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT a = 's';
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITH SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT WITHOUT SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR CHARSET utf8mb4;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR UNICODE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BYTE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY ASCII;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BINARY
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE DEFAULT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz;
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE utf8mb4_bin;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE DEFAULT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY zzz;
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'latin1'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a NCHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT PRIMARY KEY;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a(a, a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH FULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH PARTIAL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a MATCH SIMPLE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE RESTRICT ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON DELETE RESTRICT ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE CASCADE;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET NULL;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE NO ACTION;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT REFERENCES a ON UPDATE SET DEFAULT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY zzz(a);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE zzz(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`a`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP KEY a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT zzz(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT KEY zzz(a);
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL zzz(a);
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR(a);
< error
| 1210 HY000 Incorrect arguments to VECTOR INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD VECTOR zzz(a);
< error
| 1210 HY000 Incorrect arguments to VECTOR INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE;
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE;
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE VISIBLE;
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a) VISIBLE
< error
| 1210 HY000 Incorrect arguments to SPATIAL INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE VISIBLE;
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) VISIBLE
< error
| 1283 HY000 Column 'a' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY zzz(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `zzz` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING BTREE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY BTREE(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `BTREE` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a USING BTREE ON a(a);
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a TYPE BTREE ON a(a);
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) USING BTREE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) TYPE BTREE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`) KEY_BLOCK_SIZE=1
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`) COMMENT 's'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) IGNORED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) a = 's';
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz;
< error
| 1128 HY000 Function 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING RTREE(a);
< error
| 1210 HY000 Incorrect arguments to RTREE INDEX
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING HASH(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   KEY `a_2` (`a`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a IGNORED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER KEY a NOT IGNORED;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a, a);
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a, a WITHOUT OVERLAPS);
< error
| 4156 HY000 Period `a` is not found in table
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(1));
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY a(a);
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a ENUM('s', 's');
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW a AS SELECT a zzz;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a zzz;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE;
< error
| 1813 HY000 Tablespace for table 'omni_simulation/a' exists. Please DISCARD the tablespace before IMPORT
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME PARTITIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a OPTIMIZE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHECK PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REPAIR PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a COALESCE PARTITION 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT PARTITION a TO TABLE a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION(a);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION PARTITIONS 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION a INTO(a);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION a, a;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a FORCE, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS a CHECK(a);
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`),
|   CONSTRAINT `a` CHECK (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CHANGE a a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY a INT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT a;
< error
| 1091 42000 Can't DROP CONSTRAINT `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY a;
< error
| 1091 42000 Can't DROP FOREIGN KEY `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DISABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ENABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a SET DEFAULT a;
< error
| 4029 01000 Expression for field `a` is referring to uninitialized field `a`
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER a DROP DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY a TO a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET utf8mb4 COLLATE utf8mb4_bin;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SYSTEM VERSIONING;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP SYSTEM VERSIONING;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP VERSIONING;
< error
| 1091 42000 Can't DROP COLUMN `VERSIONING`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR SYSTEM_TIME;
< error
| 4124 HY000 Table `a` is not system-versioned
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PERIOD FOR a;
< error
| 1091 42000 Can't DROP PERIOD `a`; check that it exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT;
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT;
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) LOCK DEFAULT ALGORITHM DEFAULT;
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE INDEX a ON a(a) ALGORITHM DEFAULT LOCK DEFAULT;
< error
| 1061 42000 Duplicate key name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER IGNORE IGNORE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER ONLINE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP TABLE a CASCADE;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT AFTER a;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT FIRST;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME TO a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME = a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME AS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| RENAME TABLE a TO a, a TO a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ORDER BY a, a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 9223372036854775808;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 2147483648;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=2147483648
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1.5;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a MAX_ROWS 1e0;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = 9223372036854775808;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = 0x41;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a a = 2147483648;
< error
| 1911 HY000 Unknown option 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a PAGE_CHECKSUM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP INDEX a ON a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| DROP VIEW a;
< error
| 4092 42S02 Unknown VIEW: 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP IF EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE IF EXISTS a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| TRUNCATE TABLE a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a REPLACE SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| CREATE TABLE a IGNORE SELECT a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a INT DEFAULT - 1;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a DROP a.a;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a BIGINT;
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| TRUNCATE a WAIT 1;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| TRUNCATE a NOWAIT;
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| TRUNCATE a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `a` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW a AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a(a) AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a(a, a) AS SELECT a;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH CASCADED CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY a (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT a WITH LOCAL CHECK OPTION;
< error
| 1347 HY000 'omni_simulation.a' is not of type 'VIEW'
