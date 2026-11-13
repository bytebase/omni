-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSimulation default#1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4;
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4;
< error
| 1051 42S02 Unknown table 'omni_t_simulation.utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW utf8mb4;
< error
| 1051 42S02 Unknown table 'omni_t_simulation.utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `utf8mb4` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 LIKE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(LIKE utf8mb4);
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 's';
< error
| 3708 SR006 Missing mandatory attribute DEFINITION.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DEFINITION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 ORGANIZATION 's' IDENTIFIED BY 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DESCRIPTION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 PARTITION BY KEY();
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY();
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY HASH(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE COLUMNS(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST COLUMNS(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LINEAR KEY();
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY ALGORITHM = 1();
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() PARTITIONS 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY KEY(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(zzz);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4, utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4) SUBPARTITIONS 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY()(PARTITION utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4, PARTITION utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN MAXVALUE);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAXVALUE));
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR),(MAXVALUE)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(utf8mb4));
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4));
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4, SUBPARTITION utf8mb4));
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4 TABLESPACE utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 ENGINE utf8mb4);
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 NODEGROUP 1);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MAX_ROWS 1);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MIN_ROWS 1);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 DATA DIRECTORY 's');
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 INDEX DIRECTORY 's');
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 COMMENT 's');
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE NULL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MIN_ROWS 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AVG_ROW_LENGTH 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PASSWORD 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMMENT 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMPRESSION 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENCRYPTION 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTO_INCREMENT 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PACK_KEYS 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_AUTO_RECALC 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_PERSISTENT 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECKSUM 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLE_CHECKSUM 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DELAY_KEY_WRITE 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 UNION();
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHARSET utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD NO;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DATA DIRECTORY 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INDEX DIRECTORY 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLESPACE utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE DISK;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE MEMORY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONNECTION 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 KEY_BLOCK_SIZE 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 START TRANSACTION;
< error
| 3979 HY000 START TRANSACTION clause cannot be used with ALTER TABLE command.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTOEXTEND_SIZE 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT FIXED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DYNAMIC;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPRESSED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT REDUNDANT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPACT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD FIRST;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD LAST;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT, utf8mb4 INT);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(KEY(utf8mb4));
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4) REFERENCES utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT UNIQUE(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT zzz UNIQUE(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4) ENFORCED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4)
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT GENERATED ALWAYS AS(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VIRTUAL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4)
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) STORED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOLEAN;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NVARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARBINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 YEAR;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIMESTAMP;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATETIME;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYBLOB;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BLOB;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMBLOB;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGBLOB;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARBINARY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARCHAR;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYTEXT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TEXT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMTEXT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGTEXT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s');
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SET('s');
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SERIAL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 JSON;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRYCOLLECTION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POINT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOINT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LINESTRING;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTILINESTRING;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POLYGON;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOLYGON;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG CHAR VARYING;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYINT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SMALLINT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMINT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIGINT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE PRECISION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DECIMAL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NUMERIC;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FIXED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(1, 1);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL(1, 1);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME(1);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED SIGNED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNSIGNED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ZEROFILL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(18446744073709551615);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1.5);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL NULL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT NULL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT SECONDARY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT CURRENT_TIMESTAMP;
< error
| 1067 42000 Invalid default value for 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ON UPDATE CURRENT_TIMESTAMP;
< error
| 1294 HY000 Invalid ON UPDATE clause for 'utf8mb4' column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AUTO_INCREMENT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SERIAL DEFAULT VALUE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT KEY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE KEY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COMMENT 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SRID 1;
< error
| 1221 HY000 Incorrect usage of SRID and non-geometry column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT CHECK(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT FIXED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DYNAMIC;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DISK;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE MEMORY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY ASCII;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII BINARY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE BINARY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY UNICODE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BYTE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR CHARSET utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY CHARSET utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR BINARY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT PRIMARY KEY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4, utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH FULL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH PARTIAL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH SIMPLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT ON DELETE RESTRICT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT ON UPDATE RESTRICT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE CASCADE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET NULL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE NO ACTION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY zzz(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE zzz(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP KEY utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP INDEX utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT zzz(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY zzz(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEYRING(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4)
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4)
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4)
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) USING BTREE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) KEY_BLOCK_SIZE 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) COMMENT 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY zzz(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING BTREE(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY BTREE(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4 TYPE BTREE(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 USING BTREE ON utf8mb4(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) TYPE BTREE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT INVISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING RTREE(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING HASH(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4, utf8mb4) REFERENCES utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4(1));
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((1));
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4, utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((utf8mb4));
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s', 's');
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_LOAD;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE SECONDARY_LOAD;
< error
| 1146 42S02 Table 'omni_t_simulation.SECONDARY_LOAD' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REMOVE PARTITIONING;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME PARTITIONING;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, SECONDARY_LOAD;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, FORCE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD TABLESPACE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT TABLESPACE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION PARTITIONS 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION ALL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 OPTIMIZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ANALYZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECK PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REPAIR PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COALESCE PARTITION 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TRUNCATE PARTITION ALL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION utf8mb4 INTO(PARTITION utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD PARTITION ALL TABLESPACE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT PARTITION ALL TABLESPACE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_UNLOAD;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITH VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITH VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, FORCE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD(utf8mb4 INT);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHANGE utf8mb4 utf8mb4 INT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MODIFY utf8mb4 INT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP FOREIGN KEY utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PRIMARY KEY;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CHECK utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CONSTRAINT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISABLE KEYS;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENABLE KEYS;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT 's';
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT(utf8mb4);
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 DROP DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER INDEX utf8mb4 VISIBLE;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CHECK utf8mb4 ENFORCED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CONSTRAINT utf8mb4 ENFORCED;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME KEY utf8mb4 TO utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME COLUMN utf8mb4 TO utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 LOCK DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT ALGORITHM DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT LOCK DEFAULT;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP COLUMN utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 RESTRICT;
< error
| 1051 42S02 Unknown table 'omni_t_simulation.utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4
< error
| 1051 42S02 Unknown table 'omni_t_simulation.utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 CASCADE;
< error
| 1051 42S02 Unknown table 'omni_t_simulation.utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AFTER utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT FIRST;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME TO utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME = utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME AS utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4, utf8mb4 TO utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4, utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4) AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4, utf8mb4) AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4, utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 18446744073709551615;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 4294967296;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1.5;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1e0;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 IGNORE SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT + 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT - 1;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT NULL;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CHECK OPTION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CASCADED CHECK OPTION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH LOCAL CHECK OPTION;
< error
| 1146 42S02 Table 'omni_t_simulation.utf8mb4' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT);
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 LIKE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(LIKE utf8mb4);
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 's';
< error
| 3708 SR006 Missing mandatory attribute DEFINITION.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DEFINITION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 ORGANIZATION 's' IDENTIFIED BY 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DESCRIPTION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 PARTITION BY KEY();
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY HASH(utf8mb4);
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE COLUMNS(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST COLUMNS(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LINEAR KEY();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LINEAR KEY () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY ALGORITHM = 1();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() PARTITIONS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ()
| PARTITIONS 1 */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4);
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY KEY(utf8mb4);
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(zzz);
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4);
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4, utf8mb4);
< error
| 1652 HY000 Duplicate partition field name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4) SUBPARTITIONS 1;
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY()(PARTITION utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ()
| (PARTITION utf8mb4 ENGINE = InnoDB) */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4, PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN MAXVALUE);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAXVALUE));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR),(MAXVALUE)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(utf8mb4));
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4, SUBPARTITION utf8mb4));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4 TABLESPACE utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 ENGINE utf8mb4);
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 NODEGROUP 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MAX_ROWS 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MIN_ROWS 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 DATA DIRECTORY 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 INDEX DIRECTORY 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 COMMENT 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci SECONDARY_ENGINE=utf8mb4
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MIN_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MIN_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AVG_ROW_LENGTH 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci AVG_ROW_LENGTH=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PASSWORD 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='s'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMPRESSION 's';
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'COMPRESSION'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENCRYPTION 's';
< error
| 3184 HY000 Invalid encryption option.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTO_INCREMENT 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PACK_KEYS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci PACK_KEYS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_AUTO_RECALC 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_AUTO_RECALC=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_PERSISTENT 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_PERSISTENT=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_SAMPLE_PAGES=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DELAY_KEY_WRITE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci DELAY_KEY_WRITE=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 UNION();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD NO;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DATA DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INDEX DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLESPACE utf8mb4;
< error
| 3510 HY000 Tablespace utf8mb4 doesn't exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE DISK;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) /*!50100 STORAGE DISK */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE MEMORY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) /*!50100 STORAGE MEMORY */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONNECTION 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CONNECTION='s'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci KEY_BLOCK_SIZE=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 START TRANSACTION;
< error
| 3979 HY000 START TRANSACTION clause cannot be used with ALTER TABLE command.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTOEXTEND_SIZE 1;
< error
| 4025 HY000 AUTOEXTEND_SIZE value should be between 4M and 4096M
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT FIXED;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DYNAMIC;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPRESSED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT REDUNDANT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=REDUNDANT
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPACT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPACT
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD FIRST;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD LAST;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT, utf8mb4 INT);
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(KEY(utf8mb4));
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4) REFERENCES utf8mb4;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4);
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT UNIQUE(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT zzz UNIQUE(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4) ENFORCED;
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4);
< error
| 3107 HY000 Generated column can refer only to generated columns defined prior to it.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT GENERATED ALWAYS AS(utf8mb4);
< error
| 3107 HY000 Generated column can refer only to generated columns defined prior to it.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VIRTUAL;
< error
| 3107 HY000 Generated column can refer only to generated columns defined prior to it.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4)
< error
| 3107 HY000 Generated column can refer only to generated columns defined prior to it.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VISIBLE;
< error
| 3107 HY000 Generated column can refer only to generated columns defined prior to it.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) STORED;
< error
| 3107 HY000 Generated column can refer only to generated columns defined prior to it.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` float DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` bit(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOLEAN;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NVARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARBINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 YEAR;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` year DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` date DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` time DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIMESTAMP;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` timestamp NULL DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATETIME;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` datetime DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYBLOB;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` tinyblob,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BLOB;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` blob,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` geometry DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMBLOB;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumblob,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGBLOB;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` longblob,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARBINARY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumblob,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumtext,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARCHAR;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumtext,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYTEXT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` tinytext,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TEXT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` text,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMTEXT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumtext,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGTEXT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` longtext,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s');
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` enum('s') DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SET('s');
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` set('s') DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumtext,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SERIAL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 JSON;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` json DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRYCOLLECTION;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` geomcollection DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POINT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` point DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOINT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` multipoint DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LINESTRING;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` linestring DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTILINESTRING;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` multilinestring DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POLYGON;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` polygon DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOLYGON;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` multipolygon DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG CHAR VARYING;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumtext,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYINT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SMALLINT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` smallint DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMINT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` mediumint DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIGINT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` bigint DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE PRECISION;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` double DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DECIMAL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NUMERIC;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FIXED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(1, 1);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` float(1,1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL(1, 1);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` double(1,1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME(1);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` time(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED SIGNED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNSIGNED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int unsigned DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ZEROFILL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int(10) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(18446744073709551615);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1.5);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int NOT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT SECONDARY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int NOT SECONDARY DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT CURRENT_TIMESTAMP;
< error
| 1067 42000 Invalid default value for 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT(utf8mb4);
< error
| 3767 HY000 Default value expression of column 'utf8mb4' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ON UPDATE CURRENT_TIMESTAMP;
< error
| 1294 HY000 Invalid ON UPDATE clause for 'utf8mb4' column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AUTO_INCREMENT;
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SERIAL DEFAULT VALUE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT KEY;
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE KEY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL COMMENT 's',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SRID 1;
< error
| 1221 HY000 Incorrect usage of SRID and non-geometry column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT CHECK(utf8mb4);
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT FIXED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int /*!50606 COLUMN_FORMAT FIXED */ DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DYNAMIC;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int /*!50606 COLUMN_FORMAT DYNAMIC */ DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DISK;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int /*!50606 STORAGE DISK */ DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE MEMORY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int /*!50606 STORAGE MEMORY */ DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 's';
< error
| 1067 42000 Invalid default value for 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY ASCII;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII BINARY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET ucs2 COLLATE ucs2_general_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE BINARY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET ucs2 COLLATE ucs2_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET ucs2 COLLATE ucs2_general_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY UNICODE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET ucs2 COLLATE ucs2_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BYTE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR BINARY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT PRIMARY KEY;
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4, utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH FULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH PARTIAL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH SIMPLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT ON DELETE RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT ON UPDATE RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE CASCADE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE NO ACTION;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY(utf8mb4);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY zzz(utf8mb4);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE zzz(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP KEY utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP INDEX utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT zzz(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY zzz(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEYRING(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4)
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE VISIBLE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4)
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE VISIBLE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4)
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE VISIBLE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) USING BTREE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) KEY_BLOCK_SIZE 1;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) COMMENT 's';
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY zzz(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING BTREE(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY BTREE(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4 TYPE BTREE(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 USING BTREE ON utf8mb4(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) TYPE BTREE;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT INVISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL /*!80023 INVISIBLE */,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING RTREE(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING HASH(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4, utf8mb4) REFERENCES utf8mb4;
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4(1));
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((1));
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `functional_index` ((1))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4, utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((utf8mb4));
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4(utf8mb4);
< error
| 1072 42000 Key column 'utf8mb4' doesn't exist in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s', 's');
< error
| 1291 HY000 Column 'utf8mb4' has duplicated value 's' in ENUM
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_LOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE SECONDARY_LOAD;
< error
| 1146 42S02 Table 'omni_t_simulation.SECONDARY_LOAD' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME PARTITIONING;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, SECONDARY_LOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD TABLESPACE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT TABLESPACE;
< error
| 1813 HY000 Tablespace 'omni_t_simulation/utf8mb4' exists.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION PARTITIONS 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 OPTIMIZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ANALYZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECK PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REPAIR PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COALESCE PARTITION 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TRUNCATE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION utf8mb4 INTO(PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD PARTITION ALL TABLESPACE;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT PARTITION ALL TABLESPACE;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_UNLOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITH VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITH VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD(utf8mb4 INT);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHANGE utf8mb4 utf8mb4 INT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MODIFY utf8mb4 INT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP FOREIGN KEY utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CHECK utf8mb4;
< error
| 3821 HY000 Check constraint 'utf8mb4' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CONSTRAINT utf8mb4;
< error
| 3940 HY000 Constraint 'utf8mb4' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT 's';
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT(utf8mb4);
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 DROP DEFAULT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET VISIBLE;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER INDEX utf8mb4 VISIBLE;
< error
| 1176 42000 Key 'utf8mb4' doesn't exist in table 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CHECK utf8mb4 ENFORCED;
< error
| 3821 HY000 Check constraint 'utf8mb4' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CONSTRAINT utf8mb4 ENFORCED;
< error
| 3940 HY000 Constraint 'utf8mb4' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME KEY utf8mb4 TO utf8mb4;
< error
| 1176 42000 Key 'utf8mb4' doesn't exist in table 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME COLUMN utf8mb4 TO utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT ALGORITHM DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT LOCK DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP COLUMN utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 CASCADE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AFTER utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT FIRST;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `utf8mb4` int DEFAULT NULL,
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME TO utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME = utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME AS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4, utf8mb4 TO utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4, utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4) AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4, utf8mb4) AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4, utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 18446744073709551615;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 4294967296;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1.5;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1e0;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 IGNORE SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT + 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT '1',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT '1',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT - 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT '-1',
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CASCADED CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH LOCAL CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT);
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 LIKE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(LIKE utf8mb4);
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 's';
< error
| 3708 SR006 Missing mandatory attribute DEFINITION.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DEFINITION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 ORGANIZATION 's' IDENTIFIED BY 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DESCRIPTION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 PARTITION BY KEY();
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY HASH(utf8mb4);
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE COLUMNS(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST COLUMNS(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LINEAR KEY();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LINEAR KEY () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY ALGORITHM = 1();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() PARTITIONS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ()
| PARTITIONS 1 */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4);
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY KEY(utf8mb4);
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(zzz);
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4);
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4, utf8mb4);
< error
| 1652 HY000 Duplicate partition field name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4) SUBPARTITIONS 1;
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY()(PARTITION utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ()
| (PARTITION utf8mb4 ENGINE = InnoDB) */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4, PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN MAXVALUE);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAXVALUE));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR),(MAXVALUE)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(utf8mb4));
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4, SUBPARTITION utf8mb4));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4 TABLESPACE utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 ENGINE utf8mb4);
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 NODEGROUP 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MAX_ROWS 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MIN_ROWS 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 DATA DIRECTORY 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 INDEX DIRECTORY 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 COMMENT 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci SECONDARY_ENGINE=utf8mb4
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MIN_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MIN_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AVG_ROW_LENGTH 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci AVG_ROW_LENGTH=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PASSWORD 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='s'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMPRESSION 's';
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'COMPRESSION'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENCRYPTION 's';
< error
| 3184 HY000 Invalid encryption option.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTO_INCREMENT 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PACK_KEYS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci PACK_KEYS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_AUTO_RECALC 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_AUTO_RECALC=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_PERSISTENT 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_PERSISTENT=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_SAMPLE_PAGES=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DELAY_KEY_WRITE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci DELAY_KEY_WRITE=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 UNION();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD NO;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DATA DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INDEX DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLESPACE utf8mb4;
< error
| 3510 HY000 Tablespace utf8mb4 doesn't exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE DISK;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) /*!50100 STORAGE DISK */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE MEMORY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) /*!50100 STORAGE MEMORY */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONNECTION 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CONNECTION='s'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci KEY_BLOCK_SIZE=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 START TRANSACTION;
< error
| 3979 HY000 START TRANSACTION clause cannot be used with ALTER TABLE command.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTOEXTEND_SIZE 1;
< error
| 4025 HY000 AUTOEXTEND_SIZE value should be between 4M and 4096M
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT FIXED;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DYNAMIC;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPRESSED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT REDUNDANT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=REDUNDANT
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPACT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPACT
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD FIRST;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD LAST;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT, utf8mb4 INT);
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(KEY(utf8mb4));
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4);
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4) REFERENCES utf8mb4;
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4);
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT UNIQUE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT zzz UNIQUE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4) ENFORCED;
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT GENERATED ALWAYS AS(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VIRTUAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4)
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VISIBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) STORED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOLEAN;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NVARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARBINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 YEAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIMESTAMP;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATETIME;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYBLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMBLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGBLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARBINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARCHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYTEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMTEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGTEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s');
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SET('s');
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SERIAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 JSON;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRYCOLLECTION;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LINESTRING;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTILINESTRING;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POLYGON;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOLYGON;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG CHAR VARYING;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SMALLINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIGINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE PRECISION;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DECIMAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NUMERIC;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FIXED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(1, 1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL(1, 1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME(1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED SIGNED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNSIGNED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ZEROFILL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(18446744073709551615);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1.5);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT SECONDARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT CURRENT_TIMESTAMP;
< error
| 1067 42000 Invalid default value for 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ON UPDATE CURRENT_TIMESTAMP;
< error
| 1294 HY000 Invalid ON UPDATE clause for 'utf8mb4' column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AUTO_INCREMENT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SERIAL DEFAULT VALUE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT KEY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE KEY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COMMENT 's';
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DEFAULT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DEFAULT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SRID 1;
< error
| 1221 HY000 Incorrect usage of SRID and non-geometry column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT CHECK(utf8mb4);
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT VISIBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT FIXED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DYNAMIC;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DISK;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE MEMORY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 's';
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY ASCII;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY UNICODE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BYTE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR CHARSET utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY CHARSET utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT PRIMARY KEY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4, utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH FULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH PARTIAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH SIMPLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE CASCADE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE NO ACTION;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET DEFAULT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY(utf8mb4);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY zzz(utf8mb4);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE zzz(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP KEY utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP INDEX utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT zzz(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY zzz(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEYRING(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4` (`utf8mb4`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE;
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4)
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE VISIBLE;
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE;
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE VISIBLE;
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) USING BTREE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`) COMMENT 's'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY zzz(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `zzz` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING BTREE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY BTREE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `BTREE` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4 TYPE BTREE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 USING BTREE ON utf8mb4(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) TYPE BTREE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT INVISIBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING RTREE(utf8mb4);
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING HASH(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4, utf8mb4) REFERENCES utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4(1));
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((1));
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `functional_index` ((1))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4, utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((utf8mb4));
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `utf8mb4` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s', 's');
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_LOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE SECONDARY_LOAD;
< error
| 1146 42S02 Table 'omni_t_simulation.SECONDARY_LOAD' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME PARTITIONING;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, SECONDARY_LOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD TABLESPACE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT TABLESPACE;
< error
| 1813 HY000 Tablespace 'omni_t_simulation/utf8mb4' exists.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION PARTITIONS 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 OPTIMIZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ANALYZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECK PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REPAIR PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COALESCE PARTITION 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TRUNCATE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION utf8mb4 INTO(PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD PARTITION ALL TABLESPACE;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT PARTITION ALL TABLESPACE;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_UNLOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITH VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITH VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD(utf8mb4 INT);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHANGE utf8mb4 utf8mb4 INT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MODIFY utf8mb4 INT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP FOREIGN KEY utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CHECK utf8mb4;
< error
| 3821 HY000 Check constraint 'utf8mb4' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CONSTRAINT utf8mb4;
< error
| 3940 HY000 Constraint 'utf8mb4' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT 's';
< error
| 1067 42000 Invalid default value for 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT(utf8mb4);
< error
| 3767 HY000 Default value expression of column 'utf8mb4' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 DROP DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER INDEX utf8mb4 VISIBLE;
< error
| 1176 42000 Key 'utf8mb4' doesn't exist in table 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CHECK utf8mb4 ENFORCED;
< error
| 3821 HY000 Check constraint 'utf8mb4' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CONSTRAINT utf8mb4 ENFORCED;
< error
| 3940 HY000 Constraint 'utf8mb4' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME KEY utf8mb4 TO utf8mb4;
< error
| 1176 42000 Key 'utf8mb4' doesn't exist in table 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME COLUMN utf8mb4 TO utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT ALGORITHM DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT LOCK DEFAULT;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP COLUMN utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 CASCADE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AFTER utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT FIRST;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME TO utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME = utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME AS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4, utf8mb4 TO utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4, utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4) AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4, utf8mb4) AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4, utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 18446744073709551615;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 4294967296;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1.5;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1e0;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 IGNORE SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT + 1;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 1;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT - 1;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CASCADED CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH LOCAL CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1061 42000 Duplicate key name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP VIEW utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| TRUNCATE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT);
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 LIKE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(LIKE utf8mb4);
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 's';
< error
| 3708 SR006 Missing mandatory attribute DEFINITION.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DEFINITION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 ORGANIZATION 's' IDENTIFIED BY 1;
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 1 DESCRIPTION 's';
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1061 42000 Duplicate key name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1061 42000 Duplicate key name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 PARTITION BY KEY();
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 REPLACE SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY HASH(utf8mb4);
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY RANGE COLUMNS(utf8mb4);
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LIST COLUMNS(utf8mb4);
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY LINEAR KEY();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY LINEAR KEY () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY ALGORITHM = 1();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY */ /*!50611 ALGORITHM = 1 */ /*!50100 () */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() PARTITIONS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ()
| PARTITIONS 1 */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4);
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY KEY(utf8mb4);
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(zzz);
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4);
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function (prefixed columns are not considered).
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY(utf8mb4, utf8mb4);
< error
| 1652 HY000 Duplicate partition field name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY() SUBPARTITION BY HASH(utf8mb4) SUBPARTITIONS 1;
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PARTITION BY KEY()(PARTITION utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY KEY ()
| (PARTITION utf8mb4 ENGINE = InnoDB) */
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4, PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN MAXVALUE);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAXVALUE));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN((MAX_CONNECTIONS_PER_HOUR),(MAXVALUE)));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE));
< error
| 1054 42S22 Unknown column 'MAX_CONNECTIONS_PER_HOUR' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 VALUES IN(utf8mb4));
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'partition function'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4(SUBPARTITION utf8mb4, SUBPARTITION utf8mb4));
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 TABLESPACE utf8mb4 TABLESPACE utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 ENGINE utf8mb4);
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 NODEGROUP 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MAX_ROWS 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 MIN_ROWS 1);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 DATA DIRECTORY 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 INDEX DIRECTORY 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION(PARTITION utf8mb4 COMMENT 's');
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 ENGINE utf8mb4 ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE NULL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_ENGINE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci SECONDARY_ENGINE=utf8mb4
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MIN_ROWS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MIN_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AVG_ROW_LENGTH 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci AVG_ROW_LENGTH=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PASSWORD 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='s'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COMPRESSION 's';
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'COMPRESSION'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENCRYPTION 's';
< error
| 3184 HY000 Invalid encryption option.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTO_INCREMENT 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 PACK_KEYS 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci PACK_KEYS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_AUTO_RECALC 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_AUTO_RECALC=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_PERSISTENT 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_PERSISTENT=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci STATS_SAMPLE_PAGES=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STATS_SAMPLE_PAGES DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLE_CHECKSUM 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CHECKSUM=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DELAY_KEY_WRITE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci DELAY_KEY_WRITE=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 UNION();
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD NO;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DATA DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INDEX DIRECTORY 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TABLESPACE utf8mb4;
< error
| 3510 HY000 Tablespace utf8mb4 doesn't exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE DISK;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) /*!50100 STORAGE DISK */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 STORAGE MEMORY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) /*!50100 STORAGE MEMORY */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONNECTION 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci CONNECTION='s'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci KEY_BLOCK_SIZE=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 START TRANSACTION;
< error
| 3979 HY000 START TRANSACTION clause cannot be used with ALTER TABLE command.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 AUTOEXTEND_SIZE 1;
< error
| 4025 HY000 AUTOEXTEND_SIZE value should be between 4M and 4096M
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT FIXED;
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT DYNAMIC;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPRESSED;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT REDUNDANT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=REDUNDANT
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ROW_FORMAT COMPACT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPACT
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD FIRST;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 INSERT_METHOD LAST;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(utf8mb4 INT, utf8mb4 INT);
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4(KEY(utf8mb4));
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4);
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4_2` (`utf8mb4`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4) REFERENCES utf8mb4;
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4);
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT UNIQUE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `utf8mb4_2` (`utf8mb4`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CONSTRAINT zzz UNIQUE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`utf8mb4`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4) ENFORCED;
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD CHECK(utf8mb4)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT GENERATED ALWAYS AS(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VIRTUAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4)
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) VISIBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AS(utf8mb4) STORED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BOOLEAN;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NVARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 VARBINARY(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 YEAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIMESTAMP;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DATETIME;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYBLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMBLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGBLOB;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARBINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG VARCHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYTEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMTEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONGTEXT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s');
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SET('s');
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SERIAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 JSON;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 GEOMETRYCOLLECTION;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LINESTRING;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTILINESTRING;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 POLYGON;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MULTIPOLYGON;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 LONG CHAR VARYING;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARCHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NATIONAL CHAR
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR VARYING(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TINYINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 SMALLINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 MEDIUMINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 BIGINT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE PRECISION;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DOUBLE
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 DECIMAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NUMERIC;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FIXED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 FLOAT(1, 1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 REAL(1, 1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 TIME(1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED SIGNED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SIGNED
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNSIGNED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ZEROFILL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(4294967296);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(18446744073709551615);
< error
| 1439 42000 Display width out of range for column 'utf8mb4' (max = 4294967295)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1.5);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT(1);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NULL
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT NOT SECONDARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT CURRENT_TIMESTAMP;
< error
| 1067 42000 Invalid default value for 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT ON UPDATE CURRENT_TIMESTAMP;
< error
| 1294 HY000 Invalid ON UPDATE clause for 'utf8mb4' column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AUTO_INCREMENT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SERIAL DEFAULT VALUE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT KEY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE KEY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT UNIQUE
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COMMENT 's';
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLLATE utf8mb4 zzz;
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DEFAULT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DEFAULT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT SRID 1;
< error
| 1221 HY000 Incorrect usage of SRID and non-geometry column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT CHECK(utf8mb4);
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'utf8mb4_chk_1'.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT VISIBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT FIXED;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT COLUMN_FORMAT DYNAMIC;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE DISK;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT STORAGE MEMORY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 's';
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY ASCII;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR ASCII
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR UNICODE
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY UNICODE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BYTE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR CHARSET utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 CHAR BINARY CHARSET utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 NCHAR BINARY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT PRIMARY KEY;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4(utf8mb4, utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH FULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH PARTIAL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 MATCH SIMPLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE RESTRICT ON DELETE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON DELETE RESTRICT ON UPDATE RESTRICT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE CASCADE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE NO ACTION;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT REFERENCES utf8mb4 ON UPDATE SET DEFAULT;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY(utf8mb4);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PRIMARY KEY zzz(utf8mb4);
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD UNIQUE zzz(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `zzz` (`utf8mb4`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP KEY utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP INDEX utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT zzz(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEY zzz(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT KEYRING(utf8mb4);
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX utf8mb4 ON utf8mb4(utf8mb4);
< error
| 1061 42000 Duplicate key name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE;
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4)
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE VISIBLE;
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) VISIBLE
< error
| 1283 HY000 Column 'utf8mb4' cannot be part of FULLTEXT index
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE;
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE VISIBLE;
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD SPATIAL(utf8mb4) VISIBLE
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) USING BTREE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) KEY_BLOCK_SIZE 1;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) COMMENT 's';
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`) COMMENT 's'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY zzz(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `zzz` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING BTREE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY BTREE(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `BTREE` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4 TYPE BTREE(utf8mb4);
< error
| 1061 42000 Duplicate key name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE INDEX utf8mb4 USING BTREE ON utf8mb4(utf8mb4);
< error
| 1061 42000 Duplicate key name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4) TYPE BTREE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT INVISIBLE;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING RTREE(utf8mb4);
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY USING HASH(utf8mb4);
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `utf8mb4_2` (`utf8mb4`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD FOREIGN KEY(utf8mb4, utf8mb4) REFERENCES utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4(1));
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((1));
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`),
|   KEY `functional_index` ((1))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY(utf8mb4, utf8mb4);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY((utf8mb4));
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD KEY utf8mb4(utf8mb4);
< error
| 1061 42000 Duplicate key name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 ENUM('s', 's');
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_LOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE SECONDARY_LOAD;
< error
| 1146 42S02 Table 'omni_t_simulation.SECONDARY_LOAD' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REMOVE PARTITIONING;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME PARTITIONING;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, SECONDARY_LOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD TABLESPACE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT TABLESPACE;
< error
| 1813 HY000 Tablespace 'omni_t_simulation/utf8mb4' exists.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD PARTITION PARTITIONS 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 OPTIMIZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ANALYZE PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHECK PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REPAIR PARTITION ALL;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 COALESCE PARTITION 1;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 TRUNCATE PARTITION ALL;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REORGANIZE PARTITION utf8mb4 INTO(PARTITION utf8mb4);
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISCARD PARTITION ALL TABLESPACE;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 IMPORT PARTITION ALL TABLESPACE;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 SECONDARY_UNLOAD;
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITH VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 EXCHANGE PARTITION utf8mb4 WITH TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1066 42000 Not unique table/alias: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITH VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 WITHOUT VALIDATION;
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 REBUILD PARTITION utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, FORCE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 FORCE, ENGINE utf8mb4;
< error
| 1286 42000 Unknown storage engine 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALGORITHM DEFAULT, ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD(utf8mb4 INT);
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CHANGE utf8mb4 utf8mb4 INT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MODIFY utf8mb4 INT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP FOREIGN KEY utf8mb4;
< error
| 1091 42000 Can't DROP 'utf8mb4'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PRIMARY KEY;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CHECK utf8mb4;
< error
| 3821 HY000 Check constraint 'utf8mb4' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP CONSTRAINT utf8mb4;
< error
| 3940 HY000 Constraint 'utf8mb4' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DISABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ENABLE KEYS;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT 's';
< error
| 1067 42000 Invalid default value for 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET DEFAULT(utf8mb4);
< error
| 3767 HY000 Default value expression of column 'utf8mb4' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 DROP DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER utf8mb4 SET VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER INDEX utf8mb4 VISIBLE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CHECK utf8mb4 ENFORCED;
< error
| 3821 HY000 Check constraint 'utf8mb4' is not found in the table.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ALTER CONSTRAINT utf8mb4 ENFORCED;
< error
| 3940 HY000 Constraint 'utf8mb4' does not exist.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME KEY utf8mb4 TO utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME COLUMN utf8mb4 TO utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 CONVERT TO CHARSET DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 LOCK DEFAULT ALGORITHM DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP INDEX utf8mb4 ON utf8mb4 ALGORITHM DEFAULT LOCK DEFAULT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP COLUMN utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 RESTRICT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| DROP TABLE utf8mb4 CASCADE;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT AFTER utf8mb4;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT FIRST;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME TO utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME = utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 RENAME AS utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| RENAME TABLE utf8mb4 TO utf8mb4, utf8mb4 TO utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 DROP PARTITION utf8mb4, utf8mb4;
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4) AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4(utf8mb4, utf8mb4) AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ORDER BY utf8mb4, utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 18446744073709551615;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 4294967296;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=4294967295
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1.5;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 MAX_ROWS 1e0;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| TRUNCATE TABLE utf8mb4;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE TABLE utf8mb4 IGNORE SELECT utf8mb4;
< error
| 1050 42S01 Table 'utf8mb4' already exists
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT + 1;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT 1;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT - 1;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER TABLE utf8mb4 ADD utf8mb4 INT DEFAULT NULL;
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW utf8mb4 AS SELECT utf8mb4 zzz;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR"
< column "Table_type" "ENUM"
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| CREATE TABLE `utf8mb4` (
|   `b` int NOT NULL,
|   `utf8mb4` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `utf8mb4` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER ALGORITHM = MERGE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER ALGORITHM = TEMPTABLE VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER SQL SECURITY DEFINER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER SQL SECURITY INVOKER VIEW utf8mb4 AS SELECT utf8mb4;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH CASCADED CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY utf8mb4 (d))
< ok 0 0
> exec
| ALTER VIEW utf8mb4 AS SELECT utf8mb4 WITH LOCAL CHECK OPTION;
< error
| 1347 HY000 'omni_t_simulation.utf8mb4' is not VIEW
