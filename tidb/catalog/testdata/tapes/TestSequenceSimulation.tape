-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSequenceSimulation default#1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
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
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE ident
< error
| 4139 HY000 Unknown SEQUENCE: 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
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
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
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
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
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
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT BY 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident START 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident START WITH 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident MINVALUE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MINVALUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident MAXVALUE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMAXVALUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MAXVALUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident CACHE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CACHE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident CYCLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCYCLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CYCLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE CYCLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT + 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT - 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
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
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS ident NOMINVALUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE NOMINVALUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART WITH 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT BY 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident START 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident START WITH 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident MINVALUE 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident MAXVALUE 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMAXVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MAXVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident CACHE 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CACHE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident CYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE CYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT + 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT - 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS ident NOMINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE NOMINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART WITH 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `ident`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `ident` (`b`, `ident`, `d`) AS SELECT `b` AS `b`,`ident` AS `ident`,`d` AS `d` FROM `simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT BY 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident START 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident START WITH 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident MINVALUE 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident MAXVALUE 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMAXVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MAXVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident CACHE 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CACHE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident CYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE CYCLE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT + 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT - 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `ident`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `ident` (`b`, `ident`, `d`) AS SELECT `b` AS `b`,`ident` AS `ident`,`d` AS `d` FROM `simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS ident NOMINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE NOMINVALUE
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
| CREATE VIEW ident AS SELECT b, ident, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART WITH 1
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| CREATE SEQUENCE ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| DROP SEQUENCE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT BY 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident START 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident START WITH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident MINVALUE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident MAXVALUE 1
< error
| 4136 HY000 Sequence 'simulation.ident' values are conflicting
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident CACHE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 cycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache cycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT + 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT - 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS ident NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART WITH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| CREATE SEQUENCE ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| DROP SEQUENCE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| CREATE SEQUENCE ident NOMINVALUE NOMINVALUE
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by 1 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT BY 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by 1 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident START 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident START WITH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with 1 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident MINVALUE 1
< error
| 4136 HY000 Sequence 'simulation.ident' values are conflicting
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident MAXVALUE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 1 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NO MAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident CACHE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 1 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 nocache cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 nocache cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 nocycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NO CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 nocycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NOCACHE CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 nocache cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT + 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by 1 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident INCREMENT - 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -1 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS ident NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident NOMINVALUE NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE SEQUENCE ident START WITH -5 MINVALUE -100 MAXVALUE 100 INCREMENT BY -2 CACHE 50 CYCLE COMMENT 'seq'
< ok 0 0
> exec
| ALTER SEQUENCE ident RESTART WITH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| SEQUENCE
> query
| SHOW CREATE SEQUENCE `ident`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE SEQUENCE `ident` start with -5 minvalue -100 maxvalue 100 increment by -2 cache 50 cycle ENGINE=InnoDB COMMENT='seq'
