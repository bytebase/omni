-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| CREATE TABLE th (a int NOT NULL, b int, PRIMARY KEY (a), KEY k_hash (b) USING HASH, KEY k_btree (a, b) USING BTREE)
< ok 0 0
> exec
| CREATE TRIGGER tr1 BEFORE INSERT ON u FOR EACH ROW SET NEW.id = NEW.id + 1
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
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on (`src`.`id` = `src`.`u`.`id`))
> query ^
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`f``n`(1) AS `n`
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| CREATE TABLE `th` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`),
|   KEY `k_hash` (`b`) USING HASH,
|   KEY `k_btree` (`a`, `b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER'
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
| CREATE OR REPLACE VIEW `v_alias` AS select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on (`src`.`id` = `src`.`u`.`id`))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v_fn` AS select `src`.`f``n`(1) AS `n`
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET character_set_results = utf8mb4
< ok 0 0
> exec
| SET collation_connection = utf8mb4_general_ci
< ok 0 0
> exec
| SET sql_mode = 'STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER'
< ok 0 0
> exec
| CREATE TRIGGER `tr1`
| BEFORE INSERT ON `u` FOR EACH ROW
| SET NEW.id = NEW.id + 1
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
| CREATE TABLE `th` (
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`a`),
|   KEY `k_btree` (`a`, `b`),
|   KEY `k_hash` (`b`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
< ok 0 0
> exec
|
|
| CREATE TABLE `u` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
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
|   join `u` on (`src`.`id` = `u`.`id`))
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
| CREATE TRIGGER `tr1` BEFORE INSERT ON `u` FOR EACH ROW
| SET NEW.id = NEW.id + 1
< ok 0 0
> query
| SHOW CREATE TABLE `src`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int(11) NOT NULL,
|   `c``d` varchar(20) COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'it\'s a \\ path',
|   `parent_id` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `k``1` (`c``d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `parent_id` (`parent_id`) BLOCK_SIZE 16384 LOCAL,
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'table\'s \\ comment'
> query
| SHOW CREATE TABLE `dump`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int(11) NOT NULL,
|   `c``d` varchar(20) COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'it\'s a \\ path',
|   `parent_id` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `k``1` (`c``d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `parent_id` (`parent_id`) BLOCK_SIZE 16384 LOCAL,
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'table\'s \\ comment'
> query
| SHOW CREATE TABLE `src`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `dump`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `src`.`th`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| th
: bytes
| CREATE TABLE `th` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`),
|   KEY `k_hash` (`b`) BLOCK_SIZE 16384 LOCAL,
|   KEY `k_btree` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `dump`.`th`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| th
: bytes
| CREATE TABLE `th` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`),
|   KEY `k_hash` (`b`) BLOCK_SIZE 16384 LOCAL,
|   KEY `k_btree` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on (`src`.`id` = `src`.`u`.`id`))
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_alias
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on (`src`.`id` = `src`.`u`.`id`))
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
| select `src`.`f``n`(1) AS `n`
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`f``n`(1) AS `n`
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
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`n`(x int) RETURNS int(11)
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
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION `dump`.`f``n`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`n`(x int) RETURNS int(11)
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
| utf8mb4_general_ci
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| src
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr1
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| u
: bytes
| TRIGGER tr1 BEFORE INSERT ON u FOR EACH ROW SET NEW.id = NEW.id + 1
> query ^
: string
| dump
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr1
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| u
: bytes
| TRIGGER `tr1`
| BEFORE INSERT ON `u` FOR EACH ROW
| SET NEW.id = NEW.id + 1
> query
| SELECT GROUP_CONCAT(CONCAT_WS(' ', TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX, INDEX_TYPE) ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ?
: string
| src
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX, INDEX_TYPE) ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a`b k`1 1 BTREE,a`b parent_id 1 BTREE,a`b PRIMARY 1 BTREE,th k_btree 1 BTREE,th k_btree 2 BTREE,th k_hash 1 HASH,th PRIMARY 1 BTREE,u PRIMARY 1 BTREE
> query ^
: string
| dump
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX, INDEX_TYPE) ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a`b k`1 1 BTREE,a`b parent_id 1 BTREE,a`b PRIMARY 1 BTREE,th k_btree 1 BTREE,th k_btree 2 BTREE,th k_hash 1 HASH,th PRIMARY 1 BTREE,u PRIMARY 1 BTREE
> query
| SHOW CREATE TABLE `src`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int(11) NOT NULL,
|   `c``d` varchar(20) COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'it\'s a \\ path',
|   `parent_id` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `k``1` (`c``d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `parent_id` (`parent_id`) BLOCK_SIZE 16384 LOCAL,
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'table\'s \\ comment'
> query
| SHOW CREATE TABLE `sdl`.`a``b`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| a`b
: bytes
| CREATE TABLE `a``b` (
|   `id` int(11) NOT NULL,
|   `c``d` varchar(20) COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'it\'s a \\ path',
|   `parent_id` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk``x` FOREIGN KEY (`parent_id`) REFERENCES `shared`.`parent`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `k``1` (`c``d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `parent_id` (`parent_id`) BLOCK_SIZE 16384 LOCAL,
|   CONSTRAINT `ck``1` CHECK ((`id` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = COMPRESSED COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0 COMMENT = 'table\'s \\ comment'
> query
| SHOW CREATE TABLE `src`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `sdl`.`u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `src`.`th`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| th
: bytes
| CREATE TABLE `th` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`),
|   KEY `k_hash` (`b`) BLOCK_SIZE 16384 LOCAL,
|   KEY `k_btree` (`a`, `b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE `sdl`.`th`
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| th
: bytes
| CREATE TABLE `th` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`),
|   KEY `k_btree` (`a`, `b`) BLOCK_SIZE 16384 LOCAL,
|   KEY `k_hash` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| select `src`.`id` AS `id` from (`sdl`.`a``b` `src` join `sdl`.`u` on (`src`.`id` = `sdl`.`u`.`id`))
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_alias
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`id` AS `id` from (`src`.`a``b` `src` join `src`.`u` on (`src`.`id` = `src`.`u`.`id`))
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
| select `sdl`.`f``n`(1) AS `n`
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?
: string
| v_fn
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `src`.`f``n`(1) AS `n`
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
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`n`(x int) RETURNS int(11)
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
| utf8mb4_general_ci
> query
| SHOW CREATE FUNCTION `sdl`.`f``n`
< result
< column "Function" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "Create Function" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< row
: bytes
| f`n
: bytes
| STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER
: bytes
| CREATE DEFINER = `root`@`%` FUNCTION `f`n`(x int) RETURNS int(11)
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
| utf8mb4_general_ci
> query
| SELECT TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME
: string
| src
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr1
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| u
: bytes
| TRIGGER tr1 BEFORE INSERT ON u FOR EACH ROW SET NEW.id = NEW.id + 1
> query ^
: string
| sdl
< result
< column "TRIGGER_NAME" "VARCHAR"
< column "ACTION_TIMING" "VARCHAR"
< column "EVENT_MANIPULATION" "VARCHAR"
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "ACTION_STATEMENT" "TEXT"
< row
: bytes
| tr1
: bytes
| BEFORE
: bytes
| INSERT
: bytes
| u
: bytes
| TRIGGER `tr1` BEFORE INSERT ON `u` FOR EACH ROW
| SET NEW.id = NEW.id + 1
> query
| SELECT GROUP_CONCAT(CONCAT_WS(' ', TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX, INDEX_TYPE) ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ?
: string
| src
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX, INDEX_TYPE) ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a`b k`1 1 BTREE,a`b parent_id 1 BTREE,a`b PRIMARY 1 BTREE,th k_btree 1 BTREE,th k_btree 2 BTREE,th k_hash 1 HASH,th PRIMARY 1 BTREE,u PRIMARY 1 BTREE
> query ^
: string
| sdl
< result
< column "GROUP_CONCAT(CONCAT_WS(' ', TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX, INDEX_TYPE) ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a`b k`1 1 BTREE,a`b parent_id 1 BTREE,a`b PRIMARY 1 BTREE,th k_btree 1 BTREE,th k_btree 2 BTREE,th k_hash 1 HASH,th PRIMARY 1 BTREE,u PRIMARY 1 BTREE
