-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestQuerySimulation default#1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident(ident) AS(SELECT *) SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident IS NOT TRUE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident IS NOT TRUE`) AS SELECT `ident` IS NOT TRUE AS `ident IS NOT TRUE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NEXT VALUE FOR ident FROM ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NEXTVAL(ident) FROM ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident`) AS SELECT `ident` AS `ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) DELETE FROM ident
< error
| 1064 42000 …line 1 column 48 near "DELETE FROM ident" 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) DELETE FROM ident zzz
< error
| 1064 42000 …line 1 column 48 near "DELETE FROM ident zzz" 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) DELETE FROM ident USING ident
< error
| 1064 42000 …line 1 column 48 near "DELETE FROM ident USING ident" 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) DELETE FROM ident ident
< error
| 1064 42000 …line 1 column 48 near "DELETE FROM ident ident" 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) DELETE FROM ident USING ident zzz
< error
| 1064 42000 …line 1 column 48 near "DELETE FROM ident USING ident zzz" 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT @v := ident FROM ident
< error
| 1351 HY000 View's SELECT contains a variable or parameter
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident || ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident || ident`) AS SELECT `ident` OR `ident` AS `ident || ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident XOR ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident XOR ident`) AS SELECT `ident` XOR `ident` AS `ident XOR ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident && ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident && ident`) AS SELECT `ident` AND `ident` AS `ident && ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NOT ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`NOT ident`) AS SELECT NOT `ident` AS `NOT ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH(ident) AGAINST(ident) FROM ident
< error
| 1105 HY000 UnknownType: *ast.MatchAgainst
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident IS TRUE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident IS TRUE`) AS SELECT `ident` IS TRUE AS `ident IS TRUE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident IS FALSE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident IS FALSE`) AS SELECT `ident` IS FALSE AS `ident IS FALSE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident IS UNKNOWN FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident IS UNKNOWN`) AS SELECT `ident` IS NULL AS `ident IS UNKNOWN` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH(ident) AGAINST(ident IN NATURAL LANGUAGE MODE) FROM ident
< error
| 1105 HY000 UnknownType: *ast.MatchAgainst
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH(ident) AGAINST(ident IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION) FROM ident
< error
| 1105 HY000 UnknownType: *ast.MatchAgainst
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH(ident) AGAINST(ident IN BOOLEAN MODE) FROM ident
< error
| 1105 HY000 UnknownType: *ast.MatchAgainst
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH(ident) AGAINST(ident WITH QUERY EXPANSION) FROM ident
< error
| 1105 HY000 UnknownType: *ast.MatchAgainst
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident OR ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident OR ident`) AS SELECT `ident` OR `ident` AS `ident OR ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident AND ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident AND ident`) AS SELECT `ident` AND `ident` AS `ident AND ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident, ident FROM ident
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ASCII() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ascii'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ASCII(zzz) FROM ident
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ASCII(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ASCII(ident)`) AS SELECT ASCII(`ident`) AS `ASCII(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURTIME() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURTIME()`) AS SELECT CURTIME() AS `CURTIME()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURTIME(1) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURTIME(1)`) AS SELECT CURTIME(1) AS `CURTIME(1)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident IS NULL FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident IS NULL`) AS SELECT `ident` IS NULL AS `ident IS NULL` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident >= ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident >= ident`) AS SELECT `ident`>=`ident` AS `ident >= ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident >= ANY(SELECT *) FROM ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident >=(SELECT *) FROM ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident >= @v := ident FROM ident
< error
| 1351 HY000 View's SELECT contains a variable or parameter
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident > ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident > ident`) AS SELECT `ident`>`ident` AS `ident > ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident <= ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident <= ident`) AS SELECT `ident`<=`ident` AS `ident <= ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident < ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident < ident`) AS SELECT `ident`<`ident` AS `ident < ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident != ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident != ident`) AS SELECT `ident`!=`ident` AS `ident != ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident <> ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident <> ident`) AS SELECT `ident`!=`ident` AS `ident <> ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident = ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident = ident`) AS SELECT `ident`=`ident` AS `ident = ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident <=> ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident <=> ident`) AS SELECT `ident`<=>`ident` AS `ident <=> ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident BETWEEN ident AND ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident BETWEEN ident AND ident`) AS SELECT `ident` BETWEEN `ident` AND `ident` AS `ident BETWEEN ident AND ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident NOT BETWEEN ident AND ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident NOT BETWEEN ident AND ident`) AS SELECT `ident` NOT BETWEEN `ident` AND `ident` AS `ident NOT BETWEEN ident AND ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident IN(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident IN(ident)`) AS SELECT `ident` IN (`ident`) AS `ident IN(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident(ident) FROM ident
< error
| 1305 42000 FUNCTION simulation.ident does not exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident NOT IN(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident NOT IN(ident)`) AS SELECT `ident` NOT IN (`ident`) AS `ident NOT IN(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident LIKE ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident LIKE ident`) AS SELECT `ident` LIKE `ident` AS `ident LIKE ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident NOT LIKE ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident NOT LIKE ident`) AS SELECT `ident` NOT LIKE `ident` AS `ident NOT LIKE ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident ILIKE ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident ILIKE ident`) AS SELECT `ident` ILIKE `ident` AS `ident ILIKE ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident NOT ILIKE ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident NOT ILIKE ident`) AS SELECT `ident` NOT ILIKE `ident` AS `ident NOT ILIKE ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident REGEXP ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident REGEXP ident`) AS SELECT `ident` REGEXP `ident` AS `ident REGEXP ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident NOT REGEXP ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident NOT REGEXP ident`) AS SELECT `ident` NOT REGEXP `ident` AS `ident NOT REGEXP ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident >= SOME(SELECT *) FROM ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident >= ALL(SELECT *) FROM ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident IN(SELECT *) FROM ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident MEMBER OF(ident) FROM ident
< error
| 3146 22032 Invalid data type for JSON data in argument 2 to function member of; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident RLIKE ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident RLIKE ident`) AS SELECT `ident` REGEXP `ident` AS `ident RLIKE ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident.*
< error
| 1051 42S02 Unknown table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident.ident.*
< error
| 1051 42S02 Unknown table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident
< error
| 1054 42S22 Unknown column 'ident' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident zzz
< error
| 1054 42S22 Unknown column 'ident' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident ident
< error
| 1054 42S22 Unknown column 'ident' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident AS ident
< error
| 1054 42S22 Unknown column 'ident' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident 'binary'
< error
| 1054 42S22 Unknown column 'ident' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident AS 'binary'
< error
| 1054 42S22 Unknown column 'ident' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT *, *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * GROUP BY ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * GROUP BY ident WITH ROLLUP
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * HAVING ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW()
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(zzz)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(ident)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(ident, ident)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FALSE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`FALSE`) AS SELECT FALSE AS `FALSE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FAULTS FROM ident
< error
| 1054 42S22 Unknown column 'faults' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NULL FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`NULL`) AS SELECT NULL AS `NULL` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NULLS FROM ident
< error
| 1054 42S22 Unknown column 'nulls' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRUE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRUE`) AS SELECT TRUE AS `TRUE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRUE_CARD_COST FROM ident
< error
| 1054 42S22 Unknown column 'true_card_cost' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 1e0 FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`1e0`) AS SELECT 1e+00 AS `1e0` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 1.5 FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`1.5`) AS SELECT 1.5 AS `1.5` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 1 FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`1`) AS SELECT 1 AS `1` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 'binary' FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`binary`) AS SELECT _UTF8MB4'binary' AS `binary` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT _utf8 'binary' FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`binary`) AS SELECT _UTF8'binary' AS `binary` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 0x1 FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`0x1`) AS SELECT x'01' AS `0x1` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT b'1' FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`b'1'`) AS SELECT b'1' AS `b'1'` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT _utf8 0x1 FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`_utf8 0x1`) AS SELECT x'01' AS `_utf8 0x1` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT _utf8 b'1' FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`_utf8 b'1'`) AS SELECT b'1' AS `_utf8 b'1'` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 'binary' 'binary' FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`binary`) AS SELECT _UTF8MB4'binarybinary' AS `binary` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * ORDER BY ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * GROUP BY ident, ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * GROUP BY ident ASC
< error
| 1235 42000 function GROUP BY expr ASC|DESC has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * GROUP BY ident DESC
< error
| 1235 42000 function GROUP BY expr ASC|DESC has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident | ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident | ident`) AS SELECT `ident`|`ident` AS `ident | ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident & ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident & ident`) AS SELECT `ident`&`ident` AS `ident & ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident << ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident << ident`) AS SELECT `ident`<<`ident` AS `ident << ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident >> ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident >> ident`) AS SELECT `ident`>>`ident` AS `ident >> ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident + ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident + ident`) AS SELECT `ident`+`ident` AS `ident + ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident - ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident - ident`) AS SELECT `ident`-`ident` AS `ident - ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident + INTERVAL ident MICROSECOND FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident + INTERVAL ident MICROSECOND`) AS SELECT DATE_ADD(`ident`, INTERVAL `ident` MICROSECOND) AS `ident + INTERVAL ident MICROSECOND` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident - INTERVAL ident MICROSECOND FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident - INTERVAL ident MICROSECOND`) AS SELECT DATE_SUB(`ident`, INTERVAL `ident` MICROSECOND) AS `ident - INTERVAL ident MICROSECOND` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT INTERVAL ident MICROSECOND + ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`INTERVAL ident MICROSECOND + ident`) AS SELECT DATE_ADD(`ident`, INTERVAL `ident` MICROSECOND) AS `INTERVAL ident MICROSECOND + ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident * ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident * ident`) AS SELECT `ident`*`ident` AS `ident * ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident / ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident / ident`) AS SELECT `ident`/`ident` AS `ident / ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident % ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident % ident`) AS SELECT `ident`%`ident` AS `ident % ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident DIV ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident DIV ident`) AS SELECT `ident` DIV `ident` AS `ident DIV ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident MOD ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident MOD ident`) AS SELECT `ident`%`ident` AS `ident MOD ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident ^ ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident ^ ident`) AS SELECT `ident`^`ident` AS `ident ^ ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident.ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident`) AS SELECT `ident`.`ident` AS `ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident.ident.ident FROM ident
< error
| 1054 42S22 Unknown column 'ident.ident.ident' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_USER FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_USER`) AS SELECT CURRENT_USER() AS `CURRENT_USER` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident() FROM ident
< error
| 1305 42000 FUNCTION simulation.ident does not exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident COLLATE 'binary' FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident COLLATE 'binary'`) AS SELECT `ident` COLLATE binary AS `ident COLLATE 'binary'` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ? FROM ident
< error
| 1149 42000 syntax error, unexpected '?'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT @@v FROM ident
< error
| 1351 HY000 View's SELECT contains a variable or parameter
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`AVG(ident)`) AS SELECT AVG(`ident`) AS `AVG(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ! ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`! ident`) AS SELECT !`ident` AS `! ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ~ ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`~ ident`) AS SELECT ~`ident` AS `~ ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT - ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`- ident`) AS SELECT -`ident` AS `- ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT + ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`+ ident`) AS SELECT +`ident` AS `+ ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW(ident, ident) FROM ident
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT EXISTS(SELECT *) FROM ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT { ident ident } FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ident`) AS SELECT `ident` AS `ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BINARY ident FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`BINARY ident`) AS SELECT BINARY `ident` AS `BINARY ident` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS BINARY) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS BINARY)`) AS SELECT CAST(`ident` AS BINARY(11)) AS `CAST(ident AS BINARY)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE WHEN ident THEN ident END FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CASE WHEN ident THEN ident END`) AS SELECT CASE WHEN `ident` THEN `ident` END AS `CASE WHEN ident THEN ident END` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CONVERT(ident, BINARY) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CONVERT(ident, BINARY)`) AS SELECT CONVERT(`ident`, BINARY(11)) AS `CONVERT(ident, BINARY)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CONVERT(ident USING 'binary') FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CONVERT(ident USING 'binary')`) AS SELECT CONVERT(`ident` USING 'binary') AS `CONVERT(ident USING 'binary')` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DEFAULT(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`DEFAULT(ident)`) AS SELECT DEFAULT(`ident`) AS `DEFAULT(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT VALUES(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`VALUES(ident)`) AS SELECT VALUES(`ident`) AS `VALUES(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident -> 'binary' FROM ident
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_extract; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident ->> 'binary' FROM ident
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_extract; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS BINARY ARRAY) FROM ident
< error
| 1235 42000 This version of TiDB doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DISTINCT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DISTINCTROW *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ALL *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`AVG(ALL ident)`) AS SELECT AVG(`ident`) AS `AVG(ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * UNION SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * UNION ALL SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(DISTINCT ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`AVG(DISTINCT ALL ident)`) AS SELECT AVG(DISTINCT `ident`) AS `AVG(DISTINCT ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(DISTINCT ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`AVG(DISTINCT ident)`) AS SELECT AVG(DISTINCT `ident`) AS `AVG(DISTINCT ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(DISTINCTROW ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`AVG(DISTINCTROW ALL ident)`) AS SELECT AVG(DISTINCT `ident`) AS `AVG(DISTINCTROW ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CHARSET() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'charset'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COALESCE() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'coalesce'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COLLATION() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'collation'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATE() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'date'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATABASE() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`DATABASE()`) AS SELECT DATABASE() AS `DATABASE()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DAY() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'day'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HOUR() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'hour'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT IF() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'if'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT INTERVAL() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'interval'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LOG() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'log'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FORMAT() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'format'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEFT() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'left'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MICROSECOND() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'microsecond'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MINUTE() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'minute'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MONTH() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'month'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NOW() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`NOW()`) AS SELECT NOW() AS `NOW()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT POINT() FROM ident
< error
| 1305 42000 FUNCTION simulation.point does not exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT QUARTER() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'quarter'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT REPEAT() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'repeat'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT REPLACE() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'replace'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT REVERSE() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'reverse'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT RIGHT() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'right'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_COUNT() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_COUNT()`) AS SELECT ROW_COUNT() AS `ROW_COUNT()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SECOND() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'second'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIME() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'time'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMP() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'timestamp'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRUNCATE() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'truncate'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT USER() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`USER()`) AS SELECT USER() AS `USER()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WEEK() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'week'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT YEAR() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'year'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_USER() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_USER()`) AS SELECT CURRENT_USER() AS `CURRENT_USER()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_DATE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_DATE`) AS SELECT CURRENT_DATE() AS `CURRENT_DATE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_ROLE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_ROLE`) AS SELECT CURRENT_ROLE() AS `CURRENT_ROLE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT UTC_DATE FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`UTC_DATE`) AS SELECT UTC_DATE() AS `UTC_DATE` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIDB_CURRENT_TSO FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TIDB_CURRENT_TSO`) AS SELECT TIDB_CURRENT_TSO() AS `TIDB_CURRENT_TSO` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIDB_JSON FROM ident
< error
| 1054 42S22 Unknown column 'tidb_json' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_TIME`) AS SELECT CURRENT_TIME() AS `CURRENT_TIME` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIMESTAMP FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_TIMESTAMP`) AS SELECT CURRENT_TIMESTAMP() AS `CURRENT_TIMESTAMP` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LOCALTIME FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`LOCALTIME`) AS SELECT LOCALTIME() AS `LOCALTIME` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LOCALTIMESTAMP FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`LOCALTIMESTAMP`) AS SELECT LOCALTIMESTAMP() AS `LOCALTIMESTAMP` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LOCATION FROM ident
< error
| 1054 42S22 Unknown column 'location' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT UTC_TIME FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`UTC_TIME`) AS SELECT UTC_TIME() AS `UTC_TIME` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT UTC_TIMESTAMP FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`UTC_TIMESTAMP`) AS SELECT UTC_TIMESTAMP() AS `UTC_TIMESTAMP` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT UTILIZATION_LIMIT FROM ident
< error
| 1054 42S22 Unknown column 'utilization_limit' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURDATE() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURDATE()`) AS SELECT CURDATE() AS `CURDATE()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CHAR(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CHAR(ident)`) AS SELECT CHAR_FUNC(`ident`, NULL) AS `CHAR(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CHAR(ident USING 'binary') FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CHAR(ident USING 'binary')`) AS SELECT CHAR_FUNC(`ident`, 'binary') AS `CHAR(ident USING 'binary')` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATE 'binary' FROM ident
< error
| 1292 22007 Incorrect date value: 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATE FROM ident
< error
| 1054 42S22 Unknown column 'date' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIME 'binary' FROM ident
< error
| 1292 22007 Incorrect time value: 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIME FROM ident
< error
| 1054 42S22 Unknown column 'time' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMP 'binary' FROM ident
< error
| 1292 22007 Incorrect datetime value: 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMP FROM ident
< error
| 1054 42S22 Unknown column 'timestamp' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT INSERT() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'insert_func'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MOD(ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`MOD(ident, ident)`) AS SELECT `ident`%`ident` AS `MOD(ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PASSWORD() FROM ident
< error
| 1582 42000 Incorrect parameter count in the call to native function 'password_func'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SYSDATE() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`SYSDATE()`) AS SELECT SYSDATE() AS `SYSDATE()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ADDDATE(ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ADDDATE(ident, ident)`) AS SELECT ADDDATE(`ident`, INTERVAL `ident` DAY) AS `ADDDATE(ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ADDDATE(ident, INTERVAL ident MICROSECOND) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ADDDATE(ident, INTERVAL ident MICROSECOND)`) AS SELECT ADDDATE(`ident`, INTERVAL `ident` MICROSECOND) AS `ADDDATE(ident, INTERVAL ident MICROSECOND)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATE_ADD(ident, INTERVAL ident MICROSECOND) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`DATE_ADD(ident, INTERVAL ident MICROSECOND)`) AS SELECT DATE_ADD(`ident`, INTERVAL `ident` MICROSECOND) AS `DATE_ADD(ident, INTERVAL ident MICROSECOND)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT EXTRACT(MICROSECOND FROM ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`EXTRACT(MICROSECOND FROM ident)`) AS SELECT EXTRACT(MICROSECOND FROM `ident`) AS `EXTRACT(MICROSECOND FROM ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GET_FORMAT(DATE, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`GET_FORMAT(DATE, ident)`) AS SELECT GET_FORMAT(DATE, `ident`) AS `GET_FORMAT(DATE, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT POSITION(ident IN ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`POSITION(ident IN ident)`) AS SELECT POSITION(`ident` IN `ident`) AS `POSITION(ident IN ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUBSTRING(ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`SUBSTRING(ident, ident)`) AS SELECT SUBSTRING(`ident`, `ident`) AS `SUBSTRING(ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUBSTRING(ident FROM ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`SUBSTRING(ident FROM ident)`) AS SELECT SUBSTRING(`ident`, `ident`) AS `SUBSTRING(ident FROM ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUBSTRING(ident, ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`SUBSTRING(ident, ident, ident)`) AS SELECT SUBSTRING(`ident`, `ident`, `ident`) AS `SUBSTRING(ident, ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUBSTRING(ident FROM ident FOR ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`SUBSTRING(ident FROM ident FOR ident)`) AS SELECT SUBSTRING(`ident`, `ident`, `ident`) AS `SUBSTRING(ident FROM ident FOR ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMPADD(MICROSECOND, ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TIMESTAMPADD(MICROSECOND, ident, ident)`) AS SELECT TIMESTAMPADD(MICROSECOND, `ident`, `ident`) AS `TIMESTAMPADD(MICROSECOND, ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(MICROSECOND, ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TIMESTAMPDIFF(MICROSECOND, ident, ident)`) AS SELECT TIMESTAMPDIFF(MICROSECOND, `ident`, `ident`) AS `TIMESTAMPDIFF(MICROSECOND, ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRIM(ident)`) AS SELECT TRIM(`ident`) AS `TRIM(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(ident FROM ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRIM(ident FROM ident)`) AS SELECT TRIM(`ident` FROM `ident`) AS `TRIM(ident FROM ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(BOTH FROM ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRIM(BOTH FROM ident)`) AS SELECT TRIM(BOTH _UTF8MB4' ' FROM `ident`) AS `TRIM(BOTH FROM ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(BOTH ident FROM ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRIM(BOTH ident FROM ident)`) AS SELECT TRIM(BOTH `ident` FROM `ident`) AS `TRIM(BOTH ident FROM ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WEIGHT_STRING(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`WEIGHT_STRING(ident)`) AS SELECT WEIGHT_STRING(`ident`) AS `WEIGHT_STRING(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WEIGHT_STRING(ident AS CHARACTER(1)) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`WEIGHT_STRING(ident AS CHARACTER(1))`) AS SELECT WEIGHT_STRING(`ident` AS CHAR(1)) AS `WEIGHT_STRING(ident AS CHARACTER(1))` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WEIGHT_STRING(ident AS BINARY(1)) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`WEIGHT_STRING(ident AS BINARY(1))`) AS SELECT WEIGHT_STRING(`ident` AS BINARY(1)) AS `WEIGHT_STRING(ident AS BINARY(1))` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LASTVAL(ident) FROM ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRANSLATE(ident, ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRANSLATE(ident, ident, ident)`) AS SELECT TRANSLATE(`ident`, `ident`, `ident`) AS `TRANSLATE(ident, ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GET_FORMAT(DATETIME, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`GET_FORMAT(DATETIME, ident)`) AS SELECT GET_FORMAT(DATETIME, `ident`) AS `GET_FORMAT(DATETIME, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GET_FORMAT(TIME, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`GET_FORMAT(TIME, ident)`) AS SELECT GET_FORMAT(TIME, `ident`) AS `GET_FORMAT(TIME, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GET_FORMAT(TIMESTAMP, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`GET_FORMAT(TIMESTAMP, ident)`) AS SELECT GET_FORMAT(DATETIME, `ident`) AS `GET_FORMAT(TIMESTAMP, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATE_SUB(ident, INTERVAL ident MICROSECOND) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`DATE_SUB(ident, INTERVAL ident MICROSECOND)`) AS SELECT DATE_SUB(`ident`, INTERVAL `ident` MICROSECOND) AS `DATE_SUB(ident, INTERVAL ident MICROSECOND)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUBDATE(ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`SUBDATE(ident, ident)`) AS SELECT SUBDATE(`ident`, INTERVAL `ident` DAY) AS `SUBDATE(ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(BOTH zzz FROM ident) FROM ident
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(BOUND FROM ident) FROM ident
< error
| 1054 42S22 Unknown column 'bound' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(LEADING FROM ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRIM(LEADING FROM ident)`) AS SELECT TRIM(LEADING _UTF8MB4' ' FROM `ident`) AS `TRIM(LEADING FROM ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(LEADING zzz FROM ident) FROM ident
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(LEARNER FROM ident) FROM ident
< error
| 1054 42S22 Unknown column 'learner' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(TRAILING FROM ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`TRIM(TRAILING FROM ident)`) AS SELECT TRIM(TRAILING _UTF8MB4' ' FROM `ident`) AS `TRIM(TRAILING FROM ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(TRAILING zzz FROM ident) FROM ident
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(TRANSACTION FROM ident) FROM ident
< error
| 1054 42S22 Unknown column 'transaction' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SETVAL(ident, 1) FROM ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT APPROX_COUNT_DISTINCT(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`APPROX_COUNT_DISTINCT(ident)`) AS SELECT APPROX_COUNT_DISTINCT(`ident`) AS `APPROX_COUNT_DISTINCT(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT APPROX_PERCENTILE(ident) FROM ident
< error
| 1105 HY000 APPROX_PERCENTILE should take 2 arguments
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BIT_AND(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`BIT_AND(ident)`) AS SELECT BIT_AND(`ident`) AS `BIT_AND(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BIT_AND(ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`BIT_AND(ALL ident)`) AS SELECT BIT_AND(`ident`) AS `BIT_AND(ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BIT_OR(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`BIT_OR(ident)`) AS SELECT BIT_OR(`ident`) AS `BIT_OR(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BIT_OR(ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`BIT_OR(ALL ident)`) AS SELECT BIT_OR(`ident`) AS `BIT_OR(ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BIT_XOR(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`BIT_XOR(ident)`) AS SELECT BIT_XOR(`ident`) AS `BIT_XOR(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BIT_XOR(ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`BIT_XOR(ALL ident)`) AS SELECT BIT_XOR(`ident`) AS `BIT_XOR(ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COUNT(DISTINCT ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`COUNT(DISTINCT ident)`) AS SELECT COUNT(DISTINCT `ident`) AS `COUNT(DISTINCT ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COUNT(ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`COUNT(ALL ident)`) AS SELECT COUNT(`ident`) AS `COUNT(ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COUNT(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`COUNT(ident)`) AS SELECT COUNT(`ident`) AS `COUNT(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COUNT(*) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`COUNT(*)`) AS SELECT COUNT(1) AS `COUNT(*)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GROUP_CONCAT(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`GROUP_CONCAT(ident)`) AS SELECT GROUP_CONCAT(`ident` SEPARATOR ',') AS `GROUP_CONCAT(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MAX(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`MAX(ident)`) AS SELECT MAX(`ident`) AS `MAX(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MIN(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`MIN(ident)`) AS SELECT MIN(`ident`) AS `MIN(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MINUTE(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`MINUTE(ident)`) AS SELECT MINUTE(`ident`) AS `MINUTE(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUM(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`SUM(ident)`) AS SELECT SUM(`ident`) AS `SUM(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STDDEV_POP(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`STDDEV_POP(ident)`) AS SELECT STDDEV_POP(`ident`) AS `STDDEV_POP(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STDDEV_SAMP(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`STDDEV_SAMP(ident)`) AS SELECT STDDEV_SAMP(`ident`) AS `STDDEV_SAMP(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT VAR_POP(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`VAR_POP(ident)`) AS SELECT VAR_POP(`ident`) AS `VAR_POP(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT VAR_SAMP(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`VAR_SAMP(ident)`) AS SELECT VAR_SAMP(`ident`) AS `VAR_SAMP(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAYAGG(ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`JSON_ARRAYAGG(ident)`) AS SELECT JSON_ARRAYAGG(`ident`) AS `JSON_ARRAYAGG(ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAYAGG(ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`JSON_ARRAYAGG(ALL ident)`) AS SELECT JSON_ARRAYAGG(`ident`) AS `JSON_ARRAYAGG(ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECTAGG(ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`JSON_OBJECTAGG(ident, ident)`) AS SELECT JSON_OBJECTAGG(`ident`, `ident`) AS `JSON_OBJECTAGG(ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECTAGG(ALL ident, ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`JSON_OBJECTAGG(ALL ident, ident)`) AS SELECT JSON_OBJECTAGG(`ident`, `ident`) AS `JSON_OBJECTAGG(ALL ident, ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECTAGG(ident, ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`JSON_OBJECTAGG(ident, ALL ident)`) AS SELECT JSON_OBJECTAGG(`ident`, `ident`) AS `JSON_OBJECTAGG(ident, ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECTAGG(ALL ident, ALL ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`JSON_OBJECTAGG(ALL ident, ALL ident)`) AS SELECT JSON_OBJECTAGG(`ident`, `ident`) AS `JSON_OBJECTAGG(ALL ident, ALL ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GROUP_CONCAT(ident SEPARATOR 'binary') FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`GROUP_CONCAT(ident SEPARATOR 'binary')`) AS SELECT GROUP_CONCAT(`ident` SEPARATOR 'binary') AS `GROUP_CONCAT(ident SEPARATOR 'binary')` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ident.ident() FROM ident
< error
| 1305 42000 FUNCTION ident.ident does not exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_TIME()`) AS SELECT CURRENT_TIME() AS `CURRENT_TIME()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(1) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CURRENT_TIME(1)`) AS SELECT CURRENT_TIME(1) AS `CURRENT_TIME(1)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE zzz WHEN ident THEN ident END FROM ident
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE ident WHEN ident THEN ident END FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CASE ident WHEN ident THEN ident END`) AS SELECT CASE `ident` WHEN `ident` THEN `ident` END AS `CASE ident WHEN ident THEN ident END` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE WHEN ident THEN ident WHEN ident THEN ident END FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CASE WHEN ident THEN ident WHEN ident THEN ident END`) AS SELECT CASE WHEN `ident` THEN `ident` WHEN `ident` THEN `ident` END AS `CASE WHEN ident THEN ident WHEN ident THEN ident END` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE WHEN ident THEN ident ELSE ident END FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CASE WHEN ident THEN ident ELSE ident END`) AS SELECT CASE WHEN `ident` THEN `ident` ELSE `ident` END AS `CASE WHEN ident THEN ident ELSE ident END` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS CHARACTER) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS CHARACTER)`) AS SELECT CAST(`ident` AS CHAR(11)) AS `CAST(ident AS CHARACTER)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS DATE) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS DATE)`) AS SELECT CAST(`ident` AS DATE) AS `CAST(ident AS DATE)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS YEAR) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS YEAR)`) AS SELECT CAST(`ident` AS YEAR) AS `CAST(ident AS YEAR)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS DATETIME) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS DATETIME)`) AS SELECT CAST(`ident` AS DATETIME) AS `CAST(ident AS DATETIME)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS DECIMAL) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS DECIMAL)`) AS SELECT CAST(`ident` AS DECIMAL(10)) AS `CAST(ident AS DECIMAL)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS TIME) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS TIME)`) AS SELECT CAST(`ident` AS TIME) AS `CAST(ident AS TIME)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS SIGNED) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS SIGNED)`) AS SELECT CAST(`ident` AS SIGNED) AS `CAST(ident AS SIGNED)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS UNSIGNED) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS UNSIGNED)`) AS SELECT CAST(`ident` AS UNSIGNED) AS `CAST(ident AS UNSIGNED)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS JSON) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS JSON)`) AS SELECT CAST(`ident` AS JSON) AS `CAST(ident AS JSON)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS DOUBLE) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS DOUBLE)`) AS SELECT CAST(`ident` AS DOUBLE) AS `CAST(ident AS DOUBLE)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS FLOAT) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS FLOAT)`) AS SELECT CAST(`ident` AS FLOAT) AS `CAST(ident AS FLOAT)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS REAL) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS REAL)`) AS SELECT CAST(`ident` AS DOUBLE) AS `CAST(ident AS REAL)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS VECTOR) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS VECTOR)`) AS SELECT CAST(`ident` AS VECTOR) AS `CAST(ident AS VECTOR)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LOW_PRIORITY *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HIGH_PRIORITY *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DELAYED *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FROM DUAL
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`b`, `ident`, `d`) AS SELECT `simulation`.`ident`.`b` AS `b`,`simulation`.`ident`.`ident` AS `ident`,`simulation`.`ident`.`d` AS `d` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FROM ident zzz
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`b`, `ident`, `d`) AS SELECT `simulation`.`zzz`.`b` AS `b`,`simulation`.`zzz`.`ident` AS `ident`,`simulation`.`zzz`.`d` AS `d` FROM `simulation`.`ident` AS `zzz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`b`, `ident`, `d`) AS TABLE `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *)(SELECT *)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH RECURSIVE ident AS(SELECT *) SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *), ident AS(SELECT *) SELECT *
< error
| 1066 42000 Not unique table/alias: '%-.192s'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FROM ident WINDOW ident AS()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`b`, `ident`, `d`) AS SELECT `simulation`.`ident`.`b` AS `b`,`simulation`.`ident`.`ident` AS `ident`,`simulation`.`ident`.`d` AS `d` FROM `simulation`.`ident` WINDOW `ident` AS ()
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FROM ident WINDOW ident AS(), ident AS()
< error
| 3591 HY000 Window 'ident' is defined twice.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER() FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_NUMBER() OVER()`) AS SELECT ROW_NUMBER() OVER () AS `ROW_NUMBER() OVER()` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(zzz) FROM ident
< error
| 3579 HY000 Window name 'zzz' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ident) FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(PARTITION BY ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_NUMBER() OVER(PARTITION BY ident)`) AS SELECT ROW_NUMBER() OVER (PARTITION BY `ident`) AS `ROW_NUMBER() OVER(PARTITION BY ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ORDER BY ident) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_NUMBER() OVER(ORDER BY ident)`) AS SELECT ROW_NUMBER() OVER (ORDER BY `ident`) AS `ROW_NUMBER() OVER(ORDER BY ident)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING)`) AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS `ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(RANGE UNBOUNDED PRECEDING) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_NUMBER() OVER(RANGE UNBOUNDED PRECEDING)`) AS SELECT ROW_NUMBER() OVER (RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS `ROW_NUMBER() OVER(RANGE UNBOUNDED PRECEDING)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(GROUPS UNBOUNDED PRECEDING) FROM ident
< error
| 1235 42000 This version of TiDB doesn't yet support 'GROUPS'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED PRECEDING) FROM ident
< error
| 3585 HY000 Window '<unnamed window>': frame end cannot be UNBOUNDED PRECEDING.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS 1 PRECEDING) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_NUMBER() OVER(ROWS 1 PRECEDING)`) AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) AS `ROW_NUMBER() OVER(ROWS 1 PRECEDING)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS ? PRECEDING) FROM ident
< error
| 1149 42000 syntax error, unexpected '?'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS INTERVAL ident MICROSECOND PRECEDING) FROM ident
< error
| 3596 HY000 Window '<unnamed window>': INTERVAL can only be used with RANGE frames.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS CURRENT ROW) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`ROW_NUMBER() OVER(ROWS CURRENT ROW)`) AS SELECT ROW_NUMBER() OVER (ROWS BETWEEN CURRENT ROW AND CURRENT ROW) AS `ROW_NUMBER() OVER(ROWS CURRENT ROW)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN UNBOUNDED FOLLOWING AND UNBOUNDED PRECEDING) FROM ident
< error
| 3584 HY000 Window '<unnamed window>': frame start cannot be UNBOUNDED FOLLOWING.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN 1 FOLLOWING AND UNBOUNDED PRECEDING) FROM ident
< error
| 3585 HY000 Window '<unnamed window>': frame end cannot be UNBOUNDED PRECEDING.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM ident
< error
| 1149 42000 syntax error, unexpected '?'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN INTERVAL ident MICROSECOND FOLLOWING AND UNBOUNDED PRECEDING) FROM ident
< error
| 3585 HY000 Window '<unnamed window>': frame end cannot be UNBOUNDED PRECEDING.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(ident) OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT RANK() OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DENSE_RANK() OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CUME_DIST() OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PERCENT_RANK() OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(ident) OVER ident FROM ident
< error
| 1210 HY000 Incorrect arguments to ntile
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(ident) OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAG(ident) OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FIRST_VALUE(ident) OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAST_VALUE(ident) OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(ident, ident) OVER ident FROM ident
< error
| 1210 HY000 Incorrect arguments to nth_value
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(ident, 1) OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(ident, ?) OVER ident FROM ident
< error
| 1149 42000 syntax error, unexpected '?'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(ident, 1, ident) OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(ident) RESPECT NULLS OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(ident) IGNORE NULLS OVER ident FROM ident
< error
| 3579 HY000 Window name 'ident' is not defined.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(ident, ident) FROM FIRST OVER ident FROM ident
< error
| 1210 HY000 Incorrect arguments to nth_value
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(ident, ident) FROM LAST OVER ident FROM ident
< error
| 1210 HY000 Incorrect arguments to nth_value
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * LIMIT 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * LIMIT ?
< error
| 1149 42000 syntax error, unexpected '?'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FETCH FIRST ROW ONLY
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FETCH FIRST ROWS ONLY
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FETCH NEXT ROW ONLY
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FETCH FIRST 1 ROW ONLY
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * LIMIT 1, 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * LIMIT 1 OFFSET 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_SMALL_RESULT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BIG_RESULT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BUFFER_RESULT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CACHE *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CALC_FOUND_ROWS *
< error
| 1235 42000 function SQL_CALC_FOUND_ROWS has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_NO_CACHE *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary'
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR UPDATE
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR SHARE
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR UPDATE NOWAIT
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR UPDATE WAIT 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR SHARE NOWAIT
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR UPDATE SKIP LOCKED
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR SHARE SKIP LOCKED
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * LOCK IN SHARE MODE
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * FOR UPDATE OF ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *)(SELECT *) LIMIT 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) SELECT * UNION SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * UNION(SELECT *)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * UNION(SELECT *) ORDER BY ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * UNION(SELECT *) LIMIT 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * UNION(SELECT *) ORDER BY ident LIMIT 1
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * UNION SELECT * UNION SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * EXCEPT SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTERSECT SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT @v FROM ident
< error
| 1351 HY000 View's SELECT contains a variable or parameter
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS SIGNED INTEGER) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS SIGNED INTEGER)`) AS SELECT CAST(`ident` AS SIGNED) AS `CAST(ident AS SIGNED INTEGER)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(ident AS SIGNED INT) FROM ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` (`CAST(ident AS SIGNED INT)`) AS SELECT CAST(`ident` AS SIGNED) AS `CAST(ident AS SIGNED INT)` FROM `simulation`.`ident`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH ident AS(SELECT *) UPDATE ident SET ident = ident
< error
| 1064 42000 …line 1 column 48 near "UPDATE ident SET ident = ident" 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * WHERE ident
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' FIELDS TERMINATED BY 'binary'
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' FIELDS DEFINED NULL BY 'binary'
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' FIELDS DEFINED NULL BY 'binary' OPTIONALLY ENCLOSED
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' FIELDS OPTIONALLY ENCLOSED BY 0x1
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' FIELDS OPTIONALLY ENCLOSED BY b'1'
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' LINES
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' LINES STARTING BY 'binary'
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT * INTO OUTFILE 'binary' LINES TERMINATED BY 'binary'
< error
| 1350 HY000 View's SELECT contains a 'INFO' clause
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(), ROW()
< error
| 1051 42S02 Unknown table ''
