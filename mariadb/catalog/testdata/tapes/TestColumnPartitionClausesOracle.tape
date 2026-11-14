-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestColumnPartitionClausesOracle default#1
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT CONSTRAINT fk1 REFERENCES p (id), b INT CONSTRAINT REFERENCES p (id), d INT CHECK (d > 0) CONSTRAINT fk2 REFERENCES p (id) ON DELETE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL CHECK (`d` > 0),
|   KEY `fk1` (`a`),
|   KEY `b` (`b`),
|   KEY `fk2` (`d`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`b`) REFERENCES `p` (`id`),
|   CONSTRAINT `fk1` FOREIGN KEY (`a`) REFERENCES `p` (`id`),
|   CONSTRAINT `fk2` FOREIGN KEY (`d`) REFERENCES `p` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE n (a INT NOT NULL ENABLE DEFAULT 1, b INT NOT NULL ENABLE)
< ok 0 0
> query
| SHOW CREATE TABLE n
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| n
: bytes
| CREATE TABLE `n` (
|   `a` int(11) NOT NULL DEFAULT 1,
|   `b` int(11) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE i (a INT COMMENT 'c' INVISIBLE, b INT, c INT AS (b) VIRTUAL COMMENT 'g' INVISIBLE, d INT COMMENT 'd' WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING
< ok 0 0
> query
| SHOW CREATE TABLE i
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| i
: bytes
| CREATE TABLE `i` (
|   `a` int(11) DEFAULT NULL INVISIBLE COMMENT 'c',
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) GENERATED ALWAYS AS (`b`) VIRTUAL INVISIBLE COMMENT 'g',
|   `d` int(11) DEFAULT NULL WITHOUT SYSTEM VERSIONING COMMENT 'd'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE v (a INT, s TIMESTAMP(6) AS ROW START KEY, e TIMESTAMP(6) AS ROW END, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING
< ok 0 0
> query
| SHOW CREATE TABLE v
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `a` int(11) DEFAULT NULL,
|   `s` timestamp(6) GENERATED ALWAYS AS ROW START,
|   `e` timestamp(6) GENERATED ALWAYS AS ROW END,
|   PRIMARY KEY (`s`),
|   PERIOD FOR SYSTEM_TIME (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE v (a INT, s TIMESTAMP(6) AS ROW START NOT NULL ENABLE UNIQUE COMMENT 'r', e TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE, PERIOD FOR SYSTEM_TIME (s, e)) WITH SYSTEM VERSIONING
< ok 0 0
> query
| SHOW CREATE TABLE v
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `a` int(11) DEFAULT NULL,
|   `s` timestamp(6) GENERATED ALWAYS AS ROW START COMMENT 'r',
|   `e` timestamp(6) GENERATED ALWAYS AS ROW END INVISIBLE,
|   UNIQUE KEY `s` (`s`),
|   PERIOD FOR SYSTEM_TIME (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE v (a INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE v ADD a TIMESTAMP(6) AS ROW START KEY
< error
| 4124 HY000 Table `v` is not system-versioned
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE v (a INT, KEY k (a))
< ok 0 0
> exec
| ALTER TABLE v ADD s TIMESTAMP(6) AS ROW START, ADD e TIMESTAMP(6) AS ROW END, ADD PERIOD FOR SYSTEM_TIME (s, e), ADD SYSTEM VERSIONING
< ok 0 0
> query
| SHOW CREATE TABLE v
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE TABLE `v` (
|   `a` int(11) DEFAULT NULL,
|   `s` timestamp(6) GENERATED ALWAYS AS ROW START,
|   `e` timestamp(6) GENERATED ALWAYS AS ROW END,
|   KEY `k` (`a`),
|   PERIOD FOR SYSTEM_TIME (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE o (a INT x = 1)
< error
| 1911 HY000 Unknown option 'x'
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE o (a INT)
< ok 0 0
> exec
| ALTER TABLE o ADD b INT x = 's'
< error
| 1911 HY000 Unknown option 'x'
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'
< ok 0 0
> exec
| CREATE TABLE o (a INT NOT NULL DEFAULT 1 COMMENT 'c' x = 1 INVISIBLE CHECK (a > 0), b INT y = 2 z := 'it''s' Y = w, c INT q = 1 q = DEFAULT, d INT `select` = 0x41 READ_ONLY = 7)
< ok 0 0
> query
| SHOW CREATE TABLE o
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| o
: bytes
| CREATE TABLE `o` (
|   `a` int(11) NOT NULL DEFAULT 1 INVISIBLE `x`=1 COMMENT 'c' CHECK (`a` > 0),
|   `b` int(11) DEFAULT NULL `z`='it''s' `Y`=w,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL `select`=65 `READ_ONLY`=7
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'
< ok 0 0
> exec
| CREATE TABLE o (a INT x = 1, b INT y = 2)
< ok 0 0
> exec
| ALTER TABLE o MODIFY a INT, CHANGE b b INT k = 3
< ok 0 0
> query
| SHOW CREATE TABLE o
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| o
: bytes
| CREATE TABLE `o` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL `k`=3
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| SET sql_mode = 'IGNORE_BAD_TABLE_OPTIONS'
< ok 0 0
> exec
| CREATE TABLE o (a INT x = 1, b INT y = 2 z = 3)
< ok 0 0
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| ALTER TABLE o ADD c INT
< ok 0 0
> query
| SHOW CREATE TABLE o
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| o
: bytes
| CREATE TABLE `o` (
|   `a` int(11) DEFAULT NULL /* `x`=1 */,
|   `b` int(11) DEFAULT NULL /* `y`=2 `z`=3 */,
|   `c` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE q (a INT, b INT, c INT DEFAULT 1)
< ok 0 0
> exec
| ALTER TABLE q DROP q.b, ALTER q.c DROP DEFAULT, CHANGE .a z INT, ADD q.d INT
< ok 0 0
> query
| SHOW CREATE TABLE q
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| q
: bytes
| CREATE TABLE `q` (
|   `z` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE pr
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| pr
: bytes
| CREATE TABLE `pr` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`a`)
| (PARTITION `p0` VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION `p1` VALUES LESS THAN (20) ENGINE = InnoDB)
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE pr ADD PARTITION (p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> query
| SHOW CREATE TABLE pr
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| pr
: bytes
| CREATE TABLE `pr` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`a`)
| (PARTITION `p0` VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION `p1` VALUES LESS THAN (20) ENGINE = InnoDB,
|  PARTITION `p2` VALUES LESS THAN (30) ENGINE = InnoDB)
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| ALTER TABLE pr REORGANIZE PARTITION p1 INTO (p1a VALUES LESS THAN (15), p1b VALUES LESS THAN (20))
< ok 0 0
> query
| SHOW CREATE TABLE pr
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| pr
: bytes
| CREATE TABLE `pr` (
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`a`)
| (PARTITION `p0` VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION `p1a` VALUES LESS THAN (15) ENGINE = InnoDB,
|  PARTITION `p1b` VALUES LESS THAN (20) ENGINE = InnoDB)
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT PRIMARY KEY, b VARCHAR(10) COMMENT 'bb') PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN (30))
< ok 0 0
> exec
| ALTER TABLE pt CONVERT PARTITION p1 TO TABLE ct
< ok 0 0
> query
| SHOW CREATE TABLE pt
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| pt
: bytes
| CREATE TABLE `pt` (
|   `a` int(11) NOT NULL,
|   `b` varchar(10) DEFAULT NULL COMMENT 'bb',
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`a`)
| (PARTITION `p0` VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION `p2` VALUES LESS THAN (30) ENGINE = InnoDB)
> query
| SHOW CREATE TABLE ct
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ct
: bytes
| CREATE TABLE `ct` (
|   `a` int(11) NOT NULL,
|   `b` varchar(10) DEFAULT NULL COMMENT 'bb',
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT PRIMARY KEY, b VARCHAR(10) COMMENT 'bb') PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| CREATE TABLE ct (a INT PRIMARY KEY, b VARCHAR(10) COMMENT 'bb')
< ok 0 0
> exec
| ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20) WITHOUT VALIDATION
< ok 0 0
> query
| SHOW CREATE TABLE pt
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| pt
: bytes
| CREATE TABLE `pt` (
|   `a` int(11) NOT NULL,
|   `b` varchar(10) DEFAULT NULL COMMENT 'bb',
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (`a`)
| (PARTITION `p0` VALUES LESS THAN (10) ENGINE = InnoDB,
|  PARTITION `p1` VALUES LESS THAN (20) ENGINE = InnoDB)
> query
| SHOW CREATE TABLE ct
< error
| 1146 42S02 Table 'omni_testcolumnpartitionclausesoracle.ct' doesn't exist
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| CREATE TABLE ct (a INT)
< ok 0 0
> exec
| ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20)
< error
| 1736 HY000 Tables have different definitions
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT COMMENT '\n PARTITION BY x', b INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| CREATE TABLE ct (a INT COMMENT '\n PARTITION BY x', b VARCHAR(20))
< ok 0 0
> exec
| ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20)
< error
| 1736 HY000 Tables have different definitions
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE pt CONVERT TABLE ct TO PARTITION p1 VALUES LESS THAN (20)
< error
| 1146 42S02 Table 'omni_testcolumnpartitionclausesoracle.ct' doesn't exist
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| ALTER TABLE pt CONVERT PARTITION p9 TO TABLE ct
< error
| 1507 HY000 Wrong partition name or partition list
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT PRIMARY KEY) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| CREATE TABLE ct (a INT)
< ok 0 0
> exec
| ALTER TABLE pt CONVERT PARTITION p1 TO TABLE ct
< error
| 1050 42S01 Table 'ct' already exists
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT) PARTITION BY HASH (a) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE pt CONVERT PARTITION p0 TO TABLE ct
< error
| 1512 HY000 CONVERT PARTITION can only be used on RANGE/LIST partitions
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT)
< ok 0 0
> exec
| ALTER TABLE pt CONVERT PARTITION p0 TO TABLE ct
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
> exec
| CREATE TABLE pt (a INT)
< ok 0 0
> exec
| CREATE TABLE ct (a INT)
< ok 0 0
> exec
| ALTER TABLE pt CONVERT TABLE ct TO PARTITION p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET sql_mode = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS c, p, n, i, o, q, v, pr, pt, ct
< ok 0 0
