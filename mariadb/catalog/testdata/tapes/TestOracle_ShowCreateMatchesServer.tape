-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOracle_ShowCreateMatchesServer default#1
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (d DATE DEFAULT curdate(), u VARCHAR(20) DEFAULT concat('a', 'b'), z DATE DEFAULT DATE '2020-01-02', x INT DEFAULT (1 + 1))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `d` date DEFAULT curdate(),
|   `u` varchar(20) DEFAULT concat('a','b'),
|   `z` date DEFAULT '2020-01-02',
|   `x` int(11) DEFAULT (1 + 1)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, g INT DEFAULT (a), k INT DEFAULT ((a)), b INT DEFAULT (a + 1), c INT DEFAULT (-1), d INT DEFAULT ((1 + 1)))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `a` int(11) DEFAULT NULL,
|   `g` int(11) DEFAULT `a`,
|   `k` int(11) DEFAULT `a`,
|   `b` int(11) DEFAULT (`a` + 1),
|   `c` int(11) DEFAULT -1,
|   `d` int(11) DEFAULT (1 + 1)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a BIGINT UNSIGNED DEFAULT 18446744073709551614, b BIGINT DEFAULT -9223372036854775808, c DECIMAL(30,0) DEFAULT 99999999999999999999, g DECIMAL(30,0) AS (c + 99999999999999999999) VIRTUAL, h BIGINT UNSIGNED AS (a + 9223372036854775808) VIRTUAL, CHECK (a < 18446744073709551615))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `a` bigint(20) unsigned DEFAULT 18446744073709551614,
|   `b` bigint(20) DEFAULT -9223372036854775808,
|   `c` decimal(30,0) DEFAULT 99999999999999999999,
|   `g` decimal(30,0) GENERATED ALWAYS AS (`c` + 99999999999999999999) VIRTUAL,
|   `h` bigint(20) unsigned GENERATED ALWAYS AS (`a` + 9223372036854775808) VIRTUAL,
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a` < 18446744073709551615)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT) ENGINE=myisam
< ok 0 0
> exec
| ALTER TABLE t ENGINE=innodb
< ok 0 0
> exec
| CREATE TABLE p (a INT) PARTITION BY HASH (a) (PARTITION p0 ENGINE=innodb, PARTITION p1 ENGINE=INNODB)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `p`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY HASH (`a`)
| (PARTITION `p0` ENGINE = InnoDB,
|  PARTITION `p1` ENGINE = InnoDB)
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a TIME DEFAULT '1:2', b TIME(1) DEFAULT '-1 2:3:4.56', c DATETIME(3) DEFAULT '2020-1-2 3:4:5.6789', d DATE DEFAULT '2020-1-2', e DATETIME DEFAULT '2020-01-02', f TIMESTAMP(2) NULL DEFAULT '2020-01-02 03:04:05', g TIME DEFAULT '10:00:00.9', h DATETIME(1) DEFAULT '2020-01-02T03:04:05.15')
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `a` time DEFAULT '01:02:00',
|   `b` time(1) DEFAULT '-26:03:04.5',
|   `c` datetime(3) DEFAULT '2020-01-02 03:04:05.678',
|   `d` date DEFAULT '2020-01-02',
|   `e` datetime DEFAULT '2020-01-02 00:00:00',
|   `f` timestamp(2) NULL DEFAULT '2020-01-02 03:04:05.00',
|   `g` time DEFAULT '10:00:00',
|   `h` datetime(1) DEFAULT '2020-01-02 03:04:05.1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT(5) ZEROFILL DEFAULT 3, b INT ZEROFILL DEFAULT 42, c DECIMAL(6,2) ZEROFILL DEFAULT 1.5, d FLOAT(7,3) ZEROFILL DEFAULT 1.5, g INT(3) ZEROFILL DEFAULT 12345, h DOUBLE ZEROFILL DEFAULT 2)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `a` int(5) unsigned zerofill DEFAULT 00003,
|   `b` int(10) unsigned zerofill DEFAULT 0000000042,
|   `c` decimal(6,2) unsigned zerofill DEFAULT 0001.50,
|   `d` float(7,3) unsigned zerofill DEFAULT 001.500,
|   `g` int(3) unsigned zerofill DEFAULT 12345,
|   `h` double unsigned zerofill DEFAULT 0000000000000000000002
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
|q "CREATE TABLE t (a INT COMMENT 'x\ny\rz\\\\w''q', b INT COMMENT 'emoji 📝 here', KEY k (a) COMMENT 'k🔥') COMMENT 'line1\nline2 🔥'"
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `a` int(11) DEFAULT NULL COMMENT 'x\ny\rz\\w''q',
|   `b` int(11) DEFAULT NULL COMMENT 'emoji ? here',
|   KEY `k` (`a`) COMMENT 'k?'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='line1\nline2 ?'
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) ENGINE=myisam COMMENT='c' ROW_FORMAT=FIXED CHECKSUM=1 DELAY_KEY_WRITE=1 PACK_KEYS=1 AVG_ROW_LENGTH=10 MAX_ROWS=100 MIN_ROWS=5 STATS_SAMPLE_PAGES=7 STATS_AUTO_RECALC=1 STATS_PERSISTENT=0 PAGE_CHECKSUM=1
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci MIN_ROWS=5 MAX_ROWS=100 AVG_ROW_LENGTH=10 PACK_KEYS=1 STATS_PERSISTENT=0 STATS_AUTO_RECALC=1 STATS_SAMPLE_PAGES=7 CHECKSUM=1 PAGE_CHECKSUM=1 DELAY_KEY_WRITE=1 ROW_FORMAT=FIXED COMMENT='c'
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, b INT, KEY b (b) USING HASH) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE INDEX c ON t (id) USING HASH
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `id` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `b` (`b`) USING HASH,
|   KEY `c` (`id`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a TEXT, b BLOB, c VARCHAR(10), UNIQUE KEY ua (a), UNIQUE KEY ub (b(10)), UNIQUE KEY uab (a, c))
< ok 0 0
> exec
| CREATE TABLE u (a TEXT UNIQUE)
< ok 0 0
> exec
| CREATE UNIQUE INDEX ux ON t (b)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| u
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
|   `id` int(11) NOT NULL,
|   `a` text DEFAULT NULL,
|   `b` blob DEFAULT NULL,
|   `c` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ub` (`b`(10)),
|   UNIQUE KEY `ua` (`a`) USING HASH,
|   UNIQUE KEY `uab` (`a`,`c`) USING HASH,
|   UNIQUE KEY `ux` (`b`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE `u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `a` text DEFAULT NULL,
|   UNIQUE KEY `a` (`a`) USING HASH
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, s VARCHAR(10), price DECIMAL(10,2), tax DECIMAL(5,2), g1 DECIMAL(10,2) AS (price * (1 + tax / 100)) STORED, g2 INT AS (a + b * 2 - (a - b)) VIRTUAL, g3 VARCHAR(30) AS (concat(s, '-', a)) VIRTUAL, g4 INT AS (a and b) VIRTUAL, g5 INT AS (not (a = 1)) VIRTUAL, g6 INT AS (if(a > b, a, b)) VIRTUAL, g7 INT AS (a between 1 and 10) VIRTUAL, g8 INT AS (a << 2 | b & 1) VIRTUAL)
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
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
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `s` varchar(10) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `tax` decimal(5,2) DEFAULT NULL,
|   `g1` decimal(10,2) GENERATED ALWAYS AS (`price` * (1 + `tax` / 100)) STORED,
|   `g2` int(11) GENERATED ALWAYS AS (`a` + `b` * 2 - (`a` - `b`)) VIRTUAL,
|   `g3` varchar(30) GENERATED ALWAYS AS (concat(`s`,'-',`a`)) VIRTUAL,
|   `g4` int(11) GENERATED ALWAYS AS (`a` <> 0 and `b` <> 0) VIRTUAL,
|   `g5` int(11) GENERATED ALWAYS AS (`a` <> 1) VIRTUAL,
|   `g6` int(11) GENERATED ALWAYS AS (if(`a` > `b`,`a`,`b`)) VIRTUAL,
|   `g7` int(11) GENERATED ALWAYS AS (`a` between 1 and 10) VIRTUAL,
|   `g8` int(11) GENERATED ALWAYS AS (`a` << 2 | `b` & 1) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE SEQUENCE s1
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT nextval(s1), b INT DEFAULT (NEXT VALUE FOR s1), c INT DEFAULT lastval(s1), d INT DEFAULT (PREVIOUS VALUE FOR s1), e INT DEFAULT setval(s1, 5))
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| s1
: bytes
| SEQUENCE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE SEQUENCE `s1`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s1
: bytes
| CREATE SEQUENCE `s1` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
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
|   `a` int(11) DEFAULT nextval(`omni_testoracle_showcreatematchesserver`.`s1`),
|   `b` int(11) DEFAULT nextval(`omni_testoracle_showcreatematchesserver`.`s1`),
|   `c` int(11) DEFAULT lastval(`omni_testoracle_showcreatematchesserver`.`s1`),
|   `d` int(11) DEFAULT lastval(`omni_testoracle_showcreatematchesserver`.`s1`),
|   `e` int(11) DEFAULT setval(`omni_testoracle_showcreatematchesserver`.`s1`,5,1,0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_showcreatematchesserver`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_showcreatematchesserver`
< ok 1 0
> exec
| USE `omni_testoracle_showcreatematchesserver`
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE FUNCTION f(x INT) RETURNS INT DETERMINISTIC RETURN x
< ok 0 0
> exec
| CREATE VIEW v AS SELECT `omni_testoracle_showcreatematchesserver`.t.a, `omni_testoracle_showcreatematchesserver`.f(1) AS one FROM `omni_testoracle_showcreatematchesserver`.t
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_testoracle_showcreatematchesserver" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
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
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE VIEW `v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`omni_testoracle_showcreatematchesserver`.`f`(1) AS `one` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
