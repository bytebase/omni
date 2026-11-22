-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestQuerySimulation default#1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` AS `a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO @a
< error
| 1064 42000 …near 'INTO @a' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(a AS TIME) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as time) AS `CAST(a AS TIME)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' CHARSET utf8mb4
< error
| 1064 42000 …near 'INTO OUTFILE 's' CHARSET utf8mb4' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a AS CHAR(1)) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,1,65) AS `WEIGHT_STRING(a AS CHAR(1))` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LEVEL 1 ASC) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,1) AS `WEIGHT_STRING(a LEVEL 1 ASC)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LEVEL 1 DESC) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,257) AS `WEIGHT_STRING(a LEVEL 1 DESC)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LEVEL 1 REVERSE) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,65537) AS `WEIGHT_STRING(a LEVEL 1 REVERSE)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LEVEL 1) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,1) AS `WEIGHT_STRING(a LEVEL 1)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LEVEL 1 ASC REVERSE) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,65537) AS `WEIGHT_STRING(a LEVEL 1 ASC REVERSE)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LEVEL 1, 1) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,1) AS `WEIGHT_STRING(a LEVEL 1, 1)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LEVEL 1 - 1) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,1) AS `WEIGHT_STRING(a LEVEL 1 - 1)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,0,0,1) AS `WEIGHT_STRING(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a LIKE 1) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a` like 1,0,0,1) AS `WEIGHT_STRING(a LIKE 1)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS WITH a AS(SELECT a) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH a AS(SELECT a) SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH a AS(SELECT a) SELECT a INTO @a
< error
| 1064 42000 …near 'INTO @a' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a UNION SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a UNION SELECT a zzz
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a UNION(SELECT a)
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT EXISTS(SELECT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select exists(select `a`.`a` limit 1) AS `EXISTS(SELECT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT EXISTS(SELECT a zzz) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select exists(select `a`.`a` AS `zzz` limit 1) AS `EXISTS(SELECT a zzz)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a FROM a zzz
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `zzz`.`a` AS `a` from `a` `zzz`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a FROM DUAL
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FROM DUMPFILE
< error
| 1146 42S02 Table 'omni_simulation.DUMPFILE' doesn't exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRING a
< error
| 1054 42S22 Unknown column 'STRING' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN STRAIGHT_JOIN a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT STRAIGHT_JOIN STRAIGHT_JOIN zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_NO_CACHE a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_NO_CACHE zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CACHE a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CACHE zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FOR UPDATE
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a LOCK IN SHARE MODE
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO @a FOR UPDATE
< error
| 1064 42000 …near 'INTO @a FOR UPDATE' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FOR UPDATE WAIT 1
< error
| 1350 HY000 View's SELECT contains a '[NO]WAIT' clause
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FOR UPDATE NOWAIT
< error
| 1350 HY000 View's SELECT contains a '[NO]WAIT' clause
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FOR UPDATE SKIP LOCKED
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a, a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a, a zzz
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT *
< error
| 1096 HY000 No tables used
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.*
< error
| 1051 42S02 Unknown table 'omni_simulation.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a AS a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a AS 's'
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a 's'
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME(1) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime(1) AS `CURRENT_TIME(1)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_USER FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_user() AS `CURRENT_USER` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_USER() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_user() AS `CURRENT_USER()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a WHERE a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a OR a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <> 0 or `a`.`a` <> 0 AS `a OR a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a XOR a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` xor `a`.`a` AS `a XOR a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a AND a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <> 0 and `a`.`a` <> 0 AS `a AND a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NOT a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` = 0 AS `NOT a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS TRUE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is true AS `a IS TRUE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS NOT TRUE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is not true AS `a IS NOT TRUE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS FALSE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is false AS `a IS FALSE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS NOT FALSE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is not false AS `a IS NOT FALSE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS UNKNOWN FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is null AS `a IS UNKNOWN` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS NOT UNKNOWN FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is not null AS `a IS NOT UNKNOWN` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS NULL FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is null AS `a IS NULL` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IS NOT NULL FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` is not null AS `a IS NOT NULL` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a <=> a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <=> `a`.`a` AS `a <=> a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a = a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` = `a`.`a` AS `a = a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a = ALL(SELECT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` = all (select `a`.`a`) AS `a = ALL(SELECT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a =(SELECT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` = (select `a`.`a`) AS `a =(SELECT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IN(SELECT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` in (select `a`.`a`) AS `a IN(SELECT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a NOT IN(SELECT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select !(`a`.`a` in (select `a`.`a`)) AS `a NOT IN(SELECT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IN(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` = `a`.`a` AS `a IN(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a IN(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` in (`a`.`a`,`a`.`a`) AS `a IN(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a NOT IN(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <> `a`.`a` AS `a NOT IN(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a NOT IN(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` not in (`a`.`a`,`a`.`a`) AS `a NOT IN(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a BETWEEN a AND a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` between `a`.`a` and `a`.`a` AS `a BETWEEN a AND a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a NOT BETWEEN a AND a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` not between `a`.`a` and `a`.`a` AS `a NOT BETWEEN a AND a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a SOUNDS LIKE a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select soundex(`a`.`a`) = soundex(`a`.`a`) AS `a SOUNDS LIKE a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SOUNDS LIKE a FROM a
< error
| 1054 42S22 Unknown column 'SOUNDS' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a LIKE a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` like `a`.`a` AS `a LIKE a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a LIKE a ESCAPE a FROM a
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a NOT LIKE a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a`  not like `a`.`a` AS `a NOT LIKE a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a NOT LIKE a ESCAPE a FROM a
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a REGEXP a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` regexp `a`.`a` AS `a REGEXP a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a NOT REGEXP a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select !(`a`.`a` regexp `a`.`a`) AS `a NOT REGEXP a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a | a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` | `a`.`a` AS `a | a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a & a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` & `a`.`a` AS `a & a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a << a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` << `a`.`a` AS `a << a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a >> a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` >> `a`.`a` AS `a >> a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + `a`.`a` AS `a + a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a - a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` - `a`.`a` AS `a - a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a DAY FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day AS `a + INTERVAL a DAY` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a - INTERVAL a DAY FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` - interval `a`.`a` day AS `a - INTERVAL a DAY` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL a DAY + a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day AS `INTERVAL a DAY + a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT + INTERVAL a DAY + a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day AS `+ INTERVAL a DAY + a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT - INTERVAL a DAY + a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` - interval `a`.`a` day AS `- INTERVAL a DAY + a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a * a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` * `a`.`a` AS `a * a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a / a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` / `a`.`a` AS `a / a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a % a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` MOD `a`.`a` AS `a % a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a DIV a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` DIV `a`.`a` AS `a DIV a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a MOD a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` MOD `a`.`a` AS `a MOD a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a ^ a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` ^ `a`.`a` AS `a ^ a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a || a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <> 0 or `a`.`a` <> 0 AS `a || a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a && a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <> 0 and `a`.`a` <> 0 AS `a && a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ! a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` = 0 AS `! a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a >= a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` >= `a`.`a` AS `a >= a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a > a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` > `a`.`a` AS `a > a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a <= a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <= `a`.`a` AS `a <= a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a < a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` < `a`.`a` AS `a < a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a <> a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` <> `a`.`a` AS `a <> a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a = ANY(SELECT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` in (select `a`.`a`) AS `a = ANY(SELECT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a`) AS `COLUMN_CREATE(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS INT) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS int) AS `COLUMN_CREATE(a, a AS INT)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS DATE) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS date) AS `COLUMN_CREATE(a, a AS DATE)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS CHAR) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS char charset utf8mb4 collate utf8mb4_uca1400_ai_ci ) AS `COLUMN_CREATE(a, a AS CHAR)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS UNSIGNED INT) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS unsigned int) AS `COLUMN_CREATE(a, a AS UNSIGNED INT)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS DOUBLE) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS double) AS `COLUMN_CREATE(a, a AS DOUBLE)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS REAL) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS double) AS `COLUMN_CREATE(a, a AS REAL)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS FLOAT) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS double) AS `COLUMN_CREATE(a, a AS FLOAT)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS DECIMAL) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS decimal) AS `COLUMN_CREATE(a, a AS DECIMAL)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS TIME) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS time) AS `COLUMN_CREATE(a, a AS TIME)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS DATETIME) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS datetime) AS `COLUMN_CREATE(a, a AS DATETIME)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a AS NCHAR) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a` AS char charset utf8mb3 collate utf8mb3_general_ci ) AS `COLUMN_CREATE(a, a AS NCHAR)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_CREATE(a, a, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_create(`a`.`a`,`a`.`a`,`a`.`a`,`a`.`a`) AS `COLUMN_CREATE(a, a, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(`a`.`a`) AS `TRIM(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(LEADING FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ltrim(`a`.`a`) AS `TRIM(LEADING FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(LEADING a FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(leading `a`.`a` from `a`.`a`) AS `TRIM(LEADING a FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(TRAILING a FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(trailing `a`.`a` from `a`.`a`) AS `TRIM(TRAILING a FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(BOTH a FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(both `a`.`a` from `a`.`a`) AS `TRIM(BOTH a FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(TRAILING FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select rtrim(`a`.`a`) AS `TRIM(TRAILING FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(TRANSACTION FROM a) FROM a
< error
| 1054 42S22 Unknown column 'TRANSACTION' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(BOTH FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(`a`.`a`) AS `TRIM(BOTH FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM(BTREE FROM a) FROM a
< error
| 1054 42S22 Unknown column 'BTREE' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRIM(a FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select trim(both `a`.`a` from `a`.`a`) AS `TRIM(a FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_DATE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curdate() AS `CURRENT_DATE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_POS FROM a
< error
| 1054 42S22 Unknown column 'CURRENT_POS' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a() FROM a
< error
| 1305 42000 FUNCTION omni_simulation.a does not exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ASCII(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ascii(`a`.`a`) AS `ASCII(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 's' FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 's' AS `s` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ? FROM a
< error
| 1064 42000 …near '? FROM a' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(`a`.`a`) AS `AVG(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT AVG_ROW_LENGTH(a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.AVG_ROW_LENGTH does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MEDIAN(a) OVER() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select percentile_cont(0.5) within group ( order by `a`.`a` desc) over () AS `MEDIAN(a) OVER()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW(a, a) FROM a
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROWCOUNT(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.ROWCOUNT does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT { a a } FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` AS `a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MATCH a AGAINST(a) FROM a
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CAST(a AS BINARY) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as char charset binary) AS `CAST(a AS BINARY)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CASE WHEN a THEN a END FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select case when `a`.`a` then `a`.`a` end AS `CASE WHEN a THEN a END` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CASE a WHEN a THEN a END FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select case `a`.`a` when `a`.`a` then `a`.`a` end AS `CASE a WHEN a THEN a END` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CONVERT(a, BINARY) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as char charset binary) AS `CONVERT(a, BINARY)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CONVERT(a USING utf8mb4) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select convert(`a`.`a` using utf8mb4) AS `CONVERT(a USING utf8mb4)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DEFAULT(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select default(`a`.`a`) AS `DEFAULT(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DEFINER(a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.DEFINER does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT VALUE(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select value(`a`.`a`) AS `VALUE(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NEXT VALUE FOR a FROM a
< error
| 4089 42S02 'omni_simulation.a' is not a SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NEXTVAL(a) FROM a
< error
| 4089 42S02 'omni_simulation.a' is not a SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PREVIOUS VALUE FOR a FROM a
< error
| 4089 42S02 'omni_simulation.a' is not a SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LASTVAL(a) FROM a
< error
| 4089 42S02 'omni_simulation.a' is not a SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SETVAL(a, 1) FROM a
< error
| 4089 42S02 'omni_simulation.a' is not a SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SHARE(a, 1) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.SHARE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SETVAL(a, 1, 1) FROM a
< error
| 4089 42S02 'omni_simulation.a' is not a SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SHARE(a, 1, 1) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.SHARE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SETVAL(a, 1, 1, 1) FROM a
< error
| 4089 42S02 'omni_simulation.a' is not a SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SHARE(a, 1, 1, 1) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.SHARE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a COLLATE utf8mb4_bin FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` collate utf8mb4_bin AS `a COLLATE utf8mb4_bin` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BINARY a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as char charset binary) AS `BINARY a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT + a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` AS `a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT - a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select -`a`.`a` AS `- a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ~ a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select ~`a`.`a` AS `~ a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMP(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as datetime) AS `TIMESTAMP(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMP(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestamp(`a`.`a`,`a`.`a`) AS `TIMESTAMP(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHAR(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select char(`a`.`a`) AS `CHAR(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHAR(a USING utf8mb4) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select char(`a`.`a` using utf8mb4) AS `CHAR(a USING utf8mb4)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_ROLE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_role() AS `CURRENT_ROLE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as date) AS `DATE(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DAY(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select dayofmonth(`a`.`a`) AS `DAY(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT HOUR(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select hour(`a`.`a`) AS `HOUR(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INSERT(a, a, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select insert(`a`.`a`,`a`.`a`,`a`.`a`,`a`.`a`) AS `INSERT(a, a, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INSERT_METHOD(a, a, a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.INSERT_METHOD does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT INTERVAL(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select interval(`a`.`a`,`a`.`a`) AS `INTERVAL(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT INTERVAL(a, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select interval(`a`.`a`,`a`.`a`,`a`.`a`) AS `INTERVAL(a, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT LEFT(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select left(`a`.`a`,`a`.`a`) AS `LEFT(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT LESS(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.LESS does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MINUTE(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select minute(`a`.`a`) AS `MINUTE(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MONTH(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select month(`a`.`a`) AS `MONTH(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MUTEX(a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.MUTEX does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT RIGHT(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select right(`a`.`a`,`a`.`a`) AS `RIGHT(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SECOND(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select second(`a`.`a`) AS `SECOND(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SESSION_USER() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select session_user() AS `SESSION_USER()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIME(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as time) AS `TIME(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT USER() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select user() AS `USER()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT USER_RESOURCES() FROM a
< error
| 1630 42000 FUNCTION omni_simulation.USER_RESOURCES does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT YEAR(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select year(`a`.`a`) AS `YEAR(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`a`.`a`,`a`.`a`) AS `SUBSTRING(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(a FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`a`.`a`,`a`.`a`) AS `SUBSTRING(a FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(a, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`a`.`a`,`a`.`a`,`a`.`a`) AS `SUBSTRING(a, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBSTRING(a FROM a FOR a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select substr(`a`.`a`,`a`.`a`,`a`.`a`) AS `SUBSTRING(a FROM a FOR a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ADDDATE(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day AS `ADDDATE(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ADMIN(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.ADMIN does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ADDDATE(a, INTERVAL a DAY) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day AS `ADDDATE(a, INTERVAL a DAY)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE_ADD(a, INTERVAL a DAY) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day AS `DATE_ADD(a, INTERVAL a DAY)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE_SUB(a, INTERVAL a DAY) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` - interval `a`.`a` day AS `DATE_SUB(a, INTERVAL a DAY)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT EXTRACT(DAY FROM a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select extract(day from `a`.`a`) AS `EXTRACT(DAY FROM a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(DATE, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(DATE, `a`.`a`) AS `GET_FORMAT(DATE, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT LOCALTIMESTAMP FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select localtimestamp() AS `LOCALTIMESTAMP` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIMESTAMP FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select current_timestamp() AS `CURRENT_TIMESTAMP` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT POSITION(a IN a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select locate(`a`.`a`,`a`.`a`) AS `POSITION(a IN a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROWNUM() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select rownum() AS `ROWNUM()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBDATE(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` - interval `a`.`a` day AS `SUBDATE(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUBJECT(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.SUBJECT does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUBDATE(a, INTERVAL a DAY) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` - interval `a`.`a` day AS `SUBDATE(a, INTERVAL a DAY)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SYSDATE() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sysdate() AS `SYSDATE()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SYSTEM() FROM a
< error
| 1630 42000 FUNCTION omni_simulation.SYSTEM does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SYSDATE(1) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sysdate(1) AS `SYSDATE(1)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SYSTEM(1) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.SYSTEM does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMPADD(DAY, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day AS `TIMESTAMPADD(DAY, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIMESTAMPDIFF(DAY, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select timestampdiff(DAY,`a`.`a`,`a`.`a`) AS `TIMESTAMPDIFF(DAY, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRIM_ORACLE(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select oracle_schema.trim(`a`.`a`) AS `TRIM_ORACLE(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_DATE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_date() AS `UTC_DATE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_TIME FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_time() AS `UTC_TIME` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT UTC_TIMESTAMP FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select utc_timestamp() AS `UTC_TIMESTAMP` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VALIDATION FROM a
< error
| 1054 42S22 Unknown column 'VALIDATION' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COLUMN_ADD(a, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_add(`a`.`a`,`a`.`a`,`a`.`a`) AS `COLUMN_ADD(a, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_DELETE(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_add(`a`.`a`,`a`.`a`,NULL AS int) AS `COLUMN_DELETE(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COLUMN_GET(a, a AS BINARY) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select column_get(`a`.`a`,`a`.`a` as char charset binary) AS `COLUMN_GET(a, a AS BINARY)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CHARSET(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select charset(`a`.`a`) AS `CHARSET(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT IF(a, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select if(`a`.`a`,`a`.`a`,`a`.`a`) AS `IF(a, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT LAST_VALUE(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select last_value(`a`.`a`) AS `LAST_VALUE(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT LAST_VALUE(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select last_value(`a`.`a`,`a`.`a`) AS `LAST_VALUE(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MOD(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` MOD `a`.`a` AS `MOD(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MODE(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.MODE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PASSWORD(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select password(`a`.`a`) AS `PASSWORD(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT PATH(a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.PATH does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT REPEAT(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select repeat(`a`.`a`,`a`.`a`) AS `REPEAT(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REPEATABLE(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.REPEATABLE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT REPLACE(a, a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select replace(`a`.`a`,`a`.`a`,`a`.`a`) AS `REPLACE(a, a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT REPLAY(a, a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.REPLAY does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRUNCATE(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select truncate(`a`.`a`,`a`.`a`) AS `TRUNCATE(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TYPE(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.TYPE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a AS BINARY(1)) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(cast(`a`.`a` as char(1) charset binary),0,1,65) AS `WEIGHT_STRING(a AS BINARY(1))` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT WEIGHT_STRING(a, 1, 1, 1) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select weight_string(`a`.`a`,1,1,1) AS `WEIGHT_STRING(a, 1, 1, 1)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CONTEXT() FROM a
< error
| 1630 42000 FUNCTION omni_simulation.CONTEXT does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT WITHOUT() FROM a
< error
| 1630 42000 FUNCTION omni_simulation.WITHOUT does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.a() FROM a
< error
| 1305 42000 FUNCTION a.a does not exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.a.a() FROM a
< error
| 1305 42000 FUNCTION a.a.a does not exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.REPLACE() FROM a
< error
| 1630 42000 FUNCTION a.REPLACE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.REPLAY() FROM a
< error
| 1630 42000 FUNCTION a.REPLAY does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.SUBSTRING() FROM a
< error
| 1630 42000 FUNCTION a.SUBSTRING does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.SUM() FROM a
< error
| 1630 42000 FUNCTION a.SUM does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.SUBSTRING(a FROM a) FROM a
< error
| 1128 HY000 Function 'a.SUBSTRING' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.SUBSTRING(a FROM a) zzz FROM a
< error
| 1128 HY000 Function 'a.SUBSTRING' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.TRIM() FROM a
< error
| 1630 42000 FUNCTION a.TRIM does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.TRIM_ORACLE() FROM a
< error
| 1630 42000 FUNCTION a.TRIM_ORACLE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.TRIM(LEADING FROM a) FROM a
< error
| 1128 HY000 Function 'a.TRIM' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.TRIM(LEADING FROM a) zzz FROM a
< error
| 1128 HY000 Function 'a.TRIM' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH a AGAINST(a IN BOOLEAN MODE) FROM a
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH a AGAINST(a IN NATURAL LANGUAGE MODE) FROM a
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH a AGAINST(a WITH QUERY EXPANSION) FROM a
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a(zzz) FROM a
< error
| 1305 42000 FUNCTION omni_simulation.a does not exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a(a) FROM a
< error
| 1305 42000 FUNCTION omni_simulation.a does not exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a(a, a) FROM a
< error
| 1305 42000 FUNCTION omni_simulation.a does not exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT AVG(DISTINCT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(distinct `a`.`a`) AS `AVG(DISTINCT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_AND(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_and(`a`.`a`) AS `BIT_AND(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_OR(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_or(`a`.`a`) AS `BIT_OR(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT_XOR(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select bit_xor(`a`.`a`) AS `BIT_XOR(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(*) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(0) AS `COUNT(*)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT COUNT(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(`a`.`a`) AS `COUNT(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CPU(a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.CPU does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT COUNT(DISTINCT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select count(distinct `a`.`a`) AS `COUNT(DISTINCT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MIN(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select min(`a`.`a`) AS `MIN(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MINUS(a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.MINUS does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MIN(DISTINCT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select min(`a`.`a`) AS `MIN(DISTINCT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MAX(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select max(`a`.`a`) AS `MAX(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT MAX(DISTINCT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select max(`a`.`a`) AS `MAX(DISTINCT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT STD(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select std(`a`.`a`) AS `STD(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VARIANCE(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select variance(`a`.`a`) AS `VARIANCE(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT STDDEV_SAMP(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select stddev_samp(`a`.`a`) AS `STDDEV_SAMP(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VAR_SAMP(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select var_samp(`a`.`a`) AS `VAR_SAMP(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUM(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sum(`a`.`a`) AS `SUM(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT SUPER(a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.SUPER does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SUM(DISTINCT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select sum(distinct `a`.`a`) AS `SUM(DISTINCT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`a`.`a` separator ',') AS `GROUP_CONCAT(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_ARRAYAGG(a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_arrayagg(`a`.`a`) AS `JSON_ARRAYAGG(a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_OBJECTAGG(a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select json_objectagg(`a`.`a`,`a`.`a`) AS `JSON_OBJECTAGG(a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT JSON_TABLE(a, a) FROM a
< error
| 1630 42000 FUNCTION omni_simulation.JSON_TABLE does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() over () AS `ROW_NUMBER() OVER()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT AVG(a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT RANK() OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DENSE_RANK() OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PERCENT_RANK() OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CUME_DIST() OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTILE(a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FIRST_VALUE(a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAST_VALUE(a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NTH_VALUE(a, a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LEAD(a, a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAG(a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT LAG(a, a) OVER a FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT PERCENTILE_CONT(a) WITHIN GROUP(ORDER BY a) OVER() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select percentile_cont(`a`.`a`) within group ( order by `a`.`a`) over () AS `PERCENTILE_CONT(a) WITHIN GROUP(ORDER BY a) OVER()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT PERCENTILE_DISC(a) WITHIN GROUP(ORDER BY a) OVER() FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select percentile_disc(`a`.`a`) within group ( order by `a`.`a`) over () AS `PERCENTILE_DISC(a) WITHIN GROUP(ORDER BY a) OVER()` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(DISTINCT a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(distinct `a`.`a` separator ',') AS `GROUP_CONCAT(DISTINCT a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(a SEPARATOR 's') FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`a`.`a` separator 's') AS `GROUP_CONCAT(a SEPARATOR 's')` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(a ORDER BY a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`a`.`a` order by `a`.`a` ASC separator ',') AS `GROUP_CONCAT(a ORDER BY a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GROUP_CONCAT(a ORDER BY a, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select group_concat(`a`.`a` order by `a`.`a` ASC,`a`.`a` ASC separator ',') AS `GROUP_CONCAT(a ORDER BY a, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS CHAR) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS VARCHAR(2147483648)) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as char(2147483648) charset utf8mb4) AS `CAST(a AS VARCHAR(2147483648))` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS NCHAR) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as char charset utf8mb3) AS `CAST(a AS NCHAR)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS INT) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as signed) AS `CAST(a AS INT)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS DATE) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as date) AS `CAST(a AS DATE)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS SIGNED) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as signed) AS `CAST(a AS SIGNED)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS SIGNED INT) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as signed) AS `CAST(a AS SIGNED INT)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS UNSIGNED) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as unsigned) AS `CAST(a AS UNSIGNED)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS UNSIGNED INT) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as unsigned) AS `CAST(a AS UNSIGNED INT)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS DECIMAL) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as decimal(10,0)) AS `CAST(a AS DECIMAL)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS FLOAT) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as float) AS `CAST(a AS FLOAT)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS DOUBLE) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as double) AS `CAST(a AS DOUBLE)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CAST(a AS DATETIME) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select cast(`a`.`a` as datetime) AS `CAST(a AS DATETIME)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a, a FROM a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH(a) AGAINST(a) FROM a
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT MATCH a, a AGAINST(a) FROM a
< error
| 1210 HY000 Incorrect arguments to AGAINST
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CASE WHEN a THEN a WHEN a THEN a END FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select case when `a`.`a` then `a`.`a` when `a`.`a` then `a`.`a` end AS `CASE WHEN a THEN a WHEN a THEN a END` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CASE WHEN a THEN a ELSE a END FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select case when `a`.`a` then `a`.`a` else `a`.`a` end AS `CASE WHEN a THEN a ELSE a END` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a DAY_HOUR FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day_hour AS `a + INTERVAL a DAY_HOUR` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a DAY_MICROSECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day_microsecond AS `a + INTERVAL a DAY_MICROSECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a DAY_MINUTE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day_minute AS `a + INTERVAL a DAY_MINUTE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a DAY_SECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` day_second AS `a + INTERVAL a DAY_SECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a HOUR_MICROSECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` hour_microsecond AS `a + INTERVAL a HOUR_MICROSECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a HOUR_MINUTE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` hour_minute AS `a + INTERVAL a HOUR_MINUTE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a HOUR_SECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` hour_second AS `a + INTERVAL a HOUR_SECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a MINUTE_MICROSECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` minute_microsecond AS `a + INTERVAL a MINUTE_MICROSECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a MINUTE_SECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` minute_second AS `a + INTERVAL a MINUTE_SECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a SECOND_MICROSECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` second_microsecond AS `a + INTERVAL a SECOND_MICROSECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a YEAR_MONTH FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` year_month AS `a + INTERVAL a YEAR_MONTH` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a WEEK FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` week AS `a + INTERVAL a WEEK` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a HOUR FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` hour AS `a + INTERVAL a HOUR` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a MINUTE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` minute AS `a + INTERVAL a MINUTE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a MONTH FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` month AS `a + INTERVAL a MONTH` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a QUARTER FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` quarter AS `a + INTERVAL a QUARTER` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a SECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` second AS `a + INTERVAL a SECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a MICROSECOND FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` microsecond AS `a + INTERVAL a MICROSECOND` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a + INTERVAL a YEAR FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` + interval `a`.`a` year AS `a + INTERVAL a YEAR` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(TIME, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(TIME, `a`.`a`) AS `GET_FORMAT(TIME, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(DATETIME, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(DATETIME, `a`.`a`) AS `GET_FORMAT(DATETIME, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT GET_FORMAT(TIMESTAMP, a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select get_format(DATETIME, `a`.`a`) AS `GET_FORMAT(TIMESTAMP, a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT AVG(ALL a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select avg(`a`.`a`) AS `AVG(ALL a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a HAVING a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a GROUP BY a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a GROUP BY a, a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a GROUP BY a WITH ROLLUP
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(zzz) FROM a
< error
| 4009 HY000 Window specification with name 'zzz' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(a) FROM a
< error
| 4009 HY000 Window specification with name 'a' is not defined
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(PARTITION BY a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() over ( partition by `a`.`a`) AS `ROW_NUMBER() OVER(PARTITION BY a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ORDER BY a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select row_number() over ( order by `a`.`a`) AS `ROW_NUMBER() OVER(ORDER BY a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(RANGE UNBOUNDED PRECEDING) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED PRECEDING) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS CURRENT ROW) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS 's' PRECEDING) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN UNBOUNDED FOLLOWING AND UNBOUNDED PRECEDING) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS BETWEEN 's' FOLLOWING AND UNBOUNDED PRECEDING) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE CURRENT ROW) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE GROUP) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE TIES) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ROW_NUMBER() OVER(ROWS UNBOUNDED PRECEDING EXCLUDE NO OTHERS) FROM a
< error
| 4017 HY000 Window frame is not allowed with 'row_number'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a ORDER BY a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a ORDER BY a, a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a GROUP BY a ASC
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a GROUP BY a DESC
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FETCH FIRST ROW ONLY
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FETCH NEXT ROW ONLY
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a FETCH FIRST ROW WITH TIES
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a LIMIT ?
< error
| 1064 42000 …near '?' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a LIMIT 9223372036854775808
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a LIMIT 2147483648
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a LIMIT 1
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a a FOR UPDATE
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a PROCEDURE a()
< error
| 1064 42000 …near 'PROCEDURE a()' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(0x41) FROM a
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIME(2147483648) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME(2147483648)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT CURRENT_TIME(9223372036854775808) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select curtime() AS `CURRENT_TIME(9223372036854775808)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a PROCEDURE a(zzz)
< error
| 1064 42000 …near 'PROCEDURE a(zzz)' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a PROCEDURE a(a)
< error
| 1064 42000 …near 'PROCEDURE a(a)' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a PROCEDURE a(a, a)
< error
| 1064 42000 …near 'PROCEDURE a(a, a)' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO @a, @a
< error
| 1064 42000 …near 'INTO @a, @a' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO@a
< error
| 1064 42000 …near 'INTO@a' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO DUMPFILE 's'
< error
| 1064 42000 …near 'INTO DUMPFILE 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' COLUMNS TERMINATED BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' COLUMNS TERMINATED BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' COLUMNS TERMINATED BY 's' TERMINATED BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' COLUMNS TERMINATED BY 's' TERMINATED BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' COLUMNS OPTIONALLY ENCLOSED BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' COLUMNS OPTIONALLY ENCLOSED BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' COLUMNS ENCLOSED BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' COLUMNS ENCLOSED BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' COLUMNS ESCAPED BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' COLUMNS ESCAPED BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' LINES TERMINATED BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' LINES TERMINATED BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' LINES TERMINATED BY 's' TERMINATED BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' LINES TERMINATED BY 's' TERMINATED BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTO OUTFILE 's' LINES STARTING BY 's'
< error
| 1064 42000 …near 'INTO OUTFILE 's' LINES STARTING BY 's'' at line 1
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT N's' FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 's' AS `s` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 's' FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4's' AS `s` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 's' 's' FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 'ss' AS `ss` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1 FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1 AS `1` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT DATE FROM a
< error
| 1054 42S22 Unknown column 'DATE' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT NULL FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select NULL AS `NULL` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT NUMBER FROM a
< error
| 1054 42S22 Unknown column 'NUMBER' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FALSE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 0 AS `FALSE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT FAST FROM a
< error
| 1054 42S22 Unknown column 'FAST' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TRUE FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1 AS `TRUE` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TRUNCATE FROM a
< error
| 1054 42S22 Unknown column 'TRUNCATE' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT 0x41 FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 0x41 AS `0x41` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT x'41' FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select X'41' AS `x'41'` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT b'1' FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 0x01 AS `b'1'` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT _utf8mb4 0x41 FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select _utf8mb4'A' AS `_utf8mb4 0x41` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 2147483648 FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 2147483648 AS `2147483648` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 9223372036854775808 FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 9223372036854775808 AS `9223372036854775808` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1.5 FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1.5 AS `1.5` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT 1e0 FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select 1e0 AS `1e0` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT TIME FROM a
< error
| 1054 42S22 Unknown column 'TIME' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT TIMESTAMP FROM a
< error
| 1054 42S22 Unknown column 'TIMESTAMP' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS WITH RECURSIVE a AS(SELECT a) SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.a.*
< error
| 1051 42S02 Unknown table 'a.a'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a.a FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` AS `a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT a.a.a FROM a
< error
| 1054 42S22 Unknown column 'a.a.a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT VALUE(a.a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select value(`a`.`a`) AS `VALUE(a.a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VALUE(.a.a) FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select value(`a`.`a`) AS `VALUE(.a.a)` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT VALUE(a.a.a) FROM a
< error
| 1054 42S22 Unknown column 'a.a.a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT `a` FROM a
< ok 0 0
> query
| SHOW FULL TABLES FROM omni_simulation
< result
< column "Tables_in_omni_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| VIEW
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW omni_simulation.`_v`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| _v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `_v` AS select `a`.`a` AS `a` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE omni_simulation.`a`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
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
| CREATE VIEW _v AS SELECT BIT FROM a
< error
| 1054 42S22 Unknown column 'BIT' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT BIT_AND FROM a
< error
| 1054 42S22 Unknown column 'BIT_AND' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT `a`() FROM a
< error
| 1305 42000 FUNCTION omni_simulation.a does not exist
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ACTION() FROM a
< error
| 1630 42000 FUNCTION omni_simulation.ACTION does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT FOUND() FROM a
< error
| 1630 42000 FUNCTION omni_simulation.FOUND does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a INTERSECT SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a EXCEPT SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a UNION DISTINCT SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT a UNION ALL SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HIGH_PRIORITY a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HIGH_PRIORITY zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT HISTORY a
< error
| 1054 42S22 Unknown column 'HISTORY' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DISTINCT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT DISTINCT zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT UNIQUE a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT UNIQUE zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT UNKNOWN a
< error
| 1054 42S22 Unknown column 'UNKNOWN' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_SMALL_RESULT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_SMALL_RESULT zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_THREAD a
< error
| 1054 42S22 Unknown column 'SQL_THREAD' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BIG_RESULT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BIG_RESULT zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BUFFER_RESULT a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_BUFFER_RESULT zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CALC_FOUND_ROWS a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT SQL_CALC_FOUND_ROWS zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ALL a
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| DROP DATABASE IF EXISTS omni_simulation; CREATE DATABASE omni_simulation; USE omni_simulation; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW _v AS SELECT ALL zzz a
< error
| 1054 42S22 Unknown column 'zzz' in 'SELECT'
