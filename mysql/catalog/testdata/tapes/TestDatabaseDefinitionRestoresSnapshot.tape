-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDatabaseDefinitionRestoresSnapshot default#1
> exec
| CREATE DATABASE shared
< ok 1 0
> exec
| CREATE TABLE shared.parent (id int PRIMARY KEY)
< ok 0 0
> exec
| CREATE DATABASE src
< ok 1 0
> exec
| USE src
< ok 0 0
> exec
| CREATE TABLE `a``b` (id int NOT NULL, `c``d` varchar(20) COMMENT 'it''s a \\ path', parent_id int, PRIMARY KEY (id), KEY `k``1` (`c``d`), KEY parent_id (parent_id), CONSTRAINT `fk``x` FOREIGN KEY (parent_id) REFERENCES shared.parent (id), CONSTRAINT `ck``1` CHECK (id > 0)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8 COMMENT 'table''s \\ comment'
< ok 0 0
> exec
| CREATE TABLE u (id int PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE mem (a int NOT NULL, b int, PRIMARY KEY (a) USING BTREE, KEY k_hash (b), KEY k_btree (a, b) USING BTREE) ENGINE=MEMORY
< ok 0 0
> exec
| CREATE TRIGGER `tr``1` BEFORE INSERT ON `a``b` FOR EACH ROW SET NEW.`c``d` = 'x'
< ok 0 0
> exec
| CREATE FUNCTION `f``n`(x int) RETURNS int
|     DETERMINISTIC
| BEGIN
|   DECLARE y int;
|   SET y = x + 1;
|   RETURN y;
| END
< ok 0 0
> exec
| CREATE VIEW v_fn AS SELECT `f``n`(1) AS n
< ok 0 0
> exec
| CREATE VIEW v_alias AS SELECT src.id AS id FROM `a``b` AS src JOIN u ON src.id = u.id
< ok 0 0
> query
| SELECT QUOTE(COLUMN_COMMENT) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = 'a`b' AND COLUMN_NAME = 'c`d'
< result
< column "QUOTE(COLUMN_COMMENT)" "VARCHAR"
< row
: bytes
| 'it\'s a \\ path'
> query
| SELECT QUOTE(TABLE_COMMENT) FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = 'a`b'
< result
< column "QUOTE(TABLE_COMMENT)" "VARCHAR"
< row
: bytes
| 'table\'s \\ comment'
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_alias
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on((`src`.`id` = `src`.`u`.`id`)))
> query ^
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `f``n`(1) AS `n`
> exec
| CREATE DATABASE dump
< ok 1 0
> exec
| USE dump
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| CREATE VIEW `v_alias` AS SELECT
|   1 AS `id`
< ok 0 0
> exec
| CREATE VIEW `v_fn` AS SELECT
|   1 AS `n`
< ok 0 0
> exec
| CREATE TABLE `a``b` (
|   `id` int NOT NULL,
|   `c``d` varchar(20) DEFAULT NULL COMMENT 'it\'s a \\ path',
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k``1` (`c``d`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent` (`id`),
|   CONSTRAINT `ck``1` CHECK (`id` > 0)
| ) ENGINE=InnoDB row_format=COMPRESSED key_block_size=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='table\'s \\ comment'
< ok 0 0
> exec
| CREATE TABLE `u` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| CREATE TABLE `mem` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`) USING BTREE,
|   KEY `k_hash` (`b`),
|   KEY `k_btree` (`a`, `b`) USING BTREE
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE FUNCTION `f``n`(x int) RETURNS int
|     DETERMINISTIC
| BEGIN
|   DECLARE y int;
|   SET y = x + 1;
|   RETURN y;
| END
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_alias` AS select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on((`src`.`id` = `src`.`u`.`id`)))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_fn` AS select `f``n`(1) AS `n`
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| CREATE TRIGGER `tr``1`
| BEFORE INSERT ON `a``b` FOR EACH ROW
| SET NEW.`c``d` = 'x'
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> exec
| CREATE DATABASE `sdl`
< ok 1 0
> exec
| USE `sdl`
< ok 0 0
> exec
| CREATE TABLE `a``b` (
|   `id` int NOT NULL,
|   `c``d` varchar(20) DEFAULT NULL COMMENT 'it\'s a \\ path',
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k``1` (`c``d`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent` (`id`),
|   CONSTRAINT `ck``1` CHECK (`id` > 0)
| ) ENGINE=InnoDB row_format=COMPRESSED key_block_size=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='table\'s \\ comment'
< ok 0 0
> exec
|
|
| CREATE TABLE `mem` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`) USING BTREE,
|   KEY `k_btree` (`a`, `b`) USING BTREE,
|   KEY `k_hash` (`b`)
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
|
|
| CREATE TABLE `u` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> exec
|
|
| CREATE FUNCTION `f``n`(x int) RETURNS int
|     DETERMINISTIC
| BEGIN
|   DECLARE y int;
|   SET y = x + 1;
|   RETURN y;
| END
< ok 0 0
> exec
|
|
| CREATE OR REPLACE VIEW `v_alias` AS
| select
|   `src`.`id` AS `id`
| from (`a``b` `src`
|   join `u` on((`src`.`id` = `u`.`id`)))
< ok 0 0
> exec
|
|
| CREATE OR REPLACE VIEW `v_fn` AS
| select
|   `f``n`(1) AS `n`
< ok 0 0
> exec
|
|
| CREATE TRIGGER `tr``1` BEFORE INSERT ON `a``b` FOR EACH ROW
| SET NEW.`c``d` = 'x'
< ok 0 0
> query
| SHOW CREATE TABLE `src`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int NOT NULL,
|   `c``d` varchar(20) DEFAULT NULL COMMENT 'it''s a \\ path',
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k``1` (`c``d`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent` (`id`),
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8 COMMENT='table''s \\ comment'
> query
| SHOW CREATE TABLE `dump`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int NOT NULL,
|   `c``d` varchar(20) DEFAULT NULL COMMENT 'it''s a \\ path',
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k``1` (`c``d`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent` (`id`),
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8 COMMENT='table''s \\ comment'
> query
| SHOW CREATE TABLE `src`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `dump`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `src`.`mem`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| mem
: bytes
| CREATE TABLE `mem` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`) USING BTREE,
|   KEY `k_hash` (`b`),
|   KEY `k_btree` (`a`,`b`) USING BTREE
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `dump`.`mem`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| mem
: bytes
| CREATE TABLE `mem` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`) USING BTREE,
|   KEY `k_hash` (`b`),
|   KEY `k_btree` (`a`,`b`) USING BTREE
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| dump
: string
| v_alias
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on((`src`.`id` = `src`.`u`.`id`)))
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_alias
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on((`src`.`id` = `src`.`u`.`id`)))
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| dump
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `f``n`(1) AS `n`
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `f``n`(1) AS `n`
> query
| SELECT n FROM `dump`.v_fn
< result
< column "n" "INT"
< row
: int64 2
> query
| SHOW CREATE FUNCTION `src`.`f``n`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f``n`(x int) RETURNS int
|     DETERMINISTIC
| BEGIN
|   DECLARE y int;
|   SET y = x + 1;
|   RETURN y;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `dump`.`f``n`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f``n`(x int) RETURNS int
|     DETERMINISTIC
| BEGIN
|   DECLARE y int;
|   SET y = x + 1;
|   RETURN y;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT GROUP_CONCAT(CONCAT_WS(' ', TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT) ORDER BY TRIGGER_NAME) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ?
: string
| src
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT) ORDER BY TRIGGER_NAME)" "LONGTEXT"
< row
: bytes
| tr`1 BEFORE INSERT a`b SET NEW.`cc`dd` = 'x'
> query ^
: string
| dump
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT) ORDER BY TRIGGER_NAME)" "LONGTEXT"
< row
: bytes
| tr`1 BEFORE INSERT a`b SET NEW.`cc`dd` = 'x'
> query
| SHOW CREATE TABLE `src`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int NOT NULL,
|   `c``d` varchar(20) DEFAULT NULL COMMENT 'it''s a \\ path',
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k``1` (`c``d`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent` (`id`),
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8 COMMENT='table''s \\ comment'
> query
| SHOW CREATE TABLE `sdl`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int NOT NULL,
|   `c``d` varchar(20) DEFAULT NULL COMMENT 'it''s a \\ path',
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k``1` (`c``d`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent` (`id`),
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8 COMMENT='table''s \\ comment'
> query
| SHOW CREATE TABLE `src`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `sdl`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `src`.`mem`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| mem
: bytes
| CREATE TABLE `mem` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`) USING BTREE,
|   KEY `k_hash` (`b`),
|   KEY `k_btree` (`a`,`b`) USING BTREE
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE `sdl`.`mem`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| mem
: bytes
| CREATE TABLE `mem` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`) USING BTREE,
|   KEY `k_btree` (`a`,`b`) USING BTREE,
|   KEY `k_hash` (`b`)
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| sdl
: string
| v_alias
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`id` AS `id` from (`sdl`.`a``b` `src` join `sdl`.`u` on((`src`.`id` = `sdl`.`u`.`id`)))
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_alias
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on((`src`.`id` = `src`.`u`.`id`)))
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| sdl
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `f``n`(1) AS `n`
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `f``n`(1) AS `n`
> query
| SELECT n FROM `sdl`.v_fn
< result
< column "n" "INT"
< row
: int64 2
> query
| SHOW CREATE FUNCTION `src`.`f``n`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f``n`(x int) RETURNS int
|     DETERMINISTIC
| BEGIN
|   DECLARE y int;
|   SET y = x + 1;
|   RETURN y;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `sdl`.`f``n`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `f``n`(x int) RETURNS int
|     DETERMINISTIC
| BEGIN
|   DECLARE y int;
|   SET y = x + 1;
|   RETURN y;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SELECT GROUP_CONCAT(CONCAT_WS(' ', TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT) ORDER BY TRIGGER_NAME) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ?
: string
| src
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT) ORDER BY TRIGGER_NAME)" "LONGTEXT"
< row
: bytes
| tr`1 BEFORE INSERT a`b SET NEW.`cc`dd` = 'x'
> query ^
: string
| sdl
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT) ORDER BY TRIGGER_NAME)" "LONGTEXT"
< row
: bytes
| tr`1 BEFORE INSERT a`b SET NEW.`cc`dd` = 'x'
