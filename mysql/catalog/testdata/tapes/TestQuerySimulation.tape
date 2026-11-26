-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestQuerySimulation default#1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 zzz
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS TABLE STATUS
< error
| 1146 42S02 Table 'omni_t_simulation.STATUS' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS FLOAT) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as float) AS `CAST(utf8mb4 AS FLOAT)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME() FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME()` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME() zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME(1) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(1) AS `CURRENT_TIME(1)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME(1) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(1) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIMESTAMP FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select now() AS `CURRENT_TIMESTAMP` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIMESTAMP zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select now() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' CHARSET utf8mb4
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' CHAR SET utf8mb4
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4 AS CHAR(1)) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`utf8mb4`.`utf8mb4` as char(1)) AS `WEIGHT_STRING(utf8mb4 AS CHAR(1))` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 FOR UPDATE
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO utf8mb4
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO utf8mb4 FOR UPDATE
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 utf8mb4 FOR UPDATE
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FOR UPDATE INTO utf8mb4
< error
| 1064 42000 …near 'INTO utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH utf8mb4 AS(SELECT utf8mb4) SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH utf8mb4 AS(SELECT utf8mb4) SELECT utf8mb4 zzz
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 UNION SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 UNION SELECT utf8mb4 zzz
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 EXCEPT SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 EXCEPT SELECT utf8mb4 zzz
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 EXCEPT DISTINCT SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 EXCEPT ALL SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTERSECT SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTERSECT SELECT utf8mb4 zzz
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTERSECT DISTINCT SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTERSECT ALL SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW()
< error
| 3942 HY000 Each row of a VALUES clause must have at least one column, unless when used as source in an INSERT statement.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS TABLE utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`b` AS `b`,`utf8mb4`.`utf8mb4` AS `utf8mb4`,`utf8mb4`.`d` AS `d` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 FROM DUAL
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FROM DUMPFILE
< error
| 1146 42S02 Table 'omni_t_simulation.DUMPFILE' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FROM utf8mb4 zzz
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `zzz`.`utf8mb4` AS `utf8mb4` from `utf8mb4` `zzz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STREAM utf8mb4
< error
| 1054 42S22 Unknown column 'STREAM' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN STRAIGHT_JOIN utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN STRAIGHT_JOIN zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_NO_CACHE utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_NO_CACHE zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FOR UPDATE FOR UPDATE
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FOR UPDATE OF utf8mb4
< error
| 3568 HY000 Unresolved table name `utf8mb4` in locking clause.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LOCK IN SHARE MODE
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FOR SHARE
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FOR UPDATE NOWAIT
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 FOR UPDATE SKIP LOCKED
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4, utf8mb4 FROM utf8mb4
< error
| 1060 42S21 Duplicate column name 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4, utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4`,`utf8mb4`.`utf8mb4` AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT * FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`b` AS `b`,`utf8mb4`.`utf8mb4` AS `utf8mb4`,`utf8mb4`.`d` AS `d` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4.* FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`b` AS `b`,`utf8mb4`.`utf8mb4` AS `utf8mb4`,`utf8mb4`.`d` AS `d` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 AS utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 AS 's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `s` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `s` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_USER FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_user() AS `CURRENT_USER` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_USER zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_user() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_USER() FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_user() AS `CURRENT_USER()` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_USER() zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_user() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 OR utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) or (0 <> `utf8mb4`.`utf8mb4`)) AS `utf8mb4 OR utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 OR utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) or (0 <> `utf8mb4`.`utf8mb4`)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 XOR utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) xor (0 <> `utf8mb4`.`utf8mb4`)) AS `utf8mb4 XOR utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 XOR utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) xor (0 <> `utf8mb4`.`utf8mb4`)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 AND utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) and (0 <> `utf8mb4`.`utf8mb4`)) AS `utf8mb4 AND utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 AND utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) and (0 <> `utf8mb4`.`utf8mb4`)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NOT utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (0 = `utf8mb4`.`utf8mb4`) AS `NOT utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NOT utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (0 = `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NOWAIT utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'NOWAIT' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 IS TRUE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is true) AS `utf8mb4 IS TRUE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS TRUE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is true) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT TRUE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is not true) AS `utf8mb4 IS NOT TRUE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT TRUE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is not true) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS FALSE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is false) AS `utf8mb4 IS FALSE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS FALSE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is false) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT FALSE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is not false) AS `utf8mb4 IS NOT FALSE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT FALSE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) is not false) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS UNKNOWN FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is null) AS `utf8mb4 IS UNKNOWN` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 UNKNOWN FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `UNKNOWN` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS UNKNOWN zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is null) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT UNKNOWN FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is not null) AS `utf8mb4 IS NOT UNKNOWN` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT UNKNOWN zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is not null) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NULL FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is null) AS `utf8mb4 IS NULL` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NULL zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is null) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT NULL FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is not null) AS `utf8mb4 IS NOT NULL` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IS NOT NULL zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` is not null) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 = utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` = `utf8mb4`.`utf8mb4`) AS `utf8mb4 = utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 = utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` = `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 = ALL(SELECT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` = all (select `utf8mb4`.`utf8mb4`) AS `utf8mb4 = ALL(SELECT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 =(SELECT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` = (select `utf8mb4`.`utf8mb4`)) AS `utf8mb4 =(SELECT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 = ALL(SELECT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` = all (select `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IN(SELECT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` in (select `utf8mb4`.`utf8mb4`) AS `utf8mb4 IN(SELECT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IN(SELECT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` in (select `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT IN(SELECT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (not(`utf8mb4`.`utf8mb4` in (select `utf8mb4`.`utf8mb4`))) AS `utf8mb4 NOT IN(SELECT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT IN(SELECT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (not(`utf8mb4`.`utf8mb4` in (select `utf8mb4`.`utf8mb4`))) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IN(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` = `utf8mb4`.`utf8mb4`) AS `utf8mb4 IN(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IN(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` = `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IN(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` in (`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`)) AS `utf8mb4 IN(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 IN(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` in (`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT IN(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` <> `utf8mb4`.`utf8mb4`) AS `utf8mb4 NOT IN(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT IN(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` <> `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT IN(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` not in (`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`)) AS `utf8mb4 NOT IN(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT IN(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` not in (`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 MEMBER(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` member of (`utf8mb4`.`utf8mb4`) AS `utf8mb4 MEMBER(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 MEMBER(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` member of (`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 BETWEEN utf8mb4 AND utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` between `utf8mb4`.`utf8mb4` and `utf8mb4`.`utf8mb4`) AS `utf8mb4 BETWEEN utf8mb4 AND utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 BETWEEN utf8mb4 AND utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` between `utf8mb4`.`utf8mb4` and `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT BETWEEN utf8mb4 AND utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` not between `utf8mb4`.`utf8mb4` and `utf8mb4`.`utf8mb4`) AS `utf8mb4 NOT BETWEEN utf8mb4 AND utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT BETWEEN utf8mb4 AND utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` not between `utf8mb4`.`utf8mb4` and `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 SOUNDS LIKE utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (soundex(`utf8mb4`.`utf8mb4`) = soundex(`utf8mb4`.`utf8mb4`)) AS `utf8mb4 SOUNDS LIKE utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SOUNDS LIKE utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'SOUNDS' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 SOUNDS LIKE utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (soundex(`utf8mb4`.`utf8mb4`) = soundex(`utf8mb4`.`utf8mb4`)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 LIKE utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` like `utf8mb4`.`utf8mb4`) AS `utf8mb4 LIKE utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 LIKE utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` like `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 LIKE utf8mb4 ESCAPE utf8mb4 FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIKE ESCAPE utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'ESCAPE' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIKE utf8mb4 utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` like `utf8mb4`.`utf8mb4`) AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 LIKE utf8mb4 ESCAPE utf8mb4 zzz FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 NOT LIKE utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (not((`utf8mb4`.`utf8mb4` like `utf8mb4`.`utf8mb4`))) AS `utf8mb4 NOT LIKE utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT LIKE utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (not((`utf8mb4`.`utf8mb4` like `utf8mb4`.`utf8mb4`))) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT LIKE utf8mb4 ESCAPE utf8mb4 FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 NOT LIKE ESCAPE utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'ESCAPE' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 NOT LIKE utf8mb4 utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (not((`utf8mb4`.`utf8mb4` like `utf8mb4`.`utf8mb4`))) AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT LIKE utf8mb4 ESCAPE utf8mb4 zzz FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 REGEXP utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select regexp_like(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `utf8mb4 REGEXP utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 REGEXP utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select regexp_like(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT REGEXP utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (not(regexp_like(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`))) AS `utf8mb4 NOT REGEXP utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 NOT REGEXP utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (not(regexp_like(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`))) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 MEMBER OF(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` member of (`utf8mb4`.`utf8mb4`) AS `utf8mb4 MEMBER OF(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 | utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` | `utf8mb4`.`utf8mb4`) AS `utf8mb4 | utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 | utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` | `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 & utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` & `utf8mb4`.`utf8mb4`) AS `utf8mb4 & utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 & utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` & `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 << utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` << `utf8mb4`.`utf8mb4`) AS `utf8mb4 << utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 << utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` << `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 >> utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` >> `utf8mb4`.`utf8mb4`) AS `utf8mb4 >> utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 >> utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` >> `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + `utf8mb4`.`utf8mb4`) AS `utf8mb4 + utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 - utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - `utf8mb4`.`utf8mb4`) AS `utf8mb4 - utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 - utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `utf8mb4 + INTERVAL utf8mb4 DAY` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + utf8mb4 DAY FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + `utf8mb4`.`utf8mb4`) AS `DAY` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 - INTERVAL utf8mb4 DAY FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `utf8mb4 - INTERVAL utf8mb4 DAY` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 - utf8mb4 DAY FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - `utf8mb4`.`utf8mb4`) AS `DAY` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 - INTERVAL utf8mb4 DAY zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 * utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` * `utf8mb4`.`utf8mb4`) AS `utf8mb4 * utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 * utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` * `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 / utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` / `utf8mb4`.`utf8mb4`) AS `utf8mb4 / utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 / utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` / `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 % utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` % `utf8mb4`.`utf8mb4`) AS `utf8mb4 % utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 % utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` % `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 DIV utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` DIV `utf8mb4`.`utf8mb4`) AS `utf8mb4 DIV utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 DIV utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` DIV `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 MOD utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` % `utf8mb4`.`utf8mb4`) AS `utf8mb4 MOD utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 MOD utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` % `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 ^ utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` ^ `utf8mb4`.`utf8mb4`) AS `utf8mb4 ^ utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 ^ utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` ^ `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 OR zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 || utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) or (0 <> `utf8mb4`.`utf8mb4`)) AS `utf8mb4 || utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 || zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 AND zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 && utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ((0 <> `utf8mb4`.`utf8mb4`) and (0 <> `utf8mb4`.`utf8mb4`)) AS `utf8mb4 && utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 && zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ! utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (0 = `utf8mb4`.`utf8mb4`) AS `! utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ! zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 = zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 <=> utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` <=> `utf8mb4`.`utf8mb4`) AS `utf8mb4 <=> utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 <=> zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 >= utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` >= `utf8mb4`.`utf8mb4`) AS `utf8mb4 >= utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 >= zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 > utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` > `utf8mb4`.`utf8mb4`) AS `utf8mb4 > utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 > zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 <= utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` <= `utf8mb4`.`utf8mb4`) AS `utf8mb4 <= utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 <= zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 < utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` < `utf8mb4`.`utf8mb4`) AS `utf8mb4 < utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 < zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 <> utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` <> `utf8mb4`.`utf8mb4`) AS `utf8mb4 <> utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 <> zzz utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 = ANY(SELECT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` in (select `utf8mb4`.`utf8mb4`) AS `utf8mb4 = ANY(SELECT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_DATE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curdate() AS `CURRENT_DATE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_DATE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curdate() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4() FROM utf8mb4
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4() zzz FROM utf8mb4
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATABASE() FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select database() AS `DATABASE()` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATABASE() zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select database() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 COLLATE utf8mb4 FROM utf8mb4
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 COLLATE utf8mb4 zzz FROM utf8mb4
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 's' AS `s` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 's' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 's' AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ? FROM utf8mb4
< error
| 1064 42000 …near '? FROM utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ? zzz FROM utf8mb4
< error
| 1064 42000 …near '? zzz FROM utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(`utf8mb4`.`utf8mb4`) AS `AVG(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT AVG(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT + utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT + utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT - utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select -(`utf8mb4`.`utf8mb4`) AS `- utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT - utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select -(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ~ utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ~(`utf8mb4`.`utf8mb4`) AS `~ utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ~ utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ~(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ! utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (0 = `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW(utf8mb4, utf8mb4) FROM utf8mb4
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW(utf8mb4, utf8mb4) zzz FROM utf8mb4
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT EXISTS(SELECT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select exists(select `utf8mb4`.`utf8mb4`) AS `EXISTS(SELECT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT EXISTS(SELECT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select exists(select `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT { utf8mb4 utf8mb4 } FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT { utf8mb4 utf8mb4 } zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MATCH utf8mb4 AGAINST(utf8mb4) FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH utf8mb4 AGAINST(utf8mb4) zzz FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BINARY utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset binary) AS `BINARY utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BINARY utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset binary) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BINLOG utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'BINLOG' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS BINARY) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset binary) AS `CAST(utf8mb4 AS BINARY)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS BINARY) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset binary) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AT TIME ZONE 's' AS DATETIME) FROM utf8mb4
< error
| 3998 HY000 Cannot cast value to TIMESTAMP WITH TIME ZONE.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AT TIME ZONE 's' AS DATETIME) zzz FROM utf8mb4
< error
| 3998 HY000 Cannot cast value to TIMESTAMP WITH TIME ZONE.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE WHEN utf8mb4 THEN utf8mb4 END FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (case when `utf8mb4`.`utf8mb4` then `utf8mb4`.`utf8mb4` end) AS `CASE WHEN utf8mb4 THEN utf8mb4 END` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CASE WHEN utf8mb4 THEN utf8mb4 END zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (case when `utf8mb4`.`utf8mb4` then `utf8mb4`.`utf8mb4` end) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CONVERT(utf8mb4, BINARY) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset binary) AS `CONVERT(utf8mb4, BINARY)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CONVERT(utf8mb4, BINARY) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset binary) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CONVERT(utf8mb4 USING utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select convert(`utf8mb4`.`utf8mb4` using utf8mb4) AS `CONVERT(utf8mb4 USING utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CONVERT(utf8mb4 USING utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select convert(`utf8mb4`.`utf8mb4` using utf8mb4) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DEFAULT(utf8mb4) FROM utf8mb4
< error
| 1143 42000 create view command denied to user 'root'@'%' for column 'DEFAULT(utf8mb4)' in table '_v'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DEFAULT(utf8mb4) zzz FROM utf8mb4
< error
| 1143 42000 create view command denied to user 'root'@'%' for column 'zzz' in table '_v'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT VALUES(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select NULL AS `VALUES(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VALUES(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select NULL AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL utf8mb4 DAY + utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `INTERVAL utf8mb4 DAY + utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL utf8mb4 DAY + utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 -> 's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_extract(`utf8mb4`.`utf8mb4`,'s') AS `utf8mb4 -> 's'` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 -> 's' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_extract(`utf8mb4`.`utf8mb4`,'s') AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 ->> 's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_unquote(json_extract(`utf8mb4`.`utf8mb4`,'s')) AS `utf8mb4 ->> 's'` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 ->> 's' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_unquote(json_extract(`utf8mb4`.`utf8mb4`,'s')) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHAR(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select char(`utf8mb4`.`utf8mb4`) AS `CHAR(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHAR(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select char(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHAR(utf8mb4 USING utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select char(`utf8mb4`.`utf8mb4`) AS `CHAR(utf8mb4 USING utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHAR(utf8mb4 USING utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select char(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as date) AS `DATE(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as date) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DAY(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select dayofmonth(`utf8mb4`.`utf8mb4`) AS `DAY(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DAY(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select dayofmonth(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT HOUR(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select hour(`utf8mb4`.`utf8mb4`) AS `HOUR(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT HOUR(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select hour(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INSERT(utf8mb4, utf8mb4, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select insert(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `INSERT(utf8mb4, utf8mb4, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INSERT(utf8mb4, utf8mb4, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select insert(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select interval(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `INTERVAL(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select interval(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL(utf8mb4, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select interval(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `INTERVAL(utf8mb4, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL(utf8mb4, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select interval(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's') FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's') zzz FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEFT(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select left(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `LEFT(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT LEFT(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select left(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MINUTE(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select minute(`utf8mb4`.`utf8mb4`) AS `MINUTE(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MINUTE(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select minute(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MONTH(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select month(`utf8mb4`.`utf8mb4`) AS `MONTH(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MONTH(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select month(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT RIGHT(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select right(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `RIGHT(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT RIGHT(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select right(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SECOND(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select second(`utf8mb4`.`utf8mb4`) AS `SECOND(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SECOND(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select second(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIME(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as time) AS `TIME(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIME(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as time) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMP(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as datetime) AS `TIMESTAMP(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMP(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as datetime) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMP(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestamp(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMP(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMP(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestamp(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(`utf8mb4`.`utf8mb4`) AS `TRIM(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(LEADING utf8mb4 FROM utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(leading `utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `TRIM(LEADING utf8mb4 FROM utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(LEADING utf8mb4 FROM utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(leading `utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(TRAILING utf8mb4 FROM utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(trailing `utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `TRIM(TRAILING utf8mb4 FROM utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(TRAILING utf8mb4 FROM utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(trailing `utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(BOTH utf8mb4 FROM utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(both `utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `TRIM(BOTH utf8mb4 FROM utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(BOTH utf8mb4 FROM utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(both `utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(utf8mb4 FROM utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(`utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `TRIM(utf8mb4 FROM utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(utf8mb4 FROM utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(`utf8mb4`.`utf8mb4` from `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT USER() FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select user() AS `USER()` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT USER() zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select user() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT YEAR(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select year(`utf8mb4`.`utf8mb4`) AS `YEAR(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT YEAR(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select year(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ADDDATE(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `ADDDATE(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ADDDATE(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ADDDATE(utf8mb4, INTERVAL utf8mb4 DAY) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `ADDDATE(utf8mb4, INTERVAL utf8mb4 DAY)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ADDDATE(utf8mb4, INTERVAL utf8mb4 DAY) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE_ADD(utf8mb4, INTERVAL utf8mb4 DAY) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `DATE_ADD(utf8mb4, INTERVAL utf8mb4 DAY)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE_ADD(utf8mb4, INTERVAL utf8mb4 DAY) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE_SUB(utf8mb4, INTERVAL utf8mb4 DAY) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `DATE_SUB(utf8mb4, INTERVAL utf8mb4 DAY)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE_SUB(utf8mb4, INTERVAL utf8mb4 DAY) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT EXTRACT(DAY FROM utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select extract(day from `utf8mb4`.`utf8mb4`) AS `EXTRACT(DAY FROM utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT EXTRACT(DAY FROM utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select extract(day from `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(DATE, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(DATE, `utf8mb4`.`utf8mb4`) AS `GET_FORMAT(DATE, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(DATE, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(DATE, `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT POSITION(utf8mb4 IN utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select locate(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `POSITION(utf8mb4 IN utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT POSITION(utf8mb4 IN utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select locate(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBDATE(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `SUBDATE(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBDATE(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBDATE(utf8mb4, INTERVAL utf8mb4 DAY) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `SUBDATE(utf8mb4, INTERVAL utf8mb4 DAY)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBDATE(utf8mb4, INTERVAL utf8mb4 DAY) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` - interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `SUBSTRING(utf8mb4, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `SUBSTRING(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4 FROM utf8mb4 FOR utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `SUBSTRING(utf8mb4 FROM utf8mb4 FOR utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4 FROM utf8mb4 FOR utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4 FROM utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `SUBSTRING(utf8mb4 FROM utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(utf8mb4 FROM utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SYSDATE FROM utf8mb4
< error
| 1054 42S22 Unknown column 'SYSDATE' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SYSDATE zzz FROM utf8mb4
< error
| 1054 42S22 Unknown column 'SYSDATE' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMPADD(DAY, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `TIMESTAMPADD(DAY, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(DAY, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(WEEK, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` week) AS `TIMESTAMPADD(WEEK, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(HOUR, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour) AS `TIMESTAMPADD(HOUR, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(MINUTE, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` minute) AS `TIMESTAMPADD(MINUTE, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(MONTH, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` month) AS `TIMESTAMPADD(MONTH, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(QUARTER, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` quarter) AS `TIMESTAMPADD(QUARTER, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(SECOND, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` second) AS `TIMESTAMPADD(SECOND, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(MICROSECOND, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` microsecond) AS `TIMESTAMPADD(MICROSECOND, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPADD(YEAR, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` year) AS `TIMESTAMPADD(YEAR, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(DAY, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(DAY,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(DAY, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(DAY, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(DAY,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(WEEK, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(WEEK,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(WEEK, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(HOUR, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(HOUR,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(HOUR, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(MINUTE, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(MINUTE,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(MINUTE, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(MONTH, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(MONTH,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(MONTH, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(QUARTER, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(QUARTER,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(QUARTER, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(SECOND, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(SECOND,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(SECOND, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(MICROSECOND, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(MICROSECOND,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(MICROSECOND, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(YEAR, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(YEAR,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TIMESTAMPDIFF(YEAR, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_DATE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_date() AS `UTC_DATE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_DATE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_date() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_TIME FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_time() AS `UTC_TIME` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_TIME zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_time() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_TIMESTAMP FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_timestamp() AS `UTC_TIMESTAMP` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_TIMESTAMP zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_timestamp() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VALIDATION FROM utf8mb4
< error
| 1054 42S22 Unknown column 'VALIDATION' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's' RETURNING BINARY) FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ASCII(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ascii(`utf8mb4`.`utf8mb4`) AS `ASCII(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ASCII(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ascii(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHARSET(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select charset(`utf8mb4`.`utf8mb4`) AS `CHARSET(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHARSET(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select charset(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COALESCE(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select coalesce(`utf8mb4`.`utf8mb4`) AS `COALESCE(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COALESCE(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select coalesce(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLLATION(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select collation(`utf8mb4`.`utf8mb4`) AS `COLLATION(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLLATION(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select collation(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT IF(utf8mb4, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select if(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `IF(utf8mb4, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT IF(utf8mb4, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select if(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT FORMAT(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select format(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `FORMAT(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT FORMAT(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select format(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT FORMAT(utf8mb4, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select format(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `FORMAT(utf8mb4, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT FORMAT(utf8mb4, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select format(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MICROSECOND(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select microsecond(`utf8mb4`.`utf8mb4`) AS `MICROSECOND(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MICROSECOND(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select microsecond(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MOD(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` % `utf8mb4`.`utf8mb4`) AS `MOD(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MOD(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` % `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT QUARTER(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select quarter(`utf8mb4`.`utf8mb4`) AS `QUARTER(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT QUARTER(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select quarter(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REPEAT(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select repeat(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `REPEAT(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REPEAT(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select repeat(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REPLACE(utf8mb4, utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select replace(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `REPLACE(utf8mb4, utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REPLACE(utf8mb4, utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select replace(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REVERSE(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select reverse(`utf8mb4`.`utf8mb4`) AS `REVERSE(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REVERSE(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select reverse(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_COUNT() FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_count() AS `ROW_COUNT()` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_COUNT() zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_count() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRUNCATE(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select truncate(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `TRUNCATE(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRUNCATE(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select truncate(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEEK(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select week(`utf8mb4`.`utf8mb4`,0) AS `WEEK(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEEK(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select week(`utf8mb4`.`utf8mb4`,0) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEEK(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select week(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `WEEK(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEEK(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select week(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`utf8mb4`.`utf8mb4`) AS `WEIGHT_STRING(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4 AS CHAR(1)) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`utf8mb4`.`utf8mb4` as char(1)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4 AS BINARY(1)) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(cast(`utf8mb4`.`utf8mb4` as char(1) charset binary)) AS `WEIGHT_STRING(utf8mb4 AS BINARY(1))` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4 AS BINARY(1)) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(cast(`utf8mb4`.`utf8mb4` as char(1) charset binary)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4, 1, 1, 1) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`utf8mb4`.`utf8mb4` as char(1)) AS `WEIGHT_STRING(utf8mb4, 1, 1, 1)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(utf8mb4, 1, 1, 1) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`utf8mb4`.`utf8mb4` as char(1)) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GEOMETRYCOLLECTION() FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select geomcollection() AS `GEOMETRYCOLLECTION()` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GEOMETRYCOLLECTION() zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select geomcollection() AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT LINESTRING(utf8mb4) FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LINESTRING(utf8mb4) zzz FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MULTILINESTRING(utf8mb4) FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MULTILINESTRING(utf8mb4) zzz FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MULTIPOINT(utf8mb4) FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MULTIPOINT(utf8mb4) zzz FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MULTIPOLYGON(utf8mb4) FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MULTIPOLYGON(utf8mb4) zzz FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT POINT(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select point(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `POINT(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT POINT(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select point(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT POLYGON(utf8mb4, utf8mb4) FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT POLYGON(utf8mb4) FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT POLYGON(utf8mb4) zzz FROM utf8mb4
< error
| 1367 22007 Illegal non geometric '`omni_t_simulation`.`utf8mb4`.`utf8mb4`' value found during parsing
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4() FROM utf8mb4
< error
| 1305 42000 FUNCTION utf8mb4.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4() zzz FROM utf8mb4
< error
| 1305 42000 FUNCTION utf8mb4.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH utf8mb4 AGAINST(utf8mb4 IN BOOLEAN MODE) FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH utf8mb4 AGAINST(utf8mb4 IN NATURAL LANGUAGE MODE) FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH utf8mb4 AGAINST(utf8mb4 WITH QUERY EXPANSION) FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4(zzz) FROM utf8mb4
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4(utf8mb4) FROM utf8mb4
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4(utf8mb4 zzz) FROM utf8mb4
< error
| 1584 42000 Incorrect parameters in the call to stored function `utf8mb4`
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4(utf8mb4, utf8mb4) FROM utf8mb4
< error
| 1305 42000 FUNCTION omni_t_simulation.utf8mb4 does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4(utf8mb4 utf8mb4) FROM utf8mb4
< error
| 1584 42000 Incorrect parameters in the call to stored function `utf8mb4`
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4(utf8mb4, utf8mb4 zzz) FROM utf8mb4
< error
| 1584 42000 Incorrect parameters in the call to stored function `utf8mb4`
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GROUPING(utf8mb4) FROM utf8mb4
< error
| 1111 HY000 Invalid use of group function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GROUPING(utf8mb4) zzz FROM utf8mb4
< error
| 1111 HY000 Invalid use of group function
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(DISTINCT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(distinct `utf8mb4`.`utf8mb4`) AS `AVG(DISTINCT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT AVG(DISTINCT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(distinct `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_AND(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_and(`utf8mb4`.`utf8mb4`) AS `BIT_AND(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_AND(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_and(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_OR(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_or(`utf8mb4`.`utf8mb4`) AS `BIT_OR(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_OR(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_or(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_ARRAYAGG(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_arrayagg(`utf8mb4`.`utf8mb4`) AS `JSON_ARRAYAGG(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_ARRAYAGG(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_arrayagg(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_OBJECTAGG(utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_objectagg(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `JSON_OBJECTAGG(utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_OBJECTAGG(utf8mb4, utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_objectagg(`utf8mb4`.`utf8mb4`,`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ST_COLLECT(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select st_collect(`utf8mb4`.`utf8mb4`) AS `ST_COLLECT(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ST_COLLECT(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select st_collect(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ST_COLLECT(DISTINCT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select st_collect(distinct `utf8mb4`.`utf8mb4`) AS `ST_COLLECT(DISTINCT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ST_COLLECT(DISTINCT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select st_collect(distinct `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_XOR(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_xor(`utf8mb4`.`utf8mb4`) AS `BIT_XOR(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_XOR(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_xor(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(*) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(0) AS `COUNT(*)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(*) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(0) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(ALL *) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(0) AS `COUNT(ALL *)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(`utf8mb4`.`utf8mb4`) AS `COUNT(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(DISTINCT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(distinct `utf8mb4`.`utf8mb4`) AS `COUNT(DISTINCT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(DISTINCT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(distinct `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MIN(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select min(`utf8mb4`.`utf8mb4`) AS `MIN(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MIN(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select min(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MIN(DISTINCT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select min(`utf8mb4`.`utf8mb4`) AS `MIN(DISTINCT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MIN(DISTINCT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select min(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MAX(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select max(`utf8mb4`.`utf8mb4`) AS `MAX(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MAX(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select max(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MAX(DISTINCT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select max(`utf8mb4`.`utf8mb4`) AS `MAX(DISTINCT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MAX(DISTINCT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select max(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT STD(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select std(`utf8mb4`.`utf8mb4`) AS `STD(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT STD(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select std(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VARIANCE(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select variance(`utf8mb4`.`utf8mb4`) AS `VARIANCE(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VARIANCE(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select variance(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT STDDEV_SAMP(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select stddev_samp(`utf8mb4`.`utf8mb4`) AS `STDDEV_SAMP(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT STDDEV_SAMP(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select stddev_samp(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VAR_SAMP(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select var_samp(`utf8mb4`.`utf8mb4`) AS `VAR_SAMP(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VAR_SAMP(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select var_samp(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUM(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sum(`utf8mb4`.`utf8mb4`) AS `SUM(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUM(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sum(`utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUM(DISTINCT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sum(distinct `utf8mb4`.`utf8mb4`) AS `SUM(DISTINCT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUM(DISTINCT utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sum(distinct `utf8mb4`.`utf8mb4`) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`utf8mb4`.`utf8mb4` separator ',') AS `GROUP_CONCAT(utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(utf8mb4) zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`utf8mb4`.`utf8mb4` separator ',') AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT RANK() OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT RANK() OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DENSE_RANK() OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DENSE_RANK() OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CUME_DIST() OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CUME_DIST() OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PERCENT_RANK() OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PERCENT_RANK() OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(1) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(1) OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(utf8mb4) OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAG(utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAG(utf8mb4) OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAG(utf8mb4) RESPECT NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAG(utf8mb4) IGNORE NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FIRST_VALUE(utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FIRST_VALUE(utf8mb4) OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FIRST_VALUE(utf8mb4) RESPECT NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FIRST_VALUE(utf8mb4) IGNORE NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAST_VALUE(utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAST_VALUE(utf8mb4) OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAST_VALUE(utf8mb4) RESPECT NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAST_VALUE(utf8mb4) IGNORE NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(utf8mb4, utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(utf8mb4, utf8mb4) OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(utf8mb4, utf8mb4) RESPECT NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(utf8mb4, utf8mb4) IGNORE NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(utf8mb4, 1) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(?) OVER utf8mb4 FROM utf8mb4
< error
| 1064 42000 …near '?) OVER utf8mb4 FROM utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 1327 42000 Undeclared variable: utf8mb4
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(@utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 1351 HY000 View's SELECT contains a variable or parameter
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(@) OVER utf8mb4 FROM utf8mb4
< error
| 1351 HY000 View's SELECT contains a variable or parameter
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(utf8mb4, 1, utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(utf8mb4) RESPECT NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(utf8mb4) IGNORE NULLS OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(utf8mb4, utf8mb4) FROM FIRST OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(utf8mb4, utf8mb4) FROM LAST OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(utf8mb4) OVER utf8mb4 FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(utf8mb4) utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(`utf8mb4`.`utf8mb4`) AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT AVG(utf8mb4) OVER utf8mb4 zzz FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER() FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER ()  AS `ROW_NUMBER() OVER()` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER() zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER ()  AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(zzz) FROM utf8mb4
< error
| 3579 HY000 Window name 'zzz' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(utf8mb4) FROM utf8mb4
< error
| 3579 HY000 Window name 'utf8mb4' is not defined.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(PARTITION BY utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER (PARTITION BY `utf8mb4`.`utf8mb4` )  AS `ROW_NUMBER() OVER(PARTITION BY utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ORDER BY utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER (ORDER BY `utf8mb4`.`utf8mb4` )  AS `ROW_NUMBER() OVER(ORDER BY utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER (ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)  AS `ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED PRECEDING) FROM utf8mb4
< error
| 3585 HY000 Window '<unnamed window>': frame end cannot be UNBOUNDED PRECEDING.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS 1 PRECEDING) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER (ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)  AS `ROW_NUMBER() OVER(ROWS 1 PRECEDING)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS ? PRECEDING) FROM utf8mb4
< error
| 1064 42000 …near '? PRECEDING) FROM utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS INTERVAL utf8mb4 DAY PRECEDING) FROM utf8mb4
< error
| 3596 HY000 Window '<unnamed window>': INTERVAL can only be used with RANGE frames.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS CURRENT ROW) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER (ROWS BETWEEN CURRENT ROW AND CURRENT ROW)  AS `ROW_NUMBER() OVER(ROWS CURRENT ROW)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN UNBOUNDED FOLLOWING AND UNBOUNDED PRECEDING) FROM utf8mb4
< error
| 3584 HY000 Window '<unnamed window>': frame start cannot be UNBOUNDED FOLLOWING.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN 1 FOLLOWING AND UNBOUNDED PRECEDING) FROM utf8mb4
< error
| 3585 HY000 Window '<unnamed window>': frame end cannot be UNBOUNDED PRECEDING.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) FROM utf8mb4
< error
| 1064 42000 …near '? FOLLOWING AND UNBOUNDED PRECEDING) FROM utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN INTERVAL utf8mb4 DAY FOLLOWING AND UNBOUNDED PRECEDING) FROM utf8mb4
< error
| 3596 HY000 Window '<unnamed window>': INTERVAL can only be used with RANGE frames.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE CURRENT ROW) FROM utf8mb4
< error
| 1235 42000 This version of MySQL doesn't yet support 'EXCLUDE'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE GROUP) FROM utf8mb4
< error
| 1235 42000 This version of MySQL doesn't yet support 'EXCLUDE'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE TIES) FROM utf8mb4
< error
| 1235 42000 This version of MySQL doesn't yet support 'EXCLUDE'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE NO OTHERS) FROM utf8mb4
< error
| 1235 42000 This version of MySQL doesn't yet support 'EXCLUDE'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(RANGE UNBOUNDED PRECEDING) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() OVER (RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)  AS `ROW_NUMBER() OVER(RANGE UNBOUNDED PRECEDING)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(GROUPS UNBOUNDED PRECEDING) FROM utf8mb4
< error
| 1235 42000 This version of MySQL doesn't yet support 'GROUPS'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT GROUP_CONCAT(DISTINCT utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(distinct `utf8mb4`.`utf8mb4` separator ',') AS `GROUP_CONCAT(DISTINCT utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(utf8mb4 SEPARATOR 's') FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`utf8mb4`.`utf8mb4` separator 's') AS `GROUP_CONCAT(utf8mb4 SEPARATOR 's')` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(utf8mb4 ORDER BY utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`utf8mb4`.`utf8mb4` order by `utf8mb4`.`utf8mb4` ASC separator ',') AS `GROUP_CONCAT(utf8mb4 ORDER BY utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(utf8mb4 ORDER BY utf8mb4, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`utf8mb4`.`utf8mb4` order by `utf8mb4`.`utf8mb4` ASC,`utf8mb4`.`utf8mb4` ASC separator ',') AS `GROUP_CONCAT(utf8mb4 ORDER BY utf8mb4, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS CHAR) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset utf8mb4) AS `CAST(utf8mb4 AS CHAR)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS NCHAR) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset utf8mb3) AS `CAST(utf8mb4 AS NCHAR)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS NATIONAL CHAR) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as char charset utf8mb3) AS `CAST(utf8mb4 AS NATIONAL CHAR)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS SIGNED) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as signed) AS `CAST(utf8mb4 AS SIGNED)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS SIGNED INT) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as signed) AS `CAST(utf8mb4 AS SIGNED INT)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS UNSIGNED) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as unsigned) AS `CAST(utf8mb4 AS UNSIGNED)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS UNSIGNED INT) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as unsigned) AS `CAST(utf8mb4 AS UNSIGNED INT)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS DATE) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as date) AS `CAST(utf8mb4 AS DATE)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS YEAR) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `cast_as_year`(`utf8mb4`.`utf8mb4`) AS `CAST(utf8mb4 AS YEAR)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS TIME) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as time) AS `CAST(utf8mb4 AS TIME)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS DATETIME) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as datetime) AS `CAST(utf8mb4 AS DATETIME)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS DECIMAL) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as decimal(10,0)) AS `CAST(utf8mb4 AS DECIMAL)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS JSON) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as json) AS `CAST(utf8mb4 AS JSON)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS REAL) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as double) AS `CAST(utf8mb4 AS REAL)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS POINT) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as point) AS `CAST(utf8mb4 AS POINT)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS LINESTRING) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as linestring) AS `CAST(utf8mb4 AS LINESTRING)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS POLYGON) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as polygon) AS `CAST(utf8mb4 AS POLYGON)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS MULTIPOINT) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as multipoint) AS `CAST(utf8mb4 AS MULTIPOINT)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS MULTILINESTRING) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as multilinestring) AS `CAST(utf8mb4 AS MULTILINESTRING)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS MULTIPOLYGON) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as multipolygon) AS `CAST(utf8mb4 AS MULTIPOLYGON)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AS GEOMETRYCOLLECTION) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`utf8mb4`.`utf8mb4` as geometrycollection) AS `CAST(utf8mb4 AS GEOMETRYCOLLECTION)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MATCH(utf8mb4) AGAINST(utf8mb4) FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH utf8mb4, utf8mb4 AGAINST(utf8mb4) FROM utf8mb4
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE zzz WHEN utf8mb4 THEN utf8mb4 END FROM utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE utf8mb4 WHEN utf8mb4 THEN utf8mb4 END FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (case `utf8mb4`.`utf8mb4` when `utf8mb4`.`utf8mb4` then `utf8mb4`.`utf8mb4` end) AS `CASE utf8mb4 WHEN utf8mb4 THEN utf8mb4 END` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CASE WHEN utf8mb4 THEN utf8mb4 ELSE utf8mb4 END FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (case when `utf8mb4`.`utf8mb4` then `utf8mb4`.`utf8mb4` else `utf8mb4`.`utf8mb4` end) AS `CASE WHEN utf8mb4 THEN utf8mb4 ELSE utf8mb4 END` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CASE WHEN utf8mb4 THEN utf8mb4 WHEN utf8mb4 THEN utf8mb4 END FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (case when `utf8mb4`.`utf8mb4` then `utf8mb4`.`utf8mb4` when `utf8mb4`.`utf8mb4` then `utf8mb4`.`utf8mb4` end) AS `CASE WHEN utf8mb4 THEN utf8mb4 WHEN utf8mb4 THEN utf8mb4 END` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's' ERROR ON EMPTY) FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's' ERROR ON ERROR) FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's' ERROR ON EMPTY ERROR ON ERROR) FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's' NULL ON EMPTY) FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT JSON_VALUE(utf8mb4, 's' DEFAULT 's' ON EMPTY) FROM utf8mb4
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_HOUR FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_hour) AS `utf8mb4 + INTERVAL utf8mb4 DAY_HOUR` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_HOUR zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_hour) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_MICROSECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_microsecond) AS `utf8mb4 + INTERVAL utf8mb4 DAY_MICROSECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_MICROSECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_microsecond) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_MINUTE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_minute) AS `utf8mb4 + INTERVAL utf8mb4 DAY_MINUTE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_MINUTE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_minute) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_SECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_second) AS `utf8mb4 + INTERVAL utf8mb4 DAY_SECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 DAY_SECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` day_second) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR_MICROSECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour_microsecond) AS `utf8mb4 + INTERVAL utf8mb4 HOUR_MICROSECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR_MICROSECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour_microsecond) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR_MINUTE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour_minute) AS `utf8mb4 + INTERVAL utf8mb4 HOUR_MINUTE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR_MINUTE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour_minute) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR_SECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour_second) AS `utf8mb4 + INTERVAL utf8mb4 HOUR_SECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR_SECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour_second) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MINUTE_MICROSECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` minute_microsecond) AS `utf8mb4 + INTERVAL utf8mb4 MINUTE_MICROSECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MINUTE_MICROSECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` minute_microsecond) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MINUTE_SECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` minute_second) AS `utf8mb4 + INTERVAL utf8mb4 MINUTE_SECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MINUTE_SECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` minute_second) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 SECOND_MICROSECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` second_microsecond) AS `utf8mb4 + INTERVAL utf8mb4 SECOND_MICROSECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 SECOND_MICROSECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` second_microsecond) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 YEAR_MONTH FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` year_month) AS `utf8mb4 + INTERVAL utf8mb4 YEAR_MONTH` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 YEAR_MONTH zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` year_month) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 WEEK FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` week) AS `utf8mb4 + INTERVAL utf8mb4 WEEK` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 WEEK zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` week) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour) AS `utf8mb4 + INTERVAL utf8mb4 HOUR` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 HOUR zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` hour) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MINUTE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` minute) AS `utf8mb4 + INTERVAL utf8mb4 MINUTE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MINUTE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` minute) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MONTH FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` month) AS `utf8mb4 + INTERVAL utf8mb4 MONTH` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MONTH zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` month) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 QUARTER FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` quarter) AS `utf8mb4 + INTERVAL utf8mb4 QUARTER` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 QUARTER zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` quarter) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 SECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` second) AS `utf8mb4 + INTERVAL utf8mb4 SECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 SECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` second) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MICROSECOND FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` microsecond) AS `utf8mb4 + INTERVAL utf8mb4 MICROSECOND` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 MICROSECOND zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` microsecond) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 YEAR FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` year) AS `utf8mb4 + INTERVAL utf8mb4 YEAR` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 + INTERVAL utf8mb4 YEAR zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select (`utf8mb4`.`utf8mb4` + interval `utf8mb4`.`utf8mb4` year) AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(TIME, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(TIME, `utf8mb4`.`utf8mb4`) AS `GET_FORMAT(TIME, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(TIMESTAMP, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(DATETIME, `utf8mb4`.`utf8mb4`) AS `GET_FORMAT(TIMESTAMP, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(DATETIME, utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(DATETIME, `utf8mb4`.`utf8mb4`) AS `GET_FORMAT(DATETIME, utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT AVG(ALL utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(`utf8mb4`.`utf8mb4`) AS `AVG(ALL utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4 WHERE utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 HAVING utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH RECURSIVE utf8mb4 AS(SELECT utf8mb4) SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 WINDOW utf8mb4 AS()
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 WINDOW utf8mb4 AS(), utf8mb4 AS()
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 GROUP BY utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 GROUP BY utf8mb4, utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 GROUP BY utf8mb4 WITH ROLLUP
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 ORDER BY utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 ORDER BY utf8mb4, utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 ORDER BY utf8mb4 ASC
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 ORDER BY utf8mb4 DESC
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIMIT utf8mb4
< error
| 1327 42000 Undeclared variable: utf8mb4
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIMIT utf8mb4, utf8mb4
< error
| 1327 42000 Undeclared variable: utf8mb4
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIMIT utf8mb4 OFFSET utf8mb4
< error
| 1327 42000 Undeclared variable: utf8mb4
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIMIT ?
< error
| 1064 42000 …near '?' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIMIT 18446744073709551615
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIMIT 4294967296
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 LIMIT 1
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO utf8mb4, utf8mb4
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO@utf8mb4
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO@
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO DUMPFILE 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO DUMPFILE
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(), ROW()
< error
| 3942 HY000 Each row of a VALUES clause must have at least one column, unless when used as source in an INSERT statement.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(zzz)
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(utf8mb4)
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(utf8mb4, utf8mb4)
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(DEFAULT)
< error
| 3943 HY000 A VALUES clause cannot use DEFAULT values, unless used as a source in an INSERT statement.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS VALUES ROW(DEFAULT_AUTH)
< error
| 1054 42S22 Unknown column 'DEFAULT_AUTH' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH utf8mb4 AS(SELECT utf8mb4) DELETE FROM utf8mb4
< error
| 1064 42000 …near 'DELETE FROM utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' COLUMNS TERMINATED BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' COLUMNS TERMINATED BY 's' TERMINATED BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' COLUMNS OPTIONALLY ENCLOSED BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' COLUMNS ENCLOSED BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' COLUMNS ESCAPED BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' LINES TERMINATED BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' LINES TERMINATED BY 's' TERMINATED BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 INTO OUTFILE 's' LINES STARTING BY 's'
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT N's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 's' AS `s` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT N's' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 's' AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4's' AS `s` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 's' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4's' AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 's' 's' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 'ss' AS `s` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 's' 's' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 'ss' AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NULL FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select NULL AS `NULL` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NULL zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select NULL AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NULLS FROM utf8mb4
< error
| 1054 42S22 Unknown column 'NULLS' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 1 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1 AS `1` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1 AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE 's' FROM utf8mb4
< error
| 1525 HY000 Incorrect DATE value: 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATE FROM utf8mb4
< error
| 1054 42S22 Unknown column 'DATE' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATE 's' zzz FROM utf8mb4
< error
| 1525 HY000 Incorrect DATE value: 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DATETIME 's' FROM utf8mb4
< error
| 1054 42S22 Unknown column 'DATETIME' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FALSE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select false AS `FALSE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT FALSE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select false AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT FAST FROM utf8mb4
< error
| 1054 42S22 Unknown column 'FAST' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRUE FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select true AS `TRUE` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRUE zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select true AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRUNCATE FROM utf8mb4
< error
| 1054 42S22 Unknown column 'TRUNCATE' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 0x1 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 0x01 AS `0x1` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 0x1 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 0x01 AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT b'1' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 0x01 AS `b'1'` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT b'1' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 0x01 AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 0x1 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
|q "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4'\x01' AS `_utf8mb4 0x1` from `utf8mb4`"
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 0x1 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
|q "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4'\x01' AS `zzz` from `utf8mb4`"
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 b'1' FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
|q "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4'\x01' AS `_utf8mb4 b'1'` from `utf8mb4`"
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 b'1' zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
|q "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4'\x01' AS `zzz` from `utf8mb4`"
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1.5 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1.5 AS `1.5` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1.5 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1.5 AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1e0 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1e0 AS `1e0` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1e0 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1e0 AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 4294967296 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 4294967296 AS `4294967296` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 4294967296 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 4294967296 AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 18446744073709551615 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 18446744073709551615 AS `18446744073709551615` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 18446744073709551615 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 18446744073709551615 AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIME 's' FROM utf8mb4
< error
| 1525 HY000 Incorrect TIME value: 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIME FROM utf8mb4
< error
| 1054 42S22 Unknown column 'TIME' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIME 's' zzz FROM utf8mb4
< error
| 1525 HY000 Incorrect TIME value: 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMP 's' FROM utf8mb4
< error
| 1525 HY000 Incorrect DATETIME value: 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMP FROM utf8mb4
< error
| 1054 42S22 Unknown column 'TIMESTAMP' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMP 's' zzz FROM utf8mb4
< error
| 1525 HY000 Incorrect DATETIME value: 's'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMPADD 's' FROM utf8mb4
< error
| 1054 42S22 Unknown column 'TIMESTAMPADD' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(utf8mb4 AT TIME ZONE INTERVAL 's' AS DATETIME) FROM utf8mb4
< error
| 3998 HY000 Cannot cast value to TIMESTAMP WITH TIME ZONE.
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4.* FROM utf8mb4
< error
| 1051 42S02 Unknown table 'utf8mb4.utf8mb4'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4 zzz FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `zzz` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VALUES(utf8mb4.utf8mb4) FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select NULL AS `VALUES(utf8mb4.utf8mb4)` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4.utf8mb4 FROM utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4.utf8mb4.utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4 utf8mb4 FROM utf8mb4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "ENUM" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| utf8mb4
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `utf8mb4`.`utf8mb4` AS `utf8mb4` from `utf8mb4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `utf8mb4`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT utf8mb4.utf8mb4.utf8mb4 zzz FROM utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4.utf8mb4.utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 UNION DISTINCT SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT utf8mb4 UNION ALL SELECT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HIGH_PRIORITY utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HIGH_PRIORITY zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HISTOGRAM utf8mb4
< error
| 1054 42S22 Unknown column 'HISTOGRAM' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DISTINCT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DISTINCT zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_SMALL_RESULT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_SMALL_RESULT zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_THREAD utf8mb4
< error
| 1054 42S22 Unknown column 'SQL_THREAD' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BIG_RESULT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BIG_RESULT zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BUFFER_RESULT utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BUFFER_RESULT zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CALC_FOUND_ROWS utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CALC_FOUND_ROWS zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ALL utf8mb4
< error
| 1054 42S22 Unknown column 'utf8mb4' in 'field list'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE utf8mb4 (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ALL zzz utf8mb4
< error
| 1054 42S22 Unknown column 'zzz' in 'field list'
